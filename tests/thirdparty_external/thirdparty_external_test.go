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

//go:embed testdata/stripansi_source.go
var stripansiSrc string

//go:embed testdata/runewidth_source.go
var runewidthSrc string

//go:embed testdata/uniseg_source.go
var unisegSrc string

//go:embed testdata/ansi_source.go
var ansiSrc string

//go:embed testdata/unidecode_source.go
var unidecodeSrc string

//go:embed testdata/fuzzy_source.go
var fuzzySrc string

//go:embed testdata/slug_source.go
var slugSrc string

//go:embed testdata/durafmt_source.go
var durafmtSrc string

//go:embed testdata/mahonia_source.go
var mahoniaSrc string

//go:embed testdata/iso8601_source.go
var iso8601Src string

//go:embed testdata/gjson_source.go
var gjsonSrc string

//go:embed testdata/sjson_source.go
var sjsonSrc string

//go:embed testdata/carbon_source.go
var carbonSrc string

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

// --- stripansi ---

func TestThirdpartyExternal_Stripansi(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: stripansiSrc,
		cases: map[string]caseItem{
			"strip":   {funcName: "StripansiStrip", args: nil, want: "hello"},
			"plain":   {funcName: "StripansiPlain", args: nil, want: "plain text"},
			"multi":   {funcName: "StripansiMultiple", args: nil, want: "green yellow"},
		},
	})
}

// --- runewidth ---

func TestThirdpartyExternal_Runewidth(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: runewidthSrc,
		cases: map[string]caseItem{
			"rune width":  {funcName: "RunewidthRuneWidth", args: nil, want: 2},
			"fill left":   {funcName: "RunewidthFillLeft", args: nil, want: "      1234"},
			"fill right":  {funcName: "RunewidthFillRight", args: nil, want: "1234      "},
			"truncate":    {funcName: "RunewidthTruncate", args: nil, want: "Hello..."},
		},
	})

	// StringWidth and Wrap output depends on EastAsianWidth
	prog, err := gig.Build(runewidthSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("RunewidthStringWidth")
	if err != nil {
		t.Fatalf("Run RunewidthStringWidth failed: %v", err)
	}
	if got.(int) < 7 {
		t.Fatalf("RunewidthStringWidth: expected >= 7, got %v", got)
	}
	got, err = prog.Run("RunewidthWrap")
	if err != nil {
		t.Fatalf("Run RunewidthWrap failed: %v", err)
	}
	s, ok := got.(string)
	if !ok || len(s) == 0 {
		t.Fatalf("RunewidthWrap: expected non-empty string, got %v (%T)", got, got)
	}
}

// --- uniseg ---

func TestThirdpartyExternal_Uniseg(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: unisegSrc,
		cases: map[string]caseItem{
			"grapheme count":    {funcName: "UnisegGraphemeClusterCount", args: nil, want: 2},
			"has newline":       {funcName: "UnisegHasNewline", args: nil, want: true},
		},
	})

	// StringWidth and WordCount depend on locale/settings
	prog, err := gig.Build(unisegSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("UnisegStringWidth")
	if err != nil {
		t.Fatalf("Run UnisegStringWidth failed: %v", err)
	}
	if got.(int) < 7 {
		t.Fatalf("UnisegStringWidth: expected >= 7, got %v", got)
	}
	got, err = prog.Run("UnisegWordCount")
	if err != nil {
		t.Fatalf("Run UnisegWordCount failed: %v", err)
	}
	if got.(int) < 4 {
		t.Fatalf("UnisegWordCount: expected >= 4, got %v", got)
	}

	// FirstGraphemeCluster returns a non-empty string
	got, err = prog.Run("UnisegFirstGraphemeCluster")
	if err != nil {
		t.Fatalf("Run UnisegFirstGraphemeCluster failed: %v", err)
	}
	s, ok := got.(string)
	if !ok || len(s) == 0 {
		t.Fatalf("UnisegFirstGraphemeCluster: expected non-empty string, got %v (%T)", got, got)
	}
}

