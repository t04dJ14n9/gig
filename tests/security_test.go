package tests

import "testing"

func TestSecurityBanUnsafe(t *testing.T) {
	expectBuildError(t, `package main
import "unsafe"
func Compute() int { return 42 }`)
}

func TestSecurityBanReflect(t *testing.T) {
	expectBuildError(t, `package main
import "reflect"
func Compute() int { return 42 }`)
}
