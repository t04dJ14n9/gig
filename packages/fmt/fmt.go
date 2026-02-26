// Package fmt registers the Go standard library fmt package.
package fmt

import (
	"fmt"
	"reflect"

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
	pkg.AddFunction("Sprintf", fmt.Sprintf, "", directSprintf)
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
	pkg.AddFunction("Errorf", fmt.Errorf, "", directErrorf)

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

// directSprintf is a typed wrapper for fmt.Sprintf.
func directSprintf(args []value.Value) value.Value {
	if len(args) < 1 {
		return value.MakeString("")
	}
	format := args[0].String()
	if len(args) == 1 {
		return value.MakeString(fmt.Sprintf(format))
	}

	// SSA packs variadic args into a slice as the second argument
	// Check if args[1] is a slice that needs unpacking
	var ifaceArgs []any
	
	if len(args) == 2 {
		// Check if it's a packed variadic slice
		if args[1].Kind() == value.KindReflect {
			if rv, ok := args[1].ReflectValue(); ok && rv.Kind() == reflect.Slice {
				// Unpack the slice
				ifaceArgs = make([]any, rv.Len())
				for i := 0; i < rv.Len(); i++ {
					elem := rv.Index(i)
					// Unwrap interface{} if needed
					if elem.Kind() == reflect.Interface && !elem.IsNil() {
						elem = elem.Elem()
					}
					ifaceArgs[i] = elem.Interface()
				}
			} else {
				ifaceArgs = []any{args[1].Interface()}
			}
		} else if args[1].Kind() == value.KindSlice {
			// Unpack native slice
			if rv, ok := args[1].ReflectValue(); ok {
				ifaceArgs = make([]any, rv.Len())
				for i := 0; i < rv.Len(); i++ {
					ifaceArgs[i] = rv.Index(i).Interface()
				}
			}
		} else {
			// Single arg, not a slice
			ifaceArgs = []any{args[1].Interface()}
		}
	} else {
		// Multiple args passed directly
		ifaceArgs = make([]any, len(args)-1)
		for i, arg := range args[1:] {
			ifaceArgs[i] = arg.Interface()
		}
	}
	
	return value.MakeString(fmt.Sprintf(format, ifaceArgs...))
}

// directErrorf is a typed wrapper for fmt.Errorf.
func directErrorf(args []value.Value) value.Value {
	if len(args) < 1 {
		return value.FromInterface(fmt.Errorf(""))
	}
	format := args[0].String()
	if len(args) == 1 {
		return value.FromInterface(fmt.Errorf(format))
	}
	// Convert remaining args to []any
	ifaceArgs := make([]any, len(args)-1)
	for i, arg := range args[1:] {
		// Unwrap reflect.Value if needed
		if rv, ok := arg.ReflectValue(); ok {
			ifaceArgs[i] = rv.Interface()
		} else {
			ifaceArgs[i] = arg.Interface()
		}
	}
	err := fmt.Errorf(format, ifaceArgs...)
	// Return as reflect.Value since error is an interface
	return value.MakeFromReflect(reflect.ValueOf(err))
}
