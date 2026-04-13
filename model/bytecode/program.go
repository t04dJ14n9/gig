package bytecode

import (
	"go/types"
	"reflect"
	"sync"

	"github.com/t04dJ14n9/gig/model/value"
)

// CompiledProgram represents a compiled program ready for execution.
// It contains all compiled functions, constants, types, and global variables.
type CompiledProgram struct {
	// Functions maps function names to their compiled bytecode.
	Functions map[string]*CompiledFunction

	// FuncByIndex provides O(1) function lookup by index.
	// Populated at compile time so the VM can skip the FuncIndex map scan.
	FuncByIndex []*CompiledFunction

	// MethodsByName maps method name → list of compiled functions with that name.
	// Built once after compilation. Allows O(k) dispatch (k = methods with that name)
	// instead of O(n) scan over FuncByIndex.
	MethodsByName map[string][]*CompiledFunction

	// Constants is the constant pool for literal values and external references.
	Constants []any

	// PrebakedConstants is the pre-converted constant pool.
	// Built once at startup to avoid per-OpConst value.FromInterface() calls.
	PrebakedConstants []value.Value

	// IntConstants is an int-specialized constant pool.
	// For each constant that is an int64, IntConstants[i] holds the value.
	// Used by OpInt* superinstructions for zero-overhead constant access.
	IntConstants []int64

	// Globals maps global variable names to their indices.
	Globals map[string]int

	// GlobalZeroValues maps global variable index to its zero reflect.Value.
	// Used by the VM to initialize zero-valued globals to their proper zero value
	// instead of leaving them as nil value.Value{}. This enables pointer-receiver
	// method calls on value-type globals like sync.Mutex, sync.WaitGroup, etc.
	// Only populated for globals whose zero value is a non-nil struct/map/slice/chan.
	GlobalZeroValues map[int]reflect.Value

	// Types is the type pool for runtime type operations.
	Types []types.Type

	// ExternalVarValues stores external package variable values indexed by global index.
	// These are resolved at compile time and used to initialize globals in the VM.
	// The value is a pointer to the external variable (e.g., &time.UTC).
	ExternalVarValues map[int]any

	// TypeResolver resolves external types at runtime.
	// Used by the VM's typeToReflect to look up real reflect.Type for named types.
	TypeResolver TypeResolver

	// ReflectTypeCache caches types.Type → reflect.Type conversions at the
	// program level. This prevents reflect.StructOf from returning different
	// reflect.Type objects for the same types.Type across multiple VM executions,
	// which would cause "reflect.Set: value not assignable" panics.
	// Key: types.Type, Value: reflect.Type.
	ReflectTypeCache sync.Map

	// ReflectTypeNames maps reflect.Type → TypeName for interpreter-synthesized
	// struct types. Used by method dispatch to identify which named type a
	// reflect.Value belongs to, replacing the old _gig_id phantom field approach.
	// Key: reflect.Type, Value: string (bare type name, e.g. "Foo").
	ReflectTypeNames sync.Map
}

// TypeResolver resolves external types at runtime.
// This interface breaks the circular dependency between bytecode/ and importer/
// while providing type-safe access to external type lookup.
type TypeResolver interface {
	LookupExternalType(t types.Type) (reflect.Type, bool)
}

// CachedReflectType looks up a cached reflect.Type for the given types.Type.
// Returns the cached type and true, or nil and false if not cached.
func (p *CompiledProgram) CachedReflectType(t types.Type) (reflect.Type, bool) {
	if v, ok := p.ReflectTypeCache.Load(t); ok {
		return v.(reflect.Type), true
	}
	return nil, false
}

// CacheReflectType stores a types.Type → reflect.Type mapping.
// Uses LoadOrStore to handle concurrent writes safely.
func (p *CompiledProgram) CacheReflectType(t types.Type, rt reflect.Type) reflect.Type {
	actual, _ := p.ReflectTypeCache.LoadOrStore(t, rt)
	return actual.(reflect.Type)
}

// RegisterTypeName associates a reflect.Type with its bare type name (e.g. "Foo").
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
