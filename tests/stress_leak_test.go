// Package tests - stress_leak_test.go
//
// Long-running stress benchmarks to detect memory leaks in concurrent gig programs.
// Run with: go test ./tests/ -bench BenchmarkStress_MemoryLeak -benchtime 30m -timeout 6h -v
//
// These are benchmarks (not tests) so they don't run during `go test ./...`.
package tests

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

// ruleEngineSource simulates a realistic rule engine workload:
// parse input, string manipulation, math, conditional logic, stdlib calls.
const ruleEngineSource = `
import (
	"strings"
	"strconv"
	"math"
)

func EvaluateRule(userID int, name string, score float64) string {
	// String processing
	upper := strings.ToUpper(name)
	trimmed := strings.TrimSpace(upper)
	idStr := strconv.Itoa(userID)

	// Math computation
	adjusted := math.Sqrt(score) * 10.0

	// Business logic
	tier := "BRONZE"
	if adjusted > 80.0 {
		tier = "GOLD"
	} else if adjusted > 50.0 {
		tier = "SILVER"
	}

	// String building
	result := trimmed + "#" + idStr + ":" + tier
	if strings.Contains(result, "GOLD") {
		result = result + " [VIP]"
	}
	return result
}
`

// runStressMemoryLeak is the core stress test logic, used by all benchmark variants.
// duration controls how long the test runs.
func runStressMemoryLeak(b *testing.B, duration time.Duration) {
	b.Helper()

	b.Logf("Starting %v memory leak stress benchmark", duration)

	prog := buildStressMemoryProgram(b)
	defer prog.Close()

	baseline := captureStressMemoryBaseline()
	logStressMemoryBaseline(b, baseline)

	logFile := createStressMemoryLog(b)
	defer logFile.Close()

	config := stressMemoryConfig{
		concurrency:    20, // Reduced from 50 to lower mutex contention.
		reportInterval: time.Minute,
		duration:       duration,
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	counters := &stressMemoryCounters{}
	var wg sync.WaitGroup
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	stopCh := make(chan struct{})

	startStressMemoryWorkers(prog, config, counters, workerCtx, stopCh, &wg)
	monitorDone := startStressMemoryMonitor(b, ctx, logFile, baseline, config, counters)
	<-ctx.Done()

	stopStressMemoryWorkers(b, stopCh, workerCancel, &wg)
	<-monitorDone

	final := captureStressMemoryFinal()
	logStressMemoryFinalReport(b, duration, baseline, final, counters)
	checkStressMemoryFinalThresholds(b, baseline, final)
}

type stressMemoryConfig struct {
	concurrency    int
	reportInterval time.Duration
	duration       time.Duration
}

type stressMemoryBaseline struct {
	heapAlloc  uint64
	heapSys    uint64
	goroutines int
}

type stressMemoryCounters struct {
	totalOps   atomic.Int64
	ongoingOps atomic.Int64
	errors     atomic.Int64
}

type stressMemorySnapshot struct {
	heapAllocMB       float64
	heapSysMB         float64
	heapInUseMB       float64
	heapIdleMB        float64
	stackInUseMB      float64
	heapGrowthMB      float64
	heapGrowthPercent float64
	numGC             uint32
	goroutines        int
	totalOps          int64
	ongoingOps        int64
	errors            int64
}

type stressMemoryFinal struct {
	memStats   runtime.MemStats
	goroutines int
}

func buildStressMemoryProgram(b *testing.B) *gig.Program {
	b.Helper()
	prog, err := gig.Build(ruleEngineSource)
	if err != nil {
		b.Fatalf("Build error: %v", err)
	}
	return prog
}

func captureStressMemoryBaseline() stressMemoryBaseline {
	var memStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStats)
	return stressMemoryBaseline{
		heapAlloc:  memStats.HeapAlloc,
		heapSys:    memStats.HeapSys,
		goroutines: runtime.NumGoroutine(),
	}
}

