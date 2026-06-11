# Gig Backend Interview and Performance Report

This report is written for two purposes:

1. Explain Gig clearly enough to discuss it in a backend systems interview.
2. Turn current benchmark/profile evidence into a focused performance roadmap.

It is intentionally candid. Gig has strong backend-relevant design choices:
SSA-based compilation, a stack VM, direct native call wrappers, VM pooling,
context cancellation, and explicit third-party type-boundary safety. It also has
known slow paths. The strongest interview story is not "a VM is always faster
than AST walking"; it is "we measured where the VM wins, found where generic
runtime semantics dominate, and have a concrete plan to remove those costs."

## Executive Summary

Gig is a Go interpreter implemented in Go. The public API compiles source into a
`Program`, then executes named functions through a stack-based bytecode VM:

- Source is parsed and type-checked with a custom importer.
- Typed AST is lowered to `golang.org/x/tools/go/ssa`.
- SSA is compiled into Gig bytecode.
- The VM executes bytecode with a compact tagged-union `value.Value`.
- External Go packages are registered through an importer registry and called
  through generated `DirectCall` wrappers when available.
- A runner owns VM pooling, init-state snapshots, optional stateful globals, and
  timeout/cancellation behavior.

Current local performance on Apple M3 Pro:

- Gig is strong versus Yaegi on recursion, closures, and external Go calls.
- The original baseline was weak on pure loop/slice-heavy workloads such as
  `BubbleSort`, but the first P0 fix now makes `BubbleSort` faster than the
  local Yaegi result.
- The biggest confirmed `BubbleSort` bottleneck was not "VM dispatch" in
  general. SSA lowered constant-size `make([]int, N)` into `new([N]int)` plus
  slicing, which produced a reflect-backed slice and forced fused int-slice
  opcodes into the fallback path.
- A conservative constant-folding pass now removes straight-line literal
  arithmetic, known constant branches, and the resulting unreachable jump tails
  before the regular bytecode optimizer,
  improving the new `ConstFoldArithmetic` workload from about `88.5-89.9 us/op`
  to about `67-70 us/op` after the dead-code cleanup on Apple M3 Pro.
- For tight arithmetic such as `ArithSum`, the bottleneck is mostly the VM
  fetch/decode loop and integer bookkeeping.

Two high-value optimizations are implemented: the compiler recognizes synthetic
SSA make-slice lowering and emits `OpMakeSlice` directly, and it folds local
constant arithmetic before runtime execution.

## Architecture You Can Present

### Public Pipeline

The user-facing entry point is `gig.Build` in `gig.go`. It wires options,
chooses the package registry, delegates compilation, executes `init()` once,
snapshots globals, and constructs a runner:

- `gig.go`: `Build`, `Program.Run`, `Program.RunWithContext`,
  `Program.RunWithValues`
- `compiler/build.go`: source -> parse/type-check -> SSA -> bytecode
- `runner/runner.go`: VM pool, init globals, shared globals, result unwrapping
- `vm/run.go`: bytecode execution loop
- `model/value/value.go`: tagged runtime value representation

The backend story:

```text
source code
  -> parser.Parse                 checks imports, panic policy, type-checks
  -> compiler/ssa.Build           builds SSA package
  -> compiler.NewCompiler.Compile emits bytecode and metadata
  -> runner.ExecuteInit           runs init once and snapshots globals
  -> runner.RunWithValues         gets a VM from the pool
  -> vm.run                       fetch/decode/execute loop
```

This is a strong design because it separates concerns:

- The parser/type-checker owns source validity and sandbox rules.
- The compiler owns static lowering and external-call resolution.
- The bytecode model owns compact executable metadata.
- The VM owns execution semantics and hot-path specialization.
- The runner owns lifecycle and concurrency policy.
- The importer registry owns what native packages are exposed.

### Stack VM and Tagged Values

`model/value.Value` is the core runtime value. It is a 32-byte tagged union:

```go
type Value struct {
    kind Kind
    size Size
    num  int64
    obj  any
}
```

Primitives stay in `num`, avoiding reflection for common integer/bool/float
operations. Composite values, reflect payloads, native slices, functions, and
interfaces live in `obj`.

