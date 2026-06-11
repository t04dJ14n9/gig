package vm

import (
	"reflect"
	"strings"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// callCompiledMethod searches the compiled function table for a method with the
// given name and calls it. This is the fallback path for invoke (interface method)
// calls when reflection-based MethodByName fails.
func (v *vm) callCompiledMethod(methodName string, receiverTypeName string, args []value.Value) error {
	if len(args) > 0 {
		if fn, methodReceiver, ok := selectCompiledMethodCandidate(v.program, methodName, receiverTypeName, args[0]); ok {
			if len(args) > 0 && shouldPanicOnNilValueReceiver(args[0], fn) {
				v.panicking = true
				v.panicVal = value.FromInterface("runtime error: invalid memory address or nil pointer dereference")
				return nil
			}
			for i, arg := range args {
				if i == 0 {
					arg = methodReceiver
				}
				v.push(arg)
			}
			v.callCompiledFunction(fn.FuncIdx, len(args))
			return nil
		}
	}

	v.push(value.MakeNil())
	return nil
}

const (
	compiledMethodNoMatch = iota
	compiledMethodEmbeddedMatch
	compiledMethodExactMatch
)

func selectCompiledMethodCandidate(
	program *bytecode.CompiledProgram,
	methodName string,
	receiverTypeName string,
	receiver value.Value,
) (*bytecode.CompiledFunction, value.Value, bool) {
	if program == nil {
		return nil, value.MakeNil(), false
	}
	var bestFn *bytecode.CompiledFunction
	var bestReceiver value.Value
	bestScore := compiledMethodNoMatch
	for _, fn := range program.MethodsByName[methodName] {
		methodReceiver, score := receiverForCompiledMethodCandidate(methodName, receiverTypeName, receiver, fn, program)
		if score <= bestScore {
			continue
		}
		bestFn = fn
		bestReceiver = methodReceiver
		bestScore = score
		if score == compiledMethodExactMatch {
			break
		}
	}
	return bestFn, bestReceiver, bestScore != compiledMethodNoMatch
}

func shouldPanicOnNilValueReceiver(receiver value.Value, fn *bytecode.CompiledFunction) bool {
	if fn == nil || !fn.HasReceiver || fn.ReceiverIsPointer {
		return false
	}
	rv, ok := receiver.ReflectValue()
	return ok && rv.Kind() == reflect.Ptr && rv.IsNil()
}

func receiverForCompiledMethodTarget(
	methodName string,
	receiver value.Value,
	fn *bytecode.CompiledFunction,
	prog *bytecode.CompiledProgram,
) (value.Value, bool) {
	methodReceiver, score := receiverForCompiledMethodCandidate(methodName, "", receiver, fn, prog)
	return methodReceiver, score != compiledMethodNoMatch
}

func receiverForCompiledMethodCandidate(
	methodName string,
	receiverTypeName string,
	receiver value.Value,
	fn *bytecode.CompiledFunction,
	prog *bytecode.CompiledProgram,
) (value.Value, int) {
	normalized := receiverForCompiledMethod(methodName, receiver)
	if fn == nil {
		return normalized, compiledMethodExactMatch
	}
	if fn.ReceiverTypeName == "" {
		return normalized, compiledMethodNoMatch
	}
	if receiverTypeName != "" && fn.ReceiverTypeName == receiverTypeName {
		return normalized, compiledMethodExactMatch
	}
	if dyn, ok := receiver.InterpretedInterface(); ok {
		if dyn.TypeName == fn.ReceiverTypeName {
			if fn.ReceiverIsPointer && !dyn.IsPointer {
				return normalized, compiledMethodNoMatch
			}
			return normalized, compiledMethodExactMatch
		}
		if embedded, ok := embeddedReceiverForCompiledMethod(normalized, fn, prog); ok {
			return embedded, compiledMethodEmbeddedMatch
		}
		return normalized, compiledMethodNoMatch
	}
	if inferReceiverTypeName(normalized, prog) == fn.ReceiverTypeName {
		return normalized, compiledMethodExactMatch
	}
	if embedded, ok := embeddedReceiverForCompiledMethod(normalized, fn, prog); ok {
		return embedded, compiledMethodEmbeddedMatch
	}
	return normalized, compiledMethodNoMatch
}

func embeddedReceiverForCompiledMethod(receiver value.Value, fn *bytecode.CompiledFunction, prog *bytecode.CompiledProgram) (value.Value, bool) {
	rv, ok := receiver.ReflectValue()
	if !ok {
		iface := receiver.Interface()
		if iface == nil {
			return value.MakeNil(), false
		}
		rv = reflect.ValueOf(iface)
	}
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return value.MakeNil(), false
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return value.MakeNil(), false
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return value.MakeNil(), false
	}
	return embeddedReceiverFromStruct(rv, fn, prog)
}

