package benchmarks

import (
	"context"
	"embed"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	lua "github.com/yuin/gopher-lua"

	"gig"
	_ "gig/stdlib/packages"
)

// ============================================================================
// Embedded source files (Go source for gig & yaegi)
// ============================================================================

//go:embed testdata/fib.go
var goFibSrc string

//go:embed testdata/arith.go
var goArithSrc string

//go:embed testdata/bubblesort.go
var goBubbleSortSrc string

//go:embed testdata/sieve.go
var goSieveSrc string

//go:embed testdata/closure.go
var goClosureSrc string

// ============================================================================
// Embedded Lua source files
// ============================================================================

//go:embed testdata/lua_fib.lua
var luaFibSrc string

//go:embed testdata/lua_arith.lua
var luaArithSrc string

//go:embed testdata/lua_bubblesort.lua
var luaBubbleSortSrc string

//go:embed testdata/lua_sieve.lua
var luaSieveSrc string

//go:embed testdata/lua_closure.lua
var luaClosureSrc string

// External call benchmark sources
//
//go:embed testdata/extcall_directcall.go
var goExtCallDirectCallSrc string

//go:embed testdata/extcall_reflect.go
var goExtCallReflectSrc string

//go:embed testdata/extcall_method.go
var goExtCallMethodSrc string

//go:embed testdata/extcall_mixed.go
var goExtCallMixedSrc string

// Keep embed import used (for go:embed directives above)
var _ embed.FS

// ============================================================================
// Native Go implementations
// ============================================================================

func nativeFib(n int) int {
	if n <= 1 {
		return n
	}
	return nativeFib(n-1) + nativeFib(n-2)
}

func nativeArithmeticSum() int {
	sum := 0
	for i := 1; i <= 1000; i++ {
		sum += i
	}
	return sum
}

func nativeBubbleSort() int {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = 100 - i
	}
	n := len(s)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if s[j] > s[j+1] {
				s[j], s[j+1] = s[j+1], s[j]
			}
		}
	}
	return s[0] + s[99]
}

func nativeSieve() int {
	n := 1000
	sieve := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		sieve[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if sieve[i] {
			for j := i * i; j <= n; j += i {
				sieve[j] = false
			}
		}
	}
	count := 0
	for i := 2; i <= n; i++ {
		if sieve[i] {
			count++
		}
	}
	return count
}

func nativeClosureCalls() int {
	sum := 0
	adder := func(x int) int {
		sum += x
		return sum
	}
	for i := 0; i < 1000; i++ {
		adder(i)
	}
	return sum
}

// ============================================================================
// Native Go: External Call Benchmarks
// ============================================================================

func nativeExtCallDirectCall() int {
	count := 0
	for i := 0; i < 1000; i++ {
		s := strconv.Itoa(i)
		if strings.Contains(s, "5") {
			count++
		}
		_ = strings.ToUpper(s)
		_ = math.Sqrt(float64(i))
	}
	return count
}

func nativeExtCallReflect() int {
	r := strings.NewReplacer("a", "b", "c", "d")
	sum := 0
	for i := 0; i < 1000; i++ {
		s := strconv.Itoa(i)
		result := r.Replace(s)
		sum += len(result)
	}
	return sum
}

func nativeExtCallMethod() int {
	sum := 0
	for i := 0; i < 1000; i++ {
		s := strconv.Itoa(i)
		r := strings.NewReader(s)
		sum += r.Len()
	}
	return sum
}

func nativeExtCallMixed() int {
	sum := 0
	for i := 0; i < 500; i++ {
		s := strconv.Itoa(i)
		if strings.Contains(s, "3") {
			sum += len(strings.ToUpper(s))
		}
		r := strings.NewReader(s)
		sum += r.Len()
	}
	return sum
}

// ============================================================================
// Gig Benchmarks
// ============================================================================

