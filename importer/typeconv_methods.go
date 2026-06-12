package importer

import (
	"go/types"
	"reflect"
)

type methodImportContext struct {
	named         *types.Named
	isInterface   bool
	isPointerRecv bool
}

// addMethodsToNamed adds methods from a reflect.Type to a types.Named type.
// This allows the type checker to find methods on external types.
// Both value receiver and pointer receiver methods are added.
// For interface types, methods do NOT include a receiver parameter.
func addMethodsToNamed(named *types.Named, rt reflect.Type) {
	// Check if this is an interface type - interface methods don't have receiver params
	isInterface := rt.Kind() == reflect.Interface

	// Enumerate all exported methods on the value receiver
	addMethodsFromType(named, rt, isInterface, false)

	// For interface types, don't process pointer receiver methods
	if isInterface {
		return
	}

	// Also enumerate methods on the pointer receiver (*T)
	addMethodsFromType(named, reflect.PointerTo(rt), false, true)
}

// addMethodsFromType adds exported methods from methodSource to the named type.
// If skipReceiver is false (interface types), parameters start at index 0.
// If isPointerRecv is true, skips methods already present on named and uses pointer receiver.
func addMethodsFromType(named *types.Named, methodSource reflect.Type, isInterface, isPointerRecv bool) {
	ctx := methodImportContext{
		named:         named,
		isInterface:   isInterface,
		isPointerRecv: isPointerRecv,
	}
	for i := 0; i < methodSource.NumMethod(); i++ {
		method := methodSource.Method(i)
		if !ctx.shouldAddMethod(method) {
			continue
		}
		ctx.addMethod(method)
	}
}

func (ctx methodImportContext) shouldAddMethod(method reflect.Method) bool {
	if !method.IsExported() {
		return false
	}
	if ctx.isPointerRecv && ctx.hasMethod(method.Name) {
		return false
	}
	return ctx.methodTypeIsUsable(method.Type)
}

func (ctx methodImportContext) hasMethod(name string) bool {
	for j := 0; j < ctx.named.NumMethods(); j++ {
		if ctx.named.Method(j).Name() == name {
			return true
		}
	}
	return false
}

func (ctx methodImportContext) methodTypeIsUsable(methodType reflect.Type) bool {
	if methodType.Kind() != reflect.Func {
		return false
	}
	// Concrete reflected methods include receiver as input 0. Interface
	// methods do not, so a concrete method without input 0 is malformed here.
	return ctx.isInterface || methodType.NumIn() >= 1
}

func (ctx methodImportContext) addMethod(method reflect.Method) {
	sig := ctx.methodSignature(method.Type)
	fn := types.NewFunc(0, ctx.named.Obj().Pkg(), method.Name, sig)
	ctx.named.AddMethod(fn)
}

func (ctx methodImportContext) methodSignature(methodType reflect.Type) *types.Signature {
	return types.NewSignatureType(
		ctx.receiver(),
		nil, nil,
		types.NewTuple(ctx.params(methodType)...),
		types.NewTuple(reflectResults(methodType)...),
		methodType.IsVariadic(),
	)
}

func (ctx methodImportContext) receiver() *types.Var {
	// Interfaces do not encode a receiver parameter. Concrete value methods use
	// T, and pointer receiver methods use *T, matching reflect's method sets.
	if ctx.isInterface {
		return nil
	}
	recvType := types.Type(ctx.named)
	if ctx.isPointerRecv {
		recvType = types.NewPointer(ctx.named)
	}
	return types.NewVar(0, nil, "", recvType)
}

func (ctx methodImportContext) params(methodType reflect.Type) []*types.Var {
	return reflectInputs(methodType, ctx.paramStart())
}

func (ctx methodImportContext) paramStart() int {
	if ctx.isInterface {
		return 0
	}
	return 1
}

func reflectInputs(methodType reflect.Type, start int) []*types.Var {
	var params []*types.Var
	for j := start; j < methodType.NumIn(); j++ {
		params = append(params, reflectVar(methodType.In(j)))
	}
	return params
}

func reflectResults(methodType reflect.Type) []*types.Var {
	var results []*types.Var
	for j := 0; j < methodType.NumOut(); j++ {
		results = append(results, reflectVar(methodType.Out(j)))
	}
	return results
}

func reflectVar(rt reflect.Type) *types.Var {
	return types.NewVar(0, nil, "", convertReflectType(rt))
}
