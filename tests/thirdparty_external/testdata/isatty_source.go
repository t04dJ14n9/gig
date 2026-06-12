package main

import isatty "github.com/mattn/go-isatty"

func IsattyIsTerminal() bool {
	return isatty.IsTerminal(0)
}

func IsattyIsCygwinTerminal() bool {
	return isatty.IsCygwinTerminal(0)
}
