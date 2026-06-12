package main

import "github.com/axgle/mahonia"

func MahoniaGetCharset() string {
	cs := mahonia.GetCharset("utf-8")
	if cs == nil {
		return "nil"
	}
	return cs.Name
}

func MahoniaGetCharsetGbk() string {
	cs := mahonia.GetCharset("gbk")
	if cs == nil {
		return "nil"
	}
	return cs.Name
}

func MahoniaRegisterCharset() string {
	// Test that RegisterCharset is callable (it's a function, not a method)
	mahonia.RegisterCharset(&mahonia.Charset{Name: "test-charset"})
	cs := mahonia.GetCharset("test-charset")
	if cs == nil {
		return "nil"
	}
	return cs.Name
}
