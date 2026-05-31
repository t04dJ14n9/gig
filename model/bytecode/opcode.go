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

	OpPop OpCode = iota

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

	// ========================================
	// Concurrency Operations
	// ========================================

	// OpGoCall starts a new goroutine with a function call.
	// Operands: [func_idx:2, num_args:1]
	// Stack: [... args] -> [...]
	OpGoCall

	// OpGoCallExternal starts a new goroutine with an external function or method call.
	// Operands: [func_idx:2, num_args:1]
	// Stack: [... args] -> [...]
	OpGoCallExternal

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
