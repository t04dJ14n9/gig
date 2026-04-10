package compiler

import (
	"reflect"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
)

// PackageLookup provides read-only access to registered external packages.
// This is a consumer-defined interface: the compiler declares only the methods
// it actually uses, keeping it decoupled from the importer package.
//
// It embeds bytecode.TypeResolver because the compiler assigns the lookup
// as the program's TypeResolver (used by the VM at runtime).
//
// The importer.Registry satisfies this interface via Go's implicit matching.
type PackageLookup interface {
	bytecode.TypeResolver

	// LookupExternalFunc resolves an external function by package path and name.
	// Returns the function value, an optional DirectCall wrapper, and whether it was found.
	LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool)

	// LookupMethodDirectCall resolves a method DirectCall wrapper by type and method name.
	// Returns the DirectCall wrapper and whether it was found.
	LookupMethodDirectCall(typeName, methodName string) (directCall func([]value.Value) value.Value, ok bool)

	// LookupExternalVar resolves an external variable by package path and name.
	// Returns a pointer to the variable and whether it was found.
	LookupExternalVar(pkgPath, varName string) (ptr any, ok bool)

	// LookupExternalTypeByName resolves an external type by package path and type name.
	// Returns the reflect.Type and whether it was found.
	LookupExternalTypeByName(pkgPath, typeName string) (reflect.Type, bool)
}
