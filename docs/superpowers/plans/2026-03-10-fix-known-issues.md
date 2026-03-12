# Fix Known Issues Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix 4 confirmed bugs in Gig: `string([]byte)` conversion, pointer-receiver mutation, `init()` execution, and `range`-over-string rune values.

**Architecture:** Each fix is isolated to one or two files with no cross-cutting dependencies. All fixes are in the VM/iterator layer (runtime) — no changes to the compiler or bytecode format are needed for issues 1, 3, and 4. Issue 2 (pointer receiver) requires understanding the actual failure mode at runtime before coding.

**Tech Stack:** Go 1.23.1, `golang.org/x/tools/go/ssa`, `go/types`, `unicode/utf8`

---

## How to verify any fix

```bash
# Run only the known-issues tests
go test -v -run 'TestKnownIssue' ./tests/

# Run the full test suite (must stay green)
go test -race ./...

# Lint
golangci-lint run --timeout=5m
```

---

## Chunk 1: Issue 1 — `string([]byte)` conversion

### Background

`string([]byte{104,105})` should yield `"hi"`. Currently the VM's `OpConvert` handler
falls into a `default` branch that does `fmt.Sprintf("%v", val.Interface())`, which
produces `"[104 105]"`.

`[]byte` values are stored as `KindBytes` in Gig's value system (see `value/value.go`,
`MakeBytes()` / `Bytes()`). The fix is a one-liner case in the existing `switch
val.Kind()` block.

**Files:**
- Modify: `vm/ops_dispatch.go` lines ~776–788 (`OpConvert` → `types.String` case)

---

### Task 1: Fix `string([]byte)` in the VM

**Files:**
- Modify: `vm/ops_dispatch.go` (the `OpConvert` handler, `case types.String:` branch)

- [ ] **Step 1: Run the failing test to confirm baseline**

```bash
go test -v -run 'TestKnownIssue_BytesToString' ./tests/
```

Expected: FAIL — `got "[104 105]", want "hi"`

- [ ] **Step 2: Add `KindBytes` case in `OpConvert` → `types.String`**

In `vm/ops_dispatch.go`, find the block:

```go
case types.String:
    // Convert to string
    switch val.Kind() {
    case value.KindInt:
        vm.push(value.MakeString(string(rune(val.Int()))))
    case value.KindUint:
        vm.push(value.MakeString(string(byte(val.Uint()))))
    case value.KindString:
        vm.push(val)
    default:
        // Use reflection for other types
        vm.push(value.MakeString(fmt.Sprintf("%v", val.Interface())))
    }
```

Insert a new case **before** `default`:

```go
    case value.KindBytes:
        if b, ok := val.Bytes(); ok {
            vm.push(value.MakeString(string(b)))
        } else {
            vm.push(value.MakeString(""))
        }
```

The import `fmt` is already present; no new imports needed.

- [ ] **Step 3: Run the test suite**

```bash
go test -v -run 'TestKnownIssue_BytesToString' ./tests/
```

Expected: PASS (both subtests and the multi-variant table test)

- [ ] **Step 4: Run full suite to confirm no regressions**