The VM in `vm/run.go` hoists hot fields into locals:

- stack pointer
- current frame
- instruction bytes
- local slots
- integer locals
- prebaked constants

Then it inlines common opcodes:

- local load/store
- constants
- integer arithmetic and comparisons
- jumps
- returns
- selected memory operations
- superinstructions such as fused local arithmetic and int-slice get/set

This is why Gig does well on recursive calls and closure-heavy workloads: the
interpreter avoids building AST nodes at runtime and keeps most state in frames,
slots, and value tags.

### External Library Calls

External packages are registered through `importer.ExternalPackage`. Functions,
variables, constants, types, method direct calls, and interface proxies are
stored in the registry.

At compile time:

- `compiler/compile_ext.go` resolves external function and method metadata.
- `compiler/compile_external_boundary.go` checks whether custom Gig values are
  allowed to cross third-party boundaries.
- `model/bytecode` pre-resolves external calls into `ResolvedCall` entries.

At runtime:

- `vm/call_external.go` handles `OpCallExternal`.
- It pops arguments, does a lock-free lookup in `program.ExternCalls`, validates
  third-party boundaries, then chooses `DirectCall` or `reflect.Call`.
- `DirectCall` wrappers avoid `reflect.Value.Call` and convert `[]value.Value`
  straight into native Go calls.

This is a good backend interview topic because it combines:

- dependency injection through a package registry
- generated adapters for performance
- a reflect fallback for coverage
- runtime safety gates
- cancellation checks around external calls

## Measured Performance on Local M3 Pro

Environment:

- CPU: Apple M3 Pro
- OS: macOS 15.7.2
- Go: `go1.26.3 darwin/arm64`
- Command:

```bash
cd benchmarks
go test -bench='^Benchmark(Gig|Yaegi|Lua|Native)_(Fib25|ArithSum|BubbleSort|Sieve|ClosureCalls|ExtCall(DirectCall|Reflect|Method|Mixed))$' \
  -benchmem -count=3 -timeout=30m -run='^$'
```

Median-ish results from the three local runs, with the `BubbleSort` Gig cell
updated after the P0 int-slice lowering fix:

| Workload | Native Go | Gig | Yaegi | GopherLua | Gig vs Yaegi | Native/Gig |
|---|---:|---:|---:|---:|---:|---:|
| Fib25 | 215 us | 12.16 ms | 57.28 ms | 12.11 ms | Gig 4.7x faster | 56.5x |
| ArithSum | 259 ns | 32.1 us | 22.4 us | 19.1 us | Yaegi 1.4x faster | 124x |
| BubbleSort | 5.68 us | 397 us | 672 us | 358 us | Gig 1.7x faster | 70x |
| Sieve | 1.08 us | 116 us | 113 us | 101 us | roughly equal | 107x |
| ClosureCalls | 258 ns | 208 us | 462 us | 58.3 us | Gig 2.2x faster | 807x |

External calls:

| Workload | Native Go | Gig | Yaegi | Gig vs Yaegi | Native/Gig |
|---|---:|---:|---:|---:|---:|
| DirectCall | 19.2 us | 256 us | 759 us | Gig 3.0x faster | 13.4x |
| Reflect | 15.5 us | 170 us | 449 us | Gig 2.6x faster | 11.0x |
| Method | 11.4 us | 215 us | 577 us | Gig 2.7x faster | 18.8x |
| Mixed | 7.56 us | 154 us | 393 us | Gig 2.6x faster | 20.4x |

Interpretation:

- The strongest defensible claim is not "Gig beats every interpreter." It does
  not.
- The strongest claim is: Gig is substantially faster than Yaegi on recursion,
  closures, and external Go call workloads, especially where generated
  `DirectCall` wrappers avoid reflection.
- After the P0 int-slice lowering fix, Gig no longer loses the local
  `BubbleSort` benchmark to Yaegi; tight arithmetic remains the main simple-loop
  weakness.

Post-P0 BubbleSort/ArithSum check:

```bash
cd benchmarks
go test -run '^$' -bench '^BenchmarkGig_(BubbleSort|ArithSum)$' \
  -benchmem -benchtime=1s -count=5
```

