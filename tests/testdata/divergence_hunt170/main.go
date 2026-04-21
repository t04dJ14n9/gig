package divergence_hunt170

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ============================================================================
// Round 170: JSON marshaling/unmarshaling edge cases
// ============================================================================

// BasicMarshalUnmarshal tests basic JSON marshal/unmarshal
func BasicMarshalUnmarshal() string {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	p := Person{Name: "Alice", Age: 30}
	data, _ := json.Marshal(p)
	var decoded Person
	json.Unmarshal(data, &decoded)
	return fmt.Sprintf("name=%s,age=%d", decoded.Name, decoded.Age)
}

// SliceMarshalUnmarshal tests slice marshal/unmarshal
func SliceMarshalUnmarshal() string {
	numbers := []int{1, 2, 3, 4, 5}
	data, _ := json.Marshal(numbers)
	var decoded []int
	json.Unmarshal(data, &decoded)
	sum := 0
	for _, n := range decoded {
		sum += n
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapMarshalUnmarshal tests map marshal/unmarshal
func MapMarshalUnmarshal() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	data, _ := json.Marshal(m)
	var decoded map[string]int
	json.Unmarshal(data, &decoded)
	sum := 0
	for _, v := range decoded {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// NestedStructMarshal tests nested struct marshal
func NestedStructMarshal() string {
	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}
	p := Person{
		Name:    "Bob",
		Address: Address{City: "NYC", Country: "USA"},
	}
	data, _ := json.Marshal(p)
	return fmt.Sprintf("json=%s", string(data))
}

// PointerFieldMarshal tests pointer field marshal
func PointerFieldMarshal() string {
	type Container struct {
		Value *int `json:"value"`
	}
	v := 42
	c1 := Container{Value: &v}
	data1, _ := json.Marshal(c1)
	c2 := Container{Value: nil}
	data2, _ := json.Marshal(c2)
	return fmt.Sprintf("with_value=%s,nil_value=%s", string(data1), string(data2))
}

// OmitEmptyTag tests omitempty tag
func OmitEmptyTag() string {
	type Config struct {
		Name  string `json:"name,omitempty"`
		Value int    `json:"value,omitempty"`
		Empty string `json:"empty,omitempty"`
	}
	c := Config{Name: "test", Value: 0, Empty: ""}
	data, _ := json.Marshal(c)
	return fmt.Sprintf("json=%s", string(data))
}

// StringTag tests string tag for numeric types
func StringTag() string {
	type Record struct {
		ID   int    `json:"id,string"`
		Code string `json:"code"`
	}
	r := Record{ID: 12345, Code: "ABC"}
	data, _ := json.Marshal(r)
	return fmt.Sprintf("json=%s", string(data))
}

// IgnoreField tests ignored fields
func IgnoreField() string {
	type Internal struct {
		Public  string `json:"public"`
		Private string `json:"-"`
	}
	i := Internal{Public: "visible", Private: "hidden"}
	data, _ := json.Marshal(i)
	return fmt.Sprintf("json=%s", string(data))
}

// UnmarshalUnknownFields tests unmarshaling with unknown fields
func UnmarshalUnknownFields() string {
	type Minimal struct {
		ID int `json:"id"`
	}
	jsonData := `{"id": 1, "name": "unknown", "extra": "data"}`
	var m Minimal
	json.Unmarshal([]byte(jsonData), &m)
	return fmt.Sprintf("id=%d", m.ID)
}

// UnmarshalTypeMismatch tests type mismatch handling
func UnmarshalTypeMismatch() string {
	type Data struct {
		Count int `json:"count"`
	}
	jsonData := `{"count": "not_a_number"}`
	var d Data
	err := json.Unmarshal([]byte(jsonData), &d)
	return fmt.Sprintf("err=%v,count=%d", err != nil, d.Count)
}

// RawMessage tests json.RawMessage
func RawMessage() string {
	type Wrapper struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	w := Wrapper{
		Type: "person",
		Data: json.RawMessage(`{"name":"Alice"}`),
	}
	data, _ := json.Marshal(w)
	return fmt.Sprintf("json=%s", string(data))
}

// NumberHandling tests json.Number
func NumberHandling() string {
	var result interface{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(`{"value": 123.456}`)))
	decoder.UseNumber()
	decoder.Decode(&result)
	m := result.(map[string]interface{})
	num := m["value"].(json.Number)
	return fmt.Sprintf("num=%s", num.String())
}

// CustomMarshalJSON implements custom marshaler
func CustomMarshalJSON() string {
	type Timestamp struct {
		Seconds int64
	}
	t := Timestamp{Seconds: 1609459200}
	// Marshal manually
	data := []byte(fmt.Sprintf(`{"timestamp":%d}`, t.Seconds))
	return fmt.Sprintf("json=%s", string(data))
}

// MarshalEmptySlice tests empty slice marshaling
func MarshalEmptySlice() string {
	type Data struct {
		Items []string `json:"items"`
	}
	d := Data{Items: []string{}}
	data, _ := json.Marshal(d)
	return fmt.Sprintf("json=%s", string(data))
}

// MarshalNilSlice tests nil slice marshaling
func MarshalNilSlice() string {
	type Data struct {
		Items []string `json:"items"`
	}
	d := Data{Items: nil}
	data, _ := json.Marshal(d)
	return fmt.Sprintf("json=%s", string(data))
}

// EscapeHTML tests HTML escaping
func EscapeHTML() string {
	type Content struct {
		HTML string `json:"html"`
	}
	c := Content{HTML: "<script>alert('xss')</script>"}
	data, _ := json.Marshal(c)
	return fmt.Sprintf("json=%s", string(data))
}

// IndentJSON tests indented JSON
func IndentJSON() string {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	p := Person{Name: "Alice", Age: 30}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.Encode(p)
	return fmt.Sprintf("has_indent=%t", strings.Contains(buf.String(), "\n"))
}

// DecodeStream tests streaming decode
func DecodeStream() string {
	jsonData := `{"a":1}{"a":2}{"a":3}`
	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonData)))
	sum := 0
	type Item struct{ A int }
	for decoder.More() || decoder.InputOffset() < int64(len(jsonData)) {
		var item Item
		if err := decoder.Decode(&item); err != nil {
			break
		}
		sum += item.A
	}
	return fmt.Sprintf("sum=%d", sum)
}

