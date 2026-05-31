package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

// The pattern we match is always 17 bytes:
//
//	LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) <op6> <op7>
//
// where <op6>/<op7> determines the fused instruction:
//
//	DEREF + SETLOCAL(v)      -> OpIntSliceGet(s, j, v)
//	LOCAL(val) + SETDEREF    -> OpIntSliceSet(s, j, val)
//	CONST(c) + SETDEREF      -> OpIntSliceSetConst(s, j, c)
const sliceFusionPatternLen = 17

type sliceFusionCandidate struct {
	start int
	s     uint16
	j     uint16
	ptr   uint16
}

// FuseSliceOps replaces common slice access patterns with fused superinstructions.
// It matches 17-byte sequences of [LOCAL, LOCAL, INDEXADDR, SETLOCAL, LOCAL, ...]
// and replaces them with 7-byte OpIntSlice{Get,Set,SetConst}.
func FuseSliceOps(code []byte, localIsInt, localIsIntSlice []bool) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}

		candidate, ok := sliceFusionCandidateAt(code, i)
		if !ok {
			i = instrEnd
			continue
		}

		fused := candidate.fusedBytes(code, localIsInt, localIsIntSlice)
		if fused == nil {
			i = instrEnd
			continue
		}

		rewrites = append(rewrites, rewrite{i, i + sliceFusionPatternLen, fused})
		i += sliceFusionPatternLen
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}

func sliceFusionCandidateAt(code []byte, i int) (sliceFusionCandidate, bool) {
	if !hasSliceFusionPrefix(code, i) {
		return sliceFusionCandidate{}, false
	}

	candidate := sliceFusionCandidate{
		start: i,
		s:     bytecode.ReadU16(code, i+1),
		j:     bytecode.ReadU16(code, i+4),
		ptr:   bytecode.ReadU16(code, i+8),
	}
	ptrGet := bytecode.ReadU16(code, i+11)
	if candidate.ptr != ptrGet {
		return sliceFusionCandidate{}, false
	}
	if localUsedOutside(code, candidate.ptr, i, i+sliceFusionPatternLen) {
		return sliceFusionCandidate{}, false
	}
	return candidate, true
}

func hasSliceFusionPrefix(code []byte, i int) bool {
	return i+sliceFusionPatternLen <= len(code) &&
		bytecode.OpCode(code[i]) == bytecode.OpLocal &&
		bytecode.OpCode(code[i+3]) == bytecode.OpLocal &&
		bytecode.OpCode(code[i+6]) == bytecode.OpIndexAddr &&
		bytecode.OpCode(code[i+7]) == bytecode.OpSetLocal &&
		bytecode.OpCode(code[i+10]) == bytecode.OpLocal
}

func (candidate sliceFusionCandidate) fusedBytes(code []byte, localIsInt, localIsIntSlice []bool) []byte {
	switch {
	case candidate.isGetPattern(code):
		return candidate.fuseGet(code, localIsInt, localIsIntSlice)
	case candidate.isSetPattern(code):
		return candidate.fuseSet(code, localIsInt, localIsIntSlice)
	case candidate.isSetConstPattern(code):
		return candidate.fuseSetConst(code, localIsInt, localIsIntSlice)
	default:
		return nil
	}
}

func (candidate sliceFusionCandidate) isGetPattern(code []byte) bool {
	return bytecode.OpCode(code[candidate.start+13]) == bytecode.OpDeref &&
		bytecode.OpCode(code[candidate.start+14]) == bytecode.OpSetLocal
}

func (candidate sliceFusionCandidate) isSetPattern(code []byte) bool {
	return bytecode.OpCode(code[candidate.start+13]) == bytecode.OpLocal &&
		bytecode.OpCode(code[candidate.start+16]) == bytecode.OpSetDeref
}

func (candidate sliceFusionCandidate) isSetConstPattern(code []byte) bool {
	return bytecode.OpCode(code[candidate.start+13]) == bytecode.OpConst &&
		bytecode.OpCode(code[candidate.start+16]) == bytecode.OpSetDeref
}

func (candidate sliceFusionCandidate) fuseGet(code []byte, localIsInt, localIsIntSlice []bool) []byte {
	v := bytecode.ReadU16(code, candidate.start+15)
	if !candidate.hasIntSliceBaseAndIndex(localIsInt, localIsIntSlice) || !safeIdx(int(v), localIsInt) {
		return nil
	}
	return makeSliceOp(bytecode.OpIntSliceGet, candidate.s, candidate.j, v)
}

func (candidate sliceFusionCandidate) fuseSet(code []byte, localIsInt, localIsIntSlice []bool) []byte {
	val := bytecode.ReadU16(code, candidate.start+14)
	if !candidate.hasIntSliceBaseAndIndex(localIsInt, localIsIntSlice) || !safeIdx(int(val), localIsInt) {
		return nil
	}
	return makeSliceOp(bytecode.OpIntSliceSet, candidate.s, candidate.j, val)
}

func (candidate sliceFusionCandidate) fuseSetConst(code []byte, localIsInt, localIsIntSlice []bool) []byte {
	constIdx := bytecode.ReadU16(code, candidate.start+14)
	if !candidate.hasIntSliceBaseAndIndex(localIsInt, localIsIntSlice) {
		return nil
	}
	return makeSliceOp(bytecode.OpIntSliceSetConst, candidate.s, candidate.j, constIdx)
}

func (candidate sliceFusionCandidate) hasIntSliceBaseAndIndex(localIsInt, localIsIntSlice []bool) bool {
	return safeIdx(int(candidate.s), localIsIntSlice) && safeIdx(int(candidate.j), localIsInt)
}

func makeSliceOp(op bytecode.OpCode, a, b, c uint16) []byte {
	out := make([]byte, 7)
	out[0] = byte(op)
	bytecode.WriteU16(out, 1, a)
	bytecode.WriteU16(out, 3, b)
	bytecode.WriteU16(out, 5, c)
	return out
}
