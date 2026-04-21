// Package main tests external library calls through the gig interpreter.
// Tests: gjson, sjson, cast, carbon
package main

import (
	"fmt"
	"strings"
	"testing"

	_ "myapp/mydep/packages"

	"git.woa.com/youngjin/gig"
)

// ============================================================
// GJSON Tests (~25 tests)
// ============================================================

func TestExtGjsonBasic(t *testing.T) {
	source := `
package main

import "github.com/tidwall/gjson"

func GetName(jsonStr string) string {
	return gjson.Get(jsonStr, "name").String()
}

func GetAge(jsonStr string) int64 {
	return gjson.Get(jsonStr, "age").Int()
}

func GetNested(jsonStr string) string {
	return gjson.Get(jsonStr, "user.name").String()
}

func ArrayAccess(jsonStr string) string {
	return gjson.Get(jsonStr, "items.0.name").String()
}

func Exists(jsonStr string) bool {
	return gjson.Get(jsonStr, "name").Exists()
}

func NotExists(jsonStr string) bool {
	return gjson.Get(jsonStr, "nonexistent").Exists()
}

func BoolValue(jsonStr string) bool {
	return gjson.Get(jsonStr, "active").Bool()
}

func FloatValue(jsonStr string) float64 {
	return gjson.Get(jsonStr, "price").Float()
}

func UintValue(jsonStr string) uint64 {
	return gjson.Get(jsonStr, "count").Uint()
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect any
	}{
		{"simple string", "GetName", []any{`{"name": "Alice"}`}, "Alice"},
		{"simple int", "GetAge", []any{`{"age": 30}`}, int64(30)},
		{"nested", "GetNested", []any{`{"user": {"name": "Bob"}}`}, "Bob"},
		{"array access", "ArrayAccess", []any{`{"items": [{"name": "first"}]}`}, "first"},
		{"exists true", "Exists", []any{`{"name": "test"}`}, true},
		{"exists false", "NotExists", []any{`{"name": "test"}`}, false},
		{"bool true", "BoolValue", []any{`{"active": true}`}, true},
		{"bool false", "BoolValue", []any{`{"active": false}`}, false},
		{"float", "FloatValue", []any{`{"price": 19.99}`}, 19.99},
		{"uint", "UintValue", []any{`{"count": 100}`}, uint64(100)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.expect {
				t.Errorf("%s: expected %v (%T), got %v (%T)", tc.name, tc.expect, tc.expect, result, result)
			}
		})
	}
}

func TestExtGjsonPath(t *testing.T) {
	source := `
package main

import (
	"fmt"
	"github.com/tidwall/gjson"
)

func GetPath(jsonStr string, path string) string {
	return gjson.Get(jsonStr, path).String()
}

func GetMany(jsonStr string) string {
	results := gjson.GetMany(jsonStr, "name", "age")
	return results[0].String() + ":" + results[1].String()
}

func Valid(jsonStr string) bool {
	return gjson.Valid(jsonStr)
}

func ParseAndGet(jsonStr string) string {
	parsed := gjson.Parse(jsonStr)
	return parsed.Get("name").String()
}

func DeepNested(jsonStr string) string {
	return gjson.Get(jsonStr, "data.items.0.subItems.1.value").String()
}

func MultiplePaths(jsonStr string) string {
	a := gjson.Get(jsonStr, "a").String()
	b := gjson.Get(jsonStr, "b").Int()
	c := gjson.Get(jsonStr, "c").Bool()
	return fmt.Sprintf("%s:%d:%v", a, b, c)
}

func GetArrayLength(jsonStr string) int64 {
	arr := gjson.Get(jsonStr, "items").Array()
	return int64(len(arr))
}

func IsArrayCheck(jsonStr string) bool {
	return gjson.Get(jsonStr, "items").IsArray()
}

func IsObjectCheck(jsonStr string) bool {
	return gjson.Get(jsonStr, "user").IsObject()
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect any
	}{
		{"simple path", "GetPath", []any{`{"a": {"b": "c"}}`, "a.b"}, "c"},
		{"array index", "GetPath", []any{`{"items": ["x", "y", "z"]}`, "items.1"}, "y"},
		{"get many", "GetMany", []any{`{"name": "Alice", "age": 30}`}, "Alice:30"},
		{"valid JSON", "Valid", []any{`{"valid": true}`}, true},
		{"invalid JSON", "Valid", []any{`{invalid}`}, false},
		{"parse and get", "ParseAndGet", []any{`{"name": "test"}`}, "test"},
		{"deep nested", "DeepNested", []any{`{"data": {"items": [{"name": "first", "subItems": [{"value": "a"}, {"value": "b"}]}]}}`}, "b"},
		{"multiple paths", "MultiplePaths", []any{`{"a": "hello", "b": 42, "c": true}`}, "hello:42:true"},
		{"array length", "GetArrayLength", []any{`{"items": [1,2,3,4,5]}`}, int64(5)},
		{"is array", "IsArrayCheck", []any{`{"items": [1,2,3]}`}, true},
		{"is object", "IsObjectCheck", []any{`{"user": {"name": "x"}}`}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.expect {
				t.Errorf("%s: expected %v (%T), got %v (%T)", tc.name, tc.expect, tc.expect, result, result)
			}
		})
	}
}