Observed range on Apple M3 Pro:

| Workload | Before | After | Allocation Change |
|---|---:|---:|---:|
| BubbleSort | 2.54 ms/op | 395-398 us/op | 39,811 -> 7 allocs/op |
| ArithSum | 32.1 us/op | 31.5-32.0 us/op | unchanged at 6 allocs/op |

## Why a VM Can Still Lose to Yaegi

The intuitive assumption is: bytecode VM should beat AST tree walking. That is
often true, but only if the hot path stays specialized.

Yaegi is also not a naive recursive AST walker at execution time. Its
interpreter generates per-node execution closures and can run some simple
operations through direct reflect helpers. In simple loops, that can be quite
competitive.

Gig pays different costs:

- bytecode fetch/decode
- stack pointer updates
- `value.Value` tag checks
- local slot synchronization between `locals` and `intLocals`
- fallback paths for reference semantics
- reflect compatibility when values must behave like Go values at boundaries

Those costs are worth it when they buy fast calls, closures, direct external
calls, and controlled runtime behavior. They are expensive when the workload is
only a small loop doing `i++`, comparisons, and slice swaps.

## Profile Findings

### BubbleSort

Command:

```bash
cd benchmarks
go test -run '^$' -bench '^BenchmarkGig_BubbleSort$' -benchmem \
  -benchtime=3s -count=1 \
  -cpuprofile /private/tmp/gig-report-bubblesort.cpu \
  -memprofile /private/tmp/gig-report-bubblesort.mem
```

Result:

```text
BenchmarkGig_BubbleSort-12  2527302 ns/op  796942 B/op  39811 allocs/op
```

CPU profile highlights:

```text
60.40% cum  vm.(*vm).run
14.74% cum  model/value.MakeFromReflect
14.74% cum  model/value.setReflectPointerElem
```

Allocation profile highlights:

```text
85.47% flat alloc_objects  model/value.makeReflectValue
57.59% cum alloc_objects   vm.(*vm).runIntSliceGetFallbackRecovered
~42%  cum alloc_objects    vm.(*vm).runIntSliceSetFallbackRecovered
```

Root cause before the fix:

`BubbleSort` was not slow just because it is bytecode. It was slow because SSA
lowered constant-size `make([]int, 100)` into a temporary `new([100]int)` plus
`[:]`. The compiler then emitted `OpNew` + `OpSlice`, so the loop's static
int-slice access patterns were fused, but the runtime slice value was
reflect-backed. `locals[sIdx].IntSlice()` missed, and `OpIntSliceGet/Set`
fell into `runIntSliceGetFallback` and `runIntSliceSetFallback`. Those fallback
functions call `indexAddressValue`, `dereferenceValue`, and
`setDereferenceValue`, which can materialize addressable reflect values through
`value.MakeFromReflect`.

Implemented fix:

- `compiler/compile_make_slice_lowering.go` detects the synthetic SSA lowering
  only when the temporary int array allocation is marked by SSA as
  `"makeslice"` and has exactly one non-debug referrer: the slice construction.
- `compiler/compile_alloc.go` skips emitting `OpNew` for that temporary array.
- `compiler/compile_iteration.go` emits `OpMakeSlice` directly for the slice
  result.
- Ordinary arrays still use the existing reflect-backed path, preserving array
  value/copy and shared-slice semantics.

Post-fix result:

```text
BenchmarkGig_BubbleSort-12  394890-398300 ns/op  1316 B/op  7 allocs/op
```

Relevant files:

- `compiler/compile_make_slice_lowering.go` recognizes synthetic
  `make([]int, const)` lowering.
- `compiler/optimize/slice_fusion.go` matches
  `LOCAL LOCAL INDEXADDR SETLOCAL ...` patterns and emits
  `OpIntSliceGet`, `OpIntSliceSet`, and `OpIntSliceSetConst`.
- `model/bytecode/opcode.go` defines the int-slice superinstructions.
- `vm/run.go` implements the fast path when `locals[sIdx].IntSlice()` succeeds.
- `vm/run_int_slice_fallback.go` handles the slow fallback.
- `vm/reference.go` implements addressable index/deref/set semantics.

