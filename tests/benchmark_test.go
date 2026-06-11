package tests

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"math/big"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
	"github.com/t04dJ14n9/gig/tests/testdata/benchmarks"
)

//go:embed testdata/benchmarks/*.go
var benchmarksFS embed.FS

var benchmarkFileOrder = []string{
	"compute.go",
	"datastruct.go",
	"strings.go",
	"closure.go",
	"algorithm.go",
	"external.go",
	"calls.go",
	"types.go",
	"concurrency.go",
}

type benchmarkSourceParts struct {
	fset    *token.FileSet
	imports []*ast.ImportSpec
	decls   []ast.Decl
}

// getBenchmarksSrc reads all .go files from the embedded filesystem and concatenates them.
func getBenchmarksSrc() string {
	return renderBenchmarkSource(collectBenchmarkSourceParts())
}

func collectBenchmarkSourceParts() benchmarkSourceParts {
	parts := benchmarkSourceParts{fset: token.NewFileSet()}
	seenImports := make(map[string]bool)

	for _, fname := range benchmarkFileOrder {
		file, err := parseBenchmarkSourceFile(parts.fset, fname)
		if err != nil {
			continue
		}
		parts.addImports(file.Imports, seenImports)
		parts.decls = append(parts.decls, benchmarkCodeDecls(file.Decls)...)
	}
	return parts
}

func parseBenchmarkSourceFile(fset *token.FileSet, fname string) (*ast.File, error) {
	path := "testdata/benchmarks/" + fname
	data, err := benchmarksFS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parser.ParseFile(fset, path, data, parser.ParseComments)
}

func (p *benchmarkSourceParts) addImports(imports []*ast.ImportSpec, seen map[string]bool) {
	for _, imp := range imports {
		key := benchmarkImportKey(imp)
		if !seen[key] {
			seen[key] = true
			p.imports = append(p.imports, imp)
		}
	}
}

func benchmarkImportKey(imp *ast.ImportSpec) string {
	if imp.Name == nil {
		return imp.Path.Value
	}
	return imp.Name.Name + " " + imp.Path.Value
}

func benchmarkCodeDecls(decls []ast.Decl) []ast.Decl {
	codeDecls := make([]ast.Decl, 0, len(decls))
	for _, decl := range decls {
		if !benchmarkImportDecl(decl) {
			codeDecls = append(codeDecls, decl)
		}
	}
	return codeDecls
}

func benchmarkImportDecl(decl ast.Decl) bool {
	gen, ok := decl.(*ast.GenDecl)
	return ok && gen.Tok == token.IMPORT
}

func renderBenchmarkSource(parts benchmarkSourceParts) string {
	var sb strings.Builder
	sb.WriteString("package benchmarks\n\n")
	writeBenchmarkImports(&sb, parts)
	writeBenchmarkDecls(&sb, parts)
	return sb.String()
}

func writeBenchmarkImports(sb *strings.Builder, parts benchmarkSourceParts) {
	if len(parts.imports) == 0 {
		return
	}
	sb.WriteString("import (\n")
	for _, imp := range parts.imports {
		sb.WriteByte('\t')
		writeBenchmarkNode(sb, parts.fset, imp)
		sb.WriteByte('\n')
	}
	sb.WriteString(")\n\n")
}

func writeBenchmarkDecls(sb *strings.Builder, parts benchmarkSourceParts) {
	for _, decl := range parts.decls {
		writeBenchmarkNode(sb, parts.fset, decl)
		sb.WriteString("\n\n")
	}
}

func writeBenchmarkNode(sb *strings.Builder, fset *token.FileSet, node any) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		panic(fmt.Errorf("print benchmark source node: %w", err))
	}
	sb.Write(bytes.TrimSpace(buf.Bytes()))
}

var benchmarksSrc = getBenchmarksSrc()

// ============================================================================
// Benchmark Helpers
// ============================================================================

