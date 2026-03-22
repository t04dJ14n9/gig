package importer

// Package importer provides package registration, lookup, and import resolution
// for the Gig interpreter. External packages register functions, types, variables,
// and method DirectCall wrappers via the PackageRegistry interface.
//
// The global convenience functions (RegisterPackage, GetPackageByPath, etc.) delegate
// to a default global Registry instance, which is pre-populated by init() functions
// in generated package wrappers (e.g., stdlib/packages).
