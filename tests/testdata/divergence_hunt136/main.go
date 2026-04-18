package divergence_hunt136

import "fmt"

// ============================================================================
// Round 136: Deeply nested structs and field access
// ============================================================================

type Inner2 struct {
	Value int
}

type Middle struct {
	In Inner2
}

type Outer2 struct {
	Mid Middle
}

func DeepFieldAccess() string {
	o := Outer2{Mid: Middle{In: Inner2{Value: 42}}}
	return fmt.Sprintf("val=%d", o.Mid.In.Value)
}

func DeepFieldAssign() string {
	o := Outer2{Mid: Middle{In: Inner2{Value: 0}}}
	o.Mid.In.Value = 99
	return fmt.Sprintf("val=%d", o.Mid.In.Value)
}

func DeepFieldPointer() string {
	o := &Outer2{Mid: Middle{In: Inner2{Value: 7}}}
	return fmt.Sprintf("val=%d", o.Mid.In.Value)
}

type Node struct {
	Val  int
	Next *Node
}

func LinkedListTraversal() string {
	n3 := &Node{Val: 3, Next: nil}
	n2 := &Node{Val: 2, Next: n3}
	n1 := &Node{Val: 1, Next: n2}
	sum := 0
	for cur := n1; cur != nil; cur = cur.Next {
		sum += cur.Val
	}
	return fmt.Sprintf("sum=%d", sum)
}

func LinkedListCreate() string {
	head := &Node{Val: 10}
	head.Next = &Node{Val: 20}
	head.Next.Next = &Node{Val: 30}
	return fmt.Sprintf("vals=%d-%d-%d", head.Val, head.Next.Val, head.Next.Next.Val)
}

type Tree struct {
	Val   int
	Left  *Tree
	Right *Tree
}

func TreeTraversal() string {
	root := &Tree{
		Val:  5,
		Left:  &Tree{Val: 3},
		Right: &Tree{Val: 7},
	}
	return fmt.Sprintf("root=%d-left=%d-right=%d", root.Val, root.Left.Val, root.Right.Val)
}

type Address struct {
	City    string
	ZipCode int
}

type Person struct {
	Name    string
	Age     int
	Address Address
}

func NestedStructLiteral() string {
	p := Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			City:    "NYC",
			ZipCode: 10001,
		},
	}
	return fmt.Sprintf("%s-%d-%s-%d", p.Name, p.Age, p.Address.City, p.Address.ZipCode)
}

func NestedStructUpdate() string {
	p := Person{Name: "Bob", Age: 25, Address: Address{City: "LA", ZipCode: 90001}}
	p.Address.City = "SF"
	p.Age++
	return fmt.Sprintf("%s-%d-%s", p.Name, p.Age, p.Address.City)
}
