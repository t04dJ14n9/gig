// Package strings registers the Go standard library strings package.
package strings

import (
	"reflect"
	"strings"

	"gig/importer"
	"gig/value"
)

func init() {
	pkg := importer.RegisterPackage("strings", "strings")

	// Basic string operations
	pkg.AddFunction("Contains", strings.Contains, "", directContains)
	pkg.AddFunction("ContainsAny", strings.ContainsAny, "", directContainsAny)
	pkg.AddFunction("ContainsRune", strings.ContainsRune, "", directContainsRune)
	pkg.AddFunction("Count", strings.Count, "", directCount)
	pkg.AddFunction("EqualFold", strings.EqualFold, "", directEqualFold)
	pkg.AddFunction("Fields", strings.Fields, "", directFields)
	pkg.AddFunction("FieldsFunc", strings.FieldsFunc, "", nil)
	pkg.AddFunction("HasPrefix", strings.HasPrefix, "", directHasPrefix)
	pkg.AddFunction("HasSuffix", strings.HasSuffix, "", directHasSuffix)

	// Index functions
	pkg.AddFunction("Index", strings.Index, "", directIndex)
	pkg.AddFunction("IndexAny", strings.IndexAny, "", directIndexAny)
	pkg.AddFunction("IndexByte", strings.IndexByte, "", directIndexByte)
	pkg.AddFunction("IndexFunc", strings.IndexFunc, "", nil)
	pkg.AddFunction("IndexRune", strings.IndexRune, "", directIndexRune)
	pkg.AddFunction("LastIndex", strings.LastIndex, "", directLastIndex)
	pkg.AddFunction("LastIndexAny", strings.LastIndexAny, "", directLastIndexAny)
	pkg.AddFunction("LastIndexByte", strings.LastIndexByte, "", directLastIndexByte)
	pkg.AddFunction("LastIndexFunc", strings.LastIndexFunc, "", nil)

	// String manipulation
	pkg.AddFunction("Clone", strings.Clone, "", directClone)
	pkg.AddFunction("Compare", strings.Compare, "", directCompare)
	pkg.AddFunction("Cut", strings.Cut, "", nil)
	pkg.AddFunction("CutPrefix", strings.CutPrefix, "", nil)
	pkg.AddFunction("CutSuffix", strings.CutSuffix, "", nil)
	pkg.AddFunction("Join", strings.Join, "", directJoin)
	pkg.AddFunction("Map", strings.Map, "", nil)
	pkg.AddFunction("Repeat", strings.Repeat, "", directRepeat)
	pkg.AddFunction("Replace", strings.Replace, "", directReplace)
	pkg.AddFunction("ReplaceAll", strings.ReplaceAll, "", directReplaceAll)
	pkg.AddFunction("Split", strings.Split, "", directSplit)
	pkg.AddFunction("SplitAfter", strings.SplitAfter, "", directSplitAfter)
	pkg.AddFunction("SplitAfterN", strings.SplitAfterN, "", directSplitAfterN)
	pkg.AddFunction("SplitN", strings.SplitN, "", directSplitN)
	pkg.AddFunction("ToLower", strings.ToLower, "", directToLower)
	pkg.AddFunction("ToLowerSpecial", strings.ToLowerSpecial, "", nil)
	pkg.AddFunction("ToTitle", strings.ToTitle, "", directToTitle)
	pkg.AddFunction("ToTitleSpecial", strings.ToTitleSpecial, "", nil)
	pkg.AddFunction("ToUpper", strings.ToUpper, "", directToUpper)
	pkg.AddFunction("ToUpperSpecial", strings.ToUpperSpecial, "", nil)
	pkg.AddFunction("ToValidUTF8", strings.ToValidUTF8, "", nil)
	pkg.AddFunction("Trim", strings.Trim, "", directTrim)
	pkg.AddFunction("TrimFunc", strings.TrimFunc, "", nil)
	pkg.AddFunction("TrimLeft", strings.TrimLeft, "", directTrimLeft)
	pkg.AddFunction("TrimLeftFunc", strings.TrimLeftFunc, "", nil)
	pkg.AddFunction("TrimPrefix", strings.TrimPrefix, "", directTrimPrefix)
	pkg.AddFunction("TrimRight", strings.TrimRight, "", directTrimRight)
	pkg.AddFunction("TrimRightFunc", strings.TrimRightFunc, "", nil)
	pkg.AddFunction("TrimSpace", strings.TrimSpace, "", directTrimSpace)
	pkg.AddFunction("TrimSuffix", strings.TrimSuffix, "", directTrimSuffix)

	// Builder type
	pkg.AddType("Builder", nil, "efficient string builder")
	pkg.AddType("Reader", nil, "string reader")
	pkg.AddType("Replacer", nil, "string replacer")
}

// Direct wrappers for common functions (avoid reflect.Call)

func directContains(args []value.Value) value.Value {
	return value.MakeBool(strings.Contains(args[0].String(), args[1].String()))
}

func directContainsAny(args []value.Value) value.Value {
	return value.MakeBool(strings.ContainsAny(args[0].String(), args[1].String()))
}

func directContainsRune(args []value.Value) value.Value {
	return value.MakeBool(strings.ContainsRune(args[0].String(), rune(args[1].Int())))
}

func directCount(args []value.Value) value.Value {
	return value.MakeInt(int64(strings.Count(args[0].String(), args[1].String())))
}

