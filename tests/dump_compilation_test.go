package tests

import (
	"fmt"
	"testing"

	"git.woa.com/youngjin/gig"
	"git.woa.com/youngjin/gig/bytecode"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
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
	fmt.Println("=== BYTECODE DUMP ===")
	for name, fn := range bp.Functions {
		fmt.Printf("\n--- Function: %s (NumLocals=%d, NumFreeVars=%d, NumParams=%d) ---\n",
			name, fn.NumLocals, fn.NumFreeVars, fn.NumParams)
		disassemble(fn.Instructions)
	}

	// Also dump functions by index
	fmt.Println("\n=== FUNCTIONS BY INDEX ===")
	for i, fn := range bp.FuncByIndex {
		if fn != nil {
			fmt.Printf("  [%d] %s\n", i, fn.Name)
		}
	}

	// Dump constants
	fmt.Println("\n=== CONSTANTS ===")
	for i, c := range bp.Constants {
		fmt.Printf("  [%d] %v (%T)\n", i, c, c)
	}

	// Execute Filter
	fmt.Println("\n=== EXECUTION ===")
	r, err := prog.Run("Filter", []int{1, 5, 10, 15, 20}, 10)
	if err != nil {
		t.Fatalf("Filter error: %v", err)
	}
	fmt.Printf("Filter([1,5,10,15,20], 10) = %v\n", r)
	result := r.([]int)
	if len(result) != 2 || result[0] != 15 || result[1] != 20 {
		t.Fatalf("Filter wrong: got %v", result)
	}

	// Execute MakeAdder — returns a *vm.Closure, call it via RunWithValues
	r2, err := prog.Run("MakeAdder", 10)
	if err != nil {
		t.Fatalf("MakeAdder error: %v", err)
	}
	fmt.Printf("MakeAdder(10) returned closure: %T\n", r2)

	// Counter returns a closure
	r3, err := prog.Run("Counter")
	if err != nil {
		t.Fatalf("Counter error: %v", err)
	}
	fmt.Printf("Counter() returned closure: %T\n", r3)
}

func disassemble(code []byte) {
	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		name := op.String()
		width := bytecode.OperandWidth(op)

		switch width {
		case 0:
			fmt.Printf("  %04x: %s\n", i, name)
			i++
		case 1:
			if i+1 < len(code) {
				fmt.Printf("  %04x: %s %d\n", i, name, code[i+1])
			}
			i += 2
		case 2:
			if i+2 < len(code) {
				val := uint16(code[i+1])<<8 | uint16(code[i+2])
				fmt.Printf("  %04x: %s %d\n", i, name, val)
			}
			i += 3
		default:
			fmt.Printf("  %04x: %s ???\n", i, name)
			i++
		}
	}
}
