# The Value System

> A deep dive into Gig's tagged-union `Value` type — the fundamental data representation that every stack slot, local variable, constant, and function argument is built on.

## Table of Contents

1. [Overview](#overview)
2. [Memory Layout](#memory-layout)
3. [Kind System](#kind-system)
4. [Constructors](#constructors)
5. [Accessors](#accessors)
6. [Arithmetic Operations](#arithmetic-operations)
7. [Comparison and Equality](#comparison-and-equality)
8. [Bitwise Operations](#bitwise-operations)
9. [Type Conversions](#type-conversions)
10. [Container Operations](#container-operations)
11. [Pointer Operations](#pointer-operations)
12. [Channel Operations](#channel-operations)
13. [Reflection Interop](#reflection-interop)
14. [Closure Support](#closure-support)
15. [Native Fast Paths](#native-fast-paths)
16. [Design Rationale](#design-rationale)
17. [Performance Characteristics](#performance-characteristics)
18. [File Organization](#file-organization)

---

## Overview

The `Value` type is the universal data unit in the Gig interpreter. Every local variable, stack slot, constant, function argument, and return value is a `Value`. The design follows a **tagged-union** pattern — a single struct that can hold any Go type, with a type tag (`kind`) to distinguish the variant at runtime.

The key insight is the **two-tier split**:

- **Primitive types** (`bool`, `int`, `uint`, `float`, `nil`) are stored directly in the struct's numeric field. No heap allocation. No reflection. No GC pressure.
- **Composite types** (`string`, `complex`, `[]byte`, `[]int64`, closures) are stored in the `obj` field as native Go objects — still avoiding reflection where possible.
- **Complex types** (`struct`, `map`, `chan`, generic slices, interfaces) are stored in `obj` as `reflect.Value` — the fallback for types that require dynamic dispatch.

This design is inspired by how Lua, Python, and other dynamic language runtimes represent values, but adapted for Go's static type system with seamless `reflect` interop.

---

## Memory Layout

```
┌──────────────────────────────────────────────────────────┐
│                    Value (32 bytes)                       │
├────────┬─────────┬───────────────────────────────────────┤
│  kind  │ padding │                 num                    │
│ 1 byte │ 7 bytes │               8 bytes                  │
├────────┴─────────┴───────────────────────────────────────┤
│                        obj                                │
│                      16 bytes                             │
│               (interface{} / any)                         │
└──────────────────────────────────────────────────────────┘
```

```go
type Value struct {
    kind Kind    // 1 byte: type tag
    num  int64   // 8 bytes: bool, int, uint bits, float64 bits
    obj  any     // 16 bytes: string, complex128, reflect.Value, native composites, or nil
}
```

The total size is **32 bytes** on 64-bit systems. This is deliberate:

- **32 bytes = 4 machine words** — fits cleanly in cache lines and CPU registers.
- **Primitives live in `kind` + `num`** with `obj` = nil, meaning they never touch the heap and never cause GC pressure.
- **The `obj` field** is a Go `interface{}` (16 bytes: type pointer + data pointer), which can hold any Go value without allocation for small types.

---

## Kind System

The `Kind` type is a `uint8` enum that identifies the variant stored in a `Value`:

```go
type Kind uint8

const (
    KindInvalid   Kind = iota  // zero value, unintialized
    KindNil                     // nil
    KindBool                    // bool (stored in num: 0 or 1)
    KindInt                     // int, int8, int16, int32, int64 (stored in num)
    KindUint                    // uint, uint8, ..., uint64, uintptr (stored in num as bits)
    KindFloat                   // float32, float64 (stored in num as float64 bits)
    KindString                  // string (stored in obj)
    KindComplex                 // complex64, complex128 (stored in obj)
    KindPointer                 // *T (stored in obj as reflect.Value or *int64 or *Value)
    KindSlice                   // []T (stored in obj as reflect.Value, []int64, or []Value)
    KindArray                   // [N]T (stored in obj as reflect.Value)
    KindMap                     // map[K]V (stored in obj as reflect.Value)
    KindChan                    // chan T (stored in obj as reflect.Value)
    KindFunc                    // func (stored in obj as *Closure or native func)
    KindStruct                  // struct (stored in obj as reflect.Value)
    KindInterface               // interface (stored in obj as reflect.Value)
    KindReflect                 // fallback: any type stored as reflect.Value
    KindBytes                   // []byte (stored in obj as native []byte, zero reflection)
)
```

### Storage Strategy by Kind

| Kind | `num` field | `obj` field | Heap alloc? | Reflection? |
|------|-------------|-------------|-------------|-------------|
| `KindNil` | 0 | nil | No | No |
| `KindBool` | 0 or 1 | nil | No | No |
| `KindInt` | int64 value | nil | No | No |
| `KindUint` | uint64 bits | nil | No | No |
| `KindFloat` | float64 bits | nil | No | No |
| `KindString` | — | `string` | No* | No |
| `KindComplex` | — | `complex128` | Yes** | No |
| `KindBytes` | — | `[]byte` | No* | No |
| `KindSlice` | — | `[]int64` / `[]Value` / `reflect.Value` | Varies | Sometimes |
| `KindPointer` | — | `*int64` / `*Value` / `reflect.Value` | Varies | Sometimes |
| `KindFunc` | — | `*Closure` / native func | No* | No |
| `KindMap` | — | `reflect.Value` | Yes | Yes |
| `KindStruct` | — | `reflect.Value` | Yes | Yes |
| `KindChan` | — | `reflect.Value` | Yes | Yes |
| `KindArray` | — | `reflect.Value` | Yes | Yes |
| `KindInterface` | — | `reflect.Value` | Yes | Yes |
| `KindReflect` | — | `reflect.Value` | Yes | Yes |

\* Strings and slices are reference types in Go — the `obj` field holds the header, not a copy.
\** `complex128` is boxed when stored in an `interface{}`.

---

## Constructors

### Primitive Constructors (Zero Allocation)

```go
MakeNil() Value                           // → KindNil
MakeBool(b bool) Value                    // → KindBool, num = 0/1
MakeInt(i int64) Value                    // → KindInt, num = i
MakeUint(u uint64) Value                  // → KindUint, num = int64(u)
MakeFloat(f float64) Value                // → KindFloat, num = float64 bits
MakeString(s string) Value                // → KindString, obj = s
MakeComplex(real, imag float64) Value     // → KindComplex, obj = complex128
```

Each is a pure struct literal — no function calls, no allocation. The Go compiler can inline these completely.

### Specialized Constructors (Avoiding Reflection)

```go
MakeBytes(b []byte) Value                 // → KindBytes, obj = b
MakeIntSlice(s []int64) Value             // → KindSlice, obj = s
MakeIntPtr(p *int64) Value                // → KindPointer, obj = p
MakeFunc(fn any) Value                    // → KindFunc, obj = fn
MakeValueSlice(vals []Value) Value        // → KindSlice, obj = vals
```

These exist to avoid wrapping common types in `reflect.Value`. For example, `MakeIntSlice` stores the `[]int64` directly — the VM can then index, modify, and take addresses of elements without any reflection.

### Generic Constructors

```go
MakeFromReflect(rv reflect.Value) Value   // → auto-detects kind from reflect.Kind
FromInterface(v any) Value                // → type-switch fast path, fallback to MakeFromReflect
```

`FromInterface` uses a type switch to detect common Go types (`bool`, `int`, `int8`, ..., `string`, `[]byte`) and routes them to the specialized constructors. Anything not in the switch falls through to `MakeFromReflect`, which wraps the value in `reflect.Value`.

The type switch ordering is optimized for the most common types first:

```go
func FromInterface(v any) Value {
    if v == nil { return MakeNil() }
    switch val := v.(type) {
    case bool:    return MakeBool(val)
    case int:     return MakeInt(int64(val))
    case int64:   return MakeInt(val)
    case string:  return MakeString(val)
    case []byte:  return MakeBytes(val)
    // ... 12 more primitive types ...
    }
    return MakeFromReflect(reflect.ValueOf(v))
}
```

### The `MakeFromReflect` Unwrapping Logic

`MakeFromReflect` is not a simple wrapper — it **unwraps** primitives out of `reflect.Value`:

```go
func MakeFromReflect(rv reflect.Value) Value {
    switch rv.Kind() {
    case reflect.Bool:
        return MakeBool(rv.Bool())          // → KindBool, no reflect stored
    case reflect.Int, reflect.Int8, ...:
        return MakeInt(rv.Int())            // → KindInt, no reflect stored
    case reflect.String:
        return MakeString(rv.String())      // → KindString, no reflect stored
    case reflect.Slice:
        if rv.Type().Elem().Kind() == reflect.Uint8 {
            return MakeBytes(rv.Bytes())    // → KindBytes for []byte
        }
        return Value{kind: KindReflect, obj: rv}
    default:
        return Value{kind: KindReflect, obj: rv}
    }
}
```

This means that even when a value enters the system through `reflect.Value`, primitives get extracted into the efficient tagged-union form. A `reflect.Value` containing `int(42)` becomes `Value{kind: KindInt, num: 42}`, not `Value{kind: KindReflect, obj: reflect.ValueOf(42)}`.

---

## Accessors

### Typed Accessors (Zero Reflection)

```go
v.Bool() bool           // panics if kind != KindBool
v.Int() int64           // panics if kind != KindInt
v.Uint() uint64         // panics if kind != KindUint
v.Float() float64       // panics if kind != KindFloat
v.String() string       // panics if kind != KindString
v.Complex() complex128  // panics if kind != KindComplex
v.Bytes() ([]byte, bool)      // returns (nil, false) if not KindBytes
v.IntSlice() ([]int64, bool)  // returns (nil, false) if not native []int64
v.IntPtr() (*int64, bool)     // returns (nil, false) if not native *int64
v.ValueSlice() ([]Value, bool) // returns (nil, false) if not native []Value
```

### Fast-Path Accessors (No Kind Check)

For the VM's hot path, where the kind is already known from the opcode:

```go
v.RawInt() int64    // direct num field read, no kind check
v.RawBool() bool    // v.num != 0, no kind check
v.RawObj() any      // direct obj field read, no kind check
```

These are designed to be inlined by the Go compiler.

### Generic Accessor

```go
v.Interface() any   // converts to Go interface{}, may call reflect.Value.Interface()
```

The `Interface()` method performs a full kind-switch:

- Primitives: returns the Go type directly (`int64`, `float64`, `string`, etc.)
- `KindFunc`: returns the raw `obj` (closure or native function)
- `KindBytes`: returns `[]byte`
- `KindReflect` and others: calls `reflect.Value.Interface()` on the stored value

### Reflection Accessor

```go
v.ReflectValue() (reflect.Value, bool)  // returns the internal reflect.Value if stored
```

### Query Methods

```go
v.Kind() Kind       // returns the kind tag
v.IsNil() bool      // true for KindNil, or KindReflect wrapping a nil pointer/chan/map/etc.
v.IsValid() bool    // true if kind != KindInvalid (zero-value Value is invalid)
v.CanInterface() bool // true if the internal reflect.Value (if any) supports Interface()
```

---

## Arithmetic Operations

All arithmetic operations use a kind-switch with the fast path for integers first:

```go
v.Add(other Value) Value    // + (int, uint, float, string concat, complex)
v.Sub(other Value) Value    // - (int, uint, float, complex)
v.Mul(other Value) Value    // * (int, uint, float, complex)
v.Div(other Value) Value    // / (int, uint, float, complex)
v.Mod(other Value) Value    // % (int, uint, float)
v.Neg() Value               // unary - (int, float, complex)
```

### How It Works

Take `Add` as an example:

```go
func (v Value) Add(other Value) Value {
    switch v.kind {
    case KindInt:
        return MakeInt(v.num + other.Int())     // raw int64 addition
    case KindUint:
        return MakeUint(uint64(v.num) + other.Uint())
    case KindFloat:
        return MakeFloat(v.Float() + other.Float())
    case KindString:
        return MakeString(v.obj.(string) + other.obj.(string))
    case KindComplex:
        a := v.obj.(complex128)
        b := other.obj.(complex128)
        return MakeComplex(real(a)+real(b), imag(a)+imag(b))
    default:
        panic("cannot add")
    }
}
```

For `KindInt`, the entire operation is:
1. Read `v.num` (8 bytes)
2. Read `other.num` (8 bytes, via `Int()` which checks `kind == KindInt`)
3. Add (one CPU instruction)
4. Write result into a new Value struct on the stack

No heap allocation. No function pointer indirection. No reflection.

### VM Optimization

In the VM dispatch loop, arithmetic is further optimized with an inlined fast path:

```go
case bytecode.OpAdd:
    sp--
    b := stack[sp]
    sp--
    a := stack[sp]
    if a.IsInt() {
        stack[sp] = value.MakeInt(a.RawInt() + b.RawInt())  // ~3 machine instructions
    } else {
        stack[sp] = a.Add(b)  // generic path
    }
    sp++
```

`a.IsInt()` is `a.kind == KindInt` — a single byte comparison. `a.RawInt()` is `a.num` — a single field read. The entire integer addition is ~3 machine instructions plus the dispatch overhead.

---

## Comparison and Equality

```go
v.Equal(other Value) bool   // == (all types; uses reflect.DeepEqual as fallback)
v.Cmp(other Value) int      // returns -1, 0, or 1 (bool, int, uint, float, string)
```

### Equality Semantics

`Equal` follows Go's equality rules:

1. **Different kinds**: returns `false` (except nil comparison: `KindNil` vs `KindReflect` wrapping a nil value).
2. **Same kind, primitive**: direct field comparison (`v.num == other.num`).
3. **Same kind, composite**: `reflect.DeepEqual(v.Interface(), other.Interface())`.

```go
func (v Value) Equal(other Value) bool {
    if v.kind != other.kind {
        if v.kind == KindNil || other.kind == KindNil {
            return v.IsNil() && other.IsNil()
        }
        return false
    }
    switch v.kind {
    case KindNil:    return true
    case KindBool:   return v.num == other.num
    case KindInt:    return v.num == other.num
    case KindFloat:  return v.Float() == other.Float()  // NaN != NaN
    case KindString: return v.obj.(string) == other.obj.(string)
    default:         return reflect.DeepEqual(v.Interface(), other.Interface())
    }
}
```

### Comparison Semantics

`Cmp` returns a three-way comparison result. It supports `bool`, `int`, `uint`, `float`, and `string` — the types with a natural ordering. Panics for unsupported types (maps, slices, etc.).

---

## Bitwise Operations

```go
v.And(other Value) Value      // & (int, uint)
v.Or(other Value) Value       // | (int, uint)
v.Xor(other Value) Value      // ^ (int, uint)
v.AndNot(other Value) Value   // &^ (int, uint)
v.Lsh(n uint) Value           // << (int, uint)
v.Rsh(n uint) Value           // >> (int, uint)
```

All operate directly on the `num` field for `KindInt` and `KindUint`. The shift amount is always `uint` (matching Go's specification).

---

## Type Conversions

```go
v.ToInt() Value     // → KindInt (from bool, int, uint, float)
v.ToUint() Value    // → KindUint (from bool, int, uint, float)
v.ToFloat() Value   // → KindFloat (from bool, int, uint, float)
v.ToBool() Value    // → KindBool (from bool, int, uint, float, string, nil)
v.ToString() Value  // → KindString (uses fmt.Sprintf("%v", v.Interface()))
```

These correspond to Go's type conversion expressions like `int(x)`, `float64(x)`, etc. They operate within the tagged-union — no allocation for numeric conversions.

`ToBool` has extended semantics: it returns `false` for zero values, empty strings, and nil; `true` otherwise.

### Reflect-Level Conversion

```go
v.ToReflectValue(typ reflect.Type) reflect.Value
```

This is the bridge from `Value` to `reflect.Value`, used when the VM needs to pass values to external functions or store them in typed containers. It handles:

- **Primitives**: `reflect.ValueOf(v.Int()).Convert(typ)` — adapts to the target type (e.g., `int64` → `int32`)
- **Closures** (`KindFunc`): wraps the closure in `reflect.MakeFunc(typ, ...)` so it can be stored in typed function containers
- **Native int slices** (`KindSlice` with `[]int64`): converts element-by-element to the target slice type
- **Native value slices** (`KindSlice` with `[]Value`): converts each element recursively
- **Reflect values**: returns the stored `reflect.Value` directly

---

## Container Operations

### Length and Capacity

```go
v.Len() int    // string, slice, array, map, chan
v.Cap() int    // slice, array, chan
```

Both include fast paths for native `[]int64` slices (direct `len()`/`cap()` call). For other types, they delegate to `reflect.Value.Len()`/`reflect.Value.Cap()`.

### Indexing

```go
v.Index(i int) Value          // read element at index i (string, slice, array)
v.SetIndex(i int, val Value)  // write element at index i (slice, array)
```

**Index fast paths**:
- `KindString`: returns `MakeUint(uint64(s[i]))` — byte at position, matching Go semantics.
- `KindSlice` with `[]int64`: returns `MakeInt(s[i])` — direct array access, zero reflection.
- `KindSlice` with `[]Value`: returns `slice[i]` — direct access.
- Other slices/arrays: `reflect.Value.Index(i)` → `MakeFromReflect`.

**SetIndex fast paths**:
- `KindSlice` with `[]int64`: `s[i] = val.RawInt()` — direct write.
- Other: `reflect.Value.Index(i).Set(val.ToReflectValue(elemType))`.

### Map Operations

```go
v.MapIndex(k Value) Value             // read value at key k
v.SetMapIndex(k, val Value)           // write value at key k (val.IsNil() → delete)
v.MapIter(f func(key, val Value) bool) // iterate over map entries
```

All map operations go through `reflect.Value` since Go's map type requires dynamic dispatch. `SetMapIndex` with a nil value deletes the key, matching Go's `delete(m, k)` semantics.

### Struct Field Access

```go
v.Field(i int) Value              // read struct field at index i
v.SetField(i int, val Value)      // write struct field at index i
```

Both require `reflect.Value` and use `reflect.Value.Field(i)`. The field index `i` is resolved at compile time from the struct type information.

---

## Pointer Operations

```go
v.Elem() Value           // dereference: *ptr → value
v.SetElem(val Value)     // set through pointer: *ptr = val
v.Pointer() uintptr      // raw pointer address (for identity comparison)
```

### Elem Fast Paths

`Elem()` has three fast paths before falling back to reflection:

1. **`*int64` pointer** (from `OpIndexAddr` on native int slices): `return MakeInt(*ptr)`
2. **`*Value` pointer** (from `OpAddr` on locals, `OpFree` for closures): `return *ptr`
3. **`*Value` inside `reflect.Value`** (pointer to a Value struct): unwraps directly

### SetElem Fast Paths

`SetElem()` mirrors `Elem()`:

1. **`*int64`**: `*ptr = val.num` — direct memory write
2. **`*Value`**: `*ptr = val` — direct struct assignment
3. **`reflect.Value` Ptr**: handles type conversion, function wrapping, and native slice conversion

### Helper

```go
UnsafeAddrOf(v reflect.Value) unsafe.Pointer
```

Used internally by the VM to obtain settable pointers to unexported struct fields.

---

## Channel Operations

```go
v.Send(val Value)                                    // ch <- val (blocking)
v.SendContext(ctx context.Context, val Value) error   // ch <- val (with cancellation)
v.TrySend(val Value) bool                            // non-blocking send
v.Recv() (Value, bool)                               // <-ch (blocking)
v.RecvContext(ctx context.Context) (Value, bool, error) // <-ch (with cancellation)
v.TryRecv() (Value, bool)                            // non-blocking receive
v.Close()                                            // close(ch)
```

All channel operations go through `reflect.Value` since Go channels require runtime support. The `*Context` variants use `reflect.Select` with a context cancellation channel — this is one area where reflection cannot be avoided.

### Context-Aware Pattern

`SendContext` and `RecvContext` implement a two-phase strategy:

1. **Fast path**: non-blocking try (`TrySend`/`TryRecv`). If the channel has buffer space or a waiting goroutine, this succeeds immediately.
2. **Slow path**: `reflect.Select` with two cases — the channel operation and `ctx.Done()`. This enables cancellation without leaking goroutines.

---

## Reflection Interop

The value system provides bidirectional conversion between `Value` and `reflect.Value`:

### Value → reflect.Value

```go
v.ToReflectValue(typ reflect.Type) reflect.Value
```

This is the primary conversion used when the VM needs to interact with Go's reflection system (external function calls, channel operations, struct field access). It handles all kinds, with special support for:

- **Type conversion**: `MakeInt(42).ToReflectValue(reflect.TypeOf(int32(0)))` produces `reflect.ValueOf(int32(42))`
- **Closure wrapping**: converts interpreter closures to real `func` values via `reflect.MakeFunc`
- **Slice conversion**: converts `[]int64` and `[]Value` to typed slices

### reflect.Value → Value

```go
MakeFromReflect(rv reflect.Value) Value
```

Extracts primitives from `reflect.Value` into the efficient tagged-union form. Composite types remain wrapped in `reflect.Value` under `KindReflect`.

### Interface Round-Trip

```go
FromInterface(v any) Value    // any → Value
v.Interface() any             // Value → any
```

These form the bridge between the interpreter and the host Go program. `FromInterface` is the entry point for all external data; `Interface()` is the exit point for returning results.

---

## Closure Support

Closures in Gig are represented as `KindFunc` values with the closure object stored in `obj`. The value system includes a **callback mechanism** to break the circular dependency between the `value` and `vm` packages:

```go
type ClosureCaller func(closure any, args []reflect.Value, outTypes []reflect.Type) []reflect.Value

RegisterClosureCaller(caller ClosureCaller)
```

The VM registers a `ClosureCaller` at initialization. When `ToReflectValue` needs to convert a closure to a typed `func` (e.g., for storing in a `map[string]func() int`), it uses `reflect.MakeFunc` with a handler that invokes the registered caller:

```go
fn := reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
    results := closureCaller(closure, args, outTypes)
    // Convert results to match expected return types
    ...
    return out
})
```

The `outTypes` parameter enables recursive wrapping of nested closures — a closure returning a closure can be properly typed at each level.

---

## Native Fast Paths

The value system provides specialized representations for common types to avoid reflection:

### `[]byte` (KindBytes)

```go
MakeBytes(b []byte) Value
v.Bytes() ([]byte, bool)
```

`[]byte` is the most common binary type (JSON payloads, protobuf, etc.). Storing it natively means external function calls with `[]byte` parameters avoid `reflect.ValueOf` entirely.

### `[]int64` (KindSlice with native storage)

```go
MakeIntSlice(s []int64) Value
v.IntSlice() ([]int64, bool)
```

Integer slices are extremely common in computational code. Native storage means:
- `v.Index(i)` → direct array access
- `v.SetIndex(i, val)` → direct array write
- `v.Len()` → `len(s)`
- The VM's `OpIntSliceGet`/`OpIntSliceSet` superinstructions operate directly on the `[]int64`

### `*int64` (KindPointer with native storage)

```go
MakeIntPtr(p *int64) Value
v.IntPtr() (*int64, bool)
```

Used by `OpIndexAddr` on native int slices — taking the address of `s[i]` returns a `*int64` pointer, and `OpDeref`/`OpSetDeref` on it are direct memory operations.

### `*Value` (KindPointer with native storage)

Closure free variables are `*Value` pointers. `Elem()` and `SetElem()` on these are direct struct dereferences.

### `[]Value` (KindSlice with native storage)

```go
MakeValueSlice(vals []Value) Value
v.ValueSlice() ([]Value, bool)
```

Used for multi-return value packing in DirectCall wrappers. A function returning `(int, error)` packs results as `MakeValueSlice([]Value{MakeInt(r0), FromInterface(r1)})`.

---

## Design Rationale

### Why Not `reflect.Value` Everywhere?

The alternative (used by Yaegi and other Go interpreters) is to represent every value as a `reflect.Value`. This is simpler but has significant costs:

| Operation | `reflect.Value` | Gig `Value` |
|-----------|-----------------|-------------|
| Create an int | heap alloc + `reflect.ValueOf` (~15-30 ns) | struct literal (0 ns, 0 alloc) |
| Integer addition | `reflect.Value.Int()` + `reflect.ValueOf` (~40 ns, 1 alloc) | `v.num + w.num` (~1 ns, 0 alloc) |
| Compare two ints | `reflect.Value.Int()` × 2 (~20 ns) | `v.num == w.num` (~0.5 ns) |
| Pass to external func | already `reflect.Value` (0 cost) | `ToReflectValue(typ)` (~5-15 ns) |

The tagged-union trades complexity at the external call boundary (where `ToReflectValue` is needed) for dramatic speedups in the hot paths (arithmetic, comparison, local variable access).

### Why 32 Bytes?

- Smaller (24 bytes) would require bit-packing the kind tag into `num`, complicating accessor logic.
- Larger (40+ bytes) would waste cache space — the VM stack holds hundreds of values, and cache locality is critical.
- 32 bytes is **2 cache line fractions** (64-byte cache lines hold exactly 2 values), giving good spatial locality for adjacent stack accesses.

### Why Both `KindSlice` and `KindReflect` for Slices?

`KindSlice` is used for native `[]int64` and `[]Value` storage. A generic `[]string` or `[]MyStruct` goes through `KindReflect`. The distinction enables the VM to take fast paths for integer slices (the most common case in computational code) without penalizing other slice types.

### Why `KindBytes` Is a Separate Kind

`[]byte` could be stored as `KindSlice` with a `[]byte` in `obj`, but having a dedicated kind enables:
- The `Bytes()` accessor with O(1) type check
- DirectCall wrappers to extract `[]byte` without `reflect.ValueOf`
- The `MakeBytes` constructor to be a pure struct assignment

---

## Performance Characteristics

### Zero-Allocation Operations

These operations create no heap allocations:

- All primitive constructors (`MakeInt`, `MakeBool`, `MakeFloat`, etc.)
- `MakeString`, `MakeBytes`, `MakeFunc` (reference types, header-only in `obj`)
- All arithmetic on primitives (`Add`, `Sub`, `Mul`, `Div`, `Mod`, `Neg`)
- All comparisons on primitives (`Equal`, `Cmp`)
- All bitwise operations
- All type conversions between numeric types
- Typed accessors (`Int()`, `Bool()`, `Float()`, `String()`)
- `Index` and `SetIndex` on native `[]int64`
- `Elem` and `SetElem` on native `*int64` and `*Value`

### Operations That May Allocate

- `FromInterface` for non-basic types (calls `reflect.ValueOf`)
- `Interface()` on `KindReflect` (calls `reflect.Value.Interface()`)
- `ToReflectValue` for closures (uses `reflect.MakeFunc`)
- All map operations (through `reflect.Value`)
- All channel operations (through `reflect.Value`)
- `Equal` on composite types (uses `reflect.DeepEqual`)

### Benchmarks

| Operation | Time | Allocations |
|-----------|------|-------------|
| `MakeInt(42)` | ~0.3 ns | 0 |
| `MakeBytes(b)` | ~0.3 ns | 0 |
| `FromInterface(42)` | ~2 ns | 0 |
| `FromInterface([]int{1,2,3})` | ~90 ns | 3 |
| `v.Int()` | ~0.5 ns | 0 |
| `v.Interface()` (KindInt) | ~1 ns | 0 |
| `v.Interface()` (KindReflect) | ~15 ns | 0-1 |
| `MakeInt(a).Add(MakeInt(b))` | ~2 ns | 0 |
| `v.ToReflectValue(intType)` | ~10 ns | 1 |

---

## File Organization

The value package is split into 6 files by responsibility:

| File | Purpose | Key Types/Functions |
|------|---------|-------------------|
| `value.go` | Core type definition, constructors, query methods | `Value`, `Kind`, `Make*`, `FromInterface`, `MakeFromReflect` |
| `accessor.go` | Typed accessors, `Interface()`, `ToReflectValue`, closure support | `Bool()`, `Int()`, `Interface()`, `ToReflectValue()`, `ClosureCaller` |
| `arithmetic.go` | Arithmetic, comparison, bitwise operations | `Add`, `Sub`, `Mul`, `Div`, `Cmp`, `Equal`, `And`, `Or`, `Lsh`, `Rsh` |
| `convert.go` | Type conversion methods | `ToInt`, `ToUint`, `ToFloat`, `ToBool`, `ToString` |
| `container.go` | Container operations (len, cap, index, map, field, pointer) | `Len`, `Cap`, `Index`, `SetIndex`, `MapIndex`, `Field`, `Elem`, `SetElem` |
| `channel.go` | Channel operations with context support | `Send`, `Recv`, `TrySend`, `TryRecv`, `SendContext`, `RecvContext`, `Close` |
| `value_test.go` | Unit tests | Covers constructors, arithmetic, comparison, conversion, edge cases |

---

## See Also

- [`docs/gig-internals.md`](gig-internals.md) — Architecture overview with the value system in context
- [`docs/optimization-zero-reflection.md`](optimization-zero-reflection.md) — How to use the value system for zero-reflection external calls
- [`docs/optimization-int-specialization.md`](optimization-int-specialization.md) — Integer specialization using `intLocals[]` shadow array
- [`docs/optimization-directcall.md`](optimization-directcall.md) — DirectCall code generation that extracts values via typed accessors
