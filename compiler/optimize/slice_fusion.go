package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

// FuseSliceOps replaces common slice access patterns with fused superinstructions.
// It matches 17-byte sequences of [LOCAL, LOCAL, INDEXADDR, SETLOCAL, LOCAL, ...]
// and replaces them with 7-byte OpIntSlice{Get,Set,SetConst}.
func FuseSliceOps(code []byte, localIsInt, localIsIntSlice []bool) []byte {
	var rewrites []rewrite

	// The pattern we match is always 17 bytes:
	//   LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) <op6> <op7>
	// where <op6>/<op7> determines the fused instruction:
	//   DEREF + SETLOCAL(v)      -> OpIntSliceGet(s, j, v)
	//   LOCAL(val) + SETDEREF    -> OpIntSliceSet(s, j, val)
	//   CONST(c) + SETDEREF      -> OpIntSliceSetConst(s, j, c)
	const patternLen = 17

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}

		if op != bytecode.OpLocal || i+patternLen > len(code) {
			i = instrEnd
			continue
		}

		// Check the fixed prefix: LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr)
		if bytecode.OpCode(code[i+3]) != bytecode.OpLocal ||
			bytecode.OpCode(code[i+6]) != bytecode.OpIndexAddr ||
			bytecode.OpCode(code[i+7]) != bytecode.OpSetLocal ||
			bytecode.OpCode(code[i+10]) != bytecode.OpLocal {
			i = instrEnd
			continue
		}

		s := bytecode.ReadU16(code, i+1)
		j := bytecode.ReadU16(code, i+4)
		ptr := bytecode.ReadU16(code, i+8)
		ptrGet := bytecode.ReadU16(code, i+11)
		if ptr != ptrGet {
			i = instrEnd
			continue
		}

		// ptr must not be used elsewhere; it is a temporary.
		if localUsedOutside(code, ptr, i, i+patternLen) {
			i = instrEnd
			continue
		}

		op6 := bytecode.OpCode(code[i+13])
		op7 := bytecode.OpCode(code[i+14]) // for DEREF case, this is SETLOCAL opcode; for others, it is at i+16

		var fused []byte
		switch {
		case op6 == bytecode.OpDeref && op7 == bytecode.OpSetLocal:
			// s[j] read -> OpIntSliceGet(s, j, v)
			v := bytecode.ReadU16(code, i+15)
			if safeIdx(int(s), localIsIntSlice) && safeIdx(int(j), localIsInt) && safeIdx(int(v), localIsInt) {
				fused = makeSliceOp(bytecode.OpIntSliceGet, s, j, v)
			}
		case op6 == bytecode.OpLocal && bytecode.OpCode(code[i+16]) == bytecode.OpSetDeref:
			// s[j] = val -> OpIntSliceSet(s, j, val)
			val := bytecode.ReadU16(code, i+14)
			if safeIdx(int(s), localIsIntSlice) && safeIdx(int(j), localIsInt) && safeIdx(int(val), localIsInt) {
				fused = makeSliceOp(bytecode.OpIntSliceSet, s, j, val)
			}
		case op6 == bytecode.OpConst && bytecode.OpCode(code[i+16]) == bytecode.OpSetDeref:
			// s[j] = const -> OpIntSliceSetConst(s, j, constIdx)
			constIdx := bytecode.ReadU16(code, i+14)
			if safeIdx(int(s), localIsIntSlice) && safeIdx(int(j), localIsInt) {
				fused = makeSliceOp(bytecode.OpIntSliceSetConst, s, j, constIdx)
			}
		}

		if fused != nil {
			rewrites = append(rewrites, rewrite{i, i + patternLen, fused})
			i += patternLen
		} else {
			i = instrEnd
		}
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}

func makeSliceOp(op bytecode.OpCode, a, b, c uint16) []byte {
	out := make([]byte, 7)
	out[0] = byte(op)
	bytecode.WriteU16(out, 1, a)
	bytecode.WriteU16(out, 3, b)
	bytecode.WriteU16(out, 5, c)
	return out
}