// --- ansi ---

func TestThirdpartyExternal_Ansi(t *testing.T) {
	// Color returns string with ANSI escape codes wrapping the text
	prog, err := gig.Build(ansiSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("AnsiColor")
	if err != nil {
		t.Fatalf("Run AnsiColor failed: %v", err)
	}
	s, ok := got.(string)
	if !ok || !strings.Contains(s, "hello") {
		t.Fatalf("AnsiColor: expected string containing 'hello', got %q", got)
	}

	// ColorFunc produces ANSI escape codes
	got, err = prog.Run("AnsiColorFunc")
	if err != nil {
		t.Fatalf("Run AnsiColorFunc failed: %v", err)
	}
	s, ok = got.(string)
	if !ok || !strings.Contains(s, "world") {
		t.Fatalf("AnsiColorFunc: expected string containing 'world', got %q", got)
	}

	// ColorCode returns ANSI escape code string
	got, err = prog.Run("AnsiColorCode")
	if err != nil {
		t.Fatalf("Run AnsiColorCode failed: %v", err)
	}
	s, ok = got.(string)
	if !ok || !strings.HasPrefix(s, "\x1b[") {
		t.Fatalf("AnsiColorCode: expected ANSI escape code, got %q", got)
	}

	got, err = prog.Run("AnsiReset")
	if err != nil {
		t.Fatalf("Run AnsiReset failed: %v", err)
	}
	s, ok = got.(string)
	if !ok || !strings.HasPrefix(s, "\x1b[") {
		t.Fatalf("AnsiReset: expected ANSI escape code, got %q", got)
	}
}

// --- unidecode ---

func TestThirdpartyExternal_Unidecode(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: unidecodeSrc,
		cases: map[string]caseItem{
			"transliterate": {funcName: "UnidecodeUnidecode", args: nil, want: "Bei Jing kozuscek"},
			"plain":         {funcName: "UnidecodePlain", args: nil, want: "Hello World"},
		},
	})

	// Version returns a non-empty string
	prog, err := gig.Build(unidecodeSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("UnidecodeVersion")
	if err != nil {
		t.Fatalf("Run UnidecodeVersion failed: %v", err)
	}
	s, ok := got.(string)
	if !ok || len(s) == 0 {
		t.Fatalf("UnidecodeVersion: expected non-empty string, got %v (%T)", got, got)
	}
}

// --- fuzzy ---

func TestThirdpartyExternal_Fuzzy(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: fuzzySrc,
		cases: map[string]caseItem{
			"match":       {funcName: "FuzzyMatch", args: nil, want: true},
			"match fold":  {funcName: "FuzzyMatchFold", args: nil, want: true},
			"find count":  {funcName: "FuzzyFind", args: nil, want: 2},
			"lev distance": {funcName: "FuzzyLevenshteinDistance", args: nil, want: 3},
		},
	})

	// RankMatch returns a non-negative score for matches
	prog, err := gig.Build(fuzzySrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	got, err := prog.Run("FuzzyRankMatch")
	if err != nil {
		t.Fatalf("Run FuzzyRankMatch failed: %v", err)
	}
	if got.(int) <= 0 {
		t.Fatalf("FuzzyRankMatch: expected positive score, got %v", got)
	}
}

// --- slug ---

func TestThirdpartyExternal_Slug(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: slugSrc,
		cases: map[string]caseItem{
			"make":       {funcName: "SlugMake", args: nil, want: "hello-world-khello-vorld"},
			"make lang":  {funcName: "SlugMakeLang", args: nil, want: "diese-und-dass"},
			"is slug":    {funcName: "SlugIsSlug", args: nil, want: true},
			"not slug":   {funcName: "SlugIsSlugInvalid", args: nil, want: false},
			"substitute": {funcName: "SlugSubstitute", args: nil, want: "Hello Go"},
		},
	})
}

