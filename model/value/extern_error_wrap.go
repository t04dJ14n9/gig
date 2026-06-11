package value

// ErrorValue extracts a Go error from a value.Value.
// If the value is already a native Go error, returns it directly.
// If the value is an interpreter-synthesized struct with an Error() method,
// returns a gigStructWrapper that implements the error interface.
// Otherwise returns nil.
//
// This is the boundary function for generated DirectCall wrappers; use it
// whenever extracting an error-typed parameter from args.
func ErrorValue(v Value) error {
	iface := v.Interface()
	if iface == nil {
		return nil
	}
	// If it's already a Go error (e.g., from fmt.Errorf, errors.New), return as-is.
	if e, ok := iface.(error); ok {
		return e
	}
	typeName := isGigStruct(iface)
	if typeName == "" {
		return nil
	}
	// Check if the interpreted type has an Error() method.
	errorerFunc, hasError := resolveErrorer(v)
	if !hasError {
		return nil
	}
	// Also capture Stringer/GoStringer lazily so the wrapper is complete.
	stringerFunc, _ := resolveStringer(v)
	gostringerFunc, _ := resolveGoStringer(v)
	constFn := func(s string) func() string { return func() string { return s } }
	errFn := func() string { return errorerFunc() }

	var lazyStringer, lazyErrorer, lazyGoStringer func() (func() string, bool)
	if stringerFunc != nil {
		sf := constFn(stringerFunc())
		lazyStringer = func() (func() string, bool) { return sf, true }
	}
	lazyErrorer = func() (func() string, bool) { return errFn, true }
	if gostringerFunc != nil {
		gf := constFn(gostringerFunc())
		lazyGoStringer = func() (func() string, bool) { return gf, true }
	}
	return &gigStructWrapper{
		iface:          iface,
		typeName:       typeName,
		lazyStringer:   lazyStringer,
		lazyErrorer:    lazyErrorer,
		lazyGoStringer: lazyGoStringer,
	}
}

// ErrorWrap prepares a value.Value for use as a Go error.
// If the value is an interpreter-synthesized struct with an Error() method,
// returns a wrapper that implements the error interface. Otherwise returns
// the raw interface{} value.
//
// Deprecated: Use ErrorValue instead for typed error extraction.
func ErrorWrap(v Value) any {
	if e := ErrorValue(v); e != nil {
		return e
	}
	return v.Interface()
}
