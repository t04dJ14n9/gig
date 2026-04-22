package main

import "github.com/rivo/uniseg"

func UnisegGraphemeClusterCount() int {
	return uniseg.GraphemeClusterCount("🇩🇪👍🏽")
}

func UnisegStringWidth() int {
	return uniseg.StringWidth("Hello, 世界")
}

func UnisegWordCount() int {
	count := 0
	rest := "Hello, world! How are you?"
	state := 0
	for len(rest) > 0 {
		_, rest, state = uniseg.FirstWordInString(rest, state)
		count++
	}
	return count
}

func UnisegFirstGraphemeCluster() string {
	cluster, _, _, _ := uniseg.FirstGraphemeClusterInString("🇩🇪👍🏽", 0)
	return cluster
}

func UnisegHasNewline() bool {
	_, _, boundaries, _ := uniseg.FirstGraphemeClusterInString("line1\nline2", 0)
	return boundaries&uniseg.MaskLine != 0
}