func logStressMemoryBaseline(b *testing.B, baseline stressMemoryBaseline) {
	b.Helper()
	b.Logf("══════════════════════════════════════════════════════════════")
	b.Logf("BASELINE MEASUREMENTS")
	b.Logf("══════════════════════════════════════════════════════════════")
	b.Logf("  HeapAlloc:    %d MB", baseline.heapAlloc/1024/1024)
	b.Logf("  HeapSys:      %d MB", baseline.heapSys/1024/1024)
	b.Logf("  Goroutines:   %d", baseline.goroutines)
	b.Logf("══════════════════════════════════════════════════════════════")
}

func createStressMemoryLog(b *testing.B) *os.File {
	b.Helper()
	logFile, err := os.Create("stress_memory_leak.log")
	if err != nil {
		b.Fatalf("Failed to create log file: %v", err)
	}
	fmt.Fprintf(logFile, "Timestamp,HeapAllocMB,HeapSysMB,HeapInUseMB,HeapIdleMB,StackInUseMB,NumGC,Goroutines,TotalOps,OngoingOps,Errors\n")
	return logFile
}

func startStressMemoryWorkers(
	prog *gig.Program,
	config stressMemoryConfig,
	counters *stressMemoryCounters,
	workerCtx context.Context,
	stopCh <-chan struct{},
	wg *sync.WaitGroup,
) {
	for i := 0; i < config.concurrency; i++ {
		wg.Add(1)
		go runStressMemoryWorker(prog, counters, workerCtx, stopCh, wg, i)
	}
}

func runStressMemoryWorker(
	prog *gig.Program,
	counters *stressMemoryCounters,
	workerCtx context.Context,
	stopCh <-chan struct{},
	wg *sync.WaitGroup,
	workerID int,
) {
	defer wg.Done()
	opCounter := 0
	for {
		select {
		case <-workerCtx.Done():
			return
		case <-stopCh:
			return
		default:
			opCounter++
			runStressMemoryOperation(prog, counters, workerID, opCounter)
		}
	}
}

func runStressMemoryOperation(prog *gig.Program, counters *stressMemoryCounters, workerID, opCounter int) {
	counters.ongoingOps.Add(1)
	defer counters.ongoingOps.Add(-1)

	execCtx, execCancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := prog.RunWithContext(execCtx, "EvaluateRule", workerID*100000+opCounter, " test_user ", float64(25+opCounter%50))
	execCancel()

	if err != nil {
		counters.errors.Add(1)
	}
	counters.totalOps.Add(1)
}

func startStressMemoryMonitor(
	b *testing.B,
	ctx context.Context,
	logFile *os.File,
	baseline stressMemoryBaseline,
	config stressMemoryConfig,
	counters *stressMemoryCounters,
) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(config.reportInterval)
		defer ticker.Stop()

		startTime := time.Now()
		iteration := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				iteration++
				elapsed := time.Since(startTime)
				snapshot := captureStressMemorySnapshot(baseline, counters)
				writeStressMemorySnapshot(logFile, snapshot)
				logStressMemorySnapshot(b, iteration, elapsed, baseline, snapshot)
				checkStressMemoryMonitorThresholds(b, elapsed, baseline, config, snapshot)
			}
		}
	}()
	return done
}

func captureStressMemorySnapshot(baseline stressMemoryBaseline, counters *stressMemoryCounters) stressMemorySnapshot {
	var memStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStats)

	return stressMemorySnapshot{
		heapAllocMB:       float64(memStats.HeapAlloc) / 1024 / 1024,
		heapSysMB:         float64(memStats.HeapSys) / 1024 / 1024,
		heapInUseMB:       float64(memStats.HeapInuse) / 1024 / 1024,
		heapIdleMB:        float64(memStats.HeapIdle) / 1024 / 1024,
		stackInUseMB:      float64(memStats.StackInuse) / 1024 / 1024,
		heapGrowthMB:      float64(memStats.HeapAlloc-baseline.heapAlloc) / 1024 / 1024,
		heapGrowthPercent: float64(memStats.HeapAlloc-baseline.heapAlloc) / float64(baseline.heapAlloc) * 100,
		numGC:             memStats.NumGC,
		goroutines:        runtime.NumGoroutine(),
		totalOps:          counters.totalOps.Load(),
		ongoingOps:        counters.ongoingOps.Load(),
		errors:            counters.errors.Load(),
	}
}

