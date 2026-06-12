// container.go holds shared helpers for collection operations.
package value

import (
	"reflect"
	"unsafe"
)

// UnsafeAddrOf returns the unsafe.Pointer for a reflect.Value that is addressable.
// This is used internally by the VM to obtain settable pointers to unexported struct fields.
func UnsafeAddrOf(v reflect.Value) unsafe.Pointer {
	return v.Addr().UnsafePointer()
}

// convertSliceToInterface converts a slice of concrete types to slice of interface{}.
// This is needed for self-referential structs where []*T becomes []interface{} due to cycle breaking.
func convertSliceToInterface(slice reflect.Value, targetType reflect.Type) reflect.Value {
	if slice.Kind() != reflect.Slice || targetType.Kind() != reflect.Slice {
		return slice
	}
	if targetType.Elem().Kind() != reflect.Interface {
		return slice
	}
	if slice.Type().Elem().Kind() == reflect.Interface {
		return slice // Already interface slice
	}

	// Create new slice with target type
	newSlice := reflect.MakeSlice(targetType, slice.Len(), slice.Cap())
	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)
		// Wrap each element in interface{}
		// For pointers, we need to convert to interface{} explicitly
		if elem.CanInterface() {
			newSlice.Index(i).Set(reflect.ValueOf(elem.Interface()))
		} else {
			// Unexported field - this shouldn't happen but handle gracefully
			newSlice.Index(i).Set(elem)
		}
	}
	return newSlice
}
