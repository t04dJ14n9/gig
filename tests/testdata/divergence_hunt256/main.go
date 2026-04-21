package divergence_hunt256

import (
	"fmt"
)

// ============================================================================
// Round 256: Short variable declaration
// ============================================================================

// ShortDeclBasic tests basic short declaration
func ShortDeclBasic() string {
	x := 42
	return fmt.Sprintf("x=%d", x)
}

// ShortDeclMultiple tests multiple variable short declaration
func ShortDeclMultiple() string {
	x, y, z := 1, 2, 3
	return fmt.Sprintf("x=%d,y=%d,z=%d", x, y, z)
}

// ShortDeclMixedTypes tests short declaration with mixed types
func ShortDeclMixedTypes() string {
	n, s, b := 42, "hello", true
	return fmt.Sprintf("n=%d,s=%s,b=%v", n, s, b)
}

// ShortDeclRedefine tests short declaration redefining variables
func ShortDeclRedefine() string {
	x, y := 1, 2
	x, z := 3, 4
	return fmt.Sprintf("x=%d,y=%d,z=%d", x, y, z)
}

// ShortDeclInIf tests short declaration in if statement
func ShortDeclInIf() string {
	if x := 10; x > 5 {
		return fmt.Sprintf("x=%d", x)
	}
	return ""
}

// ShortDeclInFor tests short declaration in for loop
func ShortDeclInFor() string {
	result := ""
	for i := 0; i < 3; i++ {
		result += fmt.Sprintf("%d,", i)
	}
	return result
}

// ShortDeclInSwitch tests short declaration in switch
func ShortDeclInSwitch() string {
	switch x := 5; x {
	case 5:
		return fmt.Sprintf("x=%d", x)
	default:
		return "other"
	}
}

// ShortDeclSlice tests short declaration with slice
func ShortDeclSlice() string {
	s := []int{1, 2, 3}
	a, b := s[0], s[1]
	return fmt.Sprintf("a=%d,b=%d", a, b)
}

// ShortDeclFunctionCall tests short declaration from function call
func ShortDeclFunctionCall() string {
	min, max := getMinMax()
	return fmt.Sprintf("min=%d,max=%d", min, max)
}

func getMinMax() (int, int) {
	return 10, 100
}

// ShortDeclTypeInference tests type inference in short declaration
func ShortDeclTypeInference() string {
	i := 42           // int
	f := 3.14         // float64
	s := "hello"      // string
	b := true         // bool
	r := 'A'          // rune (int32)
	return fmt.Sprintf("%T,%T,%T,%T,%T", i, f, s, b, r)
}

// ShortDeclBlank tests short declaration with blank identifier
func ShortDeclBlank() string {
	x, _ := 1, 2
	_, y := 3, 4
	return fmt.Sprintf("x=%d,y=%d", x, y)
}

// ShortDeclComposite tests short declaration with composite literal
func ShortDeclComposite() string {
	m := map[string]int{"a": 1, "b": 2}
	s := struct{ X, Y int }{10, 20}
	return fmt.Sprintf("m[a]=%d,s.X=%d", m["a"], s.X)
}
