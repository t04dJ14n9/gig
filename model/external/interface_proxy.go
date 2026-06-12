package external

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

// InterfaceMethodCaller calls an interpreted method on a receiver captured by
// an interface proxy factory.
type InterfaceMethodCaller func(methodName string, args ...value.Value) (value.Value, bool)

// InterfaceProxyFactory builds a Go value that implements a native interface
// and forwards its methods back into interpreted code.
type InterfaceProxyFactory func(receiver value.Value, receiverTypeName string, call InterfaceMethodCaller) (any, bool)

// InterfaceProxyInfo describes a native interface proxy registered for an
// external package type.
type InterfaceProxyInfo struct {
	PkgPath         string
	Name            string
	InterfaceType   reflect.Type
	RequiredMethods []string
	Factory         InterfaceProxyFactory
}
