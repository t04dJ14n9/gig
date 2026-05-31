package value

import "reflect"

// GigErrorsIs implements errors.Is semantics for interpreter-defined types.
// It replicates the standard library algorithm but uses gig's method resolution
// to invoke custom Is(error) bool and Unwrap() error methods on gig types.
func GigErrorsIs(errVal Value, targetVal Value) bool {
	err := ErrorValue(errVal)
	target := ErrorValue(targetVal)
	if err == nil && target == nil {
		return true
	}
	if err == nil || target == nil {
		return err == target
	}

	for {
		// Direct comparison also handles gigStructWrapper by comparing underlying values.
		if err == target {
			return true
		}
		if gigErrorsEqual(err, target) {
			return true
		}

		// Check custom Is() method on gig types.
		if _, ok := err.(*gigStructWrapper); ok {
			if result, found := callMethodWithArgs("Is", errVal, []Value{targetVal}); found {
				if result.Kind() == KindBool && result.Bool() {
					return true
				}
			}
		} else if x, ok := err.(interface{ Is(error) bool }); ok {
			if x.Is(target) {
				return true
			}
		}

		// Unwrap.
		if _, ok := err.(*gigStructWrapper); ok {
			unwrapResult, found := callMethod(nil, "Unwrap", errVal)
			if !found {
				return false
			}
			unwrapped := ErrorValue(unwrapResult)
			if unwrapped == nil {
				return false
			}
			err = unwrapped
			errVal = unwrapResult
		} else if x, ok := err.(interface{ Unwrap() []error }); ok {
			// Multi-unwrap (errors.Join): recursively check each wrapped error.
			for _, e := range x.Unwrap() {
				if e != nil && GigErrorsIs(FromInterface(e), targetVal) {
					return true
				}
			}
			return false
		} else if x, ok := err.(interface{ Unwrap() error }); ok {
			err = x.Unwrap()
			if err == nil {
				return false
			}
			errVal = FromInterface(err)
		} else {
			return false
		}
	}
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
