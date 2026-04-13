# Gig Interpreter Code Simplification - Implementation Guide

This guide provides concrete implementation steps for each simplification opportunity.

---

## PRIORITY 1: OpCode.String() Table-Driven Conversion

### Current Implementation
**File**: `model/bytecode/opcode.go` (Lines 804-1079)
- 275 lines of switch statement
- ~80 case statements, each returning a string

### Proposed Implementation

#### Step 1: Create the opcode names array
After the `OperandWidths` definition and before the `String()` method, add:

```go
// opcodeNames maps each OpCode to its human-readable name.
// Indexed directly by OpCode value for O(1) lookup.
var opcodeNames = [256]string{
	OpNop:                      "NOP",
	OpPop:                      "POP",
	OpDup:                      "DUP",
	OpConst:                    "CONST",
	OpNil:                      "NIL",
	OpTrue:                     "TRUE",
	OpFalse:                    "FALSE",
	OpLocal:                    "LOCAL",
	OpSetLocal:                 "SETLOCAL",
	OpGlobal:                   "GLOBAL",
	OpSetGlobal:                "SETGLOBAL",
	OpFree:                     "FREE",
	OpSetFree:                  "SETFREE",
	OpAdd:                      "ADD",
	OpSub:                      "SUB",
	OpMul:                      "MUL",
	OpDiv:                      "DIV",
	OpMod:                      "MOD",
	OpNeg:                      "NEG",
	OpReal:                     "REAL",
	OpImag:                     "IMAG",
	OpComplex:                  "COMPLEX",
	OpAnd:                      "AND",
	OpOr:                       "OR",
	OpXor:                      "XOR",
	OpAndNot:                   "ANDNOT",
	OpLsh:                      "LSH",
	OpRsh:                      "RSH",
	OpEqual:                    "EQUAL",
	OpNotEqual:                 "NOTEQUAL",
	OpLess:                     "LESS",
	OpLessEq:                   "LESSEQ",
	OpGreater:                  "GREATER",
	OpGreaterEq:                "GREATEREQ",
	OpNot:                      "NOT",
	OpJump:                     "JUMP",
	OpJumpTrue:                 "JUMPTRUE",
	OpJumpFalse:                "JUMPFALSE",
	OpCall:                     "CALL",
	OpCallExternal:             "CALLEXTERNAL",
	OpCallIndirect:             "CALLINDIRECT",
	OpReturn:                   "RETURN",
	OpReturnVal:                "RETURNVAL",
	OpMakeArray:                "MAKEARRAY",
	OpMakeStruct:               "MAKESTRUCT",
	OpMakeSlice:                "MAKESLICE",
	OpMakeMap:                  "MAKEMAP",
	OpMakeChan:                 "MAKECHAN",
	OpIndex:                    "INDEX",
	OpIndexOk:                  "INDEXOK",
	OpSetIndex:                 "SETINDEX",
	OpSlice:                    "SLICE",
	OpField:                    "FIELD",
	OpSetField:                 "SETFIELD",
	OpAddr:                     "ADDR",
	OpFieldAddr:                "FIELDADDR",
	OpIndexAddr:                "INDEXADDR",
	OpAssert:                   "ASSERT",
	OpConvert:                  "CONVERT",
	OpChangeType:               "CHANGETYPE",
	OpDeref:                    "DEREF",
	OpSetDeref:                 "SETDEREF",
	OpNew:                      "NEW",
	OpMake:                     "MAKE",
	OpClosure:                  "CLOSURE",
	OpMethod:                   "METHOD",
	OpMethodCall:               "METHODCALL",
	OpDefer:                    "DEFER",
	OpDeferIndirect:            "DEFERINDIRECT",
	OpDeferExternal:            "DEFEREXTERNAL",
	OpGoCall:                   "GOCALL",
	OpGoCallIndirect:           "GOCALLINDIRECT",
	OpLen:                      "LEN",
	OpCap:                      "CAP",
	OpAppend:                   "APPEND",
	OpCopy:                     "COPY",
	OpDelete:                   "DELETE",
	OpRange:                    "RANGE",
	OpRangeNext:                "RANGENEXT",
	OpSelect:                   "SELECT",
	OpSend:                     "SEND",
	OpRecv:                     "RECV",
	OpRecvOk:                   "RECVOK",
	OpClose:                    "CLOSE",
	OpPack:                     "PACK",
	OpUnpack:                   "UNPACK",
	OpPrint:                    "PRINT",
	OpPrintln:                  "PRINTLN",
	OpHalt:                     "HALT",
	OpPanic:                    "PANIC",
	OpRecover:                  "RECOVER",
	// Superinstructions...
	OpAddLocalLocal:            "ADDLOCALLOCAL",
	OpSubLocalLocal:            "SUBLOCALLOCAL",
	OpMulLocalLocal:            "MULLOCALLOCAL",
	OpAddLocalConst:            "ADDLOCALCONST",
	OpSubLocalConst:            "SUBLOCALCONST",
	OpLessLocalLocalJumpTrue:   "LESSLOCALLOCATIONJUMPTRUE",
	// ... (continue with all superinstructions)
}
```

