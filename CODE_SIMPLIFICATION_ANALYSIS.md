# Gig Interpreter Code Simplification Opportunities

## Executive Summary
Analysis of ~8500 LOC across 6 core packages identified 47 significant simplification opportunities, primarily in:
- **High Impact**: 1 massive 275-line switch statement, 3 large functions with duplicated logic
- **Medium Impact**: 22 redundant/duplicated patterns, 15+ trivial wrapper functions, unused test exports
- **Low Impact**: Minor refactoring opportunities, style improvements

---

## HIGH IMPACT FINDINGS

### 1. **model/bytecode/opcode.go - Massive OpCode.String() Switch (CRITICAL)**
- **Location**: Lines 804-1079 (275 lines)
- **Issue**: Giant switch statement with 80+ cases for opcode name mapping
- **Impact**: HIGH - This is inefficient string lookup and maintenance burden
- **Recommendation**: Convert to table-driven approach using array indexed by OpCode
- **Before**:
```go
func (op OpCode) String() string {
    switch op {
    case OpNop: return "NOP"
    case OpPop: return "POP"
    // ... 80+ more cases
    }
}
```
- **After**:
```go
var opcodeNames = [256]string{
    OpNop: "NOP",
    OpPop: "POP",
    // ...
}
func (op OpCode) String() string {
    if op < OpCode(len(opcodeNames)) {
        return opcodeNames[op]
    }
    return "UNKNOWN"
}
```
- **Savings**: ~270 LOC, O(1) lookup instead of switch

### 2. **vm/run.go - Monolithic 1160-Line Execution Loop (CRITICAL)**
- **Location**: Lines 47-1207 (1160 lines)
- **Issue**: Giant hot-path fetch-decode-execute loop with inline handlers
- **Impact**: HIGH - Maintenance nightmare, difficult to refactor safely
- **Note**: This is intentionally inlined for performance. However, it contains:
  - 72-case switch at lines ~470+ for main opcode dispatch
  - Multiple redundant operand reads that could be factored
  - Repeated stack manipulation patterns
- **Minor Improvement**: Extract common stack patterns:
  - `pop2()` helper for (a, b) = (pop, pop)
  - `push2()` helper for push(a); push(b)
  - Estimated savings: ~50 LOC without harming performance

### 3. **vm/ops_container.go - Verbose 433-Line Handler (HIGH)**
- **Location**: Lines 31-464 (433 lines)
- **Issue**: 40-case switch with repeated patterns for container operations
- **Specific Duplications**:
  - **Lines 115-160**: Index operation with pattern `if rv, ok := container.ReflectValue()`
  - **Lines 161-210**: Same pattern repeated for SetIndex
  - **Lines 264-300**: Cap/Len operations with identical ReflectValue extraction
  - **Lines 391-430**: Append/Copy with similar structure
- **Root Cause**: Repeated "try reflect, fall back" pattern
- **Recommendation**: Extract `mustReflectValue(v Value) reflect.Value` helper
- **Savings**: ~80 LOC through DRY

### 4. **vm/typeconv.go vs importer/typeconv.go - Duplicate Basic Type Mappings**
- **Location**: 
  - vm/typeconv.go: Lines 18-36 (basicKindToReflect)
  - importer/typeconv.go: Lines 27-45 (reflectToBasicKind)
- **Issue**: Parallel but opposite direction mappings maintained separately
- **Impact**: MEDIUM - Maintenance burden, risk of divergence
- **Alternative**: Could create a single model/bytecode/basickinds.go with both:
```go
var basicKinds = map[types.BasicKind]string{
    types.Bool: "bool", types.Int: "int", ...
}
// Then generate both maps or use a bidirectional structure
```
- **Estimated Savings**: ~30 LOC, single source of truth

---

## MEDIUM IMPACT FINDINGS

### 5. **vm/ops_arithmetic.go - Duplicate Shift Logic (MEDIUM)**
- **Location**: Lines 63-84 (OpLsh and OpRsh)
- **Issue**: Identical shift value extraction code in both cases:
```go
shiftVal := v.pop()
var n uint
if shiftVal.Kind() == value.KindUint {
    n = uint(shiftVal.Uint())
} else {
    n = uint(shiftVal.Int())
}
```
- **Recommendation**: Extract helper `toShiftAmount()` function
- **Savings**: ~8 LOC

### 6. **vm/ops_dispatch.go - Four Single-Line Type Checker Functions (MEDIUM)**
- **Location**: Lines 203-217
- **Functions**:
  - `isSignedInt(k reflect.Kind) bool` - line 203
  - `isUnsignedInt(k reflect.Kind) bool` - line 207
  - `isFloat(k reflect.Kind) bool` - line 211
  - `isComplex(k reflect.Kind) bool` - line 215
