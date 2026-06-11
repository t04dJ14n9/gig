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
	return renderGigGoStringStruct(rv, g.typeName, prefix)
}

// goStringValue renders a value in Go-syntax form, detecting gig struct fields
// and rendering them with proper type names recursively.
func goStringValue(v any) string {
	if v == nil {
		return "<nil>"
	}
	rv := reflect.ValueOf(v)
	if typeName := isGigStruct(v); typeName != "" {
		if rendered, ok := formatGigGoStringStruct(rv, typeName); ok {
			return rendered
		}
	}
	if rendered, ok := formatGigGoStringSlice(rv); ok {
		return rendered
	}
	return fmt.Sprintf("%#v", v)
}

func formatGigGoStringStruct(rv reflect.Value, typeName string) (string, bool) {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "<nil>", true
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return "", false
	}
	return renderGigGoStringStruct(rv, typeName, ""), true
}

func renderGigGoStringStruct(rv reflect.Value, typeName, prefix string) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString(typeName)
	sb.WriteByte('{')
	writeGigGoStringFields(&sb, rv)
	sb.WriteByte('}')
	return sb.String()
}

func writeGigGoStringFields(sb *strings.Builder, rv reflect.Value) {
	rt := rv.Type()
	visible := 0
	for i := 0; i < rt.NumField(); i++ {
		if isGigPhantomField(rt.Field(i)) {
			continue
		}
		writeGigGoStringField(sb, rt.Field(i).Name, rv.Field(i), visible > 0)
		visible++
	}
}

func writeGigGoStringField(sb *strings.Builder, name string, field reflect.Value, needsComma bool) {
	if needsComma {
		sb.WriteString(", ")
	}
	sb.WriteString(name)
	sb.WriteByte(':')
	sb.WriteString(goStringValue(field.Interface()))
}

func formatGigGoStringSlice(rv reflect.Value) (string, bool) {
	if rv.Kind() != reflect.Slice {
		return "", false
	}
	elemTypeName := gigSliceElementTypeName(rv)
	if elemTypeName == "" {
		return "", false
	}
	return renderGigGoStringSlice(rv, elemTypeName), true
}

// gigSliceElementTypeName intentionally samples the first element only, matching
// the old formatter behavior. Empty slices fall back to fmt's native %#v output.
func gigSliceElementTypeName(rv reflect.Value) string {
	if rv.Len() == 0 {
		return ""
	}
	return isGigStruct(rv.Index(0).Interface())
}

func renderGigGoStringSlice(rv reflect.Value, elemTypeName string) string {
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
