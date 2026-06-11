package gentool

import (
	"fmt"
	"go/types"
	"strings"
)

// needsIntSliceConversion returns true if a slice type's elements are integers
// that are not int64/uint64. The VM stores these as []int64, so direct calls
// need conversion plus post-call writeback for in-place mutation.
func needsIntSliceConversion(st *types.Slice) bool {
	bt, ok := st.Elem().Underlying().(*types.Basic)
	if !ok {
		return false
	}
	if bt.Kind() == types.Byte || bt.Kind() == types.Int64 || bt.Kind() == types.Uint64 {
		return false
	}
	return bt.Info()&types.IsInteger != 0
}

// resolveSliceType unwraps a type to find the underlying *types.Slice, if any.
func resolveSliceType(t types.Type) (*types.Slice, bool) {
	st, ok := t.Underlying().(*types.Slice)
	return st, ok
}

// intSliceWriteback records info needed to generate post-call writeback code
// for a parameter converted from []int64 to a narrower integer slice.
type intSliceWriteback struct {
	argName  string
	backName string
	elemName string
}

// emitIntSliceExtraction writes extraction code for an integer slice parameter.
func emitIntSliceExtraction(b *strings.Builder, st *types.Slice, valExpr string, argName string, backName string, pkgRef string) *intSliceWriteback {
	elemName := resolveTypeName(st.Elem(), pkgRef)

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

// emitWritebacks writes post-call writeback code for integer slice parameters.
func emitWritebacks(b *strings.Builder, writebacks []*intSliceWriteback) {
	for _, wb := range writebacks {
		fmt.Fprintf(b, "\tif %s != nil { for _i, _v := range %s { %s[_i] = int64(_v) } }\n",
			wb.backName, wb.argName, wb.backName)
	}
}
