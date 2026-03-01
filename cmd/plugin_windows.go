//go:build windows

// Package main provides the gig CLI tool.
//
// This file contains the Windows-specific plugin stub.
// Go plugins are not supported on Windows, so this always returns an error.
package main

import (
	"errors"
)

// loadPluginInternal returns an error on Windows.
func (pm *PluginManager) loadPluginInternal(soPath, pkgPath string) error {
	return errors.New("plugin loading is not supported on Windows; use Linux or macOS, or WSL")
}
