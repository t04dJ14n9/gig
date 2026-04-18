package gentool

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

// --- DirectCall code generation ---

// isSprintfLike checks if a function has the signature pattern:
// func(string, ...interface{}) string — i.e., it accepts a format string
// followed by variadic args and returns a single string. For such functions,
// the generated code uses value.SprintfExtern to handle %T correctly for
// interpreter-synthesized structs. This is a general check, not package-specific.
func isSprintfLike(fi *funcInfo) bool {
	sig := fi.Sig
	if !sig.Variadic() {
		return false
	}
	params := sig.Params()
	if params.Len() < 2 {
		return false
	}
	if bt, ok := params.At(0).Type().Underlying().(*types.Basic); !ok || bt.Kind() != types.String {
		return false
	}
	sliceType := params.At(params.Len() - 1).Type().(*types.Slice)
	if !isEmptyInterface(sliceType.Elem()) {
		return false
	}
	results := sig.Results()
	if results.Len() != 1 {
		return false
	}
	if bt, ok := results.At(0).Type().Underlying().(*types.Basic); !ok || bt.Kind() != types.String {
		return false
	}
	return true
}

// isFmtPackage returns true if the package reference refers to "fmt".
// Functions in the fmt package need FmtWrap for interface{} args so that
// gig structs with String() methods are properly formatted.
func isFmtPackage(pkgRef string) bool {
	return pkgRef == "fmt"
}

// customCallOverrides maps (pkgRef, funcName) to a custom call expression
// generator. When a function matches, the generator is called instead of
// the default call expression. This is used for functions that need special
// handling (e.g., errors.As needs GigErrorsAs for interpreter type matching).
var customCallOverrides = map[string]map[string]func(argExprs []string) string{
	"errors": {
		// errors.As needs GigErrorsAs because reflect.StructOf types can't
		// implement interfaces, so reflect.Type.AssignableTo won't match
		// interpreter-defined error types.
		"As": func(argExprs []string) string {
			return fmt.Sprintf("value.GigErrorsAs(%s, %s)", argExprs[0], argExprs[1])
		},
	},
}