### ArithSum

Command:

```bash
cd benchmarks
go test -run '^$' -bench '^BenchmarkGig_ArithSum$' -benchmem \
  -benchtime=3s -count=1 \
  -cpuprofile /private/tmp/gig-report-arith.cpu
```

Result:

```text
BenchmarkGig_ArithSum-12  32132 ns/op  392 B/op  6 allocs/op
```

CPU profile highlights:

```text
58.07% flat  vm.(*vm).run
7.65% flat   vm.(*vm).run.func2
1.42% flat   model/value.truncateInt
```

Root cause:

`ArithSum` is mostly dispatch and integer bookkeeping. There is no huge reflect
fallback. The problem is instruction count and per-op overhead.

## Performance Roadmap

### P0: Remove Reflect Allocation from BubbleSort Int Slice Get/Set

Status: implemented.

Goal:

Make `BenchmarkGig_BubbleSort` stop allocating tens of thousands of objects per
run while preserving ordinary array semantics.

Current bad path:

```text
OpIntSliceGet/Set
  -> locals[sIdx].IntSlice() misses
  -> runIntSliceGetFallback / runIntSliceSetFallback
  -> indexAddressValue
  -> addressableValue
  -> value.MakeFromReflect
```

Actual cause:

The slice local was statically an int slice, but constant-size
`make([]int, N)` arrived from SSA as `new([N]int)` plus `[:]`. That produced a
reflect-backed slice even though the later loop was eligible for int-slice
superinstructions.

Implemented direction:

1. Added a compiler regression test:
   `TestBubbleSortIntSliceAccessesAreFused`.
2. Added a synthetic-make lowering pass in
   `compiler/compile_make_slice_lowering.go`.
3. Kept the optimization narrow: it fires only for SSA `"makeslice"`
   allocations with one non-debug referrer, the slice construction.
4. Verified ordinary array behavior with array/pointer/slice-heavy tests.
5. Re-ran:

```bash
go test ./vm ./compiler ./model/value
cd benchmarks
go test -run '^$' -bench '^BenchmarkGig_(BubbleSort|ArithSum)$' \
  -benchmem -benchtime=1s -count=5
```

Result:

- `BenchmarkGig_BubbleSort` dropped from about `2.54 ms/op` and
  `39,811 allocs/op` to about `395-398 us/op` and `7 allocs/op`.
- That is now faster than the earlier local Yaegi `BubbleSort` result of about
  `672 us/op`.
- `ArithSum` stayed essentially unchanged, confirming the fix targets the
  int-slice representation path rather than changing general dispatch cost.

### P1: Reduce Tight Loop Dispatch Cost

Goal:

Improve `ArithSum`, `Sieve`, and other loop-heavy integer code.

Current status:

`ArithSum` spends most CPU in `vm.run`, not reflection. That means further
optimizations need to reduce instruction count or per-instruction work.

Implemented first step:

- Added `optimize.FoldConstants`, a conservative constant-folding pre-pass.
- The pass propagates constants through locals only inside a straight-line
  bytecode region.
- It folds pure integer arithmetic and comparisons, skips divide/modulo by zero
  to preserve runtime panic behavior, and refuses to fold across jump targets.
- It also folds known boolean branches: taken branches become `OpJump`, and
  untaken branches are deleted. This preserves stack effects because the removed
  condition push and conditional-jump pop cancel each other out.
- The branch fold is guarded so standalone boolean pushes and jump-targeted
  instructions are left unchanged rather than being treated as branch sites.
- A narrow dead-code cleanup removes bytecode after unconditional jumps until
  the next jump target. That turns folded constant branches into straight-line
  bytecode without introducing full CFG optimization.
- For readability, the implementation is split by concern: pass orchestration,
  local propagation, branch folding, dead-code cleanup, rewrite construction,
  arithmetic semantics, and control-flow target discovery.
- On Apple M3 Pro, the targeted `BenchmarkGig_ConstFoldArithmetic` moved from
  about `88.5-89.9 us/op` to about `67-70 us/op` after the dead-code cleanup.

Implementation directions:

- Add larger superinstructions for common loop shapes:
  - `local = local + const`
  - `local = local + local`
  - compare-and-jump with local/const operands
  - increment/decrement local
