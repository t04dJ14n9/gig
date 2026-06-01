package benchmarks

import (
	"context"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

const extremeRuleSource = `
import (
	"strings"
	"strconv"
	"math"
)

func EvaluateRule(userID int, name string, score float64) string {
	upper := strings.ToUpper(name)
	trimmed := strings.TrimSpace(upper)
	idStr := strconv.Itoa(userID)
	adjusted := math.Sqrt(score) * 10.0
	tier := "BRONZE"
	if adjusted > 80.0 {
		tier = "GOLD"
	} else if adjusted > 50.0 {
		tier = "SILVER"
	}
	result := trimmed + "#" + idStr + ":" + tier
	if strings.Contains(result, "GOLD") {
		result = result + " [VIP]"
	}
	return result
}
`

func nativeRule(userID int, name string, score float64) string {
	upper := strings.ToUpper(name)
	trimmed := strings.TrimSpace(upper)
	idStr := strconv.Itoa(userID)
	adjusted := math.Sqrt(score) * 10.0
	tier := "BRONZE"
	if adjusted > 80.0 {
		tier = "GOLD"
	} else if adjusted > 50.0 {
		tier = "SILVER"
	}
	result := trimmed + "#" + idStr + ":" + tier
	if strings.Contains(result, "GOLD") {
		result = result + " [VIP]"
	}
	return result
}

func runStressLevel(concurrency int, duration time.Duration, fn func(gID, i int) error) (totalOps int64, totalErrors int64, heapMB float64, gcPauses uint32) {
	var ops atomic.Int64
	var errors atomic.Int64

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var memBefore, memAfter runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memBefore)

	var wg sync.WaitGroup
	for g := 0; g < concurrency; g++ {
		wg.Add(1)
		go func(gID int) {
			defer wg.Done()
			i := 0
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				if err := fn(gID, i); err != nil {
					errors.Add(1)
				}
				ops.Add(1)
				i++
			}
		}(g)
	}
	wg.Wait()

	runtime.ReadMemStats(&memAfter)
	return ops.Load(), errors.Load(),
		float64(memAfter.TotalAlloc-memBefore.TotalAlloc) / 1024 / 1024,
		memAfter.NumGC - memBefore.NumGC
}

type extremeStressConfig struct {
	levels       []int
	nativeLevels []int
	duration     time.Duration
	rounds       int
}

type extremeStressStats struct {
	throughput float64
	avgLatUs   float64
	errors     int64
	heapMB     float64
	gcPauses   uint32
}

type extremeStressRollup struct {
	medianThroughput float64
	medianLatencyUs  float64
	totalErrors      int64
	avgHeapMB        float64
	avgGCPauses      float64
}

