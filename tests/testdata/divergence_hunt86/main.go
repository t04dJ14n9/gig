package divergence_hunt86

import "encoding/json"

// ============================================================================
// Round 86: JSON marshal/unmarshal edge cases
// ============================================================================

func JsonMarshalBasic() string {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	p := Person{Name: "Alice", Age: 30}
	b, _ := json.Marshal(p)
	return string(b)
}

func JsonUnmarshalBasic() int {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	data := `{"name":"Bob","age":25}`
	var p Person
	json.Unmarshal([]byte(data), &p)
	return p.Age
}

func JsonMarshalSlice() string {
	nums := []int{1, 2, 3}
	b, _ := json.Marshal(nums)
	return string(b)
}

func JsonUnmarshalSlice() int {
	data := `[10,20,30]`
	var nums []int
	json.Unmarshal([]byte(data), &nums)
	return nums[0] + nums[1] + nums[2]
}

func JsonMarshalMap() string {
	m := map[string]int{"a": 1, "b": 2}
	b, _ := json.Marshal(m)
	return string(b)
}

func JsonMarshalNested() string {
	type Outer struct {
		Name  string `json:"name"`
		Inner struct {
			Value int `json:"value"`
		} `json:"inner"`
	}
	o := Outer{Name: "test"}
	o.Inner.Value = 42
	b, _ := json.Marshal(o)
	return string(b)
}

func JsonMarshalBool() string {
	b, _ := json.Marshal(true)
	return string(b)
}

func JsonUnmarshalBool() bool {
	data := `true`
	var v bool
	json.Unmarshal([]byte(data), &v)
	return v
}

func JsonMarshalNull() string {
	var v any
	b, _ := json.Marshal(v)
	return string(b)
}

func JsonMarshalString() string {
	b, _ := json.Marshal("hello")
	return string(b)
}

func JsonUnmarshalString() string {
	data := `"world"`
	var v string
	json.Unmarshal([]byte(data), &v)
	return v
}

func JsonRoundTrip() int {
	type Data struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	orig := Data{X: 10, Y: 20}
	b, _ := json.Marshal(orig)
	var decoded Data
	json.Unmarshal(b, &decoded)
	return decoded.X + decoded.Y
}