func directEqualFold(args []value.Value) value.Value {
	return value.MakeBool(strings.EqualFold(args[0].String(), args[1].String()))
}

func directFields(args []value.Value) value.Value {
	fields := strings.Fields(args[0].String())
	result := make([]value.Value, len(fields))
	for i, f := range fields {
		result[i] = value.MakeString(f)
	}
	return value.FromInterface(result)
}

func directHasPrefix(args []value.Value) value.Value {
	return value.MakeBool(strings.HasPrefix(args[0].String(), args[1].String()))
}

func directHasSuffix(args []value.Value) value.Value {
	return value.MakeBool(strings.HasSuffix(args[0].String(), args[1].String()))
}

func directIndex(args []value.Value) value.Value {
	return value.MakeInt(int64(strings.Index(args[0].String(), args[1].String())))
}

func directIndexAny(args []value.Value) value.Value {
	return value.MakeInt(int64(strings.IndexAny(args[0].String(), args[1].String())))
}

func directIndexByte(args []value.Value) value.Value {
	return value.MakeInt(int64(strings.IndexByte(args[0].String(), byte(args[1].Int()))))
}

func directIndexRune(args []value.Value) value.Value {
	return value.MakeInt(int64(strings.IndexRune(args[0].String(), rune(args[1].Int()))))
}

func directLastIndex(args []value.Value) value.Value {
	return value.MakeInt(int64(strings.LastIndex(args[0].String(), args[1].String())))
}

func directLastIndexAny(args []value.Value) value.Value {
	return value.MakeInt(int64(strings.LastIndexAny(args[0].String(), args[1].String())))
}

func directLastIndexByte(args []value.Value) value.Value {
	return value.MakeInt(int64(strings.LastIndexByte(args[0].String(), byte(args[1].Int()))))
}

func directClone(args []value.Value) value.Value {
	return value.MakeString(strings.Clone(args[0].String()))
}

func directCompare(args []value.Value) value.Value {
	return value.MakeInt(int64(strings.Compare(args[0].String(), args[1].String())))
}

func directJoin(args []value.Value) value.Value {
	// args[0] is []string, args[1] is sep
	if slice, ok := args[0].ReflectValue(); ok {
		if slice.Kind() == reflect.Slice {
			strs := make([]string, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				strs[i] = slice.Index(i).String()
			}
			return value.MakeString(strings.Join(strs, args[1].String()))
		}
	}
	return value.MakeString("")
}

func directRepeat(args []value.Value) value.Value {
	return value.MakeString(strings.Repeat(args[0].String(), int(args[1].Int())))
}

func directReplace(args []value.Value) value.Value {
	return value.MakeString(strings.Replace(args[0].String(), args[1].String(), args[2].String(), int(args[3].Int())))
}

func directReplaceAll(args []value.Value) value.Value {
	return value.MakeString(strings.ReplaceAll(args[0].String(), args[1].String(), args[2].String()))
}

func directSplit(args []value.Value) value.Value {
	parts := strings.Split(args[0].String(), args[1].String())
	result := make([]value.Value, len(parts))
	for i, p := range parts {
		result[i] = value.MakeString(p)
	}
	return value.FromInterface(result)
}

func directSplitAfter(args []value.Value) value.Value {
	parts := strings.SplitAfter(args[0].String(), args[1].String())
	result := make([]value.Value, len(parts))
	for i, p := range parts {
		result[i] = value.MakeString(p)
	}
	return value.FromInterface(result)
}

func directSplitAfterN(args []value.Value) value.Value {
	parts := strings.SplitAfterN(args[0].String(), args[1].String(), int(args[2].Int()))
	result := make([]value.Value, len(parts))
	for i, p := range parts {
		result[i] = value.MakeString(p)
	}
	return value.FromInterface(result)
}

func directSplitN(args []value.Value) value.Value {
	parts := strings.SplitN(args[0].String(), args[1].String(), int(args[2].Int()))
	result := make([]value.Value, len(parts))
	for i, p := range parts {
		result[i] = value.MakeString(p)
	}
	return value.FromInterface(result)
}

func directToLower(args []value.Value) value.Value {
	return value.MakeString(strings.ToLower(args[0].String()))
}

func directToTitle(args []value.Value) value.Value {
	return value.MakeString(strings.ToTitle(args[0].String()))
}

func directToUpper(args []value.Value) value.Value {
	return value.MakeString(strings.ToUpper(args[0].String()))
}

func directTrim(args []value.Value) value.Value {
	return value.MakeString(strings.Trim(args[0].String(), args[1].String()))
}

func directTrimLeft(args []value.Value) value.Value {
	return value.MakeString(strings.TrimLeft(args[0].String(), args[1].String()))
}

func directTrimPrefix(args []value.Value) value.Value {
	return value.MakeString(strings.TrimPrefix(args[0].String(), args[1].String()))
}

func directTrimRight(args []value.Value) value.Value {
	return value.MakeString(strings.TrimRight(args[0].String(), args[1].String()))
}

func directTrimSpace(args []value.Value) value.Value {
	return value.MakeString(strings.TrimSpace(args[0].String()))
}

func directTrimSuffix(args []value.Value) value.Value {
	return value.MakeString(strings.TrimSuffix(args[0].String(), args[1].String()))
}
