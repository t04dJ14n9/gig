package main

import "github.com/lithammer/fuzzysearch/fuzzy"

func FuzzyMatch() bool {
	return fuzzy.Match("whl", "wheel")
}

func FuzzyMatchFold() bool {
	return fuzzy.MatchFold("WHL", "wheel")
}

func FuzzyFind() int {
	results := fuzzy.Find("whl", []string{"wheel", "bottle", "whale"})
	return len(results)
}

func FuzzyLevenshteinDistance() int {
	return fuzzy.LevenshteinDistance("kitten", "sitting")
}

func FuzzyRankMatch() int {
	return fuzzy.RankMatch("whl", "wheel")
}
