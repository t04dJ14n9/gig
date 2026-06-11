package main

import _ "github.com/stretchr/testify/assert"

// TestifyPackageImportable tests that the testify/assert package
// can be imported in the interpreter.
func TestifyPackageImportable() string {
	return "ok"
}