func generateDirectCall(fi *funcInfo, pkgRef string) string {
	sig := fi.Sig
	params := sig.Params()
	results := sig.Results()

	if sig.Recv() != nil {
		return ""
	}

	isVariadic := sig.Variadic()
	fixedCount := params.Len()
	if isVariadic {
		fixedCount--
	}

	for i := 0; i < fixedCount; i++ {
		if !canWrapParam(params.At(i).Type()) {
			return ""
		}
	}
	if isVariadic {
		sliceType := params.At(params.Len() - 1).Type().(*types.Slice)
		elemType := sliceType.Elem()
		if !canWrapParam(elemType) && !isEmptyInterface(elemType) {
			return ""
		}
	}
	if results.Len() > 6 {
		return ""
	}

	// Determine if this function needs FmtWrap for interface{} args.
	// fmt.* functions check for Stringer/Formatter interfaces, so they need wrapping.
	useFmtWrap := isFmtPackage(pkgRef)

	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("func direct_%s_%s(args []value.Value) value.Value {\n", pkgRef, fi.Name))

	var argExprs []string
	var writebacks []*intSliceWriteback
	for i := 0; i < fixedCount; i++ {
		paramType := params.At(i).Type()
		argName := fmt.Sprintf("a%d", i)
		valExpr := fmt.Sprintf("args[%d]", i)

		// Check if this is an integer slice needing conversion + writeback
		if st, ok := resolveSliceType(paramType); ok && needsIntSliceConversion(st) {
			backName := fmt.Sprintf("_back%d", i)
			wb := emitIntSliceExtraction(b, st, valExpr, argName, backName, pkgRef)
			writebacks = append(writebacks, wb)
			argExprs = append(argExprs, argName)
			continue
		}

		expr := extractArg(paramType, valExpr, pkgRef)
		if expr == "" {
			return ""
		}
		b.WriteString(fmt.Sprintf("\t%s := %s\n", argName, expr))
		argExprs = append(argExprs, argName)
	}

	if isVariadic {
		sliceType := params.At(params.Len() - 1).Type().(*types.Slice)
		elemType := sliceType.Elem()

		if isEmptyInterface(elemType) {
			// For fmt.* functions, use FmtWrap to enable Stringer dispatch.
			// For everything else, use plain Interface() — encoding/sort/etc.
			// work natively on the raw struct.
			wrapExpr := "args[i].Interface()"
			if useFmtWrap {
				wrapExpr = "value.FmtWrap(args[i])"
			}
			b.WriteString(fmt.Sprintf("\tvarArgs := make([]interface{}, len(args)-%d)\n", fixedCount))
			b.WriteString(fmt.Sprintf("\tfor i := %d; i < len(args); i++ {\n", fixedCount))
			b.WriteString(fmt.Sprintf("\t\tvarArgs[i-%d] = %s\n", fixedCount, wrapExpr))
			b.WriteString("\t}\n")
			argExprs = append(argExprs, "varArgs...")
		} else {
			elemTypeStr := resolveTypeName(elemType, pkgRef)
			elemExtract := extractArg(elemType, "args[i]", pkgRef)
			if elemTypeStr == "" || elemExtract == "" {
				return ""
			}
			b.WriteString(fmt.Sprintf("\tvarArgs := make([]%s, len(args)-%d)\n", elemTypeStr, fixedCount))
			b.WriteString(fmt.Sprintf("\tfor i := %d; i < len(args); i++ {\n", fixedCount))
			b.WriteString(fmt.Sprintf("\t\tvarArgs[i-%d] = %s\n", fixedCount, elemExtract))
			b.WriteString("\t}\n")
			argExprs = append(argExprs, "varArgs...")
		}
	}

	// For Sprintf-like functions (string, ...interface{}) string, use
	// value.SprintfExtern to correctly handle %T with interpreter structs.
	// This is a general approach — applies to any package, not just fmt.
	var callExpr string
	if isSprintfLike(fi) {
		callExpr = fmt.Sprintf("value.SprintfExtern(%s)", strings.Join(argExprs, ", "))
	} else if override := customCallOverrides[pkgRef]; override != nil {
		if gen, ok := override[fi.Name]; ok {
			callExpr = gen(argExprs)
		} else {
			callExpr = fmt.Sprintf("%s.%s(%s)", pkgRef, fi.Name, strings.Join(argExprs, ", "))
		}
	} else {
		callExpr = fmt.Sprintf("%s.%s(%s)", pkgRef, fi.Name, strings.Join(argExprs, ", "))
	}

	switch results.Len() {
	case 0:
		fmt.Fprintf(b, "\t%s\n", callExpr)
		emitWritebacks(b, writebacks)
		b.WriteString("\treturn value.MakeNil()\n")
	case 1:
		if len(writebacks) > 0 {
			fmt.Fprintf(b, "\t_ret := %s\n", callExpr)
			emitWritebacks(b, writebacks)
			retExpr := wrapReturn(results.At(0).Type(), "_ret")
			if retExpr == "" {
				return ""
			}
			fmt.Fprintf(b, "\treturn %s\n", retExpr)
		} else {
			retExpr := wrapReturn(results.At(0).Type(), callExpr)
			if retExpr == "" {
				return ""
			}
			fmt.Fprintf(b, "\treturn %s\n", retExpr)
		}
	default:
		// Multi-return (2-4 results): r0, r1, ... := call(...)
		var retVars []string
		for i := 0; i < results.Len(); i++ {
			retVars = append(retVars, fmt.Sprintf("r%d", i))
		}
		fmt.Fprintf(b, "\t%s := %s\n", strings.Join(retVars, ", "), callExpr)
		emitWritebacks(b, writebacks)
		var wrapped []string
		for i := 0; i < results.Len(); i++ {
			w := wrapReturn(results.At(i).Type(), fmt.Sprintf("r%d", i))
			wrapped = append(wrapped, w)
		}
		fmt.Fprintf(b, "\treturn value.MakeValueSlice([]value.Value{%s})\n", strings.Join(wrapped, ", "))
	}

	b.WriteString("}\n")
	return b.String()
}

