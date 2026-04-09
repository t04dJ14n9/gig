package strange_syntax_known_issue

// This file contains test cases that are known to fail due to interpreter limitations.
// These tests are isolated so they don't cause the main test suite to fail.
//
// ARCHITECTURAL LIMITATIONS:
// 1. Self-referential structs with slice fields need special handling for type conversion
// 2. Defer calling external type methods directly doesn't work, must use closure wrapper
//
// Note: encoding/xml and encoding/gob have been dropped from gig.
// Use encoding/json for serialization instead.

import (
	"bytes"
	"encoding/base64"
)

// ============================================================================
// KNOWN ISSUE #1: Tree With Parent Reference (Self-referential struct with slices)
// LIMITATION: Self-referential structs with slice fields of self-referential types
// require special type conversion logic. The field type becomes []interface{}
// during cycle breaking, but slice literals are created as []*struct.
// WORKAROUND: Use only pointer fields for self-references, avoid slices of self-referential types.
// ============================================================================

func TreeWithParentRef() int {
	type TreeNode struct {
		Value    int
		Children []*TreeNode
		Parent   *TreeNode
	}

	root := &TreeNode{Value: 1}
	child1 := &TreeNode{Value: 2, Parent: root}
	child2 := &TreeNode{Value: 3, Parent: root}
	root.Children = []*TreeNode{child1, child2}

	sum := root.Value
	for _, child := range root.Children {
		sum += child.Value
		if child.Parent != nil {
			sum += child.Parent.Value
		}
	}
	return sum // 1 + 2 + 1 + 3 + 1 = 8
}

// ============================================================================
// KNOWN ISSUE #2: Defer with External Type Methods
// LIMITATION: defer encoder.Close() doesn't work correctly for external types.
// The defer statement captures the method call but execution may be incomplete.
// WORKAROUND: Use a closure wrapper: defer func() { encoder.Close() }()
// ============================================================================

func CombinedDeferWithIO() int {
	var buf bytes.Buffer

	_ = func() int {
		encoder := base64.NewEncoder(base64.StdEncoding, &buf)
		defer encoder.Close() // This doesn't work correctly

		encoder.Write([]byte("test"))

		// Buffer is written to even before Close
		return buf.Len()
	}()

	// After close, flush should be complete but isn't
	decoded, _ := base64.StdEncoding.DecodeString(buf.String())

	if string(decoded) == "test" {
		return 1
	}
	return 0
}
