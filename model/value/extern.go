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
// Go interfaces (fmt.Stringer, fmt.Formatter, error, fmt.GoStringer) that the
// underlying anonymous struct type cannot satisfy because reflect.StructOf
// doesn't support methods.
//
// The wrapper is transparent: it delegates all fmt verbs to the underlying
// value, and only intercepts %T (for correct type name), %v/%s (for
// String() dispatch), and %#v (for GoString() dispatch). When the interpreted
// type has an Error() method, the wrapper also implements the error interface
// so that errors.As and type assertions work correctly.
//
// Method resolution is lazy — each lazy* function resolves its corresponding
// method at most once, only when the wrapper actually needs to dispatch it.
// This avoids eagerly invoking the interpreter for methods that may never be
// needed (and whose invocation in a side-channel VM may fail).
type gigStructWrapper struct {
	iface          any                     // the underlying struct value (clean, no phantom fields)
	typeName       string                  // qualified type name from gig tag (e.g., "pkg.Type")
	lazyStringer   func() (func() string, bool)
	lazyErrorer    func() (func() string, bool)
	lazyGoStringer func() (func() string, bool)
}

// Ensure gigStructWrapper implements the relevant interfaces.
var (
	_ fmt.Stringer   = (*gigStructWrapper)(nil)
	_ fmt.Formatter  = (*gigStructWrapper)(nil)
	_ fmt.GoStringer = (*gigStructWrapper)(nil)
	_ error          = (*gigStructWrapper)(nil)
)

// tryStringer / tryErrorer / tryGoStringer return (fn, ok) where fn is the
// resolved method callable and ok indicates whether the interpreted type
// actually defined the method. They are safe to call multiple times.
func (g *gigStructWrapper) tryStringer() (func() string, bool) {
	if g.lazyStringer == nil {
		return nil, false
	}
	return g.lazyStringer()
}
func (g *gigStructWrapper) tryErrorer() (func() string, bool) {
	if g.lazyErrorer == nil {
		return nil, false
	}
	return g.lazyErrorer()
}
func (g *gigStructWrapper) tryGoStringer() (func() string, bool) {
	if g.lazyGoStringer == nil {
		return nil, false
	}
	return g.lazyGoStringer()
}

// String implements fmt.Stringer. It prefers the interpreted type's String()
// method, but falls back to Error() so that types implementing only error
// (not Stringer) still produce meaningful output when fmt calls String().
// Note: fmt.handleMethods checks error before Stringer, so fmt.Sprint(wrapper)
// calls Error() first — but String() must still work correctly when called
// directly, e.g. by code that explicitly calls .String().
func (g *gigStructWrapper) String() string {
	if f, ok := g.tryStringer(); ok {
		return f()
	}
	if f, ok := g.tryErrorer(); ok {
		return f()
	}
	return fmt.Sprint(g.iface)
}

// Error implements error. It prefers the interpreted type's Error() method,
// but falls back to String() so that types implementing only Stringer
// (not error) still produce meaningful output.
// fmt.handleMethods checks error before Stringer, so fmt.Sprint(wrapper)
// dispatches here first — ensuring error messages are never accidentally
// replaced by a decorative String() representation.
func (g *gigStructWrapper) Error() string {
	if f, ok := g.tryErrorer(); ok {
		return f()
	}
	if f, ok := g.tryStringer(); ok {
		return f()
	}
	return fmt.Sprint(g.iface)
}

// GoString implements fmt.GoStringer. Dispatches to the interpreted
// GoString() method if present; otherwise falls back to the default
// Go-syntax representation produced by Format with the '#' flag.
func (g *gigStructWrapper) GoString() string {
	if f, ok := g.tryGoStringer(); ok {
		return f()
	}
	return g.defaultGoString()
}

// defaultGoString renders the wrapped value in Go-syntax form, matching
// native fmt's %#v output: "pkg.Type{Field: value, Field: value}" for
// struct values, "&pkg.Type{...}" for struct pointers, "(*pkg.Type)(nil)"
// for nil struct pointers.
func (g *gigStructWrapper) defaultGoString() string {
	rv := reflect.ValueOf(g.iface)
	// Dereference pointers for inspection; render "&T{...}" or "(*T)(nil)".
	prefix := ""
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Sprintf("(*%s)(nil)", g.typeName)
		}
		rv = rv.Elem()
		prefix = "&"
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Sprintf("%s%#v", prefix, g.iface)
	}
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString(g.typeName)
	sb.WriteByte('{')
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(rt.Field(i).Name)
		sb.WriteByte(':')
		fmt.Fprintf(&sb, "%#v", rv.Field(i).Interface())
	}
	sb.WriteByte('}')
	return sb.String()
}

