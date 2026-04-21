package divergence_hunt231

import "fmt"

// ============================================================================
// Round 231: String indexing and byte access
// ============================================================================

func StringByteIndex() string {
	s := "Hello"
	return fmt.Sprintf("%d:%d:%d", s[0], s[1], s[4])
}

func StringByteIndexUnicode() string {
	s := "Hello, 世界"
	return fmt.Sprintf("%d:%d:%d", s[0], s[7], s[8])
}

func StringLengthBytes() string {
	s1 := "abc"
	s2 := "你好"
	return fmt.Sprintf("%d:%d", len(s1), len(s2))
}

func StringIndexOutOfBounds() string {
	result := "ok"
	defer func() {
		if r := recover(); r != nil {
			result = "panicked"
		}
	}()
	s := "hi"
	_ = s[10]
	return result
}

func StringByteLoop() string {
	s := "ABC"
	sum := 0
	for i := 0; i < len(s); i++ {
		sum += int(s[i])
	}
	return fmt.Sprintf("%d", sum)
}

func StringEmptyIndex() string {
	result := "ok"
	defer func() {
		if r := recover(); r != nil {
			result = "panicked"
		}
	}()
	s := ""
	_ = s[0]
	return result
}

func StringBackwardsIndex() string {
	s := "GoLang"
	return fmt.Sprintf("%c:%c", s[len(s)-1], s[len(s)-2])
}

func ByteSliceFromString() string {
	s := "test"
	b := []byte(s)
	return fmt.Sprintf("%d:%d:%d:%d", b[0], b[1], b[2], b[3])
}

func StringByteAssignment() string {
	// Strings are immutable in Go - assignment to string index is a compile error
	// This documents that strings cannot be modified byte-by-byte
	return "immutable"
}

func StringIndexingWithVariables() string {
	s := "programming"
	i, j := 3, 7
	return fmt.Sprintf("%c:%c", s[i], s[j])
}

func StringFirstLastByte() string {
	s := "boundary"
	return fmt.Sprintf("%d:%d", s[0], s[len(s)-1])
}
