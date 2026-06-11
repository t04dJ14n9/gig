// Package gentool provides the core dependency generation logic for gig.
//
// Gentool generates Go source files that register external packages for the interpreter.
// For each imported package, it creates wrapper code that:
//   - Registers the package with the importer
//   - Adds functions with optional DirectCall wrappers for fast dispatch
//   - Adds variables, constants, and types
//
// # Generated Code Structure
//
// For each package (e.g., "fmt"), gentool generates a file like:
//
//	package packages
//
//	import (
//	    "fmt"
//	    "github.com/t04dJ14n9/gig/importer"
//	    "github.com/t04dJ14n9/gig/model/value"
//	)
//
//	func init() {
//	    pkg := importer.RegisterPackage("fmt", "fmt")
//	    pkg.AddFunction("Sprintf", fmt.Sprintf, "", direct_fmt_Sprintf)
//	    // ... more functions
//	}
//
//	func direct_fmt_Sprintf(args []value.Value) value.Value {
//	    format := args[0].String()
//	    // ... convert arguments and call fmt.Sprintf
//	    return value.FromInterface(result)
//	}
//
// # DirectCall Wrappers
//
// DirectCall wrappers provide fast function dispatch by:
//  1. Extracting typed values from value.Value arguments
//  2. Calling the native Go function directly
//  3. Wrapping the result in a value.Value
//
// This avoids the overhead of reflect.Call for common cases.
// DirectCall wrappers are only generated for functions with supported parameter types
// (basic types, slices of basic types, and empty interfaces).
package gentool

import (
	"fmt"
	"go/importer"
	"go/token"
)

// PackageImport generates registration code for a single package.
// It reads the package's exported symbols and generates:
//   - Package registration in an init() function
//   - Function, variable, constant, and type registrations
//   - DirectCall wrappers for eligible functions
//
// outDir is the output directory for generated files.
// pkgName is the Go package name for generated files (typically "packages").
func PackageImport(path string, outDir string, pkgName string) error {
	pkg, err := importer.ForCompiler(token.NewFileSet(), "source", nil).Import(path)
	if err != nil {
		return err
	}

	currentPkgPath = path

	refs := newPackageRefs(path, pkg.Name())
	symbols := collectPackageSymbols(pkg.Scope(), refs.PkgRef)
	if symbols.empty() {
		fmt.Printf("  (skipped, nothing to register)\n")
		return nil
	}

	plan := buildPackageGenerationPlan(path, pkgName, pkg.Name(), refs, symbols)
	return writeGeneratedPackage(outDir, plan)
}
