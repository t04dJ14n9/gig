package gentool

import (
	"fmt"
	"go/types"
)

// extractUnderlyingWithPkgRef generates the Go expression to extract a value from
// a value.Value for non-named underlying types (pointers, slices, maps, etc.).
func extractUnderlyingWithPkgRef(t types.Type, valExpr string, pkgRef string) string {
	switch ut := t.(type) {
	case *types.Basic:
		return extractBasic(ut, valExpr)
	case *types.Slice:
		if _, ok := ut.Elem().Underlying().(*types.Basic); ok {
			return extractSlice(ut, valExpr, pkgRef)
		}
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Interface:
		return extractInterface(ut, valExpr, pkgRef)
	case *types.Pointer:
		return extractPointer(ut, valExpr, pkgRef)
	case *types.Struct:
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Map:
		return extractMap(ut, valExpr, pkgRef)
	case *types.Chan:
		return extractChan(ut, valExpr, pkgRef)
	case *types.Signature:
		funcTypeName := resolveFuncTypeName(ut, pkgRef)
		if funcTypeName != "" {
			return fmt.Sprintf("%s.Interface().(%s)", valExpr, funcTypeName)
		}
		return fmt.Sprintf("%s.Interface()", valExpr)
	case *types.Array:
		arrTypeName := resolveArrayTypeName(ut, pkgRef)
		if arrTypeName != "" {
			return fmt.Sprintf("%s.Interface().(%s)", valExpr, arrTypeName)
		}
		return fmt.Sprintf("%s.Interface()", valExpr)
	default:
		return ""
	}
}

func extractInterface(iface *types.Interface, valExpr string, pkgRef string) string {
	if iface.NumMethods() == 0 {
		return fmt.Sprintf("%s.Interface()", valExpr)
	}
	ifaceName := resolveInterfaceTypeName(iface, pkgRef)
	if ifaceName != "" {
		return fmt.Sprintf("%s.Interface().(%s)", valExpr, ifaceName)
	}
	return fmt.Sprintf("%s.Interface()", valExpr)
}

func extractPointer(ptr *types.Pointer, valExpr string, pkgRef string) string {
	if named, ok := ptr.Elem().(*types.Named); ok {
		qualName := resolveQualifiedName(named, pkgRef)
		if qualName != "" {
			return fmt.Sprintf("%s.Interface().(*%s)", valExpr, qualName)
		}
	}
	if bt, ok := ptr.Elem().(*types.Basic); ok {
		return fmt.Sprintf("%s.Interface().(*%s)", valExpr, bt.Name())
	}
	return fmt.Sprintf("%s.Interface()", valExpr)
}

func extractMap(mt *types.Map, valExpr string, pkgRef string) string {
	keyName := resolveTypeName(mt.Key(), pkgRef)
	elemName := resolveTypeName(mt.Elem(), pkgRef)
	if keyName != "" && elemName != "" {
		return fmt.Sprintf("%s.Interface().(map[%s]%s)", valExpr, keyName, elemName)
	}
	return fmt.Sprintf("%s.Interface()", valExpr)
}

func extractChan(ch *types.Chan, valExpr string, pkgRef string) string {
	elemName := resolveTypeName(ch.Elem(), pkgRef)
	if elemName == "" {
		return fmt.Sprintf("%s.Interface()", valExpr)
	}
	var dirStr string
	switch ch.Dir() {
	case types.SendRecv:
		dirStr = fmt.Sprintf("chan %s", elemName)
	case types.SendOnly:
		dirStr = fmt.Sprintf("chan<- %s", elemName)
	case types.RecvOnly:
		dirStr = fmt.Sprintf("<-chan %s", elemName)
	}
	if dirStr != "" {
		return fmt.Sprintf("%s.Interface().(%s)", valExpr, dirStr)
	}
	return fmt.Sprintf("%s.Interface()", valExpr)
}
