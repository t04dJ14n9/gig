package gentool

import (
	"fmt"
	"go/types"
	"strings"
)

// --- Argument extraction ---

// extractArg generates the Go expression to extract a typed value from a value.Value
// argument for passing to a native Go function. Returns "" if the type is unsupported.
func extractArg(t types.Type, valExpr string, pkgRef string) string {
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		if obj.Pkg() == nil && obj.Name() == errorTypeName {
			// Use value.ErrorValue to handle interpreter-defined types
			// with Error() method that can't satisfy the error interface
			// because reflect.StructOf types can't have methods.
			return fmt.Sprintf("value.ErrorValue(%s)", valExpr)
		}
	}

	// Handle type aliases (Go 1.23+), including builtin 'any'
	if alias, ok := t.(*types.Alias); ok {
		obj := alias.Obj()
		if obj.Pkg() == nil {
			// Builtin alias: 'error' or 'any' (= interface{})
			if obj.Name() == errorTypeName {
				// Use value.ErrorValue for the same reason as Named error above.
				return fmt.Sprintf("value.ErrorValue(%s)", valExpr)
			}
			// 'any' and other builtin aliases: pass raw value
			return fmt.Sprintf("%s.Interface()", valExpr)
		}
	}

	if named, ok := t.(*types.Named); ok {
		return extractNamedArg(named, t.Underlying(), valExpr, pkgRef)
	}

	// Handle type aliases (Go 1.23+), including builtin 'any'
	if alias, ok := t.(*types.Alias); ok {
		return extractAliasArg(alias, t.Underlying(), valExpr, pkgRef)
	}

	// For non-named types (e.g., *pkg.Type pointers), handle with pkgRef
	return extractUnderlyingWithPkgRef(t, valExpr, pkgRef)
}

// extractNamedArg handles argument extraction for named types.
func extractNamedArg(named *types.Named, underlying types.Type, valExpr string, pkgRef string) string {
	obj := named.Obj()
	pkg := obj.Pkg()

	// Same-package named type with basic underlying: cast via basic extraction
	if pkg != nil && pkg.Path() == currentPkgPath {
		if bt, ok := underlying.(*types.Basic); ok {
			basicExpr := extractBasic(bt, valExpr)
			if basicExpr == "" {
				return ""
			}
			namedName := resolveTypeName(named, pkgRef)
			if namedName == "" {
				return ""
			}
			return fmt.Sprintf("%s(%s)", namedName, basicExpr)
		}
		// For non-basic same-package types (interfaces, structs, pointers, etc.),
		// use .Interface().(TypeName) type assertion
		namedName := resolveTypeName(named, pkgRef)
		if namedName != "" {
			return fmt.Sprintf("%s.Interface().(%s)", valExpr, namedName)
		}
		return extractUnderlyingWithPkgRef(underlying, valExpr, pkgRef)
	}

	// Cross-package named type: check if it has a basic underlying type
	// If so, use a cast (e.g., time.Duration(args[i].Int())) instead of
	// .Interface().(time.Duration) which panics because the VM stores
	// named-basic types as their underlying basic kind (e.g., int64).
	if bt, ok := underlying.(*types.Basic); ok {
		basicExpr := extractBasic(bt, valExpr)
		if basicExpr == "" {
			return ""
		}
		qualifiedName := resolveQualifiedName(named, pkgRef)
		if qualifiedName != "" {
			return fmt.Sprintf("%s(%s)", qualifiedName, basicExpr)
		}
		return ""
	}

	// For non-basic cross-package types (structs, interfaces, pointers, etc.),
	// use .Interface().(pkg.Type) type assertion
	qualifiedName := resolveQualifiedName(named, pkgRef)
	if qualifiedName != "" {
		return fmt.Sprintf("%s.Interface().(%s)", valExpr, qualifiedName)
	}
	return ""
}

