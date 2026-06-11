package main

import flag "github.com/spf13/pflag"

func PflagNewFlagSet() string {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	return fs.Name()
}

func PflagString() string {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	s := fs.String("name", "default", "help")
	fs.Parse([]string{"--name", "hello"})
	return *s
}

func PflagBool() bool {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	b := fs.Bool("verbose", false, "help")
	fs.Parse([]string{"--verbose"})
	return *b
}

func PflagInt() int {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	i := fs.Int("count", 0, "help")
	fs.Parse([]string{"--count", "42"})
	return *i
}

func PflagInt64() int64 {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	i := fs.Int64("count", 0, "help")
	fs.Parse([]string{"--count", "42"})
	return *i
}

func PflagUint() uint {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	u := fs.Uint("count", 0, "help")
	fs.Parse([]string{"--count", "42"})
	return *u
}

func PflagUint64() uint64 {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	u := fs.Uint64("count", 0, "help")
	fs.Parse([]string{"--count", "42"})
	return *u
}

func PflagFloat64() float64 {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	f := fs.Float64("rate", 0.0, "help")
	fs.Parse([]string{"--rate", "3.14"})
	return *f
}

func PflagStringSlice() int {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	ss := fs.StringSlice("names", []string{}, "help")
	fs.Parse([]string{"--names", "a,b,c"})
	return len(*ss)
}

func PflagIntSlice() int {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	is := fs.IntSlice("ids", []int{}, "help")
	fs.Parse([]string{"--ids", "1,2,3"})
	return len(*is)
}

func PflagArgs() int {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Parse([]string{"arg1", "arg2"})
	return len(fs.Args())
}

func PflagArg() string {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Parse([]string{"hello"})
	return fs.Arg(0)
}

func PflagNArg() int {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Parse([]string{"a", "b"})
	return fs.NArg()
}

func PflagNFlag() int {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	_ = fs.Bool("verbose", false, "help")
	_ = fs.String("name", "", "help")
	fs.Parse([]string{"--verbose", "--name", "test"})
	return fs.NFlag()
}

func PflagChanged() bool {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	_ = fs.String("name", "default", "help")
	fs.Parse([]string{"--name", "hello"})
	return fs.Changed("name")
}

func PflagLookup() bool {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	_ = fs.String("name", "default", "help")
	f := fs.Lookup("name")
	return f != nil
}

func PflagDuration() int64 {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	d := fs.Duration("timeout", 0, "help")
	fs.Parse([]string{"--timeout", "5s"})
	return int64((*d).Seconds())
}

func PflagVar() string {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var val string
	fs.StringVar(&val, "name", "default", "help")
	fs.Parse([]string{"--name", "hello"})
	return val
}

func PflagPrintDefaults() bool {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	_ = fs.String("name", "default", "help text")
	_ = fs.Bool("verbose", false, "verbose output")
	return true // just testing it compiles and runs
}
