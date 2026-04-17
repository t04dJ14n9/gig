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
//
// Sentinel values shared between compiler and VM.
const (
	// SliceEndSentinel is emitted by the compiler as the "high" or "max" operand
	// of a slice expression when the user omits it (e.g., a[1:]). The VM
	// interprets this as "use the container's length."
	SliceEndSentinel = 0xFFFF

	// NoSourceLocal is emitted by the compiler as the source-local operand of
	// OpChangeType when the source value is not a local variable. The VM uses
	// it to skip the source-local update that shares backing arrays.
	NoSourceLocal = 0xFFFF
)

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

	// OpReal pops a complex number, pushes its real part as float64.
	OpReal

	// OpImag pops a complex number, pushes its imaginary part as float64.
	OpImag

	// OpComplex pops imag, pops real, pushes complex(real, imag).
	OpComplex

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

	// OpMakeInterface wraps a value in a Go interface, preserving type information.
	// Critical for typed nil: var p *T = nil; var e error = p → e is non-nil.
	// Stack: [... value] -> [... interface_value]
	// Operands: [iface_type_idx:2] [concrete_type_idx:2]
	OpMakeInterface

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

	// _ (placeholder: was OpSliceLen, removed as dead code - use OpLen instead)

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

	// OpChangeType performs a named-type conversion (ChangeType in SSA).
	// Unlike OpConvert, this also updates the source local variable so that
	// slice aliasing works correctly (e.g., sort.IntSlice(s) shares s's backing array).
	// Operands: [type_idx:2] [src_local:2]
	// Stack: [... value] -> [... converted_value]
	OpChangeType

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

	// OpDeferIndirect defers a closure call.
	// Operands: [num_args:1]
	// Stack: [... closure args...] -> [...]
	OpDeferIndirect

	// OpDeferExternal defers an external function or method call.
	// Operands: [func_idx:2, num_args:1]
	// Stack: [... args...] -> [...]
	OpDeferExternal

	// OpRunDefers executes all pending deferred calls synchronously.
	// This is used for named return values where defers may modify return values
	// before the function actually returns.
	OpRunDefers

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

	// OpLessEqLocalConstJumpFalse loads local[A] and const[B], jumps if NOT a<=b.
	// Operands: [local_a:2] [const_b:2] [offset:2]
	OpLessEqLocalConstJumpFalse

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

	// ========================================
	// Integer-specialized superinstructions
	// These operate on intLocals []int64 directly (8 bytes per op instead of 32).
	// Emitted by the peephole optimizer when compile-time type analysis confirms int types.
	// ========================================

	// OpIntLocalConstAddSetLocal: intLocals[C] = intLocals[A] + intConsts[B]
	// Operands: [local_a:2] [const_b:2] [local_c:2]
	OpIntLocalConstAddSetLocal

	// OpIntLocalConstSubSetLocal: intLocals[C] = intLocals[A] - intConsts[B]
	// Operands: [local_a:2] [const_b:2] [local_c:2]
	OpIntLocalConstSubSetLocal

	// OpIntLocalLocalAddSetLocal: intLocals[C] = intLocals[A] + intLocals[B]
	// Operands: [local_a:2] [local_b:2] [local_c:2]
	OpIntLocalLocalAddSetLocal

	// OpIntLessLocalConstJumpFalse: if intLocals[A] >= intConsts[B] { goto offset }
	// Operands: [local_a:2] [const_b:2] [offset:2]
	OpIntLessLocalConstJumpFalse

	// OpIntLessEqLocalConstJumpTrue: if intLocals[A] <= intConsts[B] { goto offset }
	// Operands: [local_a:2] [const_b:2] [offset:2]
	OpIntLessEqLocalConstJumpTrue

	// OpIntLessEqLocalConstJumpFalse: if intLocals[A] > intConsts[B] { goto offset }
	// Operands: [local_a:2] [const_b:2] [offset:2]
	OpIntLessEqLocalConstJumpFalse

	// OpIntLessLocalLocalJumpFalse: if intLocals[A] >= intLocals[B] { goto offset }
	// Operands: [local_a:2] [local_b:2] [offset:2]
	OpIntLessLocalLocalJumpFalse

	// OpIntGreaterLocalLocalJumpTrue: if intLocals[A] > intLocals[B] { goto offset }
	// Operands: [local_a:2] [local_b:2] [offset:2]
	OpIntGreaterLocalLocalJumpTrue

	// OpIntSetLocal: intLocals[idx] = pop().RawInt()
	// Operands: [local_idx:2]
	OpIntSetLocal

	// OpIntLocal: push(MakeInt(intLocals[idx]))
	// Operands: [local_idx:2]
	OpIntLocal

	// OpIntLessLocalConstJumpTrue: if intLocals[A] < intConsts[B] { goto offset }
	// Operands: [local_a:2] [const_b:2] [offset:2]
	OpIntLessLocalConstJumpTrue

	// OpIntLessLocalLocalJumpTrue: if intLocals[A] < intLocals[B] { goto offset }
	// Operands: [local_a:2] [local_b:2] [offset:2]
	OpIntLessLocalLocalJumpTrue

	// OpIntMoveLocal: intLocals[dst] = intLocals[src]; locals[dst] = locals[src]
	// Eliminates push+pop phi-move pattern: INTLOCAL(src) INTSETLOCAL(dst)
	// Operands: [src:2] [dst:2]
	OpIntMoveLocal

	// OpIntSliceGet: dest = intSlice[index]
	// Fuses: LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) DEREF SETLOCAL(v)
	// into a single dispatch that reads intLocals[index] from the []int64 in locals[slice].
	// Operands: [slice_local:2] [index_local:2] [dest_local:2]
	OpIntSliceGet

	// OpIntSliceSet: intSlice[index] = val
	// Fuses: LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) LOCAL(val) SETDEREF
	// into a single dispatch that writes intLocals[val] into the []int64 in locals[slice].
	// Operands: [slice_local:2] [index_local:2] [val_local:2]
	OpIntSliceSet

	// OpIntSliceSetConst: intSlice[index] = const
	// Fuses: LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) CONST(val) SETDEREF
	// into a single dispatch that writes intConsts[val] into the []int64 in locals[slice].
	// Operands: [slice_local:2] [index_local:2] [const_val:2]
	OpIntSliceSetConst

	// OpLocalLocalSubSetLocal: locals[C] = locals[A] - locals[B]
	// Operands: [local_a:2] [local_b:2] [local_c:2]
	OpLocalLocalSubSetLocal

	// OpLocalLocalMulSetLocal: locals[C] = locals[A] * locals[B]
	// Operands: [local_a:2] [local_b:2] [local_c:2]
	OpLocalLocalMulSetLocal

	// OpLocalConstMulSetLocal: locals[C] = locals[A] * consts[B]
	// Operands: [local_a:2] [const_b:2] [local_c:2]
	OpLocalConstMulSetLocal

	// OpIntLocalLocalSubSetLocal: intLocals[C] = intLocals[A] - intLocals[B]
	// Operands: [local_a:2] [local_b:2] [local_c:2]
	OpIntLocalLocalSubSetLocal

	// OpIntLocalLocalMulSetLocal: intLocals[C] = intLocals[A] * intLocals[B]
	// Operands: [local_a:2] [local_b:2] [local_c:2]
	OpIntLocalLocalMulSetLocal

	// OpIntLocalConstMulSetLocal: intLocals[C] = intLocals[A] * intConsts[B]
	// Operands: [local_a:2] [const_b:2] [local_c:2]
	OpIntLocalConstMulSetLocal
)