- Avoid creating `value.Value` for intermediate integer results when both
  operands already live in `intLocals`.
- Use `intLocals` as the primary storage for int-only locals, and update
  `locals` only at boundaries that need a full `value.Value`.
- Consider a small trace cache for stable hot loops only after simpler
  superinstructions are exhausted.

Success target:

- `BenchmarkGig_ArithSum` should move closer to Yaegi's `~22 us/op`.
- No regression in recursion or external call benchmarks.

### P2: Preserve DirectCall Coverage for Third-Party Libraries

Goal:

Keep external libraries as a performance strength.

Current evidence:

Gig beats Yaegi by `~2.6x` to `~3.0x` on the external call benchmarks on M3 Pro.
This is because registered packages can provide `DirectCall` wrappers, avoiding
`reflect.Value.Call`.

Implementation directions:

- Keep `gig gen` output broad and predictable.
- Prefer generated wrappers for exported functions and common methods.
- Track DirectCall coverage in generated packages.
- Add benchmark gates for external call workloads before large refactors.

Success target:

- External call benchmarks should not regress while loop optimizations proceed.

### P3: Improve Benchmark Trustworthiness

Goal:

Make performance claims reproducible and interview-safe.

Implementation directions:

- Add a checked-in benchmark result artifact for each major environment:
  - Apple M3 Pro, darwin/arm64
  - AMD EPYC, linux/amd64
- Use `benchstat` when available.
- Update README tables when benchmark code changes.
- Avoid claims like "Gig is faster" without naming the workload, competitor,
  environment, and command.

## External Library Type Boundary Status

Current policy:

Gig supports external library calls, but third-party boundaries are conservative.

By default, third-party packages cannot receive interpreter-defined types when
those values cross as:

- `any` / `interface{}`
- unproxied non-empty interfaces
- concrete parameters containing a Gig-defined named type
- nested containers containing Gig-defined values
- interpreted functions hidden behind interface-shaped results

Allowed:

- standard library calls, which Gig treats as trusted interpreter domains
- primitive values
- slices/maps of primitive values
- registered external package types
- third-party interface parameters with a registered interface proxy
- concrete callbacks whose result types cannot smuggle interface values

Escape hatch:

```go
prog, err := gig.Build(src, gig.WithAllowUnsafeTypePass())
```

That disables compile-time and runtime boundary checks. It should only be used
when the host library is known not to inspect type identity or method sets.

Enforcement points:

- `compiler/typecheck.go`: recursively detects user-defined named types.
- `compiler/compile_external_boundary.go`: allows registered interface proxies.
- `vm/call_boundary.go`: runtime guard for hidden interface crossings.
- `vm/call_boundary_reflect.go`: scans nested reflect values, maps, slices, and
  structs for interpreter-defined values.
- `vm/interface_boundary_proxy.go`: builds registered native proxies.
- `importer/registry_interface_proxy.go`: stores proxy metadata by name, type,
  and interface method set.

Targeted verification commands:

```bash
go test ./compiler -run 'TestThirdPartyBoundary'
go test ./vm -run 'TestExternalBoundary|TestCallExternal|TestMakeInterpretedInterfaceAdapter|Test.*Boundary'
```

Both passed during this report pass.

## Why We Cannot "Add Types" to Gig Structs

This is the key conceptual point for interviews.

Gig can synthesize a Go struct shape at runtime, but it cannot create a real
named Go type with methods in the host program.

Why:

1. Go method sets are static.

   A method belongs to a named type compiled into a Go package. You cannot attach
   new methods to a `reflect.Type` at runtime.

2. `reflect.StructOf` creates anonymous structural types.

   It can create a struct with fields and tags, but not a real package-qualified
   named type with a method set. It cannot make `reflect.Type.Name`,
   `reflect.Type.PkgPath`, and methods behave like a compiled Go type.

3. Interface satisfaction depends on method sets.

   A third-party library doing `v.(sort.Interface)`, `reflect.Type.Implements`,
   `MethodByName`, or a type switch expects the real Go method set. A
   `reflect.StructOf` value with the same fields is not enough.

