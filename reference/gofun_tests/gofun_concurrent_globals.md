# gofun 并发全局变量问题分析

## 概述

新版 gofun (`git.code.oa.com/datacenter/onefun/gofun`) 基于 SSA 编译 + 寄存器解释器，修复了旧版多个 bug，但**并发全局变量访问仍不安全**。

## 新旧 gofun 对比

| 特性 | 旧 gofun (AST) | 新 gofun (SSA) |
|------|---------------|----------------|
| 架构 | AST-Walking 解释器 | SSA 编译 + 寄存器解释器 |
| API | `Parse(src, nil)` → `Run(scope)` | `Build(src)` → `Run(funcName, args...)` |
| 整数溢出 (Bug #1) | ❌ `int(val.(int64))` 溢出 | ✅ 使用 `go/constant` 正确处理 |
| 短路求值 (Bug #3) | ❌ 先求值所有子表达式 | ✅ SSA 生成分支跳转 |
| Map v,ok (Bug #4) | ❌ 不返回存在标志 | ✅ 使用 `ssa.Lookup.CommaOk` |
| 并发安全 (Bug #6) | ❌ Scope map 无锁 | ❌ globals map 无锁 |
| sync 包 | ❌ 未注册 | ✅ 注册了类型但仍无并发保护 |

## 并发问题根因

### 新 gofun 的 Program 结构

```go
// gofun/interpreter.go
type Program struct {
    mainPkg *ssa.Package
    globals map[ssa.Value]*value.Value  // 无锁保护!
}
```

### Run() 方法

```go
func (p *Program) Run(funcName string, params ...interface{}) (interface{}, error) {
    val, _, err := p.RunWithContext(funcName, params...)
    return val, err
}

func (p *Program) RunWithContext(funcName string, params ...interface{}) (...) {
    fr := &frame{
        program: p,  // 所有 goroutine 共享同一 p.globals
        context: ctx,
    }
    ret := callSSA(fr, mainFn, args, nil)
    ...
}
```

### 全局变量读取 (frame.go)

```go
func (fr *frame) get(key ssa.Value) value.Value {
    case *ssa.Global:
        if r, ok := fr.program.globals[key]; ok {  // 无锁读 map!
            v := (*r).Interface()
            return value.ValueOf(&v)
        }
}
```

### 触发方式

```go
program, _ := gofun.Build(src)
// 多 goroutine 并发调用 → fatal error: concurrent map writes
go program.Run("Increment")
go program.Run("Increment")
// 实测: fatal error 无法被 recover 捕获
```

## Gig 的解决方案

| 机制 | 说明 |
|------|------|
| `SharedGlobals` | `sync.RWMutex` 保护的全局变量切片 |
| `GlobalRef` | 锁代理，每次读写通过锁 |
| `VMPool` | 每个 `Run()` 获取独立 VM 实例 |
| `GoroutineTracker` | goroutine 生命周期管理 |
| 值类型 Mutex | `var mu sync.Mutex` 存储为 `*sync.Mutex` |
| Closure SharedGlobals | 闭包传递给外部函数时共享全局变量 (v1.6.0) |

## 功能对比

| 功能 | 新 gofun | Gig |
|------|---------|-----|
| 并发 `Run()` 安全 | ❌ fatal: concurrent map writes | ✅ VMPool + SharedGlobals |
| `sync.Mutex` 类型注册 | ✅ 已注册 | ✅ DirectCall 包装 |
| `sync.WaitGroup` | ✅ 已注册 | ✅ DirectCall 包装 |
| `sync.Once.Do(闭包)` | 未验证 | ✅ 已修复 (v1.6.0) |
| `sync.Map` | ✅ 已注册 | ✅ DirectCall 包装 |
| goroutine 追踪 | atomic 计数 | ✅ GoroutineTracker |
| goroutine 数量限制 | ❌ 无 | ✅ 可配置 |
| 通道发送 | ✅ `reflect.Send` (阻塞) | ✅ 阻塞 + Context |
| `for` 循环 goroutine 闭包 | 未验证 | ✅ 已修复 (v1.6.0) |

## 测试运行

```bash
cd reference/gofun_tests

# Bug 验证测试 (包含新 gofun)
go test -run TestGofunVerify -v

# 性能对比
go test -bench=. -benchmem

# Gig 完整并发测试
cd ../.. && go test ./tests/ -run 'TestValueTypeMutexExact|TestRWMutex|TestOnce|TestSyncMap|TestWaitGroupGlobal' -v
```