func (g *gigStructWrapper) Format(f fmt.State, verb rune) {
	switch verb {
	case 'T':
		_, _ = fmt.Fprint(f, g.typeName)
	case 'v', 's':
		if verb == 's' || (verb == 'v' && !f.Flag('#') && !f.Flag('+')) {
			if fn, ok := g.tryStringer(); ok {
				_, _ = fmt.Fprint(f, fn())
				return
			}
			if verb == 'v' {
				// Plain %v: default struct rendering without type name.
				rv := reflect.ValueOf(g.iface)
				if rv.Kind() == reflect.Ptr {
					if rv.IsNil() {
						_, _ = fmt.Fprint(f, "<nil>")
						return
					}
					rv = rv.Elem()
					_, _ = fmt.Fprint(f, "&")
				}
				if rv.Kind() == reflect.Struct {
					rt := rv.Type()
					_, _ = fmt.Fprint(f, "{")
					for i := 0; i < rt.NumField(); i++ {
						if i > 0 {
							_, _ = fmt.Fprint(f, " ")
						}
						_, _ = fmt.Fprintf(f, "%v", rv.Field(i).Interface())
					}
					_, _ = fmt.Fprint(f, "}")
					return
				}
			}
			_, _ = fmt.Fprintf(f, "%v", g.iface)
			return
		}
		// %#v — Go-syntax representation. Prefer GoString() if defined;
		// otherwise fall back to defaultGoString which renders the struct
		// with commas between fields (matching native fmt).
		if verb == 'v' && f.Flag('#') {
			_, _ = fmt.Fprint(f, g.GoString())
			return
		}
		// %+v — fields shown with names.
		if verb == 'v' && f.Flag('+') {
			rv := reflect.ValueOf(g.iface)
			if rv.Kind() == reflect.Ptr {
				if rv.IsNil() {
					_, _ = fmt.Fprintf(f, "<nil>")
					return
				}
				rv = rv.Elem()
				_, _ = fmt.Fprint(f, "&")
			}
			if rv.Kind() != reflect.Struct {
				_, _ = fmt.Fprintf(f, "%+v", g.iface)
				return
			}
			rt := rv.Type()
			_, _ = fmt.Fprint(f, "{")
			for i := 0; i < rt.NumField(); i++ {
				if i > 0 {
					_, _ = fmt.Fprint(f, " ")
				}
				_, _ = fmt.Fprintf(f, "%s:%v", rt.Field(i).Name, rv.Field(i).Interface())
			}
			_, _ = fmt.Fprint(f, "}")
			return
		}
		_, _ = fmt.Fprintf(f, "%v", g.iface)
	default:
		_, _ = fmt.Fprintf(f, "%"+string(verb), g.iface)
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

// FmtWrap prepares a value.Value for passing to fmt.* functions.
// If the value is an interpreter-synthesized struct with compiled methods
// (e.g., String()), returns a wrapper that implements fmt.Stringer and
// fmt.Formatter. Otherwise returns the raw interface{} value.
//
// This is the boundary function for fmt.Print/Sprint/Fprintf/etc. — use it
// whenever passing interpreter values to ...interface{} variadic args in
// fmt-family functions.
//
// Method resolution (String/Error/GoString) is deferred: we capture the
// underlying Value and resolve lazily when the wrapper's corresponding method
// is actually invoked by fmt. Eager resolution can fail when the method body
// depends on program state not fully reachable from a side-channel VM.
func FmtWrap(v Value) any {
	iface := v.Interface()
	if iface == nil {
		return nil
	}

	typeName := isGigStruct(iface)
	if typeName == "" {
		return iface
	}

	captured := v

	// Lazy resolvers — each is called at most once by the wrapper, and only
	// when the corresponding fmt verb triggers its interface.
	var (
		stringerFunc   func() string
		stringerResolved bool
		errorerFunc    func() string
		errorerResolved bool
		gostringerFunc func() string
		gostringerResolved bool
	)
	lazyStringer := func() (func() string, bool) {
		if !stringerResolved {
			stringerFunc, _ = resolveStringer(captured)
			stringerResolved = true
		}
		return stringerFunc, stringerFunc != nil
	}
	lazyErrorer := func() (func() string, bool) {
		if !errorerResolved {
			errorerFunc, _ = resolveErrorer(captured)
			errorerResolved = true
		}
		return errorerFunc, errorerFunc != nil
	}
	lazyGoStringer := func() (func() string, bool) {
		if !gostringerResolved {
			gostringerFunc, _ = resolveGoStringer(captured)
			gostringerResolved = true
		}
		return gostringerFunc, gostringerFunc != nil
	}

	// Always return the wrapper for gig structs - it handles all fmt verbs correctly
	return &gigStructWrapper{
		iface:          iface,
		typeName:       typeName,
		lazyStringer:   lazyStringer,
		lazyErrorer:    lazyErrorer,
		lazyGoStringer: lazyGoStringer,
	}
}

// resolveStringer attempts to resolve the String() method for a value.
// Returns a function that can be called later, and a boolean indicating if found.
// Uses a panic recovery guard because side-channel method invocation via
// ResolveCompiledMethod can fail on method bodies that depend on features
// only wired up by a full program VM (e.g., global state, extCallCache). On
// such failure we return (nil, false) so the caller falls back to default
// formatting instead of propagating a panic through fmt.
func resolveStringer(v Value) (fn func() string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fn, ok = nil, false
		}
	}()
	// Try to call String() method via the global resolver registry
	result, found := callMethod(nil, "String", v)
	if !found {
		// If not found, try with pointer to the value (for pointer receiver methods)
		if rv, rvOK := v.ReflectValue(); rvOK && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = callMethod(nil, "String", ptrValue)
		}
	}

	if !found {
		return nil, false
	}
	str := result.String()
	return func() string { return str }, true
}

