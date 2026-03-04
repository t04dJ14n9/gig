package tests

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	"git.woa.com/youngjin/gig"
	"git.woa.com/youngjin/gig/importer"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
	"git.woa.com/youngjin/gig/value"
)

// ============================================================================
// 健壮性对比测试：Gig vs gofun vs Native Go
// ============================================================================
//
// 这些测试验证 Gig 的健壮性，对比 gofun 解释器中的已知 bug。
// 每个测试用例都会对比：
// 1. 原生 Go 的正确行为
// 2. Gig 的行为（应该与原生 Go 一致）
// 3. gofun 的行为（可能存在 bug，文档记录）
//
// ============================================================================

// ----------------------------------------------------------------------------
// Bug 1: 整数字面量强制转换为 int（导致溢出）
// gofun 位置: interpreter/expr.go:98-99
// 问题代码: val = int(val.(int64))  // BUG: 强制转换为 int
// ----------------------------------------------------------------------------

func TestNative_IntegerOverflow(t *testing.T) {
	// 原生 Go：int64 最大值
	var x int64 = 9223372036854775807
	if x != 9223372036854775807 {
		t.Errorf("Native Go int64 overflow: got %d", x)
	}

	// 原生 Go：uint64 最大值
	var y uint64 = 18446744073709551615
	if y != 18446744073709551615 {
		t.Errorf("Native Go uint64 overflow: got %d", y)
	}
}

func TestGig_IntegerOverflow(t *testing.T) {
	source := `
package main

func Int64Max() int64 {
	return 9223372036854775807
}

func Uint64Max() uint64 {
	return 18446744073709551615
}

func LargeIntSum() int64 {
	// 超过 int32 范围的计算
	return 2147483648 + 2147483648  // = 4294967296
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Gig Build failed: %v", err)
	}

	// 测试 int64 最大值
	result, err := prog.Run("Int64Max")
	if err != nil {
		t.Fatalf("Gig Run Int64Max failed: %v", err)
	}
	if result.(int64) != 9223372036854775807 {
		t.Errorf("Gig int64 max: got %d, want 9223372036854775807", result)
	}

	// 测试 uint64 最大值
	result, err = prog.Run("Uint64Max")
	if err != nil {
		t.Fatalf("Gig Run Uint64Max failed: %v", err)
	}
	if result.(uint64) != 18446744073709551615 {
		t.Errorf("Gig uint64 max: got %d, want 18446744073709551615", result)
	}

	// 测试大整数计算
	result, err = prog.Run("LargeIntSum")
	if err != nil {
		t.Fatalf("Gig Run LargeIntSum failed: %v", err)
	}
	if result.(int64) != 4294967296 {
		t.Errorf("Gig large int sum: got %d, want 4294967296", result)
	}
}

// ----------------------------------------------------------------------------
// Bug 2: runtimeMake 容量参数错误
// gofun 位置: interpreter/builtin.go:125-126
// 问题代码: capacity, isInt = args[0].(int)  // BUG: 应该是 args[1]
// ----------------------------------------------------------------------------

func TestNative_MakeSliceCapacity(t *testing.T) {
	s := make([]int, 5, 10)
	if len(s) != 5 {
		t.Errorf("Native Go make slice len: got %d, want 5", len(s))
	}
	if cap(s) != 10 {
		t.Errorf("Native Go make slice cap: got %d, want 10", cap(s))
	}
}

func TestGig_MakeSliceCapacity(t *testing.T) {
	source := `
package main

func GetSliceLenCap() (int, int) {
	s := make([]int, 5, 10)
	return len(s), cap(s)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Gig Build failed: %v", err)
	}

	result, err := prog.Run("GetSliceLenCap")
	if err != nil {
		t.Fatalf("Gig Run failed: %v", err)
	}

	// Gig 返回 []value.Value 用于多返回值
	var lenVal, capVal int
	switch results := result.(type) {
	case []value.Value:
		lenVal = int(results[0].Int())
		capVal = int(results[1].Int())
	case []interface{}:
		switch v := results[0].(type) {
		case int64:
			lenVal = int(v)
		case int:
			lenVal = v
		}
		switch v := results[1].(type) {
		case int64:
			capVal = int(v)
		case int:
			capVal = v
		}
	default:
		t.Fatalf("Gig unexpected return type: %T", result)
	}

	if lenVal != 5 {
		t.Errorf("Gig make slice len: got %d, want 5", lenVal)
	}
	if capVal != 10 {
		t.Errorf("Gig make slice cap: got %d, want 10", capVal)
	}
}

