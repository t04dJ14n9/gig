# Architectural Design: Custom Types vs Third-Party Reflection

## Executive Summary

**Go's `reflect.StructOf` cannot attach methods to synthesized types.** This is a hard
constraint of the Go runtime — there is no API, workaround, or unsafe hack to add
methods to a type created at runtime. This means:

1. Any third-party library that does `reflect.Value.MethodByName()` on a gig-defined
   type will find nothing.
2. Any interface satisfaction check (`reflect.Type.Implements()`) will fail.
3. Type assertions to user-defined interfaces will fail.

**This is fundamentally unsolvable for the general case.** The only complete solutions
are either (a) code-generating real Go types ahead of time (which defeats the purpose
of an interpreter), or (b) banning the problematic pattern.

---

## The Problem Space

### What Works Today

| Pattern | Works? | Why |
|---------|--------|-----|
| `encoding/json.Marshal(myStruct)` | ✓ | Only needs field access + struct tags |
| `fmt.Sprintf("%v", myStruct)` | ✓ | `gigStructWrapper` intercepts |
| `fmt.Sprintf("%T", myStruct)` | ✓ | `SprintfExtern` special-cases `%T` |
| `myStruct.Method()` within gig | ✓ | VM uses `MethodsByName` dispatch |
| Field read/write via reflection | ✓ | `reflect.StructOf` preserves fields |
| Struct tags (`json:"name"`) | ✓ | Tags preserved in synthesis |

### What Breaks Today (Permanently)

| Pattern | Why It Breaks | Fixable? |
|---------|---------------|----------|
| `sort.Sort(myType)` | Calls `MethodByName("Len")` on synthesized type | Only via adapter hack |
| `heap.Init(myType)` | Same — interface method dispatch via reflection | Only via adapter hack |
| `errors.As(err, &myCustomError)` | `reflect.Type` doesn't match — different type identity | No |
| `io.MultiWriter(myWriter)` | Host calls `.Write()` on interface — can't find method | Only via adapter hack |
| Any `rv.MethodByName("X")` | No methods on `reflect.StructOf` types | **Never** |
| Any `rv.Type().Implements(iface)` | No method set on synthesized types | **Never** |
| `reflect.TypeOf(x) == someType` | Synthesized type ≠ any real Go type | **Never** |

### Root Cause: Go's Reflection is a Closed World

```go
// This is all Go gives us:
reflect.StructOf(fields) → reflect.Type (no methods, ever)

// There is NO:
reflect.TypeWithMethods(fields, methods) // doesn't exist
reflect.AddMethod(type, name, impl)      // doesn't exist
reflect.RegisterInterface(type, iface)   // doesn't exist
```

