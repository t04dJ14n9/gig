# Typed Nil Interface Bug - Visual Flow Diagrams

## 1. THE BUG IN ACTION

### Real Go Behavior
```
var err *MyError = nil     // nil pointer to MyError
var e error = err          // assign to interface

Go represents e as:        (type: *MyError, value: nil_ptr)
e == nil                   → false ✓ (different from (nil, nil))
```

### Gig's Incorrect Behavior
```
var err *MyError = nil     // nil pointer to MyError
var e error = err          // assign to interface

Gig represents e as:       Value{kind: KindReflect, obj: reflect.Value{...}}
OpEqual compares e to nil:
  → calls value.Equal()
    → unwrapForComparison()  [loses type info!]
    → calls IsNil()
      → reflect.Value.IsNil() returns true
  → returns true ✗ WRONG!
```

---

## 2. CODE FLOW DURING COMPARISON (Where the Bug Happens)

```
EXECUTION: e == nil

┌─ OpEqual (vm/run.go:437)
│  ├─ stack[sp-1] = e (interface with typed nil)
│  ├─ stack[sp] = nil (KindNil)
│  └─ calls value.Equal(e, nil)
│
├─ arithmetic.go:204 Equal()
│  ├─ a = e (KindReflect)
│  ├─ b = nil (KindNil)
│  ├─ CALLS unwrapForComparison(a)
│  │  ├─ Extracts reflect.Value from a.obj
│  │  └─ Returns MakeFromReflect(rv)
│  │      └─ PROBLEM: Loses that this was a typed nil!
│  ├─ a.kind != b.kind (KindReflect != KindNil)
│  ├─ Enters nil handling block (line 217)
│  └─ RETURNS a.IsNil() && b.IsNil()
│
├─ value.go:170 IsNil() for a (KindReflect)
│  ├─ v.kind = KindReflect
│  ├─ rv = reflect.Value{Type: *MyError, IsNil: true}
│  ├─ rv.Kind() == reflect.Ptr
│  ├─ CALLS rv.IsNil()
│  │  └─ Go's reflect returns: true
│  │     (because the pointer value itself is nil)
│  └─ Returns true ✗ WRONG!
│
└─ Returns (true && true) = true ✗ BUG RESULT
```

---

## 3. THE CORE PROBLEM: IsNil() METHOD

### Current Implementation (WRONG)
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
            switch rv.Kind() {
            case reflect.Ptr:  // ← PROBLEM HERE!
                return rv.IsNil()  // Only checks if pointer value is nil
                                   // Doesn't know it came from interface!
            }
        }
    }
    return false
}
```

### What Should Happen (CORRECT Go Semantics)
```
Interface value contains (*MyError, nil_pointer):
- IsNil() should return FALSE
- Because interface has a type (*MyError), even though the value is nil

But current code:
- Calls reflect.Value.IsNil() on just the pointer
- Returns TRUE
- Loses the fact this came from an interface!
```

---

## 4. HOW TYPED NIL IS CREATED

```
COMPILATION: var e error = (*MyError)(nil)

1. compileConst() [compile_value.go:112]
   └─ Detects nil constant with type *MyError

2. Calls constTypeToReflect() to get reflect.Type of *MyError
   └─ Creates reflect.Zero(*MyError) = reflect.Value{Type: *MyError, Value: nil}

3. This is added to program Constants

4. OpConst loads it at runtime

5. Wrapped as Value{kind: KindReflect, obj: reflect.Value{...}}


PROBLEM:
At this point, we have a reflect.Value that "looks like" a nil pointer
But this was supposed to be an interface value!
The fact that it came from an interface assignment is LOST.
```

---

## 5. WHERE TYPE INFO IS LOST

```
TypeAssignment: var e error = (*MyError)(nil)
                              ↓
Interface Type: error → converted to any (typeconv.go:115-118)
                 ↓
Reflect Type: reflect.TypeFor[any]() 
              (loses that this is specifically an error interface)
                 ↓
