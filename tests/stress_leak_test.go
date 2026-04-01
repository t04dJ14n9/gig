// Package tests - stress_leak_test.go
//
// Long-running stress test to detect memory leaks in concurrent gig programs.
// Run with: go test ./tests/ -run TestStress_MemoryLeak -v -timeout 6h
//
// This test runs for 5 hours with concurrent goroutines, monitoring memory
// usage every minute and reporting any leaks.
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

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
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

// ============================================================================
// Memory Leak Detection Test (5 hours)
// ============================================================================

func TestStress_MemoryLeak(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running stress test in short mode")
	}

	// Check for environment variable to control duration
	durationStr := os.Getenv("STRESS_DURATION")
	duration := 5 * time.Hour
	if durationStr != "" {
		if d, err := time.ParseDuration(durationStr); err == nil {
			duration = d
		}
	}

	t.Logf("Starting %v memory leak stress test", duration)
	t.Logf("Set STRESS_DURATION env var to customize (e.g., STRESS_DURATION=30m go test -run TestStress_MemoryLeak)")

	// Build program once
	prog, err := gig.Build(ruleEngineSource)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	defer prog.Close()

	// Memory tracking
	var memStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStats)

	// Baseline measurements
	baselineHeapAlloc := memStats.HeapAlloc
	baselineHeapSys := memStats.HeapSys
	baselineGoroutines := runtime.NumGoroutine()

	t.Logf("══════════════════════════════════════════════════════════════")
	t.Logf("BASELINE MEASUREMENTS")
	t.Logf("══════════════════════════════════════════════════════════════")
	t.Logf("  HeapAlloc:    %d MB", baselineHeapAlloc/1024/1024)
	t.Logf("  HeapSys:      %d MB", baselineHeapSys/1024/1024)
	t.Logf("  Goroutines:   %d", baselineGoroutines)
	t.Logf("══════════════════════════════════════════════════════════════")

	// Create log file for memory tracking
	logFile, err := os.Create("stress_memory_leak.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	// Write header to log file
	fmt.Fprintf(logFile, "Timestamp,HeapAllocMB,HeapSysMB,HeapInUseMB,HeapIdleMB,StackInUseMB,NumGC,Goroutines,TotalOps,OngoingOps,Errors\n")

	// Test configuration
	concurrency := 20 // Reduced from 50 to lower mutex contention
	reportInterval := 1 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Metrics
	var totalOps atomic.Int64
	var ongoingOps atomic.Int64
	var errors atomic.Int64

	// Start worker goroutines
	var wg sync.WaitGroup
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	// Channel for graceful shutdown
	stopCh := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			opCounter := 0
			for {
				select {
				case <-workerCtx.Done():
					return
				case <-stopCh:
					return
				default:
					ongoingOps.Add(1)
					opCounter++

					// Execute with timeout to prevent stuck operations
					execCtx, execCancel := context.WithTimeout(context.Background(), 5*time.Second)
					_, err := prog.RunWithContext(execCtx, "EvaluateRule", workerID*100000+opCounter, " test_user ", float64(25+opCounter%50))
					execCancel()

					if err != nil {
						errors.Add(1)
					}
					totalOps.Add(1)
					ongoingOps.Add(-1)
				}
			}
		}(i)
	}

	// Monitor goroutine - runs every minute
	monitorDone := make(chan struct{})
	go func() {
		defer close(monitorDone)
		ticker := time.NewTicker(reportInterval)
		defer ticker.Stop()

		startTime := time.Now()
		iteration := 0

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				iteration++
				runtime.GC()
				runtime.ReadMemStats(&memStats)

				elapsed := time.Since(startTime)
				heapAllocMB := float64(memStats.HeapAlloc) / 1024 / 1024
				heapSysMB := float64(memStats.HeapSys) / 1024 / 1024
				heapInUseMB := float64(memStats.HeapInuse) / 1024 / 1024
				heapIdleMB := float64(memStats.HeapIdle) / 1024 / 1024
				stackInUseMB := float64(memStats.StackInuse) / 1024 / 1024
				numGC := memStats.NumGC
				goroutines := runtime.NumGoroutine()

				ops := totalOps.Load()
				errCount := errors.Load()
				ongoing := ongoingOps.Load()

				// Calculate growth rate
				heapGrowthMB := float64(memStats.HeapAlloc-baselineHeapAlloc) / 1024 / 1024
				heapGrowthPercent := float64(memStats.HeapAlloc-baselineHeapAlloc) / float64(baselineHeapAlloc) * 100

				// Log to file
				fmt.Fprintf(logFile, "%s,%.2f,%.2f,%.2f,%.2f,%.2f,%d,%d,%d,%d,%d\n",
					time.Now().Format(time.RFC3339),
					heapAllocMB, heapSysMB, heapInUseMB, heapIdleMB,
					stackInUseMB, numGC, goroutines, ops, ongoing, errCount)

				// Log to test output
				t.Logf("══════════════════════════════════════════════════════════════")
				t.Logf("ITERATION %d | Elapsed: %v", iteration, elapsed.Round(time.Second))
				t.Logf("══════════════════════════════════════════════════════════════")
				t.Logf("  HeapAlloc:    %.2f MB (growth: %.2f MB, %.1f%%)", heapAllocMB, heapGrowthMB, heapGrowthPercent)
				t.Logf("  HeapSys:      %.2f MB", heapSysMB)
				t.Logf("  HeapInUse:    %.2f MB", heapInUseMB)
				t.Logf("  HeapIdle:     %.2f MB", heapIdleMB)
				t.Logf("  StackInUse:   %.2f MB", stackInUseMB)
				t.Logf("  NumGC:        %d", numGC)
				t.Logf("  Goroutines:   %d (baseline: %d)", goroutines, baselineGoroutines)
				t.Logf("  TotalOps:     %d", ops)
				t.Logf("  OngoingOps:   %d", ongoing)
				t.Logf("  Errors:       %d", errCount)
				t.Logf("══════════════════════════════════════════════════════════════")

				// Leak detection thresholds
				// After 30 minutes, if heap has grown more than 100MB, warn
				if elapsed > 30*time.Minute && heapGrowthMB > 100 {
					t.Logf("⚠️  WARNING: Potential memory leak detected!")
					t.Logf("   Heap has grown %.2f MB in %v", heapGrowthMB, elapsed)
				}

				// If heap growth exceeds 200MB, fail the test
				if heapGrowthMB > 200 {
					t.Errorf("🚨 CRITICAL: Memory leak detected! Heap grew %.2f MB", heapGrowthMB)
				}

				// Check for goroutine leak
				// Expected goroutines: baseline + workers + monitor (up to 3)
				expectedGoroutines := baselineGoroutines + concurrency + 3
				// Only warn if goroutines exceed expected count significantly
				if goroutines > expectedGoroutines + 10 {
					t.Errorf("🚨 CRITICAL: Goroutine leak detected! %d goroutines (expected max %d)", goroutines, expectedGoroutines)
				}
			}
		}
	}()

	// Wait for test duration
	<-ctx.Done()

	// Graceful shutdown
	t.Logf("Stopping workers...")
	close(stopCh) // Signal workers to stop
	
	// Give workers time to finish current operations
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		t.Logf("All workers stopped gracefully")
	case <-time.After(30 * time.Second):
		t.Logf("Timeout waiting for workers, forcing shutdown")
		workerCancel()
		wg.Wait()
	}
	
	<-monitorDone

	// Final measurement
	runtime.GC()
	runtime.ReadMemStats(&memStats)

	t.Logf("\n")
	t.Logf("══════════════════════════════════════════════════════════════")
	t.Logf("FINAL REPORT")
	t.Logf("══════════════════════════════════════════════════════════════")
	t.Logf("  Duration:          %v", duration)
	t.Logf("  Total Operations:  %d", totalOps.Load())
	t.Logf("  Errors:            %d", errors.Load())
	t.Logf("  Ops/Second:        %.0f", float64(totalOps.Load())/duration.Seconds())
	t.Logf("────────────────────────────────────────────────────────────────")
	t.Logf("  Baseline HeapAlloc:  %d MB", baselineHeapAlloc/1024/1024)
	t.Logf("  Final HeapAlloc:     %d MB", memStats.HeapAlloc/1024/1024)
	t.Logf("  Heap Growth:         %.2f MB", float64(memStats.HeapAlloc-baselineHeapAlloc)/1024/1024)
	t.Logf("────────────────────────────────────────────────────────────────")
	t.Logf("  Baseline Goroutines: %d", baselineGoroutines)
	t.Logf("  Final Goroutines:    %d", runtime.NumGoroutine())
	t.Logf("  Goroutine Growth:    %d", runtime.NumGoroutine()-baselineGoroutines)
	t.Logf("────────────────────────────────────────────────────────────────")
	t.Logf("  NumGC:               %d", memStats.NumGC)
	t.Logf("  GCSys:               %d MB", memStats.GCSys/1024/1024)
	t.Logf("══════════════════════════════════════════════════════════════")
	t.Logf("Log file saved to: stress_memory_leak.log")

	// Determine pass/fail
	heapGrowthMB := float64(memStats.HeapAlloc-baselineHeapAlloc) / 1024 / 1024
	goroutineGrowth := runtime.NumGoroutine() - baselineGoroutines

	if heapGrowthMB > 50 {
		t.Errorf("FAIL: Heap growth %.2f MB exceeds 50 MB threshold", heapGrowthMB)
	}

	if goroutineGrowth > 5 {
		t.Errorf("FAIL: Goroutine growth %d exceeds 5 threshold", goroutineGrowth)
	}

	if heapGrowthMB <= 50 && goroutineGrowth <= 5 {
		t.Logf("✅ PASS: No memory leaks detected!")
	}
}

// ============================================================================
// Shorter tests for CI/CD (30 minutes)
// ============================================================================

func TestStress_MemoryLeak_Short(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	// 30-minute test for CI/CD
	origEnv := os.Getenv("STRESS_DURATION")
	os.Setenv("STRESS_DURATION", "30m")
	defer func() {
		if origEnv == "" {
			os.Unsetenv("STRESS_DURATION")
		} else {
			os.Setenv("STRESS_DURATION", origEnv)
		}
	}()

	TestStress_MemoryLeak(t)
}

// ============================================================================
// Quick sanity check (5 minutes)
// ============================================================================

func TestStress_MemoryLeak_Quick(t *testing.T) {
	// 5-minute quick test
	origEnv := os.Getenv("STRESS_DURATION")
	os.Setenv("STRESS_DURATION", "5m")
	defer func() {
		if origEnv == "" {
			os.Unsetenv("STRESS_DURATION")
		} else {
			os.Setenv("STRESS_DURATION", origEnv)
		}
	}()

	TestStress_MemoryLeak(t)
}

// ============================================================================
// Helper: analyze memory log
// ============================================================================

// Run this after the test to analyze the log:
// go run cmd/analyze_memory_log.go stress_memory_leak.log
