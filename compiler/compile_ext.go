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
	if fnVal != nil {
		if rv := reflect.ValueOf(fnVal); rv.Kind() == reflect.Func {
			rt := rv.Type()
			info.IsVariadic = rt.IsVariadic()
			info.NumIn = rt.NumIn()
		}
	}
	return info
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

	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	sig := fn.Signature
	if sig.Recv() != nil {
		methodInfo := c.lookupExternalMethodInfo(fn)
		funcIdx := c.addConstant(methodInfo)
		numArgs := len(i.Call.Args)
		c.emitCallOp(bytecode.OpCallExternal, funcIdx, numArgs)
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	}

	// Use injected PackageLookup instead of direct importer access
	extFuncInfo := c.lookupExternalFuncInfo(fn)
	if extFuncInfo == nil {
		pkgPath := ""
		if fn.Pkg != nil && fn.Pkg.Pkg != nil {
			pkgPath = fn.Pkg.Pkg.Path()
		}
		if fn.Name() == "init" && pkgPath != "" && pkgPath != "main" && pkgPath != "command-line-arguments" {
			// Imported binary packages are registered after their Go init functions
			// have already run. SSA still includes init stubs for import ordering;
			// treating those as external reflect calls produces invalid call entries.
			return
		}
		c.addError(fmt.Errorf("unresolved external function %s.%s", pkgPath, fn.Name()))
		c.emit(bytecode.OpNil)
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	}

	funcIdx := c.addConstant(extFuncInfo)
	numArgs := len(i.Call.Args)
	c.emitCallOp(bytecode.OpCallExternal, funcIdx, numArgs)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
