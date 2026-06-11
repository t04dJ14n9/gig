package compiler

import (
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// compileGlobalValue compiles an SSA Global and records the runtime zero-value metadata
// needed to initialize the VM's global slots.
func (c *compiler) compileGlobalValue(val *ssa.Global) {
	globalName, isExternal := globalBindingName(val)
	globalIdx, ok := c.globals[globalName]
	if !ok {
		globalIdx = len(c.globals)
		c.globals[globalName] = globalIdx

		c.recordGlobalZeroValue(globalIdx, val.Type())
		if isExternal && c.lookup != nil {
			c.recordExternalGlobalValue(globalIdx, val)
		}
	}
	c.emit(bytecode.OpGlobal, uint16(globalIdx))
}

func globalBindingName(val *ssa.Global) (string, bool) {
	globalName := val.Name()
	isExternal := val.Pkg != nil && val.Pkg.Pkg != nil && val.Pkg.Pkg.Path() != "main"
	if isExternal {
		globalName = val.Pkg.Pkg.Path() + "." + globalName
	}
	return globalName, isExternal
}

func (c *compiler) recordGlobalZeroValue(globalIdx int, globalType types.Type) {
	ptrType, ok := globalType.(*types.Pointer)
	if !ok {
		return
	}

	elemType := ptrType.Elem()
	switch t := elemType.(type) {
	case *types.Named:
		c.recordNamedGlobalZeroValue(globalIdx, elemType, t)
	case *types.Basic:
		if rt, ok := bytecode.BasicKindToReflectType[t.Kind()]; ok {
			c.globalZeroValues[globalIdx] = reflect.Zero(rt)
		}
	case *types.Struct:
		c.globalElemTypes[globalIdx] = elemType
	case *types.Array:
		c.globalElemTypes[globalIdx] = elemType
	}
}

func (c *compiler) recordNamedGlobalZeroValue(globalIdx int, elemType types.Type, named *types.Named) {
	obj := named.Obj()
	if obj != nil && obj.Pkg() != nil && c.lookup != nil {
		pkgPath := obj.Pkg().Path()
		typeName := obj.Name()
		if rt, found := c.lookup.LookupExternalTypeByName(pkgPath, typeName); found && rt.Kind() == reflect.Struct {
			c.globalZeroValues[globalIdx] = reflect.New(rt)
		}
	}

	if _, exists := c.globalZeroValues[globalIdx]; exists {
		return
	}
	switch named.Underlying().(type) {
	case *types.Struct:
		c.globalElemTypes[globalIdx] = elemType
	case *types.Array:
		c.globalElemTypes[globalIdx] = elemType
	}
}

func (c *compiler) recordExternalGlobalValue(globalIdx int, val *ssa.Global) {
	ptr, found := c.lookup.LookupExternalVar(val.Pkg.Pkg.Path(), val.Name())
	if !found {
		return
	}

	rv := reflect.ValueOf(ptr)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		c.externalVarValues[globalIdx] = rv.Elem().Interface()
	}
}
