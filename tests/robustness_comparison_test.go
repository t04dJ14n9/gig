package tests

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"git.woa.com/youngjin/gig"
	"git.woa.com/youngjin/gig/value"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
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
// 性能对比测试
// ----------------------------------------------------------------------------

func BenchmarkNative_Fibonacci20(b *testing.B) {
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fib(20)
	}
}

func BenchmarkGig_Fibonacci20(b *testing.B) {
	source := `
package main

func fib(n int) int {
	if n <= 1 { return n }
	return fib(n-1) + fib(n-2)
}

func Fib20() int {
	return fib(20)
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("Fib20")
	}
}

// ============================================================================
// 规则引擎性能测试
// ============================================================================
//
// 规则引擎的真实性能测试位于：
//   reference/rule_engine/sdk/benchmark_test.go
//
// 运行方式（需要配置内部包访问权限）：
//   cd reference/rule_engine
//   go test -tags=ruleengine -bench=. -benchmem ./sdk/
//
// 注意：规则引擎依赖 git.code.oa.com 内部包，需要在企业内网环境运行
//
// 规则引擎 SDK API 使用示例：
//
//   dsl, _ := sdk.NewRuleDSL([]byte(jsonStr))
//   dsl.AddGlobalVar(sdk.Var{Name: ".userInfo", Value: `{"vip": "true"}`})
//   result, _ := sdk.RunRule(ctx, *dsl)
//
// 规则引擎测试函数列表（运行于 reference/rule_engine/sdk/benchmark_test.go）：
//   - BenchmarkRuleEngine_SimpleCondition    简单条件判断
//   - BenchmarkRuleEngine_NestedConditions   嵌套条件判断
//   - BenchmarkRuleEngine_VariableAccess     变量访问
//   - BenchmarkRuleEngine_JsonParsing        JSON 解析
//   - BenchmarkRuleEngine_DSLParse           DSL 解析
//   - BenchmarkRuleEngine_TemplateGeneration 模板生成
//   - BenchmarkRuleEngine_MultipleRuleGroups 多规则组
//   - BenchmarkRuleEngine_DSLCopy            DSL 复制
//   - BenchmarkRuleEngine_FullPipeline       完整流程

// ============================================================================
// 内存分配测试
// ============================================================================

func BenchmarkNative_MemoryAlloc(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = make([]int, 100)
	}
}

func BenchmarkGig_MemoryAlloc(b *testing.B) {
	source := `
package main

func MakeSlice() []int {
	return make([]int, 100)
}
`
	prog, _ := gig.Build(source)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("MakeSlice")
	}
}

// ============================================================================
// CPU 负载测试
// ============================================================================

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
// 外部函数调用测试
// ============================================================================

func BenchmarkNative_ExternalCall(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strconv.Itoa(i)
		_ = strings.Contains(strconv.Itoa(i), "5")
	}
}

func BenchmarkGig_ExternalCall(b *testing.B) {
	source := `
package main

import (
	"strconv"
	"strings"
)

func ExternalCalls() int {
	count := 0
	for i := 0; i < 100; i++ {
		s := strconv.Itoa(i)
		if strings.Contains(s, "5") {
			count++
		}
	}
	return count
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("ExternalCalls")
	}
}

// ============================================================================
// Context 取消测试
// ============================================================================

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
		prog.RunWithContext(ctx, "LongRunning")
		cancel()
	}
}

// ============================================================================
// 综合性能对比报告
// ============================================================================

func TestPerformanceComparison(t *testing.T) {
	t.Log("========== 性能对比测试报告 ==========")
	t.Log("")
	t.Log("运行方式:")
	t.Log("  Gig/Yaegi/Gopher-Lua:")
	t.Log("    cd benchmarks && go test -bench=. -benchmem -count=3 -timeout=30m")
	t.Log("")
	t.Log("  gofun 基准测试:")
	t.Log("    go test -tags=gofun -bench=. -benchmem ./tests/gofun_benchmark_test.go")
	t.Log("")
	t.Log("  规则引擎（需要内网环境）:")
	t.Log("    cd reference/rule_engine && go test -tags=ruleengine -bench=. -benchmem ./sdk/")
	t.Log("")
	t.Log("预期结果:")
	t.Log("  - Native Go: 最快，作为基准")
	t.Log("  - Gig: 比 gofun 快 1.4-9.1 倍（大多数场景）")
	t.Log("  - Gig: 比 Yaegi 快 1.1-5.4 倍")
	t.Log("  - RuleEngine: 简单场景最优，不支持复杂逻辑")
	t.Log("")
	t.Log("文档中的性能数据来源:")
	t.Log("  - Gig/Yaegi/Gopher-Lua: benchmarks/bench_test.go")
	t.Log("  - gofun: tests/gofun_benchmark_test.go")
	t.Log("  - 规则引擎: reference/rule_engine/sdk/benchmark_test.go")
	t.Log("  - 健壮性测试: tests/robustness_comparison_test.go")
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