// ============================================================
// SJSON Tests (~20 tests)
// ============================================================

func TestSjsonSet(t *testing.T) {
	source := `
package main

import "github.com/tidwall/sjson"

func SetName(jsonStr string, name string) string {
	result, _ := sjson.Set(jsonStr, "name", name)
	return result
}

func SetAge(jsonStr string, age int) string {
	result, _ := sjson.Set(jsonStr, "age", age)
	return result
}

func SetNested(jsonStr string, city string) string {
	result, _ := sjson.Set(jsonStr, "address.city", city)
	return result
}

func SetBool(jsonStr string, active bool) string {
	result, _ := sjson.Set(jsonStr, "active", active)
	return result
}

func SetFloat(jsonStr string, price float64) string {
	result, _ := sjson.Set(jsonStr, "price", price)
	return result
}

func SetNull(jsonStr string) string {
	result, _ := sjson.Set(jsonStr, "removed", nil)
	return result
}

func SetRaw(jsonStr string, raw string) string {
	result, _ := sjson.SetRaw(jsonStr, "extra", raw)
	return result
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect string
	}{
		{"set name", "SetName", []any{`{"name": "old"}`, "new"}, `"new"`},
		{"set age", "SetAge", []any{`{}`, int(25)}, `25`},
		{"set nested", "SetNested", []any{`{}`, "Beijing"}, `"Beijing"`},
		{"set bool true", "SetBool", []any{`{}`, true}, `true`},
		{"set bool false", "SetBool", []any{`{}`, false}, `false`},
		{"set float", "SetFloat", []any{`{}`, 9.99}, `9.99`},
		{"set null", "SetNull", []any{`{}`}, `null`},
		{"set raw", "SetRaw", []any{`{}`, `{"x":1}`}, `{"x":1}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			jsonStr, ok := result.(string)
			if !ok {
				t.Fatalf("%s: expected string, got %T", tc.name, result)
			}
			// Check the set value is present in the output JSON
			if jsonStr == "" {
				t.Errorf("%s: got empty string", tc.name)
			}
			if !strings.Contains(jsonStr, tc.expect) {
				t.Errorf("%s: expected result to contain %v, got %v", tc.name, tc.expect, jsonStr)
			}
			t.Logf("%s: result = %s", tc.name, jsonStr)
		})
	}
}

func TestSjsonDelete(t *testing.T) {
	source := `
package main

import "github.com/tidwall/sjson"

func DeleteField(jsonStr string, field string) string {
	result, _ := sjson.Delete(jsonStr, field)
	return result
}

func DeleteNested(jsonStr string) string {
	result, _ := sjson.Delete(jsonStr, "address.zip")
	return result
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect string
	}{
		{"delete name", "DeleteField", []any{`{"name": "Alice", "age": 30}`, "name"}, `{ "age": 30}`},
		{"delete nested", "DeleteNested", []any{`{"address": {"city": "Beijing", "zip": "100000"}}`}, `{"address": {"city": "Beijing"}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			jsonStr, ok := result.(string)
			if !ok {
				t.Fatalf("%s: expected string, got %T", tc.name, result)
			}
			if jsonStr != tc.expect {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expect, jsonStr)
			}
			t.Logf("%s: result = %s", tc.name, jsonStr)
		})
	}
}

// ============================================================
// CAST Tests (~25 tests)
// ============================================================

