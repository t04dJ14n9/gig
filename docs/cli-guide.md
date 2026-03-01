# Gig CLI Guide

The `gig` command-line tool provides three main functionalities: dependency package generation and an interactive REPL (Read-Eval-Print-Loop) for Go code.

## Installation

```bash
# Install the CLI tool
go install github.com/t04dJ14n9/gig/cmd/gig@latest

# Or run directly (Go 1.21+)
go run github.com/t04dJ14n9/gig/cmd/gig@latest --help
```

## Commands

### `gig init` - Initialize Dependency Package

Creates a new dependency package directory with a `pkgs.go` template file.

```bash
gig init -package <name>
```

**Options:**
- `-package <name>` - Name of the dependency package to create (required)

**Example:**
```bash
gig init -package mydep
```

This creates:
```
mydep/
└── pkgs.go    # Template file for adding package imports
```

The generated `pkgs.go`:
```go
package mydep

// Edit this file to add/remove package imports.
// Then run: gig gen ./mydep

import (
	// Standard library examples
	// _ "fmt"
	// _ "strings"
	// _ "time"

	// Third-party library examples
	// _ "github.com/spf13/cast"
)
```

### `gig gen` - Generate Registration Code

Generates registration code for packages imported in the dependency package.

```bash
gig gen <dir>
```

**Arguments:**
- `<dir>` - Path to the dependency package directory

**Example:**
```bash
gig gen ./mydep
```

After editing `mydep/pkgs.go`:
```go
package mydep

import (
	_ "fmt"
	_ "strings"
	_ "github.com/spf13/cast"
)
```

Running `gig gen ./mydep` generates:
```
mydep/
├── pkgs.go
└── packages/
    ├── fmt.go
    ├── strings.go
    └── github_com_spf13_cast.go
```

Each generated file contains registration code for the package's exported symbols (functions, constants, variables, and types).

### `gig repl` - Interactive REPL

Starts an interactive Go interpreter with hot-loading support for external packages.

```bash
gig repl
```

## REPL Commands

The REPL supports special commands prefixed with `:`:

| Command | Aliases | Description |
|---------|---------|-------------|
| `:help` | `:h`, `:?` | Show help message |
| `:quit` | `:q`, `:exit` | Exit the REPL |
| `:clear` | - | Clear session state (variables, imports, declarations) |
| `:imports` | - | List imported packages |
| `:vars` | - | List captured variables |
| `:env` | - | List all variables and packages |
| `:source` | - | Show accumulated declarations |
| `:timeout <dur>` | - | Set execution timeout (e.g., `:timeout 5s`) |
| `:plugins` | - | List loaded plugins (hot-loaded packages) |

### `:env` - Environment Overview

The `:env` command provides a comprehensive view of the current session state:

```
>>> import "fmt"
>>> x := 5
>>> y := x + 3
>>> :env
Packages:
  fmt (fmt)

Variables:
  x int = 5
  y int = 8
```

## REPL Features

### Input Types

The REPL automatically classifies and handles different input types:

| Type | Example | Description |
|------|---------|-------------|
| Import | `import "fmt"` | Import packages |
| Expression | `1+1`, `x + y` | Evaluates and prints result |
| Statement | `x := 1`, `for i:=0; i<5; i++ {}` | Executes statement |
| Declaration | `func foo() {}`, `type Point struct{}` | Stores for later use |

### Multiline Input

The REPL supports multiline input for complex code:

```
>>> func fibonacci(n int) int {
...     if n <= 1 {
...         return n
...     }
...     return fibonacci(n-1) + fibonacci(n-2)
... }
>>> fibonacci(10)
55
```

