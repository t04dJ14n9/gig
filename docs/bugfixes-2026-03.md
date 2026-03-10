# Bug Fixes — March 2026

This document describes four confirmed bugs that were fixed in this release.

---

## Bug 1: `string([]byte{...})` conversion produced `"[104 105]"` instead of `"hi"`

**Symptom:** `string([]byte{104, 105})` returned the fmt-style slice representation
`"[104 105]"` instead of the UTF-8 string `"hi"`.

**Root cause:** The VM's `OpConvert` handler (in `vm/ops_dispatch.go`) for
`types.String` target did not have a `value.KindBytes` case. It fell through to the
`default` branch which used `fmt.Sprintf("%v", val.Interface())`, producing the slice
format.

**Fix:** Added a `case value.KindBytes:` branch in the `OpConvert → types.String` path
that calls `string(b)` on the underlying byte slice:

```go
case value.KindBytes:
    if b, ok := val.Bytes(); ok {
        vm.push(value.MakeString(string(b)))
    } else {
        vm.push(value.MakeString(""))
    }
```

**File changed:** `vm/ops_dispatch.go`

---

## Bug 2: Pointer-receiver methods did not persist mutations

**Symptom:** Calling a pointer-receiver method that modifies a struct field had no
visible effect. After `c.Inc(); c.Inc()`, `c.n` was still 0.

**Root cause:** Two separate issues:

1. **Method collection:** SSA methods are not `*ssa.Package` members — they hang off the
   type's method set. The compiler's `collectFuncs` loop only iterated
   `mainPkg.Members` and never saw pointer-receiver methods like `(*Counter).Inc`. As a
   result, the method was not in `c.funcIndex` and fell through to `compileExternalStaticCall`,
   where it was dispatched as an external call that did nothing.

2. **Unexported struct fields:** Even with the method found, `OpFieldAddr` called
   `field.Addr()` on the struct field to get a settable pointer. For unexported fields
   (like `n int` in a user-defined type), `reflect.Value.CanSet()` is false because the
   field is unexported. The fix uses `reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr()))`
   to create a settable `*T` pointer to the field.

**Fixes:**

*Compiler — method collection* (`compiler/compiler.go`):

```go
// Also collect methods defined on types in the package.
for _, member := range mainPkg.Members {
    t, ok := member.(*ssa.Type)
    if !ok { continue }
    for _, recv := range []types.Type{t.Type(), types.NewPointer(t.Type())} {
        mset := mainPkg.Prog.MethodSets.MethodSet(recv)
        for i := 0; i < mset.Len(); i++ {
            if fn := mainPkg.Prog.MethodValue(mset.At(i)); fn != nil && fn.Package() == mainPkg {
                collectFuncs(fn)
            }
        }
    }
}
```

*VM — `OpFieldAddr` with unexported fields* (`vm/ops_dispatch.go`):

```go
if field.CanAddr() {
    fieldPtr := reflect.NewAt(field.Type(), value.UnsafeAddrOf(field))
    vm.push(value.MakeFromReflect(fieldPtr))
}
```

*VM — `OpDeref` with unexported field pointers* (`vm/ops_dispatch.go`):

Added `rv.CanInterface()` guard before calling `rv.Interface().(*value.Value)` to avoid
a panic when dereferencing a pointer to an unexported field.

*Value layer — `UnsafeAddrOf`* (`value/container.go`):

Added helper:

```go
func UnsafeAddrOf(v reflect.Value) unsafe.Pointer {
    return v.Addr().UnsafePointer()
}
```

**Files changed:** `compiler/compiler.go`, `vm/ops_dispatch.go`, `value/container.go`

---

## Bug 3: `init()` function was never executed

**Symptom:** Package-level variables set in `init()` were zero/nil when accessed from
other functions. `initVal = 42` in init → `InitFuncResult()` returned 0.

**Root cause:** `Build()` compiled the package (which includes the `init` function) but
never executed it. Each call to `Run()` created a fresh VM with zero-initialized globals.

