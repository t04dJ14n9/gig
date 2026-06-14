package thirdparty

import (
	"container/list"
	"container/ring"
)

// container/heap tests removed: heap.Init/Push/Pop require passing
// an interpreted-defined type that implements heap.Interface (a host
// non-empty interface). Per G_iface_ban that boundary is rejected at
// compile time.

// ============================================================================
// container/list — doubly-linked list
// ============================================================================

// ContainerListPushFrontBack tests list push operations.
func ContainerListPushFrontBack() int {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushFront(0)
	l.PushBack(3)
	// List: 0, 1, 2, 3
	if l.Len() == 4 && l.Front().Value.(int) == 0 && l.Back().Value.(int) == 3 {
		return 1
	}
	return 0
}

// ContainerListRemove tests list remove.
func ContainerListRemove() int {
	l := list.New()
	l.PushBack(1)
	e2 := l.PushBack(2)
	l.PushBack(3)
	l.Remove(e2)
	if l.Len() == 2 && l.Front().Next().Value.(int) == 3 {
		return 1
	}
	return 0
}

// ContainerListMove tests list element movement.
func ContainerListMove() int {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	e := l.Front()
	l.MoveToFront(e)
	if l.Front().Value.(int) == 1 && l.Front() == e {
		return 1
	}
	return 0
}

// ContainerListIterate tests list iteration.
func ContainerListIterate() int {
	l := list.New()
	for i := 1; i <= 5; i++ {
		l.PushBack(i)
	}
	sum := 0
	for e := l.Front(); e != nil; e = e.Next() {
		sum += e.Value.(int)
	}
	if sum == 15 {
		return 1
	}
	return 0
}

// ContainerListReverseIterate tests reverse list iteration.
func ContainerListReverseIterate() int {
	l := list.New()
	for i := 1; i <= 5; i++ {
		l.PushBack(i)
	}
	sum := 0
	for e := l.Back(); e != nil; e = e.Prev() {
		sum += e.Value.(int)
	}
	if sum == 15 {
		return 1
	}
	return 0
}

// ============================================================================
// container/ring — circular ring
// ============================================================================

// ContainerRingSum tests ring iteration and sum.
func ContainerRingSum() int {
	r := ring.New(5)
	for i := 0; i < 5; i++ {
		r.Value = i + 1
		r = r.Next()
	}
	sum := 0
	r.Do(func(v any) {
		sum += v.(int)
	})
	if sum == 15 {
		return 1
	}
	return 0
}

// ContainerRingMove tests ring movement.
func ContainerRingMove() int {
	r := ring.New(3)
	r.Value = 1
	r = r.Next()
	r.Value = 2
	r = r.Next()
	r.Value = 3
	// Move forward 2 steps from value 2
	r = r.Move(2)
	if r.Value.(int) == 1 {
		return 1
	}
	return 0
}

// ContainerRingLink tests ring linking (circular join).
func ContainerRingLink() int {
	r1 := ring.New(2)
	r1.Value = 1
	r1 = r1.Next()
	r1.Value = 2

	r2 := ring.New(2)
	r2.Value = 3
	r2 = r2.Next()
	r2.Value = 4

	// Link r1 and r2
	r1.Link(r2)

	// Count total elements
	count := 0
	r1.Do(func(v any) { count++ })
	if count == 4 {
		return 1
	}
	return 0
}

// ContainerRingUnlink tests ring unlink.
func ContainerRingUnlink() int {
	r := ring.New(5)
	for i := 0; i < 5; i++ {
		r.Value = i + 1
		r = r.Next()
	}
	// Unlink 2 elements starting from r.Next()
	unlinked := r.Next().Unlink(2)
	count := 0
	unlinked.Do(func(v any) { count++ })
	if count == 2 {
		return 1
	}
	return 0
}
