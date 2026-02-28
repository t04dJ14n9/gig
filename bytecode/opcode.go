// Package bytecode defines the shared kernel types for the Gig interpreter.
//
// This package contains the bytecode instruction set, compiled program data
// structures, and dependency injection interfaces shared between the compiler
// and VM packages. It serves as the Shared Kernel in the DDD architecture,
// enabling the compiler and VM to be fully decoupled from each other.
//
// # Compilation Process
//
// The compiler translates Go SSA (Static Single Assignment) intermediate representation
// into a custom bytecode format that can be executed by the VM.
//
//  1. SSA Package Input - The compiler receives an SSA package from golang.org/x/tools/go/ssa
//  2. Function Collection - All functions (including nested/anonymous) are collected
//  3. Index Assignment - Each function is assigned a unique index for call instructions
//  4. Per-Function Compilation - Each function is compiled to bytecode:
//     - Symbol table construction for locals, parameters, and free variables
//     - Phi node slot allocation
//     - Basic block compilation in reverse postorder
//     - Jump target patching
//  5. Program Assembly - All compiled functions are combined into a Program
//
// # Bytecode Format
//
// Each instruction consists of:
//   - 1 byte opcode (see OpCode constants)
//   - 0-3 bytes of operands (see OperandWidths)
//
// Most instructions operate on a stack: pop operands, push result.
// Local variables are accessed by index into the frame's local array.
//
// # Example
//
// The Go function:
//
//	func add(a, b int) int {
//	    return a + b
//	}
//
// Compiles to approximately:
//
//	LOCAL 0      ; push local 0 (a)
//	LOCAL 1      ; push local 1 (b)
//	ADD          ; pop a, pop b, push a+b
//	RETURNVAL    ; return top of stack
package bytecode

// OpCode represents a single bytecode instruction.
// Each opcode may have 0-3 bytes of operands following it.
// See OperandWidths for the expected operand size for each opcode.
type OpCode byte

