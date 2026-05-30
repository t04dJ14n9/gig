package gentool

import (
	"fmt"
	"go/types"
	"strings"
)

type interfaceProxyInfo struct {
	TypeName          string
	FactoryName       string
	ProxyTypeName     string
	InterfaceTypeExpr string
	RequiredMethods   []string
	Code              string
}

func generateInterfaceProxy(named *types.Named, pkgRef string, typeName string) *interfaceProxyInfo {
	iface, ok := named.Underlying().(*types.Interface)
	if !ok || iface.NumMethods() == 0 {
		return nil
	}

	proxyTypeName := fmt.Sprintf("proxy_%s_%s", pkgRef, typeName)
	factoryName := fmt.Sprintf("newProxy_%s_%s", pkgRef, typeName)
	info := &interfaceProxyInfo{
		TypeName:          typeName,
		FactoryName:       factoryName,
		ProxyTypeName:     proxyTypeName,
		InterfaceTypeExpr: typeToReflectExpr(named, pkgRef),
		RequiredMethods:   make([]string, 0, iface.NumMethods()),
	}
	if info.InterfaceTypeExpr == "" {
		return nil
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("type %s struct {\n", proxyTypeName))
	b.WriteString("\tcall external.InterfaceMethodCaller\n")
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("func %s(_ value.Value, _ string, call external.InterfaceMethodCaller) (any, bool) {\n", factoryName))
	b.WriteString(fmt.Sprintf("\treturn &%s{call: call}, true\n", proxyTypeName))
	b.WriteString("}\n\n")

	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		methodCode, ok := generateInterfaceProxyMethod(proxyTypeName, method, pkgRef)
		if !ok {
			return nil
		}
		info.RequiredMethods = append(info.RequiredMethods, method.Name())
		b.WriteString(methodCode)
		b.WriteString("\n")
	}

	info.Code = b.String()
	return info
}

func generateInterfaceProxyMethod(proxyTypeName string, method *types.Func, pkgRef string) (string, bool) {
	sig, ok := method.Type().(*types.Signature)
	if !ok || sig.Variadic() || sig.Results().Len() > 1 {
		return "", false
	}

	params := sig.Params()
	paramDecls := make([]string, 0, params.Len())
	callArgs := make([]string, 0, params.Len())
	for i := 0; i < params.Len(); i++ {
		typeName := resolveTypeName(params.At(i).Type(), pkgRef)
		if typeName == "" {
			return "", false
		}
		argName := fmt.Sprintf("a%d", i)
		paramDecls = append(paramDecls, fmt.Sprintf("%s %s", argName, typeName))
		callArgs = append(callArgs, fmt.Sprintf("value.FromInterface(%s)", argName))
	}

	resultPart := ""
	resultType := types.Type(nil)
	if sig.Results().Len() == 1 {
		resultType = sig.Results().At(0).Type()
		resultName := resolveTypeName(resultType, pkgRef)
		if resultName == "" {
			return "", false
		}
		resultPart = " " + resultName
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("func (p *%s) %s(%s)%s {\n",
		proxyTypeName, method.Name(), strings.Join(paramDecls, ", "), resultPart))
	callExpr := fmt.Sprintf("p.call(%q", method.Name())
	if len(callArgs) > 0 {
		callExpr += ", " + strings.Join(callArgs, ", ")
	}
	callExpr += ")"

	if resultType == nil {
		b.WriteString(fmt.Sprintf("\t_, _ = %s\n", callExpr))
		b.WriteString("}\n")
		return b.String(), true
	}

	zero := zeroValueExpr(resultType, pkgRef)
	extract := extractArg(resultType, "result", pkgRef)
	if zero == "" || extract == "" {
		return "", false
	}
	b.WriteString(fmt.Sprintf("\tresult, ok := %s\n", callExpr))
	b.WriteString("\tif !ok {\n")
	b.WriteString(fmt.Sprintf("\t\treturn %s\n", zero))
	b.WriteString("\t}\n")
	b.WriteString(fmt.Sprintf("\treturn %s\n", extract))
	b.WriteString("}\n")
	return b.String(), true
}

func zeroValueExpr(t types.Type, pkgRef string) string {
	switch tt := t.Underlying().(type) {
	case *types.Basic:
		info := tt.Info()
		switch {
		case tt.Kind() == types.Bool:
			return "false"
		case info&types.IsString != 0:
			return `""`
		case info&(types.IsInteger|types.IsFloat|types.IsComplex) != 0:
			return "0"
		default:
			return "nil"
		}
	case *types.Interface, *types.Pointer, *types.Slice, *types.Map, *types.Chan, *types.Signature:
		return "nil"
	case *types.Struct, *types.Array:
		typeName := resolveTypeName(t, pkgRef)
		if typeName == "" {
			return ""
		}
		return typeName + "{}"
	default:
		return "nil"
	}
}

func collectInterfaceProxyImports(named *types.Named, selfPkgPath string, imports map[string]string) {
	iface, ok := named.Underlying().(*types.Interface)
	if !ok {
		return
	}
	for i := 0; i < iface.NumMethods(); i++ {
		sig, ok := iface.Method(i).Type().(*types.Signature)
		if !ok {
			continue
		}
		params := sig.Params()
		for j := 0; j < params.Len(); j++ {
			collectTypeImports(params.At(j).Type(), selfPkgPath, imports)
		}
		results := sig.Results()
		for j := 0; j < results.Len(); j++ {
			collectTypeImports(results.At(j).Type(), selfPkgPath, imports)
		}
	}
}
