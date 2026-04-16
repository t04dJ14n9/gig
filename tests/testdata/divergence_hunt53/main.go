package divergence_hunt53

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ============================================================================
// Round 53: Encoding and parsing - JSON, regex, string parsing
// ============================================================================

func JSONEncodeInt() string {
	type Data struct{ Value int `json:"value"` }
	d := Data{Value: 42}
	b, _ := json.Marshal(d)
	return string(b)
}

func JSONEncodeString() string {
	type Data struct{ Value string `json:"value"` }
	d := Data{Value: "hello"}
	b, _ := json.Marshal(d)
	return string(b)
}

func JSONDecodeInt() int {
	data := `{"value":42}`
	type Data struct{ Value int `json:"value"` }
	var d Data
	json.Unmarshal([]byte(data), &d)
	return d.Value
}

func JSONDecodeString() string {
	data := `{"value":"hello"}`
	type Data struct{ Value string `json:"value"` }
	var d Data
	json.Unmarshal([]byte(data), &d)
	return d.Value
}

func JSONSlice() int {
	data := `[1,2,3]`
	var s []int
	json.Unmarshal([]byte(data), &s)
	return s[0] + s[1] + s[2]
}

func JSONMap() int {
	data := `{"a":1,"b":2}`
	var m map[string]int
	json.Unmarshal([]byte(data), &m)
	return m["a"] + m["b"]
}

func JSONBool() bool {
	data := `true`
	var b bool
	json.Unmarshal([]byte(data), &b)
	return b
}

func JSONNull() string {
	data := `null`
	var v any
	json.Unmarshal([]byte(data), &v)
	if v == nil { return "null" }
	return "not null"
}

func RegexMatch() bool {
	ok, _ := regexp.MatchString(`^\d+$`, "12345")
	return ok
}

func RegexFind() string {
	re := regexp.MustCompile(`\d+`)
	return re.FindString("abc123def456")
}

func RegexFindAll() int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString("a1b22c333", -1)
	return len(matches)
}

func RegexReplace() string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllString("abc123def456", "NUM")
}

func RegexSplit() int {
	re := regexp.MustCompile(`\s+`)
	parts := re.Split("hello   world  foo", -1)
	return len(parts)
}

func RegexSubmatch() int {
	re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	match := re.FindStringSubmatch("user@example.com")
	return len(match)
}

func RegexNamedGroup() string {
	re := regexp.MustCompile(`(?P<first>\w+)\s+(?P<last>\w+)`)
	match := re.FindStringSubmatch("John Doe")
	return fmt.Sprintf("%s:%s", match[1], match[2])
}

func StringParse() int {
	line := "name=Alice,age=30,score=95"
	parts := strings.Split(line, ",")
	count := 0
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		_ = kv[0]
		_ = kv[1]
		count++
	}
	return count
}

func CSVParse() int {
	line := `Alice,30,"New York, NY",Engineer`
	// Simplified CSV parsing (no quote handling)
	parts := strings.Split(line, ",")
	return len(parts)
}

func TemplateParse() string {
	template := "${name} is ${age} years old"
	s := strings.ReplaceAll(template, "${name}", "Bob")
	s = strings.ReplaceAll(s, "${age}", "25")
	return s
}