// LargeNumber tests large number handling
func LargeNumber() string {
	type Data struct {
		Big int64 `json:"big"`
	}
	jsonData := `{"big": 9223372036854775807}`
	var d Data
	json.Unmarshal([]byte(jsonData), &d)
	return fmt.Sprintf("big=%d", d.Big)
}

// BooleanStringUnmarshal tests boolean string unmarshal
func BooleanStringUnmarshal() string {
	type Config struct {
		Enabled bool `json:"enabled"`
	}
	jsonData := `{"enabled": true}`
	var c Config
	json.Unmarshal([]byte(jsonData), &c)
	return fmt.Sprintf("enabled=%t", c.Enabled)
}

// FloatPrecision tests float precision
func FloatPrecision() string {
	type Data struct {
		Value float64 `json:"value"`
	}
	d := Data{Value: 0.1 + 0.2}
	data, _ := json.Marshal(d)
	var decoded Data
	json.Unmarshal(data, &decoded)
	epsilon := 0.0001
	closeEnough := decoded.Value > 0.3-epsilon && decoded.Value < 0.3+epsilon
	return fmt.Sprintf("close_enough=%t", closeEnough)
}

// MarshalInterface tests marshaling interface{}
func MarshalInterface() string {
	var data interface{} = map[string]interface{}{
		"name": "test",
		"count": 42,
		"active": true,
	}
	result, _ := json.Marshal(data)
	return fmt.Sprintf("json=%s", string(result))
}

// UnmarshalToInterface tests unmarshaling to interface{}
func UnmarshalToInterface() string {
	jsonData := `{"name":"test","count":42,"nested":{"a":1}}`
	var result interface{}
	json.Unmarshal([]byte(jsonData), &result)
	m := result.(map[string]interface{})
	name := m["name"].(string)
	count := m["count"].(float64)
	return fmt.Sprintf("name=%s,count=%.0f", name, count)
}

