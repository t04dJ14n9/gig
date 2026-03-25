package vm

import (
	"reflect"
	"sync"

	"git.woa.com/youngjin/gig/model/value"
)

// externalCallCache is a shared cache for external function lookups.
// It is shared between a parent VM and all its child goroutine VMs.
type externalCallCache struct {
	mu    sync.RWMutex
	cache []*extCallCacheEntry
}

// extCallCacheEntry caches resolved external function info for fast dispatch.
// This avoids repeated reflection lookups for external function calls.
type extCallCacheEntry struct {
	// fn is the reflect.Value of the function.
	fn reflect.Value

	// fnType is the function's type.
	fnType reflect.Type

	directCall func(args []value.Value) value.Value

	// isVariadic indicates if the function takes variadic arguments.
	isVariadic bool

	// numIn is the number of declared parameters.
	numIn int
}
