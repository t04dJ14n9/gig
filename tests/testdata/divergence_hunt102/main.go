package divergence_hunt102

import "fmt"

// ============================================================================
// Round 102: Recursive data structures - trees, graphs
// ============================================================================

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func TreeSum(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return root.Val + TreeSum(root.Left) + TreeSum(root.Right)
}

func TreeDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	l := TreeDepth(root.Left)
	r := TreeDepth(root.Right)
	if l > r {
		return l + 1
	}
	return r + 1
}

func TreeLeafCount(root *TreeNode) int {
	if root == nil {
		return 0
	}
	if root.Left == nil && root.Right == nil {
		return 1
	}
	return TreeLeafCount(root.Left) + TreeLeafCount(root.Right)
}

func TreeInorder(root *TreeNode) []int {
	var result []int
	var traverse func(*TreeNode)
	traverse = func(n *TreeNode) {
		if n == nil {
			return
		}
		traverse(n.Left)
		result = append(result, n.Val)
		traverse(n.Right)
	}
	traverse(root)
	return result
}

func BSTInsert(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return &TreeNode{Val: val}
	}
	if val < root.Val {
		root.Left = BSTInsert(root.Left, val)
	} else {
		root.Right = BSTInsert(root.Right, val)
	}
	return root
}

func TreeBuildAndSum() int {
	root := &TreeNode{Val: 5}
	BSTInsert(root, 3)
	BSTInsert(root, 7)
	BSTInsert(root, 1)
	BSTInsert(root, 4)
	return TreeSum(root)
}

func TreeBuildAndDepth() int {
	root := &TreeNode{Val: 5}
	BSTInsert(root, 3)
	BSTInsert(root, 7)
	BSTInsert(root, 1)
	return TreeDepth(root)
}

func TreeBuildAndLeaves() int {
	root := &TreeNode{Val: 5}
	BSTInsert(root, 3)
	BSTInsert(root, 7)
	BSTInsert(root, 1)
	return TreeLeafCount(root)
}

func TreeInorderResult() string {
	root := &TreeNode{Val: 5}
	BSTInsert(root, 3)
	BSTInsert(root, 7)
	BSTInsert(root, 1)
	BSTInsert(root, 4)
	return fmt.Sprintf("%v", TreeInorder(root))
}

func FibonacciTree() int {
	// Build a tree where Val = fibonacci(n)
	var build func(n int) *TreeNode
	build = func(n int) *TreeNode {
		if n <= 1 {
			return &TreeNode{Val: n}
		}
		return &TreeNode{Val: n, Left: build(n - 1), Right: build(n - 2)}
	}
	root := build(4)
	return TreeSum(root)
}
