package thirdparty

import (
	"bytes"
	"encoding/json"
)

// ============================================================================
// JSON COMPLEX STRUCTURES
// ============================================================================

// JsonWithNestedStruct tests marshaling nested struct.
func JsonWithNestedStruct() int {
	type Address struct {
		City    string
		ZipCode string
	}
	type Person struct {
		Name    string
		Address Address
	}

	p := Person{
		Name:    "John",
		Address: Address{City: "NYC", ZipCode: "10001"},
	}
	b, _ := json.Marshal(p)
	return len(b)
}

// JsonWithSliceField tests marshaling struct with slice.
func JsonWithSliceField() int {
	type Data struct {
		Values []int
		Names  []string
	}
	d := Data{Values: []int{1, 2, 3}, Names: []string{"a", "b"}}
	b, _ := json.Marshal(d)
	return len(b)
}

// JsonWithMapField tests marshaling struct with map.
func JsonWithMapField() int {
	type Data struct {
		Metadata map[string]string
	}
	d := Data{Metadata: map[string]string{"key": "value"}}
	b, _ := json.Marshal(d)
	return len(b)
}

// JsonUnmarshalToInterface tests unmarshal to interface{}.
func JsonUnmarshalToInterface() int {
	var data interface{}
	json.Unmarshal([]byte(`{"key":"value","num":42}`), &data)
	m := data.(map[string]interface{})
	return len(m)
}

// JsonStreamDecoder tests streaming JSON decode.
func JsonStreamDecoder() int {
	data := `{"a":1}{"b":2}`
	decoder := json.NewDecoder(bytes.NewReader([]byte(data)))
	count := 0
	for {
		var v map[string]int
		if err := decoder.Decode(&v); err != nil {
			break
		}
		count++
	}
	return count
}
