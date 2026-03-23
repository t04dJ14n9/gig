package compiler

import (
	"go/token"
	"go/types"
	"reflect"
	"testing"

	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/importer"
	"git.woa.com/youngjin/gig/value"
)

// ---------------------------------------------------------------------------
// Mock SSA types
// ---------------------------------------------------------------------------

// mockSSAValue implements ssa.Value for testing.
type mockSSAValue struct {
	name      string
	str       string
	typ       types.Type
	parent    *ssa.Function
	referrers []ssa.Instruction
}

func (m *mockSSAValue) Name() string           { return m.name }
func (m *mockSSAValue) String() string        { return m.str }
func (m *mockSSAValue) Type() types.Type       { return m.typ }
func (m *mockSSAValue) Parent() *ssa.Function  { return m.parent }
func (m *mockSSAValue) Referrers() *[]ssa.Instruction {
	if m.referrers == nil {
		return nil
	}
	return &m.referrers
}
func (m *mockSSAValue) Pos() token.Pos { return token.NoPos }

// mockBasicBlock implements ssa.BasicBlock for testing.
type mockBasicBlock struct {
	succs []*mockBasicBlock
	preds []*mockBasicBlock
}

// ---------------------------------------------------------------------------
// isIntType tests
// ---------------------------------------------------------------------------

