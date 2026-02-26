// gig is a CLI tool for generating gig dependency packages.
//
// Usage:
//
//	# Initialize a dependency package directory
//	gig init -package mydep
//
//	# Edit mydep/pkgs.go to add third-party libraries, then generate
//	gig gen ./mydep
//
//	# In your program, import the generated package
//	import _ "myapp/mydep/packages"
//
//	# Run directly from remote (Go 1.21+)
//	go run github.com/t04dJ14n9/gig/cmd/gig@latest init -package mydep
package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"gig/gentool"
)

// Command structure for subcommands
type command struct {
	name  string
	usage string
	run   func()
}

var commands = []command{
	{"init", "gig init -package <name>", runInit},
	{"gen", "gig gen <dir>", runGen},
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gig - generate gig dependency packages\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  gig <command> [arguments]\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		for _, cmd := range commands {
			fmt.Fprintf(os.Stderr, "  %s\n", cmd.usage)
		}
		fmt.Fprintf(os.Stderr, "\nWorkflow:\n")
		fmt.Fprintf(os.Stderr, "  1. gig init -package mydep         # Creates mydep/pkgs.go\n")
		fmt.Fprintf(os.Stderr, "  2. Edit mydep/pkgs.go              # Add third-party libraries\n")
		fmt.Fprintf(os.Stderr, "  3. gig gen ./mydep                 # Generate registration code\n")
		fmt.Fprintf(os.Stderr, "  4. import _ \"myapp/mydep/packages\"      # Use in your program\n")
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	cmdName := os.Args[1]
	for _, cmd := range commands {
		if cmd.name == cmdName {
			cmd.run()
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmdName)
	flag.Usage()
	os.Exit(1)
}

// ========== init command ==========

var flagPackage = flag.String("package", "", "Package name for the dependency (required)")

func runInit() {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	fs.StringVar(flagPackage, "package", "", "Package name for the dependency (required)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gig init -package <name>\n\n")
		fmt.Fprintf(os.Stderr, "Creates a directory <name>/ with pkgs.go containing stdlib imports.\n")
		fmt.Fprintf(os.Stderr, "Edit the file to add third-party libraries.\n\n")
		fs.PrintDefaults()
	}
	_ = fs.Parse(os.Args[2:])

	if *flagPackage == "" {
		fmt.Fprintf(os.Stderr, "Error: -package is required\n\n")
		fs.Usage()
		os.Exit(1)
	}

	// Validate package name (must be a valid Go identifier)
	if !isValidPackageName(*flagPackage) {
		fmt.Fprintf(os.Stderr, "Error: %q is not a valid Go package name\n", *flagPackage)
		os.Exit(1)
	}

	// Create directory
	pkgDir := *flagPackage
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		os.Exit(1)
	}

	// Generate pkgs.go with custom package name
	pkgsGo := generatePkgsGo(*flagPackage)
	pkgsPath := filepath.Join(pkgDir, "pkgs.go")
	if err := os.WriteFile(pkgsPath, pkgsGo, 0o666); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing pkgs.go: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s/\n", pkgDir)
	fmt.Printf("  %s\n", pkgsPath)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. Edit %s to add third-party libraries\n", pkgsPath)
	fmt.Printf("  2. Run: gig gen ./%s\n", pkgDir)
	fmt.Printf("  3. Import in your program: import _ \"%s\"\n", *flagPackage)
}

func generatePkgsGo(pkgName string) []byte {
	var b strings.Builder
	b.WriteString("// Package ")
	b.WriteString(pkgName)
	b.WriteString(" declares dependencies for gig interpreter.\n")
	b.WriteString("// Standard library packages are included by default.\n")
	b.WriteString("// Add your custom third-party library imports at the end.\n")
	b.WriteString("package ")
	b.WriteString(pkgName)
	b.WriteString("\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t// ============================================\n")
	b.WriteString("\t// Go Standard Library (provided by gig)\n")
	b.WriteString("\t// ============================================\n")
	b.WriteString("\t_ \"bytes\"\n")
	b.WriteString("\t_ \"cmp\"\n")
	b.WriteString("\t_ \"container/heap\"\n")
	b.WriteString("\t_ \"container/list\"\n")
	b.WriteString("\t_ \"container/ring\"\n")
	b.WriteString("\t_ \"context\"\n")
	b.WriteString("\t_ \"crypto/hmac\"\n")
	b.WriteString("\t_ \"crypto/sha256\"\n")
	b.WriteString("\t_ \"encoding/base64\"\n")
	b.WriteString("\t_ \"encoding/csv\"\n")
	b.WriteString("\t_ \"encoding/hex\"\n")
	b.WriteString("\t_ \"encoding/json\"\n")
	b.WriteString("\t_ \"encoding/xml\"\n")
	b.WriteString("\t_ \"errors\"\n")
	b.WriteString("\t_ \"fmt\"\n")
	b.WriteString("\t_ \"html\"\n")
	b.WriteString("\t_ \"html/template\"\n")
	b.WriteString("\t_ \"io\"\n")
	b.WriteString("\t_ \"log\"\n")
	b.WriteString("\t_ \"maps\"\n")
	b.WriteString("\t_ \"math\"\n")
	b.WriteString("\t_ \"math/rand\"\n")
	b.WriteString("\t_ \"net/http\"\n")
	b.WriteString("\t_ \"net/url\"\n")
	b.WriteString("\t_ \"os\"\n")
	b.WriteString("\t_ \"path\"\n")
	b.WriteString("\t_ \"path/filepath\"\n")
	b.WriteString("\t_ \"regexp\"\n")
	b.WriteString("\t_ \"slices\"\n")
	b.WriteString("\t_ \"sort\"\n")
	b.WriteString("\t_ \"strconv\"\n")
	b.WriteString("\t_ \"strings\"\n")
	b.WriteString("\t_ \"sync\"\n")
	b.WriteString("\t_ \"sync/atomic\"\n")
	b.WriteString("\t_ \"text/template\"\n")
	b.WriteString("\t_ \"time\"\n")
	b.WriteString("\t_ \"unicode\"\n")
	b.WriteString("\t_ \"unicode/utf8\"\n")
	b.WriteString("\t_ \"unicode/utf16\"\n")
	b.WriteString("\n")
	b.WriteString("\t// ============================================\n")
	b.WriteString("\t// Custom third-party libraries (add yours below)\n")
	b.WriteString("\t// ============================================\n")
	b.WriteString("\t// _ \"github.com/spf13/cast\"\n")
	b.WriteString("\t// _ \"github.com/tidwall/gjson\"\n")
	b.WriteString(")\n")

	// Format the code
	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		return []byte(b.String())
	}
	return formatted
}

