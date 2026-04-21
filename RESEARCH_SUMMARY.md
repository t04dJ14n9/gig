# TYPED NIL INTERFACE BUG - RESEARCH SUMMARY

## Quick Reference

**Bug:** Typed nil pointer assigned to interface incorrectly treated as nil  
**Test Case:** `tests/testdata/known_issues/main.go` - `TypedNilInterface()` function  
**Test File:** `tests/known_issue_test.go` - line 88-92  

## Files Requiring Review

### CRITICAL - These are where the bug manifests:
1. **`model/value/value.go` (Lines 169-188)** - `IsNil()` method
   - IGNORES that typed nil could come from interface
   - Calls `reflect.Value.IsNil()` without context
   - NO handling for `KindInterface` type

2. **`model/value/arithmetic.go` (Lines 203-263)** - `Equal()` method  
   - Uses `IsNil()` to check nil comparison
   - Calls `unwrapForComparison()` which loses type info
   - No special handling for interface nil semantics

### IMPORTANT - These provide context/affected:
3. **`vm/run.go` (Lines 437-448)** - `OpEqual` instruction
   - Calls `value.Equal()` for interface comparisons
   - This is the runtime bytecode execution

4. **`vm/ops_dispatch.go`** - Dispatches opcodes
   - Routes comparison operations

5. **`vm/typeconv.go` (Lines 115-118)**
   - All interfaces mapped to `any` type
   - TYPE INFORMATION LOST HERE

### COMPILATION - Shows how typed nil is created:
6. **`compiler/compile_value.go`** (Lines 417-422, 579-582)
   - `compileMakeInterface()` - just passes through
   - `compileChangeInterface()` - just passes through
   - Lines 112-151: `compileConst()` creates typed nil constants

## Root Cause Analysis

### Problem 1: Conflation of Two Different Concepts
```
"nil interface" = (Type: nil, Value: nil)         → interface IS nil
"typed nil"     = (Type: *MyError, Value: nil)    → interface is NOT nil

Gig treats both as "nil" because IsNil() only checks the VALUE,
not whether a TYPE was ever assigned.
```

### Problem 2: Type Information Lost Early
```
typeconv.go:115-118 converts ALL interfaces to reflect.TypeFor[any]()
This loses:
- That this is specifically an "error" interface
- What concrete type was assigned to it
- Whether nil was ever assigned vs a typed nil
```

### Problem 3: No Context in IsNil()
```
value.go:170-188 has NO WAY to know:
- "Did this Value come from an interface assignment?"
- "Does this reflect.Value represent a nil interface or just a nil pointer?"
- "Was a type ever associated with this value?"

It just checks: if rv.IsNil() { return true }
That's semantically WRONG for interface values!
```

### Problem 4: Interface Representation Inadequate
```
The Value struct has these fields:
- kind Kind         (8-bit tag)
- size Size         (8-bit, used for numeric type width)
- num int64         (64-bit for primitives)
- obj any           (stores complex types as interface{})

For interface values:
- kind = KindReflect (could be KindInterface, but still doesn't help)
- obj = reflect.Value{Type: *MyError, Value: nil}

NO FIELD to distinguish:
- "This is an interface value" from "This is a pointer value"
- "This came from `var e error = p`" vs "This came from `var p *MyError = nil`"
```

## Key Code Locations

### Where Typed Nil is Created (at compile time)
**File:** `compiler/compile_value.go` (lines 112-151)
```go
// When compiler sees: var e error = (*MyError)(nil)
if cnst.Value == nil {
    if rt := constTypeToReflect(cnst.Type()); rt != nil {
        v = reflect.Zero(rt)  // Creates typed nil via reflect
    }
}
// This reflect.Value is wrapped in Value{kind: KindReflect, obj: v}
```

### Where Typed Nil is Compared (at runtime)
**File:** `vm/run.go` (lines 437-448)
```go
case bytecode.OpEqual:
    // ... when comparing e == nil ...
    stack[sp] = value.MakeBool(a.Equal(b))
```

**File:** `model/value/arithmetic.go` (lines 204-263)
```go
func (v Value) Equal(other Value) bool {
    // ... unwraps interface ...
    if a.kind == KindNil || b.kind == KindNil {
        return a.IsNil() && b.IsNil()  // ← CALLS IsNil()
    }
}
```

**File:** `model/value/value.go` (lines 169-188)
```go
func (v Value) IsNil() bool {
    if v.kind == KindNil { return true }
    if v.kind == KindReflect {
        if rv, ok := v.obj.(reflect.Value); ok {
            // ...
            switch rv.Kind() {
            case reflect.Ptr:
                return rv.IsNil()  // ← THIS IS WRONG FOR INTERFACES!
                                   // Doesn't know it came from interface
            }
        }
    }
    return false
}
```

### Where Interface Type Becomes `any`
**File:** `vm/typeconv.go` (lines 115-118)
```go
case *types.Interface:
    // Interface type — use the empty interface (any) type
    // For the VM, all interfaces are represented as any
    return reflect.TypeFor[any]()  // ← TYPE INFO LOST
```

## Why This Matters

### Semantic Difference in Go
```
var p *MyError = nil
var e error = p

In Go, checking "e == nil" checks if BOTH:
1. The interface has a type (it does: *MyError) AND
2. The value is nil (it is)

Since the interface has a non-nil type, the comparison is FALSE.
This is critical for error handling:

if err == nil { /* success */ }

If err holds a typed nil, it should enter the error branch!
```

### Related Test Cases
```
tests/testdata/known_issues/main.go:
  - TypedNilInterface() - should return "not nil", returns "nil"
  - InterfaceMethodOnTypedNil() - should panic, returns "ok"
```

## Solution Approaches (Not Implemented - For Reference)

### Approach 1: Add Flag to Value
Add a field to `Value` struct to track interface typed-nil status.
Would require:
- Modifying `Value` struct (adds 1 byte)
- Updating all value constructors
- Checking flag in `IsNil()`

### Approach 2: Preserve Concrete Type
Store both the declared interface type AND the concrete type:
Would require:
- New fields in `Value` struct
- Tracking through assignment
- Checking concrete type in `IsNil()`

### Approach 3: Separate InterfaceValue Type
Create a dedicated type for interface values:
```go
type InterfaceValue struct {
    InterfaceType reflect.Type  // e.g., error (the interface)
    ConcreteType  reflect.Type  // e.g., *MyError (what was assigned)
    Value         reflect.Value // the actual value
}
```
Would require:
- New type definition
- Changes to value constructors
- Changes to comparison logic
- Changes to compilation

## Testing

**Test File:** `tests/known_issue_test.go`
**Test Data:** `tests/testdata/known_issues/main.go`

```go
func TestKnownIssues(t *testing.T) {
    issues := map[string]KnownIssue{
        "TypedNilInterface": {
            funcName: "TypedNilInterface",
            native:   func() any { return known_issues.TypedNilInterface() },
            issue:    "typed nil pointer assigned to interface is incorrectly treated as nil",
        },
        // ...
    }
}
```

Run with:
```bash
go test -run TestKnownIssues/TypedNilInterface
```

## How to Verify Bug

Create test file:
```go
package main

type MyError struct{}
func (e *MyError) Error() string { return "" }

func main() {
    var p *MyError = nil
    var e error = p
    
    // In Go: false (typed nil is NOT nil)
    // In Gig: true (BUG!)
    println(e == nil)
}
```

Run in Gig - will print `true` (wrong!)  
Run in Go - will print `false` (correct!)

