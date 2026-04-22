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

//go:embed testdata/color_source.go
var colorSrc string

//go:embed testdata/gods_source.go
var godsSrc string

//go:embed testdata/zap_source.go
var zapSrc string

//go:embed testdata/toml_source.go
var tomlSrc string

//go:embed testdata/yaml_source.go
var yamlSrc string

//go:embed testdata/cron_source.go
var cronSrc string

//go:embed testdata/argon2id_source.go
var argon2idSrc string

//go:embed testdata/now_source.go
var nowSrc string

//go:embed testdata/govalidator_source.go
var govalidatorSrc string

//go:embed testdata/testify_source.go
var testifySrc string

//go:embed testdata/isatty_source.go
var isattySrc string

//go:embed testdata/strcase_source.go
var strcaseSrc string

//go:embed testdata/xxhash_source.go
var xxhashSrc string

//go:embed testdata/humanize_source.go
var humanizeSrc string

//go:embed testdata/deepcopy_source.go
var deepcopySrc string

//go:embed testdata/mergo_source.go
var mergoSrc string

//go:embed testdata/jsoniter_source.go
var jsoniterSrc string

//go:embed testdata/pflag_source.go
var pflagSrc string

//go:embed testdata/ewma_source.go
var ewmaSrc string

//go:embed testdata/shellquote_source.go
var shellquoteSrc string

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

// --- color ---

func TestThirdpartyExternal_Color(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: colorSrc,
		cases: map[string]caseItem{
			"red string":  {funcName: "ColorRedString", args: nil, want: "hello"},
			"blue string": {funcName: "ColorBlueString", args: nil, want: "world"},
			"new":         {funcName: "ColorNew", args: nil, want: "ok"},
			"new add":     {funcName: "ColorNewAdd", args: nil, want: "ok"},
		},
	})
}

// --- gods ---

func TestThirdpartyExternal_Gods(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: godsSrc,
		cases: map[string]caseItem{
			"new size":   {funcName: "GodsHashSetNew", args: nil, want: 2},
			"values":     {funcName: "GodsHashSetValues", args: nil, want: 3},
			"contains":   {funcName: "GodsHashSetContains", args: nil, want: true},
			"remove":     {funcName: "GodsHashSetRemove", args: nil, want: 2},
		},
	})
}

// --- zap ---

func TestThirdpartyExternal_Zap(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: zapSrc,
		cases: map[string]caseItem{
			"new logger":     {funcName: "ZapNewLogger", args: nil, want: "ok"},
			"sugar logger":   {funcName: "ZapSugarLogger", args: nil, want: "ok"},
			"named logger":   {funcName: "ZapNamedLogger", args: nil, want: "ok"},
		},
	})
}

// --- toml ---

func TestThirdpartyExternal_Toml(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: tomlSrc,
		cases: map[string]caseItem{
			"unmarshal": {funcName: "TomlUnmarshal", args: nil, want: "hello"},
		},
	})

	// TomlMarshal returns a string containing the TOML — just check it contains "Name"
	prog, err := gig.Build(tomlSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("TomlMarshal")
	if err != nil {
		t.Fatalf("Run TomlMarshal failed: %v", err)
	}
	s, ok := got.(string)
	if !ok {
		t.Fatalf("TomlMarshal: expected string, got %T", got)
	}
	if !strings.Contains(s, "Name") {
		t.Fatalf("TomlMarshal: expected TOML containing 'Name', got %q", s)
	}
}

// --- yaml ---

func TestThirdpartyExternal_Yaml(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: yamlSrc,
		cases: map[string]caseItem{
			"unmarshal": {funcName: "YamlUnmarshal", args: nil, want: "Alice"},
		},
	})

	// YamlMarshal returns a string containing "name:" — just check it
	prog, err := gig.Build(yamlSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("YamlMarshal")
	if err != nil {
		t.Fatalf("Run YamlMarshal failed: %v", err)
	}
	s, ok := got.(string)
	if !ok {
		t.Fatalf("YamlMarshal: expected string, got %T", got)
	}
	if !strings.Contains(s, "name:") {
		t.Fatalf("YamlMarshal: expected YAML containing 'name:', got %q", s)
	}
}

// --- cron ---

