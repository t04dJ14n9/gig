package bytecode

import (
	"go/types"
	"reflect"
)

// TypeResolver resolves external types at runtime.
// This interface breaks the circular dependency between bytecode/ and importer/
// while providing type-safe access to external type lookup.
type TypeResolver interface {
	LookupExternalType(t types.Type) (reflect.Type, bool)
	LookupExternalTypeByName(pkgPath, typeName string) (reflect.Type, bool)
}

// CachedReflectType looks up a cached reflect.Type for the given types.Type.
// Returns the cached type and true, or nil and false if not cached.
func (p *CompiledProgram) CachedReflectType(t types.Type) (reflect.Type, bool) {
	if v, ok := p.ReflectTypeCache.Load(t); ok {
		return v.(reflect.Type), true
	}
	return nil, false
}

// CacheReflectType stores a types.Type to reflect.Type mapping.
// Uses LoadOrStore to handle concurrent writes safely.
func (p *CompiledProgram) CacheReflectType(t types.Type, rt reflect.Type) reflect.Type {
	actual, _ := p.ReflectTypeCache.LoadOrStore(t, rt)
	return actual.(reflect.Type)
}

// RegisterTypeName associates a reflect.Type with its bare type name.
// Used for method dispatch on interpreter-synthesized struct types.
func (p *CompiledProgram) RegisterTypeName(rt reflect.Type, name string) {
	p.ReflectTypeNames.Store(rt, name)
}

// LookupTypeName returns the bare type name for a reflect.Type, or "" if not registered.
func (p *CompiledProgram) LookupTypeName(rt reflect.Type) string {
	if v, ok := p.ReflectTypeNames.Load(rt); ok {
		return v.(string)
	}
	return ""
}
