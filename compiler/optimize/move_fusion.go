package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

// FuseIntMoves replaces OpIntLocal(A) OpIntSetLocal(B) pairs with OpIntMoveLocal(A, B).
func FuseIntMoves(code []byte) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}

		if op == bytecode.OpIntLocal && i+6 <= len(code) &&
			bytecode.OpCode(code[i+3]) == bytecode.OpIntSetLocal {
			src := bytecode.ReadU16(code, i+1)
			dst := bytecode.ReadU16(code, i+4)
			newInstr := make([]byte, 5)
			newInstr[0] = byte(bytecode.OpIntMoveLocal)
			bytecode.WriteU16(newInstr, 1, src)
			bytecode.WriteU16(newInstr, 3, dst)
			rewrites = append(rewrites, rewrite{i, i + 6, newInstr})
			i += 6
			continue
		}

		i = instrEnd
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}