func writeStressMemorySnapshot(logFile *os.File, snapshot stressMemorySnapshot) {
	fmt.Fprintf(logFile, "%s,%.2f,%.2f,%.2f,%.2f,%.2f,%d,%d,%d,%d,%d\n",
		time.Now().Format(time.RFC3339),
		snapshot.heapAllocMB,
		snapshot.heapSysMB,
		snapshot.heapInUseMB,
		snapshot.heapIdleMB,
		snapshot.stackInUseMB,
		snapshot.numGC,
		snapshot.goroutines,
		snapshot.totalOps,
		snapshot.ongoingOps,
		snapshot.errors)
}

func logStressMemorySnapshot(
	b *testing.B,
	iteration int,
	elapsed time.Duration,
	baseline stressMemoryBaseline,
	snapshot stressMemorySnapshot,
) {
	b.Helper()
	b.Logf("══════════════════════════════════════════════════════════════")
	b.Logf("ITERATION %d | Elapsed: %v", iteration, elapsed.Round(time.Second))
	b.Logf("══════════════════════════════════════════════════════════════")
	b.Logf("  HeapAlloc:    %.2f MB (growth: %.2f MB, %.1f%%)", snapshot.heapAllocMB, snapshot.heapGrowthMB, snapshot.heapGrowthPercent)
	b.Logf("  HeapSys:      %.2f MB", snapshot.heapSysMB)
	b.Logf("  HeapInUse:    %.2f MB", snapshot.heapInUseMB)
	b.Logf("  HeapIdle:     %.2f MB", snapshot.heapIdleMB)
	b.Logf("  StackInUse:   %.2f MB", snapshot.stackInUseMB)
	b.Logf("  NumGC:        %d", snapshot.numGC)
	b.Logf("  Goroutines:   %d (baseline: %d)", snapshot.goroutines, baseline.goroutines)
	b.Logf("  TotalOps:     %d", snapshot.totalOps)
	b.Logf("  OngoingOps:   %d", snapshot.ongoingOps)
	b.Logf("  Errors:       %d", snapshot.errors)
	b.Logf("══════════════════════════════════════════════════════════════")
}

func checkStressMemoryMonitorThresholds(
	b *testing.B,
	elapsed time.Duration,
	baseline stressMemoryBaseline,
	config stressMemoryConfig,
	snapshot stressMemorySnapshot,
) {
	b.Helper()
	if elapsed > 30*time.Minute && snapshot.heapGrowthMB > 100 {
		b.Logf("⚠️  WARNING: Potential memory leak detected!")
		b.Logf("   Heap has grown %.2f MB in %v", snapshot.heapGrowthMB, elapsed)
	}
	if snapshot.heapGrowthMB > 200 {
		b.Errorf("🚨 CRITICAL: Memory leak detected! Heap grew %.2f MB", snapshot.heapGrowthMB)
	}
	expectedGoroutines := baseline.goroutines + config.concurrency + 3
	if snapshot.goroutines > expectedGoroutines+10 {
		b.Errorf("🚨 CRITICAL: Goroutine leak detected! %d goroutines (expected max %d)", snapshot.goroutines, expectedGoroutines)
	}
}

func stopStressMemoryWorkers(
	b *testing.B,
	stopCh chan<- struct{},
	workerCancel context.CancelFunc,
	wg *sync.WaitGroup,
) {
	b.Helper()
	b.Logf("Stopping workers...")
	close(stopCh)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		b.Logf("All workers stopped gracefully")
	case <-time.After(30 * time.Second):
		b.Logf("Timeout waiting for workers, forcing shutdown")
		workerCancel()
		wg.Wait()
	}
}

func captureStressMemoryFinal() stressMemoryFinal {
	var memStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStats)
	return stressMemoryFinal{
		memStats:   memStats,
		goroutines: runtime.NumGoroutine(),
	}
}

