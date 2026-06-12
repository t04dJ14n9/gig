package pluginmgr

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// loadRegistry loads the plugin registry from disk.
func (pm *Manager) loadRegistry() {
	registryPath := filepath.Join(pm.pluginDir, "plugin_registry.json")
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return
	}
	//nolint:errchkjson // Registry file is internal, ignore unmarshal errors.
	_ = json.Unmarshal(data, &pm.registry)
}

// saveRegistry saves the plugin registry to disk.
func (pm *Manager) saveRegistry() error {
	if err := os.MkdirAll(pm.pluginDir, 0o755); err != nil {
		return err
	}
	registryPath := filepath.Join(pm.pluginDir, "plugin_registry.json")
	data := mustMarshalJSON(pm.registry)
	return os.WriteFile(registryPath, data, 0o644)
}

// mustMarshalJSON marshals to JSON, panicking on error.
func mustMarshalJSON(v any) []byte {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return data
}