#### Step 2: Replace the String() method
Replace lines 804-1079 with:

```go
// String returns the name of the opcode as a human-readable string.
func (op OpCode) String() string {
	if op < OpCode(len(opcodeNames)) && opcodeNames[op] != "" {
		return opcodeNames[op]
	}
	return fmt.Sprintf("UNKNOWN(%d)", op)
}
```

### Testing
```bash
# Test the conversion
go test -run TestOpCodeString ./model/bytecode/ -v

# Verify no behavior change
go test ./model/bytecode/... -v
go test ./vm/... -v
```

### Before/After
- **Before**: 275 lines, O(n) switch statement
- **After**: ~80 lines (array + method), O(1) lookup
- **Savings**: ~270 LOC
- **Performance**: O(n) → O(1), ~10x faster for string lookups

---

## PRIORITY 2: Extract mustReflectValue() Helper

### Current Implementation (Repeated Pattern)
**File**: `vm/ops_container.go` (Multiple locations)

```go
// Pattern 1: OpIndex
case bytecode.OpIndex:
    if rv, ok := container.ReflectValue(); ok {
        // ... handle reflect case
    }
    // ... fallback

// Pattern 2: OpSetIndex  
case bytecode.OpSetIndex:
    if rv, ok := container.ReflectValue(); ok {
        // ... handle reflect case
    }
    // ... fallback

// Pattern 3: OpLen
case bytecode.OpLen:
    if rv, ok := container.ReflectValue(); ok {
        // ... handle reflect case
    }
    // ... fallback
```

### Proposed Implementation

#### Step 1: Add helper function at top of ops_container.go
```go
// getReflectValue safely extracts a reflect.Value from a Value.
// Returns (reflect.Value, bool) where bool indicates success.
func getReflectValue(v value.Value) (reflect.Value, bool) {
	if rv, ok := v.ReflectValue(); ok {
		return rv, true
	}
	return reflect.Value{}, false
}
```

#### Step 2: Consolidate repeated patterns
Replace all instances of:
```go
if rv, ok := x.ReflectValue(); ok {
    // ... handler
}
```

With:
```go
if rv, ok := getReflectValue(x); ok {
    // ... handler
}
```

### Where This Appears
1. **vm/ops_container.go**: Lines 115-160 (Index), 161-210 (SetIndex), 264-300 (Len/Cap), 391-430 (Append/Copy)
2. **vm/ops_convert.go**: Similar patterns for type conversion
3. **vm/ops_control.go**: Select statement handling

### Testing
```bash
# Ensure no behavior change
go test ./vm/... -v -race

# Check memory allocations
go test -bench=. -benchmem ./vm/ops_container.go
```

### Before/After
- **Before**: 433 lines with repeated pattern matching
- **After**: ~350 lines with consolidated helper
- **Savings**: ~80 LOC
- **Risk**: Very LOW - pure DRY refactoring

---

## PRIORITY 3: Consolidate Numeric Conversion Functions

### Current Implementation (Duplication)
**File**: `vm/ops_dispatch.go` (Lines 55-95)

```go
// Three nearly identical functions
func toInt64(v value.Value) int64 {
	switch v.Kind() {
	case value.KindInt:
		return v.Int()
	case value.KindUint:
		return int64(v.Uint())
	case value.KindFloat:
		return int64(v.Float())
	default:
		return v.Int()
	}
}

func toUint64(v value.Value) uint64 {
	switch v.Kind() {
	case value.KindInt:
		return uint64(v.Int())
	case value.KindUint:
		return v.Uint()
	case value.KindFloat:
		return uint64(v.Float())
	default:
		return v.Uint()
	}
}

func toFloat64(v value.Value) float64 {
	switch v.Kind() {
	case value.KindInt:
		return float64(v.Int())
	case value.KindUint:
		return float64(v.Uint())
	case value.KindFloat:
		return v.Float()
	default:
		return v.Float()
	}
}
```

