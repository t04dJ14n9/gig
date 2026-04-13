# Gig Interpreter - Code Simplification Quick Reference

## Quick Stats
- **Total Opportunities Found**: 47
- **Potential LOC Savings**: ~520-600 LOC (~6-7% of core packages)
- **Analysis Coverage**: 8,500 LOC across 6 core packages
- **Detailed Report**: See `CODE_SIMPLIFICATION_ANALYSIS.md`

## The Big Three (High Impact)

### 1. OpCode.String() - 275 lines → ~50 lines
**File**: `model/bytecode/opcode.go` (lines 804-1079)
- **Current**: Giant 275-line switch statement  
- **Solution**: Table-driven with `var opcodeNames [256]string`
- **Benefit**: ~270 LOC saved, O(1) lookup, better maintainability
- **Risk**: NONE - Zero runtime impact, pure code quality

### 2. vm/ops_container.go - Extract mustReflectValue()
**File**: `vm/ops_container.go` (lines 31-464)
- **Current**: 433-line handler with repeated `if rv, ok := container.ReflectValue()`
- **Solution**: Extract helper function `mustReflectValue(v Value) reflect.Value`
- **Benefit**: ~80 LOC saved through DRY principle
- **Risk**: LOW - Simple helper extraction

### 3. vm/typeconv.go + importer/typeconv.go - Unify mappings
**Files**: `vm/typeconv.go` (18-36) + `importer/typeconv.go` (27-45)
- **Current**: Parallel mappings (basicKindToReflect vs reflectToBasicKind)
- **Solution**: Single source of truth in shared package
- **Benefit**: ~30 LOC saved, prevents divergence
- **Risk**: LOW-MEDIUM - Requires careful refactoring

---

## Medium Impact Wins

### 4. Shift Logic Consolidation
**File**: `vm/ops_arithmetic.go` (lines 63-84)
- **Savings**: ~8 LOC
- **Pattern**: Extract `toShiftAmount()` helper for OpLsh/OpRsh

### 5. Numeric Conversion Functions
**File**: `vm/ops_dispatch.go` (lines 55-95)
- **Savings**: ~30 LOC
- **Pattern**: Three similar switches (toInt64, toUint64, toFloat64)

### 6. Type Checker Functions
**File**: `vm/ops_dispatch.go` (lines 203-217)
- **Savings**: ~12 LOC
- **Pattern**: Four trivial range-check wrappers

### 7. Hot Path Patterns
**File**: `vm/run.go` (lines 47-1207)
- **Savings**: ~50 LOC
- **Pattern**: Extract `pop2()` and `push2()` helpers
- **Risk**: MEDIUM - Must benchmark for no regression

---

## Implementation Priority

### Phase 1: Do First ✅
```
1. OpCode.String() table-driven          [~270 LOC, 0% risk]
2. mustReflectValue() helper             [~80 LOC, 5% risk]
3. Numeric conversions consolidation     [~42 LOC, 5% risk]
```
**Total Effort**: ~3 hours | **LOC Saved**: ~392 | **Risk**: LOW

### Phase 2: If Time Permits
```
4. Basic kind mappings unification       [~30 LOC, 10% risk]
5. Shift logic extraction                [~8 LOC, 5% risk]
6. vm/run.go pattern extraction          [~50 LOC, 20% risk]
```
**Total Effort**: ~4 hours | **LOC Saved**: ~88 | **Risk**: MEDIUM

### Phase 3: Nice to Have
```
7. Compiler dispatch table-driven        [~50-100 LOC, 15% risk]
8. ReflectValue() caching                [Micro-opt, 5% risk]
9. Trivial wrapper removal               [~10 LOC, 5% risk]
```

---

## Detailed Checklist

### [ ] Convert OpCode.String() (PRIORITY #1)
- [ ] Create `var opcodeNames [256]string` with all names
- [ ] Replace switch with table lookup
- [ ] Handle unknown opcodes safely
- [ ] Run tests: `go test ./model/bytecode/...`
- [ ] Verify: `grep -c "OpCode.*String" *_test.go` (should still pass)

### [ ] Extract mustReflectValue() (PRIORITY #2)
- [ ] Add function in `vm/ops_container.go`
- [ ] Replace all `if rv, ok := x.ReflectValue()` patterns
- [ ] Check `ops_convert.go` for same pattern
- [ ] Run tests: `go test ./vm/...`

### [ ] Consolidate Numeric Conversions (PRIORITY #3)
- [ ] Review `ops_dispatch.go` toInt64/toUint64/toFloat64
- [ ] Extract common logic (if possible)
- [ ] Extract `toShiftAmount()` from ops_arithmetic.go
- [ ] Run tests: `go test ./vm/...`

### [ ] Unify Basic Kind Mappings (PRIORITY #4)
- [ ] Create single mapping source (maybe new file)
- [ ] Update vm/typeconv.go to use
- [ ] Update importer/typeconv.go to use
- [ ] Run tests: `go test ./vm/... ./importer/...`

---

## Risk Mitigation

**Before Each Change**:
1. Run full test suite: `go test ./...`
2. Benchmark VM performance: `go test -bench=. ./vm/...`
3. Check coverage: `go test -cover ./...`

**After Each Change**:
1. Verify all tests still pass
2. Run linter: `golangci-lint run`
3. Check for performance regressions

**Rollback Strategy**:
- Each change is isolated in a separate commit
- Easy to revert if regression detected
- Profile hot paths (run.go, ops_container.go)

---

## Testing Strategy

### For OpCode.String() - Conversion
```bash
go test ./model/bytecode/ -v -run String
# Should verify all opcodes map correctly
```

### For mustReflectValue() - Extraction
```bash
go test ./vm/ -v -race
# Run with race detector to catch any concurrency issues
```

### For Consolidations - General
```bash
go test ./... -count=100
# Run multiple times to catch flaky tests
go test -bench=. -benchmem ./vm/
# Check memory allocations don't increase
```

---

## FAQ

**Q: Will these changes affect performance?**
A: OpCode.String() will improve (O(n) → O(1)). Others have negligible impact or must be benchmarked.

**Q: Should we do all changes at once?**
A: No - implement Phase 1 first, measure, then Phase 2.

**Q: What if tests fail?**
A: Roll back the change and debug. Each commit is independent.

**Q: Are there more opportunities?**
A: Yes - see the full analysis. These are the highest-impact ones.

**Q: Can we automate any of this?**
A: Yes - OpCode.String() could be code-generated from opcode constants.

