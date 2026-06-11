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
