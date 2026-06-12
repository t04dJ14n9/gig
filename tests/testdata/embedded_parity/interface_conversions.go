package main

import "fmt"

type ifaceAlias interface {
	String() string
}

type namedString string

func (s namedString) String() string { return string(s) }

type errNode struct{ msg string }

func (e *errNode) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.msg
}

// Result compares interface assertions/conversions and typed nil behavior.
func Result() string {
	var e error = (*errNode)(nil)
	typedNil := "not-nil"
	if e == nil {
		typedNil = "nil"
	}

	var source any = namedString("typed")
	asserted, ok := source.(ifaceAlias)
	asInt := 99
	if intVal, okInt := any(asInt).(int); okInt {
		asInt = intVal * 2
	}
	if !ok {
		return "assert-failed"
	}

	var anyStringer any = namedString("go")
	stringed := []byte(anyStringer.(ifaceAlias).String())
	if len(stringed) != 2 {
		return "bad-len"
	}

	return fmt.Sprintf("%s:%s:%t:%d:%d:%s", typedNil, asserted.String(), ok, asInt, len(stringed), string(stringed))
}
