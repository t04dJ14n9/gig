package main

import "github.com/acarl005/stripansi"

func StripansiStrip() string {
	return stripansi.Strip("\x1b[31mhello\x1b[0m")
}

func StripansiPlain() string {
	return stripansi.Strip("plain text")
}

func StripansiMultiple() string {
	return stripansi.Strip("\x1b[1;32mgreen\x1b[0m \x1b[33myellow\x1b[0m")
}
