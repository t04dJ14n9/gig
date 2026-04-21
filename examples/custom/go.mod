module myapp

go 1.23.1

require (
	git.woa.com/youngjin/gig v0.0.0
	github.com/dromara/carbon/v2 v2.6.9
	github.com/spf13/cast v1.10.0
	github.com/tidwall/gjson v1.18.0
	github.com/tidwall/sjson v1.2.5
)

require (
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
)

replace git.woa.com/youngjin/gig => ../..
