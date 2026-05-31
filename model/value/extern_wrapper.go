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
	iface          any    // the underlying struct value (clean, no phantom fields)
	typeName       string // qualified type name from gig tag (e.g., "pkg.Type")
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
	// For nil pointers, Go's fmt prints "<nil>" when the type implements GoStringer.
	prefix := ""
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "<nil>"
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
	visible := 0
	for i := 0; i < rt.NumField(); i++ {
		if isGigPhantomField(rt.Field(i)) {
			continue
		}
		if visible > 0 {
			sb.WriteString(", ")
		}
		visible++
		sb.WriteString(rt.Field(i).Name)
		sb.WriteByte(':')
		fieldVal := rv.Field(i).Interface()
		// Check if the field is itself a gig struct and render with type name
		sb.WriteString(goStringValue(fieldVal))
	}
	sb.WriteByte('}')
	return sb.String()
}

// goStringValue renders a value in Go-syntax form, detecting gig struct fields
// and rendering them with proper type names recursively.
func goStringValue(v any) string {
	if v == nil {
		return "<nil>"
	}
	rv := reflect.ValueOf(v)
	// Check if this is a gig struct
	if typeName := isGigStruct(v); typeName != "" {
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return "<nil>"
			}
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			var sb strings.Builder
			sb.WriteString(typeName)
			sb.WriteByte('{')
			rt := rv.Type()
			visible := 0
			for i := 0; i < rt.NumField(); i++ {
				if isGigPhantomField(rt.Field(i)) {
					continue
				}
				if visible > 0 {
					sb.WriteString(", ")
				}
				visible++
				sb.WriteString(rt.Field(i).Name)
				sb.WriteByte(':')
				sb.WriteString(goStringValue(rv.Field(i).Interface()))
			}
			sb.WriteByte('}')
			return sb.String()
		}
	}
	// For slices of gig structs, render with type name
	if rv.Kind() == reflect.Slice {
		elemTypeName := ""
		if rv.Len() > 0 {
			elemTypeName = isGigStruct(rv.Index(0).Interface())
		}
		if elemTypeName != "" {
			var sb strings.Builder
			sb.WriteString("[]")
			sb.WriteString(elemTypeName)
			sb.WriteByte('{')
			for i := 0; i < rv.Len(); i++ {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(goStringValue(rv.Index(i).Interface()))
			}
			sb.WriteByte('}')
			return sb.String()
		}
	}
	return fmt.Sprintf("%#v", v)
}

func (g *gigStructWrapper) Format(f fmt.State, verb rune) {
	switch verb {
	case 'T':
		_, _ = fmt.Fprint(f, g.typeName)
	case 'v', 's':
		if verb == 's' || (verb == 'v' && !f.Flag('#') && !f.Flag('+')) {
			// Go's fmt checks error before Stringer for %v and %s.
			// We match that priority: try Error() first, then String().
			if fn, ok := g.tryErrorer(); ok {
				_, _ = fmt.Fprint(f, fn())
				return
			}
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
					visible := 0
					for i := 0; i < rt.NumField(); i++ {
						if isGigPhantomField(rt.Field(i)) {
							continue
						}
						if visible > 0 {
							_, _ = fmt.Fprint(f, " ")
						}
						visible++
						_, _ = fmt.Fprint(f, formatReflectPlain(rv.Field(i)))
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
			visible := 0
			for i := 0; i < rt.NumField(); i++ {
				if isGigPhantomField(rt.Field(i)) {
					continue
				}
				if visible > 0 {
					_, _ = fmt.Fprint(f, " ")
				}
				visible++
				_, _ = fmt.Fprintf(f, "%s:%s", rt.Field(i).Name, formatReflectPlain(rv.Field(i)))
			}
			_, _ = fmt.Fprint(f, "}")
			return
		}
		_, _ = fmt.Fprintf(f, "%v", g.iface)
	default:
		_, _ = fmt.Fprintf(f, "%"+string(verb), g.iface)
	}
}
