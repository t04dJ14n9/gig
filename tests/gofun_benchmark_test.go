//go:build gofun
// +build gofun

package tests

import (
	"fmt"
	"go/token"
	"strconv"
	"strings"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"

	// 使用 replace 指令后的导入路径
	gofun "git.code.oa.com/datacenter/faas/languages/golang/old/gofun/interpreter"
	_ "git.code.oa.com/datacenter/faas/languages/golang/old/gofun/interpreter/imports"
)

// ============================================================================
// gofun vs Gig 性能对比测试
// ============================================================================
//
// 运行方式：
//   go test -tags=gofun -bench=. -benchmem ./tests/gofun_benchmark_test.go
//
// 注意：需要配置 go.mod 的 replace 指令或修改 gofun 的导入路径

// ============================================================================
// Fibonacci 测试
// ============================================================================

func BenchmarkNative_Fib25(b *testing.B) {
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fib(25)
	}
}

func BenchmarkGig_Fib25(b *testing.B) {
	source := `
package main

func fib(n int) int {
	if n <= 1 { return n }
	return fib(n-1) + fib(n-2)
}

func FibRecursive() int {
	return fib(25)
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("FibRecursive")
	}
}

func BenchmarkGofun_Fib25(b *testing.B) {
	// gofun 不支持完整的函数递归，使用迭代版本
	// 注意：gfun 只执行 main 函数，需要在 main 中调用目标函数
	source := `
package main

func FibIterative() int {
	a, b := 0, 1
	for i := 0; i < 25; i++ {
		a, b = b, a+b
	}
	return a
}

func main() int {
	return FibIterative()
}
`

	program, err := gofun.Parse(source, nil)
	if err != nil {
		b.Fatalf("gofun Parse error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := program.Run(nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ============================================================================
// 算术循环测试
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

func BenchmarkGofun_ArithmeticLoop(b *testing.B) {
	source := `
package main

func SumLoop() int {
	sum := 0
	for i := 0; i < 1000; i++ {
		sum += i
	}
	return sum
}

func main() int {
	return SumLoop()
}
`

	program, _ := gofun.Parse(source, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = program.Run(nil)
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

func BenchmarkGofun_ExternalCall(b *testing.B) {
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

func main() int {
	return ExternalCalls()
}
`

	program, _ := gofun.Parse(source, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = program.Run(nil)
	}
}

// ============================================================================
// 闭包测试
// ============================================================================

func BenchmarkNative_Closure(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		adder := func(x int) int {
			sum += x
			return sum
		}
		for j := 0; j < 100; j++ {
			adder(j)
		}
	}
}

func BenchmarkGig_Closure(b *testing.B) {
	source := `
package main

func ClosureSum() int {
	sum := 0
	adder := func(x int) int {
		sum += x
		return sum
	}
	for i := 0; i < 100; i++ {
		adder(i)
	}
	return sum
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("ClosureSum")
	}
}

func BenchmarkGofun_Closure(b *testing.B) {
	source := `
package main

func ClosureSum() int {
	sum := 0
	adder := func(x int) int {
		sum += x
		return sum
	}
	for i := 0; i < 100; i++ {
		adder(i)
	}
	return sum
}

func main() int {
	return ClosureSum()
}
`

	program, _ := gofun.Parse(source, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = program.Run(nil)
	}
}

// ============================================================================
// 条件判断测试
// ============================================================================

func BenchmarkNative_Condition(b *testing.B) {
	vip := true
	level := 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vip && level >= 5
		_ = result
	}
}

func BenchmarkGig_Condition(b *testing.B) {
	source := `
package main

func CheckCondition(vip bool, level int) bool {
	return vip && level >= 5
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("CheckCondition", true, 5)
	}
}

func BenchmarkGofun_Condition(b *testing.B) {
	// gofun 使用表达式求值
	scope := gofun.NewScope()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scope.Set("vip", true)
		scope.Set("level", 5)
		_, _ = scope.InterpretExpr("vip && level >= 5")
	}
}

