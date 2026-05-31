package value

import (
	"reflect"
	"strings"
)

func isGigPhantomField(field reflect.StructField) bool {
	return field.Name == "gigType" && field.PkgPath == "gig/internal" && field.Tag.Get("gig") != ""
}

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
				// Fallback: check PkgPath of unexported fields for "#TypeName"
				for i := 0; i < elemType.NumField(); i++ {
					pkgPath := elemType.Field(i).PkgPath
					if idx := strings.LastIndex(pkgPath, "#"); idx >= 0 {
						qualName := pkgPath[idx+1:]
						if dotIdx := strings.LastIndex(qualName, "."); dotIdx >= 0 {
							return qualName[dotIdx+1:]
						}
						return qualName
					}
				}
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
		// Fallback: check PkgPath of unexported fields for "#TypeName"
		for i := 0; i < rt.NumField(); i++ {
			pkgPath := rt.Field(i).PkgPath
			if idx := strings.LastIndex(pkgPath, "#"); idx >= 0 {
				qualName := pkgPath[idx+1:]
				if dotIdx := strings.LastIndex(qualName, "."); dotIdx >= 0 {
					return qualName[dotIdx+1:]
				}
				return qualName
			}
		}
		return ""
	}
	if strings.HasPrefix(gigTag, "#") {
		return gigTag[1:]
	}
	return gigTag
}

// FmtWrap prepares a value.Value for passing to fmt.* functions.
// If the value is an interpreter-synthesized struct with compiled methods

func typeContainsGigStruct(t reflect.Type) bool {
	if t == nil {
		return false
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		return extractGigTagFromType(t) != ""
	case reflect.Slice, reflect.Array:
		return typeContainsGigStruct(t.Elem())
	case reflect.Map:
		return typeContainsGigStruct(t.Key()) || typeContainsGigStruct(t.Elem())
	case reflect.Interface:
		return false
	default:
		return false
	}
}

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