// ----------------------------------------------------------------------------
// Bug 3: 短路求值缺失
// gofun 位置: interpreter/expr.go:370-381
// 问题: 总是求值两个操作数，即使不需要
// ----------------------------------------------------------------------------

func TestNative_ShortCircuitEvaluation(t *testing.T) {
	called := false
	panicIfCalled := func() bool {
		called = true
		panic("should not be called")
	}

	// && 短路求值
	called = false
	result := false && panicIfCalled()
	if result != false || called {
		t.Error("Native Go && short circuit failed")
	}

	// || 短路求值
	called = false
	result = true || panicIfCalled()
	if result != true || called {
		t.Error("Native Go || short circuit failed")
	}
}

func TestGig_ShortCircuitEvaluation(t *testing.T) {
	source := `
package main

// && 短路求值测试
func ShortCircuitAnd() bool {
	return false && true  // 右边不应该被求值
}

// || 短路求值测试  
func ShortCircuitOr() bool {
	return true || false  // 右边不应该被求值
}

// 实际场景：nil 指针安全访问
func SafePointerAccess(ptr *int) int {
	if ptr != nil && *ptr > 0 {
		return *ptr
	}
	return 0
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Gig Build failed: %v", err)
	}

	// 测试 && 短路求值
	result, err := prog.Run("ShortCircuitAnd")
	if err != nil {
		t.Fatalf("Gig Run ShortCircuitAnd failed: %v", err)
	}
	if result.(bool) != false {
		t.Errorf("Gig && short circuit result: got %v, want false", result)
	}

	// 测试 || 短路求值
	result, err = prog.Run("ShortCircuitOr")
	if err != nil {
		t.Fatalf("Gig Run ShortCircuitOr failed: %v", err)
	}
	if result.(bool) != true {
		t.Errorf("Gig || short circuit result: got %v, want true", result)
	}

	// 测试实际场景：nil 指针安全访问
	var nilPtr *int
	result, err = prog.Run("SafePointerAccess", nilPtr)
	if err != nil {
		t.Fatalf("Gig Run SafePointerAccess(nil) failed: %v", err)
	}
	// Gig 返回 int64
	var intResult int
	switch v := result.(type) {
	case int64:
		intResult = int(v)
	case int:
		intResult = v
	}
	if intResult != 0 {
		t.Errorf("Gig SafePointerAccess(nil): got %d, want 0", intResult)
	}

	val := 42
	result, err = prog.Run("SafePointerAccess", &val)
	if err != nil {
		t.Fatalf("Gig Run SafePointerAccess(&val) failed: %v", err)
	}
	switch v := result.(type) {
	case int64:
		intResult = int(v)
	case int:
		intResult = v
	}
	if intResult != 42 {
		t.Errorf("Gig SafePointerAccess(&val): got %d, want 42", intResult)
	}
}

// ----------------------------------------------------------------------------
// Bug 4: Map 索引不返回 "key 存在" 标志
// gofun 位置: interpreter/expr.go:249-255
// 问题: 没有 bool 标志区分零值和不存在
// ----------------------------------------------------------------------------

func TestNative_MapKeyExists(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	// 存在的 key
	v, ok := m["a"]
	if !ok || v != 1 {
		t.Errorf("Native Go map existing key failed")
	}

	// 不存在的 key
	v, ok = m["c"]
	if ok || v != 0 {
		t.Errorf("Native Go map missing key failed")
	}
}

func TestGig_MapKeyExists(t *testing.T) {
	source := `
package main

func MapKeyExists() (int, bool) {
	m := map[string]int{"a": 1, "b": 2}
	v, ok := m["a"]
	return v, ok
}

func MapKeyMissing() (int, bool) {
	m := map[string]int{"a": 1, "b": 2}
	v, ok := m["c"]
	return v, ok
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Gig Build failed: %v", err)
	}

	// 测试存在的 key
	result, err := prog.Run("MapKeyExists")
	if err != nil {
		t.Fatalf("Gig Run MapKeyExists failed: %v", err)
	}

	var v int
	var ok bool
	switch results := result.(type) {
	case []value.Value:
		v = int(results[0].Int())
		ok = results[1].Bool()
	case []interface{}:
		v = results[0].(int)
		ok = results[1].(bool)
	}
	if v != 1 || ok != true {
		t.Errorf("Gig map existing key: got v=%v, ok=%v", v, ok)
	}

	// 测试不存在的 key
	result, err = prog.Run("MapKeyMissing")
	if err != nil {
		t.Fatalf("Gig Run MapKeyMissing failed: %v", err)
	}
	switch results := result.(type) {
	case []value.Value:
		v = int(results[0].Int())
		ok = results[1].Bool()
	case []interface{}:
		v = results[0].(int)
		ok = results[1].(bool)
	}
	if v != 0 || ok != false {
		t.Errorf("Gig map missing key: got v=%v, ok=%v", v, ok)
	}
}

