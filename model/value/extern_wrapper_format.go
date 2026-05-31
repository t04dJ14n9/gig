package value

import (
	"fmt"
	"reflect"
)

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
		// %#v: Go-syntax representation. Prefer GoString() if defined;
		// otherwise fall back to defaultGoString which renders the struct
		// with commas between fields, matching native fmt.
		if verb == 'v' && f.Flag('#') {
			_, _ = fmt.Fprint(f, g.GoString())
			return
		}
		// %+v: fields shown with names.
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
