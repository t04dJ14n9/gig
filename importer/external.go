// external.go implements ExternalPackage methods for adding functions,
// variables, constants, and types to a registered host package.
package importer

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/external"
)

// AddFunction adds a function to the package.
func (p *ExternalPackage) AddFunction(name string, fn any, doc string) {
	sig := funcSignature(fn)
	p.Objects[name] = &external.ExternalObject{
		Name:  name,
		Kind:  external.ObjectKindFunction,
		Value: fn,
		Type:  sig,
		Doc:   doc,
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
