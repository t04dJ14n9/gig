package value

// Kind represents the type of a Value.
type Kind uint8

const (
	KindInvalid Kind = iota
	KindNil
	KindBool
	KindInt     // int, int8, int16, int32, int64
	KindUint    // uint, uint8, uint16, uint32, uint64
	KindFloat   // float32, float64
	KindString  // string
	KindComplex // complex64, complex128
	KindPointer // *T
	KindSlice   // []T
	KindArray   // [N]T
	KindMap     // map[K]V
	KindChan    // chan T
	KindFunc    // func
	KindStruct  // struct{}
	KindInterface
	KindReflect // fallback to reflect.Value
	KindBytes   // []byte stored natively (zero reflection)
	KindExternal
)

// kindNameTable maps Kind values to their string representations.
var kindNameTable = [256]string{
	KindArray:     "array",
	KindBool:      "bool",
	KindBytes:     "bytes",
	KindChan:      "chan",
	KindComplex:   "complex",
	KindFloat:     "float",
	KindFunc:      "func",
	KindInt:       "int",
	KindInterface: "interface",
	KindInvalid:   "invalid",
	KindMap:       "map",
	KindNil:       "nil",
	KindPointer:   "pointer",
	KindReflect:   "reflect",
	KindSlice:     "slice",
	KindString:    "string",
	KindStruct:    "struct",
	KindUint:      "uint",
	KindExternal:  "external",
}

// String returns the name of the kind.
func (k Kind) String() string {
	if name := kindNameTable[k]; name != "" {
		return name
	}
	return "unknown"
}
