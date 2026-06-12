package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeClosure(frame *Frame) {
	funcIdx := frame.readUint16()
	numFree := int(frame.readByte())
	fn := v.compiledFunctionByIndex(funcIdx)
	if fn == nil {
		v.discardFreeVars(numFree)
		v.push(value.MakeNil())
		return
	}

	closure := v.prepareClosure(fn, numFree)
	v.captureFreeVars(closure, numFree)
	v.push(value.MakeFunc(closure))
}

func (v *vm) compiledFunctionByIndex(funcIdx uint16) *bytecode.CompiledFunction {
	// Function indexes are emitted by the compiler and resolved through the
	// program's dense lookup table; missing entries still consume stack inputs
	// so malformed bytecode leaves the VM stack balanced.
	if int(funcIdx) >= len(v.program.FuncByIndex) {
		return nil
	}
	return v.program.FuncByIndex[funcIdx]
}

func (v *vm) prepareClosure(fn *bytecode.CompiledFunction, numFree int) *Closure {
	closure := getClosure(fn, numFree)
	closure.Program = v.program
	closure.InitialGlobals = v.initialGlobals
	// Propagate runtime context so closures converted to Go functions (via
	// reflect.MakeFunc for sync.Once.Do etc.) use the same globals, goroutine
	// tracker, context, and external-call table as the parent VM.
	closure.Shared = v.shared
	closure.Goroutines = v.goroutines
	closure.Ctx = v.ctx
	return closure
}

func (v *vm) captureFreeVars(closure *Closure, numFree int) {
	for i := numFree - 1; i >= 0; i-- {
		captured := v.pop()
		// Captured variables live in mutable slots so OpFree and OpSetFree
		// share state across closures, including reflect-backed pointer values.
		slot := new(value.Value)
		*slot = captured
		closure.FreeVars[i] = slot
	}
}

func (v *vm) discardFreeVars(numFree int) {
	for i := 0; i < numFree; i++ {
		v.pop()
	}
}
