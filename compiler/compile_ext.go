// compile_ext.go resolves external package functions, methods, and variables.
package compiler

import (
	"fmt"
	"reflect"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
)

// compileExternalFuncValue compiles an external function reference as a value
// (not a direct call). This is used when an external function like
// strings.TrimSpace is passed as a callback argument.
// It looks up the actual Go function and stores it as a constant,
// so OpCallIndirect can call it via reflection later.
func (c *compiler) compileExternalFuncValue(fn *ssa.Function) {
	var fnVal any
	if fn.Pkg != nil && c.lookup != nil {
		pkgPath := fn.Pkg.Pkg.Path()
		if f, _, ok := c.lookup.LookupExternalFunc(pkgPath, fn.Name()); ok {
			fnVal = f
		}
	}
	if fnVal == nil {
		// Could not resolve — emit nil (best effort)
		c.emit(bytecode.OpNil)
		return
	}
	// Store the Go function as a constant; OpConst will push it as a
	// value.FromInterface which preserves the reflect.Func type.
	// OpCallIndirect's reflect.Func branch handles calling it.
	funcIdx := c.addConstant(fnVal)
	c.emit(bytecode.OpConst, uint16(funcIdx))
}

func (c *compiler) lookupExternalFuncInfo(fn *ssa.Function) *external.ExternalFuncInfo {
	if fn == nil || fn.Pkg == nil || fn.Pkg.Pkg == nil || c.lookup == nil {
		return nil
	}
	pkgPath := fn.Pkg.Pkg.Path()
	fnVal, directCall, ok := c.lookup.LookupExternalFunc(pkgPath, fn.Name())
	if !ok {
		return nil
	}
	info := &external.ExternalFuncInfo{
		PkgPath:    pkgPath,
		FuncName:   fn.Name(),
		IsStdlib:   isStdlibPath(pkgPath),
		Func:       fnVal,
		DirectCall: directCall,
	}
	attachExternalFuncReflectMetadata(info, fnVal)
	return info
}

// attachExternalFuncReflectMetadata precomputes metadata the VM needs for
// variadic DirectCall unpacking and reflect fallback. DirectCall identity stays
// owned by the registry lookup; this helper only mirrors reflect.Type facts.
func attachExternalFuncReflectMetadata(info *external.ExternalFuncInfo, fnVal any) {
	if fnVal == nil {
		return
	}
	rv := reflect.ValueOf(fnVal)
	if rv.Kind() != reflect.Func {
		return
	}
	rt := rv.Type()
	info.IsVariadic = rt.IsVariadic()
	info.NumIn = rt.NumIn()
}

func (c *compiler) lookupExternalMethodInfo(fn *ssa.Function) *external.ExternalMethodInfo {
	if fn == nil || fn.Signature == nil || fn.Signature.Recv() == nil {
		return nil
	}
	methodName := extractMethodName(fn.Name())
	pkgPath := methodOwnerPkgPath(fn)
	info := &external.ExternalMethodInfo{
		PkgPath:    pkgPath,
		MethodName: methodName,
		FuncName:   methodName,
		IsStdlib:   isStdlibPath(pkgPath),
	}
	if c.lookup != nil {
		typeName := extractReceiverTypeName(fn.Signature.Recv().Type())
		if typeName != "" {
			if dc, ok := c.lookup.LookupMethodDirectCall(typeName, methodName); ok {
				info.DirectCall = dc
			}
		}
	}
	return info
}

// compileExternalStaticCall compiles a call to an external package function.
// It uses the injected PackageLookup to resolve the function, avoiding direct importer dependency.
func (c *compiler) compileExternalStaticCall(i *ssa.Call, fn *ssa.Function, resultIdx int) {
	// Validate: reject user-defined types flowing into third-party packages.
	if origin, ok := c.externalFuncOriginForFunction(fn); ok {
		c.validateExternalFuncBoundary(origin, i.Call.Args)
	}

	c.compileExternalCallArgs(i.Call.Args)

	sig := fn.Signature
	if sig.Recv() != nil {
		c.compileResolvedExternalMethodCall(fn, len(i.Call.Args), resultIdx)
		return
	}

	// Use injected PackageLookup instead of direct importer access
	extFuncInfo := c.lookupExternalFuncInfo(fn)
	if extFuncInfo == nil {
		c.compileUnresolvedExternalFunction(fn, resultIdx)
		return
	}

	c.emitExternalCallResult(extFuncInfo, len(i.Call.Args), resultIdx)
}

func (c *compiler) compileExternalCallArgs(args []ssa.Value) {
	for _, arg := range args {
		c.compileValue(arg)
	}
}

func (c *compiler) compileResolvedExternalMethodCall(fn *ssa.Function, numArgs int, resultIdx int) {
	c.emitExternalCallResult(c.lookupExternalMethodInfo(fn), numArgs, resultIdx)
}

func (c *compiler) compileUnresolvedExternalFunction(fn *ssa.Function, resultIdx int) {
	pkgPath := externalFunctionPkgPath(fn)
	if shouldSkipUnresolvedExternalFunction(fn.Name(), pkgPath) {
		return
	}
	c.addError(fmt.Errorf("unresolved external function %s.%s", pkgPath, fn.Name()))
	c.emit(bytecode.OpNil)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

func externalFunctionPkgPath(fn *ssa.Function) string {
	if fn != nil && fn.Pkg != nil && fn.Pkg.Pkg != nil {
		return fn.Pkg.Pkg.Path()
	}
	return ""
}

func shouldSkipUnresolvedExternalFunction(funcName, pkgPath string) bool {
	// Imported binary packages are registered after their Go init functions have
	// already run. SSA still includes init stubs for import ordering; treating
	// those as external reflect calls produces invalid call entries.
	return funcName == "init" && pkgPath != "" && pkgPath != "main" && pkgPath != "command-line-arguments"
}

func (c *compiler) emitExternalCallResult(callInfo any, numArgs int, resultIdx int) {
	funcIdx := c.addConstant(callInfo)
	c.emitCallOp(bytecode.OpCallExternal, funcIdx, numArgs)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
