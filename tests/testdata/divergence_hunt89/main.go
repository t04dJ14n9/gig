package divergence_hunt89

// ============================================================================
// Round 89: Interface composition and embedding
// ============================================================================

type Reader interface {
	Read() string
}

type Writer interface {
	Write(s string)
}

type ReadWriter interface {
	Reader
	Writer
}

type TextFile struct {
	content string
}

func (t *TextFile) Read() string {
	return t.content
}

func (t *TextFile) Write(s string) {
	t.content = s
}

func InterfaceComposition() string {
	var rw ReadWriter = &TextFile{}
	rw.Write("hello")
	return rw.Read()
}

func InterfaceEmbedding() string {
	var r Reader = &TextFile{content: "world"}
	return r.Read()
}

func InterfaceAssertionComposition() string {
	var rw ReadWriter = &TextFile{}
	rw.Write("test")
	// assert back to Reader
	var r Reader = rw
	return r.Read()
}

func InterfaceSliceOfInterface() int {
	items := []any{1, "hello", true, 3.14}
	return len(items)
}

func InterfaceMapOfInterface() int {
	m := map[string]any{
		"int":    42,
		"string": "world",
		"bool":   true,
	}
	return len(m)
}

type Stringer interface {
	String() string
}

type MyString string

func (m MyString) String() string {
	return string(m)
}

func InterfaceCustomStringer() string {
	var s Stringer = MyString("hello")
	return s.String()
}

func InterfaceSliceAsAny() int {
	s := []int{1, 2, 3}
	var x any = s
	return len(x.([]int))
}

func InterfaceMapAsAny() int {
	m := map[string]int{"a": 1}
	var x any = m
	return len(x.(map[string]int))
}

func InterfaceFuncAsAny() int {
	f := func(x int) int { return x * 2 }
	var x any = f
	return x.(func(int) int)(21)
}

func NilInterfaceAssertion() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = -1
		}
	}()
	var x any
	_ = x.(int) // should panic
	return 0
}

func EmptyInterfaceTypeSwitch() string {
	var x any = 3.14
	switch x.(type) {
	case int:
		return "int"
	case float64:
		return "float64"
	case string:
		return "string"
	default:
		return "other"
	}
}
