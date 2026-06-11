package main

// Result checks defer/recover behavior.
func Result() (out string) {
	out = "start"
	defer func() {
		out = out + ":d1"
	}()
	defer func() {
		if r := recover(); r != nil {
			out = out + ":panic:" + r.(string)
		}
	}()
	panic("boom")
}
