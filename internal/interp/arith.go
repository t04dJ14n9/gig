// arith.go contains the numeric/comparison operations for the SSA
// interpreter. evalBinOp and evalUnOp implement Go's arithmetic
// semantics for the scalar Kinds (int, uint, float, bool, string) and
// produce a Value retagged to the SSA result type.
package interp

import (
	"fmt"
	"go/token"
	"reflect"

	"go/types"

	"github.com/t04dJ14n9/gig/value"
)

// evalBinOp computes x op y and converts the result to result type t.
func evalBinOp(op token.Token, x, y value.Value, t types.Type, p *program) (value.Value, error) {
	if op == token.EQL || op == token.NEQ {
		return evalEquality(op, x, y), nil
	}
	// Named primitives (e.g. time.Duration) are stored as KindReflect
	// to keep their method set. Unbox them to scalar Kinds for the
	// duration of the op, then re-tag the result via Convert.
	x = unboxNamedPrimitive(x)
	y = unboxNamedPrimitive(y)
	// Shifts: Go allows int << uint, uint << uint, etc. The shift
	// count is always an unsigned integer; coerce y to uint64 so we
	// can dispatch on x's kind alone.
	if op == token.SHL || op == token.SHR {
		var shift uint64
		switch y.Kind() {
		case value.KindUint:
			shift = y.Uint()
		case value.KindInt:
			shift = uint64(y.Int())
		default:
			return value.Value{}, fmt.Errorf("interp: shift count is %s, expected integer", y.Kind())
		}
		switch x.Kind() {
		case value.KindInt:
			return shiftInt(op, x.Int(), shift, t, p)
		case value.KindUint:
			return shiftUint(op, x.Uint(), shift, t, p)
		}
		return value.Value{}, fmt.Errorf("interp: shift on %s not supported", x.Kind())
	}
	if x.Kind() != y.Kind() {
		return value.Value{},
			fmt.Errorf("interp: binop %s: kind mismatch (%s vs %s)", op, x.Kind(), y.Kind())
	}
	switch x.Kind() {
	case value.KindInt:
		return evalIntBinOp(op, x.Int(), y.Int(), t, p)
	case value.KindUint:
		return evalUintBinOp(op, x.Uint(), y.Uint(), t, p)
	case value.KindFloat:
		return evalFloatBinOp(op, x.Float(), y.Float(), t, p)
	case value.KindComplex:
		return evalComplexBinOp(op, x.Complex(), y.Complex(), t, p)
	case value.KindBool:
		return evalBoolBinOp(op, x.Bool(), y.Bool())
	case value.KindString:
		return evalStringBinOp(op, x.Str(), y.Str())
	}
	return value.Value{}, fmt.Errorf("interp: binop %s on %s not supported", op, x.Kind())
}