func TestExtCastBasic(t *testing.T) {
	source := `
package main

import "github.com/spf13/cast"

func ToStringInt(i int) string {
	return cast.ToString(i)
}

func ToStringFloat(f float64) string {
	return cast.ToString(f)
}

func ToStringBool(b bool) string {
	return cast.ToString(b)
}

func ToInt(i interface{}) int {
	return cast.ToInt(i)
}

func ToInt64(i interface{}) int64 {
	return cast.ToInt64(i)
}

func ToFloat64(i interface{}) float64 {
	return cast.ToFloat64(i)
}

func ToBool(i interface{}) bool {
	return cast.ToBool(i)
}

func ToStringSlice(i interface{}) []string {
	return cast.ToStringSlice(i)
}

func ToIntSlice(i interface{}) []int {
	return cast.ToIntSlice(i)
}

func ToBoolSlice(i interface{}) []bool {
	return cast.ToBoolSlice(i)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect any
	}{
		{"int to string", "ToStringInt", []any{42}, "42"},
		{"float to string", "ToStringFloat", []any{3.14}, "3.14"},
		{"bool to string", "ToStringBool", []any{true}, "true"},
		{"interface int", "ToInt", []any{int64(100)}, 100},
		{"interface int64", "ToInt64", []any{float64(200.5)}, int64(200)},
		{"interface float64", "ToFloat64", []any{"3.14"}, 3.14},
		{"interface bool true", "ToBool", []any{1}, true},
		{"interface bool false", "ToBool", []any{0}, false},
		{"bool true", "ToBool", []any{true}, true},
		{"bool false", "ToBool", []any{false}, false},
		{"string slice", "ToStringSlice", []any{[]any{"a", "b"}}, []string{"a", "b"}},
		{"int slice", "ToIntSlice", []any{[]any{1, 2, 3}}, []int{1, 2, 3}},
		{"bool slice", "ToBoolSlice", []any{[]any{true, false, true}}, []bool{true, false, true}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expect) {
				t.Errorf("%s: expected %v (%T), got %v (%T)", tc.name, tc.expect, tc.expect, result, result)
			}
		})
	}
}

func TestCastAdvanced(t *testing.T) {
	source := `
package main

import (
	"fmt"
	"github.com/spf13/cast"
)

func ToStringMapString(i interface{}) string {
	m := cast.ToStringMapString(i)
	return fmt.Sprintf("%v", m)
}

func ToStringMap(i interface{}) string {
	m := cast.ToStringMap(i)
	return fmt.Sprintf("%v", m)
}

func ToDuration(d interface{}) string {
	dur := cast.ToDuration(d)
	return dur.String()
}

func ToUint(i interface{}) uint {
	return cast.ToUint(i)
}

func ToUint64(i interface{}) uint64 {
	return cast.ToUint64(i)
}

func ToInt32(i interface{}) int32 {
	return cast.ToInt32(i)
}

func ToFloat32(i interface{}) float32 {
	return cast.ToFloat32(i)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect string
	}{
		{"string map string", "ToStringMapString", []any{map[string]any{"key": "value"}}, "map[key:value]"},
		{"to duration", "ToDuration", []any{"1h30m"}, "1h30m0s"},
		{"to uint", "ToUint", []any{int64(42)}, "42"},
		{"to uint64", "ToUint64", []any{int64(100)}, "100"},
		{"to int32", "ToInt32", []any{int64(99)}, "99"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			resultStr := fmt.Sprintf("%v", result)
			if resultStr != tc.expect {
				t.Errorf("%s: expected %q, got %q", tc.name, tc.expect, resultStr)
			}
		})
	}
}

// ============================================================
// CARBON Tests (~25 tests)
// ============================================================

