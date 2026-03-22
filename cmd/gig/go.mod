module github.com/t04dJ14n9/gig/cmd/gig

go 1.23.1

require (
	github.com/t04dJ14n9/gig v0.0.0
	github.com/peterh/liner v1.2.2
)

require (
	github.com/mattn/go-runewidth v0.0.3 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
)

replace github.com/t04dJ14n9/gig => ../..

replace github.com/t04dJ14n9/gig/cmd/gig/gentool => ./gentool
