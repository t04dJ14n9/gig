package value

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// gigSequenceFormatter handles fmt formatting for slices/arrays containing gig
// structs. Native fmt would see anonymous reflect.StructOf elements with no
// methods, so we format elements through FmtWrap recursively.
type gigSequenceFormatter struct {
	rv reflect.Value
}

func (s *gigSequenceFormatter) Format(f fmt.State, verb rune) {
	if verb == 'v' && f.Flag('#') {
		_, _ = fmt.Fprint(f, s.GoString())
		return
	}
	if verb == 'v' && !f.Flag('+') {
		_, _ = fmt.Fprint(f, formatSequencePlain(s.rv))
		return
	}
	_, _ = fmt.Fprintf(f, fmt.FormatString(f, verb), s.rv.Interface())
}

func (s *gigSequenceFormatter) GoString() string {
	elemTypeName := ""
	if s.rv.Len() > 0 {
		elemTypeName = isGigStruct(s.rv.Index(0).Interface())
	}
	if elemTypeName == "" && s.rv.Type().Elem().Kind() == reflect.Struct {
		elemTypeName = extractGigTagFromType(s.rv.Type().Elem())
	}
	if elemTypeName == "" {
		return fmt.Sprintf("%#v", s.rv.Interface())
	}

	var sb strings.Builder
	if s.rv.Kind() == reflect.Array {
		sb.WriteByte('[')
		_, _ = fmt.Fprintf(&sb, "%d", s.rv.Len())
		sb.WriteByte(']')
	} else {
		sb.WriteString("[]")
	}
	sb.WriteString(elemTypeName)
	sb.WriteByte('{')
	for i := 0; i < s.rv.Len(); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		elem := s.rv.Index(i).Interface()
		// If element has GoString() via compiled method, use it
		elemVal := MakeFromReflect(s.rv.Index(i))
		if fn, ok := resolveGoStringer(elemVal); ok {
			sb.WriteString(fn())
		} else {
			sb.WriteString(goStringValue(elem))
		}
	}
	sb.WriteByte('}')
	return sb.String()
}

type gigMapFormatter struct {
	rv reflect.Value
}

func (m *gigMapFormatter) Format(f fmt.State, verb rune) {
	if verb == 'v' && !f.Flag('#') && !f.Flag('+') {
		_, _ = fmt.Fprint(f, formatMapPlain(m.rv))
		return
	}
	_, _ = fmt.Fprintf(f, fmt.FormatString(f, verb), m.rv.Interface())
}

func formatReflectPlain(rv reflect.Value) string {
	if !rv.IsValid() {
		return "<nil>"
	}
	return fmt.Sprint(wrapFmtReflect(rv))
}

func formatSequencePlain(rv reflect.Value) string {
	if !rv.IsValid() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(formatReflectPlain(rv.Index(i)))
	}
	sb.WriteByte(']')
	return sb.String()
}

func formatMapPlain(rv reflect.Value) string {
	if !rv.IsValid() {
		return "map[]"
	}
	keys := rv.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(wrapFmtReflect(keys[i])) < fmt.Sprint(wrapFmtReflect(keys[j]))
	})

	var sb strings.Builder
	sb.WriteString("map[")
	for i, key := range keys {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(formatReflectPlain(key))
		sb.WriteByte(':')
		sb.WriteString(formatReflectPlain(rv.MapIndex(key)))
	}
	sb.WriteByte(']')
	return sb.String()
}

func wrapFmtReflect(rv reflect.Value) (out any) {
	defer func() {
		if recover() != nil {
			out = fmt.Sprint(rv)
		}
	}()
	if !rv.IsValid() {
		return nil
	}
	if !rv.CanInterface() {
		return fmt.Sprint(rv)
	}
	return FmtWrap(MakeFromReflect(rv))
}
