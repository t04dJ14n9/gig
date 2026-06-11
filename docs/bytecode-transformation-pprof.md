# Go to Gig Bytecode Samples and pprof Bottleneck

Generated on Apple M3 Pro (`darwin/arm64`, Go 1.26.3) from the current checkout. The bytecode samples were produced by compiling through `gig.Build` and reading `Program.InternalProgram()`.

Commands used:

```bash
cd benchmarks
go test -run '^$' -bench '^BenchmarkGig_ArithSum$' -benchmem -benchtime=3s -count=1 -cpuprofile /private/tmp/gig-arithsum.cpu
go test -run '^$' -bench '^BenchmarkGig_ExtCallMixed$' -benchmem -benchtime=3s -count=1 -cpuprofile /private/tmp/gig-extcallmixed.cpu
go test -run '^$' -bench '^BenchmarkGig_ExtCallMixed$' -benchmem -benchtime=2s -count=1 -memprofile /private/tmp/gig-extcallmixed.mem
```

## Sample 1: external function call

Go source:

```go
package main

import "strings"

func ExternalCall(name string) bool {
	return strings.Contains(strings.ToUpper(name), "GIG")
}
```

Observed result:

```text
Run result: true
Function: ExternalCall params=1 locals=3 hasIntLocals=false bytes=27
```

Constant pool:

```text
[0] external func strings.ToUpper stdlib=true call=direct variadic=false numIn=1
[1] string = "GIG"
[2] external func strings.Contains stdlib=true call=direct variadic=false numIn=2
[3] bool = true
```

Bytecode:

```text
0000  LOCAL                            0
0003  CALLEXTERNAL                     idx=0 argc=1
0007  SETLOCAL                         1
0010  LOCAL                            1
0013  CONST                            1
0016  CALLEXTERNAL                     idx=2 argc=2
0020  SETLOCAL                         2
0023  LOCAL                            2
0026  RETURNVAL
```

What this means:

- `LOCAL 0` loads the `name` parameter.
- `CALLEXTERNAL idx=0 argc=1` calls the constant-pool entry for `strings.ToUpper`.
- `CONST 1` loads `"GIG"`.
- `CALLEXTERNAL idx=2 argc=2` calls `strings.Contains`.
- Both external calls are DirectCall-backed stdlib wrappers, so the VM avoids `reflect.Call` for these two calls.

## Sample 2: self-defined struct

Go source:

```go
package main

type User struct {
	Name string
	Age  int
}

func StructScore() int {
	u := User{Name: "ada", Age: 41}
	u.Age = u.Age + 1
	return u.Age + len(u.Name)
}
```

Observed result:

```text
Run result: 45
Function: StructScore params=0 locals=13 hasIntLocals=false bytes=129
```

Constant/type pools:

```text
Constants:
[1] string = "ada"
[2] int = 41
[3] int = 1

Types:
[0] main.User
```

Bytecode:

```text
0000  NEW                              0
0003  SETLOCAL                         0
0006  LOCAL                            0
0009  FIELDADDR                        0
0012  SETLOCAL                         1
0015  LOCAL                            0
0018  FIELDADDR                        1
0021  SETLOCAL                         2
0024  LOCAL                            1
0027  CONST                            1
0030  SETDEREF
0031  LOCAL                            2
0034  CONST                            2
0037  SETDEREF
0038  LOCAL                            0
0041  FIELDADDR                        1
0044  SETLOCAL                         3
0047  LOCAL                            3
0050  DEREF
0051  SETLOCAL                         4
0054  ADDLOCALCONST                    4, 3
0059  SETLOCAL                         5
0062  LOCAL                            0
0065  FIELDADDR                        1
0068  SETLOCAL                         6
0071  LOCAL                            6
0074  LOCAL                            5
0077  SETDEREF
0078  LOCAL                            0
0081  FIELDADDR                        1
0084  SETLOCAL                         7
0087  LOCAL                            7
0090  DEREF
0091  SETLOCAL                         8
0094  LOCAL                            0
0097  FIELDADDR                        0
0100  SETLOCAL                         9
0103  LOCAL                            9
0106  DEREF
0107  SETLOCAL                         10
0110  LOCAL                            10
0113  LEN
0114  SETLOCAL                         11
0117  ADDLOCALLOCAL                    8, 11
0122  SETLOCAL                         12
0125  LOCAL                            12
0128  RETURNVAL
```

