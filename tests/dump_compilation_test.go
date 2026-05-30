package tests

import (
	"testing"

	"github.com/t04dJ14n9/gig"
	"github.com/t04dJ14n9/gig/model/bytecode"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

func TestDumpCompilation(t *testing.T) {
	src := `package main

func Filter(nums []int, threshold int) []int {
    result := []int{}
    for _, n := range nums {
        if n > threshold {
            result = append(result, n)
        }
    }
    return result
}

func MakeAdder(base int) func(int) int {
    return func(x int) int {
        return base + x
    }
}

func Counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
`
	prog, err := gig.Build(src)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	bp := prog.InternalProgram()

	// Dump bytecode for each function
	t.Log("=== BYTECODE DUMP ===")
	for name, fn := range bp.Functions {
		t.Logf("--- Function: %s (NumLocals=%d, NumFreeVars=%d, NumParams=%d) ---",
			name, fn.NumLocals, fn.NumFreeVars, fn.NumParams)
		disassemble(t, fn.Instructions)
	}

	// Also dump functions by index
	t.Log("=== FUNCTIONS BY INDEX ===")
	for i, fn := range bp.FuncByIndex {
		if fn != nil {
			t.Logf("[%d] %s", i, fn.Name)
		}
	}

	// Dump constants
	t.Log("=== CONSTANTS ===")
	for i, c := range bp.Constants {
		t.Logf("[%d] %v (%T)", i, c, c)
	}

	// Execute Filter
	t.Log("=== EXECUTION ===")
	r, err := prog.Run("Filter", []int{1, 5, 10, 15, 20}, 10)
	if err != nil {
		t.Fatalf("Filter error: %v", err)
	}
	t.Logf("Filter([1,5,10,15,20], 10) = %v", r)
	result := r.([]int)
	if len(result) != 2 || result[0] != 15 || result[1] != 20 {
		t.Fatalf("Filter wrong: got %v", result)
	}

	// Execute MakeAdder — returns a *vm.Closure, call it via RunWithValues
	r2, err := prog.Run("MakeAdder", 10)
	if err != nil {
		t.Fatalf("MakeAdder error: %v", err)
	}
	t.Logf("MakeAdder(10) returned closure: %T", r2)

	// Counter returns a closure
	r3, err := prog.Run("Counter")
	if err != nil {
		t.Fatalf("Counter error: %v", err)
	}
	t.Logf("Counter() returned closure: %T", r3)
}

func disassemble(t *testing.T, code []byte) {
	t.Helper()
	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		name := op.String()
		width := bytecode.OperandWidth(op)

		switch width {
		case 0:
			t.Logf("%04x: %s", i, name)
			i++
		case 1:
			if i+1 < len(code) {
				t.Logf("%04x: %s %d", i, name, code[i+1])
			}
			i += 2
		case 2:
			if i+2 < len(code) {
				val := uint16(code[i+1])<<8 | uint16(code[i+2])
				t.Logf("%04x: %s %d", i, name, val)
			}
			i += 3
		default:
			t.Logf("%04x: %s ???", i, name)
			i++
		}
	}
}
