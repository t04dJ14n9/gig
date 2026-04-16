package divergence_hunt41

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ============================================================================
// Round 41: Real-world patterns - config, parsing, data transformation
// ============================================================================

func ConfigParsing() int {
	config := map[string]string{"port": "8080", "timeout": "30", "retries": "3"}
	port, _ := strconv.Atoi(config["port"])
	timeout, _ := strconv.Atoi(config["timeout"])
	retries, _ := strconv.Atoi(config["retries"])
	return port + timeout + retries
}

func CSVLineParse() int {
	line := "Alice,30,Engineer"
	parts := strings.Split(line, ",")
	return len(parts)
}

func TemplateSubstitution() string {
	template := "Hello, {name}! You are {age} years old."
	s := strings.ReplaceAll(template, "{name}", "Alice")
	s = strings.ReplaceAll(s, "{age}", "30")
	return s
}

func URLParse() string {
	url := "https://example.com/path?q=hello&r=world"
	qIdx := strings.Index(url, "?")
	if qIdx < 0 { return "" }
	query := url[qIdx+1:]
	parts := strings.Split(query, "&")
	return fmt.Sprintf("%d", len(parts))
}

func DataTransform() int {
	input := []map[string]any{
		{"name": "Alice", "score": 85},
		{"name": "Bob", "score": 92},
		{"name": "Charlie", "score": 78},
	}
	total := 0
	for _, item := range input {
		if score, ok := item["score"].(int); ok {
			total += score
		}
	}
	return total
}

func JSONConfigParse() int {
	data := `{"max_conn":100,"timeout":30,"debug":true}`
	var config map[string]any
	json.Unmarshal([]byte(data), &config)
	maxConn := int(config["max_conn"].(float64))
	return maxConn
}

func StringTemplateBuilder() string {
	var b strings.Builder
	b.WriteString("SELECT * FROM users")
	b.WriteString(" WHERE age > ")
	b.WriteString(strconv.Itoa(18))
	b.WriteString(" AND active = true")
	return b.String()
}

func NumberFormatter() string {
	format := func(n int) string {
		s := strconv.Itoa(n)
		if n < 0 {
			s = "(" + strconv.Itoa(-n) + ")"
		}
		return s
	}
	return format(42) + "," + format(-7)
}

func MapReducePattern() int {
	data := []int{1, 2, 3, 4, 5}
	// Map: double each
	doubled := make([]int, len(data))
	for i, v := range data { doubled[i] = v * 2 }
	// Reduce: sum
	sum := 0
	for _, v := range doubled { sum += v }
	return sum
}

func PipelinePattern() string {
	process := func(s string) string {
		return strings.TrimSpace(strings.ToLower(s))
	}
	return process("  HELLO WORLD  ")
}

func ErrorChainPattern() int {
	validate := func(s string) error {
		if s == "" { return fmt.Errorf("empty") }
		return nil
	}
	transform := func(s string) (string, error) {
		if err := validate(s); err != nil { return "", err }
		return strings.ToUpper(s), nil
	}
	result, err := transform("hello")
	if err != nil { return -1 }
	return len(result)
}

func BuilderPattern() string {
	type Query struct {
		Table   string
		Where   string
		OrderBy string
	}
	q := &Query{Table: "users"}
	q.Where = "age > 18"
	q.OrderBy = "name"
	return fmt.Sprintf("%s|%s|%s", q.Table, q.Where, q.OrderBy)
}

func RateLimiterPattern() int {
	max := 5
	count := 0
	for i := 0; i < 10; i++ {
		if count < max {
			count++
		}
	}
	return count
}

func RetryPattern() int {
	attempts := 3
	result := 0
	for i := 0; i < attempts; i++ {
		if i == 2 {
			result = 42
			break
		}
	}
	return result
}
