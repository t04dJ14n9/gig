package main

import "fmt"

type MyInt int

type AliasByte = byte

type MyBytes []byte

// Result exercises alias handling and named types for conversion paths.
func Result() string {
	var n MyInt = 41
	var b AliasByte = 'A'
	s := MyBytes{byte(b), 98, 99}
	copied := make(MyBytes, len(s))
	copy(copied, s)
	copied[0] = 90

	return fmt.Sprintf("%d:%c:%d:%d:%d:%v", int(n+1), b, len(copied), copied[0], copied, copied == nil)
}
