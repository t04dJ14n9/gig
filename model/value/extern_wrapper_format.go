package value

import (
	"fmt"
	"reflect"
)

func (g *gigStructWrapper) Format(f fmt.State, verb rune) {
	switch verb {
	case 'T':
		_, _ = fmt.Fprint(f, g.typeName)
	case 'v':
		if f.Flag('#') {
			// Go-syntax formatting should prefer a script-defined GoString
			// method before falling back to defaultGoString's native shape.
			_, _ = fmt.Fprint(f, g.GoString())
			return
		}
		if f.Flag('+') {
			g.formatNamedFields(f)
			return
		}
		g.formatPlainValue(f)
	case 's':
		g.formatVerbS(f)
	default:
		_, _ = fmt.Fprintf(f, "%"+string(verb), g.iface)
	}
}

func (g *gigStructWrapper) formatVerbS(f fmt.State) {
	// Native fmt gives Error() priority over String() for %v and %s. Gig
	// wrappers must keep that order because script-defined Error methods are
	// commonly used to satisfy external error formatting contracts.
	if g.formatErrorOrString(f) {
		return
	}
	_, _ = fmt.Fprintf(f, "%v", g.iface)
}

func (g *gigStructWrapper) formatPlainValue(f fmt.State) {
	// Plain %v is the hottest formatter path. Keep the Error/String checks
	// inline here so the common case stays close to native fmt's dispatch.
	if fn, ok := g.tryErrorer(); ok {
		_, _ = fmt.Fprint(f, fn())
		return
	}
	if fn, ok := g.tryStringer(); ok {
		_, _ = fmt.Fprint(f, fn())
		return
	}

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
		formatPlainStructFields(f, rv)
		return
	}
	_, _ = fmt.Fprintf(f, "%v", g.iface)
}

func (g *gigStructWrapper) formatNamedFields(f fmt.State) {
	// %+v prints the synthesized struct without the anonymous reflect type name,
	// but with field names, matching native fmt's useful debugging shape.
	rv := reflect.ValueOf(g.iface)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			_, _ = fmt.Fprint(f, "<nil>")
			return
		}
		rv = rv.Elem()
		_, _ = fmt.Fprint(f, "&")
	}
	if rv.Kind() != reflect.Struct {
		_, _ = fmt.Fprintf(f, "%+v", g.iface)
		return
	}
	formatNamedStructFields(f, rv)
}

func (g *gigStructWrapper) formatErrorOrString(f fmt.State) bool {
	if fn, ok := g.tryErrorer(); ok {
		_, _ = fmt.Fprint(f, fn())
		return true
	}
	if fn, ok := g.tryStringer(); ok {
		_, _ = fmt.Fprint(f, fn())
		return true
	}
	return false
}

func formatPlainStructFields(f fmt.State, rv reflect.Value) {
	// The hidden gig phantom field only exists to preserve script type identity;
	// it must not leak into fmt output seen by external libraries.
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
}

func formatNamedStructFields(f fmt.State, rv reflect.Value) {
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
}
