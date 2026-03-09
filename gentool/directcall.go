package gentool

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

// --- DirectCall code generation ---

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
	if results.Len() > 4 {
		return ""
	}

	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("func direct_%s_%s(args []value.Value) value.Value {\n", pkgRef, fi.Name))

	var argExprs []string
	for i := 0; i < fixedCount; i++ {
		paramType := params.At(i).Type()
		expr := extractArg(paramType, fmt.Sprintf("args[%d]", i), pkgRef)
		if expr == "" {
			return ""
		}
		argName := fmt.Sprintf("a%d", i)
		b.WriteString(fmt.Sprintf("\t%s := %s\n", argName, expr))
		argExprs = append(argExprs, argName)
	}

	if isVariadic {
		sliceType := params.At(params.Len() - 1).Type().(*types.Slice)
		elemType := sliceType.Elem()

		if isEmptyInterface(elemType) {
			b.WriteString(fmt.Sprintf("\tvarArgs := make([]interface{}, len(args)-%d)\n", fixedCount))
			b.WriteString(fmt.Sprintf("\tfor i := %d; i < len(args); i++ {\n", fixedCount))
			b.WriteString(fmt.Sprintf("\t\tvarArgs[i-%d] = args[i].Interface()\n", fixedCount))
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

	callExpr := fmt.Sprintf("%s.%s(%s)", pkgRef, fi.Name, strings.Join(argExprs, ", "))

	switch results.Len() {
	case 0:
		b.WriteString(fmt.Sprintf("\t%s\n", callExpr))
		b.WriteString("\treturn value.MakeNil()\n")
	case 1:
		retExpr := wrapReturn(results.At(0).Type(), callExpr)
		if retExpr == "" {
			return ""
		}
		b.WriteString(fmt.Sprintf("\treturn %s\n", retExpr))
	case 2:
		b.WriteString(fmt.Sprintf("\tr0, r1 := %s\n", callExpr))
		w0 := wrapReturn(results.At(0).Type(), "r0")
		w1 := wrapReturn(results.At(1).Type(), "r1")
		b.WriteString(fmt.Sprintf("\treturn value.MakeValueSlice([]value.Value{%s, %s})\n", w0, w1))
	case 3:
		b.WriteString(fmt.Sprintf("\tr0, r1, r2 := %s\n", callExpr))
		w0 := wrapReturn(results.At(0).Type(), "r0")
		w1 := wrapReturn(results.At(1).Type(), "r1")
		w2 := wrapReturn(results.At(2).Type(), "r2")
		b.WriteString(fmt.Sprintf("\treturn value.MakeValueSlice([]value.Value{%s, %s, %s})\n", w0, w1, w2))
	case 4:
		b.WriteString(fmt.Sprintf("\tr0, r1, r2, r3 := %s\n", callExpr))
		w0 := wrapReturn(results.At(0).Type(), "r0")
		w1 := wrapReturn(results.At(1).Type(), "r1")
		w2 := wrapReturn(results.At(2).Type(), "r2")
		w3 := wrapReturn(results.At(3).Type(), "r3")
		b.WriteString(fmt.Sprintf("\treturn value.MakeValueSlice([]value.Value{%s, %s, %s, %s})\n", w0, w1, w2, w3))
	default:
		return ""
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
			return canWrapUnderlying(t.Underlying())
		}

		// Cross-package named types: allow if we can extract via .Interface().(Type)
		// This covers structs, pointers, interfaces, and basic-underlying types from other packages
		return canWrapCrossPackage(t.Underlying())
	}

	// Handle type aliases (Go 1.23+), including builtin 'any'
	if alias, ok := t.(*types.Alias); ok {
		obj := alias.Obj()
		pkg := obj.Pkg()

		if pkg == nil {
			// Builtin alias (e.g. 'any' = interface{}, 'error' = interface{Error() string})
			// Allow if the underlying type is wrappable
			return canWrapUnderlying(t.Underlying())
		}

		if pkg.Path() == currentPkgPath {
			return canWrapUnderlying(t.Underlying())
		}

		return canWrapCrossPackage(t.Underlying())
	}

	return canWrapUnderlying(t.Underlying())
}

func canWrapUnderlying(t types.Type) bool {
	switch ut := t.(type) {
	case *types.Basic:
		return ut.Kind() != types.UnsafePointer && ut.Kind() != types.Invalid
	case *types.Slice:
		// Only allow []byte slices — other basic slices ([]int, []float64, etc.) are
		// stored as []int64/[]float64 in the VM which causes type assertion panics.
		// Non-basic element slices ([][]byte, []*T) are also excluded.
		if bt, ok := ut.Elem().Underlying().(*types.Basic); ok {
			return bt.Kind() == types.Byte || bt.Kind() == types.String
		}
		return false
	case *types.Interface:
		return true // support both empty and non-empty interfaces
	case *types.Pointer:
		// Only support pointers to named types or basic types (excluding unsafe.Pointer)
		if bt, ok := ut.Elem().(*types.Basic); ok {
			return bt.Kind() != types.UnsafePointer && bt.Kind() != types.Invalid
		}
		_, isNamed := ut.Elem().(*types.Named)
		return isNamed
	case *types.Struct:
		return true // support struct types
	case *types.Map:
		return true // support map types
	case *types.Chan:
		return true // support chan T via .Interface().(chan T)
	case *types.Signature:
		return true // support func(T) R via .Interface().(func(T) R)
	case *types.Array:
		return true // support [N]T via .Interface().([N]T)
	default:
		return false
	}
}

