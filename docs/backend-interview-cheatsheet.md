# Gig Backend Interview Cheat Sheet

Use this as the short version of `backend-interview-performance-report.md`.

## One-Minute Pitch

Gig is a Go interpreter built as a backend runtime system. Source code is parsed
and type-checked, lowered to `golang.org/x/tools/go/ssa`, compiled into Gig
bytecode, and executed by a stack VM with typed runtime values.

The important engineering story is not "a VM is automatically faster than
tree-walking." The stronger story is:

> I measured where the VM loses, found concrete representation and dispatch
> costs, and then made narrow compiler/VM changes that preserved semantics while
> reducing runtime work.

## Pipeline To Explain

```text
source
  -> parser.Parse                  type-checking, import policy, panic policy
  -> compiler/ssa.Build            Go SSA package
  -> compiler.Compile              bytecode, constants, types, call metadata
  -> runner.ExecuteInit            init once, snapshot globals
  -> runner.RunWithValues          pooled VM execution
  -> vm.run                        fetch/decode/execute loop
```

The clean separation is the main backend design point:

- Parser owns source safety and type checking.
- Compiler owns static lowering and optimization.
- Bytecode model owns compact executable metadata.
- VM owns runtime semantics and hot-path dispatch.
- Runner owns lifecycle, pooling, context cancellation, and globals.
- Importer owns which external packages are available.

## Performance Work Implemented

### 1. Native Int-Slice Lowering

Problem:

SSA lowered constant-size `make([]int, N)` into `new([N]int)` plus slicing. Gig
compiled that into a reflect-backed slice, so the int-slice superinstructions
could not stay on the native fast path.

Fix:

The compiler now recognizes the synthetic SSA `"makeslice"` allocation shape and
emits `OpMakeSlice` directly when it is safe.

Result on Apple M3 Pro:

- `BenchmarkGig_BubbleSort` moved from about `2.54 ms/op` and `39,811 allocs/op`
  to about `392-401 us/op` and `7 allocs/op` in the latest run.

Safe wording:

> The bottleneck was not the VM model in general; it was a representation
> mismatch introduced by SSA lowering. Fixing the lowering kept the hot loop on
> the native int-slice path.

### 2. Conservative Constant Folding

Problem:

Some literal-heavy code still executed constant arithmetic, constant branches,
and dead jump tails at runtime.

Fix:

`optimize.FoldConstants` runs before the regular dispatch-oriented optimizer. It
does three conservative things:

- Propagates constants through locals only in straight-line bytecode regions.
- Folds pure integer arithmetic and known boolean branches.
- Removes unreachable bytecode after unconditional jumps until the next jump
  target.

Safety rules:

- Do not fold divide/modulo by zero; preserve runtime panic behavior.
- Do not fold across jump targets.
- Do not treat standalone boolean pushes as branches.
- Keep this as local cleanup, not full CFG optimization.

Result on Apple M3 Pro:

- `BenchmarkGig_ConstFoldArithmetic` moved from about `88.5-89.9 us/op` to
  about `67-70 us/op`.

Safe wording:

> This is intentionally not a full SCCP or global DCE pass. It is a local
> pre-pass that removes provably constant work while keeping control-flow and
> panic semantics intact.

## Current Benchmark Snapshot

Latest focused command:

```bash
cd benchmarks
go test -run '^$' -bench '^BenchmarkGig_(ConstFoldArithmetic|ArithSum|BubbleSort)$' \
  -benchmem -benchtime=1s -count=3
```

Apple M3 Pro, latest recorded ranges:

- `BenchmarkGig_ConstFoldArithmetic`: `67.5-68.7 us/op`, `6 allocs/op`
- `BenchmarkGig_ArithSum`: `31.3-31.9 us/op`, `6 allocs/op`
- `BenchmarkGig_BubbleSort`: `388.9-399.6 us/op`, `7 allocs/op`

Use ranges, not single numbers. Mention the machine and command.
The raw output is recorded in `docs/benchmark-m3-pro-2026-06-11.md`.

## What Not To Claim

Avoid:

- "A VM is always faster than AST walking."
- "Gig is faster than Yaegi everywhere."
- "Constant folding explains the BubbleSort win."
- "This is a full optimizing compiler."

Better:

- "Gig wins on some paths because bytecode plus direct native wrappers reduce
  repeated interpretation/reflection overhead."
- "Gig was slow on some paths because representation choices still forced
  reflection or generic value dispatch."
- "The current optimizer is deliberately conservative: superinstructions,
  int-specialized locals, native int slices, constant folding, constant-branch
  folding, and narrow dead-code cleanup."

## Next Honest Roadmap

The next realistic performance work is:

- More loop superinstructions for `local = local + const` and compare/jump
  shapes.
- Range-aware slice fast paths for proven-safe loop indices.
- Broader constant propagation only after a small CFG/dataflow framework exists.
- A register VM or stack-cache design only if profiling shows dispatch remains
  the dominant cost after simpler bytecode passes.
