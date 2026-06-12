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

type symbolKind uint8

const (
	symbolFunc symbolKind = iota
	symbolConst
	symbolVar
	symbolType
)

type parsedSymbol struct {
	kind symbolKind
	name string
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
	symbol, ok := parseExportedSymbol(line)
	if !ok {
		return
	}
	symbols.add(symbol)
}

func parseExportedSymbol(line string) (parsedSymbol, bool) {
	if line == "" || strings.HasPrefix(line, "//") {
		return parsedSymbol{}, false
	}
	return parseSymbolDeclaration(line)
}

func parseSymbolDeclaration(line string) (parsedSymbol, bool) {
	switch {
	case strings.HasPrefix(line, "func "):
		name, ok := parseFunctionName(line)
		return parsedExportedSymbol(symbolFunc, name, ok)
	case strings.HasPrefix(line, "type "):
		name, ok := parseTypeName(line)
		return parsedExportedSymbol(symbolType, name, ok)
	case strings.HasPrefix(line, "const "):
		name, ok := parseFirstDeclaredName(line, "const ")
		return parsedExportedSymbol(symbolConst, name, ok)
	case strings.HasPrefix(line, "var "):
		name, ok := parseFirstDeclaredName(line, "var ")
		return parsedExportedSymbol(symbolVar, name, ok)
	}
	return parsedSymbol{}, false
}

func parsedExportedSymbol(kind symbolKind, name string, ok bool) (parsedSymbol, bool) {
	if !ok || !isExported(name) {
		return parsedSymbol{}, false
	}
	return parsedSymbol{kind: kind, name: name}, true
}

func (symbols *ExportedSymbols) add(symbol parsedSymbol) {
	switch symbol.kind {
	case symbolFunc:
		symbols.Funcs = append(symbols.Funcs, symbol.name)
	case symbolType:
		symbols.Types = append(symbols.Types, symbol.name)
	case symbolConst:
		symbols.Consts = append(symbols.Consts, symbol.name)
	case symbolVar:
		symbols.Vars = append(symbols.Vars, symbol.name)
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
