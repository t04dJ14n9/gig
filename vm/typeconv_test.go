package vm

import (
	"go/types"
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// ---------------------------------------------------------------------------
// typeToReflect tests - these test the exported behavior indirectly
// ---------------------------------------------------------------------------

func TestTypeToReflectBasicTypes(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	tests := []struct {
		name string
		typ  types.Type
		want reflect.Kind
	}{
		{"bool", types.Typ[types.Bool], reflect.Bool},
		{"int", types.Typ[types.Int], reflect.Int},
		{"int8", types.Typ[types.Int8], reflect.Int8},
		{"int16", types.Typ[types.Int16], reflect.Int16},
		{"int32", types.Typ[types.Int32], reflect.Int32},
		{"int64", types.Typ[types.Int64], reflect.Int64},
		{"uint", types.Typ[types.Uint], reflect.Uint},
		{"uint8", types.Typ[types.Uint8], reflect.Uint8},
		{"uint16", types.Typ[types.Uint16], reflect.Uint16},
		{"uint32", types.Typ[types.Uint32], reflect.Uint32},
		{"uint64", types.Typ[types.Uint64], reflect.Uint64},
		{"uintptr", types.Typ[types.Uintptr], reflect.Uintptr},
		{"float32", types.Typ[types.Float32], reflect.Float32},
		{"float64", types.Typ[types.Float64], reflect.Float64},
		{"complex64", types.Typ[types.Complex64], reflect.Complex64},
		{"complex128", types.Typ[types.Complex128], reflect.Complex128},
		{"string", types.Typ[types.String], reflect.String},
		// Note: unsafe.Pointer is not supported
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := typeToReflect(tt.typ, prog)
			if rt == nil {
				t.Fatalf("typeToReflect(%s) returned nil", tt.name)
			}
			if rt.Kind() != tt.want {
				t.Errorf("typeToReflect(%s).Kind() = %v, want %v", tt.name, rt.Kind(), tt.want)
			}
		})
	}
}

func TestTypeToReflectSlice(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	sliceType := types.NewSlice(types.Typ[types.Int])
	rt := typeToReflect(sliceType, prog)

	if rt == nil {
		t.Fatal("typeToReflect([]int) returned nil")
	}
	if rt.Kind() != reflect.Slice {
		t.Errorf("typeToReflect([]int).Kind() = %v, want Slice", rt.Kind())
	}
	if rt.Elem().Kind() != reflect.Int {
		t.Errorf("typeToReflect([]int).Elem().Kind() = %v, want Int", rt.Elem().Kind())
	}
}

func TestTypeToReflectArray(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	arrayType := types.NewArray(types.Typ[types.Int], 5)
	rt := typeToReflect(arrayType, prog)

	if rt == nil {
		t.Fatal("typeToReflect([5]int) returned nil")
	}
	if rt.Kind() != reflect.Array {
		t.Errorf("typeToReflect([5]int).Kind() = %v, want Array", rt.Kind())
	}
	if rt.Len() != 5 {
		t.Errorf("typeToReflect([5]int).Len() = %d, want 5", rt.Len())
	}
}

func TestTypeToReflectMap(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	mapType := types.NewMap(types.Typ[types.String], types.Typ[types.Int])
	rt := typeToReflect(mapType, prog)

	if rt == nil {
		t.Fatal("typeToReflect(map[string]int) returned nil")
	}
	if rt.Kind() != reflect.Map {
		t.Errorf("typeToReflect(map[string]int).Kind() = %v, want Map", rt.Kind())
	}
	if rt.Key().Kind() != reflect.String {
		t.Errorf("typeToReflect(map[string]int).Key().Kind() = %v, want String", rt.Key().Kind())
	}
	if rt.Elem().Kind() != reflect.Int {
		t.Errorf("typeToReflect(map[string]int).Elem().Kind() = %v, want Int", rt.Elem().Kind())
	}
}

func TestTypeToReflectChan(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	chanType := types.NewChan(types.SendRecv, types.Typ[types.Int])
	rt := typeToReflect(chanType, prog)

	if rt == nil {
		t.Fatal("typeToReflect(chan int) returned nil")
	}
	if rt.Kind() != reflect.Chan {
		t.Errorf("typeToReflect(chan int).Kind() = %v, want Chan", rt.Kind())
	}
}

