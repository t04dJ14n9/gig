package thirdparty

import "regexp"

// RegexpMatch tests regexp.Match.
func RegexpMatch() int {
	matched, _ := regexp.Match(`hello.*world`, []byte("hello world"))
	if matched {
		return 1
	}
	return 0
}

// RegexpCompile tests regexp.Compile with submatch.
func RegexpCompile() int {
	re := regexp.MustCompile(`(\d+)-(\d+)-(\d+)`)
	matches := re.FindStringSubmatch("2024-03-15")
	return len(matches)
}

// RegexpMustCompile tests regexp.MustCompile.
func RegexpMustCompile() int {
	re := regexp.MustCompile(`^test$`)
	if re.MatchString("test") {
		return 1
	}
	return 0
}

// RegexpFindString tests FindString.
func RegexpFindString() string {
	re := regexp.MustCompile(`f[aeiou]o`)
	return re.FindString("foo bar fao baz")
}

// RegexpFindStringSubmatch tests FindStringSubmatch.
func RegexpFindStringSubmatch() []string {
	re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	return re.FindStringSubmatch("test@example.com")
}

// RegexpFindAllString tests FindAllString.
func RegexpFindAllString() int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString("a1 b2 c3", -1)
	return len(matches)
}

// RegexpReplaceAllString tests ReplaceAllString.
func RegexpReplaceAllString() string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllString("a1 b2 c3", "#")
}

// RegexpSplit tests Split.
func RegexpSplit() int {
	re := regexp.MustCompile(`\s+`)
	parts := re.Split("hello world  test", -1)
	return len(parts)
}

// RegexpNumSubexp tests NumSubexp.
func RegexpNumSubexp() int {
	re := regexp.MustCompile(`(\w+)@(\w+)`)
	return re.NumSubexp()
}

// RegexpLongest tests Longest.
func RegexpLongest() int {
	re := regexp.MustCompile(`a(b|c)?d`)
	re.Longest()
	matches := re.FindStringSubmatch("ad")
	return len(matches)
}