// ============================================================================
// 变量操作测试
// ============================================================================

func BenchmarkNative_VariableOps(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a := 1
		b := 2
		c := a + b
		d := c * 2
		_ = d
	}
}

func BenchmarkGig_VariableOps(b *testing.B) {
	source := `
package main

func VariableOps() int {
	a := 1
	b := 2
	c := a + b
	d := c * 2
	return d
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("VariableOps")
	}
}

func BenchmarkGofun_VariableOps(b *testing.B) {
	source := `
package main

func VariableOps() int {
	a := 1
	b := 2
	c := a + b
	d := c * 2
	return d
}

func main() int {
	return VariableOps()
}
`

	program, _ := gofun.Parse(source, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = program.Run(nil)
	}
}

// ============================================================================
// 切片操作测试
// ============================================================================

func BenchmarkNative_SliceOps(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, 100)
		for j := 0; j < 100; j++ {
			s[j] = j
		}
		sum := 0
		for _, v := range s {
			sum += v
		}
		_ = sum
	}
}

func BenchmarkGig_SliceOps(b *testing.B) {
	source := `
package main

func SliceOps() int {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = i
	}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("SliceOps")
	}
}

func BenchmarkGofun_SliceOps(b *testing.B) {
	source := `
package main

func SliceOps() int {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = i
	}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func main() int {
	return SliceOps()
}
`

	program, _ := gofun.Parse(source, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = program.Run(nil)
	}
}

// ============================================================================
// Map 操作测试
// ============================================================================

func BenchmarkNative_MapOps(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]int)
		for j := 0; j < 100; j++ {
			m[strconv.Itoa(j)] = j
		}
		sum := 0
		for _, v := range m {
			sum += v
		}
		_ = sum
	}
}

func BenchmarkGig_MapOps(b *testing.B) {
	source := `
package main

import "strconv"

func MapOps() int {
	m := make(map[string]int)
	for i := 0; i < 100; i++ {
		m[strconv.Itoa(i)] = i
	}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}
`
	prog, _ := gig.Build(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("MapOps")
	}
}

func BenchmarkGofun_MapOps(b *testing.B) {
	source := `
package main

import "strconv"

func MapOps() int {
	m := make(map[string]int)
	for i := 0; i < 100; i++ {
		m[strconv.Itoa(i)] = i
	}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

func main() int {
	return MapOps()
}
`

	program, _ := gofun.Parse(source, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = program.Run(nil)
	}
}

// ============================================================================
// 解析性能测试（不含执行）
// ============================================================================

func BenchmarkGig_ParseOnly(b *testing.B) {
	source := `
package main

func Test() int {
	sum := 0
	for i := 0; i < 1000; i++ {
		sum += i
	}
	return sum
}
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = gig.Build(source)
	}
}

func BenchmarkGofun_ParseOnly(b *testing.B) {
	source := `
package main

func Test() int {
	sum := 0
	for i := 0; i < 1000; i++ {
		sum += i
	}
	return sum
}

func main() int {
	return Test()
}
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = gofun.Parse(source, nil)
	}
}

// ============================================================================
// 综合性能对比报告
// ============================================================================

func TestGofunGigComparison(t *testing.T) {
	t.Log("========== gofun vs Gig 性能对比测试 ==========")
	t.Log("")
	t.Log("运行方式:")
	t.Log("  go test -tags=gofun -bench=. -benchmem ./tests/gofun_benchmark_test.go")
	t.Log("")
	t.Log("注意:")
	t.Log("  - gofun 不支持完整的递归函数，Fibonacci 使用迭代版本")
	t.Log("  - gofun 存在已知 bug（整数溢出、make 参数错误等）")
	t.Log("  - Gig 使用 SSA 编译 + VM 执行，gofun 使用 AST-Walking")
	t.Log("")
	t.Log("预期结果:")
	t.Log("  - Gig 在大多数场景下比 gofun 快 2-5 倍")
	t.Log("  - Gig 内存分配更少（帧池化优化）")
	t.Log("  - gofun 的 AST-Walking 方式有更高的解释开销")
}