// extractAliasArg handles argument extraction for alias types.
func extractAliasArg(alias *types.Alias, underlying types.Type, valExpr string, pkgRef string) string {
	obj := alias.Obj()
	pkg := obj.Pkg()

	if pkg == nil {
		// Builtin alias (e.g. 'any' = interface{}, 'error' = interface{Error() string})
		// Already handled above for 'error'; for 'any' and others, extract as interface{}
		return extractUnderlyingWithPkgRef(underlying, valExpr, pkgRef)
	}

	if pkg.Path() == currentPkgPath {
		if bt, ok := underlying.(*types.Basic); ok {
			basicExpr := extractBasic(bt, valExpr)
			if basicExpr == "" {
				return ""
			}
			aliasName := resolveTypeName(alias, pkgRef)
			if aliasName == "" {
				return ""
			}
			return fmt.Sprintf("%s(%s)", aliasName, basicExpr)
		}
		aliasName := resolveTypeName(alias, pkgRef)
		if aliasName != "" {
			return fmt.Sprintf("%s.Interface().(%s)", valExpr, aliasName)
		}
		return extractUnderlyingWithPkgRef(underlying, valExpr, pkgRef)
	}

	// Cross-package alias: check if basic underlying for cast
	if bt, ok := underlying.(*types.Basic); ok {
		basicExpr := extractBasic(bt, valExpr)
		if basicExpr == "" {
			return ""
		}
		return fmt.Sprintf("%s.%s(%s)", sanitizePkgName(pkg.Path()), obj.Name(), basicExpr)
	}
	return fmt.Sprintf("%s.Interface().(%s.%s)", valExpr, sanitizePkgName(pkg.Path()), obj.Name())
}

// extractUnderlyingWithPkgRef generates the Go expression to extract a value from
// a value.Value for non-named underlying types (pointers, slices, maps, etc.).
func extractUnderlyingWithPkgRef(t types.Type, valExpr string, pkgRef string) string {
	switch ut := t.(type) {
	case *types.Basic:
		return extractBasic(ut, valExpr)
	case *types.Slice:
		// For basic element types, use the optimized extractSlice
		if _, ok := ut.Elem().Underlying().(*types.Basic); ok {
			return extractSlice(ut, valExpr, pkgRef)
		}
		// For non-basic element types ([][]byte, []*T, etc.), use Interface()
		// without a type assertion — the VM stores these via reflection
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Interface:
		if ut.NumMethods() == 0 {
			return fmt.Sprintf("%s.Interface()", valExpr)
		}
		// Non-empty unnamed interface: generate type assertion like
		// args[0].Interface().(interface{ Printf(string, ...interface{}) })
		ifaceName := resolveInterfaceTypeName(ut, pkgRef)
		if ifaceName != "" {
			return fmt.Sprintf("%s.Interface().(%s)", valExpr, ifaceName)
		}
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Pointer:
		// If the pointer element is a named type, generate a typed assertion like .Interface().(*pkg.Type)
		if named, ok := ut.Elem().(*types.Named); ok {
			qualName := resolveQualifiedName(named, pkgRef)
			if qualName != "" {
				return fmt.Sprintf("%s.Interface().(*%s)", valExpr, qualName)
			}
		}
		// If the pointer element is a basic type, generate .Interface().(*basicType)
		if bt, ok := ut.Elem().(*types.Basic); ok {
			return fmt.Sprintf("%s.Interface().(*%s)", valExpr, bt.Name())
		}
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Struct:
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Map:
		keyName := resolveTypeName(ut.Key(), pkgRef)
		elemName := resolveTypeName(ut.Elem(), pkgRef)
		if keyName != "" && elemName != "" {
			return fmt.Sprintf("%s.Interface().(map[%s]%s)", valExpr, keyName, elemName)
		}
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Chan:
		// chan T: extract via .Interface().(chan T)
		elemName := resolveTypeName(ut.Elem(), pkgRef)
		if elemName != "" {
			var dirStr string
			switch ut.Dir() {
			case types.SendRecv:
				dirStr = fmt.Sprintf("chan %s", elemName)
			case types.SendOnly:
				dirStr = fmt.Sprintf("chan<- %s", elemName)
			case types.RecvOnly:
				dirStr = fmt.Sprintf("<-chan %s", elemName)
			}
			if dirStr != "" {
				return fmt.Sprintf("%s.Interface().(%s)", valExpr, dirStr)
			}
		}
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Signature:
		// func(T) R: extract via .Interface().(func(T) R)
		funcTypeName := resolveFuncTypeName(ut, pkgRef)
		if funcTypeName != "" {
			return fmt.Sprintf("%s.Interface().(%s)", valExpr, funcTypeName)
		}
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Array:
		// [N]T: extract via .Interface().([N]T)
		arrTypeName := resolveArrayTypeName(ut, pkgRef)
		if arrTypeName != "" {
			return fmt.Sprintf("%s.Interface().(%s)", valExpr, arrTypeName)
		}
		return fmt.Sprintf("%s.Interface()", valExpr)
	default:
		return ""
	}
}

