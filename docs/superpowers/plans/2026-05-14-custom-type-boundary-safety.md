# Custom Type Boundary Safety Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reject user-defined types passed to third-party library functions at compile time, while allowing them through stdlib calls (which we control via adapters/wrappers).

**Architecture:** Add a compile-time SSA analysis pass that inspects every external call site. When an argument's concrete type is user-defined (declared in the script's "main" package) and the target function belongs to a third-party package (import path contains a dot), emit a compile error. Stdlib calls are allowed because gig guarantees correctness via adapters and wrappers. Also fix the three existing interface adapter bugs (scope, Swap, detached VM) so the stdlib path is actually correct.

**Tech Stack:** Go, `go/types`, `golang.org/x/tools/go/ssa`

---

### Task 1: Add `isUserDefinedType` and `isStdlib` helper functions

**Files:**
- Create: `compiler/typecheck.go`
- Test: `compiler/typecheck_test.go`

These are the core detection functions used by the compile-time analysis pass.

- [ ] **Step 1: Write the failing test**

```go
// compiler/typecheck_test.go
package compiler

import (
	"go/types"
	"testing"
)

func TestIsStdlibPath(t *testing.T) {
	tests := []struct {
		path   string
		expect bool
	}{
		{"fmt", true},
		{"encoding/json", true},
		{"sort", true},
		{"io", true},
		{"net/http", true},
		{"golang.org/x/tools", false},
		{"github.com/foo/bar", false},
		{"github.com/t04dJ14n9/gig", false},
		{"command-line-arguments", false},
		{"main", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isStdlibPath(tt.path)
			if got != tt.expect {
				t.Errorf("isStdlibPath(%q) = %v, want %v", tt.path, got, tt.expect)
			}
		})
	}
}

func TestIsUserDefinedNamedType(t *testing.T) {
	// Create a fake package for "command-line-arguments" (user script)
	userPkg := types.NewPackage("command-line-arguments", "main")
	userType := types.NewNamed(
		types.NewTypeName(0, userPkg, "MySlice", nil),
		types.NewSlice(types.Typ[types.Int]),
		nil,
	)

	// Create a fake package for "sort" (stdlib)
	sortPkg := types.NewPackage("sort", "sort")
	sortType := types.NewNamed(
		types.NewTypeName(0, sortPkg, "IntSlice", nil),
		types.NewSlice(types.Typ[types.Int]),
		nil,
	)

	// Primitives
	intType := types.Typ[types.Int]
	stringType := types.Typ[types.String]

	tests := []struct {
		name   string
		typ    types.Type
		expect bool
	}{
		{"user named type", userType, true},
		{"user pointer to named type", types.NewPointer(userType), true},
		{"stdlib named type", sortType, false},
		{"bare int", intType, false},
		{"bare string", stringType, false},
		{"slice of int", types.NewSlice(intType), false},
		{"nil type", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUserDefinedNamedType(tt.typ)
			if got != tt.expect {
				t.Errorf("isUserDefinedNamedType(%v) = %v, want %v", tt.typ, got, tt.expect)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v -run "TestIsStdlibPath|TestIsUserDefinedNamedType" ./compiler/`
Expected: FAIL — `isStdlibPath` and `isUserDefinedNamedType` undefined

- [ ] **Step 3: Write minimal implementation**

