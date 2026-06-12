package pluginmgr

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// buildPlugin compiles the wrapper as a .so shared library.
func (pm *Manager) buildPlugin(pkgPath, wrapperPath string) (string, error) {
	ctx := context.Background()
	pkgDir := filepath.Dir(wrapperPath)
	soPath := filepath.Join(pkgDir, filepath.Base(pkgPath)+".so")

	goModPath := filepath.Join(pkgDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		cmd := exec.CommandContext(ctx, "go", "mod", "init", "gig-plugin-"+filepath.Base(pkgPath))
		cmd.Dir = pkgDir
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("init go.mod: %w", err)
		}

		if gigRoot := findGigRoot(); gigRoot != "" {
			replaceCmd := exec.CommandContext(ctx, "go", "mod", "edit", "-replace", "github.com/t04dJ14n9/gig="+gigRoot)
			replaceCmd.Dir = pkgDir
			if err := replaceCmd.Run(); err != nil {
				return "", fmt.Errorf("add replace directive: %w", err)
			}
		}
	}

	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = pkgDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("go mod tidy: %w", err)
	}

	cmd = exec.CommandContext(ctx, "go", "build", "-buildmode=plugin", "-o", soPath, wrapperPath)
	cmd.Dir = pkgDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %w", ErrPluginBuildFailed, err)
	}

	return soPath, nil
}

// findGigRoot finds the root directory of the gig module.
func findGigRoot() string {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-f", "{{.Dir}}", "github.com/t04dJ14n9/gig")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// loadPlugin loads a .so file and calls its Register function.
func (pm *Manager) loadPlugin(soPath, pkgPath string) error {
	return pm.loadPluginInternal(soPath, pkgPath)
}

// ListLoaded returns a list of loaded plugin paths.
func (pm *Manager) ListLoaded() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	paths := make([]string, 0, len(pm.loaded))
	for p := range pm.loaded {
		paths = append(paths, p)
	}
	return paths
}
