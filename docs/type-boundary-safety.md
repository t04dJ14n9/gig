# Type Boundary Safety: The `reflect.StructOf` Limitation

## The Problem in One Sentence

Go's `reflect.StructOf` can create struct **fields** at runtime, but cannot attach
**methods** — and without methods, a type cannot satisfy any Go interface. This is
the single constraint that shapes gig's entire external-call architecture.

---

## 1. Why `reflect.StructOf` Cannot Attach Methods

Go's reflection API draws a hard line between data and behavior:

```
                    Fields     Methods
                    ──────     ───────
compile-time type     ✓           ✓
reflect.StructOf      ✓           ✗    ← methods cannot be added at runtime
```

There is no `reflect.MethodOf`, no `reflect.NewMethod`, and no way to mutate a
`reflect.Type`'s method set after creation. The `StructOf` function accepts only
`[]reflect.StructField` — there is no parameter for methods.

A minimal demonstration:

```go
package main

import (
    "fmt"
    "reflect"
)

type sortInterface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}

func main() {
    // Create a struct type at runtime. We can give it fields and tags...
    dynamicType := reflect.StructOf([]reflect.StructField{
        {Name: "Data", Type: reflect.TypeOf([]int{}), Tag: `gig:"#main.MySlice"`},
    })

    // ...but NumMethod() is always zero.
    fmt.Printf("NumMethod: %d\n", dynamicType.NumMethod())
    // Output: NumMethod: 0

    // This means it can never satisfy any interface:
    iface := reflect.TypeOf((*sortInterface)(nil)).Elem()
    fmt.Println(dynamicType.Implements(iface))
    // Output: false
}
```

### Why Go Made This Choice

Methods in Go are compiled into the binary as function symbols. A method call
`t.M()` compiles to a direct function call with a specific PC offset — the linker
resolves it. Adding methods at runtime would require mutable dispatch tables (like
C++ vtables), which Go deliberately avoids for simplicity, performance, and the
guarantee that all types are fully known at compile time.

The `reflect` package can describe any type and create values of any type, but it
cannot create new executable code. Method bodies are code, and code generation at
runtime is outside the scope of the `reflect` package.

---

## 2. How This Affects gig

### 2.1 The Interpreted-Type Wrapping Problem

When gig interprets user code like:

```go
// User's script (interpreted)
type ByLength []string

func (s ByLength) Len() int           { return len(s) }
func (s ByLength) Less(i, j int) bool { return len(s[i]) < len(s[j]) }
func (s ByLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func main() {
    sort.Sort(ByLength(words))
}
```

gig must create a runtime representation of `ByLength` and pass it to `sort.Sort`.
The path the value takes is:

```
User script              VM                           Go (stdlib)
───────────              ──                           ───────────
ByLength{"a","bb"}  →  Value{kind: KindReflect,   →  sort.Sort() receives
                         obj: reflect.ValueOf(        reflect.StructOf{...}
                           reflect.StructOf{...})}    with NO methods
                                                      |
                                                      v
                                              sort.Sort does:
                                                data, ok := v.(sort.Interface)
                                                // ok == false! PANIC.
```

`sort.Sort`'s first line is a type assertion to `sort.Interface`. The
`reflect.StructOf` value has zero methods, so the assertion fails, and
`sort.Sort` panics with:

```
interface conversion: struct{...} is not sort.Interface: missing method Len
```

### 2.2 The Third-Party Library Problem (General Case)

This is not limited to `sort`. **Any** third-party library that inspects types at
the Go boundary can break. Common patterns:

| Pattern | Example | What Breaks |
|---------|---------|-------------|
| Interface assertion | `sort.Sort`, `heap.Init` | `v.(Iface)` panics |
| Reflection on methods | `reflect.MethodByName` | Returns zero `Value` |
| Interface check | `reflect.Type.Implements` | Returns `false` |
| Type switch on concrete type | ORM field mapping | Never matches |
| Generic constraint satisfaction | `json.Marshal` with custom marshaler | Falls through to default |

**The fundamental issue**: gig's runtime-synthesized types are structurally
correct (right fields, right values) but behaviorally invisible (no methods).
Go's type system requires both.

### 2.3 Why Reflection-Based Libraries Are Especially Vulnerable

Libraries that use reflection heavily — ORMs, serializers, DI containers, mock
frameworks — inspect type identity and method sets. They assume types have a
compile-time origin. gig's runtime-synthesized types violate this assumption.

