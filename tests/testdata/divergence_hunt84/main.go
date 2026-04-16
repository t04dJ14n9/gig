package divergence_hunt84

// ============================================================================
// Round 84: Recursive data structures - linked list, binary tree
// ============================================================================

type ListNode struct {
	Val  int
	Next *ListNode
}

func LinkedListSum() int {
	head := &ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 3}}}
	sum := 0
	for n := head; n != nil; n = n.Next {
		sum += n.Val
	}
	return sum
}

func LinkedListLength() int {
	head := &ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 3}}}
	count := 0
	for n := head; n != nil; n = n.Next {
		count++
	}
	return count
}

func LinkedListReverse() int {
	head := &ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 3}}}
	var prev *ListNode
	for n := head; n != nil; {
		next := n.Next
		n.Next = prev
		prev = n
		n = next
	}
	return prev.Val // should be 3
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func TreeSum() int {
	root := &TreeNode{
		Val:   1,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 3, Left: &TreeNode{Val: 6}, Right: &TreeNode{Val: 7}},
	}
	var sum func(*TreeNode) int
	sum = func(n *TreeNode) int {
		if n == nil {
			return 0
		}
		return n.Val + sum(n.Left) + sum(n.Right)
	}
	return sum(root)
}

func TreeDepth() int {
	root := &TreeNode{
		Val:   1,
		Left:  &TreeNode{Val: 2},
		Right: &TreeNode{Val: 3, Right: &TreeNode{Val: 4}},
	}
	var depth func(*TreeNode) int
	depth = func(n *TreeNode) int {
		if n == nil {
			return 0
		}
		l := depth(n.Left)
		r := depth(n.Right)
		if l > r {
			return l + 1
		}
		return r + 1
	}
	return depth(root)
}

func LinkedListCreate() int {
	// Create list from slice
	vals := []int{10, 20, 30, 40, 50}
	var head *ListNode
	for i := len(vals) - 1; i >= 0; i-- {
		head = &ListNode{Val: vals[i], Next: head}
	}
	return head.Val
}

func LinkedListMiddle() int {
	head := &ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4, Next: &ListNode{Val: 5}}}}}
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow.Val
}

func TreeLeafCount() int {
	root := &TreeNode{
		Val:   1,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 3},
	}
	var leafCount func(*TreeNode) int
	leafCount = func(n *TreeNode) int {
		if n == nil {
			return 0
		}
		if n.Left == nil && n.Right == nil {
			return 1
		}
		return leafCount(n.Left) + leafCount(n.Right)
	}
	return leafCount(root)
}