// resolveErrorer attempts to resolve the Error() method for a value.
// Returns a function that can be called later, and a boolean indicating if found.
// Uses a panic recovery guard, see resolveStringer for details.
func resolveErrorer(v Value) (fn func() string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fn, ok = nil, false
		}
	}()
	result, found := callMethod(nil, "Error", v)
	if !found {
		// If not found, try with pointer to the value (for pointer receiver methods)
		if rv, rvOK := v.ReflectValue(); rvOK && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = callMethod(nil, "Error", ptrValue)
		}
	}

	if !found {
		return nil, false
	}
	str := result.String()
	return func() string { return str }, true
}

// resolveGoStringer attempts to resolve the GoString() method for a value.
// Returns a function that can be called later, and a boolean indicating if found.
// Uses a panic recovery guard, see resolveStringer for details.
func resolveGoStringer(v Value) (fn func() string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fn, ok = nil, false
		}
	}()
	result, found := callMethod(nil, "GoString", v)
	if !found {
		// If not found, try with pointer to the value (for pointer receiver methods)
		if rv, rvOK := v.ReflectValue(); rvOK && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = callMethod(nil, "GoString", ptrValue)
		}
	}

	if !found {
		return nil, false
	}
	str := result.String()
	return func() string { return str }, true
}

// ErrorValue extracts a Go error from a value.Value.
// If the value is already a native Go error, returns it directly.
// If the value is an interpreter-synthesized struct with an Error() method,
// returns a gigStructWrapper that implements the error interface.
// Otherwise returns nil.
//
// This is the boundary function for generated DirectCall wrappers —
// use it whenever extracting an error-typed parameter from args.
func ErrorValue(v Value) error {
	iface := v.Interface()
	if iface == nil {
		return nil
	}
	// If it's already a Go error (e.g., from fmt.Errorf, errors.New), return as-is
	if e, ok := iface.(error); ok {
		return e
	}
	typeName := isGigStruct(iface)
	if typeName == "" {
		return nil
	}
	// Check if the interpreted type has an Error() method
	errorerFunc, hasError := resolveErrorer(v)
	if !hasError {
		return nil
	}
	// Also capture Stringer/GoStringer lazily so the wrapper is complete.
	stringerFunc, _ := resolveStringer(v)
	gostringerFunc, _ := resolveGoStringer(v)
	constFn := func(s string) func() string { return func() string { return s } }
	errFn := func() string { return errorerFunc() }

	var lazyStringer, lazyErrorer, lazyGoStringer func() (func() string, bool)
	if stringerFunc != nil {
		sf := constFn(stringerFunc())
		lazyStringer = func() (func() string, bool) { return sf, true }
	}
	lazyErrorer = func() (func() string, bool) { return errFn, true }
	if gostringerFunc != nil {
		gf := constFn(gostringerFunc())
		lazyGoStringer = func() (func() string, bool) { return gf, true }
	}
	return &gigStructWrapper{
		iface:          iface,
		typeName:       typeName,
		lazyStringer:   lazyStringer,
		lazyErrorer:    lazyErrorer,
		lazyGoStringer: lazyGoStringer,
	}
}

// ErrorWrap prepares a value.Value for use as a Go error.
// If the value is an interpreter-synthesized struct with an Error() method,
// returns a wrapper that implements the error interface. Otherwise returns
// the raw interface{} value.
//
// Deprecated: Use ErrorValue instead for typed error extraction.
func ErrorWrap(v Value) any {
	if e := ErrorValue(v); e != nil {
		return e
	}
	return v.Interface()
}