func TestCarbonCreate(t *testing.T) {
	source := `
package main

import (
	"github.com/dromara/carbon/v2"
)

func NowDate() string {
	return carbon.Now().ToDateString()
}

func NowTime() string {
	return carbon.Now().ToTimeString()
}

func NowDateTime() string {
	return carbon.Now().ToDateTimeString()
}

func ParseDate(s string) string {
	c := carbon.Parse(s)
	return c.ToDateString()
}

func CreateFromDate(y int, m int, d int) string {
	c := carbon.CreateFromDate(y, m, d)
	return c.ToDateString()
}

func CreateFromTime(h int, m int, s int) string {
	c := carbon.CreateFromTime(h, m, s)
	return c.ToTimeString()
}

func CreateFromDateTime(y int, mo int, d int, h int, mi int, s int) string {
	c := carbon.CreateFromDateTime(y, mo, d, h, mi, s)
	return c.ToDateTimeString()
}

func Tomorrow() string {
	return carbon.Tomorrow().ToDateString()
}

func Yesterday() string {
	return carbon.Yesterday().ToDateString()
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	// Parse dates - use fixed date parsing
	result, err := prog.Run("ParseDate", "2024-01-15")
	if err != nil {
		t.Fatalf("ParseDate failed: %v", err)
	}
	if result != "2024-01-15" {
		t.Errorf("ParseDate: expected 2024-01-15, got %v", result)
	}

	result, err = prog.Run("CreateFromDate", 2024, 6, 15)
	if err != nil {
		t.Fatalf("CreateFromDate failed: %v", err)
	}
	if result != "2024-06-15" {
		t.Errorf("CreateFromDate: expected 2024-06-15, got %v", result)
	}

	result, err = prog.Run("CreateFromTime", 14, 30, 0)
	if err != nil {
		t.Fatalf("CreateFromTime failed: %v", err)
	}
	if result != "14:30:00" {
		t.Errorf("CreateFromTime: expected 14:30:00, got %v", result)
	}

	result, err = prog.Run("CreateFromDateTime", 2024, 3, 20, 10, 0, 0)
	if err != nil {
		t.Fatalf("CreateFromDateTime failed: %v", err)
	}
	if result != "2024-03-20 10:00:00" {
		t.Errorf("CreateFromDateTime: expected 2024-03-20 10:00:00, got %v", result)
	}

	// Tomorrow/Yesterday
	result, err = prog.Run("Tomorrow")
	if err != nil {
		t.Fatalf("Tomorrow failed: %v", err)
	}
	t.Logf("Tomorrow: %v", result)

	result, err = prog.Run("Yesterday")
	if err != nil {
		t.Fatalf("Yesterday failed: %v", err)
	}
	t.Logf("Yesterday: %v", result)

	// Now
	result, err = prog.Run("NowDate")
	if err != nil {
		t.Fatalf("NowDate failed: %v", err)
	}
	t.Logf("NowDate: %v", result)
}

func TestCarbonManipulation(t *testing.T) {
	source := `
package main

import (
	"github.com/dromara/carbon/v2"
)

func AddDays(s string, days int) string {
	c := carbon.Parse(s)
	return c.AddDays(days).ToDateString()
}

func SubDays(s string, days int) string {
	c := carbon.Parse(s)
	return c.SubDays(days).ToDateString()
}

func AddMonths(s string, months int) string {
	c := carbon.Parse(s)
	return c.AddMonths(months).ToDateString()
}

func AddYears(s string, years int) string {
	c := carbon.Parse(s)
	return c.AddYears(years).ToDateString()
}

func AddHours(s string, hours int) string {
	c := carbon.Parse(s)
	return c.AddHours(hours).ToTimeString()
}

func AddMinutes(s string, minutes int) string {
	c := carbon.Parse(s)
	return c.AddMinutes(minutes).ToTimeString()
}

func StartOfDay(s string) string {
	c := carbon.Parse(s)
	return c.StartOfDay().ToDateTimeString()
}

func EndOfDay(s string) string {
	c := carbon.Parse(s)
	return c.EndOfDay().ToDateTimeString()
}

func StartOfMonth(s string) string {
	c := carbon.Parse(s)
	return c.StartOfMonth().ToDateString()
}

func EndOfMonth(s string) string {
	c := carbon.Parse(s)
	return c.EndOfMonth().ToDateString()
}

func StartOfYear(s string) string {
	c := carbon.Parse(s)
	return c.StartOfYear().ToDateString()
}

func EndOfYear(s string) string {
	c := carbon.Parse(s)
	return c.EndOfYear().ToDateString()
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect any
	}{
		{"add 1 day", "AddDays", []any{"2024-01-15", int64(1)}, "2024-01-16"},
		{"add 10 days", "AddDays", []any{"2024-01-15", int64(10)}, "2024-01-25"},
		{"sub 1 day", "SubDays", []any{"2024-01-15", int64(1)}, "2024-01-14"},
		{"sub 15 days", "SubDays", []any{"2024-01-15", int64(15)}, "2023-12-31"},
		{"add 1 month", "AddMonths", []any{"2024-01-31", int64(1)}, "2024-03-02"},
		{"add 1 year", "AddYears", []any{"2024-01-15", int64(1)}, "2025-01-15"},
		{"add 2 hours", "AddHours", []any{"2024-01-15 10:00:00", int64(2)}, "12:00:00"},
		{"add 30 minutes", "AddMinutes", []any{"2024-01-15 10:00:00", int64(30)}, "10:30:00"},
		{"start of day", "StartOfDay", []any{"2024-01-15"}, "2024-01-15 00:00:00"},
		{"end of day", "EndOfDay", []any{"2024-01-15"}, "2024-01-15 23:59:59"},
		{"start of month", "StartOfMonth", []any{"2024-01-15"}, "2024-01-01"},
		{"end of month", "EndOfMonth", []any{"2024-01-15"}, "2024-01-31"},
		{"start of year", "StartOfYear", []any{"2024-06-15"}, "2024-01-01"},
		{"end of year", "EndOfYear", []any{"2024-06-15"}, "2024-12-31"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.expect {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expect, result)
			}
		})
	}
}

