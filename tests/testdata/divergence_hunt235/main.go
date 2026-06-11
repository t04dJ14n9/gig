package divergence_hunt235

import "fmt"

// ============================================================================
// Round 235: Byte slice to string conversions
// ============================================================================

func ByteSliceToString() string {
	b := []byte{'h', 'e', 'l', 'l', 'o'}
	return string(b)
}

func StringToByteSlice() string {
	s := "world"
	b := []byte(s)
	return fmt.Sprintf("%d:%d:%d", len(b), b[0], b[len(b)-1])
}

func EmptyByteSliceToString() string {
	b := []byte{}
	return string(b)
}

func NilByteSliceToString() string {
	var b []byte
	return fmt.Sprintf("len=%d,str='%s'", len(b), string(b))
}

func ByteSliceModification() string {
	b := []byte("hello")
	b[0] = 'H'
	return string(b)
}

func ByteSliceFromStringCopy() string {
	s := "original"
	b := []byte(s)
	b[0] = 'X'
	return fmt.Sprintf("s=%s,b=%s", s, string(b))
}

func RuneSliceToString() string {
	r := []rune{'H', 'e', 'l', 'l', 'o'}
	return string(r)
}

func StringToRuneSlice() string {
	s := "Hello"
	r := []rune(s)
	return fmt.Sprintf("len=%d,first=%c", len(r), r[0])
}

func UnicodeByteSlice() string {
	b := []byte("世界")
	return fmt.Sprintf("bytes=%d,string=%s", len(b), string(b))
}

func ByteSliceWithNull() string {
	b := []byte{'a', 0, 'b'}
	return fmt.Sprintf("len=%d,str=%s", len(b), string(b))
}

func ConvertAndBack() string {
	original := "test"
	b := []byte(original)
	back := string(b)
	return fmt.Sprintf("equal=%v", original == back)
}