// --- durafmt ---

func TestThirdpartyExternal_Durafmt(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: durafmtSrc,
		cases: map[string]caseItem{
			"duration": {funcName: "DurafmtDuration", args: nil, want: int64(5400000000000)},
		},
	})

	// Parse and Format produce human-readable strings
	prog, err := gig.Build(durafmtSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	for _, name := range []string{"DurafmtParse", "DurafmtParseShort", "DurafmtLimitFirstN"} {
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

// --- mahonia ---

func TestThirdpartyExternal_Mahonia(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: mahoniaSrc,
		cases: map[string]caseItem{
			"get charset utf8": {funcName: "MahoniaGetCharset", args: nil, want: "UTF-8"},
			"get charset gbk":  {funcName: "MahoniaGetCharsetGbk", args: nil, want: "GBK"},
			"register charset": {funcName: "MahoniaRegisterCharset", args: nil, want: "test-charset"},
		},
	})
}

// --- iso8601 ---

func TestThirdpartyExternal_Iso8601(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: iso8601Src,
		cases: map[string]caseItem{
			"parse":          {funcName: "Iso8601Parse", args: nil, want: "2024-06-15"},
			"parse string":   {funcName: "Iso8601ParseString", args: nil, want: "10:30:00"},
			"parse offset":   {funcName: "Iso8601ParseWithOffset", args: nil, want: "10"},
		},
	})
}

// --- gjson ---

func TestThirdpartyExternal_Gjson(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: gjsonSrc,
		cases: map[string]caseItem{
			"get name":      {funcName: "GjsonGetName", args: []any{`{"name": "Alice"}`}, want: "Alice"},
			"get age":       {funcName: "GjsonGetAge", args: []any{`{"age": 30}`}, want: int64(30)},
			"nested":        {funcName: "GjsonGetNested", args: []any{`{"user": {"name": "Bob"}}`}, want: "Bob"},
			"array access":  {funcName: "GjsonArrayAccess", args: []any{`{"items": [{"name": "first"}]}`}, want: "first"},
			"exists true":   {funcName: "GjsonExists", args: []any{`{"name": "test"}`}, want: true},
			"exists false":  {funcName: "GjsonNotExists", args: []any{`{"name": "test"}`}, want: false},
			"bool true":     {funcName: "GjsonBoolValue", args: []any{`{"active": true}`}, want: true},
			"bool false":    {funcName: "GjsonBoolValue", args: []any{`{"active": false}`}, want: false},
			"float":         {funcName: "GjsonFloatValue", args: []any{`{"price": 19.99}`}, want: 19.99},
			"uint":          {funcName: "GjsonUintValue", args: []any{`{"count": 100}`}, want: uint64(100)},
			"simple path":   {funcName: "GjsonGetPath", args: []any{`{"a": {"b": "c"}}`, "a.b"}, want: "c"},
			"array index":   {funcName: "GjsonGetPath", args: []any{`{"items": ["x", "y", "z"]}`, "items.1"}, want: "y"},
			"get many":      {funcName: "GjsonGetMany", args: []any{`{"name": "Alice", "age": 30}`}, want: "Alice:30"},
			"valid JSON":    {funcName: "GjsonValid", args: []any{`{"valid": true}`}, want: true},
			"invalid JSON":  {funcName: "GjsonValid", args: []any{`{invalid}`}, want: false},
			"parse and get": {funcName: "GjsonParseAndGet", args: []any{`{"name": "test"}`}, want: "test"},
			"deep nested":   {funcName: "GjsonDeepNested", args: []any{`{"data": {"items": [{"name": "first", "subItems": [{"value": "a"}, {"value": "b"}]}]}}`}, want: "b"},
			"multi paths":   {funcName: "GjsonMultiplePaths", args: []any{`{"a": "hello", "b": 42, "c": true}`}, want: "hello:42:true"},
			"array length":  {funcName: "GjsonGetArrayLength", args: []any{`{"items": [1,2,3,4,5]}`}, want: int64(5)},
			"is array":      {funcName: "GjsonIsArray", args: []any{`{"items": [1,2,3]}`}, want: true},
			"is object":     {funcName: "GjsonIsObject", args: []any{`{"user": {"name": "x"}}`}, want: true},
			"map values":    {funcName: "GjsonMapValues", args: []any{`{"user": {"name": "Bob", "age": 25}}`}, want: "Bob"},
		},
	})
}

