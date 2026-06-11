package main

import "github.com/mgutz/ansi"

func AnsiColor() string {
	return ansi.Color("hello", "red")
}

func AnsiColorCode() string {
	return ansi.ColorCode("red+b")
}

func AnsiColorFunc() string {
	red := ansi.ColorFunc("red")
	return red("world")
}

func AnsiReset() string {
	return ansi.ColorCode("reset")
}
