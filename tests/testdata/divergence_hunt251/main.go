package divergence_hunt251

import (
	"fmt"
)

// ============================================================================
// Round 251: Const declarations with iota
// ============================================================================

const (
	A = iota
	B
	C
)

const (
	D = iota
	E
	F
)

const (
	Flag1 = 1 << iota
	Flag2
	Flag4
	Flag8
)

const (
	_ = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

const (
	X = iota * 10
	Y
	Z
)

// IotaBasic tests basic iota enumeration
func IotaBasic() string {
	return fmt.Sprintf("A=%d,B=%d,C=%d", A, B, C)
}

// IotaReset tests iota resetting in new const block
func IotaReset() string {
	return fmt.Sprintf("D=%d,E=%d,F=%d", D, E, F)
}

// IotaBitshift tests iota with bit shift
func IotaBitshift() string {
	return fmt.Sprintf("Flag1=%d,Flag2=%d,Flag4=%d,Flag8=%d", Flag1, Flag2, Flag4, Flag8)
}

// IotaSkip tests skipping first value with iota
func IotaSkip() string {
	return fmt.Sprintf("KB=%d,MB=%d,GB=%d", KB, MB, GB)
}

// IotaExpression tests iota in expressions
func IotaExpression() string {
	return fmt.Sprintf("X=%d,Y=%d,Z=%d", X, Y, Z)
}

// IotaWithOffset tests iota with offset calculation
func IotaWithOffset() string {
	const (
		Base = 100
		One  = Base + iota
		Two
		Three
	)
	return fmt.Sprintf("One=%d,Two=%d,Three=%d", One, Two, Three)
}

// IotaStringEnum tests iota for string-like behavior
func IotaStringEnum() string {
	const (
		Low = iota
		Medium
		High
	)
	return fmt.Sprintf("Low=%d,Medium=%d,High=%d", Low, Medium, High)
}

// IotaDaysOfWeek tests iota for days enumeration
func IotaDaysOfWeek() string {
	const (
		Sunday = iota
		Monday
		Tuesday
		Wednesday
		Thursday
		Friday
		Saturday
	)
	return fmt.Sprintf("Sun=%d,Mon=%d,Tue=%d,Wed=%d", Sunday, Monday, Tuesday, Wednesday)
}

// IotaWithParentheses tests iota with parentheses
func IotaWithParentheses() string {
	const (
		P = (iota * 2)
		Q
		R
	)
	return fmt.Sprintf("P=%d,Q=%d,R=%d", P, Q, R)
}

// IotaNegativeStep tests negative iota progression
func IotaNegativeStep() string {
	const (
		Max = 10 - iota
		Nine
		Eight
	)
	return fmt.Sprintf("Max=%d,Nine=%d,Eight=%d", Max, Nine, Eight)
}

// IotaMultiplication tests iota with multiplication
func IotaMultiplication() string {
	const (
		Step5 = (iota + 1) * 5
		Step10
		Step15
	)
	return fmt.Sprintf("Step5=%d,Step10=%d,Step15=%d", Step5, Step10, Step15)
}
