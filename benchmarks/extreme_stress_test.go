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

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
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

func TestExtremeStress(t *testing.T) {
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

	levels := []int{1, 10, 100, 500, 1000, 2000, 5000, 10000}
	duration := 3 * time.Second
	rounds := 3

	// Collect stats: [level][round]
	type stats struct {
		throughput float64
		avgLatUs   float64
		errors     int64
		heapMB     float64
		gcPauses   uint32
	}

	gigStats := make([][]stats, len(levels))
	for i := range gigStats {
		gigStats[i] = make([]stats, rounds)
	}

	t.Logf("Running %d rounds × %d levels...", rounds, len(levels))

	for round := 0; round < rounds; round++ {
		for li, concurrency := range levels {
			ops, errs, heap, gc := runStressLevel(concurrency, duration, func(gID, i int) error {
			execCtx, execCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer execCancel()
			_, err := prog.RunWithContext(execCtx, "EvaluateRule", gID*10000+i, " bob ", float64(30+i%70))
			return err
			})
			tp := float64(ops) / duration.Seconds()
			gigStats[li][round] = stats{
				throughput: tp,
				avgLatUs:   1e6 / tp * float64(concurrency),
				errors:     errs,
				heapMB:     heap,
				gcPauses:   gc,
			}
		}
	}

	// Print Gig results (median of 3 rounds)
	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║              GIG EXTREME STRESS TEST (%d rounds, %v/round)              ║", rounds, duration)
	t.Logf("║  Workload: rule engine (strings + math + stdlib)  |  CPU: %2d cores        ║", runtime.NumCPU())
	t.Logf("╠══════════╦═══════════════╦════════════╦════════╦══════════╦════════════════╣")
	t.Logf("║ Goroutin ║  Throughput   ║  Avg Lat   ║ Errors ║ Heap(MB) ║   GC Pauses    ║")
	t.Logf("╠══════════╬═══════════════╬════════════╬════════╬══════════╬════════════════╣")

	for li, concurrency := range levels {
		// Take median throughput
		var tps, lats []float64
		var totalErrs int64
		var totalHeap float64
		var totalGC uint32
		for round := 0; round < rounds; round++ {
			s := gigStats[li][round]
			tps = append(tps, s.throughput)
			lats = append(lats, s.avgLatUs)
			totalErrs += s.errors
			totalHeap += s.heapMB
			totalGC += s.gcPauses
		}
		// Simple median: sort 3 values, pick middle
		medTP := median3(tps[0], tps[1], tps[2])
		medLat := median3(lats[0], lats[1], lats[2])
		avgHeap := totalHeap / float64(rounds)
		avgGC := float64(totalGC) / float64(rounds)

		t.Logf("║ %8d ║ %11.0f/s ║ %8.1f μs ║ %6d ║ %6.0f MB ║ %14.0f ║",
			concurrency, medTP, medLat, totalErrs, avgHeap, avgGC)

		if totalErrs > 0 {
			t.Errorf("  ⚠ %d errors at concurrency=%d", totalErrs, concurrency)
		}
	}
	t.Logf("╚══════════╩═══════════════╩════════════╩════════╩══════════╩════════════════╝")

	// Native baseline (just throughput, no latency comparison)
	nativeStats := make([][]stats, 4)
	nativeLevels := []int{1, 100, 1000, 10000}
	for i := range nativeStats {
		nativeStats[i] = make([]stats, rounds)
	}

	for round := 0; round < rounds; round++ {
		for li, concurrency := range nativeLevels {
			ops, _, heap, gc := runStressLevel(concurrency, duration, func(gID, i int) error {
				_ = nativeRule(gID*10000+i, " bob ", float64(30+i%70))
				return nil
			})
			tp := float64(ops) / duration.Seconds()
			nativeStats[li][round] = stats{
				throughput: tp,
				heapMB:     heap,
				gcPauses:   gc,
			}
		}
	}

	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║                    NATIVE GO BASELINE (throughput only)                     ║")
	t.Logf("╠══════════╦═══════════════╦══════════╦════════════════╗")
	t.Logf("║ Goroutin ║  Throughput   ║ Heap(MB) ║   GC Pauses    ║")
	t.Logf("╠══════════╬═══════════════╬══════════╬════════════════╣")

	for li, concurrency := range nativeLevels {
		var tps []float64
		var totalHeap float64
		var totalGC uint32
		for round := 0; round < rounds; round++ {
			s := nativeStats[li][round]
			tps = append(tps, s.throughput)
			totalHeap += s.heapMB
			totalGC += s.gcPauses
		}
		medTP := median3(tps[0], tps[1], tps[2])
		avgHeap := totalHeap / float64(rounds)
		avgGC := float64(totalGC) / float64(rounds)

		t.Logf("║ %8d ║ %11.0f/s ║ %6.0f MB ║ %14.0f ║",
			concurrency, medTP, avgHeap, avgGC)
	}
	t.Logf("╚══════════╩═══════════════╩══════════╩════════════════╝")

	// Ratio summary
	t.Logf("")
	t.Logf("Throughput ratio (Native / Gig):")
	for li, concurrency := range nativeLevels {
		var gigTP float64
		for gli, gc := range levels {
			if gc == concurrency {
				var tps []float64
				for round := 0; round < rounds; round++ {
					tps = append(tps, gigStats[gli][round].throughput)
				}
				gigTP = median3(tps[0], tps[1], tps[2])
				break
			}
		}
		if gigTP > 0 {
			var ntps []float64
			for round := 0; round < rounds; round++ {
				ntps = append(ntps, nativeStats[li][round].throughput)
			}
			nativeTP := median3(ntps[0], ntps[1], ntps[2])
			t.Logf("  %5dG: Native %.0f/s vs Gig %.0f/s = %.1fx", concurrency, nativeTP, gigTP, nativeTP/gigTP)
		}
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
	prog, err := gig.Build(statefulRuleSource, gig.WithStatefulGlobals(), gig.WithAllowPanic())
	if err != nil {
		t.Fatal(err)
	}
	defer prog.Close()

	// Verify correctness
	r, err := prog.Run("EvaluateRuleStateful", 1, " alice ", 81.0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Correctness: %v", r)

	levels := []int{1, 100, 1000, 5000, 10000}
	duration := 3 * time.Second

	type stats struct {
		throughput float64
		avgLatUs   float64
		errors     int64
		heapMB     float64
		gcPauses   uint32
	}

	t.Logf("Running stateful stress test (%d levels, %v/level)...", len(levels), duration)

	gigStats := make([]stats, len(levels))
	for li, concurrency := range levels {
		// Rebuild for each level to reset counter
		prog.Close()
		prog, err = gig.Build(statefulRuleSource, gig.WithStatefulGlobals(), gig.WithAllowPanic())
		if err != nil {
			t.Fatal(err)
		}

		ops, errs, heap, gc := runStressLevel(concurrency, duration, func(gID, i int) error {
			execCtx, execCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer execCancel()
			_, err := prog.RunWithContext(execCtx, "EvaluateRuleStateful", gID*10000+i, " bob ", float64(30+i%70))
			return err
		})
		tp := float64(ops) / duration.Seconds()
		gigStats[li] = stats{
			throughput: tp,
			avgLatUs:   1e6 / tp * float64(concurrency),
			errors:     errs,
			heapMB:     heap,
			gcPauses:   gc,
		}
	}

	// Print results
	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║           STATEFUL STRESS TEST (WithStatefulGlobals)                       ║")
	t.Logf("║  Workload: rule engine + global counter + *sync.Mutex  |  CPU: %2d cores  ║", runtime.NumCPU())
	t.Logf("╠══════════╦═══════════════╦════════════╦════════╦══════════╦════════════════╣")
	t.Logf("║ Goroutin ║  Throughput   ║  Avg Lat   ║ Errors ║ Heap(MB) ║   GC Pauses    ║")
	t.Logf("╠══════════╬═══════════════╬════════════╬════════╬══════════╬════════════════╣")

	for li, concurrency := range levels {
		s := gigStats[li]
		t.Logf("║ %8d ║ %11.0f/s ║ %8.1f μs ║ %6d ║ %6.0f MB ║ %14.0f ║",
			concurrency, s.throughput, s.avgLatUs, s.errors, s.heapMB, float64(s.gcPauses))
		if s.errors > 0 {
			t.Errorf("  ⚠ %d errors at concurrency=%d", s.errors, concurrency)
		}
	}
	t.Logf("╚══════════╩═══════════════╩════════════╩════════╩══════════╩════════════════╝")

	// Verify counter correctness with a controlled test.
	prog.Close()
	prog, err = gig.Build(statefulRuleSource, gig.WithStatefulGlobals(), gig.WithAllowPanic())
	if err != nil {
		t.Fatal(err)
	}

	const verifyConcurrency = 100
	const verifyOpsPerG = 100
	var wg sync.WaitGroup
	for g := 0; g < verifyConcurrency; g++ {
		wg.Add(1)
		go func(gID int) {
			defer wg.Done()
			for i := 0; i < verifyOpsPerG; i++ {
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
	expected := verifyConcurrency * verifyOpsPerG
	if counterVal != expected {
		t.Errorf("Counter = %d, want %d (lost %d updates)", counterVal, expected, expected-counterVal)
	} else {
		t.Logf("Counter verification: %d ops = %d (correct, no lost updates)", counterVal, expected)
	}
	prog.Close()
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

// TestConcurrentGlobals_GoNative_vs_Gig compares Go native and Gig
// behavior when concurrently accessing global variables.
// Uses moderate concurrency to keep test time reasonable.
func TestConcurrentGlobals_GoNative_vs_Gig(t *testing.T) {
	const concurrency = 50
	const opsPerGoroutine = 100
	const totalOps = concurrency * opsPerGoroutine

	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║    CONCURRENT GLOBALS: Go Native vs Gig (%dG × %d ops)               ║", concurrency, opsPerGoroutine)
	t.Logf("╚══════════════════════════════════════════════════════════════════════════════╝")

	// --- Go Native: unprotected global ---
	var nativeCounterUnprotected int
	var nativeWg sync.WaitGroup
	nativeStart := time.Now()
	for g := 0; g < concurrency; g++ {
		nativeWg.Add(1)
		go func() {
			defer nativeWg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				nativeCounterUnprotected++
			}
		}()
	}
	nativeWg.Wait()
	nativeElapsedUnprotected := time.Since(nativeStart)

	t.Logf("")
	t.Logf("┌─────────────────────────────────────────────────────┐")
	t.Logf("│ Go Native: Unprotected Global                      │")
	t.Logf("├─────────────────────────────────────────────────────┤")
	t.Logf("│ Total ops:      %10d                        │", totalOps)
	t.Logf("│ Final counter:  %10d (lost %d updates)       │", nativeCounterUnprotected, totalOps-nativeCounterUnprotected)
	t.Logf("│ Has data race:  %10v                        │", nativeCounterUnprotected != totalOps)
	t.Logf("│ Elapsed:        %10v                        │", nativeElapsedUnprotected.Round(time.Microsecond))
	t.Logf("└─────────────────────────────────────────────────────┘")

	// --- Go Native: mutex-protected global ---
	var nativeCounterMutex int
	var nativeMu sync.Mutex
	nativeStart = time.Now()
	for g := 0; g < concurrency; g++ {
		nativeWg.Add(1)
		go func() {
			defer nativeWg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				nativeMu.Lock()
				nativeCounterMutex++
				nativeMu.Unlock()
			}
		}()
	}
	nativeWg.Wait()
	nativeElapsedMutex := time.Since(nativeStart)

	t.Logf("")
	t.Logf("┌─────────────────────────────────────────────────────┐")
	t.Logf("│ Go Native: *sync.Mutex Protected Global             │")
	t.Logf("├─────────────────────────────────────────────────────┤")
	t.Logf("│ Total ops:      %10d                        │", totalOps)
	t.Logf("│ Final counter:  %10d                        │", nativeCounterMutex)
	t.Logf("│ Is precise:     %10v                        │", nativeCounterMutex == totalOps)
	t.Logf("│ Elapsed:        %10v                        │", nativeElapsedMutex.Round(time.Microsecond))
	t.Logf("└─────────────────────────────────────────────────────┘")

	// --- Go Native: atomic-protected global ---
	var nativeCounterAtomic int64
	nativeStart = time.Now()
	for g := 0; g < concurrency; g++ {
		nativeWg.Add(1)
		go func() {
			defer nativeWg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				atomic.AddInt64(&nativeCounterAtomic, 1)
			}
		}()
	}
	nativeWg.Wait()
	nativeElapsedAtomic := time.Since(nativeStart)

	t.Logf("")
	t.Logf("┌─────────────────────────────────────────────────────┐")
	t.Logf("│ Go Native: sync/atomic Protected Global             │")
	t.Logf("├─────────────────────────────────────────────────────┤")
	t.Logf("│ Total ops:      %10d                        │", totalOps)
	t.Logf("│ Final counter:  %10d                        │", nativeCounterAtomic)
	t.Logf("│ Is precise:     %10v                        │", nativeCounterAtomic == int64(totalOps))
	t.Logf("│ Elapsed:        %10v                        │", nativeElapsedAtomic.Round(time.Microsecond))
	t.Logf("└─────────────────────────────────────────────────────┘")

	// --- Gig: unprotected global (stateful mode, SharedGlobals provides RWMutex) ---
	progUnprotected, err := gig.Build(unprotectedGlobalSource, gig.WithStatefulGlobals(), gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build unprotected source: %v", err)
	}
	defer progUnprotected.Close()

	gigStart := time.Now()
	var gigWg sync.WaitGroup
	for g := 0; g < concurrency; g++ {
		gigWg.Add(1)
		go func() {
			defer gigWg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				progUnprotected.Run("Increment")
			}
		}()
	}
	gigWg.Wait()
	gigElapsedUnprotected := time.Since(gigStart)

	gigResultUnprotected, _ := progUnprotected.Run("GetCounter")
	gigCounterUnprotected, _ := gigResultUnprotected.(int)

	t.Logf("")
	t.Logf("┌─────────────────────────────────────────────────────┐")
	t.Logf("│ Gig (Stateful): Unprotected Global                 │")
	t.Logf("├─────────────────────────────────────────────────────┤")
	t.Logf("│ Total ops:      %10d                        │", totalOps)
	t.Logf("│ Final counter:  %10d (lost %d updates)       │", gigCounterUnprotected, totalOps-gigCounterUnprotected)
	t.Logf("│ Has data race:  %10v                        │", gigCounterUnprotected != totalOps)
	t.Logf("│ Elapsed:        %10v                        │", gigElapsedUnprotected.Round(time.Millisecond))
	t.Logf("│ Note: SharedGlobals RWMutex prevents torn reads,   │")
	t.Logf("│       but read-modify-write still has lost updates  │")
	t.Logf("└─────────────────────────────────────────────────────┘")

	// --- Gig: mutex-protected global (stateful mode) ---
	progMutex, err := gig.Build(mutexProtectedGlobalSource, gig.WithStatefulGlobals(), gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build mutex source: %v", err)
	}
	defer progMutex.Close()

	gigStart = time.Now()
	for g := 0; g < concurrency; g++ {
		gigWg.Add(1)
		go func() {
			defer gigWg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				progMutex.Run("Increment")
			}
		}()
	}
	gigWg.Wait()
	gigElapsedMutex := time.Since(gigStart)

	gigResultMutex, _ := progMutex.Run("GetCounter")
	gigCounterMutex, _ := gigResultMutex.(int)

	t.Logf("")
	t.Logf("┌─────────────────────────────────────────────────────┐")
	t.Logf("│ Gig (Stateful): *sync.Mutex Protected Global        │")
	t.Logf("├─────────────────────────────────────────────────────┤")
	t.Logf("│ Total ops:      %10d                        │", totalOps)
	t.Logf("│ Final counter:  %10d                        │", gigCounterMutex)
	t.Logf("│ Is precise:     %10v                        │", gigCounterMutex == totalOps)
	t.Logf("│ Elapsed:        %10v                        │", gigElapsedMutex.Round(time.Millisecond))
	t.Logf("└─────────────────────────────────────────────────────┘")

	// --- Summary table ---
	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║                         COMPARISON SUMMARY                                 ║")
	t.Logf("╠══════════════════════════════════════╦════════════╦═══════════╦════════════╣")
	t.Logf("║ Mode                                 ║ Counter    ║ Precise?  ║ Has Race?  ║")
	t.Logf("╠══════════════════════════════════════╬════════════╬═══════════╬════════════╣")
	t.Logf("║ Go Native: Unprotected               ║ %10d ║ %9v ║ %10v ║",
		nativeCounterUnprotected, nativeCounterUnprotected == totalOps, nativeCounterUnprotected != totalOps)
	t.Logf("║ Go Native: *sync.Mutex               ║ %10d ║ %9v ║ %10v ║",
		nativeCounterMutex, nativeCounterMutex == totalOps, nativeCounterMutex != totalOps)
	t.Logf("║ Go Native: sync/atomic               ║ %10d ║ %9v ║ %10v ║",
		nativeCounterAtomic, nativeCounterAtomic == int64(totalOps), nativeCounterAtomic != int64(totalOps))
	t.Logf("║ Gig Stateful: Unprotected            ║ %10d ║ %9v ║ %10v ║",
		gigCounterUnprotected, gigCounterUnprotected == totalOps, gigCounterUnprotected != totalOps)
	t.Logf("║ Gig Stateful: *sync.Mutex            ║ %10d ║ %9v ║ %10v ║",
		gigCounterMutex, gigCounterMutex == totalOps, gigCounterMutex != totalOps)
	t.Logf("╚══════════════════════════════════════╩════════════╩═══════════╩════════════╝")

	// Verify the key semantic equivalence: mutex-protected Gig == mutex-protected Go
	if gigCounterMutex == totalOps && nativeCounterMutex == totalOps {
		t.Logf("✓ Mutex-protected globals: Gig and Go native both produce precise results")
	}
}
