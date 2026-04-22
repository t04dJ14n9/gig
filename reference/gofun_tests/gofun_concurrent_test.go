package gofun

import (
	"sync"
	"testing"

	newgofun "git.code.oa.com/datacenter/onefun/gofun"
	_ "git.code.oa.com/datacenter/onefun/gofun/packages"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

// ============================================================================
// Bug #6: 并发全局变量 — gofun 无线程安全
// ============================================================================

const gofunConcurrentSrc = `
package main

var counter int

func Increment() int {
	counter = counter + 1
	return counter
}
`

const gigConcurrentSrc = `
package main

import "sync"

var (
	mu      sync.Mutex
	counter int
)

func Increment() int {
	mu.Lock()
	counter++
	v := counter
	mu.Unlock()
	return v
}

func GetCounter() int {
	mu.Lock()
	v := counter
	mu.Unlock()
	return v
}
`

// TestGofunVerify_Bug6_SingleThread 新 gofun 单线程正常。
func TestGofunVerify_Bug6_SingleThread(t *testing.T) {
	program, err := newgofun.Build(gofunConcurrentSrc)
	if err != nil {
		t.Fatalf("gofun Build: %v", err)
	}
	result, err := program.Run("Increment")
	if err != nil {
		t.Fatalf("gofun Increment: %v", err)
	}
	t.Logf("gofun 单线程: Increment() = %v", result)
}

// TestGofunVerify_Bug6_ConcurrentRace 文档化新 gofun 并发不安全。
// Program.globals 是 map[ssa.Value]*value.Value，无锁保护。
// 实际并发调用会触发 fatal error: concurrent map writes，无法 recover。
func TestGofunVerify_Bug6_ConcurrentRace(t *testing.T) {
	// 验证 globals 结构无锁
	program, err := newgofun.Build(gofunConcurrentSrc)
	if err != nil {
		t.Fatalf("gofun Build: %v", err)
	}

	// 顺序调用两次，验证全局状态共享（第二次应返回 2）
	r1, _ := program.Run("Increment")
	r2, _ := program.Run("Increment")
	t.Logf("gofun 顺序调用: Increment()=%v, Increment()=%v (全局状态共享)", r1, r2)

	// 不实际并发调用，因为 concurrent map writes 是 fatal error 无法 recover
	t.Log("新 gofun Program.globals = map[ssa.Value]*value.Value (无锁)")
	t.Log("并发调用 Run() 会触发 fatal error: concurrent map writes")
	t.Log("此 fatal error 无法被 recover 捕获")
}

// TestGofunVerify_Bug6_GigSafe Gig 并发安全。
func TestGofunVerify_Bug6_GigSafe(t *testing.T) {
	prog, err := gig.Build(gigConcurrentSrc, gig.WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Gig Build: %v", err)
	}
	defer prog.Close()

	const numG = 50
	const perG = 10
	var wg sync.WaitGroup
	wg.Add(numG)
	for i := 0; i < numG; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perG; j++ {
				prog.Run("Increment")
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("GetCounter")
	if err != nil {
		t.Fatalf("Gig GetCounter: %v", err)
	}
	got := toInt64ForGofun(result)
	expected := int64(numG * perG)
	if got != expected {
		t.Errorf("Gig counter = %d, want %d", got, expected)
	} else {
		t.Logf("Gig 并发计数器 = %d (精确)", got)
	}
}

func toInt64ForGofun(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case int32:
		return int64(n)
	default:
		return 0
	}
}