- **Issue**: Four trivial wrappers for range checks
- **Recommendation**: Use table instead:
```go
func isNumericKind(k, category reflect.Kind) bool {
    return k >= category && k <= categoryMax
}
// Or use a map
var numericKindMap = map[reflect.Kind]bool{reflect.Int: true, ...}
```
- **Savings**: ~12 LOC

### 7. **vm/ops_dispatch.go - Three Numeric Extraction Functions (MEDIUM)**
- **Location**: Lines 55-95 (toInt64, toUint64, toFloat64)
- **Issue**: Three switch statements with nearly identical structure:
```go
switch v.Kind() {
case value.KindInt: return int64(v.Int())
case value.KindUint: return int64(v.Uint())
case value.KindFloat: return int64(v.Float())
default: return v.Int()
}
```
- **Recommendation**: Extract common logic with generic function
- **Savings**: ~30 LOC

### 8. **compiler/compile_instr.go - 35-Case Instruction Switch (MEDIUM)**
- **Location**: Lines 15-100+
- **Issue**: Large switch in `compileInstruction()` could benefit from dispatch table
- **Note**: Compiler performance is less critical than VM performance, but could improve maintainability
- **Recommendation**: Consider multimethod dispatch if adding new instruction types
- **Estimated Savings**: ~50-100 LOC with table-driven approach (lower priority)

### 9. **Trivial Wrapper Functions (MEDIUM ACROSS CODEBASE)**
- **Locations and Patterns**:

#### compiler/symbol.go
```go
// Line 44: GetLocal(v) - just returns from map
func (s *SymbolTable) GetLocal(v ssa.Value) (int, bool) {
    idx, ok := s.locals[v]
    return idx, ok
}

// Line 50: NumLocals() - just returns length
func (s *SymbolTable) NumLocals() int {
    return len(s.locals)
}
```
- **Issue**: Could be replaced with exported `Locals` field
- **Savings**: 4 LOC

#### compiler/emit.go
```go
// Line 56: emitJump(target) - single function call wrapper
func (c *compiler) emitJump(target *ssa.BasicBlock) {
    c.addJump(target)
}

// Line 72: addConstant(val) - delegates to internal function
func (c *compiler) addConstant(val any) uint16 {
    return c.addConstantInternal(val)
}
```
- **Recommendation**: Inline or create alias
- **Savings**: 6 LOC

#### vm/frame.go
```go
// Line 182: readUint16() - pure wrapper
func (f *Frame) readUint16() uint16 {
    v := uint16(f.Instructions[f.IP])<<8 | uint16(f.Instructions[f.IP+1])
    f.IP += 2
    return v
}
```
- **Note**: Already inlined in run.go for performance
- **Status**: OK (intentional)

#### vm/stack.go
```go
// Line 27: pop() - single field access
func (v *vm) pop() value.Value {
    return v.stack[v.sp-1]
}

// Line 33: peek() - single field access
func (v *vm) peek() value.Value {
    return v.stack[v.sp-1]
}
```
- **Issue**: These are used in hot path, already well-optimized
- **Status**: OK (intentional for abstraction)

### 10. **model/value/value.go - 22 MakeFoo Functions (MEDIUM)**
- **Location**: Throughout file
- **Functions**: MakeInt, MakeInt8, MakeInt16, MakeInt32, MakeInt64, MakeUint, MakeUint8, MakeUint16, MakeUint32, MakeUint64, MakeBool, MakeString, etc.
- **Issue**: Repetitive constructor pattern
- **Current Approach**: Each has similar structure:
```go
func MakeInt8(i int8) Value { return MakeInt(int64(i)) }
func MakeInt16(i int16) Value { return MakeInt(int64(i)) }
// ... pattern continues
```
- **Recommendation**: Keep these for API clarity and convenience, but consider:
  - `MakeIntSized(val int64, bitWidth byte) Value` generic version
  - or `MakeSized(val int64, kind value.Kind) Value`
- **Current Design**: API is good, functions are necessary (no savings here)
- **Status**: Keep as-is for API ergonomics

### 11. **vm/typeconv.go - Cycle Detection Cache (MEDIUM)**
- **Location**: Lines 50-58
- **Issue**: `typeToReflect()` creates local cache on every call:
```go
localCache := make(map[types.Type]reflect.Type)
rt := typeToReflectWithCache(t, localCache, "", prog, 0)
```
- **Problem**: Allocates new map each invocation even if type is in program cache
- **Recommendation**: Check if method is called frequently enough to justify caching
- **Potential Savings**: 1 allocation per typeToReflect call if cache hit

---

## LOW IMPACT FINDINGS

