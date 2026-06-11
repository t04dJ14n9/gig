package gentool

import (
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageImportGeneratesInterfaceProxy(t *testing.T) {
	outDir := t.TempDir()
	if err := PackageImport("sort", outDir, "packages"); err != nil {
		t.Fatalf("PackageImport: %v", err)
	}

	src, err := os.ReadFile(filepath.Join(outDir, "sort.go"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	code := string(src)

	for _, want := range []string{
		`pkg.AddInterfaceProxy("Interface"`,
		"type proxy_sort_Interface struct",
		"func newProxy_sort_Interface",
		"func (p *proxy_sort_Interface) Len() int",
		"func (p *proxy_sort_Interface) Less(a0 int, a1 int) bool",
		"func (p *proxy_sort_Interface) Swap(a0 int, a1 int)",
	} {
		if !strings.Contains(code, want) {
			t.Fatalf("generated sort.go missing %q\n%s", want, code)
		}
	}

	rootDir, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("Abs: %v", err)
	}
	goMod := "module generatedsorttest\n\ngo 1.23.1\n\nrequire github.com/t04dJ14n9/gig v0.0.0\n\nreplace github.com/t04dJ14n9/gig => " + rootDir + "\n"
	if err := os.WriteFile(filepath.Join(outDir, "go.mod"), []byte(goMod), 0o666); err != nil {
		t.Fatalf("WriteFile(go.mod): %v", err)
	}

	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = outDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated package does not compile: %v\n%s", err, output)
	}
}

func TestWrapBasicReturnPreservesScalarWidth(t *testing.T) {
	tests := []struct {
		name string
		typ  *types.Basic
		want string
	}{
		{"int32", types.Typ[types.Int32], "value.MakeInt32(r0)"},
		{"uint32", types.Typ[types.Uint32], "value.MakeUint32(r0)"},
		{"float32", types.Typ[types.Float32], "value.MakeFloat32(r0)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wrapBasicReturn(tt.typ, "r0"); got != tt.want {
				t.Fatalf("wrapBasicReturn(%s) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