// GigErrorsAs implements errors.As semantics for interpreter-defined types.
// It mirrors the standard library's errors.As but uses the interpreter's type
// name registry for matching, since reflect.StructOf types can't implement
// interfaces and have different reflect.Type identities than named Go types.
//
// err is the error value (may be a gigStructWrapper or a native Go error).
// target is a pointer to the target type (e.g., **CustomError as interface{}).
//
// Returns true if the error (or any error in its Unwrap chain) matches target.
func GigErrorsAs(err error, target any) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}

	targetVal := reflect.ValueOf(target)
	targetType := targetVal.Type()
	if targetType.Kind() != reflect.Ptr || targetVal.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}

	elemType := targetType.Elem()

	for {
		// Try matching the current error against the target
		if gigAsMatchValue(err, elemType, targetVal) {
			return true
		}

		// Try unwrapping
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = unwrapper.Unwrap()
		if err == nil {
			return false
		}
	}
}

// gigAsMatchValue checks if an error value matches the target type.
// It handles both native Go types (via reflect.AssignableTo) and
// interpreter-defined types (via gig type name matching).
func gigAsMatchValue(err error, elemType reflect.Type, targetVal reflect.Value) bool {
	errVal := reflect.ValueOf(err)
	errType := errVal.Type()

	// Direct type match: err's type is assignable to target element type
	if errType.AssignableTo(elemType) {
		targetVal.Elem().Set(errVal)
		return true
	}

	// If err is a *gigStructWrapper, try matching by interpreter type name
	if wrapper, ok := err.(*gigStructWrapper); ok {
		// Case 1: target is **StructType (errors.As(&ce) where ce is *CustomError)
		if elemType.Kind() == reflect.Ptr {
			ptrElemType := elemType.Elem()

			// Check if the wrapper's underlying value type is assignable
			ifaceType := reflect.TypeOf(wrapper.iface)
			if ifaceType.AssignableTo(elemType) {
				targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
				return true
			}

			// Match by gig type name: compare wrapper's typeName with target's gig tag
			wrapperTypeName := extractBareTypeName(wrapper.typeName)
			targetTypeName := extractGigTagFromType(ptrElemType)
			if targetTypeName == "" {
				targetTypeName = ptrElemType.Name()
			}
			targetTypeName = extractBareTypeName(targetTypeName)

			if wrapperTypeName != "" && wrapperTypeName == targetTypeName {
				// Type names match — set the target to the underlying value
				if ifaceType.AssignableTo(elemType) {
					targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
					return true
				}
				// Try converting through interface{} if the pointer element is also a gig struct
				if ifaceType.Kind() == reflect.Ptr && ptrElemType.Kind() == reflect.Struct {
					// Both are pointers to structs — try setting the value directly
					if ifaceType.Elem().ConvertibleTo(ptrElemType) {
						converted := reflect.ValueOf(wrapper.iface).Elem().Convert(ptrElemType)
						ptr := reflect.New(ptrElemType)
						ptr.Elem().Set(converted)
						targetVal.Elem().Set(ptr)
						return true
					}
				}
			}
		}

		// Case 2: target is an interface type (e.g., error)
		if elemType.Kind() == reflect.Interface {
			if errType.Implements(elemType) {
				targetVal.Elem().Set(errVal)
				return true
			}
		}

		// Case 3: target is a struct type (value receiver, unlikely for errors)
		if elemType.Kind() == reflect.Struct {
			ifaceType := reflect.TypeOf(wrapper.iface)
			if ifaceType.AssignableTo(elemType) {
				targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
				return true
			}
		}
	}

	// For non-gig errors, try standard interface check
	if elemType.Kind() == reflect.Interface && errType.Implements(elemType) {
		targetVal.Elem().Set(errVal)
		return true
	}

	return false
}

// extractBareTypeName extracts the short type name from a qualified name.
// "known_issues.CustomError" → "CustomError"
// "CustomError" → "CustomError"
func extractBareTypeName(qualName string) string {
	if idx := strings.LastIndex(qualName, "."); idx >= 0 {
		return qualName[idx+1:]
	}
	return qualName
}

// extractGigTagFromType extracts the gig tag value from a reflect.Type.
// Returns "" if not found.
func extractGigTagFromType(rt reflect.Type) string {
	if rt.Kind() != reflect.Struct || rt.NumField() == 0 {
		return ""
	}
	tag := rt.Field(0).Tag.Get("gig")
	if strings.HasPrefix(tag, "#") {
		return tag[1:]
	}
	return tag
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
