package divergence_hunt158

import "fmt"

// ============================================================================
// Round 158: Constants and iota patterns
// ============================================================================

// IotaBasic tests basic iota usage
const (
	Zero = iota
	One
	Two
	Three
)

func IotaBasic() string {
	return fmt.Sprintf("zero=%d-one=%d-two=%d-three=%d", Zero, One, Two, Three)
}

// IotaWithValue tests iota with initial value
const (
	First = iota + 1
	Second
	Third
)

func IotaWithValue() string {
	return fmt.Sprintf("first=%d-second=%d-third=%d", First, Second, Third)
}

// IotaWithSkip tests iota with skipped value
const (
	A = iota
	B = iota
	_ = iota // skip
	D = iota
)

func IotaWithSkip() string {
	return fmt.Sprintf("a=%d-b=%d-d=%d", A, B, D)
}

// IotaBitShift tests iota with bit shifting
const (
	FlagNone = 1 << iota
	FlagRead
	FlagWrite
	FlagExecute
)

func IotaBitShift() string {
	return fmt.Sprintf("none=%d-read=%d-write=%d-execute=%d", FlagNone, FlagRead, FlagWrite, FlagExecute)
}

// IotaExpression tests iota in expression
const (
	KB = 1 << (10 * iota)
	MB
	GB
)

func IotaExpression() string {
	return fmt.Sprintf("kb=%d-mb=%d-gb=%d", KB, MB, GB)
}

// IotaMultiplePerLine tests multiple iota per line
const (
	Low, Medium, High = iota, iota + 10, iota + 100
)

func IotaMultiplePerLine() string {
	return fmt.Sprintf("low=%d-medium=%d-high=%d", Low, Medium, High)
}

// IotaReset tests iota resetting in new const block
const (
	X = iota
	Y
)

const (
	A2 = iota
	B2
)

func IotaReset() string {
	return fmt.Sprintf("x=%d-y=%d-a2=%d-b2=%d", X, Y, A2, B2)
}

// IotaStringer tests iota with String() method pattern
type Status int

const (
	Pending Status = iota
	Running
	Completed
	Failed
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "pending"
	case Running:
		return "running"
	case Completed:
		return "completed"
	case Failed:
		return "failed"
	default:
		return "unknown"
	}
}

func IotaStringer() string {
	return fmt.Sprintf("pending=%s-running=%s", Pending.String(), Running.String())
}

// UntypedConstant tests untyped constant arithmetic
const (
	BigNum   = 1 << 100
	Bigger   = BigNum >> 99
	Fraction = 3.14
)

func UntypedConstant() string {
	return fmt.Sprintf("bigger=%d-fraction=%.2f", Bigger, Fraction)
}

// ConstWithType tests typed vs untyped constants
const (
	UntypedInt    = 42
	TypedInt  int = 42
)

func ConstWithType() string {
	return fmt.Sprintf("untyped=%d-typed=%d", UntypedInt, TypedInt)
}
