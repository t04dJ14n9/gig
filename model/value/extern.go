// extern.go documents value-package helpers for crossing the interpreter/native boundary.
// Concrete behavior is split by domain:
//   - extern_format.go: fmt boundary entry point and wrap decisions
//   - extern_wrapper.go: gig struct Stringer/Formatter/error wrapper
//   - extern_collection_format.go: slice/map formatting for gig values
//   - extern_method.go: lazy String/Error/GoString method resolution
//   - extern_sprintf.go: %T-aware Sprintf replacement
//   - extern_type.go: gig reflect type-tag detection
//   - extern_errors.go: errors.Is/As/Unwrap compatibility
package value
