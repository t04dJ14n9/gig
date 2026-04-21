# Gig Go Interpreter: Typed Nil Interface Bug Research Report

## Bug Summary
In Go, when a typed nil pointer (e.g., `*MyError = nil`) is assigned to an interface variable, the interface is **NOT nil** — it has a type (the pointer type) but the value is nil. The Gig interpreter incorrectly treats such interfaces as nil.

**Example:**
```go
var err *MyError  // nil pointer
var e error = err // typed nil assigned to interface
e == nil          // should be false in Go, but Gig returns true
```

In real Go: `e == nil` → `false` (interface has type `*MyError`, so it's not nil)
In Gig: `e == nil` → `true` (bug!)

---

## Key Research Findings

### 1. INTERFACE COMPARISON WITH NIL (OpEqual)

**Location:** `vm/run.go:437-448`

The OpEqual instruction uses `value.Equal()` for non-integer comparisons:

```go
case bytecode.OpEqual:
    sp--
    b := stack[sp]
    sp--
    a := stack[sp]
    if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
        stack[sp] = value.MakeBool(a.RawInt() == b.RawInt())
    } else {
        stack[sp] = value.MakeBool(a.Equal(b))
    }
    sp++
    continue
```

**Critical Issue:** The comparison falls through to `value.Equal()` for interface types.

---

### 2. THE EQUAL METHOD (Where the Bug Lives)

**Location:** `model/value/arithmetic.go:203-263`

The `Equal()` method handles nil comparison at lines 217-219:

```go
func (v Value) Equal(other Value) bool {
    // ...unwrap interface/reflect values for comparison...
    a, b := v, other
    if a.kind == KindReflect || a.kind == KindInterface {
        a = unwrapForComparison(a)
    }
    if b.kind == KindReflect || b.kind == KindInterface {
        b = unwrapForComparison(b)
    }
    
    if a.kind != b.kind {
        // Handle nil comparison
        if a.kind == KindNil || b.kind == KindNil {
            return a.IsNil() && b.IsNil()  // ← PROBLEM: Uses IsNil()
        }
        return false
    }
    // ...rest of switch...
}
```

**The Bug:** When comparing an interface to nil (KindNil), it calls `IsNil()` on both values.

---

### 3. THE IsNil() METHOD (Root Cause)

**Location:** `model/value/value.go:169-188`

```go
func (v Value) IsNil() bool {
    if v.kind == KindNil {
        return true
    }
    if v.kind == KindReflect {
        if rv, ok := v.obj.(reflect.Value); ok {
            if !rv.IsValid() {
                return true
            }
            // Only call IsNil on types that support it
            switch rv.Kind() {
            case reflect.Chan, reflect.Func, reflect.Interface, 
                 reflect.Map, reflect.Ptr, reflect.Slice:
                return rv.IsNil()
            }
            return false
        }
    }
    return false
}
```

**Critical Problem:** When an interface value (`KindInterface`) is compared to nil:
- The Value is stored as a `KindReflect` with a `reflect.Value` inside (from line 118 in `typeconv.go`: all interfaces become `any` in reflect)
- `IsNil()` calls `rv.IsNil()` on the `reflect.Value`
- For a typed nil pointer (e.g., `*MyError = nil`), the `reflect.Value` contains a pointer type pointing to nil
- Go's `reflect.Value.IsNil()` on this returns `true`

**Missing:** There's NO special handling for `KindInterface` type in `IsNil()`. It only checks `KindReflect`.

---

### 4. HOW TYPED NIL IS REPRESENTED

**Location:** `model/value/value.go:415-463` (FromInterface and MakeFromReflect)

When a typed nil pointer is converted to an interface value:
```go
// In FromInterface (line 417-463)
// For nil constants, this is called from compile_value.go:135-137:
if cnst.Value == nil {
    if rt := constTypeToReflect(cnst.Type()); rt != nil {
        v = reflect.Zero(rt)  // Creates reflect.Value with typed nil
    }
}
```

The typed nil pointer becomes a `reflect.Value` (with kind = `reflect.Ptr`, value = nil pointer).

This `reflect.Value` is wrapped in a Gig `Value` with:
- `kind = KindReflect` 
- `obj = reflect.Value{Type: *MyError, Value: nil}`

---

### 5. INTERFACE COMPILATION

**Location:** `compiler/compile_value.go:417-422`

The MakeInterface compilation is simple:
```go
func (c *compiler) compileMakeInterface(i *ssa.MakeInterface) {
    resultIdx := c.symbolTable.AllocLocal(i)
    c.compileValue(i.X)  // Compile the value being assigned to interface
    c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
```

It just compiles the value as-is. No special handling.

**Location:** `compiler/compile_value.go:579-582` (ChangeInterface)

```go
func (c *compiler) compileChangeInterface(i *ssa.ChangeInterface) {
    resultIdx := c.symbolTable.AllocLocal(i)
    c.compileValue(i.X)  // Simply passes through
    c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
```

Again, just passes through the value.

---

### 6. HOW INTERFACES ARE CONVERTED TO reflect.Type

**Location:** `vm/typeconv.go:115-118`

```go
case *types.Interface:
    // Interface type — use the empty interface (any) type
    // For the VM, all interfaces are represented as any
    return reflect.TypeFor[any]()
```

**Key Design Decision:** All interface types (including `error`, custom interfaces, `interface{}`) are converted to `any` in the reflect type system. This loses the concrete type information that would distinguish a typed nil interface from a nil interface.

---

### 7. VALUE STORAGE FOR INTERFACES

**Location:** `model/value/value.go:145-156` (Value struct definition)

```go
type Value struct {
    kind Kind      // 1 byte
    size Size      // 1 byte (in padding)
    num  int64     // 8 bytes (bool, int, uint, float bits)
    obj  any       // 16 bytes (stores complex types)
}
```

For interface values containing a typed nil pointer:
- `kind = KindReflect` (or possibly `KindInterface`)
- `obj = reflect.Value{Type: *MyError, Value: nil}`

There's **NO separate field to track that this is a "typed nil interface"**. The distinction is lost when unwrapping.

---

### 8. THE unwrapForComparison FUNCTION (Makes Things Worse)

**Location:** `model/value/arithmetic.go:188-201`

```go
func unwrapForComparison(v Value) Value {
    rv, ok := v.obj.(reflect.Value)
    if !ok {
        return v
    }
    // Unwrap interface
    if rv.Kind() == reflect.Interface && !rv.IsNil() {
        rv = rv.Elem()
    }
    return MakeFromReflect(rv)
}
```

**Problem:** This unwraps the interface to get the concrete type, but in doing so:
- If we have a `reflect.Value` representing a nil pointer (not a nil interface), we lose the type information
- We then call `IsNil()` on the unwrapped value, which only checks the value, not the type

---

## Summary of Root Causes

1. **Missing KindInterface handling in IsNil()**: The `IsNil()` method doesn't properly handle `KindInterface` values. It only checks `KindReflect`.

2. **Conflation of "nil interface" with "typed nil"**: The code doesn't distinguish between:
   - A truly nil interface (no type, no value)
   - A typed nil interface (has type, value is nil)

3. **Use of any/reflect.Value loses type information**: All Go interfaces are converted to `reflect.TypeFor[any]()`, which doesn't preserve whether the interface contained a typed nil or was actually nil.

4. **IsNil() calls reflect.Value.IsNil() without type checking**: For a reflect.Value containing a nil pointer type, `rv.IsNil()` returns true, but in Go semantics, an interface{} holding a nil pointer is NOT nil.

5. **No explicit representation for "typed nil in interface"**: The Value type has no way to distinguish "nil interface" from "interface holding typed nil pointer".

---

## Files Involved

### Core Runtime Files
- **`vm/run.go`** (lines 437-448): OpEqual instruction uses value.Equal()
- **`vm/ops_dispatch.go`**: Dispatches opcodes to handlers
- **`vm/typeconv.go`** (lines 115-118): All interfaces → `any` reflect type

### Value Representation
- **`model/value/value.go`** 
  - Lines 169-188: IsNil() method - CRITICAL BUG LOCATION
  - Lines 415-463: FromInterface() and MakeFromReflect()
  - Lines 145-156: Value struct definition

- **`model/value/arithmetic.go`**
  - Lines 203-263: Equal() method - uses IsNil() incorrectly
  - Lines 188-201: unwrapForComparison() - loses type info

### Compilation
- **`compiler/compile_value.go`**
  - Lines 417-422: compileMakeInterface()
  - Lines 579-582: compileChangeInterface()
  - Lines 112-151: compileConst() - creates typed nil constants

### Type Assertion
- **`vm/ops_convert.go`** (lines 36-62): OpAssert handles type assertions with interface

---

## Why Go's Semantics Are Different

In real Go, interfaces are represented as a (type, value) pair:
```
interface{} = (type: *MyError, value: nil_pointer)
interface{} = nil  // (type: nil, value: nil)

// These are NOT equal:
var p *MyError = nil
var e interface{} = p
e == nil  // false — because (type: *MyError, value: nil) != (type: nil, value: nil)
```

Gig loses the type information when storing interfaces as `reflect.Value` with all interfaces mapped to `any`, making it impossible to distinguish typed nil from nil interfaces.