// canWrapCrossPackage checks if a cross-package type can be extracted via .Interface().(Type).
func canWrapCrossPackage(t types.Type) bool {
	switch t.(type) {
	case *types.Basic:
		return true
	case *types.Slice:
		return true
	case *types.Struct:
		return true
	case *types.Pointer:
		return true
	case *types.Interface:
		return true
	case *types.Map:
		return true
	case *types.Chan:
		return true // support chan T via .Interface().(chan T)
	case *types.Signature:
		return true // support func(T) R via .Interface().(func(T) R)
	case *types.Array:
		return true // support [N]T via .Interface().([N]T)
	default:
		return false
	}
}

// --- Argument extraction ---

func extractArg(t types.Type, valExpr string, pkgRef string) string {
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		if obj.Pkg() == nil && obj.Name() == errorTypeName {
			return fmt.Sprintf("%s.Interface().(error)", valExpr)
		}
	}

	// Handle type aliases (Go 1.23+), including builtin 'any'
	if alias, ok := t.(*types.Alias); ok {
		obj := alias.Obj()
		if obj.Pkg() == nil {
			// Builtin alias: 'error' or 'any' (= interface{})
			if obj.Name() == errorTypeName {
				return fmt.Sprintf("%s.Interface().(error)", valExpr)
			}
			// 'any' and other builtin aliases: extract as interface{}
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

		// Cross-package named type: use .Interface().(pkg.Type) for all cases
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

		// Cross-package alias: use .Interface()
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
		// (only []byte and []string get typed assertions; others use Interface())
		if _, ok := ut.Elem().Underlying().(*types.Basic); ok {
			return extractSlice(ut, valExpr)
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

func extractSlice(st *types.Slice, valExpr string) string {
	if bt, ok := st.Elem().Underlying().(*types.Basic); ok {
		switch bt.Kind() {
		case types.Byte:
			// Use native KindBytes accessor — zero reflection for []byte params
			// Falls back to Interface().([]byte) for KindReflect values
			return fmt.Sprintf("func() []byte { if b, ok := (%s).Bytes(); ok { return b }; return (%s).Interface().([]byte) }()", valExpr, valExpr)
		case types.String:
			// []string is stored as []string in the VM — safe to assert directly
			return fmt.Sprintf("%s.Interface().([]string)", valExpr)
		default:
			// All other basic slices ([]int, []int64, []float64, etc.) are stored
			// as their reflect type in the VM — use Interface() without assertion
			// to avoid type mismatch (e.g. VM stores []int as []int64 internally)
			return fmt.Sprintf("%s.Interface()", valExpr)
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
		return fmt.Sprintf("value.MakeUint(uint64(%s))", goExpr)
	case info&types.IsInteger != 0:
		return fmt.Sprintf("value.MakeInt(int64(%s))", goExpr)
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
	if results.Len() > 4 {
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
		recvExpr = fmt.Sprintf("args[0].Interface().(%s.%s)", pkgRef, typeName)
	}

	funcName := fmt.Sprintf("direct_method_%s_%s_%s", pkgRef, typeName, methodName)
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("func %s(args []value.Value) value.Value {\n", funcName))
	b.WriteString(fmt.Sprintf("\trecv := %s\n", recvExpr))

	var argExprs []string
	for i := 0; i < fixedCount; i++ {
		paramType := params.At(i).Type()
		expr := extractArg(paramType, fmt.Sprintf("args[%d]", i+1), pkgRef)
		if expr == "" {
			return ""
		}
		argName := fmt.Sprintf("a%d", i)
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
		b.WriteString(fmt.Sprintf("\t%s\n", callExpr))
		b.WriteString("\treturn value.MakeNil()\n")
	case 1:
		retExpr := wrapReturn(results.At(0).Type(), callExpr)
		if retExpr == "" {
			return ""
		}
		b.WriteString(fmt.Sprintf("\treturn %s\n", retExpr))
	case 2:
		b.WriteString(fmt.Sprintf("\tr0, r1 := %s\n", callExpr))
		w0 := wrapReturn(results.At(0).Type(), "r0")
		w1 := wrapReturn(results.At(1).Type(), "r1")
		b.WriteString(fmt.Sprintf("\treturn value.MakeValueSlice([]value.Value{%s, %s})\n", w0, w1))
	case 3:
		b.WriteString(fmt.Sprintf("\tr0, r1, r2 := %s\n", callExpr))
		w0 := wrapReturn(results.At(0).Type(), "r0")
		w1 := wrapReturn(results.At(1).Type(), "r1")
		w2 := wrapReturn(results.At(2).Type(), "r2")
		b.WriteString(fmt.Sprintf("\treturn value.MakeValueSlice([]value.Value{%s, %s, %s})\n", w0, w1, w2))
	default:
		return ""
	}

	b.WriteString("}\n")
	return b.String()
}