// ValidJSON tests json.Valid
func ValidJSON() string {
	valid := json.Valid([]byte(`{"key": "value"}`))
	invalid := json.Valid([]byte(`{invalid}`))
	return fmt.Sprintf("valid=%t,invalid=%t", valid, !invalid)
}

// CompactJSON tests json.Compact
func CompactJSON() string {
	jsonData := `{
		"name": "test",
		"value": 42
	}`
	var buf bytes.Buffer
	json.Compact(&buf, []byte(jsonData))
	return fmt.Sprintf("compact=%t", !strings.Contains(buf.String(), "\n"))
}

// HTMLEscape tests json.HTMLEscape
func HTMLEscape() string {
	src := []byte(`{"html": "<>&"}`)
	var buf bytes.Buffer
	json.HTMLEscape(&buf, src)
	return fmt.Sprintf("escaped=%t", strings.Contains(buf.String(), "\\u003c"))
}

// IndentTests tests json.Indent
func IndentTests() string {
	src := []byte(`{"a":1,"b":2}`)
	var buf bytes.Buffer
	json.Indent(&buf, src, "", "  ")
	return fmt.Sprintf("indented=%t", strings.Contains(buf.String(), "\n"))
}

// MarshalIntFloatString tests marshaling different types
func MarshalIntFloatString() string {
	type Data struct {
		I int     `json:"i"`
		F float64 `json:"f"`
		S string  `json:"s"`
		B bool    `json:"b"`
	}
	d := Data{I: 42, F: 3.14, S: "hello", B: true}
	data, _ := json.Marshal(d)
	return fmt.Sprintf("json=%s", string(data))
}

// TagNameCase tests tag name case sensitivity
func TagNameCase() string {
	type Data struct {
		Lower string `json:"lower"`
		Upper string `json:"UPPER"`
		Mixed string `json:"MixedCase"`
	}
	d := Data{Lower: "a", Upper: "b", Mixed: "c"}
	data, _ := json.Marshal(d)
	return fmt.Sprintf("json=%s", string(data))
}

// AnonymousStructMarshal tests anonymous struct marshal
func AnonymousStructMarshal() string {
	data := struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{
		Name:  "anon",
		Value: 42,
	}
	result, _ := json.Marshal(data)
	return fmt.Sprintf("json=%s", string(result))
}

// NullValueMarshal tests null value handling
func NullValueMarshal() string {
	type Data struct {
		Ptr *int `json:"ptr"`
	}
	d := Data{Ptr: nil}
	data, _ := json.Marshal(d)
	return fmt.Sprintf("json=%s", string(data))
}

// ArrayMarshalUnmarshal tests array (not slice) marshal
func ArrayMarshalUnmarshal() string {
	type Data struct {
		Arr [3]int `json:"arr"`
	}
	d := Data{Arr: [3]int{1, 2, 3}}
	data, _ := json.Marshal(d)
	return fmt.Sprintf("json=%s", string(data))
}

// ByteSliceMarshal tests []byte marshal (base64)
func ByteSliceMarshal() string {
	type Data struct {
		Data []byte `json:"data"`
	}
	d := Data{Data: []byte("hello world")}
	result, _ := json.Marshal(d)
	return fmt.Sprintf("starts_with_quote=%t", result[0] == '"')
}

// StringPointerUnmarshal tests string pointer unmarshal
func StringPointerUnmarshal() string {
	type Data struct {
		Name *string `json:"name"`
	}
	jsonData := `{"name": "test"}`
	var d Data
	json.Unmarshal([]byte(jsonData), &d)
	hasName := d.Name != nil && *d.Name == "test"
	return fmt.Sprintf("has_name=%t", hasName)
}

// IntPointerUnmarshal tests int pointer unmarshal
func IntPointerUnmarshal() string {
	type Data struct {
		Count *int `json:"count"`
	}
	jsonData := `{"count": 42}`
	var d Data
	json.Unmarshal([]byte(jsonData), &d)
	hasCount := d.Count != nil && *d.Count == 42
	return fmt.Sprintf("has_count=%t", hasCount)
}

// ZeroValueMarshal tests zero value marshal
func ZeroValueMarshal() string {
	type Data struct {
		S string `json:"s"`
		I int    `json:"i"`
		B bool   `json:"b"`
	}
	d := Data{} // All zero values
	data, _ := json.Marshal(d)
	return fmt.Sprintf("json=%s", string(data))
}

