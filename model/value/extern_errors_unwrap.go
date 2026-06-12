package value

// GigErrorsUnwrap implements errors.Unwrap for interpreter-defined types.
// If the error is a gig type with an Unwrap() method, invokes it via the
// compiled method resolver. Otherwise delegates to standard errors.Unwrap.
func GigErrorsUnwrap(errVal Value) Value {
	err := ErrorValue(errVal)
	if err == nil {
		return MakeNil()
	}

	// For gig types, use compiled method resolution.
	if _, ok := err.(*gigStructWrapper); ok {
		result, found := callMethod(nil, "Unwrap", errVal)
		if found {
			return result
		}
		return MakeNil()
	}

	// For native Go errors, use standard unwrap.
	if x, ok := err.(interface{ Unwrap() error }); ok {
		unwrapped := x.Unwrap()
		if unwrapped == nil {
			return MakeNil()
		}
		return FromInterface(unwrapped)
	}
	return MakeNil()
}
