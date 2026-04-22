# gofun 正确性测试结果

使用 Gig 的 305 个 correctness 测试用例（无参导出函数）跑新 gofun (onefun/gofun)。

- **通过**: 295 (96.7%)
- **Build 失败**: 0
- **Run 失败**: 10
- **总失败**: 10 (3.3%)

## 按包统计

| 包 | 总数 | 通过 | 失败 | 通过率 |
|-----|------|------|------|--------|
| advanced | 19 | 19 | 0 | 100% ✅ |
| algorithms | 10 | 10 | 0 | 100% ✅ |
| arithmetic | 10 | 10 | 0 | 100% ✅ |
| bitwise | 8 | 8 | 0 | 100% ✅ |
| closures | 7 | 7 | 0 | 100% ✅ |
| closures_advanced | 7 | 7 | 0 | 100% ✅ |
| controlflow | 12 | 12 | 0 | 100% ✅ |
| cornercases | 110 | 106 | 4 | 96% ❌ |
| edgecases | 12 | 12 | 0 | 100% ✅ |
| functions | 25 | 25 | 0 | 100% ✅ |
| leetcode_hard | 10 | 6 | 4 | 60% ❌ |
| mapadvanced | 6 | 6 | 0 | 100% ✅ |
| maps | 7 | 7 | 0 | 100% ✅ |
| multiassign | 6 | 6 | 0 | 100% ✅ |
| namedreturn | 3 | 3 | 0 | 100% ✅ |
| recursion | 6 | 6 | 0 | 100% ✅ |
| scope | 6 | 6 | 0 | 100% ✅ |
| slices | 8 | 8 | 0 | 100% ✅ |
| slicing | 7 | 7 | 0 | 100% ✅ |
| strings_pkg | 7 | 6 | 1 | 86% ❌ |
| switch | 8 | 8 | 0 | 100% ✅ |
| typeconv | 5 | 4 | 1 | 80% ❌ |
| variables | 6 | 6 | 0 | 100% ✅ |

## 失败列表


### cornercases

- `Range_EmptyString` [run]: panic: main.go:666:2: reflect: call of reflect.Value.MapKeys on string Value
- `String_LastByte` [run]: panic: main.go:258:10: unexpected instruction: *ssa.Index
- `String_SingleByteIndex` [run]: panic: main.go:253:10: unexpected instruction: *ssa.Index
- `Struct_PointerReceiver` [run]: panic: main.go:558:15: reflect.StructOf: field "value" is unexported but missing PkgPath

### leetcode_hard

- `EditDistance` [run]: panic: main.go:385:12: unexpected instruction: *ssa.Index
- `MinimumWindowSubstring` [run]: panic: main.go:418:9: unexpected instruction: *ssa.Index
- `RegularExpressionMatching` [run]: panic: main.go:132:7: unexpected instruction: *ssa.Index
- `WordLadder` [run]: panic: main.go:310:37: panic: main.go:291:8: unexpected instruction: *ssa.Index

### strings_pkg

- `Index` [run]: panic: main.go:26:17: unexpected instruction: *ssa.Index

### typeconv

- `StringToByteConversion` [run]: panic: main.go:22:15: unexpected instruction: *ssa.Index

## 错误类型统计

- **8** × `unexpected instruction: *ssa.Index`
- **1** × `reflect: call of reflect.Value.MapKeys on string Value`
- **1** × `reflect.StructOf: unexported field missing PkgPath`
