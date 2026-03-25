package peephole

import (
	"testing"

	"git.woa.com/youngjin/gig/model/bytecode"
)

func TestMatchOp(t *testing.T) {
	code := []byte{byte(bytecode.OpAdd), byte(bytecode.OpSub), byte(bytecode.OpMul)}

	if !MatchOp(code, 0, bytecode.OpAdd) {
		t.Error("MatchOp(code, 0, OpAdd) = false, want true")
	}
	if MatchOp(code, 0, bytecode.OpSub) {
		t.Error("MatchOp(code, 0, OpSub) = true, want false")
	}
	if !MatchOp(code, 1, bytecode.OpSub) {
		t.Error("MatchOp(code, 1, OpSub) = false, want true")
	}
	if !MatchOp(code, 2, bytecode.OpMul) {
		t.Error("MatchOp(code, 2, OpMul) = false, want true")
	}

	// Out of bounds
	if MatchOp(code, 3, bytecode.OpAdd) {
		t.Error("MatchOp(code, 3, OpAdd) = true, want false")
	}
	// Note: MatchOp does not handle negative offsets, so we don't test that case
}

func TestMake1Op(t *testing.T) {
	result := Make1Op(bytecode.OpLocal, 0x1234)

	if len(result) != 3 {
		t.Fatalf("Make1Op length = %d, want 3", len(result))
	}
	if result[0] != byte(bytecode.OpLocal) {
		t.Errorf("result[0] = 0x%02x, want 0x%02x", result[0], bytecode.OpLocal)
	}
	if result[1] != 0x12 || result[2] != 0x34 {
		t.Errorf("result[1:3] = [0x%02x, 0x%02x], want [0x12, 0x34]", result[1], result[2])
	}
}

func TestMake2Op(t *testing.T) {
	result := Make2Op(bytecode.OpAdd, 0x1234, 0x5678)

	if len(result) != 5 {
		t.Fatalf("Make2Op length = %d, want 5", len(result))
	}
	if result[0] != byte(bytecode.OpAdd) {
		t.Errorf("result[0] = 0x%02x, want 0x%02x", result[0], bytecode.OpAdd)
	}
	if result[1] != 0x12 || result[2] != 0x34 {
		t.Errorf("result[1:3] = [0x%02x, 0x%02x], want [0x12, 0x34]", result[1], result[2])
	}
	if result[3] != 0x56 || result[4] != 0x78 {
		t.Errorf("result[3:5] = [0x%02x, 0x%02x], want [0x56, 0x78]", result[3], result[4])
	}
}

func TestMake3Op(t *testing.T) {
	result := Make3Op(bytecode.OpSub, 0x1111, 0x2222, 0x3333)

	if len(result) != 7 {
		t.Fatalf("Make3Op length = %d, want 7", len(result))
	}
	if result[0] != byte(bytecode.OpSub) {
		t.Errorf("result[0] = 0x%02x, want 0x%02x", result[0], bytecode.OpSub)
	}
	if result[1] != 0x11 || result[2] != 0x11 {
		t.Errorf("result[1:3] = [0x%02x, 0x%02x], want [0x11, 0x11]", result[1], result[2])
	}
	if result[3] != 0x22 || result[4] != 0x22 {
		t.Errorf("result[3:5] = [0x%02x, 0x%02x], want [0x22, 0x22]", result[3], result[4])
	}
	if result[5] != 0x33 || result[6] != 0x33 {
		t.Errorf("result[5:7] = [0x%02x, 0x%02x], want [0x33, 0x33]", result[5], result[6])
	}
}

func TestPatternsReturnsRegisteredPatterns(t *testing.T) {
	// Patterns() should return the registered patterns
	patterns := Patterns()
	if patterns == nil {
		t.Error("Patterns() returned nil")
	}
	// The exact number depends on how many patterns are registered
	// We just verify it's not nil and has expected behavior
	t.Logf("Number of registered patterns: %d", len(patterns))
}

func TestRegister(t *testing.T) {
	initialLen := len(Patterns())

	// Create a simple test pattern
	testPattern := &testPatternImpl{}

	Register(testPattern)

	if len(Patterns()) != initialLen+1 {
		t.Errorf("After Register: len(Patterns()) = %d, want %d", len(Patterns()), initialLen+1)
	}
}

// testPatternImpl is a simple pattern for testing.
type testPatternImpl struct{}

func (p *testPatternImpl) Match(code []byte, i int) (consumed int, newBytes []byte, ok bool) {
	return 0, nil, false
}