**Fix:**

1. Added `InitialGlobals []value.Value` field to `bytecode.Program` to store the
   post-init global state snapshot.

2. In `Build()`, after compilation, check if the user-defined `init#1` function exists.
   If so, run the `init` wrapper on a temporary VM and snapshot the resulting globals:

```go
if _, hasInit := compiled.Functions["init#1"]; hasInit {
    initVM := vm.New(compiled)
    if _, err := initVM.Execute("init", initCtx); err != nil {
        return nil, fmt.Errorf("executing init(): %w", err)
    }
    snap := make([]value.Value, len(initVM.Globals()))
    copy(snap, initVM.Globals())
    compiled.InitialGlobals = snap
}
```

   Note: The check uses `"init#1"` (the user-defined body), not `"init"` (the SSA
   wrapper), to avoid running init for programs that only have the SSA-generated wrapper
   with no user code. This prevents panics from auto-imported package init calls.

3. `vm.New()` and `vm.Reset()` now copy `InitialGlobals` into the VM's globals slice.

**Sub-bug: nil-slice append via `KindInvalid` globals**

An additional issue was that `OpAppend` for `append(nilSlice, elem...)` checked
`slice.IsNil()` to detect nil slices. However, zero-initialized global variables have
`kind = KindInvalid` (the Go zero value for the `Kind` type), not `KindNil`. `IsNil()`
returns false for `KindInvalid`, causing the append to fall through and return the
invalid value unchanged.

**Fix:** Changed the nil-slice branch to also handle `KindInvalid`:

```go
} else if slice.IsNil() || slice.Kind() == value.KindInvalid {
```

**Files changed:** `gig.go`, `bytecode/bytecode.go`, `vm/vm.go`, `vm/ops_dispatch.go`

---

## Bug 4: `range`-over-string always yielded rune value 0

**Symptom:** `for _, r := range "abc" { sum += int(r) }` returned 0 instead of
`97 + 98 + 99 = 294`. The rune variable was always 0.

**Root cause:** Three separate issues:

1. **Iterator used byte indexing:** `iterator.next()` handled strings in the same branch
   as slices/arrays, calling `it.collection.Index(it.index)` which returns a raw byte
   (`uint8`), not a rune.

2. **Index pre-increment:** `OpRangeNext` did `iter.index++` before calling `next()`.
   This was fine for fixed-width elements but conflicted with the variable-width rune
   advancement needed for strings.

3. **Untyped string constants compiled to nil:** The compiler's `compileConst` function
   handled `types.String` but not `types.UntypedString`. String literals used in `range`
   expressions (e.g., `range "abc"`) have type `untyped string` in SSA, so they compiled
   to `nil` in the constants pool.

**Fixes:**

*Compiler — untyped string constants* (`compiler/compile_value.go`):

Added `types.UntypedString` to the string case in `compileConst`:

```go
case types.String, types.UntypedString:
    if cnst.Value != nil {
        v = constant.StringVal(cnst.Value)
    } else {
        v = ""
    }
```

Similarly added `types.UntypedInt`, `types.UntypedRune`, `types.UntypedFloat`, and
`types.UntypedBool` to their respective cases.

*VM — iterator* (`vm/iterator.go`):

Rewrote `next()` to handle strings separately using `unicode/utf8.DecodeRuneInString`
and advance by rune byte size. The index increment was moved inside `next()` for all
types, removing the external `iter.index++` from `OpRangeNext`. Iterator initial index
changed from `-1` to `0`.

```go
case value.KindString:
    s := it.collection.String()
    if it.index >= len(s) {
        return value.MakeNil(), value.MakeNil(), false
    }
    r, size := utf8.DecodeRuneInString(s[it.index:])
    key = value.MakeInt(int64(it.index))
    val = value.MakeInt(int64(r))
    it.index += size
    return key, val, true
```

**Files changed:** `compiler/compile_value.go`, `vm/iterator.go`, `vm/ops_dispatch.go`