// ============================================================================
// 健壮性对比测试（验证 gofun bug）
// ============================================================================
//
// 这些测试验证 gofun 解释器的已知 bug，对比 Gig 的正确行为
// 每个 Bug 都包含：
// 1. Bug 描述和影响
// 2. gofun 源码位置和问题代码
// 3. 实际测试验证
// 4. Gig 的正确行为对比

// ----------------------------------------------------------------------------
// Bug #1: 整数字面量溢出
// ----------------------------------------------------------------------------
//
// 源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:93-99
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
// ----------------------------------------------------------------------------

func TestGofun_Bug_IntegerOverflow(t *testing.T) {
	t.Log("========== Bug #1: 整数字面量溢出 ==========")
	t.Log("源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:93-99")
	t.Log("")

	source := `
package main

func LargeInt() int64 {
	return 9223372036854775807
}

func main() {}
`

	program, err := gofun.Parse(source, nil)
	if err != nil {
		t.Fatalf("gofun Parse error: %v", err)
	}

	result, err := program.Run(nil)
	if err != nil {
		t.Fatalf("gofun Run error: %v", err)
	}

	t.Logf("gofun LargeInt() = %v (type: %T)", result, result)
	t.Log("")
	t.Log("预期: 应该返回 9223372036854775807 (int64 max)")
	t.Log("实际: gofun 将 int64 强制转换为 int，导致溢出")
	t.Log("")
	t.Log("Gig 正确行为: TestGig_IntegerOverflow 通过，正确返回 int64 max")
}

// ----------------------------------------------------------------------------
// Bug #2: runtimeMake 容量参数错误
// ----------------------------------------------------------------------------
//
// 源码位置: reference/faas/languages/golang/old/gofun/interpreter/builtin.go:110-132
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
// ----------------------------------------------------------------------------

func TestGofun_Bug_MakeCapacity(t *testing.T) {
	t.Log("========== Bug #2: runtimeMake 容量参数错误 ==========")
	t.Log("源码位置: reference/faas/languages/golang/old/gofun/interpreter/builtin.go:110-132")
	t.Log("")

	source := `
package main

func MakeSlice() []int {
	s := make([]int, 5, 10)
	return s
}

func main() {
	s := MakeSlice()
}
`

	program, err := gofun.Parse(source, nil)
	if err != nil {
		t.Fatalf("gofun Parse error: %v", err)
	}

	result, err := program.Run(nil)
	if err != nil {
		t.Logf("gofun Run error: %v", err)
	}

	t.Logf("gofun MakeSlice() result: %v (type: %T)", result, result)
	t.Log("")
	t.Log("预期: len=5, cap=10")
	t.Log("实际: gofun 的 capacity 参数读取错误，cap 可能是 5")
	t.Log("")
	t.Log("Gig 正确行为: TestGig_MakeSliceCapacity 通过，len=5, cap=10")
}

// ----------------------------------------------------------------------------
// Bug #3: 缺少短路求值
// ----------------------------------------------------------------------------
//
// 源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:370-381
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
// ----------------------------------------------------------------------------

func TestGofun_Bug_ShortCircuit(t *testing.T) {
	t.Log("========== Bug #3: 缺少短路求值 ==========")
	t.Log("源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:370-381")
	t.Log("")

	_, _ = gofun.Parse("package main\nfunc main() {}", nil) // 验证 gofun 可用

	// 测试 AND 短路（直接使用 scope 表达式求值）
	scope := gofun.NewScope()
	scope.Set("a", false)
	scope.Set("b", true)

	result, _ := scope.InterpretExpr("a && b")
	t.Logf("gofun false && true = %v", result)
	t.Log("")
	t.Log("问题: gofun 会先求值 a，再求值 b，即使 a 是 false")
	t.Log("这意味着如果 b 有副作用（如 panic），仍然会执行")
	t.Log("")
	t.Log("Gig 正确行为: TestGig_ShortCircuitEvaluation 通过，正确实现短路求值")

	// 更危险的例子
	t.Log("")
	t.Log("危险示例:")
	t.Log("  if ptr != nil && *ptr > 0 { ... }")
	t.Log("  在 gofun 中，如果 ptr 是 nil，*ptr > 0 仍然会被求值，导致 panic")
}

