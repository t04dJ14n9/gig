// Package optimize implements a 4-pass bytecode optimization pipeline.
//
// The passes run in order after SSA→bytecode compilation:
//  1. Peephole — pattern-based superinstruction fusion (17 rules)
//  2. Slice fusion — OpIntSliceGet/Set/SetConst for native []int64
//  3. Int specialization — generic ops → OpInt* variants for int-typed locals
//  4. Move fusion — OpIntLocal+OpIntSetLocal → OpIntMoveLocal
//
// See docs/optimization-report.md for detailed performance analysis.
package optimize

// Optimize applies all optimization passes to compiled bytecode in the correct order.
// localIsInt/constIsInt/localIsIntSlice flag which slots hold int-typed values.
// Returns the optimized code and whether int-specialized opcodes were emitted.
func Optimize(code []byte, localIsInt, constIsInt, localIsIntSlice []bool) ([]byte, bool) {
	code = Peephole(code)
	code = FuseSliceOps(code, localIsInt, localIsIntSlice)
	code, hasInt := IntSpecialize(code, localIsInt, constIsInt)
	code = FuseIntMoves(code)
	return code, hasInt
}