```bash
go test -race ./...
```

Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add vm/ops_dispatch.go
git commit -m "fix: convert KindBytes to string using string([]byte) not fmt.Sprintf"
```

---

## Chunk 2: Issue 3 — `init()` not executed

### Background

`prog.Run("Compute")` never calls `init()`. The compiled `init` function exists in
`program.Functions["init"]` (the compiler does compile it) but `ExecuteWithValues` in
`vm/vm.go` simply looks up the requested function and runs it — there is no `init` call
anywhere in the execution path.

The fix must run `init()` exactly once per **program instance**, not per VM (VMs are
pooled and reused across calls). The `Program` struct in `gig.go` is the right owner.

**Files:**
- Modify: `gig.go` — `Program` struct + `RunWithContext` method (or the `Build` function)
- Modify: `vm/vm.go` — no change needed (use `ExecuteWithValues` directly on a fresh VM)

### Design choice

Running `init()` at `Build()` time (before the VMPool is created) is the cleanest
approach:
- `init` runs exactly once.
- Global state set by `init` is baked into the initial `globals` slice inside the
  `Program`.
- The `VMPool` is then created from the `Program` whose globals are already initialised.
- No `sync.Once` or `initExecuted` field needed.

The `Build()` function in `gig.go` is the right place: after `compiler.Compile()` and
before returning `&Program{...}`, run `init` if present.

---

### Task 2: Execute `init()` at `Build()` time

**Files:**
- Modify: `gig.go` (function `Build`, roughly lines 108–189)

- [ ] **Step 1: Run the failing test to confirm baseline**

```bash
go test -v -run 'TestKnownIssue_InitFunc' ./tests/
```

Expected: FAIL — `got nil result, want 42`

- [ ] **Step 2: Add init() execution in `Build()`**

In `gig.go`, inside `Build()`, after the call to `compiler.Compile(compiled, ...)` and
before `return &Program{...}`, insert:

```go
// Execute init() once so package-level globals are initialised.
if initFn, ok := compiled.Functions["init"]; ok {
    initVM := vm.New(compiled)
    initVM.ctx = context.Background()  // need to set ctx before using run()
    // Build a minimal frame for init (no args, no locals beyond the function's own)
    // Use Execute which handles frame creation internally.
    ctx := context.Background()
    if _, err := initVM.Execute("init", ctx); err != nil {
        return nil, fmt.Errorf("init() failed: %w", err)
    }
    _ = initFn
}
```

Wait — `vm.New(compiled)` gives a VM whose globals are all zero. After `initVM.Execute("init", ctx)` runs, the mutations to globals are in `initVM.globals`. We need to copy those globals back into the `Program` so all subsequent VMs start with initialised globals.

The `Program` struct (defined in `gig.go`) wraps the `*bytecode.Program`. Globals live in `vm.globals` (the VM's own slice), not in `bytecode.Program`. Each new VM created by `VMPool.New` creates its own `globals` slice initialised to zero.

**Revised approach:** store a `initialGlobals []value.Value` on `Program` (or on `bytecode.Program`), then in `VMPool.Get()` (or `vm.New()`) copy the initial globals into each new VM.

**Simplest correct approach:**

1. Add `InitialGlobals []value.Value` to `bytecode.Program`.
2. After compiling, run `init()` on a fresh VM, then snapshot `vm.globals` into `compiled.InitialGlobals`.
3. Change `vm.New(program)` to copy `program.InitialGlobals` (if set) as the starting globals instead of a zero slice.

**Detailed steps:**

- [ ] **Step 2a: Add `InitialGlobals` field to `bytecode.Program`**

File: `bytecode/program.go`

Find the `Program` struct and add:

```go
// InitialGlobals holds the global variable state after init() has run.
// New VMs copy this slice as their starting globals.
// Nil means init() was not present or not yet run.
InitialGlobals []Value
```

(Note: `bytecode.Program` uses `value.Value` but lives in a different package — check the actual import. The field type should match `vm.globals` which is `[]value.Value`.)

- [ ] **Step 2b: Run `init()` at Build time and snapshot globals**

File: `gig.go`, in `Build()`, after `compiler.Compile(...)` succeeds, before `return`:

```go
// If the program has an init() function, run it once and snapshot the
// resulting globals so all subsequent VMs start with initialised state.
if _, hasInit := compiled.Functions["init"]; hasInit {
    initVM := vm.New(compiled)
    ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
    defer cancel()
    if _, err := initVM.Execute("init", ctx); err != nil {
        return nil, fmt.Errorf("executing init(): %w", err)
    }
    // Snapshot the globals after init.
    snap := make([]value.Value, len(initVM.Globals()))
    copy(snap, initVM.Globals())
    compiled.InitialGlobals = snap
}
```

This requires exposing `vm.Globals()` (or making `globals` accessible):

```go
// Globals returns the VM's global variable slice (for snapshotting after init).
func (vm *VM) Globals() []value.Value {
    return vm.globals
}
```

Add this to `vm/vm.go`.

- [ ] **Step 2c: Copy `InitialGlobals` in `vm.New()`**

File: `vm/vm.go`, in `New()`:

```go
func New(program *bytecode.Program) *VM {
    globals := make([]value.Value, len(program.Globals))
    // Initialise with post-init snapshot if available.
    if len(program.InitialGlobals) == len(globals) {
        copy(globals, program.InitialGlobals)
    }
    return &VM{
        program: program,
        stack:   make([]value.Value, 1024),
        sp:      0,
        frames:  make([]*Frame, 64),
        fp:      0,
        globals: globals,
        extCallCache: &externalCallCache{
            cache: make([]*extCallCacheEntry, len(program.Constants)),
        },
    }
}
```

Also update `vm.Reset()` to restore from `InitialGlobals` instead of zeroing:

```go
func (vm *VM) Reset() {
    vm.sp = 0
    vm.fp = 0
    vm.panicking = false
    vm.panicVal = value.MakeNil()
    vm.ctx = nil
    vm.globalsPtr = nil
    // Restore globals to post-init snapshot (or zero if no init).
    if len(vm.program.InitialGlobals) == len(vm.globals) {
        copy(vm.globals, vm.program.InitialGlobals)
    } else {
        for i := range vm.globals {
            vm.globals[i] = value.Value{}
        }
    }
}
```

- [ ] **Step 3: Check `bytecode/program.go` for the actual `Program` struct shape**

Read `bytecode/program.go` to confirm where `Globals` is declared and the correct
`Value` import path before editing.

- [ ] **Step 4: Run the failing tests**

```bash
go test -v -run 'TestKnownIssue_InitFunc' ./tests/
```

Expected: PASS for both `TestKnownIssue_InitFuncExecuted` and `TestKnownIssue_InitFuncSideEffect`

- [ ] **Step 5: Run full suite**

```bash
go test -race ./...
```

Expected: all PASS

- [ ] **Step 6: Commit**

```bash
git add bytecode/program.go vm/vm.go gig.go
git commit -m "fix: execute init() at Build time and snapshot globals for each VM"
```

---

## Chunk 3: Issue 4 — `range`-over-string rune values

### Background

Go's `for i, r := range str` yields `(byteIndex int, runeValue rune)`. The current
`iterator.next()` in `vm/iterator.go` handles strings in the same branch as slices/arrays:

```go
case value.KindSlice, value.KindArray, value.KindString:
    key = value.MakeInt(int64(it.index))
    val = it.collection.Index(it.index)
    return key, val, true