func TestCarbonComparison(t *testing.T) {
	source := `
package main

import "github.com/dromara/carbon/v2"

func IsWeekend(s string) bool {
	c := carbon.Parse(s)
	return c.IsWeekend()
}

func IsWeekday(s string) bool {
	c := carbon.Parse(s)
	return c.IsWeekday()
}

func IsLeapYear(s string) bool {
	c := carbon.Parse(s)
	return c.IsLeapYear()
}

func IsFuture(s string) bool {
	c := carbon.Parse(s)
	return c.IsFuture()
}

func IsPast(s string) bool {
	c := carbon.Parse(s)
	return c.IsPast()
}

func DayOfWeek(s string) int {
	c := carbon.Parse(s)
	return c.DayOfWeek()
}

func DayOfYear(s string) int {
	c := carbon.Parse(s)
	return c.DayOfYear()
}

func WeekOfYear(s string) int {
	c := carbon.Parse(s)
	return c.WeekOfYear()
}

func Month(s string) int {
	c := carbon.Parse(s)
	return c.Month()
}

func Year(s string) int {
	c := carbon.Parse(s)
	return c.Year()
}

func DaysInMonth(s string) int {
	c := carbon.Parse(s)
	return c.DaysInMonth()
}

func ToTimestamp(s string) int64 {
	c := carbon.Parse(s)
	return c.Timestamp()
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect any
	}{
		// 2024-01-15 is a Monday (weekday)
		{"is weekend false", "IsWeekend", []any{"2024-01-15"}, false},
		// 2024-01-13 is a Saturday (weekend)
		{"is weekend true", "IsWeekend", []any{"2024-01-13"}, true},
		{"is weekday true", "IsWeekday", []any{"2024-01-15"}, true},
		{"is leap year 2024", "IsLeapYear", []any{"2024-01-01"}, true},
		{"is leap year 2023", "IsLeapYear", []any{"2023-01-01"}, false},
		{"day of week", "DayOfWeek", []any{"2024-01-15"}, 1},  // Monday
		{"day of year", "DayOfYear", []any{"2024-01-15"}, 15},
		{"month", "Month", []any{"2024-01-15"}, 1},
		{"year", "Year", []any{"2024-01-15"}, 2024},
		{"days in month jan", "DaysInMonth", []any{"2024-01-15"}, 31},
		{"days in month feb leap", "DaysInMonth", []any{"2024-02-15"}, 29},
		{"days in month feb non-leap", "DaysInMonth", []any{"2023-02-15"}, 28},
		{"timestamp", "ToTimestamp", []any{"2024-01-01 00:00:00"}, int64(1704067200)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.expect {
				t.Errorf("%s: expected %v (%T), got %v (%T)", tc.name, tc.expect, tc.expect, result, result)
			}
		})
	}
}

func TestCarbonFormat(t *testing.T) {
	source := `
package main

import "github.com/dromara/carbon/v2"

func ToRfc3339(s string) string {
	c := carbon.Parse(s)
	return c.ToRfc3339String()
}

func ToRfc850(s string) string {
	c := carbon.Parse(s)
	return c.ToRfc850String()
}

func ToAtom(s string) string {
	c := carbon.Parse(s)
	return c.ToAtomString()
}

func ToCookie(s string) string {
	c := carbon.Parse(s)
	return c.ToCookieString()
}

func ToIso8601(s string) string {
	c := carbon.Parse(s)
	return c.ToIso8601String()
}

func ToRss(s string) string {
	c := carbon.Parse(s)
	return c.ToRssString()
}

func LayoutFormat(s string, layout string) string {
	c := carbon.Parse(s)
	return c.Layout(layout)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect any
	}{
		{"rfc3339", "ToRfc3339", []any{"2024-01-15 10:30:00"}, "2024-01-15T10:30:00Z"},
		{"atom", "ToAtom", []any{"2024-01-15 10:30:00"}, "2024-01-15T10:30:00Z"},
		{"iso8601", "ToIso8601", []any{"2024-01-15 10:30:00"}, "2024-01-15T10:30:00+00:00"},
		{"layout", "LayoutFormat", []any{"2024-01-15", "2006/01/02"}, "2024/01/15"},
		{"layout time", "LayoutFormat", []any{"2024-01-15 14:30:00", "15:04:05"}, "14:30:00"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.expect {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expect, result)
			}
		})
	}
}

