package external

import "github.com/t04dJ14n9/gig/model/value"

// ExternalMethodInfo contains method dispatch information.
// It is stored in the constant pool and used by OpCallExternal.
type ExternalMethodInfo struct {
	// PkgPath is the Go import path for the package that owns the receiver type.
	PkgPath string

	// MethodName is the name of the method to call.
	MethodName string

	// FuncName is the exported method name used for diagnostics.
	FuncName string

	// IsStdlib records whether PkgPath belongs to the Go standard library.
	// It lets the VM skip repeated string classification in method hot paths.
	IsStdlib bool

	// ReceiverTypeName is the fully qualified name of the receiver type
	// (e.g., "GetterImpl", "AdderStruct"). Used by callCompiledMethod
	// to disambiguate when multiple compiled methods share the same name.
	// Empty string means "match any receiver" (backward compatible).
	ReceiverTypeName string

	// DirectCall is an optional typed wrapper that avoids reflect.Call for this method.
	DirectCall func(args []value.Value) value.Value
}