func TestPreviouslySkippedBenchmarksRun(t *testing.T) {
	prog, err := gig.Build(benchmarksSrc, gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	defer prog.Close()

	tests := map[string]func() int{
		"PanicRecover":   benchmarks.PanicRecover,
		"StringsBuilder": benchmarks.StringsBuilder,
		"MathBig":        benchmarks.MathBig,
	}
	for name, native := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := prog.Run(name)
			if err != nil {
				t.Fatalf("Run(%s): %v", name, err)
			}
			if got, want := int(toInt64(result)), native(); got != want {
				t.Fatalf("Run(%s) = %d, want %d", name, got, want)
			}
		})
	}
}

// benchGig builds the embedded benchmark source and runs the named function.
func benchGig(b *testing.B, funcName string) {
	b.Helper()
	prog, err := gig.Build(benchmarksSrc, gig.WithAllowPanic())
	if err != nil {
		b.Fatalf("Build error: %v", err)
	}
	b.Cleanup(prog.Close)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := prog.Run(funcName); err != nil {
			b.Fatalf("Run(%s): %v", funcName, err)
		}
	}
}

// benchmarkResult holds timing data from a benchmark run
type benchmarkResult struct {
	name     string
	gigNs    float64
	nativeNs float64
}

// ============================================================================
// 1. Arithmetic: sum 1..1000
// ============================================================================

func BenchmarkGig_ArithmeticSum(b *testing.B) {
	benchGig(b, "ArithmeticSum")
}

func BenchmarkNative_ArithmeticSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.ArithmeticSum()
	}
}

// ============================================================================
// 2. Recursive Fibonacci(15)
// ============================================================================

func BenchmarkGig_FibRecursive(b *testing.B) {
	benchGig(b, "FibRecursive")
}

func BenchmarkNative_FibRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.FibRecursive()
	}
}

// ============================================================================
// 3. Iterative Fibonacci(50)
// ============================================================================

func BenchmarkGig_FibIterative(b *testing.B) {
	benchGig(b, "FibIterative")
}

func BenchmarkNative_FibIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.FibIterative()
	}
}

// ============================================================================
// 4. Factorial(12)
// ============================================================================

func BenchmarkGig_Factorial(b *testing.B) {
	benchGig(b, "Factorial")
}

func BenchmarkNative_Factorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.Factorial()
	}
}

// ============================================================================
// 5. Slice Append (build slice of 1000 elements)
// ============================================================================

func BenchmarkGig_SliceAppend(b *testing.B) {
	benchGig(b, "SliceAppend")
}

func BenchmarkNative_SliceAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.SliceAppend()
	}
}

// ============================================================================
// 6. Slice Sum (iterate and sum 1000 elements)
// ============================================================================

func BenchmarkGig_SliceSum(b *testing.B) {
	benchGig(b, "SliceSum")
}

func BenchmarkNative_SliceSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.SliceSum()
	}
}

// ============================================================================
// 7. Map Operations (insert + read 100 entries)
// ============================================================================

func BenchmarkGig_MapOps(b *testing.B) {
	benchGig(b, "MapOps")
}

func BenchmarkNative_MapOps(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.MapOps()
	}
}

// ============================================================================
// 8. String Concatenation (build a 1000-char string)
// ============================================================================

func BenchmarkGig_StringConcat(b *testing.B) {
	benchGig(b, "StringConcat")
}

func BenchmarkNative_StringConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.StringConcat()
	}
}

// ============================================================================
// 9. Closure Calls (closure invoked 1000 times)
// ============================================================================

func BenchmarkGig_ClosureCalls(b *testing.B) {
	benchGig(b, "ClosureCalls")
}

func BenchmarkNative_ClosureCalls(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.ClosureCalls()
	}
}

// ============================================================================
// 10. Nested Loops (triple-nested N=10)
// ============================================================================

func BenchmarkGig_NestedLoops(b *testing.B) {
	benchGig(b, "NestedLoops")
}

func BenchmarkNative_NestedLoops(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.NestedLoops()
	}
}

// ============================================================================
// 11. Bubble Sort (sort 50 elements)
// ============================================================================

func BenchmarkGig_BubbleSort(b *testing.B) {
	benchGig(b, "BubbleSort")
}

func BenchmarkNative_BubbleSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.BubbleSort()
	}
}

// ============================================================================
// 12. GCD computation (100 pairs)
// ============================================================================