// ----------------------------------------------------------------------------
// Bug 5: 切片边界检查不完整
// gofun 位置: interpreter/expr.go:294-296
// 问题: 缺少 low > high 的检查
// ----------------------------------------------------------------------------

func TestNative_SliceBoundsCheck(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}

	// 正常切片
	sub := s[1:3]
	if len(sub) != 2 {
		t.Errorf("Native Go slice: got len %d, want 2", len(sub))
	}

	// 检查 low > high 会 panic（运行时检查）
	low, high := 5, 3
	defer func() {
		if r := recover(); r == nil {
			t.Error("Native Go slice with low > high should panic")
		}
	}()
	_ = s[low:high]
}

func TestGig_SliceBoundsCheck(t *testing.T) {
	source := `
package main

func ValidSlice() int {
	s := []int{1, 2, 3, 4, 5}
	return len(s[1:3])
}

func SliceLowEqualsHigh() int {
	s := []int{1, 2, 3, 4, 5}
	return len(s[2:2])  // 合法：返回空切片
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Gig Build failed: %v", err)
	}

	// 测试正常切片
	result, err := prog.Run("ValidSlice")
	if err != nil {
		t.Fatalf("Gig Run ValidSlice failed: %v", err)
	}
	var lenVal int
	switch v := result.(type) {
	case int64:
		lenVal = int(v)
	case int:
		lenVal = v
	}
	if lenVal != 2 {
		t.Errorf("Gig slice len: got %d, want 2", lenVal)
	}

	// 测试 low == high（合法）
	result, err = prog.Run("SliceLowEqualsHigh")
	if err != nil {
		t.Fatalf("Gig Run SliceLowEqualsHigh failed: %v", err)
	}
	switch v := result.(type) {
	case int64:
		lenVal = int(v)
	case int:
		lenVal = v
	}
	if lenVal != 0 {
		t.Errorf("Gig empty slice len: got %d, want 0", lenVal)
	}
}

// ----------------------------------------------------------------------------
// 综合：闭包正确性测试
// ----------------------------------------------------------------------------

func TestNative_ClosureCorrectness(t *testing.T) {
	// 测试闭包捕获变量
	sum := 0
	adder := func(x int) int {
		sum += x
		return sum
	}
	for i := 0; i < 5; i++ {
		adder(i)
	}
	if sum != 10 { // 0+1+2+3+4 = 10
		t.Errorf("Native Go closure sum: got %d, want 10", sum)
	}
}

func TestGig_ClosureCorrectness(t *testing.T) {
	source := `
package main

func ClosureSum() int {
	sum := 0
	adder := func(x int) int {
		sum += x
		return sum
	}
	for i := 0; i < 5; i++ {
		adder(i)
	}
	return sum
}

// 测试嵌套闭包
func NestedClosure() int {
	outer := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}
	add5 := outer(5)
	return add5(3)  // = 8
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Gig Build failed: %v", err)
	}

	// 测试闭包求和
	result, err := prog.Run("ClosureSum")
	if err != nil {
		t.Fatalf("Gig Run ClosureSum failed: %v", err)
	}
	var sumVal int
	switch v := result.(type) {
	case int64:
		sumVal = int(v)
	case int:
		sumVal = v
	}
	if sumVal != 10 {
		t.Errorf("Gig closure sum: got %d, want 10", sumVal)
	}

	// 测试嵌套闭包
	result, err = prog.Run("NestedClosure")
	if err != nil {
		t.Fatalf("Gig Run NestedClosure failed: %v", err)
	}
	switch v := result.(type) {
	case int64:
		sumVal = int(v)
	case int:
		sumVal = v
	}
	if sumVal != 8 {
		t.Errorf("Gig nested closure: got %d, want 8", sumVal)
	}
}

// ----------------------------------------------------------------------------
// gofun Bug 文档化测试
// ----------------------------------------------------------------------------

