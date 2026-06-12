package divergence_hunt240

import "fmt"

// ============================================================================
// Round 240: GoStringer interface
// ============================================================================

type Config struct {
	Host string
	Port int
}

func (c Config) GoString() string {
	return fmt.Sprintf("Config{Host: %q, Port: %d}", c.Host, c.Port)
}

type Coordinate struct {
	Lat, Long float64
}

func (c Coordinate) GoString() string {
	return fmt.Sprintf("Coordinate{Lat: %.4f, Long: %.4f}", c.Lat, c.Long)
}

type Color struct {
	R, G, B uint8
}

func (c Color) GoString() string {
	return fmt.Sprintf("Color{R: %d, G: %d, B: %d}", c.R, c.G, c.B)
}

func GoStringerBasic() string {
	c := Config{Host: "localhost", Port: 8080}
	return fmt.Sprintf("%#v", c)
}

func GoStringerStruct() string {
	coord := Coordinate{Lat: 35.6762, Long: 139.6503}
	return fmt.Sprintf("%#v", coord)
}

func GoStringerWithColor() string {
	c := Color{R: 255, G: 128, B: 0}
	return fmt.Sprintf("%#v", c)
}

func GoStringerVsStringer() string {
	type Item struct {
		Name  string
		Value int
	}
	i := Item{Name: "test", Value: 42}
	regular := fmt.Sprintf("%v", i)
	gosyntax := fmt.Sprintf("%#v", i)
	return fmt.Sprintf("regular=%s,gosyntax=%s", regular, gosyntax)
}

func GoStringerInSlice() string {
	configs := []Config{
		{Host: "localhost", Port: 8080},
		{Host: "example.com", Port: 443},
	}
	return fmt.Sprintf("%#v", configs)
}

func GoStringerNilCheck() string {
	var c *Config
	return fmt.Sprintf("%#v", c)
}

func GoStringerEmptyStruct() string {
	type Empty struct{}
	e := Empty{}
	return fmt.Sprintf("%#v", e)
}

func GoStringerNestedStruct() string {
	type Inner struct {
		Value int
	}
	type Outer struct {
		I     Inner
		Label string
	}
	o := Outer{I: Inner{Value: 42}, Label: "test"}
	return fmt.Sprintf("%#v", o)
}

func GoStringerWithPointer() string {
	type Node struct {
		Value int
	}
	n := &Node{Value: 100}
	return fmt.Sprintf("%#v", n)
}

func GoStringerFormatComparison() string {
	c := Config{Host: "test", Port: 1234}
	v := fmt.Sprintf("%v", c)
	plusV := fmt.Sprintf("%+v", c)
	sharpV := fmt.Sprintf("%#v", c)
	return fmt.Sprintf("v=%s,plus=%s,sharp=%s", v, plusV, sharpV)
}
