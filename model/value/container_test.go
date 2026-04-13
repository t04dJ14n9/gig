package value

import (
	"reflect"
	"testing"
)

// ---------------------------------------------------------------------------
// Len tests
// ---------------------------------------------------------------------------

func TestLenString(t *testing.T) {
	v := MakeString("hello")
	if v.Len() != 5 {
		t.Errorf("Len(%q) = %d, want 5", v.String(), v.Len())
	}
}

func TestLenEmptyString(t *testing.T) {
	v := MakeString("")
	if v.Len() != 0 {
		t.Errorf("Len(empty string) = %d, want 0", v.Len())
	}
}

func TestLenSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	v := FromInterface(slice)
	if v.Len() != 5 {
		t.Errorf("Len(slice) = %d, want 5", v.Len())
	}
}

func TestLenInt64Slice(t *testing.T) {
	// Native int64 slice fast path
	slice := []int64{10, 20, 30}
	v := Value{kind: KindSlice, obj: slice}
	if v.Len() != 3 {
		t.Errorf("Len([]int64) = %d, want 3", v.Len())
	}
}

func TestLenMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	v := FromInterface(m)
	if v.Len() != 3 {
		t.Errorf("Len(map) = %d, want 3", v.Len())
	}
}

func TestLenChan(t *testing.T) {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	v := FromInterface(ch)
	if v.Len() != 2 {
		t.Errorf("Len(chan) = %d, want 2", v.Len())
	}
}

// ---------------------------------------------------------------------------
// Cap tests
// ---------------------------------------------------------------------------

func TestCapSlice(t *testing.T) {
	slice := make([]int, 3, 10)
	v := FromInterface(slice)
	if v.Cap() != 10 {
		t.Errorf("Cap(slice) = %d, want 10", v.Cap())
	}
}

func TestCapInt64Slice(t *testing.T) {
	// Native int64 slice fast path
	slice := make([]int64, 2, 8)
	v := Value{kind: KindSlice, obj: slice}
	if v.Cap() != 8 {
		t.Errorf("Cap([]int64) = %d, want 8", v.Cap())
	}
}

func TestCapChan(t *testing.T) {
	ch := make(chan int, 5)
	v := FromInterface(ch)
	if v.Cap() != 5 {
		t.Errorf("Cap(chan) = %d, want 5", v.Cap())
	}
}

// ---------------------------------------------------------------------------
// Index tests
// ---------------------------------------------------------------------------

func TestIndexString(t *testing.T) {
	v := MakeString("hello")
	for i := 0; i < 5; i++ {
		elem := v.Index(i)
		if elem.Kind() != KindUint {
			t.Errorf("Index(%d).Kind() = %v, want KindUint", i, elem.Kind())
		}
		if elem.Uint() != uint64("hello"[i]) {
			t.Errorf("Index(%d) = %d, want %d", i, elem.Uint(), "hello"[i])
		}
	}
}

func TestIndexSlice(t *testing.T) {
	slice := []int{10, 20, 30}
	v := FromInterface(slice)

	elem0 := v.Index(0)
	if elem0.Int() != 10 {
		t.Errorf("Index(0) = %d, want 10", elem0.Int())
	}

	elem2 := v.Index(2)
	if elem2.Int() != 30 {
		t.Errorf("Index(2) = %d, want 30", elem2.Int())
	}
}

func TestIndexInt64Slice(t *testing.T) {
	// Native int64 slice fast path
	slice := []int64{100, 200, 300}
	v := Value{kind: KindSlice, obj: slice}

	elem := v.Index(1)
	if elem.Int() != 200 {
		t.Errorf("Index(1) = %d, want 200", elem.Int())
	}
}

func TestIndexValueSlice(t *testing.T) {
	// Native []Value slice
	slice := []Value{MakeInt(1), MakeInt(2), MakeInt(3)}
	v := Value{kind: KindSlice, obj: slice}

	elem := v.Index(0)
	if elem.Int() != 1 {
		t.Errorf("Index(0) = %d, want 1", elem.Int())
	}
}

// ---------------------------------------------------------------------------
// SetIndex tests
// ---------------------------------------------------------------------------

func TestSetIndexSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	v := FromInterface(slice)

	v.SetIndex(1, MakeInt(20))
	if slice[1] != 20 {
		t.Errorf("SetIndex failed: slice[1] = %d, want 20", slice[1])
	}
}

func TestSetIndexInt64Slice(t *testing.T) {
	// Native int64 slice fast path
	slice := []int64{100, 200, 300}
	v := Value{kind: KindSlice, obj: slice}

	v.SetIndex(0, MakeInt(999))
	if slice[0] != 999 {
		t.Errorf("SetIndex failed: slice[0] = %d, want 999", slice[0])
	}
}

func TestSetIndexValueSlice(t *testing.T) {
	// Native []Value slice
	slice := []Value{MakeInt(1), MakeInt(2), MakeInt(3)}
	v := Value{kind: KindSlice, obj: slice}

	v.SetIndex(2, MakeInt(33))
	if slice[2].Int() != 33 {
		t.Errorf("SetIndex failed: slice[2] = %d, want 33", slice[2].Int())
	}
}

// ---------------------------------------------------------------------------
// MapIndex tests
// ---------------------------------------------------------------------------

func TestMapIndex(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	v := FromInterface(m)

	key := MakeString("a")
	val := v.MapIndex(key)
	if val.Int() != 1 {
		t.Errorf("MapIndex(\"a\") = %d, want 1", val.Int())
	}

	key2 := MakeString("b")
	val2 := v.MapIndex(key2)
	if val2.Int() != 2 {
		t.Errorf("MapIndex(\"b\") = %d, want 2", val2.Int())
	}
}

