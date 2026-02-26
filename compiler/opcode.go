// Package compiler provides SSA-to-bytecode compilation.
package compiler

// OpCode represents a single bytecode instruction.
type OpCode byte

const (
	// Stack operations
	OpNop OpCode = iota // no operation
	OpPop               // pop top of stack
	OpDup               // duplicate top of stack

	// Constants and locals
	OpConst     // push constant from pool [const_idx:2]
	OpNil       // push nil
	OpTrue      // push true
	OpFalse     // push false
	OpLocal     // push local variable [local_idx:2]
	OpSetLocal  // pop and set local variable [local_idx:2]
	OpGlobal    // push global variable [global_idx:2]
	OpSetGlobal // pop and set global variable [global_idx:2]
	OpFree      // push free variable (closure) [free_idx:1]
	OpSetFree   // pop and set free variable [free_idx:1]

	// Arithmetic
	OpAdd // pop b, pop a, push a + b
	OpSub // pop b, pop a, push a - b
	OpMul // pop b, pop a, push a * b
	OpDiv // pop b, pop a, push a / b
	OpMod // pop b, pop a, push a % b
	OpNeg // pop a, push -a

	// Bitwise
	OpAnd    // pop b, pop a, push a & b
	OpOr     // pop b, pop a, push a | b
	OpXor    // pop b, pop a, push a ^ b
	OpAndNot // pop b, pop a, push a &^ b
	OpLsh    // pop b, pop a, push a << b
	OpRsh    // pop b, pop a, push a >> b

	// Comparison
	OpEqual     // pop b, pop a, push a == b
	OpNotEqual  // pop b, pop a, push a != b
	OpLess      // pop b, pop a, push a < b
	OpLessEq    // pop b, pop a, push a <= b
	OpGreater   // pop b, pop a, push a > b
	OpGreaterEq // pop b, pop a, push a >= b

	// Logical
	OpNot // pop a, push !a

	// Control flow
	OpJump      // unconditional jump [offset:2]
	OpJumpTrue  // jump if true [offset:2]
	OpJumpFalse // jump if false [offset:2]
	OpCall      // call function [num_args:1]
	OpReturn    // return from function
	OpReturnVal // pop and return value

	// Container operations
	OpMakeSlice  // make slice: pop cap, pop len, pop typeIdx from stack
	OpMakeMap    // make map: pop size, pop typeIdx from stack
	OpMakeChan   // make chan: pop size, pop typeIdx from stack
	OpMakeArray  // make array [type_idx:2]
	OpMakeStruct // make struct [type_idx:2]

	// Index operations
	OpIndex    // pop key, pop container, push container[key]
	OpSetIndex // pop val, pop key, pop container, container[key] = val
	OpSlice    // slice operation: pop max, pop high, pop low, pop container from stack
	OpSliceLen // get length

	// Map operations
	OpMapIter     // push map iterator
	OpMapIterNext // pop iter, push (key, value, ok)

	// Struct operations
	OpField    // pop struct, push field [field_idx:2]
	OpSetField // pop val, pop struct, struct.field = val [field_idx:2]

	// Pointer operations
	OpAddr       // push address of local [local_idx:2]
	OpIndexAddr  // pop index, pop slice/array, push &slice[index]
	OpDeref      // pop pointer, push *pointer
	OpSetDeref   // pop val, pop pointer, *pointer = val

	// Interface operations
	OpAssert  // type assertion [type_idx:2]
	OpConvert // type conversion [type_idx:2]

	// Function operations
	OpClosure    // create closure [func_idx:2, num_free:1]
	OpMethod     // method value [method_idx:2]
	OpMethodCall // method call [method_idx:2, num_args:1]

	// Goroutine and channel
	OpGo      // start goroutine
	OpSend    // channel send
	OpRecv    // channel recv
	OpTrySend // non-blocking send
	OpTryRecv // non-blocking recv
	OpClose   // close channel
	OpSelect  // select statement

	// Defer/recover
	OpDefer   // defer function call [func_idx:2]
	OpRecover // recover from panic

	// Range
	OpRange     // range over slice/map/string
	OpRangeNext // next iteration

	// Builtin functions
	OpLen     // len()
	OpCap     // cap()
	OpAppend  // append()
	OpCopy    // copy()
	OpDelete  // delete()
	OpPanic   // panic() - should be banned but VM needs to handle it
	OpPrint   // print()
	OpPrintln // println()
	OpNew     // new()
	OpMake    // make() - generic

	// External function call
	OpCallExternal // call external function [func_idx:2, num_args:1]

	// Indirect call (closures, function values)
	OpCallIndirect // call function value on stack [num_args:1]

	// Tuple/multi-value operations
	OpPack   // pack N values from stack into a slice [count:2]
	OpUnpack // unpack a slice onto the stack

	// Halt
	OpHalt // stop execution
)

