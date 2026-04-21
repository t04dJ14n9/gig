package thirdparty_external

import (
	_ "embed"
	"reflect"
	"strings"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "thirdpartytests/mydep/packages"
)

//go:embed testdata/uuid_source.go
var uuidSrc string

//go:embed testdata/decimal_source.go
var decimalSrc string

//go:embed testdata/semver_source.go
var semverSrc string

type caseItem struct {
	funcName string
	args     []any
	want     any
}

type sourceSet struct {
	src   string
	cases map[string]caseItem
}

func runSourceSet(t *testing.T, set sourceSet) {
	t.Helper()
	prog, err := gig.Build(set.src)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	for name, tc := range set.cases {
		t.Run(name, func(t *testing.T) {
			got, err := prog.Run(tc.funcName, tc.args...)
			if err != nil {
				t.Fatalf("Run %s failed: %v", tc.funcName, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("%s: got %v (%T), want %v (%T)", tc.funcName, got, got, tc.want, tc.want)
			}
		})
	}
}

func TestThirdpartyExternal_UUID(t *testing.T) {
	set := sourceSet{
		src: uuidSrc,
		cases: map[string]caseItem{
			"round trip": {funcName: "UUIDRoundTrip", args: []any{"550e8400-e29b-41d4-a716-446655440000"}, want: "550e8400-e29b-41d4-a716-446655440000"},
			"urn prefix": {funcName: "UUIDURNPrefix", args: nil, want: "urn:uuid:"},
		},
	}
	runSourceSet(t, set)

	// format-only assertion for random uuid string
	prog, err := gig.Build(uuidSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("UUIDNewString")
	if err != nil {
		t.Fatalf("Run UUIDNewString failed: %v", err)
	}
	s, ok := got.(string)
	if !ok {
		t.Fatalf("UUIDNewString: expected string, got %T", got)
	}
	if len(s) != 36 || !strings.Contains(s, "-") {
		t.Fatalf("UUIDNewString: unexpected format %q", s)
	}
}

func TestThirdpartyExternal_Decimal(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: decimalSrc,
		cases: map[string]caseItem{
			"add": {funcName: "DecimalAdd", args: []any{"1.20", "3.05"}, want: "4.25"},
			"sum": {funcName: "DecimalSum", args: nil, want: "6"},
			"avg": {funcName: "DecimalAvg", args: nil, want: "4.67"},
		},
	})
}

func TestThirdpartyExternal_Semver(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: semverSrc,
		cases: map[string]caseItem{
			"check":     {funcName: "SemverCheck", args: []any{"1.4.5", ">= 1.2, < 2.0"}, want: true},
			"inc patch": {funcName: "SemverIncPatch", args: []any{"2.3.4"}, want: "2.3.5"},
			"compare":   {funcName: "SemverCompare", args: []any{"1.2.0", "1.3.0"}, want: -1},
		},
	})
}
