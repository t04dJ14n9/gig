// Package register provides the public API for registering external packages.
//
// This package provides a simplified API for registering external packages
// compared to the lower-level importer package. It wraps importer.ExternalPackage
// to provide a cleaner interface for users.
//
// # Example Usage
//
//	pkg := register.AddPackage("mypkg", "mypkg")
//	pkg.NewFunction("Process", ProcessFunc, "Process data")
//	pkg.NewVar("Config", &config, "Global configuration")
//	pkg.NewConst("Version", "1.0.0", "Package version")
package register

import (
	"gig/importer"
	"gig/value"
)

// AddPackage registers a new external package with the given import path and name.
// Returns an ExternalPackage for adding functions, variables, constants, and types.
func AddPackage(path, name string) *ExternalPackage {
	return &ExternalPackage{importer.RegisterPackage(path, name)}
}

// ExternalPackage wraps importer.ExternalPackage to provide a public API.
type ExternalPackage struct {
	inner *importer.ExternalPackage
}

// NewFunction adds a function to the package.
// The fn parameter must be a function value.
// The doc parameter is optional documentation.
func (p *ExternalPackage) NewFunction(name string, fn any, doc string) {
	p.inner.AddFunction(name, fn, doc, nil)
}

// NewFunctionDirect adds a function with a direct-call wrapper for fast dispatch.
// The directCall function should convert Value arguments to native types, call the
// function, and wrap the result, avoiding the overhead of reflect.Call.
func (p *ExternalPackage) NewFunctionDirect(name string, fn any, doc string, directCall func([]value.Value) value.Value) {
	p.inner.AddFunction(name, fn, doc, directCall)
}

// NewVar adds a mutable variable to the package.
// The ptr parameter must be a pointer to the variable.
func (p *ExternalPackage) NewVar(name string, ptr any, doc string) {
	p.inner.AddVariable(name, ptr, doc)
}

// NewConst adds an immutable constant to the package.
// The val parameter is the constant value.
func (p *ExternalPackage) NewConst(name string, val any, doc string) {
	p.inner.AddConstant(name, val, doc)
}

// NewType adds a named type to the package.
// The typ parameter should be a reflect.Type or a value of the type.
func (p *ExternalPackage) NewType(name string, typ any, doc string) {
	// typ should be a reflect.Type or a value of the type
	p.inner.AddType(name, nil, doc) // Will be overridden by type registration
}

// GetPackageByPath returns a registered package by its import path.
// Returns nil if no package with the given path is registered.
func GetPackageByPath(path string) *importer.ExternalPackage {
	return importer.GetPackageByPath(path)
}

// GetPackageByName returns a registered package by its name.
// Returns nil if no package with the given name is registered.
func GetPackageByName(name string) *importer.ExternalPackage {
	return importer.GetPackageByName(name)
}

// GetAllPackages returns all registered packages, keyed by import path.
func GetAllPackages() map[string]*importer.ExternalPackage {
	return importer.GetAllPackages()
}
