package main

import (
	"time"

	"github.com/hako/durafmt"
)

func DurafmtParse() string {
	d, _ := durafmt.ParseString("1h30m")
	return d.String()
}

func DurafmtParseShort() string {
	d, _ := durafmt.ParseStringShort("1h30m")
	return d.String()
}

func DurafmtDuration() time.Duration {
	d := durafmt.Parse(90 * time.Minute)
	return d.Duration()
}

func DurafmtLimitFirstN() string {
	d := durafmt.Parse(365*24*time.Hour + 2*time.Hour + 30*time.Minute)
	return d.LimitFirstN(1).String()
}