Value Wrapping: Value{kind: KindReflect, obj: rv}
                (has no field to mark "this is a typed nil in interface")
                 ↓
Comparison: e == nil
            ↓
Result: IsNil() incorrectly returns true
```

---

## 6. FILES AND LINE NUMBERS - EXECUTION PATH

```
User Code: var e error = err; e == nil
                         ↓
                 [COMPILATION]
                         ↓
compiler/compile_value.go:579-582
    compileChangeInterface(i *ssa.ChangeInterface)
        ↓ emits OpSetLocal
                         ↓
                    [RUNTIME]
                         ↓
vm/run.go:437-448
    OpEqual handler
        calls value.Equal(a, b)
                         ↓
model/value/arithmetic.go:204-263
    Equal() method
        calls unwrapForComparison()
        calls IsNil() ← BUG HERE!
                         ↓
model/value/value.go:169-188
    IsNil() method
        calls reflect.Value.IsNil() ← WRONG SEMANTICS
        returns true ✗
```

---

## 7. INTERFACE REPRESENTATION COMPARISON

### Go's Representation (Correct)
```
Interface{} = (Type, Value)

Example 1: var e error = (*MyError)(nil)
  e = (Type: *MyError, Value: nil_pointer)
  e == nil → FALSE (not nil interface!)

Example 2: var e error = nil
  e = (Type: nil, Value: nil)
  e == nil → TRUE (nil interface)
```

### Gig's Representation (Broken)
```
Interface{} = Value{kind, obj}

Example 1: var e error = (*MyError)(nil)
  e = Value{
    kind: KindReflect,
    obj: reflect.Value{
      Type: *MyError,    ← Type info exists in reflect.Value
      Value: nil_pointer
    }
  }
  e == nil → TRUE ✗ (should be FALSE!)

Example 2: var e error = nil
  e = Value{
    kind: KindNil,
    obj: nil
  }
  e == nil → TRUE ✓ (correct by accident)
```

**Problem:** Case 1 and Case 2 are indistinguishable in IsNil()!
- Both end up calling reflect.Value.IsNil() which returns true
- But in Go, they should be different!

---

## 8. KEY DESIGN FLAW

### The Assumption (Wrong)
```
"If reflect.Value.IsNil() returns true, then the interface is nil"
```

### The Reality (Correct)
```
"If reflect.Value.IsNil() returns true, AND the Value came from a nil interface
 (not from an interface assignment of a typed nil), THEN the interface is nil"
```

The code has no way to distinguish these two cases because:
1. All interfaces are mapped to `any` type
2. Type information is not tracked separately
3. IsNil() only checks the reflect.Value, not where it came from

---

## 9. THE SOLUTION (High Level)

### Option A: Track Interface Origin
```
Add field to Value: isInterfaceTypedNil bool

When assigning typed nil to interface:
  Value{
    kind: KindReflect,
    obj: reflect.Value{Type: *MyError, Value: nil},
    isInterfaceTypedNil: true  ← NEW FIELD
  }

IsNil() checks this flag:
  if isInterfaceTypedNil { return false }  ← Can't be nil!
```

### Option B: Store Interface Type Metadata
```
Preserve both the interface type AND the value type:
  Value{
    kind: KindInterface,
    interfaceType: *types.Interface,     ← NEW
    concreteType: *MyError,              ← NEW
    obj: reflect.Value{...}
  }

IsNil() checks if concreteType is non-nil:
  if concreteType != nil { return false }  ← Has type, can't be nil!
```

### Option C: Create InterfaceValue Wrapper
```
Create separate type for interface values:
  type InterfaceValue struct {
    Type  reflect.Type   // e.g., *MyError
    Value reflect.Value  // the nil pointer
  }

Value{
  kind: KindInterface,
  obj: InterfaceValue{...}
}

IsNil() checks InterfaceValue.Type:
  if iface.Type != nil { return false }  ← Has type!
```

