package gentool

import (
	"fmt"
	"go/types"
)

// isSprintfLike checks if a function has the signature pattern:
// func(string, ...interface{}) string. Such functions use value.SprintfExtern
// so %T handles interpreter-synthesized structs correctly.
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
// fmt functions need FmtWrap for interface{} args so gig structs with String()
// methods are formatted through the interpreter method resolver.
func isFmtPackage(pkgRef string) bool {
	return pkgRef == "fmt"
}

// customCallOverrides maps (pkgRef, funcName) to a custom call expression
// generator. These are narrow shims for host APIs whose native reflection rules
// cannot observe interpreter-synthesized type identity.
var customCallOverrides = map[string]map[string]func(argExprs []string) string{
	"errors": {
		// errors.As needs GigErrorsAs because reflect.StructOf types cannot
		// implement interfaces, so reflect.Type.AssignableTo will not match
		// interpreter-defined error types.
		"As": func(argExprs []string) string {
			return fmt.Sprintf("value.GigErrorsAs(%s, %s)", argExprs[0], argExprs[1])
		},
	},
}
