// Package fmt registers the Go standard library fmt package.
package fmt

import (
	"fmt"

	"gig/importer"
	"gig/value"
)

func init() {
	pkg := importer.RegisterPackage("fmt", "fmt")

	// Print functions
	pkg.AddFunction("Print", fmt.Print, "", directPrint)
	pkg.AddFunction("Println", fmt.Println, "", directPrintln)
	pkg.AddFunction("Printf", fmt.Printf, "", nil) // variadic, use reflect
	pkg.AddFunction("Sprint", fmt.Sprint, "", directSprint)
	pkg.AddFunction("Sprintln", fmt.Sprintln, "", directSprintln)
	pkg.AddFunction("Sprintf", fmt.Sprintf, "", nil) // variadic, use reflect
	pkg.AddFunction("Fprint", fmt.Fprint, "", nil)
	pkg.AddFunction("Fprintln", fmt.Fprintln, "", nil)
	pkg.AddFunction("Fprintf", fmt.Fprintf, "", nil)

	// Scan functions
	pkg.AddFunction("Scan", fmt.Scan, "", nil)
	pkg.AddFunction("Scanln", fmt.Scanln, "", nil)
	pkg.AddFunction("Scanf", fmt.Scanf, "", nil)
	pkg.AddFunction("Sscan", fmt.Sscan, "", nil)
	pkg.AddFunction("Sscanln", fmt.Sscanln, "", nil)
	pkg.AddFunction("Sscanf", fmt.Sscanf, "", nil)

	// Error functions
	pkg.AddFunction("Errorf", fmt.Errorf, "", nil)

	// Formatter type
	pkg.AddType("Formatter", nil, "interface for custom formatting")
}

// directPrint is a typed wrapper for fmt.Print.
func directPrint(args []value.Value) value.Value {
	ifaceArgs := make([]any, len(args))
	for i, arg := range args {
		ifaceArgs[i] = arg.Interface()
	}
	n, err := fmt.Print(ifaceArgs...)
	return value.FromInterface([]any{n, err})
}

// directPrintln is a typed wrapper for fmt.Println.
func directPrintln(args []value.Value) value.Value {
	ifaceArgs := make([]any, len(args))
	for i, arg := range args {
		ifaceArgs[i] = arg.Interface()
	}
	n, err := fmt.Println(ifaceArgs...)
	return value.FromInterface([]any{n, err})
}

// directSprint is a typed wrapper for fmt.Sprint.
func directSprint(args []value.Value) value.Value {
	ifaceArgs := make([]any, len(args))
	for i, arg := range args {
		ifaceArgs[i] = arg.Interface()
	}
	return value.MakeString(fmt.Sprint(ifaceArgs...))
}

// directSprintln is a typed wrapper for fmt.Sprintln.
func directSprintln(args []value.Value) value.Value {
	ifaceArgs := make([]any, len(args))
	for i, arg := range args {
		ifaceArgs[i] = arg.Interface()
	}
	return value.MakeString(fmt.Sprintln(ifaceArgs...))
}