// --- Type info structs ---
func canWrapParam(t types.Type) bool {
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		pkg := obj.Pkg()

		if pkg == nil {
			return obj.Name() == errorTypeName
		}

		if pkg.Path() == currentPkgPath {
			return canWrapType(t.Underlying(), false)
		}

		// Cross-package named types: allow if we can extract via .Interface().(Type)
		return canWrapType(t.Underlying(), true)
	}

	// Handle type aliases (Go 1.23+), including builtin 'any'
	if alias, ok := t.(*types.Alias); ok {
		obj := alias.Obj()
		pkg := obj.Pkg()

		if pkg == nil {
			return canWrapType(t.Underlying(), false)
		}

		if pkg.Path() == currentPkgPath {
			return canWrapType(t.Underlying(), false)
		}

		return canWrapType(t.Underlying(), true)
	}

	return canWrapType(t.Underlying(), false)
}

// canWrapType checks whether a type can be wrapped in a DirectCall.
// If crossPkg is true, all representable types are allowed (extracted via .Interface()).
// If false, stricter checks apply (e.g., slices only with basic element types).
func canWrapType(t types.Type, crossPkg bool) bool {
	switch ut := t.(type) {
	case *types.Basic:
		return ut.Kind() != types.UnsafePointer && ut.Kind() != types.Invalid
	case *types.Slice:
		if crossPkg {
			return true
		}
		// Same-package: allow slices with basic element types only
		if _, ok := ut.Elem().Underlying().(*types.Basic); ok {
			return true
		}
		return false
	case *types.Interface:
		return true
	case *types.Pointer:
		if crossPkg {
			return true
		}
		if bt, ok := ut.Elem().(*types.Basic); ok {
			return bt.Kind() != types.UnsafePointer && bt.Kind() != types.Invalid
		}
		_, isNamed := ut.Elem().(*types.Named)
		return isNamed
	case *types.Struct:
		return true
	case *types.Map:
		return true
	case *types.Chan:
		return true
	case *types.Signature:
		return true
	case *types.Array:
		return true
	default:
		return false
	}
}

// --- Argument extraction ---

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
		underlying := t.Underlying()
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

	// Handle type aliases (Go 1.23+), including builtin 'any'
	if alias, ok := t.(*types.Alias); ok {
		underlying := t.Underlying()
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

	// For non-named types (e.g., *pkg.Type pointers), handle with pkgRef
	return extractUnderlyingWithPkgRef(t, valExpr, pkgRef)
}

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
		return fmt.Sprintf("%s.Complex()", valExpr)
	case info&types.IsString != 0:
		return fmt.Sprintf("%s.String()", valExpr)
	case kind == types.UnsafePointer:
		return ""
	default:
		return fmt.Sprintf("%s.Interface()", valExpr)
	}
}

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
	b.WriteString(fmt.Sprintf("\tvar %s []int64\n", backName))
	b.WriteString(fmt.Sprintf("\tvar %s []%s\n", argName, elemName))
	b.WriteString(fmt.Sprintf("\tif _s, _ok := %s.IntSlice(); _ok {\n", valExpr))
	b.WriteString(fmt.Sprintf("\t\t%s = _s\n", backName))
	b.WriteString(fmt.Sprintf("\t\t%s = make([]%s, len(_s))\n", argName, elemName))
	b.WriteString(fmt.Sprintf("\t\tfor _i, _v := range _s { %s[_i] = %s(_v) }\n", argName, elemName))
	b.WriteString("\t} else {\n")
	b.WriteString(fmt.Sprintf("\t\t%s = %s.Interface().([]%s)\n", argName, valExpr, elemName))
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
		b.WriteString(fmt.Sprintf("\tif %s != nil { for _i, _v := range %s { %s[_i] = int64(_v) } }\n",
			wb.backName, wb.argName, wb.backName))
	}
}

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

