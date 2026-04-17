package divergence_hunt112

import (
	"fmt"
	"regexp"
)

// ============================================================================
// Round 112: Regexp matching and replacement
// ============================================================================

func RegexpMatch() string {
	matched, _ := regexp.MatchString(`^\d+$`, "12345")
	return fmt.Sprintf("%v", matched)
}

func RegexpMatchFail() string {
	matched, _ := regexp.MatchString(`^\d+$`, "abc")
	return fmt.Sprintf("%v", matched)
}

func RegexpFindString() string {
	re := regexp.MustCompile(`\d+`)
	return re.FindString("abc123def")
}

func RegexpFindAllString() string {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString("a1b23c456", -1)
	return fmt.Sprintf("%v", matches)
}

func RegexpReplaceAllString() string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllString("a1b23c456", "X")
}

func RegexpSplit() string {
	re := regexp.MustCompile(`[,;]`)
	parts := re.Split("a,b;c,d", -1)
	return fmt.Sprintf("%v", parts)
}

func RegexpSubmatch() string {
	re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	matches := re.FindStringSubmatch("user@example.com")
	return fmt.Sprintf("%d", len(matches))
}

func RegexpReplaceAllStringFunc() string {
	re := regexp.MustCompile(`\d+`)
	result := re.ReplaceAllStringFunc("a1b23c", func(s string) string {
		return fmt.Sprintf("[%s]", s)
	})
	return result
}

func RegexpFindStringIndex() string {
	re := regexp.MustCompile(`\d+`)
	loc := re.FindStringIndex("abc123def")
	if loc != nil {
		return fmt.Sprintf("%d:%d", loc[0], loc[1])
	}
	return "not found"
}

func RegexpCompileMust() string {
	re := regexp.MustCompile(`hello.*world`)
	return fmt.Sprintf("%v", re.MatchString("hello beautiful world"))
}
