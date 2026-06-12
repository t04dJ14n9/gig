package value

import (
	"fmt"
	"reflect"
	"strings"
)

func (g *gigStructError) defaultGoString() string {
	rv := reflect.ValueOf(g.iface)
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
		if visible > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(rt.Field(i).Name)
		sb.WriteByte(':')
		sb.WriteString(fmt.Sprintf("%#v", rv.Field(i).Interface()))
		visible++
	}
}

func isGigPhantomField(field reflect.StructField) bool {
	return field.Name == "gigType" && field.PkgPath == "gig/internal" && field.Tag.Get("gig") != ""
}
