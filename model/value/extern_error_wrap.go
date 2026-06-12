package value

// ErrorValue extracts a Go error from a Value. Interpreter-synthesized structs
// with an Error method are wrapped because reflect.StructOf types cannot carry
// Go methods.
func ErrorValue(v Value) error {
	iface := v.Interface()
	if iface == nil {
		return nil
	}
	if e, ok := iface.(error); ok {
		return e
	}
	typeName := isGigStruct(iface)
	if typeName == "" {
		return nil
	}
	errorerFunc, hasError := resolveErrorer(v)
	if !hasError {
		return nil
	}
	stringerFunc := resolveStringer(v)
	gostringerFunc := resolveGoStringer(v)
	return &gigStructError{
		iface:      iface,
		typeName:   typeName,
		stringer:   stringerFunc,
		errorer:    errorerFunc,
		gostringer: gostringerFunc,
	}
}

// ErrorWrap prepares a Value for use as a Go error, returning the raw interface
// value when the Value does not represent an interpreted error.
func ErrorWrap(v Value) any {
	if e := ErrorValue(v); e != nil {
		return e
	}
	return v.Interface()
}
