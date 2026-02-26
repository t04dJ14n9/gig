// Package time registers the Go standard library time package.
package time

import (
	"time"

	"gig/importer"
	"gig/value"
)

func init() {
	pkg := importer.RegisterPackage("time", "time")

	// Time creation
	pkg.AddFunction("Now", time.Now, "", directNow)
	pkg.AddFunction("Date", time.Date, "", nil)
	pkg.AddFunction("Unix", time.Unix, "", nil)
	pkg.AddFunction("UnixMilli", time.UnixMilli, "", nil)
	pkg.AddFunction("UnixMicro", time.UnixMicro, "", nil)
	pkg.AddFunction("Parse", time.Parse, "", nil)
	pkg.AddFunction("ParseInLocation", time.ParseInLocation, "", nil)

	// Duration
	pkg.AddFunction("ParseDuration", time.ParseDuration, "", nil)
	pkg.AddFunction("Since", time.Since, "", nil)
	pkg.AddFunction("Until", time.Until, "", nil)
	pkg.AddFunction("Sleep", time.Sleep, "", nil)
	pkg.AddFunction("Tick", time.Tick, "", nil)
	pkg.AddFunction("After", time.After, "", nil)
	pkg.AddFunction("AfterFunc", time.AfterFunc, "", nil)
	pkg.AddFunction("NewTimer", time.NewTimer, "", nil)
	pkg.AddFunction("NewTicker", time.NewTicker, "", nil)

	// Duration constants
	pkg.AddConstant("Nanosecond", time.Nanosecond, "nanosecond duration")
	pkg.AddConstant("Microsecond", time.Microsecond, "microsecond duration")
	pkg.AddConstant("Millisecond", time.Millisecond, "millisecond duration")
	pkg.AddConstant("Second", time.Second, "second duration")
	pkg.AddConstant("Minute", time.Minute, "minute duration")
	pkg.AddConstant("Hour", time.Hour, "hour duration")

	// Layout constants
	pkg.AddConstant("Layout", time.Layout, "reference layout")
	pkg.AddConstant("ANSIC", time.ANSIC, "ANSIC format")
	pkg.AddConstant("UnixDate", time.UnixDate, "Unix date format")
	pkg.AddConstant("RubyDate", time.RubyDate, "Ruby date format")
	pkg.AddConstant("RFC822", time.RFC822, "RFC822 format")
	pkg.AddConstant("RFC822Z", time.RFC822Z, "RFC822 with zone")
	pkg.AddConstant("RFC850", time.RFC850, "RFC850 format")
	pkg.AddConstant("RFC1123", time.RFC1123, "RFC1123 format")
	pkg.AddConstant("RFC1123Z", time.RFC1123Z, "RFC1123 with zone")
	pkg.AddConstant("RFC3339", time.RFC3339, "RFC3339 format")
	pkg.AddConstant("RFC3339Nano", time.RFC3339Nano, "RFC3339 with nanoseconds")
	pkg.AddConstant("Kitchen", time.Kitchen, "kitchen format")
	pkg.AddConstant("Stamp", time.Stamp, "stamp format")
	pkg.AddConstant("StampMilli", time.StampMilli, "stamp with milliseconds")
	pkg.AddConstant("StampMicro", time.StampMicro, "stamp with microseconds")
	pkg.AddConstant("StampNano", time.StampNano, "stamp with nanoseconds")
	pkg.AddConstant("DateTime", time.DateTime, "datetime format")
	pkg.AddConstant("DateOnly", time.DateOnly, "date only format")
	pkg.AddConstant("TimeOnly", time.TimeOnly, "time only format")

	// Types
	pkg.AddType("Time", nil, "time instant")
	pkg.AddType("Duration", nil, "time duration")
	pkg.AddType("Location", nil, "time location")
	pkg.AddType("Timer", nil, "timer")
	pkg.AddType("Ticker", nil, "ticker")
	pkg.AddType("Month", nil, "month")
	pkg.AddType("Weekday", nil, "weekday")

	// UTC and Local
	pkg.AddVariable("UTC", &time.UTC, "UTC location")
	pkg.AddVariable("Local", &time.Local, "local location")
}

func directNow(args []value.Value) value.Value {
	return value.FromInterface(time.Now())
}