// EmbeddedStructMarshal tests embedded struct marshal
func EmbeddedStructMarshal() string {
	type Inner struct {
		Value int `json:"value"`
	}
	type Outer struct {
		Inner
		Name string `json:"name"`
	}
	o := Outer{Inner: Inner{Value: 42}, Name: "test"}
	data, _ := json.Marshal(o)
	return fmt.Sprintf("json=%s", string(data))
}

// MapStringInterfaceMarshal tests map[string]interface{} marshal
func MapStringInterfaceMarshal() string {
	m := map[string]interface{}{
		"str":   "hello",
		"num":   42,
		"float": 3.14,
		"bool":  true,
		"null":  nil,
	}
	data, _ := json.Marshal(m)
	return fmt.Sprintf("length=%d", len(data))
}

// StructWithJSONTagDash tests struct with JSON tag "-"
func StructWithJSONTagDash() string {
	type Data struct {
		Keep    string `json:"keep"`
		Ignore  string `json:"-"`
		Ignore2 int    `json:"-"`
	}
	d := Data{Keep: "visible", Ignore: "hidden", Ignore2: 123}
	data, _ := json.Marshal(d)
	return fmt.Sprintf("contains_hidden=%t", strings.Contains(string(data), "hidden"))
}

// ParseStringToFloat tests parsing string as number
func ParseStringToFloat() string {
	s := "3.14159"
	f, _ := strconv.ParseFloat(s, 64)
	return fmt.Sprintf("pi=%.5f", f)
}

// ParseIntTests tests parsing int
func ParseIntTests() string {
	s := "42"
	i, _ := strconv.Atoi(s)
	return fmt.Sprintf("int=%d", i)
}

// ItoaTests tests int to string
func ItoaTests() string {
	i := 42
	s := strconv.Itoa(i)
	return fmt.Sprintf("str=%s", s)
}

// FormatFloatTests tests formatting float
func FormatFloatTests() string {
	f := 3.14159
	s := strconv.FormatFloat(f, 'f', 2, 64)
	return fmt.Sprintf("formatted=%s", s)
}

// QuoteString tests quoting string
func QuoteString() string {
	s := "hello\nworld"
	q := strconv.Quote(s)
	return fmt.Sprintf("quoted=%s", q)
}

// UnquoteString tests unquoting string
func UnquoteString() string {
	q := `"hello\nworld"`
	s, _ := strconv.Unquote(q)
	return fmt.Sprintf("unquoted=%s", s)
}

// ParseBoolTests tests parsing bool
func ParseBoolTests() string {
	trues := []string{"true", "TRUE", "True", "1"}
	falses := []string{"false", "FALSE", "False", "0"}
	results := fmt.Sprintf("trues=%d,falses=%d", len(trues), len(falses))
	return results
}

// AppendIntTests tests appending int
func AppendIntTests() string {
	buf := []byte("value=")
	buf = strconv.AppendInt(buf, 42, 10)
	return string(buf)
}

// AppendFloatTests tests appending float
func AppendFloatTests() string {
	buf := []byte("pi=")
	buf = strconv.AppendFloat(buf, 3.14, 'f', 2, 64)
	return string(buf)
}

// AppendBoolTests tests appending bool
func AppendBoolTests() string {
	buf := []byte("active=")
	buf = strconv.AppendBool(buf, true)
	return string(buf)
}

// AppendQuoteTests tests appending quoted string
func AppendQuoteTests() string {
	buf := []byte("msg=")
	buf = strconv.AppendQuote(buf, "hello")
	return string(buf)
}

// IsPrintTests tests IsPrint
func IsPrintTests() string {
	printable := strconv.IsPrint('a')
	notPrintable := strconv.IsPrint('\x01')
	return fmt.Sprintf("printable=%t,not_printable=%t", printable, !notPrintable)
}

// CanBackquoteTests tests CanBackquote
func CanBackquoteTests() string {
	can := strconv.CanBackquote("hello")
	cannot := strconv.CanBackquote("hello`world")
	return fmt.Sprintf("can=%t,cannot=%t", can, !cannot)
}