Multiline mode activates automatically when input has unclosed brackets (`{}`, `()`, `[]`) or backticks (`` ` ``).

### Tab Completion

The REPL provides intelligent tab completion for:

- **Commands**: Type `:` then `Tab` to cycle through commands
- **Variables**: Complete variable names from the current session
- **Packages**: Complete package names from imports and known packages
- **Package Symbols**: After typing `pkg.`, complete exported symbols (for hot-loaded packages)

**Known packages include**: `fmt`, `strings`, `strconv`, `math`, `time`, `bytes`, `errors`, `sort`, `regexp`, `json`, `xml`, `http`, `url`, and 40+ more.

### Variable Persistence

Variables defined with short declaration (`:=`) persist across statements:

```
>>> x := 10
>>> y := 20
>>> x + y
30
>>> :vars
Variables:
  x int = 10
  y int = 20
```

### Hot-Loading External Packages

The REPL can hot-load external Go packages using Go's plugin system (Linux and macOS only):

```
>>> import "github.com/spf13/cast"
Imported: github.com/spf13/cast (hot-loaded)
>>> cast.ToString(123)
"123"
>>> cast.ToInt("456")
456
```

**Platform Support:**
- ✅ Linux - Full support
- ✅ macOS - Full support
- ❌ Windows - Not supported (returns error)

**How it works:**
1. Downloads the package with `go get`
2. Generates wrapper code using `go doc` for symbol discovery
3. Compiles as a `.so` plugin with `go build -buildmode=plugin`
4. Loads the plugin and registers symbols with gig

**Plugin Cache:**
Plugins are cached in `~/.gig/plugins/`:
```
~/.gig/plugins/
├── github.com/
│   └── spf13/
│       └── cast/
│           ├── cast.so      # Compiled plugin
│           └── wrapper.go   # Generated wrapper code
├── go.mod                   # Plugin module file
└── plugin_registry.json     # Plugin metadata
```

## Examples

### Basic REPL Session

```
>>> import "fmt"
Imported: fmt
>>> import "strings"
Imported: strings
>>> name := "world"
>>> fmt.Sprintf("Hello, %s!", strings.ToUpper(name))
"Hello, WORLD!"
>>> :env
Packages:
  fmt (fmt)
  strings (strings)

Variables:
  name string = "world"
```

### Working with Functions

```
>>> func add(a, b int) int { return a + b }

>>> func multiply(a, b int) int { return a * b }

>>> add(2, 3)
5
>>> multiply(4, 5)
20
>>> :source
Accumulated declarations:
func add(a, b int) int { return a + b }

func multiply(a, b int) int { return a * b }
```

### Custom Types

```
>>> type Point struct { X, Y int }

>>> p := Point{X: 10, Y: 20}
>>> p.X + p.Y
30
>>> :vars
Variables:
  p main.Point = {X:10 Y:20}
```

### Using Third-Party Packages

```
>>> import "github.com/spf13/cast"
Imported: github.com/spf13/cast (hot-loaded)
>>> cast.ToInt("123")
123
>>> cast.ToBool("true")
true
>>> cast.ToString(456.78)
"456.78"
>>> :plugins
Loaded plugins:
  github.com/spf13/cast
```

### Control Flow

```
>>> sum := 0
>>> for i := 1; i <= 10; i++ { sum += i }

>>> sum
55
>>> :vars
Variables:
  sum int = 55
```

### Using Context and Timeout

```
>>> :timeout 5s
Timeout set to: 5s
>>> :timeout
Current timeout: 5s
```

The timeout applies to expression evaluation and statement execution, preventing infinite loops from hanging the REPL.

## Troubleshooting

### Plugin Loading Errors

**"plugin was built with a different version of package"**

This occurs when the plugin uses a different version of gig than the host. The plugin manager automatically handles this by using a `replace` directive in the plugin's `go.mod`.

**"plugin loading requires Linux or macOS"**

Windows does not support Go's plugin system. Consider:
- Using WSL (Windows Subsystem for Linux)
- Pre-generating registration code with `gig gen`
- Using the built-in stdlib packages

### Import Errors

**"package not found; run 'go get <package>' first"**

The package isn't installed. Run:
```bash
go get <package-path>
```

Then retry the import in the REPL.

## Configuration

### Timeout

Default execution timeout is 10 seconds. Change with:
```
:timeout 30s
```

### Known Packages

The REPL has a built-in mapping of common package names to import paths:

| Name | Import Path |
|------|-------------|
| `fmt` | `fmt` |
| `strings` | `strings` |
| `json` | `encoding/json` |
| `http` | `net/http` |
| `url` | `net/url` |
| `rand` | `math/rand` |
| ... | ... |

This allows automatic import detection when you use package functions without explicit imports.

## See Also

- [README.md](../README.md) - Main documentation
- [gig-internals.md](gig-internals.md) - Internal architecture
- [examples/](../examples/) - Example programs