// Opcode constants define the bytecode instruction set.
// They are grouped by category:
//   - Stack operations: push, pop, duplicate values
//   - Constants/locals: access constant pool and local variables
//   - Arithmetic: add, sub, mul, div, mod, neg
//   - Bitwise: and, or, xor, shifts
//   - Comparison: equal, not-equal, less, greater, etc.
//   - Control flow: jumps, calls, returns
//   - Container ops: slice, map, array, struct operations
//   - Pointer ops: address-of, dereference
//   - Interface ops: type assertions, conversions
//   - Concurrency: goroutines, channels
//   - Builtins: len, cap, append, copy, etc.
const (
	// ========================================
	// Stack Operations
	// ========================================

	// OpNop is a no-operation instruction.
	// Used as a placeholder or for debugging.
	OpNop OpCode = iota

	// OpPop discards the top value from the stack.
	OpPop

	// OpDup duplicates the top value on the stack.
	OpDup

	// ========================================
	// Constants and Locals
	// ========================================

	// OpConst pushes a constant from the constant pool.
	// Operands: [const_idx:2] - 2-byte index into constant pool
	OpConst

	// OpNil pushes a nil value onto the stack.
	OpNil

	// OpTrue pushes the boolean value true.
	OpTrue

	// OpFalse pushes the boolean value false.
	OpFalse

	// OpLocal pushes a local variable onto the stack.
	// Operands: [local_idx:2] - 2-byte index into local variable array
	OpLocal

	// OpSetLocal pops a value and stores it in a local variable.
	// Operands: [local_idx:2] - 2-byte index into local variable array
	OpSetLocal

	// OpGlobal pushes a global variable onto the stack.
	// Operands: [global_idx:2] - 2-byte index into global variable array
	OpGlobal

	// OpSetGlobal pops a value and stores it in a global variable.
	// Operands: [global_idx:2] - 2-byte index into global variable array
	OpSetGlobal

	// OpFree pushes a free variable (closure capture) onto the stack.
	// Operands: [free_idx:1] - 1-byte index into free variable array
	OpFree

	// OpSetFree pops a value and stores it in a free variable.
	// Operands: [free_idx:1] - 1-byte index into free variable array
	OpSetFree

	// ========================================
	// Arithmetic Operations
	// ========================================

	// OpAdd pops b, pops a, pushes a + b.
	// Works for int, uint, float, string (concatenation), complex.
	OpAdd

	// OpSub pops b, pops a, pushes a - b.
	// Works for int, uint, float, complex.
	OpSub

	// OpMul pops b, pops a, pushes a * b.
	// Works for int, uint, float, complex.
	OpMul

	// OpDiv pops b, pops a, pushes a / b.
	// Works for int, uint, float, complex.
	OpDiv

	// OpMod pops b, pops a, pushes a % b.
	// Works for int, uint, float.
	OpMod

	// OpNeg pops a, pushes -a.
	// Works for int, float, complex.
	OpNeg

	// ========================================
	// Bitwise Operations
	// ========================================

	// OpAnd pops b, pops a, pushes a & b.
	// Works for int, uint.
	OpAnd

	// OpOr pops b, pops a, pushes a | b.
	// Works for int, uint.
	OpOr

	// OpXor pops b, pops a, pushes a ^ b.
	// Works for int, uint.
	OpXor

	// OpAndNot pops b, pops a, pushes a &^ b (bit clear).
	// Works for int, uint.
	OpAndNot

	// OpLsh pops n, pops a, pushes a << n.
	// Works for int, uint.
	OpLsh

	// OpRsh pops n, pops a, pushes a >> n.
	// Works for int, uint.
	OpRsh

	// ========================================
	// Comparison Operations
	// ========================================

	// OpEqual pops b, pops a, pushes a == b.
	OpEqual

	// OpNotEqual pops b, pops a, pushes a != b.
	OpNotEqual

	// OpLess pops b, pops a, pushes a < b.
	OpLess

	// OpLessEq pops b, pops a, pushes a <= b.
	OpLessEq

	// OpGreater pops b, pops a, pushes a > b.
	OpGreater

	// OpGreaterEq pops b, pops a, pushes a >= b.
	OpGreaterEq

	// ========================================
	// Logical Operations
	// ========================================

	// OpNot pops a, pushes !a.
	// Works for bool.
	OpNot

	// ========================================
	// Control Flow
	// ========================================

	// OpJump jumps to a bytecode offset.
	// Operands: [offset:2] - 2-byte target instruction offset
	OpJump

	// OpJumpTrue pops a condition; if true, jumps to offset.
	// Operands: [offset:2] - 2-byte target instruction offset
	OpJumpTrue

	// OpJumpFalse pops a condition; if false, jumps to offset.
	// Operands: [offset:2] - 2-byte target instruction offset
	OpJumpFalse

	// OpCall calls a compiled function.
	// Operands: [func_idx:2] [num_args:1]
	OpCall

	// OpReturn returns from a function with no value.
	OpReturn

	// OpReturnVal pops a value and returns it.
	OpReturnVal

	// ========================================
	// Container Operations
	// ========================================

	// OpMakeSlice creates a new slice.
	// Stack: [... typeIdx len cap] -> [... slice]
	OpMakeSlice

	// OpMakeMap creates a new map.
	// Stack: [... typeIdx size] -> [... map]
	OpMakeMap

	// OpMakeChan creates a new channel.
	// Stack: [... typeIdx size] -> [... chan]
	OpMakeChan

	// OpMakeArray creates a new array (rarely used).
	// Operands: [type_idx:2]
	OpMakeArray

	// OpMakeStruct creates a new struct (rarely used).
	// Operands: [type_idx:2]
	OpMakeStruct

	// ========================================
	// Index Operations
	// ========================================

	// OpIndex indexes into a container.
	// Stack: [... container key] -> [... value]
	OpIndex

	// OpIndexOk indexes with comma-ok (for maps).
	// Stack: [... container key] -> [... (value, ok) tuple]
	OpIndexOk

	// OpSetIndex sets an element in a container.
	// Stack: [... container key value] -> [...]
	OpSetIndex

	// OpSlice slices a slice/array/string.
	// Stack: [... container low high max] -> [... sliced]
	OpSlice

	// OpSliceLen gets the length (legacy, use OpLen).
	OpSliceLen

	// ========================================
	// Map Operations
	// ========================================

	// OpMapIter creates a map iterator.
	OpMapIter

	// OpMapIterNext advances a map iterator.
	OpMapIterNext

	// ========================================
	// Struct Operations
	// ========================================

	// OpField accesses a struct field.
	// Operands: [field_idx:2]
	// Stack: [... struct] -> [... field_value]
	OpField

	// OpSetField sets a struct field.
	// Operands: [field_idx:2]
	// Stack: [... struct value] -> [...]
	OpSetField

	// ========================================
	// Pointer Operations
	// ========================================

	// OpAddr pushes the address of a local variable.
	// Operands: [local_idx:2]
	OpAddr

	// OpFieldAddr pushes the address of a struct field.
	// Operands: [field_idx:2]
	// Stack: [... struct_ptr] -> [... field_ptr]
	OpFieldAddr

	// OpIndexAddr pushes the address of a slice/array element.
	// Stack: [... container index] -> [... element_ptr]
	OpIndexAddr

	// OpDeref dereferences a pointer.
	// Stack: [... ptr] -> [... *ptr]
	OpDeref

	// OpSetDeref sets the value pointed to.
	// Stack: [... ptr value] -> [...]
	OpSetDeref

	// ========================================
	// Interface Operations
	// ========================================

	// OpAssert performs a type assertion.
	// Operands: [type_idx:2]
	// Stack: [... interface] -> [... (value, ok) tuple]
	OpAssert

	// OpConvert performs a type conversion.
	// Operands: [type_idx:2]
	// Stack: [... value] -> [... converted_value]
	OpConvert

	// ========================================
	// Function Operations
	// ========================================

	// OpClosure creates a closure with captured variables.
	// Operands: [func_idx:2] [num_free:1]
	OpClosure

	// OpMethod gets a method value.
	// Operands: [method_idx:2]
	OpMethod

	// OpMethodCall calls a method.
	// Operands: [method_idx:2] [num_args:1]
	OpMethodCall

	// ========================================
	// Concurrency Operations
	// ========================================

	// OpGoCall starts a new goroutine with a function call.
	// Operands: [func_idx:2, num_args:1]
	// Stack: [... args] -> [...]
	OpGoCall

	// OpGoCallIndirect starts a new goroutine with a closure call.
	// Operands: [num_args:1]
	// Stack: [... closure args...] -> [...]
	OpGoCallIndirect

	// OpSend sends a value on a channel.
	// Stack: [... ch value] -> [...]
	OpSend

	// OpRecv receives a value from a channel.
	// Stack: [... ch] -> [... value]
	OpRecv

	// OpRecvOk receives a value from a channel with comma-ok.
	// Stack: [... ch] -> [... (value, ok) tuple]
	OpRecvOk

	// OpTrySend sends non-blocking.
	OpTrySend

	// OpTryRecv receives non-blocking.
	OpTryRecv

	// OpClose closes a channel.
	// Stack: [... ch] -> [...]
	OpClose

	// OpSelect performs a select statement.
	OpSelect

	// ========================================
	// Defer/Recover
	// ========================================

	// OpDefer defers a function call.
	// Operands: [func_idx:2]
	OpDefer

	// OpRecover recovers from a panic.
	OpRecover

	// ========================================
	// Range Operations
	// ========================================

	// OpRange creates an iterator for range loops.
	// Stack: [... collection] -> [... iterator]
	OpRange

	// OpRangeNext advances an iterator.
	// Stack: [... iterator] -> [... (ok, key, value) tuple]
	OpRangeNext

	// ========================================
	// Builtin Functions
	// ========================================

	// OpLen returns the length of a string, slice, array, map, or channel.
	// Stack: [... value] -> [... len]
	OpLen

	// OpCap returns the capacity of a slice, array, or channel.
	// Stack: [... value] -> [... cap]
	OpCap

	// OpAppend appends elements to a slice.
	// Stack: [... slice elem] -> [... new_slice]
	OpAppend

	// OpCopy copies elements between slices.
	// Stack: [... dst src] -> [... n]
	OpCopy

	// OpDelete deletes a key from a map.
	// Stack: [... map key] -> [...]
	OpDelete

	// OpPanic triggers a panic (banned but handled for error messages).
	// Stack: [... message] -> [panic]
	OpPanic

	// OpPrint prints values.
	// Operands: [count:1]
	OpPrint

	// OpPrintln prints values with newlines.
	// Operands: [count:1]
	OpPrintln

	// OpNew allocates a new pointer.
	// Operands: [type_idx:2]
	OpNew

	// OpMake allocates with make (generic).
	// Operands: [type_idx:2] [size_idx:2]
	OpMake

	// ========================================
	// External Function Calls
	// ========================================

	// OpCallExternal calls an external (native Go) function.
	// Operands: [func_idx:2] [num_args:1]
	OpCallExternal

	// OpCallIndirect calls a function value (closure, function variable).
	// Operands: [num_args:1]
	// Stack: [... func_value args...] -> [... result]
	OpCallIndirect

	// ========================================
	// Tuple/Multi-value Operations
	// ========================================

	// OpPack packs N values into a slice.
	// Operands: [count:2]
	OpPack

	// OpUnpack unpacks a slice onto the stack.
	OpUnpack

	// OpHalt stops execution (for debugging).
	OpHalt

	// ========================================
	// Superinstructions (fused ops for hot loops)
	// ========================================

	// OpAddLocalLocal pops nothing; loads locals[A] and locals[B], pushes a+b (int fast path).
	// Operands: [local_a:2] [local_b:2]
	OpAddLocalLocal

	// OpSubLocalLocal loads locals[A] and locals[B], pushes a-b (int fast path).
	// Operands: [local_a:2] [local_b:2]
	OpSubLocalLocal

	// OpMulLocalLocal loads locals[A] and locals[B], pushes a*b (int fast path).
	// Operands: [local_a:2] [local_b:2]
	OpMulLocalLocal

	// OpAddLocalConst loads local[A] and const[B], pushes a+b (int fast path).
	// Operands: [local_a:2] [const_b:2]
	OpAddLocalConst

	// OpSubLocalConst loads local[A] and const[B], pushes a-b (int fast path).
	// Operands: [local_a:2] [const_b:2]
	OpSubLocalConst

	// OpLessLocalLocal loads locals[A] and locals[B], pushes a<b (int) then JumpTrue.
	// Operands: [local_a:2] [local_b:2] [offset:2]
	OpLessLocalLocalJumpTrue

	// OpLessLocalConstJumpTrue loads local[A] and const[B], jumps if a<b.
	// Operands: [local_a:2] [const_b:2] [offset:2]
	OpLessLocalConstJumpTrue

	// OpLessEqLocalConstJumpTrue loads local[A] and const[B], jumps if a<=b.
	// Operands: [local_a:2] [const_b:2] [offset:2]
	OpLessEqLocalConstJumpTrue

	// OpGreaterLocalLocalJumpTrue loads locals[A] and locals[B], jumps if a>b.
	// Operands: [local_a:2] [local_b:2] [offset:2]
	OpGreaterLocalLocalJumpTrue

	// OpLessLocalLocalJumpFalse loads locals[A] and locals[B], jumps if NOT a<b.
	// Operands: [local_a:2] [local_b:2] [offset:2]
	OpLessLocalLocalJumpFalse

	// OpLessLocalConstJumpFalse loads local[A] and const[B], jumps if NOT a<b.
	// Operands: [local_a:2] [const_b:2] [offset:2]
	OpLessLocalConstJumpFalse

	// OpAddSetLocal pops two values, adds, stores to local[A].
	// Operands: [local_a:2]
	OpAddSetLocal

	// OpSubSetLocal pops two values, subs, stores to local[A].
	// Operands: [local_a:2]
	OpSubSetLocal

	// OpLocalLocalAddSetLocal loads locals[A]+locals[B], stores to local[C].
	// Operands: [local_a:2] [local_b:2] [local_c:2]
	OpLocalLocalAddSetLocal

	// OpLocalConstAddSetLocal loads local[A]+const[B], stores to local[C].
	// Operands: [local_a:2] [const_b:2] [local_c:2]
	OpLocalConstAddSetLocal

	// OpLocalConstSubSetLocal loads local[A]-const[B], stores to local[C].
	// Operands: [local_a:2] [const_b:2] [local_c:2]
	OpLocalConstSubSetLocal
)

