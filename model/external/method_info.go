package external

import "github.com/t04dJ14n9/gig/model/value"

// ExternalMethodInfo contains method dispatch information.
// It is stored in the constant pool and used by OpMethodCall.
type ExternalMethodInfo struct {
	// MethodName is the name of the method to call.
	MethodName string

	// ReceiverTypeName is the fully qualified name of the receiver type
	// (e.g., "GetterImpl", "AdderStruct"). Used by callCompiledMethod
	// to disambiguate when multiple compiled methods share the same name.
	// Empty string means "match any receiver" (backward compatible).
	ReceiverTypeName string

	// DirectCall is an optional typed wrapper that avoids reflect.Call for this method.
	DirectCall func(args []value.Value) value.Value
}
