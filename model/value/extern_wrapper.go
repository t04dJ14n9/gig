package value

import "fmt"

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
