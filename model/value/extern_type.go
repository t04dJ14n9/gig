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
	return gigStructNameFromType(baseReflectType(reflect.TypeOf(v)))
}

func baseReflectType(rt reflect.Type) reflect.Type {
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt
}

func gigStructNameFromType(rt reflect.Type) string {
	if rt.Kind() != reflect.Struct {
		return ""
	}
	if rt.NumField() == 0 {
		return ""
	}
	gigTag := rt.Field(0).Tag.Get("gig")
	if gigTag != "" {
		return normalizeGigTag(gigTag)
	}
	return gigStructNameFromPkgPath(rt)
}

func normalizeGigTag(gigTag string) string {
	if strings.HasPrefix(gigTag, "#") {
		return gigTag[1:]
	}
	return gigTag
}

func gigStructNameFromPkgPath(rt reflect.Type) string {
	// Some synthesized structs carry the script type name in the unexported
	// field PkgPath because reflect.StructOf cannot attach real methods.
	for i := 0; i < rt.NumField(); i++ {
		pkgPath := rt.Field(i).PkgPath
		if idx := strings.LastIndex(pkgPath, "#"); idx >= 0 {
			return extractBareTypeName(pkgPath[idx+1:])
		}
	}
	return ""
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
	return normalizeGigTag(rt.Field(0).Tag.Get("gig"))
}