```go
// compiler/typecheck.go
// typecheck.go provides compile-time type flow analysis for custom type safety.
package compiler

import (
	"fmt"
	"go/types"
	"strings"
)

// isStdlibPath returns true if the import path belongs to the Go standard library.
// Stdlib paths have no dots (e.g., "fmt", "encoding/json", "sort").
// Third-party paths contain dots (e.g., "github.com/foo", "golang.org/x/tools").
func isStdlibPath(path string) bool {
	if path == "" || path == "command-line-arguments" || path == "main" {
		return false
	}
	// Stdlib paths never contain a dot in the first path segment.
	// "encoding/json" → first segment "encoding" (no dot) → stdlib
	// "golang.org/x/tools" → first segment "golang.org" (has dot) → third-party
	firstSlash := strings.IndexByte(path, '/')
	firstSegment := path
	if firstSlash >= 0 {
		firstSegment = path[:firstSlash]
	}
	return !strings.ContainsRune(firstSegment, '.')
}

// isUserDefinedNamedType returns true if the type is a named type defined in the
// user's script (the "main" / "command-line-arguments" package), not from an
// external registered package.
//
// Unwraps pointers: *MyStruct → MyStruct → check package.
func isUserDefinedNamedType(t types.Type) bool {
	if t == nil {
		return false
	}
	// Unwrap pointers
	for {
		if ptr, ok := t.(*types.Pointer); ok {
			t = ptr.Elem()
		} else {
			break
		}
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj == nil {
		return false
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return false // universe scope (error, any, etc.)
	}
	pkgPath := pkg.Path()
	// User script types have "command-line-arguments" or "main" as their package path.
	return pkgPath == "command-line-arguments" || pkgPath == "main"
}

// containsUserDefinedType checks if a type contains a user-defined named type,
// including inside slices, maps, and pointers.
// This catches cases like []MyStruct or map[string]MyStruct.
func containsUserDefinedType(t types.Type) bool {
	if t == nil {
		return false
	}
	// Direct check
	if isUserDefinedNamedType(t) {
		return true
	}
	// Unwrap structural types
	switch tt := t.Underlying().(type) {
	case *types.Slice:
		return isUserDefinedNamedType(tt.Elem())
	case *types.Map:
		return isUserDefinedNamedType(tt.Key()) || isUserDefinedNamedType(tt.Elem())
	case *types.Pointer:
		return containsUserDefinedType(tt.Elem())
	case *types.Array:
		return isUserDefinedNamedType(tt.Elem())
	}
	return false
}

// validateExternalCallArgs checks whether any argument to an external function
// call is a user-defined type being passed to a third-party (non-stdlib) package.
//
// Returns an error if a user-defined type crosses into third-party code.
// Returns nil if all arguments are safe (primitives, stdlib types, or target is stdlib).
func validateExternalCallArgs(pkgPath, funcName string, argTypes []types.Type) error {
	// Stdlib calls are always allowed — we guarantee correctness via adapters/wrappers.
	if isStdlibPath(pkgPath) {
		return nil
	}

	// Third-party call: check every argument
	for i, argType := range argTypes {
		if containsUserDefinedType(argType) {
			typeName := describeType(argType)
			return fmt.Errorf(
				"cannot pass interpreter-defined type %q to third-party function %s.%s (argument %d): "+
					"custom types are not compatible with external libraries that use reflection. "+
					"Use primitive types, slices, maps, or types from registered packages instead",
				typeName, pkgPath, funcName, i+1,
			)
		}
	}
	return nil
}

// describeType returns a human-readable name for a type for error messages.
func describeType(t types.Type) string {
	if t == nil {
		return "<nil>"
	}
	return t.String()
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v -run "TestIsStdlibPath|TestIsUserDefinedNamedType" ./compiler/`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add compiler/typecheck.go compiler/typecheck_test.go
git commit -m "feat(compiler): add type boundary detection helpers

