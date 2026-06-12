package benchmarks

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
)

// ============================================================================
// External Calls
// ============================================================================

func ExternalSprintf() int {
	s := ""
	for i := 0; i < 100; i++ {
		s = fmt.Sprintf("%d", i)
	}
	return len(s)
}

func ExternalStrings() int {
	s := ""
	for i := 0; i < 100; i++ {
		s = strings.ToUpper("hello world test string")
	}
	return len(s)
}

func SortInts() int {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = 100 - i
	}
	sort.Ints(s)
	return s[0] + s[99]
}

func StringsBuilder() int {
	var sb strings.Builder
	for j := 0; j < 100; j++ {
		sb.WriteString("hello")
		sb.WriteString("world")
	}
	return sb.Len()
}

func MathBig() int {
	a := big.NewInt(1)
	bv := big.NewInt(1)
	for j := 0; j < 100; j++ {
		a.Add(a, bv)
		bv.Sub(a, bv)
	}
	return int(a.Int64() % 1000)
}

type Data struct {
	Name string
	Age  int
	City string
}

func JsonMarshal() int {
	d := Data{Name: "John", Age: 30, City: "NYC"}
	s, _ := json.Marshal(d)
	return len(s)
}
