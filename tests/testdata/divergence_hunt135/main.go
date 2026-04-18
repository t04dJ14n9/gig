package divergence_hunt135

import "fmt"

// ============================================================================
// Round 135: Const, iota, and type safety
// ============================================================================

const (
	A = iota
	B
	C
)

func IotaBasic() string {
	return fmt.Sprintf("%d-%d-%d", A, B, C)
}

const (
	Sunday    = 0
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
)

func ConstExplicit() string {
	return fmt.Sprintf("%d-%d", Sunday, Wednesday)
}

const (
	KB = 1 << (10 * iota)
	MB
	GB
)

func IotaExpression() string {
	return fmt.Sprintf("KB=%d-MB=%d-GB=%d", KB, MB, GB)
}

func ConstUntyped() string {
	const x = 10
	var i int = x
	var f float64 = x
	return fmt.Sprintf("i=%d-f=%.1f", i, f)
}

type Status int

const (
	StatusOK    Status = 0
	StatusError Status = 1
	StatusPanic Status = 2
)

func ConstTyped() string {
	s := StatusOK
	return fmt.Sprintf("status=%d", s)
}

func ConstString() string {
	const greeting = "hello"
	return greeting
}

func ConstBool() string {
	const flag = true
	if flag {
		return "true"
	}
	return "false"
}

func ConstExpression() string {
	const (
		a = 10
		b = a * 2
		c = b + a
	)
	return fmt.Sprintf("%d-%d-%d", a, b, c)
}

const (
	_  = iota
	KB2 = 1 << (10 * iota)
	MB2
)

func IotaSkip() string {
	return fmt.Sprintf("KB2=%d-MB2=%d", KB2, MB2)
}
