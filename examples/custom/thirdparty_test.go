// Package main demonstrates comprehensive testing of third-party libraries in GIG.
package main

import (
	"fmt"
	"testing"

	_ "myapp/mydep/packages"

	"git.woa.com/youngjin/gig"
)

// TestGjsonBasic tests basic gjson functionality
func TestGjsonBasic(t *testing.T) {
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

// TestGjsonPath tests gjson path parsing
func TestGjsonPath(t *testing.T) {
	source := `
package main

import "github.com/tidwall/gjson"

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

func ValidInvalid(jsonStr string) bool {
	return gjson.Valid(jsonStr)
}

func ParseAndGet(jsonStr string) string {
	parsed := gjson.Parse(jsonStr)
	return parsed.Get("name").String()
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
		{"invalid JSON", "ValidInvalid", []any{`{invalid}`}, false},
		{"parse and get", "ParseAndGet", []any{`{"name": "test"}`}, "test"},
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

// TestCastBasic tests basic cast functionality
func TestCastBasic(t *testing.T) {
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
		{"interface bool", "ToBool", []any{1}, true},
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
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expect, result)
			}
		})
	}
}

// TestMixedLibraries tests using gjson and cast in the same program
func TestMixedLibraries(t *testing.T) {
	// Test gjson alone
	source1 := `
package main

import (
	"fmt"
	"github.com/tidwall/gjson"
)

func ProcessData(jsonStr string, field string) string {
	val := gjson.Get(jsonStr, field)
	return val.String()
}

func ChainOps(data string) string {
	name := gjson.Get(data, "user.name").String()
	age := gjson.Get(data, "user.age").Int()
	active := gjson.Get(data, "user.active").Bool()
	return fmt.Sprintf("%s:%d:%v", name, age, active)
}

func Calculate(jsonStr string) string {
	price := gjson.Get(jsonStr, "price").Float()
	qty := gjson.Get(jsonStr, "quantity").Int()
	total := price * float64(qty)
	tax := total * 0.1
	return fmt.Sprintf("%.2f", total+tax)
}
`
	prog1, err := gig.Build(source1)
	if err != nil {
		t.Fatal(err)
	}

	tests1 := []struct {
		name string
		fn   string
		args []any
		want string
	}{
		{"number field", "ProcessData", []any{`{"count": 42}`, "count"}, "42"},
		{"string field", "ProcessData", []any{`{"name": "test"}`, "name"}, "test"},
		{"chain operations", "ChainOps", []any{`{"user": {"name": "Alice", "age": 30, "active": true}}`}, "Alice:30:true"},
		{"calculate with tax", "Calculate", []any{`{"price": 100.5, "quantity": 2}`}, "221.10"},
	}

	for _, tc := range tests1 {
		t.Run("gjson/"+tc.name, func(t *testing.T) {
			result, err := prog1.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.want {
				t.Errorf("%s: expected %q, got %q", tc.name, tc.want, result)
			}
		})
	}

	// Test cast alone - type conversions
	source2 := `
package main

import (
	"fmt"
	"github.com/spf13/cast"
)

func ConvertAndFormat(val interface{}) string {
	i := cast.ToInt(val)
	f := cast.ToFloat64(val)
	b := cast.ToBool(val)
	return fmt.Sprintf("int:%d,float:%.2f,bool:%v", i, f, b)
}

func SliceConversion(vals interface{}) string {
	si := cast.ToStringSlice(vals)
	return fmt.Sprintf("%v", si)
}
`
	prog2, err := gig.Build(source2)
	if err != nil {
		t.Fatal(err)
	}

	tests2 := []struct {
		name string
		fn   string
		args []any
		want string
	}{
		{"convert int", "ConvertAndFormat", []any{42}, "int:42,float:42.00,bool:true"},
		{"convert float", "ConvertAndFormat", []any{3.14}, "int:3,float:3.14,bool:true"},
		{"convert string", "ConvertAndFormat", []any{"10"}, "int:10,float:10.00,bool:false"},
		{"slice conversion", "SliceConversion", []any{[]any{"a", "b", "c"}}, "[a b c]"},
	}

	for _, tc := range tests2 {
		t.Run("cast/"+tc.name, func(t *testing.T) {
			result, err := prog2.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.want {
				t.Errorf("%s: expected %q, got %q", tc.name, tc.want, result)
			}
		})
	}
}

// TestRealWorldAPI simulates real-world API response parsing
func TestRealWorldAPI(t *testing.T) {
	source := `
package main

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
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

// ExtractUserFields extracts multiple fields using GetMany
func ExtractUserFields(jsonStr string) string {
	fields := gjson.GetMany(jsonStr, "name", "email", "age")
	name := fields[0].String()
	email := fields[1].String()
	age := cast.ToInt(fields[2].Int())
	return fmt.Sprintf("%s<%s>:%d", name, email, age)
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

// ProcessOrder processes order data
func ProcessOrder(jsonStr string) string {
	orderID := gjson.Get(jsonStr, "order.id").String()
	customer := gjson.Get(jsonStr, "order.customer.name").String()
	items := cast.ToInt(gjson.Get(jsonStr, "order.items.#").Int())
	shipping := gjson.Get(jsonStr, "order.shipping.cost").Float()
	return fmt.Sprintf("Order#%s|%s|items:%d|ship:%.2f", orderID, customer, items, shipping)
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

	fieldsAPI := `{"name": "Alice", "email": "alice@test.com", "age": 28}`

	orderAPI := `{
		"order": {
			"id": "ORD-2024-001",
			"customer": {"name": "Bob Smith"},
			"items": [{}, {}, {}],
			"shipping": {"cost": 15.99}
		}
	}`

	tests := []struct {
		name string
		fn   string
		args []any
		want string
	}{
		{"parse user", "ParseUserResponse", []any{userAPI}, "12345|John Doe|john@example.com|30|true"},
		{"extract fields", "ExtractUserFields", []any{fieldsAPI}, "Alice<alice@test.com>:28"},
		{"valid order", "ValidateAndTransform", []any{`{"valid": true, "amount": 100, "taxRate": 0.1, "quantity": 2}`}, "220.00"},
		{"invalid order", "ValidateAndTransform", []any{`{"valid": false, "amount": 100, "taxRate": 0.1, "quantity": 2}`}, "INVALID"},
		{"process order", "ProcessOrder", []any{orderAPI}, "Order#ORD-2024-001|Bob Smith|items:3|ship:15.99"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.want {
				t.Errorf("%s: expected %q, got %q", tc.name, tc.want, result)
			}
		})
	}
}

// TestComplexJson scenarios with complex JSON parsing
func TestComplexJson(t *testing.T) {
	source := `
package main

import (
	"fmt"
	"github.com/tidwall/gjson"
)

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

func NumberOperations(jsonStr string) string {
	x := gjson.Get(jsonStr, "x").Int()
	y := gjson.Get(jsonStr, "y").Int()
	add := x + y
	sub := x - y
	mul := x * y
	return fmt.Sprintf("add:%d,sub:%d,mul:%d", add, sub, mul)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatal(err)
	}

	complexJSON := `{
		"data": {
			"items": [
				{
					"name": "first",
					"subItems": [
						{"value": "a"},
						{"value": "b"}
					]
				}
			]
		}
	}`

	typesJSON := `{
		"value": "string",
		"nullField": null,
		"x": 10,
		"y": 3
	}`

	tests := []struct {
		name string
		fn   string
		args []any
		want any
	}{
		{"deep nested", "DeepNested", []any{complexJSON}, "b"},
		{"multiple paths", "MultiplePaths", []any{`{"a": "hello", "b": 42, "c": true}`}, "hello:42:true"},
		{"array length", "GetArrayLength", []any{`{"items": [1,2,3,4,5]}`}, int64(5)},
		{"number operations", "NumberOperations", []any{typesJSON}, "add:13,sub:7,mul:30"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prog.Run(tc.fn, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if result != tc.want {
				t.Errorf("%s: expected %q, got %q", tc.name, tc.want, result)
			}
		})
	}
}

func TestMain(m *testing.M) {
	fmt.Println("=== Third-Party Library Integration Tests ===")
	fmt.Println("Libraries: github.com/tidwall/gjson, github.com/spf13/cast")
	fmt.Println()
	m.Run()
}