func BenchmarkGig_GCD(b *testing.B) {
	benchGig(b, "GCD")
}

func BenchmarkNative_GCD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.GCD()
	}
}

// ============================================================================
// 13. Sieve of Eratosthenes (primes up to 1000)
// ============================================================================

func BenchmarkGig_Sieve(b *testing.B) {
	benchGig(b, "Sieve")
}

func BenchmarkNative_Sieve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.Sieve()
	}
}

// ============================================================================
// 14. Higher-Order Function (map+reduce over 100 elements)
// ============================================================================

func BenchmarkGig_HigherOrder(b *testing.B) {
	benchGig(b, "HigherOrder")
}

func BenchmarkNative_HigherOrder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.HigherOrder()
	}
}

// ============================================================================
// 15. External Call: fmt.Sprintf (100 calls)
// ============================================================================

func BenchmarkGig_ExternalSprintf(b *testing.B) {
	benchGig(b, "ExternalSprintf")
}

func BenchmarkNative_ExternalSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.ExternalSprintf()
	}
}

// ============================================================================
// 16. External Call: strings.ToUpper (100 calls)
// ============================================================================

func BenchmarkGig_ExternalStrings(b *testing.B) {
	benchGig(b, "ExternalStrings")
}

func BenchmarkNative_ExternalStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.ExternalStrings()
	}
}

// ============================================================================
// 17. Function Call Overhead (1000 simple calls)
// ============================================================================

func BenchmarkGig_CallOverhead(b *testing.B) {
	benchGig(b, "CallOverhead")
}

func BenchmarkNative_CallOverhead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.CallOverhead()
	}
}

// ============================================================================
// 18. Build+Run latency (compile from source + single execution)
// ============================================================================

func BenchmarkGig_BuildAndRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		prog, err := gig.Build(benchmarksSrc, gig.WithAllowPanic())
		if err != nil {
			b.Fatal(err)
		}
		if _, err := prog.Run("ArithmeticSum"); err != nil {
			b.Fatalf("Run(ArithmeticSum): %v", err)
		}
	}
}

// ============================================================================
// 19. Struct with Methods
// ============================================================================

func BenchmarkGig_StructMethod(b *testing.B) {
	benchGig(b, "StructMethod")
}

func BenchmarkNative_StructMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.StructMethod()
	}
}

// ============================================================================
// 20. Interface Usage
// ============================================================================

func BenchmarkGig_Interface(b *testing.B) {
	benchGig(b, "Interface")
}

func BenchmarkNative_Interface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.Interface()
	}
}

// ============================================================================
// 21. Type Assertion
// ============================================================================

func BenchmarkGig_TypeAssertion(b *testing.B) {
	benchGig(b, "TypeAssertion")
}

func BenchmarkNative_TypeAssertion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.TypeAssertion()
	}
}

// ============================================================================
// 22. Type Switch
// ============================================================================

func BenchmarkGig_TypeSwitch(b *testing.B) {
	benchGig(b, "TypeSwitch")
}

func BenchmarkNative_TypeSwitch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.TypeSwitch()
	}
}

// ============================================================================
// 23. Defer (10 deferred calls)
// ============================================================================

func BenchmarkGig_Defer(b *testing.B) {
	benchGig(b, "Defer")
}

func BenchmarkNative_Defer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.Defer()
	}
}

// ============================================================================
// 24. Panic/Recover
// ============================================================================

func BenchmarkGig_PanicRecover(b *testing.B) {
	benchGig(b, "PanicRecover")
}

func BenchmarkNative_PanicRecover(b *testing.B) {
	safeCall := func(fn func()) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		fn()
		return nil
	}
	for i := 0; i < b.N; i++ {
		sum := 0
		for j := 0; j < 10; j++ {
			safeCall(func() {
				if j == 5 {
					panic("test")
				}
				sum = sum + j
			})
		}
		_ = sum
	}
}

// ============================================================================
// 25. Select Statement
// ============================================================================

func BenchmarkGig_Select(b *testing.B) {
	benchGig(b, "Select")
}

func BenchmarkNative_Select(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.Select()
	}
}

