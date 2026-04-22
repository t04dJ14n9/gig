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

// --- Shared helpers for DirectCall generation ---

// resolveNamedBasicCast checks if a type is a named or alias type wrapping a basic
// type and returns the cast expression for variadic element conversion.
// Returns e.g. "color.Attribute(args[i].Int())" or "".
func resolveNamedBasicCast(elemType types.Type, pkgRef string) string {
	if named, ok := elemType.(*types.Named); ok {
		if bt, ok := named.Underlying().(*types.Basic); ok {
			basicExpr := extractBasic(bt, "args[i]")
			qualName := resolveQualifiedName(named, pkgRef)
			if basicExpr != "" && qualName != "" {
				return fmt.Sprintf("%s(%s)", qualName, basicExpr)
			}
		}
	}
	if alias, ok := elemType.(*types.Alias); ok {
		if bt, ok := alias.Underlying().(*types.Basic); ok {
			basicExpr := extractBasic(bt, "args[i]")
			obj := alias.Obj()
			pkg := obj.Pkg()
			if basicExpr != "" && pkg != nil {
				return fmt.Sprintf("%s.%s(%s)", sanitizePkgName(pkg.Path()), obj.Name(), basicExpr)
			}
		}
	}
	return ""
}

// variadicConfig controls variadic argument code generation.
type variadicConfig struct {
	elemType  types.Type
	pkgRef    string
	fixedCnt  int       // number of fixed params (before variadic)
	argsOff   int       // args index offset (0 for functions, 1 for methods after receiver)
	isMethod  bool      // true for method DirectCalls (different variadic pattern)
	useFmtWrap bool     // true for fmt.* functions (FmtWrap for Stringer dispatch)
}

// emitVariadicArgs generates variadic argument extraction code into the builder.
// Returns the arg expression to append (e.g. "varArgs...") or "" on error.
func emitVariadicArgs(b *strings.Builder, cfg variadicConfig) string {
	elemType := cfg.elemType
	pkgRef := cfg.pkgRef
	fixedCnt := cfg.fixedCnt
	argsOff := cfg.argsOff

	// Named types wrapping interface{} (e.g., type DecodeHookFunc interface{})
	// must use the named type for the slice, not []interface{}, because Go
	// does not allow implicit conversion between them.
	if isEmptyInterface(elemType) && !isNamedOrAlias(elemType) {
		if cfg.isMethod {
			// Method: compiler passes variadic args as a single slice at args[fixedCnt+1]
			sliceIdx := fixedCnt + argsOff
			fmt.Fprintf(b, "\tvar varArgs []interface{}\n")
			fmt.Fprintf(b, "\tif len(args) > %d {\n", sliceIdx)
			fmt.Fprintf(b, "\t\tif sl, ok := args[%d].Interface().([]interface{}); ok {\n", sliceIdx)
			fmt.Fprintf(b, "\t\t\tvarArgs = sl\n")
			fmt.Fprintf(b, "\t\t} else if sl, ok := args[%d].Interface().([]string); ok {\n", sliceIdx)
			fmt.Fprintf(b, "\t\t\tvarArgs = make([]interface{}, len(sl))\n")
			fmt.Fprintf(b, "\t\t\tfor i, s := range sl { varArgs[i] = s }\n")
			fmt.Fprintf(b, "\t\t}\n")
			fmt.Fprintf(b, "\t}\n")
		} else {
			// Function: args are spread across args[fixedCnt:]
			wrapExpr := "args[i].Interface()"
			if cfg.useFmtWrap {
				wrapExpr = "value.FmtWrap(args[i])"
			}
			fmt.Fprintf(b, "\tvarArgs := make([]interface{}, len(args)-%d)\n", fixedCnt)
			fmt.Fprintf(b, "\tfor i := %d; i < len(args); i++ {\n", fixedCnt)
			fmt.Fprintf(b, "\t\tvarArgs[i-%d] = %s\n", fixedCnt, wrapExpr)
			b.WriteString("\t}\n")
		}
		return "varArgs..."
	}

	elemTypeStr := resolveTypeName(elemType, pkgRef)
	if elemTypeStr == "" {
		return ""
	}

	namedBasicCast := resolveNamedBasicCast(elemType, pkgRef)

	if cfg.isMethod {
		// Method: compiler passes variadic args as a single slice at args[fixedCnt+1]
		sliceIdx := fixedCnt + argsOff
		fmt.Fprintf(b, "\tvar varArgs []%s\n", elemTypeStr)
		fmt.Fprintf(b, "\tif len(args) > %d {\n", sliceIdx)
		fmt.Fprintf(b, "\t\tif sl, ok := args[%d].Interface().([]%s); ok {\n", sliceIdx, elemTypeStr)
		fmt.Fprintf(b, "\t\t\tvarArgs = sl\n")
		fmt.Fprintf(b, "\t\t} else {\n")
		emitVariadicElementLoop(b, elemTypeStr, namedBasicCast, sliceIdx)
		b.WriteString("\t\t}\n")
		fmt.Fprintf(b, "\t}\n")
	} else {
		// Function: args are spread across args[fixedCnt:]
		if namedBasicCast != "" {
			fmt.Fprintf(b, "\tvarArgs := make([]%s, 0, len(args)-%d)\n", elemTypeStr, fixedCnt)
			fmt.Fprintf(b, "\tfor i := %d; i < len(args); i++ {\n", fixedCnt)
			fmt.Fprintf(b, "\t\tif args[i].IsValid() {\n")
			fmt.Fprintf(b, "\t\t\tvarArgs = append(varArgs, %s)\n", namedBasicCast)
			b.WriteString("\t\t}\n")
			b.WriteString("\t}\n")
		} else {
			fmt.Fprintf(b, "\tvarArgs := make([]%s, 0, len(args)-%d)\n", elemTypeStr, fixedCnt)
			fmt.Fprintf(b, "\tfor i := %d; i < len(args); i++ {\n", fixedCnt)
			fmt.Fprintf(b, "\t\tif v := args[i].Interface(); v != nil {\n")
			fmt.Fprintf(b, "\t\t\tvarArgs = append(varArgs, v.(%s))\n", elemTypeStr)
			b.WriteString("\t\t}\n")
			b.WriteString("\t}\n")
		}
	}

	return "varArgs..."
}