// operandWidthTable is a lookup table for opcode operand widths.
// Using a fixed-size array gives O(1) access with no hash overhead, unlike a map.
// Index is the opcode byte value; value is the total operand byte width.
var operandWidthTable = buildOperandWidthTable()

func buildOperandWidthTable() [256]int {
	var t [256]int

	t[OpConst] = 2
	t[OpLocal] = 2
	t[OpSetLocal] = 2
	t[OpGlobal] = 2
	t[OpSetGlobal] = 2
	t[OpFree] = 1
	t[OpSetFree] = 1
	t[OpJump] = 2
	t[OpJumpTrue] = 2
	t[OpJumpFalse] = 2
	t[OpCall] = 3 // func_idx(2) + num_args(1)
	t[OpMakeArray] = 2
	t[OpMakeStruct] = 2
	t[OpField] = 2
	t[OpSetField] = 2
	t[OpAddr] = 2
	t[OpFieldAddr] = 2
	t[OpAssert] = 2
	t[OpConvert] = 2
	t[OpChangeType] = 4
	t[OpClosure] = 3 // func_idx(2) + num_free(1)
	t[OpMethod] = 2
	t[OpMethodCall] = 3 // method_idx(2) + num_args(1)
	t[OpDefer] = 2
	t[OpDeferIndirect] = 2  // num_args(2)
	t[OpDeferExternal] = 3  // func_idx(2) + num_args(1)
	t[OpCallExternal] = 3   // func_idx(2) + num_args(1)
	t[OpCallIndirect] = 1   // num_args(1)
	t[OpGoCall] = 3         // func_idx(2) + num_args(1)
	t[OpGoCallIndirect] = 1 // num_args(1)
	t[OpSelect] = 2         // meta_idx(2)
	t[OpPack] = 2           // count(2)
	t[OpNew] = 2
	t[OpMake] = 4    // type_idx(2) + size_idx(2)
	t[OpMakeInterface] = 4 // iface_type_idx(2) + concrete_type_idx(2)
	t[OpPrint] = 1   // count(1)
	t[OpPrintln] = 1 // count(1)

	// Superinstruction operand widths
	t[OpAddLocalLocal] = 4             // local_a(2) + local_b(2)
	t[OpSubLocalLocal] = 4             // local_a(2) + local_b(2)
	t[OpMulLocalLocal] = 4             // local_a(2) + local_b(2)
	t[OpAddLocalConst] = 4             // local_a(2) + const_b(2)
	t[OpSubLocalConst] = 4             // local_a(2) + const_b(2)
	t[OpLessLocalLocalJumpTrue] = 6    // local_a(2) + local_b(2) + offset(2)
	t[OpLessLocalConstJumpTrue] = 6    // local_a(2) + const_b(2) + offset(2)
	t[OpLessEqLocalConstJumpTrue] = 6  // local_a(2) + const_b(2) + offset(2)
	t[OpGreaterLocalLocalJumpTrue] = 6 // local_a(2) + local_b(2) + offset(2)
	t[OpLessLocalLocalJumpFalse] = 6   // local_a(2) + local_b(2) + offset(2)
	t[OpLessLocalConstJumpFalse] = 6   // local_a(2) + const_b(2) + offset(2)
	t[OpLessEqLocalConstJumpFalse] = 6 // local_a(2) + const_b(2) + offset(2)
	t[OpAddSetLocal] = 2               // local_a(2)
	t[OpSubSetLocal] = 2               // local_a(2)
	t[OpLocalLocalAddSetLocal] = 6     // local_a(2) + local_b(2) + local_c(2)
	t[OpLocalConstAddSetLocal] = 6     // local_a(2) + const_b(2) + local_c(2)
	t[OpLocalConstSubSetLocal] = 6     // local_a(2) + const_b(2) + local_c(2)

	// Integer-specialized operand widths
	t[OpIntLocalConstAddSetLocal] = 6     // local_a(2) + const_b(2) + local_c(2)
	t[OpIntLocalConstSubSetLocal] = 6     // local_a(2) + const_b(2) + local_c(2)
	t[OpIntLocalLocalAddSetLocal] = 6     // local_a(2) + local_b(2) + local_c(2)
	t[OpIntLessLocalConstJumpFalse] = 6   // local_a(2) + const_b(2) + offset(2)
	t[OpIntLessEqLocalConstJumpTrue] = 6  // local_a(2) + const_b(2) + offset(2)
	t[OpIntLessEqLocalConstJumpFalse] = 6 // local_a(2) + const_b(2) + offset(2)
	t[OpIntLessLocalLocalJumpFalse] = 6   // local_a(2) + local_b(2) + offset(2)
	t[OpIntGreaterLocalLocalJumpTrue] = 6 // local_a(2) + local_b(2) + offset(2)
	t[OpIntSetLocal] = 2                  // local_idx(2)
	t[OpIntLocal] = 2                     // local_idx(2)
	t[OpIntLessLocalConstJumpTrue] = 6    // local_a(2) + const_b(2) + offset(2)
	t[OpIntLessLocalLocalJumpTrue] = 6    // local_a(2) + local_b(2) + offset(2)
	t[OpIntMoveLocal] = 4                 // src(2) + dst(2)
	t[OpIntSliceGet] = 6                  // slice_local(2) + index_local(2) + dest_local(2)
	t[OpIntSliceSet] = 6                  // slice_local(2) + index_local(2) + val_local(2)
	t[OpIntSliceSetConst] = 6             // slice_local(2) + index_local(2) + const_val(2)

	// New superinstruction operand widths
	t[OpLocalLocalSubSetLocal] = 6    // local_a(2) + local_b(2) + local_c(2)
	t[OpLocalLocalMulSetLocal] = 6    // local_a(2) + local_b(2) + local_c(2)
	t[OpLocalConstMulSetLocal] = 6    // local_a(2) + const_b(2) + local_c(2)
	t[OpIntLocalLocalSubSetLocal] = 6 // local_a(2) + local_b(2) + local_c(2)
	t[OpIntLocalLocalMulSetLocal] = 6 // local_a(2) + local_b(2) + local_c(2)
	t[OpIntLocalConstMulSetLocal] = 6 // local_a(2) + const_b(2) + local_c(2)

	return t
}

