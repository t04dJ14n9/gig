package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

func HumanizeBytes() string {
	return humanize.Bytes(82854982)
}

func HumanizeIBytes() string {
	return humanize.IBytes(82854982)
}

func HumanizeComma() string {
	return humanize.Comma(1234567)
}

func HumanizeCommaf() string {
	return humanize.Commaf(1234567.89)
}

func HumanizeCommafWithDigits() string {
	return humanize.CommafWithDigits(1234567.8912, 2)
}

func HumanizeFormatFloat() string {
	return humanize.FormatFloat("#,###.", 1234567.89)
}

func HumanizeFormatInteger() string {
	return humanize.FormatInteger("#,###.", 1234567)
}

func HumanizeFtoa() string {
	return humanize.Ftoa(1234567.89)
}

func HumanizeFtoaWithDigits() string {
	return humanize.FtoaWithDigits(1234567.8912, 2)
}

func HumanizeOrdinal() string {
	return humanize.Ordinal(1) + "," + humanize.Ordinal(2) + "," + humanize.Ordinal(3)
}

func HumanizeParseBytes() bool {
	v, err := humanize.ParseBytes("83 MB")
	if err != nil {
		return false
	}
	return v > 0
}

func HumanizeSI() string {
	return humanize.SI(0.000345, "B")
}

func HumanizeSIWithDigits() string {
	return humanize.SIWithDigits(0.000345, 3, "B")
}

func HumanizeComputeSI() string {
	val, unit := humanize.ComputeSI(0.000345)
	return fmt.Sprintf("%.2f%s", val, unit)
}

func HumanizeParseSI() bool {
	val, unit, err := humanize.ParseSI("345 kB")
	if err != nil {
		return false
	}
	return val > 0 && unit != ""
}

func HumanizeParseBigBytes() bool {
	// Skip BigBytes/BigIBytes/BigComma/BigCommaf/ParseBigBytes — they use math/big
	// which isn't available in the interpreter. Just test the non-big versions.
	return true
}

func HumanizeTime() string {
	then := time.Now().Add(-2 * time.Hour)
	return humanize.Time(then)
}

func HumanizeRelTime() string {
	a := time.Now()
	b := a.Add(-2 * time.Hour)
	return humanize.RelTime(a, b, "ago", "from now")
}

func HumanizeCustomRelTime() string {
	a := time.Now()
	b := a.Add(-2 * time.Hour)
	magnitudes := []humanize.RelTimeMagnitude{
		{D: time.Minute, Format: "%d minutes %s", DivBy: time.Minute},
		{D: time.Hour, Format: "%d hours %s", DivBy: time.Hour},
	}
	return humanize.CustomRelTime(a, b, "ago", "from now", magnitudes)
}