// String returns the name of the opcode.
func (op OpCode) String() string {
	switch op {
	case OpNop:
		return "NOP"
	case OpPop:
		return "POP"
	case OpDup:
		return "DUP"
	case OpConst:
		return "CONST"
	case OpNil:
		return "NIL"
	case OpTrue:
		return "TRUE"
	case OpFalse:
		return "FALSE"
	case OpLocal:
		return "LOCAL"
	case OpSetLocal:
		return "SETLOCAL"
	case OpGlobal:
		return "GLOBAL"
	case OpSetGlobal:
		return "SETGLOBAL"
	case OpFree:
		return "FREE"
	case OpSetFree:
		return "SETFREE"
	case OpAdd:
		return "ADD"
	case OpSub:
		return "SUB"
	case OpMul:
		return "MUL"
	case OpDiv:
		return "DIV"
	case OpMod:
		return "MOD"
	case OpNeg:
		return "NEG"
	case OpAnd:
		return "AND"
	case OpOr:
		return "OR"
	case OpXor:
		return "XOR"
	case OpAndNot:
		return "ANDNOT"
	case OpLsh:
		return "LSH"
	case OpRsh:
		return "RSH"
	case OpEqual:
		return "EQUAL"
	case OpNotEqual:
		return "NOTEQUAL"
	case OpLess:
		return "LESS"
	case OpLessEq:
		return "LESSEQ"
	case OpGreater:
		return "GREATER"
	case OpGreaterEq:
		return "GREATEREQ"
	case OpNot:
		return "NOT"
	case OpJump:
		return "JUMP"
	case OpJumpTrue:
		return "JUMPTRUE"
	case OpJumpFalse:
		return "JUMPFALSE"
	case OpCall:
		return "CALL"
	case OpReturn:
		return "RETURN"
	case OpReturnVal:
		return "RETURNVAL"
	case OpMakeSlice:
		return "MAKESLICE"
	case OpMakeMap:
		return "MAKEMAP"
	case OpMakeChan:
		return "MAKECHAN"
	case OpMakeArray:
		return "MAKEARRAY"
	case OpMakeStruct:
		return "MAKESTRUCT"
	case OpIndex:
		return "INDEX"
	case OpSetIndex:
		return "SETINDEX"
	case OpSlice:
		return "SLICE"
	case OpSliceLen:
		return "SLICELEN"
	case OpMapIter:
		return "MAPITER"
	case OpMapIterNext:
		return "MAPITERNEXT"
	case OpField:
		return "FIELD"
	case OpSetField:
		return "SETFIELD"
	case OpAddr:
		return "ADDR"
	case OpIndexAddr:
		return "INDEXADDR"
	case OpDeref:
		return "DEREF"
	case OpSetDeref:
		return "SETDEREF"
	case OpAssert:
		return "ASSERT"
	case OpConvert:
		return "CONVERT"
	case OpClosure:
		return "CLOSURE"
	case OpMethod:
		return "METHOD"
	case OpMethodCall:
		return "METHODCALL"
	case OpGo:
		return "GO"
	case OpSend:
		return "SEND"
	case OpRecv:
		return "RECV"
	case OpTrySend:
		return "TRYSEND"
	case OpTryRecv:
		return "TRYRECV"
	case OpClose:
		return "CLOSE"
	case OpSelect:
		return "SELECT"
	case OpDefer:
		return "DEFER"
	case OpRecover:
		return "RECOVER"
	case OpRange:
		return "RANGE"
	case OpRangeNext:
		return "RANGENEXT"
	case OpLen:
		return "LEN"
	case OpCap:
		return "CAP"
	case OpAppend:
		return "APPEND"
	case OpCopy:
		return "COPY"
	case OpDelete:
		return "DELETE"
	case OpPanic:
		return "PANIC"
	case OpPrint:
		return "PRINT"
	case OpPrintln:
		return "PRINTLN"
	case OpNew:
		return "NEW"
	case OpMake:
		return "MAKE"
	case OpCallExternal:
		return "CALLEXTERNAL"
	case OpCallIndirect:
		return "CALLINDIRECT"
	case OpPack:
		return "PACK"
	case OpUnpack:
		return "UNPACK"
	case OpHalt:
		return "HALT"
	default:
		return "UNKNOWN"
	}
}

// OperandWidths maps opcodes to their operand widths.
// 0 means no operands, 1 means 1-byte operand, 2 means 2-byte operand.
var OperandWidths = map[OpCode]int{
	OpConst:        2,
	OpLocal:        2,
	OpSetLocal:     2,
	OpGlobal:       2,
	OpSetGlobal:    2,
	OpFree:         1,
	OpSetFree:      1,
	OpJump:         2,
	OpJumpTrue:     2,
	OpJumpFalse:    2,
	OpCall:         3, // func_idx(2) + num_args(1)
	OpMakeArray:    2,
	OpMakeStruct:   2,
	OpField:        2,
	OpSetField:     2,
	OpAddr:         2,
	OpAssert:       2,
	OpConvert:      2,
	OpClosure:      3, // func_idx(2) + num_free(1)
	OpMethod:       2,
	OpMethodCall:   3, // method_idx(2) + num_args(1)
	OpDefer:        2,
	OpCallExternal: 3, // func_idx(2) + num_args(1)
	OpCallIndirect: 1, // num_args(1)
	OpPack:         2, // count(2)
	OpNew:          2,
	OpMake:         4, // type_idx(2) + size_idx(2)
}

// ReadUint16 reads a 2-byte operand from the bytecode.
func ReadUint16(code []byte, ip int) uint16 {
	return uint16(code[ip])<<8 | uint16(code[ip+1])
}

// WriteUint16 writes a 2-byte operand to the bytecode.
func WriteUint16(code []byte, offset int, val uint16) {
	code[offset] = byte(val >> 8)
	code[offset+1] = byte(val)
}
