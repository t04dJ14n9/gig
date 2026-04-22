package main

import (
	"fmt"
	gjson "github.com/tidwall/gjson"
)

// GjsonGetName extracts "name" field from JSON string
func GjsonGetName(jsonStr string) string {
	return gjson.Get(jsonStr, "name").String()
}

// GjsonGetAge extracts "age" field as int64
func GjsonGetAge(jsonStr string) int64 {
	return gjson.Get(jsonStr, "age").Int()
}

// GjsonGetNested extracts nested "user.name" path
func GjsonGetNested(jsonStr string) string {
	return gjson.Get(jsonStr, "user.name").String()
}

// GjsonArrayAccess extracts array element via "items.0.name"
func GjsonArrayAccess(jsonStr string) string {
	return gjson.Get(jsonStr, "items.0.name").String()
}

// GjsonExists checks if a key exists
func GjsonExists(jsonStr string) bool {
	return gjson.Get(jsonStr, "name").Exists()
}

// GjsonNotExists checks that a nonexistent key returns false
func GjsonNotExists(jsonStr string) bool {
	return gjson.Get(jsonStr, "nonexistent").Exists()
}

// GjsonBoolValue extracts a bool field
func GjsonBoolValue(jsonStr string) bool {
	return gjson.Get(jsonStr, "active").Bool()
}

// GjsonFloatValue extracts a float64 field
func GjsonFloatValue(jsonStr string) float64 {
	return gjson.Get(jsonStr, "price").Float()
}

// GjsonUintValue extracts a uint64 field
func GjsonUintValue(jsonStr string) uint64 {
	return gjson.Get(jsonStr, "count").Uint()
}

// GjsonGetPath extracts value by dynamic path
func GjsonGetPath(jsonStr string, path string) string {
	return gjson.Get(jsonStr, path).String()
}

// GjsonGetMany extracts multiple fields and concatenates
func GjsonGetMany(jsonStr string) string {
	results := gjson.GetMany(jsonStr, "name", "age")
	return results[0].String() + ":" + results[1].String()
}

// GjsonValid checks if JSON is valid
func GjsonValid(jsonStr string) bool {
	return gjson.Valid(jsonStr)
}

// GjsonParseAndGet parses JSON then gets a nested field
func GjsonParseAndGet(jsonStr string) string {
	parsed := gjson.Parse(jsonStr)
	return parsed.Get("name").String()
}

// GjsonDeepNested tests deep nested path access
func GjsonDeepNested(jsonStr string) string {
	return gjson.Get(jsonStr, "data.items.0.subItems.1.value").String()
}

// GjsonMultiplePaths extracts multiple typed values and formats
func GjsonMultiplePaths(jsonStr string) string {
	a := gjson.Get(jsonStr, "a").String()
	b := gjson.Get(jsonStr, "b").Int()
	c := gjson.Get(jsonStr, "c").Bool()
	return fmt.Sprintf("%s:%d:%v", a, b, c)
}

// GjsonGetArrayLength returns length of a JSON array
func GjsonGetArrayLength(jsonStr string) int64 {
	arr := gjson.Get(jsonStr, "items").Array()
	return int64(len(arr))
}

// GjsonIsArray checks if a field is a JSON array
func GjsonIsArray(jsonStr string) bool {
	return gjson.Get(jsonStr, "items").IsArray()
}

// GjsonIsObject checks if a field is a JSON object
func GjsonIsObject(jsonStr string) bool {
	return gjson.Get(jsonStr, "user").IsObject()
}

// GjsonMapValues extracts map values from JSON object
func GjsonMapValues(jsonStr string) string {
	m := gjson.Get(jsonStr, "user").Map()
	name := m["name"].String()
	return name
}
