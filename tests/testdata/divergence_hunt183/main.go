package divergence_hunt183

import (
	"fmt"
	"unicode"
)

// ============================================================================
// Round 183: Unicode category checks
// ============================================================================

func IsDigit() string {
	c1 := unicode.IsDigit('5')
	c2 := unicode.IsDigit('a')
	c3 := unicode.IsDigit('٥') // Arabic-Indic digit
	return fmt.Sprintf("%v:%v:%v", c1, c2, c3)
}

func IsLetter() string {
	c1 := unicode.IsLetter('A')
	c2 := unicode.IsLetter('z')
	c3 := unicode.IsLetter('9')
	return fmt.Sprintf("%v:%v:%v", c1, c2, c3)
}

func IsSpace() string {
	c1 := unicode.IsSpace(' ')
	c2 := unicode.IsSpace('\t')
	c3 := unicode.IsSpace('\n')
	c4 := unicode.IsSpace('x')
	return fmt.Sprintf("%v:%v:%v:%v", c1, c2, c3, c4)
}

func IsUpper() string {
	c1 := unicode.IsUpper('A')
	c2 := unicode.IsUpper('a')
	c3 := unicode.IsUpper('5')
	return fmt.Sprintf("%v:%v:%v", c1, c2, c3)
}

func IsLower() string {
	c1 := unicode.IsLower('a')
	c2 := unicode.IsLower('A')
	c3 := unicode.IsLower('5')
	return fmt.Sprintf("%v:%v:%v", c1, c2, c3)
}

func IsNumber() string {
	c1 := unicode.IsNumber('5')
	c2 := unicode.IsNumber('²') // Superscript 2
	c3 := unicode.IsNumber('a')
	return fmt.Sprintf("%v:%v:%v", c1, c2, c3)
}

func IsPunct() string {
	c1 := unicode.IsPunct('.')
	c2 := unicode.IsPunct('!')
	c3 := unicode.IsPunct('a')
	return fmt.Sprintf("%v:%v:%v", c1, c2, c3)
}

func ToUpper() string {
	r1 := unicode.ToUpper('a')
	r2 := unicode.ToUpper('A')
	r3 := unicode.ToUpper('ä') // a with diaeresis
	return fmt.Sprintf("%c:%c:%c", r1, r2, r3)
}

func ToLower() string {
	r1 := unicode.ToLower('A')
	r2 := unicode.ToLower('a')
	r3 := unicode.ToLower('Ä') // A with diaeresis
	return fmt.Sprintf("%c:%c:%c", r1, r2, r3)
}

func SimpleFold() string {
	r := unicode.SimpleFold('A')
	return fmt.Sprintf("%c", r)
}

func IsGraphic() string {
	c1 := unicode.IsGraphic('A')
	c2 := unicode.IsGraphic(' ')
	c3 := unicode.IsGraphic('\n')
	return fmt.Sprintf("%v:%v:%v", c1, c2, c3)
}

func IsControl() string {
	c1 := unicode.IsControl('\n')
	c2 := unicode.IsControl('\t')
	c3 := unicode.IsControl('A')
	return fmt.Sprintf("%v:%v:%v", c1, c2, c3)
}
