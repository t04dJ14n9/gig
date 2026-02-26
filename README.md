# Gig - Go Interpreter in Go

Gig is a high-performance Go interpreter written in Go, featuring SSA-to-bytecode compilation and a stack-based virtual machine.

## Features

- **SSA-based compilation**: Uses `golang.org/x/tools/go/ssa` for intermediate representation
- **Stack-based VM**: Efficient bytecode execution with minimal overhead
- **Tagged-union Value system**: Zero reflection overhead for primitive types
- **Security**: Bans `unsafe`, `reflect`, and `panic` in interpreted code
- **Extensible**: Support for registering external Go packages (40+ stdlib packages built-in)

## Installation

```bash
go get github.com/t04dJ14n9/gig
```

## Quick Start

### Option 1: Use Built-in Standard Library (Recommended)

Gig comes with 40+ standard library packages pre-registered. Just import `gig/packages`:

```go
package main

import (
    "fmt"
    _ "gig/packages" // Import gig's built-in stdlib
    "gig"
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

**Built-in packages include**: `fmt`, `strings`, `strconv`, `math`, `time`, `bytes`, `errors`, `sort`, `regexp`, `encoding/json`, `encoding/base64`, `net/url`, and 30+ more.

### Option 2: Use Custom Dependencies

If you need third-party libraries or a subset of standard library, use the `gig` CLI tool:

#### Step 1: Install the CLI

```bash
# Install the CLI tool
go install github.com/t04dJ14n9/gig/cmd/gig@latest

# Or run directly (Go 1.21+)
go run github.com/t04dJ14n9/gig/cmd/gig@latest --help
```

#### Step 2: Initialize a dependency package

```bash
# Create a dependency package named "mydep"
gig init -package mydep
```

This creates:
```
mydep/
└── pkgs.go    # Edit this to add/remove packages
```

#### Step 3: Customize dependencies

Edit `mydep/pkgs.go` to add third-party libraries:

```go
package mydep

import (
    // Standard library (keep what you need)
    _ "fmt"
    _ "strings"
    _ "time"

    // Third-party libraries
    _ "github.com/spf13/cast"
    _ "github.com/tidwall/gjson"
)
```

#### Step 4: Generate registration code

```bash
# Generate registration code from pkgs.go
gig gen ./mydep
```

This generates:
```
mydep/
├── pkgs.go
└── packages/
    ├── fmt.go
    ├── strings.go
    ├── github_com_spf13_cast.go
    └── github_com_tidwall_gjson.go
```

#### Step 5: Use in your program

```go
package main

import (
    "fmt"
    _ "myapp/mydep/packages" // Your custom dependency package
    "gig"
)

func main() {
    source := `
package main

import "github.com/tidwall/gjson"

func GetJsonValue(json string, path string) string {
    return gjson.Get(json, path).String()
}
`

    prog, _ := gig.Build(source)
    result, _ := prog.Run("GetJsonValue", `{"name":"Alice"}`, "name")
    fmt.Println(result) // Output: Alice
}
```

## API Reference

### Building and Running

```go
// Build parses and compiles Go source code
prog, err := gig.Build(source string) (*Program, error)

// Run executes a function by name
result, err := prog.Run(funcName string, args ...interface{}) (interface{}, error)

// RunWithContext executes with context for cancellation (ctx is first parameter)
result, err := prog.RunWithContext(ctx context.Context, funcName string, args ...interface{}) (interface{}, error)
```

### Registering Packages (Advanced)

```go
import "gig/register"

// Register a package manually (usually done via generated code)
pkg := register.RegisterPackage("mypkg", "mypkg")
pkg.AddFunction("MyFunc", MyFunc, "", directCall_MyFunc)
pkg.AddConstant("MyConst", MyConst, "")
pkg.AddVariable("MyVar", &MyVar, "")
pkg.AddType("MyType", reflect.TypeOf(MyType{}), "")
```

## Examples

See the `examples/` directory:

- **`examples/simple/`** - Using gig with built-in stdlib (easiest)
- **`examples/custom/`** - Using gig with custom dependencies

Run examples:

```bash
# Simple example (uses built-in stdlib)
cd gig/examples/simple
go run main.go

# Custom example
cd gig/examples/custom
go run main.go
```

## gig CLI Commands

```bash
# Initialize a dependency package
gig init -package <name>

# Generate registration code
gig gen <dir>

# Examples
gig init -package mydep         # Creates mydep/pkgs.go
gig gen ./mydep                 # Generates registration code in myapp/mydep/packages/
```

## Supported Features

- ✅ Arithmetic operations
- ✅ Variables and assignments
- ✅ Control flow (if/else, for loops, switch)
- ✅ Functions and recursion
- ✅ Multiple return values
- ✅ Closures
- ✅ String operations
- ✅ Slices and arrays
- ✅ Maps
- ✅ Structs and methods
- ✅ Interfaces
- ✅ Goroutines (basic)
- ✅ Context-based timeouts
- ✅ External Go function calls

## Security

Gig enforces security by banning certain imports:
- `unsafe` - Memory safety
- `reflect` - Type safety
- `panic` usage - Controlled execution

## Architecture

### Compiler
- Parses Go source code using `go/parser`
- Type-checks using `go/types`
- Builds SSA representation using `golang.org/x/tools/go/ssa`
- Compiles SSA to custom bytecode (~70 opcodes)

### Virtual Machine
- Stack-based execution
- Support for functions, closures, and recursion
- Context-based timeout control
- Concurrent execution support

### Value System
- Tagged-union design for efficient primitive operations
- Direct operations for `int`, `float64`, `string`, `bool`
- Fallback to `interface{}` for complex types

## Changelog

### v0.2.0 - External Type Method Support

**Bug Fix**: Methods on external (registered) types are now fully supported.

Previously, calling methods on external types like `gjson.Get(json, path).String()` would fail with a type-check error because methods were not registered on `types.Named` types. This has been fixed across three layers:

- **`importer/importer.go`**: Added `addMethodsToNamed()` — when converting `reflect.Type` to `types.Named`, all exported methods (both value and pointer receivers) are now enumerated and added via `named.AddMethod()`. This allows the Go type checker to resolve method calls on external types.

- **`compiler/compiler.go`**: Added `ExternalMethodInfo` and updated `compileExternalStaticCall` to detect method calls (`sig.Recv() != nil`) and emit them with method dispatch metadata instead of looking up a static function object.

- **`vm/vm.go`**: Added `callExternalMethod()` which dispatches method calls on external types via `reflect.Value.MethodByName()`, handling variadic arguments, pointer receivers, and multi-return values.

## License

MIT License