Example: a popular validation library might do:

```go
func Validate(v any) error {
    val := reflect.ValueOf(v)
    if val.Type().Implements(reflect.TypeOf((*Validator)(nil)).Elem()) {
        return v.(Validator).Validate()  // ← never reached for gig types
    }
    // fall through to struct field validation
}
```

The `Implements` check returns `false` because the gig-synthesized struct has
no methods. The library silently skips custom validation logic.

---

## 3. gig's Solution: Defense in Depth

gig handles this at three layers, each backing up the one below it:

```
┌─────────────────────────────────────────────────────────┐
│ Layer 1: Compile-Time Gate                              │
│ Reject user-defined types → third-party functions.      │
│ Stdlib calls are allowed (gig guarantees correctness).  │
├─────────────────────────────────────────────────────────┤
│ Layer 2: Interface Adapters (stdlib only)               │
│ For sort.Interface and heap.Interface, wrap interpreted │
│ types in a Go-native adapter that dispatches callbacks  │
│ back into the VM.                                       │
├─────────────────────────────────────────────────────────┤
│ Layer 3: Escape Hatch                                   │
│ WithAllowUnsafeTypePass() disables the compile-time     │
│ check for users who accept the risk.                    │
└─────────────────────────────────────────────────────────┘
```

### 3.1 Layer 1: Compile-Time Gate

At compile time, gig inspects every external function call (`OpCallExternal`). If:

1. The target package is **third-party** (import path contains a dot, e.g.,
   `github.com/foo/bar`), AND
2. Any argument type is **user-defined** (declared in the script's `main` or
   `command-line-arguments` package)

...then compilation fails with an error:

```
cannot pass interpreter-defined type "ByLength" to third-party function
github.com/foo/bar.Process (argument 1): custom types are not compatible
with external libraries that use reflection. Use primitive types, slices,
maps, or types from registered packages instead.
```

**How stdlib vs. third-party is determined**:

```go
func isStdlibPath(path string) bool {
    // Stdlib paths have no dot in the first segment:
    //   "fmt"           → stdlib
    //   "encoding/json" → stdlib  (first segment "encoding", no dot)
    //   "sort"          → stdlib
    //
    // Third-party paths have a dot:
    //   "github.com/foo"    → third-party (first segment "github.com")
    //   "golang.org/x/tools" → third-party (first segment "golang.org")
    firstSlash := strings.IndexByte(path, '/')
    firstSegment := path
    if firstSlash >= 0 {
        firstSegment = path[:firstSlash]
    }
    return !strings.ContainsRune(firstSegment, '.')
}
```

**What's always allowed**:
- Primitives (`int`, `string`, `float64`, `bool`, etc.) → any function
- Slices, maps of primitives → any function
- Types from registered packages (`sort.IntSlice`, `time.Time`) → any function
- User-defined types → stdlib functions (gig provides adapters)

**What's rejected**:
- User-defined named types (`type MyStruct struct{...}`) → third-party functions
- Pointers to user-defined types → third-party functions
- Slices/maps of user-defined types → third-party functions

### 3.2 Layer 2: Interface Adapters (stdlib)

For stdlib interfaces whose implementations are callback-driven (`sort.Interface`,
`heap.Interface`), gig creates a Go-native adapter at the VM boundary:

```
                    ┌──────────────────────────────┐
                    │ interpretedInterfaceAdapter  │
                    │                              │
sort.Sort(adapter)  │  Len() int  ──→ call("Len")  │──→ VM executes
       │            │  Less(i,j)  ──→ call("Less") │    compiled method
       ▼            │  Swap(i,j)  ──→ call("Swap") │    via temp VM
   Len() called ────│  Push(x)    ──→ call("Push") │
                    │  Pop() any  ──→ call("Pop")  │
                    └──────────────────────────────┘
```

The adapter satisfies Go's `sort.Interface` at the Go type level (it's a
statically-compiled Go struct with real methods). But each method call dispatches
back into the interpreter to execute the user's compiled method body. This means:

- **`sort.Sort` sees a valid `sort.Interface`** — the type assertion passes.
- **User's method bodies execute** — `Len()`, `Less(i,j)`, `Swap(i,j)` run as
  compiled bytecode in a temporary VM.
