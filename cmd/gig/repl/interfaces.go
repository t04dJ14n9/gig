package repl

import "github.com/t04dJ14n9/gig/cmd/gig/pluginmgr"

type packagePluginManager interface {
	LoadPackage(pkgPath string) error
	GetSymbols(pkgPath string) []string
	ListLoaded() []string
}

var newPluginManager = func() packagePluginManager {
	return pluginmgr.NewManager()
}
