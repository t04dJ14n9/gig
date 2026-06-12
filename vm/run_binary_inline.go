package vm

import "github.com/t04dJ14n9/gig/model/value"

// These predicates keep compound hot-path guards out of the opcode loop while
// leaving the arithmetic itself in place. They are intentionally tiny so the Go
// compiler can inline them inside run().
func runBothInts(a, b value.Value) bool {
	return a.Kind() == value.KindInt && b.Kind() == value.KindInt
}

func runSameSizedInts(a, b value.Value) bool {
	return runBothInts(a, b) && a.RawSize() == b.RawSize()
}
