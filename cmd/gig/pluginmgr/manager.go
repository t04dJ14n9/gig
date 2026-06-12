// Package pluginmgr provides hot-loading of external Go packages via the
// Go plugin system. This allows users to import third-party libraries
// directly in the REPL without manual code generation.
//
// # Platform Support
//
// Plugin loading is only supported on Linux and macOS.
// Windows users will receive an error message suggesting alternatives.
//
// # Architecture
//
// The plugin system works by:
//  1. Detecting new import statements in REPL
//  2. Downloading packages with `go get`
//  3. Generating wrapper code (similar to gentool)
//  4. Compiling as a `.so` shared library
//  5. Loading with `plugin.Open()`
//  6. Calling the registration function to register with gig
//
// # Cache Directory
//
// Plugins are cached in ~/.gig/plugins/ with the following structure:
//
//	~/.gig/plugins/
//	├── github.com/
//	│   └── spf13/
//	│       └── cast/
//	│           ├── cast.so      # Compiled plugin
//	│           └── wrapper.go   # Generated wrapper code
//	└── go.mod                   # Plugin module file
package pluginmgr

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrPluginNotSupported = errors.New("plugin loading requires Linux or macOS; Windows is not supported")
	ErrPackageNotFound    = errors.New("package not found; run 'go get <package>' first")
	ErrPluginBuildFailed  = errors.New("failed to build plugin")
)

// Manager manages hot-loading of external packages.
type Manager struct {
	mu        sync.RWMutex
	pluginDir string                      // ~/.gig/plugins
	loaded    map[string]bool             // Set of loaded package paths
	registry  map[string]pluginMetadata   // Package path -> metadata
	symbols   map[string]*ExportedSymbols // Package path -> cached symbols
}

// pluginMetadata stores information about a loaded plugin.
type pluginMetadata struct {
	Path      string `json:"path"`       // Package import path
	SOPath    string `json:"so_path"`    // Path to .so file
	BuildTime string `json:"build_time"` // When it was built
}

// NewManager creates a new plugin manager.
func NewManager() *Manager {
	pluginDir := getPluginDir()
	pm := &Manager{
		pluginDir: pluginDir,
		loaded:    make(map[string]bool),
		registry:  make(map[string]pluginMetadata),
		symbols:   make(map[string]*ExportedSymbols),
	}
	pm.loadRegistry()
	return pm
}

// getPluginDir returns the plugin cache directory.
func getPluginDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir()
	}
	return filepath.Join(homeDir, ".gig", "plugins")
}
