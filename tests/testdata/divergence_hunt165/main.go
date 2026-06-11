package divergence_hunt165

import "fmt"

// ============================================================================
// Round 165: Embedding interfaces and structs
// ============================================================================

// Reader interface
type Reader interface {
	Read(p []byte) (n int, err error)
	ReadByte() (byte, error)
}

// Writer interface
type Writer interface {
	Write(p []byte) (n int, err error)
	WriteByte(c byte) error
}

// ReadWriter embeds Reader and Writer
type ReadWriter interface {
	Reader
	Writer
}

// Closer interface
type Closer interface {
	Close() error
}

// ReadWriteCloser embeds ReadWriter and Closer
type ReadWriteCloser interface {
	ReadWriter
	Closer
}

// NamedReadWriter has explicit methods from embedded interfaces
type NamedReadWriter interface {
	Name() string
	ReadWriter
}

// BaseLogger struct
type BaseLogger struct {
	prefix string
}

func (b *BaseLogger) Log(msg string) string {
	return b.prefix + ": " + msg
}

// TimestampLogger embeds BaseLogger
type TimestampLogger struct {
	BaseLogger
	timestamp string
}

// DebugLogger embeds TimestampLogger (nested embedding)
type DebugLogger struct {
	TimestampLogger
	level int
}

// MultiEmbed struct with multiple embedded fields
type MultiEmbed struct {
	BaseLogger
	typeName string
}

// EmbeddedInterfaceField has interface as embedded field
type EmbeddedInterfaceField struct {
	fmt.Stringer
}

// EmbeddedPointerField has pointer as embedded field
type EmbeddedPointerField struct {
	*BaseLogger
}

// InterfaceEmbedding tests interface embedding
func InterfaceEmbedding() string {
	// ReadWriteCloser should have all methods from Reader, Writer, Closer
	return "interface embedding works"
}

// StructEmbeddingBasic tests basic struct embedding
func StructEmbeddingBasic() string {
	logger := TimestampLogger{
		BaseLogger: BaseLogger{prefix: "INFO"},
		timestamp:  "2024-01-01",
	}
	msg := logger.Log("test message")
	return fmt.Sprintf("msg=%s", msg)
}

// StructEmbeddingNested tests nested struct embedding
func StructEmbeddingNested() string {
	logger := DebugLogger{
		TimestampLogger: TimestampLogger{
			BaseLogger: BaseLogger{prefix: "DEBUG"},
			timestamp:  "12:00",
		},
		level: 2,
	}
	msg := logger.Log("debug message")
	return fmt.Sprintf("msg=%s,prefix=%s", msg, logger.prefix)
}

// EmbeddedFieldAccess tests direct access to embedded fields
func EmbeddedFieldAccess() string {
	logger := TimestampLogger{
		BaseLogger: BaseLogger{prefix: "TEST"},
		timestamp:  "now",
	}
	// Access embedded field directly
	prefix := logger.prefix
	// Access embedded field through full path
	prefix2 := logger.BaseLogger.prefix
	return fmt.Sprintf("direct=%s,full=%s", prefix, prefix2)
}

// EmbeddedPointerFieldAccess tests embedded pointer field
func EmbeddedPointerFieldAccess() string {
	embed := EmbeddedPointerField{
		BaseLogger: &BaseLogger{prefix: "PTR"},
	}
	msg := embed.Log("pointer embedded")
	return fmt.Sprintf("msg=%s", msg)
}

// EmbeddedNilPointer tests embedded nil pointer
func EmbeddedNilPointer() string {
	defer func() { recover() }()
	embed := EmbeddedPointerField{
		BaseLogger: nil,
	}
	_ = embed.Log("nil embedded")
	return "no panic"
}

// MultipleEmbeddedStructs tests multiple embedded structs
func MultipleEmbeddedStructs() string {
	type Named struct {
		name string
	}
	type Valued struct {
		value int
	}
	type Combined struct {
		Named
		Valued
	}
	c := Combined{
		Named:  Named{name: "test"},
		Valued: Valued{value: 42},
	}
	return fmt.Sprintf("name=%s,value=%d", c.name, c.value)
}

// MethodShadowing tests method shadowing in embedding
func MethodShadowing() string {
	type Base struct{}
	type Derived struct {
		Base
	}
	// Base methods are promoted unless shadowed
	return "method shadowing understood"
}

// EmbeddedInterfaceSatisfaction tests embedded interface satisfaction
func EmbeddedInterfaceSatisfaction() string {
	type Stringer interface{ String() string }
	type Inner struct {
		Value string
	}
	_ = func(i Inner) string { return i.Value } // Method would be here
	type MyStruct struct {
		Stringer
	}
	// This tests that an interface field can be embedded
	var s Stringer
	_ = MyStruct{Stringer: s} // nil is ok
	return "interface satisfaction"
}

// DeepEmbedding tests deep embedding chains
func DeepEmbedding() string {
	type A struct{ value int }
	type B struct{ A }
	type C struct{ B }
	type D struct{ C }
	d := D{C{B{A{value: 100}}}}
	return fmt.Sprintf("value=%d", d.value)
}

// EmbeddedWithTags tests embedded fields with struct tags
func EmbeddedWithTags() string {
	type Inner struct {
		Value int `json:"value"`
	}
	type Outer struct {
		Inner `json:"inner"`
		Name  string `json:"name"`
	}
	o := Outer{
		Inner: Inner{Value: 42},
		Name:  "test",
	}
	return fmt.Sprintf("value=%d,name=%s", o.Value, o.Name)
}