// emitVariadicElementLoop emits the element conversion loop for method variadic args
// when the slice couldn't be directly asserted.
func emitVariadicElementLoop(b *strings.Builder, elemTypeStr string, namedBasicCast string, startIdx int) {
	if namedBasicCast != "" {
		fmt.Fprintf(b, "\t\t\tfor i := %d; i < len(args); i++ {\n", startIdx)
		fmt.Fprintf(b, "\t\t\t\tif args[i].IsValid() {\n")
		fmt.Fprintf(b, "\t\t\t\t\tvarArgs = append(varArgs, %s)\n", namedBasicCast)
		b.WriteString("\t\t\t\t}\n")
		b.WriteString("\t\t\t}\n")
	} else {
		fmt.Fprintf(b, "\t\t\tfor i := %d; i < len(args); i++ {\n", startIdx)
		fmt.Fprintf(b, "\t\t\t\tif v := args[i].Interface(); v != nil {\n")
		fmt.Fprintf(b, "\t\t\t\t\tvarArgs = append(varArgs, v.(%s))\n", elemTypeStr)
		b.WriteString("\t\t\t\t}\n")
		b.WriteString("\t\t\t}\n")
	}
}

// emitCallResults generates the result wrapping code for a DirectCall function.
// Returns "" if any return type is unsupported.
func emitCallResults(b *strings.Builder, results *types.Tuple, callExpr string, writebacks []*intSliceWriteback) string {
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
		// Multi-return (2-6 results): r0, r1, ... := call(...)
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
	return "ok" // non-empty string indicates success
}

// --- Function DirectCall generation ---

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

	b := &strings.Builder{}
	fmt.Fprintf(b, "func direct_%s_%s(args []value.Value) value.Value {\n", pkgRef, fi.Name)

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
		fmt.Fprintf(b, "\t%s := %s\n", argName, expr)
		argExprs = append(argExprs, argName)
	}

	if isVariadic {
		sliceType := params.At(params.Len() - 1).Type().(*types.Slice)
		elemType := sliceType.Elem()
		varArgExpr := emitVariadicArgs(b, variadicConfig{
			elemType:   elemType,
			pkgRef:     pkgRef,
			fixedCnt:   fixedCount,
			argsOff:    0,
			isMethod:   false,
			useFmtWrap: isFmtPackage(pkgRef),
		})
		if varArgExpr == "" {
			return ""
		}
		argExprs = append(argExprs, varArgExpr)
	}

	// For Sprintf-like functions (string, ...interface{}) string, use
	// value.SprintfExtern to correctly handle %T with interpreter structs.
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

	if emitCallResults(b, results, callExpr, writebacks) == "" {
		return ""
	}

	b.WriteString("}\n")
	return b.String()
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
	fmt.Fprintf(b, "func %s(args []value.Value) value.Value {\n", funcName)
	fmt.Fprintf(b, "\trecv := %s\n", recvExpr)

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
		fmt.Fprintf(b, "\t%s := %s\n", argName, expr)
		argExprs = append(argExprs, argName)
	}

	if isVariadic {
		sliceType := params.At(params.Len() - 1).Type().(*types.Slice)
		elemType := sliceType.Elem()
		varArgExpr := emitVariadicArgs(b, variadicConfig{
			elemType: elemType,
			pkgRef:   pkgRef,
			fixedCnt: fixedCount,
			argsOff:  1, // offset by 1 for receiver
			isMethod: true,
		})
		if varArgExpr == "" {
			return ""
		}
		argExprs = append(argExprs, varArgExpr)
	}

	callExpr := fmt.Sprintf("recv.%s(%s)", methodName, strings.Join(argExprs, ", "))

	if emitCallResults(b, results, callExpr, writebacks) == "" {
		return ""
	}

	b.WriteString("}\n")
	return b.String()
}