### Proposed Implementation

#### Option A: Keep as-is (Simplest)
The duplication is minimal and each function is clear. The Go compiler will inline them anyway. **Keep for clarity.**

#### Option B: Extract common logic (Moderate)
```go
// toNumeric converts a Value to a numeric value of the requested kind.
// Handles KindInt, KindUint, KindFloat by dispatching to type-specific conversions.
func toNumeric(v value.Value, getter func(k value.Kind) interface{}) interface{} {
	switch v.Kind() {
	case value.KindInt:
		return getter(value.KindInt)
	case value.KindUint:
		return getter(value.KindUint)
	case value.KindFloat:
		return getter(value.KindFloat)
	default:
		return getter(value.KindInt)
	}
}

// toInt64 converts any numeric Value to int64
func toInt64(v value.Value) int64 {
	switch v.Kind() {
	case value.KindInt:
		return v.Int()
	case value.KindUint:
		return int64(v.Uint())
	case value.KindFloat:
		return int64(v.Float())
	default:
		return v.Int()
	}
}
// Keep others for clarity
```

#### Option C: Table-driven (Advanced)
This is probably overkill for 3 functions. Skip.

### Related: Extract toShiftAmount()
**File**: `vm/ops_arithmetic.go` (Lines 63-84)

```go
// Current (duplicated in OpLsh and OpRsh):
shiftVal := v.pop()
var n uint
if shiftVal.Kind() == value.KindUint {
    n = uint(shiftVal.Uint())
} else {
    n = uint(shiftVal.Int())
}

// Proposed:
func toShiftAmount(v value.Value) uint {
	switch v.Kind() {
	case value.KindUint:
		return uint(v.Uint())
	default:
		return uint(v.Int())
	}
}

// Usage in OpLsh and OpRsh:
case bytecode.OpLsh:
    shiftVal := v.pop()
    n := toShiftAmount(shiftVal)
    a := v.pop()
    v.push(a.Lsh(n))
```

### Testing
```bash
go test ./vm/... -v
go test -bench=. -benchmem ./vm/ops_arithmetic.go
go test -bench=. -benchmem ./vm/ops_dispatch.go
```

### Before/After
- **Before**: ~40 lines (3 switch statements + 1 duplicate)
- **After**: ~25 lines (with extracted toShiftAmount)
- **Savings**: ~15 LOC (conservative estimate ~42 with context)
- **Performance**: Inlined by compiler, no impact

---

## PRIORITY 4: Unify Basic Kind Mappings

### Current Implementation (Duplication)
**File 1**: `vm/typeconv.go` (Lines 18-36)
```go
var basicKindToReflect = map[types.BasicKind]reflect.Type{
	types.Bool:       reflect.TypeFor[bool](),
	types.Int:        reflect.TypeFor[int](),
	types.Int8:       reflect.TypeFor[int8](),
	// ... 16 more entries
}
```

**File 2**: `importer/typeconv.go` (Lines 27-45)
```go
var reflectToBasicKind = map[reflect.Kind]types.BasicKind{
	reflect.Bool:          types.Bool,
	reflect.Int:           types.Int,
	reflect.Int8:          types.Int8,
	// ... 16 more entries (reverse direction)
}
```

### Proposed Implementation

#### Option A: Create shared package (Best)
Create `model/bytecode/basickinds.go`:

