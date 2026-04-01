// Package tests - pool_bench_test.go
//
// Benchmark for VMPool performance comparison
package tests

import (
	"context"
	"sync"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
)

// Uses ruleEngineSource from stress_leak_test.go

func BenchmarkVMPoolConcurrent(b *testing.B) {
	prog, err := gig.Build(ruleEngineSource)
	if err != nil {
		b.Fatal(err)
	}
	defer prog.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ctx := context.Background()
			_, err := prog.RunWithContext(ctx, "EvaluateRule", i, " user ", float64(25+i%50))
			if err != nil {
				b.Errorf("Run error: %v", err)
			}
			i++
		}
	})
}

func BenchmarkVMPoolSerial(b *testing.B) {
	prog, err := gig.Build(ruleEngineSource)
	if err != nil {
		b.Fatal(err)
	}
	defer prog.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := prog.RunWithContext(ctx, "EvaluateRule", i, " user ", float64(25+i%50))
		if err != nil {
			b.Errorf("Run error: %v", err)
		}
	}
}

func BenchmarkVMPoolConcurrent10(b *testing.B) {
	prog, err := gig.Build(ruleEngineSource)
	if err != nil {
		b.Fatal(err)
	}
	defer prog.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					ctx := context.Background()
					_, err := prog.RunWithContext(ctx, "EvaluateRule", id*1000+j, " user ", float64(25+j%50))
					if err != nil {
						b.Errorf("Run error: %v", err)
					}
				}
			}(i)
		}
		wg.Wait()
	}
}