// unknownOpName is the string returned for unrecognized opcodes.
const unknownOpName = "UNKNOWN"

// opNameTable maps OpCode values to their string representations.
// Built at init time for O(1) lookup performance.
var opNameTable [256]string

func init() { //nolint:gochecknoinits,decorder // table init placed after var declarations for readability
	// Initialize all entries to "UNKNOWN"
	for i := range opNameTable {
		opNameTable[i] = unknownOpName
	}

	// Populate known opcodes
	opNameTable[OpAdd] = "ADD"
	opNameTable[OpAddLocalConst] = "ADDLOCALCONST"
	opNameTable[OpAddLocalLocal] = "ADDLOCALLOCAL"
	opNameTable[OpAddSetLocal] = "ADDSETLOCAL"
	opNameTable[OpAddr] = "ADDR"
	opNameTable[OpAnd] = "AND"
	opNameTable[OpAndNot] = "ANDNOT"
	opNameTable[OpAppend] = "APPEND"
	opNameTable[OpAssert] = "ASSERT"
	opNameTable[OpCall] = "CALL"
	opNameTable[OpCallExternal] = "CALLEXTERNAL"
	opNameTable[OpCallIndirect] = "CALLINDIRECT"
	opNameTable[OpCap] = "CAP"
	opNameTable[OpChangeType] = "CHANGETYPE"
	opNameTable[OpClose] = "CLOSE"
	opNameTable[OpClosure] = "CLOSURE"
	opNameTable[OpComplex] = "COMPLEX"
	opNameTable[OpConst] = "CONST"
	opNameTable[OpConvert] = "CONVERT"
	opNameTable[OpCopy] = "COPY"
	opNameTable[OpDefer] = "DEFER"
	opNameTable[OpDeferExternal] = "DEFEREXTERNAL"
	opNameTable[OpDeferIndirect] = "DEFERINDIRECT"
	opNameTable[OpDelete] = "DELETE"
	opNameTable[OpDeref] = "DEREF"
	opNameTable[OpDiv] = "DIV"
	opNameTable[OpDup] = "DUP"
	opNameTable[OpEqual] = "EQUAL"
	opNameTable[OpFalse] = "FALSE"
	opNameTable[OpField] = "FIELD"
	opNameTable[OpFieldAddr] = "FIELDADDR"
	opNameTable[OpFree] = "FREE"
	opNameTable[OpGlobal] = "GLOBAL"
	opNameTable[OpGoCall] = "GOCALL"
	opNameTable[OpGoCallIndirect] = "GOCALLINDIRECT"
	opNameTable[OpGreater] = "GREATER"
	opNameTable[OpGreaterEq] = "GREATEREQ"
	opNameTable[OpGreaterLocalLocalJumpTrue] = "GREATERLOCALLOCALJUMPTRUE"
	opNameTable[OpHalt] = "HALT"
	opNameTable[OpImag] = "IMAG"
	opNameTable[OpIndex] = "INDEX"
	opNameTable[OpIndexAddr] = "INDEXADDR"
	opNameTable[OpIndexOk] = "INDEXOK"
	opNameTable[OpIntGreaterLocalLocalJumpTrue] = "INTGREATERLOCALLOCALJUMPTRUE"
	opNameTable[OpIntLessEqLocalConstJumpFalse] = "INTLESSEQLOCALCONSTJUMPFALSE"
	opNameTable[OpIntLessEqLocalConstJumpTrue] = "INTLESSEQLOCALCONSTJUMPTRUE"
	opNameTable[OpIntLessLocalConstJumpFalse] = "INTLESSLOCALCONSTJUMPFALSE"
	opNameTable[OpIntLessLocalConstJumpTrue] = "INTLESSLOCALCONSTJUMPTRUE"
	opNameTable[OpIntLessLocalLocalJumpFalse] = "INTLESSLOCALLOCALJUMPFALSE"
	opNameTable[OpIntLessLocalLocalJumpTrue] = "INTLESSLOCALLOCALJUMPTRUE"
	opNameTable[OpIntLocal] = "INTLOCAL"
	opNameTable[OpIntLocalConstAddSetLocal] = "INTLOCALCONSTADDSETLOCAL"
	opNameTable[OpIntLocalConstMulSetLocal] = "INTLOCALCONSTMULSETLOCAL"
	opNameTable[OpIntLocalConstSubSetLocal] = "INTLOCALCONSTSUBSETLOCAL"
	opNameTable[OpIntLocalLocalAddSetLocal] = "INTLOCALLOCALADDSETLOCAL"
	opNameTable[OpIntLocalLocalMulSetLocal] = "INTLOCALLOCALMULSETLOCAL"
	opNameTable[OpIntLocalLocalSubSetLocal] = "INTLOCALLOCALSUBSETLOCAL"
	opNameTable[OpIntMoveLocal] = "INTMOVELOCAL"
	opNameTable[OpIntSetLocal] = "INTSETLOCAL"
	opNameTable[OpIntSliceGet] = "INTSLICEGET"
	opNameTable[OpIntSliceSet] = "INTSLICESET"
	opNameTable[OpIntSliceSetConst] = "INTSLICESETCONST"
	opNameTable[OpJump] = "JUMP"
	opNameTable[OpJumpFalse] = "JUMPFALSE"
	opNameTable[OpJumpTrue] = "JUMPTRUE"
	opNameTable[OpLen] = "LEN"
	opNameTable[OpLess] = "LESS"
	opNameTable[OpLessEq] = "LESSEQ"
	opNameTable[OpLessEqLocalConstJumpFalse] = "LESSEQLOCALCONSTJUMPFALSE"
	opNameTable[OpLessEqLocalConstJumpTrue] = "LESSEQLOCALCONSTJUMPTRUE"
	opNameTable[OpLessLocalConstJumpFalse] = "LESSLOCALCONSTJUMPFALSE"
	opNameTable[OpLessLocalConstJumpTrue] = "LESSLOCALCONSTJUMPTRUE"
	opNameTable[OpLessLocalLocalJumpFalse] = "LESSLOCALLOCALJUMPFALSE"
	opNameTable[OpLessLocalLocalJumpTrue] = "LESSLOCALLOCALJUMPTRUE"
	opNameTable[OpLocal] = "LOCAL"
	opNameTable[OpLocalConstAddSetLocal] = "LOCALCONSTADDSETLOCAL"
	opNameTable[OpLocalConstMulSetLocal] = "LOCALCONSTMULSETLOCAL"
	opNameTable[OpLocalConstSubSetLocal] = "LOCALCONSTSUBSETLOCAL"
	opNameTable[OpLocalLocalAddSetLocal] = "LOCALLOCALADDSETLOCAL"
	opNameTable[OpLocalLocalMulSetLocal] = "LOCALLOCALMULSETLOCAL"
	opNameTable[OpLocalLocalSubSetLocal] = "LOCALLOCALSUBSETLOCAL"
	opNameTable[OpLsh] = "LSH"
	opNameTable[OpMake] = "MAKE"
	opNameTable[OpMakeArray] = "MAKEARRAY"
	opNameTable[OpMakeChan] = "MAKECHAN"
	opNameTable[OpMakeInterface] = "MAKEINTERFACE"
	opNameTable[OpMakeMap] = "MAKEMAP"
	opNameTable[OpMakeSlice] = "MAKESLICE"
	opNameTable[OpMakeStruct] = "MAKESTRUCT"
	opNameTable[OpMapIter] = "MAPITER"
	opNameTable[OpMapIterNext] = "MAPITERNEXT"
	opNameTable[OpMethod] = "METHOD"
	opNameTable[OpMethodCall] = "METHODCALL"
	opNameTable[OpMod] = "MOD"
	opNameTable[OpMul] = "MUL"
	opNameTable[OpMulLocalLocal] = "MULLOCALLOCAL"
	opNameTable[OpNeg] = "NEG"
	opNameTable[OpNew] = "NEW"
	opNameTable[OpNil] = "NIL"
	opNameTable[OpNop] = "NOP"
	opNameTable[OpNot] = "NOT"
	opNameTable[OpNotEqual] = "NOTEQUAL"
	opNameTable[OpOr] = "OR"
	opNameTable[OpPack] = "PACK"
	opNameTable[OpPanic] = "PANIC"
	opNameTable[OpPop] = "POP"
	opNameTable[OpPrint] = "PRINT"
	opNameTable[OpPrintln] = "PRINTLN"
	opNameTable[OpRange] = "RANGE"
	opNameTable[OpRangeNext] = "RANGENEXT"
	opNameTable[OpReal] = "REAL"
	opNameTable[OpRecover] = "RECOVER"
	opNameTable[OpRecv] = "RECV"
	opNameTable[OpRecvOk] = "RECVOK"
	opNameTable[OpReturn] = "RETURN"
	opNameTable[OpReturnVal] = "RETURNVAL"
	opNameTable[OpRsh] = "RSH"
	opNameTable[OpRunDefers] = "RUNDEFERS"
	opNameTable[OpSelect] = "SELECT"
	opNameTable[OpSend] = "SEND"
	opNameTable[OpSetDeref] = "SETDEREF"
	opNameTable[OpSetField] = "SETFIELD"
	opNameTable[OpSetFree] = "SETFREE"
	opNameTable[OpSetGlobal] = "SETGLOBAL"
	opNameTable[OpSetIndex] = "SETINDEX"
	opNameTable[OpSetLocal] = "SETLOCAL"
	opNameTable[OpSlice] = "SLICE"
	opNameTable[OpSub] = "SUB"
	opNameTable[OpSubLocalConst] = "SUBLOCALCONST"
	opNameTable[OpSubLocalLocal] = "SUBLOCALLOCAL"
	opNameTable[OpSubSetLocal] = "SUBSETLOCAL"
	opNameTable[OpTrue] = "TRUE"
	opNameTable[OpTryRecv] = "TRYRECV"
	opNameTable[OpTrySend] = "TRYSEND"
	opNameTable[OpUnpack] = "UNPACK"
	opNameTable[OpXor] = "XOR"
}

// String returns the name of the opcode as a human-readable string.
func (op OpCode) String() string {
	if int(op) < len(opNameTable) {
		return opNameTable[op]
	}
	return unknownOpName
}

// OperandWidth returns the operand byte width for an opcode using O(1) array lookup.
func OperandWidth(op OpCode) int {
	return operandWidthTable[op]
}