// --- sjson ---

func TestThirdpartyExternal_Sjson(t *testing.T) {
	// Set operations — check the result contains the expected value
	prog, err := gig.Build(sjsonSrc)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	setCases := []struct {
		name     string
		funcName string
		args     []any
		contains string
	}{
		{"set name", "SjsonSetName", []any{`{"name": "old"}`, "new"}, `"new"`},
		{"set age", "SjsonSetAge", []any{`{}`, int(25)}, `25`},
		{"set nested", "SjsonSetNested", []any{`{}`, "Beijing"}, `"Beijing"`},
		{"set bool true", "SjsonSetBool", []any{`{}`, true}, `true`},
		{"set bool false", "SjsonSetBool", []any{`{}`, false}, `false`},
		{"set float", "SjsonSetFloat", []any{`{}`, 9.99}, `9.99`},
		{"set null", "SjsonSetNull", []any{`{}`}, `null`},
		{"set raw", "SjsonSetRaw", []any{`{}`, `{"x":1}`}, `{"x":1}`},
	}

	for _, tc := range setCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := prog.Run(tc.funcName, tc.args...)
			if err != nil {
				t.Fatalf("Run %s failed: %v", tc.funcName, err)
			}
			s, ok := got.(string)
			if !ok {
				t.Fatalf("%s: expected string, got %T", tc.name, got)
			}
			if !strings.Contains(s, tc.contains) {
				t.Errorf("%s: expected result to contain %q, got %q", tc.name, tc.contains, s)
			}
		})
	}

	// Delete operations — exact match
	runSourceSet(t, sourceSet{
		src: sjsonSrc,
		cases: map[string]caseItem{
			"delete name":     {funcName: "SjsonDeleteField", args: []any{`{"name": "Alice", "age": 30}`, "name"}, want: `{ "age": 30}`},
			"delete nested":   {funcName: "SjsonDeleteNested", args: []any{`{"address": {"city": "Beijing", "zip": "100000"}}`}, want: `{"address": {"city": "Beijing"}}`},
		},
	})
}

// --- carbon ---

