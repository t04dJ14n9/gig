// Package sort registers the Go standard library sort package.
package sort

import (
	"sort"

	"gig/importer"
)

func init() {
	pkg := importer.RegisterPackage("sort", "sort")

	// Basic sort functions
	pkg.AddFunction("Ints", sort.Ints, "", nil)
	pkg.AddFunction("Float64s", sort.Float64s, "", nil)
	pkg.AddFunction("Strings", sort.Strings, "", nil)
	pkg.AddFunction("IntsAreSorted", sort.IntsAreSorted, "", nil)
	pkg.AddFunction("Float64sAreSorted", sort.Float64sAreSorted, "", nil)
	pkg.AddFunction("StringsAreSorted", sort.StringsAreSorted, "", nil)

	// Search functions
	pkg.AddFunction("Search", sort.Search, "", nil)
	pkg.AddFunction("SearchInts", sort.SearchInts, "", nil)
	pkg.AddFunction("SearchFloat64s", sort.SearchFloat64s, "", nil)
	pkg.AddFunction("SearchStrings", sort.SearchStrings, "", nil)

	// Generic sort
	pkg.AddFunction("Slice", sort.Slice, "", nil)
	pkg.AddFunction("SliceIsSorted", sort.SliceIsSorted, "", nil)
	pkg.AddFunction("SliceStable", sort.SliceStable, "", nil)
	pkg.AddFunction("Sort", sort.Sort, "", nil)
	pkg.AddFunction("Stable", sort.Stable, "", nil)
	pkg.AddFunction("IsSorted", sort.IsSorted, "", nil)
	pkg.AddFunction("Reverse", sort.Reverse, "", nil)

	// Types
	pkg.AddType("IntSlice", nil, "int slice for sorting")
	pkg.AddType("Float64Slice", nil, "float64 slice for sorting")
	pkg.AddType("StringSlice", nil, "string slice for sorting")
}
