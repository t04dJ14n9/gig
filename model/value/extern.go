// extern.go provides general-purpose wrapping of interpreter values for
// crossing into external Go code. When a synthesized struct has compiled
// methods (e.g., String()), the wrapper implements the corresponding Go
// interfaces (fmt.Stringer, fmt.Formatter, error) so that standard library and
// third-party code can discover them via type assertion / reflection.
//
// DESIGN: Only method-shaped host boundaries need a wrapper, because
// reflect.StructOf types can't have methods. Encoding packages (json, etc.)
// work natively on the raw struct via struct tags and reflection —
// wrapping them would actually *break* native encoding by intercepting it.
package value

import (
	"fmt"
	"reflect"
	"strings"
)

// gigStructError wraps an interpreter-synthesized struct value to implement
// Go interfaces (fmt.Stringer, fmt.Formatter, error, fmt.GoStringer) that the underlying anonymous
// struct type cannot satisfy because reflect.StructOf doesn't support methods.
//
// The wrapper is transparent: it delegates all fmt verbs to the underlying
// value, and only intercepts %T (for correct type name) and %v/%s (for
// String() dispatch).
type gigStructError struct {
	iface      any
	typeName   string
	stringer   func() string
	errorer    func() string
	gostringer func() string
}

// Ensure gigStructError implements the relevant interfaces.
var (
	_ fmt.Stringer   = (*gigStructError)(nil)
	_ fmt.Formatter  = (*gigStructError)(nil)
	_ fmt.GoStringer = (*gigStructError)(nil)
	_ error          = (*gigStructError)(nil)
)

func (g *gigStructError) String() string {
	if g.stringer != nil {
		return g.stringer()
	}
	if g.errorer != nil {
		return g.errorer()
	}
	return fmt.Sprint(g.iface)
}

func (g *gigStructError) Error() string {
	if g.errorer != nil {
		return g.errorer()
	}
	if g.stringer != nil {
		return g.stringer()
	}
	return fmt.Sprint(g.iface)
}

func (g *gigStructError) GoString() string {
	if g.gostringer != nil {
		return g.gostringer()
	}
	return g.defaultGoString()
}

func (g *gigStructError) Format(f fmt.State, verb rune) {
	switch verb {
	case 'T':
		_, _ = fmt.Fprint(f, g.typeName)
	case 'v', 's':
		if g.stringer != nil && (verb == 's' || (verb == 'v' && !f.Flag('#'))) {
			_, _ = fmt.Fprint(f, g.stringer())
			return
		}
		// Delegate to the underlying value — it's clean (no phantom fields)
		if verb == 'v' && f.Flag('#') {
			_, _ = fmt.Fprintf(f, "%s{", g.typeName)
			rv := reflect.ValueOf(g.iface)
			rt := rv.Type()
			for i := 0; i < rt.NumField(); i++ {
				if i > 0 {
					_, _ = fmt.Fprint(f, " ")
				}
				_, _ = fmt.Fprintf(f, "%s:%v", rt.Field(i).Name, rv.Field(i).Interface())
			}
			_, _ = fmt.Fprint(f, "}")
		} else if verb == 'v' && f.Flag('+') {
			rv := reflect.ValueOf(g.iface)
			rt := rv.Type()
			_, _ = fmt.Fprint(f, "{")
			for i := 0; i < rt.NumField(); i++ {
				if i > 0 {
					_, _ = fmt.Fprint(f, " ")
				}
				_, _ = fmt.Fprintf(f, "%s:%v", rt.Field(i).Name, rv.Field(i).Interface())
			}
			_, _ = fmt.Fprint(f, "}")
		} else {
			_, _ = fmt.Fprintf(f, "%v", g.iface)
		}
	default:
		_, _ = fmt.Fprintf(f, "%"+string(verb), g.iface)
	}
}

func asGigStructError(err error) (*gigStructError, bool) {
	wrapper, ok := err.(*gigStructError) //nolint:errorlint // Internal exact wrapper check; native wrapping is handled by caller traversal.
	return wrapper, ok
}

// isGigStruct checks if a Go value is an interpreter-synthesized struct
// by looking for the "gig" struct tag on its first field.
// Returns the qualified type name (e.g., "pkg.TypeName") or "" if not a gig struct.
// Handles struct values, pointers to structs, and multiple levels of pointers (**T, ***T, etc.).
func isGigStruct(v any) string {
	if v == nil {
		return ""
	}
	return gigStructNameFromType(baseReflectType(reflect.TypeOf(v)))
}

func baseReflectType(rt reflect.Type) reflect.Type {
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt
}

func gigStructNameFromType(rt reflect.Type) string {
	if rt.Kind() != reflect.Struct {
		return ""
	}
	if rt.NumField() == 0 {
		return ""
	}
	gigTag := rt.Field(0).Tag.Get("gig")
	if gigTag != "" {
		return normalizeGigTag(gigTag)
	}
	return gigStructNameFromPkgPath(rt)
}

func normalizeGigTag(gigTag string) string {
	if strings.HasPrefix(gigTag, "#") {
		return gigTag[1:]
	}
	return gigTag
}

