package compiler

import "golang.org/x/tools/go/ssa"

func (c *compiler) externalFuncOriginForFunction(fn *ssa.Function) (externalFuncOrigin, bool) {
	// A function already indexed by this compiler is script-owned and must not
	// be treated as an external package boundary. Only functions outside the
	// compiled graph are candidates for third-party validation.
	if fn == nil || fn.Signature == nil {
		return externalFuncOrigin{}, false
	}
	if _, known := c.funcIndex[fn]; known {
		return externalFuncOrigin{}, false
	}

	if fn.Signature.Recv() != nil {
		return externalMethodOrigin(fn)
	}
	return externalPackageFunctionOrigin(fn)
}

func externalMethodOrigin(fn *ssa.Function) (externalFuncOrigin, bool) {
	// External method values often have package ownership encoded on the
	// receiver type rather than fn.Pkg. Resolve that receiver owner before
	// deciding whether custom Gig values may cross the boundary.
	pkgPath := methodOwnerPkgPath(fn)
	if isLocalExternalOriginPath(pkgPath) {
		return externalFuncOrigin{}, false
	}
	return externalFuncOrigin{
		PkgPath:  pkgPath,
		FuncName: extractMethodName(fn.Name()),
		Sig:      fn.Signature,
	}, true
}

func externalPackageFunctionOrigin(fn *ssa.Function) (externalFuncOrigin, bool) {
	if fn.Pkg == nil || fn.Pkg.Pkg == nil {
		return externalFuncOrigin{}, false
	}
	pkgPath := fn.Pkg.Pkg.Path()
	if isLocalExternalOriginPath(pkgPath) {
		return externalFuncOrigin{}, false
	}
	return externalFuncOrigin{
		PkgPath:  pkgPath,
		FuncName: fn.Name(),
		Sig:      fn.Signature,
	}, true
}

func isLocalExternalOriginPath(pkgPath string) bool {
	return pkgPath == "" || pkgPath == "main" || pkgPath == "command-line-arguments"
}
