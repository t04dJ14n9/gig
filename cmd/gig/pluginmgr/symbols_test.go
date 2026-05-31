package pluginmgr

import (
	"reflect"
	"testing"
)

func TestAddExportedSymbolClassifiesSupportedDeclarations(t *testing.T) {
	symbols := &ExportedSymbols{}

	for _, line := range []string{
		"func Parse(input string) (Value, error)",
		"type Decoder struct {",
		"const MaxSize = 10",
		"var DefaultClient = &Client{}",
	} {
		addExportedSymbol(symbols, line)
	}

	if got, want := symbols.Funcs, []string{"Parse"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Funcs = %v, want %v", got, want)
	}
	if got, want := symbols.Types, []string{"Decoder"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Types = %v, want %v", got, want)
	}
	if got, want := symbols.Consts, []string{"MaxSize"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Consts = %v, want %v", got, want)
	}
	if got, want := symbols.Vars, []string{"DefaultClient"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Vars = %v, want %v", got, want)
	}
}

func TestAddExportedSymbolSkipsUnsupportedDeclarations(t *testing.T) {
	symbols := &ExportedSymbols{}

	for _, line := range []string{
		"",
		"// comment",
		"func Decode[T any](input string) T",
		"func (d Decoder) Decode(input string) error",
		"type Number interface { ~int }",
		"const localValue = 10",
		"var privateClient = &Client{}",
	} {
		addExportedSymbol(symbols, line)
	}

	if len(symbols.Funcs) != 0 || len(symbols.Types) != 0 || len(symbols.Consts) != 0 || len(symbols.Vars) != 0 {
		t.Fatalf("symbols = %+v, want no supported exported symbols", symbols)
	}
}

func TestSanitizePkgNameForImport(t *testing.T) {
	tests := map[string]string{
		"github.com/spf13/cast": "github_com_spf13_cast",
		"golang.org/x/text":     "golang_org_x_text",
		"3rd.party/pkg-name":    "pkg_3rd_party_pkg_name",
	}

	for input, want := range tests {
		if got := sanitizePkgNameForImport(input); got != want {
			t.Fatalf("sanitizePkgNameForImport(%q) = %q, want %q", input, got, want)
		}
	}
}