func TestThirdpartyExternal_Cron(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: cronSrc,
		cases: map[string]caseItem{
			"new cron":    {funcName: "CronNew", args: nil, want: "ok"},
			"parse cron":  {funcName: "CronParserParse", args: nil, want: "ok"},
		},
	})
}

// --- argon2id ---

func TestThirdpartyExternal_Argon2id(t *testing.T) {
	// Argon2idCreateHash returns first 10 chars of hash — just check it starts with $
	prog, err := gig.Build(argon2idSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("Argon2idCreateHash")
	if err != nil {
		t.Fatalf("Run Argon2idCreateHash failed: %v", err)
	}
	s, ok := got.(string)
	if !ok {
		t.Fatalf("Argon2idCreateHash: expected string, got %T", got)
	}
	if !strings.HasPrefix(s, "$") {
		t.Fatalf("Argon2idCreateHash: expected hash starting with '$', got %q", s)
	}

	// Argon2idCheckPassword returns true
	runSourceSet(t, sourceSet{
		src: argon2idSrc,
		cases: map[string]caseItem{
			"check password": {funcName: "Argon2idCheckPassword", args: nil, want: true},
		},
	})
}

// --- now ---

func TestThirdpartyExternal_Now(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: nowSrc,
		cases: map[string]caseItem{
			"beginning of day": {funcName: "NowBeginningOfDay", args: nil, want: "2024-06-15"},
			"end of day":       {funcName: "NowEndOfDay", args: nil, want: "23:59:59"},
		},
	})
}

// --- govalidator ---

func TestThirdpartyExternal_Govalidator(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: govalidatorSrc,
		cases: map[string]caseItem{
			"is email":  {funcName: "GovalidatorIsEmail", args: nil, want: true},
			"is url":    {funcName: "GovalidatorIsURL", args: nil, want: true},
			"is alpha":  {funcName: "GovalidatorIsAlpha", args: nil, want: true},
			"to string": {funcName: "GovalidatorToString", args: nil, want: "42"},
		},
	})
}

// --- testify ---

func TestThirdpartyExternal_Testify(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: testifySrc,
		cases: map[string]caseItem{
			"importable": {funcName: "TestifyPackageImportable", args: nil, want: "ok"},
		},
	})
}

// --- isatty ---

func TestThirdpartyExternal_Isatty(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: isattySrc,
		cases: map[string]caseItem{
			"is terminal":       {funcName: "IsattyIsTerminal", args: nil, want: false},
			"is cygwin terminal": {funcName: "IsattyIsCygwinTerminal", args: nil, want: false},
		},
	})
}

// --- strcase ---

func TestThirdpartyExternal_Strcase(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: strcaseSrc,
		cases: map[string]caseItem{
			"to camel":              {funcName: "StrcaseToCamel", args: nil, want: "HelloWorld"},
			"to lower camel":        {funcName: "StrcaseToLowerCamel", args: nil, want: "helloWorld"},
			"to snake":              {funcName: "StrcaseToSnake", args: nil, want: "hello_world"},
			"to screaming snake":    {funcName: "StrcaseToScreamingSnake", args: nil, want: "HELLO_WORLD"},
			"to kebab":              {funcName: "StrcaseToKebab", args: nil, want: "hello-world"},
			"to screaming kebab":    {funcName: "StrcaseToScreamingKebab", args: nil, want: "HELLO-WORLD"},
			"to delimited":          {funcName: "StrcaseToDelimited", args: nil, want: "hello.world"},
			"to screaming delimited": {funcName: "StrcaseToScreamingDelimited", args: nil, want: "HELLO.WORLD"},
			"to snake with ignore":  {funcName: "StrcaseToSnakeWithIgnore", args: nil, want: "hello_world_api"},
			"configure acronym":     {funcName: "StrcaseConfigureAcronym", args: nil, want: "MyApiKey"},
		},
	})
}

// --- xxhash ---

