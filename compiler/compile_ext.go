// compile_ext.go resolves external package functions, methods, and variables.
package compiler

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
)

// extractMethodName strips SSA receiver qualification from a method name.
// SSA names look like "(*Type).Method" or "pkgpath.Method"; this extracts just "Method".
func extractMethodName(ssaName string) string {
	name := ssaName
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		name = name[idx+1:]
	}
	if idx := strings.LastIndex(name, ")"); idx >= 0 {
		rest := name[idx+1:]
		if len(rest) > 0 && rest[0] == '.' {
			name = rest[1:]
		}
	}
	return name
}

func methodOwnerPkgPath(fn *ssa.Function) string {
	if fn == nil || fn.Signature == nil {
		return ""
	}
	if fn.Pkg != nil && fn.Pkg.Pkg != nil {
		return fn.Pkg.Pkg.Path()
	}
	recv := fn.Signature.Recv()
	if recv == nil {
		return ""
	}
	recvType := recv.Type()
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}
	named, ok := recvType.(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return ""
	}
	return named.Obj().Pkg().Path()
}

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
	info := &external.ExternalMethodInfo{
		PkgPath:    methodOwnerPkgPath(fn),
		MethodName: methodName,
		FuncName:   methodName,
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

func externalBoundaryArgType(arg ssa.Value) types.Type {
	switch v := arg.(type) {
	case *ssa.MakeInterface:
		return externalBoundaryArgType(v.X)
	case *ssa.ChangeInterface:
		return externalBoundaryArgType(v.X)
	case *ssa.Convert:
		return externalBoundaryArgType(v.X)
	default:
		return arg.Type()
	}
}

func callTargetType(sig *types.Signature, argIndex int) types.Type {
	if sig == nil || sig.Params() == nil || argIndex < 0 {
		return nil
	}
	if sig.Recv() != nil {
		if argIndex == 0 {
			return sig.Recv().Type()
		}
		return callParamTargetType(sig, argIndex-1)
	}
	return callParamTargetType(sig, argIndex)
}

func callParamTargetType(sig *types.Signature, argIndex int) types.Type {
	if sig == nil || sig.Params() == nil || argIndex < 0 {
		return nil
	}
	params := sig.Params()
	if argIndex < params.Len() {
		return params.At(argIndex).Type()
	}
	if sig.Variadic() && params.Len() > 0 {
		return params.At(params.Len() - 1).Type()
	}
	return nil
}

func (c *compiler) externalTargetAllowsInterfaceProxy(sig *types.Signature, argIndex int) bool {
	return c.externalTargetHasInterfaceProxy(callTargetType(sig, argIndex))
}

func (c *compiler) externalTargetHasInterfaceProxy(t types.Type) bool {
	if t == nil {
		return false
	}
	iface, ok := t.Underlying().(*types.Interface)
	if !ok || iface.NumMethods() == 0 {
		return false
	}
	if c.lookup == nil {
		return false
	}
	nameLookup, hasNameLookup := c.lookup.(interface {
		LookupInterfaceProxy(string, string) (*external.InterfaceProxyInfo, bool)
	})
	if hasNameLookup {
		if named, ok := t.(*types.Named); ok {
			obj := named.Obj()
			if obj != nil && obj.Pkg() != nil {
				if _, ok := nameLookup.LookupInterfaceProxy(obj.Pkg().Path(), obj.Name()); ok {
					return true
				}
			}
		}
	}
	return false
}

// compileIndirectCall compiles an indirect call (closure or function value).
func (c *compiler) compileIndirectCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.validateExternalFuncValueBoundary(i.Call.Value, i.Call.Args)

	c.compileValue(i.Call.Value)

	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	numArgs := len(i.Call.Args)
	c.emit(bytecode.OpCallIndirect, uint16(numArgs))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

func (c *compiler) validateExternalFuncValueBoundary(callee ssa.Value, args []ssa.Value) {
	if c.allowUnsafeTypePass {
		return
	}
	for _, origin := range c.externalFuncOrigins(callee) {
		c.validateExternalFuncBoundary(origin, args)
	}
}

func (c *compiler) validateExternalFuncBoundary(origin externalFuncOrigin, args []ssa.Value) {
	if c.allowUnsafeTypePass {
		return
	}
	callArgs := make([]externalCallArg, len(args))
	for idx, arg := range args {
		callArgs[idx] = externalCallArg{
			SourceType:          externalBoundaryArgType(arg),
			AllowInterfaceProxy: c.externalTargetAllowsInterfaceProxy(origin.Sig, idx),
		}
	}
	if err := validateExternalCallBoundary(origin.PkgPath, origin.FuncName, callArgs); err != nil {
		c.addError(err)
	}
}

// extractReceiverTypeName extracts the package-path-qualified type name from a receiver type.
// For pointer receivers like *Reader, it unwraps the pointer.
// Returns "pkgPath.TypeName" (e.g., "encoding/json.Encoder") for use as a DirectCall lookup key.
func extractReceiverTypeName(recvType types.Type) string {
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}
	named, ok := recvType.(*types.Named)
	if !ok {
		return ""
	}
	obj := named.Obj()
	pkg := obj.Pkg()
	if pkg != nil {
		return pkg.Path() + "." + obj.Name()
	}
	return obj.Name()
}

// extractNamedType unwraps pointer types to find the underlying *types.Named type.
// Returns nil if the type is not named (e.g., interface types without a name).
func extractNamedType(t types.Type) *types.Named {
	for {
		switch tt := t.(type) {
		case *types.Named:
			return tt
		case *types.Pointer:
			t = tt.Elem()
		default:
			return nil
		}
	}
}

// extractReceiverShortName extracts the unqualified type name from a receiver type.
// For pointer receivers like *Reader, it unwraps the pointer.
// Returns just the type name (e.g., "Reader"), without package path.
func extractReceiverShortName(recvType types.Type) string {
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}
	if named, ok := recvType.(*types.Named); ok {
		return named.Obj().Name()
	}
	return ""
}

// compileReturn compiles a Return instruction.
func (c *compiler) compileReturn(i *ssa.Return) {
	if len(i.Results) == 0 {
		c.emit(bytecode.OpReturn)
		return
	}

	if len(i.Results) == 1 {
		c.compileValue(i.Results[0])
		c.emit(bytecode.OpReturnVal)
		return
	}

	for _, result := range i.Results {
		c.compileValue(result)
	}

	c.emit(bytecode.OpPack, uint16(len(i.Results)))
	c.emit(bytecode.OpReturnVal)
}
