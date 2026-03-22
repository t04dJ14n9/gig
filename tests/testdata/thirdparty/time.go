package thirdparty

import "time"

// TimeNow tests time.Now.
func TimeNow() int {
	now := time.Now()
	return now.Year()
}

// TimeParse tests time.Parse.
func TimeParse() int {
	t, _ := time.Parse("2006-01-02", "2024-03-15")
	return int(t.Month())
}

// TimeFormat tests time.Format.
func TimeFormat() string {
	t, _ := time.Parse("2006-01-02", "2024-03-15")
	return t.Format("January")
}

// TimeAdd tests time.Add.
func TimeAdd() int {
	t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t.Add(24 * time.Hour)
	return t2.Day()
}

// TimeSub tests time.Sub.
func TimeSub() int {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	return int(t2.Sub(t1).Hours())
}

// TimeBefore tests time.Before.
func TimeBefore() int {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	if t2.Before(t1) {
		return 0
	}
	return 1
}

// TimeAfter tests time.After.
func TimeAfter() int {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	if t1.After(t2) {
		return 0
	}
	return 1
}

// TimeEqual tests time.Equal.
func TimeEqual() int {
	t1 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	if t1.Equal(t2) {
		return 1
	}
	return 0
}

// TimeUnix tests time.Unix.
func TimeUnix() int64 {
	t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return t.Unix()
}

// TimeUnixMilli tests time.UnixMilli.
func TimeUnixMilli() int64 {
	t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return t.UnixMilli()
}

// TimeDuration tests time.Duration operations.
func TimeDuration() int {
	d := 5*time.Minute + 30*time.Second
	return int(d.Seconds())
}

// TimeSleep tests time.Sleep (just check it doesn't panic).
func TimeSleep() int {
	time.Sleep(1 * time.Millisecond)
	return 1
}
