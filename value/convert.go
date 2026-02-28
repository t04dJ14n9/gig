package value

import "fmt"

// --- Type Conversions ---

// ToInt converts to int.
func (v Value) ToInt() Value {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		if v.Bool() {
			return MakeInt(1)
		}
		return MakeInt(0)
	case KindInt:
		return v
	case KindUint:
		return MakeInt(int64(v.Uint()))
	case KindFloat:
		return MakeInt(int64(v.Float()))
	default:
		panic(fmt.Sprintf("cannot convert %v to int", v.kind))
	}
}

// ToUint converts to uint.
func (v Value) ToUint() Value {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		if v.Bool() {
			return MakeUint(1)
		}
		return MakeUint(0)
	case KindInt:
		return MakeUint(uint64(v.num))
	case KindUint:
		return v
	case KindFloat:
		return MakeUint(uint64(v.Float()))
	default:
		panic(fmt.Sprintf("cannot convert %v to uint", v.kind))
	}
}

// ToFloat converts to float.
func (v Value) ToFloat() Value {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		if v.Bool() {
			return MakeFloat(1.0)
		}
		return MakeFloat(0.0)
	case KindInt:
		return MakeFloat(float64(v.num))
	case KindUint:
		return MakeFloat(float64(v.Uint()))
	case KindFloat:
		return v
	default:
		panic(fmt.Sprintf("cannot convert %v to float", v.kind))
	}
}

// ToBool converts to bool.
func (v Value) ToBool() Value {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		return v
	case KindInt:
		return MakeBool(v.num != 0)
	case KindUint:
		return MakeBool(v.num != 0)
	case KindFloat:
		return MakeBool(v.Float() != 0)
	case KindString:
		return MakeBool(v.obj.(string) != "")
	default:
		return MakeBool(!v.IsNil())
	}
}

// ToString converts to string representation.
func (v Value) ToString() Value {
	return MakeString(fmt.Sprintf("%v", v.Interface()))
}
