package divergence_hunt173

import (
	"container/heap"
	"fmt"
)

// ============================================================================
// Round 173: Heap interface operations
// ============================================================================

// IntHeap implements heap.Interface
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func HeapInit() string {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	return fmt.Sprintf("%v", *h)
}

func HeapPush() string {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	return fmt.Sprintf("%v", *h)
}

func HeapPop() string {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	result := heap.Pop(h).(int)
	return fmt.Sprintf("%d:%v", result, *h)
}

func HeapRemove() string {
	h := &IntHeap{1, 2, 5, 3}
	heap.Init(h)
	result := heap.Remove(h, 1).(int)
	return fmt.Sprintf("%d:%v", result, *h)
}

func HeapFix() string {
	h := &IntHeap{1, 5, 3, 4, 2}
	heap.Init(h)
	(*h)[2] = 6
	heap.Fix(h, 2)
	return fmt.Sprintf("%v", *h)
}

func HeapMultiplePushPop() string {
	h := &IntHeap{}
	heap.Init(h)
	values := []int{3, 1, 4, 1, 5, 9, 2, 6}
	for _, v := range values {
		heap.Push(h, v)
	}
	result := ""
	for h.Len() > 0 {
		result += fmt.Sprintf("%d ", heap.Pop(h).(int))
	}
	return result
}

// MaxHeap implements a max heap
type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] } // Max heap
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func MaxHeapOperations() string {
	h := &MaxHeap{}
	heap.Init(h)
	heap.Push(h, 10)
	heap.Push(h, 30)
	heap.Push(h, 20)
	result := heap.Pop(h).(int)
	return fmt.Sprintf("%d", result)
}

// Item for priority queue
type Item struct {
	Value    string
	Priority int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func PriorityQueueOperations() string {
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{"task1", 3})
	heap.Push(pq, &Item{"task2", 1})
	heap.Push(pq, &Item{"task3", 2})
	result := ""
	for pq.Len() > 0 {
		item := heap.Pop(pq).(*Item)
		result += fmt.Sprintf("%s:%d ", item.Value, item.Priority)
	}
	return result
}

func HeapEmpty() string {
	h := &IntHeap{}
	heap.Init(h)
	return fmt.Sprintf("len=%d", h.Len())
}

func HeapSingleElement() string {
	h := &IntHeap{42}
	heap.Init(h)
	return fmt.Sprintf("%v", *h)
}
