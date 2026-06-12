package divergence_hunt111

import (
	"fmt"
	"time"
)

// ============================================================================
// Round 111: Time formatting and parsing
// ============================================================================

func TimeFormat() string {
	t := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	return t.Format("2006-01-02")
}

func TimeParse() string {
	t, _ := time.Parse("2006-01-02", "2024-06-15")
	return fmt.Sprintf("%d", t.Year())
}

func TimeNow() string {
	// Just verify time.Now() returns something reasonable
	t := time.Now()
	return fmt.Sprintf("%v", t.Year() > 2020)
}

func TimeAdd() string {
	t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t.Add(24 * time.Hour)
	return t2.Format("2006-01-02")
}

func TimeSub() string {
	t1 := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	d := t1.Sub(t2)
	return fmt.Sprintf("%d", int(d.Hours()/24))
}

func TimeUnix() string {
	t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return fmt.Sprintf("%v", t.Unix() > 0)
}

func TimeWeekday() string {
	t := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC) // Saturday
	return t.Weekday().String()
}

func TimeBefore() string {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	return fmt.Sprintf("%v", t1.Before(t2))
}

func TimeFormatCustom() string {
	t := time.Date(2024, 3, 5, 14, 30, 0, 0, time.UTC)
	return t.Format("15:04:05")
}

func TimeDateComponents() string {
	t := time.Date(2024, 6, 15, 10, 30, 45, 0, time.UTC)
	return fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second())
}