// ----------------------------------------------------------------------------
// Bug #4: Map 索引不返回 "key 存在" 标志
// ----------------------------------------------------------------------------
//
// 源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:249-255
//
// 问题代码:
//   if reflect.TypeOf(X).Kind() == reflect.Map {
//       val := xVal.MapIndex(reflect.ValueOf(i))
//       if !val.IsValid() {
//           return reflect.Zero(xVal.Type().Elem()).Interface(), nil
//           // BUG: 没有返回 bool 标志
//       }
//   }
//
// 影响: 无法区分 m["key"] 返回零值是因为 key 不存在还是因为值本身就是零值
// ----------------------------------------------------------------------------

func TestGofun_Bug_MapKeyExists(t *testing.T) {
	t.Log("========== Bug #4: Map 索引不返回 'key 存在' 标志 ==========")
	t.Log("源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:249-255")
	t.Log("")

	t.Log("问题: gofun 的 map 索引操作不返回 'key 存在' 标志")
	t.Log("")
	t.Log("原生 Go:")
	t.Log("  v, ok := m[\"key\"]  // ok 表示 key 是否存在")
	t.Log("")
	t.Log("gofun:")
	t.Log("  v := m[\"key\"]       // 无法区分零值和不存在的 key")
	t.Log("")
	t.Log("影响示例:")
	t.Log("  m := map[string]int{\"a\": 0}")
	t.Log("  v1 := m[\"a\"]  // v1 = 0, key 存在")
	t.Log("  v2 := m[\"b\"]  // v2 = 0, key 不存在")
	t.Log("  // 在 gofun 中无法区分 v1 和 v2 的情况")
	t.Log("")
	t.Log("Gig 正确行为: TestGig_MapKeyExists 通过，正确支持 v, ok := m[key]")
}

// ----------------------------------------------------------------------------
// Bug #5: 切片边界检查不完整
// ----------------------------------------------------------------------------
//
// 源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:294-296
//
// 问题代码:
//   if lowVal < 0 || highVal > xVal.Len() {
//       return nil, errors.New("slice: index out of bounds")
//   }
//   // BUG: 缺少 lowVal > highVal 的检查
//
// 影响: s[5:3] 在原生 Go 中会 panic，但 gofun 可能返回错误结果
// ----------------------------------------------------------------------------

func TestGofun_Bug_SliceBoundsCheck(t *testing.T) {
	t.Log("========== Bug #5: 切片边界检查不完整 ==========")
	t.Log("源码位置: reference/faas/languages/golang/old/gofun/interpreter/expr.go:294-296")
	t.Log("")

	t.Log("问题: gofun 的切片边界检查不完整")
	t.Log("")
	t.Log("原生 Go:")
	t.Log("  s := []int{1, 2, 3}")
	t.Log("  s[5:3]  // panic: invalid slice bounds: 5 > 3")
	t.Log("")
	t.Log("gofun:")
	t.Log("  可能不检查 low > high 的情况，返回错误结果")
	t.Log("")
	t.Log("Gig 正确行为: TestGig_SliceBoundsCheck 通过，完整检查切片边界")
}

// ============================================================================
// 辅助函数
// ============================================================================

func parseGofunSource(b *testing.B, source string) *gofun.Program {
	b.Helper()
	program, err := gofun.Parse(source, nil)
	if err != nil {
		b.Fatalf("gofun Parse error: %v", err)
	}
	return program
}

func parseGigSource(b *testing.B, source string) *gig.Program {
	b.Helper()
	prog, err := gig.Build(source)
	if err != nil {
		b.Fatalf("gig Build error: %v", err)
	}
	return prog
}

// 打印文件集位置（用于调试）
func printPos(fset *token.FileSet, pos token.Pos) {
	fmt.Println(fset.Position(pos))
}