// --- Return value wrapping ---

func wrapReturn(t types.Type, goExpr string) string {
	// Unwrap named/alias types to check their underlying type
	underlying := t.Underlying()

	// Basic type (or named type with basic underlying): use typed Make* constructors
	if bt, ok := underlying.(*types.Basic); ok {
		// If it's a named type wrapping a basic, cast to the basic type first
		if _, isNamed := t.(*types.Named); isNamed {
			basicName := bt.Name()
			return wrapBasicReturn(bt, fmt.Sprintf("%s(%s)", basicName, goExpr))
		}
		if _, isAlias := t.(*types.Alias); isAlias {
			basicName := bt.Name()
			return wrapBasicReturn(bt, fmt.Sprintf("%s(%s)", basicName, goExpr))
		}
		return wrapBasicReturn(bt, goExpr)
	}

	// []byte: use MakeBytes for zero-reflection
	if st, ok := underlying.(*types.Slice); ok {
		if bt, ok := st.Elem().Underlying().(*types.Basic); ok && bt.Kind() == types.Byte {
			return fmt.Sprintf("value.MakeBytes([]byte(%s))", goExpr)
		}
	}

	// error interface: use FromInterface (handles nil correctly)
	if named, ok := t.(*types.Named); ok {
		if named.Obj().Pkg() == nil && named.Obj().Name() == errorTypeName {
			return fmt.Sprintf("value.FromInterface(%s)", goExpr)
		}
	}

	return fmt.Sprintf("value.FromInterface(%s)", goExpr)
}

func wrapBasicReturn(bt *types.Basic, goExpr string) string {
	info := bt.Info()
	kind := bt.Kind()

	switch {
	case kind == types.Bool:
		return fmt.Sprintf("value.MakeBool(%s)", goExpr)
	case info&types.IsInteger != 0 && info&types.IsUnsigned != 0:
		// Use sized Make* constructors to preserve the exact width across the
		// VM boundary.  MakeUint uses SizePtr (returns uint via Interface()),
		// while MakeUint64 uses Size64 (returns uint64 via Interface()).
		switch kind {
		case types.Uint64:
			return fmt.Sprintf("value.MakeUint64(%s)", goExpr)
		default:
			return fmt.Sprintf("value.MakeUint(uint64(%s))", goExpr)
		}
	case info&types.IsInteger != 0:
		// Use sized Make* constructors to preserve the exact width across the
		// VM boundary.  MakeInt uses SizePtr (returns int via Interface()),
		// while MakeInt64 uses Size64 (returns int64 via Interface()).
		switch kind {
		case types.Int64:
			return fmt.Sprintf("value.MakeInt64(%s)", goExpr)
		default:
			return fmt.Sprintf("value.MakeInt(int64(%s))", goExpr)
		}
	case info&types.IsFloat != 0:
		return fmt.Sprintf("value.MakeFloat(float64(%s))", goExpr)
	case info&types.IsComplex != 0:
		return fmt.Sprintf("value.FromInterface(%s)", goExpr)
	case info&types.IsString != 0:
		return fmt.Sprintf("value.MakeString(string(%s))", goExpr)
	default:
		return fmt.Sprintf("value.FromInterface(%s)", goExpr)
	}
}

// --- Method DirectCall generation ---

// methodDirectCallInfo holds generated method DirectCall code and registration info.
type methodDirectCallInfo struct {
	TypeName   string
	MethodName string
	FuncName   string
	Code       string
}

