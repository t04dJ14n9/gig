package external

import (
	"go/types"

	"github.com/t04dJ14n9/gig/model/value"
)

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

	DirectCall func([]value.Value) value.Value
}
