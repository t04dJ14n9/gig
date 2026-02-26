package tests

import "testing"

func TestExternalFmtSprintf(t *testing.T) {
	runStr(t, `package main
import "fmt"
func Compute() string {
	return fmt.Sprintf("hello %d", 42)
}`, "hello 42")
}

func TestExternalFmtSprintfMulti(t *testing.T) {
	runStr(t, `package main
import "fmt"
func Compute() string {
	return fmt.Sprintf("%s is %d years old", "Alice", 30)
}`, "Alice is 30 years old")
}

func TestExternalStringsToUpper(t *testing.T) {
	runStr(t, `package main
import "strings"
func Compute() string {
	return strings.ToUpper("hello world")
}`, "HELLO WORLD")
}

func TestExternalStringsToLower(t *testing.T) {
	runStr(t, `package main
import "strings"
func Compute() string {
	return strings.ToLower("HELLO")
}`, "hello")
}

func TestExternalStringsContains(t *testing.T) {
	runInt(t, `package main
import "strings"
func Compute() int {
	if strings.Contains("hello world", "world") {
		return 1
	}
	return 0
}`, 1)
}

func TestExternalStringsReplace(t *testing.T) {
	runStr(t, `package main
import "strings"
func Compute() string {
	return strings.ReplaceAll("foo bar foo", "foo", "baz")
}`, "baz bar baz")
}

func TestExternalStringsHasPrefix(t *testing.T) {
	runInt(t, `package main
import "strings"
func Compute() int {
	if strings.HasPrefix("hello world", "hello") {
		return 1
	}
	return 0
}`, 1)
}

func TestExternalStrcovItoa(t *testing.T) {
	runStr(t, `package main
import "strconv"
func Compute() string {
	return strconv.Itoa(42)
}`, "42")
}

func TestExternalStrcovAtoi(t *testing.T) {
	runInt(t, `package main
import "strconv"
func Compute() int {
	n, _ := strconv.Atoi("123")
	return n
}`, 123)
}
