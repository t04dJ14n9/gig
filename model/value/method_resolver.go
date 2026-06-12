package value

import "sync"

// MethodResolverFunc is a callback for calling compiled methods on interpreted types.
// It receives a method name and receiver value, and returns the result if found.
type MethodResolverFunc func(methodName string, receiver Value) (Value, bool)

// MethodWithArgsResolverFunc is like MethodResolverFunc but supports passing arguments.
// Used for methods like Is(error) bool that require parameters beyond the receiver.
type MethodWithArgsResolverFunc func(methodName string, receiver Value, args []Value) (Value, bool)

// methodResolverRegistry is a thread-safe global registry of per-program method resolvers.
// This allows fmt DirectCall wrappers (which lack VM context) to resolve compiled methods
// on interpreted types. Each program registers its resolver on creation and unregisters
// on cleanup. Using sync.Map eliminates the data race that the old single-global approach had.
var methodResolverRegistry sync.Map // map[uintptr]MethodResolverFunc

// methodWithArgsResolverRegistry stores resolvers that support argument passing.
var methodWithArgsResolverRegistry sync.Map // map[uintptr]MethodWithArgsResolverFunc

// RegisterMethodResolver registers a method resolver for a program identified by key.
// The key should be a unique identifier per program (e.g., uintptr of program pointer).
func RegisterMethodResolver(key uintptr, resolver MethodResolverFunc) {
	methodResolverRegistry.Store(key, resolver)
}

// RegisterMethodWithArgsResolver registers a method-with-args resolver for a program.
func RegisterMethodWithArgsResolver(key uintptr, resolver MethodWithArgsResolverFunc) {
	methodWithArgsResolverRegistry.Store(key, resolver)
}

// UnregisterMethodResolver removes a method resolver for the given program key.
func UnregisterMethodResolver(key uintptr) {
	methodResolverRegistry.Delete(key)
	methodWithArgsResolverRegistry.Delete(key)
}

// callMethod attempts to call a compiled method on the receiver using the given resolver.
// If resolver is nil, it falls back to searching all registered per-program resolvers.
// Returns (result, true) if the method was found and called, or (zero, false) otherwise.
func callMethod(resolver MethodResolverFunc, methodName string, receiver Value) (Value, bool) {
	if resolver != nil {
		return resolver(methodName, receiver)
	}
	// Fallback: try all registered resolvers (for fmt DirectCall wrappers lacking VM context)
	var result Value
	var found bool
	methodResolverRegistry.Range(func(_, v any) bool {
		if r, ok := v.(MethodResolverFunc); ok {
			result, found = r(methodName, receiver)
			if found {
				return false // stop iteration
			}
		}
		return true // continue
	})
	if found {
		return result, true
	}
	return MakeNil(), false
}

// callMethodWithArgs calls a compiled method with additional arguments beyond the receiver.
// Used for methods like Is(error) bool that require parameters.
func callMethodWithArgs(methodName string, receiver Value, args []Value) (Value, bool) {
	var result Value
	var found bool
	methodWithArgsResolverRegistry.Range(func(_, v any) bool {
		if r, ok := v.(MethodWithArgsResolverFunc); ok {
			result, found = r(methodName, receiver, args)
			if found {
				return false
			}
		}
		return true
	})
	if found {
		return result, true
	}
	return MakeNil(), false
}
