package main

import (
	"fmt"

	"github.com/relvacode/iso8601"
)

func Iso8601Parse() string {
	t, err := iso8601.Parse([]byte("2024-06-15T10:30:00Z"))
	if err != nil {
		return err.Error()
	}
	return t.Format("2006-01-02")
}

func Iso8601ParseString() string {
	t, err := iso8601.ParseString("2024-06-15T10:30:00+08:00")
	if err != nil {
		return err.Error()
	}
	return t.Format("15:04:05")
}

func Iso8601ParseWithOffset() string {
	t, err := iso8601.ParseString("2024-06-15T10:30:00+09:00")
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%d", t.Hour())
}
