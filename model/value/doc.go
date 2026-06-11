// Package value implements a tagged-union Value system for high-performance interpretation.
//
// The Value type is the fundamental data unit in the Gig interpreter. It uses a tagged-union
// design that stores primitive types (bool, int, uint, float) directly in the num field,
// avoiding allocation and reflection overhead for common operations.
//
// # Design Philosophy
//
// The Value type is designed for:
//   - Zero allocation for primitive values
//   - Fast arithmetic and comparison operations without reflection
//   - Seamless interop with Go's reflect package for complex types
//   - Type safety through explicit kind checking
//
// # Memory Layout
//
// The Value struct is 32 bytes on 64-bit systems:
//   - kind: 1 byte (type tag)
//   - size: 1 byte (original Go type bit-width for numeric kinds) + 6 bytes padding
//   - num: 8 bytes (stores bool, int, uint bits, float bits)
//   - obj: 16 bytes (interface for string, complex, reflect.Value, composite types)
//
// The size field records the original Go type (e.g. int8 vs int32 vs int64) so that
// Interface() can return the exact Go type declared in the user's source code.
// It lives in the padding gap between kind and num, so it adds zero extra memory.
//
// Primitives (int, float, bool, uint, nil) are stored entirely in kind+size+num with obj=nil,
// so they never cause GC pressure.
//
// # Kind Types
//
//   - KindNil: null value
//   - KindBool: boolean (stored in num)
//   - KindInt: signed integers (stored in num)
//   - KindUint: unsigned integers (stored in num as bits)
//   - KindFloat: floating point (stored in num as float64 bits)
//   - KindString: string (stored in obj)
//   - KindComplex: complex number (stored in obj as complex128)
//   - KindPointer, KindSlice, KindArray, KindMap, KindChan, KindFunc, KindStruct, KindInterface:
//     stored in obj as reflect.Value or native Go value
//   - KindReflect: fallback for types not directly supported
//   - KindExternal: raw host Go pointer stored without eager reflect.Value boxing
//
// # Example Usage
//
//	// Create values
//	i := value.MakeInt(42)
//	s := value.MakeString("hello")
//	f := value.MakeFloat(3.14)
//
//	// Arithmetic
//	sum := i.Add(value.MakeInt(8)) // sum.Int() == 50
//
//	// Comparison
//	if i.Cmp(value.MakeInt(40)) > 0 {
//	    fmt.Println("42 > 40")
//	}
//
//	// Convert to/from interface{}
//	v := value.FromInterface(myStruct)
//	obj := v.Interface()
package value
