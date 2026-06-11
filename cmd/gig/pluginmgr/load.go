package pluginmgr

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/t04dJ14n9/gig"
)

// LoadPackage attempts to load an external package.
func (pm *Manager) LoadPackage(pkgPath string) error {
	if gig.GetPackageByPath(pkgPath) != nil {
		return nil
	}
	if runtime.GOOS == "windows" {
		return ErrPluginNotSupported
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.loaded[pkgPath] {
		return nil
	}

	if meta, ok := pm.registry[pkgPath]; ok {
		if _, err := os.Stat(meta.SOPath); err == nil {
			if err := pm.loadPlugin(meta.SOPath, pkgPath); err == nil {
				pm.loaded[pkgPath] = true
				return nil
			}
		}
	}

	return pm.buildAndLoad(pkgPath)
}

// buildAndLoad downloads, builds, and loads a package as a plugin.
func (pm *Manager) buildAndLoad(pkgPath string) error {
	if err := pm.downloadPackage(pkgPath); err != nil {
		return err
	}

	wrapperPath, err := pm.generateWrapper(pkgPath)
	if err != nil {
		return fmt.Errorf("generate wrapper: %w", err)
	}

	soPath, err := pm.buildPlugin(pkgPath, wrapperPath)
	if err != nil {
		return fmt.Errorf("build plugin: %w", err)
	}

	if err := pm.loadPlugin(soPath, pkgPath); err != nil {
		return fmt.Errorf("load plugin: %w", err)
	}

	pm.loaded[pkgPath] = true
	pm.registry[pkgPath] = pluginMetadata{
		Path:      pkgPath,
		SOPath:    soPath,
		BuildTime: fmt.Sprintf("%d", getCurrentTime()),
	}
	_ = pm.saveRegistry()
	return nil
}

// downloadPackage ensures the package is available locally.
func (pm *Manager) downloadPackage(pkgPath string) error {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "go", "list", pkgPath)
	if err := cmd.Run(); err == nil {
		return nil
	}

	cmd = exec.CommandContext(ctx, "go", "get", pkgPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", ErrPackageNotFound, pkgPath)
	}
	return nil
}

// getCurrentTime returns the current Unix timestamp.
func getCurrentTime() int64 {
	return time.Now().Unix()
}
