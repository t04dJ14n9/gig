package main

import "github.com/cespare/xxhash/v2"

func XxhashSum64() uint64 {
	return xxhash.Sum64([]byte("hello"))
}

func XxhashSum64String() uint64 {
	return xxhash.Sum64String("hello")
}

func XxhashNew() uint64 {
	h := xxhash.New()
	h.WriteString("hello")
	return h.Sum64()
}

func XxhashNewWithSeed() uint64 {
	h := xxhash.NewWithSeed(42)
	h.WriteString("hello")
	return h.Sum64()
}

func XxhashDigestSize() int {
	h := xxhash.New()
	return h.Size()
}

func XxhashDigestBlockSize() int {
	h := xxhash.New()
	return h.BlockSize()
}

func XxhashDigestWrite() int {
	h := xxhash.New()
	n, _ := h.Write([]byte("hello"))
	return n
}

func XxhashDigestWriteString() int {
	h := xxhash.New()
	n, _ := h.WriteString("hello")
	return n
}

func XxhashDigestReset() uint64 {
	h := xxhash.New()
	h.WriteString("hello")
	h.Reset()
	h.WriteString("world")
	return h.Sum64()
}

func XxhashDigestResetWithSeed() uint64 {
	h := xxhash.New()
	h.WriteString("hello")
	h.ResetWithSeed(123)
	h.WriteString("hello")
	return h.Sum64()
}

func XxhashDigestSum() int {
	h := xxhash.New()
	h.WriteString("hello")
	b := h.Sum(nil)
	return len(b)
}

func XxhashDigestMarshalBinary() bool {
	h := xxhash.New()
	h.WriteString("hello")
	data, err := h.MarshalBinary()
	if err != nil {
		return false
	}
	return len(data) > 0
}

func XxhashDigestUnmarshalBinary() uint64 {
	h1 := xxhash.New()
	h1.WriteString("hello")
	data, _ := h1.MarshalBinary()
	h2 := xxhash.New()
	h2.UnmarshalBinary(data)
	return h2.Sum64()
}