// resolveQualifiedName returns the fully-qualified Go type name for a named type,
// suitable for use in a type assertion like args[i].Interface().(pkg.Type).
// For pointer types it returns *pkg.Type, etc.
func resolveQualifiedName(named *types.Named, pkgRef string) string {
	obj := named.Obj()
	pkg := obj.Pkg()
	if pkg == nil {
		// Builtin type
		return obj.Name()
	}
	if pkg.Path() == currentPkgPath {
		return fmt.Sprintf("%s.%s", pkgRef, obj.Name())
	}
	// Cross-package: use the import alias
	return fmt.Sprintf("%s.%s", sanitizePkgName(pkg.Path()), obj.Name())
}

// extractBasic generates the Go expression to extract a basic typed value from a value.Value.
// Returns "" for unsupported types (e.g., UnsafePointer).
func extractBasic(bt *types.Basic, valExpr string) string {
	info := bt.Info()
	kind := bt.Kind()

	switch {
	case kind == types.Bool:
		return fmt.Sprintf("%s.Bool()", valExpr)
	case info&types.IsInteger != 0 && info&types.IsUnsigned != 0:
		switch kind {
		case types.Uint8:
			return fmt.Sprintf("byte(%s.Uint())", valExpr)
		case types.Uint16:
			return fmt.Sprintf("uint16(%s.Uint())", valExpr)
		case types.Uint32:
			return fmt.Sprintf("uint32(%s.Uint())", valExpr)
		case types.Uint64:
			return fmt.Sprintf("%s.Uint()", valExpr)
		case types.Uint:
			return fmt.Sprintf("uint(%s.Uint())", valExpr)
		case types.Uintptr:
			return fmt.Sprintf("uintptr(%s.Uint())", valExpr)
		default:
			return fmt.Sprintf("uint(%s.Uint())", valExpr)
		}
	case info&types.IsInteger != 0:
		switch kind {
		case types.Int8:
			return fmt.Sprintf("int8(%s.Int())", valExpr)
		case types.Int16:
			return fmt.Sprintf("int16(%s.Int())", valExpr)
		case types.Int32:
			return fmt.Sprintf("int32(%s.Int())", valExpr)
		case types.Int64:
			return fmt.Sprintf("%s.Int()", valExpr)
		case types.Int:
			return fmt.Sprintf("int(%s.Int())", valExpr)
		default:
			return fmt.Sprintf("int(%s.Int())", valExpr)
		}
	case info&types.IsFloat != 0:
		switch kind {
		case types.Float32:
			return fmt.Sprintf("float32(%s.Float())", valExpr)
		case types.Float64:
			return fmt.Sprintf("%s.Float()", valExpr)
		default:
			return fmt.Sprintf("%s.Float()", valExpr)
		}
	case info&types.IsComplex != 0:
		switch kind {
		case types.Complex64:
			return fmt.Sprintf("complex64(%s.Complex())", valExpr)
		default:
			return fmt.Sprintf("%s.Complex()", valExpr)
		}
	case info&types.IsString != 0:
		return fmt.Sprintf("%s.String()", valExpr)
	case kind == types.UnsafePointer:
		return ""
	default:
		return fmt.Sprintf("%s.Interface()", valExpr)
	}
}

