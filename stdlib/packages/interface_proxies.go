package packages

import (
	container_heap "container/heap"
	"reflect"
	"sort"

	"github.com/t04dJ14n9/gig/importer"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

func init() {
	importer.AddInterfaceProxy(
		"sort",
		"Interface",
		reflect.TypeOf((*sort.Interface)(nil)).Elem(),
		[]string{"Len", "Less", "Swap"},
		newSortInterfaceProxy,
	)
	importer.AddInterfaceProxy(
		"container/heap",
		"Interface",
		reflect.TypeOf((*container_heap.Interface)(nil)).Elem(),
		[]string{"Len", "Less", "Swap", "Push", "Pop"},
		newHeapInterfaceProxy,
	)
}

type sortInterfaceProxy struct {
	call external.InterfaceMethodCaller
}

func newSortInterfaceProxy(_ value.Value, _ string, call external.InterfaceMethodCaller) (any, bool) {
	return &sortInterfaceProxy{call: call}, true
}

func (p *sortInterfaceProxy) Len() int {
	result, ok := p.call("Len")
	if !ok {
		return 0
	}
	return int(result.Int())
}

func (p *sortInterfaceProxy) Less(i, j int) bool {
	result, ok := p.call("Less", value.MakeInt(int64(i)), value.MakeInt(int64(j)))
	if !ok {
		return false
	}
	return result.Bool()
}

func (p *sortInterfaceProxy) Swap(i, j int) {
	_, _ = p.call("Swap", value.MakeInt(int64(i)), value.MakeInt(int64(j)))
}

type heapInterfaceProxy struct {
	call external.InterfaceMethodCaller
}

func newHeapInterfaceProxy(_ value.Value, _ string, call external.InterfaceMethodCaller) (any, bool) {
	return &heapInterfaceProxy{call: call}, true
}

func (p *heapInterfaceProxy) Len() int {
	result, ok := p.call("Len")
	if !ok {
		return 0
	}
	return int(result.Int())
}

func (p *heapInterfaceProxy) Less(i, j int) bool {
	result, ok := p.call("Less", value.MakeInt(int64(i)), value.MakeInt(int64(j)))
	if !ok {
		return false
	}
	return result.Bool()
}

func (p *heapInterfaceProxy) Swap(i, j int) {
	_, _ = p.call("Swap", value.MakeInt(int64(i)), value.MakeInt(int64(j)))
}

func (p *heapInterfaceProxy) Push(x any) {
	_, _ = p.call("Push", value.FromInterface(x))
}

func (p *heapInterfaceProxy) Pop() any {
	result, ok := p.call("Pop")
	if !ok {
		return nil
	}
	return result.Interface()
}
