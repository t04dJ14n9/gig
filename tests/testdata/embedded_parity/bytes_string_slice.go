package main

import "fmt"

// Result checks byte/string/slice conversion stability.
func Result() string {
	b := []byte("go")
	s := string(b)
	copied := append([]byte(nil), b...)
	copied[0] = 'z'
	r := []rune(s)

	return fmt.Sprintf("%s:%d:%v:%q:%d", string(copied), len(r), r, string([]byte{b[0], b[1]}), int(r[0]))
}