What this means:

- `NEW 0` allocates a `main.User` value using type-pool entry `0`.
- Struct literal initialization lowers to field addresses plus stores: `FIELDADDR 0`/`FIELDADDR 1`, then `SETDEREF`.
- `u.Age = u.Age + 1` lowers to `FIELDADDR 1`, `DEREF`, `ADDLOCALCONST`, then `SETDEREF`.
- `len(u.Name)` lowers through `FIELDADDR 0`, `DEREF`, then `LEN`.

## CPU pprof: pure VM path

Benchmark:

```text
BenchmarkGig_ArithSum-12  108548  33894 ns/op  393 B/op  6 allocs/op
```

`go tool pprof -top /private/tmp/gig-arithsum.cpu`:

```text
flat    flat%   cum     cum%   symbol
2.47s   67.67%  2.87s   78.63% github.com/t04dJ14n9/gig/vm.(*vm).run
0.30s    8.22%  0.30s    8.22% github.com/t04dJ14n9/gig/vm.(*vm).run.func2 (readU16)
0.08s    2.19%  0.08s    2.19% github.com/t04dJ14n9/gig/model/value.truncateInt
```

Main bottleneck for numeric code is still the interpreter loop itself:

- per-instruction loop/check overhead in `vm.(*vm).run`;
- operand decode in the inlined `readU16` helper;
- local/stack movement and tagged `value.Value` updates.

This is why Gig can still be slow on some paths even with bytecode: bytecode removes AST walking, but it still executes many tiny VM instructions per Go operation.

## CPU/alloc pprof: external-call path

Benchmark:

```text
BenchmarkGig_ExtCallMixed-12  23986  148891 ns/op  126205 B/op  4258 allocs/op
```

CPU profile:

```text
flat    flat%   cum     cum%   symbol
0.72s    7.27%  2.64s   26.67% github.com/t04dJ14n9/gig/vm.(*vm).run
0.05s    0.51%  1.68s   16.97% github.com/t04dJ14n9/gig/vm.(*vm).callExternal
0.02s    0.20%  0.91s    9.19% github.com/t04dJ14n9/gig/vm.(*vm).callResolvedExternal
```

Allocation profile:

```text
flat       flat%   cum        cum%    symbol
2122.59MB  67.41%  3132.61MB  99.49% github.com/t04dJ14n9/gig/vm.(*vm).callExternal
415.01MB   13.18%   415.01MB  13.18% strings.NewReader
279.50MB    8.88%   279.50MB   8.88% github.com/t04dJ14n9/gig/model/value.MakeString
279.01MB    8.86%   279.01MB   8.86% github.com/t04dJ14n9/gig/model/value.makeReflectValue
```

Line-level allocation evidence:

```text
vm/call_external.go:22  args := make([]value.Value, numArgs)
flat: 2.07GB, cum: 2.07GB
```

Main bottleneck for external/string-heavy code is allocation at the call boundary:

- every `CALLEXTERNAL` allocates a new `[]value.Value` argument slice;
- DirectCall avoids `reflect.Call`, but not argument packaging;
- return wrapping for strings/readers adds `value.MakeString` and reflect-backed value allocations;
- the allocations drive GC/runtime time, visible in the external-call CPU profile.

The follow-up optimization is tracked in `docs/external-call-allocation-optimization-report.md`. The final implementation uses a zero-copy slice over the VM operand stack; a local `[8]value.Value` buffer still escaped through the `DirectCall` function field.