// ============================================================
// Mixed Library Tests (~15 tests)
// ============================================================

func TestMixedGjsonSjson(t *testing.T) {
	source := `
package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func GetAndModify(jsonStr string, field string, val interface{}) string {
	original := gjson.Get(jsonStr, field).String()
	result, _ := sjson.Set(jsonStr, field, fmt.Sprintf("modified_%v", val))
	return original + "->" + result
}

func SetAndVerify(jsonStr string) string {
	result, _ := sjson.Set(jsonStr, "status", "active")
	name := gjson.Get(result, "name").String()
	status := gjson.Get(result, "status").String()
	return name + ":" + status
}

func DeleteAndCheck(jsonStr string) string {
	result, _ := sjson.Delete(jsonStr, "age")
	hasAge := gjson.Get(result, "age").Exists()
	name := gjson.Get(result, "name").String()
	return fmt.Sprintf("%v:%s", hasAge, name)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	result, err := prog.Run("SetAndVerify", `{"name": "Alice"}`)
	if err != nil {
		t.Fatalf("SetAndVerify failed: %v", err)
	}
	resultStr, ok := result.(string)
	if !ok {
		t.Fatalf("SetAndVerify: expected string, got %T", result)
	}
	if resultStr != "Alice:active" {
		t.Fatalf("SetAndVerify: expected %q, got %q", "Alice:active", resultStr)
	}
	t.Logf("SetAndVerify: %v", resultStr)

	result, err = prog.Run("DeleteAndCheck", `{"name": "Bob", "age": 30}`)
	if err != nil {
		t.Fatalf("DeleteAndCheck failed: %v", err)
	}
	t.Logf("DeleteAndCheck: %v", result)
}

func TestMixedGjsonCast(t *testing.T) {
	source := `
package main

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

func ParseAndConvert(jsonStr string) string {
	name := gjson.Get(jsonStr, "name").String()
	age := cast.ToInt(gjson.Get(jsonStr, "age").Int())
	active := cast.ToBool(gjson.Get(jsonStr, "active").Bool())
	return fmt.Sprintf("%s:%d:%v", name, age, active)
}

func ConvertFields(jsonStr string) string {
	id := cast.ToString(gjson.Get(jsonStr, "id").Int())
	score := cast.ToString(gjson.Get(jsonStr, "score").Float())
	return id + ":" + score
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect string
	}{
		{"parse and convert", "ParseAndConvert", []any{`{"name": "Alice", "age": 30, "active": true}`}, "Alice:30:true"},
		{"convert fields", "ConvertFields", []any{`{"id": 123, "score": 95.5}`}, "123:95.5"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.expect {
				t.Errorf("%s: expected %q, got %q", tc.name, tc.expect, result)
			}
		})
	}
}