func isValidPackageName(name string) bool {
	if name == "" {
		return false
	}
	// Simple check: must start with letter, contain only letters, digits, underscores
	if name[0] < 'a' || name[0] > 'z' {
		if name[0] < 'A' || name[0] > 'Z' {
			return false
		}
	}
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

// ========== gen command ==========

func runGen() {
	fs := flag.NewFlagSet("gen", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gig gen <dir>\n\n")
		fmt.Fprintf(os.Stderr, "Generates dependency registration code from <dir>/pkgs.go.\n")
		fmt.Fprintf(os.Stderr, "The generated files will be placed in <dir>/packages/.\n\n")
		fmt.Fprintf(os.Stderr, "After generation, import the package in your program:\n")
		fmt.Fprintf(os.Stderr, "  import _ \"<package-name>/packages\"\n")
	}
	_ = fs.Parse(os.Args[2:])

	if len(fs.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Error: directory argument required\n\n")
		fs.Usage()
		os.Exit(1)
	}

	pkgDir := fs.Arg(0)
	pkgsPath := filepath.Join(pkgDir, "pkgs.go")

	// Check if pkgs.go exists
	if _, err := os.Stat(pkgsPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s not found\n", pkgsPath)
		fmt.Fprintf(os.Stderr, "Run 'gig init -package %s' first\n", filepath.Base(pkgDir))
		os.Exit(1)
	}

	// Parse pkgs.go to get package name and imports
	importPaths, pkgName, err := parsePkgsGo(pkgsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %s\n", pkgsPath, err)
		os.Exit(1)
	}

	if len(importPaths) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no imports found in %s\n", pkgsPath)
		os.Exit(1)
	}

	fmt.Printf("Generating dependency package %q:\n", pkgName)
	fmt.Printf("  source: %s\n", pkgsPath)
	fmt.Printf("  output: %s/packages/\n", pkgDir)
	fmt.Printf("  packages: %d\n\n", len(importPaths))

	// Create packages subdirectory
	packagesDir := filepath.Join(pkgDir, "packages")
	if err := os.MkdirAll(packagesDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating packages directory: %s\n", err)
		os.Exit(1)
	}

	// Generate registration files (package name is "packages")
	const generatedPkgName = "packages"
	var generatedCount int
	for _, path := range importPaths {
		fmt.Printf("importing %s\n", path)
		err := gentool.PackageImport(path, packagesDir, generatedPkgName)
		if err != nil {
			fmt.Printf("Error importing %s: %s\n", path, err.Error())
			continue
		}
		generatedCount++
	}

	fmt.Printf("\n✓ Generated %d packages\n", generatedCount)
	fmt.Printf("\nDone! Add this to your program:\n")
	fmt.Printf("  import _ %q\n", pkgName+"/packages")
}

func sanitizePkgName(path string) string {
	return strings.NewReplacer(
		"/", "_",
		"-", "_",
		".", "_",
	).Replace(path)
}

// parsePkgsGo reads a pkgs.go file and extracts the package name and import paths.
func parsePkgsGo(path string) ([]string, string, error) {
	// Use gentool's parser for imports
	imports, err := gentool.ParsePkgsFile(path)
	if err != nil {
		return nil, "", err
	}

	// Extract package name by reading the file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	// Simple regex-like search for package name
	lines := strings.Split(string(content), "\n")
	var pkgName string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package ") {
			pkgName = strings.TrimSpace(strings.TrimPrefix(line, "package "))
			// Remove any trailing comments
			if idx := strings.Index(pkgName, "//"); idx >= 0 {
				pkgName = strings.TrimSpace(pkgName[:idx])
			}
			break
		}
	}

	if pkgName == "" {
		return nil, "", fmt.Errorf("could not find package declaration")
	}

	return imports, pkgName, nil
}
