package peephole

import "git.woa.com/youngjin/gig/bytecode"

// arithSetPattern fuses arithOp SETLOCAL(A) (4 bytes → 3 bytes).
type arithSetPattern struct {
	arith bytecode.OpCode
	fused bytecode.OpCode
}

func init() {
	Register(
		arithSetPattern{bytecode.OpAdd, bytecode.OpAddSetLocal},
		arithSetPattern{bytecode.OpSub, bytecode.OpSubSetLocal},
	)
}

func (p arithSetPattern) Match(code []byte, i int) (int, []byte, bool) {
	const size = 4
	if bytecode.OpCode(code[i]) != p.arith || i+size > len(code) {
		return 0, nil, false
	}
	if bytecode.OpCode(code[i+1]) != bytecode.OpSetLocal {
		return 0, nil, false
	}
	a := ReadU16(code, i+2)
	return size, Make1Op(p.fused, a), true
}
