package main

import (
	"fmt"
	"github.com/jinzhu/now"
	"time"
)

func NowBeginningOfDay() string {
	t := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	bod := now.With(t).BeginningOfDay()
	return fmt.Sprintf("%04d-%02d-%02d", bod.Year(), bod.Month(), bod.Day())
}

func NowEndOfDay() string {
	t := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	eod := now.With(t).EndOfDay()
	return fmt.Sprintf("%02d:%02d:%02d", eod.Hour(), eod.Minute(), eod.Second())
}