The Go team has explicitly rejected proposals to add method attachment to
`reflect.StructOf` (see golang/go#16522). Their position: methods require compile-time
dispatch tables (itabs) that cannot be safely created at runtime.

---

## Analysis of Possible Approaches

### Approach 1: Per-Interface Adapter Pattern (Current — Limited)

**What it does**: For specific known interfaces (`sort.Interface`, `heap.Interface`),
intercept at the `OpMakeInterface` boundary and replace the value with a Go struct
that implements the interface by delegating back to the interpreter.

**Coverage**: Only works for interfaces you explicitly code adapters for.

**Problems** (from the Codex review):
- Adapter replaces concrete type identity (type assertions fail)
- Temp VM has no access to caller's globals/context
- Over-broad matching (shape-based instead of exact)
- Must be hand-written for every interface you want to support

**Verdict**: Useful for a small set of stdlib interfaces but fundamentally doesn't
scale. Every new interface requires a new adapter. User-defined interfaces are
impossible to adapt (you'd need to generate code at runtime).

### Approach 2: `reflect.MakeFunc` + Interface Proxy (Theoretically Possible, Practically Impossible)

**Idea**: At runtime, for each user-defined type that needs to satisfy an interface:
1. Create a new struct type with a field for each interface method
2. Each field is a `func` created via `reflect.MakeFunc` that delegates to the interpreter
3. Wrap the whole thing in the interface

**Why it fails**:
- `reflect.MakeFunc` creates functions, not methods on a type
- You cannot make a `reflect.StructOf` type satisfy an interface
- There's no way to create a type that `reflect.Type.Implements()` returns true for
- Even if you could, `reflect.TypeOf(proxy) != reflect.TypeOf(original)` breaks identity

### Approach 3: CGo Itab Injection (Unsafe, Fragile, Wrong)

**Idea**: Use `unsafe.Pointer` to construct an interface header (itab + data pointer)
manually, pointing method slots to interpreted dispatch trampolines.

**Why it's wrong**:
- Go's interface dispatch tables (itabs) are internal runtime structures
- Layout changes between Go versions (and even between compiler flags)
- Would break with every Go upgrade
- Violates Go's memory safety guarantees
- Not portable across GOOS/GOARCH

### Approach 4: Code Generation (Pre-Compilation)

**Idea**: Before interpretation, scan the user's code for custom types, generate real
Go source code that defines those types with methods, compile it as a plugin, and
load it at runtime.

**Why it's impractical for gig**:
- Requires `go build` at runtime (seconds of latency, needs Go toolchain)
- Plugin support is limited (Linux/macOS only, no Windows)
- Defeats the purpose of a lightweight embedded interpreter
- Cannot handle dynamically constructed types

### Approach 5: Accept the Limitation and Enforce Safety (Recommended)

**Thesis**: The boundary between interpreted and compiled code is a **type erasure
boundary**. Custom types defined in gig scripts can never be true Go types.
Rather than fighting this with fragile hacks, we should:

1. **Clearly define what works** (field access, encoding, formatting)
2. **Clearly define what doesn't** (interface satisfaction via reflection)
3. **Provide a compile-time option to ban dangerous patterns**
4. **Fix the adapter for the handful of stdlib interfaces we explicitly support**

---

## Recommended Architecture

### Tier 1: Things That Just Work (No Changes Needed)

These patterns work because they only need field-level reflection:

- `encoding/json` (Marshal/Unmarshal)
- `encoding/xml`
- `encoding/gob`
- `database/sql.Scan` (struct scanning)
- `fmt.Sprintf` (via gigStructWrapper)
- Template engines (`text/template`, `html/template`)
- Struct-based validation libraries
- YAML/TOML parsers

**No action required.** These cover the majority of real-world use cases.

### Tier 2: Explicitly Supported Interfaces (Fixed Adapter)

For a curated set of stdlib interfaces where the pattern is:
"pass a custom type to a stdlib function that calls methods on it via reflection"

| Interface | Functions | Status |
|-----------|-----------|--------|
| `sort.Interface` | `sort.Sort`, `sort.Stable`, `sort.Reverse` | Fix adapter (Phase 1-3 from prior plan) |
| `heap.Interface` | `heap.Init`, `heap.Push`, `heap.Pop` | Fix adapter |
| `fmt.Stringer` | `fmt.Sprint`, `fmt.Fprintf` | Already works via wrapper |
| `error` | `errors.As`, `fmt.Errorf` | Partially works via wrapper |
| `io.Reader/Writer` | `io.Copy`, `io.MultiWriter` | Add adapter |

**Architecture for Tier 2**:

```
┌─────────────────────────────────────────────────┐
│              INTERPRETER WORLD                   │
│                                                  │
│  type MySlice []int                              │
│  func (s MySlice) Len() int { ... }             │
│  func (s MySlice) Less(i,j) bool { ... }        │
│  func (s MySlice) Swap(i,j) { ... }             │
│                                                  │
│  sort.Sort(MySlice{3,1,2})                       │
│       │                                          │
└───────┼──────────────────────────────────────────┘
        │ OpCallExternal for sort.Sort
        ▼
┌─────────────────────────────────────────────────┐
│         BOUNDARY INTERCEPTION LAYER              │
│                                                  │
│  Detect: argument goes to sort.Sort(Interface)   │
│  Detect: concrete type has compiled Len/Less/Swap│
│                                                  │
│  Create: sortAdapter{                            │
│    callerVM: v,          ← live VM reference     │
│    receiver: mySlice,    ← original value        │
│    program:  v.program,  ← for method lookup     │
│  }                                               │
│                                                  │
│  Call: sort.Sort(sortAdapter)                     │
│       │                                          │
└───────┼──────────────────────────────────────────┘
        │ Host calls adapter.Len(), adapter.Less(), etc.
        ▼
┌─────────────────────────────────────────────────┐
│              HOST WORLD (sort pkg)               │
│                                                  │
│  sortAdapter implements sort.Interface natively  │
│  Each method → callCompiledMethod on callerVM    │
│  callerVM has correct globals, ctx, goroutines   │
│                                                  │
└─────────────────────────────────────────────────┘
```

Key changes from current adapter:
1. **Adapter is injected at the call site** (in `callExternal` when we detect the
   target function accepts `sort.Interface`), NOT in `OpMakeInterface`
2. **Original concrete type identity preserved** — `OpMakeInterface` never touches it
3. **Caller VM threaded through** — no detached temp VM
4. **Errors propagated** — no panic swallowing

### Tier 3: Banned Patterns (Compile-Time Enforcement)

For patterns that **cannot be fixed** — where third-party code introspects the runtime
type of a gig value and the synthesized type is inherently wrong:

```go
// These will NEVER work correctly:
reflect.TypeOf(myStruct)              // returns synthesized type, not "real" type
errors.As(err, &myCustomError)        // reflect.Type mismatch
rv.Type().Implements(someInterface)   // always false
rv.MethodByName("CustomMethod")       // always invalid

// But users don't know this — they expect interpreter == compiler.
```

**Solution: `WithStrictTypeMode()` build option**

When enabled, the compiler statically detects when a user-defined type would cross
the host boundary in a way that requires type identity or method discovery:

```go
// Compile-time analysis:
// 1. For each user-defined type T, track all sites where a value of type T
//    flows into an external function call argument
// 2. Check the external function's parameter type:
//    - interface{}/any → OK (only field access needed)
//    - Named interface (error, io.Reader, sort.Interface) → check Tier 2 support
//    - Unsupported interface → COMPILE ERROR

prog, err := gig.Build(source, gig.WithStrictTypeMode())
// Error: "type 'MyWriter' passed to io.Copy as io.Reader — custom types cannot
//         satisfy external interfaces (use WithAllowUnsafeTypePass to override)"
```

### Tier 4: Escape Hatch for Power Users

For advanced users who understand the limitations and want to proceed anyway:

```go
gig.Build(source, gig.WithAllowUnsafeTypePass())
```

This disables the compile-time check and lets values flow across boundaries with no
guarantees. The user accepts that reflection-based code may see wrong types.

---

## Implementation Roadmap

### Phase 1: Fix Adapter Architecture (Addresses sort/heap)

**Move adapter injection from `OpMakeInterface` to `callExternal`.**

Instead of replacing the concrete type during interface boxing, intercept at the
actual call site where the value crosses into host code:

```go
// In callExternal, BEFORE calling the host function:
func (v *vm) callExternal(funcIdx, numArgs int) error {
    args := popArgs(numArgs)
    
    // NEW: Check if any argument needs an adapter
    for i, arg := range args {
        if adapter, ok := v.tryCreateAdapter(funcIdx, i, arg); ok {
            args[i] = adapter
        }
    }
    
    // Proceed with normal dispatch
    ...
}
```

The `tryCreateAdapter` function:
1. Looks up the external function's parameter type
2. If the parameter is a supported interface (sort.Interface, heap.Interface, io.Reader, etc.)
3. And the argument is a gig-synthesized type with compiled methods matching the interface
4. Creates an adapter that delegates to the caller VM

**Key insight**: The adapter is ephemeral — it exists only for the duration of the
external call. It doesn't replace the value in the interpreter's state. Type identity
is preserved everywhere except inside the host function call.

### Phase 2: Add Compile-Time Type Flow Analysis

Add a static analysis pass after SSA construction:

```go
// In compiler, after all functions are compiled:
if cfg.strictTypeMode {
    for _, fn := range allFunctions {
        for _, block := range fn.Blocks {
            for _, instr := range block.Instrs {
                if call, ok := instr.(*ssa.Call); ok {
                    if isExternalCall(call) {
                        checkTypeFlowSafety(call)
                    }
                }
            }
        }
    }
}
```

The `checkTypeFlowSafety` function examines each argument to an external call:
- Is the argument's type user-defined?
- Does the external function's parameter type require methods?
- Is that interface in our Tier 2 supported list?
- If not → compile error

### Phase 3: `WithDisableCustomTypes()` (Nuclear Option)

For users who want absolute safety:
- Ban ALL `type` declarations in the script
- Only allow primitives, slices, maps, and external types
- Guarantees no type-identity issues at any boundary

### Phase 4: Documentation and Diagnostics

Add runtime diagnostics that detect when a user hits a Tier 3 failure:

```go
// When sort.Sort panics or returns wrong results:
// "gig: warning: type 'MySlice' was passed to sort.Sort but the adapter could
//  not dispatch method 'Swap'. Consider using WithStrictTypeMode() to catch
//  these issues at compile time."
```

---

## Decision Matrix: Which Tier Does Each Pattern Fall Into?

```
User code                        External function          Tier   Action
─────────────────────────────────────────────────────────────────────────
type T struct{X int}             json.Marshal(t)            1      Works
type T struct{X int}             fmt.Sprint(t)              1      Works (wrapper)
type T []int; Len/Less/Swap      sort.Sort(t)               2      Adapter
type T []int; Push/Pop/Len/...   heap.Init(&t)              2      Adapter  
type T struct; Write([]byte)     io.Copy(dst, t)            2      Adapter (new)
type T struct; Error() string    errors.As(err, &t)         3      Ban or unsafe-pass
type T struct; Read([]byte)      customLib.Process(t)       3      Ban or unsafe-pass
type T struct                    reflect.TypeOf(t)          3      Ban or unsafe-pass
```

---

## What "100% Correctness" Actually Means

**100% correctness against compiled Go is impossible for custom types that cross the
reflection boundary.** This is a fundamental theorem, not an engineering limitation.

What IS achievable:

| Scope | Correctness | How |
|-------|-------------|-----|
| Pure interpreter code (no host boundary crossing) | 100% | Already works |
| Custom types + encoding (json, xml, gob) | 100% | Already works |
| Custom types + fmt | 100% | Already works (wrapper) |
| Custom types + sort/heap | 100% | Fix adapter (Phase 1) |
| Custom types + io.Reader/Writer | 100% | Add adapter (Phase 1) |
| Custom types + arbitrary reflect-heavy libraries | **Impossible** | Ban at compile time |
| No custom types (primitives + external types only) | 100% | `WithDisableCustomTypes()` |

**The honest answer**: If you need 100% correctness for ALL Go code that ANY third-party
library might run on your values, the only option is `WithDisableCustomTypes()`.

---

## Recommended Default Configuration

```go
// For maximum safety (recommended for production rule engines):
prog, err := gig.Build(source,
    gig.WithStrictTypeMode(),    // Compile-time check for unsafe boundary crossings
)

// For maximum freedom (recommended for prototyping/scripting):
prog, err := gig.Build(source)  // Default: custom types allowed, no checks

// For absolute safety (no type issues possible):
prog, err := gig.Build(source,
    gig.WithDisableCustomTypes(),  // No custom types at all
)
```

---

## Summary

| Approach | Coverage | Complexity | Recommended? |
|----------|----------|-----------|--------------|
| Do nothing (current) | ~60% | None | No — known bugs |
| Fix adapter + move to call site | ~85% | Medium | **Yes — Phase 1** |
| Add compile-time flow analysis | ~95% | High | **Yes — Phase 2** |
| `WithDisableCustomTypes()` | 100% | Low | **Yes — Phase 3** |
| Full type synthesis (impossible) | 100% | ∞ | No — Go limitation |

The recommended path is: **Fix the adapter → Add flow analysis → Offer nuclear option.**
This gives users a progressive safety ladder from "just works for most cases" to
"guaranteed correct by construction."
