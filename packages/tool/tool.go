// tool.go generates gig package registration files with typed DirectCall wrappers.
// This is the internal tool used by gig for generating its own standard library packages.
//
// For external users, use the CLI tool:
//
//	go run github.com/t04dJ14n9/gig/cmd/gig@latest init -package mydep
//
//go:generate go run tool.go
package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"gig/gentool"
	"gig/packages/tool/pkgs"
)

var sourceDir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	sourceDir = filepath.Dir(filename)
}

func main() {
	// Default mode: use built-in stdlib list, output to sibling directories
	importPaths := pkgs.ImportPkgs
	outBase := filepath.Dir(sourceDir)
	modPrefix := "gig"

	fmt.Printf("Generating gig internal stdlib packages (%d packages):\n\n", len(importPaths))

	var generatedPkgs []string
	for _, path := range importPaths {
		fmt.Printf("importing %s\n", path)
		err := gentool.PackageImport(path, outBase, modPrefix)
		if err != nil {
			fmt.Printf("Error importing %s: %s\n", path, err.Error())
			continue
		}
		generatedPkgs = append(generatedPkgs, sanitizePkgName(path))
	}

	fmt.Printf("\n✓ Generated %d packages\n", len(generatedPkgs))
}

func sanitizePkgName(path string) string {
	return strings.NewReplacer(
		"/", "_",
		"-", "_",
		".", "_",
	).Replace(path)
}
