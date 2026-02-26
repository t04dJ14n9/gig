# Gig - Go Interpreter in Go

Gig is a high-performance Go interpreter written in Go, featuring SSA-to-bytecode compilation and a stack-based virtual machine.

## Features

- **SSA-based compilation**: Uses `golang.org/x/tools/go/ssa` for intermediate representation
- **Stack-based VM**: Efficient bytecode execution with minimal overhead
- **Tagged-union Value system**: Zero reflection overhead for primitive types
- **Security**: Bans `unsafe`, `reflect`, and `panic` in interpreted code
- **Extensible**: Support for registering external Go packages

## Installation

```bash
go get github.com/t04dJ14n9/GIG
```

## Quick Start

```go
package main

import (
    "fmt"
    "gig"
)

func main() {
    source := `
package main

func Compute() int {
    sum := 0
    for i := 1; i <= 10; i++ {
        sum = sum + i
    }
    return sum
}
`
    prog, err := gig.Build(source)
    if err != nil {
        panic(err)
    }

    result, err := prog.Run("Compute")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Result: %v\n", result) // Output: Result: 55
}
```

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
- Fallback to `reflect.Value` for complex types

## Supported Features

- ✅ Arithmetic operations
- ✅ Variables and assignments
- ✅ Control flow (if/else, for loops)
- ✅ Functions and recursion
- ✅ Multiple return values
- ✅ String operations
- ✅ Context-based timeouts
- 🚧 Slices and arrays (in progress)
- 🚧 Maps (in progress)
- 🚧 External function calls (in progress)

## Security

Gig enforces security by banning certain imports:
- `unsafe` - Memory safety
- `reflect` - Type safety
- `panic` usage - Controlled execution

## License

MIT License