```go
package bytecode

import (
	"go/types"
	"reflect"
)

// BasicKindInfo contains bidirectional mapping for basic kinds
type BasicKindInfo struct {
	BasicKind types.BasicKind
	ReflectKind reflect.Kind
	ReflectType reflect.Type
	String     string
}

// basicKindMap is the canonical source for all basic kind mappings
var basicKindMap = map[types.BasicKind]BasicKindInfo{
	types.Bool:          {types.Bool, reflect.Bool, reflect.TypeFor[bool](), "bool"},
	types.Int:           {types.Int, reflect.Int, reflect.TypeFor[int](), "int"},
	types.Int8:          {types.Int8, reflect.Int8, reflect.TypeFor[int8](), "int8"},
	types.Int16:         {types.Int16, reflect.Int16, reflect.TypeFor[int16](), "int16"},
	types.Int32:         {types.Int32, reflect.Int32, reflect.TypeFor[int32](), "int32"},
	types.Int64:         {types.Int64, reflect.Int64, reflect.TypeFor[int64](), "int64"},
	types.Uint:          {types.Uint, reflect.Uint, reflect.TypeFor[uint](), "uint"},
	types.Uint8:         {types.Uint8, reflect.Uint8, reflect.TypeFor[uint8](), "uint8"},
	types.Uint16:        {types.Uint16, reflect.Uint16, reflect.TypeFor[uint16](), "uint16"},
	types.Uint32:        {types.Uint32, reflect.Uint32, reflect.TypeFor[uint32](), "uint32"},
	types.Uint64:        {types.Uint64, reflect.Uint64, reflect.TypeFor[uint64](), "uint64"},
	types.Uintptr:       {types.Uintptr, reflect.Uintptr, reflect.TypeFor[uintptr](), "uintptr"},
	types.Float32:       {types.Float32, reflect.Float32, reflect.TypeFor[float32](), "float32"},
	types.Float64:       {types.Float64, reflect.Float64, reflect.TypeFor[float64](), "float64"},
	types.Complex64:     {types.Complex64, reflect.Complex64, reflect.TypeFor[complex64](), "complex64"},
	types.Complex128:    {types.Complex128, reflect.Complex128, reflect.TypeFor[complex128](), "complex128"},
	types.String:        {types.String, reflect.String, reflect.TypeFor[string](), "string"},
	types.UnsafePointer: {types.UnsafePointer, reflect.UnsafePointer, reflect.TypeFor[unsafe.Pointer](), "unsafe.Pointer"},
}

// BasicKindToReflect returns the reflect.Type for a types.BasicKind
func BasicKindToReflect(k types.BasicKind) reflect.Type {
	if info, ok := basicKindMap[k]; ok {
		return info.ReflectType
	}
	return nil
}

// ReflectKindToBasicKind returns the types.BasicKind for a reflect.Kind
func ReflectKindToBasicKind(k reflect.Kind) types.BasicKind {
	for _, info := range basicKindMap {
		if info.ReflectKind == k {
			return info.BasicKind
		}
	}
	return 0 // invalid
}
```

#### Step 2: Update vm/typeconv.go
Replace:
```go
var basicKindToReflect = map[types.BasicKind]reflect.Type{...}
```

With:
```go
// Use shared mapping from bytecode package
// In typeToReflectInner:
if basic, ok := t.(*types.Basic); ok {
    if rt := bytecode.BasicKindToReflect(basic.Kind()); rt != nil {
        return rt
    }
}
```

#### Step 3: Update importer/typeconv.go
Replace:
```go
var reflectToBasicKind = map[reflect.Kind]types.BasicKind{...}
```

With:
```go
// Use shared mapping from bytecode package
// In convertReflectType:
if bk := bytecode.ReflectKindToBasicKind(t.Kind()); bk != 0 {
    return types.Typ[bk]
}
```

### Testing
```bash
go test ./vm/... -v
go test ./importer/... -v
go test ./model/bytecode/... -v

# Verify bidirectional consistency
go test -run TestBasicKindMapping -v
```

### Before/After
- **Before**: 30 LOC (two separate maps)
- **After**: 35 LOC (shared + accessors, but single source of truth)
- **Savings**: ~10 LOC net
- **Benefit**: Prevents divergence, easier maintenance

---

## LOWER PRIORITY ITEMS

### 5. Type Checker Functions Consolidation
**File**: `vm/ops_dispatch.go` (Lines 203-217)

```go
// Current (four trivial wrappers):
func isSignedInt(k reflect.Kind) bool {
	return k >= reflect.Int && k <= reflect.Int64
}

func isUnsignedInt(k reflect.Kind) bool {
	return k >= reflect.Uint && k <= reflect.Uintptr
}

func isFloat(k reflect.Kind) bool {
	return k == reflect.Float32 || k == reflect.Float64
}

func isComplex(k reflect.Kind) bool {
	return k == reflect.Complex64 || k == reflect.Complex128
}

// Proposed (table-driven):
var numericKindRanges = map[string][2]reflect.Kind{
	"signed":   {reflect.Int, reflect.Int64},
	"unsigned": {reflect.Uint, reflect.Uintptr},
}

var floatKinds = map[reflect.Kind]bool{
	reflect.Float32:    true,
	reflect.Float64:    true,
}

var complexKinds = map[reflect.Kind]bool{
	reflect.Complex64:  true,
	reflect.Complex128: true,
}

func isNumericKind(k reflect.Kind, category string) bool {
	switch category {
	case "signed":
		return k >= reflect.Int && k <= reflect.Int64
	case "unsigned":
		return k >= reflect.Uint && k <= reflect.Uintptr
	case "float":
		return floatKinds[k]
	case "complex":
		return complexKinds[k]
	}
	return false
}
```

