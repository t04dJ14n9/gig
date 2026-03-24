package main

import (
	"flag"
	"fmt"
	"os"
)

// runREPL implements the "gig repl" subcommand.
func runREPL(fs *flag.FlagSet, args []string) error {
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gig repl\n\n")
		fmt.Fprintf(os.Stderr, "Starts an interactive Go REPL (Read-Eval-Print Loop).\n")
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  >>> 1 + 2\n")
		fmt.Fprintf(os.Stderr, "  3\n")
		fmt.Fprintf(os.Stderr, "  >>> fmt.Sprintf(\"Hello, %%s!\", \"World\")\n")
		fmt.Fprintf(os.Stderr, "  \"Hello, World!\"\n")
		fmt.Fprintf(os.Stderr, "  >>> :help\n")
		fmt.Fprintf(os.Stderr, "  (shows available commands)\n")
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	session := NewSession()
	session.Run()
	return nil
}
