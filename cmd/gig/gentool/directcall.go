package gentool

import (
	"fmt"
	"go/types"
	"strings"
)

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

	// Sprintf-like functions need the interpreter-aware formatter for %T.
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