func TestMixedCarbonCast(t *testing.T) {
	source := `
package main

import (
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

func ParseAndFormat(s string) string {
	c := carbon.Parse(s)
	year := cast.ToString(c.Year())
	month := cast.ToString(c.Month())
	day := cast.ToString(c.Day())
	return year + "-" + month + "-" + day
}

func TimestampConversion(s string) string {
	c := carbon.Parse(s)
	ts := cast.ToString(c.Timestamp())
	return ts
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect string
	}{
		{"parse and format", "ParseAndFormat", []any{"2024-06-15"}, "2024-6-15"},
		{"timestamp conversion", "TimestampConversion", []any{"2024-01-01 00:00:00"}, "1704067200"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.expect {
				t.Errorf("%s: expected %q, got %q", tc.name, tc.expect, result)
			}
		})
	}
}

func TestExtRealWorldAPI(t *testing.T) {
	source := `
package main

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/dromara/carbon/v2"
)

// ParseUserResponse parses user API response
func ParseUserResponse(jsonStr string) string {
	id := gjson.Get(jsonStr, "data.user.id").Int()
	name := gjson.Get(jsonStr, "data.user.name").String()
	email := gjson.Get(jsonStr, "data.user.email").String()
	age := cast.ToInt(gjson.Get(jsonStr, "data.user.age").Int())
	active := gjson.Get(jsonStr, "data.user.active").Bool()
	return fmt.Sprintf("%d|%s|%s|%d|%v", id, name, email, age, active)
}

// ValidateAndTransform validates and transforms API data
func ValidateAndTransform(jsonStr string) string {
	valid := gjson.Get(jsonStr, "valid").Bool()
	if !valid {
		return "INVALID"
	}
	amount := gjson.Get(jsonStr, "amount").Float()
	taxRate := gjson.Get(jsonStr, "taxRate").Float()
	quantity := cast.ToInt(gjson.Get(jsonStr, "quantity").Int())
	total := amount * float64(quantity)
	tax := total * taxRate
	grandTotal := total + tax
	return fmt.Sprintf("%.2f", grandTotal)
}

// FormatOrderDate formats an order date
func FormatOrderDate(jsonStr string) string {
	dateStr := gjson.Get(jsonStr, "order.date").String()
	c := carbon.Parse(dateStr)
	return c.ToDateString()
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	userAPI := `{
		"data": {
			"user": {
				"id": 12345,
				"name": "John Doe",
				"email": "john@example.com",
				"age": 30,
				"active": true
			}
		}
	}`

	tests := []struct {
		name   string
		fn     string
		args   []any
		expect string
	}{
		{"parse user", "ParseUserResponse", []any{userAPI}, "12345|John Doe|john@example.com|30|true"},
		{"valid order", "ValidateAndTransform", []any{`{"valid": true, "amount": 100, "taxRate": 0.1, "quantity": 2}`}, "220.00"},
		{"invalid order", "ValidateAndTransform", []any{`{"valid": false, "amount": 100, "taxRate": 0.1, "quantity": 2}`}, "INVALID"},
		{"format order date", "FormatOrderDate", []any{`{"order": {"date": "2024-06-15 10:30:00"}}`}, "2024-06-15"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.expect {
				t.Errorf("%s: expected %q, got %q", tc.name, tc.expect, result)
			}
		})
	}
}

func TestGjsonMethodChain(t *testing.T) {
	source := `
package main

import (
	"fmt"
	"github.com/tidwall/gjson"
)

func ForEachKey(jsonStr string) string {
	result := ""
	gjson.Get(jsonStr, "items").ForEach(func(key, value gjson.Result) bool {
		result += fmt.Sprintf("%s:%s|", key.String(), value.Get("name").String())
		return true
	})
	return result
}

func MapValues(jsonStr string) string {
	m := gjson.Get(jsonStr, "user").Map()
	name := m["name"].String()
	return name
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	result, err := prog.Run("MapValues", `{"user": {"name": "Bob", "age": 25}}`)
	if err != nil {
		t.Fatalf("MapValues failed: %v", err)
	}
	if result != "Bob" {
		t.Errorf("MapValues: expected Bob, got %v", result)
	}

	// ForEach is tricky - it involves callback, test if it works
	result, err = prog.Run("ForEachKey", `{"items": [{"name": "a"}, {"name": "b"}]}`)
	if err != nil {
		t.Fatalf("ForEachKey failed: %v", err)
	}
	t.Logf("ForEachKey: %v", result)
}