func benchGig(b *testing.B, src, funcName string) {
	b.Helper()
	prog, err := gig.Build(src)
	if err != nil {
		b.Fatalf("gig Build error: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := prog.Run(funcName)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGig_Fib25(b *testing.B)        { benchGig(b, goFibSrc, "FibRecursive") }
func BenchmarkGig_ArithSum(b *testing.B)     { benchGig(b, goArithSrc, "ArithmeticSum") }
func BenchmarkGig_BubbleSort(b *testing.B)   { benchGig(b, goBubbleSortSrc, "BubbleSort") }
func BenchmarkGig_Sieve(b *testing.B)        { benchGig(b, goSieveSrc, "Sieve") }
func BenchmarkGig_ClosureCalls(b *testing.B) { benchGig(b, goClosureSrc, "ClosureCalls") }

// Gig: External call benchmarks
func BenchmarkGig_ExtCallDirectCall(b *testing.B) { benchGig(b, goExtCallDirectCallSrc, "ExtCallDirectCall") }
func BenchmarkGig_ExtCallReflect(b *testing.B)    { benchGig(b, goExtCallReflectSrc, "ExtCallReflect") }
func BenchmarkGig_ExtCallMethod(b *testing.B)     { benchGig(b, goExtCallMethodSrc, "ExtCallMethod") }
func BenchmarkGig_ExtCallMixed(b *testing.B)      { benchGig(b, goExtCallMixedSrc, "ExtCallMixed") }

// ============================================================================
// Yaegi Benchmarks
// ============================================================================

func benchYaegi(b *testing.B, src, funcName string) {
	b.Helper()
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	_, err := i.Eval(src)
	if err != nil {
		b.Fatalf("yaegi Eval error: %v", err)
	}
	fn, err := i.Eval(funcName)
	if err != nil {
		b.Fatalf("yaegi func lookup error: %v", err)
	}
	callable := fn.Interface().(func() int)
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		_ = callable()
	}
}

func BenchmarkYaegi_Fib25(b *testing.B)        { benchYaegi(b, goFibSrc, "FibRecursive") }
func BenchmarkYaegi_ArithSum(b *testing.B)     { benchYaegi(b, goArithSrc, "ArithmeticSum") }
func BenchmarkYaegi_BubbleSort(b *testing.B)   { benchYaegi(b, goBubbleSortSrc, "BubbleSort") }
func BenchmarkYaegi_Sieve(b *testing.B)        { benchYaegi(b, goSieveSrc, "Sieve") }
func BenchmarkYaegi_ClosureCalls(b *testing.B) { benchYaegi(b, goClosureSrc, "ClosureCalls") }

// Yaegi: External call benchmarks
func BenchmarkYaegi_ExtCallDirectCall(b *testing.B) { benchYaegi(b, goExtCallDirectCallSrc, "ExtCallDirectCall") }
func BenchmarkYaegi_ExtCallReflect(b *testing.B)    { benchYaegi(b, goExtCallReflectSrc, "ExtCallReflect") }
func BenchmarkYaegi_ExtCallMethod(b *testing.B)     { benchYaegi(b, goExtCallMethodSrc, "ExtCallMethod") }
func BenchmarkYaegi_ExtCallMixed(b *testing.B)      { benchYaegi(b, goExtCallMixedSrc, "ExtCallMixed") }

// ============================================================================
// GopherLua Benchmarks
// ============================================================================

func benchLuaReuse(b *testing.B, src string) {
	b.Helper()
	L := lua.NewState()
	defer L.Close()
	fn, err := L.LoadString(src)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L.Push(fn)
		if err := L.PCall(0, lua.MultRet, nil); err != nil {
			b.Fatal(err)
		}
		L.Pop(L.GetTop())
	}
}

func BenchmarkLua_Fib25(b *testing.B)        { benchLuaReuse(b, luaFibSrc) }
func BenchmarkLua_ArithSum(b *testing.B)     { benchLuaReuse(b, luaArithSrc) }
func BenchmarkLua_BubbleSort(b *testing.B)   { benchLuaReuse(b, luaBubbleSortSrc) }
func BenchmarkLua_Sieve(b *testing.B)        { benchLuaReuse(b, luaSieveSrc) }
func BenchmarkLua_ClosureCalls(b *testing.B) { benchLuaReuse(b, luaClosureSrc) }

// ============================================================================
// Expr Benchmarks (expression evaluation - different domain)
// ============================================================================

func BenchmarkExpr_ArithExpr(b *testing.B) {
	env := map[string]any{
		"values": makeRange(1000),
	}
	prog, err := expr.Compile(`len(filter(values, {# > 500}))`, expr.Env(env))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = expr.Run(prog, env)
	}
}

