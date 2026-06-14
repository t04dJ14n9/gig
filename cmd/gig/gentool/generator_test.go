package gentool

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageImportGeneratesDirectWrappers(t *testing.T) {
	dir := t.TempDir()
	if err := PackageImport("math", dir, "packages"); err != nil {
		t.Fatalf("PackageImport: %v", err)
	}
	src, err := os.ReadFile(filepath.Join(dir, "math.go"))
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	code := string(src)
	checks := []string{
		`pkg.AddFunction("Sqrt", math.Sqrt, "", directCallMathSqrt)`,
		`func directCallMathSqrt(args []value.Value) ([]value.Value, error)`,
		`r0 := math.Sqrt(a0)`,
		`return directResultsMath(r0)`,
		`func directCallMathModf(args []value.Value) ([]value.Value, error)`,
		`r0, r1 := math.Modf(a0)`,
		`return directResultsMath(r0, r1)`,
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Fatalf("generated code missing %q:\n%s", want, code)
		}
	}
}

func TestPackageImportGeneratesThirdPartyDirectWrappers(t *testing.T) {
	dir := t.TempDir()
	if err := PackageImport("golang.org/x/mod/semver", dir, "packages"); err != nil {
		t.Fatalf("PackageImport: %v", err)
	}
	src, err := os.ReadFile(filepath.Join(dir, "golang_org_x_mod_semver.go"))
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	code := string(src)
	checks := []string{
		`pkg.AddFunction("Compare", golang_org_x_mod_semver.Compare, "", directCallGolangOrgXModSemverCompare)`,
		`func directCallGolangOrgXModSemverCompare(args []value.Value) ([]value.Value, error)`,
		`r0 := golang_org_x_mod_semver.Compare(a0, a1)`,
		`return directResultsGolangOrgXModSemver(r0)`,
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Fatalf("generated code missing %q:\n%s", want, code)
		}
	}
}