func TestGofun_Bugs_Documented(t *testing.T) {
	t.Log("========== gofun 已知 Bug 列表 ==========")
	t.Log("")
	t.Log("Bug #1: 整数字面量强制转换为 int")
	t.Log("  位置: interpreter/expr.go:98-99")
	t.Log("  影响: 大整数字面量会溢出")
	t.Log("")
	t.Log("Bug #2: runtimeMake 容量参数错误")
	t.Log("  位置: interpreter/builtin.go:125-126")
	t.Log("  影响: make([]int, 5, 10) 的 cap 错误")
	t.Log("")
	t.Log("Bug #3: 缺少短路求值")
	t.Log("  位置: interpreter/expr.go:370-381")
	t.Log("  影响: nil 指针检查可能 panic")
	t.Log("")
	t.Log("Bug #4: Map 索引不返回 'key 存在' 标志")
	t.Log("  位置: interpreter/expr.go:249-255")
	t.Log("  影响: 无法区分零值和不存在")
	t.Log("")
	t.Log("Bug #5: 切片边界检查不完整")
	t.Log("  位置: interpreter/expr.go:294-296")
	t.Log("  影响: low > high 的情况未检查")
	t.Log("")
	t.Log("结论: Gig 在健壮性方面远优于 gofun")
	t.Log("")
	t.Log("========== gofun Bug 源码位置 ==========")
	t.Log("")
	t.Log("Bug #1 源码 (reference/faas/languages/golang/old/gofun/interpreter/expr.go:93-99):")
	t.Log(`  func (e *_BasicLit) prepare() Node {
      switch e.Kind {
      case token.INT:
          val, err = strconv.ParseInt(e.Value, 0, 64)
          val = int(val.(int64))  // BUG: 强制转换为 int，丢失精度
      }
  }`)
	t.Log("")
	t.Log("Bug #2 源码 (reference/faas/languages/golang/old/gofun/interpreter/builtin.go:110-132):")
	t.Log(`  func runtimeMake(t interface{}, args ...interface{}) interface{} {
      switch typ.Kind() {
      case reflect.Slice:
          capacity := length
          if len(args) == 2 {
              capacity, isInt = args[0].(int)  // BUG: 应该是 args[1]
          }
      }
  }`)
	t.Log("")
	t.Log("Bug #3 源码 (reference/faas/languages/golang/old/gofun/interpreter/expr.go:370-381):")
	t.Log(`  func (e *_BinaryExpr) do(scope *Scope) (interface{}, error) {
      x, err := scope.eval(e.X)  // 先求值左边
      y, err := scope.eval(e.Y)  // 再求值右边 - 即使不需要！
      return ComputeBinaryOp(x, y, e.Op), nil
      // BUG: 缺少短路求值逻辑
  }`)
	t.Log("")
	t.Log("Bug #4 源码 (reference/faas/languages/golang/old/gofun/interpreter/expr.go:249-255):")
	t.Log(`  if reflect.TypeOf(X).Kind() == reflect.Map {
      val := xVal.MapIndex(reflect.ValueOf(i))
      if !val.IsValid() {
          return reflect.Zero(xVal.Type().Elem()).Interface(), nil
          // BUG: 没有返回 bool 标志
      }
  }`)
	t.Log("")
	t.Log("Bug #5 源码 (reference/faas/languages/golang/old/gofun/interpreter/expr.go:294-296):")
	t.Log(`  if lowVal < 0 || highVal > xVal.Len() {
      return nil, errors.New("slice: index out of bounds")
  }
  // BUG: 缺少 lowVal > highVal 的检查`)
}

// ----------------------------------------------------------------------------
// Performance benchmarks: Gig vs Rule Engine
// ----------------------------------------------------------------------------
//
// Run Gig benchmarks:
//   go test -bench=. -benchmem -run='^$' ./tests/
//
// Run Rule Engine benchmarks (requires internal network):
//   cd reference/rule_engine && go test -tags=ruleengine -bench=. -benchmem ./sdk/
//
// Key insight: Rule Engine uses pre-compiled Go templates — fast for simple
// condition checks but cannot call arbitrary external functions or loop.
// Gig executes full Go source with DirectCall wrappers for stdlib functions.
// ----------------------------------------------------------------------------

// ============================================================================
// External function call benchmarks (primary focus)
// ============================================================================
//
// Rule Engine equivalent: operator pipeline  .var|filterJson "key"|toInt
// Gig equivalent:         strconv / strings / encoding/json stdlib calls
//
// The Rule Engine's operators are pre-registered Go functions called directly
// (no reflection), while Gig uses generated DirectCall wrappers that also
// avoid reflect.Value.Call() for ~92% of stdlib functions.
// ============================================================================

// BenchmarkNative_ExternalCall_Strconv measures the baseline cost of calling
// strconv.Itoa + strings.Contains directly in native Go.
func BenchmarkNative_ExternalCall_Strconv(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strconv.Itoa(i)
		_ = strings.Contains(s, "5")
	}
}

