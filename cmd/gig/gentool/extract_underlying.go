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
		return extractUnderlyingSlice(ut, valExpr, pkgRef)
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
		return extractResolvedTypeAssertion(resolveFuncTypeName(ut, pkgRef), valExpr)
	case *types.Array:
		return extractResolvedTypeAssertion(resolveArrayTypeName(ut, pkgRef), valExpr)
	default:
		return ""
	}
}

func extractUnderlyingSlice(st *types.Slice, valExpr string, pkgRef string) string {
	// extractSlice only has specialized zero-reflection paths for basic element
	// slices. Composite element slices stay on Interface(), matching the prior
	// fallback and avoiding generated casts for shapes resolveTypeName cannot
	// express safely.
	if _, ok := st.Elem().Underlying().(*types.Basic); ok {
		return extractSlice(st, valExpr, pkgRef)
	}
	return fmt.Sprintf("%s.Interface()", valExpr)
}

func extractResolvedTypeAssertion(typeName string, valExpr string) string {
	if typeName != "" {
		return fmt.Sprintf("%s.Interface().(%s)", valExpr, typeName)
	}
	return fmt.Sprintf("%s.Interface()", valExpr)
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
