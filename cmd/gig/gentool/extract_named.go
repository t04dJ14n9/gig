package gentool

import (
	"fmt"
	"go/types"
)

// extractNamedArg handles argument extraction for named types.
func extractNamedArg(named *types.Named, underlying types.Type, valExpr string, pkgRef string) string {
	obj := named.Obj()
	pkg := obj.Pkg()

	if pkg != nil && pkg.Path() == currentPkgPath {
		if bt, ok := underlying.(*types.Basic); ok {
			basicExpr := extractBasic(bt, valExpr)
			if basicExpr == "" {
				return ""
			}
			namedName := resolveTypeName(named, pkgRef)
			if namedName == "" {
				return ""
			}
			return fmt.Sprintf("%s(%s)", namedName, basicExpr)
		}
		namedName := resolveTypeName(named, pkgRef)
		if namedName != "" {
			return fmt.Sprintf("%s.Interface().(%s)", valExpr, namedName)
		}
		return extractUnderlyingWithPkgRef(underlying, valExpr, pkgRef)
	}

	if bt, ok := underlying.(*types.Basic); ok {
		basicExpr := extractBasic(bt, valExpr)
		if basicExpr == "" {
			return ""
		}
		qualifiedName := resolveQualifiedName(named, pkgRef)
		if qualifiedName != "" {
			return fmt.Sprintf("%s(%s)", qualifiedName, basicExpr)
		}
		return ""
	}

	qualifiedName := resolveQualifiedName(named, pkgRef)
	if qualifiedName != "" {
		return fmt.Sprintf("%s.Interface().(%s)", valExpr, qualifiedName)
	}
	return ""
}

// extractAliasArg handles argument extraction for alias types.
func extractAliasArg(alias *types.Alias, underlying types.Type, valExpr string, pkgRef string) string {
	obj := alias.Obj()
	pkg := obj.Pkg()

	if pkg == nil {
		return extractUnderlyingWithPkgRef(underlying, valExpr, pkgRef)
	}

	if pkg.Path() == currentPkgPath {
		if bt, ok := underlying.(*types.Basic); ok {
			basicExpr := extractBasic(bt, valExpr)
			if basicExpr == "" {
				return ""
			}
			aliasName := resolveTypeName(alias, pkgRef)
			if aliasName == "" {
				return ""
			}
			return fmt.Sprintf("%s(%s)", aliasName, basicExpr)
		}
		aliasName := resolveTypeName(alias, pkgRef)
		if aliasName != "" {
			return fmt.Sprintf("%s.Interface().(%s)", valExpr, aliasName)
		}
		return extractUnderlyingWithPkgRef(underlying, valExpr, pkgRef)
	}

	if bt, ok := underlying.(*types.Basic); ok {
		basicExpr := extractBasic(bt, valExpr)
		if basicExpr == "" {
			return ""
		}
		return fmt.Sprintf("%s.%s(%s)", sanitizePkgName(pkg.Path()), obj.Name(), basicExpr)
	}
	return fmt.Sprintf("%s.Interface().(%s.%s)", valExpr, sanitizePkgName(pkg.Path()), obj.Name())
}

// resolveQualifiedName returns the fully-qualified Go type name for a named type,
// suitable for use in a type assertion like args[i].Interface().(pkg.Type).
func resolveQualifiedName(named *types.Named, pkgRef string) string {
	obj := named.Obj()
	pkg := obj.Pkg()
	if pkg == nil {
		return obj.Name()
	}
	if pkg.Path() == currentPkgPath {
		return fmt.Sprintf("%s.%s", pkgRef, obj.Name())
	}
	return fmt.Sprintf("%s.%s", sanitizePkgName(pkg.Path()), obj.Name())
}