func TestTypeToReflectPointer(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	ptrType := types.NewPointer(types.Typ[types.Int])
	rt := typeToReflect(ptrType, prog)

	if rt == nil {
		t.Fatal("typeToReflect(*int) returned nil")
	}
	if rt.Kind() != reflect.Ptr {
		t.Errorf("typeToReflect(*int).Kind() = %v, want Ptr", rt.Kind())
	}
	if rt.Elem().Kind() != reflect.Int {
		t.Errorf("typeToReflect(*int).Elem().Kind() = %v, want Int", rt.Elem().Kind())
	}
}

func TestTypeToReflectInterface(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	ifaceType := types.NewInterfaceType(nil, nil)
	rt := typeToReflect(ifaceType, prog)

	if rt == nil {
		t.Fatal("typeToReflect(interface{}) returned nil")
	}
	if rt.Kind() != reflect.Interface {
		t.Errorf("typeToReflect(interface{}).Kind() = %v, want Interface", rt.Kind())
	}
}

func TestTypeToReflectNamed(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	// Create a named type "MyInt" with underlying type int
	named := types.NewNamed(
		types.NewTypeName(0, nil, "MyInt", nil),
		types.Typ[types.Int],
		nil,
	)

	rt := typeToReflect(named, prog)

	if rt == nil {
		t.Fatal("typeToReflect(MyInt) returned nil")
	}
	// Named int type should convert to int
	if rt.Kind() != reflect.Int {
		t.Errorf("typeToReflect(MyInt).Kind() = %v, want Int", rt.Kind())
	}
}

func TestTypeToReflectStruct(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	// Create a struct type with two fields
	fields := []*types.Var{
		types.NewField(0, nil, "X", types.Typ[types.Int], false),
		types.NewField(0, nil, "Y", types.Typ[types.String], false),
	}
	structType := types.NewStruct(fields, nil)

	rt := typeToReflect(structType, prog)

	if rt == nil {
		t.Fatal("typeToReflect(struct{X int; Y string}) returned nil")
	}
	if rt.Kind() != reflect.Struct {
		t.Errorf("typeToReflect(struct).Kind() = %v, want Struct", rt.Kind())
	}
	if rt.NumField() != 2 {
		t.Errorf("typeToReflect(struct).NumField() = %d, want 2", rt.NumField())
	}
}

func TestTypeToReflectSignature(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	// Create a function type: func(int, string) (int, error)
	params := types.NewTuple(
		types.NewVar(0, nil, "", types.Typ[types.Int]),
		types.NewVar(0, nil, "", types.Typ[types.String]),
	)
	results := types.NewTuple(
		types.NewVar(0, nil, "", types.Typ[types.Int]),
	)
	sig := types.NewSignatureType(nil, nil, nil, params, results, false)

	rt := typeToReflect(sig, prog)

	if rt == nil {
		t.Fatal("typeToReflect(func) returned nil")
	}
	if rt.Kind() != reflect.Func {
		t.Errorf("typeToReflect(func).Kind() = %v, want Func", rt.Kind())
	}
	if rt.NumIn() != 2 {
		t.Errorf("typeToReflect(func).NumIn() = %d, want 2", rt.NumIn())
	}
	if rt.NumOut() != 1 {
		t.Errorf("typeToReflect(func).NumOut() = %d, want 1", rt.NumOut())
	}
}

func TestTypeToReflectVariadic(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	// Create a variadic function type: func(int, ...string)
	params := types.NewTuple(
		types.NewVar(0, nil, "", types.Typ[types.Int]),
		types.NewVar(0, nil, "", types.NewSlice(types.Typ[types.String])),
	)
	sig := types.NewSignatureType(nil, nil, nil, params, nil, true)

	rt := typeToReflect(sig, prog)

	if rt == nil {
		t.Fatal("typeToReflect(variadic func) returned nil")
	}
	if !rt.IsVariadic() {
		t.Error("typeToReflect(variadic func).IsVariadic() = false, want true")
	}
}

func TestTypeToReflectWithNamedPackage(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	pkg := types.NewPackage("example.com/test", "testpkg")
	named := types.NewNamed(
		types.NewTypeName(0, pkg, "MyStruct", nil),
		types.NewStruct([]*types.Var{
			types.NewField(0, nil, "Value", types.Typ[types.Int], false),
		}, nil),
		nil,
	)

	rt := typeToReflect(named, prog)

	if rt == nil {
		t.Fatal("typeToReflect(named struct) returned nil")
	}
	if rt.Kind() != reflect.Struct {
		t.Errorf("typeToReflect(named struct).Kind() = %v, want Struct", rt.Kind())
	}
}

