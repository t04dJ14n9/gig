package vm

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
)

// ResolveCompiledMethod finds a compiled method in the program's function table
// and executes it with the given receiver.
func ResolveCompiledMethod(program *bytecode.CompiledProgram, methodName string, receiver value.Value) (value.Value, bool) {
	rv, ok := receiver.ReflectValue()
	if !ok {
		return value.MakeNil(), false
	}
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}

	// Extract the type name from the _gig_id field for matching
	receiverTypeName := ""
	if rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			sf := rt.Field(i)
			if sf.Name == "_gig_id" {
				if idx := strings.LastIndex(sf.PkgPath, "#"); idx >= 0 {
					qualName := sf.PkgPath[idx+1:]
					if dotIdx := strings.LastIndex(qualName, "."); dotIdx >= 0 {
						receiverTypeName = qualName[dotIdx+1:]
					} else {
						receiverTypeName = qualName
					}
				}
				break
			}
		}
	}

	if receiverTypeName == "" {
		return value.MakeNil(), false
	}

	// Search the compiled function table for a method with matching name and receiver type
	for _, fn := range program.MethodsByName[methodName] {
		if fn.ReceiverTypeName != receiverTypeName {
			continue
		}
		// Found the method! Execute it with a temporary VM.
		tempVM := &vm{
			program: program,
			stack:   make([]value.Value, deferVMStackSize),
			sp:      0,
			frames:  make([]*Frame, initialFrameDepth),
			fp:      0,
			globals: make([]value.Value, len(program.Globals)),
			ctx:     context.Background(),
			extCallCache: &externalCallCache{
				cache: make([]*extCallCacheEntry, len(program.Constants)),
			},
		}
		// Note: tempVM does not have initialGlobals since resolveCompiledMethod
		// is called without a VM context. This is acceptable because method resolution
		// only needs to execute the method, not full program init.
		tempVM.callFunction(fn, []value.Value{receiver}, nil)
		result, err := tempVM.run()
		if err != nil {
			return value.MakeNil(), false
		}
		return result, true
	}
	return value.MakeNil(), false
}

// VMPool is a thread-safe pool of VMs for a given program.
type VMPool struct {
	mu    sync.Mutex
	vms   []*vm // available VMs
	newVM func() *vm
}

// NewVMPool creates a VM pool for the given program.
func NewVMPool(program *bytecode.CompiledProgram, initialGlobals []value.Value, goroutines *GoroutineTracker) *VMPool {
	return &VMPool{
		newVM: func() *vm {
			return newVM(program, initialGlobals, goroutines)
		},
	}
}

// Get returns an idle VM from the pool.
func (p *VMPool) Get() VM {
	p.mu.Lock()
	if len(p.vms) > 0 {
		v := p.vms[len(p.vms)-1]
		p.vms = p.vms[:len(p.vms)-1]
		p.mu.Unlock()
		return v
	}
	p.mu.Unlock()
	return p.newVM()
}

// Put returns a VM to the pool for reuse.
func (p *VMPool) Put(x VM) {
	x.Reset()
	p.mu.Lock()
	p.vms = append(p.vms, x.(*vm))
	p.mu.Unlock()
}