// String returns the name of the opcode as a human-readable string.
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
	case OpFieldAddr:
		return "FIELDADDR"
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
	case OpGoCall:
		return "GOCALL"
	case OpGoCallIndirect:
		return "GOCALLINDIRECT"
	case OpSend:
		return "SEND"
	case OpRecv:
		return "RECV"
	case OpRecvOk:
		return "RECVOK"
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
	case OpAddLocalLocal:
		return "ADDLOCALLOCAL"
	case OpSubLocalLocal:
		return "SUBLOCALLOCAL"
	case OpMulLocalLocal:
		return "MULLOCALLOCAL"
	case OpAddLocalConst:
		return "ADDLOCALCONST"
	case OpSubLocalConst:
		return "SUBLOCALCONST"
	case OpLessLocalLocalJumpTrue:
		return "LESSLOCALLOCALJUMPTRUE"
	case OpLessLocalConstJumpTrue:
		return "LESSLOCALCONSTJUMPTRUE"
	case OpLessEqLocalConstJumpTrue:
		return "LESSEQLOCALCONSTJUMPTRUE"
	case OpGreaterLocalLocalJumpTrue:
		return "GREATERLOCALLOCALJUMPTRUE"
	case OpLessLocalLocalJumpFalse:
		return "LESSLOCALLOCALJUMPFALSE"
	case OpLessLocalConstJumpFalse:
		return "LESSLOCALCONSTJUMPFALSE"
	case OpAddSetLocal:
		return "ADDSETLOCAL"
	case OpSubSetLocal:
		return "SUBSETLOCAL"
	case OpLocalLocalAddSetLocal:
		return "LOCALLOCALADDSETLOCAL"
	case OpLocalConstAddSetLocal:
		return "LOCALCONSTADDSETLOCAL"
	case OpLocalConstSubSetLocal:
		return "LOCALCONSTSUBSETLOCAL"
	default:
		return "UNKNOWN"
	}
}

