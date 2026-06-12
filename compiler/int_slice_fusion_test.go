package compiler

import (
	"testing"

	"github.com/t04dJ14n9/gig/importer"
	"github.com/t04dJ14n9/gig/model/bytecode"
)

func TestBubbleSortIntSliceAccessesAreFused(t *testing.T) {
	src := `package main

func BubbleSort() int {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = 100 - i
	}
	n := len(s)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if s[j] > s[j+1] {
				tmp := s[j]
				s[j] = s[j+1]
				s[j+1] = tmp
			}
		}
	}
	return s[0] + s[99]
}
`
	prog, err := Build(src, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	fn := prog.Program.Functions["BubbleSort"]
	if fn == nil {
		t.Fatal("BubbleSort function was not compiled")
	}

	counts := countBytecodeOps(fn.Instructions)
	if got := counts[bytecode.OpIntSliceGet]; got == 0 {
		t.Fatalf("OpIntSliceGet count = 0, want at least one fused int-slice read")
	}
	if got := counts[bytecode.OpIntSliceSet]; got == 0 {
		t.Fatalf("OpIntSliceSet count = 0, want at least one fused int-slice write")
	}
	if got := counts[bytecode.OpMakeSlice]; got != 1 {
		t.Fatalf("OpMakeSlice count = %d, want 1 native int-slice allocation", got)
	}
	if got := counts[bytecode.OpNew]; got != 0 {
		t.Fatalf("OpNew count = %d, want 0 synthetic int-array allocations", got)
	}
	if got := counts[bytecode.OpSlice]; got != 0 {
		t.Fatalf("OpSlice count = %d, want 0 synthetic int-array slices", got)
	}
}

func countBytecodeOps(code []byte) map[bytecode.OpCode]int {
	counts := make(map[bytecode.OpCode]int)
	for i := 0; i < len(code); {
		op := bytecode.OpCode(code[i])
		counts[op]++
		i += 1 + bytecode.OperandWidth(op)
	}
	return counts
}