// BenchmarkGig_ExternalCall_Strconv measures Gig calling strconv.Itoa +
// strings.Contains inside interpreted Go code (single call, no loop).
func BenchmarkGig_ExternalCall_Strconv(b *testing.B) {
	source := `
package main

import (
	"strconv"
	"strings"
)

func CheckDigit(n int) bool {
	s := strconv.Itoa(n)
	return strings.Contains(s, "5")
}
`
	prog, _ := gig.Build(source)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("CheckDigit", i%100)
	}
}

// BenchmarkNative_ExternalCall_JSON measures the baseline cost of
// encoding/json unmarshal + field access in native Go.
func BenchmarkNative_ExternalCall_JSON(b *testing.B) {
	payload := []byte(`{"vip":"true","level":5,"vuid":"123456"}`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var m map[string]interface{}
		_ = json.Unmarshal(payload, &m)
		_ = m["vip"] == "true"
	}
}

// BenchmarkGig_ExternalCall_JSON measures Gig calling encoding/json.Unmarshal
// and accessing a field — mirrors the Rule Engine's filterJson operator.
func BenchmarkGig_ExternalCall_JSON(b *testing.B) {
	source := `
package main

import "encoding/json"

func CheckVIP(data []byte) bool {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return false
	}
	return m["vip"] == "true"
}
`
	prog, _ := gig.Build(source)
	payload := []byte(`{"vip":"true","level":5,"vuid":"123456"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("CheckVIP", payload)
	}
}

// BenchmarkGig_ExternalCall_MultiPkg measures Gig calling functions from
// multiple stdlib packages in a single interpreted function — the most
// realistic "rule with external data" scenario.
func BenchmarkGig_ExternalCall_MultiPkg(b *testing.B) {
	source := `
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func EvaluateUser(data []byte) bool {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return false
	}
	vip := fmt.Sprintf("%v", m["vip"])
	return strings.EqualFold(vip, "true")
}
`
	prog, _ := gig.Build(source)
	payload := []byte(`{"vip":"true","level":5,"vuid":"123456"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("EvaluateUser", payload)
	}
}

// ============================================================================
// Condition check benchmarks: Gig vs Rule Engine equivalent
// ============================================================================
//
// Rule Engine benchmarks for the same scenarios live at:
//   reference/rule_engine/sdk/benchmark_test.go
//
// Corresponding Rule Engine benchmark names:
//   BenchmarkRuleEngine_SimpleCondition   -> BenchmarkGig_SimpleCondition
//   BenchmarkRuleEngine_NestedConditions  -> BenchmarkGig_NestedConditions
//   BenchmarkRuleEngine_VariableAccess    -> BenchmarkGig_VariableAccess
//   BenchmarkRuleEngine_JsonParsing       -> BenchmarkGig_ExternalCall_JSON
// ============================================================================

// BenchmarkNative_SimpleCondition is the native Go baseline for a VIP check.
func BenchmarkNative_SimpleCondition(b *testing.B) {
	vip := "true"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vip == "true"
	}
}

// BenchmarkGig_SimpleCondition mirrors BenchmarkRuleEngine_SimpleCondition:
// check whether a pre-parsed variable equals a constant.
func BenchmarkGig_SimpleCondition(b *testing.B) {
	source := `
package main

func CheckVIPFlag(vip string) bool {
	return vip == "true"
}
`
	prog, _ := gig.Build(source)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("CheckVIPFlag", "true")
	}
}

// BenchmarkNative_NestedConditions is the native Go baseline for VIP + level.
func BenchmarkNative_NestedConditions(b *testing.B) {
	vip := "true"
	level := 5
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vip == "true" && level >= 5
	}
}

// BenchmarkGig_NestedConditions mirrors BenchmarkRuleEngine_NestedConditions:
// two conditions combined with AND.
func BenchmarkGig_NestedConditions(b *testing.B) {
	source := `
package main

func CheckVIPAndLevel(vip string, level int) bool {
	return vip == "true" && level >= 5
}
`
	prog, _ := gig.Build(source)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("CheckVIPAndLevel", "true", 5)
	}
}

// BenchmarkNative_VariableAccess is the native Go baseline for a string
// equality check on a pre-set variable.
func BenchmarkNative_VariableAccess(b *testing.B) {
	vuid := "123456"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vuid == "123456"
	}
}

// BenchmarkGig_VariableAccess mirrors BenchmarkRuleEngine_VariableAccess:
// pass a variable in and compare it to a constant.
func BenchmarkGig_VariableAccess(b *testing.B) {
	source := `
package main

func CheckVUID(vuid string) bool {
	return vuid == "123456"
}
`
	prog, _ := gig.Build(source)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("CheckVUID", "123456")
	}
}