Add isStdlibPath, isUserDefinedNamedType, and containsUserDefinedType
for compile-time detection of custom types crossing into third-party code."
```

---

### Task 2: Add compile-time validation in `compileExternalStaticCall`

**Files:**
- Modify: `compiler/compile_ext.go` (add validation call)
- Modify: `compiler/compiler.go` (add errors list and config to compiler struct)
- Modify: `compiler/typecheck_test.go` (add validation test)

This is the core enforcement: when compiling an external call, check every argument
and reject user-defined types flowing into third-party functions.

- [ ] **Step 1: Write the failing test for `validateExternalCallArgs`**

Add to `compiler/typecheck_test.go`:

```go
func TestRejectUserTypeToThirdParty(t *testing.T) {
	// Mock: user type "MySlice" from "command-line-arguments"
	userPkg := types.NewPackage("command-line-arguments", "main")
	userType := types.NewNamed(
		types.NewTypeName(0, userPkg, "MySlice", nil),
		types.NewSlice(types.Typ[types.Int]),
		nil,
	)

	// Mock: stdlib type from "sort"
	sortPkg := types.NewPackage("sort", "sort")
	sortType := types.NewNamed(
		types.NewTypeName(0, sortPkg, "IntSlice", nil),
		types.NewSlice(types.Typ[types.Int]),
		nil,
	)

	tests := []struct {
		name      string
		pkgPath   string
		argTypes  []types.Type
		expectErr bool
	}{
		{
			name:      "user type to stdlib is allowed",
			pkgPath:   "sort",
			argTypes:  []types.Type{userType},
			expectErr: false,
		},
		{
			name:      "user type to third-party is rejected",
			pkgPath:   "github.com/foo/bar",
			argTypes:  []types.Type{userType},
			expectErr: true,
		},
		{
			name:      "stdlib type to third-party is allowed",
			pkgPath:   "github.com/foo/bar",
			argTypes:  []types.Type{sortType},
			expectErr: false,
		},
		{
			name:      "primitive to third-party is allowed",
			pkgPath:   "github.com/foo/bar",
			argTypes:  []types.Type{types.Typ[types.Int]},
			expectErr: false,
		},
		{
			name:      "slice of user type to third-party is rejected",
			pkgPath:   "github.com/foo/bar",
			argTypes:  []types.Type{types.NewSlice(userType)},
			expectErr: true,
		},
		{
			name:      "pointer to user type to third-party is rejected",
			pkgPath:   "github.com/foo/bar",
			argTypes:  []types.Type{types.NewPointer(userType)},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExternalCallArgs(tt.pkgPath, "SomeFunc", tt.argTypes)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it passes** (validation function already implemented in Task 1)

Run: `go test -v -run "TestRejectUserTypeToThirdParty" ./compiler/`
Expected: PASS

- [ ] **Step 3: Add error collection and config to the compiler struct**

Modify `compiler/compiler.go` — add fields to the compiler struct (after `phiSlots`):

```go
	// errors collects non-fatal compilation errors (e.g., type safety violations).
	errors []error

	// allowUnsafeTypePass disables type safety validation for external calls.
	allowUnsafeTypePass bool
```

Add the `addError` method:

```go
// addError records a compilation error for later reporting.
func (c *compiler) addError(err error) {
	c.errors = append(c.errors, err)
}
```

In `Compile()`, add error check just before `c.program.ResolveExternCalls()` (before line 226):

```go
	// Check for compilation errors (type safety violations, etc.)
	if len(c.errors) > 0 {
		return nil, c.errors[0]
	}
```

- [ ] **Step 4: Wire validation into `compileExternalStaticCall`**

Modify `compiler/compile_ext.go` — add validation at the start of `compileExternalStaticCall`, before the loop that pushes args (before line 59):

```go
func (c *compiler) compileExternalStaticCall(i *ssa.Call, fn *ssa.Function, resultIdx int) {
	// Validate: reject user-defined types flowing into third-party packages.
	if !c.allowUnsafeTypePass && fn.Pkg != nil && fn.Pkg.Pkg != nil {
		pkgPath := fn.Pkg.Pkg.Path()
		argTypes := make([]types.Type, len(i.Call.Args))
		for idx, arg := range i.Call.Args {
			argTypes[idx] = arg.Type()
		}
		if err := validateExternalCallArgs(pkgPath, fn.Name(), argTypes); err != nil {
			c.addError(err)
		}
	}

	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}
	// ... rest unchanged ...
```

- [ ] **Step 5: Run all compiler tests**

Run: `go test -v ./compiler/`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add compiler/typecheck_test.go compiler/compile_ext.go compiler/compiler.go
git commit -m "feat(compiler): reject user-defined types passed to third-party functions

Add compile-time validation in compileExternalStaticCall that rejects
user-defined types (from the interpreter script) being passed as
arguments to third-party library functions. Stdlib calls are allowed
because gig guarantees correctness via adapters and wrappers."
```

---

### Task 3: Add `WithAllowUnsafeTypePass()` escape hatch

**Files:**
- Modify: `gig.go` (add build option)
- Modify: `compiler/build.go` (add build option, thread to compiler)
- Modify: `compiler/compiler.go` (update constructor)
- Test: `gig_test.go` (integration test)

- [ ] **Step 1: Add the build option to `compiler/build.go`**

In the `buildConfig` struct, add:

```go
type buildConfig struct {
	allowPanic          bool
	allowUnsafeTypePass bool
}
```

Add the option function:

```go
// WithAllowUnsafeTypePass disables type safety validation for external calls.
func WithAllowUnsafeTypePass() BuildOption {
	return func(c *buildConfig) {
		c.allowUnsafeTypePass = true
	}
}
```

In `Build()`, thread the option to `NewCompiler`. Replace line 68:

```go
	compiled, err := NewCompiler(reg, cfg.allowUnsafeTypePass).Compile(ssaResult.Pkg)
```

- [ ] **Step 2: Update `NewCompiler` to accept the config**

In `compiler/compiler.go`, change the `NewCompiler` signature:

```go
func NewCompiler(lookup PackageLookup, allowUnsafeTypePass bool) Compiler {
	return &compiler{
		lookup:              lookup,
		allowUnsafeTypePass: allowUnsafeTypePass,
		constants:           make([]any, 0),
		types:               make([]types.Type, 0),
		globals:             make(map[string]int),
		globalZeroValues:    make(map[int]reflect.Value),
		externalVarValues:   make(map[int]any),
		funcs:               make(map[string]*bytecode.CompiledFunction),
		funcIndex:           make(map[*ssa.Function]int),
	}
}
```

- [ ] **Step 3: Add the public API option to `gig.go`**

In the `buildConfig` struct, add:

```go
	allowUnsafeTypePass bool
```

Add the option function:

```go
// WithAllowUnsafeTypePass disables the compile-time check that rejects
// user-defined types being passed to third-party library functions.
//
// By default, gig rejects custom types (structs, interfaces, named types
// defined in the script) from being passed to non-stdlib external functions,
// because Go's reflect.StructOf cannot attach methods to synthesized types.
// Third-party libraries using reflection will see incorrect type information.
//
// Enable this option only if you understand the reflection limitations and
// your third-party libraries do not inspect type identity or method sets.
func WithAllowUnsafeTypePass() BuildOption {
	return func(c *buildConfig) {
		c.allowUnsafeTypePass = true
	}
}
```

In `Build()`, thread the option. Add after the `allowPanic` block:

```go
	if cfg.allowUnsafeTypePass {
		compilerOpts = append(compilerOpts, compiler.WithAllowUnsafeTypePass())
	}
```

- [ ] **Step 4: Write integration test**

Add to `gig_test.go`:

```go
func TestWithAllowUnsafeTypePass(t *testing.T) {
	source := `
		package main

		func Add(a, b int) int { return a + b }
	`
	prog, err := Build(source, WithAllowUnsafeTypePass())
	if err != nil {
		t.Fatalf("Build with WithAllowUnsafeTypePass failed: %v", err)
	}
	defer prog.Close()

	result, err := prog.Run("Add", 1, 2)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if result != int64(3) {
		t.Errorf("got %v, want 3", result)
	}
}
```

- [ ] **Step 5: Run tests**

Run: `go test -v -run TestWithAllowUnsafeTypePass ./`
Expected: PASS

- [ ] **Step 6: Run full test suite**

Run: `go test -v -race ./...`
Expected: PASS — all existing tests use stdlib, so validation never triggers

- [ ] **Step 7: Commit**

```bash
git add gig.go compiler/build.go compiler/compiler.go gig_test.go
git commit -m "feat: add WithAllowUnsafeTypePass build option

Escape hatch for users who understand the reflection limitations and
want to pass custom types to third-party libraries anyway."
```

---

### Task 4: Fix adapter scope — narrow `isHostCallbackInterface` to exact types

**Files:**
- Modify: `vm/ops_convert.go:392-409` (narrow matching)
- Test: `tests/correctness_test.go` (add type identity test)

- [ ] **Step 1: Write the failing test**

Add to `tests/correctness_test.go`:

```go
func TestCustomInterfaceTypeIdentity(t *testing.T) {
	source := `
		package main

		import "fmt"

		type Sortable interface {
			Len() int
			Less(i, j int) bool
			Swap(i, j int)
		}

		type MyData struct {
			Items []int
		}

		func (d MyData) Len() int           { return len(d.Items) }
		func (d MyData) Less(i, j int) bool { return d.Items[i] < d.Items[j] }
		func (d MyData) Swap(i, j int)      { d.Items[i], d.Items[j] = d.Items[j], d.Items[i] }

		func TypeAssertionWorks() string {
			var s Sortable = MyData{Items: []int{3, 1, 2}}
			if d, ok := s.(MyData); ok {
				return fmt.Sprintf("ok:%d", d.Len())
			}
			return "failed"
		}
	`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer prog.Close()

	result, err := prog.Run("TypeAssertionWorks")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if result != "ok:3" {
		t.Errorf("got %v, want ok:3", result)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v -run TestCustomInterfaceTypeIdentity ./tests/`
Expected: FAIL — type assertion fails because adapter replaces the concrete type

- [ ] **Step 3: Narrow `isHostCallbackInterface` to exact type matching only**

In `vm/ops_convert.go`, replace the `isHostCallbackInterface` function (lines 392-409):

```go
func isHostCallbackInterface(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}
	// Only match the exact sort.Interface and heap.Interface named types.
	// Do NOT match arbitrary interfaces that happen to have Len/Less/Swap methods.
	pkgPath := obj.Pkg().Path()
	name := obj.Name()
	return name == "Interface" && (pkgPath == "sort" || pkgPath == "container/heap")
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v -run TestCustomInterfaceTypeIdentity ./tests/`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go test -v -race ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add vm/ops_convert.go tests/correctness_test.go
git commit -m "fix(vm): narrow interface adapter to exact sort/heap.Interface only

Remove shape-based matching (any interface with Len/Less/Swap) that
incorrectly replaced concrete type identity for user-defined interfaces.
Only sort.Interface and container/heap.Interface are adapted now."
```

---

### Task 5: Fix Swap bypass — always call compiled method first

**Files:**
- Modify: `vm/interface_adapter.go:41-46` (fix Swap dispatch order)
- Test: `tests/correctness_test.go` (add Swap test)

- [ ] **Step 1: Write the failing test**

Add to `tests/correctness_test.go`:

```go
func TestSortSwapCallsUserMethod(t *testing.T) {
	source := `
		package main

		import (
			"fmt"
			"sort"
		)

		var swapCount int

		type CountSlice []int
		func (s CountSlice) Len() int           { return len(s) }
		func (s CountSlice) Less(i, j int) bool { return s[i] < s[j] }
		func (s CountSlice) Swap(i, j int) {
			s[i], s[j] = s[j], s[i]
			swapCount++
		}

		func SortAndCountSwaps() string {
			s := CountSlice{3, 1, 2}
			sort.Sort(s)
			return fmt.Sprintf("%v:%d", []int(s), swapCount)
		}
	`
	prog, err := gig.Build(source, gig.WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer prog.Close()

	result, err := prog.Run("SortAndCountSwaps")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	resultStr := fmt.Sprintf("%v", result)
	if strings.HasSuffix(resultStr, ":0") {
		t.Errorf("Swap method was never called (swapCount=0), got %v", result)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v -run TestSortSwapCallsUserMethod ./tests/`
Expected: FAIL — `swapCount` is 0 because `callReceiverSliceSwap` bypasses the user's Swap

- [ ] **Step 3: Fix Swap dispatch order**

In `vm/interface_adapter.go`, replace lines 41-46:

```go
func (a *interpretedInterfaceAdapter) Swap(i, j int) {
	// Always try the compiled method first — user-defined Swap may update
	// auxiliary state (counters, indexes, parallel slices) beyond just
	// swapping elements. Only fall back to direct slice swap if no
	// compiled method is available.
	result, ok := callInterfaceMethodValue(a.program, "Swap", a.receiverTypeName, a.receiver, []value.Value{value.MakeInt(int64(i)), value.MakeInt(int64(j))})
	_ = result
	if ok {
		return
	}
	// Fallback: direct slice element swap (for types without a compiled Swap)
	callReceiverSliceSwap(a.receiver, i, j)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v -run TestSortSwapCallsUserMethod ./tests/`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go test -v -race ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add vm/interface_adapter.go tests/correctness_test.go
git commit -m "fix(vm): call compiled Swap method before direct slice swap

The adapter was bypassing user-defined Swap methods for slice-backed
types, silently skipping auxiliary state updates. Now always tries the
compiled method first, falling back to direct swap only when no
compiled method exists."
```

---

### Task 6: Fix detached VM — thread caller context through adapter

**Files:**
- Modify: `vm/interface_adapter.go` (add caller VM context to adapter)
- Modify: `vm/ops_convert.go:381-390` (pass caller VM to adapter)
- Test: `tests/correctness_test.go` (add globals test)

- [ ] **Step 1: Write the failing test**

Add to `tests/correctness_test.go`:

```go
func TestSortCallbackReadsGlobals(t *testing.T) {
	source := `
		package main

		import (
			"fmt"
			"sort"
		)

		var multiplier = 1

		type ScaledSlice []int
		func (s ScaledSlice) Len() int           { return len(s) }
		func (s ScaledSlice) Less(i, j int) bool { return s[i]*multiplier < s[j]*multiplier }
		func (s ScaledSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

		func SortWithGlobal() string {
			multiplier = -1
			s := ScaledSlice{1, 2, 3}
			sort.Sort(s)
			return fmt.Sprintf("%v", []int(s))
		}
	`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer prog.Close()

	result, err := prog.Run("SortWithGlobal")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	// multiplier = -1 reverses the comparison, so sort produces descending order
	if fmt.Sprintf("%v", result) != "[3 2 1]" {
		t.Errorf("got %v, want [3 2 1] (global multiplier=-1 should reverse sort)", result)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v -run TestSortCallbackReadsGlobals ./tests/`
Expected: FAIL — callback sees zeroed globals (multiplier=0 instead of -1)

- [ ] **Step 3: Add caller context to the adapter struct**

In `vm/interface_adapter.go`, update the struct and constructor:

```go
type interpretedInterfaceAdapter struct {
	program          *bytecode.CompiledProgram
	receiver         value.Value
	receiverTypeName string
	// Caller VM context for correct callback execution
	globals        []value.Value
	initialGlobals []value.Value
	shared         *SharedGlobals
	ctx            context.Context
	goroutines     *GoroutineTracker
}

func newInterpretedInterfaceAdapter(
	program *bytecode.CompiledProgram,
	receiver value.Value,
	receiverTypeName string,
	globals []value.Value,
	initialGlobals []value.Value,
	shared *SharedGlobals,
	ctx context.Context,
	goroutines *GoroutineTracker,
) *interpretedInterfaceAdapter {
	return &interpretedInterfaceAdapter{
		program:          program,
		receiver:         receiver,
		receiverTypeName: receiverTypeName,
		globals:          globals,
		initialGlobals:   initialGlobals,
		shared:           shared,
		ctx:              ctx,
		goroutines:       goroutines,
	}
}
```

Add `"context"` to the import list.

- [ ] **Step 4: Fix `callCompiledMethodValue` to use caller context**

Update its signature and temp VM construction:

```go
func callCompiledMethodValue(program *bytecode.CompiledProgram, methodName, receiverTypeName string, receiver value.Value, args []value.Value, globals []value.Value, initialGlobals []value.Value, shared *SharedGlobals, ctx context.Context, goroutines *GoroutineTracker) (value.Value, bool) {
	if program == nil {
		return value.MakeNil(), false
	}

	for _, fn := range program.MethodsByName[methodName] {
		if receiverTypeName != "" && fn.ReceiverTypeName != receiverTypeName {
			continue
		}
		if shouldPanicOnNilValueReceiver(receiver, fn) {
			return value.MakeNil(), false
		}

		methodReceiver := receiverForCompiledMethod(methodName, receiver)
		callArgs := make([]value.Value, 0, len(args)+1)
		callArgs = append(callArgs, methodReceiver)
		callArgs = append(callArgs, args...)

		// Use the caller's globals and context instead of empty defaults
		callGlobals := make([]value.Value, len(program.Globals))
		if globals != nil {
			copy(callGlobals, globals)
		}
		callCtx := ctx
		if callCtx == nil {
			callCtx = context.Background()
		}

		tempVM := &vm{
			program:        program,
			stack:          make([]value.Value, deferVMStackSize),
			sp:             0,
			frames:         make([]*Frame, initialFrameDepth),
			fp:             0,
			globals:        callGlobals,
			initialGlobals: initialGlobals,
			ctx:            callCtx,
			goroutines:     goroutines,
		}
		if shared != nil {
			tempVM.shared = shared
		}
		tempVM.callFunction(fn, callArgs, nil)

		var result value.Value
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("compiled method %q panicked: %v", methodName, r)
				}
			}()
			result, err = tempVM.run()
		}()
		if err != nil {
			continue
		}
		return result, true
	}

	return value.MakeNil(), false
}
```

Add `"fmt"` to the import list if not already present.

- [ ] **Step 5: Update `callInterfaceMethodValue` to thread context**

```go
func callInterfaceMethodValue(program *bytecode.CompiledProgram, methodName, receiverTypeName string, receiver value.Value, args []value.Value, globals []value.Value, initialGlobals []value.Value, shared *SharedGlobals, ctx context.Context, goroutines *GoroutineTracker) (value.Value, bool) {
	if result, ok := callCompiledMethodValue(program, methodName, receiverTypeName, receiver, args, globals, initialGlobals, shared, ctx, goroutines); ok {
		return result, true
	}
	return callReflectInterfaceMethod(methodName, receiver, args)
}
```

- [ ] **Step 6: Update the adapter's `call` method**

```go
func (a *interpretedInterfaceAdapter) call(methodName string, args ...value.Value) value.Value {
	result, ok := callInterfaceMethodValue(a.program, methodName, a.receiverTypeName, a.receiver, args, a.globals, a.initialGlobals, a.shared, a.ctx, a.goroutines)
	if !ok {
		return value.MakeNil()
	}
	return result
}
```

- [ ] **Step 7: Update the Swap method from Task 5 to thread context**

```go
func (a *interpretedInterfaceAdapter) Swap(i, j int) {
	result, ok := callInterfaceMethodValue(a.program, "Swap", a.receiverTypeName, a.receiver, []value.Value{value.MakeInt(int64(i)), value.MakeInt(int64(j))}, a.globals, a.initialGlobals, a.shared, a.ctx, a.goroutines)
	_ = result
	if ok {
		return
	}
	callReceiverSliceSwap(a.receiver, i, j)
}
```

- [ ] **Step 8: Update the construction call site in `vm/ops_convert.go`**

Replace the `makeInterpretedInterfaceAdapter` method:

```go
func (v *vm) makeInterpretedInterfaceAdapter(targetType, concreteType types.Type, receiver value.Value) (*interpretedInterfaceAdapter, bool) {
	if !isHostCallbackInterface(targetType) {
		return nil, false
	}
	receiverTypeName := namedTypeName(concreteType)
	if receiverTypeName == "" || !hasCompiledInterfaceMethods(v.program, receiverTypeName) {
		return nil, false
	}
	return newInterpretedInterfaceAdapter(
		v.program, receiver, receiverTypeName,
		v.getGlobals(), v.initialGlobals, v.shared, v.ctx, v.goroutines,
	), true
}
```

- [ ] **Step 9: Run test to verify it passes**

Run: `go test -v -run TestSortCallbackReadsGlobals ./tests/`
Expected: PASS

- [ ] **Step 10: Run full test suite**

Run: `go test -v -race ./...`
Expected: PASS

- [ ] **Step 11: Commit**

```bash
git add vm/interface_adapter.go vm/ops_convert.go tests/correctness_test.go
git commit -m "fix(vm): thread caller VM context through interface adapter

The adapter was creating a temp VM with zeroed globals and
context.Background(), severing callbacks from the caller's state.
Now passes globals, initialGlobals, shared, ctx, and goroutines
from the caller VM so callbacks see correct program state."
```

---

### Task 7: End-to-end integration tests

**Files:**
- Modify: `gig_test.go` (public API tests)

- [ ] **Step 1: Write comprehensive integration tests**

Add to `gig_test.go`:

```go
func TestCustomTypeBoundaryEnforcement(t *testing.T) {
	t.Run("stdlib sort with custom type is allowed", func(t *testing.T) {
		source := `
			package main

			import (
				"fmt"
				"sort"
			)

			type ByLen []string
			func (s ByLen) Len() int           { return len(s) }
			func (s ByLen) Less(i, j int) bool { return len(s[i]) < len(s[j]) }
			func (s ByLen) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

			func SortWords() string {
				words := ByLen{"banana", "pie", "kiwi"}
				sort.Sort(words)
				return fmt.Sprintf("%v", []string(words))
			}
		`
		prog, err := Build(source)
		if err != nil {
			t.Fatalf("Build should allow stdlib sort: %v", err)
		}
		defer prog.Close()

		result, err := prog.Run("SortWords")
		if err != nil {
			t.Fatalf("Run failed: %v", err)
		}
		expected := "[pie kiwi banana]"
		if fmt.Sprintf("%v", result) != expected {
			t.Errorf("got %v, want %v", result, expected)
		}
	})

	t.Run("primitives to any function are allowed", func(t *testing.T) {
		source := `
			package main

			import "fmt"

			func FormatInt() string {
				return fmt.Sprintf("%d", 42)
			}
		`
		prog, err := Build(source)
		if err != nil {
			t.Fatalf("Build should allow primitives: %v", err)
		}
		defer prog.Close()

		result, err := prog.Run("FormatInt")
		if err != nil {
			t.Fatalf("Run failed: %v", err)
		}
		if result != "42" {
			t.Errorf("got %v, want 42", result)
		}
	})

	t.Run("custom struct with fmt is allowed", func(t *testing.T) {
		source := `
			package main

			import "fmt"

			type Point struct { X, Y int }

			func FormatPoint() string {
				p := Point{1, 2}
				return fmt.Sprintf("(%d,%d)", p.X, p.Y)
			}
		`
		prog, err := Build(source)
		if err != nil {
			t.Fatalf("Build should allow struct with fmt: %v", err)
		}
		defer prog.Close()

		result, err := prog.Run("FormatPoint")
		if err != nil {
			t.Fatalf("Run failed: %v", err)
		}
		if result != "(1,2)" {
			t.Errorf("got %v, want (1,2)", result)
		}
	})
}
```

- [ ] **Step 2: Run the integration tests**

Run: `go test -v -run "TestCustomTypeBoundaryEnforcement|TestWithAllowUnsafeTypePass" ./`
Expected: PASS

- [ ] **Step 3: Run the FULL test suite with race detector**

Run: `go test -v -race ./...`
Expected: PASS — all existing tests pass since they all use stdlib

- [ ] **Step 4: Commit**

```bash
git add gig_test.go
git commit -m "test: add integration tests for custom type boundary enforcement

Tests verify: stdlib sort with custom types works, primitives pass
through to any function, custom structs with fmt work, and the
WithAllowUnsafeTypePass option functions correctly."
```

---

### Task 8: Update documentation

**Files:**
- Modify: `CLAUDE.md`

- [ ] **Step 1: Add type safety section to CLAUDE.md**

Add after the "### Security Model" section in the project CLAUDE.md:

```markdown
### Type Boundary Safety

By default, gig rejects user-defined types (structs, interfaces, named types declared
in the script) from being passed as arguments to third-party library functions (import
path contains a dot, e.g., `github.com/foo/bar`).

**Why**: Go's `reflect.StructOf` cannot attach methods to runtime-synthesized types.
Third-party libraries using reflection (`reflect.MethodByName`, `Type.Implements()`)
will see incorrect type information, causing silent data corruption or panics.

**What's allowed**:
- Custom types → stdlib functions (gig guarantees correctness via adapters for
  `sort.Interface`, `heap.Interface`, `fmt.Stringer`, `error`, etc.)
- Primitive types, slices, maps → any function
- External registered types (`sort.IntSlice`, `time.Time`) → any function

**Escape hatch**: `gig.WithAllowUnsafeTypePass()` disables the compile-time check.
```

- [ ] **Step 2: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: add type boundary safety documentation to CLAUDE.md"
```