func TestTypeToReflectStructWithUnexportedFields(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	pkg := types.NewPackage("example.com/test", "testpkg")
	fields := []*types.Var{
		types.NewField(0, pkg, "Exported", types.Typ[types.Int], false),
		types.NewField(0, pkg, "unexported", types.Typ[types.String], false),
	}
	structType := types.NewStruct(fields, nil)

	rt := typeToReflect(structType, prog)

	if rt == nil {
		t.Fatal("typeToReflect(struct with unexported fields) returned nil")
	}
	if rt.NumField() != 2 {
		t.Errorf("typeToReflect(struct).NumField() = %d, want 2", rt.NumField())
	}
}

// ---------------------------------------------------------------------------
// typeToReflect caching tests - verify cache is used
// ---------------------------------------------------------------------------

func TestTypeToReflectCaching(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	intType := types.Typ[types.Int]

	// First call
	rt1 := typeToReflect(intType, prog)

	// Second call should return the same cached type
	rt2 := typeToReflect(intType, prog)

	if rt1 != rt2 {
		t.Error("typeToReflect should return cached type on second call")
	}
}

// ---------------------------------------------------------------------------
// typeToReflect edge cases
// ---------------------------------------------------------------------------

func TestTypeToReflectNil(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	rt := typeToReflect(nil, prog)
	if rt != nil {
		t.Errorf("typeToReflect(nil) = %v, want nil", rt)
	}
}

func TestTypeToReflectEmptyStruct(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	structType := types.NewStruct(nil, nil)

	rt := typeToReflect(structType, prog)

	// Empty struct should return struct{} type
	expected := reflect.TypeOf(struct{}{})
	if rt != expected {
		t.Errorf("typeToReflect(empty struct) = %v, want %v", rt, expected)
	}
}

// ---------------------------------------------------------------------------
// typeToReflect with complex nested types
// ---------------------------------------------------------------------------

func TestTypeToReflectNestedSlice(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	// [][]int
	nestedSlice := types.NewSlice(types.NewSlice(types.Typ[types.Int]))
	rt := typeToReflect(nestedSlice, prog)

	if rt == nil {
		t.Fatal("typeToReflect([][]int) returned nil")
	}
	if rt.Kind() != reflect.Slice {
		t.Errorf("Kind = %v, want Slice", rt.Kind())
	}
	if rt.Elem().Kind() != reflect.Slice {
		t.Errorf("Elem().Kind() = %v, want Slice", rt.Elem().Kind())
	}
	if rt.Elem().Elem().Kind() != reflect.Int {
		t.Errorf("Elem().Elem().Kind() = %v, want Int", rt.Elem().Elem().Kind())
	}
}

func TestTypeToReflectMapWithStructKey(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	fields := []*types.Var{
		types.NewField(0, nil, "Name", types.Typ[types.String], false),
	}
	structType := types.NewStruct(fields, nil)
	mapType := types.NewMap(structType, types.Typ[types.Int])

	rt := typeToReflect(mapType, prog)

	if rt == nil {
		t.Fatal("typeToReflect(map[struct]int) returned nil")
	}
	if rt.Kind() != reflect.Map {
		t.Errorf("Kind = %v, want Map", rt.Kind())
	}
	if rt.Key().Kind() != reflect.Struct {
		t.Errorf("Key().Kind() = %v, want Struct", rt.Key().Kind())
	}
}

func TestTypeToReflectStructWithFuncField(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Types: []types.Type{},
	}

	params := types.NewTuple(
		types.NewVar(0, nil, "", types.Typ[types.Int]),
	)
	results := types.NewTuple(
		types.NewVar(0, nil, "", types.Typ[types.String]),
	)
	funcType := types.NewSignatureType(nil, nil, nil, params, results, false)

	fields := []*types.Var{
		types.NewField(0, nil, "Callback", funcType, false),
	}
	structType := types.NewStruct(fields, nil)

	rt := typeToReflect(structType, prog)

	if rt == nil {
		t.Fatal("typeToReflect(struct with func field) returned nil")
	}
	if rt.Kind() != reflect.Struct {
		t.Errorf("Kind = %v, want Struct", rt.Kind())
	}
	if rt.Field(0).Type.Kind() != reflect.Func {
		t.Errorf("Field(0).Type.Kind() = %v, want Func", rt.Field(0).Type.Kind())
	}
}
