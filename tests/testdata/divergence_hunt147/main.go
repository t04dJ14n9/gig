package divergence_hunt147

import (
	"fmt"
	"strings"
)

// ============================================================================
// Round 147: String manipulation and conversion patterns
// ============================================================================

func StringsJoin() string {
	parts := []string{"hello", "world", "test"}
	return strings.Join(parts, "-")
}

func StringsSplit() string {
	s := "a,b,c"
	parts := strings.Split(s, ",")
	return fmt.Sprintf("%v", parts)
}

func StringsContains() string {
	s := "hello world"
	if strings.Contains(s, "world") {
		return "found"
	}
	return "not-found"
}

func StringsHasPrefix() string {
	s := "golang"
	if strings.HasPrefix(s, "go") {
		return "yes"
	}
	return "no"
}

func StringsHasSuffix() string {
	s := "test.go"
	if strings.HasSuffix(s, ".go") {
		return "go-file"
	}
	return "other"
}

func StringsTrimSpace() string {
	s := "  hello  "
	return fmt.Sprintf("[%s]", strings.TrimSpace(s))
}

func StringsReplace() string {
	s := "hello world"
	result := strings.Replace(s, "world", "go", 1)
	return result
}

func StringsToUpper() string {
	return strings.ToUpper("hello")
}

func StringsRepeat() string {
	return strings.Repeat("ab", 3)
}

func StringsCount() string {
	s := "mississippi"
	return fmt.Sprintf("count=%d", strings.Count(s, "s"))
}

func StringsIndex() string {
	s := "hello world"
	idx := strings.Index(s, "world")
	return fmt.Sprintf("idx=%d", idx)
}
