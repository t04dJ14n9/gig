package main

import (
	"fmt"
	carbon "github.com/dromara/carbon/v2"
)

// CarbonParseDate parses a date string and returns date string
func CarbonParseDate(s string) string {
	c := carbon.Parse(s)
	return c.ToDateString()
}

// CarbonCreateFromDate creates date from year/month/day
func CarbonCreateFromDate(y int, m int, d int) string {
	c := carbon.CreateFromDate(y, m, d)
	return c.ToDateString()
}

// CarbonCreateFromTime creates time from hour/minute/second
func CarbonCreateFromTime(h int, m int, s int) string {
	c := carbon.CreateFromTime(h, m, s)
	return c.ToTimeString()
}

// CarbonCreateFromDateTime creates datetime from all components
func CarbonCreateFromDateTime(y int, mo int, d int, h int, mi int, s int) string {
	c := carbon.CreateFromDateTime(y, mo, d, h, mi, s)
	return c.ToDateTimeString()
}

// CarbonAddDays adds days to a parsed date
func CarbonAddDays(s string, days int) string {
	c := carbon.Parse(s)
	return c.AddDays(days).ToDateString()
}

// CarbonSubDays subtracts days from a parsed date
func CarbonSubDays(s string, days int) string {
	c := carbon.Parse(s)
	return c.SubDays(days).ToDateString()
}

// CarbonAddMonths adds months to a parsed date
func CarbonAddMonths(s string, months int) string {
	c := carbon.Parse(s)
	return c.AddMonths(months).ToDateString()
}

// CarbonAddYears adds years to a parsed date
func CarbonAddYears(s string, years int) string {
	c := carbon.Parse(s)
	return c.AddYears(years).ToDateString()
}

// CarbonAddHours adds hours to a parsed datetime
func CarbonAddHours(s string, hours int) string {
	c := carbon.Parse(s)
	return c.AddHours(hours).ToTimeString()
}

// CarbonAddMinutes adds minutes to a parsed datetime
func CarbonAddMinutes(s string, minutes int) string {
	c := carbon.Parse(s)
	return c.AddMinutes(minutes).ToTimeString()
}

// CarbonStartOfDay returns start of day
func CarbonStartOfDay(s string) string {
	c := carbon.Parse(s)
	return c.StartOfDay().ToDateTimeString()
}

// CarbonEndOfDay returns end of day
func CarbonEndOfDay(s string) string {
	c := carbon.Parse(s)
	return c.EndOfDay().ToDateTimeString()
}

// CarbonStartOfMonth returns start of month
func CarbonStartOfMonth(s string) string {
	c := carbon.Parse(s)
	return c.StartOfMonth().ToDateString()
}

// CarbonEndOfMonth returns end of month
func CarbonEndOfMonth(s string) string {
	c := carbon.Parse(s)
	return c.EndOfMonth().ToDateString()
}

// CarbonStartOfYear returns start of year
func CarbonStartOfYear(s string) string {
	c := carbon.Parse(s)
	return c.StartOfYear().ToDateString()
}

// CarbonEndOfYear returns end of year
func CarbonEndOfYear(s string) string {
	c := carbon.Parse(s)
	return c.EndOfYear().ToDateString()
}

// CarbonIsWeekend checks if date is weekend
func CarbonIsWeekend(s string) bool {
	c := carbon.Parse(s)
	return c.IsWeekend()
}

// CarbonIsLeapYear checks if year is leap year
func CarbonIsLeapYear(s string) bool {
	c := carbon.Parse(s)
	return c.IsLeapYear()
}

// CarbonDayOfWeek returns day of week (1=Monday)
func CarbonDayOfWeek(s string) int {
	c := carbon.Parse(s)
	return c.DayOfWeek()
}

// CarbonDayOfYear returns day of year
func CarbonDayOfYear(s string) int {
	c := carbon.Parse(s)
	return c.DayOfYear()
}

// CarbonMonth returns month number
func CarbonMonth(s string) int {
	c := carbon.Parse(s)
	return c.Month()
}

// CarbonYear returns year
func CarbonYear(s string) int {
	c := carbon.Parse(s)
	return c.Year()
}

// CarbonDaysInMonth returns days in month
func CarbonDaysInMonth(s string) int {
	c := carbon.Parse(s)
	return c.DaysInMonth()
}

// CarbonToTimestamp returns unix timestamp
func CarbonToTimestamp(s string) int64 {
	c := carbon.Parse(s)
	return c.Timestamp()
}

// CarbonToRfc3339 formats as RFC3339
func CarbonToRfc3339(s string) string {
	c := carbon.Parse(s)
	return c.ToRfc3339String()
}

// CarbonToIso8601 formats as ISO8601
func CarbonToIso8601(s string) string {
	c := carbon.Parse(s)
	return c.ToIso8601String()
}

// CarbonLayoutFormat formats with custom layout
func CarbonLayoutFormat(s string, layout string) string {
	c := carbon.Parse(s)
	return c.Layout(layout)
}

// CarbonParseAndFormat uses cast + carbon
func CarbonParseAndFormat(s string) string {
	c := carbon.Parse(s)
	return fmt.Sprintf("%d-%d-%d", c.Year(), c.Month(), c.Day())
}

// CarbonTimestampConversion converts timestamp to string
func CarbonTimestampConversion(s string) string {
	c := carbon.Parse(s)
	return fmt.Sprintf("%d", c.Timestamp())
}