func TestExtremeStress(t *testing.T) {
	cfg := defaultExtremeStressConfig()
	prog, err := gig.Build(extremeRuleSource)
	if err != nil {
		t.Fatal(err)
	}
	defer prog.Close()

	r, err := prog.Run("EvaluateRule", 1, " alice ", 81.0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Correctness: %v", r)

	gigStats := runGigExtremeStress(t, prog, cfg)
	logGigExtremeResults(t, cfg, gigStats)
	nativeStats := runNativeExtremeStress(cfg)
	logNativeExtremeResults(t, cfg, nativeStats)
	logExtremeThroughputRatios(t, cfg, gigStats, nativeStats)
}

func defaultExtremeStressConfig() extremeStressConfig {
	return extremeStressConfig{
		levels:       []int{1, 10, 100, 500, 1000, 2000, 5000, 10000},
		nativeLevels: []int{1, 100, 1000, 10000},
		duration:     3 * time.Second,
		rounds:       3,
	}
}

func newExtremeStressTable(levels []int, rounds int) [][]extremeStressStats {
	stats := make([][]extremeStressStats, len(levels))
	for i := range stats {
		stats[i] = make([]extremeStressStats, rounds)
	}
	return stats
}

func runGigExtremeStress(t *testing.T, prog *gig.Program, cfg extremeStressConfig) [][]extremeStressStats {
	t.Helper()
	stats := newExtremeStressTable(cfg.levels, cfg.rounds)
	t.Logf("Running %d rounds × %d levels...", cfg.rounds, len(cfg.levels))

	for round := 0; round < cfg.rounds; round++ {
		for li, concurrency := range cfg.levels {
			stats[li][round] = runGigExtremeStressLevel(prog, concurrency, cfg.duration)
		}
	}
	return stats
}

func runGigExtremeStressLevel(prog *gig.Program, concurrency int, duration time.Duration) extremeStressStats {
	ops, errs, heap, gc := runStressLevel(concurrency, duration, func(gID, i int) error {
		execCtx, execCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer execCancel()
		_, err := prog.RunWithContext(execCtx, "EvaluateRule", gID*10000+i, " bob ", float64(30+i%70))
		return err
	})
	return extremeStressStatsFromRun(concurrency, duration, ops, errs, heap, gc)
}

func runNativeExtremeStress(cfg extremeStressConfig) [][]extremeStressStats {
	stats := newExtremeStressTable(cfg.nativeLevels, cfg.rounds)
	for round := 0; round < cfg.rounds; round++ {
		for li, concurrency := range cfg.nativeLevels {
			stats[li][round] = runNativeExtremeStressLevel(concurrency, cfg.duration)
		}
	}
	return stats
}

func runNativeExtremeStressLevel(concurrency int, duration time.Duration) extremeStressStats {
	ops, _, heap, gc := runStressLevel(concurrency, duration, func(gID, i int) error {
		_ = nativeRule(gID*10000+i, " bob ", float64(30+i%70))
		return nil
	})
	return extremeStressStatsFromRun(concurrency, duration, ops, 0, heap, gc)
}

func extremeStressStatsFromRun(concurrency int, duration time.Duration, ops, errs int64, heap float64, gc uint32) extremeStressStats {
	throughput := float64(ops) / duration.Seconds()
	return extremeStressStats{
		throughput: throughput,
		avgLatUs:   1e6 / throughput * float64(concurrency),
		errors:     errs,
		heapMB:     heap,
		gcPauses:   gc,
	}
}

func logGigExtremeResults(t *testing.T, cfg extremeStressConfig, gigStats [][]extremeStressStats) {
	t.Helper()
	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║              GIG EXTREME STRESS TEST (%d rounds, %v/round)              ║", cfg.rounds, cfg.duration)
	t.Logf("║  Workload: rule engine (strings + math + stdlib)  |  CPU: %2d cores        ║", runtime.NumCPU())
	t.Logf("╠══════════╦═══════════════╦════════════╦════════╦══════════╦════════════════╣")
	t.Logf("║ Goroutin ║  Throughput   ║  Avg Lat   ║ Errors ║ Heap(MB) ║   GC Pauses    ║")
	t.Logf("╠══════════╬═══════════════╬════════════╬════════╬══════════╬════════════════╣")

	for li, concurrency := range cfg.levels {
		rollup := rollupExtremeStress(gigStats[li], cfg.rounds)
		t.Logf("║ %8d ║ %11.0f/s ║ %8.1f μs ║ %6d ║ %6.0f MB ║ %14.0f ║",
			concurrency, rollup.medianThroughput, rollup.medianLatencyUs, rollup.totalErrors, rollup.avgHeapMB, rollup.avgGCPauses)

		if rollup.totalErrors > 0 {
			t.Errorf("  ⚠ %d errors at concurrency=%d", rollup.totalErrors, concurrency)
		}
	}
	t.Logf("╚══════════╩═══════════════╩════════════╩════════╩══════════╩════════════════╝")
}

func logNativeExtremeResults(t *testing.T, cfg extremeStressConfig, nativeStats [][]extremeStressStats) {
	t.Helper()
	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║                    NATIVE GO BASELINE (throughput only)                     ║")
	t.Logf("╠══════════╦═══════════════╦══════════╦════════════════╗")
	t.Logf("║ Goroutin ║  Throughput   ║ Heap(MB) ║   GC Pauses    ║")
	t.Logf("╠══════════╬═══════════════╬══════════╬════════════════╣")

	for li, concurrency := range cfg.nativeLevels {
		rollup := rollupExtremeStress(nativeStats[li], cfg.rounds)
		t.Logf("║ %8d ║ %11.0f/s ║ %6.0f MB ║ %14.0f ║",
			concurrency, rollup.medianThroughput, rollup.avgHeapMB, rollup.avgGCPauses)
	}
	t.Logf("╚══════════╩═══════════════╩══════════╩════════════════╝")
}

func logExtremeThroughputRatios(
	t *testing.T,
	cfg extremeStressConfig,
	gigStats [][]extremeStressStats,
	nativeStats [][]extremeStressStats,
) {
	t.Helper()
	t.Logf("")
	t.Logf("Throughput ratio (Native / Gig):")
	for li, concurrency := range cfg.nativeLevels {
		gigTP := gigExtremeThroughputAt(cfg, gigStats, concurrency)
		if gigTP > 0 {
			nativeTP := rollupExtremeStress(nativeStats[li], cfg.rounds).medianThroughput
			t.Logf("  %5dG: Native %.0f/s vs Gig %.0f/s = %.1fx", concurrency, nativeTP, gigTP, nativeTP/gigTP)
		}
	}
}

func gigExtremeThroughputAt(cfg extremeStressConfig, gigStats [][]extremeStressStats, concurrency int) float64 {
	for li, level := range cfg.levels {
		if level == concurrency {
			return rollupExtremeStress(gigStats[li], cfg.rounds).medianThroughput
		}
	}
	return 0
}

func rollupExtremeStress(samples []extremeStressStats, rounds int) extremeStressRollup {
	var tps, lats []float64
	var totalErrs int64
	var totalHeap float64
	var totalGC uint32
	for _, s := range samples {
		tps = append(tps, s.throughput)
		lats = append(lats, s.avgLatUs)
		totalErrs += s.errors
		totalHeap += s.heapMB
		totalGC += s.gcPauses
	}
	return extremeStressRollup{
		medianThroughput: median3(tps[0], tps[1], tps[2]),
		medianLatencyUs:  median3(lats[0], lats[1], lats[2]),
		totalErrors:      totalErrs,
		avgHeapMB:        totalHeap / float64(rounds),
		avgGCPauses:      float64(totalGC) / float64(rounds),
	}
}

func median3(a, b, c float64) float64 {
	if a > b {
		a, b = b, a
	}
	if b > c {
		b, c = c, b
	}
	if a > b {
		a, b = b, a
	}
	return b
}

// ============================================================================
// Stateful Stress Test: concurrent execution with persistent globals
// ============================================================================

// statefulRuleSource uses global variables protected by *sync.Mutex.
// Must use pointer form (*sync.Mutex) due to Gig limitation: value-type
// struct method calls on globals are not supported.
const statefulRuleSource = `
import (
	"strings"
	"strconv"
	"math"
	"sync"
)

var counter int
var mu *sync.Mutex

func init() {
	mu = &sync.Mutex{}
	counter = 0
}

func EvaluateRuleStateful(userID int, name string, score float64) string {
	upper := strings.ToUpper(name)
	trimmed := strings.TrimSpace(upper)
	idStr := strconv.Itoa(userID)
	adjusted := math.Sqrt(score) * 10.0
	tier := "BRONZE"
	if adjusted > 80.0 {
		tier = "GOLD"
	} else if adjusted > 50.0 {
		tier = "SILVER"
	}
	result := trimmed + "#" + idStr + ":" + tier
	if strings.Contains(result, "GOLD") {
		result = result + " [VIP]"
	}
	mu.Lock()
	counter++
	mu.Unlock()
	return result
}

func GetCounter() int {
	mu.Lock()
	defer mu.Unlock()
	return counter
}
`

// TestStatefulStress tests concurrent stateful execution with WithStatefulGlobals().
// Uses 5 concurrency levels with 1 round each to keep test time reasonable.
func TestStatefulStress(t *testing.T) {
	cfg := defaultStatefulStressConfig()
	verifyStatefulRuleCorrectness(t)
	gigStats := runStatefulStressTable(t, cfg)
	logStatefulStressResults(t, cfg, gigStats)
	verifyStatefulCounter(t, cfg)
}

type statefulStressConfig struct {
	levels            []int
	duration          time.Duration
	verifyConcurrency int
	verifyOpsPerG     int
}

func defaultStatefulStressConfig() statefulStressConfig {
	return statefulStressConfig{
		levels:            []int{1, 100, 1000, 5000, 10000},
		duration:          3 * time.Second,
		verifyConcurrency: 100,
		verifyOpsPerG:     100,
	}
}

func buildStatefulStressProgram(t *testing.T) *gig.Program {
	t.Helper()
	prog, err := gig.Build(statefulRuleSource, gig.WithStatefulGlobals(), gig.WithAllowPanic())
	if err != nil {
		t.Fatal(err)
	}
	return prog
}

func verifyStatefulRuleCorrectness(t *testing.T) {
	t.Helper()
	prog := buildStatefulStressProgram(t)
	defer prog.Close()

	result, err := prog.Run("EvaluateRuleStateful", 1, " alice ", 81.0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Correctness: %v", result)
}

func runStatefulStressTable(t *testing.T, cfg statefulStressConfig) []extremeStressStats {
	t.Helper()
	t.Logf("Running stateful stress test (%d levels, %v/level)...", len(cfg.levels), cfg.duration)

	stats := make([]extremeStressStats, len(cfg.levels))
	for li, concurrency := range cfg.levels {
		stats[li] = runStatefulStressLevel(t, concurrency, cfg.duration)
	}
	return stats
}

func runStatefulStressLevel(t *testing.T, concurrency int, duration time.Duration) extremeStressStats {
	t.Helper()
	prog := buildStatefulStressProgram(t)
	defer prog.Close()

	ops, errs, heap, gc := runStressLevel(concurrency, duration, func(gID, i int) error {
		execCtx, execCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer execCancel()
		_, err := prog.RunWithContext(execCtx, "EvaluateRuleStateful", gID*10000+i, " bob ", float64(30+i%70))
		return err
	})
	return extremeStressStatsFromRun(concurrency, duration, ops, errs, heap, gc)
}

func logStatefulStressResults(t *testing.T, cfg statefulStressConfig, gigStats []extremeStressStats) {
	t.Helper()
	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║           STATEFUL STRESS TEST (WithStatefulGlobals)                       ║")
	t.Logf("║  Workload: rule engine + global counter + *sync.Mutex  |  CPU: %2d cores  ║", runtime.NumCPU())
	t.Logf("╠══════════╦═══════════════╦════════════╦════════╦══════════╦════════════════╣")
	t.Logf("║ Goroutin ║  Throughput   ║  Avg Lat   ║ Errors ║ Heap(MB) ║   GC Pauses    ║")
	t.Logf("╠══════════╬═══════════════╬════════════╬════════╬══════════╬════════════════╣")

	for li, concurrency := range cfg.levels {
		s := gigStats[li]
		t.Logf("║ %8d ║ %11.0f/s ║ %8.1f μs ║ %6d ║ %6.0f MB ║ %14.0f ║",
			concurrency, s.throughput, s.avgLatUs, s.errors, s.heapMB, float64(s.gcPauses))
		if s.errors > 0 {
			t.Errorf("  ⚠ %d errors at concurrency=%d", s.errors, concurrency)
		}
	}
	t.Logf("╚══════════╩═══════════════╩════════════╩════════╩══════════╩════════════════╝")
}

func verifyStatefulCounter(t *testing.T, cfg statefulStressConfig) {
	t.Helper()
	prog := buildStatefulStressProgram(t)
	defer prog.Close()

	var wg sync.WaitGroup
	for g := 0; g < cfg.verifyConcurrency; g++ {
		wg.Add(1)
		go func(gID int) {
			defer wg.Done()
			for i := 0; i < cfg.verifyOpsPerG; i++ {
				prog.Run("EvaluateRuleStateful", gID*1000+i, "test", 50.0)
			}
		}(g)
	}
	wg.Wait()

	result, err := prog.Run("GetCounter")
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	counterVal, ok := result.(int)
	if !ok {
		t.Fatalf("GetCounter returned %T, want int", result)
	}
	expected := cfg.verifyConcurrency * cfg.verifyOpsPerG
	if counterVal != expected {
		t.Errorf("Counter = %d, want %d (lost %d updates)", counterVal, expected, expected-counterVal)
	} else {
		t.Logf("Counter verification: %d ops = %d (correct, no lost updates)", counterVal, expected)
	}
}

// ============================================================================
// Go Native vs Gig Concurrent Globals Comparison
// ============================================================================

const unprotectedGlobalSource = `
var counter int

func Increment() int {
	counter++
	return counter
}

func GetCounter() int {
	return counter
}
`

const mutexProtectedGlobalSource = `
import "sync"

var counter int
var mu *sync.Mutex

func init() {
	mu = &sync.Mutex{}
}

func Increment() int {
	mu.Lock()
	counter++
	mu.Unlock()
	return counter
}

func GetCounter() int {
	mu.Lock()
	defer mu.Unlock()
	return counter
}
`

type concurrentGlobalsConfig struct {
	concurrency     int
	opsPerGoroutine int
	totalOps        int
}

type counterRunResult struct {
	counter int64
	elapsed time.Duration
}

// TestConcurrentGlobals_GoNative_vs_Gig compares Go native and Gig
// behavior when concurrently accessing global variables.
// Uses moderate concurrency to keep test time reasonable.
func TestConcurrentGlobals_GoNative_vs_Gig(t *testing.T) {
	cfg := defaultConcurrentGlobalsConfig()
	logConcurrentGlobalsIntro(t, cfg)

	nativeUnprotected := runNativeUnprotectedGlobal(cfg)
	logUnprotectedGlobalReport(t, "Go Native: Unprotected Global", cfg, nativeUnprotected, time.Microsecond)

	nativeMutex := runNativeMutexGlobal(cfg)
	logPreciseGlobalReport(t, "Go Native: *sync.Mutex Protected Global", cfg, nativeMutex, time.Microsecond)

	nativeAtomic := runNativeAtomicGlobal(cfg)
	logPreciseGlobalReport(t, "Go Native: sync/atomic Protected Global", cfg, nativeAtomic, time.Microsecond)

	gigUnprotected := runGigGlobalCounter(t, unprotectedGlobalSource, cfg)
	logUnprotectedGlobalReport(
		t,
		"Gig (Stateful): Unprotected Global",
		cfg,
		gigUnprotected,
		time.Millisecond,
		"│ Note: SharedGlobals RWMutex prevents torn reads,   │",
		"│       but read-modify-write still has lost updates  │",
	)

	gigMutex := runGigGlobalCounter(t, mutexProtectedGlobalSource, cfg)
	logPreciseGlobalReport(t, "Gig (Stateful): *sync.Mutex Protected Global", cfg, gigMutex, time.Millisecond)

	logConcurrentGlobalsSummary(t, cfg, nativeUnprotected, nativeMutex, nativeAtomic, gigUnprotected, gigMutex)
	logMutexProtectedEquivalence(t, cfg, nativeMutex, gigMutex)
}

func defaultConcurrentGlobalsConfig() concurrentGlobalsConfig {
	const concurrency = 50
	const opsPerGoroutine = 100
	return concurrentGlobalsConfig{
		concurrency:     concurrency,
		opsPerGoroutine: opsPerGoroutine,
		totalOps:        concurrency * opsPerGoroutine,
	}
}

func logConcurrentGlobalsIntro(t *testing.T, cfg concurrentGlobalsConfig) {
	t.Helper()
	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║    CONCURRENT GLOBALS: Go Native vs Gig (%dG × %d ops)               ║", cfg.concurrency, cfg.opsPerGoroutine)
	t.Logf("╚══════════════════════════════════════════════════════════════════════════════╝")
}

func runNativeUnprotectedGlobal(cfg concurrentGlobalsConfig) counterRunResult {
	var counter int
	var wg sync.WaitGroup
	start := time.Now()
	for g := 0; g < cfg.concurrency; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < cfg.opsPerGoroutine; i++ {
				counter++
			}
		}()
	}
	wg.Wait()
	return counterRunResult{counter: int64(counter), elapsed: time.Since(start)}
}

