package gentool

import (
	"fmt"
	"go/types"
	"strings"
)

// emitCallResults generates result wrapping code for a DirectCall function.
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
	return "ok"
}
