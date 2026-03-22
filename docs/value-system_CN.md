# 值系统（Value System）

> 深入解析 Gig 的标签联合体（tagged-union）`Value` 类型——解释器中每一个栈槽、局部变量、常量和函数参数的底层数据表示。

## 目录

1. [概述](#概述)
2. [内存布局](#内存布局)
3. [Kind 类型系统](#kind-类型系统)
4. [构造器](#构造器)
5. [访问器](#访问器)
6. [算术运算](#算术运算)
7. [比较与相等](#比较与相等)
8. [位运算](#位运算)
9. [类型转换](#类型转换)
10. [容器操作](#容器操作)
11. [指针操作](#指针操作)
12. [通道操作](#通道操作)
13. [反射互操作](#反射互操作)
14. [闭包支持](#闭包支持)
15. [原生快速路径](#原生快速路径)
16. [设计思考](#设计思考)
17. [性能特性](#性能特性)
18. [文件组织](#文件组织)

---

## 概述

`Value` 类型是 Gig 解释器中的通用数据单元。解释器中的每一个局部变量、栈槽、常量、函数参数和返回值都是一个 `Value`。其设计遵循 **标签联合体（tagged-union）** 模式——用一个结构体存储任意 Go 类型，并通过类型标签（`kind`）在运行时区分具体变体。

核心设计理念是 **双层分离**：

- **原始类型**（`bool`、`int`、`uint`、`float`、`nil`）直接存储在结构体的数值字段中。无堆分配、无反射、无 GC 压力。
- **复合类型**（`string`、`complex`、`[]byte`、`[]int64`、闭包）存储在 `obj` 字段中，作为原生 Go 对象——尽可能避免反射。
- **复杂类型**（`struct`、`map`、`chan`、通用切片、接口）存储在 `obj` 中，作为 `reflect.Value`——需要动态派发的类型的兜底方案。

这种设计灵感来自 Lua、Python 等动态语言运行时的值表示方式，但针对 Go 的静态类型系统做了适配，实现了与 `reflect` 的无缝互操作。

---

## 内存布局

```
┌──────────────────────────────────────────────────────────┐
│                    Value (32 字节)                         │
├────────┬─────────┬───────────────────────────────────────┤
│  kind  │  填充    │                 num                    │
│ 1 字节  │ 7 字节   │               8 字节                    │
├────────┴─────────┴───────────────────────────────────────┤
│                        obj                                │
│                      16 字节                               │
│               (interface{} / any)                         │
└──────────────────────────────────────────────────────────┘
```

```go
type Value struct {
    kind Kind    // 1 字节：类型标签
    num  int64   // 8 字节：存储 bool、int、uint 位模式、float64 位模式
    obj  any     // 16 字节：存储 string、complex128、reflect.Value、原生复合类型，或 nil
}
```

总大小在 64 位系统上为 **32 字节**。这是刻意为之：

- **32 字节 = 4 个机器字** —— 能整齐地放入缓存行和 CPU 寄存器。
- **原始类型仅使用 `kind` + `num`**，`obj` 为 nil，意味着永远不会触及堆内存，永远不会产生 GC 压力。
- **`obj` 字段** 是 Go 的 `interface{}`（16 字节：类型指针 + 数据指针），可以持有任何 Go 值，对于小型类型无需额外分配。

---

## Kind 类型系统

`Kind` 类型是一个 `uint8` 枚举，标识 `Value` 中存储的具体变体：

```go
type Kind uint8

const (
    KindInvalid   Kind = iota  // 零值，未初始化
    KindNil                     // nil
    KindBool                    // bool（存储在 num 中：0 或 1）
    KindInt                     // int, int8, int16, int32, int64（存储在 num 中）
    KindUint                    // uint, uint8, ..., uint64, uintptr（存储在 num 中的位模式）
    KindFloat                   // float32, float64（存储在 num 中的 float64 位模式）
    KindString                  // string（存储在 obj 中）
    KindComplex                 // complex64, complex128（存储在 obj 中）
    KindPointer                 // *T（obj 中为 reflect.Value、*int64 或 *Value）
    KindSlice                   // []T（obj 中为 reflect.Value、[]int64 或 []Value）
    KindArray                   // [N]T（obj 中为 reflect.Value）
    KindMap                     // map[K]V（obj 中为 reflect.Value）
    KindChan                    // chan T（obj 中为 reflect.Value）
    KindFunc                    // func（obj 中为 *Closure 或原生函数）
    KindStruct                  // struct（obj 中为 reflect.Value）
    KindInterface               // interface（obj 中为 reflect.Value）
    KindReflect                 // 兜底：任何类型存储为 reflect.Value
    KindBytes                   // []byte（obj 中为原生 []byte，零反射）
)
```

### 各 Kind 的存储策略

| Kind | `num` 字段 | `obj` 字段 | 堆分配？ | 需要反射？ |
|------|-----------|-----------|---------|----------|
| `KindNil` | 0 | nil | 否 | 否 |
| `KindBool` | 0 或 1 | nil | 否 | 否 |
| `KindInt` | int64 值 | nil | 否 | 否 |
| `KindUint` | uint64 位模式 | nil | 否 | 否 |
| `KindFloat` | float64 位模式 | nil | 否 | 否 |
| `KindString` | — | `string` | 否* | 否 |
| `KindComplex` | — | `complex128` | 是** | 否 |
| `KindBytes` | — | `[]byte` | 否* | 否 |
| `KindSlice` | — | `[]int64` / `[]Value` / `reflect.Value` | 视情况 | 有时 |
| `KindPointer` | — | `*int64` / `*Value` / `reflect.Value` | 视情况 | 有时 |
| `KindFunc` | — | `*Closure` / 原生函数 | 否* | 否 |
| `KindMap` | — | `reflect.Value` | 是 | 是 |
| `KindStruct` | — | `reflect.Value` | 是 | 是 |
| `KindChan` | — | `reflect.Value` | 是 | 是 |
| `KindArray` | — | `reflect.Value` | 是 | 是 |
| `KindInterface` | — | `reflect.Value` | 是 | 是 |
| `KindReflect` | — | `reflect.Value` | 是 | 是 |

\* 字符串和切片在 Go 中是引用类型——`obj` 字段持有的是头部信息，而非数据拷贝。
\** `complex128` 存储在 `interface{}` 中时会被装箱。

---

## 构造器

### 原始类型构造器（零分配）

```go
MakeNil() Value                           // → KindNil
MakeBool(b bool) Value                    // → KindBool, num = 0/1
MakeInt(i int64) Value                    // → KindInt, num = i
MakeUint(u uint64) Value                  // → KindUint, num = int64(u)
MakeFloat(f float64) Value                // → KindFloat, num = float64 位模式
MakeString(s string) Value                // → KindString, obj = s
MakeComplex(real, imag float64) Value     // → KindComplex, obj = complex128
```

每个构造器都是纯粹的结构体字面量——没有函数调用、没有内存分配。Go 编译器可以完全内联这些操作。

### 特化构造器（避免反射）

```go
MakeBytes(b []byte) Value                 // → KindBytes, obj = b
MakeIntSlice(s []int64) Value             // → KindSlice, obj = s
MakeIntPtr(p *int64) Value                // → KindPointer, obj = p
MakeFunc(fn any) Value                    // → KindFunc, obj = fn
MakeValueSlice(vals []Value) Value        // → KindSlice, obj = vals
```

这些构造器的存在是为了避免将常见类型包装在 `reflect.Value` 中。例如，`MakeIntSlice` 直接存储 `[]int64`——VM 随后可以在不使用任何反射的情况下进行索引、修改和取地址操作。

### 通用构造器

```go
MakeFromReflect(rv reflect.Value) Value   // → 从 reflect.Kind 自动检测 kind
FromInterface(v any) Value                // → type switch 快速路径，回退到 MakeFromReflect
```

`FromInterface` 使用 type switch 来检测常见的 Go 类型（`bool`、`int`、`int8`、...、`string`、`[]byte`），并将它们路由到特化构造器。不在 switch 中的类型会回退到 `MakeFromReflect`，后者将值包装在 `reflect.Value` 中。

type switch 的排序针对最常见的类型做了优化：

```go
func FromInterface(v any) Value {
    if v == nil { return MakeNil() }
    switch val := v.(type) {
    case bool:    return MakeBool(val)
    case int:     return MakeInt(int64(val))
    case int64:   return MakeInt(val)
    case string:  return MakeString(val)
    case []byte:  return MakeBytes(val)
    // ... 还有 12 种原始类型 ...
    }
    return MakeFromReflect(reflect.ValueOf(v))
}
```

### MakeFromReflect 的解包逻辑

`MakeFromReflect` 不是简单的包装器——它会将原始类型从 `reflect.Value` 中 **解包** 出来：

```go
func MakeFromReflect(rv reflect.Value) Value {
    switch rv.Kind() {
    case reflect.Bool:
        return MakeBool(rv.Bool())          // → KindBool，不存储 reflect
    case reflect.Int, reflect.Int8, ...:
        return MakeInt(rv.Int())            // → KindInt，不存储 reflect
    case reflect.String:
        return MakeString(rv.String())      // → KindString，不存储 reflect
    case reflect.Slice:
        if rv.Type().Elem().Kind() == reflect.Uint8 {
            return MakeBytes(rv.Bytes())    // → KindBytes（针对 []byte）
        }
        return Value{kind: KindReflect, obj: rv}
    default:
        return Value{kind: KindReflect, obj: rv}
    }
}
```

这意味着即使值通过 `reflect.Value` 进入系统，原始类型也会被提取到高效的标签联合体形式。一个包含 `int(42)` 的 `reflect.Value` 会变成 `Value{kind: KindInt, num: 42}`，而不是 `Value{kind: KindReflect, obj: reflect.ValueOf(42)}`。

---

## 访问器

### 类型化访问器（零反射）

```go
v.Bool() bool           // kind != KindBool 时 panic
v.Int() int64           // kind != KindInt 时 panic
v.Uint() uint64         // kind != KindUint 时 panic
v.Float() float64       // kind != KindFloat 时 panic
v.String() string       // kind != KindString 时 panic
v.Complex() complex128  // kind != KindComplex 时 panic
v.Bytes() ([]byte, bool)      // 非 KindBytes 时返回 (nil, false)
v.IntSlice() ([]int64, bool)  // 非原生 []int64 时返回 (nil, false)
v.IntPtr() (*int64, bool)     // 非原生 *int64 时返回 (nil, false)
v.ValueSlice() ([]Value, bool) // 非原生 []Value 时返回 (nil, false)
```

### 快速路径访问器（无 Kind 检查）

供 VM 热路径使用，此时 kind 已从操作码中确定：

```go
v.RawInt() int64    // 直接读 num 字段，无 kind 检查
v.RawBool() bool    // v.num != 0，无 kind 检查
v.RawObj() any      // 直接读 obj 字段，无 kind 检查
```

这些方法被设计为可被 Go 编译器内联。

### 通用访问器

```go
v.Interface() any   // 转换为 Go interface{}，可能调用 reflect.Value.Interface()
```

`Interface()` 方法执行完整的 kind-switch：

- 原始类型：直接返回 Go 类型（`int64`、`float64`、`string` 等）
- `KindFunc`：返回原始 `obj`（闭包或原生函数）
- `KindBytes`：返回 `[]byte`
- `KindReflect` 及其他：对存储的值调用 `reflect.Value.Interface()`

### 反射访问器

```go
v.ReflectValue() (reflect.Value, bool)  // 如果内部存储了 reflect.Value 则返回它
```

### 查询方法

```go
v.Kind() Kind       // 返回 kind 标签
v.IsNil() bool      // KindNil 时为 true，或 KindReflect 包装了 nil 指针/chan/map 等时为 true
v.IsValid() bool    // kind != KindInvalid 时为 true（零值 Value 是 invalid 的）
v.CanInterface() bool // 内部的 reflect.Value（如果有）支持 Interface() 时为 true
```

---

## 算术运算

所有算术运算都使用 kind-switch，整数快速路径优先：

```go
v.Add(other Value) Value    // +（int、uint、float、字符串拼接、complex）
v.Sub(other Value) Value    // -（int、uint、float、complex）
v.Mul(other Value) Value    // *（int、uint、float、complex）
v.Div(other Value) Value    // /（int、uint、float、complex）
v.Mod(other Value) Value    // %（int、uint、float）
v.Neg() Value               // 一元 -（int、float、complex）
```

### 工作原理

以 `Add` 为例：

```go
func (v Value) Add(other Value) Value {
    switch v.kind {
    case KindInt:
        return MakeInt(v.num + other.Int())     // 原生 int64 加法
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

对于 `KindInt`，整个操作是：
1. 读 `v.num`（8 字节）
2. 读 `other.num`（8 字节，通过 `Int()` 检查 `kind == KindInt`）
3. 相加（一条 CPU 指令）
4. 将结果写入栈上的新 Value 结构体

无堆分配、无函数指针间接调用、无反射。

### VM 优化

在 VM 分发循环中，算术运算通过内联快速路径进一步优化：

```go
case bytecode.OpAdd:
    sp--
    b := stack[sp]
    sp--
    a := stack[sp]
    if a.IsInt() {
        stack[sp] = value.MakeInt(a.RawInt() + b.RawInt())  // 约 3 条机器指令
    } else {
        stack[sp] = a.Add(b)  // 通用路径
    }
    sp++
```

`a.IsInt()` 即 `a.kind == KindInt`——单字节比较。`a.RawInt()` 即 `a.num`——单字段读取。整个整数加法约 3 条机器指令加上分发开销。

---

## 比较与相等

```go
v.Equal(other Value) bool   // ==（所有类型；复合类型回退到 reflect.DeepEqual）
v.Cmp(other Value) int      // 返回 -1、0 或 1（bool、int、uint、float、string）
```

### 相等性语义

`Equal` 遵循 Go 的相等性规则：

1. **不同 kind**：返回 `false`（nil 比较除外：`KindNil` 与包装了 nil 值的 `KindReflect` 的比较）。
2. **相同 kind，原始类型**：直接字段比较（`v.num == other.num`）。
3. **相同 kind，复合类型**：`reflect.DeepEqual(v.Interface(), other.Interface())`。

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

### 比较语义

`Cmp` 返回三路比较结果。支持 `bool`、`int`、`uint`、`float` 和 `string`——具有自然排序的类型。对不支持的类型（map、slice 等）会 panic。

---

## 位运算

```go
v.And(other Value) Value      // &（int、uint）
v.Or(other Value) Value       // |（int、uint）
v.Xor(other Value) Value      // ^（int、uint）
v.AndNot(other Value) Value   // &^（int、uint）
v.Lsh(n uint) Value           // <<（int、uint）
v.Rsh(n uint) Value           // >>（int、uint）
```

全部直接操作 `KindInt` 和 `KindUint` 的 `num` 字段。移位量始终为 `uint`（与 Go 规范一致）。

---

## 类型转换

```go
v.ToInt() Value     // → KindInt（从 bool、int、uint、float）
v.ToUint() Value    // → KindUint（从 bool、int、uint、float）
v.ToFloat() Value   // → KindFloat（从 bool、int、uint、float）
v.ToBool() Value    // → KindBool（从 bool、int、uint、float、string、nil）
v.ToString() Value  // → KindString（使用 fmt.Sprintf("%v", v.Interface())）
```

这些对应 Go 中的类型转换表达式，如 `int(x)`、`float64(x)` 等。它们在标签联合体内部操作——数值转换无需分配内存。

`ToBool` 有扩展语义：零值、空字符串和 nil 返回 `false`；其他情况返回 `true`。

### 反射级别的转换

```go
v.ToReflectValue(typ reflect.Type) reflect.Value
```

这是从 `Value` 到 `reflect.Value` 的桥梁，当 VM 需要将值传递给外部函数或存储到类型化容器时使用。它处理：

- **原始类型**：`reflect.ValueOf(v.Int()).Convert(typ)` —— 适配目标类型（如 `int64` → `int32`）
- **闭包**（`KindFunc`）：通过 `reflect.MakeFunc(typ, ...)` 将闭包包装为真正的 Go 函数
- **原生 int 切片**（`KindSlice` 含 `[]int64`）：逐元素转换为目标切片类型
- **原生 value 切片**（`KindSlice` 含 `[]Value`）：递归转换每个元素
- **reflect 值**：直接返回存储的 `reflect.Value`

---

## 容器操作

### 长度和容量

```go
v.Len() int    // string、slice、array、map、chan
v.Cap() int    // slice、array、chan
```

两者都包含原生 `[]int64` 切片的快速路径（直接调用 `len()`/`cap()`）。其他类型委托给 `reflect.Value.Len()`/`reflect.Value.Cap()`。

### 索引

```go
v.Index(i int) Value          // 读取索引 i 处的元素（string、slice、array）
v.SetIndex(i int, val Value)  // 设置索引 i 处的元素（slice、array）
```

**Index 快速路径**：
- `KindString`：返回 `MakeUint(uint64(s[i]))` —— 位置处的字节，符合 Go 语义。
- `KindSlice` 含 `[]int64`：返回 `MakeInt(s[i])` —— 直接数组访问，零反射。
- `KindSlice` 含 `[]Value`：返回 `slice[i]` —— 直接访问。
- 其他切片/数组：`reflect.Value.Index(i)` → `MakeFromReflect`。

**SetIndex 快速路径**：
- `KindSlice` 含 `[]int64`：`s[i] = val.RawInt()` —— 直接写入。
- 其他：`reflect.Value.Index(i).Set(val.ToReflectValue(elemType))`。

### Map 操作

```go
v.MapIndex(k Value) Value             // 读取键 k 处的值
v.SetMapIndex(k, val Value)           // 设置键 k 处的值（val.IsNil() → 删除）
v.MapIter(f func(key, val Value) bool) // 遍历 map 条目
```

所有 map 操作都通过 `reflect.Value`，因为 Go 的 map 类型需要动态派发。`SetMapIndex` 传入 nil 值时会删除键，与 Go 的 `delete(m, k)` 语义一致。

### 结构体字段访问

```go
v.Field(i int) Value              // 读取索引 i 处的结构体字段
v.SetField(i int, val Value)      // 设置索引 i 处的结构体字段
```

两者都需要 `reflect.Value`，使用 `reflect.Value.Field(i)`。字段索引 `i` 在编译时从结构体类型信息中解析。

---

## 指针操作

```go
v.Elem() Value           // 解引用：*ptr → value
v.SetElem(val Value)     // 通过指针赋值：*ptr = val
v.Pointer() uintptr      // 原始指针地址（用于身份比较）
```

### Elem 快速路径

`Elem()` 在回退到反射之前有三条快速路径：

1. **`*int64` 指针**（来自原生 int 切片的 `OpIndexAddr`）：`return MakeInt(*ptr)`
2. **`*Value` 指针**（来自局部变量的 `OpAddr`，闭包的 `OpFree`）：`return *ptr`
3. **`reflect.Value` 中的 `*Value`**（指向 Value 结构体的指针）：直接解包

### SetElem 快速路径

`SetElem()` 与 `Elem()` 对应：

1. **`*int64`**：`*ptr = val.num` —— 直接内存写入
2. **`*Value`**：`*ptr = val` —— 直接结构体赋值
3. **`reflect.Value` Ptr**：处理类型转换、函数包装和原生切片转换

### 辅助函数

```go
UnsafeAddrOf(v reflect.Value) unsafe.Pointer
```

VM 内部使用，用于获取未导出结构体字段的可设置指针。

---

## 通道操作

```go
v.Send(val Value)                                    // ch <- val（阻塞）
v.SendContext(ctx context.Context, val Value) error   // ch <- val（支持取消）
v.TrySend(val Value) bool                            // 非阻塞发送
v.Recv() (Value, bool)                               // <-ch（阻塞）
v.RecvContext(ctx context.Context) (Value, bool, error) // <-ch（支持取消）
v.TryRecv() (Value, bool)                            // 非阻塞接收
v.Close()                                            // close(ch)
```

所有通道操作都通过 `reflect.Value`，因为 Go 通道需要运行时支持。带 `*Context` 的变体使用 `reflect.Select` 配合上下文取消通道——这是无法避免反射的场景之一。

### 上下文感知模式

`SendContext` 和 `RecvContext` 实现了两阶段策略：

1. **快速路径**：非阻塞尝试（`TrySend`/`TryRecv`）。如果通道有缓冲空间或有等待的 goroutine，这会立即成功。
2. **慢速路径**：`reflect.Select` 带两个 case——通道操作和 `ctx.Done()`。这实现了取消而不泄漏 goroutine。

---

## 反射互操作

值系统提供了 `Value` 和 `reflect.Value` 之间的双向转换：

### Value → reflect.Value

```go
v.ToReflectValue(typ reflect.Type) reflect.Value
```

这是 VM 需要与 Go 反射系统交互时（外部函数调用、通道操作、结构体字段访问）使用的主要转换方法。它处理所有 kind，对以下情况有特殊支持：

- **类型转换**：`MakeInt(42).ToReflectValue(reflect.TypeOf(int32(0)))` 产生 `reflect.ValueOf(int32(42))`
- **闭包包装**：通过 `reflect.MakeFunc` 将解释器闭包转换为真正的 `func` 值
- **切片转换**：将 `[]int64` 和 `[]Value` 转换为类型化切片

### reflect.Value → Value

```go
MakeFromReflect(rv reflect.Value) Value
```

从 `reflect.Value` 中提取原始类型到高效的标签联合体形式。复合类型保持包装在 `KindReflect` 下的 `reflect.Value` 中。

### Interface 往返

```go
FromInterface(v any) Value    // any → Value
v.Interface() any             // Value → any
```

这构成了解释器与宿主 Go 程序之间的桥梁。`FromInterface` 是所有外部数据的入口；`Interface()` 是返回结果的出口。

---

## 闭包支持

Gig 中的闭包表示为 `KindFunc` 值，闭包对象存储在 `obj` 中。值系统包含一个 **回调机制**，用于打破 `value` 和 `vm` 包之间的循环依赖：

```go
type ClosureCaller func(closure any, args []reflect.Value, outTypes []reflect.Type) []reflect.Value

RegisterClosureCaller(caller ClosureCaller)
```

VM 在初始化时注册一个 `ClosureCaller`。当 `ToReflectValue` 需要将闭包转换为类型化的 `func`（例如存储在 `map[string]func() int` 中）时，它使用 `reflect.MakeFunc` 配合一个调用已注册 caller 的处理器：

```go
fn := reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
    results := closureCaller(closure, args, outTypes)
    // 将结果转换为匹配预期的返回类型
    ...
    return out
})
```

`outTypes` 参数支持嵌套闭包的递归包装——返回闭包的闭包可以在每一层被正确类型化。

---

## 原生快速路径

值系统为常见类型提供了特化表示以避免反射：

### `[]byte`（KindBytes）

```go
MakeBytes(b []byte) Value
v.Bytes() ([]byte, bool)
```

`[]byte` 是最常见的二进制类型（JSON 负载、protobuf 等）。原生存储意味着带 `[]byte` 参数的外部函数调用完全避免 `reflect.ValueOf`。

### `[]int64`（KindSlice 含原生存储）

```go
MakeIntSlice(s []int64) Value
v.IntSlice() ([]int64, bool)
```

整数切片在计算代码中极其常见。原生存储意味着：
- `v.Index(i)` → 直接数组访问
- `v.SetIndex(i, val)` → 直接数组写入
- `v.Len()` → `len(s)`
- VM 的 `OpIntSliceGet`/`OpIntSliceSet` 超级指令直接操作 `[]int64`

### `*int64`（KindPointer 含原生存储）

```go
MakeIntPtr(p *int64) Value
v.IntPtr() (*int64, bool)
```

由原生 int 切片的 `OpIndexAddr` 使用——取 `s[i]` 的地址返回 `*int64` 指针，之后的 `OpDeref`/`OpSetDeref` 是直接内存操作。

### `*Value`（KindPointer 含原生存储）

闭包的自由变量是 `*Value` 指针。对它们的 `Elem()` 和 `SetElem()` 是直接的结构体解引用。

### `[]Value`（KindSlice 含原生存储）

```go
MakeValueSlice(vals []Value) Value
v.ValueSlice() ([]Value, bool)
```

用于 DirectCall 包装器中的多返回值打包。返回 `(int, error)` 的函数将结果打包为 `MakeValueSlice([]Value{MakeInt(r0), FromInterface(r1)})`。

---

## 设计思考

### 为什么不是所有值都用 `reflect.Value`？

另一种方案（Yaegi 和其他 Go 解释器使用的）是将每个值都表示为 `reflect.Value`。这更简单，但有显著的代价：

| 操作 | `reflect.Value` | Gig `Value` |
|------|----------------|-------------|
| 创建一个 int | 堆分配 + `reflect.ValueOf`（约 15-30 ns） | 结构体字面量（0 ns，0 分配） |
| 整数加法 | `reflect.Value.Int()` + `reflect.ValueOf`（约 40 ns，1 次分配） | `v.num + w.num`（约 1 ns，0 分配） |
| 比较两个 int | `reflect.Value.Int()` × 2（约 20 ns） | `v.num == w.num`（约 0.5 ns） |
| 传递给外部函数 | 已经是 `reflect.Value`（0 开销） | `ToReflectValue(typ)`（约 5-15 ns） |

标签联合体用外部调用边界的复杂性（需要 `ToReflectValue`）换取了热路径上的巨大加速（算术、比较、局部变量访问）。

### 为什么是 32 字节？

- 更小（24 字节）需要将 kind 标签位打包到 `num` 中，使访问器逻辑复杂化。
- 更大（40+ 字节）会浪费缓存空间——VM 栈持有数百个值，缓存局部性至关重要。
- 32 字节是 **2 个缓存行片段**（64 字节缓存行恰好容纳 2 个值），为相邻栈访问提供良好的空间局部性。

### 为什么切片同时有 `KindSlice` 和 `KindReflect`？

`KindSlice` 用于原生 `[]int64` 和 `[]Value` 存储。通用的 `[]string` 或 `[]MyStruct` 通过 `KindReflect`。这种区分使 VM 可以为整数切片（计算代码中最常见的情况）走快速路径，而不会影响其他切片类型。

### 为什么 `KindBytes` 是独立的 Kind

`[]byte` 可以作为 `KindSlice` 存储在 `obj` 中，但拥有专用的 kind 使得：
- `Bytes()` 访问器可以用 O(1) 的类型检查
- DirectCall 包装器可以在不使用 `reflect.ValueOf` 的情况下提取 `[]byte`
- `MakeBytes` 构造器可以是纯粹的结构体赋值

---

## 性能特性

### 零分配操作

以下操作不产生堆分配：

- 所有原始类型构造器（`MakeInt`、`MakeBool`、`MakeFloat` 等）
- `MakeString`、`MakeBytes`、`MakeFunc`（引用类型，`obj` 中仅存头部）
- 原始类型的所有算术运算（`Add`、`Sub`、`Mul`、`Div`、`Mod`、`Neg`）
- 原始类型的所有比较（`Equal`、`Cmp`）
- 所有位运算
- 数值类型间的所有类型转换
- 类型化访问器（`Int()`、`Bool()`、`Float()`、`String()`）
- 原生 `[]int64` 上的 `Index` 和 `SetIndex`
- 原生 `*int64` 和 `*Value` 上的 `Elem` 和 `SetElem`

### 可能产生分配的操作

- 非基础类型的 `FromInterface`（调用 `reflect.ValueOf`）
- `KindReflect` 上的 `Interface()`（调用 `reflect.Value.Interface()`）
- 闭包的 `ToReflectValue`（使用 `reflect.MakeFunc`）
- 所有 map 操作（通过 `reflect.Value`）
- 所有通道操作（通过 `reflect.Value`）
- 复合类型的 `Equal`（使用 `reflect.DeepEqual`）

### 基准测试

| 操作 | 耗时 | 分配次数 |
|------|------|---------|
| `MakeInt(42)` | 约 0.3 ns | 0 |
| `MakeBytes(b)` | 约 0.3 ns | 0 |
| `FromInterface(42)` | 约 2 ns | 0 |
| `FromInterface([]int{1,2,3})` | 约 90 ns | 3 |
| `v.Int()` | 约 0.5 ns | 0 |
| `v.Interface()`（KindInt） | 约 1 ns | 0 |
| `v.Interface()`（KindReflect） | 约 15 ns | 0-1 |
| `MakeInt(a).Add(MakeInt(b))` | 约 2 ns | 0 |
| `v.ToReflectValue(intType)` | 约 10 ns | 1 |

---

## 文件组织

value 包按职责分为 6 个文件：

| 文件 | 职责 | 关键类型/函数 |
|------|------|-------------|
| `value.go` | 核心类型定义、构造器、查询方法 | `Value`、`Kind`、`Make*`、`FromInterface`、`MakeFromReflect` |
| `accessor.go` | 类型化访问器、`Interface()`、`ToReflectValue`、闭包支持 | `Bool()`、`Int()`、`Interface()`、`ToReflectValue()`、`ClosureCaller` |
| `arithmetic.go` | 算术、比较、位运算 | `Add`、`Sub`、`Mul`、`Div`、`Cmp`、`Equal`、`And`、`Or`、`Lsh`、`Rsh` |
| `convert.go` | 类型转换方法 | `ToInt`、`ToUint`、`ToFloat`、`ToBool`、`ToString` |
| `container.go` | 容器操作（len、cap、index、map、field、pointer） | `Len`、`Cap`、`Index`、`SetIndex`、`MapIndex`、`Field`、`Elem`、`SetElem` |
| `channel.go` | 支持上下文的通道操作 | `Send`、`Recv`、`TrySend`、`TryRecv`、`SendContext`、`RecvContext`、`Close` |
| `value_test.go` | 单元测试 | 覆盖构造器、算术、比较、转换、边界情况 |

---

## 参见

- [`docs/gig-internals.md`](gig-internals.md) — 架构概览，包含值系统的上下文
- [`docs/gig-internals_CN.md`](gig-internals_CN.md) — 架构概览（中文版）
- [`docs/optimization-zero-reflection.md`](optimization-zero-reflection.md) — 如何利用值系统实现零反射外部调用
- [`docs/optimization-int-specialization.md`](optimization-int-specialization.md) — 使用 `intLocals[]` 影子数组的整数特化
- [`docs/optimization-directcall.md`](optimization-directcall.md) — 通过类型化访问器提取值的 DirectCall 代码生成
