package divergence_hunt224

import "fmt"

// ============================================================================
// Round 224: Map with interface keys
// ============================================================================

// MapInterfaceKeyBasic tests basic interface{} key operations
func MapInterfaceKeyBasic() string {
	m := map[interface{}]string{
		1:       "int",
		"hello": "string",
		3.14:    "float64",
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapInterfaceKeyLookup tests interface{} key lookup
func MapInterfaceKeyLookup() string {
	m := map[interface{}]int{
		42:        1,
		"key":     2,
		true:      3,
	}
	v1, _ := m[42]
	v2, _ := m["key"]
	v3, _ := m[true]
	return fmt.Sprintf("sum=%d", v1+v2+v3)
}

// MapInterfaceKeyMixedTypes tests mixed type keys
func MapInterfaceKeyMixedTypes() string {
	type MyKey struct {
		V int
	}
	m := map[interface{}]string{
		10:          "int",
		"str":       "string",
		MyKey{V: 5}: "struct",
		2.5:         "float",
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapInterfaceKeyNil tests nil as interface key
func MapInterfaceKeyNil() string {
	m := map[interface{}]int{}
	var p *int = nil
	m[p] = 100
	v, ok := m[(*int)(nil)]
	return fmt.Sprintf("v=%d,ok=%t", v, ok)
}

// MapInterfaceKeySlicePtr uses slice pointer as key
func MapInterfaceKeySlicePtr() string {
	m := map[interface{}]int{}
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	m[&s1] = 1
	m[&s2] = 2
	return fmt.Sprintf("len=%d", len(m))
}

// MapInterfaceKeyMapPtr uses map pointer as key
func MapInterfaceKeyMapPtr() string {
	m := map[interface{}]string{}
	inner1 := map[string]int{"a": 1}
	inner2 := map[string]int{"a": 1}
	m[&inner1] = "first"
	m[&inner2] = "second"
	return fmt.Sprintf("len=%d", len(m))
}

// MapInterfaceKeyIterate tests iterating over interface{} keyed map
func MapInterfaceKeyIterate() string {
	m := map[interface{}]int{
		1:     10,
		"two": 20,
		3.0:   30,
	}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapInterfaceKeyDelete tests deleting interface{} keys
func MapInterfaceKeyDelete() string {
	m := map[interface{}]string{
		"key1": "val1",
		"key2": "val2",
		100:    "val3",
	}
	delete(m, "key1")
	delete(m, 100)
	return fmt.Sprintf("len=%d", len(m))
}

// MapInterfaceKeyTypeSwitch demonstrates type switching on keys
func MapInterfaceKeyTypeSwitch() string {
	m := map[interface{}]int{
		1:     10,
		"two": 20,
		true:  30,
	}
	intSum, strCount, boolSum := 0, 0, 0
	for k, v := range m {
		switch k.(type) {
		case int:
			intSum += v
		case string:
			strCount++
		case bool:
			boolSum += v
		}
	}
	return fmt.Sprintf("int=%d,str=%d,bool=%d", intSum, strCount, boolSum)
}

// MapInterfaceKeyOverwrite tests overwriting with different types
func MapInterfaceKeyOverwrite() string {
	m := map[interface{}]int{}
	m[1] = 10
	m[int8(1)] = 20
	m[int16(1)] = 30
	return fmt.Sprintf("len=%d", len(m))
}

// MapInterfaceKeyCommaOk tests comma-ok with interface keys
func MapInterfaceKeyCommaOk() string {
	m := map[interface{}]string{"exists": "yes"}
	_, ok1 := m["exists"]
	_, ok2 := m["missing"]
	return fmt.Sprintf("ok1=%t,ok2=%t", ok1, ok2)
}