// unboxNamedPrimitive converts a KindReflect value whose underlying is
// a primitive scalar (e.g. time.Duration -> int64) into the matching
// scalar Kind. Returns the original value unchanged otherwise.
func unboxNamedPrimitive(v value.Value) value.Value {
	if v.Kind() != value.KindReflect {
		return v
	}
	rv, ok := v.Reflect()
	if !ok || !rv.IsValid() {
		return v
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.MakeInt(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.MakeUint(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return value.MakeFloat(rv.Float())
	case reflect.Bool:
		return value.MakeBool(rv.Bool())
	case reflect.String:
		return value.MakeString(rv.String())
	}
	return v
}

func shiftInt(op token.Token, x int64, n uint64, t types.Type, p *program) (value.Value, error) {
	var r int64
	switch op {
	case token.SHL:
		r = x << n
	case token.SHR:
		r = x >> n
	default:
		return value.Value{}, fmt.Errorf("interp: not a shift op: %s", op)
	}
	return p.converter.Convert(value.MakeInt(r), t, p.resolver)
}

func shiftUint(op token.Token, x, n uint64, t types.Type, p *program) (value.Value, error) {
	var r uint64
	switch op {
	case token.SHL:
		r = x << n
	case token.SHR:
		r = x >> n
	default:
		return value.Value{}, fmt.Errorf("interp: not a shift op: %s", op)
	}
	return p.converter.Convert(value.MakeUint(r), t, p.resolver)
}

func evalComplexBinOp(op token.Token, x, y complex128, t types.Type, p *program) (value.Value, error) {
	switch op {
	case token.ADD:
		r := x + y
		return p.converter.Convert(value.MakeComplex(real(r), imag(r)), t, p.resolver)
	case token.SUB:
		r := x - y
		return p.converter.Convert(value.MakeComplex(real(r), imag(r)), t, p.resolver)
	case token.MUL:
		r := x * y
		return p.converter.Convert(value.MakeComplex(real(r), imag(r)), t, p.resolver)
	case token.QUO:
		r := x / y
		return p.converter.Convert(value.MakeComplex(real(r), imag(r)), t, p.resolver)
	}
	return value.Value{}, fmt.Errorf("interp: complex binop %s not supported", op)
}

func evalIntBinOp(op token.Token, x, y int64, t types.Type, p *program) (value.Value, error) {
	switch op {
	case token.ADD:
		return p.converter.Convert(value.MakeInt(x+y), t, p.resolver)
	case token.SUB:
		return p.converter.Convert(value.MakeInt(x-y), t, p.resolver)
	case token.MUL:
		return p.converter.Convert(value.MakeInt(x*y), t, p.resolver)
	case token.QUO:
		if y == 0 {
			return value.Value{}, fmt.Errorf("interp: integer divide by zero")
		}
		return p.converter.Convert(value.MakeInt(x/y), t, p.resolver)
	case token.REM:
		if y == 0 {
			return value.Value{}, fmt.Errorf("interp: integer modulo by zero")
		}
		return p.converter.Convert(value.MakeInt(x%y), t, p.resolver)
	case token.AND:
		return p.converter.Convert(value.MakeInt(x&y), t, p.resolver)
	case token.OR:
		return p.converter.Convert(value.MakeInt(x|y), t, p.resolver)
	case token.XOR:
		return p.converter.Convert(value.MakeInt(x^y), t, p.resolver)
	case token.AND_NOT:
		return p.converter.Convert(value.MakeInt(x&^y), t, p.resolver)
	case token.SHL:
		return p.converter.Convert(value.MakeInt(x<<uint64(y)), t, p.resolver)
	case token.SHR:
		return p.converter.Convert(value.MakeInt(x>>uint64(y)), t, p.resolver)
	case token.LSS:
		return value.MakeBool(x < y), nil
	case token.LEQ:
		return value.MakeBool(x <= y), nil
	case token.GTR:
		return value.MakeBool(x > y), nil
	case token.GEQ:
		return value.MakeBool(x >= y), nil
	}
	return value.Value{}, fmt.Errorf("interp: int binop %s not supported", op)
}

func evalUintBinOp(op token.Token, x, y uint64, t types.Type, p *program) (value.Value, error) {
	switch op {
	case token.ADD:
		return p.converter.Convert(value.MakeUint(x+y), t, p.resolver)
	case token.SUB:
		return p.converter.Convert(value.MakeUint(x-y), t, p.resolver)
	case token.MUL:
		return p.converter.Convert(value.MakeUint(x*y), t, p.resolver)
	case token.QUO:
		if y == 0 {
			return value.Value{}, fmt.Errorf("interp: integer divide by zero")
		}
		return p.converter.Convert(value.MakeUint(x/y), t, p.resolver)
	case token.REM:
		if y == 0 {
			return value.Value{}, fmt.Errorf("interp: integer modulo by zero")
		}
		return p.converter.Convert(value.MakeUint(x%y), t, p.resolver)
	case token.AND:
		return p.converter.Convert(value.MakeUint(x&y), t, p.resolver)
	case token.OR:
		return p.converter.Convert(value.MakeUint(x|y), t, p.resolver)
	case token.XOR:
		return p.converter.Convert(value.MakeUint(x^y), t, p.resolver)
	case token.AND_NOT:
		return p.converter.Convert(value.MakeUint(x&^y), t, p.resolver)
	case token.SHL:
		return p.converter.Convert(value.MakeUint(x<<y), t, p.resolver)
	case token.SHR:
		return p.converter.Convert(value.MakeUint(x>>y), t, p.resolver)
	case token.LSS:
		return value.MakeBool(x < y), nil
	case token.LEQ:
		return value.MakeBool(x <= y), nil
	case token.GTR:
		return value.MakeBool(x > y), nil
	case token.GEQ:
		return value.MakeBool(x >= y), nil
	}
	return value.Value{}, fmt.Errorf("interp: uint binop %s not supported", op)
}

func evalFloatBinOp(op token.Token, x, y float64, t types.Type, p *program) (value.Value, error) {
	switch op {
	case token.ADD:
		return p.converter.Convert(value.MakeFloat(x+y), t, p.resolver)
	case token.SUB:
		return p.converter.Convert(value.MakeFloat(x-y), t, p.resolver)
	case token.MUL:
		return p.converter.Convert(value.MakeFloat(x*y), t, p.resolver)
	case token.QUO:
		return p.converter.Convert(value.MakeFloat(x/y), t, p.resolver)
	case token.LSS:
		return value.MakeBool(x < y), nil
	case token.LEQ:
		return value.MakeBool(x <= y), nil
	case token.GTR:
		return value.MakeBool(x > y), nil
	case token.GEQ:
		return value.MakeBool(x >= y), nil
	}
	return value.Value{}, fmt.Errorf("interp: float binop %s not supported", op)
}

func evalBoolBinOp(op token.Token, x, y bool) (value.Value, error) {
	switch op {
	case token.LAND:
		return value.MakeBool(x && y), nil
	case token.LOR:
		return value.MakeBool(x || y), nil
	}
	return value.Value{}, fmt.Errorf("interp: bool binop %s not supported", op)
}

func evalStringBinOp(op token.Token, x, y string) (value.Value, error) {
	switch op {
	case token.ADD:
		return value.MakeString(x + y), nil
	case token.LSS:
		return value.MakeBool(x < y), nil
	case token.LEQ:
		return value.MakeBool(x <= y), nil
	case token.GTR:
		return value.MakeBool(x > y), nil
	case token.GEQ:
		return value.MakeBool(x >= y), nil
	}
	return value.Value{}, fmt.Errorf("interp: string binop %s not supported", op)
}

// evalEquality covers EQL/NEQ across all comparable Kinds. Compares via
// Interface() round-trip so types that disagree (e.g. comparing a typed
// nil with KindNil) end up using Go's native == over any.
func evalEquality(op token.Token, x, y value.Value) value.Value {
	var equal bool
	xb, xIsBox := x.InterfaceBox()
	yb, yIsBox := y.InterfaceBox()
	switch {
	case xIsBox && yIsBox:
		// Both interface boxes: compare reflect.Value boxes via
		// Interface() so == respects (type, value) pair semantics.
		if !xb.IsValid() || !yb.IsValid() {
			equal = xb.IsValid() == yb.IsValid()
		} else if xb.IsNil() || yb.IsNil() {
			equal = xb.IsNil() && yb.IsNil()
		} else {
			equal = xb.Interface() == yb.Interface()
		}
	case xIsBox:
		// Comparing an interface box to a non-box (typically nil).
		if y.IsNil() {
			equal = xb.IsNil()
		} else {
			equal = xb.IsValid() && !xb.IsNil() && xb.Elem().Interface() == y.Interface()
		}
	case yIsBox:
		if x.IsNil() {
			equal = yb.IsNil()
		} else {
			equal = yb.IsValid() && !yb.IsNil() && yb.Elem().Interface() == x.Interface()
		}
	default:
		if x.IsNil() || y.IsNil() {
			equal = x.IsNil() == y.IsNil()
		} else {
			equal = x.Interface() == y.Interface()
		}
	}
	if op == token.NEQ {
		equal = !equal
	}
	return value.MakeBool(equal)
}

// evalUnOp computes (op x) and converts to t. MUL (deref) is handled
// in runUnOp and never reaches here.
func evalUnOp(op token.Token, x value.Value, t types.Type, p *program) (value.Value, error) {
	switch x.Kind() {
	case value.KindInt:
		switch op {
		case token.SUB:
			return p.converter.Convert(value.MakeInt(-x.Int()), t, p.resolver)
		case token.XOR:
			return p.converter.Convert(value.MakeInt(^x.Int()), t, p.resolver)
		}
	case value.KindUint:
		switch op {
		case token.SUB:
			return p.converter.Convert(value.MakeUint(-x.Uint()), t, p.resolver)
		case token.XOR:
			return p.converter.Convert(value.MakeUint(^x.Uint()), t, p.resolver)
		}
	case value.KindFloat:
		if op == token.SUB {
			return p.converter.Convert(value.MakeFloat(-x.Float()), t, p.resolver)
		}
	case value.KindBool:
		if op == token.NOT {
			return value.MakeBool(!x.Bool()), nil
		}
	}
	return value.Value{}, fmt.Errorf("interp: unop %s on %s not supported", op, x.Kind())
}
