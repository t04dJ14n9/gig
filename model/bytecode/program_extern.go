package bytecode

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

// ResolvedCall is a pre-resolved external function/method entry.
// It is built once by ResolveExternCalls after compilation and is immutable
// afterwards. Because it lives on CompiledProgram, all VMs sharing the same
// program can access it without locks.
type ResolvedCall struct {
	// PkgPath is the Go import path for the package that owns the function.
	PkgPath string

	// FuncName is the exported function name used for diagnostics.
	FuncName string

	// IsStdlib is true when PkgPath belongs to the Go standard library.
	// Boundary validation is skipped for stdlib calls; storing the decision here
	// avoids repeated string scans in tight external-call loops.
	IsStdlib bool

	// DirectCall is the fast-path wrapper. Nil means use reflect.
	DirectCall func(args []value.Value) value.Value

	// Fn is the reflect.Value of the function (slow path).
	Fn reflect.Value

	// FnType is the function's reflect.Type (slow path).
	FnType reflect.Type

	// IsVariadic indicates whether the function takes variadic arguments.
	IsVariadic bool

	// NumIn is the number of declared parameters.
	NumIn int

	// MethodName is set for method calls (ExternalMethodInfo), empty for functions.
	MethodName string
}

// ResolveExternCalls pre-resolves external functions and methods in the constant pool.
func (p *CompiledProgram) ResolveExternCalls() {
	p.ExternCalls = make([]*ResolvedCall, len(p.Constants))
	for i, c := range p.Constants {
		p.ExternCalls[i] = ResolveConstant(c)
	}
}

// ResolveConstant creates a ResolvedCall from a constant pool entry.
// Exported so the VM's fallback path can use it.
func ResolveConstant(c any) *ResolvedCall {
	if c == nil {
		return nil
	}

	switch info := c.(type) {
	case *external.ExternalFuncInfo:
		return resolveExternalFunc(info)
	case *external.ExternalMethodInfo:
		return &ResolvedCall{
			DirectCall: info.DirectCall,
			MethodName: info.MethodName,
			IsStdlib:   info.IsStdlib,
		}
	default:
		return resolveLegacyFunc(c)
	}
}

func resolveExternalFunc(info *external.ExternalFuncInfo) *ResolvedCall {
	rc := &ResolvedCall{
		PkgPath:    info.PkgPath,
		FuncName:   info.FuncName,
		IsStdlib:   info.IsStdlib,
		DirectCall: info.DirectCall,
		IsVariadic: info.IsVariadic,
		NumIn:      info.NumIn,
	}
	if info.Func == nil {
		return rc
	}

	rc.Fn = reflect.ValueOf(info.Func)
	if rc.Fn.Kind() != reflect.Func {
		return rc
	}

	rc.FnType = rc.Fn.Type()
	if !rc.IsVariadic && rc.FnType.IsVariadic() {
		rc.IsVariadic = true
		rc.NumIn = rc.FnType.NumIn()
	}
	return rc
}

func resolveLegacyFunc(c any) *ResolvedCall {
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Func {
		return nil
	}
	ft := rv.Type()
	return &ResolvedCall{
		Fn:         rv,
		FnType:     ft,
		IsVariadic: ft.IsVariadic(),
		NumIn:      ft.NumIn(),
		IsStdlib:   true,
	}
}
