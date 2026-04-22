package gofun

import (
	"testing"

	newgofun "git.code.oa.com/datacenter/onefun/gofun"
	_ "git.code.oa.com/datacenter/onefun/gofun/packages"

	"github.com/t04dJ14n9/gig"
	"github.com/t04dJ14n9/gig/model/value"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

// ============================================================================
// gofun (onefun/gofun) Bug 实际验证测试
// ============================================================================
//
// 新 gofun 基于 SSA 编译 + 寄存器解释器，修复了旧版部分 bug。
// 以下测试验证哪些问题在新版仍存在。
//
// API:
//   gofun.Build(src) -> *Program
//   program.Run(funcName, args...) -> (interface{}, error)
// ============================================================================

// --- Bug #1: 整数字面量溢出 ---
// 旧 gofun: val = int(val.(int64)) 强制转换
// 新 gofun: 基于 SSA，使用 constValue + go/constant，应已修复

func TestGofunVerify_Bug1_IntegerOverflow(t *testing.T) {
	src := `
package main

func Int64Max() int64 {
	return 9223372036854775807
}

func LargeIntSum() int64 {
	return 2147483648 + 2147483648
}
`
	// 新 gofun
	program, err := newgofun.Build(src)
	if err != nil {
		t.Fatalf("gofun Build: %v", err)
	}
	r1, err := program.Run("Int64Max")
	if err != nil {
		t.Logf("gofun Int64Max 错误: %v", err)
	} else {
		t.Logf("gofun Int64Max() = %v (type: %T)", r1, r1)
	}
	r2, err := program.Run("LargeIntSum")
	if err != nil {
		t.Logf("gofun LargeIntSum 错误: %v", err)
	} else {
		t.Logf("gofun LargeIntSum() = %v (期望: 4294967296)", r2)
	}

	// Gig 对比
	gigProg, _ := gig.Build(src)
	gigR1, _ := gigProg.Run("Int64Max")
	gigR2, _ := gigProg.Run("LargeIntSum")
	t.Logf("Gig Int64Max() = %v", gigR1)
	t.Logf("Gig LargeIntSum() = %v", gigR2)
}

// --- Bug #2: make 容量参数 ---
// 新 gofun: 使用 SSA MakeSlice 指令，len/cap 分别获取，应已修复

func TestGofunVerify_Bug2_MakeCapacity(t *testing.T) {
	src := `
package main

func GetLenCap() (int, int) {
	s := make([]int, 5, 10)
	return len(s), cap(s)
}
`
	program, err := newgofun.Build(src)
	if err != nil {
		t.Fatalf("gofun Build: %v", err)
	}
	result, err := program.Run("GetLenCap")
	if err != nil {
		t.Logf("gofun GetLenCap 错误: %v", err)
	} else {
		t.Logf("gofun GetLenCap() = %v (期望: 5, 10)", result)
	}

	// Gig 对比
	gigProg, _ := gig.Build(src)
	gigR, _ := gigProg.Run("GetLenCap")
	if values, ok := gigR.([]value.Value); ok && len(values) == 2 {
		t.Logf("Gig GetLenCap() = %v, %v", values[0].Interface(), values[1].Interface())
	}
}

// --- Bug #3: 短路求值 ---
// 新 gofun: SSA 编译器将 && / || 编译为分支跳转，天然支持短路
// 但需验证是否正确实现

func TestGofunVerify_Bug3_ShortCircuit(t *testing.T) {
	src := `
package main

func ShortCircuitAnd() bool {
	return false && true
}

func SafeNilCheck(ptr *int) int {
	if ptr != nil && *ptr > 0 {
		return *ptr
	}
	return 0
}
`
	program, err := newgofun.Build(src)
	if err != nil {
		t.Fatalf("gofun Build: %v", err)
	}
	r1, err := program.Run("ShortCircuitAnd")
	if err != nil {
		t.Logf("gofun ShortCircuitAnd 错误: %v", err)
	} else {
		t.Logf("gofun ShortCircuitAnd() = %v", r1)
	}
	r2, err := program.Run("SafeNilCheck", (*int)(nil))
	if err != nil {
		t.Logf("gofun SafeNilCheck(nil) 错误: %v", err)
	} else {
		t.Logf("gofun SafeNilCheck(nil) = %v (期望: 0)", r2)
	}

	// Gig 对比
	gigProg, _ := gig.Build(src)
	gigR1, _ := gigProg.Run("ShortCircuitAnd")
	gigR2, _ := gigProg.Run("SafeNilCheck", (*int)(nil))
	t.Logf("Gig ShortCircuitAnd() = %v", gigR1)
	t.Logf("Gig SafeNilCheck(nil) = %v", gigR2)
}

// --- Bug #4: Map 索引存在标志 ---
// 新 gofun: 使用 ssa.Lookup 指令 + CommaOk，应已支持

func TestGofunVerify_Bug4_MapKeyExists(t *testing.T) {
	src := `
package main

func MapKeyExists() (int, bool) {
	m := map[string]int{"a": 1, "b": 0}
	v, ok := m["b"]
	return v, ok
}

func MapKeyMissing() (int, bool) {
	m := map[string]int{"a": 1, "b": 0}
	v, ok := m["c"]
	return v, ok
}
`
	program, err := newgofun.Build(src)
	if err != nil {
		t.Fatalf("gofun Build: %v", err)
	}
	r1, err := program.Run("MapKeyExists")
	if err != nil {
		t.Logf("gofun MapKeyExists 错误: %v", err)
	} else {
		t.Logf("gofun MapKeyExists() = %v (期望: 0, true)", r1)
	}
	r2, err := program.Run("MapKeyMissing")
	if err != nil {
		t.Logf("gofun MapKeyMissing 错误: %v", err)
	} else {
		t.Logf("gofun MapKeyMissing() = %v (期望: 0, false)", r2)
	}

	// Gig 对比
	gigProg, _ := gig.Build(src)
	gigR1, _ := gigProg.Run("MapKeyExists")
	gigR2, _ := gigProg.Run("MapKeyMissing")
	if values, ok := gigR1.([]value.Value); ok && len(values) == 2 {
		t.Logf("Gig MapKeyExists() = %v, %v", values[0].Interface(), values[1].Interface())
	}
	if values, ok := gigR2.([]value.Value); ok && len(values) == 2 {
		t.Logf("Gig MapKeyMissing() = %v, %v", values[0].Interface(), values[1].Interface())
	}
}

// --- Bug #5: 切片边界检查 ---
// 新 gofun: SSA 的 Slice 指令直接调用 reflect.Slice，由 Go 运行时检查边界

func TestGofunVerify_Bug5_SliceBoundsCheck(t *testing.T) {
	src := `
package main

func ValidSlice() []int {
	s := []int{1, 2, 3, 4, 5}
	return s[1:3]
}

func EmptySlice() []int {
	s := []int{1, 2, 3, 4, 5}
	return s[2:2]
}
`
	program, err := newgofun.Build(src)
	if err != nil {
		t.Fatalf("gofun Build: %v", err)
	}
	r1, err := program.Run("ValidSlice")
	if err != nil {
		t.Logf("gofun ValidSlice 错误: %v", err)
	} else {
		t.Logf("gofun ValidSlice() = %v", r1)
	}
	r2, err := program.Run("EmptySlice")
	if err != nil {
		t.Logf("gofun EmptySlice 错误: %v", err)
	} else {
		t.Logf("gofun EmptySlice() = %v", r2)
	}

	// Gig 对比
	gigProg, _ := gig.Build(src)
	gigR1, _ := gigProg.Run("ValidSlice")
	gigR2, _ := gigProg.Run("EmptySlice")
	t.Logf("Gig ValidSlice() = %v", gigR1)
	t.Logf("Gig EmptySlice() = %v", gigR2)
}
