package tests

import (
	"container/heap"
	_ "embed"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
	known_issues "github.com/t04dJ14n9/gig/tests/testdata/known_issues"
)

//go:embed testdata/known_issues/main.go
var resolvedIssuesSrc string

type resolvedIssue struct {
	funcName string
	native   func() any
	issue    string
}

func runResolvedIssueTest(t *testing.T, prog *gig.Program, name string, tc resolvedIssue) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		nativeResult := tc.native()

		var interpResult any
		var interpErr error
		panicked := false
		var panicVal any

		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
					panicVal = r
				}
			}()
			interpResult, interpErr = prog.Run(tc.funcName)
		}()

		if panicked {
			t.Errorf("REGRESSION (panic): %s\n  interpreter panicked: %v\n  native returned:      %v (%T)",
				tc.issue, panicVal, nativeResult, nativeResult)
			return
		}

		if interpErr != nil {
			t.Errorf("REGRESSION (error): %s\n  interpreter error: %v\n  native returned:   %v (%T)",
				tc.issue, interpErr, nativeResult, nativeResult)
			return
		}

		if !reflect.DeepEqual(interpResult, nativeResult) {
			t.Errorf("REGRESSION (mismatch): %s\n  interpreter returned: %v (%T)\n  native returned:      %v (%T)",
				tc.issue, interpResult, interpResult, nativeResult, nativeResult)
		}
	})
}

func TestResolvedIssues(t *testing.T) {
	issues := map[string]resolvedIssue{
		"InterfaceWithNilConcrete": {
			funcName: "InterfaceWithNilConcrete",
			native: func() any {
				var p *known_issues.PointerReceiver
				var s known_issues.Stringer = p
				return s.String()
			},
			issue: "Interface holding nil pointer is incorrectly treated as nil",
		},
		"NestedNilReceiver": {
			funcName: "NestedNilReceiver",
			native: func() any {
				type Inner struct{ Name string }
				type Outer struct{ *Inner }
				outer := Outer{}
				defer func() { recover() }()
				return outer.Name
			},
			issue: "Accessing promoted field on nil embedded pointer panics",
		},
		"SortByLength": {
			funcName: "SortByLength",
			native: func() any {
				words := []string{"apple", "pie", "banana", "kiwi"}
				sort.Sort(known_issues.ByLength(words))
				return fmt.Sprintf("%v", words)
			},
			issue: "sort.Sort with custom sort.Interface fails (reflection can't see methods)",
		},
		"SortReverse": {
			funcName: "SortReverse",
			native: func() any {
				nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
				sort.Sort(known_issues.Reverse{Interface: sort.IntSlice(nums)})
				return fmt.Sprintf("%v", nums)
			},
			issue: "sort.Reverse wrapper fails (reflection can't see embedded interface methods)",
		},
		"SortEmbeddedInterfaceDescending": {
			funcName: "SortEmbeddedInterfaceDescending",
			native: func() any {
				nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
				sort.Sort(known_issues.Descending{Interface: sort.IntSlice(nums)})
				return fmt.Sprintf("%v", nums)
			},
			issue: "embedded sort.Interface wrappers must dispatch compiled Less regardless of receiver type name",
		},
		"HeapInit": {
			funcName: "HeapInit",
			native: func() any {
				h := &known_issues.IntHeap{2, 1, 5}
				heap.Init(h)
				return fmt.Sprintf("%v", *h)
			},
			issue: "heap.Init fails (reflection can't see heap.Interface methods)",
		},
		"HeapPush": {
			funcName: "HeapPush",
			native: func() any {
				h := &known_issues.IntHeap{2, 1, 5}
				heap.Init(h)
				heap.Push(h, 3)
				return fmt.Sprintf("%v", *h)
			},
			issue: "heap.Push fails (reflection can't see heap.Interface methods)",
		},
		"HeapPop": {
			funcName: "HeapPop",
			native: func() any {
				h := &known_issues.IntHeap{2, 1, 5}
				heap.Init(h)
				result := heap.Pop(h).(int)
				return fmt.Sprintf("%d:%v", result, *h)
			},
			issue: "heap.Pop fails (reflection can't see heap.Interface methods)",
		},
	}

	if len(issues) == 0 {
		t.Log("No resolved issues")
		return
	}

	prog, err := gig.Build(resolvedIssuesSrc, gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	for name, tc := range issues {
		runResolvedIssueTest(t, prog, name, tc)
	}
}
