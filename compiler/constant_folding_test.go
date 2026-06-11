package compiler

import (
	"context"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/vm"
)

func TestBuildFoldsLocalConstantArithmetic(t *testing.T) {
	prog := compileBuild(t, `func F() int {
	a := 3
	b := 4
	return a*b + 5
}`)

	fn := prog.Functions["F"]
	if fn == nil {
		t.Fatal("F function was not compiled")
	}

	counts := countBytecodeOps(fn.Instructions)
	for _, arithmeticOp := range foldedArithmeticOps() {
		if counts[arithmeticOp] != 0 {
			t.Fatalf("compiled F still has %v after constant folding; counts=%v", arithmeticOp, counts)
		}
	}

	result, err := vm.New(prog).Execute("F", context.Background())
	if err != nil {
		t.Fatalf("Execute(F) error = %v", err)
	}
	if result.Int() != 17 {
		t.Fatalf("F() = %d, want 17", result.Int())
	}
}

func TestBuildFoldsConstantIfBranch(t *testing.T) {
	prog := compileBuild(t, `func F() int {
	a := 3
	b := 4
	if a*b > 10 {
		return 1
	}
	return 2
}`)

	fn := prog.Functions["F"]
	if fn == nil {
		t.Fatal("F function was not compiled")
	}

	counts := countBytecodeOps(fn.Instructions)
	if counts[bytecode.OpJumpTrue] != 0 || counts[bytecode.OpJumpFalse] != 0 {
		t.Fatalf("compiled F still has conditional jumps after constant branch folding; counts=%v", counts)
	}
	if counts[bytecode.OpJump] != 0 {
		t.Fatalf("compiled F still has unconditional jumps after dead-code cleanup; counts=%v", counts)
	}

	result, err := vm.New(prog).Execute("F", context.Background())
	if err != nil {
		t.Fatalf("Execute(F) error = %v", err)
	}
	if result.Int() != 1 {
		t.Fatalf("F() = %d, want 1", result.Int())
	}
}

func foldedArithmeticOps() []bytecode.OpCode {
	return []bytecode.OpCode{
		bytecode.OpAdd,
		bytecode.OpSub,
		bytecode.OpMul,
		bytecode.OpDiv,
		bytecode.OpMod,
		bytecode.OpAddLocalLocal,
		bytecode.OpSubLocalLocal,
		bytecode.OpMulLocalLocal,
		bytecode.OpAddLocalConst,
		bytecode.OpSubLocalConst,
		bytecode.OpAddSetLocal,
		bytecode.OpSubSetLocal,
		bytecode.OpLocalConstAddSetLocal,
		bytecode.OpLocalConstSubSetLocal,
		bytecode.OpLocalConstMulSetLocal,
		bytecode.OpLocalLocalAddSetLocal,
		bytecode.OpLocalLocalSubSetLocal,
		bytecode.OpLocalLocalMulSetLocal,
		bytecode.OpIntLocalConstAddSetLocal,
		bytecode.OpIntLocalConstSubSetLocal,
		bytecode.OpIntLocalConstMulSetLocal,
		bytecode.OpIntLocalLocalAddSetLocal,
		bytecode.OpIntLocalLocalSubSetLocal,
		bytecode.OpIntLocalLocalMulSetLocal,
	}
}