// OperandWidths maps opcodes to their operand widths.
// 0 means no operands, 1 means 1-byte operand, 2 means 2-byte operand.
var OperandWidths = map[OpCode]int{
	OpConst:          2,
	OpLocal:          2,
	OpSetLocal:       2,
	OpGlobal:         2,
	OpSetGlobal:      2,
	OpFree:           1,
	OpSetFree:        1,
	OpJump:           2,
	OpJumpTrue:       2,
	OpJumpFalse:      2,
	OpCall:           3, // func_idx(2) + num_args(1)
	OpMakeArray:      2,
	OpMakeStruct:     2,
	OpField:          2,
	OpSetField:       2,
	OpAddr:           2,
	OpFieldAddr:      2,
	OpAssert:         2,
	OpConvert:        2,
	OpClosure:        3, // func_idx(2) + num_free(1)
	OpMethod:         2,
	OpMethodCall:     3, // method_idx(2) + num_args(1)
	OpDefer:          2,
	OpCallExternal:   3, // func_idx(2) + num_args(1)
	OpCallIndirect:   1, // num_args(1)
	OpGoCall:         3, // func_idx(2) + num_args(1)
	OpGoCallIndirect: 1, // num_args(1)
	OpSelect:         2, // meta_idx(2)
	OpPack:           2, // count(2)
	OpNew:            2,
	OpMake:           4, // type_idx(2) + size_idx(2)
	OpPrint:          1, // count(1)
	OpPrintln:        1, // count(1)

	// Superinstruction operand widths
	OpAddLocalLocal:             4, // local_a(2) + local_b(2)
	OpSubLocalLocal:             4, // local_a(2) + local_b(2)
	OpMulLocalLocal:             4, // local_a(2) + local_b(2)
	OpAddLocalConst:             4, // local_a(2) + const_b(2)
	OpSubLocalConst:             4, // local_a(2) + const_b(2)
	OpLessLocalLocalJumpTrue:    6, // local_a(2) + local_b(2) + offset(2)
	OpLessLocalConstJumpTrue:    6, // local_a(2) + const_b(2) + offset(2)
	OpLessEqLocalConstJumpTrue:  6, // local_a(2) + const_b(2) + offset(2)
	OpGreaterLocalLocalJumpTrue: 6, // local_a(2) + local_b(2) + offset(2)
	OpLessLocalLocalJumpFalse:   6, // local_a(2) + local_b(2) + offset(2)
	OpLessLocalConstJumpFalse:   6, // local_a(2) + const_b(2) + offset(2)
	OpAddSetLocal:               2, // local_a(2)
	OpSubSetLocal:               2, // local_a(2)
	OpLocalLocalAddSetLocal:     6, // local_a(2) + local_b(2) + local_c(2)
	OpLocalConstAddSetLocal:     6, // local_a(2) + const_b(2) + local_c(2)
	OpLocalConstSubSetLocal:     6, // local_a(2) + const_b(2) + local_c(2)
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

// SelectMeta stores metadata for a select statement.
// It is stored in the constant pool and referenced by OpSelect.
type SelectMeta struct {
	NumStates int    // total number of select cases (excluding default)
	Blocking  bool   // true if no default branch
	Dirs      []bool // direction for each case: true=send, false=recv
	NumRecv   int    // number of recv cases (determines result tuple size)
}
