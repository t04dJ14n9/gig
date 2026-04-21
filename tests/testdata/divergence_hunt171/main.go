package divergence_hunt171

import (
	"fmt"
	"time"
)

// ============================================================================
// Round 171: Time package operations
// ============================================================================

func TimeNow() string {
	t := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	return t.Format(time.RFC3339)
}

func TimeAdd() string {
	t := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	later := t.Add(time.Hour * 2)
	return later.Format("15:04")
}

func TimeSub() string {
	t1 := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	d := t1.Sub(t2)
	return fmt.Sprintf("%d", int(d.Hours()))
}

func TimeBeforeAfter() string {
	t1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	if t1.Before(t2) && t2.After(t1) {
		return "true"
	}
	return "false"
}

func TimeUnix() string {
	t := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	return fmt.Sprintf("%d", t.Unix())
}

func TimeYearMonthDay() string {
	t := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func TimeHourMinuteSecond() string {
	t := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

func TimeWeekday() string {
	t := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Monday
	return t.Weekday().String()
}

func TimeTruncate() string {
	t := time.Date(2024, 1, 15, 14, 37, 45, 0, time.UTC)
	truncated := t.Truncate(time.Hour)
	return truncated.Format("15:04:05")
}

func DurationString() string {
	d := time.Hour*2 + time.Minute*30
	return d.String()
}

func DurationHours() string {
	d := time.Hour*3 + time.Minute*30
	return fmt.Sprintf("%.1f", d.Hours())
}

func DurationMinutes() string {
	d := time.Hour + time.Minute*30
	return fmt.Sprintf("%.0f", d.Minutes())
}
