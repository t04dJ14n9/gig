package main

import "fmt"

// Result checks byte/string/slice conversion stability.
func Result() string {
	b := []byte("go")
	s := string(b)
	copied := append([]byte(nil), b...)
	copied[0] = 'z'
	r := []rune(s)
	buf := make([]byte, 2, 5)
	copy(buf, "hi")
	window := buf[:1:2]
	spread := append([]byte("x"), "yz"...)

	return fmt.Sprintf("%s:%d:%v:%q:%d:%d:%d:%d:%s:%s",
		string(copied), len(r), r, string([]byte{b[0], b[1]}), int(r[0]),
		cap(buf), len(window), cap(window), string(window), string(spread))
}
