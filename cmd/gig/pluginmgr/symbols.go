package pluginmgr

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"unicode"
)

// ExportedSymbols holds discovered exported symbols from a package.
type ExportedSymbols struct {
	Funcs  []string
	Consts []string
	Vars   []string
	Types  []string
}

// getExportedSymbols uses go doc to discover exported package symbols.
func (pm *Manager) getExportedSymbols(pkgPath string) (*ExportedSymbols, error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "go", "doc", "-short", pkgPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("go doc failed: %w", err)
	}

	symbols := &ExportedSymbols{}
	for _, line := range strings.Split(string(output), "\n") {
		addExportedSymbol(symbols, strings.TrimSpace(line))
	}

	if len(symbols.Funcs) == 0 && len(symbols.Types) == 0 && len(symbols.Consts) == 0 && len(symbols.Vars) == 0 {
		return nil, fmt.Errorf("no symbols discovered")
	}
	return symbols, nil
}

// GetSymbols returns the exported symbols for a package path.
func (pm *Manager) GetSymbols(pkgPath string) []string {
	pm.mu.RLock()
	symbols, ok := pm.symbols[pkgPath]
	pm.mu.RUnlock()

	if !ok || symbols == nil {
		return nil
	}

	var result []string
	result = append(result, symbols.Funcs...)
	result = append(result, symbols.Types...)
	result = append(result, symbols.Consts...)
	result = append(result, symbols.Vars...)
	return result
}

func addExportedSymbol(symbols *ExportedSymbols, line string) {
	if line == "" || strings.HasPrefix(line, "//") {
		return
	}
	switch {
	case strings.HasPrefix(line, "func "):
		if name, ok := parseFunctionName(line); ok && isExported(name) {
			symbols.Funcs = append(symbols.Funcs, name)
		}
	case strings.HasPrefix(line, "type "):
		if name, ok := parseTypeName(line); ok && isExported(name) {
			symbols.Types = append(symbols.Types, name)
		}
	case strings.HasPrefix(line, "const "):
		if name, ok := parseFirstDeclaredName(line, "const "); ok && isExported(name) {
			symbols.Consts = append(symbols.Consts, name)
		}
	case strings.HasPrefix(line, "var "):
		if name, ok := parseFirstDeclaredName(line, "var "); ok && isExported(name) {
			symbols.Vars = append(symbols.Vars, name)
		}
	}
}

func parseFunctionName(line string) (string, bool) {
	rest := strings.TrimPrefix(line, "func ")
	for _, c := range rest {
		if c == '[' {
			return "", false
		}
		if c == '(' {
			break
		}
	}
	if idx := strings.Index(rest, "("); idx >= 0 {
		return strings.TrimSpace(rest[:idx]), true
	}
	return "", false
}

func parseTypeName(line string) (string, bool) {
	rest := strings.TrimPrefix(line, "type ")
	if strings.Contains(rest, "interface{") || strings.Contains(rest, "interface {") {
		return "", false
	}
	for i, c := range rest {
		if c == ' ' || c == '[' || c == '{' {
			return strings.TrimSpace(rest[:i]), true
		}
	}
	fields := strings.Fields(rest)
	if len(fields) == 0 {
		return "", false
	}
	return strings.TrimSpace(fields[0]), true
}

func parseFirstDeclaredName(line, prefix string) (string, bool) {
	fields := strings.Fields(strings.TrimPrefix(line, prefix))
	if len(fields) == 0 {
		return "", false
	}
	return fields[0], true
}

// isExported checks if a name is exported.
func isExported(name string) bool {
	if len(name) == 0 {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}