func TestThirdpartyExternal_Xxhash(t *testing.T) {
	// Sum64 returns deterministic hash values for the same input
	prog, err := gig.Build(xxhashSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Test Sum64 and Sum64String produce same result
	got64, err := prog.Run("XxhashSum64")
	if err != nil {
		t.Fatalf("Run XxhashSum64 failed: %v", err)
	}
	got64s, err := prog.Run("XxhashSum64String")
	if err != nil {
		t.Fatalf("Run XxhashSum64String failed: %v", err)
	}
	if got64 != got64s {
		t.Fatalf("Sum64 and Sum64String mismatch: %v vs %v", got64, got64s)
	}
	if got64.(uint64) == 0 {
		t.Fatalf("Sum64 returned 0")
	}

	runSourceSet(t, sourceSet{
		src: xxhashSrc,
		cases: map[string]caseItem{
			"size":        {funcName: "XxhashDigestSize", args: nil, want: 8},
			"block size":  {funcName: "XxhashDigestBlockSize", args: nil, want: 32},
			"write":       {funcName: "XxhashDigestWrite", args: nil, want: 5},
			"write string": {funcName: "XxhashDigestWriteString", args: nil, want: 5},
			"sum len":     {funcName: "XxhashDigestSum", args: nil, want: 8},
			"marshal":     {funcName: "XxhashDigestMarshalBinary", args: nil, want: true},
		},
	})
}

// --- humanize ---

func TestThirdpartyExternal_Humanize(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: humanizeSrc,
		cases: map[string]caseItem{
			"bytes":          {funcName: "HumanizeBytes", args: nil, want: "83 MB"},
			"ibytes":         {funcName: "HumanizeIBytes", args: nil, want: "79 MiB"},
			"comma":          {funcName: "HumanizeComma", args: nil, want: "1,234,567"},
			"ordinal":        {funcName: "HumanizeOrdinal", args: nil, want: "1st,2nd,3rd"},
			"parse bytes":    {funcName: "HumanizeParseBytes", args: nil, want: true},
			"parse si":       {funcName: "HumanizeParseSI", args: nil, want: true},
			"compute si":     {funcName: "HumanizeComputeSI", args: nil, want: "345.00µ"},
			"parse big bytes": {funcName: "HumanizeParseBigBytes", args: nil, want: true},
		},
	})

	// Test functions that produce variable output
	prog, err := gig.Build(humanizeSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	for _, name := range []string{"HumanizeCommaf", "HumanizeCommafWithDigits",
		"HumanizeFormatFloat", "HumanizeFormatInteger", "HumanizeFtoa",
		"HumanizeFtoaWithDigits", "HumanizeSI", "HumanizeSIWithDigits",
		"HumanizeTime", "HumanizeRelTime", "HumanizeCustomRelTime"} {
		got, err := prog.Run(name)
		if err != nil {
			t.Fatalf("Run %s failed: %v", name, err)
		}
		s, ok := got.(string)
		if !ok || len(s) == 0 {
			t.Fatalf("%s: expected non-empty string, got %v (%T)", name, got, got)
		}
	}
}

// --- deepcopy ---

func TestThirdpartyExternal_Deepcopy(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: deepcopySrc,
		cases: map[string]caseItem{
			"copy":  {funcName: "DeepcopyCopy", args: nil, want: "value"},
			"iface": {funcName: "DeepcopyIface", args: nil, want: "value"},
		},
	})
}

// --- mergo ---

func TestThirdpartyExternal_Mergo(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: mergoSrc,
		cases: map[string]caseItem{
			"merge":             {funcName: "MergoMerge", args: nil, want: "Alice"},
			"merge override":    {funcName: "MergoMergeWithOverride", args: nil, want: "Bob"},
			"map":               {funcName: "MergoMap", args: nil, want: "3"},
		},
	})
}

// --- jsoniter ---

