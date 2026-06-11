package optimize

// FoldConstants reduces straight-line constant expressions before the regular
// peephole pipeline. It is deliberately local: state is reset at control-flow
// boundaries, and folds never cross jump targets.
func FoldConstants(code []byte, constants *[]any) []byte {
	for {
		propagated := propagateLocalConstants(code)
		folded, changed := foldConstantStackOps(propagated, constants)
		folded, branchChanged := foldConstantBranches(folded, *constants)
		folded, deadChanged := removeUnreachableAfterJumps(folded)
		code = folded
		if !changed && !branchChanged && !deadChanged {
			return code
		}
	}
}
