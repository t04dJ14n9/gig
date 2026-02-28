# DirectCall: Eliminating reflect.Value.Call() for External Functions & Methods

## Problem

When Gig calls functions from external Go packages (e.g., `strings.Contains`, `fmt.Sprintf`),
it must cross the boundary between the interpreter's `value.Value` representation and Go's
native types. Before this optimization, **every external call** went through:

1. Convert `[]value.Value` args → `[]reflect.Value` (allocation-heavy)
2. `reflect.Value.Call()` (slow: safety checks, type validation, `runtime.call`)
3. Convert `[]reflect.Value` results → `[]value.Value`

This made external function calls the single slowest operation in the VM — roughly
**50x overhead** compared to native Go calls.

## Solution: Generated Typed Wrappers (DirectCall)

### Function DirectCall

For each external function with compatible parameter types, we generate a typed Go wrapper
at code-generation time. For example, `strings.Contains(s, substr string) bool` gets:

```go
func directcall_Contains(args []value.Value) value.Value {
    a0 := args[0].String()
    a1 := args[1].String()
    r0 := strings.Contains(a0, a1)
    return value.FromBool(r0)
}
```

No `reflect.Value` allocation. No `reflect.Value.Call()`. Just direct typed extraction
from `value.Value`, a native Go function call, and direct result wrapping.

### Method DirectCall (NEW)

Extended the same approach to **methods on external types**. For example,
`(*bytes.Buffer).WriteString(s string) (int, error)` gets:

```go
func directcall_method_Buffer_WriteString(args []value.Value) value.Value {
    recv := args[0].Interface().(*bytes.Buffer)
    a1 := args[1].String()
    r0, r1 := recv.WriteString(a1)
    // ... return tuple
}
```

The receiver is extracted via `.Interface().(T)` type assertion — no `reflect.MethodByName`
lookup at runtime.

## Architecture

### Code Generation Pipeline (`gentool/`)

```
Package type info (go/types)
    │
    ├── directcall.go: generateDirectCall()       → function wrappers
    ├── directcall.go: generateMethodDirectCalls() → method wrappers
    ├── resolve.go:    collectCrossPkgImports()    → import resolution
    │                  collectMethodImports()
    └── generator.go:  orchestration + output
           │
           ▼
    stdlib/packages/*.go  (generated, 1162 wrappers)
```

### Parameter Type Support

| Type Category | Example | Extraction |
|---|---|---|
| Basic types | `int`, `string`, `bool`, `float64` | `.Int()`, `.String()`, `.Bool()`, `.Float()` |
| Named (same-pkg) | `Regexp`, `Template` | `.Interface().(TypeName)` |
| Named (cross-pkg) | `time.Time`, `io.Reader` | `.Interface().(pkg.Type)` |
| Pointer to named | `*bytes.Buffer`, `*http.Request` | `.Interface().(*pkg.Type)` |
| Pointer to basic | `*int32`, `*int64` | `.Interface().(*int32)` |
| Slice types | `[]byte`, `[]string` | `.Bytes()`, `.Interface().([]string)` |
| Map types | `map[string]bool` | `.Interface().(map[string]bool)` |
| Empty interface | `any` / `interface{}` | `.Interface()` |
| Error interface | `error` | Converted via `value.ErrorFromValue()` |

### Compile-Time Resolution

The compiler resolves DirectCall wrappers at compile time (`compiler/compile_ext.go`),
storing them in `ExternalFuncInfo.DirectCall` and `ExternalMethodInfo.DirectCall`. At
runtime, the VM checks `DirectCall != nil` and calls it directly — zero map lookups.

### SSA External Method Pkg=nil Problem

For external package methods in SSA, `fn.Pkg`, `fn.Object().Pkg()`, and even
`named.Obj().Pkg()` are all `nil` because external types lack proper package binding.
We solved this by keying the method DirectCall registry on `typeName.methodName`
(without package path), since type names are sufficiently unique across stdlib packages.

## Coverage

| Category | Wrappers | Coverage |
|---|---|---|
| Function DirectCall | 619 / 671 | 92.2% |
| Method DirectCall | 543 | All eligible methods |
| **Total** | **1,162** | — |

Remaining ~8% of functions use parameter types that can't be statically wrapped
(e.g., `unsafe.Pointer`, variadic with complex element types).

## Benchmark Results

### External Call Benchmarks (5 runs, `benchstat`)

| Benchmark | Baseline (ns/op) | Optimized (ns/op) | Speedup | Memory Δ | Allocs Δ |
|---|---|---|---|---|---|
| ExtCallReflect | 1,319,800 | 359,100 | **3.7x** (−72.8%) | −62.9% | −64.6% |
| ExtCallMethod | 1,216,000 | 460,100 | **2.6x** (−62.2%) | −49.0% | −50.3% |
| ExtCallMixed | 730,300 | 330,500 | **2.2x** (−54.8%) | −39.4% | −45.1% |
| ExtCallDirectCall | 588,000 | 583,500 | ~same | ~ | ~ |

ExtCallDirectCall was already using function DirectCall in the baseline — the improvement
there came in previous work. The massive gains in the other three benchmarks come from
**method DirectCall** and **expanded function DirectCall coverage** (from ~460 to 619 wrappers).

### Gig vs Yaegi (post-optimization)

| Benchmark | Gig (ns/op) | Yaegi (ns/op) | Gig advantage |
|---|---|---|---|
| ExtCallDirectCall | 583,500 | 1,551,000 | **2.7x faster** |
| ExtCallReflect | 359,100 | 1,001,500 | **2.8x faster** |
| ExtCallMethod | 460,100 | 1,214,000 | **2.6x faster** |
| ExtCallMixed | 330,500 | 845,900 | **2.6x faster** |

### Core VM Benchmarks (no regression)

All core VM benchmarks (Fib25, ArithSum, BubbleSort, Sieve, ClosureCalls) show no
statistically significant change — the optimization is purely additive.

## Files Changed

| File | Role |
|---|---|
| `gentool/directcall.go` | Core wrapper generation for functions and methods |
| `gentool/generator.go` | Orchestrates generation, outputs method wrappers |
| `gentool/resolve.go` | Cross-package import collection for method signatures |
| `bytecode/bytecode.go` | Added `ExternalMethodInfo.DirectCall` field |
| `compiler/compile_ext.go` | Compile-time DirectCall resolution for methods |
| `vm/call.go` | Runtime fast path: `DirectCall != nil` → call directly |
| `importer/register.go` | Method DirectCall registry (`AddMethodDirectCall` / `LookupMethodDirectCall`) |
| `gig.go` | `packageLookupAdapter` wiring for method DirectCall |
| `stdlib/packages/*.go` | 20 regenerated packages with 1,162 total wrappers |
| `benchmarks/bench_test.go` | 12 new benchmarks (4 Gig + 4 Yaegi + 4 Native) |
