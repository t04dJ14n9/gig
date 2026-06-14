module github.com/t04dJ14n9/gig/cmd/gig

go 1.23.1

require (
	github.com/peterh/liner v1.2.2
	github.com/t04dJ14n9/gig v1.7.1-0.20260612135619-58884bff1673
)

require (
	github.com/mattn/go-runewidth v0.0.3 // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
)

// Use the local repository for all gig modules. This keeps the CLI in
// lock-step with the in-tree v2 SSA backend during refactor.
replace github.com/t04dJ14n9/gig => ../..
