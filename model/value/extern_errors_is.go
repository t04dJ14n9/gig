package value

import "reflect"

// GigErrorsIs implements errors.Is semantics for interpreter-defined types.
// It replicates the standard library algorithm but uses gig's method resolution
// to invoke custom Is(error) bool and Unwrap() error methods on gig types.
func GigErrorsIs(errVal Value, targetVal Value) bool {
	state := newGigErrorsIsState(errVal, targetVal)
	if result, handled := state.nilResult(); handled {
		return result
	}

	for {
		if state.matchesTarget() {
			return true
		}
		if state.matchesCustomIs() {
			return true
		}
		if matched, ok := state.advanceUnwrap(); matched || !ok {
			return matched
		}
	}
}

type gigErrorsIsState struct {
	errVal    Value
	targetVal Value
	err       error
	target    error
}

func newGigErrorsIsState(errVal Value, targetVal Value) gigErrorsIsState {
	return gigErrorsIsState{
		errVal:    errVal,
		targetVal: targetVal,
		err:       ErrorValue(errVal),
		target:    ErrorValue(targetVal),
	}
}

func (s gigErrorsIsState) nilResult() (bool, bool) {
	if s.err == nil && s.target == nil {
		return true, true
	}
	if s.err == nil || s.target == nil {
		return s.err == s.target, true
	}
	return false, false
}

func (s gigErrorsIsState) matchesTarget() bool {
	// Direct comparison also handles comparable native errors.
	return s.err == s.target || gigErrorsEqual(s.err, s.target)
}

func (s gigErrorsIsState) matchesCustomIs() bool {
	if _, ok := s.err.(*gigStructWrapper); ok {
		return gigCustomIs(s.errVal, s.targetVal)
	}
	if x, ok := s.err.(interface{ Is(error) bool }); ok {
		return x.Is(s.target)
	}
	return false
}

func gigCustomIs(errVal Value, targetVal Value) bool {
	result, found := callMethodWithArgs("Is", errVal, []Value{targetVal})
	return found && result.Kind() == KindBool && result.Bool()
}

func (s *gigErrorsIsState) advanceUnwrap() (bool, bool) {
	if _, ok := s.err.(*gigStructWrapper); ok {
		return false, s.advanceGigUnwrap()
	}
	if x, ok := s.err.(interface{ Unwrap() []error }); ok {
		return s.matchesJoinedErrors(x.Unwrap()), false
	}
	if x, ok := s.err.(interface{ Unwrap() error }); ok {
		return false, s.advanceNativeUnwrap(x.Unwrap())
	}
	return false, false
}

func (s *gigErrorsIsState) advanceGigUnwrap() bool {
	unwrapResult, found := callMethod(nil, "Unwrap", s.errVal)
	if !found {
		return false
	}
	unwrapped := ErrorValue(unwrapResult)
	if unwrapped == nil {
		return false
	}
	s.err = unwrapped
	s.errVal = unwrapResult
	return true
}

func (s gigErrorsIsState) matchesJoinedErrors(errs []error) bool {
	for _, err := range errs {
		if err != nil && GigErrorsIs(FromInterface(err), s.targetVal) {
			return true
		}
	}
	return false
}

func (s *gigErrorsIsState) advanceNativeUnwrap(unwrapped error) bool {
	if unwrapped == nil {
		return false
	}
	s.err = unwrapped
	s.errVal = FromInterface(unwrapped)
	return true
}

// gigErrorsEqual compares two errors for equality, handling gigStructWrapper.
// Two gigStructWrappers are equal if they wrap the same type and underlying value.
func gigErrorsEqual(a, b error) bool {
	wa, aIsGig := a.(*gigStructWrapper)
	wb, bIsGig := b.(*gigStructWrapper)
	if aIsGig && bIsGig {
		return wa.typeName == wb.typeName && reflect.DeepEqual(wa.iface, wb.iface)
	}
	return false
}
