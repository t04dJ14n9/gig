// extern.go provides general-purpose wrapping of interpreter values for
// crossing into external Go code. When a synthesized struct has compiled
// methods (e.g., String()), the wrapper implements the corresponding Go
// interfaces (fmt.Stringer, fmt.Formatter) so that standard library and
// third-party code can discover them via type assertion / reflection.
//
// DESIGN: Only fmt.Stringer and fmt.Formatter need a wrapper, because
// reflect.StructOf types can't have methods. Encoding packages (json, etc.)
// work natively on the raw struct via struct tags and reflection —
// wrapping them would actually *break* native encoding by intercepting it.
package value

import (
	"fmt"
	"reflect"
	"strings"
)

// gigStructWrapper wraps an interpreter-synthesized struct value to implement
// Go interfaces (fmt.Stringer, fmt.Formatter) that the underlying anonymous
// struct type cannot satisfy because reflect.StructOf doesn't support methods.
//
// The wrapper is transparent: it delegates all fmt verbs to the underlying
// value, and only intercepts %T (for correct type name) and %v/%s (for
// String() dispatch).
type gigStructWrapper struct {
	iface     any           // the underlying struct value (clean, no phantom fields)
	typeName  string        // qualified type name from gig tag (e.g., "pkg.Type")
	stringer  func() string // nil if no String() method
	hasMethod bool          // true if String() method exists
}

// Ensure gigStructWrapper implements the relevant interfaces.
var (
	_ fmt.Stringer  = (*gigStructWrapper)(nil)
	_ fmt.Formatter = (*gigStructWrapper)(nil)
)

func (g *gigStructWrapper) String() string {
	if g.stringer != nil {
		return g.stringer()
	}
	return fmt.Sprint(g.iface)
}

func (g *gigStructWrapper) Format(f fmt.State, verb rune) {
	switch verb {
	case 'T':
		fmt.Fprint(f, g.typeName)
	case 'v', 's':
		if g.stringer != nil && (verb == 's' || (verb == 'v' && !f.Flag('#'))) {
			fmt.Fprint(f, g.stringer())
			return
		}
		// Delegate to the underlying value — it's clean (no phantom fields)
		if verb == 'v' && f.Flag('#') {
			fmt.Fprintf(f, "%s{", g.typeName)
			rv := reflect.ValueOf(g.iface)
			rt := rv.Type()
			for i := 0; i < rt.NumField(); i++ {
				if i > 0 {
					fmt.Fprint(f, " ")
				}
				fmt.Fprintf(f, "%s:%v", rt.Field(i).Name, rv.Field(i).Interface())
			}
			fmt.Fprint(f, "}")
		} else if verb == 'v' && f.Flag('+') {
			rv := reflect.ValueOf(g.iface)
			rt := rv.Type()
			fmt.Fprint(f, "{")
			for i := 0; i < rt.NumField(); i++ {
				if i > 0 {
					fmt.Fprint(f, " ")
				}
				fmt.Fprintf(f, "%s:%v", rt.Field(i).Name, rv.Field(i).Interface())
			}
			fmt.Fprint(f, "}")
		} else {
			fmt.Fprintf(f, "%v", g.iface)
		}
	default:
		fmt.Fprintf(f, "%"+string(verb), g.iface)
	}
}

// isGigStruct checks if a Go value is an interpreter-synthesized struct
// by looking for the "gig" struct tag on its first field.
// Returns the qualified type name (e.g., "pkg.TypeName") or "" if not a gig struct.
// Handles struct values, pointers to structs, and multiple levels of pointers (**T, ***T, etc.).
func isGigStruct(v any) string {
	if v == nil {
		return ""
	}
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	// Handle multiple levels of pointers: **T, ***T, etc.
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			elemType := rt.Elem()
			for elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			if elemType.Kind() != reflect.Struct || elemType.NumField() == 0 {
				return ""
			}
			gigTag := elemType.Field(0).Tag.Get("gig")
			if gigTag == "" {
				return ""
			}
			if strings.HasPrefix(gigTag, "#") {
				return gigTag[1:]
			}
			return gigTag
		}
		rv = rv.Elem()
		rt = rv.Type()
	}

	if rv.Kind() != reflect.Struct {
		return ""
	}
	rt = rv.Type()
	if rt.NumField() == 0 {
		return ""
	}
	gigTag := rt.Field(0).Tag.Get("gig")
	if gigTag == "" {
		return ""
	}
	if strings.HasPrefix(gigTag, "#") {
		return gigTag[1:]
	}
	return gigTag
}

// ExternWrap is a pass-through that returns the raw interface{} value.
// For most external Go code (json, container, sync, etc.),
// the raw struct is what they need — struct tags and reflection work natively.
//
// Use FmtWrap instead when passing values to fmt.* functions that check
// for fmt.Stringer/fmt.Formatter interfaces.
func ExternWrap(v Value) any {
	return v.Interface()
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
		return iface
	}

	typeName := isGigStruct(iface)
	if typeName == "" {
		return iface
	}

	// Check if the interpreted type has a String() method via the global resolver registry
	stringerFunc, hasStringer := resolveStringer(v)

	// Always return the wrapper for gig structs - it handles all fmt verbs correctly
	return &gigStructWrapper{
		iface:     iface,
		typeName:  typeName,
		stringer:  stringerFunc,
		hasMethod: hasStringer,
	}
}

// resolveStringer attempts to resolve the String() method for a value.
// Returns a function that can be called later, and a boolean indicating if found.
func resolveStringer(v Value) (func() string, bool) {
	// Try to call String() method via the global resolver registry
	result, found := CallMethod(nil, "String", v)
	if !found {
		// If not found, try with pointer to the value (for pointer receiver methods)
		if rv, ok := v.ReflectValue(); ok && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = CallMethod(nil, "String", ptrValue)
		}
	}

	if !found {
		return nil, false
	}
	str := result.String()
	return func() string { return str }, true
}

// SprintfExtern is a general-purpose fmt.Sprintf replacement that correctly
// handles %T for gigStructWrapper values. Go's fmt.Sprintf("%T") bypasses
// fmt.Formatter entirely and uses reflect.TypeOf().String(), so we must
// intercept %T ourselves.
func SprintfExtern(format string, args ...any) string {
	// Fast path: no %T in format string — use standard fmt.Sprintf
	if !strings.Contains(format, "%T") {
		return fmt.Sprintf(format, args...)
	}
	// Slow path: replace %T for gigStructWrapper args with their type name
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
					if w, ok := args[argIdx].(*gigStructWrapper); ok {
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