```

`it.index` advances by `+1` in `ops_dispatch.go`'s `OpRangeNext` handler, which is
correct for byte-by-byte iteration but wrong for rune-by-rune iteration. Also,
`value.Value.Index()` on a string returns the raw byte (`uint8`) at that position, not
the decoded rune.

The fix splits the string case out of the shared branch and uses `unicode/utf8` to
decode one rune at a time, advancing by `runeSize` bytes each step.

**Files:**
- Modify: `vm/iterator.go` (split string out of the shared slice/string branch)
- Possibly modify: `vm/ops_dispatch.go` (`OpRangeNext` increments `iter.index++` — after
  the fix this increment must be removed from `OpRangeNext` and done inside
  `iterator.next()` instead, since rune sizes vary)

**Design note on the `index++` location:**

`OpRangeNext` does `iter.index++` **before** calling `iter.next()`. The current code:

```go
iter.index++
key, val, ok := iter.next()
```

After the fix, for strings the index must advance by `runeSize` (1–4 bytes), not 1.
Two options:
1. Move the increment entirely inside `iterator.next()` for all types.
2. Keep `index++` in `OpRangeNext` for non-string types, add a separate `stringNext()`.

**Option 1 is cleaner.** Change `iterator.next()` to also advance the index, and remove
`iter.index++` from `OpRangeNext`.

---

### Task 3: Fix `range`-over-string rune iteration

**Files:**
- Modify: `vm/iterator.go`
- Modify: `vm/ops_dispatch.go` (remove `iter.index++`)

- [ ] **Step 1: Run the failing tests to confirm baseline**

```bash
go test -v -run 'TestKnownIssue_RangeString' ./tests/
```

Expected: FAIL — `got 0, want 294`

- [ ] **Step 2: Refactor `iterator.next()` to own the index increment**

In `vm/iterator.go`, change `next()` so it advances `it.index` itself.
This is simpler and allows per-type advance amounts.

Add `unicode/utf8` to the import block.

Replace the current `next()` body with:

```go
func (it *iterator) next() (key, val value.Value, ok bool) {
    switch it.collection.Kind() {
    case value.KindString:
        s := it.collection.String()
        if it.index >= len(s) {
            return value.MakeNil(), value.MakeNil(), false
        }
        r, size := utf8.DecodeRuneInString(s[it.index:])
        key = value.MakeInt(int64(it.index))
        val = value.MakeInt(int64(r)) // rune as int (matches Go spec: rune == int32)
        it.index += size
        return key, val, true

    case value.KindSlice, value.KindArray:
        if it.index >= it.collection.Len() {
            return value.MakeNil(), value.MakeNil(), false
        }
        key = value.MakeInt(int64(it.index))
        val = it.collection.Index(it.index)
        it.index++
        return key, val, true

    case value.KindMap:
        if it.mapIter == nil {
            if rv, isValid := it.collection.ReflectValue(); isValid {
                it.mapIter = rv.MapRange()
            } else {
                return value.MakeNil(), value.MakeNil(), false
            }
        }
        if !it.mapIter.Next() {
            return value.MakeNil(), value.MakeNil(), false
        }
        key = value.MakeFromReflect(it.mapIter.Key())
        val = value.MakeFromReflect(it.mapIter.Value())
        return key, val, true

    case value.KindChan:
        val, ok = it.collection.Recv()
        return value.MakeNil(), val, ok

    default:
        // reflect fallback
        if rv, isValid := it.collection.ReflectValue(); isValid {
            switch rv.Kind() {
            case reflect.String:
                s := rv.String()
                if it.index >= len(s) {
                    return value.MakeNil(), value.MakeNil(), false
                }
                r, size := utf8.DecodeRuneInString(s[it.index:])
                key = value.MakeInt(int64(it.index))
                val = value.MakeInt(int64(r))
                it.index += size
                return key, val, true
            case reflect.Slice, reflect.Array:
                if it.index >= rv.Len() {
                    return value.MakeNil(), value.MakeNil(), false
                }
                key = value.MakeInt(int64(it.index))
                val = value.MakeFromReflect(rv.Index(it.index))
                it.index++
                return key, val, true
            case reflect.Map:
                if it.mapIter == nil {
                    it.mapIter = rv.MapRange()
                }
                if !it.mapIter.Next() {
                    return value.MakeNil(), value.MakeNil(), false
                }
                key = value.MakeFromReflect(it.mapIter.Key())
                val = value.MakeFromReflect(it.mapIter.Value())
                return key, val, true
            case reflect.Chan:
                v, ok := rv.Recv()
                if !ok {
                    return value.MakeNil(), value.MakeNil(), false
                }
                return value.MakeNil(), value.MakeFromReflect(v), true
            }
        }
        return value.MakeNil(), value.MakeNil(), false
    }
}
```

Note: `value.Value.String()` returns the string content for `KindString`. Confirm with
`grep -n 'func.*String()' value/value.go`.

- [ ] **Step 3: Remove `iter.index++` from `OpRangeNext` in `vm/ops_dispatch.go`**

Find the `OpRangeNext` case (around line 1092):

```go
case bytecode.OpRangeNext:
    iterVal := vm.pop()
    iter, ok := iterVal.Interface().(*iterator)
    if !ok {
        tuple := []value.Value{value.MakeBool(false), value.MakeNil(), value.MakeNil()}
        vm.push(value.FromInterface(tuple))
        return nil
    }
    iter.index++             // <-- REMOVE THIS LINE
    key, val, ok := iter.next()
    tuple := []value.Value{value.MakeBool(ok), key, val}
    vm.push(value.FromInterface(tuple))
