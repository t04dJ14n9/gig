package value

// GigErrorsUnwrap implements errors.Unwrap for interpreter-defined types.
func GigErrorsUnwrap(errVal Value) Value {
	err := ErrorValue(errVal)
	if err == nil {
		return MakeNil()
	}

	if _, ok := asGigStructError(err); ok {
		result, found := callMethod("Unwrap", errVal)
		if found {
			return result
		}
		return MakeNil()
	}

	if x, ok := err.(interface{ Unwrap() error }); ok {
		unwrapped := x.Unwrap()
		if unwrapped == nil {
			return MakeNil()
		}
		return FromInterface(unwrapped)
	}
	return MakeNil()
}