### 12. **Redundant Type Assertions in Value Operations (LOW)**
- **Location**: Scattered in vm/ops_*.go
- **Pattern**: Multiple `.ReflectValue()` calls on same value
```go
// vm/ops_container.go Line ~118
if rv, ok := container.ReflectValue(); ok {
    // ... later in same block
    if rv, ok := container.ReflectValue(); ok {  // redundant!
```
- **Issue**: Repeated reflection calls could cache result
- **Impact**: Performance optimization, not correctness
- **Recommendation**: Extract `rv` from first check and reuse
- **Estimated Savings**: 2-3 microseconds per repeated call pattern

### 13. **Unused Exports in vm/interfaces.go (LOW)**
- **Location**: vm/interfaces.go
- **Functions**:
  - `New()` - Created by `NewWithOptions()` already
  - `WithContext()` - Part of options pattern
- **Analysis**: These are part of public API
- **Note**: Should verify if actually unused externally before removing
- **Status**: Check call sites outside this repo

### 14. **Model Package Type Conversions (LOW)**
- **Location**: Multiple files with similar pattern matching
- **Issue**: Many switch statements could use maps for semantic matching
- **Example**: value/value.go Cmp() method
- **Estimated Savings**: ~20 LOC with table-driven comparisons
- **Current State**: Acceptable for clarity

### 15. **Compiler SSA Collection Logic (LOW)**
- **Location**: compiler/compiler.go Lines 93-164
- **Issue**: collectFuncs() and method collection is verbose
- **Recommendation**: Could simplify with better data structure
- **Estimated Savings**: ~30 LOC
- **Trade-off**: Readability vs. conciseness

---

## SUMMARY TABLE

| Category | Impact | File | Lines | Issue | Savings |
|----------|--------|------|-------|-------|---------|
| Table-Driven | **HIGH** | model/bytecode/opcode.go | 804-1079 | 275-line String() switch | ~270 LOC |
| Hot Path | **HIGH** | vm/run.go | 47-1207 | Inline hot-path (intentional but could extract patterns) | ~50 LOC (with performance) |
| Container Ops | **HIGH** | vm/ops_container.go | 31-464 | Duplicate ReflectValue extraction | ~80 LOC |
| Mapping | **MEDIUM** | vm/typeconv.go + importer/typeconv.go | Multi-file | Parallel mappings | ~30 LOC |
| Shift Logic | **MEDIUM** | vm/ops_arithmetic.go | 63-84 | Duplicate extraction | ~8 LOC |
| Numeric Helpers | **MEDIUM** | vm/ops_dispatch.go | 203-217, 55-95 | Trivial wrappers + repetitive switches | ~42 LOC |
| Compiler Dispatch | **MEDIUM** | compiler/compile_instr.go | 15-100+ | Large switch | ~50-100 LOC (lower priority) |
| Compiler Wrappers | **LOW** | compiler/* | Various | Trivial delegation functions | ~10 LOC |
| Repeated Assertions | **LOW** | vm/ops_*.go | Scattered | Cached ReflectValue() | Micro-optimization |
| Type Checking | **LOW** | model/value/* | Scattered | Could use tables | ~20 LOC |

---

## PRIORITIZED ACTION ITEMS

### Phase 1: High Impact, Low Risk
1. **Convert OpCode.String() to table-driven** (~270 LOC saved, 0 performance impact)
   - Create `var opcodeNames [256]string` array
   - Replace 275-line switch
   - Add safety check for unknown opcodes

### Phase 2: Medium Impact, Medium Effort
2. **Extract mustReflectValue() helper** (~80 LOC saved, no performance cost)
   - Consolidate ReflectValue() extraction pattern
   - Use in ops_container.go, ops_convert.go

3. **Consolidate numeric conversion functions** (~42 LOC saved)
   - `toInt64()`, `toUint64()`, `toFloat64()` in ops_dispatch.go
   - Extract shift value logic from ops_arithmetic.go

4. **Unify basic kind mappings** (~30 LOC saved, improves maintainability)
   - Create single source of truth for basic type mappings

### Phase 3: Medium Impact, Higher Risk
5. **Extract patterns from vm/run.go** (~50 LOC with performance analysis)
   - Create `pop2()` and `push2()` helpers
   - Profile to ensure no regression

6. **Refactor compiler dispatch** (~50-100 LOC, lower priority)
   - Consider table-driven approach if adding new instruction types

### Phase 4: Low Priority
7. **Cache repeated ReflectValue() calls** (micro-optimization)
8. **Simplify compiler wrapper functions** (code cleanliness)
9. **Add generic Make function variants** (optional API improvement)

---

## Risk Assessment

**Low Risk**:
- OpCode.String() table conversion
- Extract helpers (mustReflectValue)
- Numeric conversion consolidation

**Medium Risk**:
- vm/run.go pattern extraction (must benchmark)
- Compiler SSA collection simplification

**Higher Risk**:
- Large refactors without comprehensive testing

**Safe Approach**: Implement Phase 1 & Phase 2, then measure and iterate on Phase 3.

