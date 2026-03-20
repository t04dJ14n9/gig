# Bug Fix Report — gig interpreter

Date: 2026-03-20
Branch: `feature/dev_youngjin`

---

## Resolved Issue 34 — json.Encoder method dispatch collision

**Bug ID**: Bug 8 (formerly tracked in `tests/known_issue_test.go`)

**Severity**: High — causes panic at runtime

**Affected programs**: Any program using both `encoding/json` and `encoding/xml`
(or any two packages that define types with the same method name)

**Symptom**:
```
panic: interface {} is *json.Encoder, not *xml.Encoder
```

**Root cause**: Two-layer collision in method dispatch and type registration.

### Layer 1 — DirectCall cache key collision

`AddMethodDirectCall` in `importer/register.go` used **bare type name** as key:

```go
// BEFORE (broken)
key := typeName + "." + methodName   // e.g., "Encoder.Encode"
methodDirectCalls[key] = dc
```

When `encoding/json.Encoder` and `encoding/xml.Encoder` both call
`AddMethodDirectCall("Encoder", "Encode", ...)`, they write to the **same key**
`"Encoder.Encode"`. The last `init()` to run wins, overwriting the other.

### Layer 2 — Missing package in type importer

`convertReflectType` in `importer/typeconv.go` created `types.Named` types
with a `nil` package:

```go
// BEFORE (broken)
typeName := types.NewTypeName(0, nil, rt.Name(), nil)
//    nil package — compiler could never know this was "encoding/json.Encoder"
```

The compiler's `extractReceiverTypeName` could not extract a package-qualified
name because the package was always `nil`.

### Fix (3 files)

| File | Change |
|------|--------|
| `importer/typeconv.go` | `getOrCreateTypesPackage(pkgPath)` — creates/caches `*types.Package` from `reflect.Type.PkgPath()`. Named types now carry the correct package. |
| `importer/register.go` | Key changed to `"pkgPath.TypeName.MethodName"` (e.g., `"encoding/json.Encoder.Encode"`) |
| `compiler/compile_ext.go` | `extractReceiverTypeName` returns `pkg.Path() + "." + obj.Name()` (e.g., `"encoding/json.Encoder"`) |

**Impact**: All method dispatch lookups via `LookupMethodDirectCall` now use
fully-qualified keys. The fix also resolves the same collision for any other
pair of packages that share a type name (e.g., `encoding/gob.Encoder`).

**Regression test**: `TestCorrectnessResolvedIssue/JsonEncodeResolved`

---

## VMPool concurrent race condition

**Bug ID**: (no known issue tracker entry — silently fixed)

**Severity**: Critical — causes memory corruption and infinite recursion

**Affected programs**: Any program where multiple goroutines call `Run()`
concurrently on the same `*gig.Program` (or on programs sharing a cache key)

**Symptom**: Stack overflow, crash, or silent memory corruption in concurrent
workloads.

**Root cause**: `sync.Pool` is not safe for objects that may be in use by
another goroutine at the time of `Get()`.

```go
// BEFORE (broken)
type VMPool struct {
    pool sync.Pool  // Get() can return a VM still being used!
}
```

When goroutine A calls `Run()` and gets a VM from the pool, goroutine B can
call `Get()` and receive the **same VM instance** before goroutine A calls
`Put()`. Since `VM` state (`stack`, `frames`, `pc`, `sp`, `fp`) is mutated
during execution, both goroutines corrupt the same VM simultaneously.

### Fix (`vm/vm.go`)

Replaced `sync.Pool` with a **mutex-protected slice pool**:

```go
// AFTER (fixed)
type VMPool struct {
    mu    sync.Mutex
    vms   []*VM   // available VMs only
    newVM func() *VM
}

func (p *VMPool) Get() *VM {
    p.mu.Lock()
    if len(p.vms) > 0 {
        vm := p.vms[len(p.vms)-1]
        p.vms = p.vms[:len(p.vms)-1]
        p.mu.Unlock()
        return vm
    }
    p.mu.Unlock()
    return p.newVM()
}
```

Additionally, `VM.Reset()` now clears all frames to prevent stale frame
references from a previous execution leaking into the next:

```go
for i := range vm.frames {
    vm.frames[i] = nil
}
```

**Why not `sync.Pool`?**
`sync.Pool` is designed for objects that are **stateless between uses** (e.g.,
`bytes.Buffer` of a known size). A `VM` is stateful during execution — it
holds a stack, program counter, goroutine state, and deferred functions.
`sync.Pool.Get()` makes no guarantee that a previous borrower has finished
with the object.

**Why not per-goroutine VMs?**
That would require `VMPool.Get()` to allocate a new VM on every call when
called concurrently, which defeats the purpose of pooling. The mutex solution
serializes `Get/Put` access (not execution), so VMs are reused **safely**.

---

## Unary `^` (bitwise NOT) crash — discovered via fuzzing

**Bug ID**: (no known issue — discovered by `FuzzBitwise`)

**Severity**: High — causes panic at runtime

**Symptom**: `panic: runtime error: index out of range [-1]` in `vm.(*VM).pop`

**Root cause**: `compileUnOp` in `compiler/compile_instr.go` treated unary `^`
(bitwise NOT) the same as binary `^` (XOR), emitting `OpXor` directly.
`OpXor` is a **binary** opcode — it pops two values and pushes one.
For unary `^a`, only the operand `a` is on the stack, so `OpXor` pops two
values from a one-value stack → panic.

**Fix**: For unary `^`, emit the operand, then push an all-ones constant,
then `OpXor`. For each integer width:

| Type | All-ones constant |
|------|-------------------|
| `uint8` | `uint8(0xFF)` |
| `uint16` | `uint16(0xFFFF)` |
| `uint32` | `uint32(0xFFFFFFFF)` |
| `uint64` | `uint64(^uint64(0))` |
| `intN` | `intN(-1)` (all bits set in two's complement) |

The constant is added to the constant pool as a typed Go value so
`FromInterface` correctly stores it as `KindUint` or `KindInt` to match
the operand's value kind.

**Regression test**: `FuzzBitwise` (fuzz target in `tests/fuzz_test.go`)

---

## Previously resolved issues (for reference)

| ID | Issue | Resolution | Commit |
|----|-------|-----------|--------|
| Bug 1 | sort named-type conversion | `OpChangeType` opcode | ab36e12 |
| Bug 2 | time.Duration DirectCall | gentool fix | abe45f0 |
| Bug 3 | fmt.Stringer not called | `gigStructFormatter` + method resolver | ab36e12 |
| Bug 4 | `%T` wrong type name | qualified `_gig_id` PkgPath | ab36e12 |
| Bug 5 | `%v` leaks `_gig_id` | `gigStructFormatter` skips `_gig_id` | ab36e12 |
| Bug 6 | int64/uint64 narrowing | `MakeInt64`/`MakeUint64` in wrappers | abe45f0 |
| Bug 7 | bytes.Buffer.Cap() wrong | `make([]byte,len)+copy` instead of `[]byte(str)` | 5768f9b |
| Bug 8 | json.Encoder dispatch collision | package-qualified DirectCall keys | 686e971 |
| — | unary ^ (bitwise NOT) crash | push all-ones const then XOR | (in working tree, pending commit) |
