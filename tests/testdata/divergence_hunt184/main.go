package divergence_hunt184

import (
	"fmt"
	"regexp"
)

// ============================================================================
// Round 184: Regexp basic patterns
// ============================================================================

func MatchString() string {
	matched, _ := regexp.MatchString("^[a-z]+$", "hello")
	return fmt.Sprintf("%v", matched)
}

func MatchStringDigit() string {
	matched, _ := regexp.MatchString(`^\d+$`, "12345")
	return fmt.Sprintf("%v", matched)
}

func MatchStringNotMatch() string {
	matched, _ := regexp.MatchString("^[a-z]+$", "Hello123")
	return fmt.Sprintf("%v", matched)
}

func MatchStringWord() string {
	matched, _ := regexp.MatchString(`\w+`, "test_word")
	return fmt.Sprintf("%v", matched)
}

func CompileAndMatch() string {
	re := regexp.MustCompile(`foo.?`)
	matched := re.MatchString("seafood")
	return fmt.Sprintf("%v", matched)
}

func FindString() string {
	re := regexp.MustCompile(`foo.?`)
	found := re.FindString("seafood fool")
	return fmt.Sprintf("%s", found)
}

func FindAllString() string {
	re := regexp.MustCompile(`a.`)
	found := re.FindAllString("paranormal", -1)
	return fmt.Sprintf("%d", len(found))
}

func ReplaceAllString() string {
	re := regexp.MustCompile(`a(x*)b`)
	result := re.ReplaceAllString("-ab-axxb-", "T")
	return fmt.Sprintf("%s", result)
}

func SplitString() string {
	re := regexp.MustCompile(`[.,]`)
	parts := re.Split("a.b,c", -1)
	return fmt.Sprintf("%d", len(parts))
}

func QuoteMeta() string {
	quoted := regexp.QuoteMeta(`[foo]`)
	return fmt.Sprintf("%s", quoted)
}

func FindStringIndex() string {
	re := regexp.MustCompile(`ab*`)
	loc := re.FindStringIndex("cabbbc")
	if loc == nil {
		return fmt.Sprintf("nil")
	}
	return fmt.Sprintf("%d:%d", loc[0], loc[1])
}
