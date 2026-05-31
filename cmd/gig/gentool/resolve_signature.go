package gentool

import (
	"fmt"
	"go/types"
	"strings"
)

// resolveFuncTypeName returns the Go type string for a function signature,
// e.g. "func(string) bool" or "func(int, int) int".
func resolveFuncTypeName(sig *types.Signature, pkgRef string) string {
	params := sig.Params()
	results := sig.Results()

	var paramStrs []string
	for i := 0; i < params.Len(); i++ {
		pName := resolveTypeName(params.At(i).Type(), pkgRef)
		if pName == "" {
			return ""
		}
		paramStrs = append(paramStrs, pName)
	}

	switch results.Len() {
	case 0:
		return fmt.Sprintf("func(%s)", strings.Join(paramStrs, ", "))
	case 1:
		rName := resolveTypeName(results.At(0).Type(), pkgRef)
		if rName == "" {
			return ""
		}
		return fmt.Sprintf("func(%s) %s", strings.Join(paramStrs, ", "), rName)
	default:
		var retStrs []string
		for i := 0; i < results.Len(); i++ {
			rName := resolveTypeName(results.At(i).Type(), pkgRef)
			if rName == "" {
				return ""
			}
			retStrs = append(retStrs, rName)
		}
		return fmt.Sprintf("func(%s) (%s)", strings.Join(paramStrs, ", "), strings.Join(retStrs, ", "))
	}
}

// resolveInterfaceTypeName returns the Go type string for an unnamed interface
// with methods, e.g. "interface{ Printf(string, ...interface{}) }".
func resolveInterfaceTypeName(iface *types.Interface, pkgRef string) string {
	var methodStrs []string
	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)
		sig := m.Type().(*types.Signature)
		methodStr := resolveMethodSigStr(m.Name(), sig, pkgRef)
		if methodStr == "" {
			return ""
		}
		methodStrs = append(methodStrs, methodStr)
	}
	return "interface{ " + strings.Join(methodStrs, "; ") + " }"
}

// resolveMethodSigStr returns the Go method signature string, e.g. "Printf(string, ...interface{})".
func resolveMethodSigStr(name string, sig *types.Signature, pkgRef string) string {
	params := sig.Params()
	results := sig.Results()

	var paramStrs []string
	for i := 0; i < params.Len(); i++ {
		pType := params.At(i).Type()
		if sig.Variadic() && i == params.Len()-1 {
			sliceType, ok := pType.(*types.Slice)
			if !ok {
				return ""
			}
			elemName := resolveTypeName(sliceType.Elem(), pkgRef)
			if elemName == "" {
				return ""
			}
			paramStrs = append(paramStrs, "..."+elemName)
		} else {
			pName := resolveTypeName(pType, pkgRef)
			if pName == "" {
				return ""
			}
			paramStrs = append(paramStrs, pName)
		}
	}

	var resultPart string
	switch results.Len() {
	case 0:
		resultPart = ""
	case 1:
		rName := resolveTypeName(results.At(0).Type(), pkgRef)
		if rName == "" {
			return ""
		}
		resultPart = " " + rName
	default:
		var retStrs []string
		for i := 0; i < results.Len(); i++ {
			rName := resolveTypeName(results.At(i).Type(), pkgRef)
			if rName == "" {
				return ""
			}
			retStrs = append(retStrs, rName)
		}
		resultPart = " (" + strings.Join(retStrs, ", ") + ")"
	}

	return fmt.Sprintf("%s(%s)%s", name, strings.Join(paramStrs, ", "), resultPart)
}
