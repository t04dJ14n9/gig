// external.go implements ExternalPackage methods for adding functions, variables,
// constants, types, and method DirectCall wrappers to a registered package.
package importer

import (
	"reflect"

	"git.woa.com/youngjin/gig/model/external"
	"git.woa.com/youngjin/gig/model/value"
)

// AddFunction adds a function to the package.
func (p *ExternalPackage) AddFunction(name string, fn any, doc string, directCall func([]value.Value) value.Value) {
	sig := funcSignature(fn)
	p.Objects[name] = &external.ExternalObject{
		Name:       name,
		Kind:       external.ObjectKindFunction,
		Value:      fn,
		Type:       sig,
		Doc:        doc,
		DirectCall: directCall,
	}
}

// AddVariable adds a variable to the package.
func (p *ExternalPackage) AddVariable(name string, ptr any, doc string) {
	typ := typeOf(reflect.TypeOf(ptr).Elem())
	p.Objects[name] = &external.ExternalObject{
		Name:  name,
		Kind:  external.ObjectKindVariable,
		Value: ptr,
		Type:  typ,
		Doc:   doc,
	}
}

// AddConstant adds a constant to the package.
func (p *ExternalPackage) AddConstant(name string, val any, doc string) {
	typ := typeOf(reflect.TypeOf(val))
	p.Objects[name] = &external.ExternalObject{
		Name:  name,
		Kind:  external.ObjectKindConstant,
		Value: val,
		Type:  typ,
		Doc:   doc,
	}
}

// AddType adds a named type to the package.
func (p *ExternalPackage) AddType(name string, typ reflect.Type, doc string) {
	if typ == nil {
		return
	}
	p.Types[name] = typ
	p.Objects[name] = &external.ExternalObject{
		Name:  name,
		Kind:  external.ObjectKindType,
		Value: reflect.Zero(typ).Interface(),
		Type:  typeOf(typ),
		Doc:   doc,
	}
}

// AddMethodDirectCall registers a DirectCall wrapper for a method on a type in this package.
// It uses the package's owning registry instance.
func (p *ExternalPackage) AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value) {
	if p.registry != nil {
		p.registry.AddMethodDirectCall(p.Path+"."+typeName, methodName, dc)
	}
}