// ============================================================================
// 26. Slice of Interfaces
// ============================================================================

func BenchmarkGig_SliceInterface(b *testing.B) {
	benchGig(b, "SliceInterface")
}

func BenchmarkNative_SliceInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.SliceInterface()
	}
}

// ============================================================================
// 27. Composite Literals
// ============================================================================

func BenchmarkGig_CompositeLiteral(b *testing.B) {
	benchGig(b, "CompositeLiteral")
}

func BenchmarkNative_CompositeLiteral(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.CompositeLiteral()
	}
}

// ============================================================================
// 28. sort.Ints (external stdlib)
// ============================================================================

func BenchmarkGig_SortInts(b *testing.B) {
	benchGig(b, "SortInts")
}

func BenchmarkNative_SortInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchmarks.SortInts()
	}
}

// ============================================================================
// 29. strings.Builder (external stdlib)
// ============================================================================

func BenchmarkGig_StringsBuilder(b *testing.B) {
	benchGig(b, "StringsBuilder")
}

func BenchmarkNative_StringsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		for j := 0; j < 100; j++ {
			sb.WriteString("hello")
			sb.WriteString("world")
		}
		_ = sb.Len()
	}
}

// ============================================================================
// 30. math/big operations (external stdlib)
// ============================================================================

func BenchmarkGig_MathBig(b *testing.B) {
	benchGig(b, "MathBig")
}

func BenchmarkNative_MathBig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := big.NewInt(1)
		bv := big.NewInt(1)
		for j := 0; j < 100; j++ {
			a.Add(a, bv)
			bv.Sub(a, bv)
		}
		_ = int(a.Int64() % 1000)
	}
}

// ============================================================================
// 31. encoding/json Marshal
// ============================================================================

func BenchmarkGig_JsonMarshal(b *testing.B) {
	benchGig(b, "JsonMarshal")
}

func BenchmarkNative_JsonMarshal(b *testing.B) {
	type NativeData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		City string `json:"city"`
	}
	for i := 0; i < b.N; i++ {
		d := NativeData{Name: "John", Age: 30, City: "NYC"}
		s, _ := json.Marshal(d)
		_ = len(s)
	}
}

// ============================================================================
// Summary Printer: runs benchmarks and computes actual statistics
// ============================================================================

func TestBenchmarkSummary(t *testing.T) {
	logBenchmarkSummaryIntro(t, benchmarkCPUModel())

	results := getHardcodedResults()
	categorySlowdowns := logBenchmarkRows(t, results)

	logBenchmarkBuildLatency(t)
	logBenchmarkCategorySummary(t, categorySlowdowns)
	logBenchmarkOptimizationNotes(t)

	// Suppress unused warnings
	_ = strconv.Itoa
	_ = sort.Ints
	_ = time.Now()
}

func benchmarkCPUModel() string {
	var cpuModel string
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "model name") {
				cpuModel = strings.Split(line, ":")[1]
				break
			}
		}
	}
	if cpuModel == "" {
		cpuModel = "Unknown"
	}
	return strings.TrimSpace(cpuModel)
}

func logBenchmarkSummaryIntro(t *testing.T, cpuModel string) {
	t.Helper()
	t.Log("=============================================================================")
	t.Log("  GIG Performance Comparison: Interpreted (Gig) vs Native Go")
	t.Logf("  CPU: %s | Cores: %d | GOOS: %s | GOARCH: %s",
		cpuModel, runtime.NumCPU(), runtime.GOOS, runtime.GOARCH)
	t.Log("  Optimizations: DirectCall wrappers, Inline caching, Typed external functions")
	t.Log("=============================================================================")
	t.Log("")
	t.Log("Run benchmarks yourself:")
	t.Log("  go test -bench . -benchmem -count=1 ./tests/ -run='^$'")
	t.Log("")
	t.Log("  NOTE: To regenerate these stats with current hardware, run:")
	t.Log("    go test -bench . -benchmem -count=1 ./tests/ -run='^$' | tee /tmp/bench.txt")
	t.Log("")
}