// extractSlice generates the Go expression to extract a slice from a value.Value.
func extractSlice(st *types.Slice, valExpr string, pkgRef string) string {
	if bt, ok := st.Elem().Underlying().(*types.Basic); ok {
		switch bt.Kind() {
		case types.Byte:
			// Use native KindBytes accessor — zero reflection for []byte params
			// Falls back to Interface() for KindReflect values, with nil-safe assertion
			return fmt.Sprintf("func() []byte { if b, ok := (%s).Bytes(); ok { return b }; v := (%s).Interface(); if v == nil { return nil }; return v.([]byte) }()", valExpr, valExpr)
		case types.String:
			// []string is stored as []string in the VM — safe to assert directly
			return fmt.Sprintf("%s.Interface().([]string)", valExpr)
		default:
			// For integer slices that need conversion ([]int, []int32, etc.),
			// the caller MUST use emitIntSliceExtraction instead of this function
			// so that writeback code can be generated. This branch handles
			// the simple assertion cases ([]int64, []float64, etc.).
			elemName := resolveTypeName(st.Elem(), pkgRef)
			if elemName == "" {
				return fmt.Sprintf("%s.Interface()", valExpr)
			}
			return fmt.Sprintf("%s.Interface().([]%s)", valExpr, elemName)
		}
	}
	return fmt.Sprintf("%s.Interface()", valExpr)
}

// --- Integer slice conversion and writeback ---

// needsIntSliceConversion returns true if a slice type's elements are integers
// that are NOT int64/uint64 (i.e., the VM stores them as []int64 but the Go
// function expects a different int width like []int, []int32, []uint16, etc.).
// These parameters need conversion AND post-call writeback for in-place mutation.
func needsIntSliceConversion(st *types.Slice) bool {
	bt, ok := st.Elem().Underlying().(*types.Basic)
	if !ok {
		return false
	}
	// Exclude byte (uint8) — the VM has native KindBytes support.
	// Exclude int64/uint64 — stored natively, assertion works directly.
	if bt.Kind() == types.Byte || bt.Kind() == types.Int64 || bt.Kind() == types.Uint64 {
		return false
	}
	return bt.Info()&types.IsInteger != 0
}

// resolveSliceType unwraps a type (including named types and aliases) to find
// the underlying *types.Slice, if any. Returns (slice, true) or (nil, false).
func resolveSliceType(t types.Type) (*types.Slice, bool) {
	st, ok := t.Underlying().(*types.Slice)
	return st, ok
}

// intSliceWriteback records info needed to generate post-call writeback code
// for a parameter that was converted from []int64 to a narrower integer slice.
type intSliceWriteback struct {
	argName  string // e.g. "a0"
	backName string // e.g. "_back0"
	elemName string // e.g. "int", "math_big.Word"
}

// emitIntSliceExtraction writes multi-line extraction code for an integer slice
// parameter that needs conversion from the VM's native []int64 representation.
// It returns an intSliceWriteback struct for generating post-call writeback.
func emitIntSliceExtraction(b *strings.Builder, st *types.Slice, valExpr string, argName string, backName string, pkgRef string) *intSliceWriteback {
	elemName := resolveTypeName(st.Elem(), pkgRef)

	// Declare backing store reference and target slice
	fmt.Fprintf(b, "\tvar %s []int64\n", backName)
	fmt.Fprintf(b, "\tvar %s []%s\n", argName, elemName)
	fmt.Fprintf(b, "\tif _s, _ok := %s.IntSlice(); _ok {\n", valExpr)
	fmt.Fprintf(b, "\t\t%s = _s\n", backName)
	fmt.Fprintf(b, "\t\t%s = make([]%s, len(_s))\n", argName, elemName)
	fmt.Fprintf(b, "\t\tfor _i, _v := range _s { %s[_i] = %s(_v) }\n", argName, elemName)
	b.WriteString("\t} else {\n")
	fmt.Fprintf(b, "\t\t%s = %s.Interface().([]%s)\n", argName, valExpr, elemName)
	b.WriteString("\t}\n")

	return &intSliceWriteback{
		argName:  argName,
		backName: backName,
		elemName: elemName,
	}
}

// emitWritebacks writes post-call writeback code for all integer slice parameters
// that were converted from []int64.
func emitWritebacks(b *strings.Builder, writebacks []*intSliceWriteback) {
	for _, wb := range writebacks {
		fmt.Fprintf(b, "\tif %s != nil { for _i, _v := range %s { %s[_i] = int64(_v) } }\n",
			wb.backName, wb.argName, wb.backName)
	}
}
