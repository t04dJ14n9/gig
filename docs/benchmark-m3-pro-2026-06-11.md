# M3 Pro Benchmark Evidence - 2026-06-11

This file records the focused benchmark command used for the interview summary.
Use it as a reproducibility note, not a global claim that Gig is faster on every
workload.

## Environment

- Date: 2026-06-11
- Machine: Apple M3 Pro
- OS/arch: `darwin/arm64`
- Package: `gig-benchmarks`

## Command

```bash
cd benchmarks
go test -run '^$' -bench '^BenchmarkGig_(ConstFoldArithmetic|ArithSum|BubbleSort)$' \
  -benchmem -benchtime=1s -count=3
```

## Raw Output

```text
goos: darwin
goarch: arm64
pkg: gig-benchmarks
cpu: Apple M3 Pro
BenchmarkGig_ArithSum-12               	   36862	     31906 ns/op	     396 B/op	       6 allocs/op
BenchmarkGig_ArithSum-12               	   38508	     31288 ns/op	     395 B/op	       6 allocs/op
BenchmarkGig_ArithSum-12               	   37671	     31310 ns/op	     393 B/op	       6 allocs/op
BenchmarkGig_BubbleSort-12             	    3165	    399605 ns/op	    1316 B/op	       7 allocs/op
BenchmarkGig_BubbleSort-12             	    2877	    390403 ns/op	    1317 B/op	       7 allocs/op
BenchmarkGig_BubbleSort-12             	    3177	    388880 ns/op	    1316 B/op	       7 allocs/op
BenchmarkGig_ConstFoldArithmetic-12    	   17826	     68280 ns/op	     394 B/op	       6 allocs/op
BenchmarkGig_ConstFoldArithmetic-12    	   17907	     68729 ns/op	     394 B/op	       6 allocs/op
BenchmarkGig_ConstFoldArithmetic-12    	   17798	     67524 ns/op	     397 B/op	       6 allocs/op
PASS
ok  	gig-benchmarks	14.266s
```

## Range Summary

- `BenchmarkGig_ConstFoldArithmetic`: `67.5-68.7 us/op`, `6 allocs/op`
- `BenchmarkGig_ArithSum`: `31.3-31.9 us/op`, `6 allocs/op`
- `BenchmarkGig_BubbleSort`: `388.9-399.6 us/op`, `7 allocs/op`

## Interview Wording

Say:

> On my Apple M3 Pro, with `go test` and `count=3`, the targeted constant-fold
> benchmark runs around 68 microseconds after the conservative constant folding,
> branch folding, and narrow dead-code cleanup work.

Do not say:

> Gig is always faster than Yaegi.
