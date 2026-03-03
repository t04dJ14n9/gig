module git.woa.com/youngjin/gig

go 1.23.1

require (
	git.code.oa.com/datacenter/faas/languages/golang/old/gofun v0.0.0
	github.com/peterh/liner v1.2.2
	golang.org/x/tools v0.30.0
)

require (
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-runewidth v0.0.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

// 本地 gofun 解释器引用（用于性能对比测试）
// 运行 gofun 测试时使用: go test -tags=gofun ./tests/gofun_benchmark_test.go
replace git.code.oa.com/datacenter/faas/languages/golang/old/gofun => ./reference/faas/languages/golang/old/gofun
