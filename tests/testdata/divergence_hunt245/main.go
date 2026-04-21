package divergence_hunt245

import (
	"fmt"
)

// ============================================================================
// Round 245: Nested panic/recover
// ============================================================================

// NestedPanicRecoverSimple tests simple nested panic/recover
func NestedPanicRecoverSimple() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("outer:%v", r)
			}
		}()
		func() {
			defer func() {
				if r := recover(); r != nil {
					result = fmt.Sprintf("inner:%v", r)
				}
			}()
			panic("inner panic")
		}()
	}()
	return result
}

// NestedPanicRecoverBothPanic tests both inner and outer panic
func NestedPanicRecoverBothPanic() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("outer:%v", r)
			}
		}()
		func() {
			defer func() {
				if r := recover(); r != nil {
					result = fmt.Sprintf("inner:%v", r)
					panic("outer panic")
				}
			}()
			panic("inner panic")
		}()
	}()
	return result
}

// TripleNestedPanicRecover tests triple nesting
func TripleNestedPanicRecover() string {
	results := []string{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				results = append(results, fmt.Sprintf("L1:%v", r))
			}
		}()
		func() {
			defer func() {
				if r := recover(); r != nil {
					results = append(results, fmt.Sprintf("L2:%v", r))
				}
			}()
			func() {
				defer func() {
					if r := recover(); r != nil {
						results = append(results, fmt.Sprintf("L3:%v", r))
					}
				}()
				panic("level3")
			}()
		}()
	}()
	return fmt.Sprintf("%v", results)
}

// NestedPanicRecoverPartial tests partial recovery in nested functions
func NestedPanicRecoverPartial() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("caught:%v", r)
			}
		}()
		func() {
			// inner function does NOT recover
			panic("not caught inner")
		}()
	}()
	return result
}

// NestedDeferPanicRecover tests nested defers with panic
func NestedDeferPanicRecover() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("outer:%v", r)
			}
		}()
		defer func() {
			panic("defer panic")
		}()
		panic("original")
	}()
	return result
}

// NestedPanicRecoverLoop tests panic/recover in nested loop functions
func NestedPanicRecoverLoop() string {
	results := []string{}
	for i := 0; i < 3; i++ {
		func(n int) {
			defer func() {
				if r := recover(); r != nil {
					results = append(results, fmt.Sprintf("%d:%v", n, r))
				}
			}()
			if n == 1 {
				panic("one")
			}
		}(i)
	}
	return fmt.Sprintf("%v", results)
}

// DeepNestedRecoverChain tests deep nesting with recover chain
func DeepNestedRecoverChain() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("L1:%v", r)
			}
		}()
		func() {
			defer func() {
				if r := recover(); r != nil {
					panic(fmt.Sprintf("L2:%v", r))
				}
			}()
			func() {
				defer func() {
					if r := recover(); r != nil {
						panic(fmt.Sprintf("L3:%v", r))
					}
				}()
				panic("start")
			}()
		}()
	}()
	return result
}

// NestedPanicDifferentTypes tests panic with different types at different levels
func NestedPanicDifferentTypes() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("outer:%T:%v", r, r)
			}
		}()
		func() {
			defer func() {
				if r := recover(); r != nil {
					panic(42)
				}
			}()
			panic("string panic")
		}()
	}()
	return result
}

// NestedPanicNamedReturn tests nested panic/recover with named return
func NestedPanicNamedReturn() string {
	result := func() (r string) {
		defer func() {
			if rec := recover(); rec != nil {
				r = fmt.Sprintf("outer:%v", rec)
			}
		}()
		func() (inner string) {
			defer func() {
				if rec := recover(); rec != nil {
					inner = fmt.Sprintf("inner:%v", rec)
				}
			}()
			panic("nested")
		}()
		return "no panic"
	}()
	return result
}
