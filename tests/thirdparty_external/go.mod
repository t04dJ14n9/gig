module thirdpartytests

go 1.23.1

require (
	git.woa.com/youngjin/gig v0.0.0
	github.com/Masterminds/semver/v3 v3.3.1
	github.com/google/uuid v1.6.0
	github.com/shopspring/decimal v1.4.0
)

require golang.org/x/tools v0.30.0 // indirect

replace git.woa.com/youngjin/gig => ../..