func BenchmarkExpr_MapAccess(b *testing.B) {
	env := map[string]any{
		"user": map[string]any{
			"name":   "Alice",
			"age":    30,
			"scores": []int{95, 87, 92, 88, 96},
		},
	}
	prog, err := expr.Compile(`user.name == "Alice" && user.age > 20 && len(user.scores) > 3`, expr.Env(env))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = expr.Run(prog, env)
	}
}

func BenchmarkExpr_Conditional(b *testing.B) {
	env := map[string]any{
		"x": 42,
		"y": 17,
	}
	prog, err := expr.Compile(`x > y ? x * 2 + y : y * 2 + x`, expr.Env(env))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = expr.Run(prog, env)
	}
}

// ============================================================================
// Native Go Benchmarks
// ============================================================================

func BenchmarkNative_Fib25(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nativeFib(25)
	}
}

func BenchmarkNative_ArithSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nativeArithmeticSum()
	}
}

func BenchmarkNative_BubbleSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nativeBubbleSort()
	}
}

func BenchmarkNative_Sieve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nativeSieve()
	}
}

func BenchmarkNative_ClosureCalls(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nativeClosureCalls()
	}
}

// Native: External call benchmarks
func BenchmarkNative_ExtCallDirectCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nativeExtCallDirectCall()
	}
}

func BenchmarkNative_ExtCallReflect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nativeExtCallReflect()
	}
}

func BenchmarkNative_ExtCallMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nativeExtCallMethod()
	}
}

func BenchmarkNative_ExtCallMixed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = nativeExtCallMixed()
	}
}

// ============================================================================
// Helper
// ============================================================================

func makeRange(n int) []any {
	s := make([]any, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}
	return s
}

// ============================================================================
// Summary Test: Run all benchmarks and print comparison table
// ============================================================================

func TestCrossInterpreterComparison(t *testing.T) {
	var cpuModel string
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "model name") {
				cpuModel = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
				break
			}
		}
	}
	if cpuModel == "" {
		cpuModel = "Unknown"
	}

	t.Log("")
	t.Log("╔═══════════════════════════════════════════════════════════════════════════════╗")
	t.Log("║            Cross-Interpreter Benchmark Comparison                            ║")
	t.Logf("║  CPU: %-70s ║", cpuModel)
	t.Logf("║  Cores: %d | GOOS: %s | GOARCH: %s                                    ║",
		runtime.NumCPU(), runtime.GOOS, runtime.GOARCH)
	t.Log("╚═══════════════════════════════════════════════════════════════════════════════╝")
	t.Log("")
	t.Log("Run benchmarks yourself:")
	t.Log("  cd benchmarks && go test -bench=. -benchmem -count=3 -timeout=30m")
	t.Log("")
	t.Log("NOTE: This test only prints instructions. Run the benchmarks with -bench flag.")
	t.Log("")

	t.Run("VerifyGig", func(t *testing.T) {
		prog, err := gig.Build(goFibSrc)
		if err != nil {
			t.Fatal(err)
		}
		result, err := prog.RunWithContext(context.Background(), "FibRecursive")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("gig fib(25) = %v (type %T)", result, result)
		got := fmt.Sprintf("%v", result)
		if got != "75025" {
			t.Errorf("gig fib(25) = %v, want 75025", result)
		}
	})

	t.Run("VerifyYaegi", func(t *testing.T) {
		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)
		_, err := i.Eval(goFibSrc)
		if err != nil {
			t.Fatal(err)
		}
		fn, err := i.Eval("FibRecursive")
		if err != nil {
			t.Fatal(err)
		}
		result := fn.Interface().(func() int)()
		if result != 75025 {
			t.Errorf("yaegi fib(25) = %d, want 75025", result)
		}
	})

	t.Run("VerifyLua", func(t *testing.T) {
		L := lua.NewState()
		defer L.Close()
		if err := L.DoString(luaFibSrc); err != nil {
			t.Fatal(err)
		}
		result := L.Get(-1)
		t.Logf("lua fib(25) = %s", result.String())
		if result.String() != "75025" {
			t.Errorf("lua fib(25) = %s, want 75025", result.String())
		}
	})

	t.Run("VerifyNative", func(t *testing.T) {
		result := nativeFib(25)
		if result != 75025 {
			t.Errorf("native fib(25) = %d, want 75025", result)
		}
	})

	_ = fmt.Sprintf // suppress import
}
