//go:build windows

// Package pluginmgr provides hot-loading of external Go packages.
//
// This file contains the Windows-specific plugin stub.
// Go plugins are not supported on Windows, so this always returns an error.
package pluginmgr

import (
	"errors"
)

// loadPluginInternal returns an error on Windows.
func (pm *Manager) loadPluginInternal(soPath, pkgPath string) error {
	return errors.New("plugin loading is not supported on Windows; use Linux or macOS, or WSL")
}
