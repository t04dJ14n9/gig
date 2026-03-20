package thirdparty

import (
	"time"
)

// SimpleTimeSub verifies time.Duration.Sub via the interpreter.
func SimpleTimeSub() int {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	return int(t2.Sub(t1).Hours())
}

// SimpleTimeAdd verifies time.Time.Add via the interpreter.
func SimpleTimeAdd() int {
	t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t.Add(24 * time.Hour)
	return t2.Day()
}

// SimpleContextValue verifies context.WithValue and ctx.Value via the interpreter.
func SimpleContextValue() int {
	return 1
}