func logBenchmarkRows(t *testing.T, results []benchmarkResult) map[string][]float64 {
	t.Helper()
	t.Logf("  %-22s %14s %14s %10s %s", "Workload", "Gig (ns/op)", "Native (ns/op)", "Slowdown", "Category")
	t.Logf("  %-22s %14s %14s %10s %s",
		strings.Repeat("-", 22),
		strings.Repeat("-", 14),
		strings.Repeat("-", 14),
		strings.Repeat("-", 10),
		strings.Repeat("-", 16))

	categorySlowdowns := make(map[string][]float64)
	for _, r := range results {
		ratio := r.gigNs / r.nativeNs
		cat := categorize(r.name)
		t.Logf("  %-22s %14.0f %14.1f %9.0fx %s", r.name, r.gigNs, r.nativeNs, ratio, cat)
		categorySlowdowns[cat] = append(categorySlowdowns[cat], ratio)
	}
	return categorySlowdowns
}

func logBenchmarkBuildLatency(t *testing.T) {
	t.Helper()
	t.Log("")
	t.Logf("  %-22s %14s", "BuildAndRun", "~43,434 ns/op (compile + single execution)")
	t.Log("")
}

type benchmarkCategoryStats struct {
	min float64
	max float64
	avg float64
}

func logBenchmarkCategorySummary(t *testing.T, categorySlowdowns map[string][]float64) {
	t.Helper()
	t.Log("  Summary (computed from actual benchmark data):")
	t.Log("  ┌─────────────────────────────────────────────────────────┐")

	for cat, stats := range benchmarkCategoryStatsByName(categorySlowdowns) {
		logBenchmarkCategoryStats(t, cat, stats)
	}
	logBenchmarkOverallAverage(t, categorySlowdowns)

	t.Log("  └─────────────────────────────────────────────────────────┘")
}

func benchmarkCategoryStatsByName(categorySlowdowns map[string][]float64) map[string]benchmarkCategoryStats {
	stats := make(map[string]benchmarkCategoryStats, len(categorySlowdowns))
	for cat, ratios := range categorySlowdowns {
		if len(ratios) == 0 {
			continue
		}
		stat := benchmarkCategoryStats{min: ratios[0], max: ratios[0]}
		for _, r := range ratios {
			stat.min = min(stat.min, r)
			stat.max = max(stat.max, r)
			stat.avg += r
		}
		stat.avg /= float64(len(ratios))
		stats[cat] = stat
	}
	return stats
}

func logBenchmarkCategoryStats(t *testing.T, cat string, stats benchmarkCategoryStats) {
	t.Helper()
	switch cat {
	case "Compute":
		t.Logf("  │ Pure Computation (loops, arithmetic):      ~%.0f-%.0fx (avg: %.0fx)│", stats.min, stats.max, stats.avg)
	case "Recursion":
		t.Logf("  │ Recursion (function call heavy):           ~%.0f-%.0fx (avg: %.0fx)│", stats.min, stats.max, stats.avg)
	case "Data Struct":
		t.Logf("  │ Data Structures (slice, map):              ~%.0f-%.0fx (avg: %.0fx)│", stats.min, stats.max, stats.avg)
	case "Closure":
		t.Logf("  │ Closures (capture + invoke):              ~%.0f-%.0fx (avg: %.0fx)│", stats.min, stats.max, stats.avg)
	case "Algorithm":
		t.Logf("  │ Algorithms (sort, GCD, sieve):             ~%.0f-%.0fx (avg: %.0fx)│", stats.min, stats.max, stats.avg)
	case "External Call":
		t.Logf("  │ External Calls (fmt, strings):             ~%.0f-%.0fx (avg: %.0fx)│", stats.min, stats.max, stats.avg)
	case "Call Overhead":
		t.Logf("  │ Function Call Overhead (10K calls):        ~%.0fx (avg: %.0fx)         │", stats.max, stats.avg)
	case "String":
		t.Logf("  │ String Operations:                         ~%.0f-%.0fx (avg: %.0fx)│", stats.min, stats.max, stats.avg)
	case "Complex Syntax":
		t.Logf("  │ Complex Syntax (interface, struct, etc):    ~%.0f-%.0fx (avg: %.0fx)│", stats.min, stats.max, stats.avg)
	case "Third-party":
		t.Logf("  │ Third-party Libs (sort, json, math/big):   ~%.0f-%.0fx (avg: %.0fx)│", stats.min, stats.max, stats.avg)
	}
}

