package gofun

import (
	"strconv"
	"testing"

	newgofun "git.code.oa.com/datacenter/onefun/gofun"
	_ "git.code.oa.com/datacenter/onefun/gofun/packages"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

// ============================================================================
// 新 gofun (onefun/gofun) vs Gig 性能对比测试
// ============================================================================
//
// 运行: cd reference/gofun_tests && go test -bench=. -benchmem
//
// 新 gofun API:
//   gofun.Build(src) -> *Program
//   program.Run(funcName, args...) -> (interface{}, error)

// --- Fibonacci ---

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
		prog.Run("FibRecursive")
	}
}

func BenchmarkGofun_Fib25(b *testing.B) {
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
	program, err := newgofun.Build(source)
	if err != nil {
		b.Fatalf("gofun Build: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		program.Run("FibRecursive")
	}
}

// --- 算术循环 ---

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
		prog.Run("SumLoop")
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
`
	program, _ := newgofun.Build(source)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		program.Run("SumLoop")
	}
}

// --- 外部函数调用 ---

func BenchmarkNative_ExternalCall(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strconv.Itoa(i)
	}
}

func BenchmarkGig_ExternalCall(b *testing.B) {
	source := `
package main

import "strconv"

func ExternalCalls() int {
	count := 0
	for i := 0; i < 100; i++ {
		s := strconv.Itoa(i)
		if len(s) > 1 {
			count++
		}
	}
	return count
}
`
	prog, _ := gig.Build(source)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prog.Run("ExternalCalls")
	}
}

func BenchmarkGofun_ExternalCall(b *testing.B) {
	source := `
package main

import "strconv"

func ExternalCalls() int {
	count := 0
	for i := 0; i < 100; i++ {
		s := strconv.Itoa(i)
		if len(s) > 1 {
			count++
		}
	}
	return count
}
`
	program, _ := newgofun.Build(source)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		program.Run("ExternalCalls")
	}
}

// --- 闭包 ---

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
		prog.Run("ClosureSum")
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
`
	program, _ := newgofun.Build(source)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		program.Run("ClosureSum")
	}
}

// --- 变量操作 ---

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
		prog.Run("VariableOps")
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
`
	program, _ := newgofun.Build(source)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		program.Run("VariableOps")
	}
}

// --- 切片操作 ---

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
		prog.Run("SliceOps")
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
`
	program, _ := newgofun.Build(source)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		program.Run("SliceOps")
	}
}

// --- 编译性能 ---

func BenchmarkGig_BuildOnly(b *testing.B) {
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
		gig.Build(source)
	}
}

func BenchmarkGofun_BuildOnly(b *testing.B) {
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
		newgofun.Build(source)
	}
}