func TestMapIndexNotFound(t *testing.T) {
	m := map[string]int{"a": 1}
	v := FromInterface(m)

	key := MakeString("nonexistent")
	val := v.MapIndex(key)
	// Should return zero value for element type
	if val.Int() != 0 {
		t.Errorf("MapIndex(nonexistent) = %d, want 0 (zero value)", val.Int())
	}
}

func TestMapIndexIntKey(t *testing.T) {
	m := map[int]string{1: "one", 2: "two"}
	v := FromInterface(m)

	key := MakeInt(1)
	val := v.MapIndex(key)
	if val.String() != "one" {
		t.Errorf("MapIndex(1) = %q, want \"one\"", val.String())
	}
}

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// Field tests
// ---------------------------------------------------------------------------

type testStruct struct {
	X int
	Y string
	z float64 // unexported
}

func TestField(t *testing.T) {
	s := testStruct{X: 42, Y: "hello", z: 3.14}
	v := FromInterface(&s).Elem()

	xField := v.Field(0)
	if xField.Int() != 42 {
		t.Errorf("Field(0) = %d, want 42", xField.Int())
	}

	yField := v.Field(1)
	if yField.String() != "hello" {
		t.Errorf("Field(1) = %q, want \"hello\"", yField.String())
	}
}

func TestFieldUnexported(t *testing.T) {
	s := testStruct{X: 42, Y: "hello", z: 3.14}
	v := FromInterface(&s).Elem()

	zField := v.Field(2)
	if zField.Float() != 3.14 {
		t.Errorf("Field(2) = %f, want 3.14", zField.Float())
	}
}

// ---------------------------------------------------------------------------
// SetField tests
// ---------------------------------------------------------------------------

func TestSetField(t *testing.T) {
	s := testStruct{X: 0, Y: ""}
	v := FromInterface(&s).Elem()

	v.SetField(0, MakeInt(100))
	if s.X != 100 {
		t.Errorf("SetField(0) failed: s.X = %d, want 100", s.X)
	}

	v.SetField(1, MakeString("world"))
	if s.Y != "world" {
		t.Errorf("SetField(1) failed: s.Y = %q, want \"world\"", s.Y)
	}
}

// ---------------------------------------------------------------------------
// Elem tests
// ---------------------------------------------------------------------------

func TestElemPointer(t *testing.T) {
	x := 42
	ptr := &x
	v := FromInterface(ptr)

	elem := v.Elem()
	if elem.Int() != 42 {
		t.Errorf("Elem() = %d, want 42", elem.Int())
	}
}

func TestElemPointerToInt64(t *testing.T) {
	// Fast path: *int64 pointer
	x := int64(100)
	v := Value{kind: KindReflect, obj: reflect.ValueOf(&x)}

	elem := v.Elem()
	if elem.Int() != 100 {
		t.Errorf("Elem() = %d, want 100", elem.Int())
	}
}

func TestElemValuePointer(t *testing.T) {
	// Fast path: *Value pointer
	inner := MakeInt(42)
	v := Value{kind: KindReflect, obj: reflect.ValueOf(&inner)}

	elem := v.Elem()
	if elem.Int() != 42 {
		t.Errorf("Elem() = %d, want 42", elem.Int())
	}
}

func TestElemInterface(t *testing.T) {
	var iface any = int64(99)
	v := FromInterface(&iface).Elem()

	elem := v.Elem()
	if elem.Int() != 99 {
		t.Errorf("Elem() = %d, want 99", elem.Int())
	}
}

// ---------------------------------------------------------------------------
// SetElem tests
// ---------------------------------------------------------------------------

func TestSetElemPointer(t *testing.T) {
	x := 0
	ptr := &x
	v := FromInterface(ptr)

	v.SetElem(MakeInt(123))
	if x != 123 {
		t.Errorf("SetElem failed: x = %d, want 123", x)
	}
}

func TestSetElemPointerToInt64(t *testing.T) {
	// Fast path: *int64 pointer
	x := int64(0)
	v := Value{kind: KindReflect, obj: reflect.ValueOf(&x)}

	v.SetElem(MakeInt(456))
	if x != 456 {
		t.Errorf("SetElem failed: x = %d, want 456", x)
	}
}

func TestSetElemValuePointer(t *testing.T) {
	// Fast path: *Value pointer
	inner := MakeInt(0)
	v := Value{kind: KindReflect, obj: reflect.ValueOf(&inner)}

	v.SetElem(MakeInt(789))
	if inner.Int() != 789 {
		t.Errorf("SetElem failed: inner = %d, want 789", inner.Int())
	}
}

// ---------------------------------------------------------------------------
// Edge cases and panics
// ---------------------------------------------------------------------------

func TestLenPanicOnInvalidKind(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Len on invalid kind should panic")
		}
	}()

	v := MakeInt(42)
	v.Len() // Should panic
}

func TestCapPanicOnInvalidKind(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Cap on invalid kind should panic")
		}
	}()

	v := MakeInt(42)
	v.Cap() // Should panic
}

func TestIndexPanicOnInvalidKind(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Index on invalid kind should panic")
		}
	}()

	v := MakeInt(42)
	v.Index(0) // Should panic
}

func TestUnsafeAddrOf(t *testing.T) {
	s := struct{ x int }{x: 42}
	rv := reflect.ValueOf(&s).Elem().Field(0)

	ptr := UnsafeAddrOf(rv)
	if ptr == nil {
		t.Error("UnsafeAddrOf returned nil")
	}
}