```

Delete the `iter.index++` line. The iterator now owns its own advancement.

**IMPORTANT:** The iterator is initialised with `index: -1` in `OpRange`:

```go
vm.push(value.FromInterface(&iterator{collection: collection, index: -1}))
```

Since the old code did `iter.index++` first (bringing -1 → 0), and the new code starts
in `next()` with `it.index` at whatever was set during construction, change the initial
value to `0` in `OpRange`:

```go
vm.push(value.FromInterface(&iterator{collection: collection, index: 0}))
```

- [ ] **Step 4: Run the range-string tests**

```bash
go test -v -run 'TestKnownIssue_RangeString' ./tests/
```

Expected: PASS for all three subtests

- [ ] **Step 5: Run existing range tests to confirm slices/maps still work**

```bash
go test -v -run 'TestCompiler_ForRange|TestCompiler_RangeOverMap' ./tests/
```

Expected: all PASS

- [ ] **Step 6: Run full suite**

```bash
go test -race ./...
```

Expected: all PASS

- [ ] **Step 7: Commit**

```bash
git add vm/iterator.go vm/ops_dispatch.go
git commit -m "fix: range-over-string now yields correct rune values using UTF-8 decoding"
```

---

## Chunk 4: Issue 2 — Pointer-receiver method mutation

### Background

```go
type Counter struct{ n int }
func (c *Counter) Inc() { c.n++ }
func Compute() int {
    c := &Counter{}
    c.Inc()
    c.Inc()
    return c.n  // want 2, got 0
}
```

This is the most complex fix. The root cause requires investigation.

**What SSA generates for `c.Inc()`:**

In Go's SSA form, `c.Inc()` where `c *Counter` becomes:
```
t0 = &Counter{} : *Counter        // OpNew
t1 = *t0 : Counter                // OpDeref (to get the struct value for the call)
// OR:
t1 = Inc(t0)                      // direct call passing *Counter
```

Actually for pointer receivers in SSA, the call IS `t1 = (Counter).Inc(t0)` where `t0`
is `*Counter`. The pointer is passed directly as the first argument. Inside `Inc`, the
SSA operations are:
```
t0 = *c : Counter                // load the struct from the pointer
t1 = t0.n + 1 : int              // add 1 to n
t2 = &t0.n : *int                // FieldAddr to get address of n
*t2 = t1                         // store back — OpSetDeref
```

The issue is likely that `OpFieldAddr` on `structPtr` (a `*Counter` stored as a
`reflect.Value`) requires `CanAddr()` to be true. If the `Counter` was allocated with
`reflect.New`, the pointer IS addressable. But if `structPtr` is a copy stored in
`frame.locals`, mutations go to the copy.

**Investigation step:** Before writing code, use a debug print or a test that captures
the SSA output to understand exactly what instructions are generated and where the
mutation is lost.

---

### Task 4: Investigate pointer-receiver mutation

**Files:**
- Read only (investigation): `compiler/compile_instr.go`, `compiler/compile_value.go`,
  the VM's `OpFieldAddr`, `OpSetDeref`, `OpDeref` handlers

- [ ] **Step 1: Run the failing test**

```bash
go test -v -run 'TestKnownIssue_PointerReceiver' ./tests/
```

Expected: FAIL

- [ ] **Step 2: Add a debug test that prints the compiled bytecode for the Counter example**

Temporarily add a test in `tests/known_issues_test.go` or create a throwaway
`tests/debug_test.go`:

```go
func TestDebug_PointerReceiverBytecode(t *testing.T) {
    source := `
type Counter struct{ n int }
func (c *Counter) Inc() { c.n++ }
func Compute() int {
    c := &Counter{}
    c.Inc()
    c.Inc()
    return c.n
}
`
    prog, err := gig.Build(source)
    if err != nil {
        t.Fatalf("Build error: %v", err)
    }
    // Print all functions and their bytecode
    for name, fn := range prog.InternalProgram().Functions {
        t.Logf("=== %s ===", name)
        t.Logf("NumLocals=%d Instructions=%v", fn.NumLocals, fn.Instructions)
    }
}
```

This requires exposing `prog.InternalProgram()` or accessing the bytecode directly.
If that's not possible, use `t.Log(prog)` and look at any available debug method.

Alternatively, look for a `Disassemble()` or similar function in `bytecode/`.

- [ ] **Step 3: Trace the actual failure**

Based on the bytecode dump (or code reading), identify:
- Is `OpFieldAddr` receiving a valid `reflect.Value` pointer?
- Is `CanAddr()` true for the field?
- Does `OpSetDeref` successfully write the value?
- After `Inc()` returns, does `c.n` read back as 0 or the updated value?

Key file: `vm/ops_dispatch.go` lines 610–632 (`OpFieldAddr`) and 696–699 (`OpSetDeref`).

- [ ] **Step 4: Implement the fix based on investigation findings**

The most likely root causes and their fixes:

**Case A: The struct is passed by value to `Inc()` (not by pointer)**

If `compile_instr.go` emits `OpDeref` on `c` before pushing it as an argument to
`Inc()`, the callee gets a copy. Fix: don't dereference pointer receivers before the
call.

SSA guarantees that pointer-receiver methods receive the pointer itself as first arg, so
if the compiler is emitting an extra deref, remove it.

**Case B: `OpNew` creates a struct with `reflect.New`, but the Value returned stores a
`reflect.Value` of kind `Ptr` (not `Struct`). Inside `Inc()`, the locals slot receives
this pointer. `OpFieldAddr` then works on this pointer via `reflect.Ptr → Elem() →
Field()`. If `CanAddr()` is false, `vm.push(value.MakeFromReflect(field))` pushes a
non-addressable field copy.**

Fix in `OpFieldAddr`: if the field is not addressable, use `reflect.New` to create a
pointer to a copy, or detect that we're in a method with a pointer receiver and use
a different approach.

Actually the real fix is: in `OpFieldAddr`, if `s.Field(i).CanAddr()` is false, the
struct itself is a non-addressable copy. We must NOT be calling `OpFieldAddr` on a copy
in the first place — the issue is earlier.

**Case C: The pointer receiver is passed correctly, but `vm.Reset()` is called between
calls to `Inc()` and `Compute()`.**

This should not happen since all three calls are within a single `prog.Run("Compute")`.

After investigation, implement the targeted fix. Common fix patterns:

For Case A, in `compiler/compile_instr.go`, check before emitting `OpDeref` that
the value being dereffed is actually being used as a value (not a pointer receiver).

For Case B, in `vm/ops_dispatch.go` `OpFieldAddr`, change the non-addressable path:

```go
if field.CanAddr() {
    vm.push(value.MakeFromReflect(field.Addr()))
} else {
    // Field is not addressable — create a new pointer and use indirection.
    // This happens when the struct was not allocated on the heap.
    // Promote to heap by creating a new reflect.Value.
    newField := reflect.New(field.Type())
    newField.Elem().Set(field)
    vm.push(value.MakeFromReflect(newField))
}
```

Note: this alone won't fix writeback to the original struct. If the struct is a copy,
any `OpSetDeref` on `newField` will write to the copy, not the original.

The fundamental fix for Case B is to ensure that when a `*T` is created with `OpNew`,
the underlying struct is stored in a heap-allocated, addressable location from the start.
`reflect.New(structType)` already does this — `reflect.New` returns a `*T` where the
`T` is heap-allocated and addressable. If `OpNew` already does this correctly, the
issue must be elsewhere.

**Investigation is mandatory before writing fix code.**

- [ ] **Step 5: Run tests after fix**

```bash
go test -v -run 'TestKnownIssue_PointerReceiver' ./tests/
go test -race ./...
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add <changed files>
git commit -m "fix: pointer-receiver methods now persist mutations to struct fields"
```

---

## Chunk 5: Final validation

- [ ] **Run all known-issues tests**

```bash
go test -v -run 'TestKnownIssue' ./tests/
```

Expected: all 9 tests PASS

- [ ] **Run the full test suite with race detector**

```bash
go test -race ./...
```

Expected: all PASS

- [ ] **Run benchmarks to confirm no performance regression**

```bash
go test -bench=. -benchmem -count=3 -run='^$' ./tests/
```

Compare against pre-fix numbers; no significant (>10%) regression expected.

- [ ] **Run linter**

```bash
golangci-lint run --timeout=5m
```

Expected: no new issues

- [ ] **Final commit summary**

All four bugs are independently fixable. Recommended commit order:
1. `string([]byte)` (30 min — one-liner fix)
2. `init()` execution (1–2 hours — multi-file but clear path)
3. `range`-over-string (1 hour — iterator refactor)
4. Pointer receiver (2–4 hours — requires investigation first)