// generateMethodDirectCalls generates DirectCall wrappers for all eligible methods of a named type.
func generateMethodDirectCalls(named *types.Named, pkgRef string, typeName string) []*methodDirectCallInfo {
	var results []*methodDirectCallInfo

	methodSets := []*types.MethodSet{
		types.NewMethodSet(named),
		types.NewMethodSet(types.NewPointer(named)),
	}

	seen := make(map[string]bool)

	for _, mset := range methodSets {
		for i := 0; i < mset.Len(); i++ {
			sel := mset.At(i)
			if len(sel.Index()) != 1 {
				continue
			}
			fn := sel.Obj().(*types.Func)
			methodName := fn.Name()
			if seen[methodName] || !ast.IsExported(methodName) {
				continue
			}
			seen[methodName] = true

			sig := fn.Type().(*types.Signature)
			code := generateSingleMethodDirectCall(sig, pkgRef, typeName, methodName)
			if code != "" {
				funcName := fmt.Sprintf("direct_method_%s_%s_%s", pkgRef, typeName, methodName)
				results = append(results, &methodDirectCallInfo{
					TypeName:   typeName,
					MethodName: methodName,
					FuncName:   funcName,
					Code:       code,
				})
			}
		}
	}

	return results
}

// generateSingleMethodDirectCall generates a DirectCall wrapper for a single method.
// args[0] is the receiver, args[1:] are method arguments.
func generateSingleMethodDirectCall(sig *types.Signature, pkgRef string, typeName string, methodName string) string {
	params := sig.Params()
	results := sig.Results()
	recv := sig.Recv()
	if recv == nil {
		return ""
	}

	isVariadic := sig.Variadic()
	fixedCount := params.Len()
	if isVariadic {
		fixedCount--
	}
	for i := 0; i < fixedCount; i++ {
		if !canWrapParam(params.At(i).Type()) {
			return ""
		}
	}
	if isVariadic {
		sliceType := params.At(params.Len() - 1).Type().(*types.Slice)
		elemType := sliceType.Elem()
		if !canWrapParam(elemType) && !isEmptyInterface(elemType) {
			return ""
		}
	}
	if results.Len() > 6 {
		return ""
	}

	recvType := recv.Type()
	isPtr := false
	if _, ok := recvType.(*types.Pointer); ok {
		isPtr = true
	}

	var recvExpr string
	if isPtr {
		recvExpr = fmt.Sprintf("args[0].Interface().(*%s.%s)", pkgRef, typeName)
	} else {
		// Check if the receiver's named type has a basic underlying type.
		// If so, use a cast (e.g., time.Duration(args[0].Int())) instead of
		// .Interface().(time.Duration) which panics because the VM stores
		// named-basic types as their underlying basic kind.
		namedType := recvType
		if pt, ok := namedType.(*types.Pointer); ok {
			namedType = pt.Elem()
		}
		if bt, ok := namedType.Underlying().(*types.Basic); ok {
			basicExpr := extractBasic(bt, "args[0]")
			if basicExpr != "" {
				recvExpr = fmt.Sprintf("%s.%s(%s)", pkgRef, typeName, basicExpr)
			} else {
				recvExpr = fmt.Sprintf("args[0].Interface().(%s.%s)", pkgRef, typeName)
			}
		} else {
			recvExpr = fmt.Sprintf("args[0].Interface().(%s.%s)", pkgRef, typeName)
		}
	}

	funcName := fmt.Sprintf("direct_method_%s_%s_%s", pkgRef, typeName, methodName)
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("func %s(args []value.Value) value.Value {\n", funcName))
	b.WriteString(fmt.Sprintf("\trecv := %s\n", recvExpr))

	var argExprs []string
	var writebacks []*intSliceWriteback
	for i := 0; i < fixedCount; i++ {
		paramType := params.At(i).Type()
		argName := fmt.Sprintf("a%d", i)
		valExpr := fmt.Sprintf("args[%d]", i+1)

		// Check if this is an integer slice needing conversion + writeback
		if st, ok := resolveSliceType(paramType); ok && needsIntSliceConversion(st) {
			backName := fmt.Sprintf("_back%d", i)
			wb := emitIntSliceExtraction(b, st, valExpr, argName, backName, pkgRef)
			writebacks = append(writebacks, wb)
			argExprs = append(argExprs, argName)
			continue
		}

		expr := extractArg(paramType, valExpr, pkgRef)
		if expr == "" {
			return ""
		}
		b.WriteString(fmt.Sprintf("\t%s := %s\n", argName, expr))
		argExprs = append(argExprs, argName)
	}

	if isVariadic {
		sliceType := params.At(params.Len() - 1).Type().(*types.Slice)
		elemType := sliceType.Elem()

		if isEmptyInterface(elemType) {
			b.WriteString(fmt.Sprintf("\tvarArgs := make([]interface{}, len(args)-%d)\n", fixedCount+1))
			b.WriteString(fmt.Sprintf("\tfor i := %d; i < len(args); i++ {\n", fixedCount+1))
			b.WriteString(fmt.Sprintf("\t\tvarArgs[i-%d] = args[i].Interface()\n", fixedCount+1))
			b.WriteString("\t}\n")
			argExprs = append(argExprs, "varArgs...")
		} else {
			elemTypeStr := resolveTypeName(elemType, pkgRef)
			elemExtract := extractArg(elemType, "args[i]", pkgRef)
			if elemTypeStr == "" || elemExtract == "" {
				return ""
			}
			b.WriteString(fmt.Sprintf("\tvarArgs := make([]%s, len(args)-%d)\n", elemTypeStr, fixedCount+1))
			b.WriteString(fmt.Sprintf("\tfor i := %d; i < len(args); i++ {\n", fixedCount+1))
			b.WriteString(fmt.Sprintf("\t\tvarArgs[i-%d] = %s\n", fixedCount+1, elemExtract))
			b.WriteString("\t}\n")
			argExprs = append(argExprs, "varArgs...")
		}
	}

	callExpr := fmt.Sprintf("recv.%s(%s)", methodName, strings.Join(argExprs, ", "))

	switch results.Len() {
	case 0:
		fmt.Fprintf(b, "\t%s\n", callExpr)
		emitWritebacks(b, writebacks)
		b.WriteString("\treturn value.MakeNil()\n")
	case 1:
		if len(writebacks) > 0 {
			fmt.Fprintf(b, "\t_ret := %s\n", callExpr)
			emitWritebacks(b, writebacks)
			retExpr := wrapReturn(results.At(0).Type(), "_ret")
			if retExpr == "" {
				return ""
			}
			fmt.Fprintf(b, "\treturn %s\n", retExpr)
		} else {
			retExpr := wrapReturn(results.At(0).Type(), callExpr)
			if retExpr == "" {
				return ""
			}
			fmt.Fprintf(b, "\treturn %s\n", retExpr)
		}
	default:
		// Multi-return: r0, r1, ... := recv.Method(...)
		var retVars []string
		for i := 0; i < results.Len(); i++ {
			retVars = append(retVars, fmt.Sprintf("r%d", i))
		}
		fmt.Fprintf(b, "\t%s := %s\n", strings.Join(retVars, ", "), callExpr)
		emitWritebacks(b, writebacks)
		var wrapped []string
		for i := 0; i < results.Len(); i++ {
			w := wrapReturn(results.At(i).Type(), fmt.Sprintf("r%d", i))
			wrapped = append(wrapped, w)
		}
		fmt.Fprintf(b, "\treturn value.MakeValueSlice([]value.Value{%s})\n", strings.Join(wrapped, ", "))
	}

	b.WriteString("}\n")
	return b.String()
}