func TestThirdpartyExternal_Jsoniter(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: jsoniterSrc,
		cases: map[string]caseItem{
			"marshal":        {funcName: "JsoniterMarshal", args: nil, want: `{"key":"value"}`},
			"unmarshal":      {funcName: "JsoniterUnmarshal", args: nil, want: "value"},
			"unmarshal from": {funcName: "JsoniterUnmarshalFromString", args: nil, want: "value"},
			"valid":          {funcName: "JsoniterValid", args: nil, want: true},
			"marshal to str": {funcName: "JsoniterMarshalToString", args: nil, want: `{"key":"value"}`},
			"get":            {funcName: "JsoniterGet", args: nil, want: "value"},
			"wrap":           {funcName: "JsoniterWrap", args: nil, want: "3.14E+00"},
			"wrap string":    {funcName: "JsoniterWrapString", args: nil, want: "hello"},
			"wrap int64":     {funcName: "JsoniterWrapInt64", args: nil, want: "42"},
			"wrap uint64":    {funcName: "JsoniterWrapUint64", args: nil, want: "42"},
			"wrap float64":   {funcName: "JsoniterWrapFloat64", args: nil, want: "3.14E+00"},
			"wrap int32":     {funcName: "JsoniterWrapInt32", args: nil, want: "42"},
			"wrap uint32":    {funcName: "JsoniterWrapUint32", args: nil, want: "42"},
		},
	})

	// Test MarshalIndent separately (output is indented, harder to exact-match)
	prog, err := gig.Build(jsoniterSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	for _, name := range []string{"JsoniterMarshalIndent", "JsoniterConfigDefaultMarshal",
		"JsoniterConfigFastestMarshal", "JsoniterConfigCompatibleMarshal"} {
		got, err := prog.Run(name)
		if err != nil {
			t.Fatalf("Run %s failed: %v", name, err)
		}
		s, ok := got.(string)
		if !ok || !strings.Contains(s, "key") {
			t.Fatalf("%s: expected JSON containing 'key', got %q", name, got)
		}
	}
}

// --- pflag ---

func TestThirdpartyExternal_Pflag(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: pflagSrc,
		cases: map[string]caseItem{
			"new flagset":   {funcName: "PflagNewFlagSet", args: nil, want: "test"},
			"string":        {funcName: "PflagString", args: nil, want: "hello"},
			"bool":          {funcName: "PflagBool", args: nil, want: true},
			"int":           {funcName: "PflagInt", args: nil, want: 42},
			"int64":         {funcName: "PflagInt64", args: nil, want: int64(42)},
			"uint":          {funcName: "PflagUint", args: nil, want: uint(42)},
			"uint64":        {funcName: "PflagUint64", args: nil, want: uint64(42)},
			"float64":       {funcName: "PflagFloat64", args: nil, want: 3.14},
			"string slice":  {funcName: "PflagStringSlice", args: nil, want: 3},
			"int slice":     {funcName: "PflagIntSlice", args: nil, want: 3},
			"args":          {funcName: "PflagArgs", args: nil, want: 2},
			"arg":           {funcName: "PflagArg", args: nil, want: "hello"},
			"n arg":         {funcName: "PflagNArg", args: nil, want: 2},
			"n flag":        {funcName: "PflagNFlag", args: nil, want: 2},
			"changed":       {funcName: "PflagChanged", args: nil, want: true},
			"lookup":        {funcName: "PflagLookup", args: nil, want: true},
			"duration":      {funcName: "PflagDuration", args: nil, want: int64(5)},
			"var":           {funcName: "PflagVar", args: nil, want: "hello"},
			"print defaults": {funcName: "PflagPrintDefaults", args: nil, want: true},
		},
	})
}

// --- ewma ---

func TestThirdpartyExternal_Ewma(t *testing.T) {
	prog, err := gig.Build(ewmaSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Test that EWMA values are positive and reasonable
	for _, name := range []string{"EwmaNewMovingAverage", "EwmaNewMovingAverageWithAge",
		"EwmaSimpleEWMA", "EwmaSimpleEWMASet", "EwmaMovingAverageInterface"} {
		got, err := prog.Run(name)
		if err != nil {
			t.Fatalf("Run %s failed: %v", name, err)
		}
		f, ok := got.(float64)
		if !ok {
			t.Fatalf("%s: expected float64, got %T", name, got)
		}
		if name == "EwmaSimpleEWMASet" && f != 42.0 {
			t.Fatalf("%s: expected 42.0, got %f", name, f)
		}
		if f < 0 {
			t.Fatalf("%s: got negative value %f", name, f)
		}
	}
}

// --- shellquote ---

func TestThirdpartyExternal_Shellquote(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: shellquoteSrc,
		cases: map[string]caseItem{
			"split": {funcName: "ShellquoteSplit", args: nil, want: 2},
		},
	})

	// Test Join separately (exact output may vary)
	prog, err := gig.Build(shellquoteSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("ShellquoteJoin")
	if err != nil {
		t.Fatalf("Run ShellquoteJoin failed: %v", err)
	}
	s, ok := got.(string)
	if !ok || len(s) == 0 {
		t.Fatalf("ShellquoteJoin: expected non-empty string, got %v (%T)", got, got)
	}
}
