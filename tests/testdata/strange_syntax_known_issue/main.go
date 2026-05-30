package strange_syntax_known_issue

// This file contains regression tests for interpreter limitations that have
// been fixed. The directory name is historical.
//
// Note: encoding/xml and encoding/gob have been dropped from gig.
// Use encoding/json for serialization instead.

import (
	"bytes"
	"encoding/base64"
)

// ============================================================================
// RESOLVED ISSUE #1: Tree With Parent Reference (self-referential struct with slices)
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
// RESOLVED ISSUE #2: Defer with External Type Methods
// ============================================================================

func CombinedDeferWithIO() int {
	var buf bytes.Buffer

	_ = func() int {
		encoder := base64.NewEncoder(base64.StdEncoding, &buf)
		defer encoder.Close()

		encoder.Write([]byte("test"))

		// Buffer is written to even before Close
		return buf.Len()
	}()

	// After close, flush should be complete.
	decoded, _ := base64.StdEncoding.DecodeString(buf.String())

	if string(decoded) == "test" {
		return 1
	}
	return 0
}
