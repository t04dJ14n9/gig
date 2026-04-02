// constants.go defines named constants replacing magic numbers across the VM.
package vm

import "github.com/t04dJ14n9/gig/model/bytecode"

// VM-wide constants extracted from magic numbers across the codebase.
const (
	// initialStackSize is the starting size of the operand stack for main
	// and goroutine child VMs.
	initialStackSize = 1024

	// deferVMStackSize is the starting stack size for child VMs that execute
	// deferred functions and method resolution.
	deferVMStackSize = 256

	// contextCheckInterval is the number of instructions between context
	// cancellation checks. Must be a power of two for efficient masking.
	contextCheckInterval = 1024

	// contextCheckMask is used for bitwise AND to check if it's time for a
	// context cancellation check (contextCheckInterval - 1).
	contextCheckMask = contextCheckInterval - 1

	// sliceEndSentinel is the sentinel value meaning "use the container's
	// length" in slice operations (OpSlice high operand).
	sliceEndSentinel = bytecode.SliceEndSentinel

	// noSourceLocalSentinel is the sentinel value meaning "no source local"
	// in convert operations.
	noSourceLocalSentinel = bytecode.NoSourceLocal

	// maxStackSize is the hard ceiling for the operand stack. Each slot is
	// a 32-byte value.Value, so 1<<20 slots = 32 MB per VM.
	maxStackSize = 1 << 20

	// defaultMaxGoroutines is the default limit on concurrent interpreter
	// goroutines per program.
	defaultMaxGoroutines = 10000
)
