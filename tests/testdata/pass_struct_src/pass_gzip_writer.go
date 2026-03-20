package main

import "compress/gzip"

func Test(w *gzip.Writer) int {
	n, err := w.Write([]byte("compressed by gig"))
	if err != nil {
		return -1
	}
	err = w.Close()
	if err != nil {
		return -2
	}
	return n
}
