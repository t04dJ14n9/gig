package main

import (
	cast "github.com/spf13/cast"
)

// CastToString tests basic type conversion to string
func CastToString(v interface{}) string {
	return cast.ToString(v)
}

// CastToInt tests conversion to int — returns int, not int64
func CastToInt(v interface{}) int {
	return cast.ToInt(v)
}

// CastToFloat64 tests conversion to float64
func CastToFloat64(v interface{}) float64 {
	return cast.ToFloat64(v)
}

// CastToBool tests conversion to bool
func CastToBool(v interface{}) bool {
	return cast.ToBool(v)
}
