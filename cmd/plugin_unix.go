//go:build !windows

// Package main provides the gig CLI tool.
//
// This file contains the plugin loading implementation for Unix-like systems
// (Linux and macOS). It uses the Go plugin package to dynamically load
// shared libraries (.so files).
package main

import (
	"fmt"
	"plugin"
)

// loadPluginInternal implements plugin loading for Unix-like systems.
func (pm *PluginManager) loadPluginInternal(soPath, pkgPath string) error {
	// Open the plugin
	p, err := plugin.Open(soPath)
	if err != nil {
		return fmt.Errorf("open plugin: %w", err)
	}

	// Look up the Register symbol
	sym, err := p.Lookup("Register")
	if err != nil {
		return fmt.Errorf("lookup Register: %w", err)
	}

	// Call the Register function
	register, ok := sym.(func())
	if !ok {
		return fmt.Errorf("Register is not a function")
	}
	register()

	return nil
}
