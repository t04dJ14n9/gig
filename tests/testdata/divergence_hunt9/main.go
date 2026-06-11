package divergence_hunt9

import (
	"encoding/json"
	"math"
	"regexp"
)

// ============================================================================
// Round 9: JSON, regex, more math, encoding, parsing
// ============================================================================

// JSONMarshal tests JSON marshaling
func JSONMarshal() string {
	type Point struct{ X, Y int }
	p := Point{X: 1, Y: 2}
	b, _ := json.Marshal(p)
	return string(b)
}

// JSONUnmarshal tests JSON unmarshaling
func JSONUnmarshal() int {
	type Point struct{ X, Y int }
	data := `{"X":10,"Y":20}`
	var p Point
	json.Unmarshal([]byte(data), &p)
	return p.X + p.Y
}

// JSONMarshalMap tests JSON with map
func JSONMarshalMap() string {
	m := map[string]int{"a": 1, "b": 2}
	b, _ := json.Marshal(m)
	return string(b)
}

// RegexMatch tests regexp.MatchString
func RegexMatch() bool {
	ok, _ := regexp.MatchString(`^hello\s+world$`, "hello world")
	return ok
}

// RegexFind tests regexp.FindString
func RegexFind() string {
	re := regexp.MustCompile(`\d+`)
	return re.FindString("abc123def")
}

// RegexFindAll tests regexp.FindAllString
func RegexFindAll() int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString("a1b22c333", -1)
	return len(matches)
}

// RegexReplace tests regexp.ReplaceAllString
func RegexReplace() string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllString("abc123def456", "NUM")
}

// MathMod tests math.Mod
func MathMod() float64 { return math.Mod(7.5, 2.5) }

// MathLog tests math.Log
func MathLog() float64 { return math.Log(math.E) }

// MathExp tests math.Exp
func MathExp() float64 { return math.Exp(1) }

// MathRound tests math.Round
func MathRound() float64 {
	return math.Round(3.5) + math.Round(3.4)
}

// MathTrunc tests math.Trunc
func MathTrunc() float64 { return math.Trunc(3.7) }

// MathRemainder tests math.Remainder
func MathRemainder() float64 { return math.Remainder(7, 3) }

// MathCopysign tests math.Copysign
func MathCopysign() float64 { return math.Copysign(3, -1) }

// JSONMarshalSlice tests JSON with slice
func JSONMarshalSlice() string {
	s := []int{1, 2, 3}
	b, _ := json.Marshal(s)
	return string(b)
}

// JSONUnmarshalSlice tests JSON unmarshal into slice
func JSONUnmarshalSlice() int {
	data := `[1,2,3]`
	var s []int
	json.Unmarshal([]byte(data), &s)
	return s[0] + s[1] + s[2]
}

// RegexSplit tests regexp.Split
func RegexSplit() int {
	re := regexp.MustCompile(`\s+`)
	parts := re.Split("hello   world  foo", -1)
	return len(parts)
}

// RegexSubmatch tests regexp.FindStringSubmatch
func RegexSubmatch() int {
	re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	match := re.FindStringSubmatch("user@example.com")
	return len(match)
}

// MathHypot tests math.Hypot
func MathHypot() float64 { return math.Hypot(3, 4) }

// MathPow10 tests math.Pow10
func MathPow10() float64 { return math.Pow10(3) }

// MathSignbit tests math.Signbit
func MathSignbit() bool { return math.Signbit(-1.0) && !math.Signbit(1.0) }
