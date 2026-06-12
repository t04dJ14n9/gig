package vm

import "github.com/t04dJ14n9/gig/model/value"

func (v *vm) executeRange() {
	// Create an iterator for the collection
	collection := v.pop()
	v.push(value.FromInterface(&iterator{collection: collection, index: 0}))
}

func (v *vm) executeRangeNext() {
	// Advance iterator and push a tuple (ok, key, value)
	iterVal := v.pop()
	iter, ok := iterVal.Interface().(*iterator)
	if !ok {
		// Return tuple (false, nil, nil)
		tuple := []value.Value{value.MakeBool(false), value.MakeNil(), value.MakeNil()}
		v.push(value.FromInterface(tuple))
		return
	}
	key, val, iterOk := iter.next()
	// SSA Next returns (ok, key, value) as a tuple
	tuple := []value.Value{value.MakeBool(iterOk), key, val}
	v.push(value.FromInterface(tuple))
}
