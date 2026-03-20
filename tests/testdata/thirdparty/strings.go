package thirdparty

import "strings"

// StringsBuilder tests strings.Builder.
func StringsBuilder() int {
	var builder strings.Builder
	builder.WriteString("hello")
	builder.WriteString(" ")
	builder.WriteString("world")
	return builder.Len()
}

// StringsBuilderString tests Builder.String().
func StringsBuilderString() string {
	var builder strings.Builder
	builder.WriteString("test")
	return builder.String()
}

// StringsBuilderGrow tests Builder.Grow().
func StringsBuilderGrow() int {
	var builder strings.Builder
	builder.Grow(100)
	builder.WriteString("x")
	return builder.Cap()
}

// StringsMap tests strings.Map.
func StringsMap() int {
	result := strings.Map(func(r rune) rune {
		return r + 1
	}, "abcd")
	if result == "bcde" {
		return 1
	}
	return 0
}

// StringsRepeat tests strings.Repeat.
func StringsRepeat() int {
	if strings.Repeat("ab", 3) == "ababab" {
		return 1
	}
	return 0
}

// StringsRepeatCount returns the length of a repeated string.
func StringsRepeatCount() int {
	return len(strings.Repeat("x", 100))
}

// StringsIndexAny tests strings.IndexAny.
func StringsIndexAny() int {
	return strings.IndexAny("hello", "aeiou")
}

// StringsIndexFunc tests strings.IndexFunc with a closure.
func StringsIndexFunc() int {
	f := func(r rune) bool {
		return r == 'x'
	}
	return strings.IndexFunc("abcxdef", f)
}

// StringsTitle tests strings.Title.
func StringsTitle() int {
	if strings.Title("hello world") == "Hello World" {
		return 1
	}
	return 0
}

// StringsToTitle tests strings.ToTitle.
func StringsToTitle() int {
	if strings.ToTitle("hello") == "HELLO" {
		return 1
	}
	return 0
}

// StringsToValidUTF8 tests strings.ToValidUTF8.
func StringsToValidUTF8() int {
	result := strings.ToValidUTF8("hello\xC0\xC1world", "?")
	count := strings.Count(result, "?")
	if count == 2 {
		return 1
	}
	return 0
}

// StringsTrimLeft tests strings.TrimLeft.
func StringsTrimLeft() int {
	if strings.TrimLeft("xxhelloxx", "x") == "helloxx" {
		return 1
	}
	return 0
}

// StringsTrimRight tests strings.TrimRight.
func StringsTrimRight() int {
	if strings.TrimRight("xxhelloxx", "x") == "xxhello" {
		return 1
	}
	return 0
}

// StringsTrimFunc tests strings.TrimFunc with a closure.
func StringsTrimFunc() int {
	if strings.TrimFunc("  hello  ", func(r rune) bool {
		return r == ' '
	}) == "hello" {
		return 1
	}
	return 0
}

// StringsIndexFuncTest tests strings.IndexFunc with digit detection.
func StringsIndexFuncTest() int {
	f := func(r rune) bool {
		return r >= '0' && r <= '9'
	}
	return strings.IndexFunc("abc123def", f)
}

// StringsCut tests strings.Cut.
func StringsCut() int {
	before, after, found := strings.Cut("hello world", " ")
	if found && before == "hello" && after == "world" {
		return 1
	}
	return 0
}

// StringsCutPrefix tests strings.TrimPrefix.
func StringsCutPrefix() int {
	s := strings.TrimPrefix("hello world", "hello ")
	if s == "world" {
		return 1
	}
	return 0
}

// StringsCutSuffix tests strings.TrimSuffix.
func StringsCutSuffix() int {
	s := strings.TrimSuffix("hello world", " world")
	if s == "hello" {
		return 1
	}
	return 0
}
