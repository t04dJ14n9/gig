package tests

import (
	"testing"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

func TestErrorsAsMatchesInterpretedPointerStruct(t *testing.T) {
	const src = `
package main

import "errors"

type errorImpl struct {
	msg string
}

func (e *errorImpl) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.msg
}

func Main() string {
	err := &errorImpl{"test"}
	var target *errorImpl
	if errors.As(err, &target) {
		return target.Error()
	}
	return "no match"
}
`

	prog, err := gig.Build(src)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	defer prog.Close()

	got, err := prog.Run("Main")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got != "test" {
		t.Fatalf("Main() = %v, want test", got)
	}
}
