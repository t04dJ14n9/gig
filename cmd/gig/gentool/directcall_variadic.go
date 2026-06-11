package gentool

import (
	"fmt"
	"go/types"
	"strings"
)

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
	elemType   types.Type
	pkgRef     string
	fixedCnt   int  // number of fixed params before the variadic parameter
	argsOff    int  // args index offset; 0 for functions, 1 for methods after receiver
	isMethod   bool // method DirectCalls receive variadic args as a single slice
	useFmtWrap bool // fmt.* functions need FmtWrap for Stringer dispatch
}

// emitVariadicArgs generates variadic argument extraction code into the builder.
// Returns the argument expression to append (e.g. "varArgs...") or "" on error.
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
			// Method: compiler passes variadic args as a single slice at args[fixedCnt+1].
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
			// Function: args are spread across args[fixedCnt:].
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
		// Method: compiler passes variadic args as a single slice at args[fixedCnt+1].
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
		// Function: args are spread across args[fixedCnt:].
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

// emitVariadicElementLoop emits element conversion for method variadic args when
// the slice could not be directly asserted.
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
