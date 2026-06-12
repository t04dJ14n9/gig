package divergence_hunt233

import "fmt"

// ============================================================================
// Round 233: String immutability
// ============================================================================

func StringCannotBeModified() string {
	// Strings are immutable in Go - assignment to string index is a compile error
	// String modification requires conversion to []byte or []rune
	return "immutable"
}

func StringReassignmentIsAllowed() string {
	s := "first"
	s = "second"
	s = "third"
	return s
}

func StringConcatenationCreatesNew() string {
	s1 := "hello"
	s2 := s1
	s1 += " world"
	return fmt.Sprintf("s1=%s,s2=%s", s1, s2)
}

func ByteSliceModifiable() string {
	b := []byte("hello")
	b[0] = 'H'
	return string(b)
}

func StringToByteSliceCopy() string {
	s := "immutable"
	b := []byte(s)
	b[0] = 'X'
	return fmt.Sprintf("original=%s,modified=%s", s, string(b))
}

func StringPassedByValue() string {
	s := "original"
	modify := func(str string) {
		str = "modified"
	}
	modify(s)
	return s
}

func StringInStructImmutable() string {
	type Person struct {
		Name string
	}
	p := Person{Name: "Alice"}
	original := p.Name
	defer func() {
		recover()
	}()
	_ = original
	return "immutable"
}

func StringAppendDoesNotModify() string {
	s := "base"
	_ = s + "suffix"
	return s
}

func StringConstantImmutability() string {
	const s = "constant"
	_ = s
	return "constant_unchanged"
}

func StringSliceRemainsImmutable() string {
	s := "hello world"
	sub := s[0:5]
	_ = sub
	return "immutable"
}

func MultipleReferencesSameString() string {
	s1 := "shared"
	s2 := s1
	s3 := s1
	return fmt.Sprintf("%s:%s:%s", s1, s2, s3)
}
