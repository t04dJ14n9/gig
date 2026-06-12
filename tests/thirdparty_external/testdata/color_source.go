package main

import "github.com/fatih/color"

func ColorRedString() string {
	return color.RedString("hello")
}

func ColorBlueString() string {
	return color.BlueString("world")
}

func ColorNew() string {
	c := color.New(color.FgCyan)
	if c == nil {
		return "nil"
	}
	return "ok"
}

func ColorNewAdd() string {
	c := color.New(color.FgRed).Add(color.Bold)
	if c == nil {
		return "nil"
	}
	return "ok"
}
