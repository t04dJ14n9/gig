// type_assert.go implements ssa.TypeAssert and ssa.Panic.
package interp

import (
	"fmt"
	"reflect"

	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/value"
)

// runPanic implements `panic(x)` as an SSA instruction (distinct from
// the panic builtin used as a function value).
func (p *program) runPanic(fr *frame, instr *ssa.Panic) (continuation, []value.Value, error) {
	v, err := p.readValue(fr, instr.X)
	if err != nil {
		return contNext, nil, err
	}
	panic(v.Interface())
}

// runTypeAssert implements `x.(T)` and `x, ok := x.(T)`. The interpreter
// uses reflect.Type.AssignableTo for the type check; on success it
// returns the asserted value (or a tuple in the comma-ok form), on
// failure either panics or returns the zero/false tuple.
func (p *program) runTypeAssert(fr *frame, instr *ssa.TypeAssert) (continuation, []value.Value, error) {
	x, err := p.readValue(fr, instr.X)
	if err != nil {
		return contNext, nil, err
	}
	rv, err := p.reflectOf(x, nil)
	if err != nil {
		return contNext, nil, err
	}
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	dst, err := p.resolver.ResolveType(instr.AssertedType)
	if err != nil {
		return contNext, nil, err
	}

	assignable := false
	if rv.IsValid() {
		// Type assertion compares the dynamic type to the asserted
		// type. Go's runtime only allows the exact type or types
		// implementing the asserted interface — not generic
		// convertibility (Go would never let `(any(42.0)).(int)`
		// succeed even though float64 is ConvertibleTo int).
		if dst.Kind() == reflect.Interface {
			assignable = rv.Type().AssignableTo(dst) || rv.Type().Implements(dst)
		} else {
			assignable = rv.Type() == dst || rv.Type().AssignableTo(dst)
		}
	}

	switch {
	case instr.CommaOk:
		// Build a synthetic (T, bool) tuple.
		tt, ok := instr.Type().(*types.Tuple)
		if !ok {
			return contNext, nil, fmt.Errorf("interp: TypeAssert CommaOk type not a tuple: %s", instr.Type())
		}
		holderType, err := p.resolver.ResolveType(tt)
		if err != nil {
			return contNext, nil, err
		}
		holder := reflect.New(holderType).Elem()
		if assignable {
			converted := rv
			if rv.Type() != dst {
				converted = rv.Convert(dst)
			}
			holder.Field(0).Set(converted)
			holder.Field(1).SetBool(true)
		} else {
			holder.Field(1).SetBool(false)
		}
		fr.setCell(instr, reflectValue(holder))
	default:
		if !assignable {
			if !rv.IsValid() {
				panic(fmt.Errorf("interface conversion: interface is nil, not %s", dst))
			}
			panic(fmt.Errorf("interface conversion: %s is not %s", rv.Type(), dst))
		}
		converted := rv
		if rv.Type() != dst {
			converted = rv.Convert(dst)
		}
		out, err := p.converter.FromReflect(converted)
		if err != nil {
			return contNext, nil, err
		}
		fr.setCell(instr, out)
	}
	return contNext, nil, nil
}
