package gentool

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

// ParsePkgsFile parses a Go source file and extracts all import paths.
// The file should contain blank imports (_ "package/path") for packages to register.
// Returns a list of import paths or an error if parsing fails.
func ParsePkgsFile(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", filePath, err)
	}

	var paths []string
	for _, imp := range f.Imports {
		p := strings.Trim(imp.Path.Value, `"`)
		if p == "" {
			continue
		}
		paths = append(paths, p)
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no imports found in %s", filePath)
	}
	return paths, nil
}
