package value

import (
	"fmt"
	"reflect"
	"strings"
)

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
		// Check if the field is itself a gig struct and render with type name.
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
	// Check if this is a gig struct.
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
	// For slices of gig structs, render with type name.
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
