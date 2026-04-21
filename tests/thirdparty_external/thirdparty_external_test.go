package thirdparty_external

import (
	_ "embed"
	"reflect"
	"strings"
	"testing"

	"github.com/t04dJ14n9/gig"
	_ "thirdpartytests/mydep/packages"
)

//go:embed testdata/uuid_source.go
var uuidSrc string

//go:embed testdata/decimal_source.go
var decimalSrc string

//go:embed testdata/semver_source.go
var semverSrc string

//go:embed testdata/logrus_source.go
var logrusSrc string

//go:embed testdata/cast_source.go
var castSrc string

//go:embed testdata/jwt_source.go
var jwtSrc string

//go:embed testdata/mapstructure_source.go
var mapstructureSrc string

//go:embed testdata/validator_source.go
var validatorSrc string

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

// --- uuid ---

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

// --- decimal ---

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

// --- semver ---

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

// --- logrus ---

func TestThirdpartyExternal_Logrus(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: logrusSrc,
		cases: map[string]caseItem{
			"parse debug":  {funcName: "LogrusLevelParse", args: []any{"debug"}, want: "debug"},
			"parse info":   {funcName: "LogrusLevelParse", args: []any{"info"}, want: "info"},
			"parse error":  {funcName: "LogrusLevelParse", args: []any{"error"}, want: "error"},
			"get level":    {funcName: "LogrusGetLevel", args: nil, want: "info"},
		},
	})

	// LogrusNewLogger returns JSON with "hello" msg — check it contains that
	prog, err := gig.Build(logrusSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("LogrusNewLogger")
	if err != nil {
		t.Fatalf("Run LogrusNewLogger failed: %v", err)
	}
	s, ok := got.(string)
	if !ok {
		t.Fatalf("LogrusNewLogger: expected string, got %T", got)
	}
	if !strings.Contains(s, "hello") {
		t.Fatalf("LogrusNewLogger: expected JSON containing 'hello', got %q", s)
	}
}

// --- cast ---

func TestThirdpartyExternal_Cast(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: castSrc,
		cases: map[string]caseItem{
			"to string":      {funcName: "CastToString", args: []any{42}, want: "42"},
			"to int":         {funcName: "CastToInt", args: []any{"123"}, want: 123},
			"to float64":     {funcName: "CastToFloat64", args: []any{"3.14"}, want: float64(3.14)},
			"to bool true":   {funcName: "CastToBool", args: []any{1}, want: true},
			"to bool false":  {funcName: "CastToBool", args: []any{0}, want: false},
		},
	})
}

// --- jwt ---

func TestThirdpartyExternal_JWT(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: jwtSrc,
		cases: map[string]caseItem{
			"signing method": {funcName: "JWTGetSigningMethod", args: nil, want: "HS256"},
		},
	})
}

// --- mapstructure ---

func TestThirdpartyExternal_Mapstructure(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: mapstructureSrc,
		cases: map[string]caseItem{
			"decode name": {funcName: "MapstructureDecode", args: []any{map[string]any{"Name": "Alice", "Age": 30}}, want: "Alice"},
			"weak decode": {funcName: "MapstructureWeakDecode", args: []any{map[string]any{"Port": "8080"}}, want: 8080},
		},
	})
}

// --- validator ---

func TestThirdpartyExternal_Validator(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: validatorSrc,
		cases: map[string]caseItem{
			"valid var":    {funcName: "ValidatorVarValid", args: nil, want: true},
			"invalid var":  {funcName: "ValidatorVarInvalid", args: nil, want: true},
		},
	})
}