func logStressMemoryFinalReport(
	b *testing.B,
	duration time.Duration,
	baseline stressMemoryBaseline,
	final stressMemoryFinal,
	counters *stressMemoryCounters,
) {
	b.Helper()
	b.Logf("\n")
	b.Logf("══════════════════════════════════════════════════════════════")
	b.Logf("FINAL REPORT")
	b.Logf("══════════════════════════════════════════════════════════════")
	b.Logf("  Duration:          %v", duration)
	b.Logf("  Total Operations:  %d", counters.totalOps.Load())
	b.Logf("  Errors:            %d", counters.errors.Load())
	b.Logf("  Ops/Second:        %.0f", float64(counters.totalOps.Load())/duration.Seconds())
	b.Logf("────────────────────────────────────────────────────────────────")
	b.Logf("  Baseline HeapAlloc:  %d MB", baseline.heapAlloc/1024/1024)
	b.Logf("  Final HeapAlloc:     %d MB", final.memStats.HeapAlloc/1024/1024)
	b.Logf("  Heap Growth:         %.2f MB", stressHeapGrowthMB(baseline, final))
	b.Logf("────────────────────────────────────────────────────────────────")
	b.Logf("  Baseline Goroutines: %d", baseline.goroutines)
	b.Logf("  Final Goroutines:    %d", final.goroutines)
	b.Logf("  Goroutine Growth:    %d", stressGoroutineGrowth(baseline, final))
	b.Logf("────────────────────────────────────────────────────────────────")
	b.Logf("  NumGC:               %d", final.memStats.NumGC)
	b.Logf("  GCSys:               %d MB", final.memStats.GCSys/1024/1024)
	b.Logf("══════════════════════════════════════════════════════════════")
	b.Logf("Log file saved to: stress_memory_leak.log")
}

func checkStressMemoryFinalThresholds(b *testing.B, baseline stressMemoryBaseline, final stressMemoryFinal) {
	b.Helper()
	heapGrowthMB := stressHeapGrowthMB(baseline, final)
	goroutineGrowth := stressGoroutineGrowth(baseline, final)

	if heapGrowthMB > 50 {
		b.Errorf("FAIL: Heap growth %.2f MB exceeds 50 MB threshold", heapGrowthMB)
	}

	if goroutineGrowth > 5 {
		b.Errorf("FAIL: Goroutine growth %d exceeds 5 threshold", goroutineGrowth)
	}

	if heapGrowthMB <= 50 && goroutineGrowth <= 5 {
		b.Logf("✅ PASS: No memory leaks detected!")
	}
}

func stressHeapGrowthMB(baseline stressMemoryBaseline, final stressMemoryFinal) float64 {
	return float64(final.memStats.HeapAlloc-baseline.heapAlloc) / 1024 / 1024
}

func stressGoroutineGrowth(baseline stressMemoryBaseline, final stressMemoryFinal) int {
	return final.goroutines - baseline.goroutines
}

// ============================================================================
// Benchmark variants (only run with -bench flag)
// ============================================================================

// BenchmarkStress_MemoryLeak runs a 5-hour memory leak stress test.
// Usage: go test ./tests/ -bench BenchmarkStress_MemoryLeak$ -benchtime 1x -timeout 6h -v
func BenchmarkStress_MemoryLeak(b *testing.B) {
	durationStr := os.Getenv("STRESS_DURATION")
	duration := 5 * time.Hour
	if durationStr != "" {
		if d, err := time.ParseDuration(durationStr); err == nil {
			duration = d
		}
	}
	for i := 0; i < b.N; i++ {
		runStressMemoryLeak(b, duration)
	}
}

// BenchmarkStress_MemoryLeak_Short runs a 30-minute memory leak stress test.
// Usage: go test ./tests/ -bench BenchmarkStress_MemoryLeak_Short -benchtime 1x -timeout 1h -v
func BenchmarkStress_MemoryLeak_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runStressMemoryLeak(b, 30*time.Minute)
	}
}

// BenchmarkStress_MemoryLeak_Quick runs a 5-minute memory leak stress test.
// Usage: go test ./tests/ -bench BenchmarkStress_MemoryLeak_Quick -benchtime 1x -timeout 10m -v
func BenchmarkStress_MemoryLeak_Quick(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runStressMemoryLeak(b, 5*time.Minute)
	}
}