func TestThirdpartyExternal_Carbon(t *testing.T) {
	runSourceSet(t, sourceSet{
		src: carbonSrc,
		cases: map[string]caseItem{
			"parse date":       {funcName: "CarbonParseDate", args: []any{"2024-01-15"}, want: "2024-01-15"},
			"create from date": {funcName: "CarbonCreateFromDate", args: []any{int64(2024), int64(6), int64(15)}, want: "2024-06-15"},
			"create from time": {funcName: "CarbonCreateFromTime", args: []any{int64(14), int64(30), int64(0)}, want: "14:30:00"},
			"create datetime":  {funcName: "CarbonCreateFromDateTime", args: []any{int64(2024), int64(3), int64(20), int64(10), int64(0), int64(0)}, want: "2024-03-20 10:00:00"},
			"add 1 day":        {funcName: "CarbonAddDays", args: []any{"2024-01-15", int64(1)}, want: "2024-01-16"},
			"add 10 days":      {funcName: "CarbonAddDays", args: []any{"2024-01-15", int64(10)}, want: "2024-01-25"},
			"sub 1 day":        {funcName: "CarbonSubDays", args: []any{"2024-01-15", int64(1)}, want: "2024-01-14"},
			"sub 15 days":      {funcName: "CarbonSubDays", args: []any{"2024-01-15", int64(15)}, want: "2023-12-31"},
			"add 1 month":      {funcName: "CarbonAddMonths", args: []any{"2024-01-31", int64(1)}, want: "2024-03-02"},
			"add 1 year":       {funcName: "CarbonAddYears", args: []any{"2024-01-15", int64(1)}, want: "2025-01-15"},
			"add 2 hours":      {funcName: "CarbonAddHours", args: []any{"2024-01-15 10:00:00", int64(2)}, want: "12:00:00"},
			"add 30 min":       {funcName: "CarbonAddMinutes", args: []any{"2024-01-15 10:00:00", int64(30)}, want: "10:30:00"},
			"start of day":     {funcName: "CarbonStartOfDay", args: []any{"2024-01-15"}, want: "2024-01-15 00:00:00"},
			"end of day":       {funcName: "CarbonEndOfDay", args: []any{"2024-01-15"}, want: "2024-01-15 23:59:59"},
			"start of month":   {funcName: "CarbonStartOfMonth", args: []any{"2024-01-15"}, want: "2024-01-01"},
			"end of month":     {funcName: "CarbonEndOfMonth", args: []any{"2024-01-15"}, want: "2024-01-31"},
			"start of year":    {funcName: "CarbonStartOfYear", args: []any{"2024-06-15"}, want: "2024-01-01"},
			"end of year":      {funcName: "CarbonEndOfYear", args: []any{"2024-06-15"}, want: "2024-12-31"},
			"is weekend true":  {funcName: "CarbonIsWeekend", args: []any{"2024-01-13"}, want: true},
			"is weekend false": {funcName: "CarbonIsWeekend", args: []any{"2024-01-15"}, want: false},
			"is leap 2024":     {funcName: "CarbonIsLeapYear", args: []any{"2024-01-01"}, want: true},
			"is leap 2023":     {funcName: "CarbonIsLeapYear", args: []any{"2023-01-01"}, want: false},
			"day of week":      {funcName: "CarbonDayOfWeek", args: []any{"2024-01-15"}, want: 1},
			"day of year":      {funcName: "CarbonDayOfYear", args: []any{"2024-01-15"}, want: 15},
			"month":            {funcName: "CarbonMonth", args: []any{"2024-01-15"}, want: 1},
			"year":             {funcName: "CarbonYear", args: []any{"2024-01-15"}, want: 2024},
			"days in jan":      {funcName: "CarbonDaysInMonth", args: []any{"2024-01-15"}, want: 31},
			"days in feb leap": {funcName: "CarbonDaysInMonth", args: []any{"2024-02-15"}, want: 29},
			"days in feb norm": {funcName: "CarbonDaysInMonth", args: []any{"2023-02-15"}, want: 28},
			"timestamp":        {funcName: "CarbonToTimestamp", args: []any{"2024-01-01 00:00:00"}, want: int64(1704067200)},
			"rfc3339":          {funcName: "CarbonToRfc3339", args: []any{"2024-01-15 10:30:00"}, want: "2024-01-15T10:30:00Z"},
			"iso8601":          {funcName: "CarbonToIso8601", args: []any{"2024-01-15 10:30:00"}, want: "2024-01-15T10:30:00+00:00"},
			"layout":           {funcName: "CarbonLayoutFormat", args: []any{"2024-01-15", "2006/01/02"}, want: "2024/01/15"},
			"layout time":      {funcName: "CarbonLayoutFormat", args: []any{"2024-01-15 14:30:00", "15:04:05"}, want: "14:30:00"},
			"parse and fmt":    {funcName: "CarbonParseAndFormat", args: []any{"2024-06-15"}, want: "2024-6-15"},
			"ts conversion":    {funcName: "CarbonTimestampConversion", args: []any{"2024-01-01 00:00:00"}, want: "1704067200"},
		},
	})
}
