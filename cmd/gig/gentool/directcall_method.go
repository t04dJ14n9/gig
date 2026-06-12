package gentool

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

// methodDirectCallInfo holds generated method DirectCall code and registration info.
type methodDirectCallInfo struct {
	TypeName   string
	MethodName string
	FuncName   string
	Code       string
}

type methodDirectCallShape struct {
	params     *types.Tuple
	results    *types.Tuple
	recv       *types.Var
	fixedCount int
	isVariadic bool
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
	shape, ok := analyzeMethodDirectCallSignature(sig)
	if !ok {
		return ""
	}

	recvExpr := methodReceiverExpr(shape.recv.Type(), pkgRef, typeName)

	funcName := fmt.Sprintf("direct_method_%s_%s_%s", pkgRef, typeName, methodName)
	b := &strings.Builder{}
	fmt.Fprintf(b, "func %s(args []value.Value) value.Value {\n", funcName)
	fmt.Fprintf(b, "\trecv := %s\n", recvExpr)

	argExprs, writebacks, ok := emitMethodFixedArgs(b, shape.params, shape.fixedCount, pkgRef)
	if !ok {
		return ""
	}
	argExprs, ok = emitMethodVariadicArg(b, shape, argExprs, pkgRef)
	if !ok {
		return ""
	}

	callExpr := fmt.Sprintf("recv.%s(%s)", methodName, strings.Join(argExprs, ", "))

	if emitCallResults(b, shape.results, callExpr, writebacks) == "" {
		return ""
	}

	b.WriteString("}\n")
	return b.String()
}

func analyzeMethodDirectCallSignature(sig *types.Signature) (methodDirectCallShape, bool) {
	shape := methodDirectCallShape{
		params:     sig.Params(),
		results:    sig.Results(),
		recv:       sig.Recv(),
		isVariadic: sig.Variadic(),
	}
	if shape.recv == nil {
		return methodDirectCallShape{}, false
	}

	shape.fixedCount = shape.params.Len()
	if shape.isVariadic {
		shape.fixedCount--
	}
	if !fixedDirectCallParamsSupported(shape.params, shape.fixedCount) {
		return methodDirectCallShape{}, false
	}
	if shape.isVariadic && !variadicDirectCallParamSupported(shape.params) {
		return methodDirectCallShape{}, false
	}
	if shape.results.Len() > 6 {
		return methodDirectCallShape{}, false
	}
	return shape, true
}

func fixedDirectCallParamsSupported(params *types.Tuple, fixedCount int) bool {
	for i := 0; i < fixedCount; i++ {
		if !canWrapParam(params.At(i).Type()) {
			return false
		}
	}
	return true
}

func variadicDirectCallParamSupported(params *types.Tuple) bool {
	sliceType := params.At(params.Len() - 1).Type().(*types.Slice)
	elemType := sliceType.Elem()
	return canWrapParam(elemType) || isEmptyInterface(elemType)
}

func methodReceiverExpr(recvType types.Type, pkgRef, typeName string) string {
	if _, ok := recvType.(*types.Pointer); ok {
		return fmt.Sprintf("args[0].Interface().(*%s.%s)", pkgRef, typeName)
	}
	return valueMethodReceiverExpr(recvType, pkgRef, typeName)
}

func valueMethodReceiverExpr(recvType types.Type, pkgRef, typeName string) string {
	// Named-basic receivers are stored by the VM as their underlying basic kind,
	// so generated wrappers rebuild the named value before dispatching.
	if bt, ok := recvType.Underlying().(*types.Basic); ok {
		if basicExpr := extractBasic(bt, "args[0]"); basicExpr != "" {
			return fmt.Sprintf("%s.%s(%s)", pkgRef, typeName, basicExpr)
		}
	}
	return fmt.Sprintf("args[0].Interface().(%s.%s)", pkgRef, typeName)
}

func emitMethodFixedArgs(
	b *strings.Builder,
	params *types.Tuple,
	fixedCount int,
	pkgRef string,
) ([]string, []*intSliceWriteback, bool) {
	var argExprs []string
	var writebacks []*intSliceWriteback
	for i := 0; i < fixedCount; i++ {
		paramType := params.At(i).Type()
		argName := fmt.Sprintf("a%d", i)
		valExpr := fmt.Sprintf("args[%d]", i+1)

		if st, ok := resolveSliceType(paramType); ok && needsIntSliceConversion(st) {
			backName := fmt.Sprintf("_back%d", i)
			wb := emitIntSliceExtraction(b, st, valExpr, argName, backName, pkgRef)
			writebacks = append(writebacks, wb)
			argExprs = append(argExprs, argName)
			continue
		}

		expr := extractArg(paramType, valExpr, pkgRef)
		if expr == "" {
			return nil, nil, false
		}
		fmt.Fprintf(b, "\t%s := %s\n", argName, expr)
		argExprs = append(argExprs, argName)
	}
	return argExprs, writebacks, true
}

func emitMethodVariadicArg(
	b *strings.Builder,
	shape methodDirectCallShape,
	argExprs []string,
	pkgRef string,
) ([]string, bool) {
	if !shape.isVariadic {
		return argExprs, true
	}

	sliceType := shape.params.At(shape.params.Len() - 1).Type().(*types.Slice)
	varArgExpr := emitVariadicArgs(b, variadicConfig{
		elemType: sliceType.Elem(),
		pkgRef:   pkgRef,
		fixedCnt: shape.fixedCount,
		argsOff:  1,
		isMethod: true,
	})
	if varArgExpr == "" {
		return nil, false
	}
	return append(argExprs, varArgExpr), true
}
