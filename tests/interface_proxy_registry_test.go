package tests

import (
	"container/heap"
	"reflect"
	"sort"
	"testing"

	"github.com/t04dJ14n9/gig/importer"
	"github.com/t04dJ14n9/gig/model/external"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

type interfaceProxyLookup interface {
	LookupInterfaceProxyByType(reflect.Type) (*external.InterfaceProxyInfo, bool)
}

func TestStdlibRegistersInterfaceProxies(t *testing.T) {
	lookup := importer.GlobalRegistry().(interfaceProxyLookup)

	for name, ifaceType := range map[string]reflect.Type{
		"sort.Interface":           reflect.TypeOf((*sort.Interface)(nil)).Elem(),
		"container/heap.Interface": reflect.TypeOf((*heap.Interface)(nil)).Elem(),
	} {
		t.Run(name, func(t *testing.T) {
			info, ok := lookup.LookupInterfaceProxyByType(ifaceType)
			if !ok {
				t.Fatalf("missing proxy for %s", name)
			}
			if info.Factory == nil {
				t.Fatalf("proxy for %s has nil factory", name)
			}
		})
	}
}