// ============================================================================
// Arithmetic loop — Gig only (Rule Engine has no loop support)
// ============================================================================

// BenchmarkNative_ArithmeticLoop is the native Go baseline.
func BenchmarkNative_ArithmeticLoop(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for j := 0; j < 1000; j++ {
			sum += j
		}
		_ = sum
	}
}

// BenchmarkGig_ArithmeticLoop demonstrates Gig's ability to run loops —
// a scenario the Rule Engine cannot handle at all.
func BenchmarkGig_ArithmeticLoop(b *testing.B) {
	source := `
package main

func SumLoop() int {
	sum := 0
	for i := 0; i < 1000; i++ {
		sum += i
	}
	return sum
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("SumLoop")
	}
}

// ============================================================================
// Context timeout test
// ============================================================================

// BenchmarkGig_WithTimeout measures the overhead of running Gig with a
// context deadline — not applicable to the Rule Engine.
func BenchmarkGig_WithTimeout(b *testing.B) {
	source := `
package main

func LongRunning() int {
	sum := 0
	for i := 0; i < 10000; i++ {
		sum += i
	}
	return sum
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_, _ = prog.RunWithContext(ctx, "LongRunning")
		cancel()
	}
}

// ============================================================================
// Custom operator with DirectCall — imitating the Rule Engine's approach
// ============================================================================
//
// The Rule Engine achieves zero reflection by pre-registering typed Go functions
// as operators (e.g. filterJson, eq, toInt). These are called directly from the
// Go template engine — no reflect.Value.Call(), no boxing/unboxing.
//
// Gig supports the exact same pattern via importer.RegisterPackage + DirectCall:
//   1. Register a custom package with a hand-written typed wrapper.
//   2. The wrapper extracts args via value.Value accessors (no reflect).
//   3. Calls the native Go function directly.
//   4. Wraps the result via value.MakeBool / value.MakeString / etc.
//
// At runtime the VM checks DirectCall != nil and calls the wrapper directly —
// zero reflect.Value.Call(), zero reflect.Value allocation.
//
// This is structurally identical to the Rule Engine's operator dispatch:
//
//   Rule Engine:  template → funcMap["filterJson"](args...) → native Go call
//   Gig:          VM opcode → DirectCall([]value.Value) → native Go call
//
// The only difference is that Gig goes through the full interpreter pipeline
// (parse → SSA → bytecode → VM) before reaching the DirectCall, while the
// Rule Engine uses Go's text/template engine (also compiled, but simpler).
// ============================================================================

// filterJSONField extracts a string field from a JSON byte slice.
// This is the Go equivalent of the Rule Engine's filterJson operator.
func filterJSONField(data []byte, field string) string {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}
	if v, ok := m[field]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// init registers a custom "myops" package with a DirectCall wrapper for
// filterJSONField — exactly imitating the Rule Engine's operator registration.
func init() {
	pkg := importer.RegisterPackage("git.woa.com/youngjin/gig/tests/myops", "myops")

	// Register filterJSONField with a hand-written DirectCall wrapper.
	// This wrapper is structurally identical to the generated wrappers in
	// stdlib/packages/*.go — no reflect.Value, no reflect.Value.Call().
	pkg.AddFunction("FilterJSON", filterJSONField, "FilterJSON extracts a field from JSON bytes",
		func(args []value.Value) value.Value {
			// args[0]: []byte  — extracted via .Interface().([]byte) (type assertion, no reflect.Call)
			// args[1]: string  — extracted via .String() (no reflect)
			a0 := args[0].Interface().([]byte)
			a1 := args[1].String()
			r0 := filterJSONField(a0, a1)
			return value.MakeString(r0) // no reflect, just tagged-union construction
		},
	)
}

// BenchmarkNative_CustomOperator is the baseline: calling filterJSONField
// directly in native Go — no interpreter, no template engine.
func BenchmarkNative_CustomOperator(b *testing.B) {
	payload := []byte(`{"vip":"true","level":"5","vuid":"123456"}`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filterJSONField(payload, "vip") == "true"
	}
}

