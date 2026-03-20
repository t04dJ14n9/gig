package main

import "encoding/csv"

func Test(w *csv.Writer) int {
	w.Write([]string{"name", "age"})
	w.Write([]string{"Alice", "30"})
	w.Flush()
	if w.Error() != nil {
		return -1
	}
	return 1
}
