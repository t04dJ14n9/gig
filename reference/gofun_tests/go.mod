module github.com/t04dJ14n9/gig/reference/gofun_tests

go 1.23.1

require (
	git.code.oa.com/datacenter/onefun v1.0.2
	github.com/t04dJ14n9/gig v0.0.0
)

require (
	git.code.oa.com/gcloud_storage_group/tcaplus-go-api v0.2.0 // indirect
	git.code.oa.com/tsf4g/tdrcom v0.0.0-20200426021242-6f024d6c8199 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

replace (
	git.code.oa.com/datacenter/onefun/gofun => ../faas/vendor/git.code.oa.com/datacenter/onefun/gofun
	github.com/t04dJ14n9/gig => ../..
	github.com/cihub/seelog => ../faas/vendor/github.com/cihub/seelog
)
