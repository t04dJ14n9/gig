package main

import (
	sjson "github.com/tidwall/sjson"
)

// SjsonSetName sets a string value at a key
func SjsonSetName(jsonStr string, name string) string {
	result, _ := sjson.Set(jsonStr, "name", name)
	return result
}

// SjsonSetAge sets an int value at a key
func SjsonSetAge(jsonStr string, age int) string {
	result, _ := sjson.Set(jsonStr, "age", age)
	return result
}

// SjsonSetNested sets a nested path value
func SjsonSetNested(jsonStr string, city string) string {
	result, _ := sjson.Set(jsonStr, "address.city", city)
	return result
}

// SjsonSetBool sets a bool value
func SjsonSetBool(jsonStr string, active bool) string {
	result, _ := sjson.Set(jsonStr, "active", active)
	return result
}

// SjsonSetFloat sets a float64 value
func SjsonSetFloat(jsonStr string, price float64) string {
	result, _ := sjson.Set(jsonStr, "price", price)
	return result
}

// SjsonSetNull sets null at a key
func SjsonSetNull(jsonStr string) string {
	result, _ := sjson.Set(jsonStr, "removed", nil)
	return result
}

// SjsonSetRaw sets raw JSON at a key
func SjsonSetRaw(jsonStr string, raw string) string {
	result, _ := sjson.SetRaw(jsonStr, "extra", raw)
	return result
}

// SjsonDeleteField deletes a top-level field
func SjsonDeleteField(jsonStr string, field string) string {
	result, _ := sjson.Delete(jsonStr, field)
	return result
}

// SjsonDeleteNested deletes a nested field
func SjsonDeleteNested(jsonStr string) string {
	result, _ := sjson.Delete(jsonStr, "address.zip")
	return result
}
