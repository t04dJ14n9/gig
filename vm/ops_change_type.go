package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeChangeType(frame *Frame) {
	typeIdx := frame.readUint16()
	srcLocalIdx := frame.readUint16()
	targetType := v.program.Types[typeIdx]
	val := v.pop()

	// Named-type conversion (e.g., []int -> sort.IntSlice).
	if named, ok := targetType.(*types.Named); ok {
		targetRT := typeToReflect(named, v.program)
		if targetRT != nil {
			// Get a reflect.Value, using the target type so ToReflectValue
			// handles element-type conversion (e.g., []int64 -> []int).
			rv := val.ToReflectValue(targetRT)
			if rv.IsValid() {
				// If ToReflectValue returned the underlying type, Convert to the named type.
				if rv.Type() != targetRT && rv.Type().ConvertibleTo(targetRT) {
					rv = rv.Convert(targetRT)
				}
				// For slices: update the source local to share the same backing array.
				// This ensures that sort.IntSlice(s) and s refer to the same data,
				// matching Go's semantics where ChangeType on slices shares memory.
				if srcLocalIdx != noSourceLocalSentinel && rv.Kind() == reflect.Slice {
					if int(srcLocalIdx) < len(frame.locals) {
						// Create a view of the same backing array as the underlying slice type.
						// e.g., for sort.IntSlice -> create a []int sharing the same backing.
						underlyingRV := rv.Convert(reflect.SliceOf(rv.Type().Elem()))
						frame.locals[srcLocalIdx] = value.MakeFromReflect(underlyingRV)
					}
				}
				v.push(value.MakeFromReflect(rv))
			} else {
				v.push(val)
			}
		} else {
			// Named type not in external registry (interpreted type) — pass through.
			v.push(val)
		}
	} else {
		// Not a named type, fall back to simple pass-through.
		v.push(val)
	}
}