func gigStructNameFromPkgPath(rt reflect.Type) string {
	for i := 0; i < rt.NumField(); i++ {
		pkgPath := rt.Field(i).PkgPath
		if idx := strings.LastIndex(pkgPath, "#"); idx >= 0 {
			return extractBareTypeName(pkgPath[idx+1:])
		}
	}
	return ""
}

// FmtWrap prepares a value.Value for passing to fmt.* functions.
// If the value is an interpreter-synthesized struct with compiled methods
// (e.g., String()), returns a wrapper that implements fmt.Stringer and
// fmt.Formatter. Otherwise returns the raw interface{} value.
//
// This is the boundary function for fmt.Print/Sprint/Fprintf/etc. — use it
// whenever passing interpreter values to ...interface{} variadic args in
// fmt-family functions.
func FmtWrap(v Value) any {
	iface := v.Interface()
	if iface == nil {
		return nil
	}

	typeName := isGigStruct(iface)
	if typeName == "" {
		return iface
	}

	// Check if the interpreted type has a String() method via the global resolver registry
	stringerFunc := resolveStringer(v)
	errorerFunc, _ := resolveErrorer(v)
	gostringerFunc := resolveGoStringer(v)

	// Always return the wrapper for gig structs - it handles all fmt verbs correctly
	return &gigStructError{
		iface:      iface,
		typeName:   typeName,
		stringer:   stringerFunc,
		errorer:    errorerFunc,
		gostringer: gostringerFunc,
	}
}

// resolveStringer attempts to resolve the String() method for a value.
// It returns nil when the interpreted type does not define String().
func resolveStringer(v Value) func() string {
	defer func() {
		_ = recover()
	}()
	// Try to call String() method via the global resolver registry
	result, found := callMethod("String", v)
	if !found {
		// If not found, try with pointer to the value (for pointer receiver methods)
		if rv, ok := v.ReflectValue(); ok && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = callMethod("String", ptrValue)
		}
	}

	if !found {
		return nil
	}
	str := result.String()
	return func() string { return str }
}

func resolveErrorer(v Value) (func() string, bool) {
	defer func() {
		_ = recover()
	}()
	result, found := callMethod("Error", v)
	if !found {
		if rv, ok := v.ReflectValue(); ok && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = callMethod("Error", ptrValue)
		}
	}
	if !found {
		return nil, false
	}
	str := result.String()
	return func() string { return str }, true
}

func resolveGoStringer(v Value) func() string {
	defer func() {
		_ = recover()
	}()
	result, found := callMethod("GoString", v)
	if !found {
		if rv, ok := v.ReflectValue(); ok && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = callMethod("GoString", ptrValue)
		}
	}
	if !found {
		return nil
	}
	str := result.String()
	return func() string { return str }
}

func extractBareTypeName(qualName string) string {
	if idx := strings.LastIndex(qualName, "."); idx >= 0 {
		return qualName[idx+1:]
	}
	return qualName
}

func extractGigTagFromType(rt reflect.Type) string {
	if rt.Kind() != reflect.Struct || rt.NumField() == 0 {
		return ""
	}
	return normalizeGigTag(rt.Field(0).Tag.Get("gig"))
}

// SprintfExtern is a general-purpose fmt.Sprintf replacement that correctly
// handles %T for gigStructError values. Go's fmt.Sprintf("%T") bypasses
// fmt.Formatter entirely and uses reflect.TypeOf().String(), so we must
// intercept %T ourselves.
func SprintfExtern(format string, args ...any) string {
	// Fast path: no %T in format string — use standard fmt.Sprintf
	if !strings.Contains(format, "%T") {
		return fmt.Sprintf(format, args...)
	}
	// Slow path: replace %T for gigStructError args with their type name
	var result strings.Builder
	argIdx := 0
	i := 0
	for i < len(format) {
		if format[i] == '%' {
			if i+1 < len(format) && format[i+1] == '%' {
				result.WriteString("%%")
				i += 2
				continue
			}
			j := i + 1
			// Skip flags
			for j < len(format) && (format[j] == '-' || format[j] == '+' || format[j] == '#' || format[j] == ' ' || format[j] == '0') {
				j++
			}
			// Skip width
			for j < len(format) && format[j] >= '0' && format[j] <= '9' {
				j++
			}
			// Skip precision
			if j < len(format) && format[j] == '.' {
				j++
				for j < len(format) && format[j] >= '0' && format[j] <= '9' {
					j++
				}
			}
			if j < len(format) {
				verb := format[j]
				if verb == 'T' && argIdx < len(args) {
					if w, ok := args[argIdx].(*gigStructError); ok {
						result.WriteString(w.typeName)
						argIdx++
						i = j + 1
						continue
					}
				}
				if argIdx < len(args) {
					result.WriteString(fmt.Sprintf(format[i:j+1], args[argIdx]))
					argIdx++
				} else {
					result.WriteString(format[i : j+1])
				}
				i = j + 1
				continue
			}
		}
		result.WriteByte(format[i])
		i++
	}
	return result.String()
}