func embeddedReceiverFromStruct(structVal reflect.Value, fn *bytecode.CompiledFunction, prog *bytecode.CompiledProgram) (value.Value, bool) {
	structType := structVal.Type()
	for i := 0; i < structVal.NumField(); i++ {
		structField := structType.Field(i)
		if !structField.Anonymous && structField.Tag.Get("gig_embed") != "1" {
			continue
		}
		field := structVal.Field(i)
		fieldType := structField.Type
		baseType := fieldType
		if baseType.Kind() == reflect.Ptr {
			baseType = baseType.Elem()
		}
		if resolveTypeName(baseType, prog) != fn.ReceiverTypeName {
			continue
		}
		if fn.ReceiverIsPointer {
			if field.Kind() == reflect.Ptr {
				return value.MakeFromReflect(field), true
			}
			if field.CanAddr() {
				return value.MakeFromReflect(field.Addr()), true
			}
			ptr := reflect.New(field.Type())
			ptr.Elem().Set(field)
			return value.MakeFromReflect(ptr), true
		}
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return value.MakeNil(), false
			}
			return value.MakeFromReflect(field.Elem()), true
		}
		return value.MakeFromReflect(field), true
	}
	return value.MakeNil(), false
}

// inferReceiverTypeName tries to extract a type name from a runtime value.Value receiver.
func inferReceiverTypeName(receiver value.Value, prog *bytecode.CompiledProgram) string {
	if dyn, ok := receiver.InterpretedInterface(); ok {
		return dyn.TypeName
	}
	rv, ok := receiver.ReflectValue()
	if !ok {
		return ""
	}
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	var t reflect.Type
	if rv.Kind() == reflect.Ptr {
		t = rv.Type().Elem()
	} else if rv.IsValid() {
		t = rv.Type()
	} else {
		return ""
	}
	return resolveTypeName(t, prog)
}

// resolveTypeName returns a human-readable type name, trying (in order):
// 1. reflect.Type.Name() (works for named types)
// 2. Program-level ReflectTypeNames registry
// 3. Scanning unexported struct field PkgPath for the "#" suffix heuristic
func resolveTypeName(t reflect.Type, prog *bytecode.CompiledProgram) string {
	if t.Name() != "" {
		return t.Name()
	}
	if prog != nil {
		if name := prog.LookupTypeName(t); name != "" {
			return name
		}
	}
	return pkgPathTypeName(t)
}

// pkgPathTypeName scans unexported struct fields for a PkgPath containing "#",
// which embeds the original package path + type name (e.g. "pkg/path#TypeName").
// Returns the type name portion, or "" if not found.
func pkgPathTypeName(t reflect.Type) string {
	if t.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < t.NumField(); i++ {
		pkgPath := t.Field(i).PkgPath
		if idx := strings.LastIndex(pkgPath, "#"); idx >= 0 {
			qualName := pkgPath[idx+1:]
			if dotIdx := strings.LastIndex(qualName, "."); dotIdx >= 0 {
				return qualName[dotIdx+1:]
			}
			return qualName
		}
	}
	return ""
}
