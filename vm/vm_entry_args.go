package vm

import (
	"fmt"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) validateAndPrepareEntryArgs(fn *bytecode.CompiledFunction, args []value.Value) ([]value.Value, error) {
	if err := validateEntryArgCount(fn, len(args)); err != nil {
		return nil, err
	}
	if err := v.validateEntryArgTypes(fn, args); err != nil {
		return nil, err
	}
	return v.prepareEntryArgs(fn, args)
}

func validateEntryArgCount(fn *bytecode.CompiledFunction, got int) error {
	if fn.IsVariadic {
		min := fn.NumParams - 1
		if min < 0 {
			min = 0
		}
		if got < min {
			return fmt.Errorf("function %q expects at least %d arguments, got %d", fn.Name, min, got)
		}
		return nil
	}
	if got != fn.NumParams {
		return fmt.Errorf("function %q expects %d arguments, got %d", fn.Name, fn.NumParams, got)
	}
	return nil
}

func (v *vm) validateEntryArgTypes(fn *bytecode.CompiledFunction, args []value.Value) error {
	if len(fn.ParamTypes) == 0 {
		return nil
	}

	fixedCount := fn.NumParams
	if fn.IsVariadic {
		fixedCount--
	}
	if fixedCount > len(fn.ParamTypes) {
		fixedCount = len(fn.ParamTypes)
	}

	for i := 0; i < fixedCount; i++ {
		paramType := typeToReflect(fn.ParamTypes[i], v.program)
		if paramType == nil {
			continue
		}
		if err := validateEntryArgType(args[i], paramType); err != nil {
			return fmt.Errorf("function %q argument %d: %w", fn.Name, i, err)
		}
	}
	return nil
}

func validateEntryArgType(arg value.Value, paramType reflect.Type) (err error) {
	if arg.Kind() == value.KindNil && !entryTypeAcceptsNil(paramType) {
		return fmt.Errorf("cannot use nil as %s", paramType)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("cannot convert %s to %s: %v", arg.Kind(), paramType, r)
		}
	}()

	rv := arg.ToReflectValue(paramType)
	if !rv.IsValid() {
		return nil
	}
	if rv.Type().AssignableTo(paramType) || rv.Type().ConvertibleTo(paramType) {
		return nil
	}
	return fmt.Errorf("cannot use %s as %s", rv.Type(), paramType)
}

func (v *vm) prepareEntryArgs(fn *bytecode.CompiledFunction, args []value.Value) ([]value.Value, error) {
	if !fn.IsVariadic || fn.NumParams == 0 || fn.VariadicParamType == nil {
		return args, nil
	}

	fixedCount := fn.NumParams - 1
	if len(args) < fixedCount {
		return args, nil
	}

	variadicType := typeToReflect(fn.VariadicParamType, v.program)
	if variadicType == nil || variadicType.Kind() != reflect.Slice {
		return args, nil
	}

	prepared := make([]value.Value, fn.NumParams)
	copy(prepared, args[:fixedCount])

	if len(args) == fn.NumParams && isEntryVariadicSliceArg(args[fixedCount], variadicType) {
		prepared[fixedCount] = args[fixedCount]
		return prepared, nil
	}

	packed, err := packEntryVariadicArgs(args[fixedCount:], variadicType)
	if err != nil {
		return nil, fmt.Errorf("function %q variadic argument: %w", fn.Name, err)
	}
	prepared[fixedCount] = packed
	return prepared, nil
}

func isEntryVariadicSliceArg(arg value.Value, variadicType reflect.Type) bool {
	if rv, ok := arg.ReflectValue(); ok && rv.IsValid() && rv.Kind() == reflect.Slice {
		return rv.Type().AssignableTo(variadicType) || rv.Type().ConvertibleTo(variadicType)
	}
	if arg.Kind() == value.KindBytes {
		return reflect.TypeOf([]byte(nil)).AssignableTo(variadicType)
	}
	if _, ok := arg.IntSlice(); ok {
		elemKind := variadicType.Elem().Kind()
		return elemKind == reflect.Int || elemKind == reflect.Int64
	}
	if _, ok := arg.ValueSlice(); ok {
		return reflect.TypeOf([]value.Value(nil)).AssignableTo(variadicType)
	}
	return false
}

func packEntryVariadicArgs(args []value.Value, variadicType reflect.Type) (value.Value, error) {
	elemType := variadicType.Elem()
	slice := reflect.MakeSlice(variadicType, len(args), len(args))
	for i, arg := range args {
		elem, err := entryVariadicElement(arg, elemType)
		if err != nil {
			return value.MakeNil(), fmt.Errorf("element %d: %w", i, err)
		}
		slice.Index(i).Set(elem)
	}
	return value.MakeFromReflect(slice), nil
}

func entryVariadicElement(arg value.Value, elemType reflect.Type) (rv reflect.Value, err error) {
	if arg.Kind() == value.KindNil && !entryTypeAcceptsNil(elemType) {
		return reflect.Value{}, fmt.Errorf("cannot use nil as %s", elemType)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("cannot convert %s to %s: %v", arg.Kind(), elemType, r)
		}
	}()

	rv = arg.ToReflectValue(elemType)
	if !rv.IsValid() {
		return reflect.Zero(elemType), nil
	}
	if rv.Type().AssignableTo(elemType) {
		return rv, nil
	}
	if rv.Type().ConvertibleTo(elemType) {
		return rv.Convert(elemType), nil
	}
	return reflect.Value{}, fmt.Errorf("cannot use %s as %s", rv.Type(), elemType)
}

func entryTypeAcceptsNil(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}
