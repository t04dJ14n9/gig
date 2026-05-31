// VM call dispatch is split by call domain.
//
//   - call_runtime.go: call frames, cancellation checks, closure conversion, and variadic unpacking.
//   - call_external.go: native function dispatch, DirectCall, and reflect.Call.
//   - call_boundary.go: third-party boundary validation for interpreter-defined values.
//   - call_external_method.go: native method dispatch and reflect method lookup.
//   - call_compiled_method.go: compiled Gig method fallback and receiver selection.
package vm