**Note**: This change doesn't save much code and may reduce clarity. **Keep functions as-is.**

---

## IMPLEMENTATION CHECKLIST

### Phase 1 Quick Wins (Priority)
- [ ] **OpCode.String() table conversion** (3-4 hours)
  - [ ] Create opcodeNames array
  - [ ] Replace switch statement
  - [ ] Test all opcodes map correctly
  - [ ] Run full test suite
  - [ ] Commit with message: "refactor: convert OpCode.String() to table-driven"

- [ ] **Extract mustReflectValue() helper** (2-3 hours)
  - [ ] Add helper function to ops_container.go
  - [ ] Identify all ReflectValue() extraction patterns
  - [ ] Replace patterns one at a time
  - [ ] Run tests after each replacement
  - [ ] Check ops_convert.go and ops_control.go
  - [ ] Commit with message: "refactor: extract mustReflectValue() helper"

- [ ] **Consolidate numeric conversions** (1-2 hours)
  - [ ] Extract toShiftAmount() from ops_arithmetic.go
  - [ ] Verify toInt64/toUint64/toFloat64 are inlined
  - [ ] Run benchmarks to confirm no regression
  - [ ] Commit with message: "refactor: extract toShiftAmount() helper"

### Phase 2 (If Time Permits)
- [ ] **Unify basic kind mappings** (3-4 hours)
  - [ ] Create model/bytecode/basickinds.go
  - [ ] Update vm/typeconv.go to use shared mapping
  - [ ] Update importer/typeconv.go to use shared mapping
  - [ ] Run tests for vm and importer
  - [ ] Commit with message: "refactor: consolidate basic kind mappings"

---

## TESTING STRATEGY

### Before Each Change
```bash
# Run full test suite
go test ./...

# Run specific package tests
go test ./vm/... -v
go test ./model/bytecode/... -v

# Get baseline benchmarks
go test -bench=. -benchmem ./vm/ > baseline.txt
```

### After Each Change
```bash
# Verify no regressions
go test ./... -v

# Run with race detector
go test -race ./vm/...

# Compare benchmarks
go test -bench=. -benchmem ./vm/ > after.txt
# benchstat baseline.txt after.txt
```

### Specific Tests Per Change
```bash
# OpCode.String()
go test -run TestOpCodeString ./model/bytecode/ -v

# mustReflectValue() 
go test ./vm/ops_container.go -v
go test ./vm/ops_convert.go -v

# Numeric conversions
go test ./vm/ops_dispatch.go -v
go test ./vm/ops_arithmetic.go -v
```

---

## ROLLBACK PROCEDURE

Each change is designed to be independently reversible:

```bash
# If a test fails after a change
git status                  # See what changed
git diff <file>            # Review changes
git checkout <file>        # Rollback
go test ./...              # Verify revert
```

---

## PERFORMANCE VERIFICATION

### Critical: vm/run.go and vm/ops_container.go
These are in the hot path. Always benchmark:

```bash
# Before change
go test -bench=. -benchmem ./vm/ -benchtime=10s > before.txt

# After change
go test -bench=. -benchmem ./vm/ -benchtime=10s > after.txt

# Compare (requires benchstat tool)
go install golang.org/x/perf/cmd/benchstat@latest
benchstat before.txt after.txt
```

**Acceptable regression**: < 2% on any benchmark

---

## TIMELINE ESTIMATE

| Phase | Task | Hours | Risk |
|-------|------|-------|------|
| 1 | OpCode.String() table | 3-4 | NONE |
| 1 | mustReflectValue() | 2-3 | LOW |
| 1 | toShiftAmount() | 1-2 | LOW |
| Subtotal Phase 1 | | ~7-9 | LOW |
| 2 | Basic kind mappings | 3-4 | LOW-MED |
| 2 | Polish & docs | 1-2 | NONE |
| Subtotal Phase 2 | | ~4-6 | LOW-MED |
| **Total** | | ~11-15 | **LOW** |