func TestIsIntType(t *testing.T) {
	tests := []struct {
		name string
		t    types.Type
		want bool
	}{
		{"nil", nil, false},
		{"int", types.Typ[types.Int], true},
		{"int8", types.Typ[types.Int8], true},
		{"int16", types.Typ[types.Int16], true},
		{"int32", types.Typ[types.Int32], true},
		{"int64", types.Typ[types.Int64], true},
		{"uint", types.Typ[types.Uint], false},
		{"uint8", types.Typ[types.Uint8], false},
		{"uint16", types.Typ[types.Uint16], false},
		{"uint32", types.Typ[types.Uint32], false},
		{"uint64", types.Typ[types.Uint64], false},
		{"float32", types.Typ[types.Float32], false},
		{"float64", types.Typ[types.Float64], false},
		{"bool", types.Typ[types.Bool], false},
		{"string", types.Typ[types.String], false},
		{"complex64", types.Typ[types.Complex64], false},
		{"complex128", types.Typ[types.Complex128], false},
		{"unsafe.Pointer", types.Typ[types.UnsafePointer], false},
		{"uintptr", types.Typ[types.Uintptr], false},
		{"any", types.Typ[types.UntypedNil], false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIntType(tt.t)
			if got != tt.want {
				t.Errorf("isIntType(%v) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}

func TestIsIntTypeNamed(t *testing.T) {
	// Test named types that are underlying int types
	// Create a named int type using types.NewNamed
	intType := types.Typ[types.Int]
	namedInt := types.NewNamed(
		types.NewTypeName(0, nil, "MyInt", nil),
		intType,
		nil,
	)

	if !isIntType(namedInt) {
		t.Error("isIntType(named int) = false, want true")
	}
}

func TestIsIntTypePtr(t *testing.T) {
	// Pointer types are not int types
	ptrInt := types.NewPointer(types.Typ[types.Int])
	if isIntType(ptrInt) {
		t.Error("isIntType(*int) = true, want false")
	}
}

// ---------------------------------------------------------------------------
// isIntSliceType tests
// ---------------------------------------------------------------------------

func TestIsIntSliceType(t *testing.T) {
	tests := []struct {
		name string
		t    types.Type
		want bool
	}{
		{"nil", nil, false},
		{"[]int", types.NewSlice(types.Typ[types.Int]), true},
		{"[]int64", types.NewSlice(types.Typ[types.Int64]), true},
		{"[]int8", types.NewSlice(types.Typ[types.Int8]), false},
		{"[]int16", types.NewSlice(types.Typ[types.Int16]), false},
		{"[]int32", types.NewSlice(types.Typ[types.Int32]), false},
		{"[]uint", types.NewSlice(types.Typ[types.Uint]), false},
		{"[]uint64", types.NewSlice(types.Typ[types.Uint64]), false},
		{"[]float32", types.NewSlice(types.Typ[types.Float32]), false},
		{"[]float64", types.NewSlice(types.Typ[types.Float64]), false},
		{"[]string", types.NewSlice(types.Typ[types.String]), false},
		{"[]bool", types.NewSlice(types.Typ[types.Bool]), false},
		{"[]byte", types.NewSlice(types.Typ[types.Uint8]), false},
		{"map[int]int", types.NewMap(types.Typ[types.Int], types.Typ[types.Int]), false},
		{"chan int", types.NewChan(types.SendRecv, types.Typ[types.Int]), false},
		{"*[0]int", types.NewArray(types.Typ[types.Int], 0), false},
		{"int", types.Typ[types.Int], false},
		{"string", types.Typ[types.String], false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIntSliceType(tt.t)
			if got != tt.want {
				t.Errorf("isIntSliceType(%v) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}

func TestIsIntSliceTypeNamed(t *testing.T) {
	// Test named slice types with int element
	intSlice := types.NewSlice(types.Typ[types.Int])
	namedSlice := types.NewNamed(
		types.NewTypeName(0, nil, "IntSlice", nil),
		intSlice,
		nil,
	)

	if !isIntSliceType(namedSlice) {
		t.Error("isIntSliceType(named []int) = false, want true")
	}
}

func TestIsIntSliceTypePtr(t *testing.T) {
	// Pointer types are not slice types
	ptrSlice := types.NewPointer(types.NewSlice(types.Typ[types.Int]))
	if isIntSliceType(ptrSlice) {
		t.Error("isIntSliceType(*[]int) = true, want false")
	}
}

// ---------------------------------------------------------------------------
// extractReceiverTypeName tests
// ---------------------------------------------------------------------------

func TestExtractReceiverTypeName(t *testing.T) {
	tests := []struct {
		name    string
		recv    types.Type
		want    string
	}{
		{
			name:    "nil",
			recv:    nil,
			want:    "",
		},
		{
			name:    "int (non-pointer, non-named)",
			recv:    types.Typ[types.Int],
			want:    "",
		},
		{
			name:    "*int pointer to basic",
			recv:    types.NewPointer(types.Typ[types.Int]),
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractReceiverTypeName(tt.recv)
			if got != tt.want {
				t.Errorf("extractReceiverTypeName(%v) = %q, want %q", tt.recv, got, tt.want)
			}
		})
	}
}

func TestExtractReceiverTypeNameNamed(t *testing.T) {
	// Create a named type
	namedType := types.NewNamed(
		types.NewTypeName(0, nil, "MyStruct", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	got := extractReceiverTypeName(namedType)
	if got != "MyStruct" {
		t.Errorf("extractReceiverTypeName(named) = %q, want %q", got, "MyStruct")
	}
}

func TestExtractReceiverTypeNamePointerToNamed(t *testing.T) {
	// Create a named type
	namedType := types.NewNamed(
		types.NewTypeName(0, nil, "MyStruct", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	ptrType := types.NewPointer(namedType)

	got := extractReceiverTypeName(ptrType)
	if got != "MyStruct" {
		t.Errorf("extractReceiverTypeName(*named) = %q, want %q", got, "MyStruct")
	}
}

func TestExtractReceiverTypeNameWithPackage(t *testing.T) {
	// Create a named type with a package
	pkg := types.NewPackage("my/pkg", "pkg")
	namedType := types.NewNamed(
		types.NewTypeName(0, pkg, "MyStruct", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	got := extractReceiverTypeName(namedType)
	if got != "my/pkg.MyStruct" {
		t.Errorf("extractReceiverTypeName(pkg.named) = %q, want %q", got, "my/pkg.MyStruct")
	}
}

func TestExtractReceiverTypeNamePointerWithPackage(t *testing.T) {
	// Create a named type with a package
	pkg := types.NewPackage("my/pkg", "pkg")
	namedType := types.NewNamed(
		types.NewTypeName(0, pkg, "MyStruct", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	ptrType := types.NewPointer(namedType)

	got := extractReceiverTypeName(ptrType)
	if got != "my/pkg.MyStruct" {
		t.Errorf("extractReceiverTypeName(*pkg.named) = %q, want %q", got, "my/pkg.MyStruct")
	}
}

// ---------------------------------------------------------------------------
// extractNamedType tests
// ---------------------------------------------------------------------------

func TestExtractNamedType(t *testing.T) {
	tests := []struct {
		name string
		t    types.Type
		want bool
	}{
		{
			name: "nil",
			t:    nil,
			want: false,
		},
		{
			name: "int",
			t:    types.Typ[types.Int],
			want: false,
		},
		{
			name: "[]int slice",
			t:    types.NewSlice(types.Typ[types.Int]),
			want: false,
		},
		{
			name: "map[int]int",
			t:    types.NewMap(types.Typ[types.Int], types.Typ[types.Int]),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractNamedType(tt.t)
			if (got != nil) != tt.want {
				t.Errorf("extractNamedType(%v) = %v, want nil=%v", tt.t, got, !tt.want)
			}
		})
	}
}

func TestExtractNamedTypeDirect(t *testing.T) {
	namedType := types.NewNamed(
		types.NewTypeName(0, nil, "MyStruct", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	got := extractNamedType(namedType)
	if got != namedType {
		t.Errorf("extractNamedType(named) = %v, want %v", got, namedType)
	}
}

func TestExtractNamedTypePointer(t *testing.T) {
	namedType := types.NewNamed(
		types.NewTypeName(0, nil, "MyStruct", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	ptrType := types.NewPointer(namedType)

	got := extractNamedType(ptrType)
	if got != namedType {
		t.Errorf("extractNamedType(*named) = %v, want %v", got, namedType)
	}
}

func TestExtractNamedTypeDoublePointer(t *testing.T) {
	// ExtractNamedType unwraps ALL pointer levels until it finds a Named type.
	// So **MyStruct -> *MyStruct -> MyStruct (Named) -> returns MyStruct
	namedType := types.NewNamed(
		types.NewTypeName(0, nil, "MyStruct", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	ptrType := types.NewPointer(namedType)
	doublePtr := types.NewPointer(ptrType)

	got := extractNamedType(doublePtr)
	if got != namedType {
		t.Errorf("extractNamedType(**named) = %v, want %v", got, namedType)
	}
}

// ---------------------------------------------------------------------------
// reversePostorder tests
// ---------------------------------------------------------------------------

func TestReversePostorderEmpty(t *testing.T) {
	fn := &ssa.Function{}
	got := reversePostorder(fn)
	if got != nil {
		t.Errorf("reversePostorder(empty fn) = %v, want nil", got)
	}
}

func TestReversePostorderSingleBlockNoSuccs(t *testing.T) {
	// Single block with no successors
	block := &ssa.BasicBlock{}
	fn := &ssa.Function{
		Blocks: []*ssa.BasicBlock{block},
	}

	got := reversePostorder(fn)
	if len(got) != 1 {
		t.Errorf("reversePostorder(single block) len = %d, want 1", len(got))
	}
}

func TestReversePostorderFirstBlockOnly(t *testing.T) {
	// Only first block (index 0) is visited if blocks aren't connected via Succs
	// This is expected behavior since reversePostorder traverses via Succs
	block0 := &ssa.BasicBlock{}
	block1 := &ssa.BasicBlock{}
	fn := &ssa.Function{
		Blocks: []*ssa.BasicBlock{block0, block1},
	}

	got := reversePostorder(fn)
	// Only block0 is reachable from fn.Blocks[0] with no Succs set
	if len(got) != 1 {
		t.Errorf("reversePostorder(len = %d, want 1 (only first block reachable))", len(got))
	}
	if got[0] != block0 {
		t.Errorf("reversePostorder[0] != block0")
	}
}

// ---------------------------------------------------------------------------
// SymbolTable with mock SSA values
// ---------------------------------------------------------------------------

func TestSymbolTableAllocLocal(t *testing.T) {
	st := NewSymbolTable()

	// Create mock SSA values
	v1 := &mockSSAValue{name: "v1", typ: types.Typ[types.Int]}
	v2 := &mockSSAValue{name: "v2", typ: types.Typ[types.String]}
	v3 := &mockSSAValue{name: "v3", typ: types.Typ[types.Int64]}

	// Test AllocLocal
	idx1 := st.AllocLocal(v1)
	if idx1 != 0 {
		t.Errorf("AllocLocal(v1) = %d, want 0", idx1)
	}

	idx2 := st.AllocLocal(v2)
	if idx2 != 1 {
		t.Errorf("AllocLocal(v2) = %d, want 1", idx2)
	}

	idx3 := st.AllocLocal(v3)
	if idx3 != 2 {
		t.Errorf("AllocLocal(v3) = %d, want 2", idx3)
	}

	// AllocLocal again for same value should return same index
	idx1Again := st.AllocLocal(v1)
	if idx1Again != idx1 {
		t.Errorf("AllocLocal(v1 again) = %d, want %d", idx1Again, idx1)
	}

	// Test NumLocals
	if n := st.NumLocals(); n != 3 {
		t.Errorf("NumLocals() = %d, want 3", n)
	}
}

func TestSymbolTableGetLocal(t *testing.T) {
	st := NewSymbolTable()

	v1 := &mockSSAValue{name: "v1", typ: types.Typ[types.Int]}
	v2 := &mockSSAValue{name: "v2", typ: types.Typ[types.String]}

	// Get before alloc should return false
	if idx, ok := st.GetLocal(v1); ok || idx != 0 {
		t.Errorf("GetLocal(v1 before alloc) = (%d, %v), want (0, false)", idx, ok)
	}

	st.AllocLocal(v1)
	st.AllocLocal(v2)

	// Get after alloc
	idx, ok := st.GetLocal(v1)
	if !ok || idx != 0 {
		t.Errorf("GetLocal(v1) = (%d, %v), want (0, true)", idx, ok)
	}

	idx2, ok2 := st.GetLocal(v2)
	if !ok2 || idx2 != 1 {
		t.Errorf("GetLocal(v2) = (%d, %v), want (1, true)", idx2, ok2)
	}

	// Get non-existent
	unknown := &mockSSAValue{name: "unknown", typ: types.Typ[types.Int]}
	if _, ok := st.GetLocal(unknown); ok {
		t.Error("GetLocal(unknown) = true, want false")
	}
}

func TestSymbolTableFreeVars(t *testing.T) {
	st := NewSymbolTable()

	// Test that freeVars map is initialized
	if st.freeVars == nil {
		t.Error("freeVars map is nil")
	}

	fv1 := &mockSSAValue{name: "fv1", typ: types.Typ[types.Int]}
	fv2 := &mockSSAValue{name: "fv2", typ: types.Typ[types.String]}

	// Manually set free vars (simulating what compileFunction does)
	st.freeVars[fv1] = 0
	st.freeVars[fv2] = 1

	if idx, ok := st.freeVars[fv1]; !ok || idx != 0 {
		t.Errorf("freeVars[fv1] = (%d, %v), want (0, true)", idx, ok)
	}

	if idx, ok := st.freeVars[fv2]; !ok || idx != 1 {
		t.Errorf("freeVars[fv2] = (%d, %v), want (1, true)", idx, ok)
	}
}

// ---------------------------------------------------------------------------
// Compiler interface and NewCompiler
// ---------------------------------------------------------------------------



func TestNewCompilerFields(t *testing.T) {
	lookup := &mockLookup{}
	c := NewCompiler(lookup).(*compiler)

	if c.lookup != lookup {
		t.Error("compiler.lookup != lookup")
	}
	if c.constants == nil {
		t.Error("compiler.constants is nil")
	}
	if c.types == nil {
		t.Error("compiler.types is nil")
	}
	if c.globals == nil {
		t.Error("compiler.globals is nil")
	}
	if c.externalVarValues == nil {
		t.Error("compiler.externalVarValues is nil")
	}
	if c.funcs == nil {
		t.Error("compiler.funcs is nil")
	}
	if c.funcIndex == nil {
		t.Error("compiler.funcIndex is nil")
	}
}

func TestCompilerInterfaceContract(t *testing.T) {
	lookup := &mockLookup{}
	c := NewCompiler(lookup)

	defer func() {
		if r := recover(); r == nil {
			t.Log("Compile(nil) returned without panic (may return error)")
		}
	}()

	_, _ = c.Compile(nil)
}

func TestPackageLevelCompile(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Log("Compile(nil, nil) returned without panic")
		}
	}()

	_, _ = Compile(&mockLookup{}, nil)
}

// ---------------------------------------------------------------------------
// Mock PackageLookup
// ---------------------------------------------------------------------------

type mockLookup struct{}

func (m *mockLookup) LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool) {
	return nil, nil, false
}

func (m *mockLookup) LookupMethodDirectCall(typeName, methodName string) (directCall func([]value.Value) value.Value, ok bool) {
	return nil, false
}

func (m *mockLookup) LookupExternalVar(pkgPath, varName string) (ptr any, ok bool) {
	return nil, false
}

func (m *mockLookup) LookupExternalType(t types.Type) (reflect.Type, bool) {
	return nil, false
}

// Verify mockLookup satisfies the interface.
var _ importer.PackageLookup = (*mockLookup)(nil)

// ---------------------------------------------------------------------------
// jumpInfo struct
// ---------------------------------------------------------------------------

func TestJumpInfoStruct(t *testing.T) {
	block := &ssa.BasicBlock{}
	ji := jumpInfo{
		offset:      10,
		targetBlock: block,
	}

	if ji.offset != 10 {
		t.Errorf("jumpInfo.offset = %d, want 10", ji.offset)
	}
	if ji.targetBlock != block {
		t.Errorf("jumpInfo.targetBlock = %v, want %v", ji.targetBlock, block)
	}
}

// ---------------------------------------------------------------------------
// isIntSliceType edge cases
// ---------------------------------------------------------------------------

func TestIsIntSliceTypeEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		t    types.Type
		want bool
	}{
		{"nil", nil, false},
		{"[]int", types.NewSlice(types.Typ[types.Int]), true},
		{"[]int64", types.NewSlice(types.Typ[types.Int64]), true},
		{"[]uint", types.NewSlice(types.Typ[types.Uint]), false},
		{"[]uint64", types.NewSlice(types.Typ[types.Uint64]), false},
		{"[]float64", types.NewSlice(types.Typ[types.Float64]), false},
		{"[]string", types.NewSlice(types.Typ[types.String]), false},
		{"[]bool", types.NewSlice(types.Typ[types.Bool]), false},
		{"[]*int (pointer slice)", types.NewSlice(types.NewPointer(types.Typ[types.Int])), false},
		{"[]interface{}", types.NewSlice(types.Typ[types.UntypedNil]), false},
		{"[]unnamed int", types.NewSlice(types.Typ[types.Int]), true},
		{"named []int", newNamedSlice("MyIntSlice", types.NewSlice(types.Typ[types.Int])), true},
		{"named []uint", newNamedSlice("MyUintSlice", types.NewSlice(types.Typ[types.Uint])), false},
		{"map[int]int", types.NewMap(types.Typ[types.Int], types.Typ[types.Int]), false},
		{"chan int", types.NewChan(types.SendRecv, types.Typ[types.Int]), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIntSliceType(tt.t)
			if got != tt.want {
				t.Errorf("isIntSliceType(%v) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}

// newNamedSlice creates a named slice type for testing.
func newNamedSlice(name string, slice *types.Slice) types.Type {
	return types.NewNamed(
		types.NewTypeName(0, nil, name, nil),
		slice,
		nil,
	)
}

// ---------------------------------------------------------------------------
// SymbolTable edge cases
// ---------------------------------------------------------------------------

func TestSymbolTableNumLocalsEmpty(t *testing.T) {
	st := NewSymbolTable()
	if n := st.NumLocals(); n != 0 {
		t.Errorf("NewSymbolTable().NumLocals() = %d, want 0", n)
	}
}

func TestSymbolTableMultipleLocalsSameType(t *testing.T) {
	st := NewSymbolTable()

	// Allocate multiple locals of the same type
	v1 := &mockSSAValue{name: "a", typ: types.Typ[types.Int]}
	v2 := &mockSSAValue{name: "b", typ: types.Typ[types.Int]}
	v3 := &mockSSAValue{name: "c", typ: types.Typ[types.Int]}

	idx1 := st.AllocLocal(v1)
	idx2 := st.AllocLocal(v2)
	idx3 := st.AllocLocal(v3)

	if idx1 != 0 || idx2 != 1 || idx3 != 2 {
		t.Errorf("AllocLocal indices = %d, %d, %d, want 0, 1, 2", idx1, idx2, idx3)
	}
	if n := st.NumLocals(); n != 3 {
		t.Errorf("NumLocals() = %d, want 3", n)
	}
}

func TestSymbolTableLocalLookupConsistency(t *testing.T) {
	st := NewSymbolTable()

	v := &mockSSAValue{name: "x", typ: types.Typ[types.Float64]}

	idx := st.AllocLocal(v)
	got, ok := st.GetLocal(v)

	if !ok {
		t.Error("GetLocal after AllocLocal returned ok=false")
	}
	if got != idx {
		t.Errorf("GetLocal(v) = %d, AllocLocal(v) = %d", got, idx)
	}
}

// ---------------------------------------------------------------------------
// Compiler creation and configuration
// ---------------------------------------------------------------------------

func TestNewCompilerWithNilLookup(t *testing.T) {
	// Should not panic
	c := NewCompiler(nil)
	if c == nil {
		t.Fatal("NewCompiler(nil) returned nil")
	}
}



// ---------------------------------------------------------------------------
// Compiler.Interface implementation verification
// ---------------------------------------------------------------------------



// ---------------------------------------------------------------------------
// Emit helper tests (reversePostorder edge cases)
// ---------------------------------------------------------------------------

func TestReversePostorderNilBlocks(t *testing.T) {
	fn := &ssa.Function{
		Blocks: nil,
	}
	got := reversePostorder(fn)
	if got != nil {
		t.Errorf("reversePostorder(nil blocks) = %v, want nil", got)
	}
}

func TestReversePostorderMultipleDisconnectedBlocks(t *testing.T) {
	// Only block 0 is reachable - block 1 is orphaned
	block0 := &ssa.BasicBlock{Index: 0}
	block1 := &ssa.BasicBlock{Index: 1}
	fn := &ssa.Function{
		Blocks: []*ssa.BasicBlock{block0, block1},
	}
	// block0 has no successors, so only block0 is visited

	got := reversePostorder(fn)
	if len(got) != 1 {
		t.Errorf("reversePostorder(disconnected) len = %d, want 1", len(got))
	}
	if got[0] != block0 {
		t.Errorf("reversePostorder[0] = %v, want block0", got[0])
	}
}

// ---------------------------------------------------------------------------
// extractNamedType edge cases
// ---------------------------------------------------------------------------

func TestExtractNamedTypeNil(t *testing.T) {
	got := extractNamedType(nil)
	if got != nil {
		t.Errorf("extractNamedType(nil) = %v, want nil", got)
	}
}

func TestExtractNamedTypeBasic(t *testing.T) {
	got := extractNamedType(types.Typ[types.Int])
	if got != nil {
		t.Errorf("extractNamedType(int) = %v, want nil", got)
	}
}

func TestExtractNamedTypeInterface(t *testing.T) {
	got := extractNamedType(types.Typ[types.UntypedNil])
	if got != nil {
		t.Errorf("extractNamedType(interface{}) = %v, want nil", got)
	}
}

func TestExtractNamedTypeSlice(t *testing.T) {
	slice := types.NewSlice(types.Typ[types.Int])
	got := extractNamedType(slice)
	if got != nil {
		t.Errorf("extractNamedType([]int) = %v, want nil", got)
	}
}

func TestExtractNamedTypeMap(t *testing.T) {
	m := types.NewMap(types.Typ[types.Int], types.Typ[types.String])
	got := extractNamedType(m)
	if got != nil {
		t.Errorf("extractNamedType(map[int]string) = %v, want nil", got)
	}
}

func TestExtractNamedTypeChan(t *testing.T) {
	ch := types.NewChan(types.SendRecv, types.Typ[types.Int])
	got := extractNamedType(ch)
	if got != nil {
		t.Errorf("extractNamedType(chan int) = %v, want nil", got)
	}
}

func TestExtractNamedTypePointerToSlice(t *testing.T) {
	// *[]int is not a named type
	ptr := types.NewPointer(types.NewSlice(types.Typ[types.Int]))
	got := extractNamedType(ptr)
	if got != nil {
		t.Errorf("extractNamedType(*[]int) = %v, want nil", got)
	}
}

func TestExtractNamedTypeTriplePointer(t *testing.T) {
	// ***int unwraps to int (basic), then returns nil
	namedType := types.NewNamed(
		types.NewTypeName(0, nil, "MyInt", nil),
		types.Typ[types.Int],
		nil,
	)
	ptr1 := types.NewPointer(namedType)
	ptr2 := types.NewPointer(ptr1)
	ptr3 := types.NewPointer(ptr2)

	// extractNamedType unwraps all pointers and returns Named
	got := extractNamedType(ptr3)
	if got != namedType {
		t.Errorf("extractNamedType(***MyInt) = %v, want %v", got, namedType)
	}
}

// ---------------------------------------------------------------------------
// extractReceiverTypeName edge cases
// ---------------------------------------------------------------------------

func TestExtractReceiverTypeNameSlice(t *testing.T) {
	slice := types.NewSlice(types.Typ[types.Int])
	got := extractReceiverTypeName(slice)
	if got != "" {
		t.Errorf("extractReceiverTypeName([]int) = %q, want empty", got)
	}
}

func TestExtractReceiverTypeNameMap(t *testing.T) {
	m := types.NewMap(types.Typ[types.Int], types.Typ[types.Int])
	got := extractReceiverTypeName(m)
	if got != "" {
		t.Errorf("extractReceiverTypeName(map[int]int) = %q, want empty", got)
	}
}

func TestExtractReceiverTypeNameChan(t *testing.T) {
	ch := types.NewChan(types.SendRecv, types.Typ[types.Int])
	got := extractReceiverTypeName(ch)
	if got != "" {
		t.Errorf("extractReceiverTypeName(chan int) = %q, want empty", got)
	}
}