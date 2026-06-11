package gentool

import (
	"go/types"
	"strings"
)

type packageRefs struct {
	ImportAlias string
	PkgRef      string
	GoPkgName   string
}

func newPackageRefs(path, pkgName string) packageRefs {
	// Generated wrappers need two names for packages with slashes: a legal Go
	// import alias and the expression prefix used in emitted registrations.
	// Standard packages can keep their real package name for both.
	refs := packageRefs{
		PkgRef:    pkgName,
		GoPkgName: sanitizePkgName(path),
	}
	if strings.Contains(path, "/") {
		refs.ImportAlias = sanitizePkgName(path)
		refs.PkgRef = refs.ImportAlias
	}
	return refs
}

type packageGenerationPlan struct {
	Path     string
	PkgName  string
	GoName   string
	Refs     packageRefs
	Symbols  packageSymbols
	Adapters generatedAdapters
}

func buildPackageGenerationPlan(path, pkgName, goName string, refs packageRefs, symbols packageSymbols) packageGenerationPlan {
	// The plan is the handoff between package analysis and source emission.
	// Keeping it explicit makes import requirements, generated helpers, and
	// registration data reviewable without reading string-builder code.
	return packageGenerationPlan{
		Path:     path,
		PkgName:  pkgName,
		GoName:   goName,
		Refs:     refs,
		Symbols:  symbols,
		Adapters: buildGeneratedAdapters(path, refs.PkgRef, symbols),
	}
}

type generatedAdapters struct {
	NeedReflect       bool
	HasDirectCalls    bool
	CrossPkgImports   map[string]string
	MethodDirectCalls []*methodDirectCallInfo
	InterfaceProxies  []*interfaceProxyInfo
}

func buildGeneratedAdapters(path, pkgRef string, symbols packageSymbols) generatedAdapters {
	// Direct calls, method direct calls, and interface proxies all affect the
	// import list. Collect them together before emission so generated files are
	// formatted once with a complete dependency set.
	adapters := generatedAdapters{
		NeedReflect:     symbols.needReflect(pkgRef),
		CrossPkgImports: make(map[string]string),
	}
	for _, fi := range symbols.Funcs {
		if fi.DirectCall == "" {
			continue
		}
		adapters.HasDirectCalls = true
		collectCrossPkgImports(fi.Sig, path, adapters.CrossPkgImports)
	}
	adapters.collectTypeAdapters(path, pkgRef, symbols.Types)
	return adapters
}

func (s packageSymbols) needReflect(pkgRef string) bool {
	for _, ti := range s.Types {
		if typeToReflectExpr(ti.Obj.Type(), pkgRef) != "" {
			return true
		}
	}
	return false
}

func (a *generatedAdapters) collectTypeAdapters(path, pkgRef string, typesToRegister []*typeInfo) {
	// Named interfaces produce interface proxy factories; concrete named types
	// produce method direct-call wrappers where supported. The distinction must
	// happen on the named type, not only on the registered type name string.
	for _, ti := range typesToRegister {
		named, ok := ti.Obj.Type().(*types.Named)
		if !ok {
			continue
		}
		if _, isIface := named.Underlying().(*types.Interface); isIface {
			a.collectInterfaceProxy(path, pkgRef, ti.Name, named)
			continue
		}
		a.collectMethodDirectCalls(path, pkgRef, ti.Name, named)
	}
}

func (a *generatedAdapters) collectInterfaceProxy(path, pkgRef, typeName string, named *types.Named) {
	proxy := generateInterfaceProxy(named, pkgRef, typeName)
	if proxy == nil {
		return
	}
	a.InterfaceProxies = append(a.InterfaceProxies, proxy)
	collectInterfaceProxyImports(named, path, a.CrossPkgImports)
}

func (a *generatedAdapters) collectMethodDirectCalls(path, pkgRef, typeName string, named *types.Named) {
	methodDCs := generateMethodDirectCalls(named, pkgRef, typeName)
	if len(methodDCs) == 0 {
		return
	}
	a.HasDirectCalls = true
	a.MethodDirectCalls = append(a.MethodDirectCalls, methodDCs...)
	for _, mdc := range methodDCs {
		collectMethodImports(named, mdc.MethodName, path, a.CrossPkgImports)
	}
}
