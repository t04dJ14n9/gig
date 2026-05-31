package compiler

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// binOpMap maps Go token operators to bytecode opcodes.
var binOpMap = map[token.Token]bytecode.OpCode{
	token.ADD:     bytecode.OpAdd,
	token.SUB:     bytecode.OpSub,
	token.MUL:     bytecode.OpMul,
	token.QUO:     bytecode.OpDiv,
	token.REM:     bytecode.OpMod,
	token.AND:     bytecode.OpAnd,
	token.OR:      bytecode.OpOr,
	token.XOR:     bytecode.OpXor,
	token.AND_NOT: bytecode.OpAndNot,
	token.SHL:     bytecode.OpLsh,
	token.SHR:     bytecode.OpRsh,
	token.EQL:     bytecode.OpEqual,
	token.NEQ:     bytecode.OpNotEqual,
	token.LSS:     bytecode.OpLess,
	token.LEQ:     bytecode.OpLessEq,
	token.GTR:     bytecode.OpGreater,
	token.GEQ:     bytecode.OpGreaterEq,
	token.LAND:    bytecode.OpAnd,
	token.LOR:     bytecode.OpOr,
}

// compileBinOp compiles a binary operation.
func (c *compiler) compileBinOp(i *ssa.BinOp) {
	c.compileBinaryOpWithSetLocal(i.X, i.Y, i, binOpMap[i.Op])
}

// compileUnOp compiles a UnOp instruction.
func (c *compiler) compileUnOp(i *ssa.UnOp) {
	c.compileValue(i.X)

	var op bytecode.OpCode
	switch i.Op { //nolint:exhaustive
	case token.ADD:
		resultIdx := c.symbolTable.AllocLocal(i)
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	case token.SUB:
		op = bytecode.OpNeg
	case token.NOT:
		op = bytecode.OpNot
	case token.XOR:
		// Unary ^ (bitwise NOT) is not a single-stack op — it requires
		// pushing an all-ones constant then XORing. e.g. ^x = x ^ allOnes.
		// We must NOT emit OpXor directly here; that would pop 2 values
		// when only 1 (the operand) is on the stack, causing a panic.
		allOnes := allOnesConstant(i.X.Type())
		c.emit(bytecode.OpConst, uint16(c.addConstant(allOnes)))
		c.emit(bytecode.OpXor)
		resultIdx := c.symbolTable.AllocLocal(i)
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	case token.ARROW:
		if i.CommaOk {
			op = bytecode.OpRecvOk
		} else {
			op = bytecode.OpRecv
		}
	case token.MUL:
		op = bytecode.OpDeref
	}

	c.emit(op)

	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// allOnesConstant returns the "all ones" value for a given numeric type,
// used to implement unary ^ (bitwise NOT) as x ^ allOnes.
// Returns an any so it can be passed to addConstant, which calls FromInterface
// to create the correct value.Kind (KindUint for unsigned, KindInt for signed).
func allOnesConstant(t types.Type) any {
	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return int64(-1)
	}
	if mask, ok := allOnesByBasicKind[basic.Kind()]; ok {
		return mask
	}
	return int64(-1) // fallback
}

// allOnesByBasicKind is the unary-^ mask table. Values intentionally keep the
// original Go width so addConstant records the same signed/unsigned Value kind.
var allOnesByBasicKind = map[types.BasicKind]any{
	types.Uint8:   uint8(0xFF),
	types.Uint16:  uint16(0xFFFF),
	types.Uint32:  uint32(0xFFFFFFFF),
	types.Uint:    uint(^uint(0)),
	types.Uint64:  uint64(^uint64(0)),
	types.Uintptr: uintptr(^uintptr(0)),
	types.Int8:    int8(-1),
	types.Int16:   int16(-1),
	types.Int32:   int32(-1),
	types.Int:     int(^int(0)),
	types.Int64:   int64(-1),
}