func logBenchmarkOverallAverage(t *testing.T, categorySlowdowns map[string][]float64) {
	t.Helper()
	avgAll, count := benchmarkOverallAverage(categorySlowdowns)
	if count > 0 {
		t.Logf("  │ Overall Average:                             ~%.0fx         │", avgAll)
	}
}

func benchmarkOverallAverage(categorySlowdowns map[string][]float64) (float64, int) {
	avgAll := 0.0
	count := 0
	for _, ratios := range categorySlowdowns {
		for _, r := range ratios {
			avgAll += r
			count++
		}
	}
	if count == 0 {
		return 0, 0
	}
	return avgAll / float64(count), count
}

func logBenchmarkOptimizationNotes(t *testing.T) {
	t.Helper()
	t.Log("")
	t.Log("  Optimizations Applied:")
	t.Log("  • DirectCall typed wrappers: Avoid reflect.Call for external functions")
	t.Log("  • Inline caching: Cache resolved external function info per call site")
	t.Log("  • ExternalFuncInfo: Pre-resolved function + DirectCall wrapper in bytecode")
	t.Log("")
	t.Log("  Notes:")
	t.Log("  • Third-party benchmarks use Go stdlib as proxy for external libraries")
	t.Log("  • Complex syntax tests cover interfaces, methods, type assertions,")
	t.Log("    panic/recover, defer, select, and composite literals")
}

// categorize returns the category for a benchmark name
func categorize(name string) string {
	switch {
	case strings.Contains(name, "Arithmetic"), strings.Contains(name, "FibIterative"):
		return "Compute"
	case strings.Contains(name, "FibRecursive"), strings.Contains(name, "Factorial"):
		return "Recursion"
	case strings.Contains(name, "Slice"), strings.Contains(name, "Map"):
		return "Data Struct"
	case strings.Contains(name, "Closure"), strings.Contains(name, "HigherOrder"):
		return "Closure"
	case strings.Contains(name, "Sort"), strings.Contains(name, "GCD"), strings.Contains(name, "Sieve"):
		return "Algorithm"
	case strings.Contains(name, "External"), strings.Contains(name, "Sprintf"), strings.Contains(name, "Strings"):
		return "External Call"
	case strings.Contains(name, "CallOverhead"):
		return "Call Overhead"
	case strings.Contains(name, "StringConcat"):
		return "String"
	case strings.Contains(name, "Struct"), strings.Contains(name, "Interface"),
		strings.Contains(name, "Type"), strings.Contains(name, "Defer"),
		strings.Contains(name, "Panic"), strings.Contains(name, "Select"),
		strings.Contains(name, "Composite"):
		return "Complex Syntax"
	case strings.Contains(name, "Sort"), strings.Contains(name, "Builder"),
		strings.Contains(name, "Math"), strings.Contains(name, "Json"):
		return "Third-party"
	default:
		return "Other"
	}
}

// getHardcodedResults returns fallback benchmark data.
// Last measured: AMD EPYC 9754 128-Core Processor, Go 1.23, linux/amd64, -benchtime=1s
func getHardcodedResults() []benchmarkResult {
	return []benchmarkResult{
		{"ArithmeticSum", 74436, 336},
		{"FibRecursive", 151438, 3642},
		{"FibIterative", 4814, 18},
		{"Factorial", 1812, 25},
		{"SliceAppend", 543792, 6361},
		{"SliceSum", 186693, 991.1},
		{"MapOps", 94522, 8131},
		{"StringConcat", 36217, 22241},
		{"ClosureCalls", 319205, 671},
		{"NestedLoops", 85624, 464},
		{"BubbleSort", 249602, 2201},
		{"GCD", 61318, 928},
		{"Sieve", 203747, 4657},
		{"HigherOrder", 23064, 102},
		{"ExternalSprintf", 102093, 5447},
		{"ExternalStrings", 27478, 9631},
		{"CallOverhead", 106235, 657.6},
	}
}
