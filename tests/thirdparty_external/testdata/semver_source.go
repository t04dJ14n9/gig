package main

import "github.com/Masterminds/semver/v3"

func SemverCheck(ver, rule string) bool {
	v, err := semver.NewVersion(ver)
	if err != nil {
		return false
	}
	c, err := semver.NewConstraint(rule)
	if err != nil {
		return false
	}
	return c.Check(v)
}

func SemverIncPatch(ver string) string {
	v, _ := semver.NewVersion(ver)
	return v.IncPatch().String()
}

func SemverCompare(a, b string) int {
	v1, _ := semver.NewVersion(a)
	v2, _ := semver.NewVersion(b)
	return v1.Compare(v2)
}
