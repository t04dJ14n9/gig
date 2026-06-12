package repl

import (
	"fmt"
	"go/parser"
	"go/token"
	"strconv"

	"github.com/t04dJ14n9/gig"
)

// handleImport processes import statements.
func (s *Session) handleImport(input string) {
	fset := token.NewFileSet()
	testSource := "package main\n" + input
	file, err := parser.ParseFile(fset, "test.go", testSource, parser.ImportsOnly)
	if err != nil {
		fmt.Printf("Error parsing import: %v\n", err)
		return
	}

	for _, imp := range file.Imports {
		path, _ := strconv.Unquote(imp.Path.Value)

		if gig.GetPackageByPath(path) != nil {
			s.imports[path] = true
			fmt.Printf("Imported: %s\n", path)
			continue
		}

		if s.pluginMgr == nil {
			s.pluginMgr = newPluginManager()
		}

		if err := s.pluginMgr.LoadPackage(path); err != nil {
			fmt.Printf("Error loading package %s: %v\n", path, err)
			continue
		}

		s.imports[path] = true
		fmt.Printf("Imported: %s (hot-loaded)\n", path)
	}
}

// pkgNameFromPath returns the package name for a given import path.
func pkgNameFromPath(path string) (string, bool) {
	for name, p := range knownPackages {
		if p == path {
			return name, true
		}
	}
	if pkg := gig.GetPackageByPath(path); pkg != nil {
		return pkg.Name, true
	}
	return "", false
}