func runNativeMutexGlobal(cfg concurrentGlobalsConfig) counterRunResult {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup
	start := time.Now()
	for g := 0; g < cfg.concurrency; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < cfg.opsPerGoroutine; i++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return counterRunResult{counter: int64(counter), elapsed: time.Since(start)}
}

func runNativeAtomicGlobal(cfg concurrentGlobalsConfig) counterRunResult {
	var counter int64
	var wg sync.WaitGroup
	start := time.Now()
	for g := 0; g < cfg.concurrency; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < cfg.opsPerGoroutine; i++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}
	wg.Wait()
	return counterRunResult{counter: counter, elapsed: time.Since(start)}
}

func runGigGlobalCounter(t *testing.T, source string, cfg concurrentGlobalsConfig) counterRunResult {
	t.Helper()
	prog, err := gig.Build(source, gig.WithStatefulGlobals(), gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build source: %v", err)
	}
	defer prog.Close()

	var wg sync.WaitGroup
	start := time.Now()
	for g := 0; g < cfg.concurrency; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < cfg.opsPerGoroutine; i++ {
				prog.Run("Increment")
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("GetCounter")
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	counter, ok := result.(int)
	if !ok {
		t.Fatalf("GetCounter returned %T, want int", result)
	}
	return counterRunResult{counter: int64(counter), elapsed: time.Since(start)}
}

func logUnprotectedGlobalReport(
	t *testing.T,
	title string,
	cfg concurrentGlobalsConfig,
	result counterRunResult,
	round time.Duration,
	noteLines ...string,
) {
	t.Helper()
	t.Logf("")
	t.Logf("┌─────────────────────────────────────────────────────┐")
	t.Logf("│ %-51s │", title)
	t.Logf("├─────────────────────────────────────────────────────┤")
	t.Logf("│ Total ops:      %10d                        │", cfg.totalOps)
	t.Logf("│ Final counter:  %10d (lost %d updates)       │", result.counter, int64(cfg.totalOps)-result.counter)
	t.Logf("│ Has data race:  %10v                        │", result.counter != int64(cfg.totalOps))
	t.Logf("│ Elapsed:        %10v                        │", result.elapsed.Round(round))
	for _, line := range noteLines {
		t.Logf("%s", line)
	}
	t.Logf("└─────────────────────────────────────────────────────┘")
}

func logPreciseGlobalReport(t *testing.T, title string, cfg concurrentGlobalsConfig, result counterRunResult, round time.Duration) {
	t.Helper()
	t.Logf("")
	t.Logf("┌─────────────────────────────────────────────────────┐")
	t.Logf("│ %-51s │", title)
	t.Logf("├─────────────────────────────────────────────────────┤")
	t.Logf("│ Total ops:      %10d                        │", cfg.totalOps)
	t.Logf("│ Final counter:  %10d                        │", result.counter)
	t.Logf("│ Is precise:     %10v                        │", result.counter == int64(cfg.totalOps))
	t.Logf("│ Elapsed:        %10v                        │", result.elapsed.Round(round))
	t.Logf("└─────────────────────────────────────────────────────┘")
}

func logConcurrentGlobalsSummary(
	t *testing.T,
	cfg concurrentGlobalsConfig,
	nativeUnprotected, nativeMutex, nativeAtomic, gigUnprotected, gigMutex counterRunResult,
) {
	t.Helper()
	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║                         COMPARISON SUMMARY                                 ║")
	t.Logf("╠══════════════════════════════════════╦════════════╦═══════════╦════════════╣")
	t.Logf("║ Mode                                 ║ Counter    ║ Precise?  ║ Has Race?  ║")
	t.Logf("╠══════════════════════════════════════╬════════════╬═══════════╬════════════╣")
	t.Logf("║ Go Native: Unprotected               ║ %10d ║ %9v ║ %10v ║",
		nativeUnprotected.counter, nativeUnprotected.counter == int64(cfg.totalOps), nativeUnprotected.counter != int64(cfg.totalOps))
	t.Logf("║ Go Native: *sync.Mutex               ║ %10d ║ %9v ║ %10v ║",
		nativeMutex.counter, nativeMutex.counter == int64(cfg.totalOps), nativeMutex.counter != int64(cfg.totalOps))
	t.Logf("║ Go Native: sync/atomic               ║ %10d ║ %9v ║ %10v ║",
		nativeAtomic.counter, nativeAtomic.counter == int64(cfg.totalOps), nativeAtomic.counter != int64(cfg.totalOps))
	t.Logf("║ Gig Stateful: Unprotected            ║ %10d ║ %9v ║ %10v ║",
		gigUnprotected.counter, gigUnprotected.counter == int64(cfg.totalOps), gigUnprotected.counter != int64(cfg.totalOps))
	t.Logf("║ Gig Stateful: *sync.Mutex            ║ %10d ║ %9v ║ %10v ║",
		gigMutex.counter, gigMutex.counter == int64(cfg.totalOps), gigMutex.counter != int64(cfg.totalOps))
	t.Logf("╚══════════════════════════════════════╩════════════╩═══════════╩════════════╝")
}

func logMutexProtectedEquivalence(t *testing.T, cfg concurrentGlobalsConfig, nativeMutex, gigMutex counterRunResult) {
	t.Helper()
	if gigMutex.counter == int64(cfg.totalOps) && nativeMutex.counter == int64(cfg.totalOps) {
		t.Logf("✓ Mutex-protected globals: Gig and Go native both produce precise results")
	}
}
