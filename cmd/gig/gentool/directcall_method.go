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
		// Named-basic receivers are stored by the VM as their underlying basic kind.
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
			argsOff:  1,
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
