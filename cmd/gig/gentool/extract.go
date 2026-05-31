package gentool

import (
	"fmt"
	"go/types"
)

// --- Argument extraction ---

// extractArg generates the Go expression to extract a typed value from a value.Value
// argument for passing to a native Go function. Returns "" if the type is unsupported.
func extractArg(t types.Type, valExpr string, pkgRef string) string {
	if named, ok := t.(*types.Named); ok {
		if isBuiltinError(named.Obj()) {
			// Use value.ErrorValue to handle interpreter-defined types with Error()
			// methods that cannot satisfy error through reflect.StructOf.
			return fmt.Sprintf("value.ErrorValue(%s)", valExpr)
		}
	}

	if alias, ok := t.(*types.Alias); ok {
		obj := alias.Obj()
		if obj.Pkg() == nil {
			if obj.Name() == errorTypeName {
				return fmt.Sprintf("value.ErrorValue(%s)", valExpr)
			}
			return fmt.Sprintf("%s.Interface()", valExpr)
		}
	}

	if named, ok := t.(*types.Named); ok {
		return extractNamedArg(named, t.Underlying(), valExpr, pkgRef)
	}
	if alias, ok := t.(*types.Alias); ok {
		return extractAliasArg(alias, t.Underlying(), valExpr, pkgRef)
	}
	return extractUnderlyingWithPkgRef(t, valExpr, pkgRef)
}

func isBuiltinError(obj types.Object) bool {
	return obj.Pkg() == nil && obj.Name() == errorTypeName
}
