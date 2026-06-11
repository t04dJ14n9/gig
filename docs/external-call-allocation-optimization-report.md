# External Call Argument Allocation Optimization

Date: 2026-06-11
Machine: Apple M3 Pro, `darwin/arm64`, Go 1.26.3

## Summary

`vm.callExternal` used to allocate a fresh `[]value.Value` for every external call:

```go
args := make([]value.Value, numArgs)
for i := numArgs - 1; i >= 0; i-- {
	args[i] = v.pop()
}
```

That made external-call-heavy code allocate heavily even when the external target used a DirectCall wrapper and avoided `reflect.Call`.

The fix now slices the VM operand stack directly:

```go
argStart := v.sp - numArgs
args := v.stack[argStart:v.sp]
v.sp = argStart
```

This is allocation-free for small arities and for larger arities too, because call arguments are already contiguous on the VM stack in source order.

## Why not the `[8]value.Value` buffer

I first tried the same small-buffer pattern used by `OpCallIndirect`:

```go
var argsBuf [8]value.Value
args := argsBuf[:numArgs]
```

The regression test still reported one allocation per call. The reason is that `args` is passed through a function field (`DirectCall func([]value.Value) value.Value`), so escape analysis cannot prove the callee will not retain the slice. The local array therefore still escaped.

Using a stack-backed slice avoids creating new backing storage at all.

## Regression Test

Added:

```text
TestCallExternalSmallArityDoesNotAllocateArgSlice
```

Red run before the final fix:

```text
callExternal allocations per small-arity call = 1, want 0
```

Green run after the fix:

```bash
go test ./vm -run '^TestCallExternalSmallArityDoesNotAllocateArgSlice$' -count=1
```

```text
ok  	github.com/t04dJ14n9/gig/vm	0.331s
```

## Benchmarks

Before baseline from the M3 Pro pprof run:

```text
BenchmarkGig_ExtCallMixed-12  23986  148891 ns/op  126205 B/op  4258 allocs/op
```

After, five runs:

```text
BenchmarkGig_ExtCallMixed-12  30458  119437 ns/op  40533 B/op  2082 allocs/op
BenchmarkGig_ExtCallMixed-12  31200  118221 ns/op  40537 B/op  2082 allocs/op
BenchmarkGig_ExtCallMixed-12  31622  114249 ns/op  40534 B/op  2082 allocs/op
BenchmarkGig_ExtCallMixed-12  31653  113304 ns/op  40528 B/op  2082 allocs/op
BenchmarkGig_ExtCallMixed-12  30885  117880 ns/op  40537 B/op  2082 allocs/op
```

Average after:

```text
116618 ns/op
40534 B/op
2082 allocs/op
```

Delta against the before baseline:

```text
runtime:     21.68% faster
bytes/op:    67.88% lower
allocs/op:   51.10% lower
```

## pprof Evidence

Before allocation profile:

```text
flat       flat%   cum        cum%    symbol
2122.59MB  67.41%  3132.61MB  99.49% github.com/t04dJ14n9/gig/vm.(*vm).callExternal
```

Before line-level evidence:

```text
vm/call_external.go:22  args := make([]value.Value, numArgs)
flat: 2.07GB, cum: 2.07GB
```

After allocation profile:

```text
flat      flat%   symbol
499.02MB  40.36% strings.NewReader
343.51MB  27.78% github.com/t04dJ14n9/gig/model/value.MakeString
332.01MB  26.85% github.com/t04dJ14n9/gig/model/value.makeReflectValue
35.50MB    2.87% internal/strconv.FormatInt
```

After line-level `callExternal` profile:

```text
ROUTINE github.com/t04dJ14n9/gig/vm.(*vm).callExternal
flat: 0
```

`callExternal` still has cumulative allocation because it calls stdlib wrappers that allocate real result objects, but the argument packaging allocation is gone.

## Verification

Commands run:

```bash
go test ./vm -run '^TestCallExternalSmallArityDoesNotAllocateArgSlice$' -count=1
go test ./vm -count=1
go test ./... -count=1
cd benchmarks && go test -run '^$' -bench '^BenchmarkGig_ExtCallMixed$' -benchmem -benchtime=3s -count=5
cd benchmarks && go test -run '^$' -bench '^BenchmarkGig_ExtCallMixed$' -benchmem -benchtime=2s -count=1 -memprofile /private/tmp/gig-extcallmixed-after.mem -cpuprofile /private/tmp/gig-extcallmixed-after.cpu
```

Results:

```text
go test ./vm: pass
go test ./...: pass
BenchmarkGig_ExtCallMixed after: 113304-119437 ns/op, 40528-40537 B/op, 2082 allocs/op
```

## Remaining Bottleneck

The next bottleneck is no longer `callExternal` argument packaging. The allocation profile now points at real external-call payload/result costs:

- `strings.NewReader`
- `value.MakeString`
- `value.makeReflectValue`
- `strconv.Itoa` string creation

Further optimization should focus on value wrapping and stdlib DirectCall return paths, not the external-call argument slice.
