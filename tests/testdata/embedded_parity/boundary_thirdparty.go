package main

import "example.com/gig-embedded-parity/host"

type ScriptStruct struct {
	Name string
}

func Result() string {
	return host.AcceptAny(ScriptStruct{Name: "x"})
}
