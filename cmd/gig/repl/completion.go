// Package repl provides the core REPL session management.
//
// This file contains tab completion logic for the REPL session.
package repl

import (
	"strings"

	"git.woa.com/youngjin/gig"
)

// completer provides tab completion for the REPL.
func (s *Session) completer(line string) []string {
	var completions []string

	// Get the word being typed (last word in line)
	line = strings.TrimRight(line, " \t")
	lastSpace := strings.LastIndexAny(line, " \t(")
	var prefix string
	var wordStart int
	if lastSpace >= 0 {
		prefix = line[:lastSpace+1]
		wordStart = lastSpace + 1
	} else {
		wordStart = 0
	}

	word := line[wordStart:]

	// If word is empty, no completions
	if word == "" {
		return completions
	}

	// Check if completing a selector expression (e.g., fmt.)
	if strings.Contains(word, ".") {
		completions = s.completeSelector(word, prefix)
	} else {
		// Complete simple identifiers
		completions = s.completeIdentifier(word, prefix)
	}

	return completions
}

// completeSelector completes package.Symbol expressions.
func (s *Session) completeSelector(word, prefix string) []string {
	var completions []string

	// Split word into package and partial symbol
	dotIdx := strings.LastIndex(word, ".")
	pkgName := word[:dotIdx]
	partialSymbol := word[dotIdx+1:]

	// Get the package symbols
	symbols := s.getPackageSymbols(pkgName)

	// Filter by prefix and add completions
	for _, sym := range symbols {
		if strings.HasPrefix(sym, partialSymbol) {
			completion := prefix + pkgName + "." + sym
			completions = append(completions, completion)
		}
	}

	return completions
}

// completeIdentifier completes simple identifiers (variables, packages, commands).
func (s *Session) completeIdentifier(word, prefix string) []string {
	var completions []string

	// Check if this is a command (starts with :)
	if strings.HasPrefix(word, ":") {
		commands := []string{":help", ":quit", ":clear", ":imports", ":vars", ":env", ":source", ":timeout", ":plugins"}
		for _, cmd := range commands {
			if strings.HasPrefix(cmd, word) {
				completions = append(completions, prefix+cmd)
			}
		}
		return completions
	}

	// Add variable names
	for name := range s.vars {
		if strings.HasPrefix(name, word) {
			completions = append(completions, prefix+name)
		}
	}

	// Add package names (from imports)
	for path := range s.imports {
		pkgName := path
		if idx := strings.LastIndex(path, "/"); idx >= 0 {
			pkgName = path[idx+1:]
		}
		// Check for known packages with different names
		if name, ok := pkgNameFromPath(path); ok {
			pkgName = name
		}
		if strings.HasPrefix(pkgName, word) {
			completions = append(completions, prefix+pkgName)
		}
	}

	// Add known packages
	for name := range knownPackages {
		if strings.HasPrefix(name, word) {
			completions = append(completions, prefix+name)
		}
	}

	// Add registered packages
	for _, pkg := range gig.GetAllPackages() {
		if strings.HasPrefix(pkg.Name, word) {
			completions = append(completions, prefix+pkg.Name)
		}
	}

	return completions
}

// getPackageSymbols returns exported symbols for a package.
func (s *Session) getPackageSymbols(pkgName string) []string {
	var symbols []string

	// Check registered packages first
	if pkg := gig.GetPackageByName(pkgName); pkg != nil {
		// We need to get symbols from the package
		// Since gig doesn't expose symbols directly, we'll use the plugin manager
		if s.pluginMgr != nil {
			symbols = s.pluginMgr.GetSymbols(pkg.Path)
		}
	}

	// Check built-in stdlib packages
	for path := range s.imports {
		name, ok := pkgNameFromPath(path)
		if ok && name == pkgName {
			if s.pluginMgr != nil {
				syms := s.pluginMgr.GetSymbols(path)
				symbols = append(symbols, syms...)
			}
		}
	}

	return symbols
}
