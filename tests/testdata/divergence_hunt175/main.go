package divergence_hunt175

import (
	"container/ring"
	"fmt"
)

// ============================================================================
// Round 175: Ring container operations
// ============================================================================

func RingNew() string {
	r := ring.New(3)
	return fmt.Sprintf("len=%d", r.Len())
}

func RingInit() string {
	r := ring.New(3)
	for i := 0; i < r.Len(); i++ {
		r.Value = i + 1
		r = r.Next()
	}
	result := ""
	r.Do(func(x any) {
		result += fmt.Sprintf("%v ", x)
	})
	return result
}

func RingNext() string {
	r := ring.New(3)
	r.Value = 1
	r.Next().Value = 2
	r.Next().Next().Value = 3
	return fmt.Sprintf("%v", r.Next().Value)
}

func RingPrev() string {
	r := ring.New(3)
	r.Value = 1
	r.Next().Value = 2
	r.Next().Next().Value = 3
	return fmt.Sprintf("%v", r.Prev().Value)
}

func RingMove() string {
	r := ring.New(5)
	for i := 0; i < r.Len(); i++ {
		r.Value = i + 1
		r = r.Next()
	}
	r = r.Move(2)
	return fmt.Sprintf("%v", r.Value)
}

func RingMoveNegative() string {
	r := ring.New(5)
	for i := 0; i < r.Len(); i++ {
		r.Value = i + 1
		r = r.Next()
	}
	r = r.Move(-1)
	return fmt.Sprintf("%v", r.Value)
}

func RingLink() string {
	r1 := ring.New(2)
	r1.Value = 1
	r1.Next().Value = 2

	r2 := ring.New(2)
	r2.Value = 3
	r2.Next().Value = 4

	r1.Link(r2)
	result := ""
	r1.Do(func(x any) {
		result += fmt.Sprintf("%v ", x)
	})
	return result
}

func RingUnlink() string {
	r := ring.New(5)
	for i := 0; i < r.Len(); i++ {
		r.Value = i + 1
		r = r.Next()
	}
	removed := r.Unlink(2)
	result := ""
	removed.Do(func(x any) {
		result += fmt.Sprintf("%v ", x)
	})
	return result
}

func RingDo() string {
	r := ring.New(4)
	for i := 0; i < r.Len(); i++ {
		r.Value = i + 1
		r = r.Next()
	}
	sum := 0
	r.Do(func(x any) {
		sum += x.(int)
	})
	return fmt.Sprintf("%d", sum)
}

func RingDoEmpty() string {
	r := ring.New(0)
	count := 0
	r.Do(func(x any) {
		count++
	})
	return fmt.Sprintf("count=%d", count)
}

func RingSingleElement() string {
	r := ring.New(1)
	r.Value = 42
	return fmt.Sprintf("%v next=%v", r.Value, r.Next().Value)
}

func RingCircular() string {
	r := ring.New(3)
	for i := 0; i < r.Len(); i++ {
		r.Value = i + 1
		r = r.Next()
	}
	// Go around the ring twice
	result := ""
	for i := 0; i < 6; i++ {
		result += fmt.Sprintf("%v ", r.Value)
		r = r.Next()
	}
	return result
}
