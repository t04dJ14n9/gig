package main

import "github.com/mattn/go-runewidth"

func RunewidthRuneWidth() int {
	return runewidth.RuneWidth('界')
}

func RunewidthStringWidth() int {
	return runewidth.StringWidth("Hello, 世界")
}

func RunewidthFillLeft() string {
	return runewidth.FillLeft("1234", 10)
}

func RunewidthFillRight() string {
	return runewidth.FillRight("1234", 10)
}

func RunewidthTruncate() string {
	return runewidth.Truncate("Hello, 世界", 8, "...")
}

func RunewidthWrap() string {
	return runewidth.Wrap("Hello, 世界!", 8)
}

func RunewidthIsAmbiguousWidth() bool {
	return runewidth.IsAmbiguousWidth('á')
}
