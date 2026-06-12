package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

type stackConst struct {
	idx uint16
	ok  bool
}

func propagateLocalConstants(code []byte) []byte {
	out := append([]byte(nil), code...)
	targets := jumpTargetSet(code)
	localConsts := make(map[uint16]uint16)
	var stack []stackConst

	for i := 0; i < len(out); {
		op := bytecode.OpCode(out[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(out) {
			break
		}
		if targets[i] {
			clearLocalState(localConsts, &stack)
		}
		handleConstantPropagation(out, i, op, localConsts, &stack)
		i = instrEnd
	}
	return out
}

func handleConstantPropagation(
	code []byte,
	i int,
	op bytecode.OpCode,
	localConsts map[uint16]uint16,
	stack *[]stackConst,
) {
	switch op {
	case bytecode.OpConst:
		pushConst(stack, bytecode.ReadU16(code, i+1))
	case bytecode.OpLocal:
		propagateLocalLoad(code, i, localConsts, stack)
	case bytecode.OpSetLocal:
		storeLocalConst(code, i, localConsts, stack)
	case bytecode.OpPop:
		popConst(stack)
	case bytecode.OpJump, bytecode.OpJumpTrue, bytecode.OpJumpFalse:
		clearLocalState(localConsts, stack)
	case bytecode.OpAddr, bytecode.OpSetDeref, bytecode.OpSetIndex, bytecode.OpSetGlobal, bytecode.OpSetFree:
		clearLocalState(localConsts, stack)
	default:
		*stack = nil
	}
}

func propagateLocalLoad(code []byte, i int, localConsts map[uint16]uint16, stack *[]stackConst) {
	local := bytecode.ReadU16(code, i+1)
	if constIdx, ok := localConsts[local]; ok {
		code[i] = byte(bytecode.OpConst)
		bytecode.WriteU16(code, i+1, constIdx)
		pushConst(stack, constIdx)
		return
	}
	*stack = append(*stack, stackConst{})
}

func storeLocalConst(code []byte, i int, localConsts map[uint16]uint16, stack *[]stackConst) {
	local := bytecode.ReadU16(code, i+1)
	top := popConst(stack)
	if top.ok {
		localConsts[local] = top.idx
		return
	}
	delete(localConsts, local)
}

func pushConst(stack *[]stackConst, idx uint16) {
	*stack = append(*stack, stackConst{idx: idx, ok: true})
}

func popConst(stack *[]stackConst) stackConst {
	if len(*stack) == 0 {
		return stackConst{}
	}
	last := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return last
}

func clearLocalState(localConsts map[uint16]uint16, stack *[]stackConst) {
	for local := range localConsts {
		delete(localConsts, local)
	}
	*stack = nil
}
