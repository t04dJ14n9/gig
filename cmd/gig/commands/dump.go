// Package commands implements the CLI subcommands for the gig tool.
package commands

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

// RunDump implements the "gig dump" subcommand.
func RunDump(fs *flag.FlagSet, args []string) error {
	allowPanic := fs.Bool("allow-panic", false, "allow panic/recover/defer compilation while dumping")
	allowUnsafeTypePass := fs.Bool("unsafe-type-pass", false, "disable third-party external type boundary checks while dumping")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gig dump [flags] <file|->\n\n")
		fmt.Fprintf(os.Stderr, "Compiles Gig source and prints readable SSA plus bytecode.\n")
		fmt.Fprintf(os.Stderr, "Use '-' to read source from stdin.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("source file argument required")
	}

	source, err := readDumpSource(fs.Arg(0))
	if err != nil {
		return err
	}

	var opts []gig.BuildOption
	if *allowPanic {
		opts = append(opts, gig.WithAllowPanic())
	}
	if *allowUnsafeTypePass {
		opts = append(opts, gig.WithAllowUnsafeTypePass())
	}
	dump, err := gig.DebugDump(string(source), opts...)
	if err != nil {
		return err
	}
	fmt.Print(dump)
	return nil
}

func readDumpSource(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	return source, nil
}
