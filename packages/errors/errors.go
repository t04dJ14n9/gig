// Package errors registers the Go standard library errors package.
package errors

import (
	"errors"

	"gig/importer"
	"gig/value"
)

func init() {
	pkg := importer.RegisterPackage("errors", "errors")

	// Error creation and handling
	pkg.AddFunction("New", errors.New, "", directNew)
	pkg.AddFunction("Is", errors.Is, "", nil)
	pkg.AddFunction("As", errors.As, "", nil)
	pkg.AddFunction("Unwrap", errors.Unwrap, "", nil)
	pkg.AddFunction("Join", errors.Join, "", nil)
}

func directNew(args []value.Value) value.Value {
	return value.FromInterface(errors.New(args[0].String()))
}