4. Type identity matters.

   Go has nominal type identity for named types. A Gig-defined `type User struct`
   is not the same runtime type as any native Go `User`, and it cannot become one
   just by carrying a string name in metadata.

5. Plugins/code generation are different tradeoffs.

   In theory, generating Go code and compiling a plugin could produce real Go
   types. That would be slow, platform-limited, unsafe for sandboxed execution,
   and incompatible with the current embedded interpreter model.

What Gig does instead:

- Internally, it carries dynamic type metadata in `value.InterpretedInterfaceValue`.
- For common formatting/error cases, `gigStructWrapper` implements interfaces
  such as `fmt.Stringer`, `fmt.Formatter`, `fmt.GoStringer`, and `error`.
- For third-party interfaces, it requires explicit registered proxies that
  implement the real host interface and forward method calls back into Gig.
- For unsafe deployments, `WithAllowUnsafeTypePass` lets users opt out.

This is not a weakness in the compiler design. It is a deliberate boundary
between Go's static nominal type system and an embedded interpreter that creates
types dynamically.

## How to Present Gig on a CV

Use wording that is specific and defensible:

- Built a Go interpreter using SSA lowering, custom bytecode, and a stack-based
  virtual machine.
- Implemented a 32-byte tagged-union value representation to avoid reflection
  for primitive operations.
- Added VM hot-path inlining and superinstructions for arithmetic, comparisons,
  branches, returns, and slice access.
- Designed a package registry and generated `DirectCall` wrappers so interpreted
  code can call native Go libraries without `reflect.Value.Call` on common paths.
- Implemented context cancellation, panic/recover handling, VM pooling, and
  optional stateful globals for embedded backend workloads.
- Added third-party boundary validation to prevent runtime type-identity bugs
  when dynamically synthesized structs cross into reflection-heavy libraries.
- Benchmarked against Yaegi, GopherLua, and native Go; identified and profiled
  remaining slow paths with pprof.

Avoid:

- "Gig is faster than Yaegi" as a blanket claim.
- "VM is always faster than AST walking."
- "Gig structs are real Go types."

Better:

- "Gig is faster than Yaegi on recursive, closure, and external-call workloads
  in our local benchmarks; after fixing synthetic int-slice lowering, the local
  BubbleSort benchmark is also faster than Yaegi."
- "The project taught me where VM design helps and where runtime representation
  costs dominate."

## Interview Explanation: One-Minute Version

Gig compiles Go source to SSA and then to a compact stack bytecode. The VM uses a
tagged `Value` type so integers and booleans do not go through reflection. I
added generated DirectCall wrappers for external packages, so common native Go
calls avoid `reflect.Value.Call`, which is why external-call benchmarks are
around 2.6x to 3.0x faster than Yaegi on my M3 Pro.

The main weakness is not the VM idea; it is where runtime representation falls
back from specialized values into reflect-compatible reference semantics. The
best example was BubbleSort: the profile showed most allocations came from
`makeReflectValue` through int-slice get/set fallback paths. I fixed that by
recognizing SSA's synthetic `new([N]int)` plus `[:]` lowering for
`make([]int, N)` and emitting `OpMakeSlice` directly. BubbleSort dropped from
about `2.54 ms/op` and `39,811 allocs/op` to about `0.397 ms/op` and
`7 allocs/op` on my M3 Pro.

For external libraries, we intentionally ban Gig-defined structs crossing into
third-party `any` or unproxied interface parameters. Go does not let an
interpreter attach methods to runtime-created `reflect.StructOf` types, so
third-party reflection would see the wrong method set and type identity. We use
wrappers/proxies where we can make the adaptation explicit.

## Next Work Items

1. Add a benchmark guard for `BenchmarkGig_BubbleSort` allocations so the
   synthetic-make lowering does not regress.
2. Re-run the full cross-interpreter benchmark table and update README
   performance tables with environment labels.
3. Add focused tests for named slice conversion aliasing, especially
   `sort.IntSlice(s)` sharing with `s`.
4. Re-profile `BubbleSort` after the fix to identify the next remaining hot
   opcode group.
5. Revisit tight-loop superinstructions for `ArithSum` and `Sieve`.
