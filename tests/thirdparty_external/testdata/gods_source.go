package main

import "github.com/emirpasic/gods/sets/hashset"

func GodsHashSetNew() int {
	s := hashset.New()
	s.Add(1)
	s.Add(2)
	s.Add(2)
	return s.Size()
}

func GodsHashSetValues() int {
	s := hashset.New()
	s.Add("a", "b", "c")
	return len(s.Values())
}

func GodsHashSetContains() bool {
	s := hashset.New()
	s.Add(10, 20, 30)
	return s.Contains(20)
}

func GodsHashSetRemove() int {
	s := hashset.New()
	s.Add(1, 2, 3)
	s.Remove(2)
	return s.Size()
}