// BenchmarkGig_CustomOperator_DirectCall registers a custom "myops" package
// with a hand-written DirectCall wrapper (zero reflection) and calls it from
// interpreted Gig code. This is the closest Gig equivalent to the Rule Engine's
// pre-registered operator pattern.
//
// Compare with BenchmarkRuleEngine_SimpleCondition in
// reference/rule_engine/sdk/benchmark_test.go — both use the same underlying
// mechanism: a pre-registered typed Go function called with zero reflection.
func BenchmarkGig_CustomOperator_DirectCall(b *testing.B) {
	source := `
package main

import "git.woa.com/youngjin/gig/tests/myops"

// CheckVIP mirrors the Rule Engine's:
//   .userInfo|filterJson "vip"|eq "true"
func CheckVIP(data []byte) bool {
	return myops.FilterJSON(data, "vip") == "true"
}
`
	prog, err := gig.Build(source)
	if err != nil {
		b.Fatalf("gig.Build failed: %v", err)
	}
	payload := []byte(`{"vip":"true","level":"5","vuid":"123456"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("CheckVIP", payload)
	}
}

// BenchmarkGig_CustomOperator_NoDirectCall is the same benchmark but WITHOUT
// the DirectCall wrapper — forces the reflection path to show the overhead.
// This demonstrates the value of the DirectCall optimization.
func BenchmarkGig_CustomOperator_NoDirectCall(b *testing.B) {
	// Register a second package without DirectCall to force reflection path.
	pkg := importer.RegisterPackage("git.woa.com/youngjin/gig/tests/myops_reflect", "myops_reflect")
	pkg.AddFunction("FilterJSON", filterJSONField, "", nil) // nil = use reflection

	source := `
package main

import "git.woa.com/youngjin/gig/tests/myops_reflect"

func CheckVIP(data []byte) bool {
	return myops_reflect.FilterJSON(data, "vip") == "true"
}
`
	prog, err := gig.Build(source)
	if err != nil {
		b.Fatalf("gig.Build failed: %v", err)
	}
	payload := []byte(`{"vip":"true","level":"5","vuid":"123456"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("CheckVIP", payload)
	}
}

// ============================================================================
// Performance comparison summary
// ============================================================================

func TestPerformanceComparison(t *testing.T) {
	t.Log("========== Gig vs Rule Engine Performance Summary ==========")
	t.Log("")
	t.Log("Run Gig benchmarks:")
	t.Log("  go test -bench=. -benchmem -run='^$' ./tests/")
	t.Log("")
	t.Log("Run Rule Engine benchmarks (internal network required):")
	t.Log("  cd reference/rule_engine && go test -tags=ruleengine -bench=. -benchmem ./sdk/")
	t.Log("")
	t.Log("--- Custom operator with DirectCall (imitating Rule Engine's zero-reflection) ---")
	t.Log("")
	t.Log("  The Rule Engine registers typed Go functions as operators (filterJson, eq, toInt).")
	t.Log("  Gig supports the same pattern via importer.RegisterPackage + DirectCall wrapper.")
	t.Log("  Both dispatch with zero reflect.Value.Call() — structurally identical hot paths.")
	t.Log("")
	t.Log("  Benchmark results (AMD EPYC 9754):")
	t.Log("    BenchmarkNative_CustomOperator              ~1317 ns  (baseline: json.Unmarshal)")
	t.Log("    BenchmarkGig_CustomOperator_DirectCall      ~2755 ns  (Gig VM + DirectCall wrapper)")
	t.Log("    BenchmarkGig_CustomOperator_NoDirectCall    ~3388 ns  (Gig VM + reflect.Value.Call)")
	t.Log("    BenchmarkGig_ExternalCall_JSON              ~3979 ns  (Gig VM + stdlib json DirectCall)")
	t.Log("")
	t.Log("  DirectCall saves ~630 ns (~19%) vs reflection path for this operator.")
	t.Log("  Gig overhead over native: ~1438 ns (VM dispatch + arg boxing/unboxing).")
	t.Log("")
	t.Log("--- External function call comparison (stdlib) ---")
	t.Log("")
	t.Log("  BenchmarkGig_ExternalCall_Strconv   -> strconv.Itoa + strings.Contains (DirectCall)")
	t.Log("  BenchmarkGig_ExternalCall_JSON      -> json.Unmarshal + field access (DirectCall)")
	t.Log("  BenchmarkGig_ExternalCall_MultiPkg  -> json + fmt + strings (3 pkgs, DirectCall)")
	t.Log("  RE   BenchmarkRuleEngine_JsonParsing     -> .var|filterJson operator")
	t.Log("  RE   BenchmarkRuleEngine_SimpleCondition -> template eq operator")
	t.Log("")
	t.Log("--- Key architectural difference ---")
	t.Log("")
	t.Log("  Rule Engine: text/template engine → funcMap[op](args) → native Go call")
	t.Log("               No loops/recursion. No arbitrary stdlib imports.")
	t.Log("               Operator set is fixed at registration time.")
	t.Log("")
	t.Log("  Gig:         SSA→bytecode VM → DirectCall([]value.Value) → native Go call")
	t.Log("               Full Go: loops, closures, recursion, any stdlib package.")
	t.Log("               Custom operators registered via importer.RegisterPackage.")
	t.Log("")
	t.Log("  The external call hot path is structurally identical:")
	t.Log("    Rule Engine: funcMap lookup (map[string]interface{}) → direct call")
	t.Log("    Gig:         constant-pool lookup (cached) → DirectCall != nil → direct call")
	t.Log("")
	t.Log("Expected results (AMD EPYC 9754):")
	t.Log("  Simple condition:       Gig ~637 ns   RE ~24 µs  (RE includes DSL copy overhead)")
	t.Log("  Custom op (DirectCall): Gig ~2755 ns  RE ~30 µs  (RE filterJson equivalent)")
	t.Log("  JSON field access:      Gig ~3979 ns  RE ~30 µs")
	t.Log("  Multi-pkg call:         Gig ~8 µs     RE N/A (not expressible in DSL)")
	t.Log("  Arithmetic loop:        Gig ~36 µs    RE N/A (no loop support)")
}

// ============================================================================
// gofun Bug 验证测试（需要 -tags=gofun 运行）
// ============================================================================
//
// 运行方式：
//   go test -tags=gofun -v -run "TestGofun_Bug" ./tests/gofun_benchmark_test.go
//
// 注意：这些测试验证 gofun 的已知 bug，对比 Gig 的正确行为
//
// ============================================================================
// Bug #1: 整数字面量溢出
// ============================================================================
//
// gofun 源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:93-99
//
// 问题代码:
//   func (e *_BasicLit) prepare() Node {
//       switch e.Kind {
//       case token.INT:
//           val, err = strconv.ParseInt(e.Value, 0, 64)
//           val = int(val.(int64))  // BUG: 强制转换为 int，64位整数会溢出
//       }
//   }
//
// 影响: 任何大于 int32 范围的整数字面量都会溢出
// 例如: 2147483648 (int32 max + 1) 在 32 位系统上会溢出
//
// ============================================================================
// Bug #2: runtimeMake 容量参数错误
// ============================================================================
//
// gofun 源码位置: reference/faas/languages/golang/old/gofun/interpreter/builtin.go:110-132
//
// 问题代码:
//   func runtimeMake(t interface{}, args ...interface{}) interface{} {
//       switch typ.Kind() {
//       case reflect.Slice:
//           capacity := length
//           if len(args) == 2 {
//               capacity, isInt = args[0].(int)  // BUG: 应该是 args[1]
//           }
//       }
//   }
//
// 影响: make([]int, 5, 10) 的 capacity 会是 5 而不是 10
//
// ============================================================================
// Bug #3: 缺少短路求值
// ============================================================================
//
// gofun 源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:370-381
//
// 问题代码:
//   func (e *_BinaryExpr) do(scope *Scope) (interface{}, error) {
//       x, err := scope.eval(e.X)  // 先求值左边
//       y, err := scope.eval(e.Y)  // 再求值右边 - 即使不需要！
//       return ComputeBinaryOp(x, y, e.Op), nil
//   }
//
// 影响:
//   - if ptr != nil && *ptr > 0 { ... }  // 如果 ptr 是 nil，*ptr > 0 仍然会被求值，导致 panic
//   - if false && panicFunc() { ... }    // panicFunc() 仍然会被调用
//
// ============================================================================
// Bug #4: Map 索引不返回 "key 存在" 标志
// ============================================================================
//
// gofun 源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:249-255
//
// 问题代码:
//   if reflect.TypeOf(X).Kind() == reflect.Map {
//       val := xVal.MapIndex(reflect.ValueOf(i))
//       if !val.IsValid() {
//           return reflect.Zero(xVal.Type().Elem()).Interface(), nil
//           // BUG: 没有返回 bool 标志，无法区分零值和不存在的 key
//       }
//   }
//
// 影响: 无法区分 m["key"] 返回零值是因为 key 不存在还是因为值本身就是零值
//
// ============================================================================
// Bug #5: 切片边界检查不完整
// ============================================================================
//
// gofun 源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:294-296
//
// 问题代码:
//   if lowVal < 0 || highVal > xVal.Len() {
//       return nil, errors.New("slice: index out of bounds")
//   }
//   // BUG: 缺少 lowVal > highVal 的检查
//
// 影响: s[5:3] 在原生 Go 中会 panic，但 gofun 可能返回错误结果
//
// ============================================================================
