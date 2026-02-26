# Gig Examples

[![中文](https://img.shields.io/badge/lang-中文-red.svg)](README_CN.md) [![English](https://img.shields.io/badge/lang-English-blue.svg)](README.md)

This directory contains examples demonstrating how to use Gig, a high-performance Go interpreter.

## Examples Overview

| Example | Description | Difficulty |
|---------|-------------|------------|
| [simple](./simple) | Using Gig with built-in standard library | Beginner |
| [custom](./custom) | Using Gig with custom/third-party dependencies | Intermediate |

## Quick Start

### Simple Example (Recommended for First-Time Users)

The simplest way to use Gig is with its built-in standard library support. Just import `gig/packages` and you have access to 40+ standard library packages.

```bash
cd simple
go run main.go
```

### Custom Example (For Third-Party Libraries)

When you need third-party libraries or want to minimize dependencies, use the CLI tool to generate registration code.

```bash
cd custom
go run main.go
```

---

## Simple Example Details

**Location:** `./simple/`

**Use Case:** Quick prototyping, scripting, rule engines with standard library only.

**Key Features Demonstrated:**
- Basic computation with loops and variables
- Using standard library packages (`fmt`, `strings`, `math`, `time`)
- Context-based timeout control
- Multi-function programs

### Code Structure

```
simple/
├── go.mod
├── go.sum
└── main.go
```

### Usage Pattern

```go
import (
    "gig"
    _ "gig/packages" // Import built-in stdlib (40+ packages)
)

func main() {
    source := `
    package main
    
    import "fmt"
    import "strings"
    
    func Greet(name string) string {
        return fmt.Sprintf("Hello, %s!", strings.ToUpper(name))
    }
    `
    
    prog, err := gig.Build(source)
    if err != nil {
        panic(err)
    }
    
    result, err := prog.Run("Greet", "world")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result) // Output: Hello, WORLD!
}
```

### Built-in Packages

Gig's built-in standard library includes:

| Category | Packages |
|----------|----------|
| **I/O** | `fmt`, `io`, `bufio`, `bytes`, `strings` |
| **Encoding** | `encoding/json`, `encoding/base64`, `encoding/hex`, `encoding/xml` |
| **Text** | `strings`, `strconv`, `text/template`, `regexp` |
| **Math** | `math`, `math/rand` |
| **Time** | `time` |
| **Collections** | `sort`, `container/list`, `container/heap` |
| **Network** | `net/url`, `net/http` (partial) |
| **Crypto** | `crypto/md5`, `crypto/sha1`, `crypto/sha256` |
| **Other** | `errors`, `sync`, `context`, `path`, `path/filepath`, `os` (partial) |

### Context Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

result, err := prog.RunWithContext(ctx, "SlowFunction", args...)
if errors.Is(err, context.DeadlineExceeded) {
    // Handle timeout
}
```

---

## Custom Example Details

**Location:** `./custom/`

**Use Case:** Production applications requiring third-party libraries, minimal dependency footprint, or custom package subsets.

**Key Features Demonstrated:**
- Using third-party libraries (github.com/tidwall/gjson)
- Method calls on external types (`.String()`, `.Int()`, `.Bool()`)
- Custom dependency package generation workflow

### Code Structure

```
custom/
├── go.mod
├── go.sum
├── main.go
└── mydep/
    ├── pkgs.go              # Dependency declarations
    └── packages/            # Generated registration code
        ├── fmt.go
        ├── strings.go
        ├── github_com_tidwall_gjson.go
        └── ...
```

### Setup Workflow

#### Step 1: Install the CLI

```bash
go install github.com/t04dJ14n9/gig/cmd/gig@latest
```

#### Step 2: Initialize a Dependency Package

```bash
gig init -package mydep
```

This creates:
```
mydep/
└── pkgs.go    # Edit this to declare dependencies
```

#### Step 3: Declare Dependencies

Edit `mydep/pkgs.go`:

```go
package mydep

import (
    // Standard library (keep what you need)
    _ "fmt"
    _ "strings"
    _ "time"

    // Third-party libraries
    _ "github.com/tidwall/gjson"
    _ "github.com/spf13/cast"
)
```

#### Step 4: Generate Registration Code

```bash
gig gen ./mydep
```

This generates:
```
mydep/
├── pkgs.go
└── packages/
    ├── fmt.go
    ├── strings.go
    ├── time.go
    ├── github_com_tidwall_gjson.go
    └── github_com_spf13_cast.go
```

#### Step 5: Use in Your Program

```go
import (
    "gig"
    _ "myapp/mydep/packages" // Your custom dependency package
)

func main() {
    source := `
    package main
    
    import "github.com/tidwall/gjson"
    
    func GetValue(json string, path string) string {
        return gjson.Get(json, path).String()
    }
    `
    
    prog, _ := gig.Build(source)
    result, _ := prog.Run("GetValue", `{"name":"Alice"}`, "name")
    fmt.Println(result) // Output: Alice
}
```

### External Type Methods

Gig supports calling methods on external types:

```go
source := `
package main

import "github.com/tidwall/gjson"

func GetUserAge(json string) int64 {
    return gjson.Get(json, "age").Int()  // Method call on gjson.Result
}
`
```

Supported method patterns:
- Value receiver methods: `.String()`, `.Int()`, `.Bool()`, `.Float()`
- Pointer receiver methods: `.Scan()`, `.ForEach()`
- Chained method calls: `gjson.Get(json, "arr").Array()[0].String()`

---

## Comparison: Simple vs Custom

| Aspect | Simple | Custom |
|--------|--------|--------|
| **Setup** | None (just import) | CLI tool + code generation |
| **Dependencies** | 40+ stdlib packages | Only what you need |
| **Third-party libs** | Not supported | Fully supported |
| **Binary size** | Larger | Smaller (tree-shaking) |
| **Build time** | Faster (pre-built) | Slower (code generation) |
| **Use case** | Prototyping, scripts | Production, embedded |

---

## API Reference

### Building

```go
// Build parses and compiles Go source code
prog, err := gig.Build(source string) (*Program, error)
```

### Running

```go
// Run executes a function by name (default 10s timeout)
result, err := prog.Run(funcName string, args ...interface{}) (interface{}, error)

// RunWithContext executes with custom context for cancellation
result, err := prog.RunWithContext(ctx context.Context, funcName string, args ...interface{}) (interface{}, error)
```

### Package Registration

```go
import "gig/register"

pkg := register.AddPackage("mypkg", "mypkg")
pkg.NewFunction("MyFunc", MyFunc, "documentation")
pkg.NewConst("MyConst", MyConst, "documentation")
pkg.NewVar("MyVar", &MyVar, "documentation")
```

---

## Security Considerations

Gig enforces security by banning certain imports in interpreted code:

| Banned | Reason |
|--------|--------|
| `unsafe` | Memory safety |
| `reflect` | Type safety |
| `panic` | Controlled execution |

---

## Supported Language Features

| Feature | Status |
|---------|--------|
| Arithmetic operations | Fully supported |
| Variables and assignments | Fully supported |
| Control flow (if/else, for, switch) | Fully supported |
| Functions and recursion | Fully supported |
| Multiple return values | Fully supported |
| Closures | Fully supported |
| String operations | Fully supported |
| Slices and arrays | Fully supported |
| Maps | Fully supported |
| Structs and methods | Fully supported |
| Interfaces | Fully supported |
| Goroutines | Basic support |
| Channels | Basic support |
| Context-based timeouts | Fully supported |
| External Go function calls | Fully supported |

---

## Troubleshooting

### "package not registered"

**Problem:** You're importing a package that Gig doesn't know about.

**Solution:**
- For stdlib: Make sure you imported `_ "gig/packages"`
- For third-party: Use the CLI to generate registration code

### "method not found on external type"

**Problem:** Calling a method on an external type fails.

**Solution:** Ensure the method is exported and the type is properly registered via `gig gen`.

### Timeout errors

**Problem:** Function execution times out.

**Solution:**
- Increase timeout: `context.WithTimeout(ctx, longerDuration)`
- Check for infinite loops in your interpreted code

### Global variable issues

**Problem:** Global variables aren't working as expected.

**Solution:** This is a known limitation. Use function-local variables or pass data as arguments.
