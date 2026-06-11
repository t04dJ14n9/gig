package gentool

import (
	"fmt"
	"go/types"
	"strings"
)

// --- Function DirectCall generation ---

type functionDirectCallShape struct {
	params     *types.Tuple
	results    *types.Tuple
	fixedCount int
	isVariadic bool
}

func generateDirectCall(fi *funcInfo, pkgRef string) string {
	shape, ok := analyzeFunctionDirectCallSignature(fi.Sig)
	if !ok {
		return ""
	}

	b := &strings.Builder{}
	fmt.Fprintf(b, "func direct_%s_%s(args []value.Value) value.Value {\n", pkgRef, fi.Name)

	argExprs, writebacks, ok := emitFunctionFixedArgs(b, shape.params, shape.fixedCount, pkgRef)
	if !ok {
		return ""
	}
	argExprs, ok = emitFunctionVariadicArg(b, shape, argExprs, pkgRef)
	if !ok {
		return ""
	}

	callExpr := functionCallExpr(fi, pkgRef, argExprs)
	if emitCallResults(b, shape.results, callExpr, writebacks) == "" {
		return ""
	}

	b.WriteString("}\n")
	return b.String()
}

func analyzeFunctionDirectCallSignature(sig *types.Signature) (functionDirectCallShape, bool) {
	shape := functionDirectCallShape{
		params:     sig.Params(),
		results:    sig.Results(),
		isVariadic: sig.Variadic(),
	}
	if sig.Recv() != nil {
		return functionDirectCallShape{}, false
	}

	shape.fixedCount = shape.params.Len()
	if shape.isVariadic {
		shape.fixedCount--
	}
	if !fixedDirectCallParamsSupported(shape.params, shape.fixedCount) {
		return functionDirectCallShape{}, false
	}
	if shape.isVariadic && !variadicDirectCallParamSupported(shape.params) {
		return functionDirectCallShape{}, false
	}
	if shape.results.Len() > 6 {
		return functionDirectCallShape{}, false
	}
	return shape, true
}

func emitFunctionFixedArgs(
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
		valExpr := fmt.Sprintf("args[%d]", i)

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

func emitFunctionVariadicArg(
	b *strings.Builder,
	shape functionDirectCallShape,
	argExprs []string,
	pkgRef string,
) ([]string, bool) {
	if !shape.isVariadic {
		return argExprs, true
	}

	sliceType := shape.params.At(shape.params.Len() - 1).Type().(*types.Slice)
	varArgExpr := emitVariadicArgs(b, variadicConfig{
		elemType:   sliceType.Elem(),
		pkgRef:     pkgRef,
		fixedCnt:   shape.fixedCount,
		argsOff:    0,
		isMethod:   false,
		useFmtWrap: isFmtPackage(pkgRef),
	})
	if varArgExpr == "" {
		return nil, false
	}
	return append(argExprs, varArgExpr), true
}

func functionCallExpr(fi *funcInfo, pkgRef string, argExprs []string) string {
	joinedArgs := strings.Join(argExprs, ", ")
	// Sprintf-like functions need the interpreter-aware formatter for %T.
	if isSprintfLike(fi) {
		return fmt.Sprintf("value.SprintfExtern(%s)", joinedArgs)
	}
	if override := customCallOverrides[pkgRef]; override != nil {
		if gen, ok := override[fi.Name]; ok {
			return gen(argExprs)
		}
	}
	return fmt.Sprintf("%s.%s(%s)", pkgRef, fi.Name, joinedArgs)
}
