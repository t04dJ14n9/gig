// Package register provides the public API for registering external packages.
package register

import (
	"gig/importer"
	"gig/value"
)

// AddPackage registers a new external package with the given path and name.
// Returns an ExternalPackage for adding objects.
func AddPackage(path, name string) *ExternalPackage {
	return &ExternalPackage{importer.RegisterPackage(path, name)}
}

// ExternalPackage wraps importer.ExternalPackage for public API.
type ExternalPackage struct {
	inner *importer.ExternalPackage
}

// NewFunction adds a function to the package.
// fn must be a function value.
// doc is optional documentation.
func (p *ExternalPackage) NewFunction(name string, fn any, doc string) {
	p.inner.AddFunction(name, fn, doc, nil)
}

// NewFunctionDirect adds a function with a direct-call wrapper.
// The directCall function should implement the function without using reflect.Call.
func (p *ExternalPackage) NewFunctionDirect(name string, fn any, doc string, directCall func([]value.Value) value.Value) {
	p.inner.AddFunction(name, fn, doc, directCall)
}

// NewVar adds a variable to the package.
// ptr must be a pointer to the variable.
func (p *ExternalPackage) NewVar(name string, ptr any, doc string) {
	p.inner.AddVariable(name, ptr, doc)
}

// NewConst adds a constant to the package.
func (p *ExternalPackage) NewConst(name string, val any, doc string) {
	p.inner.AddConstant(name, val, doc)
}

// NewType adds a type to the package.
func (p *ExternalPackage) NewType(name string, typ any, doc string) {
	// typ should be a reflect.Type or a value of the type
	p.inner.AddType(name, nil, doc) // Will be overridden by type registration
}

// GetPackageByPath returns a registered package by its import path.
func GetPackageByPath(path string) *importer.ExternalPackage {
	return importer.GetPackageByPath(path)
}

// GetPackageByName returns a registered package by its name.
func GetPackageByName(name string) *importer.ExternalPackage {
	return importer.GetPackageByName(name)
}

// GetAllPackages returns all registered packages.
func GetAllPackages() map[string]*importer.ExternalPackage {
	return importer.GetAllPackages()
}
