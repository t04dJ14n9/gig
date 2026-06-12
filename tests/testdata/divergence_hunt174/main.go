package divergence_hunt174

import (
	"container/list"
	"fmt"
)

// ============================================================================
// Round 174: List container operations
// ============================================================================

func ListNew() string {
	l := list.New()
	return fmt.Sprintf("len=%d", l.Len())
}

func ListPushBack() string {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	return fmt.Sprintf("len=%d", l.Len())
}

func ListPushFront() string {
	l := list.New()
	l.PushFront(3)
	l.PushFront(2)
	l.PushFront(1)
	return fmt.Sprintf("len=%d", l.Len())
}

func ListFrontBack() string {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	front := l.Front().Value
	back := l.Back().Value
	return fmt.Sprintf("front=%v back=%v", front, back)
}

func ListIterateForward() string {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	result := ""
	for e := l.Front(); e != nil; e = e.Next() {
		result += fmt.Sprintf("%v ", e.Value)
	}
	return result
}

func ListIterateBackward() string {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	result := ""
	for e := l.Back(); e != nil; e = e.Prev() {
		result += fmt.Sprintf("%v ", e.Value)
	}
	return result
}

func ListInsertAfter() string {
	l := list.New()
	l.PushBack(1)
	l.PushBack(3)
	second := l.Front().Next()
	l.InsertAfter(2, second)
	result := ""
	for e := l.Front(); e != nil; e = e.Next() {
		result += fmt.Sprintf("%v ", e.Value)
	}
	return result
}

func ListInsertBefore() string {
	l := list.New()
	l.PushBack(1)
	l.PushBack(3)
	second := l.Front().Next()
	l.InsertBefore(2, second)
	result := ""
	for e := l.Front(); e != nil; e = e.Next() {
		result += fmt.Sprintf("%v ", e.Value)
	}
	return result
}

func ListRemove() string {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	second := l.Front().Next()
	l.Remove(second)
	result := ""
	for e := l.Front(); e != nil; e = e.Next() {
		result += fmt.Sprintf("%v ", e.Value)
	}
	return result
}

func ListMoveToFront() string {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.MoveToFront(l.Back())
	result := ""
	for e := l.Front(); e != nil; e = e.Next() {
		result += fmt.Sprintf("%v ", e.Value)
	}
	return result
}

func ListMoveToBack() string {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.MoveToBack(l.Front())
	result := ""
	for e := l.Front(); e != nil; e = e.Next() {
		result += fmt.Sprintf("%v ", e.Value)
	}
	return result
}

func ListMixedOperations() string {
	l := list.New()
	l.PushBack(1)
	l.PushFront(0)
	l.PushBack(2)
	l.Remove(l.Front())
	l.PushBack(3)
	return fmt.Sprintf("len=%d front=%v back=%v", l.Len(), l.Front().Value, l.Back().Value)
}
