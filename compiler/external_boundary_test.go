package compiler

import (
	"strings"
	"testing"

	"github.com/t04dJ14n9/gig/importer"
)

func TestBuildRejectsInterpretedStructPassedToThirdPartyInterface(t *testing.T) {
	reg := importer.NewRegistry()
	pkg := reg.RegisterPackage("example.com/thirdparty", "thirdparty")
	pkg.AddFunction("Accept", func(any) bool { return true }, "", nil)

	const src = `
package main

import "example.com/thirdparty"

type Rule struct {
	Name string
}

func Main() bool {
	return thirdparty.Accept(Rule{Name: "x"})
}
`

	_, err := Build(src, reg)
	if err == nil {
		t.Fatal("Build succeeded, want external boundary error")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("Build error = %v, want interpreter-defined type boundary error", err)
	}
}