- **User's state is visible** — the adapter threads the caller's globals,
  context, and goroutine tracker through to the callback VM.

The adapter is created at the `OpConvertInterface` instruction when the compiler
emits a conversion from a user-defined concrete type to `sort.Interface` or
`container/heap.Interface`.

**Key design decisions**:

1. **Exact type matching only.** The adapter is created only for `sort.Interface`
   and `container/heap.Interface` — NOT for any interface that happens to have
   `Len`/`Less`/`Swap` methods. This prevents the adapter from shadowing
   user-defined interfaces that happen to have the same shape.

2. **Swap tries compiled method first, falls back to direct slice swap.** This
   handles types with auxiliary state (counters, parallel arrays) that must be
   updated by `Swap`.

3. **Value-receiver vs. pointer-receiver disambiguation.** For `heap.Interface`,
   `Len`/`Less`/`Swap` are typically value-receiver methods on a named slice,
   while `Push`/`Pop` are pointer-receiver methods. The adapter dereferences
   the receiver for value-receiver methods and keeps the pointer for
   pointer-receiver methods.

### 3.3 Layer 3: Escape Hatch

```go
prog, err := gig.Build(source, gig.WithAllowUnsafeTypePass())
```

This disables the compile-time check. Use cases:
- The third-party library does NOT use reflection on argument types
- The user has verified type compatibility manually
- Prototyping/MVP where the risk is acceptable

---

## 4. Why Not Fix It at the Runtime Level?

Several approaches were considered and rejected:

### 4.1 Make `reflect.StructOf` Types Implement Interfaces (Impossible)

Go's runtime does not support adding methods to dynamically-created types. This
is a language-level constraint, not something gig can work around without
modifying the Go runtime itself (via a patched Go compiler or a plugin that hooks
into the runtime type system).

### 4.2 Generate Go Code and Compile It (Too Slow)

gig could generate `.go` files with wrapper types, compile them with `go build`,
and load them as plugins. This would add seconds of latency per `gig.Build()`
call, defeating the purpose of an embedded interpreter.

### 4.3 cgo/Assembly Method Injection (Too Fragile)

The internal layout of `runtime._type` and method tables is not part of Go's
public API. Patching them via unsafe pointer manipulation would break on every
Go release and is not viable for a library used in production.

### 4.4 Require Users to Pre-register All Types (Too Restrictive)

If users had to define every possible struct and interface in a `.go` file compiled
into the host binary, gig would lose its dynamic-code-loading value proposition.

---

## 5. What This Means for Users

### Safe Patterns (Always Work)

```go
// Pass primitives to any function
fmt.Sprintf("%d", 42)

// Pass slices/maps of primitives to any function
json.Marshal([]string{"a", "b"})

// Use custom types with stdlib functions
type MySorter []int
// ... implement sort.Interface ...
sort.Sort(MySorter{3, 1, 2})  // ← gig provides adapter

// Use types from registered packages with any function
t := time.Now()
formatted := t.Format("2006-01-02")
```

### Unsafe Patterns (Rejected at Compile Time)

```go
// User-defined struct passed to third-party ORM
type User struct { Name string; Age int }
db.Insert(User{"Alice", 30})  // ← compile error

// User-defined interface passed to third-party DI container
type Validator interface { Validate() error }
container.Register(ValidatorImpl{})  // ← compile error
```

### Escape Hatch

```go
// Use WithAllowUnsafeTypePass() if you know the library
// doesn't inspect types via reflection:
prog, err := gig.Build(source, gig.WithAllowUnsafeTypePass())
```

---

## 6. Related Documents

- [Gig Internals](gig-internals.md) — comprehensive architecture guide
- [External Call Architecture](external-call-architecture.md) — proposal for a
  unified boundary conversion approach (future direction)
- [CLAUDE.md](../CLAUDE.md) — project overview and commands

---

## 7. References

- Go issue [#16522](https://github.com/golang/go/issues/16522) — proposal for
  `reflect.MethodOf` (declined)
- Go issue [#4146](https://github.com/golang/go/issues/4146) — runtime type
  mutation discussion
- `reflect.StructOf` [documentation](https://pkg.go.dev/reflect#StructOf) —
  creates struct types with fields only, no methods
- `reflect.MakeFunc` [documentation](https://pkg.go.dev/reflect#MakeFunc) —
  creates function values, but cannot be attached as methods to a type
