package repl

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// handleCommand processes REPL commands.
func (s *Session) handleCommand(input string) {
	cmd := strings.TrimSpace(strings.TrimPrefix(input, ":"))
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "help", "h", "?":
		s.showHelp()
	case "quit", "q", "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "clear":
		s.clear()
		fmt.Println("Session cleared.")
	case "imports":
		s.showImports()
	case "vars":
		s.showVars()
	case "source":
		s.showSource()
	case "timeout":
		s.setTimeout(parts)
	case "plugins":
		s.showPlugins()
	case "env":
		s.showEnv()
	default:
		fmt.Printf("Unknown command: :%s\n", parts[0])
		fmt.Println("Type :help for available commands.")
	}
}

// showHelp displays available commands.
func (s *Session) showHelp() {
	fmt.Print(`
Commands:
  :help, :h, :?    Show this help message
  :quit, :q        Exit the REPL
  :clear           Clear the session (reset all state)
  :imports         Show imported packages
  :vars            Show captured variables
  :env             Show all variables and packages
  :source          Show accumulated declarations
  :timeout <dur>   Set execution timeout (e.g., :timeout 5s)
  :plugins         Show loaded plugins

Input Types:
  Import       -> import "fmt"
  Expression   -> 1+1, x + y
  Statement    -> x := 1, for i:=0; i<5; i++ {}
  Declaration  -> func foo() {}, type Point struct{}

Examples:
  >>> import "fmt"
  >>> x := 5
  >>> y := x + 3
  >>> x * y
  40
  >>> :vars
  Variables:
    x int = 5
    y int = 8
`)
}

// clear resets the session state.
func (s *Session) clear() {
	s.declarations.Reset()
	s.imports = make(map[string]bool)
	s.vars = make(map[string]varInfo)
	s.exprCounter = 0
}

// showImports displays imported packages.
func (s *Session) showImports() {
	if len(s.imports) == 0 {
		fmt.Println("No packages imported.")
		return
	}
	fmt.Println("Imported packages:")
	for pkg := range s.imports {
		fmt.Printf("  %s\n", pkg)
	}
}

// showVars displays captured variables.
func (s *Session) showVars() {
	if len(s.vars) == 0 {
		fmt.Println("No variables defined.")
		return
	}
	fmt.Println("Variables:")
	for name, info := range s.vars {
		fmt.Printf("  %s %s = %v\n", name, info.typeName, info.value)
	}
}

// showSource displays accumulated source code.
func (s *Session) showSource() {
	src := s.declarations.String()
	if src == "" {
		fmt.Println("No declarations in session.")
		return
	}
	fmt.Println("Accumulated declarations:")
	fmt.Println(src)
}

// setTimeout sets the execution timeout.
func (s *Session) setTimeout(parts []string) {
	if len(parts) < 2 {
		fmt.Printf("Current timeout: %v\n", s.timeout)
		return
	}
	dur, err := time.ParseDuration(parts[1])
	if err != nil {
		fmt.Printf("Invalid duration: %s\n", parts[1])
		return
	}
	s.timeout = dur
	fmt.Printf("Timeout set to: %v\n", dur)
}

// showPlugins displays loaded plugins.
func (s *Session) showPlugins() {
	if s.pluginMgr == nil {
		fmt.Println("No plugins loaded.")
		return
	}
	plugins := s.pluginMgr.ListLoaded()
	if len(plugins) == 0 {
		fmt.Println("No plugins loaded.")
		return
	}
	fmt.Println("Loaded plugins:")
	for _, p := range plugins {
		fmt.Printf("  %s\n", p)
	}
}

// showEnv displays all variables and packages in the current session.
func (s *Session) showEnv() {
	fmt.Println("Packages:")
	if len(s.imports) == 0 {
		fmt.Println("  (none)")
	} else {
		pkgs := make([]string, 0, len(s.imports))
		for pkg := range s.imports {
			pkgs = append(pkgs, pkg)
		}
		sort.Strings(pkgs)
		for _, pkg := range pkgs {
			pkgName := pkg
			if idx := strings.LastIndex(pkg, "/"); idx >= 0 {
				pkgName = pkg[idx+1:]
			}
			if name, ok := pkgNameFromPath(pkg); ok {
				pkgName = name
			}
			fmt.Printf("  %s (%s)\n", pkgName, pkg)
		}
	}

	fmt.Println()
	fmt.Println("Variables:")
	if len(s.vars) == 0 {
		fmt.Println("  (none)")
		return
	}

	names := make([]string, 0, len(s.vars))
	for name := range s.vars {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		info := s.vars[name]
		fmt.Printf("  %s %s = %v\n", name, info.typeName, info.value)
	}
}
