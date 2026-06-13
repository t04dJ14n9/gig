package external

import "go/types"

// ObjectKind represents the kind of external object (function, variable, constant, or type).
type ObjectKind int

const (
	// ObjectKindInvalid indicates an invalid or uninitialized object.
	ObjectKindInvalid ObjectKind = iota

	// ObjectKindFunction indicates a function object.
	ObjectKindFunction

	// ObjectKindVariable indicates a mutable variable object.
	ObjectKindVariable

	// ObjectKindConstant indicates an immutable constant object.
	ObjectKindConstant

	// ObjectKindType indicates a named type object.
	ObjectKindType
)

// ExternalObject represents a function, variable, constant, or type from an external package.
// It stores the Go value and type information needed for the interpreter.
//
// In the legacy bytecode VM, this struct also carried a DirectCall
// wrapper that converted between the legacy model/value.Value type and
// the host signature. The v2 SSA interpreter dispatches via reflect
// directly, so DirectCall is no longer needed.
type ExternalObject struct {
	// Name is the object's identifier (e.g., "Sprintf", "NoError").
	Name string

	// Kind indicates whether this is a function, variable, constant, or type.
	Kind ObjectKind

	// Value is the Go value:
	//   - Function: the function value
	//   - Variable: pointer to the variable
	//   - Constant: the constant value
	//   - Type: zero value of the type
	Value any

	// Type is the Go types.Type representation.
	Type types.Type

	// Doc is optional documentation text.
	Doc string
}
