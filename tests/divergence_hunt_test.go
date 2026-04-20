// Package tests - divergence_hunt_test.go
//
// Divergence hunt tests: compare interpreted execution with native Go results.
// Uses //go:embed to load source from testdata/ directories, same pattern as correctness_test.go.
package tests

import (
	"context"
	_ "embed"
	"reflect"
	"testing"
	"time"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt1"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt2"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt3"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt4"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt5"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt6"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt7"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt8"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt9"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt10"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt11"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt12"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt13"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt14"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt15"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt16"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt17"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt18"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt19"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt20"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt21"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt22"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt23"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt24"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt25"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt26"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt27"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt28"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt29"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt30"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt31"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt32"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt33"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt34"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt35"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt36"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt37"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt38"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt39"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt40"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt41"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt42"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt43"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt44"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt45"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt46"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt47"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt48"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt49"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt50"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt51"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt52"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt53"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt54"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt55"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt56"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt57"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt58"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt59"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt60"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt61"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt62"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt63"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt64"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt65"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt66"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt67"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt68"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt69"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt70"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt71"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt72"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt73"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt74"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt75"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt76"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt77"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt78"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt79"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt80"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt81"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt82"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt83"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt84"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt85"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt86"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt87"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt88"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt89"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt90"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt91"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt92"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt93"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt94"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt95"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt96"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt97"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt98"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt99"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt100"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt101"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt102"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt103"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt104"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt105"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt106"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt107"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt108"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt109"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt110"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt111"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt112"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt113"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt114"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt115"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt116"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt117"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt118"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt119"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt120"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt121"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt122"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt123"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt124"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt125"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt126"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt127"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt128"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt129"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt130"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt131"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt132"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt133"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt134"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt135"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt136"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt137"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt138"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt139"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt140"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt141"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt142"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt143"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt144"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt145"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt146"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt147"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt148"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt149"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt150"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt151"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt152"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt153"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt154"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt155"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt156"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt157"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt158"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt159"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt160"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt161"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt162"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt163"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt164"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt165"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt166"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt167"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt168"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt169"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt170"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt171"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt172"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt173"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt174"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt175"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt176"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt177"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt178"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt179"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt180"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt181"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt182"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt183"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt184"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt185"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt186"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt187"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt188"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt189"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt190"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt191"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt192"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt193"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt194"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt195"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt196"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt197"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt198"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt199"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt200"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt201"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt202"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt203"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt204"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt205"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt206"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt207"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt208"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt209"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt210"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt211"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt212"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt213"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt214"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt215"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt216"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt217"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt218"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt219"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt220"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt221"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt222"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt223"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt224"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt225"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt226"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt227"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt228"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt229"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt230"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt231"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt232"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt233"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt234"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt235"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt236"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt237"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt238"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt239"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt240"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt241"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt242"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt243"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt244"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt245"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt246"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt247"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt248"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt249"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt250"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt251"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt252"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt253"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt254"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt255"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt256"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt257"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt258"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt259"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt260"
)

//go:embed testdata/divergence_hunt1/main.go
var divergenceHunt1Src string

//go:embed testdata/divergence_hunt2/main.go
var divergenceHunt2Src string

//go:embed testdata/divergence_hunt3/main.go
var divergenceHunt3Src string

//go:embed testdata/divergence_hunt4/main.go
var divergenceHunt4Src string

//go:embed testdata/divergence_hunt5/main.go
var divergenceHunt5Src string

//go:embed testdata/divergence_hunt6/main.go
var divergenceHunt6Src string

//go:embed testdata/divergence_hunt7/main.go
var divergenceHunt7Src string

//go:embed testdata/divergence_hunt8/main.go
var divergenceHunt8Src string

//go:embed testdata/divergence_hunt9/main.go
var divergenceHunt9Src string

//go:embed testdata/divergence_hunt10/main.go
var divergenceHunt10Src string

//go:embed testdata/divergence_hunt11/main.go
var divergenceHunt11Src string

//go:embed testdata/divergence_hunt12/main.go
var divergenceHunt12Src string

//go:embed testdata/divergence_hunt13/main.go
var divergenceHunt13Src string

//go:embed testdata/divergence_hunt14/main.go
var divergenceHunt14Src string

//go:embed testdata/divergence_hunt15/main.go
var divergenceHunt15Src string

//go:embed testdata/divergence_hunt16/main.go
var divergenceHunt16Src string

//go:embed testdata/divergence_hunt17/main.go
var divergenceHunt17Src string

//go:embed testdata/divergence_hunt18/main.go
var divergenceHunt18Src string

//go:embed testdata/divergence_hunt19/main.go
var divergenceHunt19Src string

//go:embed testdata/divergence_hunt20/main.go
var divergenceHunt20Src string

//go:embed testdata/divergence_hunt21/main.go
var divergenceHunt21Src string

//go:embed testdata/divergence_hunt22/main.go
var divergenceHunt22Src string

//go:embed testdata/divergence_hunt23/main.go
var divergenceHunt23Src string

//go:embed testdata/divergence_hunt24/main.go
var divergenceHunt24Src string

//go:embed testdata/divergence_hunt25/main.go
var divergenceHunt25Src string

//go:embed testdata/divergence_hunt26/main.go
var divergenceHunt26Src string

//go:embed testdata/divergence_hunt27/main.go
var divergenceHunt27Src string

//go:embed testdata/divergence_hunt28/main.go
var divergenceHunt28Src string

//go:embed testdata/divergence_hunt29/main.go
var divergenceHunt29Src string

//go:embed testdata/divergence_hunt30/main.go
var divergenceHunt30Src string

//go:embed testdata/divergence_hunt31/main.go
var divergenceHunt31Src string

//go:embed testdata/divergence_hunt32/main.go
var divergenceHunt32Src string

//go:embed testdata/divergence_hunt33/main.go
var divergenceHunt33Src string

//go:embed testdata/divergence_hunt34/main.go
var divergenceHunt34Src string

//go:embed testdata/divergence_hunt35/main.go
var divergenceHunt35Src string

//go:embed testdata/divergence_hunt36/main.go
var divergenceHunt36Src string

//go:embed testdata/divergence_hunt37/main.go
var divergenceHunt37Src string

//go:embed testdata/divergence_hunt38/main.go
var divergenceHunt38Src string

//go:embed testdata/divergence_hunt39/main.go
var divergenceHunt39Src string

//go:embed testdata/divergence_hunt40/main.go
var divergenceHunt40Src string

//go:embed testdata/divergence_hunt41/main.go
var divergenceHunt41Src string

//go:embed testdata/divergence_hunt42/main.go
var divergenceHunt42Src string

//go:embed testdata/divergence_hunt43/main.go
var divergenceHunt43Src string

//go:embed testdata/divergence_hunt44/main.go
var divergenceHunt44Src string

//go:embed testdata/divergence_hunt45/main.go
var divergenceHunt45Src string

//go:embed testdata/divergence_hunt46/main.go
var divergenceHunt46Src string

//go:embed testdata/divergence_hunt47/main.go
var divergenceHunt47Src string

//go:embed testdata/divergence_hunt48/main.go
var divergenceHunt48Src string

//go:embed testdata/divergence_hunt49/main.go
var divergenceHunt49Src string

//go:embed testdata/divergence_hunt50/main.go
var divergenceHunt50Src string

//go:embed testdata/divergence_hunt51/main.go
var divergenceHunt51Src string

//go:embed testdata/divergence_hunt52/main.go
var divergenceHunt52Src string

//go:embed testdata/divergence_hunt53/main.go
var divergenceHunt53Src string

//go:embed testdata/divergence_hunt54/main.go
var divergenceHunt54Src string

//go:embed testdata/divergence_hunt55/main.go
var divergenceHunt55Src string

//go:embed testdata/divergence_hunt56/main.go
var divergenceHunt56Src string

//go:embed testdata/divergence_hunt57/main.go
var divergenceHunt57Src string

//go:embed testdata/divergence_hunt58/main.go
var divergenceHunt58Src string

//go:embed testdata/divergence_hunt59/main.go
var divergenceHunt59Src string

//go:embed testdata/divergence_hunt60/main.go
var divergenceHunt60Src string

//go:embed testdata/divergence_hunt61/main.go
var divergenceHunt61Src string

//go:embed testdata/divergence_hunt62/main.go
var divergenceHunt62Src string

//go:embed testdata/divergence_hunt63/main.go
var divergenceHunt63Src string

//go:embed testdata/divergence_hunt64/main.go
var divergenceHunt64Src string

//go:embed testdata/divergence_hunt65/main.go
var divergenceHunt65Src string

//go:embed testdata/divergence_hunt66/main.go
var divergenceHunt66Src string

//go:embed testdata/divergence_hunt67/main.go
var divergenceHunt67Src string

//go:embed testdata/divergence_hunt68/main.go
var divergenceHunt68Src string

//go:embed testdata/divergence_hunt69/main.go
var divergenceHunt69Src string

//go:embed testdata/divergence_hunt70/main.go
var divergenceHunt70Src string

//go:embed testdata/divergence_hunt71/main.go
var divergenceHunt71Src string

//go:embed testdata/divergence_hunt72/main.go
var divergenceHunt72Src string

//go:embed testdata/divergence_hunt73/main.go
var divergenceHunt73Src string

//go:embed testdata/divergence_hunt74/main.go
var divergenceHunt74Src string

//go:embed testdata/divergence_hunt75/main.go
var divergenceHunt75Src string

//go:embed testdata/divergence_hunt76/main.go
var divergenceHunt76Src string

//go:embed testdata/divergence_hunt77/main.go
var divergenceHunt77Src string

//go:embed testdata/divergence_hunt78/main.go
var divergenceHunt78Src string

//go:embed testdata/divergence_hunt79/main.go
var divergenceHunt79Src string

//go:embed testdata/divergence_hunt80/main.go
var divergenceHunt80Src string

//go:embed testdata/divergence_hunt81/main.go
var divergenceHunt81Src string

//go:embed testdata/divergence_hunt82/main.go
var divergenceHunt82Src string

//go:embed testdata/divergence_hunt83/main.go
var divergenceHunt83Src string

//go:embed testdata/divergence_hunt84/main.go
var divergenceHunt84Src string

//go:embed testdata/divergence_hunt85/main.go
var divergenceHunt85Src string

//go:embed testdata/divergence_hunt86/main.go
var divergenceHunt86Src string

//go:embed testdata/divergence_hunt87/main.go
var divergenceHunt87Src string

//go:embed testdata/divergence_hunt88/main.go
var divergenceHunt88Src string

//go:embed testdata/divergence_hunt89/main.go
var divergenceHunt89Src string

//go:embed testdata/divergence_hunt90/main.go
var divergenceHunt90Src string

//go:embed testdata/divergence_hunt91/main.go
var divergenceHunt91Src string

//go:embed testdata/divergence_hunt92/main.go
var divergenceHunt92Src string

//go:embed testdata/divergence_hunt93/main.go
var divergenceHunt93Src string

//go:embed testdata/divergence_hunt94/main.go
var divergenceHunt94Src string

//go:embed testdata/divergence_hunt95/main.go
var divergenceHunt95Src string

//go:embed testdata/divergence_hunt96/main.go
var divergenceHunt96Src string

//go:embed testdata/divergence_hunt97/main.go
var divergenceHunt97Src string

//go:embed testdata/divergence_hunt98/main.go
var divergenceHunt98Src string

//go:embed testdata/divergence_hunt99/main.go
var divergenceHunt99Src string

//go:embed testdata/divergence_hunt100/main.go
var divergenceHunt100Src string

//go:embed testdata/divergence_hunt101/main.go
var divergenceHunt101Src string

//go:embed testdata/divergence_hunt102/main.go
var divergenceHunt102Src string

//go:embed testdata/divergence_hunt103/main.go
var divergenceHunt103Src string

//go:embed testdata/divergence_hunt104/main.go
var divergenceHunt104Src string

//go:embed testdata/divergence_hunt105/main.go
var divergenceHunt105Src string

//go:embed testdata/divergence_hunt106/main.go
var divergenceHunt106Src string

//go:embed testdata/divergence_hunt107/main.go
var divergenceHunt107Src string

//go:embed testdata/divergence_hunt108/main.go
var divergenceHunt108Src string

//go:embed testdata/divergence_hunt109/main.go
var divergenceHunt109Src string

//go:embed testdata/divergence_hunt110/main.go
var divergenceHunt110Src string

//go:embed testdata/divergence_hunt111/main.go
var divergenceHunt111Src string

//go:embed testdata/divergence_hunt112/main.go
var divergenceHunt112Src string

//go:embed testdata/divergence_hunt113/main.go
var divergenceHunt113Src string

//go:embed testdata/divergence_hunt114/main.go
var divergenceHunt114Src string

//go:embed testdata/divergence_hunt115/main.go
var divergenceHunt115Src string

//go:embed testdata/divergence_hunt116/main.go
var divergenceHunt116Src string

//go:embed testdata/divergence_hunt117/main.go
var divergenceHunt117Src string

//go:embed testdata/divergence_hunt118/main.go
var divergenceHunt118Src string

//go:embed testdata/divergence_hunt119/main.go
var divergenceHunt119Src string

//go:embed testdata/divergence_hunt120/main.go
var divergenceHunt120Src string

//go:embed testdata/divergence_hunt121/main.go
var divergenceHunt121Src string

//go:embed testdata/divergence_hunt122/main.go
var divergenceHunt122Src string

//go:embed testdata/divergence_hunt123/main.go
var divergenceHunt123Src string

//go:embed testdata/divergence_hunt124/main.go
var divergenceHunt124Src string

//go:embed testdata/divergence_hunt125/main.go
var divergenceHunt125Src string

//go:embed testdata/divergence_hunt126/main.go
var divergenceHunt126Src string

//go:embed testdata/divergence_hunt127/main.go
var divergenceHunt127Src string

//go:embed testdata/divergence_hunt128/main.go
var divergenceHunt128Src string

//go:embed testdata/divergence_hunt129/main.go
var divergenceHunt129Src string

//go:embed testdata/divergence_hunt130/main.go
var divergenceHunt130Src string

//go:embed testdata/divergence_hunt131/main.go
var divergenceHunt131Src string

//go:embed testdata/divergence_hunt132/main.go
var divergenceHunt132Src string

//go:embed testdata/divergence_hunt133/main.go
var divergenceHunt133Src string

//go:embed testdata/divergence_hunt134/main.go
var divergenceHunt134Src string

//go:embed testdata/divergence_hunt135/main.go
var divergenceHunt135Src string

//go:embed testdata/divergence_hunt136/main.go
var divergenceHunt136Src string

//go:embed testdata/divergence_hunt137/main.go
var divergenceHunt137Src string

//go:embed testdata/divergence_hunt138/main.go
var divergenceHunt138Src string

//go:embed testdata/divergence_hunt139/main.go
var divergenceHunt139Src string

//go:embed testdata/divergence_hunt140/main.go
var divergenceHunt140Src string

//go:embed testdata/divergence_hunt141/main.go
var divergenceHunt141Src string

//go:embed testdata/divergence_hunt142/main.go
var divergenceHunt142Src string

//go:embed testdata/divergence_hunt143/main.go
var divergenceHunt143Src string

//go:embed testdata/divergence_hunt144/main.go
var divergenceHunt144Src string

//go:embed testdata/divergence_hunt145/main.go
var divergenceHunt145Src string

//go:embed testdata/divergence_hunt146/main.go
var divergenceHunt146Src string

//go:embed testdata/divergence_hunt147/main.go
var divergenceHunt147Src string

//go:embed testdata/divergence_hunt148/main.go
var divergenceHunt148Src string

//go:embed testdata/divergence_hunt149/main.go
var divergenceHunt149Src string

//go:embed testdata/divergence_hunt150/main.go
var divergenceHunt150Src string

//go:embed testdata/divergence_hunt151/main.go
var divergenceHunt151Src string

//go:embed testdata/divergence_hunt152/main.go
var divergenceHunt152Src string

//go:embed testdata/divergence_hunt153/main.go
var divergenceHunt153Src string

//go:embed testdata/divergence_hunt154/main.go
var divergenceHunt154Src string

//go:embed testdata/divergence_hunt155/main.go
var divergenceHunt155Src string

//go:embed testdata/divergence_hunt156/main.go
var divergenceHunt156Src string

//go:embed testdata/divergence_hunt157/main.go
var divergenceHunt157Src string

//go:embed testdata/divergence_hunt158/main.go
var divergenceHunt158Src string

//go:embed testdata/divergence_hunt159/main.go
var divergenceHunt159Src string

//go:embed testdata/divergence_hunt160/main.go
var divergenceHunt160Src string
//go:embed testdata/divergence_hunt161/main.go
var divergenceHunt161Src string

//go:embed testdata/divergence_hunt162/main.go
var divergenceHunt162Src string

//go:embed testdata/divergence_hunt163/main.go
var divergenceHunt163Src string

//go:embed testdata/divergence_hunt164/main.go
var divergenceHunt164Src string

//go:embed testdata/divergence_hunt165/main.go
var divergenceHunt165Src string

//go:embed testdata/divergence_hunt166/main.go
var divergenceHunt166Src string

//go:embed testdata/divergence_hunt167/main.go
var divergenceHunt167Src string

//go:embed testdata/divergence_hunt168/main.go
var divergenceHunt168Src string

//go:embed testdata/divergence_hunt169/main.go
var divergenceHunt169Src string

//go:embed testdata/divergence_hunt170/main.go
var divergenceHunt170Src string

//go:embed testdata/divergence_hunt171/main.go
var divergenceHunt171Src string

//go:embed testdata/divergence_hunt172/main.go
var divergenceHunt172Src string

//go:embed testdata/divergence_hunt173/main.go
var divergenceHunt173Src string

//go:embed testdata/divergence_hunt174/main.go
var divergenceHunt174Src string

//go:embed testdata/divergence_hunt175/main.go
var divergenceHunt175Src string

//go:embed testdata/divergence_hunt176/main.go
var divergenceHunt176Src string

//go:embed testdata/divergence_hunt177/main.go
var divergenceHunt177Src string

//go:embed testdata/divergence_hunt178/main.go
var divergenceHunt178Src string

//go:embed testdata/divergence_hunt179/main.go
var divergenceHunt179Src string

//go:embed testdata/divergence_hunt180/main.go
var divergenceHunt180Src string

//go:embed testdata/divergence_hunt181/main.go
var divergenceHunt181Src string

//go:embed testdata/divergence_hunt182/main.go
var divergenceHunt182Src string

//go:embed testdata/divergence_hunt183/main.go
var divergenceHunt183Src string

//go:embed testdata/divergence_hunt184/main.go
var divergenceHunt184Src string

//go:embed testdata/divergence_hunt185/main.go
var divergenceHunt185Src string

//go:embed testdata/divergence_hunt186/main.go
var divergenceHunt186Src string

//go:embed testdata/divergence_hunt187/main.go
var divergenceHunt187Src string

//go:embed testdata/divergence_hunt188/main.go
var divergenceHunt188Src string

//go:embed testdata/divergence_hunt189/main.go
var divergenceHunt189Src string

//go:embed testdata/divergence_hunt190/main.go
var divergenceHunt190Src string

//go:embed testdata/divergence_hunt191/main.go
var divergenceHunt191Src string

//go:embed testdata/divergence_hunt192/main.go
var divergenceHunt192Src string

//go:embed testdata/divergence_hunt193/main.go
var divergenceHunt193Src string

//go:embed testdata/divergence_hunt194/main.go
var divergenceHunt194Src string

//go:embed testdata/divergence_hunt195/main.go
var divergenceHunt195Src string

//go:embed testdata/divergence_hunt196/main.go
var divergenceHunt196Src string

//go:embed testdata/divergence_hunt197/main.go
var divergenceHunt197Src string

//go:embed testdata/divergence_hunt198/main.go
var divergenceHunt198Src string

//go:embed testdata/divergence_hunt199/main.go
var divergenceHunt199Src string

//go:embed testdata/divergence_hunt200/main.go
var divergenceHunt200Src string

//go:embed testdata/divergence_hunt201/main.go
var divergenceHunt201Src string

//go:embed testdata/divergence_hunt202/main.go
var divergenceHunt202Src string

//go:embed testdata/divergence_hunt203/main.go
var divergenceHunt203Src string

//go:embed testdata/divergence_hunt204/main.go
var divergenceHunt204Src string

//go:embed testdata/divergence_hunt205/main.go
var divergenceHunt205Src string

//go:embed testdata/divergence_hunt206/main.go
var divergenceHunt206Src string

//go:embed testdata/divergence_hunt207/main.go
var divergenceHunt207Src string

//go:embed testdata/divergence_hunt208/main.go
var divergenceHunt208Src string

//go:embed testdata/divergence_hunt209/main.go
var divergenceHunt209Src string

//go:embed testdata/divergence_hunt210/main.go
var divergenceHunt210Src string

//go:embed testdata/divergence_hunt211/main.go
var divergenceHunt211Src string

//go:embed testdata/divergence_hunt212/main.go
var divergenceHunt212Src string

//go:embed testdata/divergence_hunt213/main.go
var divergenceHunt213Src string

//go:embed testdata/divergence_hunt214/main.go
var divergenceHunt214Src string

//go:embed testdata/divergence_hunt215/main.go
var divergenceHunt215Src string

//go:embed testdata/divergence_hunt216/main.go
var divergenceHunt216Src string

//go:embed testdata/divergence_hunt217/main.go
var divergenceHunt217Src string

//go:embed testdata/divergence_hunt218/main.go
var divergenceHunt218Src string

//go:embed testdata/divergence_hunt219/main.go
var divergenceHunt219Src string

//go:embed testdata/divergence_hunt220/main.go
var divergenceHunt220Src string

//go:embed testdata/divergence_hunt221/main.go
var divergenceHunt221Src string

//go:embed testdata/divergence_hunt222/main.go
var divergenceHunt222Src string

//go:embed testdata/divergence_hunt223/main.go
var divergenceHunt223Src string

//go:embed testdata/divergence_hunt224/main.go
var divergenceHunt224Src string

//go:embed testdata/divergence_hunt225/main.go
var divergenceHunt225Src string

//go:embed testdata/divergence_hunt226/main.go
var divergenceHunt226Src string

//go:embed testdata/divergence_hunt227/main.go
var divergenceHunt227Src string

//go:embed testdata/divergence_hunt228/main.go
var divergenceHunt228Src string

//go:embed testdata/divergence_hunt229/main.go
var divergenceHunt229Src string

//go:embed testdata/divergence_hunt230/main.go
var divergenceHunt230Src string

//go:embed testdata/divergence_hunt231/main.go
var divergenceHunt231Src string

//go:embed testdata/divergence_hunt232/main.go
var divergenceHunt232Src string

//go:embed testdata/divergence_hunt233/main.go
var divergenceHunt233Src string

//go:embed testdata/divergence_hunt234/main.go
var divergenceHunt234Src string

//go:embed testdata/divergence_hunt235/main.go
var divergenceHunt235Src string

//go:embed testdata/divergence_hunt236/main.go
var divergenceHunt236Src string

//go:embed testdata/divergence_hunt237/main.go
var divergenceHunt237Src string

//go:embed testdata/divergence_hunt238/main.go
var divergenceHunt238Src string

//go:embed testdata/divergence_hunt239/main.go
var divergenceHunt239Src string

//go:embed testdata/divergence_hunt240/main.go
var divergenceHunt240Src string

//go:embed testdata/divergence_hunt241/main.go
var divergenceHunt241Src string

//go:embed testdata/divergence_hunt242/main.go
var divergenceHunt242Src string

//go:embed testdata/divergence_hunt243/main.go
var divergenceHunt243Src string

//go:embed testdata/divergence_hunt244/main.go
var divergenceHunt244Src string

//go:embed testdata/divergence_hunt245/main.go
var divergenceHunt245Src string

//go:embed testdata/divergence_hunt246/main.go
var divergenceHunt246Src string

//go:embed testdata/divergence_hunt247/main.go
var divergenceHunt247Src string

//go:embed testdata/divergence_hunt248/main.go
var divergenceHunt248Src string

//go:embed testdata/divergence_hunt249/main.go
var divergenceHunt249Src string

//go:embed testdata/divergence_hunt250/main.go
var divergenceHunt250Src string

//go:embed testdata/divergence_hunt251/main.go
var divergenceHunt251Src string

//go:embed testdata/divergence_hunt252/main.go
var divergenceHunt252Src string

//go:embed testdata/divergence_hunt253/main.go
var divergenceHunt253Src string

//go:embed testdata/divergence_hunt254/main.go
var divergenceHunt254Src string

//go:embed testdata/divergence_hunt255/main.go
var divergenceHunt255Src string

//go:embed testdata/divergence_hunt256/main.go
var divergenceHunt256Src string

//go:embed testdata/divergence_hunt257/main.go
var divergenceHunt257Src string

//go:embed testdata/divergence_hunt258/main.go
var divergenceHunt258Src string

//go:embed testdata/divergence_hunt259/main.go
var divergenceHunt259Src string

//go:embed testdata/divergence_hunt260/main.go
var divergenceHunt260Src string

// divergenceTestCase is like testCase but with explicit expected value.
// This is used for divergence hunting where we compare interpreter output
// against native Go execution.
type divergenceTestCase struct {
	funcName    string
	args        []any
	native      any  // native function, called via reflection
	knownIssue  bool // if true, test is skipped (known bug tracked for fixing)
}

// divergenceTestSet is a set of divergence test cases sharing one source file.
type divergenceTestSet struct {
	src       string
	tests     map[string]divergenceTestCase
	buildOpts []gig.BuildOption
}

// perTestTimeout is the maximum time a single divergence sub-test is allowed to
// run before it's declared hung. This is a belt-and-braces watchdog: gig.Program
// already has its own ctx-based timeout, but tests that block in native code
// (e.g. a deadlocked sync.Mutex.Lock) can't be preempted by ctx cancellation.
// When that happens we fail the sub-test fast instead of hanging the whole suite.
const perTestTimeout = 8 * time.Second

// runProgramWithWatchdog calls prog.Run with a watchdog goroutine so a hung
// interpreter call fails the sub-test within perTestTimeout instead of blocking
// the whole go test run. The underlying goroutine is leaked on timeout (we can't
// force-kill a goroutine stuck in native code) but the test process is short-lived
// so it's acceptable.
func runProgramWithWatchdog(prog *gig.Program, funcName string, args ...any) (any, error, bool) {
	type result struct {
		val any
		err error
	}
	ch := make(chan result, 1)
	ctx, cancel := context.WithTimeout(context.Background(), perTestTimeout)
	defer cancel()
	go func() {
		val, err := prog.RunWithContext(ctx, funcName, args...)
		ch <- result{val, err}
	}()
	select {
	case r := <-ch:
		return r.val, r.err, false
	case <-ctx.Done():
		return nil, nil, true
	}
}

// runDivergenceTestSet compiles the source once and runs each test,
// comparing interpreter output with native Go execution.
func runDivergenceTestSet(t *testing.T, set divergenceTestSet) {
	t.Helper()
	t.Parallel()
	// If every test in the set is marked knownIssue, skip the entire set so
	// that a Build error from the side effect of multiple init() or similar
	// doesn't fail the test before sub-test skip flags are honored.
	if len(set.tests) > 0 {
		allKnown := true
		for _, tc := range set.tests {
			if !tc.knownIssue {
				allKnown = false
				break
			}
		}
		if allKnown {
			t.Skip("All tests in this set are marked knownIssue — skipping Build to avoid init-time panics")
			return
		}
	}

	prog, err := gig.Build(set.src, set.buildOpts...)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	defer prog.Close()

	for name, tc := range set.tests {
		t.Run(name, func(t *testing.T) {
			if tc.knownIssue {
				t.Skip("KNOWN ISSUE - pending fix")
			}

			// Run interpreter with per-test watchdog
			interpResult, interpErr, timedOut := runProgramWithWatchdog(prog, tc.funcName, tc.args...)
			if timedOut {
				t.Errorf("DIVERGENCE (timeout): interpreter hung for >%s running %s", perTestTimeout, tc.funcName)
				return
			}
			if interpErr != nil {
				t.Errorf("DIVERGENCE (error): %v", interpErr)
				return
			}

			// Run native
			if tc.native == nil {
				t.Fatalf("native function is nil for %s", name)
			}
			nativeResult := callNative(tc.native, tc.args)

			// Compare
			if !reflect.DeepEqual(interpResult, nativeResult) {
				t.Errorf("DIVERGENCE (mismatch): interp=%v (%T), native=%v (%T)",
					interpResult, interpResult, nativeResult, nativeResult)
			}
		})
	}
}

func TestDivergenceHunt1(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt1Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"NilSliceCompare":    {funcName: "NilSliceCompare", native: divergence_hunt1.NilSliceCompare},
			"NilMapCompare":      {funcName: "NilMapCompare", native: divergence_hunt1.NilMapCompare},
			"NilChanCompare":     {funcName: "NilChanCompare", native: divergence_hunt1.NilChanCompare},
			"ComplexArith":       {funcName: "ComplexArith", native: divergence_hunt1.ComplexArith},
			"StringIndexByte":    {funcName: "StringIndexByte", native: divergence_hunt1.StringIndexByte},
			"IntOverflow":        {funcName: "IntOverflow", native: divergence_hunt1.IntOverflow},
			"DeferModify":        {funcName: "DeferModify", native: divergence_hunt1.DeferModify},
			"TypeAssertPanic":    {funcName: "TypeAssertPanic", native: divergence_hunt1.TypeAssertPanic},
			"Complex64Arith":     {funcName: "Complex64Arith", native: divergence_hunt1.Complex64Arith},
			"SliceBoundsPanic":   {funcName: "SliceBoundsPanic", native: divergence_hunt1.SliceBoundsPanic},
			"NilPointerDeref":    {funcName: "NilPointerDeref", native: divergence_hunt1.NilPointerDeref},
			"NilMapWrite":        {funcName: "NilMapWrite", native: divergence_hunt1.NilMapWrite},
			"DivZeroPanicTest":   {funcName: "DivZeroPanicTest", native: divergence_hunt1.DivZeroPanicTest},
			"UintOverflow":       {funcName: "UintOverflow", native: divergence_hunt1.UintOverflow},
			"Int8Negative":       {funcName: "Int8Negative", native: divergence_hunt1.Int8Negative},
			"NaNCompare":         {funcName: "NaNCompare", native: divergence_hunt1.NaNCompare},
			"MapNilLookup":       {funcName: "MapNilLookup", native: divergence_hunt1.MapNilLookup},
			"SliceCopy":          {funcName: "SliceCopy", native: divergence_hunt1.SliceCopy},
			"RuneLiteral":        {funcName: "RuneLiteral", native: divergence_hunt1.RuneLiteral},
			"NilInterfaceAssert": {funcName: "NilInterfaceAssert", native: divergence_hunt1.NilInterfaceAssert},
			"SortInts":           {funcName: "SortInts", native: divergence_hunt1.SortInts},
			"StringsJoin":        {funcName: "StringsJoin", native: divergence_hunt1.StringsJoin},
			"StringsSplit":       {funcName: "StringsSplit", native: divergence_hunt1.StringsSplit},
			"StringsContains":    {funcName: "StringsContains", native: divergence_hunt1.StringsContains},
			"StrconvRoundTrip":   {funcName: "StrconvRoundTrip", native: divergence_hunt1.StrconvRoundTrip},
			"FmtSprintf":         {funcName: "FmtSprintf", native: divergence_hunt1.FmtSprintf},
			"PanicInDefer":       {funcName: "PanicInDefer", native: divergence_hunt1.PanicInDefer},
			"MultipleRecoverCalls": {funcName: "MultipleRecoverCalls", native: divergence_hunt1.MultipleRecoverCalls},
			"BoolToStrconv":      {funcName: "BoolToStrconv", native: divergence_hunt1.BoolToStrconv},
			"FloatToStrconv":     {funcName: "FloatToStrconv", native: divergence_hunt1.FloatToStrconv},
			"StringsReplace":     {funcName: "StringsReplace", native: divergence_hunt1.StringsReplace},
			"StringsHasPrefix":   {funcName: "StringsHasPrefix", native: divergence_hunt1.StringsHasPrefix},
			"StringsTrim":        {funcName: "StringsTrim", native: divergence_hunt1.StringsTrim},
			"MapIntKey":          {funcName: "MapIntKey", native: divergence_hunt1.MapIntKey},
			"CapSlice":           {funcName: "CapSlice", native: divergence_hunt1.CapSlice},
			"ByteSliceIndex":     {funcName: "ByteSliceIndex", native: divergence_hunt1.ByteSliceIndex},
			"DeferMultipleOrder": {funcName: "DeferMultipleOrder", native: divergence_hunt1.DeferMultipleOrder},
			"ErrorTypeAssertion": {funcName: "ErrorTypeAssertion", native: divergence_hunt1.ErrorTypeAssertion},
			"RecursiveFactorial": {funcName: "RecursiveFactorial", native: divergence_hunt1.RecursiveFactorial},
			"ClosureCounter":     {funcName: "ClosureCounter", native: divergence_hunt1.ClosureCounter},
			"BitwiseAnd":         {funcName: "BitwiseAnd", native: divergence_hunt1.BitwiseAnd},
			"BitwiseOr":          {funcName: "BitwiseOr", native: divergence_hunt1.BitwiseOr},
			"BitwiseXor":         {funcName: "BitwiseXor", native: divergence_hunt1.BitwiseXor},
			"BitwiseShift":       {funcName: "BitwiseShift", native: divergence_hunt1.BitwiseShift},
			"Float64Arith":       {funcName: "Float64Arith", native: divergence_hunt1.Float64Arith},
			"PanicIntValue":      {funcName: "PanicIntValue", native: divergence_hunt1.PanicIntValue},
			"DoublePanic":        {funcName: "DoublePanic", native: divergence_hunt1.DoublePanic},
			"DeferModifyAfterPanic": {funcName: "DeferModifyAfterPanic", native: divergence_hunt1.DeferModifyAfterPanic},
			"SliceOfStructs":     {funcName: "SliceOfStructs", native: divergence_hunt1.SliceOfStructs},
			"ForBreak":           {funcName: "ForBreak", native: divergence_hunt1.ForBreak},
			"NestedLoop":         {funcName: "NestedLoop", native: divergence_hunt1.NestedLoop},
			"StringCompareOps":   {funcName: "StringCompareOps", native: divergence_hunt1.StringCompareOps},
			"MapCommaOkMissing":  {funcName: "MapCommaOkMissing", native: divergence_hunt1.MapCommaOkMissing},
			"SwitchDefault":      {funcName: "SwitchDefault", native: divergence_hunt1.SwitchDefault},
			"VariadicFunc":       {funcName: "VariadicFunc", native: divergence_hunt1.VariadicFunc},
			"TypeSwitch":         {funcName: "TypeSwitch", native: divergence_hunt1.TypeSwitch},
			"StructEmbedding":    {funcName: "StructEmbedding", native: divergence_hunt1.StructEmbedding},
			"ChannelBuffered":    {funcName: "ChannelBuffered", native: divergence_hunt1.ChannelBuffered},
		},
	})
}

func TestDivergenceHunt2(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt2Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"MapLen":                    {funcName: "MapLen", native: divergence_hunt2.MapLen},
			"MapDelete":                 {funcName: "MapDelete", native: divergence_hunt2.MapDelete},
			"MapOverwrite":              {funcName: "MapOverwrite", native: divergence_hunt2.MapOverwrite},
			"SliceNilAppend":            {funcName: "SliceNilAppend", native: divergence_hunt2.SliceNilAppend},
			"SliceGrow":                 {funcName: "SliceGrow", native: divergence_hunt2.SliceGrow},
			"StringLen":                 {funcName: "StringLen", native: divergence_hunt2.StringLen},
			"StringConcat":              {funcName: "StringConcat", native: divergence_hunt2.StringConcat},
			"IntConversion":             {funcName: "IntConversion", native: divergence_hunt2.IntConversion},
			"UintConversion":            {funcName: "UintConversion", native: divergence_hunt2.UintConversion},
			"MultiReturnSwap":           {funcName: "MultiReturnSwap", native: divergence_hunt2.MultiReturnSwap},
			"BlankIdentifier":           {funcName: "BlankIdentifier", native: divergence_hunt2.BlankIdentifier},
			"NilSliceLen":               {funcName: "NilSliceLen", native: divergence_hunt2.NilSliceLen},
			"NilMapLen":                 {funcName: "NilMapLen", native: divergence_hunt2.NilMapLen},
			"PointerDeref":              {funcName: "PointerDeref", native: divergence_hunt2.PointerDeref},
			"PointerAssign":             {funcName: "PointerAssign", native: divergence_hunt2.PointerAssign},
			"SliceOfPointers":           {funcName: "SliceOfPointers", native: divergence_hunt2.SliceOfPointers},
			"MapIteration":              {funcName: "MapIteration", native: divergence_hunt2.MapIteration},
			"StringRange":               {funcName: "StringRange", native: divergence_hunt2.StringRange},
			"FloatConversion":           {funcName: "FloatConversion", native: divergence_hunt2.FloatConversion},
			"ByteSliceAppend":           {funcName: "ByteSliceAppend", native: divergence_hunt2.ByteSliceAppend},
			"ByteSliceWrite":            {funcName: "ByteSliceWrite", native: divergence_hunt2.ByteSliceWrite},
			"StructCompare":             {funcName: "StructCompare", native: divergence_hunt2.StructCompare},
			"ArrayLen":                  {funcName: "ArrayLen", native: divergence_hunt2.ArrayLen},
			"ArrayValue":                {funcName: "ArrayValue", native: divergence_hunt2.ArrayValue},
			"StringIndexOutOfRange":     {funcName: "StringIndexOutOfRange", native: divergence_hunt2.StringIndexOutOfRange},
			"MapKeyIntFloat":            {funcName: "MapKeyIntFloat", native: divergence_hunt2.MapKeyIntFloat},
			"ShortVarDecl":              {funcName: "ShortVarDecl", native: divergence_hunt2.ShortVarDecl},
			"MultipleShortVar":          {funcName: "MultipleShortVar", native: divergence_hunt2.MultipleShortVar},
			"SliceThreeIndex":           {funcName: "SliceThreeIndex", native: divergence_hunt2.SliceThreeIndex},
			"NilFuncCall":               {funcName: "NilFuncCall", native: divergence_hunt2.NilFuncCall},
			"StringByteSliceConversion": {funcName: "StringByteSliceConversion", native: divergence_hunt2.StringByteSliceConversion},
		},
	})
}

func TestDivergenceHunt3(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt3Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"StringBuilder":     {funcName: "StringBuilder", native: divergence_hunt3.StringBuilder},
			"ConstBlock":        {funcName: "ConstBlock", native: divergence_hunt3.ConstBlock},
			"IotaEnum":          {funcName: "IotaEnum", native: divergence_hunt3.IotaEnum},
			"MultipleAssign":    {funcName: "MultipleAssign", native: divergence_hunt3.MultipleAssign},
			"NestedMap":         {funcName: "NestedMap", native: divergence_hunt3.NestedMap},
			"RuneIteration":     {funcName: "RuneIteration", native: divergence_hunt3.RuneIteration},
			"StringIndexRune":   {funcName: "StringIndexRune", native: divergence_hunt3.StringIndexRune},
			"StringCount":       {funcName: "StringCount", native: divergence_hunt3.StringCount},
			"MapBoolKey":        {funcName: "MapBoolKey", native: divergence_hunt3.MapBoolKey},
			"SliceReverse":      {funcName: "SliceReverse", native: divergence_hunt3.SliceReverse},
			"StructMethod":      {funcName: "StructMethod", native: divergence_hunt3.StructMethod},
			"InterfaceEmpty":    {funcName: "InterfaceEmpty", native: divergence_hunt3.InterfaceEmpty},
			"InterfaceNil":      {funcName: "InterfaceNil", native: divergence_hunt3.InterfaceNil},
			"SliceOfInterface":  {funcName: "SliceOfInterface", native: divergence_hunt3.SliceOfInterface},
			"MapWithStructValue": {funcName: "MapWithStructValue", native: divergence_hunt3.MapWithStructValue},
			"StringFields":      {funcName: "StringFields", native: divergence_hunt3.StringFields},
			"StringRepeat":      {funcName: "StringRepeat", native: divergence_hunt3.StringRepeat},
			"StringMap":         {funcName: "StringMap", native: divergence_hunt3.StringMap},
			"MapStructKey":      {funcName: "MapStructKey", native: divergence_hunt3.MapStructKey},
			"SliceMinMax":       {funcName: "SliceMinMax", native: divergence_hunt3.SliceMinMax},
			"NestedIf":          {funcName: "NestedIf", native: divergence_hunt3.NestedIf},
			"StringToLower":     {funcName: "StringToLower", native: divergence_hunt3.StringToLower},
			"StringToUpper":     {funcName: "StringToUpper", native: divergence_hunt3.StringToUpper},
			"ContinueLoop":      {funcName: "ContinueLoop", native: divergence_hunt3.ContinueLoop},
			"LabeledBreak":      {funcName: "LabeledBreak", native: divergence_hunt3.LabeledBreak},
			"SliceMakeZero":     {funcName: "SliceMakeZero", native: divergence_hunt3.SliceMakeZero},
			"ArrayIteration":    {funcName: "ArrayIteration", native: divergence_hunt3.ArrayIteration},
			"Float32Arith":      {funcName: "Float32Arith", native: divergence_hunt3.Float32Arith},
			"Int8Arith":         {funcName: "Int8Arith", native: divergence_hunt3.Int8Arith},
			"Uint16Arith":       {funcName: "Uint16Arith", native: divergence_hunt3.Uint16Arith},
		},
	})
}

func TestDivergenceHunt4(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt4Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"Float64NaN":         {funcName: "Float64NaN", native: divergence_hunt4.Float64NaN},
			"Float64Inf":         {funcName: "Float64Inf", native: divergence_hunt4.Float64Inf},
			"Float64NegZero":     {funcName: "Float64NegZero", native: divergence_hunt4.Float64NegZero},
			"Int16Conversion":    {funcName: "Int16Conversion", native: divergence_hunt4.Int16Conversion},
			"Uint32Conversion":   {funcName: "Uint32Conversion", native: divergence_hunt4.Uint32Conversion},
			"FloatToIntTruncation": {funcName: "FloatToIntTruncation", native: divergence_hunt4.FloatToIntTruncation},
			"NegativeFloatToInt": {funcName: "NegativeFloatToInt", native: divergence_hunt4.NegativeFloatToInt},
			"StrconvAtoi":        {funcName: "StrconvAtoi", native: divergence_hunt4.StrconvAtoi},
			"StrconvItoa":        {funcName: "StrconvItoa", native: divergence_hunt4.StrconvItoa},
			"StrconvFormatInt":   {funcName: "StrconvFormatInt", native: divergence_hunt4.StrconvFormatInt},
			"StrconvParseFloat":  {funcName: "StrconvParseFloat", native: divergence_hunt4.StrconvParseFloat},
			"MathAbs":            {funcName: "MathAbs", native: divergence_hunt4.MathAbs},
			"MathMax":            {funcName: "MathMax", native: divergence_hunt4.MathMax},
			"MathMin":            {funcName: "MathMin", native: divergence_hunt4.MathMin},
			"MathPow":            {funcName: "MathPow", native: divergence_hunt4.MathPow},
			"MathSqrt":           {funcName: "MathSqrt", native: divergence_hunt4.MathSqrt},
			"MathCeil":           {funcName: "MathCeil", native: divergence_hunt4.MathCeil},
			"MathFloor":          {funcName: "MathFloor", native: divergence_hunt4.MathFloor},
			"IntMin":             {funcName: "IntMin", native: divergence_hunt4.IntMin},
			"IntMax":             {funcName: "IntMax", native: divergence_hunt4.IntMax},
			"UintptrSize":        {funcName: "UintptrSize", native: divergence_hunt4.UintptrSize},
			"ByteArith":          {funcName: "ByteArith", native: divergence_hunt4.ByteArith},
			"Int32Overflow":      {funcName: "Int32Overflow", native: divergence_hunt4.Int32Overflow},
			"Uint8Wrap":          {funcName: "Uint8Wrap", native: divergence_hunt4.Uint8Wrap},
			"ComplexConj":        {funcName: "ComplexConj", native: divergence_hunt4.ComplexConj},
			"Float32Precision":   {funcName: "Float32Precision", native: divergence_hunt4.Float32Precision},
			"MapLenAfterDelete":  {funcName: "MapLenAfterDelete", native: divergence_hunt4.MapLenAfterDelete},
			"SliceCapAfterAppend": {funcName: "SliceCapAfterAppend", native: divergence_hunt4.SliceCapAfterAppend},
			"StringFromRunes":    {funcName: "StringFromRunes", native: divergence_hunt4.StringFromRunes},
			"RuneToInt":          {funcName: "RuneToInt", native: divergence_hunt4.RuneToInt},
			"BoolToInt":          {funcName: "BoolToInt", native: divergence_hunt4.BoolToInt},
		},
	})
}

func TestDivergenceHunt5(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt5Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"ErrorNew":            {funcName: "ErrorNew", native: divergence_hunt5.ErrorNew},
			"Errorf":              {funcName: "Errorf", native: divergence_hunt5.Errorf},
			"FmtPrintln":          {funcName: "FmtPrintln", native: divergence_hunt5.FmtPrintln},
			"FmtIntWidth":         {funcName: "FmtIntWidth", native: divergence_hunt5.FmtIntWidth},
			"FmtFloat":            {funcName: "FmtFloat", native: divergence_hunt5.FmtFloat},
			"FmtBool":             {funcName: "FmtBool", native: divergence_hunt5.FmtBool},
			"FmtHex":              {funcName: "FmtHex", native: divergence_hunt5.FmtHex},
			"FmtOctal":            {funcName: "FmtOctal", native: divergence_hunt5.FmtOctal},
			"FmtBinary":           {funcName: "FmtBinary", native: divergence_hunt5.FmtBinary},
			"FmtChar":             {funcName: "FmtChar", native: divergence_hunt5.FmtChar},
			"FmtStringWidth":      {funcName: "FmtStringWidth", native: divergence_hunt5.FmtStringWidth},
			"SliceFilter":         {funcName: "SliceFilter", native: divergence_hunt5.SliceFilter},
			"SliceMap":            {funcName: "SliceMap", native: divergence_hunt5.SliceMap},
			"ClosureSum":          {funcName: "ClosureSum", native: divergence_hunt5.ClosureSum},
			"ClosureCapture":      {funcName: "ClosureCapture", native: divergence_hunt5.ClosureCapture},
			"InterfaceSlice":      {funcName: "InterfaceSlice", native: divergence_hunt5.InterfaceSlice},
			"MultipleReturnIgnore": {funcName: "MultipleReturnIgnore", native: divergence_hunt5.MultipleReturnIgnore},
			"NamedReturn":         {funcName: "NamedReturn", native: divergence_hunt5.NamedReturn},
			"NamedReturnBare":     {funcName: "NamedReturnBare", native: divergence_hunt5.NamedReturnBare},
			"StringJoinInts":      {funcName: "StringJoinInts", native: divergence_hunt5.StringJoinInts},
			"MapStringSlice":      {funcName: "MapStringSlice", native: divergence_hunt5.MapStringSlice},
			"NestedStruct":        {funcName: "NestedStruct", native: divergence_hunt5.NestedStruct},
			"StructLiteral":       {funcName: "StructLiteral", native: divergence_hunt5.StructLiteral},
			"StructPointer":       {funcName: "StructPointer", native: divergence_hunt5.StructPointer},
			"DeferReturn":         {funcName: "DeferReturn", native: divergence_hunt5.DeferReturn},
			"DeferClosure":        {funcName: "DeferClosure", native: divergence_hunt5.DeferClosure},
			"StringEqual":         {funcName: "StringEqual", native: divergence_hunt5.StringEqual},
			"MapLookup":           {funcName: "MapLookup", native: divergence_hunt5.MapLookup},
		},
	})
}

func TestDivergenceHunt6(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt6Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"ChannelClose":           {funcName: "ChannelClose", native: divergence_hunt6.ChannelClose},
			"ChannelSelect":          {funcName: "ChannelSelect", native: divergence_hunt6.ChannelSelect},
			"ChannelNilBlock":        {funcName: "ChannelNilBlock", native: divergence_hunt6.ChannelNilBlock},
			"FuncAsValue":            {funcName: "FuncAsValue", native: divergence_hunt6.FuncAsValue},
			"HigherOrderFunc":        {funcName: "HigherOrderFunc", native: divergence_hunt6.HigherOrderFunc},
			"ClosureOverLoop":        {funcName: "ClosureOverLoop", native: divergence_hunt6.ClosureOverLoop},
			"RecursiveFib":           {funcName: "RecursiveFib", native: divergence_hunt6.RecursiveFib},
			"PartialApplication":     {funcName: "PartialApplication", native: divergence_hunt6.PartialApplication},
			"FunctionSlice":          {funcName: "FunctionSlice", native: divergence_hunt6.FunctionSlice},
			"MapFunc":                {funcName: "MapFunc", native: divergence_hunt6.MapFunc},
			"ChannelBufferLen":       {funcName: "ChannelBufferLen", native: divergence_hunt6.ChannelBufferLen},
			"ChannelBufferCap":       {funcName: "ChannelBufferCap", native: divergence_hunt6.ChannelBufferCap},
			"SelectDefault":          {funcName: "SelectDefault", native: divergence_hunt6.SelectDefault},
			"MultiReturnFunc":        {funcName: "MultiReturnFunc", native: divergence_hunt6.MultiReturnFunc},
			"NestedClosure":          {funcName: "NestedClosure", native: divergence_hunt6.NestedClosure},
			"ClosureReturnFunc":      {funcName: "ClosureReturnFunc", native: divergence_hunt6.ClosureReturnFunc},
			"ChannelReceiveOnClosed": {funcName: "ChannelReceiveOnClosed", native: divergence_hunt6.ChannelReceiveOnClosed},
			"FuncTypeDeclaration":    {funcName: "FuncTypeDeclaration", native: divergence_hunt6.FuncTypeDeclaration},
			"VariadicSpread":         {funcName: "VariadicSpread", native: divergence_hunt6.VariadicSpread},
			"InterfaceMethod":        {funcName: "InterfaceMethod", native: divergence_hunt6.InterfaceMethod},
			"StringConversion":       {funcName: "StringConversion", native: divergence_hunt6.StringConversion},
		},
	})
}

func TestDivergenceHunt7(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt7Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"SortInts":            {funcName: "SortInts", native: divergence_hunt7.SortInts},
			"SortStrings":         {funcName: "SortStrings", native: divergence_hunt7.SortStrings},
			"SortFloat64s":        {funcName: "SortFloat64s", native: divergence_hunt7.SortFloat64s},
			"SliceDelete":         {funcName: "SliceDelete", native: divergence_hunt7.SliceDelete},
			"SliceInsert":         {funcName: "SliceInsert", native: divergence_hunt7.SliceInsert},
			"SliceContains":       {funcName: "SliceContains", native: divergence_hunt7.SliceContains},
			"MapKeys":             {funcName: "MapKeys", native: divergence_hunt7.MapKeys},
			"MapValues":           {funcName: "MapValues", native: divergence_hunt7.MapValues},
			"StructWithMethods":   {funcName: "StructWithMethods", native: divergence_hunt7.StructWithMethods},
			"PointerReceiverMethod": {funcName: "PointerReceiverMethod", native: divergence_hunt7.PointerReceiverMethod},
			"TypeAssertion":       {funcName: "TypeAssertion", native: divergence_hunt7.TypeAssertion},
			"TypeAssertionString": {funcName: "TypeAssertionString", native: divergence_hunt7.TypeAssertionString},
			"TypeAssertionFail":   {funcName: "TypeAssertionFail", native: divergence_hunt7.TypeAssertionFail},
			"InterfaceTypeSwitch": {funcName: "InterfaceTypeSwitch", native: divergence_hunt7.InterfaceTypeSwitch},
			"SliceDedupe":         {funcName: "SliceDedupe", native: divergence_hunt7.SliceDedupe},
			"MapMerge":            {funcName: "MapMerge", native: divergence_hunt7.MapMerge},
			"StructSliceSort":     {funcName: "StructSliceSort", native: divergence_hunt7.StructSliceSort},
			"MapInvert":           {funcName: "MapInvert", native: divergence_hunt7.MapInvert},
			"NestedInterface":     {funcName: "NestedInterface", native: divergence_hunt7.NestedInterface},
			"SliceFlatten":        {funcName: "SliceFlatten", native: divergence_hunt7.SliceFlatten},
			"IntSliceSortCustom":  {funcName: "IntSliceSortCustom", native: divergence_hunt7.IntSliceSortCustom},
			"MapCountValues":      {funcName: "MapCountValues", native: divergence_hunt7.MapCountValues},
		},
	})
}

func TestDivergenceHunt8(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt8Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"MutexBasic":            {funcName: "MutexBasic", native: divergence_hunt8.MutexBasic},
			"OnceBasic":             {funcName: "OnceBasic", native: divergence_hunt8.OnceBasic},
			"SliceOfSlice":          {funcName: "SliceOfSlice", native: divergence_hunt8.SliceOfSlice},
			"MapOfMap":              {funcName: "MapOfMap", native: divergence_hunt8.MapOfMap},
			"StructWithSlice":       {funcName: "StructWithSlice", native: divergence_hunt8.StructWithSlice},
			"StructWithMap":         {funcName: "StructWithMap", native: divergence_hunt8.StructWithMap},
			"NestedSliceAppend":     {funcName: "NestedSliceAppend", native: divergence_hunt8.NestedSliceAppend},
			"DeepStruct":            {funcName: "DeepStruct", native: divergence_hunt8.DeepStruct},
			"SliceOfStructAppend":   {funcName: "SliceOfStructAppend", native: divergence_hunt8.SliceOfStructAppend},
			"MapWithSliceValue":     {funcName: "MapWithSliceValue", native: divergence_hunt8.MapWithSliceValue},
			"MutexInDefer":          {funcName: "MutexInDefer", native: divergence_hunt8.MutexInDefer},
			"RWMutexBasic":          {funcName: "RWMutexBasic", native: divergence_hunt8.RWMutexBasic},
			"StructWithFunc":        {funcName: "StructWithFunc", native: divergence_hunt8.StructWithFunc},
			"StructWithPointer":     {funcName: "StructWithPointer", native: divergence_hunt8.StructWithPointer},
			"SliceGrowPattern":      {funcName: "SliceGrowPattern", native: divergence_hunt8.SliceGrowPattern},
			"MapGrowPattern":        {funcName: "MapGrowPattern", native: divergence_hunt8.MapGrowPattern},
			"CompositeLiteralNested": {funcName: "CompositeLiteralNested", native: divergence_hunt8.CompositeLiteralNested},
		},
	})
}

func TestDivergenceHunt9(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt9Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"JSONMarshal":       {funcName: "JSONMarshal", native: divergence_hunt9.JSONMarshal},
			"JSONUnmarshal":     {funcName: "JSONUnmarshal", native: divergence_hunt9.JSONUnmarshal},
			"JSONMarshalMap":    {funcName: "JSONMarshalMap", native: divergence_hunt9.JSONMarshalMap},
			"RegexMatch":        {funcName: "RegexMatch", native: divergence_hunt9.RegexMatch},
			"RegexFind":         {funcName: "RegexFind", native: divergence_hunt9.RegexFind},
			"RegexFindAll":      {funcName: "RegexFindAll", native: divergence_hunt9.RegexFindAll},
			"RegexReplace":      {funcName: "RegexReplace", native: divergence_hunt9.RegexReplace},
			"MathMod":           {funcName: "MathMod", native: divergence_hunt9.MathMod},
			"MathLog":           {funcName: "MathLog", native: divergence_hunt9.MathLog},
			"MathExp":           {funcName: "MathExp", native: divergence_hunt9.MathExp},
			"MathRound":         {funcName: "MathRound", native: divergence_hunt9.MathRound},
			"MathTrunc":         {funcName: "MathTrunc", native: divergence_hunt9.MathTrunc},
			"MathRemainder":     {funcName: "MathRemainder", native: divergence_hunt9.MathRemainder},
			"MathCopysign":      {funcName: "MathCopysign", native: divergence_hunt9.MathCopysign},
			"JSONMarshalSlice":  {funcName: "JSONMarshalSlice", native: divergence_hunt9.JSONMarshalSlice},
			"JSONUnmarshalSlice": {funcName: "JSONUnmarshalSlice", native: divergence_hunt9.JSONUnmarshalSlice},
			"RegexSplit":        {funcName: "RegexSplit", native: divergence_hunt9.RegexSplit},
			"RegexSubmatch":     {funcName: "RegexSubmatch", native: divergence_hunt9.RegexSubmatch},
			"MathHypot":         {funcName: "MathHypot", native: divergence_hunt9.MathHypot},
			"MathPow10":         {funcName: "MathPow10", native: divergence_hunt9.MathPow10},
			"MathSignbit":       {funcName: "MathSignbit", native: divergence_hunt9.MathSignbit},
		},
	})
}

func TestDivergenceHunt10(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt10Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"BinarySearch":       {funcName: "BinarySearch", native: divergence_hunt10.BinarySearch},
			"StackPattern":       {funcName: "StackPattern", native: divergence_hunt10.StackPattern},
			"QueuePattern":       {funcName: "QueuePattern", native: divergence_hunt10.QueuePattern},
			"TwoSum":             {funcName: "TwoSum", native: divergence_hunt10.TwoSum},
			"IsPalindrome":       {funcName: "IsPalindrome", native: divergence_hunt10.IsPalindrome},
			"FizzBuzz":           {funcName: "FizzBuzz", native: divergence_hunt10.FizzBuzz},
			"FmtVerb":            {funcName: "FmtVerb", native: divergence_hunt10.FmtVerb},
			"FmtWidthPrecision":  {funcName: "FmtWidthPrecision", native: divergence_hunt10.FmtWidthPrecision},
			"NestedMapLookup":    {funcName: "NestedMapLookup", native: divergence_hunt10.NestedMapLookup},
			"StructSliceFilter":  {funcName: "StructSliceFilter", native: divergence_hunt10.StructSliceFilter},
			"GCD":                {funcName: "GCD", native: divergence_hunt10.GCD},
			"LCM":                {funcName: "LCM", native: divergence_hunt10.LCM},
			"Power":              {funcName: "Power", native: divergence_hunt10.Power},
			"CountDigits":        {funcName: "CountDigits", native: divergence_hunt10.CountDigits},
			"ReverseInt":         {funcName: "ReverseInt", native: divergence_hunt10.ReverseInt},
			"FibIterative":       {funcName: "FibIterative", native: divergence_hunt10.FibIterative},
			"PrimeCheck":         {funcName: "PrimeCheck", native: divergence_hunt10.PrimeCheck},
			"FactorialIterative": {funcName: "FactorialIterative", native: divergence_hunt10.FactorialIterative},
			"CountingSort":       {funcName: "CountingSort", native: divergence_hunt10.CountingSort},
			"PrefixSum":          {funcName: "PrefixSum", native: divergence_hunt10.PrefixSum},
			"StringAnagram":      {funcName: "StringAnagram", native: divergence_hunt10.StringAnagram},
		},
	})
}

func TestDivergenceHunt11(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt11Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"DeferInLoop":          {funcName: "DeferInLoop", native: divergence_hunt11.DeferInLoop},
			"DeferAndPanicOrder":   {funcName: "DeferAndPanicOrder", native: divergence_hunt11.DeferAndPanicOrder},
			"RecoverInFunction":    {funcName: "RecoverInFunction", native: divergence_hunt11.RecoverInFunction},
			"PanicWithStruct":      {funcName: "PanicWithStruct", native: divergence_hunt11.PanicWithStruct},
			"NamedReturnWithDefer": {funcName: "NamedReturnWithDefer", native: divergence_hunt11.NamedReturnWithDefer},
			"MultipleDeferModify":  {funcName: "MultipleDeferModify", native: divergence_hunt11.MultipleDeferModify},
			"DeferWithArgument":    {funcName: "DeferWithArgument", native: divergence_hunt11.DeferWithArgument},
			"PanicNilValue":        {funcName: "PanicNilValue", native: divergence_hunt11.PanicNilValue},
			"ClosureReturnFunc":    {funcName: "ClosureReturnFunc", native: divergence_hunt11.ClosureReturnFunc},
			"FmtSprintfMulti":      {funcName: "FmtSprintfMulti", native: divergence_hunt11.FmtSprintfMulti},
			"FmtErrorf":            {funcName: "FmtErrorf", native: divergence_hunt11.FmtErrorf},
			"NestedDeferRecover":   {funcName: "NestedDeferRecover", native: divergence_hunt11.NestedDeferRecover},
			"DeferWithMethod":      {funcName: "DeferWithMethod", native: divergence_hunt11.DeferWithMethod},
			"ClosureCaptureSlice":  {funcName: "ClosureCaptureSlice", native: divergence_hunt11.ClosureCaptureSlice},
			"ClosureCaptureMap":    {funcName: "ClosureCaptureMap", native: divergence_hunt11.ClosureCaptureMap},
			"MultiplePanicRecover": {funcName: "MultiplePanicRecover", native: divergence_hunt11.MultiplePanicRecover},
			"DeferRecoverReturnsValue": {funcName: "DeferRecoverReturnsValue", native: divergence_hunt11.DeferRecoverReturnsValue},
			"SliceAppendInClosure": {funcName: "SliceAppendInClosure", native: divergence_hunt11.SliceAppendInClosure},
			"MapWriteInClosure":    {funcName: "MapWriteInClosure", native: divergence_hunt11.MapWriteInClosure},
			"DeferChain":           {funcName: "DeferChain", native: divergence_hunt11.DeferChain},
			"RecoverReturnsNilAfter": {funcName: "RecoverReturnsNilAfter", native: divergence_hunt11.RecoverReturnsNilAfter},
		},
	})
}

func TestDivergenceHunt12(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt12Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"JSONNestedStruct":     {funcName: "JSONNestedStruct", native: divergence_hunt12.JSONNestedStruct},
			"JSONSliceOfStructs":   {funcName: "JSONSliceOfStructs", native: divergence_hunt12.JSONSliceOfStructs},
			"JSONUnmarshalIntoMap": {funcName: "JSONUnmarshalIntoMap", native: divergence_hunt12.JSONUnmarshalIntoMap},
			"StringTitle":          {funcName: "StringTitle", native: divergence_hunt12.StringTitle},
			"StringEqualFold":      {funcName: "StringEqualFold", native: divergence_hunt12.StringEqualFold},
			"StringIndex":          {funcName: "StringIndex", native: divergence_hunt12.StringIndex},
			"StringLastIndex":      {funcName: "StringLastIndex", native: divergence_hunt12.StringLastIndex},
			"StringIndexAny":       {funcName: "StringIndexAny", native: divergence_hunt12.StringIndexAny},
			"StringNewReplacer":    {funcName: "StringNewReplacer", native: divergence_hunt12.StringNewReplacer},
			"StringBuilderGrow":    {funcName: "StringBuilderGrow", native: divergence_hunt12.StringBuilderGrow},
			"SortSliceStable":      {funcName: "SortSliceStable", native: divergence_hunt12.SortSliceStable},
			"SortSearch":           {funcName: "SortSearch", native: divergence_hunt12.SortSearch},
			"FmtSprintfBoolean":    {funcName: "FmtSprintfBoolean", native: divergence_hunt12.FmtSprintfBoolean},
			"FmtSprintfFloat":      {funcName: "FmtSprintfFloat", native: divergence_hunt12.FmtSprintfFloat},
			"FmtSprintfInt":        {funcName: "FmtSprintfInt", native: divergence_hunt12.FmtSprintfInt},
			"FmtSprintfString":     {funcName: "FmtSprintfString", native: divergence_hunt12.FmtSprintfString},
			"JSONMarshalBool":      {funcName: "JSONMarshalBool", native: divergence_hunt12.JSONMarshalBool},
			"JSONUnmarshalBool":    {funcName: "JSONUnmarshalBool", native: divergence_hunt12.JSONUnmarshalBool},
			"JSONMarshalNil":       {funcName: "JSONMarshalNil", native: divergence_hunt12.JSONMarshalNil},
			"SliceMinMaxInt":       {funcName: "SliceMinMaxInt", native: divergence_hunt12.SliceMinMaxInt},
			"StringCountSubstring": {funcName: "StringCountSubstring", native: divergence_hunt12.StringCountSubstring},
			"MapHasKey":            {funcName: "MapHasKey", native: divergence_hunt12.MapHasKey},
		},
	})
}

func TestDivergenceHunt13(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt13Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"StructZeroValue":              {funcName: "StructZeroValue", native: divergence_hunt13.StructZeroValue},
			"StructPointerNil":             {funcName: "StructPointerNil", native: divergence_hunt13.StructPointerNil},
			"StructCopyOnAssign":           {funcName: "StructCopyOnAssign", native: divergence_hunt13.StructCopyOnAssign},
			"StructFieldAccess":            {funcName: "StructFieldAccess", native: divergence_hunt13.StructFieldAccess},
			"InterfaceNilComparison":       {funcName: "InterfaceNilComparison", native: divergence_hunt13.InterfaceNilComparison},
			"InterfaceTypedNil":            {funcName: "InterfaceTypedNil", native: divergence_hunt13.InterfaceTypedNil},
			"TypeAssertionWithBool":        {funcName: "TypeAssertionWithBool", native: divergence_hunt13.TypeAssertionWithBool},
			"MultipleTypeAssertions":       {funcName: "MultipleTypeAssertions", native: divergence_hunt13.MultipleTypeAssertions},
			"PointerToStruct":              {funcName: "PointerToStruct", native: divergence_hunt13.PointerToStruct},
			"PointerToStructModify":        {funcName: "PointerToStructModify", native: divergence_hunt13.PointerToStructModify},
			"StructAsMapValue":             {funcName: "StructAsMapValue", native: divergence_hunt13.StructAsMapValue},
			"StructInSlice":                {funcName: "StructInSlice", native: divergence_hunt13.StructInSlice},
			"IntTypeAlias":                 {funcName: "IntTypeAlias", native: divergence_hunt13.IntTypeAlias},
			"StringTypeAlias":              {funcName: "StringTypeAlias", native: divergence_hunt13.StringTypeAlias},
			"SliceOfAlias":                 {funcName: "SliceOfAlias", native: divergence_hunt13.SliceOfAlias},
			"NestedTypeDefinitions":        {funcName: "NestedTypeDefinitions", native: divergence_hunt13.NestedTypeDefinitions},
			"FmtStruct":                    {funcName: "FmtStruct", native: divergence_hunt13.FmtStruct},
			"FmtPointer":                   {funcName: "FmtPointer", native: divergence_hunt13.FmtPointer},
			"ConversionBetweenNumericTypes": {funcName: "ConversionBetweenNumericTypes", native: divergence_hunt13.ConversionBetweenNumericTypes},
			"UnsignedToSigned":             {funcName: "UnsignedToSigned", native: divergence_hunt13.UnsignedToSigned},
		},
	})
}

func TestDivergenceHunt14(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt14Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"FloatAddPrecision":       {funcName: "FloatAddPrecision", native: divergence_hunt14.FloatAddPrecision},
			"FloatMultiplyPrecision":  {funcName: "FloatMultiplyPrecision", native: divergence_hunt14.FloatMultiplyPrecision},
			"FloatDivPrecision":       {funcName: "FloatDivPrecision", native: divergence_hunt14.FloatDivPrecision},
			"FloatNegative":           {funcName: "FloatNegative", native: divergence_hunt14.FloatNegative},
			"FloatZeroDivision":       {funcName: "FloatZeroDivision", native: divergence_hunt14.FloatZeroDivision},
			"FloatNaNArithmetic":      {funcName: "FloatNaNArithmetic", native: divergence_hunt14.FloatNaNArithmetic},
			"FloatInfArithmetic":      {funcName: "FloatInfArithmetic", native: divergence_hunt14.FloatInfArithmetic},
			"FloatComparisonPrecision": {funcName: "FloatComparisonPrecision", native: divergence_hunt14.FloatComparisonPrecision},
			"IntDivisionTruncation":   {funcName: "IntDivisionTruncation", native: divergence_hunt14.IntDivisionTruncation},
			"IntModulo":               {funcName: "IntModulo", native: divergence_hunt14.IntModulo},
			"NegativeModulo":          {funcName: "NegativeModulo", native: divergence_hunt14.NegativeModulo},
			"Float32NaN":              {funcName: "Float32NaN", native: divergence_hunt14.Float32NaN},
			"Float32Inf":              {funcName: "Float32Inf", native: divergence_hunt14.Float32Inf},
			"MathSin":                 {funcName: "MathSin", native: divergence_hunt14.MathSin},
			"MathCos":                 {funcName: "MathCos", native: divergence_hunt14.MathCos},
			"MathTan":                 {funcName: "MathTan", native: divergence_hunt14.MathTan},
			"MathAtan2":               {funcName: "MathAtan2", native: divergence_hunt14.MathAtan2},
			"MathLog2":                {funcName: "MathLog2", native: divergence_hunt14.MathLog2},
			"MathLog10":               {funcName: "MathLog10", native: divergence_hunt14.MathLog10},
			"FmtFloatFormat":          {funcName: "FmtFloatFormat", native: divergence_hunt14.FmtFloatFormat},
			"FmtIntFormat":            {funcName: "FmtIntFormat", native: divergence_hunt14.FmtIntFormat},
			"FloatMaxMin":             {funcName: "FloatMaxMin", native: divergence_hunt14.FloatMaxMin},
			"Float32Limits":           {funcName: "Float32Limits", native: divergence_hunt14.Float32Limits},
			"ComplexMagnitude":        {funcName: "ComplexMagnitude", native: divergence_hunt14.ComplexMagnitude},
		},
	})
}

func TestDivergenceHunt15(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt15Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"WordCount":              {funcName: "WordCount", native: divergence_hunt15.WordCount},
			"TopKElements":           {funcName: "TopKElements", native: divergence_hunt15.TopKElements},
			"FlattenAndSum":          {funcName: "FlattenAndSum", native: divergence_hunt15.FlattenAndSum},
			"FrequencyCount":         {funcName: "FrequencyCount", native: divergence_hunt15.FrequencyCount},
			"ReverseString":          {funcName: "ReverseString", native: divergence_hunt15.ReverseString},
			"StringPermutationCheck": {funcName: "StringPermutationCheck", native: divergence_hunt15.StringPermutationCheck},
			"MatrixSum":              {funcName: "MatrixSum", native: divergence_hunt15.MatrixSum},
			"MatrixTranspose":        {funcName: "MatrixTranspose", native: divergence_hunt15.MatrixTranspose},
			"JSONEncodeDecode":       {funcName: "JSONEncodeDecode", native: divergence_hunt15.JSONEncodeDecode},
			"StringCompression":      {funcName: "StringCompression", native: divergence_hunt15.StringCompression},
			"UniqueElements":         {funcName: "UniqueElements", native: divergence_hunt15.UniqueElements},
			"IntersectSlices":        {funcName: "IntersectSlices", native: divergence_hunt15.IntersectSlices},
			"MergeSortedSlices":      {funcName: "MergeSortedSlices", native: divergence_hunt15.MergeSortedSlices},
			"MovingAverage":          {funcName: "MovingAverage", native: divergence_hunt15.MovingAverage},
			"SpiralMatrix":           {funcName: "SpiralMatrix", native: divergence_hunt15.SpiralMatrix},
			"FmtStructFormatting":    {funcName: "FmtStructFormatting", native: divergence_hunt15.FmtStructFormatting},
		},
	})
}

func TestDivergenceHunt16(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt16Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"SwitchNoCase":             {funcName: "SwitchNoCase", native: divergence_hunt16.SwitchNoCase},
			"SwitchMultipleCases":      {funcName: "SwitchMultipleCases", native: divergence_hunt16.SwitchMultipleCases},
			"SwitchWithInit":           {funcName: "SwitchWithInit", native: divergence_hunt16.SwitchWithInit},
			"NestedSwitch":             {funcName: "NestedSwitch", native: divergence_hunt16.NestedSwitch},
			"ForRangeWithIndex":        {funcName: "ForRangeWithIndex", native: divergence_hunt16.ForRangeWithIndex},
			"ForRangeWithValue":        {funcName: "ForRangeWithValue", native: divergence_hunt16.ForRangeWithValue},
			"ForRangeMap":              {funcName: "ForRangeMap", native: divergence_hunt16.ForRangeMap},
			"IfElseChain":              {funcName: "IfElseChain", native: divergence_hunt16.IfElseChain},
			"NestedIfElse":             {funcName: "NestedIfElse", native: divergence_hunt16.NestedIfElse},
			"InfiniteLoopBreak":        {funcName: "InfiniteLoopBreak", native: divergence_hunt16.InfiniteLoopBreak},
			"ForLoopContinue":          {funcName: "ForLoopContinue", native: divergence_hunt16.ForLoopContinue},
			"LoopWithMultipleBreaks":   {funcName: "LoopWithMultipleBreaks", native: divergence_hunt16.LoopWithMultipleBreaks},
			"SwitchExpression":         {funcName: "SwitchExpression", native: divergence_hunt16.SwitchExpression},
			"ForRangeString":           {funcName: "ForRangeString", native: divergence_hunt16.ForRangeString},
			"ForRangeEmptySlice":       {funcName: "ForRangeEmptySlice", native: divergence_hunt16.ForRangeEmptySlice},
			"DoubleLoop":               {funcName: "DoubleLoop", native: divergence_hunt16.DoubleLoop},
			"LoopAccumulator":          {funcName: "LoopAccumulator", native: divergence_hunt16.LoopAccumulator},
			"SwitchFallthroughSimulated": {funcName: "SwitchFallthroughSimulated", native: divergence_hunt16.SwitchFallthroughSimulated},
			"EarlyReturn":              {funcName: "EarlyReturn", native: divergence_hunt16.EarlyReturn},
			"LoopWithEarlyReturn":       {funcName: "LoopWithEarlyReturn", native: divergence_hunt16.LoopWithEarlyReturn},
		},
	})
}

func TestDivergenceHunt17(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt17Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"InterfaceComposition":    {funcName: "InterfaceComposition", native: divergence_hunt17.InterfaceComposition},
			"InterfaceEmpty":          {funcName: "InterfaceEmpty", native: divergence_hunt17.InterfaceEmpty},
			"InterfaceSlice":          {funcName: "InterfaceSlice", native: divergence_hunt17.InterfaceSlice},
			"InterfaceMap":            {funcName: "InterfaceMap", native: divergence_hunt17.InterfaceMap},
			"StructMethodOnPointer":   {funcName: "StructMethodOnPointer", native: divergence_hunt17.StructMethodOnPointer},
			"StructMethodOnValue":     {funcName: "StructMethodOnValue", native: divergence_hunt17.StructMethodOnValue},
			"MethodChain":             {funcName: "MethodChain", native: divergence_hunt17.MethodChain},
			"PolymorphismPattern":     {funcName: "PolymorphismPattern", native: divergence_hunt17.PolymorphismPattern},
			"NullableInterface":       {funcName: "NullableInterface", native: divergence_hunt17.NullableInterface},
			"InterfaceTypeAssertion":  {funcName: "InterfaceTypeAssertion", native: divergence_hunt17.InterfaceTypeAssertion},
			"EmbeddedStructAccess":    {funcName: "EmbeddedStructAccess", native: divergence_hunt17.EmbeddedStructAccess},
			"NestedStructAccess":      {funcName: "NestedStructAccess", native: divergence_hunt17.NestedStructAccess},
			"StructSliceMethods":      {funcName: "StructSliceMethods", native: divergence_hunt17.StructSliceMethods},
			"FmtInterface":            {funcName: "FmtInterface", native: divergence_hunt17.FmtInterface},
			"FmtNilInterface":         {funcName: "FmtNilInterface", native: divergence_hunt17.FmtNilInterface},
			"StructComparison":        {funcName: "StructComparison", native: divergence_hunt17.StructComparison},
			"InterfaceEquality":       {funcName: "InterfaceEquality", native: divergence_hunt17.InterfaceEquality},
			"InterfaceInequality":     {funcName: "InterfaceInequality", native: divergence_hunt17.InterfaceInequality},
		},
	})
}

func TestDivergenceHunt18(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt18Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"StringToIntConversion":  {funcName: "StringToIntConversion", native: divergence_hunt18.StringToIntConversion},
			"IntToStringConversion":  {funcName: "IntToStringConversion", native: divergence_hunt18.IntToStringConversion},
			"FloatToStringConversion": {funcName: "FloatToStringConversion", native: divergence_hunt18.FloatToStringConversion},
			"StringToFloatConversion": {funcName: "StringToFloatConversion", native: divergence_hunt18.StringToFloatConversion},
			"BoolToStringConversion": {funcName: "BoolToStringConversion", native: divergence_hunt18.BoolToStringConversion},
			"StringToBoolConversion": {funcName: "StringToBoolConversion", native: divergence_hunt18.StringToBoolConversion},
			"StringSplitJoin":        {funcName: "StringSplitJoin", native: divergence_hunt18.StringSplitJoin},
			"StringTrimSpace":        {funcName: "StringTrimSpace", native: divergence_hunt18.StringTrimSpace},
			"StringTrimPrefix":       {funcName: "StringTrimPrefix", native: divergence_hunt18.StringTrimPrefix},
			"StringTrimSuffix":       {funcName: "StringTrimSuffix", native: divergence_hunt18.StringTrimSuffix},
			"StringReplaceAll":       {funcName: "StringReplaceAll", native: divergence_hunt18.StringReplaceAll},
			"StringBuilderPattern":   {funcName: "StringBuilderPattern", native: divergence_hunt18.StringBuilderPattern},
			"StringRuneConversion":   {funcName: "StringRuneConversion", native: divergence_hunt18.StringRuneConversion},
			"RuneToStringConversion": {funcName: "RuneToStringConversion", native: divergence_hunt18.RuneToStringConversion},
			"StringByteConversion":   {funcName: "StringByteConversion", native: divergence_hunt18.StringByteConversion},
			"ByteToStringConversion": {funcName: "ByteToStringConversion", native: divergence_hunt18.ByteToStringConversion},
			"FmtSprintfComplex":      {funcName: "FmtSprintfComplex", native: divergence_hunt18.FmtSprintfComplex},
			"FmtSprintfPadding":      {funcName: "FmtSprintfPadding", native: divergence_hunt18.FmtSprintfPadding},
			"StringPadLeft":          {funcName: "StringPadLeft", native: divergence_hunt18.StringPadLeft},
			"StringPadRight":         {funcName: "StringPadRight", native: divergence_hunt18.StringPadRight},
			"CamelCaseSplit":         {funcName: "CamelCaseSplit", native: divergence_hunt18.CamelCaseSplit},
			"StringReverseWords":     {funcName: "StringReverseWords", native: divergence_hunt18.StringReverseWords},
		},
	})
}

func TestDivergenceHunt19(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt19Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"EmptySliceOperations":   {funcName: "EmptySliceOperations", native: divergence_hunt19.EmptySliceOperations},
			"EmptyMapOperations":     {funcName: "EmptyMapOperations", native: divergence_hunt19.EmptyMapOperations},
			"EmptyStringOperations":  {funcName: "EmptyStringOperations", native: divergence_hunt19.EmptyStringOperations},
			"ZeroValueInt":           {funcName: "ZeroValueInt", native: divergence_hunt19.ZeroValueInt},
			"ZeroValueFloat":         {funcName: "ZeroValueFloat", native: divergence_hunt19.ZeroValueFloat},
			"ZeroValueBool":          {funcName: "ZeroValueBool", native: divergence_hunt19.ZeroValueBool},
			"ZeroValueString":        {funcName: "ZeroValueString", native: divergence_hunt19.ZeroValueString},
			"ZeroValueSlice":         {funcName: "ZeroValueSlice", native: divergence_hunt19.ZeroValueSlice},
			"ZeroValueMap":           {funcName: "ZeroValueMap", native: divergence_hunt19.ZeroValueMap},
			"ZeroValuePointer":       {funcName: "ZeroValuePointer", native: divergence_hunt19.ZeroValuePointer},
			"NilSliceAppend":         {funcName: "NilSliceAppend", native: divergence_hunt19.NilSliceAppend},
			"NilMapRead":             {funcName: "NilMapRead", native: divergence_hunt19.NilMapRead},
			"NilSliceRange":          {funcName: "NilSliceRange", native: divergence_hunt19.NilSliceRange},
			"NilMapRange":            {funcName: "NilMapRange", native: divergence_hunt19.NilMapRange},
			"NilChannelRead":         {funcName: "NilChannelRead", native: divergence_hunt19.NilChannelRead},
			"SliceBoundary":          {funcName: "SliceBoundary", native: divergence_hunt19.SliceBoundary},
			"MapBoundary":            {funcName: "MapBoundary", native: divergence_hunt19.MapBoundary},
			"ErrorHandlingPattern":   {funcName: "ErrorHandlingPattern", native: divergence_hunt19.ErrorHandlingPattern},
			"MultipleErrorCheck":     {funcName: "MultipleErrorCheck", native: divergence_hunt19.MultipleErrorCheck},
			"NilFuncVariable":        {funcName: "NilFuncVariable", native: divergence_hunt19.NilFuncVariable},
			"EmptyInterfaceContains": {funcName: "EmptyInterfaceContains", native: divergence_hunt19.EmptyInterfaceContains},
			"StructZeroValueFields":  {funcName: "StructZeroValueFields", native: divergence_hunt19.StructZeroValueFields},
		},
	})
}

func TestDivergenceHunt20(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt20Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"StudentGradeSystem": {funcName: "StudentGradeSystem", native: divergence_hunt20.StudentGradeSystem},
			"TextProcessing":     {funcName: "TextProcessing", native: divergence_hunt20.TextProcessing},
			"DataTransform":      {funcName: "DataTransform", native: divergence_hunt20.DataTransform},
			"InventorySystem":    {funcName: "InventorySystem", native: divergence_hunt20.InventorySystem},
			"JSONProcessing":     {funcName: "JSONProcessing", native: divergence_hunt20.JSONProcessing},
			"StringProcessing":   {funcName: "StringProcessing", native: divergence_hunt20.StringProcessing},
			"SortAndSearch":      {funcName: "SortAndSearch", native: divergence_hunt20.SortAndSearch},
			"MatrixOperations":   {funcName: "MatrixOperations", native: divergence_hunt20.MatrixOperations},
			"FmtTable":           {funcName: "FmtTable", native: divergence_hunt20.FmtTable},
			"Histogram":          {funcName: "Histogram", native: divergence_hunt20.Histogram},
			"ParseAndCompute":    {funcName: "ParseAndCompute", native: divergence_hunt20.ParseAndCompute},
			"SetOperations":      {funcName: "SetOperations", native: divergence_hunt20.SetOperations},
			"GroupBy":            {funcName: "GroupBy", native: divergence_hunt20.GroupBy},
			"RunningSum":         {funcName: "RunningSum", native: divergence_hunt20.RunningSum},
			"SlidingWindow":      {funcName: "SlidingWindow", native: divergence_hunt20.SlidingWindow},
		},
	})
}

func TestDivergenceHunt21(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt21Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"MapIterateSum":     {funcName: "MapIterateSum", native: divergence_hunt21.MapIterateSum},
			"SliceRotateLeft":   {funcName: "SliceRotateLeft", native: divergence_hunt21.SliceRotateLeft},
			"SliceRotateRight":  {funcName: "SliceRotateRight", native: divergence_hunt21.SliceRotateRight},
			"SliceChunk":        {funcName: "SliceChunk", native: divergence_hunt21.SliceChunk},
			"MapFilterSlice":    {funcName: "MapFilterSlice", native: divergence_hunt21.MapFilterSlice},
			"ReducePattern":     {funcName: "ReducePattern", native: divergence_hunt21.ReducePattern},
			"ZipSlices":         {funcName: "ZipSlices", native: divergence_hunt21.ZipSlices},
			"SliceCompact":      {funcName: "SliceCompact", native: divergence_hunt21.SliceCompact},
			"MapMergeOverwrite": {funcName: "MapMergeOverwrite", native: divergence_hunt21.MapMergeOverwrite},
			"SlicePartition":    {funcName: "SlicePartition", native: divergence_hunt21.SlicePartition},
			"NestedMapAccess":   {funcName: "NestedMapAccess", native: divergence_hunt21.NestedMapAccess},
			"FlattenMap":        {funcName: "FlattenMap", native: divergence_hunt21.FlattenMap},
			"MapKeySlice":       {funcName: "MapKeySlice", native: divergence_hunt21.MapKeySlice},
			"SliceSlidingWindow": {funcName: "SliceSlidingWindow", native: divergence_hunt21.SliceSlidingWindow},
			"MultiLevelSlice":   {funcName: "MultiLevelSlice", native: divergence_hunt21.MultiLevelSlice},
		},
	})
}

func TestDivergenceHunt22(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt22Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"JSONMarshalInt":     {funcName: "JSONMarshalInt", native: divergence_hunt22.JSONMarshalInt},
			"JSONMarshalString":  {funcName: "JSONMarshalString", native: divergence_hunt22.JSONMarshalString},
			"JSONMarshalFloat":   {funcName: "JSONMarshalFloat", native: divergence_hunt22.JSONMarshalFloat},
			"JSONUnmarshalInt":   {funcName: "JSONUnmarshalInt", native: divergence_hunt22.JSONUnmarshalInt},
			"JSONUnmarshalString": {funcName: "JSONUnmarshalString", native: divergence_hunt22.JSONUnmarshalString},
			"JSONUnmarshalFloat": {funcName: "JSONUnmarshalFloat", native: divergence_hunt22.JSONUnmarshalFloat},
			"JSONUnmarshalArray": {funcName: "JSONUnmarshalArray", native: divergence_hunt22.JSONUnmarshalArray},
			"FmtVerbP":           {funcName: "FmtVerbP", native: divergence_hunt22.FmtVerbP},
			"FmtVerbT":           {funcName: "FmtVerbT", native: divergence_hunt22.FmtVerbT},
			"FmtVerbV":           {funcName: "FmtVerbV", native: divergence_hunt22.FmtVerbV},
			"FmtVerbPlusV":       {funcName: "FmtVerbPlusV", native: divergence_hunt22.FmtVerbPlusV},
			"FmtVerbHashV":       {funcName: "FmtVerbHashV", native: divergence_hunt22.FmtVerbHashV},
			"FmtSprintfPointer":  {funcName: "FmtSprintfPointer", native: divergence_hunt22.FmtSprintfPointer},
			"ErrorWrap":          {funcName: "ErrorWrap", native: divergence_hunt22.ErrorWrap},
			"ErrorIs":            {funcName: "ErrorIs", native: divergence_hunt22.ErrorIs},
			"JSONNestedMap":      {funcName: "JSONNestedMap", native: divergence_hunt22.JSONNestedMap},
			"JSONStructTag":      {funcName: "JSONStructTag", native: divergence_hunt22.JSONStructTag},
			"JSONOmitEmpty":      {funcName: "JSONOmitEmpty", native: divergence_hunt22.JSONOmitEmpty},
			"FmtWidthInt":        {funcName: "FmtWidthInt", native: divergence_hunt22.FmtWidthInt},
			"FmtFloatScientific": {funcName: "FmtFloatScientific", native: divergence_hunt22.FmtFloatScientific},
		},
	})
}

func TestDivergenceHunt23(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt23Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NewInt": {funcName: "NewInt", native: divergence_hunt23.NewInt}, "NewStruct": {funcName: "NewStruct", native: divergence_hunt23.NewStruct}, "MakeSliceLen": {funcName: "MakeSliceLen", native: divergence_hunt23.MakeSliceLen}, "MakeSliceLenCap": {funcName: "MakeSliceLenCap", native: divergence_hunt23.MakeSliceLenCap}, "MakeMapSize": {funcName: "MakeMapSize", native: divergence_hunt23.MakeMapSize}, "PointerSwap": {funcName: "PointerSwap", native: divergence_hunt23.PointerSwap}, "StructPointerNew": {funcName: "StructPointerNew", native: divergence_hunt23.StructPointerNew}, "SliceOfNew": {funcName: "SliceOfNew", native: divergence_hunt23.SliceOfNew}, "PointerToSlice": {funcName: "PointerToSlice", native: divergence_hunt23.PointerToSlice}, "PointerToMap": {funcName: "PointerToMap", native: divergence_hunt23.PointerToMap}, "DoublePointer": {funcName: "DoublePointer", native: divergence_hunt23.DoublePointer}, "PointerArithmeticSim": {funcName: "PointerArithmeticSim", native: divergence_hunt23.PointerArithmeticSim}, "NewArray": {funcName: "NewArray", native: divergence_hunt23.NewArray}, "SliceFromArray": {funcName: "SliceFromArray", native: divergence_hunt23.SliceFromArray}, "SliceFromArrayPointer": {funcName: "SliceFromArrayPointer", native: divergence_hunt23.SliceFromArrayPointer}, "MapPointer": {funcName: "MapPointer", native: divergence_hunt23.MapPointer}, "StructPointerMethod": {funcName: "StructPointerMethod", native: divergence_hunt23.StructPointerMethod}, "PointerComparison": {funcName: "PointerComparison", native: divergence_hunt23.PointerComparison}, "NilPointerComparison": {funcName: "NilPointerComparison", native: divergence_hunt23.NilPointerComparison},
	}})
}
func TestDivergenceHunt24(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt24Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SortAndDedupe": {funcName: "SortAndDedupe", native: divergence_hunt24.SortAndDedupe}, "WordFrequency": {funcName: "WordFrequency", native: divergence_hunt24.WordFrequency}, "CSVLikeParsing": {funcName: "CSVLikeParsing", native: divergence_hunt24.CSVLikeParsing}, "HistogramFromData": {funcName: "HistogramFromData", native: divergence_hunt24.HistogramFromData}, "FlattenJSON": {funcName: "FlattenJSON", native: divergence_hunt24.FlattenJSON}, "StringTokenize": {funcName: "StringTokenize", native: divergence_hunt24.StringTokenize}, "MatrixRowColSum": {funcName: "MatrixRowColSum", native: divergence_hunt24.MatrixRowColSum}, "StringTemplate": {funcName: "StringTemplate", native: divergence_hunt24.StringTemplate}, "MapTransformKeys": {funcName: "MapTransformKeys", native: divergence_hunt24.MapTransformKeys}, "SlicePartitionPoint": {funcName: "SlicePartitionPoint", native: divergence_hunt24.SlicePartitionPoint}, "NestedLoopBreak": {funcName: "NestedLoopBreak", native: divergence_hunt24.NestedLoopBreak}, "RecursiveSum": {funcName: "RecursiveSum", native: divergence_hunt24.RecursiveSum}, "ReverseSliceInPlace": {funcName: "ReverseSliceInPlace", native: divergence_hunt24.ReverseSliceInPlace}, "MapToSlice": {funcName: "MapToSlice", native: divergence_hunt24.MapToSlice}, "StringDiff": {funcName: "StringDiff", native: divergence_hunt24.StringDiff}, "FmtSlice": {funcName: "FmtSlice", native: divergence_hunt24.FmtSlice},
	}})
}
func TestDivergenceHunt25(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt25Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"DeferStack": {funcName: "DeferStack", native: divergence_hunt25.DeferStack}, "DeferInClosure": {funcName: "DeferInClosure", native: divergence_hunt25.DeferInClosure}, "RecoverInNestedDefer": {funcName: "RecoverInNestedDefer", native: divergence_hunt25.RecoverInNestedDefer}, "MultipleRecover": {funcName: "MultipleRecover", native: divergence_hunt25.MultipleRecover}, "DeferClosureCapture": {funcName: "DeferClosureCapture", native: divergence_hunt25.DeferClosureCapture}, "DeferClosureCopy": {funcName: "DeferClosureCopy", native: divergence_hunt25.DeferClosureCopy}, "PanicInDeferRecover": {funcName: "PanicInDeferRecover", native: divergence_hunt25.PanicInDeferRecover}, "DeferModifyNamedReturn": {funcName: "DeferModifyNamedReturn", native: divergence_hunt25.DeferModifyNamedReturn}, "NestedPanicRecover": {funcName: "NestedPanicRecover", native: divergence_hunt25.NestedPanicRecover}, "ClosureWithDefer": {funcName: "ClosureWithDefer", native: divergence_hunt25.ClosureWithDefer}, "RecursiveWithDefer": {funcName: "RecursiveWithDefer", native: divergence_hunt25.RecursiveWithDefer}, "PanicRecoverTypeSwitch": {funcName: "PanicRecoverTypeSwitch", native: divergence_hunt25.PanicRecoverTypeSwitch}, "DeferMultipleModifies": {funcName: "DeferMultipleModifies", native: divergence_hunt25.DeferMultipleModifies}, "RecoverReturnsPanicValue": {funcName: "RecoverReturnsPanicValue", native: divergence_hunt25.RecoverReturnsPanicValue}, "DeferInMethod": {funcName: "DeferInMethod", native: divergence_hunt25.DeferInMethod}, "ClosureState": {funcName: "ClosureState", native: divergence_hunt25.ClosureState}, "ClosureSharedState": {funcName: "ClosureSharedState", native: divergence_hunt25.ClosureSharedState}, "FmtDefer": {funcName: "FmtDefer", native: divergence_hunt25.FmtDefer},
	}})
}
func TestDivergenceHunt26(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt26Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Int8Range": {funcName: "Int8Range", native: divergence_hunt26.Int8Range}, "Int8MinRange": {funcName: "Int8MinRange", native: divergence_hunt26.Int8MinRange}, "Uint8Max": {funcName: "Uint8Max", native: divergence_hunt26.Uint8Max}, "Int16Range": {funcName: "Int16Range", native: divergence_hunt26.Int16Range}, "Uint16Max": {funcName: "Uint16Max", native: divergence_hunt26.Uint16Max}, "Float32Smallest": {funcName: "Float32Smallest", native: divergence_hunt26.Float32Smallest}, "Complex64Basic": {funcName: "Complex64Basic", native: divergence_hunt26.Complex64Basic}, "Complex128Basic": {funcName: "Complex128Basic", native: divergence_hunt26.Complex128Basic}, "RuneType": {funcName: "RuneType", native: divergence_hunt26.RuneType}, "ByteType": {funcName: "ByteType", native: divergence_hunt26.ByteType}, "StringType": {funcName: "StringType", native: divergence_hunt26.StringType}, "BoolType": {funcName: "BoolType", native: divergence_hunt26.BoolType}, "IntType": {funcName: "IntType", native: divergence_hunt26.IntType}, "Int64Type": {funcName: "Int64Type", native: divergence_hunt26.Int64Type}, "UintType": {funcName: "UintType", native: divergence_hunt26.UintType}, "Uint64Type": {funcName: "Uint64Type", native: divergence_hunt26.Uint64Type}, "Float64Type": {funcName: "Float64Type", native: divergence_hunt26.Float64Type}, "Float32Type": {funcName: "Float32Type", native: divergence_hunt26.Float32Type}, "TypeConversionChain": {funcName: "TypeConversionChain", native: divergence_hunt26.TypeConversionChain}, "UnsignedConversion": {funcName: "UnsignedConversion", native: divergence_hunt26.UnsignedConversion}, "SignedToUnsigned": {funcName: "SignedToUnsigned", native: divergence_hunt26.SignedToUnsigned}, "FloatToIntTrunc": {funcName: "FloatToIntTrunc", native: divergence_hunt26.FloatToIntTrunc}, "IntToFloatPrecise": {funcName: "IntToFloatPrecise", native: divergence_hunt26.IntToFloatPrecise}, "StringToSlice": {funcName: "StringToSlice", native: divergence_hunt26.StringToSlice}, "SliceToString": {funcName: "SliceToString", native: divergence_hunt26.SliceToString},
	}})
}
func TestDivergenceHunt27(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt27Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringSort": {funcName: "StringSort", native: divergence_hunt27.StringSort}, "StringUnique": {funcName: "StringUnique", native: divergence_hunt27.StringUnique}, "StringIsDigit": {funcName: "StringIsDigit", native: divergence_hunt27.StringIsDigit}, "StringIsAlpha": {funcName: "StringIsAlpha", native: divergence_hunt27.StringIsAlpha}, "StringToUpperLower": {funcName: "StringToUpperLower", native: divergence_hunt27.StringToUpperLower}, "StringCapitalize": {funcName: "StringCapitalize", native: divergence_hunt27.StringCapitalize}, "StringCountWords": {funcName: "StringCountWords", native: divergence_hunt27.StringCountWords}, "StringReverseWords": {funcName: "StringReverseWords", native: divergence_hunt27.StringReverseWords}, "FmtInteger": {funcName: "FmtInteger", native: divergence_hunt27.FmtInteger}, "FmtHexInt": {funcName: "FmtHexInt", native: divergence_hunt27.FmtHexInt}, "FmtOctalInt": {funcName: "FmtOctalInt", native: divergence_hunt27.FmtOctalInt}, "FmtBinaryInt": {funcName: "FmtBinaryInt", native: divergence_hunt27.FmtBinaryInt}, "FmtCharFromInt": {funcName: "FmtCharFromInt", native: divergence_hunt27.FmtCharFromInt}, "FmtUnicode": {funcName: "FmtUnicode", native: divergence_hunt27.FmtUnicode}, "SortIntSliceDesc": {funcName: "SortIntSliceDesc", native: divergence_hunt27.SortIntSliceDesc}, "SortFloatSliceDesc": {funcName: "SortFloatSliceDesc", native: divergence_hunt27.SortFloatSliceDesc}, "StringJoinWithSep": {funcName: "StringJoinWithSep", native: divergence_hunt27.StringJoinWithSep}, "StringSplitN": {funcName: "StringSplitN", native: divergence_hunt27.StringSplitN}, "StringRepeatN": {funcName: "StringRepeatN", native: divergence_hunt27.StringRepeatN}, "StringMapFunc": {funcName: "StringMapFunc", native: divergence_hunt27.StringMapFunc},
	}})
}
func TestDivergenceHunt28(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt28Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ChannelSendRecv": {funcName: "ChannelSendRecv", native: divergence_hunt28.ChannelSendRecv}, "ChannelBuffered": {funcName: "ChannelBuffered", native: divergence_hunt28.ChannelBuffered}, "ChannelCloseRange": {funcName: "ChannelCloseRange", native: divergence_hunt28.ChannelCloseRange}, "ChannelSelectTwo": {funcName: "ChannelSelectTwo", native: divergence_hunt28.ChannelSelectTwo}, "ChannelSelectDefault2": {funcName: "ChannelSelectDefault2", native: divergence_hunt28.ChannelSelectDefault2}, "ChannelNilSelect": {funcName: "ChannelNilSelect", native: divergence_hunt28.ChannelNilSelect}, "ChannelLen": {funcName: "ChannelLen", native: divergence_hunt28.ChannelLen}, "ChannelCap2": {funcName: "ChannelCap2", native: divergence_hunt28.ChannelCap2}, "ChannelRecvAfterClose": {funcName: "ChannelRecvAfterClose", native: divergence_hunt28.ChannelRecvAfterClose}, "ChannelDirection": {funcName: "ChannelDirection", native: divergence_hunt28.ChannelDirection}, "SelectMultipleReady": {funcName: "SelectMultipleReady", native: divergence_hunt28.SelectMultipleReady}, "ChannelAsSignal": {funcName: "ChannelAsSignal", native: divergence_hunt28.ChannelAsSignal},
	}})
}
func TestDivergenceHunt29(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt29Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SimpleError": {funcName: "SimpleError", native: divergence_hunt29.SimpleError}, "ErrorWithFormat": {funcName: "ErrorWithFormat", native: divergence_hunt29.ErrorWithFormat}, "ValidatePositive": {funcName: "ValidatePositive", native: divergence_hunt29.ValidatePositive}, "ValidateRange": {funcName: "ValidateRange", native: divergence_hunt29.ValidateRange}, "ErrorPropagation": {funcName: "ErrorPropagation", native: divergence_hunt29.ErrorPropagation}, "ErrorInDefer": {funcName: "ErrorInDefer", native: divergence_hunt29.ErrorInDefer}, "MultiErrorCollect": {funcName: "MultiErrorCollect", native: divergence_hunt29.MultiErrorCollect}, "ErrorTypeAssertion": {funcName: "ErrorTypeAssertion", native: divergence_hunt29.ErrorTypeAssertion}, "PanicWithFmtError": {funcName: "PanicWithFmtError", native: divergence_hunt29.PanicWithFmtError}, "NilErrorCheck": {funcName: "NilErrorCheck", native: divergence_hunt29.NilErrorCheck}, "ErrorStringMethods": {funcName: "ErrorStringMethods", native: divergence_hunt29.ErrorStringMethods}, "ValidateStruct": {funcName: "ValidateStruct", native: divergence_hunt29.ValidateStruct}, "ErrorInClosure": {funcName: "ErrorInClosure", native: divergence_hunt29.ErrorInClosure}, "FmtErrorfWrap": {funcName: "FmtErrorfWrap", native: divergence_hunt29.FmtErrorfWrap},
	}})
}
func TestDivergenceHunt30(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt30Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Comprehensive1": {funcName: "Comprehensive1", native: divergence_hunt30.Comprehensive1}, "Comprehensive2": {funcName: "Comprehensive2", native: divergence_hunt30.Comprehensive2}, "Comprehensive3": {funcName: "Comprehensive3", native: divergence_hunt30.Comprehensive3}, "Comprehensive4": {funcName: "Comprehensive4", native: divergence_hunt30.Comprehensive4}, "Comprehensive5": {funcName: "Comprehensive5", native: divergence_hunt30.Comprehensive5}, "Comprehensive6": {funcName: "Comprehensive6", native: divergence_hunt30.Comprehensive6}, "Comprehensive7": {funcName: "Comprehensive7", native: divergence_hunt30.Comprehensive7}, "Comprehensive8": {funcName: "Comprehensive8", native: divergence_hunt30.Comprehensive8}, "Comprehensive9": {funcName: "Comprehensive9", native: divergence_hunt30.Comprehensive9}, 		"Comprehensive10": {funcName: "Comprehensive10", native: divergence_hunt30.Comprehensive10},
	}})
}
func TestDivergenceHunt31(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt31Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ValueReceiverNoMutation": {funcName: "ValueReceiverNoMutation", native: divergence_hunt31.ValueReceiverNoMutation}, "PointerReceiverChain": {funcName: "PointerReceiverChain", native: divergence_hunt31.PointerReceiverChain}, "ValueReceiverCopy": {funcName: "ValueReceiverCopy", native: divergence_hunt31.ValueReceiverCopy}, "StructMethodOnLiteral": {funcName: "StructMethodOnLiteral", native: divergence_hunt31.StructMethodOnLiteral}, "NestedMethodCall": {funcName: "NestedMethodCall", native: divergence_hunt31.NestedMethodCall}, "MethodValueVsPointer": {funcName: "MethodValueVsPointer", native: divergence_hunt31.MethodValueVsPointer}, "MethodOnValueStruct": {funcName: "MethodOnValueStruct", native: divergence_hunt31.MethodOnValueStruct}, "InterfaceMethodCall": {funcName: "InterfaceMethodCall", native: divergence_hunt31.InterfaceMethodCall}, "InterfaceMethodOnPointer": {funcName: "InterfaceMethodOnPointer", native: divergence_hunt31.InterfaceMethodOnPointer}, "MethodReturnsMultipleValues": {funcName: "MethodReturnsMultipleValues", native: divergence_hunt31.MethodReturnsMultipleValues}, "StructWithBoolMethod": {funcName: "StructWithBoolMethod", native: divergence_hunt31.StructWithBoolMethod}, "FmtStructWithMethods": {funcName: "FmtStructWithMethods", native: divergence_hunt31.FmtStructWithMethods}, "StructSliceWithMethods": {funcName: "StructSliceWithMethods", native: divergence_hunt31.StructSliceWithMethods}, "EmbedStructMethod": {funcName: "EmbedStructMethod", native: divergence_hunt31.EmbedStructMethod},
	}})
}
func TestDivergenceHunt32(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt32Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ArrayLenCap": {funcName: "ArrayLenCap", native: divergence_hunt32.ArrayLenCap}, "ArrayCopyValue": {funcName: "ArrayCopyValue", native: divergence_hunt32.ArrayCopyValue}, "ArrayPointerModify": {funcName: "ArrayPointerModify", native: divergence_hunt32.ArrayPointerModify}, "ArrayAsArg": {funcName: "ArrayAsArg", native: divergence_hunt32.ArrayAsArg}, "ArrayPointerAsArg": {funcName: "ArrayPointerAsArg", native: divergence_hunt32.ArrayPointerAsArg}, "ArrayIteration": {funcName: "ArrayIteration", native: divergence_hunt32.ArrayIteration}, "ArrayIndexAccess": {funcName: "ArrayIndexAccess", native: divergence_hunt32.ArrayIndexAccess}, "ArrayZeroValue": {funcName: "ArrayZeroValue", native: divergence_hunt32.ArrayZeroValue}, "ArrayOfString": {funcName: "ArrayOfString", native: divergence_hunt32.ArrayOfString}, "ArrayOfStruct": {funcName: "ArrayOfStruct", native: divergence_hunt32.ArrayOfStruct}, "SliceFromArray": {funcName: "SliceFromArray", native: divergence_hunt32.SliceFromArray}, "ArrayComparison": {funcName: "ArrayComparison", native: divergence_hunt32.ArrayComparison}, "MultiDimensionalArray": {funcName: "MultiDimensionalArray", native: divergence_hunt32.MultiDimensionalArray}, "ArrayInStruct": {funcName: "ArrayInStruct", native: divergence_hunt32.ArrayInStruct}, "FmtArray": {funcName: "FmtArray", native: divergence_hunt32.FmtArray}, "ArrayLiteralPartial": {funcName: "ArrayLiteralPartial", native: divergence_hunt32.ArrayLiteralPartial}, "ArrayLiteralIndex": {funcName: "ArrayLiteralIndex", native: divergence_hunt32.ArrayLiteralIndex},
	}})
}
func TestDivergenceHunt33(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt33Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Int8Overflow": {funcName: "Int8Overflow", native: divergence_hunt33.Int8Overflow}, "Int8Underflow": {funcName: "Int8Underflow", native: divergence_hunt33.Int8Underflow}, "Int16Overflow": {funcName: "Int16Overflow", native: divergence_hunt33.Int16Overflow}, "Uint16Overflow": {funcName: "Uint16Overflow", native: divergence_hunt33.Uint16Overflow}, "Uint32Arith": {funcName: "Uint32Arith", native: divergence_hunt33.Uint32Arith}, "Int32Arith": {funcName: "Int32Arith", native: divergence_hunt33.Int32Arith}, "ShiftLeft8": {funcName: "ShiftLeft8", native: divergence_hunt33.ShiftLeft8}, "ShiftRight8": {funcName: "ShiftRight8", native: divergence_hunt33.ShiftRight8}, "NegateInt8": {funcName: "NegateInt8", native: divergence_hunt33.NegateInt8}, "NegateInt16": {funcName: "NegateInt16", native: divergence_hunt33.NegateInt16}, "MixedIntArith": {funcName: "MixedIntArith", native: divergence_hunt33.MixedIntArith}, "IntDivTruncation": {funcName: "IntDivTruncation", native: divergence_hunt33.IntDivTruncation}, "IntModNegative": {funcName: "IntModNegative", native: divergence_hunt33.IntModNegative}, "UintDivTruncation": {funcName: "UintDivTruncation", native: divergence_hunt33.UintDivTruncation}, "UintMod": {funcName: "UintMod", native: divergence_hunt33.UintMod}, "BitwiseAndNot": {funcName: "BitwiseAndNot", native: divergence_hunt33.BitwiseAndNot}, "BitwiseXor": {funcName: "BitwiseXor", native: divergence_hunt33.BitwiseXor}, "BitwiseOr": {funcName: "BitwiseOr", native: divergence_hunt33.BitwiseOr}, "ComplexShift": {funcName: "ComplexShift", native: divergence_hunt33.ComplexShift}, "ShiftWithUintAmount": {funcName: "ShiftWithUintAmount", native: divergence_hunt33.ShiftWithUintAmount},
	}})
}
func TestDivergenceHunt34(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt34Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TypedNilSlice": {funcName: "TypedNilSlice", native: divergence_hunt34.TypedNilSlice}, "TypedNilMap": {funcName: "TypedNilMap", native: divergence_hunt34.TypedNilMap}, "TypedNilPointer": {funcName: "TypedNilPointer", native: divergence_hunt34.TypedNilPointer}, "TypedNilFunc": {funcName: "TypedNilFunc", native: divergence_hunt34.TypedNilFunc}, "TypedNilChan": {funcName: "TypedNilChan", native: divergence_hunt34.TypedNilChan}, "InterfaceEqualSame": {funcName: "InterfaceEqualSame", native: divergence_hunt34.InterfaceEqualSame}, "InterfaceEqualDifferent": {funcName: "InterfaceEqualDifferent", native: divergence_hunt34.InterfaceEqualDifferent}, "InterfaceEqualNil": {funcName: "InterfaceEqualNil", native: divergence_hunt34.InterfaceEqualNil}, "TypeSwitchMultiCase": {funcName: "TypeSwitchMultiCase", native: divergence_hunt34.TypeSwitchMultiCase}, "TypeSwitchUintFamily": {funcName: "TypeSwitchUintFamily", native: divergence_hunt34.TypeSwitchUintFamily}, "TypeSwitchFloatFamily": {funcName: "TypeSwitchFloatFamily", native: divergence_hunt34.TypeSwitchFloatFamily}, "AssertToSliceType": {funcName: "AssertToSliceType", native: divergence_hunt34.AssertToSliceType}, "AssertToMapType": {funcName: "AssertToMapType", native: divergence_hunt34.AssertToMapType}, "AssertToFuncType": {funcName: "AssertToFuncType", native: divergence_hunt34.AssertToFuncType}, "FmtTypedNil": {funcName: "FmtTypedNil", native: divergence_hunt34.FmtTypedNil}, "FmtNilInterface": {funcName: "FmtNilInterface", native: divergence_hunt34.FmtNilInterface}, "InterfaceSliceOfTypeSwitch": {funcName: "InterfaceSliceOfTypeSwitch", native: divergence_hunt34.InterfaceSliceOfTypeSwitch}, "NestedTypeSwitch": {funcName: "NestedTypeSwitch", native: divergence_hunt34.NestedTypeSwitch},
	}})
}
func TestDivergenceHunt35(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt35Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"IotaShift": {funcName: "IotaShift", native: divergence_hunt35.IotaShift}, "IotaBitmask": {funcName: "IotaBitmask", native: divergence_hunt35.IotaBitmask}, "ConstExpression": {funcName: "ConstExpression", native: divergence_hunt35.ConstExpression}, "TypeAliasBasic": {funcName: "TypeAliasBasic", native: divergence_hunt35.TypeAliasBasic}, "TypeAliasString": {funcName: "TypeAliasString", native: divergence_hunt35.TypeAliasString}, "TypeAliasArith": {funcName: "TypeAliasArith", native: divergence_hunt35.TypeAliasArith}, "TypeAliasComparison": {funcName: "TypeAliasComparison", native: divergence_hunt35.TypeAliasComparison}, "NestedTypeAlias": {funcName: "NestedTypeAlias", native: divergence_hunt35.NestedTypeAlias}, "ConstBlockBlank": {funcName: "ConstBlockBlank", native: divergence_hunt35.ConstBlockBlank}, "ConstWithString": {funcName: "ConstWithString", native: divergence_hunt35.ConstWithString}, "TypeAliasSlice": {funcName: "TypeAliasSlice", native: divergence_hunt35.TypeAliasSlice}, "TypeAliasMap": {funcName: "TypeAliasMap", native: divergence_hunt35.TypeAliasMap}, "ConstExpressionFloat": {funcName: "ConstExpressionFloat", native: divergence_hunt35.ConstExpressionFloat}, "TypeAliasConversion": {funcName: "TypeAliasConversion", native: divergence_hunt35.TypeAliasConversion}, "ConstArithComplex": {funcName: "ConstArithComplex", native: divergence_hunt35.ConstArithComplex}, "IotaSkip": {funcName: "IotaSkip", native: divergence_hunt35.IotaSkip}, "ConstBitwiseOps": {funcName: "ConstBitwiseOps", native: divergence_hunt35.ConstBitwiseOps},
	}})
}
func TestDivergenceHunt36(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt36Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringByteLen": {funcName: "StringByteLen", native: divergence_hunt36.StringByteLen}, "StringRuneLen": {funcName: "StringRuneLen", native: divergence_hunt36.StringRuneLen}, "StringByteIndex": {funcName: "StringByteIndex", native: divergence_hunt36.StringByteIndex}, "StringSliceMultiByte": {funcName: "StringSliceMultiByte", native: divergence_hunt36.StringSliceMultiByte}, "RuneFromInt": {funcName: "RuneFromInt", native: divergence_hunt36.RuneFromInt}, "StringFromBytes": {funcName: "StringFromBytes", native: divergence_hunt36.StringFromBytes}, "BytesFromString": {funcName: "BytesFromString", native: divergence_hunt36.BytesFromString}, "RuneSliceFromString": {funcName: "RuneSliceFromString", native: divergence_hunt36.RuneSliceFromString}, "StringFromRuneSlice": {funcName: "StringFromRuneSlice", native: divergence_hunt36.StringFromRuneSlice}, "StrconvAtoiNegative": {funcName: "StrconvAtoiNegative", native: divergence_hunt36.StrconvAtoiNegative}, "StrconvItoaNegative": {funcName: "StrconvItoaNegative", native: divergence_hunt36.StrconvItoaNegative}, "StrconvFormatUint": {funcName: "StrconvFormatUint", native: divergence_hunt36.StrconvFormatUint}, "StrconvFormatIntBase": {funcName: "StrconvFormatIntBase", native: divergence_hunt36.StrconvFormatIntBase}, "StringRangeRuneIndex": {funcName: "StringRangeRuneIndex", native: divergence_hunt36.StringRangeRuneIndex}, "StringCompareOps": {funcName: "StringCompareOps", native: divergence_hunt36.StringCompareOps}, "StringConcatMulti": {funcName: "StringConcatMulti", native: divergence_hunt36.StringConcatMulti}, "StringEmptyLen": {funcName: "StringEmptyLen", native: divergence_hunt36.StringEmptyLen}, "StringMultiByteIndex": {funcName: "StringMultiByteIndex", native: divergence_hunt36.StringMultiByteIndex}, "RuneValue": {funcName: "RuneValue", native: divergence_hunt36.RuneValue},
	}})
}
func TestDivergenceHunt37(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt37Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"PanicIntRecover": {funcName: "PanicIntRecover", native: divergence_hunt37.PanicIntRecover}, "PanicStringRecover": {funcName: "PanicStringRecover", native: divergence_hunt37.PanicStringRecover}, "PanicFloatRecover": {funcName: "PanicFloatRecover", native: divergence_hunt37.PanicFloatRecover}, "PanicBoolRecover": {funcName: "PanicBoolRecover", native: divergence_hunt37.PanicBoolRecover}, "PanicInt32Recover": {funcName: "PanicInt32Recover", native: divergence_hunt37.PanicInt32Recover}, "PanicUint8Recover": {funcName: "PanicUint8Recover", native: divergence_hunt37.PanicUint8Recover}, "PanicSliceRecover": {funcName: "PanicSliceRecover", native: divergence_hunt37.PanicSliceRecover}, "PanicMapRecover": {funcName: "PanicMapRecover", native: divergence_hunt37.PanicMapRecover}, "RecoverInMultipleDefers": {funcName: "RecoverInMultipleDefers", native: divergence_hunt37.RecoverInMultipleDefers}, "RecoverTypeSwitch": {funcName: "RecoverTypeSwitch", native: divergence_hunt37.RecoverTypeSwitch}, "PanicInNestedFunc": {funcName: "PanicInNestedFunc", native: divergence_hunt37.PanicInNestedFunc}, "PanicWithNilInterface": {funcName: "PanicWithNilInterface", native: divergence_hunt37.PanicWithNilInterface},
	}})
}
func TestDivergenceHunt38(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt38Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"DeepEmbedding": {funcName: "DeepEmbedding", native: divergence_hunt38.DeepEmbedding}, "EmbeddingFieldAccess": {funcName: "EmbeddingFieldAccess", native: divergence_hunt38.EmbeddingFieldAccess}, "EmbeddedMethodAccess": {funcName: "EmbeddedMethodAccess", native: divergence_hunt38.EmbeddedMethodAccess}, "NestedStructField": {funcName: "NestedStructField", native: divergence_hunt38.NestedStructField}, "StructWithSliceField": {funcName: "StructWithSliceField", native: divergence_hunt38.StructWithSliceField}, "StructWithMapField": {funcName: "StructWithMapField", native: divergence_hunt38.StructWithMapField}, "StructWithPointerField": {funcName: "StructWithPointerField", native: divergence_hunt38.StructWithPointerField}, "StructWithFuncField": {funcName: "StructWithFuncField", native: divergence_hunt38.StructWithFuncField}, "StructWithArrayField": {funcName: "StructWithArrayField", native: divergence_hunt38.StructWithArrayField}, "StructWithChanField": {funcName: "StructWithChanField", native: divergence_hunt38.StructWithChanField}, "StructWithInterfaceField": {funcName: "StructWithInterfaceField", native: divergence_hunt38.StructWithInterfaceField}, "StructJSONRoundTrip": {funcName: "StructJSONRoundTrip", native: divergence_hunt38.StructJSONRoundTrip}, "StructSliceJSONRoundTrip": {funcName: "StructSliceJSONRoundTrip", native: divergence_hunt38.StructSliceJSONRoundTrip}, "FmtNestedStruct": {funcName: "FmtNestedStruct", native: divergence_hunt38.FmtNestedStruct}, "StructComparisonEqual": {funcName: "StructComparisonEqual", native: divergence_hunt38.StructComparisonEqual}, "StructComparisonNotEqual": {funcName: "StructComparisonNotEqual", native: divergence_hunt38.StructComparisonNotEqual},
	}})
}
func TestDivergenceHunt39(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt39Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ResliceAlias": {funcName: "ResliceAlias", native: divergence_hunt39.ResliceAlias}, "ResliceCap": {funcName: "ResliceCap", native: divergence_hunt39.ResliceCap}, "ThreeIndexSlice": {funcName: "ThreeIndexSlice", native: divergence_hunt39.ThreeIndexSlice}, "ThreeIndexSliceNoAlias": {funcName: "ThreeIndexSliceNoAlias", native: divergence_hunt39.ThreeIndexSliceNoAlias}, "AppendNil": {funcName: "AppendNil", native: divergence_hunt39.AppendNil}, "AppendToEmpty": {funcName: "AppendToEmpty", native: divergence_hunt39.AppendToEmpty}, "AppendSliceSpread": {funcName: "AppendSliceSpread", native: divergence_hunt39.AppendSliceSpread}, "CopySlice": {funcName: "CopySlice", native: divergence_hunt39.CopySlice}, "CopyPartial": {funcName: "CopyPartial", native: divergence_hunt39.CopyPartial}, "NilSliceLenCap": {funcName: "NilSliceLenCap", native: divergence_hunt39.NilSliceLenCap}, "EmptySliceLenCap": {funcName: "EmptySliceLenCap", native: divergence_hunt39.EmptySliceLenCap}, "NilSliceCompare": {funcName: "NilSliceCompare", native: divergence_hunt39.NilSliceCompare}, "EmptySliceNotNIl": {funcName: "EmptySliceNotNIl", native: divergence_hunt39.EmptySliceNotNIl}, "SliceMakeWithCap": {funcName: "SliceMakeWithCap", native: divergence_hunt39.SliceMakeWithCap}, "SliceMakeWithLen": {funcName: "SliceMakeWithLen", native: divergence_hunt39.SliceMakeWithLen}, "SliceOfString": {funcName: "SliceOfString", native: divergence_hunt39.SliceOfString}, "SliceOfBool": {funcName: "SliceOfBool", native: divergence_hunt39.SliceOfBool}, "ByteSliceOperations": {funcName: "ByteSliceOperations", native: divergence_hunt39.ByteSliceOperations}, "SliceDeletePattern": {funcName: "SliceDeletePattern", native: divergence_hunt39.SliceDeletePattern}, "JSONRoundTripSlice": {funcName: "JSONRoundTripSlice", native: divergence_hunt39.JSONRoundTripSlice}, "StringSliceJoin": {funcName: "StringSliceJoin", native: divergence_hunt39.StringSliceJoin},
	}})
}
func TestDivergenceHunt40(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt40Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapIntKey": {funcName: "MapIntKey", native: divergence_hunt40.MapIntKey}, "MapFloatKey": {funcName: "MapFloatKey", native: divergence_hunt40.MapFloatKey}, "MapBoolKey": {funcName: "MapBoolKey", native: divergence_hunt40.MapBoolKey}, "MapStructKey": {funcName: "MapStructKey", native: divergence_hunt40.MapStructKey}, "MapStringKey": {funcName: "MapStringKey", native: divergence_hunt40.MapStringKey}, "MapWithSliceValue": {funcName: "MapWithSliceValue", native: divergence_hunt40.MapWithSliceValue}, "MapWithMapValue": {funcName: "MapWithMapValue", native: divergence_hunt40.MapWithMapValue}, "MapDeleteAndLen": {funcName: "MapDeleteAndLen", native: divergence_hunt40.MapDeleteAndLen}, "MapDeleteNonExistent": {funcName: "MapDeleteNonExistent", native: divergence_hunt40.MapDeleteNonExistent}, "MapNilDelete": {funcName: "MapNilDelete", native: divergence_hunt40.MapNilDelete}, "MapNilLookup": {funcName: "MapNilLookup", native: divergence_hunt40.MapNilLookup}, "MapCommaOkPresent": {funcName: "MapCommaOkPresent", native: divergence_hunt40.MapCommaOkPresent}, "MapCommaOkMissing": {funcName: "MapCommaOkMissing", native: divergence_hunt40.MapCommaOkMissing}, "MapOverwrite": {funcName: "MapOverwrite", native: divergence_hunt40.MapOverwrite}, "MapIterationSum": {funcName: "MapIterationSum", native: divergence_hunt40.MapIterationSum}, "MapMakeWithSize": {funcName: "MapMakeWithSize", native: divergence_hunt40.MapMakeWithSize}, "MapEmptyLiteral": {funcName: "MapEmptyLiteral", native: divergence_hunt40.MapEmptyLiteral}, "JSONRoundTripMap": {funcName: "JSONRoundTripMap", native: divergence_hunt40.JSONRoundTripMap}, "FmtMap": {funcName: "FmtMap", native: divergence_hunt40.FmtMap}, "SortMapKeys": {funcName: "SortMapKeys", native: divergence_hunt40.SortMapKeys}, "MapStringJoin": {funcName: "MapStringJoin", native: divergence_hunt40.MapStringJoin},
	}})
}
func TestDivergenceHunt41(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt41Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ConfigParsing": {funcName: "ConfigParsing", native: divergence_hunt41.ConfigParsing}, "CSVLineParse": {funcName: "CSVLineParse", native: divergence_hunt41.CSVLineParse}, "TemplateSubstitution": {funcName: "TemplateSubstitution", native: divergence_hunt41.TemplateSubstitution}, "URLParse": {funcName: "URLParse", native: divergence_hunt41.URLParse}, "DataTransform": {funcName: "DataTransform", native: divergence_hunt41.DataTransform}, "JSONConfigParse": {funcName: "JSONConfigParse", native: divergence_hunt41.JSONConfigParse}, "StringTemplateBuilder": {funcName: "StringTemplateBuilder", native: divergence_hunt41.StringTemplateBuilder}, "NumberFormatter": {funcName: "NumberFormatter", native: divergence_hunt41.NumberFormatter}, "MapReducePattern": {funcName: "MapReducePattern", native: divergence_hunt41.MapReducePattern}, "PipelinePattern": {funcName: "PipelinePattern", native: divergence_hunt41.PipelinePattern}, "ErrorChainPattern": {funcName: "ErrorChainPattern", native: divergence_hunt41.ErrorChainPattern}, "BuilderPattern": {funcName: "BuilderPattern", native: divergence_hunt41.BuilderPattern}, "RateLimiterPattern": {funcName: "RateLimiterPattern", native: divergence_hunt41.RateLimiterPattern}, "RetryPattern": {funcName: "RetryPattern", native: divergence_hunt41.RetryPattern},
	}})
}
func TestDivergenceHunt42(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt42Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Float32AddPrecision": {funcName: "Float32AddPrecision", native: divergence_hunt42.Float32AddPrecision}, "Float64ToIntTrunc": {funcName: "Float64ToIntTrunc", native: divergence_hunt42.Float64ToIntTrunc}, "Float64NegativeTrunc": {funcName: "Float64NegativeTrunc", native: divergence_hunt42.Float64NegativeTrunc}, "IntToFloatRoundTrip": {funcName: "IntToFloatRoundTrip", native: divergence_hunt42.IntToFloatRoundTrip}, "LargeIntToFloat": {funcName: "LargeIntToFloat", native: divergence_hunt42.LargeIntToFloat}, "Uint64ToFloat": {funcName: "Uint64ToFloat", native: divergence_hunt42.Uint64ToFloat}, "Float32ToFloat64": {funcName: "Float32ToFloat64", native: divergence_hunt42.Float32ToFloat64}, "Float64ToFloat32": {funcName: "Float64ToFloat32", native: divergence_hunt42.Float64ToFloat32}, "StrconvParseFloat32": {funcName: "StrconvParseFloat32", native: divergence_hunt42.StrconvParseFloat32}, "StrconvParseFloat64": {funcName: "StrconvParseFloat64", native: divergence_hunt42.StrconvParseFloat64}, "MathRoundEven": {funcName: "MathRoundEven", native: divergence_hunt42.MathRoundEven}, "MathRoundNegative": {funcName: "MathRoundNegative", native: divergence_hunt42.MathRoundNegative}, "FloatCompareNaN": {funcName: "FloatCompareNaN", native: divergence_hunt42.FloatCompareNaN}, "FloatCompareInf": {funcName: "FloatCompareInf", native: divergence_hunt42.FloatCompareInf}, "FloatNegativeZero": {funcName: "FloatNegativeZero", native: divergence_hunt42.FloatNegativeZero}, "FloatAddInf": {funcName: "FloatAddInf", native: divergence_hunt42.FloatAddInf}, "FloatNaNCompare": {funcName: "FloatNaNCompare", native: divergence_hunt42.FloatNaNCompare}, "Int8ToInt16Promotion": {funcName: "Int8ToInt16Promotion", native: divergence_hunt42.Int8ToInt16Promotion}, "Uint8Addition": {funcName: "Uint8Addition", native: divergence_hunt42.Uint8Addition}, "JSONFloatPrecision": {funcName: "JSONFloatPrecision", native: divergence_hunt42.JSONFloatPrecision}, "FmtFloatPrecision": {funcName: "FmtFloatPrecision", native: divergence_hunt42.FmtFloatPrecision},
	}})
}
func TestDivergenceHunt43(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt43Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ClosureCaptureValue": {funcName: "ClosureCaptureValue", native: divergence_hunt43.ClosureCaptureValue}, "ClosureCapturePointer": {funcName: "ClosureCapturePointer", native: divergence_hunt43.ClosureCapturePointer}, "ClosureModifyCaptured": {funcName: "ClosureModifyCaptured", native: divergence_hunt43.ClosureModifyCaptured}, "ClosureReturnFunc": {funcName: "ClosureReturnFunc", native: divergence_hunt43.ClosureReturnFunc}, "ClosureCurry": {funcName: "ClosureCurry", native: divergence_hunt43.ClosureCurry}, "ClosureCounter": {funcName: "ClosureCounter", native: divergence_hunt43.ClosureCounter}, "ClosureAccumulator": {funcName: "ClosureAccumulator", native: divergence_hunt43.ClosureAccumulator}, "ClosureOverSlice": {funcName: "ClosureOverSlice", native: divergence_hunt43.ClosureOverSlice}, "ClosureOverMap": {funcName: "ClosureOverMap", native: divergence_hunt43.ClosureOverMap}, "ClosureOverLoopCopy": {funcName: "ClosureOverLoopCopy", native: divergence_hunt43.ClosureOverLoopCopy}, "ClosureOverLoopNoCopy": {funcName: "ClosureOverLoopNoCopy", native: divergence_hunt43.ClosureOverLoopNoCopy}, "ClosurePartialApplication": {funcName: "ClosurePartialApplication", native: divergence_hunt43.ClosurePartialApplication}, "ClosureFilter": {funcName: "ClosureFilter", native: divergence_hunt43.ClosureFilter}, "ClosureMapFunc": {funcName: "ClosureMapFunc", native: divergence_hunt43.ClosureMapFunc}, "ClosureReduce": {funcName: "ClosureReduce", native: divergence_hunt43.ClosureReduce}, "ClosureInStruct": {funcName: "ClosureInStruct", native: divergence_hunt43.ClosureInStruct}, "ClosureStringProcessor": {funcName: "ClosureStringProcessor", native: divergence_hunt43.ClosureStringProcessor}, "FmtClosure": {funcName: "FmtClosure", native: divergence_hunt43.FmtClosure},
	}})
}
func TestDivergenceHunt44(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt44Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BinarySearch": {funcName: "BinarySearch", native: divergence_hunt44.BinarySearch}, "BubbleSort": {funcName: "BubbleSort", native: divergence_hunt44.BubbleSort}, "InsertionSort": {funcName: "InsertionSort", native: divergence_hunt44.InsertionSort}, "TreeDepth": {funcName: "TreeDepth", native: divergence_hunt44.TreeDepth}, "GraphBFS": {funcName: "GraphBFS", native: divergence_hunt44.GraphBFS}, "LongestCommonSubstrLen": {funcName: "LongestCommonSubstrLen", native: divergence_hunt44.LongestCommonSubstrLen}, "TopKFrequent": {funcName: "TopKFrequent", native: divergence_hunt44.TopKFrequent}, "TwoSum": {funcName: "TwoSum", native: divergence_hunt44.TwoSum}, "MergeSort": {funcName: "MergeSort", native: divergence_hunt44.MergeSort}, "SlidingWindowMax": {funcName: "SlidingWindowMax", native: divergence_hunt44.SlidingWindowMax}, "JSONDataPipeline": {funcName: "JSONDataPipeline", native: divergence_hunt44.JSONDataPipeline}, "StringCompression": {funcName: "StringCompression", native: divergence_hunt44.StringCompression}, "WordFrequency": {funcName: "WordFrequency", native: divergence_hunt44.WordFrequency},
	}})
}
func TestDivergenceHunt45(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt45Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"DeferNamedReturn": {funcName: "DeferNamedReturn", native: divergence_hunt45.DeferNamedReturn}, "DeferMultipleNamedReturn": {funcName: "DeferMultipleNamedReturn", native: divergence_hunt45.DeferMultipleNamedReturn}, "DeferCaptureByValue": {funcName: "DeferCaptureByValue", native: divergence_hunt45.DeferCaptureByValue}, "DeferCaptureByRef": {funcName: "DeferCaptureByRef", native: divergence_hunt45.DeferCaptureByRef}, "DeferInLoop": {funcName: "DeferInLoop", native: divergence_hunt45.DeferInLoop}, "DeferOrder": {funcName: "DeferOrder", native: divergence_hunt45.DeferOrder}, "DeferModifyBeforeReturn": {funcName: "DeferModifyBeforeReturn", native: divergence_hunt45.DeferModifyBeforeReturn}, "DeferWithRecover": {funcName: "DeferWithRecover", native: divergence_hunt45.DeferWithRecover}, "DeferAfterRecover": {funcName: "DeferAfterRecover", native: divergence_hunt45.DeferAfterRecover}, "DeferWithNilRecover": {funcName: "DeferWithNilRecover", native: divergence_hunt45.DeferWithNilRecover}, "NestedDeferRecover": {funcName: "NestedDeferRecover", native: divergence_hunt45.NestedDeferRecover}, "DeferExternalFunc": {funcName: "DeferExternalFunc", native: divergence_hunt45.DeferExternalFunc}, "DeferReturnOverride": {funcName: "DeferReturnOverride", native: divergence_hunt45.DeferReturnOverride}, "DeferConditional": {funcName: "DeferConditional", native: divergence_hunt45.DeferConditional}, "FmtDeferCapture": {funcName: "FmtDeferCapture", native: divergence_hunt45.FmtDeferCapture},
	}})
}
func TestDivergenceHunt46(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt46Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BufferedChannelSendRecv": {funcName: "BufferedChannelSendRecv", native: divergence_hunt46.BufferedChannelSendRecv}, "BufferedChannelLenCap": {funcName: "BufferedChannelLenCap", native: divergence_hunt46.BufferedChannelLenCap}, "ChannelCloseAndRange": {funcName: "ChannelCloseAndRange", native: divergence_hunt46.ChannelCloseAndRange}, "ChannelRecvAfterClose": {funcName: "ChannelRecvAfterClose", native: divergence_hunt46.ChannelRecvAfterClose}, "SelectTwoChannels": {funcName: "SelectTwoChannels", native: divergence_hunt46.SelectTwoChannels}, "SelectDefault": {funcName: "SelectDefault", native: divergence_hunt46.SelectDefault}, "SelectNilChannel": {funcName: "SelectNilChannel", native: divergence_hunt46.SelectNilChannel}, "ChannelDirection": {funcName: "ChannelDirection", native: divergence_hunt46.ChannelDirection}, "ChannelAsSignal": {funcName: "ChannelAsSignal", native: divergence_hunt46.ChannelAsSignal}, "ChannelStruct": {funcName: "ChannelStruct", native: divergence_hunt46.ChannelStruct}, "ChannelSlice": {funcName: "ChannelSlice", native: divergence_hunt46.ChannelSlice}, "ChannelMap": {funcName: "ChannelMap", native: divergence_hunt46.ChannelMap}, "JSONThroughChannel": {funcName: "JSONThroughChannel", native: divergence_hunt46.JSONThroughChannel}, "MultipleSelects": {funcName: "MultipleSelects", native: divergence_hunt46.MultipleSelects}, "FmtChannel": {funcName: "FmtChannel", native: divergence_hunt46.FmtChannel}, "SortThroughChannel": {funcName: "SortThroughChannel", native: divergence_hunt46.SortThroughChannel}, "StringsThroughChannel": {funcName: "StringsThroughChannel", native: divergence_hunt46.StringsThroughChannel},
	}})
}
func TestDivergenceHunt47(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt47Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SimpleErrorCheck": {funcName: "SimpleErrorCheck", native: divergence_hunt47.SimpleErrorCheck}, "ErrorPropagation": {funcName: "ErrorPropagation", native: divergence_hunt47.ErrorPropagation}, "ErrorInClosure": {funcName: "ErrorInClosure", native: divergence_hunt47.ErrorInClosure}, "ErrorChain": {funcName: "ErrorChain", native: divergence_hunt47.ErrorChain}, "ValidationError": {funcName: "ValidationError", native: divergence_hunt47.ValidationError}, "MultiErrorCollect": {funcName: "MultiErrorCollect", native: divergence_hunt47.MultiErrorCollect}, "PanicInsteadOfError": {funcName: "PanicInsteadOfError", native: divergence_hunt47.PanicInsteadOfError}, "ErrorTypeAssertion": {funcName: "ErrorTypeAssertion", native: divergence_hunt47.ErrorTypeAssertion}, "JSONUnmarshalError": {funcName: "JSONUnmarshalError", native: divergence_hunt47.JSONUnmarshalError}, "FmtErrorfWrap": {funcName: "FmtErrorfWrap", native: divergence_hunt47.FmtErrorfWrap}, "ErrorStringMethod": {funcName: "ErrorStringMethod", native: divergence_hunt47.ErrorStringMethod}, "SortWithValidation": {funcName: "SortWithValidation", native: divergence_hunt47.SortWithValidation}, "StringsErrorCheck": {funcName: "StringsErrorCheck", native: divergence_hunt47.StringsErrorCheck},
	}})
}
func TestDivergenceHunt48(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt48Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TrimSpace": {funcName: "TrimSpace", native: divergence_hunt48.TrimSpace}, "TrimPrefix": {funcName: "TrimPrefix", native: divergence_hunt48.TrimPrefix}, "TrimSuffix": {funcName: "TrimSuffix", native: divergence_hunt48.TrimSuffix}, "SplitN": {funcName: "SplitN", native: divergence_hunt48.SplitN}, "SplitAfter": {funcName: "SplitAfter", native: divergence_hunt48.SplitAfter}, "ReplaceN": {funcName: "ReplaceN", native: divergence_hunt48.ReplaceN}, "ReplaceAll": {funcName: "ReplaceAll", native: divergence_hunt48.ReplaceAll}, "Repeat": {funcName: "Repeat", native: divergence_hunt48.Repeat}, "Contains": {funcName: "Contains", native: divergence_hunt48.Contains}, "ContainsAny": {funcName: "ContainsAny", native: divergence_hunt48.ContainsAny}, "HasPrefix": {funcName: "HasPrefix", native: divergence_hunt48.HasPrefix}, "HasSuffix": {funcName: "HasSuffix", native: divergence_hunt48.HasSuffix}, "IndexFunc": {funcName: "IndexFunc", native: divergence_hunt48.IndexFunc}, "TitleCase": {funcName: "TitleCase", native: divergence_hunt48.TitleCase}, "ToTitle": {funcName: "ToTitle", native: divergence_hunt48.ToTitle}, "MapFunc": {funcName: "MapFunc", native: divergence_hunt48.MapFunc}, "BuilderString": {funcName: "BuilderString", native: divergence_hunt48.BuilderString}, "BuilderLen": {funcName: "BuilderLen", native: divergence_hunt48.BuilderLen}, "StrconvQuote": {funcName: "StrconvQuote", native: divergence_hunt48.StrconvQuote}, "FmtStringOps": {funcName: "FmtStringOps", native: divergence_hunt48.FmtStringOps},
	}})
}
func TestDivergenceHunt49(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt49Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NewInt": {funcName: "NewInt", native: divergence_hunt49.NewInt}, "NewStruct": {funcName: "NewStruct", native: divergence_hunt49.NewStruct}, "AddressOf": {funcName: "AddressOf", native: divergence_hunt49.AddressOf}, "AddressOfModify": {funcName: "AddressOfModify", native: divergence_hunt49.AddressOfModify}, "PointerToSlice": {funcName: "PointerToSlice", native: divergence_hunt49.PointerToSlice}, "PointerToMap": {funcName: "PointerToMap", native: divergence_hunt49.PointerToMap}, "PointerToStruct": {funcName: "PointerToStruct", native: divergence_hunt49.PointerToStruct}, "DoublePointer": {funcName: "DoublePointer", native: divergence_hunt49.DoublePointer}, "NilPointerComparison": {funcName: "NilPointerComparison", native: divergence_hunt49.NilPointerComparison}, "PointerComparison": {funcName: "PointerComparison", native: divergence_hunt49.PointerComparison}, "PointerSlice": {funcName: "PointerSlice", native: divergence_hunt49.PointerSlice}, "PointerArray": {funcName: "PointerArray", native: divergence_hunt49.PointerArray}, "StructPointerMethod": {funcName: "StructPointerMethod", native: divergence_hunt49.StructPointerMethod}, "JSONPointerRoundTrip": {funcName: "JSONPointerRoundTrip", native: divergence_hunt49.JSONPointerRoundTrip}, "FmtPointer": {funcName: "FmtPointer", native: divergence_hunt49.FmtPointer}, "PointerSwap": {funcName: "PointerSwap", native: divergence_hunt49.PointerSwap},
	}})
}
func TestDivergenceHunt50(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt50Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StudentRanking": {funcName: "StudentRanking", native: divergence_hunt50.StudentRanking}, "TextAnalyzer": {funcName: "TextAnalyzer", native: divergence_hunt50.TextAnalyzer}, "ShoppingCart": {funcName: "ShoppingCart", native: divergence_hunt50.ShoppingCart}, "JSONAPIResponse": {funcName: "JSONAPIResponse", native: divergence_hunt50.JSONAPIResponse}, "MatrixRotate": {funcName: "MatrixRotate", native: divergence_hunt50.MatrixRotate}, "DataDedup": {funcName: "DataDedup", native: divergence_hunt50.DataDedup}, "StringTemplate": {funcName: "StringTemplate", native: divergence_hunt50.StringTemplate}, "GroupByCategory": {funcName: "GroupByCategory", native: divergence_hunt50.GroupByCategory}, "LRUPrototype": {funcName: "LRUPrototype", native: divergence_hunt50.LRUPrototype}, "FibonacciMemo": {funcName: "FibonacciMemo", native: divergence_hunt50.FibonacciMemo}, "FmtTable": {funcName: "FmtTable", native: divergence_hunt50.FmtTable}, "Pipeline": {funcName: "Pipeline", native: divergence_hunt50.Pipeline},
	}})
}
func TestDivergenceHunt51(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt51Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"IntToInt8": {funcName: "IntToInt8", native: divergence_hunt51.IntToInt8}, "IntToInt16": {funcName: "IntToInt16", native: divergence_hunt51.IntToInt16}, "IntToUint": {funcName: "IntToUint", native: divergence_hunt51.IntToUint}, "UintToInt": {funcName: "UintToInt", native: divergence_hunt51.UintToInt}, "Float32ToInt": {funcName: "Float32ToInt", native: divergence_hunt51.Float32ToInt}, "Float64ToInt": {funcName: "Float64ToInt", native: divergence_hunt51.Float64ToInt}, "IntToFloat32": {funcName: "IntToFloat32", native: divergence_hunt51.IntToFloat32}, "IntToFloat64": {funcName: "IntToFloat64", native: divergence_hunt51.IntToFloat64}, "RuneToString": {funcName: "RuneToString", native: divergence_hunt51.RuneToString}, "IntRuneToString": {funcName: "IntRuneToString", native: divergence_hunt51.IntRuneToString}, "BoolToInt": {funcName: "BoolToInt", native: divergence_hunt51.BoolToInt}, "ByteToString": {funcName: "ByteToString", native: divergence_hunt51.ByteToString}, "BytesToString": {funcName: "BytesToString", native: divergence_hunt51.BytesToString}, "StringToBytes": {funcName: "StringToBytes", native: divergence_hunt51.StringToBytes}, "RunesToString": {funcName: "RunesToString", native: divergence_hunt51.RunesToString}, "StringToRunes": {funcName: "StringToRunes", native: divergence_hunt51.StringToRunes}, "Int64ToInt32": {funcName: "Int64ToInt32", native: divergence_hunt51.Int64ToInt32}, "Uint32ToInt32": {funcName: "Uint32ToInt32", native: divergence_hunt51.Uint32ToInt32}, "FmtConversion": {funcName: "FmtConversion", native: divergence_hunt51.FmtConversion},
	}})
}
func TestDivergenceHunt52(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt52Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SortInts": {funcName: "SortInts", native: divergence_hunt52.SortInts}, "SortStrings": {funcName: "SortStrings", native: divergence_hunt52.SortStrings}, "SortFloat64s": {funcName: "SortFloat64s", native: divergence_hunt52.SortFloat64s}, "SortIntSlice": {funcName: "SortIntSlice", native: divergence_hunt52.SortIntSlice}, "SortReverse": {funcName: "SortReverse", native: divergence_hunt52.SortReverse}, "SortStructSlice": {funcName: "SortStructSlice", native: divergence_hunt52.SortStructSlice}, "SortSliceDesc": {funcName: "SortSliceDesc", native: divergence_hunt52.SortSliceDesc}, "SortStable": {funcName: "SortStable", native: divergence_hunt52.SortStable}, "SortSearch": {funcName: "SortSearch", native: divergence_hunt52.SortSearch}, "SortSearchString": {funcName: "SortSearchString", native: divergence_hunt52.SortSearchString}, "SortFloat64Search": {funcName: "SortFloat64Search", native: divergence_hunt52.SortFloat64Search}, "SortIsSorted": {funcName: "SortIsSorted", native: divergence_hunt52.SortIsSorted}, "SortEmptySlice": {funcName: "SortEmptySlice", native: divergence_hunt52.SortEmptySlice}, "SortSingleElement": {funcName: "SortSingleElement", native: divergence_hunt52.SortSingleElement}, "SortDuplicate": {funcName: "SortDuplicate", native: divergence_hunt52.SortDuplicate},
	}})
}
func TestDivergenceHunt53(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt53Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"JSONEncodeInt": {funcName: "JSONEncodeInt", native: divergence_hunt53.JSONEncodeInt}, "JSONEncodeString": {funcName: "JSONEncodeString", native: divergence_hunt53.JSONEncodeString}, "JSONDecodeInt": {funcName: "JSONDecodeInt", native: divergence_hunt53.JSONDecodeInt}, "JSONDecodeString": {funcName: "JSONDecodeString", native: divergence_hunt53.JSONDecodeString}, "JSONSlice": {funcName: "JSONSlice", native: divergence_hunt53.JSONSlice}, "JSONMap": {funcName: "JSONMap", native: divergence_hunt53.JSONMap}, "JSONBool": {funcName: "JSONBool", native: divergence_hunt53.JSONBool}, "JSONNull": {funcName: "JSONNull", native: divergence_hunt53.JSONNull}, "RegexMatch": {funcName: "RegexMatch", native: divergence_hunt53.RegexMatch}, "RegexFind": {funcName: "RegexFind", native: divergence_hunt53.RegexFind}, "RegexFindAll": {funcName: "RegexFindAll", native: divergence_hunt53.RegexFindAll}, "RegexReplace": {funcName: "RegexReplace", native: divergence_hunt53.RegexReplace}, "RegexSplit": {funcName: "RegexSplit", native: divergence_hunt53.RegexSplit}, "RegexSubmatch": {funcName: "RegexSubmatch", native: divergence_hunt53.RegexSubmatch}, "RegexNamedGroup": {funcName: "RegexNamedGroup", native: divergence_hunt53.RegexNamedGroup}, "StringParse": {funcName: "StringParse", native: divergence_hunt53.StringParse}, "CSVParse": {funcName: "CSVParse", native: divergence_hunt53.CSVParse}, "TemplateParse": {funcName: "TemplateParse", native: divergence_hunt53.TemplateParse},
	}})
}
func TestDivergenceHunt54(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt54Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MathAbs": {funcName: "MathAbs", native: divergence_hunt54.MathAbs}, "MathCeil": {funcName: "MathCeil", native: divergence_hunt54.MathCeil}, "MathFloor": {funcName: "MathFloor", native: divergence_hunt54.MathFloor}, "MathRound": {funcName: "MathRound", native: divergence_hunt54.MathRound}, "MathMax": {funcName: "MathMax", native: divergence_hunt54.MathMax}, "MathMin": {funcName: "MathMin", native: divergence_hunt54.MathMin}, "MathPow": {funcName: "MathPow", native: divergence_hunt54.MathPow}, "MathSqrt": {funcName: "MathSqrt", native: divergence_hunt54.MathSqrt}, "MathMod": {funcName: "MathMod", native: divergence_hunt54.MathMod}, "MathLog": {funcName: "MathLog", native: divergence_hunt54.MathLog}, "MathLog2": {funcName: "MathLog2", native: divergence_hunt54.MathLog2}, "MathLog10": {funcName: "MathLog10", native: divergence_hunt54.MathLog10}, "MathExp": {funcName: "MathExp", native: divergence_hunt54.MathExp}, "MathSin": {funcName: "MathSin", native: divergence_hunt54.MathSin}, "MathCos": {funcName: "MathCos", native: divergence_hunt54.MathCos}, "MathHypot": {funcName: "MathHypot", native: divergence_hunt54.MathHypot}, "MathIsNaN": {funcName: "MathIsNaN", native: divergence_hunt54.MathIsNaN}, "MathIsInf": {funcName: "MathIsInf", native: divergence_hunt54.MathIsInf}, "MathSignbit": {funcName: "MathSignbit", native: divergence_hunt54.MathSignbit}, "StrconvAtoi": {funcName: "StrconvAtoi", native: divergence_hunt54.StrconvAtoi}, "StrconvItoa": {funcName: "StrconvItoa", native: divergence_hunt54.StrconvItoa}, "StrconvFormatFloat": {funcName: "StrconvFormatFloat", native: divergence_hunt54.StrconvFormatFloat}, "StrconvParseFloat": {funcName: "StrconvParseFloat", native: divergence_hunt54.StrconvParseFloat}, "FmtFloat": {funcName: "FmtFloat", native: divergence_hunt54.FmtFloat}, "FmtInt": {funcName: "FmtInt", native: divergence_hunt54.FmtInt},
	}})
}
func TestDivergenceHunt55(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt55Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SwitchBasic": {funcName: "SwitchBasic", native: divergence_hunt55.SwitchBasic}, "SwitchDefault": {funcName: "SwitchDefault", native: divergence_hunt55.SwitchDefault}, "SwitchMultiCase": {funcName: "SwitchMultiCase", native: divergence_hunt55.SwitchMultiCase}, "SwitchWithInit": {funcName: "SwitchWithInit", native: divergence_hunt55.SwitchWithInit}, "SwitchExpression": {funcName: "SwitchExpression", native: divergence_hunt55.SwitchExpression}, "SwitchFallthrough": {funcName: "SwitchFallthrough", native: divergence_hunt55.SwitchFallthrough}, "SwitchNoMatch": {funcName: "SwitchNoMatch", native: divergence_hunt55.SwitchNoMatch}, "NestedSwitch": {funcName: "NestedSwitch", native: divergence_hunt55.NestedSwitch}, "LabeledBreak": {funcName: "LabeledBreak", native: divergence_hunt55.LabeledBreak}, "LabeledContinue": {funcName: "LabeledContinue", native: divergence_hunt55.LabeledContinue}, "ForBreak": {funcName: "ForBreak", native: divergence_hunt55.ForBreak}, "ForContinue": {funcName: "ForContinue", native: divergence_hunt55.ForContinue}, "InfiniteLoopBreak": {funcName: "InfiniteLoopBreak", native: divergence_hunt55.InfiniteLoopBreak}, "RangeBreak": {funcName: "RangeBreak", native: divergence_hunt55.RangeBreak}, "RangeContinue": {funcName: "RangeContinue", native: divergence_hunt55.RangeContinue}, "FmtSwitch": {funcName: "FmtSwitch", native: divergence_hunt55.FmtSwitch}, "StringsSwitch": {funcName: "StringsSwitch", native: divergence_hunt55.StringsSwitch},
	}})
}
func TestDivergenceHunt56(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt56Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SalesReport": {funcName: "SalesReport", native: divergence_hunt56.SalesReport}, "SalesByRegion": {funcName: "SalesByRegion", native: divergence_hunt56.SalesByRegion}, "TopProducts": {funcName: "TopProducts", native: divergence_hunt56.TopProducts}, "DataCleaning": {funcName: "DataCleaning", native: divergence_hunt56.DataCleaning}, "DataNormalization": {funcName: "DataNormalization", native: divergence_hunt56.DataNormalization}, "JSONDataExport": {funcName: "JSONDataExport", native: divergence_hunt56.JSONDataExport}, "PivotTable": {funcName: "PivotTable", native: divergence_hunt56.PivotTable}, "PercentileCalc": {funcName: "PercentileCalc", native: divergence_hunt56.PercentileCalc}, "MovingAverage": {funcName: "MovingAverage", native: divergence_hunt56.MovingAverage}, "FrequencyDistribution": {funcName: "FrequencyDistribution", native: divergence_hunt56.FrequencyDistribution}, "DataMerge": {funcName: "DataMerge", native: divergence_hunt56.DataMerge}, "StringReport": {funcName: "StringReport", native: divergence_hunt56.StringReport},
	}})
}
func TestDivergenceHunt57(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt57Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MutexBasic": {funcName: "MutexBasic", native: divergence_hunt57.MutexBasic}, "MutexInDefer": {funcName: "MutexInDefer", native: divergence_hunt57.MutexInDefer}, "RWMutexBasic": {funcName: "RWMutexBasic", native: divergence_hunt57.RWMutexBasic}, "OnceBasic": {funcName: "OnceBasic", native: divergence_hunt57.OnceBasic}, "MutexCounter": {funcName: "MutexCounter", native: divergence_hunt57.MutexCounter}, "OnceInClosure": {funcName: "OnceInClosure", native: divergence_hunt57.OnceInClosure}, "RWMutexMultipleReaders": {funcName: "RWMutexMultipleReaders", native: divergence_hunt57.RWMutexMultipleReaders}, "MutexSwapPattern": {funcName: "MutexSwapPattern", native: divergence_hunt57.MutexSwapPattern}, "MutexMapProtect": {funcName: "MutexMapProtect", native: divergence_hunt57.MutexMapProtect}, "OnceLazyInit": {funcName: "OnceLazyInit", native: divergence_hunt57.OnceLazyInit}, "FmtMutex": {funcName: "FmtMutex", native: divergence_hunt57.FmtMutex}, "MutexNestedLock": {funcName: "MutexNestedLock", native: divergence_hunt57.MutexNestedLock},
	}})
}
func TestDivergenceHunt58(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt58Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"FmtSprintfInt": {funcName: "FmtSprintfInt", native: divergence_hunt58.FmtSprintfInt}, "FmtSprintfFloat": {funcName: "FmtSprintfFloat", native: divergence_hunt58.FmtSprintfFloat}, "FmtSprintfString": {funcName: "FmtSprintfString", native: divergence_hunt58.FmtSprintfString}, "FmtSprintfBool": {funcName: "FmtSprintfBool", native: divergence_hunt58.FmtSprintfBool}, "FmtSprintfWidth": {funcName: "FmtSprintfWidth", native: divergence_hunt58.FmtSprintfWidth}, "FmtSprintfHex": {funcName: "FmtSprintfHex", native: divergence_hunt58.FmtSprintfHex}, "FmtSprintfOctal": {funcName: "FmtSprintfOctal", native: divergence_hunt58.FmtSprintfOctal}, "FmtSprintfBinary": {funcName: "FmtSprintfBinary", native: divergence_hunt58.FmtSprintfBinary}, "FmtSprintfChar": {funcName: "FmtSprintfChar", native: divergence_hunt58.FmtSprintfChar}, "FmtSprintfPadZero": {funcName: "FmtSprintfPadZero", native: divergence_hunt58.FmtSprintfPadZero}, "FmtSprintfQuoted": {funcName: "FmtSprintfQuoted", native: divergence_hunt58.FmtSprintfQuoted}, "FmtSprintfDefault": {funcName: "FmtSprintfDefault", native: divergence_hunt58.FmtSprintfDefault}, "FmtErrorf": {funcName: "FmtErrorf", native: divergence_hunt58.FmtErrorf}, "StrconvAtoiPositive": {funcName: "StrconvAtoiPositive", native: divergence_hunt58.StrconvAtoiPositive}, "StrconvAtoiNegative": {funcName: "StrconvAtoiNegative", native: divergence_hunt58.StrconvAtoiNegative}, "StrconvItoaPositive": {funcName: "StrconvItoaPositive", native: divergence_hunt58.StrconvItoaPositive}, "StrconvItoaNegative": {funcName: "StrconvItoaNegative", native: divergence_hunt58.StrconvItoaNegative}, "StrconvFormatBool": {funcName: "StrconvFormatBool", native: divergence_hunt58.StrconvFormatBool}, "StrconvParseBool": {funcName: "StrconvParseBool", native: divergence_hunt58.StrconvParseBool}, "StrconvFormatFloat": {funcName: "StrconvFormatFloat", native: divergence_hunt58.StrconvFormatFloat}, "StrconvParseFloat": {funcName: "StrconvParseFloat", native: divergence_hunt58.StrconvParseFloat}, "StringBuilderConcat": {funcName: "StringBuilderConcat", native: divergence_hunt58.StringBuilderConcat},
	}})
}
func TestDivergenceHunt59(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt59Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MultiReturnSwap": {funcName: "MultiReturnSwap", native: divergence_hunt59.MultiReturnSwap}, "MultiReturnDivMod": {funcName: "MultiReturnDivMod", native: divergence_hunt59.MultiReturnDivMod}, "MultiReturnMinMax": {funcName: "MultiReturnMinMax", native: divergence_hunt59.MultiReturnMinMax}, "BlankIdentifier": {funcName: "BlankIdentifier", native: divergence_hunt59.BlankIdentifier}, "BlankInMultiReturn": {funcName: "BlankInMultiReturn", native: divergence_hunt59.BlankInMultiReturn}, "VariadicSum": {funcName: "VariadicSum", native: divergence_hunt59.VariadicSum}, "VariadicSpread": {funcName: "VariadicSpread", native: divergence_hunt59.VariadicSpread}, "VariadicEmpty": {funcName: "VariadicEmpty", native: divergence_hunt59.VariadicEmpty}, "VariadicWithRegular": {funcName: "VariadicWithRegular", native: divergence_hunt59.VariadicWithRegular}, "NamedReturnBare": {funcName: "NamedReturnBare", native: divergence_hunt59.NamedReturnBare}, "NamedReturnWithDefer": {funcName: "NamedReturnWithDefer", native: divergence_hunt59.NamedReturnWithDefer}, "MultipleNamedReturn": {funcName: "MultipleNamedReturn", native: divergence_hunt59.MultipleNamedReturn}, "MultiReturnWithInterface": {funcName: "MultiReturnWithInterface", native: divergence_hunt59.MultiReturnWithInterface}, "MultiReturnInLoop": {funcName: "MultiReturnInLoop", native: divergence_hunt59.MultiReturnInLoop}, "BlankInLoop": {funcName: "BlankInLoop", native: divergence_hunt59.BlankInLoop}, "MultiReturnError": {funcName: "MultiReturnError", native: divergence_hunt59.MultiReturnError},
	}})
}
func TestDivergenceHunt60(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt60Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Comprehensive1": {funcName: "Comprehensive1", native: divergence_hunt60.Comprehensive1}, "Comprehensive2": {funcName: "Comprehensive2", native: divergence_hunt60.Comprehensive2}, "Comprehensive3": {funcName: "Comprehensive3", native: divergence_hunt60.Comprehensive3}, "Comprehensive4": {funcName: "Comprehensive4", native: divergence_hunt60.Comprehensive4}, "Comprehensive5": {funcName: "Comprehensive5", native: divergence_hunt60.Comprehensive5}, "Comprehensive6": {funcName: "Comprehensive6", native: divergence_hunt60.Comprehensive6}, "Comprehensive7": {funcName: "Comprehensive7", native: divergence_hunt60.Comprehensive7}, "Comprehensive8": {funcName: "Comprehensive8", native: divergence_hunt60.Comprehensive8}, "Comprehensive9": {funcName: "Comprehensive9", native: divergence_hunt60.Comprehensive9}, "Comprehensive10": {funcName: "Comprehensive10", native: divergence_hunt60.Comprehensive10},
	}})
}
func TestDivergenceHunt61(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt61Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Uint8Overflow": {funcName: "Uint8Overflow", native: divergence_hunt61.Uint8Overflow}, "Uint8Underflow": {funcName: "Uint8Underflow", native: divergence_hunt61.Uint8Underflow}, "Int8Overflow": {funcName: "Int8Overflow", native: divergence_hunt61.Int8Overflow}, "Int8Underflow": {funcName: "Int8Underflow", native: divergence_hunt61.Int8Underflow}, "Uint16Overflow": {funcName: "Uint16Overflow", native: divergence_hunt61.Uint16Overflow}, "Uint32Overflow": {funcName: "Uint32Overflow", native: divergence_hunt61.Uint32Overflow}, "Int16Overflow": {funcName: "Int16Overflow", native: divergence_hunt61.Int16Overflow}, "Int16Underflow": {funcName: "Int16Underflow", native: divergence_hunt61.Int16Underflow}, "IntNegateMin": {funcName: "IntNegateMin", native: divergence_hunt61.IntNegateMin}, "UintMulOverflow": {funcName: "UintMulOverflow", native: divergence_hunt61.UintMulOverflow}, "IntDivTruncation": {funcName: "IntDivTruncation", native: divergence_hunt61.IntDivTruncation}, "IntModNegative": {funcName: "IntModNegative", native: divergence_hunt61.IntModNegative}, "ShiftLeftLarge": {funcName: "ShiftLeftLarge", native: divergence_hunt61.ShiftLeftLarge}, "ShiftRightSigned": {funcName: "ShiftRightSigned", native: divergence_hunt61.ShiftRightSigned}, "UintConvertNegative": {funcName: "UintConvertNegative", native: divergence_hunt61.UintConvertNegative}, "IntConvertLargeUint": {funcName: "IntConvertLargeUint", native: divergence_hunt61.IntConvertLargeUint}, "FloatTruncateToInt": {funcName: "FloatTruncateToInt", native: divergence_hunt61.FloatTruncateToInt}, "FloatTruncateNegToInt": {funcName: "FloatTruncateNegToInt", native: divergence_hunt61.FloatTruncateNegToInt}, "ComplexRealImag": {funcName: "ComplexRealImag", native: divergence_hunt61.ComplexRealImag},
	}})
}
func TestDivergenceHunt62(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt62Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TypedNilError": {funcName: "TypedNilError", native: divergence_hunt62.TypedNilError}, "NilInterfaceCheck": {funcName: "NilInterfaceCheck", native: divergence_hunt62.NilInterfaceCheck}, "NilSliceVsEmptySlice": {funcName: "NilSliceVsEmptySlice", native: divergence_hunt62.NilSliceVsEmptySlice}, "NilMapVsEmptyMap": {funcName: "NilMapVsEmptyMap", native: divergence_hunt62.NilMapVsEmptyMap}, "NilChanVsMakeChan": {funcName: "NilChanVsMakeChan", native: divergence_hunt62.NilChanVsMakeChan}, "NilFuncCheck": {funcName: "NilFuncCheck", native: divergence_hunt62.NilFuncCheck}, "NilPointerCheck": {funcName: "NilPointerCheck", native: divergence_hunt62.NilPointerCheck}, "TypeAssertNil": {funcName: "TypeAssertNil", native: divergence_hunt62.TypeAssertNil}, "TypeAssertTypedNil": {funcName: "TypeAssertTypedNil", native: divergence_hunt62.TypeAssertTypedNil}, "FmtTypedNil": {funcName: "FmtTypedNil", native: divergence_hunt62.FmtTypedNil},
	}})
}
func TestDivergenceHunt63(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt63Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapNaNKey": {funcName: "MapNaNKey", native: divergence_hunt63.MapNaNKey}, "MapNaNKeyLookup": {funcName: "MapNaNKeyLookup", native: divergence_hunt63.MapNaNKeyLookup}, "MapStructKey": {funcName: "MapStructKey", native: divergence_hunt63.MapStructKey}, "MapArrayKey": {funcName: "MapArrayKey", native: divergence_hunt63.MapArrayKey}, "MapDeleteDuringRange": {funcName: "MapDeleteDuringRange", native: divergence_hunt63.MapDeleteDuringRange}, "MapDeleteAndReadd": {funcName: "MapDeleteAndReadd", native: divergence_hunt63.MapDeleteAndReadd}, "MapNilDelete": {funcName: "MapNilDelete", native: divergence_hunt63.MapNilDelete}, "MapLenAfterDelete": {funcName: "MapLenAfterDelete", native: divergence_hunt63.MapLenAfterDelete}, "MapZeroValueAccess": {funcName: "MapZeroValueAccess", native: divergence_hunt63.MapZeroValueAccess}, "MapBoolKey": {funcName: "MapBoolKey", native: divergence_hunt63.MapBoolKey}, "MapStringKeyEmpty": {funcName: "MapStringKeyEmpty", native: divergence_hunt63.MapStringKeyEmpty}, "MapIntKeyZero": {funcName: "MapIntKeyZero", native: divergence_hunt63.MapIntKeyZero}, "MapNestedMap": {funcName: "MapNestedMap", native: divergence_hunt63.MapNestedMap}, "MapOverwritePreservesType": {funcName: "MapOverwritePreservesType", native: divergence_hunt63.MapOverwritePreservesType}, "MapCommaOkDelete": {funcName: "MapCommaOkDelete", native: divergence_hunt63.MapCommaOkDelete},
	}})
}
func TestDivergenceHunt64(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt64Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ThreeIndexSlice": {funcName: "ThreeIndexSlice", native: divergence_hunt64.ThreeIndexSlice}, "ThreeIndexSliceFull": {funcName: "ThreeIndexSliceFull", native: divergence_hunt64.ThreeIndexSliceFull}, "SliceAppendWithinCap": {funcName: "SliceAppendWithinCap", native: divergence_hunt64.SliceAppendWithinCap}, "SliceAppendNil": {funcName: "SliceAppendNil", native: divergence_hunt64.SliceAppendNil}, "SliceCopyCount": {funcName: "SliceCopyCount", native: divergence_hunt64.SliceCopyCount}, "SliceCopyFromSub": {funcName: "SliceCopyFromSub", native: divergence_hunt64.SliceCopyFromSub}, "SliceAppendGrow": {funcName: "SliceAppendGrow", native: divergence_hunt64.SliceAppendGrow}, "SliceEmptySubslice": {funcName: "SliceEmptySubslice", native: divergence_hunt64.SliceEmptySubslice}, "SliceCapAfterAppend": {funcName: "SliceCapAfterAppend", native: divergence_hunt64.SliceCapAfterAppend}, "SliceMakeZeroLen": {funcName: "SliceMakeZeroLen", native: divergence_hunt64.SliceMakeZeroLen}, "SliceOfString": {funcName: "SliceOfString", native: divergence_hunt64.SliceOfString}, "SliceOfBool": {funcName: "SliceOfBool", native: divergence_hunt64.SliceOfBool}, "SliceOverlappingCopy": {funcName: "SliceOverlappingCopy", native: divergence_hunt64.SliceOverlappingCopy}, "SliceDoubleAppend": {funcName: "SliceDoubleAppend", native: divergence_hunt64.SliceDoubleAppend},
	}})
}
func TestDivergenceHunt65(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt65Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ClosureCounter": {funcName: "ClosureCounter", native: divergence_hunt65.ClosureCounter}, "ClosureSharedState": {funcName: "ClosureSharedState", native: divergence_hunt65.ClosureSharedState}, "ClosureChain": {funcName: "ClosureChain", native: divergence_hunt65.ClosureChain}, "ClosureOverLoopVar": {funcName: "ClosureOverLoopVar", native: divergence_hunt65.ClosureOverLoopVar}, "RecursiveClosure": {funcName: "RecursiveClosure", native: divergence_hunt65.RecursiveClosure}, "ClosureReturnClosure": {funcName: "ClosureReturnClosure", native: divergence_hunt65.ClosureReturnClosure}, "ClosureCaptureSlice": {funcName: "ClosureCaptureSlice", native: divergence_hunt65.ClosureCaptureSlice}, "ClosureCaptureMap": {funcName: "ClosureCaptureMap", native: divergence_hunt65.ClosureCaptureMap}, "ClosureMultipleReturns": {funcName: "ClosureMultipleReturns", native: divergence_hunt65.ClosureMultipleReturns}, "ClosureCurry": {funcName: "ClosureCurry", native: divergence_hunt65.ClosureCurry}, "ClosureCaptureModify": {funcName: "ClosureCaptureModify", native: divergence_hunt65.ClosureCaptureModify}, "ClosureNoCapture": {funcName: "ClosureNoCapture", native: divergence_hunt65.ClosureNoCapture}, "ClosureAsArg": {funcName: "ClosureAsArg", native: divergence_hunt65.ClosureAsArg}, "ClosureSliceMap": {funcName: "ClosureSliceMap", native: divergence_hunt65.ClosureSliceMap},
	}})
}
func TestDivergenceHunt66(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt66Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RecoverInNestedDefer": {funcName: "RecoverInNestedDefer", native: divergence_hunt66.RecoverInNestedDefer}, "PanicAfterRecover": {funcName: "PanicAfterRecover", native: divergence_hunt66.PanicAfterRecover}, "RecoverReturnsNilAfterCall": {funcName: "RecoverReturnsNilAfterCall", native: divergence_hunt66.RecoverReturnsNilAfterCall}, "MultipleRecoverSameDefer": {funcName: "MultipleRecoverSameDefer", native: divergence_hunt66.MultipleRecoverSameDefer}, "RecoverOnlyInDefer": {funcName: "RecoverOnlyInDefer", native: divergence_hunt66.RecoverOnlyInDefer}, "NestedPanicRecover": {funcName: "NestedPanicRecover", native: divergence_hunt66.NestedPanicRecover}, "PanicString": {funcName: "PanicString", native: divergence_hunt66.PanicString}, "PanicNilInterface": {funcName: "PanicNilInterface", native: divergence_hunt66.PanicNilInterface}, "DeferPanicOrder": {funcName: "DeferPanicOrder", native: divergence_hunt66.DeferPanicOrder}, "RecoverTypeAssertion": {funcName: "RecoverTypeAssertion", native: divergence_hunt66.RecoverTypeAssertion},
	}})
}
func TestDivergenceHunt67(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt67Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringLenBytes": {funcName: "StringLenBytes", native: divergence_hunt67.StringLenBytes}, "StringLenMultiByte": {funcName: "StringLenMultiByte", native: divergence_hunt67.StringLenMultiByte}, "RuneCount": {funcName: "RuneCount", native: divergence_hunt67.RuneCount}, "StringIndexByte": {funcName: "StringIndexByte", native: divergence_hunt67.StringIndexByte}, "StringSlice": {funcName: "StringSlice", native: divergence_hunt67.StringSlice}, "StringConcatEmpty": {funcName: "StringConcatEmpty", native: divergence_hunt67.StringConcatEmpty}, "StringConcatMulti": {funcName: "StringConcatMulti", native: divergence_hunt67.StringConcatMulti}, "StringCompare": {funcName: "StringCompare", native: divergence_hunt67.StringCompare}, "StringEqual": {funcName: "StringEqual", native: divergence_hunt67.StringEqual}, "StringEmptyCompare": {funcName: "StringEmptyCompare", native: divergence_hunt67.StringEmptyCompare}, "RuneValue": {funcName: "RuneValue", native: divergence_hunt67.RuneValue}, "RuneChineseValue": {funcName: "RuneChineseValue", native: divergence_hunt67.RuneChineseValue}, "StringFromBytes": {funcName: "StringFromBytes", native: divergence_hunt67.StringFromBytes}, "StringToBytes": {funcName: "StringToBytes", native: divergence_hunt67.StringToBytes}, "StringFromRunes": {funcName: "StringFromRunes", native: divergence_hunt67.StringFromRunes}, "StringRangeIndex": {funcName: "StringRangeIndex", native: divergence_hunt67.StringRangeIndex}, "StringsRepeat": {funcName: "StringsRepeat", native: divergence_hunt67.StringsRepeat}, "StringsTrimCutset": {funcName: "StringsTrimCutset", native: divergence_hunt67.StringsTrimCutset}, "StringContainsEmpty": {funcName: "StringContainsEmpty", native: divergence_hunt67.StringContainsEmpty},
	}})
}
func TestDivergenceHunt68(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt68Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"PointerReceiverModify": {funcName: "PointerReceiverModify", native: divergence_hunt68.PointerReceiverModify}, "ValueReceiverNoModify": {funcName: "ValueReceiverNoModify", native: divergence_hunt68.ValueReceiverNoModify}, "MixedReceiver": {funcName: "MixedReceiver", native: divergence_hunt68.MixedReceiver}, "StructLiteral": {funcName: "StructLiteral", native: divergence_hunt68.StructLiteral}, "StructZeroValue": {funcName: "StructZeroValue", native: divergence_hunt68.StructZeroValue}, "StructPointerLiteral": {funcName: "StructPointerLiteral", native: divergence_hunt68.StructPointerLiteral}, "StructFieldAssign": {funcName: "StructFieldAssign", native: divergence_hunt68.StructFieldAssign}, "StructPointerFieldAssign": {funcName: "StructPointerFieldAssign", native: divergence_hunt68.StructPointerFieldAssign}, "StructNested": {funcName: "StructNested", native: divergence_hunt68.StructNested}, "StructNestedFieldAssign": {funcName: "StructNestedFieldAssign", native: divergence_hunt68.StructNestedFieldAssign}, "StructMethodChain": {funcName: "StructMethodChain", native: divergence_hunt68.StructMethodChain}, "StructCopySemantics": {funcName: "StructCopySemantics", native: divergence_hunt68.StructCopySemantics}, "StructPointerCopy": {funcName: "StructPointerCopy", native: divergence_hunt68.StructPointerCopy}, "MethodOnLiteral": {funcName: "MethodOnLiteral", native: divergence_hunt68.MethodOnLiteral},
	}})
}
func TestDivergenceHunt69(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt69Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ChannelBasic": {funcName: "ChannelBasic", native: divergence_hunt69.ChannelBasic}, "ChannelBuffered": {funcName: "ChannelBuffered", native: divergence_hunt69.ChannelBuffered}, "ChannelClose": {funcName: "ChannelClose", native: divergence_hunt69.ChannelClose}, "ChannelClosedReadZero": {funcName: "ChannelClosedReadZero", native: divergence_hunt69.ChannelClosedReadZero}, "ChannelLen": {funcName: "ChannelLen", native: divergence_hunt69.ChannelLen}, "ChannelCap": {funcName: "ChannelCap", native: divergence_hunt69.ChannelCap}, "SelectBasic": {funcName: "SelectBasic", native: divergence_hunt69.SelectBasic}, "SelectDefault": {funcName: "SelectDefault", native: divergence_hunt69.SelectDefault}, "ChannelCloseAndRange": {funcName: "ChannelCloseAndRange", native: divergence_hunt69.ChannelCloseAndRange}, "NilChannelBlocks": {funcName: "NilChannelBlocks", native: divergence_hunt69.NilChannelBlocks}, "ChannelDirection": {funcName: "ChannelDirection", native: divergence_hunt69.ChannelDirection}, "ChannelSelectMultiple": {funcName: "ChannelSelectMultiple", native: divergence_hunt69.ChannelSelectMultiple},
	}})
}
func TestDivergenceHunt70(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt70Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NamedTypeMethod": {funcName: "NamedTypeMethod", native: divergence_hunt70.NamedTypeMethod}, "NamedTypeConversion": {funcName: "NamedTypeConversion", native: divergence_hunt70.NamedTypeConversion}, "NamedTypeArith": {funcName: "NamedTypeArith", native: divergence_hunt70.NamedTypeArith}, "TypeAliasConversion": {funcName: "TypeAliasConversion", native: divergence_hunt70.TypeAliasConversion}, "NamedStringType": {funcName: "NamedStringType", native: divergence_hunt70.NamedStringType}, "NamedBoolType": {funcName: "NamedBoolType", native: divergence_hunt70.NamedBoolType}, "NamedSliceType": {funcName: "NamedSliceType", native: divergence_hunt70.NamedSliceType}, "NamedMapType": {funcName: "NamedMapType", native: divergence_hunt70.NamedMapType}, "NamedFuncType": {funcName: "NamedFuncType", native: divergence_hunt70.NamedFuncType}, "NamedPointerType": {funcName: "NamedPointerType", native: divergence_hunt70.NamedPointerType}, "NamedTypeCompare": {funcName: "NamedTypeCompare", native: divergence_hunt70.NamedTypeCompare}, "NamedTypeLessThan": {funcName: "NamedTypeLessThan", native: divergence_hunt70.NamedTypeLessThan},
	}})
}
func TestDivergenceHunt71(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt71Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InterfaceSatisfaction": {funcName: "InterfaceSatisfaction", native: divergence_hunt71.InterfaceSatisfaction}, "PointerReceiverInterface": {funcName: "PointerReceiverInterface", native: divergence_hunt71.PointerReceiverInterface}, "InterfaceNilCheck": {funcName: "InterfaceNilCheck", native: divergence_hunt71.InterfaceNilCheck}, "InterfaceTypeSwitch": {funcName: "InterfaceTypeSwitch", native: divergence_hunt71.InterfaceTypeSwitch}, "InterfaceAssertionOk": {funcName: "InterfaceAssertionOk", native: divergence_hunt71.InterfaceAssertionOk}, "InterfaceAssertionFail": {funcName: "InterfaceAssertionFail", native: divergence_hunt71.InterfaceAssertionFail}, "EmptyInterface": {funcName: "EmptyInterface", native: divergence_hunt71.EmptyInterface}, "InterfaceSlice": {funcName: "InterfaceSlice", native: divergence_hunt71.InterfaceSlice}, "InterfaceMap": {funcName: "InterfaceMap", native: divergence_hunt71.InterfaceMap}, "InterfaceMethodCall": {funcName: "InterfaceMethodCall", native: divergence_hunt71.InterfaceMethodCall}, "InterfaceAsField": {funcName: "InterfaceAsField", native: divergence_hunt71.InterfaceAsField},
	}})
}
func TestDivergenceHunt72(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt72Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Float64NaN": {funcName: "Float64NaN", native: divergence_hunt72.Float64NaN}, "Float64Inf": {funcName: "Float64Inf", native: divergence_hunt72.Float64Inf}, "Float64NegInf": {funcName: "Float64NegInf", native: divergence_hunt72.Float64NegInf}, "Float64Zero": {funcName: "Float64Zero", native: divergence_hunt72.Float64Zero}, "Float64NegZero": {funcName: "Float64NegZero", native: divergence_hunt72.Float64NegZero}, "Float64NaNNotEqual": {funcName: "Float64NaNNotEqual", native: divergence_hunt72.Float64NaNNotEqual}, "Float64InfArith": {funcName: "Float64InfArith", native: divergence_hunt72.Float64InfArith}, "Float64InfSubInf": {funcName: "Float64InfSubInf", native: divergence_hunt72.Float64InfSubInf}, "Float64ZeroDiv": {funcName: "Float64ZeroDiv", native: divergence_hunt72.Float64ZeroDiv}, "Float32Precision": {funcName: "Float32Precision", native: divergence_hunt72.Float32Precision}, "Float64Truncation": {funcName: "Float64Truncation", native: divergence_hunt72.Float64Truncation}, "Float64NegativeTruncation": {funcName: "Float64NegativeTruncation", native: divergence_hunt72.Float64NegativeTruncation}, "Float64Mod": {funcName: "Float64Mod", native: divergence_hunt72.Float64Mod}, "Float64Pow": {funcName: "Float64Pow", native: divergence_hunt72.Float64Pow}, "Float64Sqrt": {funcName: "Float64Sqrt", native: divergence_hunt72.Float64Sqrt}, "Float64Abs": {funcName: "Float64Abs", native: divergence_hunt72.Float64Abs}, "Float64Max": {funcName: "Float64Max", native: divergence_hunt72.Float64Max}, "Float64Min": {funcName: "Float64Min", native: divergence_hunt72.Float64Min},
	}})
}
func TestDivergenceHunt73(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt73Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SwitchNoExpression": {funcName: "SwitchNoExpression", native: divergence_hunt73.SwitchNoExpression}, "SwitchMultiCase": {funcName: "SwitchMultiCase", native: divergence_hunt73.SwitchMultiCase}, "SwitchFallthrough": {funcName: "SwitchFallthrough", native: divergence_hunt73.SwitchFallthrough}, "SwitchWithInit": {funcName: "SwitchWithInit", native: divergence_hunt73.SwitchWithInit}, "SwitchString": {funcName: "SwitchString", native: divergence_hunt73.SwitchString}, "SwitchEmpty": {funcName: "SwitchEmpty", native: divergence_hunt73.SwitchEmpty}, "SwitchOnlyDefault": {funcName: "SwitchOnlyDefault", native: divergence_hunt73.SwitchOnlyDefault}, "TypeSwitchWithDefault": {funcName: "TypeSwitchWithDefault", native: divergence_hunt73.TypeSwitchWithDefault}, "TypeSwitchNil": {funcName: "TypeSwitchNil", native: divergence_hunt73.TypeSwitchNil}, "SwitchBreak": {funcName: "SwitchBreak", native: divergence_hunt73.SwitchBreak}, "SwitchNested": {funcName: "SwitchNested", native: divergence_hunt73.SwitchNested}, "SwitchBool": {funcName: "SwitchBool", native: divergence_hunt73.SwitchBool}, "SwitchInterface": {funcName: "SwitchInterface", native: divergence_hunt73.SwitchInterface},
	}})
}
func TestDivergenceHunt74(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt74Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"VariadicDirect": {funcName: "VariadicDirect", native: divergence_hunt74.VariadicDirect}, "VariadicEmpty": {funcName: "VariadicEmpty", native: divergence_hunt74.VariadicEmpty}, "VariadicSpread": {funcName: "VariadicSpread", native: divergence_hunt74.VariadicSpread}, "VariadicWithPrefix": {funcName: "VariadicWithPrefix", native: divergence_hunt74.VariadicWithPrefix}, "VariadicString": {funcName: "VariadicString", native: divergence_hunt74.VariadicString}, "VariadicInterface": {funcName: "VariadicInterface", native: divergence_hunt74.VariadicInterface}, "VariadicNilSpread": {funcName: "VariadicNilSpread", native: divergence_hunt74.VariadicNilSpread}, "VariadicInClosure": {funcName: "VariadicInClosure", native: divergence_hunt74.VariadicInClosure}, "VariadicAppend": {funcName: "VariadicAppend", native: divergence_hunt74.VariadicAppend}, "VariadicAppendSpread": {funcName: "VariadicAppendSpread", native: divergence_hunt74.VariadicAppendSpread}, "VariadicFmt": {funcName: "VariadicFmt", native: divergence_hunt74.VariadicFmt}, "VariadicReturnSlice": {funcName: "VariadicReturnSlice", native: divergence_hunt74.VariadicReturnSlice},
	}})
}
func TestDivergenceHunt75(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt75Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"EmbeddingPromotion": {funcName: "EmbeddingPromotion", native: divergence_hunt75.EmbeddingPromotion}, "EmbeddingFieldAccess": {funcName: "EmbeddingFieldAccess", native: divergence_hunt75.EmbeddingFieldAccess}, "EmbeddingExplicitBase": {funcName: "EmbeddingExplicitBase", native: divergence_hunt75.EmbeddingExplicitBase}, "NestedEmbedding": {funcName: "NestedEmbedding", native: divergence_hunt75.NestedEmbedding}, "ShadowingEmbed": {funcName: "ShadowingEmbed", native: divergence_hunt75.ShadowingEmbed}, "ShadowingExplicit": {funcName: "ShadowingExplicit", native: divergence_hunt75.ShadowingExplicit}, "EmbeddingInterface": {funcName: "EmbeddingInterface", native: divergence_hunt75.EmbeddingInterface}, "EmbeddingLiteral": {funcName: "EmbeddingLiteral", native: divergence_hunt75.EmbeddingLiteral}, "EmbeddingFieldAssign": {funcName: "EmbeddingFieldAssign", native: divergence_hunt75.EmbeddingFieldAssign}, "DoubleEmbedding": {funcName: "DoubleEmbedding", native: divergence_hunt75.DoubleEmbedding}, "EmbeddingMethodPromotion": {funcName: "EmbeddingMethodPromotion", native: divergence_hunt75.EmbeddingMethodPromotion},
	}})
}
func TestDivergenceHunt76(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt76Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"PointerBasic": {funcName: "PointerBasic", native: divergence_hunt76.PointerBasic}, "PointerAssign": {funcName: "PointerAssign", native: divergence_hunt76.PointerAssign}, "PointerToStruct": {funcName: "PointerToStruct", native: divergence_hunt76.PointerToStruct}, "PointerNilCheck": {funcName: "PointerNilCheck", native: divergence_hunt76.PointerNilCheck}, "PointerReassign": {funcName: "PointerReassign", native: divergence_hunt76.PointerReassign}, "PointerAsArg": {funcName: "PointerAsArg", native: divergence_hunt76.PointerAsArg}, "PointerReturn": {funcName: "PointerReturn", native: divergence_hunt76.PointerReturn}, "PointerDerefAssign": {funcName: "PointerDerefAssign", native: divergence_hunt76.PointerDerefAssign}, "StructPointerMethod": {funcName: "StructPointerMethod", native: divergence_hunt76.StructPointerMethod}, "PointerToPointer": {funcName: "PointerToPointer", native: divergence_hunt76.PointerToPointer}, "PointerArray": {funcName: "PointerArray", native: divergence_hunt76.PointerArray}, "PointerSliceElem": {funcName: "PointerSliceElem", native: divergence_hunt76.PointerSliceElem}, "PointerSwap": {funcName: "PointerSwap", native: divergence_hunt76.PointerSwap}, "NewPointer": {funcName: "NewPointer", native: divergence_hunt76.NewPointer},
	}})
}
func TestDivergenceHunt77(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt77Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SliceOfStructLiteral": {funcName: "SliceOfStructLiteral", native: divergence_hunt77.SliceOfStructLiteral}, "MapOfStructLiteral": {funcName: "MapOfStructLiteral", native: divergence_hunt77.MapOfStructLiteral}, "NestedSliceLiteral": {funcName: "NestedSliceLiteral", native: divergence_hunt77.NestedSliceLiteral}, "ArrayLiteral": {funcName: "ArrayLiteral", native: divergence_hunt77.ArrayLiteral}, "ArrayAutoLen": {funcName: "ArrayAutoLen", native: divergence_hunt77.ArrayAutoLen}, "MapLiteralEmpty": {funcName: "MapLiteralEmpty", native: divergence_hunt77.MapLiteralEmpty}, "SliceLiteralEmpty": {funcName: "SliceLiteralEmpty", native: divergence_hunt77.SliceLiteralEmpty}, "StructLiteralPositional": {funcName: "StructLiteralPositional", native: divergence_hunt77.StructLiteralPositional}, "StructLiteralNamed": {funcName: "StructLiteralNamed", native: divergence_hunt77.StructLiteralNamed}, "SliceOfMap": {funcName: "SliceOfMap", native: divergence_hunt77.SliceOfMap}, "MapKeyStruct": {funcName: "MapKeyStruct", native: divergence_hunt77.MapKeyStruct}, "NestedMapLiteral": {funcName: "NestedMapLiteral", native: divergence_hunt77.NestedMapLiteral}, "PointerStructLiteral": {funcName: "PointerStructLiteral", native: divergence_hunt77.PointerStructLiteral}, "SliceOfPointer": {funcName: "SliceOfPointer", native: divergence_hunt77.SliceOfPointer},
	}})
}
func TestDivergenceHunt78(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt78Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"IntToInt8": {funcName: "IntToInt8", native: divergence_hunt78.IntToInt8}, "IntToUint": {funcName: "IntToUint", native: divergence_hunt78.IntToUint}, "UintToInt": {funcName: "UintToInt", native: divergence_hunt78.UintToInt}, "Float64ToInt": {funcName: "Float64ToInt", native: divergence_hunt78.Float64ToInt}, "Float64ToUint": {funcName: "Float64ToUint", native: divergence_hunt78.Float64ToUint}, "IntToFloat64": {funcName: "IntToFloat64", native: divergence_hunt78.IntToFloat64}, "Int8ToInt16": {funcName: "Int8ToInt16", native: divergence_hunt78.Int8ToInt16}, "Uint8ToUint16": {funcName: "Uint8ToUint16", native: divergence_hunt78.Uint8ToUint16}, "Int32ToInt8": {funcName: "Int32ToInt8", native: divergence_hunt78.Int32ToInt8}, "Float32ToFloat64": {funcName: "Float32ToFloat64", native: divergence_hunt78.Float32ToFloat64}, "Float64ToFloat32": {funcName: "Float64ToFloat32", native: divergence_hunt78.Float64ToFloat32}, "ByteToInt": {funcName: "ByteToInt", native: divergence_hunt78.ByteToInt}, "RuneToInt": {funcName: "RuneToInt", native: divergence_hunt78.RuneToInt}, "SliceToInterface": {funcName: "SliceToInterface", native: divergence_hunt78.SliceToInterface}, "StringToByteSlice": {funcName: "StringToByteSlice", native: divergence_hunt78.StringToByteSlice}, "ByteSliceToString": {funcName: "ByteSliceToString", native: divergence_hunt78.ByteSliceToString}, "IntSliceToFloatSlice": {funcName: "IntSliceToFloatSlice", native: divergence_hunt78.IntSliceToFloatSlice},
	}})
}
func TestDivergenceHunt79(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt79Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RangeSlice": {funcName: "RangeSlice", native: divergence_hunt79.RangeSlice}, "RangeSliceIndex": {funcName: "RangeSliceIndex", native: divergence_hunt79.RangeSliceIndex}, "RangeStringRunes": {funcName: "RangeStringRunes", native: divergence_hunt79.RangeStringRunes}, "RangeMapKeys": {funcName: "RangeMapKeys", native: divergence_hunt79.RangeMapKeys}, "RangeModifySlice": {funcName: "RangeModifySlice", native: divergence_hunt79.RangeModifySlice}, "RangeWithBreak": {funcName: "RangeWithBreak", native: divergence_hunt79.RangeWithBreak}, "RangeWithContinue": {funcName: "RangeWithContinue", native: divergence_hunt79.RangeWithContinue}, "RangeEmptySlice": {funcName: "RangeEmptySlice", native: divergence_hunt79.RangeEmptySlice}, "RangeNilSlice": {funcName: "RangeNilSlice", native: divergence_hunt79.RangeNilSlice}, "RangeNilMap": {funcName: "RangeNilMap", native: divergence_hunt79.RangeNilMap}, "RangeChannel": {funcName: "RangeChannel", native: divergence_hunt79.RangeChannel}, "RangeArray": {funcName: "RangeArray", native: divergence_hunt79.RangeArray}, "RangeMultiByteString": {funcName: "RangeMultiByteString", native: divergence_hunt79.RangeMultiByteString}, "RangeStringIndexRune": {funcName: "RangeStringIndexRune", native: divergence_hunt79.RangeStringIndexRune},
	}})
}
func TestDivergenceHunt80(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt80Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SwapInt": {funcName: "SwapInt", native: divergence_hunt80.SwapInt}, "MultiReturnAssign": {funcName: "MultiReturnAssign", native: divergence_hunt80.MultiReturnAssign}, "BlankAssign": {funcName: "BlankAssign", native: divergence_hunt80.BlankAssign}, "MultiAssignSameVar": {funcName: "MultiAssignSameVar", native: divergence_hunt80.MultiAssignSameVar}, "MultiAssignSwap": {funcName: "MultiAssignSwap", native: divergence_hunt80.MultiAssignSwap}, "AssignMapAccess": {funcName: "AssignMapAccess", native: divergence_hunt80.AssignMapAccess}, "AssignTypeAssertionString": {funcName: "AssignTypeAssertionString", native: divergence_hunt80.AssignTypeAssertionString}, "MultiReturnBlank": {funcName: "MultiReturnBlank", native: divergence_hunt80.MultiReturnBlank}, "SwapSliceElements": {funcName: "SwapSliceElements", native: divergence_hunt80.SwapSliceElements}, "AssignStructFields": {funcName: "AssignStructFields", native: divergence_hunt80.AssignStructFields}, "MultiAssignExpression": {funcName: "MultiAssignExpression", native: divergence_hunt80.MultiAssignExpression}, "AssignPointerDeref": {funcName: "AssignPointerDeref", native: divergence_hunt80.AssignPointerDeref}, "NestedMultiReturn": {funcName: "NestedMultiReturn", native: divergence_hunt80.NestedMultiReturn},
	}})
}
func TestDivergenceHunt81(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt81Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"AppendToNil": {funcName: "AppendToNil", native: divergence_hunt81.AppendToNil}, "AppendMultiple": {funcName: "AppendMultiple", native: divergence_hunt81.AppendMultiple}, "AppendSlice": {funcName: "AppendSlice", native: divergence_hunt81.AppendSlice}, "AppendEmptySlice": {funcName: "AppendEmptySlice", native: divergence_hunt81.AppendEmptySlice}, "AppendNilSlice": {funcName: "AppendNilSlice", native: divergence_hunt81.AppendNilSlice}, "CopyBasic": {funcName: "CopyBasic", native: divergence_hunt81.CopyBasic}, "CopyLargerDst": {funcName: "CopyLargerDst", native: divergence_hunt81.CopyLargerDst}, "CopySlice": {funcName: "CopySlice", native: divergence_hunt81.CopySlice}, "CopyPartial": {funcName: "CopyPartial", native: divergence_hunt81.CopyPartial}, "AppendBool": {funcName: "AppendBool", native: divergence_hunt81.AppendBool}, "AppendString": {funcName: "AppendString", native: divergence_hunt81.AppendString}, "AppendFloat": {funcName: "AppendFloat", native: divergence_hunt81.AppendFloat}, "CopyStringSlice": {funcName: "CopyStringSlice", native: divergence_hunt81.CopyStringSlice}, "AppendGrow": {funcName: "AppendGrow", native: divergence_hunt81.AppendGrow},
	}})
}
func TestDivergenceHunt82(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt82Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ErrorAsInterface": {funcName: "ErrorAsInterface", native: divergence_hunt82.ErrorAsInterface}, "ErrorNilCheck": {funcName: "ErrorNilCheck", native: divergence_hunt82.ErrorNilCheck}, "ErrorTypeAssertion": {funcName: "ErrorTypeAssertion", native: divergence_hunt82.ErrorTypeAssertion}, "ErrorPointerAssertion": {funcName: "ErrorPointerAssertion", native: divergence_hunt82.ErrorPointerAssertion}, "ErrorPointerDoesNotMatchValue": {funcName: "ErrorPointerDoesNotMatchValue", native: divergence_hunt82.ErrorPointerDoesNotMatchValue}, "ErrorValueDoesNotMatchPointer": {funcName: "ErrorValueDoesNotMatchPointer", native: divergence_hunt82.ErrorValueDoesNotMatchPointer}, "FmtErrorf": {funcName: "FmtErrorf", native: divergence_hunt82.FmtErrorf}, "ErrorInMultiReturn": {funcName: "ErrorInMultiReturn", native: divergence_hunt82.ErrorInMultiReturn}, "ErrorInMultiReturnFail": {funcName: "ErrorInMultiReturnFail", native: divergence_hunt82.ErrorInMultiReturnFail}, "ErrorSlice": {funcName: "ErrorSlice", native: divergence_hunt82.ErrorSlice},
	}})
}
func TestDivergenceHunt83(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt83Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"GlobalConstAccess": {funcName: "GlobalConstAccess", native: divergence_hunt83.GlobalConstAccess}, "GlobalIota": {funcName: "GlobalIota", native: divergence_hunt83.GlobalIota},
	}})
}
func TestDivergenceHunt84(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt84Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"LinkedListSum": {funcName: "LinkedListSum", native: divergence_hunt84.LinkedListSum}, "LinkedListLength": {funcName: "LinkedListLength", native: divergence_hunt84.LinkedListLength}, "TreeSum": {funcName: "TreeSum", native: divergence_hunt84.TreeSum}, "TreeDepth": {funcName: "TreeDepth", native: divergence_hunt84.TreeDepth}, "LinkedListCreate": {funcName: "LinkedListCreate", native: divergence_hunt84.LinkedListCreate}, "LinkedListMiddle": {funcName: "LinkedListMiddle", native: divergence_hunt84.LinkedListMiddle}, "TreeLeafCount": {funcName: "TreeLeafCount", native: divergence_hunt84.TreeLeafCount},
	}})
}
func TestDivergenceHunt85(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt85Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BuilderBasic": {funcName: "BuilderBasic", native: divergence_hunt85.BuilderBasic}, "BuilderLen": {funcName: "BuilderLen", native: divergence_hunt85.BuilderLen}, "BuilderGrow": {funcName: "BuilderGrow", native: divergence_hunt85.BuilderGrow}, "BuilderReset": {funcName: "BuilderReset", native: divergence_hunt85.BuilderReset}, "BuilderWriteByte": {funcName: "BuilderWriteByte", native: divergence_hunt85.BuilderWriteByte}, "BuilderWriteString": {funcName: "BuilderWriteString", native: divergence_hunt85.BuilderWriteString}, "BuilderCap": {funcName: "BuilderCap", native: divergence_hunt85.BuilderCap}, "BuilderEmpty": {funcName: "BuilderEmpty", native: divergence_hunt85.BuilderEmpty}, "BuilderLarge": {funcName: "BuilderLarge", native: divergence_hunt85.BuilderLarge}, "StringConcatMany": {funcName: "StringConcatMany", native: divergence_hunt85.StringConcatMany}, "StringJoin": {funcName: "StringJoin", native: divergence_hunt85.StringJoin}, "StringRepeat": {funcName: "StringRepeat", native: divergence_hunt85.StringRepeat}, "StringReplace": {funcName: "StringReplace", native: divergence_hunt85.StringReplace}, "StringReplaceAll": {funcName: "StringReplaceAll", native: divergence_hunt85.StringReplaceAll}, "StringContains": {funcName: "StringContains", native: divergence_hunt85.StringContains},
	}})
}
func TestDivergenceHunt86(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt86Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"JsonMarshalBasic": {funcName: "JsonMarshalBasic", native: divergence_hunt86.JsonMarshalBasic}, "JsonUnmarshalBasic": {funcName: "JsonUnmarshalBasic", native: divergence_hunt86.JsonUnmarshalBasic}, "JsonMarshalSlice": {funcName: "JsonMarshalSlice", native: divergence_hunt86.JsonMarshalSlice}, "JsonUnmarshalSlice": {funcName: "JsonUnmarshalSlice", native: divergence_hunt86.JsonUnmarshalSlice}, "JsonMarshalMap": {funcName: "JsonMarshalMap", native: divergence_hunt86.JsonMarshalMap}, "JsonMarshalNested": {funcName: "JsonMarshalNested", native: divergence_hunt86.JsonMarshalNested}, "JsonMarshalBool": {funcName: "JsonMarshalBool", native: divergence_hunt86.JsonMarshalBool}, "JsonUnmarshalBool": {funcName: "JsonUnmarshalBool", native: divergence_hunt86.JsonUnmarshalBool}, "JsonMarshalNull": {funcName: "JsonMarshalNull", native: divergence_hunt86.JsonMarshalNull}, "JsonMarshalString": {funcName: "JsonMarshalString", native: divergence_hunt86.JsonMarshalString}, "JsonUnmarshalString": {funcName: "JsonUnmarshalString", native: divergence_hunt86.JsonUnmarshalString}, "JsonRoundTrip": {funcName: "JsonRoundTrip", native: divergence_hunt86.JsonRoundTrip},
	}})
}
func TestDivergenceHunt87(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt87Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BitwiseAnd": {funcName: "BitwiseAnd", native: divergence_hunt87.BitwiseAnd}, "BitwiseOr": {funcName: "BitwiseOr", native: divergence_hunt87.BitwiseOr}, "BitwiseXor": {funcName: "BitwiseXor", native: divergence_hunt87.BitwiseXor}, "BitwiseNot": {funcName: "BitwiseNot", native: divergence_hunt87.BitwiseNot}, "BitwiseShiftLeft": {funcName: "BitwiseShiftLeft", native: divergence_hunt87.BitwiseShiftLeft}, "BitwiseShiftRight": {funcName: "BitwiseShiftRight", native: divergence_hunt87.BitwiseShiftRight}, "BitwiseShiftLeftOverflow": {funcName: "BitwiseShiftLeftOverflow", native: divergence_hunt87.BitwiseShiftLeftOverflow}, "BitwiseAndNot": {funcName: "BitwiseAndNot", native: divergence_hunt87.BitwiseAndNot}, "BitMask": {funcName: "BitMask", native: divergence_hunt87.BitMask}, "BitSet": {funcName: "BitSet", native: divergence_hunt87.BitSet}, "BitClear": {funcName: "BitClear", native: divergence_hunt87.BitClear}, "BitToggle": {funcName: "BitToggle", native: divergence_hunt87.BitToggle}, "BitCheck": {funcName: "BitCheck", native: divergence_hunt87.BitCheck}, "ShiftByVariable": {funcName: "ShiftByVariable", native: divergence_hunt87.ShiftByVariable}, "Uint8BitOps": {funcName: "Uint8BitOps", native: divergence_hunt87.Uint8BitOps}, "IntBitSign": {funcName: "IntBitSign", native: divergence_hunt87.IntBitSign}, "BitCount": {funcName: "BitCount", native: divergence_hunt87.BitCount}, "ReverseBits": {funcName: "ReverseBits", native: divergence_hunt87.ReverseBits},
	}})
}
func TestDivergenceHunt88(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt88Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MethodValue": {funcName: "MethodValue", native: divergence_hunt88.MethodValue}, "MethodValuePointer": {funcName: "MethodValuePointer", native: divergence_hunt88.MethodValuePointer}, "MethodCall": {funcName: "MethodCall", native: divergence_hunt88.MethodCall}, "MethodCallPointer": {funcName: "MethodCallPointer", native: divergence_hunt88.MethodCallPointer}, "MethodValueString": {funcName: "MethodValueString", native: divergence_hunt88.MethodValueString}, "MethodValueModify": {funcName: "MethodValueModify", native: divergence_hunt88.MethodValueModify}, "MethodOnStructLiteral": {funcName: "MethodOnStructLiteral", native: divergence_hunt88.MethodOnStructLiteral}, "MethodValueInLoop": {funcName: "MethodValueInLoop", native: divergence_hunt88.MethodValueInLoop}, "MethodValueReturn": {funcName: "MethodValueReturn", native: divergence_hunt88.MethodValueReturn}, "MethodValueChain": {funcName: "MethodValueChain", native: divergence_hunt88.MethodValueChain},
	}})
}
func TestDivergenceHunt89(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt89Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InterfaceComposition": {funcName: "InterfaceComposition", native: divergence_hunt89.InterfaceComposition}, "InterfaceEmbedding": {funcName: "InterfaceEmbedding", native: divergence_hunt89.InterfaceEmbedding}, "InterfaceAssertionComposition": {funcName: "InterfaceAssertionComposition", native: divergence_hunt89.InterfaceAssertionComposition}, "InterfaceSliceOfInterface": {funcName: "InterfaceSliceOfInterface", native: divergence_hunt89.InterfaceSliceOfInterface}, "InterfaceMapOfInterface": {funcName: "InterfaceMapOfInterface", native: divergence_hunt89.InterfaceMapOfInterface}, "InterfaceCustomStringer": {funcName: "InterfaceCustomStringer", native: divergence_hunt89.InterfaceCustomStringer}, "InterfaceSliceAsAny": {funcName: "InterfaceSliceAsAny", native: divergence_hunt89.InterfaceSliceAsAny}, "InterfaceMapAsAny": {funcName: "InterfaceMapAsAny", native: divergence_hunt89.InterfaceMapAsAny}, "InterfaceFuncAsAny": {funcName: "InterfaceFuncAsAny", native: divergence_hunt89.InterfaceFuncAsAny}, "NilInterfaceAssertion": {funcName: "NilInterfaceAssertion", native: divergence_hunt89.NilInterfaceAssertion}, "EmptyInterfaceTypeSwitch": {funcName: "EmptyInterfaceTypeSwitch", native: divergence_hunt89.EmptyInterfaceTypeSwitch},
	}})
}
func TestDivergenceHunt90(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt90Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Comprehensive1": {funcName: "Comprehensive1", native: divergence_hunt90.Comprehensive1}, "Comprehensive2": {funcName: "Comprehensive2", native: divergence_hunt90.Comprehensive2}, "Comprehensive3": {funcName: "Comprehensive3", native: divergence_hunt90.Comprehensive3}, "Comprehensive4": {funcName: "Comprehensive4", native: divergence_hunt90.Comprehensive4}, "Comprehensive5": {funcName: "Comprehensive5", native: divergence_hunt90.Comprehensive5}, "Comprehensive6": {funcName: "Comprehensive6", native: divergence_hunt90.Comprehensive6}, "Comprehensive7": {funcName: "Comprehensive7", native: divergence_hunt90.Comprehensive7}, "Comprehensive8": {funcName: "Comprehensive8", native: divergence_hunt90.Comprehensive8}, "Comprehensive9": {funcName: "Comprehensive9", native: divergence_hunt90.Comprehensive9}, "Comprehensive10": {funcName: "Comprehensive10", native: divergence_hunt90.Comprehensive10},
	}})
}
func TestDivergenceHunt91(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt91Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ValueReceiver": {funcName: "ValueReceiver", native: divergence_hunt91.ValueReceiver}, "PointerReceiver": {funcName: "PointerReceiver", native: divergence_hunt91.PointerReceiver}, "EmbeddedMethod": {funcName: "EmbeddedMethod", native: divergence_hunt91.EmbeddedMethod}, "EmbeddedReset": {funcName: "EmbeddedReset", native: divergence_hunt91.EmbeddedReset}, "MethodOnLiteral": {funcName: "MethodOnLiteral", native: divergence_hunt91.MethodOnLiteral}, "MethodChain": {funcName: "MethodChain", native: divergence_hunt91.MethodChain}, "ValueCopySemantics": {funcName: "ValueCopySemantics", native: divergence_hunt91.ValueCopySemantics}, "PointerSharedSemantics": {funcName: "PointerSharedSemantics", native: divergence_hunt91.PointerSharedSemantics}, "ReceiverOnStructLiteral": {funcName: "ReceiverOnStructLiteral", native: divergence_hunt91.ReceiverOnStructLiteral}, "EmbeddedPromoteMethod": {funcName: "EmbeddedPromoteMethod", native: divergence_hunt91.EmbeddedPromoteMethod},
	}})
}
func TestDivergenceHunt92(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt92Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BufferedSendRecv": {funcName: "BufferedSendRecv", native: divergence_hunt92.BufferedSendRecv}, "BufferedLenCap": {funcName: "BufferedLenCap", native: divergence_hunt92.BufferedLenCap}, "CloseChannel": {funcName: "CloseChannel", native: divergence_hunt92.CloseChannel}, "CloseAndRecv": {funcName: "CloseAndRecv", native: divergence_hunt92.CloseAndRecv}, "SelectBasic": {funcName: "SelectBasic", native: divergence_hunt92.SelectBasic}, "SelectDefault": {funcName: "SelectDefault", native: divergence_hunt92.SelectDefault}, "ChannelNilBlock": {funcName: "ChannelNilBlock", native: divergence_hunt92.ChannelNilBlock}, "NilChannelSelect": {funcName: "NilChannelSelect", native: divergence_hunt92.NilChannelSelect}, "BufferedStringChan": {funcName: "BufferedStringChan", native: divergence_hunt92.BufferedStringChan}, "ChannelOfStruct": {funcName: "ChannelOfStruct", native: divergence_hunt92.ChannelOfStruct}, "ChannelDirection": {funcName: "ChannelDirection", native: divergence_hunt92.ChannelDirection}, "SelectMultipleReady": {funcName: "SelectMultipleReady", native: divergence_hunt92.SelectMultipleReady}, "CloseRangeSum": {funcName: "CloseRangeSum", native: divergence_hunt92.CloseRangeSum},
	}})
}
func TestDivergenceHunt93(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt93Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InterfaceSlice": {funcName: "InterfaceSlice", native: divergence_hunt93.InterfaceSlice}, "InterfaceMap": {funcName: "InterfaceMap", native: divergence_hunt93.InterfaceMap}, "InterfaceParam": {funcName: "InterfaceParam", native: divergence_hunt93.InterfaceParam}, "InterfaceReturn": {funcName: "InterfaceReturn", native: divergence_hunt93.InterfaceReturn}, "InterfaceNil": {funcName: "InterfaceNil", native: divergence_hunt93.InterfaceNil}, "InterfaceTypedNil": {funcName: "InterfaceTypedNil", native: divergence_hunt93.InterfaceTypedNil}, "InterfaceSliceOfInterface": {funcName: "InterfaceSliceOfInterface", native: divergence_hunt93.InterfaceSliceOfInterface}, "InterfaceSliceTypeAssert": {funcName: "InterfaceSliceTypeAssert", native: divergence_hunt93.InterfaceSliceTypeAssert}, "DoubleInterfaceEmbedding": {funcName: "DoubleInterfaceEmbedding", native: divergence_hunt93.DoubleInterfaceEmbedding},
	}})
}
func TestDivergenceHunt94(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt94Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TypeSwitchBasic": {funcName: "TypeSwitchBasic", native: divergence_hunt94.TypeSwitchBasic}, "TypeSwitchMultiple": {funcName: "TypeSwitchMultiple", native: divergence_hunt94.TypeSwitchMultiple}, "TypeAssertionCommaOk": {funcName: "TypeAssertionCommaOk", native: divergence_hunt94.TypeAssertionCommaOk}, "TypeAssertionPanicSafe": {funcName: "TypeAssertionPanicSafe", native: divergence_hunt94.TypeAssertionPanicSafe}, "NestedTypeSwitch": {funcName: "NestedTypeSwitch", native: divergence_hunt94.NestedTypeSwitch}, "TypeSwitchWithNil": {funcName: "TypeSwitchWithNil", native: divergence_hunt94.TypeSwitchWithNil}, "AssertToInterface": {funcName: "AssertToInterface", native: divergence_hunt94.AssertToInterface}, "AssertSliceTypes": {funcName: "AssertSliceTypes", native: divergence_hunt94.AssertSliceTypes}, "AssertMapType": {funcName: "AssertMapType", native: divergence_hunt94.AssertMapType}, "TypeSwitchFallthrough": {funcName: "TypeSwitchFallthrough", native: divergence_hunt94.TypeSwitchFallthrough},
	}})
}
func TestDivergenceHunt95(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt95Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"DeferArgEval": {funcName: "DeferArgEval", native: divergence_hunt95.DeferArgEval}, "DeferArgCapture": {funcName: "DeferArgCapture", native: divergence_hunt95.DeferArgCapture}, "DeferModifyReturn": {funcName: "DeferModifyReturn", native: divergence_hunt95.DeferModifyReturn}, "StackedDefers": {funcName: "StackedDefers", native: divergence_hunt95.StackedDefers}, "DeferInLoop": {funcName: "DeferInLoop", native: divergence_hunt95.DeferInLoop}, "DeferClosureCapture": {funcName: "DeferClosureCapture", native: divergence_hunt95.DeferClosureCapture}, "DeferWithRecover": {funcName: "DeferWithRecover", native: divergence_hunt95.DeferWithRecover}, "MultipleDefersOrder": {funcName: "MultipleDefersOrder", native: divergence_hunt95.MultipleDefersOrder}, "DeferReturnOrder": {funcName: "DeferReturnOrder", native: divergence_hunt95.DeferReturnOrder}, "DeferWithMethod": {funcName: "DeferWithMethod", native: divergence_hunt95.DeferWithMethod}, "DeferClosureArgVsCapture": {funcName: "DeferClosureArgVsCapture", native: divergence_hunt95.DeferClosureArgVsCapture},
	}})
}
func TestDivergenceHunt96(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt96Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SliceDelete": {funcName: "SliceDelete", native: divergence_hunt96.SliceDelete}, "SliceInsert": {funcName: "SliceInsert", native: divergence_hunt96.SliceInsert}, "SliceFilter": {funcName: "SliceFilter", native: divergence_hunt96.SliceFilter}, "SliceReverse": {funcName: "SliceReverse", native: divergence_hunt96.SliceReverse}, "SliceUnique": {funcName: "SliceUnique", native: divergence_hunt96.SliceUnique}, "SliceFlatten": {funcName: "SliceFlatten", native: divergence_hunt96.SliceFlatten}, "SliceBatch": {funcName: "SliceBatch", native: divergence_hunt96.SliceBatch}, "SliceClone": {funcName: "SliceClone", native: divergence_hunt96.SliceClone}, "SliceAppendGrow": {funcName: "SliceAppendGrow", native: divergence_hunt96.SliceAppendGrow}, "SliceCut": {funcName: "SliceCut", native: divergence_hunt96.SliceCut},
	}})
}
func TestDivergenceHunt97(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt97Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapDeleteKey": {funcName: "MapDeleteKey", native: divergence_hunt97.MapDeleteKey}, "MapDeleteNonExistent": {funcName: "MapDeleteNonExistent", native: divergence_hunt97.MapDeleteNonExistent}, "MapDoubleDelete": {funcName: "MapDoubleDelete", native: divergence_hunt97.MapDoubleDelete}, "MapClear": {funcName: "MapClear", native: divergence_hunt97.MapClear}, "MapAccessMissing": {funcName: "MapAccessMissing", native: divergence_hunt97.MapAccessMissing}, "MapSetDefault": {funcName: "MapSetDefault", native: divergence_hunt97.MapSetDefault}, "MapCountValues": {funcName: "MapCountValues", native: divergence_hunt97.MapCountValues}, "MapInvert": {funcName: "MapInvert", native: divergence_hunt97.MapInvert}, "MapMerge": {funcName: "MapMerge", native: divergence_hunt97.MapMerge}, "MapKeys": {funcName: "MapKeys", native: divergence_hunt97.MapKeys}, "MapNestedAccess": {funcName: "MapNestedAccess", native: divergence_hunt97.MapNestedAccess},
	}})
}
func TestDivergenceHunt98(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt98Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RuneCount": {funcName: "RuneCount", native: divergence_hunt98.RuneCount}, "ByteLen": {funcName: "ByteLen", native: divergence_hunt98.ByteLen}, "RuneAt": {funcName: "RuneAt", native: divergence_hunt98.RuneAt}, "StringFromRunes": {funcName: "StringFromRunes", native: divergence_hunt98.StringFromRunes}, "StringSliceByte": {funcName: "StringSliceByte", native: divergence_hunt98.StringSliceByte}, "StringConcat": {funcName: "StringConcat", native: divergence_hunt98.StringConcat}, "StringRangeRunes": {funcName: "StringRangeRunes", native: divergence_hunt98.StringRangeRunes}, "RuneSliceModify": {funcName: "RuneSliceModify", native: divergence_hunt98.RuneSliceModify}, "MultiByteIndex": {funcName: "MultiByteIndex", native: divergence_hunt98.MultiByteIndex}, "StringCompare": {funcName: "StringCompare", native: divergence_hunt98.StringCompare}, "StringPrefixSuffix": {funcName: "StringPrefixSuffix", native: divergence_hunt98.StringPrefixSuffix}, "EmptyString": {funcName: "EmptyString", native: divergence_hunt98.EmptyString},
	}})
}
func TestDivergenceHunt99(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt99Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"OverrideMethod": {funcName: "OverrideMethod", native: divergence_hunt99.OverrideMethod}, "PromotedMethod": {funcName: "PromotedMethod", native: divergence_hunt99.PromotedMethod}, "DirectBaseMethod": {funcName: "DirectBaseMethod", native: divergence_hunt99.DirectBaseMethod}, "DeepEmbedding": {funcName: "DeepEmbedding", native: divergence_hunt99.DeepEmbedding}, "DeepSetViaBase": {funcName: "DeepSetViaBase", native: divergence_hunt99.DeepSetViaBase}, "TripleEmbedding": {funcName: "TripleEmbedding", native: divergence_hunt99.TripleEmbedding}, "EmbeddedLiteral": {funcName: "EmbeddedLiteral", native: divergence_hunt99.EmbeddedLiteral}, "OverrideVsPromote": {funcName: "OverrideVsPromote", native: divergence_hunt99.OverrideVsPromote},
	}})
}
func TestDivergenceHunt100(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt100Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicPanicRecover": {funcName: "BasicPanicRecover", native: divergence_hunt100.BasicPanicRecover}, "PanicInt": {funcName: "PanicInt", native: divergence_hunt100.PanicInt}, "PanicStruct": {funcName: "PanicStruct", native: divergence_hunt100.PanicStruct}, "NestedPanicRecover": {funcName: "NestedPanicRecover", native: divergence_hunt100.NestedPanicRecover}, "PanicInDefer": {funcName: "PanicInDefer", native: divergence_hunt100.PanicInDefer}, "NoPanicReturn": {funcName: "NoPanicReturn", native: divergence_hunt100.NoPanicReturn}, "RecoverWithoutPanic": {funcName: "RecoverWithoutPanic", native: divergence_hunt100.RecoverWithoutPanic}, "PanicNilInterface": {funcName: "PanicNilInterface", native: divergence_hunt100.PanicNilInterface}, "PanicSliceBounds": {funcName: "PanicSliceBounds", native: divergence_hunt100.PanicSliceBounds}, "PanicNilMap": {funcName: "PanicNilMap", native: divergence_hunt100.PanicNilMap}, "PanicNilPointer": {funcName: "PanicNilPointer", native: divergence_hunt100.PanicNilPointer},
	}})
}
func TestDivergenceHunt101(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt101Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"VariadicSumDirect": {funcName: "VariadicSumDirect", native: divergence_hunt101.VariadicSumDirect}, "VariadicConcatDirect": {funcName: "VariadicConcatDirect", native: divergence_hunt101.VariadicConcatDirect}, "VariadicEmpty": {funcName: "VariadicEmpty", native: divergence_hunt101.VariadicEmpty}, "VariadicFromSlice": {funcName: "VariadicFromSlice", native: divergence_hunt101.VariadicFromSlice}, "VariadicInterface": {funcName: "VariadicInterface", native: divergence_hunt101.VariadicInterface}, "VariadicNil": {funcName: "VariadicNil", native: divergence_hunt101.VariadicNil}, "VariadicStrings": {funcName: "VariadicStrings", native: divergence_hunt101.VariadicStrings}, "VariadicIntfType": {funcName: "VariadicIntfType", native: divergence_hunt101.VariadicIntfType}, "VariadicAppend": {funcName: "VariadicAppend", native: divergence_hunt101.VariadicAppend}, "VariadicSpread": {funcName: "VariadicSpread", native: divergence_hunt101.VariadicSpread},
	}})
}
func TestDivergenceHunt102(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt102Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TreeBuildAndSum": {funcName: "TreeBuildAndSum", native: divergence_hunt102.TreeBuildAndSum}, "TreeBuildAndDepth": {funcName: "TreeBuildAndDepth", native: divergence_hunt102.TreeBuildAndDepth}, "TreeBuildAndLeaves": {funcName: "TreeBuildAndLeaves", native: divergence_hunt102.TreeBuildAndLeaves}, "TreeInorderResult": {funcName: "TreeInorderResult", native: divergence_hunt102.TreeInorderResult}, "FibonacciTree": {funcName: "FibonacciTree", native: divergence_hunt102.FibonacciTree},
	}})
}
func TestDivergenceHunt103(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt103Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MultiReturnBlank": {funcName: "MultiReturnBlank", native: divergence_hunt103.MultiReturnBlank}, "MultiReturnAll": {funcName: "MultiReturnAll", native: divergence_hunt103.MultiReturnAll}, "ErrorReturn": {funcName: "ErrorReturn", native: divergence_hunt103.ErrorReturn}, "ErrorReturnFail": {funcName: "ErrorReturnFail", native: divergence_hunt103.ErrorReturnFail}, "NamedReturnBare": {funcName: "NamedReturnBare", native: divergence_hunt103.NamedReturnBare}, "NamedReturnOverride": {funcName: "NamedReturnOverride", native: divergence_hunt103.NamedReturnOverride}, "SwapValues": {funcName: "SwapValues", native: divergence_hunt103.SwapValues}, "MultiAssignExpression": {funcName: "MultiAssignExpression", native: divergence_hunt103.MultiAssignExpression}, "BlankInLoop": {funcName: "BlankInLoop", native: divergence_hunt103.BlankInLoop},
	}})
}
func TestDivergenceHunt104(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt104Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NamedBareReturn": {funcName: "NamedBareReturn", native: divergence_hunt104.NamedBareReturn}, "NamedBareReturnModify": {funcName: "NamedBareReturnModify", native: divergence_hunt104.NamedBareReturnModify}, "NamedBareReturnConditional": {funcName: "NamedBareReturnConditional", native: divergence_hunt104.NamedBareReturnConditional}, "NamedMultiBareReturn": {funcName: "NamedMultiBareReturn", native: divergence_hunt104.NamedMultiBareReturn}, "NamedReturnWithDefer": {funcName: "NamedReturnWithDefer", native: divergence_hunt104.NamedReturnWithDefer}, "NamedReturnDeferChain": {funcName: "NamedReturnDeferChain", native: divergence_hunt104.NamedReturnDeferChain}, "NamedReturnZeroValue": {funcName: "NamedReturnZeroValue", native: divergence_hunt104.NamedReturnZeroValue}, "NamedReturnPartial": {funcName: "NamedReturnPartial", native: divergence_hunt104.NamedReturnPartial}, "NamedReturnLoop": {funcName: "NamedReturnLoop", native: divergence_hunt104.NamedReturnLoop}, "NamedReturnClosure": {funcName: "NamedReturnClosure", native: divergence_hunt104.NamedReturnClosure},
	}})
}
func TestDivergenceHunt105(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt105Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NestedSliceLiteral": {funcName: "NestedSliceLiteral", native: divergence_hunt105.NestedSliceLiteral}, "MapLiteralWithStruct": {funcName: "MapLiteralWithStruct", native: divergence_hunt105.MapLiteralWithStruct}, "SliceOfMap": {funcName: "SliceOfMap", native: divergence_hunt105.SliceOfMap}, "StructWithSlice": {funcName: "StructWithSlice", native: divergence_hunt105.StructWithSlice}, "NestedMapLiteral": {funcName: "NestedMapLiteral", native: divergence_hunt105.NestedMapLiteral}, "SliceOfFunc": {funcName: "SliceOfFunc", native: divergence_hunt105.SliceOfFunc}, "EmptyCompositeLiterals": {funcName: "EmptyCompositeLiterals", native: divergence_hunt105.EmptyCompositeLiterals}, "PointerStructLiteral": {funcName: "PointerStructLiteral", native: divergence_hunt105.PointerStructLiteral}, "NestedStructLiteral": {funcName: "NestedStructLiteral", native: divergence_hunt105.NestedStructLiteral}, "ArrayLiteral": {funcName: "ArrayLiteral", native: divergence_hunt105.ArrayLiteral},
	}})
}
func TestDivergenceHunt106(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt106Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SwitchBasic": {funcName: "SwitchBasic", native: divergence_hunt106.SwitchBasic}, "SwitchDefault": {funcName: "SwitchDefault", native: divergence_hunt106.SwitchDefault}, "SwitchMultipleValues": {funcName: "SwitchMultipleValues", native: divergence_hunt106.SwitchMultipleValues}, "SwitchNoExpression": {funcName: "SwitchNoExpression", native: divergence_hunt106.SwitchNoExpression}, "SwitchFallthrough": {funcName: "SwitchFallthrough", native: divergence_hunt106.SwitchFallthrough}, "SwitchInLoop": {funcName: "SwitchInLoop", native: divergence_hunt106.SwitchInLoop}, "SwitchBreak": {funcName: "SwitchBreak", native: divergence_hunt106.SwitchBreak}, "SwitchString": {funcName: "SwitchString", native: divergence_hunt106.SwitchString}, "SwitchWithInit": {funcName: "SwitchWithInit", native: divergence_hunt106.SwitchWithInit}, "NestedSwitch": {funcName: "NestedSwitch", native: divergence_hunt106.NestedSwitch},
	}})
}
func TestDivergenceHunt107(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt107Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"PointerBasic": {funcName: "PointerBasic", native: divergence_hunt107.PointerBasic}, "PointerModify": {funcName: "PointerModify", native: divergence_hunt107.PointerModify}, "PointerToStruct": {funcName: "PointerToStruct", native: divergence_hunt107.PointerToStruct}, "PointerSwap": {funcName: "PointerSwap", native: divergence_hunt107.PointerSwap}, "NewKeyword": {funcName: "NewKeyword", native: divergence_hunt107.NewKeyword}, "NewStruct": {funcName: "NewStruct", native: divergence_hunt107.NewStruct}, "PointerSlice": {funcName: "PointerSlice", native: divergence_hunt107.PointerSlice}, "NilPointerCheck": {funcName: "NilPointerCheck", native: divergence_hunt107.NilPointerCheck}, "PointerAsParam": {funcName: "PointerAsParam", native: divergence_hunt107.PointerAsParam}, "PointerReturn": {funcName: "PointerReturn", native: divergence_hunt107.PointerReturn},
	}})
}
func TestDivergenceHunt108(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt108Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ClosureCounter": {funcName: "ClosureCounter", native: divergence_hunt108.ClosureCounter}, "ClosureCapture": {funcName: "ClosureCapture", native: divergence_hunt108.ClosureCapture}, "ClosureMultiCapture": {funcName: "ClosureMultiCapture", native: divergence_hunt108.ClosureMultiCapture}, "ClosureInLoop": {funcName: "ClosureInLoop", native: divergence_hunt108.ClosureInLoop}, "ClosureModifyOuter": {funcName: "ClosureModifyOuter", native: divergence_hunt108.ClosureModifyOuter}, "ClosureReturnClosure": {funcName: "ClosureReturnClosure", native: divergence_hunt108.ClosureReturnClosure}, "ClosureSlice": {funcName: "ClosureSlice", native: divergence_hunt108.ClosureSlice}, "ClosureAsParam": {funcName: "ClosureAsParam", native: divergence_hunt108.ClosureAsParam}, "ClosureCaptureSlice": {funcName: "ClosureCaptureSlice", native: divergence_hunt108.ClosureCaptureSlice}, "ClosureStacked": {funcName: "ClosureStacked", native: divergence_hunt108.ClosureStacked},
	}})
}
func TestDivergenceHunt109(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt109Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SortIntSlice": {funcName: "SortIntSlice", native: divergence_hunt109.SortIntSlice}, "SortStringSlice": {funcName: "SortStringSlice", native: divergence_hunt109.SortStringSlice}, "SortByLen": {funcName: "SortByLen", native: divergence_hunt109.SortByLen}, "SortStructByField": {funcName: "SortStructByField", native: divergence_hunt109.SortStructByField}, "SortReverse": {funcName: "SortReverse", native: divergence_hunt109.SortReverse}, "SortFloatSlice": {funcName: "SortFloatSlice", native: divergence_hunt109.SortFloatSlice}, "SortStable": {funcName: "SortStable", native: divergence_hunt109.SortStable}, "SortIsSorted": {funcName: "SortIsSorted", native: divergence_hunt109.SortIsSorted}, "SortEmpty": {funcName: "SortEmpty", native: divergence_hunt109.SortEmpty}, "SortSingleElement": {funcName: "SortSingleElement", native: divergence_hunt109.SortSingleElement},
	}})
}
func TestDivergenceHunt110(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt110Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ErrorBasic": {funcName: "ErrorBasic", native: divergence_hunt110.ErrorBasic}, "ErrorFmtErrorf": {funcName: "ErrorFmtErrorf", native: divergence_hunt110.ErrorFmtErrorf}, "ErrorWrapUnwrap": {funcName: "ErrorWrapUnwrap", native: divergence_hunt110.ErrorWrapUnwrap}, "ErrorIs": {funcName: "ErrorIs", native: divergence_hunt110.ErrorIs}, "ErrorAs": {funcName: "ErrorAs", native: divergence_hunt110.ErrorAs}, "ErrorChainIs": {funcName: "ErrorChainIs", native: divergence_hunt110.ErrorChainIs}, "ErrorNilIs": {funcName: "ErrorNilIs", native: divergence_hunt110.ErrorNilIs}, "ErrorTypeAssertion": {funcName: "ErrorTypeAssertion", native: divergence_hunt110.ErrorTypeAssertion}, "ErrorMultiWrap": {funcName: "ErrorMultiWrap", native: divergence_hunt110.ErrorMultiWrap}, "ErrorUnwrapNil": {funcName: "ErrorUnwrapNil", native: divergence_hunt110.ErrorUnwrapNil},
	}})
}
func TestDivergenceHunt111(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt111Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TimeFormat": {funcName: "TimeFormat", native: divergence_hunt111.TimeFormat}, "TimeParse": {funcName: "TimeParse", native: divergence_hunt111.TimeParse}, "TimeNow": {funcName: "TimeNow", native: divergence_hunt111.TimeNow}, "TimeAdd": {funcName: "TimeAdd", native: divergence_hunt111.TimeAdd}, "TimeSub": {funcName: "TimeSub", native: divergence_hunt111.TimeSub}, "TimeUnix": {funcName: "TimeUnix", native: divergence_hunt111.TimeUnix}, "TimeWeekday": {funcName: "TimeWeekday", native: divergence_hunt111.TimeWeekday}, "TimeBefore": {funcName: "TimeBefore", native: divergence_hunt111.TimeBefore}, "TimeFormatCustom": {funcName: "TimeFormatCustom", native: divergence_hunt111.TimeFormatCustom}, "TimeDateComponents": {funcName: "TimeDateComponents", native: divergence_hunt111.TimeDateComponents},
	}})
}
func TestDivergenceHunt112(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt112Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RegexpMatch": {funcName: "RegexpMatch", native: divergence_hunt112.RegexpMatch}, "RegexpMatchFail": {funcName: "RegexpMatchFail", native: divergence_hunt112.RegexpMatchFail}, "RegexpFindString": {funcName: "RegexpFindString", native: divergence_hunt112.RegexpFindString}, "RegexpFindAllString": {funcName: "RegexpFindAllString", native: divergence_hunt112.RegexpFindAllString}, "RegexpReplaceAllString": {funcName: "RegexpReplaceAllString", native: divergence_hunt112.RegexpReplaceAllString}, "RegexpSplit": {funcName: "RegexpSplit", native: divergence_hunt112.RegexpSplit}, "RegexpSubmatch": {funcName: "RegexpSubmatch", native: divergence_hunt112.RegexpSubmatch}, "RegexpReplaceAllStringFunc": {funcName: "RegexpReplaceAllStringFunc", native: divergence_hunt112.RegexpReplaceAllStringFunc}, "RegexpFindStringIndex": {funcName: "RegexpFindStringIndex", native: divergence_hunt112.RegexpFindStringIndex}, "RegexpCompileMust": {funcName: "RegexpCompileMust", native: divergence_hunt112.RegexpCompileMust},
	}})
}
func TestDivergenceHunt113(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt113Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"HexEncode": {funcName: "HexEncode", native: divergence_hunt113.HexEncode}, "HexDecode": {funcName: "HexDecode", native: divergence_hunt113.HexDecode}, "HexRoundtrip": {funcName: "HexRoundtrip", native: divergence_hunt113.HexRoundtrip}, "Base64Encode": {funcName: "Base64Encode", native: divergence_hunt113.Base64Encode}, "Base64Decode": {funcName: "Base64Decode", native: divergence_hunt113.Base64Decode}, "Base64Roundtrip": {funcName: "Base64Roundtrip", native: divergence_hunt113.Base64Roundtrip}, "Base64URLEncoding": {funcName: "Base64URLEncoding", native: divergence_hunt113.Base64URLEncoding}, "HexEmpty": {funcName: "HexEmpty", native: divergence_hunt113.HexEmpty}, "Base64Empty": {funcName: "Base64Empty", native: divergence_hunt113.Base64Empty}, "HexEncodeNumbers": {funcName: "HexEncodeNumbers", native: divergence_hunt113.HexEncodeNumbers},
	}})
}
func TestDivergenceHunt114(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt114Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MathAbs": {funcName: "MathAbs", native: divergence_hunt114.MathAbs}, "MathMax": {funcName: "MathMax", native: divergence_hunt114.MathMax}, "MathMin": {funcName: "MathMin", native: divergence_hunt114.MathMin}, "MathCeil": {funcName: "MathCeil", native: divergence_hunt114.MathCeil}, "MathFloor": {funcName: "MathFloor", native: divergence_hunt114.MathFloor}, "MathRound": {funcName: "MathRound", native: divergence_hunt114.MathRound}, "MathPow": {funcName: "MathPow", native: divergence_hunt114.MathPow}, "MathSqrt": {funcName: "MathSqrt", native: divergence_hunt114.MathSqrt}, "IntOverflow": {funcName: "IntOverflow", native: divergence_hunt114.IntOverflow}, "FloatPrecision": {funcName: "FloatPrecision", native: divergence_hunt114.FloatPrecision}, "IntegerDivision": {funcName: "IntegerDivision", native: divergence_hunt114.IntegerDivision}, "UintRange": {funcName: "UintRange", native: divergence_hunt114.UintRange}, "NegativeModulo": {funcName: "NegativeModulo", native: divergence_hunt114.NegativeModulo}, "FloatToIntTruncation": {funcName: "FloatToIntTruncation", native: divergence_hunt114.FloatToIntTruncation},
	}})
}
func TestDivergenceHunt115(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt115Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"GlobalRead": {funcName: "GlobalRead", native: divergence_hunt115.GlobalRead, knownIssue: true}, "GlobalModify": {funcName: "GlobalModify", native: divergence_hunt115.GlobalModify}, "GlobalSliceRead": {funcName: "GlobalSliceRead", native: divergence_hunt115.GlobalSliceRead}, "GlobalSliceLen": {funcName: "GlobalSliceLen", native: divergence_hunt115.GlobalSliceLen}, "GlobalMapRead": {funcName: "GlobalMapRead", native: divergence_hunt115.GlobalMapRead}, "GlobalMapLen": {funcName: "GlobalMapLen", native: divergence_hunt115.GlobalMapLen}, "GlobalStringRead": {funcName: "GlobalStringRead", native: divergence_hunt115.GlobalStringRead}, "GlobalInitValues": {funcName: "GlobalInitValues", native: divergence_hunt115.GlobalInitValues, knownIssue: true},
	}})
}
func TestDivergenceHunt116(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt116Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InterfaceNilCompare": {funcName: "InterfaceNilCompare", native: divergence_hunt116.InterfaceNilCompare}, "TypedNilInterface": {funcName: "TypedNilInterface", native: divergence_hunt116.TypedNilInterface}, "NilInterfaceTypeAssert": {funcName: "NilInterfaceTypeAssert", native: divergence_hunt116.NilInterfaceTypeAssert}, "NilInterfaceTypeSwitch": {funcName: "NilInterfaceTypeSwitch", native: divergence_hunt116.NilInterfaceTypeSwitch}, "EmptyInterfaceVsNil": {funcName: "EmptyInterfaceVsNil", native: divergence_hunt116.EmptyInterfaceVsNil}, "NilSliceVsNilInterface": {funcName: "NilSliceVsNilInterface", native: divergence_hunt116.NilSliceVsNilInterface}, "NilMapVsNilInterface": {funcName: "NilMapVsNilInterface", native: divergence_hunt116.NilMapVsNilInterface}, "NilFuncVsNilInterface": {funcName: "NilFuncVsNilInterface", native: divergence_hunt116.NilFuncVsNilInterface}, "NilChanVsNilInterface": {funcName: "NilChanVsNilInterface", native: divergence_hunt116.NilChanVsNilInterface}, "InterfaceReturnNil": {funcName: "InterfaceReturnNil", native: divergence_hunt116.InterfaceReturnNil},
	}})
}
func TestDivergenceHunt117(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt117Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StructEqual": {funcName: "StructEqual", native: divergence_hunt117.StructEqual}, "StructNotEqual": {funcName: "StructNotEqual", native: divergence_hunt117.StructNotEqual}, "StructCopy": {funcName: "StructCopy", native: divergence_hunt117.StructCopy}, "StructPointerEqual": {funcName: "StructPointerEqual", native: divergence_hunt117.StructPointerEqual}, "StructPointerSame": {funcName: "StructPointerSame", native: divergence_hunt117.StructPointerSame}, "StructWithSlice": {funcName: "StructWithSlice", native: divergence_hunt117.StructWithSlice}, "StructWithMap": {funcName: "StructWithMap", native: divergence_hunt117.StructWithMap}, "StructNested": {funcName: "StructNested", native: divergence_hunt117.StructNested}, "StructZeroValue": {funcName: "StructZeroValue", native: divergence_hunt117.StructZeroValue}, "StructSliceOfPointers": {funcName: "StructSliceOfPointers", native: divergence_hunt117.StructSliceOfPointers},
	}})
}
func TestDivergenceHunt118(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt118Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"FuncVariable": {funcName: "FuncVariable", native: divergence_hunt118.FuncVariable}, "FuncParam": {funcName: "FuncParam", native: divergence_hunt118.FuncParam}, "FuncReturn": {funcName: "FuncReturn", native: divergence_hunt118.FuncReturn}, "FuncSlice": {funcName: "FuncSlice", native: divergence_hunt118.FuncSlice}, "FuncMap": {funcName: "FuncMap", native: divergence_hunt118.FuncMap}, "FuncChaining": {funcName: "FuncChaining", native: divergence_hunt118.FuncChaining}, "FuncCompose": {funcName: "FuncCompose", native: divergence_hunt118.FuncCompose}, "FuncAsField": {funcName: "FuncAsField", native: divergence_hunt118.FuncAsField}, "FuncComparison": {funcName: "FuncComparison", native: divergence_hunt118.FuncComparison}, "FuncNilCheck": {funcName: "FuncNilCheck", native: divergence_hunt118.FuncNilCheck},
	}})
}
func TestDivergenceHunt119(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt119Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ArrayBasic": {funcName: "ArrayBasic", native: divergence_hunt119.ArrayBasic}, "ArrayLiteral": {funcName: "ArrayLiteral", native: divergence_hunt119.ArrayLiteral}, "ArrayAutoLen": {funcName: "ArrayAutoLen", native: divergence_hunt119.ArrayAutoLen}, "ArrayCopy": {funcName: "ArrayCopy", native: divergence_hunt119.ArrayCopy}, "ArrayRange": {funcName: "ArrayRange", native: divergence_hunt119.ArrayRange}, "ArrayPointer": {funcName: "ArrayPointer", native: divergence_hunt119.ArrayPointer}, "ArrayCompare": {funcName: "ArrayCompare", native: divergence_hunt119.ArrayCompare}, "ArrayNotEqual": {funcName: "ArrayNotEqual", native: divergence_hunt119.ArrayNotEqual}, "ArrayZeroValue": {funcName: "ArrayZeroValue", native: divergence_hunt119.ArrayZeroValue}, "ArrayOfStruct": {funcName: "ArrayOfStruct", native: divergence_hunt119.ArrayOfStruct}, "ArrayMultiDim": {funcName: "ArrayMultiDim", native: divergence_hunt119.ArrayMultiDim}, "ArrayLenCap": {funcName: "ArrayLenCap", native: divergence_hunt119.ArrayLenCap},
	}})
}
func TestDivergenceHunt120(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt120Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Integration1": {funcName: "Integration1", native: divergence_hunt120.Integration1}, "Integration2": {funcName: "Integration2", native: divergence_hunt120.Integration2}, "Integration3": {funcName: "Integration3", native: divergence_hunt120.Integration3}, "Integration4": {funcName: "Integration4", native: divergence_hunt120.Integration4}, "Integration5": {funcName: "Integration5", native: divergence_hunt120.Integration5}, "Integration6": {funcName: "Integration6", native: divergence_hunt120.Integration6}, "Integration7": {funcName: "Integration7", native: divergence_hunt120.Integration7}, "Integration8": {funcName: "Integration8", native: divergence_hunt120.Integration8}, "Integration9": {funcName: "Integration9", native: divergence_hunt120.Integration9}, "Integration10": {funcName: "Integration10", native: divergence_hunt120.Integration10},
	}})
}
func TestDivergenceHunt121(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt121Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ChanSelectDefault": {funcName: "ChanSelectDefault", native: divergence_hunt121.ChanSelectDefault}, "ChanSelectReady": {funcName: "ChanSelectReady", native: divergence_hunt121.ChanSelectReady}, "ChanBufferedSend": {funcName: "ChanBufferedSend", native: divergence_hunt121.ChanBufferedSend}, "ChanNilBlock": {funcName: "ChanNilBlock", native: divergence_hunt121.ChanNilBlock}, "ChanClosedReceive": {funcName: "ChanClosedReceive", native: divergence_hunt121.ChanClosedReceive}, "ChanClosedEmpty": {funcName: "ChanClosedEmpty", native: divergence_hunt121.ChanClosedEmpty}, "ChanSelectMultiReady": {funcName: "ChanSelectMultiReady", native: divergence_hunt121.ChanSelectMultiReady}, "ChanLenCap": {funcName: "ChanLenCap", native: divergence_hunt121.ChanLenCap}, "ChanSelectWithAssign": {funcName: "ChanSelectWithAssign", native: divergence_hunt121.ChanSelectWithAssign},
	}})
}
func TestDivergenceHunt122(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt122Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InterfaceEmbedMethod": {funcName: "InterfaceEmbedMethod", native: divergence_hunt122.InterfaceEmbedMethod}, "InterfaceEmbedInterface": {funcName: "InterfaceEmbedInterface", native: divergence_hunt122.InterfaceEmbedInterface}, "InterfaceEmbedFieldAccess": {funcName: "InterfaceEmbedFieldAccess", native: divergence_hunt122.InterfaceEmbedFieldAccess}, "InterfaceEmbedPromoted": {funcName: "InterfaceEmbedPromoted", native: divergence_hunt122.InterfaceEmbedPromoted}, "InterfaceEmbedOverride": {funcName: "InterfaceEmbedOverride", native: divergence_hunt122.InterfaceEmbedOverride}, "InterfaceNilCheck": {funcName: "InterfaceNilCheck", native: divergence_hunt122.InterfaceNilCheck}, "InterfaceNilTypedCheck": {funcName: "InterfaceNilTypedCheck", native: divergence_hunt122.InterfaceNilTypedCheck}, "InterfaceStructLiteral": {funcName: "InterfaceStructLiteral", native: divergence_hunt122.InterfaceStructLiteral},
	}})
}
func TestDivergenceHunt123(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt123Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SliceAppendNil": {funcName: "SliceAppendNil", native: divergence_hunt123.SliceAppendNil}, "SliceAppendExpand": {funcName: "SliceAppendExpand", native: divergence_hunt123.SliceAppendExpand}, "SliceCopyCount": {funcName: "SliceCopyCount", native: divergence_hunt123.SliceCopyCount}, "SliceCopyOverlap": {funcName: "SliceCopyOverlap", native: divergence_hunt123.SliceCopyOverlap}, "SliceDeleteElement": {funcName: "SliceDeleteElement", native: divergence_hunt123.SliceDeleteElement}, "SliceThreeIndex": {funcName: "SliceThreeIndex", native: divergence_hunt123.SliceThreeIndex}, "SliceNilAppend": {funcName: "SliceNilAppend", native: divergence_hunt123.SliceNilAppend}, "SliceNilCopy": {funcName: "SliceNilCopy", native: divergence_hunt123.SliceNilCopy}, "SliceAppendSlice": {funcName: "SliceAppendSlice", native: divergence_hunt123.SliceAppendSlice}, "SliceCapAfterAppend": {funcName: "SliceCapAfterAppend", native: divergence_hunt123.SliceCapAfterAppend},
	}})
}
func TestDivergenceHunt124(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt124Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapLiteral": {funcName: "MapLiteral", native: divergence_hunt124.MapLiteral}, "MapDeleteLen": {funcName: "MapDeleteLen", native: divergence_hunt124.MapDeleteLen}, "MapNilWrite": {funcName: "MapNilWrite", native: divergence_hunt124.MapNilWrite}, "MapZeroValue": {funcName: "MapZeroValue", native: divergence_hunt124.MapZeroValue}, "MapOkCheck": {funcName: "MapOkCheck", native: divergence_hunt124.MapOkCheck}, "MapSortedKeys": {funcName: "MapSortedKeys", native: divergence_hunt124.MapSortedKeys}, "MapIntKey": {funcName: "MapIntKey", native: divergence_hunt124.MapIntKey}, "MapNestedMap": {funcName: "MapNestedMap", native: divergence_hunt124.MapNestedMap}, "MapUpdateValue": {funcName: "MapUpdateValue", native: divergence_hunt124.MapUpdateValue}, "MapBoolValue": {funcName: "MapBoolValue", native: divergence_hunt124.MapBoolValue},
	}})
}
func TestDivergenceHunt125(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt125Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringLenBytes": {funcName: "StringLenBytes", native: divergence_hunt125.StringLenBytes}, "StringLenRunes": {funcName: "StringLenRunes", native: divergence_hunt125.StringLenRunes}, "StringRuneAt": {funcName: "StringRuneAt", native: divergence_hunt125.StringRuneAt}, "StringRangeRunes": {funcName: "StringRangeRunes", native: divergence_hunt125.StringRangeRunes}, "StringConcat": {funcName: "StringConcat", native: divergence_hunt125.StringConcat}, "StringCompare": {funcName: "StringCompare", native: divergence_hunt125.StringCompare}, "StringSliceBytes": {funcName: "StringSliceBytes", native: divergence_hunt125.StringSliceBytes}, "StringByteConversion": {funcName: "StringByteConversion", native: divergence_hunt125.StringByteConversion}, "StringRuneConversion": {funcName: "StringRuneConversion", native: divergence_hunt125.StringRuneConversion}, "StringEmptyCheck": {funcName: "StringEmptyCheck", native: divergence_hunt125.StringEmptyCheck},
	}})
}
func TestDivergenceHunt126(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt126Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SwitchFallthrough": {funcName: "SwitchFallthrough", native: divergence_hunt126.SwitchFallthrough}, "SwitchNoFallthrough": {funcName: "SwitchNoFallthrough", native: divergence_hunt126.SwitchNoFallthrough}, "SwitchDefaultOnly": {funcName: "SwitchDefaultOnly", native: divergence_hunt126.SwitchDefaultOnly}, "SwitchCaseOrder": {funcName: "SwitchCaseOrder", native: divergence_hunt126.SwitchCaseOrder}, "SwitchStringCase": {funcName: "SwitchStringCase", native: divergence_hunt126.SwitchStringCase}, "SwitchNoCaseNoDefault": {funcName: "SwitchNoCaseNoDefault", native: divergence_hunt126.SwitchNoCaseNoDefault}, "SwitchBreakExplicit": {funcName: "SwitchBreakExplicit", native: divergence_hunt126.SwitchBreakExplicit}, "SwitchMultiCase": {funcName: "SwitchMultiCase", native: divergence_hunt126.SwitchMultiCase}, "SwitchInLoop": {funcName: "SwitchInLoop", native: divergence_hunt126.SwitchInLoop},
	}})
}
func TestDivergenceHunt127(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt127Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TypeSwitchBasic": {funcName: "TypeSwitchBasic", native: divergence_hunt127.TypeSwitchBasic}, "TypeSwitchString": {funcName: "TypeSwitchString", native: divergence_hunt127.TypeSwitchString}, "TypeSwitchNil": {funcName: "TypeSwitchNil", native: divergence_hunt127.TypeSwitchNil}, "TypeAssertionOk": {funcName: "TypeAssertionOk", native: divergence_hunt127.TypeAssertionOk}, "TypeAssertionFail": {funcName: "TypeAssertionFail", native: divergence_hunt127.TypeAssertionFail}, "TypeAssertionPanicFree": {funcName: "TypeAssertionPanicFree", native: divergence_hunt127.TypeAssertionPanicFree}, "TypeSwitchMultiCase": {funcName: "TypeSwitchMultiCase", native: divergence_hunt127.TypeSwitchMultiCase}, "TypeSwitchStruct": {funcName: "TypeSwitchStruct", native: divergence_hunt127.TypeSwitchStruct}, "TypeAssertionChain": {funcName: "TypeAssertionChain", native: divergence_hunt127.TypeAssertionChain},
	}})
}
func TestDivergenceHunt128(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt128Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SwapVariables": {funcName: "SwapVariables", native: divergence_hunt128.SwapVariables}, "MultiReturnAssign": {funcName: "MultiReturnAssign", native: divergence_hunt128.MultiReturnAssign}, "BlankAssign": {funcName: "BlankAssign", native: divergence_hunt128.BlankAssign}, "MultiAssignExpression": {funcName: "MultiAssignExpression", native: divergence_hunt128.MultiAssignExpression}, "MultiAssignSwap": {funcName: "MultiAssignSwap", native: divergence_hunt128.MultiAssignSwap}, "MultiAssignMap": {funcName: "MultiAssignMap", native: divergence_hunt128.MultiAssignMap}, "NestedMultiReturn": {funcName: "NestedMultiReturn", native: divergence_hunt128.NestedMultiReturn}, "AssignDifferentTypes": {funcName: "AssignDifferentTypes", native: divergence_hunt128.AssignDifferentTypes}, "MultiAssignStruct": {funcName: "MultiAssignStruct", native: divergence_hunt128.MultiAssignStruct},
	}})
}
func TestDivergenceHunt129(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt129Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StructTagJSON": {funcName: "StructTagJSON", native: divergence_hunt129.StructTagJSON}, "StructTagUnmarshal": {funcName: "StructTagUnmarshal", native: divergence_hunt129.StructTagUnmarshal}, "StructTagOmitEmpty": {funcName: "StructTagOmitEmpty", native: divergence_hunt129.StructTagOmitEmpty}, "StructNestedJSON": {funcName: "StructNestedJSON", native: divergence_hunt129.StructNestedJSON}, "StructMapJSON": {funcName: "StructMapJSON", native: divergence_hunt129.StructMapJSON}, "StructSliceJSON": {funcName: "StructSliceJSON", native: divergence_hunt129.StructSliceJSON}, "StructSliceOfStructs": {funcName: "StructSliceOfStructs", native: divergence_hunt129.StructSliceOfStructs}, "StructBoolJSON": {funcName: "StructBoolJSON", native: divergence_hunt129.StructBoolJSON}, "StructNilJSON": {funcName: "StructNilJSON", native: divergence_hunt129.StructNilJSON},
	}})
}
func TestDivergenceHunt130(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt130Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"DeferStackOrder": {funcName: "DeferStackOrder", native: divergence_hunt130.DeferStackOrder}, "DeferModifyReturn": {funcName: "DeferModifyReturn", native: divergence_hunt130.DeferModifyReturn}, "DeferNamedReturn": {funcName: "DeferNamedReturn", native: divergence_hunt130.DeferNamedReturn}, "DeferCaptureValue": {funcName: "DeferCaptureValue", native: divergence_hunt130.DeferCaptureValue}, "DeferCapturePointer": {funcName: "DeferCapturePointer", native: divergence_hunt130.DeferCapturePointer}, "RecoverBasic": {funcName: "RecoverBasic", native: divergence_hunt130.RecoverBasic}, "RecoverInDefer": {funcName: "RecoverInDefer", native: divergence_hunt130.RecoverInDefer}, "RecoverNoPanic": {funcName: "RecoverNoPanic", native: divergence_hunt130.RecoverNoPanic}, "DeferMultipleRecovers": {funcName: "DeferMultipleRecovers", native: divergence_hunt130.DeferMultipleRecovers}, "DeferPanicInDefer": {funcName: "DeferPanicInDefer", native: divergence_hunt130.DeferPanicInDefer},
	}})
}
func TestDivergenceHunt131(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt131Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"GoroutineWaitGroup": {funcName: "GoroutineWaitGroup", native: divergence_hunt131.GoroutineWaitGroup}, "GoroutineOnce": {funcName: "GoroutineOnce", native: divergence_hunt131.GoroutineOnce}, "GoroutineChannelSum": {funcName: "GoroutineChannelSum", native: divergence_hunt131.GoroutineChannelSum}, "GoroutineMutex": {funcName: "GoroutineMutex", native: divergence_hunt131.GoroutineMutex}, "GoroutineSelectTimeout": {funcName: "GoroutineSelectTimeout", native: divergence_hunt131.GoroutineSelectTimeout}, "GoroutineSendReceive": {funcName: "GoroutineSendReceive", native: divergence_hunt131.GoroutineSendReceive}, "GoroutineCloseSignal": {funcName: "GoroutineCloseSignal", native: divergence_hunt131.GoroutineCloseSignal}, "GoroutinePanicRecover": {funcName: "GoroutinePanicRecover", native: divergence_hunt131.GoroutinePanicRecover},
	}})
}
func TestDivergenceHunt132(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt132Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SliceGrowFromEmpty": {funcName: "SliceGrowFromEmpty", native: divergence_hunt132.SliceGrowFromEmpty}, "SliceGrowWithCap": {funcName: "SliceGrowWithCap", native: divergence_hunt132.SliceGrowWithCap}, "SliceReslice": {funcName: "SliceReslice", native: divergence_hunt132.SliceReslice}, "SliceResliceCap": {funcName: "SliceResliceCap", native: divergence_hunt132.SliceResliceCap}, "SliceAppendBeyondCap": {funcName: "SliceAppendBeyondCap", native: divergence_hunt132.SliceAppendBeyondCap}, "SliceMakeZeroLen": {funcName: "SliceMakeZeroLen", native: divergence_hunt132.SliceMakeZeroLen}, "SliceNilVsEmpty": {funcName: "SliceNilVsEmpty", native: divergence_hunt132.SliceNilVsEmpty}, "SliceOfString": {funcName: "SliceOfString", native: divergence_hunt132.SliceOfString}, "SliceBool": {funcName: "SliceBool", native: divergence_hunt132.SliceBool}, "SliceStructLiteral": {funcName: "SliceStructLiteral", native: divergence_hunt132.SliceStructLiteral},
	}})
}
func TestDivergenceHunt133(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt133Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RecursiveClosure": {funcName: "RecursiveClosure", native: divergence_hunt133.RecursiveClosure}, "ClosureCounter": {funcName: "ClosureCounter", native: divergence_hunt133.ClosureCounter}, "ClosureCapture": {funcName: "ClosureCapture", native: divergence_hunt133.ClosureCapture}, "ClosureParamCapture": {funcName: "ClosureParamCapture", native: divergence_hunt133.ClosureParamCapture}, "ClosureSliceCapture": {funcName: "ClosureSliceCapture", native: divergence_hunt133.ClosureSliceCapture}, "ClosureSliceCaptureNoCopy": {funcName: "ClosureSliceCaptureNoCopy", native: divergence_hunt133.ClosureSliceCaptureNoCopy}, "MutualRecursion": {funcName: "MutualRecursion", native: divergence_hunt133.MutualRecursion}, "ClosureAsParam": {funcName: "ClosureAsParam", native: divergence_hunt133.ClosureAsParam}, "ClosureReturnClosure": {funcName: "ClosureReturnClosure", native: divergence_hunt133.ClosureReturnClosure},
	}})
}
func TestDivergenceHunt134(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt134Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NamedReturnBasic": {funcName: "NamedReturnBasic", native: divergence_hunt134.NamedReturnBasic}, "NamedReturnOverride": {funcName: "NamedReturnOverride", native: divergence_hunt134.NamedReturnOverride}, "NamedReturnDefer": {funcName: "NamedReturnDefer", native: divergence_hunt134.NamedReturnDefer}, "NamedReturnDeferDouble": {funcName: "NamedReturnDeferDouble", native: divergence_hunt134.NamedReturnDeferDouble}, "NamedReturnMulti": {funcName: "NamedReturnMulti", native: divergence_hunt134.NamedReturnMulti}, "NamedReturnDeferMulti": {funcName: "NamedReturnDeferMulti", native: divergence_hunt134.NamedReturnDeferMulti}, "NamedReturnShadow": {funcName: "NamedReturnShadow", native: divergence_hunt134.NamedReturnShadow}, "NamedReturnZeroValue": {funcName: "NamedReturnZeroValue", native: divergence_hunt134.NamedReturnZeroValue}, "NamedReturnPanicRecover": {funcName: "NamedReturnPanicRecover", native: divergence_hunt134.NamedReturnPanicRecover}, "NamedReturnDeferModify": {funcName: "NamedReturnDeferModify", native: divergence_hunt134.NamedReturnDeferModify},
	}})
}
func TestDivergenceHunt135(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt135Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"IotaBasic": {funcName: "IotaBasic", native: divergence_hunt135.IotaBasic}, "ConstExplicit": {funcName: "ConstExplicit", native: divergence_hunt135.ConstExplicit}, "IotaExpression": {funcName: "IotaExpression", native: divergence_hunt135.IotaExpression}, "ConstUntyped": {funcName: "ConstUntyped", native: divergence_hunt135.ConstUntyped}, "ConstTyped": {funcName: "ConstTyped", native: divergence_hunt135.ConstTyped}, "ConstString": {funcName: "ConstString", native: divergence_hunt135.ConstString}, "ConstBool": {funcName: "ConstBool", native: divergence_hunt135.ConstBool}, "ConstExpression": {funcName: "ConstExpression", native: divergence_hunt135.ConstExpression}, "IotaSkip": {funcName: "IotaSkip", native: divergence_hunt135.IotaSkip},
	}})
}
func TestDivergenceHunt136(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt136Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"DeepFieldAccess": {funcName: "DeepFieldAccess", native: divergence_hunt136.DeepFieldAccess}, "DeepFieldAssign": {funcName: "DeepFieldAssign", native: divergence_hunt136.DeepFieldAssign}, "DeepFieldPointer": {funcName: "DeepFieldPointer", native: divergence_hunt136.DeepFieldPointer}, "LinkedListTraversal": {funcName: "LinkedListTraversal", native: divergence_hunt136.LinkedListTraversal}, "LinkedListCreate": {funcName: "LinkedListCreate", native: divergence_hunt136.LinkedListCreate}, "TreeTraversal": {funcName: "TreeTraversal", native: divergence_hunt136.TreeTraversal}, "NestedStructLiteral": {funcName: "NestedStructLiteral", native: divergence_hunt136.NestedStructLiteral}, "NestedStructUpdate": {funcName: "NestedStructUpdate", native: divergence_hunt136.NestedStructUpdate},
	}})
}
func TestDivergenceHunt137(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt137Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MethodValueExpr": {funcName: "MethodValueExpr", native: divergence_hunt137.MethodValueExpr}, "MethodCallDirect": {funcName: "MethodCallDirect", native: divergence_hunt137.MethodCallDirect}, "MethodValueReceiver": {funcName: "MethodValueReceiver", native: divergence_hunt137.MethodValueReceiver}, "MethodPtrReceiver": {funcName: "MethodPtrReceiver", native: divergence_hunt137.MethodPtrReceiver}, "MethodOnLiteral": {funcName: "MethodOnLiteral", native: divergence_hunt137.MethodOnLiteral}, "MethodWithArgs": {funcName: "MethodWithArgs", native: divergence_hunt137.MethodWithArgs}, "MethodStringer": {funcName: "MethodStringer", native: divergence_hunt137.MethodStringer}, "MethodStackPushPop": {funcName: "MethodStackPushPop", native: divergence_hunt137.MethodStackPushPop},
	}})
}
func TestDivergenceHunt138(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt138Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NewInt": {funcName: "NewInt", native: divergence_hunt138.NewInt}, "NewStruct": {funcName: "NewStruct", native: divergence_hunt138.NewStruct}, "AddressOf": {funcName: "AddressOf", native: divergence_hunt138.AddressOf}, "NilPointerCheck": {funcName: "NilPointerCheck", native: divergence_hunt138.NilPointerCheck}, "PointerSwap": {funcName: "PointerSwap", native: divergence_hunt138.PointerSwap}, "PointerToSlice": {funcName: "PointerToSlice", native: divergence_hunt138.PointerToSlice}, "PointerToMap": {funcName: "PointerToMap", native: divergence_hunt138.PointerToMap}, "PointerStructMethod": {funcName: "PointerStructMethod", native: divergence_hunt138.PointerStructMethod}, "DoublePointer": {funcName: "DoublePointer", native: divergence_hunt138.DoublePointer}, "PointerArray": {funcName: "PointerArray", native: divergence_hunt138.PointerArray},
	}})
}
func TestDivergenceHunt139(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt139Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InterfaceNilComparison": {funcName: "InterfaceNilComparison", native: divergence_hunt139.InterfaceNilComparison}, "InterfaceTypedNil": {funcName: "InterfaceTypedNil", native: divergence_hunt139.InterfaceTypedNil}, "InterfaceNilTypeAssertion": {funcName: "InterfaceNilTypeAssertion", native: divergence_hunt139.InterfaceNilTypeAssertion}, "InterfaceSliceOfNil": {funcName: "InterfaceSliceOfNil", native: divergence_hunt139.InterfaceSliceOfNil}, "InterfaceMapNilValue": {funcName: "InterfaceMapNilValue", native: divergence_hunt139.InterfaceMapNilValue}, "InterfaceFuncReturn": {funcName: "InterfaceFuncReturn", native: divergence_hunt139.InterfaceFuncReturn}, "InterfaceStructMethodNil": {funcName: "InterfaceStructMethodNil", native: divergence_hunt139.InterfaceStructMethodNil}, "InterfaceEmptySlice": {funcName: "InterfaceEmptySlice", native: divergence_hunt139.InterfaceEmptySlice}, "InterfaceNonNilCheck": {funcName: "InterfaceNonNilCheck", native: divergence_hunt139.InterfaceNonNilCheck}, "InterfaceNilSliceAssign": {funcName: "InterfaceNilSliceAssign", native: divergence_hunt139.InterfaceNilSliceAssign},
	}})
}
func TestDivergenceHunt140(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt140Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"EmbeddedFieldAccess": {funcName: "EmbeddedFieldAccess", native: divergence_hunt140.EmbeddedFieldAccess}, "EmbeddedMethodCall": {funcName: "EmbeddedMethodCall", native: divergence_hunt140.EmbeddedMethodCall}, "EmbeddedFieldExplicit": {funcName: "EmbeddedFieldExplicit", native: divergence_hunt140.EmbeddedFieldExplicit}, "EmbeddedChain": {funcName: "EmbeddedChain", native: divergence_hunt140.EmbeddedChain}, "EmbeddedChainMethod": {funcName: "EmbeddedChainMethod", native: divergence_hunt140.EmbeddedChainMethod}, "EmbeddedOverride": {funcName: "EmbeddedOverride", native: divergence_hunt140.EmbeddedOverride}, "EmbeddedPointer": {funcName: "EmbeddedPointer", native: divergence_hunt140.EmbeddedPointer}, "EmbeddedBoolField": {funcName: "EmbeddedBoolField", native: divergence_hunt140.EmbeddedBoolField},
	}})
}
func TestDivergenceHunt141(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt141Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ClosureLoopCapture": {funcName: "ClosureLoopCapture", native: divergence_hunt141.ClosureLoopCapture}, "ClosureLoopDeferred": {funcName: "ClosureLoopDeferred", native: divergence_hunt141.ClosureLoopDeferred}, "ClosureShadowVar": {funcName: "ClosureShadowVar", native: divergence_hunt141.ClosureShadowVar}, "ClosureMutateOuter": {funcName: "ClosureMutateOuter", native: divergence_hunt141.ClosureMutateOuter}, "ClosureMultipleCaptures": {funcName: "ClosureMultipleCaptures", native: divergence_hunt141.ClosureMultipleCaptures}, "ClosureReturned": {funcName: "ClosureReturned", native: divergence_hunt141.ClosureReturned}, "ClosureSliceAppend": {funcName: "ClosureSliceAppend", native: divergence_hunt141.ClosureSliceAppend}, "ClosureMapModify": {funcName: "ClosureMapModify", native: divergence_hunt141.ClosureMapModify}, "ClosureNested": {funcName: "ClosureNested", native: divergence_hunt141.ClosureNested},
	}})
}
func TestDivergenceHunt142(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt142Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ThreeIndexBasic": {funcName: "ThreeIndexBasic", native: divergence_hunt142.ThreeIndexBasic}, "ThreeIndexAppendNoGrow": {funcName: "ThreeIndexAppendNoGrow", native: divergence_hunt142.ThreeIndexAppendNoGrow}, "ThreeIndexFullSlice": {funcName: "ThreeIndexFullSlice", native: divergence_hunt142.ThreeIndexFullSlice}, "AppendCopyPattern": {funcName: "AppendCopyPattern", native: divergence_hunt142.AppendCopyPattern}, "SliceInsertMiddle": {funcName: "SliceInsertMiddle", native: divergence_hunt142.SliceInsertMiddle}, "SliceFilter": {funcName: "SliceFilter", native: divergence_hunt142.SliceFilter}, "SliceReverse": {funcName: "SliceReverse", native: divergence_hunt142.SliceReverse}, "SliceClone": {funcName: "SliceClone", native: divergence_hunt142.SliceClone}, "SliceStackPattern": {funcName: "SliceStackPattern", native: divergence_hunt142.SliceStackPattern},
	}})
}
func TestDivergenceHunt143(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt143Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapInterfaceValues": {funcName: "MapInterfaceValues", native: divergence_hunt143.MapInterfaceValues}, "MapInterfaceTypeSwitch": {funcName: "MapInterfaceTypeSwitch", native: divergence_hunt143.MapInterfaceTypeSwitch}, "MapInterfaceAssertion": {funcName: "MapInterfaceAssertion", native: divergence_hunt143.MapInterfaceAssertion}, "MapStringSlice": {funcName: "MapStringSlice", native: divergence_hunt143.MapStringSlice}, "MapStringFunc": {funcName: "MapStringFunc", native: divergence_hunt143.MapStringFunc}, "MapDeleteAndRead": {funcName: "MapDeleteAndRead", native: divergence_hunt143.MapDeleteAndRead}, "MapLengthAfterDelete": {funcName: "MapLengthAfterDelete", native: divergence_hunt143.MapLengthAfterDelete}, "MapNilVsEmptyAccess": {funcName: "MapNilVsEmptyAccess", native: divergence_hunt143.MapNilVsEmptyAccess}, "MapCompositeLiteral": {funcName: "MapCompositeLiteral", native: divergence_hunt143.MapCompositeLiteral},
	}})
}
func TestDivergenceHunt144(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt144Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StructEqual": {funcName: "StructEqual", native: divergence_hunt144.StructEqual}, "StructNotEqual": {funcName: "StructNotEqual", native: divergence_hunt144.StructNotEqual}, "StructZeroValue": {funcName: "StructZeroValue", native: divergence_hunt144.StructZeroValue}, "StructCopy": {funcName: "StructCopy", native: divergence_hunt144.StructCopy}, "StructPointerCopy": {funcName: "StructPointerCopy", native: divergence_hunt144.StructPointerCopy}, "StructStringField": {funcName: "StructStringField", native: divergence_hunt144.StructStringField}, "StructBoolField": {funcName: "StructBoolField", native: divergence_hunt144.StructBoolField}, "StructSliceField": {funcName: "StructSliceField", native: divergence_hunt144.StructSliceField}, "StructMapField": {funcName: "StructMapField", native: divergence_hunt144.StructMapField}, "StructEmbeddedCompare": {funcName: "StructEmbeddedCompare", native: divergence_hunt144.StructEmbeddedCompare},
	}})
}
func TestDivergenceHunt145(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt145Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BreakBasic": {funcName: "BreakBasic", native: divergence_hunt145.BreakBasic}, "ContinueBasic": {funcName: "ContinueBasic", native: divergence_hunt145.ContinueBasic}, "LabeledBreak": {funcName: "LabeledBreak", native: divergence_hunt145.LabeledBreak}, "LabeledContinue": {funcName: "LabeledContinue", native: divergence_hunt145.LabeledContinue}, "RangeBreak": {funcName: "RangeBreak", native: divergence_hunt145.RangeBreak}, "RangeContinue": {funcName: "RangeContinue", native: divergence_hunt145.RangeContinue}, "NestedLoopBreak": {funcName: "NestedLoopBreak", native: divergence_hunt145.NestedLoopBreak}, "SwitchBreakInLoop": {funcName: "SwitchBreakInLoop", native: divergence_hunt145.SwitchBreakInLoop},
	}})
}
func TestDivergenceHunt146(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt146Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ErrorAsStructPointer": {funcName: "ErrorAsStructPointer", native: divergence_hunt146.ErrorAsStructPointer}, "ErrorChainUnwrap": {funcName: "ErrorChainUnwrap", native: divergence_hunt146.ErrorChainUnwrap}, "ErrorNilInterface": {funcName: "ErrorNilInterface", native: divergence_hunt146.ErrorNilInterface}, "ErrorTypeAssertion": {funcName: "ErrorTypeAssertion", native: divergence_hunt146.ErrorTypeAssertion}, "ErrorInterfaceAssertion": {funcName: "ErrorInterfaceAssertion", native: divergence_hunt146.ErrorInterfaceAssertion}, "ErrorSpecificMethod": {funcName: "ErrorSpecificMethod", native: divergence_hunt146.ErrorSpecificMethod}, "ErrorMultiWrap": {funcName: "ErrorMultiWrap", native: divergence_hunt146.ErrorMultiWrap},
	}})
}
func TestDivergenceHunt147(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt147Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringsJoin": {funcName: "StringsJoin", native: divergence_hunt147.StringsJoin}, "StringsSplit": {funcName: "StringsSplit", native: divergence_hunt147.StringsSplit}, "StringsContains": {funcName: "StringsContains", native: divergence_hunt147.StringsContains}, "StringsHasPrefix": {funcName: "StringsHasPrefix", native: divergence_hunt147.StringsHasPrefix}, "StringsHasSuffix": {funcName: "StringsHasSuffix", native: divergence_hunt147.StringsHasSuffix}, "StringsTrimSpace": {funcName: "StringsTrimSpace", native: divergence_hunt147.StringsTrimSpace}, "StringsReplace": {funcName: "StringsReplace", native: divergence_hunt147.StringsReplace}, "StringsToUpper": {funcName: "StringsToUpper", native: divergence_hunt147.StringsToUpper}, "StringsRepeat": {funcName: "StringsRepeat", native: divergence_hunt147.StringsRepeat}, "StringsCount": {funcName: "StringsCount", native: divergence_hunt147.StringsCount}, "StringsIndex": {funcName: "StringsIndex", native: divergence_hunt147.StringsIndex},
	}})
}
func TestDivergenceHunt148(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt148Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ErrorsNewCheck": {funcName: "ErrorsNewCheck", native: divergence_hunt148.ErrorsNewCheck}, "ErrorsIsMatch": {funcName: "ErrorsIsMatch", native: divergence_hunt148.ErrorsIsMatch}, "ErrorsAsInterface": {funcName: "ErrorsAsInterface", native: divergence_hunt148.ErrorsAsInterface}, "ErrorsUnwrapNil": {funcName: "ErrorsUnwrapNil", native: divergence_hunt148.ErrorsUnwrapNil}, "ErrorfWrapUnwrap": {funcName: "ErrorfWrapUnwrap", native: divergence_hunt148.ErrorfWrapUnwrap}, "ErrorfMultiWrap": {funcName: "ErrorfMultiWrap", native: divergence_hunt148.ErrorfMultiWrap}, "ErrorJoin": {funcName: "ErrorJoin", native: divergence_hunt148.ErrorJoin}, "ErrorJoinIs": {funcName: "ErrorJoinIs", native: divergence_hunt148.ErrorJoinIs}, "ErrorNilIs": {funcName: "ErrorNilIs", native: divergence_hunt148.ErrorNilIs},
	}})
}
func TestDivergenceHunt149(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt149Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"VariadicSum": {funcName: "VariadicSum", native: divergence_hunt149.VariadicSum}, "VariadicStringJoin": {funcName: "VariadicStringJoin", native: divergence_hunt149.VariadicStringJoin}, "VariadicEmpty": {funcName: "VariadicEmpty", native: divergence_hunt149.VariadicEmpty}, "VariadicInterface": {funcName: "VariadicInterface", native: divergence_hunt149.VariadicInterface}, "VariadicSpread": {funcName: "VariadicSpread", native: divergence_hunt149.VariadicSpread}, "VariadicFmt": {funcName: "VariadicFmt", native: divergence_hunt149.VariadicFmt}, "VariadicSort": {funcName: "VariadicSort", native: divergence_hunt149.VariadicSort}, "VariadicStrconv": {funcName: "VariadicStrconv", native: divergence_hunt149.VariadicStrconv}, "VariadicPrintf": {funcName: "VariadicPrintf", native: divergence_hunt149.VariadicPrintf},
	}})
}
func TestDivergenceHunt150(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt150Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"IntegrationStructMethod": {funcName: "IntegrationStructMethod", native: divergence_hunt150.IntegrationStructMethod}, "IntegrationStructMutation": {funcName: "IntegrationStructMutation", native: divergence_hunt150.IntegrationStructMutation}, "IntegrationSliceMapFilter": {funcName: "IntegrationSliceMapFilter", native: divergence_hunt150.IntegrationSliceMapFilter}, "IntegrationErrorChain": {funcName: "IntegrationErrorChain", native: divergence_hunt150.IntegrationErrorChain}, "IntegrationStringProcess": {funcName: "IntegrationStringProcess", native: divergence_hunt150.IntegrationStringProcess}, "IntegrationClosureCounter": {funcName: "IntegrationClosureCounter", native: divergence_hunt150.IntegrationClosureCounter}, "IntegrationDeferRecover": {funcName: "IntegrationDeferRecover", native: divergence_hunt150.IntegrationDeferRecover}, "IntegrationPointerChain": {funcName: "IntegrationPointerChain", native: divergence_hunt150.IntegrationPointerChain}, "IntegrationTypeSwitch": {funcName: "IntegrationTypeSwitch", native: divergence_hunt150.IntegrationTypeSwitch}, "IntegrationNamedReturn": {funcName: "IntegrationNamedReturn", native: divergence_hunt150.IntegrationNamedReturn},
	}})
}
func TestDivergenceHunt151(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt151Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SelectBasic": {funcName: "SelectBasic", native: divergence_hunt151.SelectBasic}, "SelectDefaultOnly": {funcName: "SelectDefaultOnly", native: divergence_hunt151.SelectDefaultOnly}, "SelectNoDefault": {funcName: "SelectNoDefault", native: divergence_hunt151.SelectNoDefault}, "SelectMultipleChannels": {funcName: "SelectMultipleChannels", native: divergence_hunt151.SelectMultipleChannels}, "ChannelBufferedCapacity": {funcName: "ChannelBufferedCapacity", native: divergence_hunt151.ChannelBufferedCapacity}, "ChannelNilSend": {funcName: "ChannelNilSend", native: divergence_hunt151.ChannelNilSend}, "ChannelNilReceive": {funcName: "ChannelNilReceive", native: divergence_hunt151.ChannelNilReceive}, "ChannelCloseCheck": {funcName: "ChannelCloseCheck", native: divergence_hunt151.ChannelCloseCheck}, "ChannelRange": {funcName: "ChannelRange", native: divergence_hunt151.ChannelRange}, "ChannelSendReceiveOrder": {funcName: "ChannelSendReceiveOrder", native: divergence_hunt151.ChannelSendReceiveOrder},
	}})
}
func TestDivergenceHunt152(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt152Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SortIntsReverse": {funcName: "SortIntsReverse", native: divergence_hunt152.SortIntsReverse}, "SortStrings": {funcName: "SortStrings", native: divergence_hunt152.SortStrings}, "SortFloat64s": {funcName: "SortFloat64s", native: divergence_hunt152.SortFloat64s}, "SortSearch": {funcName: "SortSearch", native: divergence_hunt152.SortSearch}, "SortSearchInts": {funcName: "SortSearchInts", native: divergence_hunt152.SortSearchInts}, "SliceIsSorted": {funcName: "SliceIsSorted", native: divergence_hunt152.SliceIsSorted}, "SliceSortStable": {funcName: "SliceSortStable", native: divergence_hunt152.SliceSortStable}, "SliceCut": {funcName: "SliceCut", native: divergence_hunt152.SliceCut}, "SliceInsert": {funcName: "SliceInsert", native: divergence_hunt152.SliceInsert}, "SliceDeleteUnordered": {funcName: "SliceDeleteUnordered", native: divergence_hunt152.SliceDeleteUnordered}, "SliceCompact": {funcName: "SliceCompact", native: divergence_hunt152.SliceCompact},
	}})
}
func TestDivergenceHunt153(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt153Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BitsLeadingZeros": {funcName: "BitsLeadingZeros", native: divergence_hunt153.BitsLeadingZeros}, "BitsTrailingZeros": {funcName: "BitsTrailingZeros", native: divergence_hunt153.BitsTrailingZeros}, "BitsOnesCount": {funcName: "BitsOnesCount", native: divergence_hunt153.BitsOnesCount}, "BitsRotate": {funcName: "BitsRotate", native: divergence_hunt153.BitsRotate}, "BitsReverse": {funcName: "BitsReverse", native: divergence_hunt153.BitsReverse}, "BitsLen": {funcName: "BitsLen", native: divergence_hunt153.BitsLen}, "MathAbs": {funcName: "MathAbs", native: divergence_hunt153.MathAbs}, "MathMinMax": {funcName: "MathMinMax", native: divergence_hunt153.MathMinMax}, "MathFloorCeil": {funcName: "MathFloorCeil", native: divergence_hunt153.MathFloorCeil}, "MathRound": {funcName: "MathRound", native: divergence_hunt153.MathRound}, "MathTrunc": {funcName: "MathTrunc", native: divergence_hunt153.MathTrunc}, "MathMod": {funcName: "MathMod", native: divergence_hunt153.MathMod}, "MathPow": {funcName: "MathPow", native: divergence_hunt153.MathPow}, "MathSqrt": {funcName: "MathSqrt", native: divergence_hunt153.MathSqrt}, "MathNaNCheck": {funcName: "MathNaNCheck", native: divergence_hunt153.MathNaNCheck}, "MathInfCheck": {funcName: "MathInfCheck", native: divergence_hunt153.MathInfCheck}, "MathCopySign": {funcName: "MathCopySign", native: divergence_hunt153.MathCopySign}, "MathPiE": {funcName: "MathPiE", native: divergence_hunt153.MathPiE},
	}})
}
func TestDivergenceHunt154(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt154Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"AnonymousStructLiteral": {funcName: "AnonymousStructLiteral", native: divergence_hunt154.AnonymousStructLiteral}, "AnonymousStructSlice": {funcName: "AnonymousStructSlice", native: divergence_hunt154.AnonymousStructSlice}, "AnonymousStructMap": {funcName: "AnonymousStructMap", native: divergence_hunt154.AnonymousStructMap}, "FunctionTypeComparison": {funcName: "FunctionTypeComparison", native: divergence_hunt154.FunctionTypeComparison}, "FunctionTypeAssignment": {funcName: "FunctionTypeAssignment", native: divergence_hunt154.FunctionTypeAssignment}, "FunctionTypeReturn": {funcName: "FunctionTypeReturn", native: divergence_hunt154.FunctionTypeReturn}, "FunctionTypeParam": {funcName: "FunctionTypeParam", native: divergence_hunt154.FunctionTypeParam}, "FunctionTypeSlice": {funcName: "FunctionTypeSlice", native: divergence_hunt154.FunctionTypeSlice}, "FunctionTypeMap": {funcName: "FunctionTypeMap", native: divergence_hunt154.FunctionTypeMap}, "AnonymousStructNested": {funcName: "AnonymousStructNested", native: divergence_hunt154.AnonymousStructNested},
	}})
}
func TestDivergenceHunt155(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt155Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InterfaceSatisfaction": {funcName: "InterfaceSatisfaction", native: divergence_hunt155.InterfaceSatisfaction}, "EmptyInterface": {funcName: "EmptyInterface", native: divergence_hunt155.EmptyInterface}, "InterfaceNilComparison": {funcName: "InterfaceNilComparison", native: divergence_hunt155.InterfaceNilComparison}, "InterfaceTypeSwitch": {funcName: "InterfaceTypeSwitch", native: divergence_hunt155.InterfaceTypeSwitch}, "InterfaceTypeAssertion": {funcName: "InterfaceTypeAssertion", native: divergence_hunt155.InterfaceTypeAssertion}, "InterfaceSlice": {funcName: "InterfaceSlice", native: divergence_hunt155.InterfaceSlice}, "InterfaceMap": {funcName: "InterfaceMap", native: divergence_hunt155.InterfaceMap}, "InterfaceEmbedding": {funcName: "InterfaceEmbedding", native: divergence_hunt155.InterfaceEmbedding},
	}})
}
func TestDivergenceHunt156(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt156Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapZeroValue": {funcName: "MapZeroValue", native: divergence_hunt156.MapZeroValue}, "MapEmptyVsNil": {funcName: "MapEmptyVsNil", native: divergence_hunt156.MapEmptyVsNil}, "MapMakeWithCapacity": {funcName: "MapMakeWithCapacity", native: divergence_hunt156.MapMakeWithCapacity}, "MapStringKey": {funcName: "MapStringKey", native: divergence_hunt156.MapStringKey}, "MapIntKey": {funcName: "MapIntKey", native: divergence_hunt156.MapIntKey}, "MapStructKey": {funcName: "MapStructKey", native: divergence_hunt156.MapStructKey}, "MapArrayKey": {funcName: "MapArrayKey", native: divergence_hunt156.MapArrayKey}, "MapPointerKey": {funcName: "MapPointerKey", native: divergence_hunt156.MapPointerKey}, "MapInterfaceKey": {funcName: "MapInterfaceKey", native: divergence_hunt156.MapInterfaceKey}, "MapSliceValue": {funcName: "MapSliceValue", native: divergence_hunt156.MapSliceValue}, "MapMapValue": {funcName: "MapMapValue", native: divergence_hunt156.MapMapValue},
	}})
}
func TestDivergenceHunt157(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt157Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Utf8RuneCount": {funcName: "Utf8RuneCount", native: divergence_hunt157.Utf8RuneCount}, "Utf8DecodeRune": {funcName: "Utf8DecodeRune", native: divergence_hunt157.Utf8DecodeRune}, "Utf8Valid": {funcName: "Utf8Valid", native: divergence_hunt157.Utf8Valid}, "Utf8RuneLen": {funcName: "Utf8RuneLen", native: divergence_hunt157.Utf8RuneLen}, "StringBuilderBasic": {funcName: "StringBuilderBasic", native: divergence_hunt157.StringBuilderBasic}, "StringBuilderGrow": {funcName: "StringBuilderGrow", native: divergence_hunt157.StringBuilderGrow}, "StringBuilderByte": {funcName: "StringBuilderByte", native: divergence_hunt157.StringBuilderByte}, "StringBuilderRune": {funcName: "StringBuilderRune", native: divergence_hunt157.StringBuilderRune}, "StringCompare": {funcName: "StringCompare", native: divergence_hunt157.StringCompare}, "StringEqualFold": {funcName: "StringEqualFold", native: divergence_hunt157.StringEqualFold}, "StringIndexAny": {funcName: "StringIndexAny", native: divergence_hunt157.StringIndexAny}, "StringLastIndex": {funcName: "StringLastIndex", native: divergence_hunt157.StringLastIndex}, "StringCutPrefix": {funcName: "StringCutPrefix", native: divergence_hunt157.StringCutPrefix}, "StringCutSuffix": {funcName: "StringCutSuffix", native: divergence_hunt157.StringCutSuffix}, "StringClone": {funcName: "StringClone", native: divergence_hunt157.StringClone}, "StringCount": {funcName: "StringCount", native: divergence_hunt157.StringCount},
	}})
}
func TestDivergenceHunt158(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt158Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"IotaBasic": {funcName: "IotaBasic", native: divergence_hunt158.IotaBasic}, "IotaWithValue": {funcName: "IotaWithValue", native: divergence_hunt158.IotaWithValue}, "IotaWithSkip": {funcName: "IotaWithSkip", native: divergence_hunt158.IotaWithSkip}, "IotaBitShift": {funcName: "IotaBitShift", native: divergence_hunt158.IotaBitShift}, "IotaExpression": {funcName: "IotaExpression", native: divergence_hunt158.IotaExpression}, "IotaMultiplePerLine": {funcName: "IotaMultiplePerLine", native: divergence_hunt158.IotaMultiplePerLine}, "IotaReset": {funcName: "IotaReset", native: divergence_hunt158.IotaReset}, "IotaStringer": {funcName: "IotaStringer", native: divergence_hunt158.IotaStringer}, "UntypedConstant": {funcName: "UntypedConstant", native: divergence_hunt158.UntypedConstant}, "ConstWithType": {funcName: "ConstWithType", native: divergence_hunt158.ConstWithType},
	}})
}
func TestDivergenceHunt159(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt159Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NamedReturnWithDefer": {funcName: "NamedReturnWithDefer", native: divergence_hunt159.NamedReturnWithDefer}, "NamedReturnWithDeferAndValue": {funcName: "NamedReturnWithDeferAndValue", native: divergence_hunt159.NamedReturnWithDeferAndValue}, "MultipleDefersExecutionOrder": {funcName: "MultipleDefersExecutionOrder", native: divergence_hunt159.MultipleDefersExecutionOrder}, "DeferWithArguments": {funcName: "DeferWithArguments", native: divergence_hunt159.DeferWithArguments}, "DeferInLoopLastValue": {funcName: "DeferInLoopLastValue", native: divergence_hunt159.DeferInLoopLastValue}, "PanicWithString": {funcName: "PanicWithString", native: divergence_hunt159.PanicWithString}, "PanicWithError": {funcName: "PanicWithError", native: divergence_hunt159.PanicWithError}, "PanicWithInt": {funcName: "PanicWithInt", native: divergence_hunt159.PanicWithInt}, "NestedDeferPanic": {funcName: "NestedDeferPanic", native: divergence_hunt159.NestedDeferPanic}, "RecoverOnlyInDefer": {funcName: "RecoverOnlyInDefer", native: divergence_hunt159.RecoverOnlyInDefer}, "PanicNilInterface": {funcName: "PanicNilInterface", native: divergence_hunt159.PanicNilInterface}, "DeferReturnValue": {funcName: "DeferReturnValue", native: divergence_hunt159.DeferReturnValue},
	}})
}
func TestDivergenceHunt160(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt160Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TypeAliasConversion": {funcName: "TypeAliasConversion", native: divergence_hunt160.TypeAliasConversion}, "TypeAliasString": {funcName: "TypeAliasString", native: divergence_hunt160.TypeAliasString}, "ByteSliceToString": {funcName: "ByteSliceToString", native: divergence_hunt160.ByteSliceToString}, "StringToByteSlice": {funcName: "StringToByteSlice", native: divergence_hunt160.StringToByteSlice}, "RuneSliceToString": {funcName: "RuneSliceToString", native: divergence_hunt160.RuneSliceToString}, "IntToFloatConversion": {funcName: "IntToFloatConversion", native: divergence_hunt160.IntToFloatConversion}, "FloatToIntConversion": {funcName: "FloatToIntConversion", native: divergence_hunt160.FloatToIntConversion}, "UintToIntConversion": {funcName: "UintToIntConversion", native: divergence_hunt160.UintToIntConversion}, "IntToUintConversion": {funcName: "IntToUintConversion", native: divergence_hunt160.IntToUintConversion}, "LargeIntToUint8": {funcName: "LargeIntToUint8", native: divergence_hunt160.LargeIntToUint8}, "NegativeIntToUint": {funcName: "NegativeIntToUint", native: divergence_hunt160.NegativeIntToUint}, "FloatSpecialToInt": {funcName: "FloatSpecialToInt", native: divergence_hunt160.FloatSpecialToInt}, "ByteToIntConversion": {funcName: "ByteToIntConversion", native: divergence_hunt160.ByteToIntConversion}, "Int8ToInt16Conversion": {funcName: "Int8ToInt16Conversion", native: divergence_hunt160.Int8ToInt16Conversion}, "ComplexConversion": {funcName: "ComplexConversion", native: divergence_hunt160.ComplexConversion}, "Complex64ToComplex128": {funcName: "Complex64ToComplex128", native: divergence_hunt160.Complex64ToComplex128},
	}})
}

func TestDivergenceHunt161(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt161Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RecursiveFactorial": {funcName: "RecursiveFactorial", native: divergence_hunt161.RecursiveFactorial}, "RecursiveFibonacci": {funcName: "RecursiveFibonacci", native: divergence_hunt161.RecursiveFibonacci}, "RecursiveSum": {funcName: "RecursiveSum", native: divergence_hunt161.RecursiveSum}, "RecursiveMax": {funcName: "RecursiveMax", native: divergence_hunt161.RecursiveMax}, "RecursiveReverse": {funcName: "RecursiveReverse", native: divergence_hunt161.RecursiveReverse}, "RecursiveTreeTraversal": {funcName: "RecursiveTreeTraversal", native: divergence_hunt161.RecursiveTreeTraversal}, "RecursivePower": {funcName: "RecursivePower", native: divergence_hunt161.RecursivePower}, "RecursiveGCD": {funcName: "RecursiveGCD", native: divergence_hunt161.RecursiveGCD}, "RecursivePalindrome": {funcName: "RecursivePalindrome", native: divergence_hunt161.RecursivePalindrome}, "RecursiveBinarySearch": {funcName: "RecursiveBinarySearch", native: divergence_hunt161.RecursiveBinarySearch}, "RecursiveMergeSort": {funcName: "RecursiveMergeSort", native: divergence_hunt161.RecursiveMergeSort}, "TailRecursiveFactorial": {funcName: "TailRecursiveFactorial", native: divergence_hunt161.TailRecursiveFactorial},
	}})
}
func TestDivergenceHunt162(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt162Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NilValueReceiver": {funcName: "NilValueReceiver", native: divergence_hunt162.NilValueReceiver, knownIssue: true},
		"NilPointerReceiverSafe": {funcName: "NilPointerReceiverSafe", native: divergence_hunt162.NilPointerReceiverSafe},
		"NonNilValueReceiver": {funcName: "NonNilValueReceiver", native: divergence_hunt162.NonNilValueReceiver},
		"NonNilPointerReceiver": {funcName: "NonNilPointerReceiver", native: divergence_hunt162.NonNilPointerReceiver},
		"NilMixedReceiverValueMethod": {funcName: "NilMixedReceiverValueMethod", native: divergence_hunt162.NilMixedReceiverValueMethod, knownIssue: true},
		"NilMixedReceiverPointerMethod": {funcName: "NilMixedReceiverPointerMethod", native: divergence_hunt162.NilMixedReceiverPointerMethod},
		"InterfaceWithNilConcrete": {funcName: "InterfaceWithNilConcrete", native: divergence_hunt162.InterfaceWithNilConcrete, knownIssue: true},
		"NilInterfaceCall": {funcName: "NilInterfaceCall", native: divergence_hunt162.NilInterfaceCall, knownIssue: true},
		"NestedNilReceiver": {funcName: "NestedNilReceiver", native: divergence_hunt162.NestedNilReceiver, knownIssue: true},
		"SliceOfInterfacesWithNil": {funcName: "SliceOfInterfacesWithNil", native: divergence_hunt162.SliceOfInterfacesWithNil},
		"MapWithNilValues": {funcName: "MapWithNilValues", native: divergence_hunt162.MapWithNilValues, knownIssue: true},
	}})
}
func TestDivergenceHunt163(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt163Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MinInt": {funcName: "MinInt", native: divergence_hunt163.MinInt}, "MinFloat64": {funcName: "MinFloat64", native: divergence_hunt163.MinFloat64}, "MinString": {funcName: "MinString", native: divergence_hunt163.MinString}, "GenericStackPattern": {funcName: "GenericStackPattern", native: divergence_hunt163.GenericStackPattern}, "GenericStackPatternString": {funcName: "GenericStackPatternString", native: divergence_hunt163.GenericStackPatternString}, "GenericMapPattern": {funcName: "GenericMapPattern", native: divergence_hunt163.GenericMapPattern}, "GenericFilterPattern": {funcName: "GenericFilterPattern", native: divergence_hunt163.GenericFilterPattern}, "GenericReducePattern": {funcName: "GenericReducePattern", native: divergence_hunt163.GenericReducePattern}, "ComparableConstraintPattern": {funcName: "ComparableConstraintPattern", native: divergence_hunt163.ComparableConstraintPattern}, "OrderedConstraintPattern": {funcName: "OrderedConstraintPattern", native: divergence_hunt163.OrderedConstraintPattern}, "GenericCachePattern": {funcName: "GenericCachePattern", native: divergence_hunt163.GenericCachePattern},
	}})
}
func TestDivergenceHunt164(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt164Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ValueReceiverNoModify": {funcName: "ValueReceiverNoModify", native: divergence_hunt164.ValueReceiverNoModify},
		"PointerReceiverModifies": {funcName: "PointerReceiverModifies", native: divergence_hunt164.PointerReceiverModifies},
		"PointerReceiverViaValue": {funcName: "PointerReceiverViaValue", native: divergence_hunt164.PointerReceiverViaValue},
		"ValueReceiverOnPointer": {funcName: "ValueReceiverOnPointer", native: divergence_hunt164.ValueReceiverOnPointer},
		"MethodSetDifference": {funcName: "MethodSetDifference", native: divergence_hunt164.MethodSetDifference},
		"InterfaceSatisfactionValue": {funcName: "InterfaceSatisfactionValue", native: divergence_hunt164.InterfaceSatisfactionValue},
		"InterfaceSatisfactionPointer": {funcName: "InterfaceSatisfactionPointer", native: divergence_hunt164.InterfaceSatisfactionPointer},
		"NilPointerReceiverWithValueMethod": {funcName: "NilPointerReceiverWithValueMethod", native: divergence_hunt164.NilPointerReceiverWithValueMethod, knownIssue: true},
		"NilPointerReceiverWithPointerMethod": {funcName: "NilPointerReceiverWithPointerMethod", native: divergence_hunt164.NilPointerReceiverWithPointerMethod, knownIssue: true},
		"AssignmentToInterface": {funcName: "AssignmentToInterface", native: divergence_hunt164.AssignmentToInterface},
		"SliceOfStructsWithMethods": {funcName: "SliceOfStructsWithMethods", native: divergence_hunt164.SliceOfStructsWithMethods},
		"MapOfStructsWithMethods": {funcName: "MapOfStructsWithMethods", native: divergence_hunt164.MapOfStructsWithMethods},
		"EmbeddedFieldMethodPromotion": {funcName: "EmbeddedFieldMethodPromotion", native: divergence_hunt164.EmbeddedFieldMethodPromotion},
	}})
}
func TestDivergenceHunt165(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt165Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InterfaceEmbedding": {funcName: "InterfaceEmbedding", native: divergence_hunt165.InterfaceEmbedding}, "StructEmbeddingBasic": {funcName: "StructEmbeddingBasic", native: divergence_hunt165.StructEmbeddingBasic}, "StructEmbeddingNested": {funcName: "StructEmbeddingNested", native: divergence_hunt165.StructEmbeddingNested}, "EmbeddedFieldAccess": {funcName: "EmbeddedFieldAccess", native: divergence_hunt165.EmbeddedFieldAccess}, "EmbeddedPointerFieldAccess": {funcName: "EmbeddedPointerFieldAccess", native: divergence_hunt165.EmbeddedPointerFieldAccess}, "EmbeddedNilPointer": {funcName: "EmbeddedNilPointer", native: divergence_hunt165.EmbeddedNilPointer}, "MultipleEmbeddedStructs": {funcName: "MultipleEmbeddedStructs", native: divergence_hunt165.MultipleEmbeddedStructs}, "MethodShadowing": {funcName: "MethodShadowing", native: divergence_hunt165.MethodShadowing}, "EmbeddedInterfaceSatisfaction": {funcName: "EmbeddedInterfaceSatisfaction", native: divergence_hunt165.EmbeddedInterfaceSatisfaction}, "DeepEmbedding": {funcName: "DeepEmbedding", native: divergence_hunt165.DeepEmbedding}, "EmbeddedWithTags": {funcName: "EmbeddedWithTags", native: divergence_hunt165.EmbeddedWithTags},
	}})
}
func TestDivergenceHunt166(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt166Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicGoroutineLaunch": {funcName: "BasicGoroutineLaunch", native: divergence_hunt166.BasicGoroutineLaunch}, "MultipleGoroutines": {funcName: "MultipleGoroutines", native: divergence_hunt166.MultipleGoroutines}, "WaitGroupBasic": {funcName: "WaitGroupBasic", native: divergence_hunt166.WaitGroupBasic}, "WaitGroupWithData": {funcName: "WaitGroupWithData", native: divergence_hunt166.WaitGroupWithData}, "ChannelSynchronization": {funcName: "ChannelSynchronization", native: divergence_hunt166.ChannelSynchronization}, "BufferedChannelGoroutines": {funcName: "BufferedChannelGoroutines", native: divergence_hunt166.BufferedChannelGoroutines}, "GoroutineClosureCapture": {funcName: "GoroutineClosureCapture", native: divergence_hunt166.GoroutineClosureCapture}, "SelectStatement": {funcName: "SelectStatement", native: divergence_hunt166.SelectStatement}, "FanOutPattern": {funcName: "FanOutPattern", native: divergence_hunt166.FanOutPattern}, "FanInPattern": {funcName: "FanInPattern", native: divergence_hunt166.FanInPattern}, "GoroutineWithError": {funcName: "GoroutineWithError", native: divergence_hunt166.GoroutineWithError}, "TimeoutPattern": {funcName: "TimeoutPattern", native: divergence_hunt166.TimeoutPattern},
	}})
}
func TestDivergenceHunt167(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt167Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicMutex": {funcName: "BasicMutex", native: divergence_hunt167.BasicMutex}, "RWMutexBasic": {funcName: "RWMutexBasic", native: divergence_hunt167.RWMutexBasic}, "DeferredUnlock": {funcName: "DeferredUnlock", native: divergence_hunt167.DeferredUnlock}, "MutexOperations": {funcName: "MutexOperations", native: divergence_hunt167.MutexOperations}, "MutexLoadStore": {funcName: "MutexLoadStore", native: divergence_hunt167.MutexLoadStore}, "OncePattern": {funcName: "OncePattern", native: divergence_hunt167.OncePattern}, "PoolPattern": {funcName: "PoolPattern", native: divergence_hunt167.PoolPattern}, "CondPattern": {funcName: "CondPattern", native: divergence_hunt167.CondPattern}, "MapPattern": {funcName: "MapPattern", native: divergence_hunt167.MapPattern}, "ProtectedStruct": {funcName: "ProtectedStruct", native: divergence_hunt167.ProtectedStruct}, "TryLockPattern": {funcName: "TryLockPattern", native: divergence_hunt167.TryLockPattern},
	}})
}
func TestDivergenceHunt168(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt168Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BackgroundContext": {funcName: "BackgroundContext", native: divergence_hunt168.BackgroundContext}, "TODOContext": {funcName: "TODOContext", native: divergence_hunt168.TODOContext}, "WithValue": {funcName: "WithValue", native: divergence_hunt168.WithValue}, "WithCancel": {funcName: "WithCancel", native: divergence_hunt168.WithCancel}, "WithDeadline": {funcName: "WithDeadline", native: divergence_hunt168.WithDeadline}, "WithTimeout": {funcName: "WithTimeout", native: divergence_hunt168.WithTimeout}, "NestedContext": {funcName: "NestedContext", native: divergence_hunt168.NestedContext}, "ContextValueOverride": {funcName: "ContextValueOverride", native: divergence_hunt168.ContextValueOverride}, "ContextPropagation": {funcName: "ContextPropagation", native: divergence_hunt168.ContextPropagation}, "SelectWithContext": {funcName: "SelectWithContext", native: divergence_hunt168.SelectWithContext}, "ContextInStruct": {funcName: "ContextInStruct", native: divergence_hunt168.ContextInStruct}, "MultipleValues": {funcName: "MultipleValues", native: divergence_hunt168.MultipleValues},
	}})
}
func TestDivergenceHunt169(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt169Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"URLParsing": {funcName: "URLParsing", native: divergence_hunt169.URLParsing}, "URLQuery": {funcName: "URLQuery", native: divergence_hunt169.URLQuery}, "URLString": {funcName: "URLString", native: divergence_hunt169.URLString}, "URLEscaping": {funcName: "URLEscaping", native: divergence_hunt169.URLEscaping}, "URLPathEscaping": {funcName: "URLPathEscaping", native: divergence_hunt169.URLPathEscaping}, "URLUserInfo": {funcName: "URLUserInfo", native: divergence_hunt169.URLUserInfo}, "URLIsAbs": {funcName: "URLIsAbs", native: divergence_hunt169.URLIsAbs}, "URLHostname": {funcName: "URLHostname", native: divergence_hunt169.URLHostname}, "URLRequestURI": {funcName: "URLRequestURI", native: divergence_hunt169.URLRequestURI}, "URLResolveReference": {funcName: "URLResolveReference", native: divergence_hunt169.URLResolveReference}, "URLValues": {funcName: "URLValues", native: divergence_hunt169.URLValues}, "URLEncode": {funcName: "URLEncode", native: divergence_hunt169.URLEncode},
	}})
}
func TestDivergenceHunt170(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt170Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicMarshalUnmarshal": {funcName: "BasicMarshalUnmarshal", native: divergence_hunt170.BasicMarshalUnmarshal}, "SliceMarshalUnmarshal": {funcName: "SliceMarshalUnmarshal", native: divergence_hunt170.SliceMarshalUnmarshal}, "MapMarshalUnmarshal": {funcName: "MapMarshalUnmarshal", native: divergence_hunt170.MapMarshalUnmarshal}, "NestedStructMarshal": {funcName: "NestedStructMarshal", native: divergence_hunt170.NestedStructMarshal}, "PointerFieldMarshal": {funcName: "PointerFieldMarshal", native: divergence_hunt170.PointerFieldMarshal}, "OmitEmptyTag": {funcName: "OmitEmptyTag", native: divergence_hunt170.OmitEmptyTag}, "StringTag": {funcName: "StringTag", native: divergence_hunt170.StringTag}, "IgnoreField": {funcName: "IgnoreField", native: divergence_hunt170.IgnoreField}, "UnmarshalUnknownFields": {funcName: "UnmarshalUnknownFields", native: divergence_hunt170.UnmarshalUnknownFields}, "UnmarshalTypeMismatch": {funcName: "UnmarshalTypeMismatch", native: divergence_hunt170.UnmarshalTypeMismatch}, "RawMessage": {funcName: "RawMessage", native: divergence_hunt170.RawMessage}, "NumberHandling": {funcName: "NumberHandling", native: divergence_hunt170.NumberHandling}, "CustomMarshalJSON": {funcName: "CustomMarshalJSON", native: divergence_hunt170.CustomMarshalJSON}, "MarshalEmptySlice": {funcName: "MarshalEmptySlice", native: divergence_hunt170.MarshalEmptySlice}, "MarshalNilSlice": {funcName: "MarshalNilSlice", native: divergence_hunt170.MarshalNilSlice}, "EscapeHTML": {funcName: "EscapeHTML", native: divergence_hunt170.EscapeHTML}, "IndentJSON": {funcName: "IndentJSON", native: divergence_hunt170.IndentJSON}, "DecodeStream": {funcName: "DecodeStream", native: divergence_hunt170.DecodeStream}, "LargeNumber": {funcName: "LargeNumber", native: divergence_hunt170.LargeNumber}, "BooleanStringUnmarshal": {funcName: "BooleanStringUnmarshal", native: divergence_hunt170.BooleanStringUnmarshal}, "FloatPrecision": {funcName: "FloatPrecision", native: divergence_hunt170.FloatPrecision}, "MarshalInterface": {funcName: "MarshalInterface", native: divergence_hunt170.MarshalInterface}, "UnmarshalToInterface": {funcName: "UnmarshalToInterface", native: divergence_hunt170.UnmarshalToInterface}, "ValidJSON": {funcName: "ValidJSON", native: divergence_hunt170.ValidJSON}, "CompactJSON": {funcName: "CompactJSON", native: divergence_hunt170.CompactJSON}, "HTMLEscape": {funcName: "HTMLEscape", native: divergence_hunt170.HTMLEscape}, "IndentTests": {funcName: "IndentTests", native: divergence_hunt170.IndentTests}, "MarshalIntFloatString": {funcName: "MarshalIntFloatString", native: divergence_hunt170.MarshalIntFloatString}, "TagNameCase": {funcName: "TagNameCase", native: divergence_hunt170.TagNameCase}, "AnonymousStructMarshal": {funcName: "AnonymousStructMarshal", native: divergence_hunt170.AnonymousStructMarshal}, "NullValueMarshal": {funcName: "NullValueMarshal", native: divergence_hunt170.NullValueMarshal}, "ArrayMarshalUnmarshal": {funcName: "ArrayMarshalUnmarshal", native: divergence_hunt170.ArrayMarshalUnmarshal}, "ByteSliceMarshal": {funcName: "ByteSliceMarshal", native: divergence_hunt170.ByteSliceMarshal}, "StringPointerUnmarshal": {funcName: "StringPointerUnmarshal", native: divergence_hunt170.StringPointerUnmarshal}, "IntPointerUnmarshal": {funcName: "IntPointerUnmarshal", native: divergence_hunt170.IntPointerUnmarshal}, "ZeroValueMarshal": {funcName: "ZeroValueMarshal", native: divergence_hunt170.ZeroValueMarshal}, "EmbeddedStructMarshal": {funcName: "EmbeddedStructMarshal", native: divergence_hunt170.EmbeddedStructMarshal}, "MapStringInterfaceMarshal": {funcName: "MapStringInterfaceMarshal", native: divergence_hunt170.MapStringInterfaceMarshal}, "StructWithJSONTagDash": {funcName: "StructWithJSONTagDash", native: divergence_hunt170.StructWithJSONTagDash}, "ParseStringToFloat": {funcName: "ParseStringToFloat", native: divergence_hunt170.ParseStringToFloat}, "ParseIntTests": {funcName: "ParseIntTests", native: divergence_hunt170.ParseIntTests}, "ItoaTests": {funcName: "ItoaTests", native: divergence_hunt170.ItoaTests}, "FormatFloatTests": {funcName: "FormatFloatTests", native: divergence_hunt170.FormatFloatTests}, "QuoteString": {funcName: "QuoteString", native: divergence_hunt170.QuoteString}, "UnquoteString": {funcName: "UnquoteString", native: divergence_hunt170.UnquoteString}, "ParseBoolTests": {funcName: "ParseBoolTests", native: divergence_hunt170.ParseBoolTests}, "AppendIntTests": {funcName: "AppendIntTests", native: divergence_hunt170.AppendIntTests}, "AppendFloatTests": {funcName: "AppendFloatTests", native: divergence_hunt170.AppendFloatTests}, "AppendBoolTests": {funcName: "AppendBoolTests", native: divergence_hunt170.AppendBoolTests}, "AppendQuoteTests": {funcName: "AppendQuoteTests", native: divergence_hunt170.AppendQuoteTests}, "IsPrintTests": {funcName: "IsPrintTests", native: divergence_hunt170.IsPrintTests}, "CanBackquoteTests": {funcName: "CanBackquoteTests", native: divergence_hunt170.CanBackquoteTests},
	}})
}
func TestDivergenceHunt171(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt171Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TimeNow": {funcName: "TimeNow", native: divergence_hunt171.TimeNow}, "TimeAdd": {funcName: "TimeAdd", native: divergence_hunt171.TimeAdd}, "TimeSub": {funcName: "TimeSub", native: divergence_hunt171.TimeSub}, "TimeBeforeAfter": {funcName: "TimeBeforeAfter", native: divergence_hunt171.TimeBeforeAfter}, "TimeUnix": {funcName: "TimeUnix", native: divergence_hunt171.TimeUnix}, "TimeYearMonthDay": {funcName: "TimeYearMonthDay", native: divergence_hunt171.TimeYearMonthDay}, "TimeHourMinuteSecond": {funcName: "TimeHourMinuteSecond", native: divergence_hunt171.TimeHourMinuteSecond}, "TimeWeekday": {funcName: "TimeWeekday", native: divergence_hunt171.TimeWeekday}, "TimeTruncate": {funcName: "TimeTruncate", native: divergence_hunt171.TimeTruncate}, "DurationString": {funcName: "DurationString", native: divergence_hunt171.DurationString}, "DurationHours": {funcName: "DurationHours", native: divergence_hunt171.DurationHours}, "DurationMinutes": {funcName: "DurationMinutes", native: divergence_hunt171.DurationMinutes},
	}})
}
func TestDivergenceHunt172(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt172Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SortByLength": {funcName: "SortByLength", native: divergence_hunt172.SortByLength, knownIssue: true},
		"SortStructByField": {funcName: "SortStructByField", native: divergence_hunt172.SortStructByField, knownIssue: true},
		"SortReverse": {funcName: "SortReverse", native: divergence_hunt172.SortReverse, knownIssue: true},
		"SortInts": {funcName: "SortInts", native: divergence_hunt172.SortInts},
		"SortStrings": {funcName: "SortStrings", native: divergence_hunt172.SortStrings},
		"SortFloats": {funcName: "SortFloats", native: divergence_hunt172.SortFloats},
		"SortSearchInts": {funcName: "SortSearchInts", native: divergence_hunt172.SortSearchInts},
		"SortIsSorted": {funcName: "SortIsSorted", native: divergence_hunt172.SortIsSorted},
		"SortSlice": {funcName: "SortSlice", native: divergence_hunt172.SortSlice},
		"SortStable": {funcName: "SortStable", native: divergence_hunt172.SortStable},
	}})
}
func TestDivergenceHunt173(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt173Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"HeapInit": {funcName: "HeapInit", native: divergence_hunt173.HeapInit, knownIssue: true},
		"HeapPush": {funcName: "HeapPush", native: divergence_hunt173.HeapPush, knownIssue: true},
		"HeapPop": {funcName: "HeapPop", native: divergence_hunt173.HeapPop, knownIssue: true},
		"HeapRemove": {funcName: "HeapRemove", native: divergence_hunt173.HeapRemove, knownIssue: true},
		"HeapFix": {funcName: "HeapFix", native: divergence_hunt173.HeapFix, knownIssue: true},
		"HeapMultiplePushPop": {funcName: "HeapMultiplePushPop", native: divergence_hunt173.HeapMultiplePushPop, knownIssue: true},
		"MaxHeapOperations": {funcName: "MaxHeapOperations", native: divergence_hunt173.MaxHeapOperations, knownIssue: true},
		"PriorityQueueOperations": {funcName: "PriorityQueueOperations", native: divergence_hunt173.PriorityQueueOperations, knownIssue: true},
		"HeapEmpty": {funcName: "HeapEmpty", native: divergence_hunt173.HeapEmpty, knownIssue: true},
		"HeapSingleElement": {funcName: "HeapSingleElement", native: divergence_hunt173.HeapSingleElement, knownIssue: true},
	}})
}
func TestDivergenceHunt174(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt174Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ListNew": {funcName: "ListNew", native: divergence_hunt174.ListNew}, "ListPushBack": {funcName: "ListPushBack", native: divergence_hunt174.ListPushBack}, "ListPushFront": {funcName: "ListPushFront", native: divergence_hunt174.ListPushFront}, "ListFrontBack": {funcName: "ListFrontBack", native: divergence_hunt174.ListFrontBack}, "ListIterateForward": {funcName: "ListIterateForward", native: divergence_hunt174.ListIterateForward}, "ListIterateBackward": {funcName: "ListIterateBackward", native: divergence_hunt174.ListIterateBackward}, "ListInsertAfter": {funcName: "ListInsertAfter", native: divergence_hunt174.ListInsertAfter}, "ListInsertBefore": {funcName: "ListInsertBefore", native: divergence_hunt174.ListInsertBefore}, "ListRemove": {funcName: "ListRemove", native: divergence_hunt174.ListRemove}, "ListMoveToFront": {funcName: "ListMoveToFront", native: divergence_hunt174.ListMoveToFront}, "ListMoveToBack": {funcName: "ListMoveToBack", native: divergence_hunt174.ListMoveToBack}, "ListMixedOperations": {funcName: "ListMixedOperations", native: divergence_hunt174.ListMixedOperations},
	}})
}
func TestDivergenceHunt175(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt175Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RingNew": {funcName: "RingNew", native: divergence_hunt175.RingNew}, "RingInit": {funcName: "RingInit", native: divergence_hunt175.RingInit}, "RingNext": {funcName: "RingNext", native: divergence_hunt175.RingNext}, "RingPrev": {funcName: "RingPrev", native: divergence_hunt175.RingPrev}, "RingMove": {funcName: "RingMove", native: divergence_hunt175.RingMove}, "RingMoveNegative": {funcName: "RingMoveNegative", native: divergence_hunt175.RingMoveNegative}, "RingLink": {funcName: "RingLink", native: divergence_hunt175.RingLink}, "RingUnlink": {funcName: "RingUnlink", native: divergence_hunt175.RingUnlink}, "RingDo": {funcName: "RingDo", native: divergence_hunt175.RingDo}, "RingDoEmpty": {funcName: "RingDoEmpty", native: divergence_hunt175.RingDoEmpty}, "RingSingleElement": {funcName: "RingSingleElement", native: divergence_hunt175.RingSingleElement}, "RingCircular": {funcName: "RingCircular", native: divergence_hunt175.RingCircular},
	}})
}
func TestDivergenceHunt176(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt176Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BinarySearchIntsFound": {funcName: "BinarySearchIntsFound", native: divergence_hunt176.BinarySearchIntsFound}, "BinarySearchIntsNotFound": {funcName: "BinarySearchIntsNotFound", native: divergence_hunt176.BinarySearchIntsNotFound}, "BinarySearchIntsBefore": {funcName: "BinarySearchIntsBefore", native: divergence_hunt176.BinarySearchIntsBefore}, "BinarySearchIntsAfter": {funcName: "BinarySearchIntsAfter", native: divergence_hunt176.BinarySearchIntsAfter}, "BinarySearchStrings": {funcName: "BinarySearchStrings", native: divergence_hunt176.BinarySearchStrings}, "BinarySearchFloats": {funcName: "BinarySearchFloats", native: divergence_hunt176.BinarySearchFloats}, "BinarySearchIntsFunc": {funcName: "BinarySearchIntsFunc", native: divergence_hunt176.BinarySearchIntsFunc}, "BinarySearchStringsFunc": {funcName: "BinarySearchStringsFunc", native: divergence_hunt176.BinarySearchStringsFunc}, "BinarySearchEmpty": {funcName: "BinarySearchEmpty", native: divergence_hunt176.BinarySearchEmpty}, "BinarySearchSingleElement": {funcName: "BinarySearchSingleElement", native: divergence_hunt176.BinarySearchSingleElement}, "BinarySearchDuplicates": {funcName: "BinarySearchDuplicates", native: divergence_hunt176.BinarySearchDuplicates}, "BinarySearchFindRange": {funcName: "BinarySearchFindRange", native: divergence_hunt176.BinarySearchFindRange},
	}})
}
func TestDivergenceHunt177(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt177Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BytesCompare": {funcName: "BytesCompare", native: divergence_hunt177.BytesCompare}, "BytesEqual": {funcName: "BytesEqual", native: divergence_hunt177.BytesEqual}, "BytesContains": {funcName: "BytesContains", native: divergence_hunt177.BytesContains}, "BytesIndex": {funcName: "BytesIndex", native: divergence_hunt177.BytesIndex}, "BytesLastIndex": {funcName: "BytesLastIndex", native: divergence_hunt177.BytesLastIndex}, "BytesCount": {funcName: "BytesCount", native: divergence_hunt177.BytesCount}, "BytesReplace": {funcName: "BytesReplace", native: divergence_hunt177.BytesReplace}, "BytesReplaceAll": {funcName: "BytesReplaceAll", native: divergence_hunt177.BytesReplaceAll}, "BytesRepeat": {funcName: "BytesRepeat", native: divergence_hunt177.BytesRepeat}, "BytesToUpper": {funcName: "BytesToUpper", native: divergence_hunt177.BytesToUpper}, "BytesToLower": {funcName: "BytesToLower", native: divergence_hunt177.BytesToLower}, "BytesTrim": {funcName: "BytesTrim", native: divergence_hunt177.BytesTrim}, "BytesTrimSpace": {funcName: "BytesTrimSpace", native: divergence_hunt177.BytesTrimSpace}, "BytesTrimPrefix": {funcName: "BytesTrimPrefix", native: divergence_hunt177.BytesTrimPrefix}, "BytesTrimSuffix": {funcName: "BytesTrimSuffix", native: divergence_hunt177.BytesTrimSuffix}, "BytesSplit": {funcName: "BytesSplit", native: divergence_hunt177.BytesSplit}, "BytesJoin": {funcName: "BytesJoin", native: divergence_hunt177.BytesJoin}, "BytesHasPrefix": {funcName: "BytesHasPrefix", native: divergence_hunt177.BytesHasPrefix}, "BytesHasSuffix": {funcName: "BytesHasSuffix", native: divergence_hunt177.BytesHasSuffix}, "BytesFields": {funcName: "BytesFields", native: divergence_hunt177.BytesFields},
	}})
}
func TestDivergenceHunt178(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt178Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BufferNew": {funcName: "BufferNew", native: divergence_hunt178.BufferNew}, "BufferNewString": {funcName: "BufferNewString", native: divergence_hunt178.BufferNewString}, "BufferWriteString": {funcName: "BufferWriteString", native: divergence_hunt178.BufferWriteString}, "BufferWriteByte": {funcName: "BufferWriteByte", native: divergence_hunt178.BufferWriteByte}, "BufferWriteRune": {funcName: "BufferWriteRune", native: divergence_hunt178.BufferWriteRune}, "BufferWrite": {funcName: "BufferWrite", native: divergence_hunt178.BufferWrite}, "BufferLen": {funcName: "BufferLen", native: divergence_hunt178.BufferLen}, "BufferCap": {funcName: "BufferCap", native: divergence_hunt178.BufferCap}, "BufferReset": {funcName: "BufferReset", native: divergence_hunt178.BufferReset}, "BufferBytes": {funcName: "BufferBytes", native: divergence_hunt178.BufferBytes}, "BufferGrow": {funcName: "BufferGrow", native: divergence_hunt178.BufferGrow}, "BufferRead": {funcName: "BufferRead", native: divergence_hunt178.BufferRead}, "BufferReadString": {funcName: "BufferReadString", native: divergence_hunt178.BufferReadString}, "BufferReadByte": {funcName: "BufferReadByte", native: divergence_hunt178.BufferReadByte}, "BufferNext": {funcName: "BufferNext", native: divergence_hunt178.BufferNext}, "BufferUnreadByte": {funcName: "BufferUnreadByte", native: divergence_hunt178.BufferUnreadByte}, "BufferTruncate": {funcName: "BufferTruncate", native: divergence_hunt178.BufferTruncate}, "BufferAvailable": {funcName: "BufferAvailable", native: divergence_hunt178.BufferAvailable}, "BufferAvailableBuffer": {funcName: "BufferAvailableBuffer", native: divergence_hunt178.BufferAvailableBuffer}, "BufferStringEmpty": {funcName: "BufferStringEmpty", native: divergence_hunt178.BufferStringEmpty},
	}})
}
func TestDivergenceHunt179(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt179Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ReaderFromString": {funcName: "ReaderFromString", native: divergence_hunt179.ReaderFromString}, "ReaderReadAt": {funcName: "ReaderReadAt", native: divergence_hunt179.ReaderReadAt}, "ReaderSeek": {funcName: "ReaderSeek", native: divergence_hunt179.ReaderSeek}, "ReaderSize": {funcName: "ReaderSize", native: divergence_hunt179.ReaderSize}, "ReaderLen": {funcName: "ReaderLen", native: divergence_hunt179.ReaderLen}, "WriterToBuffer": {funcName: "WriterToBuffer", native: divergence_hunt179.WriterToBuffer}, "WriterWriteByte": {funcName: "WriterWriteByte", native: divergence_hunt179.WriterWriteByte}, "WriterWriteRune": {funcName: "WriterWriteRune", native: divergence_hunt179.WriterWriteRune}, "ReaderWriterCombined": {funcName: "ReaderWriterCombined", native: divergence_hunt179.ReaderWriterCombined}, "ReadFull": {funcName: "ReadFull", native: divergence_hunt179.ReadFull}, "WriteString": {funcName: "WriteString", native: divergence_hunt179.WriteString}, "CopyReaderToWriter": {funcName: "CopyReaderToWriter", native: divergence_hunt179.CopyReaderToWriter, knownIssue: true}, "MultiWrite": {funcName: "MultiWrite", native: divergence_hunt179.MultiWrite}, "LimitedRead": {funcName: "LimitedRead", native: divergence_hunt179.LimitedRead}, "DiscardRead": {funcName: "DiscardRead", native: divergence_hunt179.DiscardRead}, "PipeSimulation": {funcName: "PipeSimulation", native: divergence_hunt179.PipeSimulation}, "TeeReaderSimulation": {funcName: "TeeReaderSimulation", native: divergence_hunt179.TeeReaderSimulation, knownIssue: true}, "SectionReaderSimulation": {funcName: "SectionReaderSimulation", native: divergence_hunt179.SectionReaderSimulation},
	}})
}
func TestDivergenceHunt180(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt180Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"FmtVerbInt": {funcName: "FmtVerbInt", native: divergence_hunt180.FmtVerbInt}, "FmtVerbFloat": {funcName: "FmtVerbFloat", native: divergence_hunt180.FmtVerbFloat}, "FmtVerbString": {funcName: "FmtVerbString", native: divergence_hunt180.FmtVerbString}, "FmtVerbBool": {funcName: "FmtVerbBool", native: divergence_hunt180.FmtVerbBool}, "FmtVerbChar": {funcName: "FmtVerbChar", native: divergence_hunt180.FmtVerbChar}, "FmtVerbBase": {funcName: "FmtVerbBase", native: divergence_hunt180.FmtVerbBase}, "FmtVerbWidth": {funcName: "FmtVerbWidth", native: divergence_hunt180.FmtVerbWidth}, "FmtVerbPrecision": {funcName: "FmtVerbPrecision", native: divergence_hunt180.FmtVerbPrecision}, "FmtVerbPointer": {funcName: "FmtVerbPointer", native: divergence_hunt180.FmtVerbPointer}, "FmtVerbType": {funcName: "FmtVerbType", native: divergence_hunt180.FmtVerbType}, "FmtVerbSlice": {funcName: "FmtVerbSlice", native: divergence_hunt180.FmtVerbSlice}, "FmtVerbStruct": {funcName: "FmtVerbStruct", native: divergence_hunt180.FmtVerbStruct}, "FmtSprintVsSprintln": {funcName: "FmtSprintVsSprintln", native: divergence_hunt180.FmtSprintVsSprintln}, "FmtFprintfSimulation": {funcName: "FmtFprintfSimulation", native: divergence_hunt180.FmtFprintfSimulation}, "FmtErrorf": {funcName: "FmtErrorf", native: divergence_hunt180.FmtErrorf}, "FmtComplex": {funcName: "FmtComplex", native: divergence_hunt180.FmtComplex}, "FmtUint": {funcName: "FmtUint", native: divergence_hunt180.FmtUint}, "FmtHexPointer": {funcName: "FmtHexPointer", native: divergence_hunt180.FmtHexPointer}, "FmtPadding": {funcName: "FmtPadding", native: divergence_hunt180.FmtPadding}, "FmtEmptyInterface": {funcName: "FmtEmptyInterface", native: divergence_hunt180.FmtEmptyInterface},
	}})
}
func TestDivergenceHunt181(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt181Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"AtoiBasic": {funcName: "AtoiBasic", native: divergence_hunt181.AtoiBasic}, "AtoiNegative": {funcName: "AtoiNegative", native: divergence_hunt181.AtoiNegative}, "AtoiInvalid": {funcName: "AtoiInvalid", native: divergence_hunt181.AtoiInvalid}, "ParseIntBase10": {funcName: "ParseIntBase10", native: divergence_hunt181.ParseIntBase10}, "ParseIntBase16": {funcName: "ParseIntBase16", native: divergence_hunt181.ParseIntBase16}, "ParseIntBase2": {funcName: "ParseIntBase2", native: divergence_hunt181.ParseIntBase2}, "FormatIntBase10": {funcName: "FormatIntBase10", native: divergence_hunt181.FormatIntBase10}, "FormatIntBase16": {funcName: "FormatIntBase16", native: divergence_hunt181.FormatIntBase16}, "FormatIntBase2": {funcName: "FormatIntBase2", native: divergence_hunt181.FormatIntBase2}, "ParseBool": {funcName: "ParseBool", native: divergence_hunt181.ParseBool}, "ItoaBasic": {funcName: "ItoaBasic", native: divergence_hunt181.ItoaBasic}, "ItoaNegative": {funcName: "ItoaNegative", native: divergence_hunt181.ItoaNegative},
	}})
}
func TestDivergenceHunt182(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt182Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BuilderBasic": {funcName: "BuilderBasic", native: divergence_hunt182.BuilderBasic}, "BuilderWriteByte": {funcName: "BuilderWriteByte", native: divergence_hunt182.BuilderWriteByte}, "BuilderWriteRune": {funcName: "BuilderWriteRune", native: divergence_hunt182.BuilderWriteRune}, "BuilderGrow": {funcName: "BuilderGrow", native: divergence_hunt182.BuilderGrow}, "BuilderReset": {funcName: "BuilderReset", native: divergence_hunt182.BuilderReset}, "BuilderMultipleWrites": {funcName: "BuilderMultipleWrites", native: divergence_hunt182.BuilderMultipleWrites}, "BuilderMixedWrites": {funcName: "BuilderMixedWrites", native: divergence_hunt182.BuilderMixedWrites}, "BuilderEmpty": {funcName: "BuilderEmpty", native: divergence_hunt182.BuilderEmpty}, "BuilderNested": {funcName: "BuilderNested", native: divergence_hunt182.BuilderNested}, "BuilderLargeString": {funcName: "BuilderLargeString", native: divergence_hunt182.BuilderLargeString},
	}})
}
func TestDivergenceHunt183(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt183Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"IsDigit": {funcName: "IsDigit", native: divergence_hunt183.IsDigit}, "IsLetter": {funcName: "IsLetter", native: divergence_hunt183.IsLetter}, "IsSpace": {funcName: "IsSpace", native: divergence_hunt183.IsSpace}, "IsUpper": {funcName: "IsUpper", native: divergence_hunt183.IsUpper}, "IsLower": {funcName: "IsLower", native: divergence_hunt183.IsLower}, "IsNumber": {funcName: "IsNumber", native: divergence_hunt183.IsNumber}, "IsPunct": {funcName: "IsPunct", native: divergence_hunt183.IsPunct}, "ToUpper": {funcName: "ToUpper", native: divergence_hunt183.ToUpper}, "ToLower": {funcName: "ToLower", native: divergence_hunt183.ToLower}, "SimpleFold": {funcName: "SimpleFold", native: divergence_hunt183.SimpleFold}, "IsGraphic": {funcName: "IsGraphic", native: divergence_hunt183.IsGraphic}, "IsControl": {funcName: "IsControl", native: divergence_hunt183.IsControl},
	}})
}
func TestDivergenceHunt184(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt184Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MatchString": {funcName: "MatchString", native: divergence_hunt184.MatchString}, "MatchStringDigit": {funcName: "MatchStringDigit", native: divergence_hunt184.MatchStringDigit}, "MatchStringNotMatch": {funcName: "MatchStringNotMatch", native: divergence_hunt184.MatchStringNotMatch}, "MatchStringWord": {funcName: "MatchStringWord", native: divergence_hunt184.MatchStringWord}, "CompileAndMatch": {funcName: "CompileAndMatch", native: divergence_hunt184.CompileAndMatch}, "FindString": {funcName: "FindString", native: divergence_hunt184.FindString}, "FindAllString": {funcName: "FindAllString", native: divergence_hunt184.FindAllString}, "ReplaceAllString": {funcName: "ReplaceAllString", native: divergence_hunt184.ReplaceAllString}, "SplitString": {funcName: "SplitString", native: divergence_hunt184.SplitString}, "QuoteMeta": {funcName: "QuoteMeta", native: divergence_hunt184.QuoteMeta}, "FindStringIndex": {funcName: "FindStringIndex", native: divergence_hunt184.FindStringIndex},
	}})
}
func TestDivergenceHunt185(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt185Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NewInt": {funcName: "NewInt", native: divergence_hunt185.NewInt}, "IntAdd": {funcName: "IntAdd", native: divergence_hunt185.IntAdd}, "IntSub": {funcName: "IntSub", native: divergence_hunt185.IntSub}, "IntMul": {funcName: "IntMul", native: divergence_hunt185.IntMul}, "IntDiv": {funcName: "IntDiv", native: divergence_hunt185.IntDiv}, "IntMod": {funcName: "IntMod", native: divergence_hunt185.IntMod}, "IntAbs": {funcName: "IntAbs", native: divergence_hunt185.IntAbs}, "IntNeg": {funcName: "IntNeg", native: divergence_hunt185.IntNeg}, "IntCmp": {funcName: "IntCmp", native: divergence_hunt185.IntCmp}, "IntSet": {funcName: "IntSet", native: divergence_hunt185.IntSet}, "IntLargeNumber": {funcName: "IntLargeNumber", native: divergence_hunt185.IntLargeNumber},
	}})
}
func TestDivergenceHunt186(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt186Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SliceCopyBasic": {funcName: "SliceCopyBasic", native: divergence_hunt186.SliceCopyBasic}, "SliceCopyPartial": {funcName: "SliceCopyPartial", native: divergence_hunt186.SliceCopyPartial}, "SliceCopyOverlap": {funcName: "SliceCopyOverlap", native: divergence_hunt186.SliceCopyOverlap}, "SliceCopyString": {funcName: "SliceCopyString", native: divergence_hunt186.SliceCopyString, knownIssue: true}, "SliceClear": {funcName: "SliceClear", native: divergence_hunt186.SliceClear}, "SliceAppendGrow": {funcName: "SliceAppendGrow", native: divergence_hunt186.SliceAppendGrow}, "SliceClip": {funcName: "SliceClip", native: divergence_hunt186.SliceClip}, "SliceDeleteElement": {funcName: "SliceDeleteElement", native: divergence_hunt186.SliceDeleteElement}, "SliceInsertElement": {funcName: "SliceInsertElement", native: divergence_hunt186.SliceInsertElement}, "SliceReverse": {funcName: "SliceReverse", native: divergence_hunt186.SliceReverse}, "SliceFilter": {funcName: "SliceFilter", native: divergence_hunt186.SliceFilter},
	}})
}
func TestDivergenceHunt187(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt187Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapLength": {funcName: "MapLength", native: divergence_hunt187.MapLength}, "MapContainsKey": {funcName: "MapContainsKey", native: divergence_hunt187.MapContainsKey}, "MapDelete": {funcName: "MapDelete", native: divergence_hunt187.MapDelete}, "MapDeleteNonExistent": {funcName: "MapDeleteNonExistent", native: divergence_hunt187.MapDeleteNonExistent}, "MapGetZeroValue": {funcName: "MapGetZeroValue", native: divergence_hunt187.MapGetZeroValue}, "MapOverwrite": {funcName: "MapOverwrite", native: divergence_hunt187.MapOverwrite}, "MapAddKey": {funcName: "MapAddKey", native: divergence_hunt187.MapAddKey}, "MapNil": {funcName: "MapNil", native: divergence_hunt187.MapNil}, "MapEmpty": {funcName: "MapEmpty", native: divergence_hunt187.MapEmpty}, "MapKeyTypes": {funcName: "MapKeyTypes", native: divergence_hunt187.MapKeyTypes}, "MapValueSlice": {funcName: "MapValueSlice", native: divergence_hunt187.MapValueSlice}, "MapClearAll": {funcName: "MapClearAll", native: divergence_hunt187.MapClearAll},
	}})
}
func TestDivergenceHunt188(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt188Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ChannelBidirectional": {funcName: "ChannelBidirectional", native: divergence_hunt188.ChannelBidirectional}, "ChannelSendOnly": {funcName: "ChannelSendOnly", native: divergence_hunt188.ChannelSendOnly}, "ChannelReceiveOnly": {funcName: "ChannelReceiveOnly", native: divergence_hunt188.ChannelReceiveOnly}, "ChannelBufferedFull": {funcName: "ChannelBufferedFull", native: divergence_hunt188.ChannelBufferedFull}, "ChannelRange": {funcName: "ChannelRange", native: divergence_hunt188.ChannelRange}, "ChannelCloseCheck": {funcName: "ChannelCloseCheck", native: divergence_hunt188.ChannelCloseCheck}, "ChannelSelectDefault": {funcName: "ChannelSelectDefault", native: divergence_hunt188.ChannelSelectDefault}, "ChannelNil": {funcName: "ChannelNil", native: divergence_hunt188.ChannelNil}, "ChannelMakeZero": {funcName: "ChannelMakeZero", native: divergence_hunt188.ChannelMakeZero}, "ChannelSendReceiveOrder": {funcName: "ChannelSendReceiveOrder", native: divergence_hunt188.ChannelSendReceiveOrder},
	}})
}
func TestDivergenceHunt189(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt189Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicTypeAssertion": {funcName: "BasicTypeAssertion", native: divergence_hunt189.BasicTypeAssertion}, "TypeAssertionFail": {funcName: "TypeAssertionFail", native: divergence_hunt189.TypeAssertionFail}, "TypeAssertionString": {funcName: "TypeAssertionString", native: divergence_hunt189.TypeAssertionString}, "TypeAssertionFloat64": {funcName: "TypeAssertionFloat64", native: divergence_hunt189.TypeAssertionFloat64}, "TypeAssertionBool": {funcName: "TypeAssertionBool", native: divergence_hunt189.TypeAssertionBool}, "TypeAssertionSlice": {funcName: "TypeAssertionSlice", native: divergence_hunt189.TypeAssertionSlice}, "TypeAssertionMap": {funcName: "TypeAssertionMap", native: divergence_hunt189.TypeAssertionMap}, "TypeAssertionNil": {funcName: "TypeAssertionNil", native: divergence_hunt189.TypeAssertionNil}, "TypeAssertionInterface": {funcName: "TypeAssertionInterface", native: divergence_hunt189.TypeAssertionInterface}, "TypeAssertionPointer": {funcName: "TypeAssertionPointer", native: divergence_hunt189.TypeAssertionPointer}, "TypeAssertionChained": {funcName: "TypeAssertionChained", native: divergence_hunt189.TypeAssertionChained},
	}})
}
func TestDivergenceHunt190(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt190Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicTypeSwitch": {funcName: "BasicTypeSwitch", native: divergence_hunt190.BasicTypeSwitch}, "TypeSwitchString": {funcName: "TypeSwitchString", native: divergence_hunt190.TypeSwitchString}, "TypeSwitchMultipleTypes": {funcName: "TypeSwitchMultipleTypes", native: divergence_hunt190.TypeSwitchMultipleTypes}, "TypeSwitchBool": {funcName: "TypeSwitchBool", native: divergence_hunt190.TypeSwitchBool}, "TypeSwitchSlice": {funcName: "TypeSwitchSlice", native: divergence_hunt190.TypeSwitchSlice}, "TypeSwitchMap": {funcName: "TypeSwitchMap", native: divergence_hunt190.TypeSwitchMap}, "TypeSwitchNil": {funcName: "TypeSwitchNil", native: divergence_hunt190.TypeSwitchNil}, "TypeSwitchPointer": {funcName: "TypeSwitchPointer", native: divergence_hunt190.TypeSwitchPointer}, "TypeSwitchInterface": {funcName: "TypeSwitchInterface", native: divergence_hunt190.TypeSwitchInterface, knownIssue: true}, "TypeSwitchDefault": {funcName: "TypeSwitchDefault", native: divergence_hunt190.TypeSwitchDefault}, "TypeSwitchFunction": {funcName: "TypeSwitchFunction", native: divergence_hunt190.TypeSwitchFunction},
	}})
}
func TestDivergenceHunt191(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt191Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NaNEquality": {funcName: "NaNEquality", native: divergence_hunt191.NaNEquality}, "NaNIsNaN": {funcName: "NaNIsNaN", native: divergence_hunt191.NaNIsNaN}, "PositiveInf": {funcName: "PositiveInf", native: divergence_hunt191.PositiveInf}, "NegativeInf": {funcName: "NegativeInf", native: divergence_hunt191.NegativeInf}, "InfArithmetic": {funcName: "InfArithmetic", native: divergence_hunt191.InfArithmetic}, "InfComparison": {funcName: "InfComparison", native: divergence_hunt191.InfComparison}, "ZeroSign": {funcName: "ZeroSign", native: divergence_hunt191.ZeroSign}, "NaNPropagation": {funcName: "NaNPropagation", native: divergence_hunt191.NaNPropagation}, "InfDivision": {funcName: "InfDivision", native: divergence_hunt191.InfDivision}, "FiniteCheck": {funcName: "FiniteCheck", native: divergence_hunt191.FiniteCheck},
	}})
}
func TestDivergenceHunt192(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt192Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ComplexAddition": {funcName: "ComplexAddition", native: divergence_hunt192.ComplexAddition}, "ComplexSubtraction": {funcName: "ComplexSubtraction", native: divergence_hunt192.ComplexSubtraction}, "ComplexMultiplication": {funcName: "ComplexMultiplication", native: divergence_hunt192.ComplexMultiplication}, "ComplexDivision": {funcName: "ComplexDivision", native: divergence_hunt192.ComplexDivision}, "ComplexConjugate": {funcName: "ComplexConjugate", native: divergence_hunt192.ComplexConjugate}, "ComplexMagnitude": {funcName: "ComplexMagnitude", native: divergence_hunt192.ComplexMagnitude}, "Complex64Operations": {funcName: "Complex64Operations", native: divergence_hunt192.Complex64Operations}, "ComplexZero": {funcName: "ComplexZero", native: divergence_hunt192.ComplexZero}, "ComplexPureReal": {funcName: "ComplexPureReal", native: divergence_hunt192.ComplexPureReal}, "ComplexPureImag": {funcName: "ComplexPureImag", native: divergence_hunt192.ComplexPureImag}, "ComplexComparison": {funcName: "ComplexComparison", native: divergence_hunt192.ComplexComparison},
	}})
}
func TestDivergenceHunt193(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt193Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BitwiseAndSigned": {funcName: "BitwiseAndSigned", native: divergence_hunt193.BitwiseAndSigned}, "BitwiseOrSigned": {funcName: "BitwiseOrSigned", native: divergence_hunt193.BitwiseOrSigned}, "BitwiseXorSigned": {funcName: "BitwiseXorSigned", native: divergence_hunt193.BitwiseXorSigned}, "BitwiseNotSigned": {funcName: "BitwiseNotSigned", native: divergence_hunt193.BitwiseNotSigned}, "SignBitIsolation": {funcName: "SignBitIsolation", native: divergence_hunt193.SignBitIsolation}, "SignExtension": {funcName: "SignExtension", native: divergence_hunt193.SignExtension}, "AbsoluteValueBitwise": {funcName: "AbsoluteValueBitwise", native: divergence_hunt193.AbsoluteValueBitwise}, "MinMaxBitwise": {funcName: "MinMaxBitwise", native: divergence_hunt193.MinMaxBitwise}, "SwapXor": {funcName: "SwapXor", native: divergence_hunt193.SwapXor}, "ClearLowestBit": {funcName: "ClearLowestBit", native: divergence_hunt193.ClearLowestBit}, "IsolateLowestBit": {funcName: "IsolateLowestBit", native: divergence_hunt193.IsolateLowestBit},
	}})
}
func TestDivergenceHunt194(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt194Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"VariableLeftShift": {funcName: "VariableLeftShift", native: divergence_hunt194.VariableLeftShift}, "VariableRightShift": {funcName: "VariableRightShift", native: divergence_hunt194.VariableRightShift}, "SignedRightShift": {funcName: "SignedRightShift", native: divergence_hunt194.SignedRightShift}, "ShiftByZero": {funcName: "ShiftByZero", native: divergence_hunt194.ShiftByZero}, "LargeShiftCount": {funcName: "LargeShiftCount", native: divergence_hunt194.LargeShiftCount}, "ShiftAssignment": {funcName: "ShiftAssignment", native: divergence_hunt194.ShiftAssignment}, "ShiftInExpression": {funcName: "ShiftInExpression", native: divergence_hunt194.ShiftInExpression}, "ShiftThenMask": {funcName: "ShiftThenMask", native: divergence_hunt194.ShiftThenMask}, "ShiftBounds": {funcName: "ShiftBounds", native: divergence_hunt194.ShiftBounds}, "CircularShiftPattern": {funcName: "CircularShiftPattern", native: divergence_hunt194.CircularShiftPattern},
	}})
}
func TestDivergenceHunt195(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt195Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Uint8Overflow": {funcName: "Uint8Overflow", native: divergence_hunt195.Uint8Overflow}, "Uint16Overflow": {funcName: "Uint16Overflow", native: divergence_hunt195.Uint16Overflow}, "Int8PositiveOverflow": {funcName: "Int8PositiveOverflow", native: divergence_hunt195.Int8PositiveOverflow}, "Int8NegativeOverflow": {funcName: "Int8NegativeOverflow", native: divergence_hunt195.Int8NegativeOverflow}, "Int16Overflow": {funcName: "Int16Overflow", native: divergence_hunt195.Int16Overflow}, "MultiplicationOverflow": {funcName: "MultiplicationOverflow", native: divergence_hunt195.MultiplicationOverflow}, "CompoundOverflow": {funcName: "CompoundOverflow", native: divergence_hunt195.CompoundOverflow}, "UnderflowSubtraction": {funcName: "UnderflowSubtraction", native: divergence_hunt195.UnderflowSubtraction}, "IncrementOverflow": {funcName: "IncrementOverflow", native: divergence_hunt195.IncrementOverflow}, "DecrementUnderflow": {funcName: "DecrementUnderflow", native: divergence_hunt195.DecrementUnderflow}, "AdditiveChainOverflow": {funcName: "AdditiveChainOverflow", native: divergence_hunt195.AdditiveChainOverflow},
	}})
}
func TestDivergenceHunt196(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt196Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ArrayCopyByValue": {funcName: "ArrayCopyByValue", native: divergence_hunt196.ArrayCopyByValue},
		"ArrayPassByValue": {funcName: "ArrayPassByValue", native: divergence_hunt196.ArrayPassByValue},
		"ArrayReturnByValue": {funcName: "ArrayReturnByValue", native: divergence_hunt196.ArrayReturnByValue},
		"ArrayEquality": {funcName: "ArrayEquality", native: divergence_hunt196.ArrayEquality},
		"ArrayOfStructs": {funcName: "ArrayOfStructs", native: divergence_hunt196.ArrayOfStructs},
		"ArraySliceConversion": {funcName: "ArraySliceConversion", native: divergence_hunt196.ArraySliceConversion, knownIssue: true},
		"ArrayIndexing": {funcName: "ArrayIndexing", native: divergence_hunt196.ArrayIndexing},
		"ArrayLength": {funcName: "ArrayLength", native: divergence_hunt196.ArrayLength},
		"ArrayLiteral": {funcName: "ArrayLiteral", native: divergence_hunt196.ArrayLiteral},
		"ArrayIteration": {funcName: "ArrayIteration", native: divergence_hunt196.ArrayIteration},
		"ArrayOfArrays": {funcName: "ArrayOfArrays", native: divergence_hunt196.ArrayOfArrays},
	}})
}
func TestDivergenceHunt197(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt197Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicPointer": {funcName: "BasicPointer", native: divergence_hunt197.BasicPointer}, "PointerAssignment": {funcName: "PointerAssignment", native: divergence_hunt197.PointerAssignment}, "PointerToPointer": {funcName: "PointerToPointer", native: divergence_hunt197.PointerToPointer}, "PointerComparison": {funcName: "PointerComparison", native: divergence_hunt197.PointerComparison}, "PointerZeroValue": {funcName: "PointerZeroValue", native: divergence_hunt197.PointerZeroValue}, "PointerToArray": {funcName: "PointerToArray", native: divergence_hunt197.PointerToArray}, "PointerToStruct": {funcName: "PointerToStruct", native: divergence_hunt197.PointerToStruct}, "PointerArithmeticSimulated": {funcName: "PointerArithmeticSimulated", native: divergence_hunt197.PointerArithmeticSimulated}, "PointerInStruct": {funcName: "PointerInStruct", native: divergence_hunt197.PointerInStruct}, "PointerSwap": {funcName: "PointerSwap", native: divergence_hunt197.PointerSwap}, "PointerToInterface": {funcName: "PointerToInterface", native: divergence_hunt197.PointerToInterface},
	}})
}
func TestDivergenceHunt198(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt198Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TypeAliasConversion": {funcName: "TypeAliasConversion", native: divergence_hunt198.TypeAliasConversion}, "TypeAliasComparison": {funcName: "TypeAliasComparison", native: divergence_hunt198.TypeAliasComparison}, "UnderlyingType": {funcName: "UnderlyingType", native: divergence_hunt198.UnderlyingType}, "ByteSliceToStringConversion": {funcName: "ByteSliceToStringConversion", native: divergence_hunt198.ByteSliceToStringConversion}, "StringToByteSliceConversion": {funcName: "StringToByteSliceConversion", native: divergence_hunt198.StringToByteSliceConversion}, "RuneSliceConversion": {funcName: "RuneSliceConversion", native: divergence_hunt198.RuneSliceConversion}, "StringToRuneSlice": {funcName: "StringToRuneSlice", native: divergence_hunt198.StringToRuneSlice}, "BinaryDataRepresentation": {funcName: "BinaryDataRepresentation", native: divergence_hunt198.BinaryDataRepresentation}, "StructLayoutExploration": {funcName: "StructLayoutExploration", native: divergence_hunt198.StructLayoutExploration}, "SliceHeaderConcept": {funcName: "SliceHeaderConcept", native: divergence_hunt198.SliceHeaderConcept}, "StringImmutabilityPattern": {funcName: "StringImmutabilityPattern", native: divergence_hunt198.StringImmutabilityPattern},
	}})
}
func TestDivergenceHunt199(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt199Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"TypeSwitchDynamicDispatch": {funcName: "TypeSwitchDynamicDispatch", native: divergence_hunt199.TypeSwitchDynamicDispatch}, "TypeAssertDynamicDispatch": {funcName: "TypeAssertDynamicDispatch", native: divergence_hunt199.TypeAssertDynamicDispatch}, "InterfaceHoldingDifferentTypes": {funcName: "InterfaceHoldingDifferentTypes", native: divergence_hunt199.InterfaceHoldingDifferentTypes}, "NilInterfaceValue": {funcName: "NilInterfaceValue", native: divergence_hunt199.NilInterfaceValue}, "InterfaceWithNilPointer": {funcName: "InterfaceWithNilPointer", native: divergence_hunt199.InterfaceWithNilPointer}, "EmptyInterfaceIdentity": {funcName: "EmptyInterfaceIdentity", native: divergence_hunt199.EmptyInterfaceIdentity}, "DynamicMethodDispatch": {funcName: "DynamicMethodDispatch", native: divergence_hunt199.DynamicMethodDispatch}, "InterfaceSlice": {funcName: "InterfaceSlice", native: divergence_hunt199.InterfaceSlice}, "InterfaceMap": {funcName: "InterfaceMap", native: divergence_hunt199.InterfaceMap}, "NestedInterface": {funcName: "NestedInterface", native: divergence_hunt199.NestedInterface}, "InterfaceEquality": {funcName: "InterfaceEquality", native: divergence_hunt199.InterfaceEquality},
	}})
}
func TestDivergenceHunt200(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt200Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InterfaceNilComparison": {funcName: "InterfaceNilComparison", native: divergence_hunt200.InterfaceNilComparison},
		"InterfaceTypedNil": {funcName: "InterfaceTypedNil", native: divergence_hunt200.InterfaceTypedNil},
		"InterfaceValueExtraction": {funcName: "InterfaceValueExtraction", native: divergence_hunt200.InterfaceValueExtraction},
		"InterfaceTypeAssertionPanic": {funcName: "InterfaceTypeAssertionPanic", native: divergence_hunt200.InterfaceTypeAssertionPanic},
		"EmptyInterfaceStorage": {funcName: "EmptyInterfaceStorage", native: divergence_hunt200.EmptyInterfaceStorage},
		"InterfaceMethodSet": {funcName: "InterfaceMethodSet", native: divergence_hunt200.InterfaceMethodSet},
		"InterfaceEmbedding": {funcName: "InterfaceEmbedding", native: divergence_hunt200.InterfaceEmbedding},
		"InterfaceAssignmentCompatibility": {funcName: "InterfaceAssignmentCompatibility", native: divergence_hunt200.InterfaceAssignmentCompatibility},
		"InterfaceComparisonWithDifferentTypes": {funcName: "InterfaceComparisonWithDifferentTypes", native: divergence_hunt200.InterfaceComparisonWithDifferentTypes, knownIssue: true},
		"InterfacePointerReceiver": {funcName: "InterfacePointerReceiver", native: divergence_hunt200.InterfacePointerReceiver},
		"InterfaceValueReceiver": {funcName: "InterfaceValueReceiver", native: divergence_hunt200.InterfaceValueReceiver},
	}})
}
func TestDivergenceHunt201(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt201Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ArrayCopyByValue": {funcName: "ArrayCopyByValue", native: divergence_hunt201.ArrayCopyByValue},
		"ArrayEquality": {funcName: "ArrayEquality", native: divergence_hunt201.ArrayEquality},
		"ArrayInequality": {funcName: "ArrayInequality", native: divergence_hunt201.ArrayInequality},
		"ArraySliceShare": {funcName: "ArraySliceShare", native: divergence_hunt201.ArraySliceShare, knownIssue: true},
		"ArrayOfArrays": {funcName: "ArrayOfArrays", native: divergence_hunt201.ArrayOfArrays},
		"ArrayZeroValue": {funcName: "ArrayZeroValue", native: divergence_hunt201.ArrayZeroValue},
		"ArrayLenCap": {funcName: "ArrayLenCap", native: divergence_hunt201.ArrayLenCap},
		"ArrayRange": {funcName: "ArrayRange", native: divergence_hunt201.ArrayRange},
		"ArrayLiteral": {funcName: "ArrayLiteral", native: divergence_hunt201.ArrayLiteral},
		"ArrayIndexAccess": {funcName: "ArrayIndexAccess", native: divergence_hunt201.ArrayIndexAccess},
		"ArrayPointerDeref": {funcName: "ArrayPointerDeref", native: divergence_hunt201.ArrayPointerDeref},
	}})
}
func TestDivergenceHunt202(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt202Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicEmbedding": {funcName: "BasicEmbedding", native: divergence_hunt202.BasicEmbedding}, "EmbeddedFieldAccess": {funcName: "EmbeddedFieldAccess", native: divergence_hunt202.EmbeddedFieldAccess}, "ShadowedField": {funcName: "ShadowedField", native: divergence_hunt202.ShadowedField}, "DoubleEmbedded": {funcName: "DoubleEmbedded", native: divergence_hunt202.DoubleEmbedded}, "DeepShadowing": {funcName: "DeepShadowing", native: divergence_hunt202.DeepShadowing}, "EmbeddedPointer": {funcName: "EmbeddedPointer", native: divergence_hunt202.EmbeddedPointer}, "NamedEmbeddedAccess": {funcName: "NamedEmbeddedAccess", native: divergence_hunt202.NamedEmbeddedAccess}, "EmbeddedMethodAccess": {funcName: "EmbeddedMethodAccess", native: divergence_hunt202.EmbeddedMethodAccess},
	}})
}
func TestDivergenceHunt203(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt203Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicPromotion": {funcName: "BasicPromotion", native: divergence_hunt203.BasicPromotion}, "MethodOverride": {funcName: "MethodOverride", native: divergence_hunt203.MethodOverride}, "EmbeddedMethodCall": {funcName: "EmbeddedMethodCall", native: divergence_hunt203.EmbeddedMethodCall}, "InterfaceWithPromoted": {funcName: "InterfaceWithPromoted", native: divergence_hunt203.InterfaceWithPromoted}, "DeepEmbedding": {funcName: "DeepEmbedding", native: divergence_hunt203.DeepEmbedding}, "PromotedValueMethod": {funcName: "PromotedValueMethod", native: divergence_hunt203.PromotedValueMethod}, "EmbeddedPointerMethod": {funcName: "EmbeddedPointerMethod", native: divergence_hunt203.EmbeddedPointerMethod},
	}})
}
func TestDivergenceHunt204(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt204Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ValueReceiverNoMutate": {funcName: "ValueReceiverNoMutate", native: divergence_hunt204.ValueReceiverNoMutate}, "PointerReceiverMutates": {funcName: "PointerReceiverMutates", native: divergence_hunt204.PointerReceiverMutates}, "ValueReceiverOnPointer": {funcName: "ValueReceiverOnPointer", native: divergence_hunt204.ValueReceiverOnPointer}, "PointerReceiverOnValue": {funcName: "PointerReceiverOnValue", native: divergence_hunt204.PointerReceiverOnValue}, "MethodSetOnValue": {funcName: "MethodSetOnValue", native: divergence_hunt204.MethodSetOnValue}, "ContainerValueReceiver": {funcName: "ContainerValueReceiver", native: divergence_hunt204.ContainerValueReceiver}, "InterfaceWithPointer": {funcName: "InterfaceWithPointer", native: divergence_hunt204.InterfaceWithPointer}, "MixedReceivers": {funcName: "MixedReceivers", native: divergence_hunt204.MixedReceivers}, "CopyInValueReceiver": {funcName: "CopyInValueReceiver", native: divergence_hunt204.CopyInValueReceiver},
	}})
}
func TestDivergenceHunt205(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt205Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NilInterface": {funcName: "NilInterface", native: divergence_hunt205.NilInterface},
		"NilPointerInInterface": {funcName: "NilPointerInInterface", native: divergence_hunt205.NilPointerInInterface},
		"TypedNilInterface": {funcName: "TypedNilInterface", native: divergence_hunt205.TypedNilInterface},
		"NonNilInterface": {funcName: "NonNilInterface", native: divergence_hunt205.NonNilInterface},
		"NilErrorInterface": {funcName: "NilErrorInterface", native: divergence_hunt205.NilErrorInterface},
		"NilPointerAsError": {funcName: "NilPointerAsError", native: divergence_hunt205.NilPointerAsError},
		"InterfaceTypeWithNil": {funcName: "InterfaceTypeWithNil", native: divergence_hunt205.InterfaceTypeWithNil},
		"CompareNilInterfaces": {funcName: "CompareNilInterfaces", native: divergence_hunt205.CompareNilInterfaces},
		"ReturnNilInterface": {funcName: "ReturnNilInterface", native: divergence_hunt205.ReturnNilInterface},
		"ReturnNilPointerInterface": {funcName: "ReturnNilPointerInterface", native: divergence_hunt205.ReturnNilPointerInterface, knownIssue: true},
		"NilReturnComparison": {funcName: "NilReturnComparison", native: divergence_hunt205.NilReturnComparison},
	}})
}
func TestDivergenceHunt206(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt206Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicTypeAssertion": {funcName: "BasicTypeAssertion", native: divergence_hunt206.BasicTypeAssertion},
		"TypeAssertionFail": {funcName: "TypeAssertionFail", native: divergence_hunt206.TypeAssertionFail},
		"TypeAssertionPanic": {funcName: "TypeAssertionPanic", native: divergence_hunt206.TypeAssertionPanic},
		"TypeSwitchWithInterface": {funcName: "TypeSwitchWithInterface", native: divergence_hunt206.TypeSwitchWithInterface},
		"NestedInterfaceAssertion": {funcName: "NestedInterfaceAssertion", native: divergence_hunt206.NestedInterfaceAssertion, knownIssue: true},
		"InterfaceToConcrete": {funcName: "InterfaceToConcrete", native: divergence_hunt206.InterfaceToConcrete},
		"PointerTypeAssertion": {funcName: "PointerTypeAssertion", native: divergence_hunt206.PointerTypeAssertion},
		"NilInterfaceAssertion": {funcName: "NilInterfaceAssertion", native: divergence_hunt206.NilInterfaceAssertion},
		"MultipleTypeAssertions": {funcName: "MultipleTypeAssertions", native: divergence_hunt206.MultipleTypeAssertions},
		"AssertToInterface": {funcName: "AssertToInterface", native: divergence_hunt206.AssertToInterface},
	}})
}
func TestDivergenceHunt207(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt207Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"EmptyStructSize": {funcName: "EmptyStructSize", native: divergence_hunt207.EmptyStructSize}, "EmptyStructEquality": {funcName: "EmptyStructEquality", native: divergence_hunt207.EmptyStructEquality}, "ZeroSizeArray": {funcName: "ZeroSizeArray", native: divergence_hunt207.ZeroSizeArray}, "EmptyStructSlice": {funcName: "EmptyStructSlice", native: divergence_hunt207.EmptyStructSlice}, "EmptyStructMap": {funcName: "EmptyStructMap", native: divergence_hunt207.EmptyStructMap}, "EmptyStructChannel": {funcName: "EmptyStructChannel", native: divergence_hunt207.EmptyStructChannel}, "EmptyStructMethod": {funcName: "EmptyStructMethod", native: divergence_hunt207.EmptyStructMethod}, "EmptyStructInStruct": {funcName: "EmptyStructInStruct", native: divergence_hunt207.EmptyStructInStruct}, "EmptyStructPointer": {funcName: "EmptyStructPointer", native: divergence_hunt207.EmptyStructPointer}, "EmptyStructInterface": {funcName: "EmptyStructInterface", native: divergence_hunt207.EmptyStructInterface},
	}})
}
func TestDivergenceHunt208(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt208Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SimpleStructAccess": {funcName: "SimpleStructAccess", native: divergence_hunt208.SimpleStructAccess}, "PackedStructAccess": {funcName: "PackedStructAccess", native: divergence_hunt208.PackedStructAccess}, "AlignedStructAccess": {funcName: "AlignedStructAccess", native: divergence_hunt208.AlignedStructAccess}, "BoolStructAccess": {funcName: "BoolStructAccess", native: divergence_hunt208.BoolStructAccess}, "PointerStructAccess": {funcName: "PointerStructAccess", native: divergence_hunt208.PointerStructAccess}, "SliceStructAccess": {funcName: "SliceStructAccess", native: divergence_hunt208.SliceStructAccess}, "StringStructAccess": {funcName: "StringStructAccess", native: divergence_hunt208.StringStructAccess}, "NestedStructAccess": {funcName: "NestedStructAccess", native: divergence_hunt208.NestedStructAccess}, "StructArrayAccess": {funcName: "StructArrayAccess", native: divergence_hunt208.StructArrayAccess}, "StructSliceAccess": {funcName: "StructSliceAccess", native: divergence_hunt208.StructSliceAccess}, "StructEquality": {funcName: "StructEquality", native: divergence_hunt208.StructEquality},
	}})
}
func TestDivergenceHunt209(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt209Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"FunctionNilComparison": {funcName: "FunctionNilComparison", native: divergence_hunt209.FunctionNilComparison}, "FunctionNonNilComparison": {funcName: "FunctionNonNilComparison", native: divergence_hunt209.FunctionNonNilComparison}, "FunctionSameVar": {funcName: "FunctionSameVar", native: divergence_hunt209.FunctionSameVar}, "FunctionDifferentVars": {funcName: "FunctionDifferentVars", native: divergence_hunt209.FunctionDifferentVars}, "FunctionInStruct": {funcName: "FunctionInStruct", native: divergence_hunt209.FunctionInStruct}, "FunctionInMap": {funcName: "FunctionInMap", native: divergence_hunt209.FunctionInMap}, "FunctionSlice": {funcName: "FunctionSlice", native: divergence_hunt209.FunctionSlice}, "FunctionAsInterface": {funcName: "FunctionAsInterface", native: divergence_hunt209.FunctionAsInterface}, "FunctionVariableCapture": {funcName: "FunctionVariableCapture", native: divergence_hunt209.FunctionVariableCapture}, "FunctionAssignedToVar": {funcName: "FunctionAssignedToVar", native: divergence_hunt209.FunctionAssignedToVar}, "ClosureComparison": {funcName: "ClosureComparison", native: divergence_hunt209.ClosureComparison},
	}})
}
func TestDivergenceHunt210(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt210Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"LoopCaptureClassic": {funcName: "LoopCaptureClassic", native: divergence_hunt210.LoopCaptureClassic}, "LoopCaptureFixed": {funcName: "LoopCaptureFixed", native: divergence_hunt210.LoopCaptureFixed}, "ClosureModifiesCaptured": {funcName: "ClosureModifiesCaptured", native: divergence_hunt210.ClosureModifiesCaptured}, "MultipleClosuresShareVar": {funcName: "MultipleClosuresShareVar", native: divergence_hunt210.MultipleClosuresShareVar}, "ClosureCapturesSlice": {funcName: "ClosureCapturesSlice", native: divergence_hunt210.ClosureCapturesSlice}, "ClosureCapturesMap": {funcName: "ClosureCapturesMap", native: divergence_hunt210.ClosureCapturesMap}, "NestedClosure": {funcName: "NestedClosure", native: divergence_hunt210.NestedClosure}, "ClosureAsReturn": {funcName: "ClosureAsReturn", native: divergence_hunt210.ClosureAsReturn}, "ClosureReceivesItself": {funcName: "ClosureReceivesItself", native: divergence_hunt210.ClosureReceivesItself}, "DeferredClosureCapture": {funcName: "DeferredClosureCapture", native: divergence_hunt210.DeferredClosureCapture}, "RangeCaptureIssue": {funcName: "RangeCaptureIssue", native: divergence_hunt210.RangeCaptureIssue},
	}})
}
func TestDivergenceHunt211(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt211Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SelectWithTimeout": {funcName: "SelectWithTimeout", native: divergence_hunt211.SelectWithTimeout}, "SelectTimeoutNotTriggered": {funcName: "SelectTimeoutNotTriggered", native: divergence_hunt211.SelectTimeoutNotTriggered}, "SelectMultipleWithTimeout": {funcName: "SelectMultipleWithTimeout", native: divergence_hunt211.SelectMultipleWithTimeout}, "SelectDefaultWithTimeout": {funcName: "SelectDefaultWithTimeout", native: divergence_hunt211.SelectDefaultWithTimeout}, "NestedSelectTimeout": {funcName: "NestedSelectTimeout", native: divergence_hunt211.NestedSelectTimeout}, "SelectTimeoutChannelClosed": {funcName: "SelectTimeoutChannelClosed", native: divergence_hunt211.SelectTimeoutChannelClosed}, "SelectTimeoutZeroDuration": {funcName: "SelectTimeoutZeroDuration", native: divergence_hunt211.SelectTimeoutZeroDuration}, "BufferedChannelSelectTimeout": {funcName: "BufferedChannelSelectTimeout", native: divergence_hunt211.BufferedChannelSelectTimeout}, "SelectTimeoutReuse": {funcName: "SelectTimeoutReuse", native: divergence_hunt211.SelectTimeoutReuse}, "LongTimeoutNotTriggered": {funcName: "LongTimeoutNotTriggered", native: divergence_hunt211.LongTimeoutNotTriggered},
	}})
}
func TestDivergenceHunt212(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt212Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RangeOverClosedChannel": {funcName: "RangeOverClosedChannel", native: divergence_hunt212.RangeOverClosedChannel}, "RangeOverEmptyClosedChannel": {funcName: "RangeOverEmptyClosedChannel", native: divergence_hunt212.RangeOverEmptyClosedChannel}, "CloseChannelReceiveValues": {funcName: "CloseChannelReceiveValues", native: divergence_hunt212.CloseChannelReceiveValues}, "CloseChannelWithGoroutine": {funcName: "CloseChannelWithGoroutine", native: divergence_hunt212.CloseChannelWithGoroutine, knownIssue: true}, "DoubleClosePanic": {funcName: "DoubleClosePanic", native: divergence_hunt212.DoubleClosePanic}, "CloseNilChannelPanic": {funcName: "CloseNilChannelPanic", native: divergence_hunt212.CloseNilChannelPanic}, "RangeWithBreak": {funcName: "RangeWithBreak", native: divergence_hunt212.RangeWithBreak}, "RangeStringChannel": {funcName: "RangeStringChannel", native: divergence_hunt212.RangeStringChannel}, "RangeStructChannel": {funcName: "RangeStructChannel", native: divergence_hunt212.RangeStructChannel}, "PartialReceiveThenRange": {funcName: "PartialReceiveThenRange", native: divergence_hunt212.PartialReceiveThenRange}, "CloseSentinelPattern": {funcName: "CloseSentinelPattern", native: divergence_hunt212.CloseSentinelPattern, knownIssue: true},
	}})
}
func TestDivergenceHunt213(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt213Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BufferedChannelCapacity": {funcName: "BufferedChannelCapacity", native: divergence_hunt213.BufferedChannelCapacity}, "UnbufferedChannelCapacity": {funcName: "UnbufferedChannelCapacity", native: divergence_hunt213.UnbufferedChannelCapacity}, "BufferedChannelLength": {funcName: "BufferedChannelLength", native: divergence_hunt213.BufferedChannelLength}, "BufferedChannelBlocking": {funcName: "BufferedChannelBlocking", native: divergence_hunt213.BufferedChannelBlocking}, "UnbufferedChannelSynchronization": {funcName: "UnbufferedChannelSynchronization", native: divergence_hunt213.UnbufferedChannelSynchronization}, "BufferedVsUnbufferedSend": {funcName: "BufferedVsUnbufferedSend", native: divergence_hunt213.BufferedVsUnbufferedSend}, "BufferedChannelDrainOrder": {funcName: "BufferedChannelDrainOrder", native: divergence_hunt213.BufferedChannelDrainOrder}, "ZeroBufferChannel": {funcName: "ZeroBufferChannel", native: divergence_hunt213.ZeroBufferChannel}, "BufferedStringChannel": {funcName: "BufferedStringChannel", native: divergence_hunt213.BufferedStringChannel}, "BufferedChannelCloseSemantics": {funcName: "BufferedChannelCloseSemantics", native: divergence_hunt213.BufferedChannelCloseSemantics}, "MixedBufferedUnbuffered": {funcName: "MixedBufferedUnbuffered", native: divergence_hunt213.MixedBufferedUnbuffered}, "BufferedChannelOverwriteProtection": {funcName: "BufferedChannelOverwriteProtection", native: divergence_hunt213.BufferedChannelOverwriteProtection},
	}})
}
func TestDivergenceHunt214(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt214Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NilChannelIsNil": {funcName: "NilChannelIsNil", native: divergence_hunt214.NilChannelIsNil}, "NilChannelSendBlock": {funcName: "NilChannelSendBlock", native: divergence_hunt214.NilChannelSendBlock}, "NilChannelReceiveBlock": {funcName: "NilChannelReceiveBlock", native: divergence_hunt214.NilChannelReceiveBlock}, "NilChannelInSelect": {funcName: "NilChannelInSelect", native: divergence_hunt214.NilChannelInSelect}, "NilChannelClosePanic": {funcName: "NilChannelClosePanic", native: divergence_hunt214.NilChannelClosePanic}, "NilChannelLenCap": {funcName: "NilChannelLenCap", native: divergence_hunt214.NilChannelLenCap}, "NilChannelTypeAssertion": {funcName: "NilChannelTypeAssertion", native: divergence_hunt214.NilChannelTypeAssertion}, "NilChannelAssignAndUse": {funcName: "NilChannelAssignAndUse", native: divergence_hunt214.NilChannelAssignAndUse}, "NilChannelComparison": {funcName: "NilChannelComparison", native: divergence_hunt214.NilChannelComparison}, "NilChannelInSlice": {funcName: "NilChannelInSlice", native: divergence_hunt214.NilChannelInSlice}, "NilChannelSelectMultiple": {funcName: "NilChannelSelectMultiple", native: divergence_hunt214.NilChannelSelectMultiple},
	}})
}
func TestDivergenceHunt215(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt215Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SendOnlyChannel": {funcName: "SendOnlyChannel", native: divergence_hunt215.SendOnlyChannel}, "ReceiveOnlyChannel": {funcName: "ReceiveOnlyChannel", native: divergence_hunt215.ReceiveOnlyChannel}, "BidirectionalAsSendOnly": {funcName: "BidirectionalAsSendOnly", native: divergence_hunt215.BidirectionalAsSendOnly}, "BidirectionalAsReceiveOnly": {funcName: "BidirectionalAsReceiveOnly", native: divergence_hunt215.BidirectionalAsReceiveOnly}, "DirectionConversion": {funcName: "DirectionConversion", native: divergence_hunt215.DirectionConversion}, "FunctionWithDirectionalChannels": {funcName: "FunctionWithDirectionalChannels", native: divergence_hunt215.FunctionWithDirectionalChannels}, "PipelinePattern": {funcName: "PipelinePattern", native: divergence_hunt215.PipelinePattern}, "CannotReceiveFromSendOnly": {funcName: "CannotReceiveFromSendOnly", native: divergence_hunt215.CannotReceiveFromSendOnly}, "CannotSendToReceiveOnly": {funcName: "CannotSendToReceiveOnly", native: divergence_hunt215.CannotSendToReceiveOnly}, "DirectionalChannelInStruct": {funcName: "DirectionalChannelInStruct", native: divergence_hunt215.DirectionalChannelInStruct, knownIssue: true}, "CloseSendOnlyChannel": {funcName: "CloseSendOnlyChannel", native: divergence_hunt215.CloseSendOnlyChannel},
	}})
}
func TestDivergenceHunt216(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt216Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"WaitGroupBasic": {funcName: "WaitGroupBasic", native: divergence_hunt216.WaitGroupBasic}, "WaitGroupWithParam": {funcName: "WaitGroupWithParam", native: divergence_hunt216.WaitGroupWithParam}, "WaitGroupMultipleAdd": {funcName: "WaitGroupMultipleAdd", native: divergence_hunt216.WaitGroupMultipleAdd}, "WaitGroupAddInsideGoroutine": {funcName: "WaitGroupAddInsideGoroutine", native: divergence_hunt216.WaitGroupAddInsideGoroutine}, "WaitGroupWithMutex": {funcName: "WaitGroupWithMutex", native: divergence_hunt216.WaitGroupWithMutex}, "NestedWaitGroup": {funcName: "NestedWaitGroup", native: divergence_hunt216.NestedWaitGroup}, "WaitGroupWithErrorResult": {funcName: "WaitGroupWithErrorResult", native: divergence_hunt216.WaitGroupWithErrorResult}, "WaitGroupReuse": {funcName: "WaitGroupReuse", native: divergence_hunt216.WaitGroupReuse}, "WaitGroupWithChannelResult": {funcName: "WaitGroupWithChannelResult", native: divergence_hunt216.WaitGroupWithChannelResult}, "WaitGroupDeferDone": {funcName: "WaitGroupDeferDone", native: divergence_hunt216.WaitGroupDeferDone}, "WaitGroupZeroValue": {funcName: "WaitGroupZeroValue", native: divergence_hunt216.WaitGroupZeroValue},
	}})
}
func TestDivergenceHunt217(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt217Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"OnceBasic": {funcName: "OnceBasic", native: divergence_hunt217.OnceBasic}, "OnceWithGoroutines": {funcName: "OnceWithGoroutines", native: divergence_hunt217.OnceWithGoroutines}, "OnceWithDataInitialization": {funcName: "OnceWithDataInitialization", native: divergence_hunt217.OnceWithDataInitialization}, "OnceWithPanic": {funcName: "OnceWithPanic", native: divergence_hunt217.OnceWithPanic, knownIssue: true}, "OncePerInstance": {funcName: "OncePerInstance", native: divergence_hunt217.OncePerInstance}, "OnceWithComplexInitialization": {funcName: "OnceWithComplexInitialization", native: divergence_hunt217.OnceWithComplexInitialization}, "MultipleOnceVariables": {funcName: "MultipleOnceVariables", native: divergence_hunt217.MultipleOnceVariables}, "OnceWithMutexCombo": {funcName: "OnceWithMutexCombo", native: divergence_hunt217.OnceWithMutexCombo}, "OnceInStructSlice": {funcName: "OnceInStructSlice", native: divergence_hunt217.OnceInStructSlice}, "OnceWithLazyLoading": {funcName: "OnceWithLazyLoading", native: divergence_hunt217.OnceWithLazyLoading}, "OnceWithChannelClose": {funcName: "OnceWithChannelClose", native: divergence_hunt217.OnceWithChannelClose},
	}})
}
func TestDivergenceHunt218(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt218Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SyncMapBasic": {funcName: "SyncMapBasic", native: divergence_hunt218.SyncMapBasic}, "SyncMapLoadOrStore": {funcName: "SyncMapLoadOrStore", native: divergence_hunt218.SyncMapLoadOrStore}, "SyncMapDelete": {funcName: "SyncMapDelete", native: divergence_hunt218.SyncMapDelete}, "SyncMapRange": {funcName: "SyncMapRange", native: divergence_hunt218.SyncMapRange, knownIssue: true}, "SyncMapRangeEarlyExit": {funcName: "SyncMapRangeEarlyExit", native: divergence_hunt218.SyncMapRangeEarlyExit}, "SyncMapLoadAndDelete": {funcName: "SyncMapLoadAndDelete", native: divergence_hunt218.SyncMapLoadAndDelete}, "SyncMapSwap": {funcName: "SyncMapSwap", native: divergence_hunt218.SyncMapSwap}, "SyncMapCompareAndSwap": {funcName: "SyncMapCompareAndSwap", native: divergence_hunt218.SyncMapCompareAndSwap}, "SyncMapCompareAndDelete": {funcName: "SyncMapCompareAndDelete", native: divergence_hunt218.SyncMapCompareAndDelete}, "SyncMapWithGoroutines": {funcName: "SyncMapWithGoroutines", native: divergence_hunt218.SyncMapWithGoroutines}, "SyncMapTypeSafety": {funcName: "SyncMapTypeSafety", native: divergence_hunt218.SyncMapTypeSafety},
	}})
}
func TestDivergenceHunt219(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt219Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MutexAddInt32": {funcName: "MutexAddInt32", native: divergence_hunt219.MutexAddInt32}, "MutexAddInt64": {funcName: "MutexAddInt64", native: divergence_hunt219.MutexAddInt64}, "MutexLoadStoreInt32": {funcName: "MutexLoadStoreInt32", native: divergence_hunt219.MutexLoadStoreInt32}, "MutexLoadStoreInt64": {funcName: "MutexLoadStoreInt64", native: divergence_hunt219.MutexLoadStoreInt64}, "MutexSwapInt32": {funcName: "MutexSwapInt32", native: divergence_hunt219.MutexSwapInt32}, "MutexCompareAndSwapInt32": {funcName: "MutexCompareAndSwapInt32", native: divergence_hunt219.MutexCompareAndSwapInt32}, "ProtectedInt32": {funcName: "ProtectedInt32", native: divergence_hunt219.ProtectedInt32}, "ProtectedInt64": {funcName: "ProtectedInt64", native: divergence_hunt219.ProtectedInt64}, "MutexUint32": {funcName: "MutexUint32", native: divergence_hunt219.MutexUint32}, "MutexUint64": {funcName: "MutexUint64", native: divergence_hunt219.MutexUint64}, "MutexCounterPattern": {funcName: "MutexCounterPattern", native: divergence_hunt219.MutexCounterPattern}, "MutexFlagPattern": {funcName: "MutexFlagPattern", native: divergence_hunt219.MutexFlagPattern}, "MutexMaxValue": {funcName: "MutexMaxValue", native: divergence_hunt219.MutexMaxValue},
	}})
}
func TestDivergenceHunt220(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt220Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RWMutexBasic": {funcName: "RWMutexBasic", native: divergence_hunt220.RWMutexBasic}, "RWMutexMultipleReaders": {funcName: "RWMutexMultipleReaders", native: divergence_hunt220.RWMutexMultipleReaders, knownIssue: true}, "RWMutexWriterExclusion": {funcName: "RWMutexWriterExclusion", native: divergence_hunt220.RWMutexWriterExclusion}, "RWMutexDeferredUnlock": {funcName: "RWMutexDeferredUnlock", native: divergence_hunt220.RWMutexDeferredUnlock}, "RWMutexDeferredRUnlock": {funcName: "RWMutexDeferredRUnlock", native: divergence_hunt220.RWMutexDeferredRUnlock}, "RWMutexPromote": {funcName: "RWMutexPromote", native: divergence_hunt220.RWMutexPromote}, "RWMutexNestedReadLock": {funcName: "RWMutexNestedReadLock", native: divergence_hunt220.RWMutexNestedReadLock}, "RWMutexWithWaitGroup": {funcName: "RWMutexWithWaitGroup", native: divergence_hunt220.RWMutexWithWaitGroup}, "RWMutexTryLock": {funcName: "RWMutexTryLock", native: divergence_hunt220.RWMutexTryLock}, "RWMutexTryRLock": {funcName: "RWMutexTryRLock", native: divergence_hunt220.RWMutexTryRLock}, "RWMutexInStruct": {funcName: "RWMutexInStruct", native: divergence_hunt220.RWMutexInStruct}, "RWMutexReadPrefer": {funcName: "RWMutexReadPrefer", native: divergence_hunt220.RWMutexReadPrefer},
	}})
}
func TestDivergenceHunt221(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt221Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapIterationSum": {funcName: "MapIterationSum", native: divergence_hunt221.MapIterationSum},
		"MapIterationCount": {funcName: "MapIterationCount", native: divergence_hunt221.MapIterationCount},
		"MapIterationKeys": {funcName: "MapIterationKeys", native: divergence_hunt221.MapIterationKeys},
		"MapIterationValues": {funcName: "MapIterationValues", native: divergence_hunt221.MapIterationValues},
		"MapIterationOrder": {funcName: "MapIterationOrder", native: divergence_hunt221.MapIterationOrder},
		"MapIterationBreak": {funcName: "MapIterationBreak", native: divergence_hunt221.MapIterationBreak, knownIssue: true},
		"MapIterationContinue": {funcName: "MapIterationContinue", native: divergence_hunt221.MapIterationContinue},
		"MapIterationNested": {funcName: "MapIterationNested", native: divergence_hunt221.MapIterationNested},
		"MapIterationWithModify": {funcName: "MapIterationWithModify", native: divergence_hunt221.MapIterationWithModify},
		"MapIterationEmpty": {funcName: "MapIterationEmpty", native: divergence_hunt221.MapIterationEmpty},
		"MapIterationNil": {funcName: "MapIterationNil", native: divergence_hunt221.MapIterationNil},
	}})
}
func TestDivergenceHunt222(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt222Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapDeleteDuringIteration": {funcName: "MapDeleteDuringIteration", native: divergence_hunt222.MapDeleteDuringIteration}, "MapDeleteOtherDuringIteration": {funcName: "MapDeleteOtherDuringIteration", native: divergence_hunt222.MapDeleteOtherDuringIteration}, "MapDeleteAllDuringIteration": {funcName: "MapDeleteAllDuringIteration", native: divergence_hunt222.MapDeleteAllDuringIteration}, "MapDeleteAndAddDuringIteration": {funcName: "MapDeleteAndAddDuringIteration", native: divergence_hunt222.MapDeleteAndAddDuringIteration}, "MapSafeDeleteCollect": {funcName: "MapSafeDeleteCollect", native: divergence_hunt222.MapSafeDeleteCollect}, "MapDeleteNonExistent": {funcName: "MapDeleteNonExistent", native: divergence_hunt222.MapDeleteNonExistent}, "MapDeleteThenLookup": {funcName: "MapDeleteThenLookup", native: divergence_hunt222.MapDeleteThenLookup}, "MapIterDeleteWithComplicatedKeys": {funcName: "MapIterDeleteWithComplicatedKeys", native: divergence_hunt222.MapIterDeleteWithComplicatedKeys}, "MapDeletePreservesIteration": {funcName: "MapDeletePreservesIteration", native: divergence_hunt222.MapDeletePreservesIteration}, "MapConditionalDeleteDuringIteration": {funcName: "MapConditionalDeleteDuringIteration", native: divergence_hunt222.MapConditionalDeleteDuringIteration},
	}})
}
func TestDivergenceHunt223(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt223Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapStructKeyBasic": {funcName: "MapStructKeyBasic", native: divergence_hunt223.MapStructKeyBasic}, "MapStructKeyLookup": {funcName: "MapStructKeyLookup", native: divergence_hunt223.MapStructKeyLookup}, "MapStructKeyInsert": {funcName: "MapStructKeyInsert", native: divergence_hunt223.MapStructKeyInsert}, "MapStructKeyDelete": {funcName: "MapStructKeyDelete", native: divergence_hunt223.MapStructKeyDelete}, "MapStructKeyIterate": {funcName: "MapStructKeyIterate", native: divergence_hunt223.MapStructKeyIterate}, "MapStructKeyWithStringField": {funcName: "MapStructKeyWithStringField", native: divergence_hunt223.MapStructKeyWithStringField}, "MapStructKeyZeroValue": {funcName: "MapStructKeyZeroValue", native: divergence_hunt223.MapStructKeyZeroValue}, "MapStructKeyFloat": {funcName: "MapStructKeyFloat", native: divergence_hunt223.MapStructKeyFloat}, "MapStructKeyComposite": {funcName: "MapStructKeyComposite", native: divergence_hunt223.MapStructKeyComposite}, "MapStructKeyUpdate": {funcName: "MapStructKeyUpdate", native: divergence_hunt223.MapStructKeyUpdate}, "MapStructKeyCommaOk": {funcName: "MapStructKeyCommaOk", native: divergence_hunt223.MapStructKeyCommaOk},
	}})
}
func TestDivergenceHunt224(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt224Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapInterfaceKeyBasic": {funcName: "MapInterfaceKeyBasic", native: divergence_hunt224.MapInterfaceKeyBasic}, "MapInterfaceKeyLookup": {funcName: "MapInterfaceKeyLookup", native: divergence_hunt224.MapInterfaceKeyLookup}, "MapInterfaceKeyMixedTypes": {funcName: "MapInterfaceKeyMixedTypes", native: divergence_hunt224.MapInterfaceKeyMixedTypes}, "MapInterfaceKeyNil": {funcName: "MapInterfaceKeyNil", native: divergence_hunt224.MapInterfaceKeyNil}, "MapInterfaceKeySlicePtr": {funcName: "MapInterfaceKeySlicePtr", native: divergence_hunt224.MapInterfaceKeySlicePtr}, "MapInterfaceKeyMapPtr": {funcName: "MapInterfaceKeyMapPtr", native: divergence_hunt224.MapInterfaceKeyMapPtr}, "MapInterfaceKeyIterate": {funcName: "MapInterfaceKeyIterate", native: divergence_hunt224.MapInterfaceKeyIterate}, "MapInterfaceKeyDelete": {funcName: "MapInterfaceKeyDelete", native: divergence_hunt224.MapInterfaceKeyDelete}, "MapInterfaceKeyTypeSwitch": {funcName: "MapInterfaceKeyTypeSwitch", native: divergence_hunt224.MapInterfaceKeyTypeSwitch}, "MapInterfaceKeyOverwrite": {funcName: "MapInterfaceKeyOverwrite", native: divergence_hunt224.MapInterfaceKeyOverwrite}, "MapInterfaceKeyCommaOk": {funcName: "MapInterfaceKeyCommaOk", native: divergence_hunt224.MapInterfaceKeyCommaOk},
	}})
}
func TestDivergenceHunt225(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt225Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MapPointerKeyBasic": {funcName: "MapPointerKeyBasic", native: divergence_hunt225.MapPointerKeyBasic}, "MapPointerKeyLookup": {funcName: "MapPointerKeyLookup", native: divergence_hunt225.MapPointerKeyLookup}, "MapPointerKeyDifferentAddresses": {funcName: "MapPointerKeyDifferentAddresses", native: divergence_hunt225.MapPointerKeyDifferentAddresses}, "MapPointerKeyFromSlice": {funcName: "MapPointerKeyFromSlice", native: divergence_hunt225.MapPointerKeyFromSlice, knownIssue: true}, "MapPointerKeyStruct": {funcName: "MapPointerKeyStruct", native: divergence_hunt225.MapPointerKeyStruct}, "MapPointerKeyArray": {funcName: "MapPointerKeyArray", native: divergence_hunt225.MapPointerKeyArray}, "MapPointerKeyIterate": {funcName: "MapPointerKeyIterate", native: divergence_hunt225.MapPointerKeyIterate}, "MapPointerKeyDelete": {funcName: "MapPointerKeyDelete", native: divergence_hunt225.MapPointerKeyDelete}, "MapPointerKeyNil": {funcName: "MapPointerKeyNil", native: divergence_hunt225.MapPointerKeyNil}, "MapPointerKeyModifyValue": {funcName: "MapPointerKeyModifyValue", native: divergence_hunt225.MapPointerKeyModifyValue}, "MapPointerKeyComposite": {funcName: "MapPointerKeyComposite", native: divergence_hunt225.MapPointerKeyComposite},
	}})
}
func TestDivergenceHunt226(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt226Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SliceAppendGrowth": {funcName: "SliceAppendGrowth", native: divergence_hunt226.SliceAppendGrowth},
		"SliceAppendFromNil": {funcName: "SliceAppendFromNil", native: divergence_hunt226.SliceAppendFromNil, knownIssue: true},
		"SliceAppendManyElements": {funcName: "SliceAppendManyElements", native: divergence_hunt226.SliceAppendManyElements},
		"SliceAppendToFull": {funcName: "SliceAppendToFull", native: divergence_hunt226.SliceAppendToFull},
		"SliceAppendEmptySlice": {funcName: "SliceAppendEmptySlice", native: divergence_hunt226.SliceAppendEmptySlice},
		"SliceAppendByteGrowth": {funcName: "SliceAppendByteGrowth", native: divergence_hunt226.SliceAppendByteGrowth, knownIssue: true},
		"SliceAppendInLoop": {funcName: "SliceAppendInLoop", native: divergence_hunt226.SliceAppendInLoop},
		"SliceAppendStringToBytes": {funcName: "SliceAppendStringToBytes", native: divergence_hunt226.SliceAppendStringToBytes, knownIssue: true},
		"SliceAppendSelf": {funcName: "SliceAppendSelf", native: divergence_hunt226.SliceAppendSelf},
		"SliceAppendCapacityReuse": {funcName: "SliceAppendCapacityReuse", native: divergence_hunt226.SliceAppendCapacityReuse},
		"SliceAppendLargeBatch": {funcName: "SliceAppendLargeBatch", native: divergence_hunt226.SliceAppendLargeBatch},
	}})
}
func TestDivergenceHunt227(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt227Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SliceLenCapBasic": {funcName: "SliceLenCapBasic", native: divergence_hunt227.SliceLenCapBasic}, "SliceLenCapAfterAppend": {funcName: "SliceLenCapAfterAppend", native: divergence_hunt227.SliceLenCapAfterAppend}, "SliceLenCapReslice": {funcName: "SliceLenCapReslice", native: divergence_hunt227.SliceLenCapReslice}, "SliceLenCapNil": {funcName: "SliceLenCapNil", native: divergence_hunt227.SliceLenCapNil}, "SliceLenCapEmptyLiteral": {funcName: "SliceLenCapEmptyLiteral", native: divergence_hunt227.SliceLenCapEmptyLiteral}, "SliceLenCapSlicingBounds": {funcName: "SliceLenCapSlicingBounds", native: divergence_hunt227.SliceLenCapSlicingBounds}, "SliceLenCapFullSlice": {funcName: "SliceLenCapFullSlice", native: divergence_hunt227.SliceLenCapFullSlice}, "SliceLenCapZeroLength": {funcName: "SliceLenCapZeroLength", native: divergence_hunt227.SliceLenCapZeroLength}, "SliceLenCapByteSlice": {funcName: "SliceLenCapByteSlice", native: divergence_hunt227.SliceLenCapByteSlice, knownIssue: true}, "SliceLenCapStringConversion": {funcName: "SliceLenCapStringConversion", native: divergence_hunt227.SliceLenCapStringConversion, knownIssue: true}, "SliceLenCapAfterClear": {funcName: "SliceLenCapAfterClear", native: divergence_hunt227.SliceLenCapAfterClear},
	}})
}
func TestDivergenceHunt228(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt228Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SliceSharingModify": {funcName: "SliceSharingModify", native: divergence_hunt228.SliceSharingModify}, "SliceSharingIndependence": {funcName: "SliceSharingIndependence", native: divergence_hunt228.SliceSharingIndependence}, "SliceSharingCapacityCheck": {funcName: "SliceSharingCapacityCheck", native: divergence_hunt228.SliceSharingCapacityCheck}, "SliceSharingCopy": {funcName: "SliceSharingCopy", native: divergence_hunt228.SliceSharingCopy}, "SliceSharingMultipleRefs": {funcName: "SliceSharingMultipleRefs", native: divergence_hunt228.SliceSharingMultipleRefs}, "SliceSharingResliceIndependence": {funcName: "SliceSharingResliceIndependence", native: divergence_hunt228.SliceSharingResliceIndependence}, "SliceSharingFullReslice": {funcName: "SliceSharingFullReslice", native: divergence_hunt228.SliceSharingFullReslice}, "SliceSharingPointerSharing": {funcName: "SliceSharingPointerSharing", native: divergence_hunt228.SliceSharingPointerSharing}, "SliceSharingNestedSlice": {funcName: "SliceSharingNestedSlice", native: divergence_hunt228.SliceSharingNestedSlice}, "SliceSharingAppendAndModify": {funcName: "SliceSharingAppendAndModify", native: divergence_hunt228.SliceSharingAppendAndModify}, "SliceSharingMakeCopy": {funcName: "SliceSharingMakeCopy", native: divergence_hunt228.SliceSharingMakeCopy},
	}})
}
func TestDivergenceHunt229(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt229Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NilVsEmptyLenCap": {funcName: "NilVsEmptyLenCap", native: divergence_hunt229.NilVsEmptyLenCap}, "NilVsEmptyNilCheck": {funcName: "NilVsEmptyNilCheck", native: divergence_hunt229.NilVsEmptyNilCheck}, "NilVsEmptyAppend": {funcName: "NilVsEmptyAppend", native: divergence_hunt229.NilVsEmptyAppend}, "NilVsEmptyIteration": {funcName: "NilVsEmptyIteration", native: divergence_hunt229.NilVsEmptyIteration}, "NilVsEmptyIndexingPanics": {funcName: "NilVsEmptyIndexingPanics", native: divergence_hunt229.NilVsEmptyIndexingPanics}, "NilVsEmptyJSONMarshal": {funcName: "NilVsEmptyJSONMarshal", native: divergence_hunt229.NilVsEmptyJSONMarshal}, "NilVsEmptySliceExpr": {funcName: "NilVsEmptySliceExpr", native: divergence_hunt229.NilVsEmptySliceExpr}, "NilVsEmptySlicingNil": {funcName: "NilVsEmptySlicingNil", native: divergence_hunt229.NilVsEmptySlicingNil}, "NilVsEmptyCopyToNil": {funcName: "NilVsEmptyCopyToNil", native: divergence_hunt229.NilVsEmptyCopyToNil}, "NilVsEmptyEqual": {funcName: "NilVsEmptyEqual", native: divergence_hunt229.NilVsEmptyEqual}, "NilVsEmptyRangeNil": {funcName: "NilVsEmptyRangeNil", native: divergence_hunt229.NilVsEmptyRangeNil}, "NilVsEmptyLenOnly": {funcName: "NilVsEmptyLenOnly", native: divergence_hunt229.NilVsEmptyLenOnly},
	}})
}
func TestDivergenceHunt230(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt230Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ThreeIndexBasic": {funcName: "ThreeIndexBasic", native: divergence_hunt230.ThreeIndexBasic},
		"ThreeIndexFullCapacity": {funcName: "ThreeIndexFullCapacity", native: divergence_hunt230.ThreeIndexFullCapacity},
		"ThreeIndexLimits": {funcName: "ThreeIndexLimits", native: divergence_hunt230.ThreeIndexLimits},
		"ThreeIndexAppendIndependence": {funcName: "ThreeIndexAppendIndependence", native: divergence_hunt230.ThreeIndexAppendIndependence},
		"ThreeIndexPanicLow": {funcName: "ThreeIndexPanicLow", native: divergence_hunt230.ThreeIndexPanicLow},
		"ThreeIndexPanicMax": {funcName: "ThreeIndexPanicMax", native: divergence_hunt230.ThreeIndexPanicMax},
		"ThreeIndexByteSlice": {funcName: "ThreeIndexByteSlice", native: divergence_hunt230.ThreeIndexByteSlice, knownIssue: true},
		"ThreeIndexStringSlice": {funcName: "ThreeIndexStringSlice", native: divergence_hunt230.ThreeIndexStringSlice},
		"ThreeIndexNestedSlice": {funcName: "ThreeIndexNestedSlice", native: divergence_hunt230.ThreeIndexNestedSlice},
		"ThreeIndexThenReslice": {funcName: "ThreeIndexThenReslice", native: divergence_hunt230.ThreeIndexThenReslice},
		"ThreeIndexPreserveOriginal": {funcName: "ThreeIndexPreserveOriginal", native: divergence_hunt230.ThreeIndexPreserveOriginal},
		"ThreeIndexCapacityControl": {funcName: "ThreeIndexCapacityControl", native: divergence_hunt230.ThreeIndexCapacityControl},
	}})
}
func TestDivergenceHunt231(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt231Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringByteIndex": {funcName: "StringByteIndex", native: divergence_hunt231.StringByteIndex}, "StringByteIndexUnicode": {funcName: "StringByteIndexUnicode", native: divergence_hunt231.StringByteIndexUnicode}, "StringLengthBytes": {funcName: "StringLengthBytes", native: divergence_hunt231.StringLengthBytes}, "StringIndexOutOfBounds": {funcName: "StringIndexOutOfBounds", native: divergence_hunt231.StringIndexOutOfBounds}, "StringByteLoop": {funcName: "StringByteLoop", native: divergence_hunt231.StringByteLoop}, "StringEmptyIndex": {funcName: "StringEmptyIndex", native: divergence_hunt231.StringEmptyIndex}, "StringBackwardsIndex": {funcName: "StringBackwardsIndex", native: divergence_hunt231.StringBackwardsIndex}, "ByteSliceFromString": {funcName: "ByteSliceFromString", native: divergence_hunt231.ByteSliceFromString}, "StringByteAssignment": {funcName: "StringByteAssignment", native: divergence_hunt231.StringByteAssignment}, "StringIndexingWithVariables": {funcName: "StringIndexingWithVariables", native: divergence_hunt231.StringIndexingWithVariables}, "StringFirstLastByte": {funcName: "StringFirstLastByte", native: divergence_hunt231.StringFirstLastByte},
	}})
}
func TestDivergenceHunt232(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt232Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"RangeOverStringASCII": {funcName: "RangeOverStringASCII", native: divergence_hunt232.RangeOverStringASCII}, "RangeOverStringUnicode": {funcName: "RangeOverStringUnicode", native: divergence_hunt232.RangeOverStringUnicode}, "RangeStringIndexValues": {funcName: "RangeStringIndexValues", native: divergence_hunt232.RangeStringIndexValues}, "CountRunesInString": {funcName: "CountRunesInString", native: divergence_hunt232.CountRunesInString}, "RuneSliceFromString": {funcName: "RuneSliceFromString", native: divergence_hunt232.RuneSliceFromString}, "StringFromRuneSlice": {funcName: "StringFromRuneSlice", native: divergence_hunt232.StringFromRuneSlice}, "IterateEmptyString": {funcName: "IterateEmptyString", native: divergence_hunt232.IterateEmptyString}, "UnicodeByteVsRuneCount": {funcName: "UnicodeByteVsRuneCount", native: divergence_hunt232.UnicodeByteVsRuneCount}, "RangeWithEmoji": {funcName: "RangeWithEmoji", native: divergence_hunt232.RangeWithEmoji}, "RuneComparison": {funcName: "RuneComparison", native: divergence_hunt232.RuneComparison}, "RuneArithmetic": {funcName: "RuneArithmetic", native: divergence_hunt232.RuneArithmetic},
	}})
}
func TestDivergenceHunt233(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt233Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringCannotBeModified": {funcName: "StringCannotBeModified", native: divergence_hunt233.StringCannotBeModified}, "StringReassignmentIsAllowed": {funcName: "StringReassignmentIsAllowed", native: divergence_hunt233.StringReassignmentIsAllowed}, "StringConcatenationCreatesNew": {funcName: "StringConcatenationCreatesNew", native: divergence_hunt233.StringConcatenationCreatesNew}, "ByteSliceModifiable": {funcName: "ByteSliceModifiable", native: divergence_hunt233.ByteSliceModifiable}, "StringToByteSliceCopy": {funcName: "StringToByteSliceCopy", native: divergence_hunt233.StringToByteSliceCopy}, "StringPassedByValue": {funcName: "StringPassedByValue", native: divergence_hunt233.StringPassedByValue}, "StringInStructImmutable": {funcName: "StringInStructImmutable", native: divergence_hunt233.StringInStructImmutable}, "StringAppendDoesNotModify": {funcName: "StringAppendDoesNotModify", native: divergence_hunt233.StringAppendDoesNotModify}, "StringConstantImmutability": {funcName: "StringConstantImmutability", native: divergence_hunt233.StringConstantImmutability}, "StringSliceRemainsImmutable": {funcName: "StringSliceRemainsImmutable", native: divergence_hunt233.StringSliceRemainsImmutable}, "MultipleReferencesSameString": {funcName: "MultipleReferencesSameString", native: divergence_hunt233.MultipleReferencesSameString},
	}})
}
func TestDivergenceHunt234(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt234Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicConcatenation": {funcName: "BasicConcatenation", native: divergence_hunt234.BasicConcatenation}, "ConcatenationInLoop": {funcName: "ConcatenationInLoop", native: divergence_hunt234.ConcatenationInLoop}, "ConcatenationWithNumbers": {funcName: "ConcatenationWithNumbers", native: divergence_hunt234.ConcatenationWithNumbers}, "ConcatenationWithBooleans": {funcName: "ConcatenationWithBooleans", native: divergence_hunt234.ConcatenationWithBooleans}, "MixedTypeConcatenation": {funcName: "MixedTypeConcatenation", native: divergence_hunt234.MixedTypeConcatenation}, "StringBuilderPattern": {funcName: "StringBuilderPattern", native: divergence_hunt234.StringBuilderPattern}, "JoinWithSeparator": {funcName: "JoinWithSeparator", native: divergence_hunt234.JoinWithSeparator}, "ConcatenationWithRunes": {funcName: "ConcatenationWithRunes", native: divergence_hunt234.ConcatenationWithRunes}, "EmptyStringConcatenation": {funcName: "EmptyStringConcatenation", native: divergence_hunt234.EmptyStringConcatenation}, "UnicodeConcatenation": {funcName: "UnicodeConcatenation", native: divergence_hunt234.UnicodeConcatenation}, "RepeatedConcatenation": {funcName: "RepeatedConcatenation", native: divergence_hunt234.RepeatedConcatenation},
	}})
}
func TestDivergenceHunt235(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt235Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ByteSliceToString": {funcName: "ByteSliceToString", native: divergence_hunt235.ByteSliceToString}, "StringToByteSlice": {funcName: "StringToByteSlice", native: divergence_hunt235.StringToByteSlice}, "EmptyByteSliceToString": {funcName: "EmptyByteSliceToString", native: divergence_hunt235.EmptyByteSliceToString}, "NilByteSliceToString": {funcName: "NilByteSliceToString", native: divergence_hunt235.NilByteSliceToString}, "ByteSliceModification": {funcName: "ByteSliceModification", native: divergence_hunt235.ByteSliceModification}, "ByteSliceFromStringCopy": {funcName: "ByteSliceFromStringCopy", native: divergence_hunt235.ByteSliceFromStringCopy}, "RuneSliceToString": {funcName: "RuneSliceToString", native: divergence_hunt235.RuneSliceToString}, "StringToRuneSlice": {funcName: "StringToRuneSlice", native: divergence_hunt235.StringToRuneSlice}, "UnicodeByteSlice": {funcName: "UnicodeByteSlice", native: divergence_hunt235.UnicodeByteSlice}, "ByteSliceWithNull": {funcName: "ByteSliceWithNull", native: divergence_hunt235.ByteSliceWithNull}, "ConvertAndBack": {funcName: "ConvertAndBack", native: divergence_hunt235.ConvertAndBack},
	}})
}
func TestDivergenceHunt236(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt236Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringEquality": {funcName: "StringEquality", native: divergence_hunt236.StringEquality}, "StringInequality": {funcName: "StringInequality", native: divergence_hunt236.StringInequality}, "StringLessThan": {funcName: "StringLessThan", native: divergence_hunt236.StringLessThan}, "StringGreaterThan": {funcName: "StringGreaterThan", native: divergence_hunt236.StringGreaterThan}, "StringLessOrEqual": {funcName: "StringLessOrEqual", native: divergence_hunt236.StringLessOrEqual}, "StringCaseSensitivity": {funcName: "StringCaseSensitivity", native: divergence_hunt236.StringCaseSensitivity}, "EmptyStringComparison": {funcName: "EmptyStringComparison", native: divergence_hunt236.EmptyStringComparison}, "StringLengthComparison": {funcName: "StringLengthComparison", native: divergence_hunt236.StringLengthComparison}, "StringLexicographicOrder": {funcName: "StringLexicographicOrder", native: divergence_hunt236.StringLexicographicOrder}, "UnicodeStringComparison": {funcName: "UnicodeStringComparison", native: divergence_hunt236.UnicodeStringComparison}, "StringPrefixComparison": {funcName: "StringPrefixComparison", native: divergence_hunt236.StringPrefixComparison},
	}})
}
func TestDivergenceHunt237(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt237Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"UnicodeStringLength": {funcName: "UnicodeStringLength", native: divergence_hunt237.UnicodeStringLength}, "UnicodeRuneAccess": {funcName: "UnicodeRuneAccess", native: divergence_hunt237.UnicodeRuneAccess}, "UnicodeStringIndexing": {funcName: "UnicodeStringIndexing", native: divergence_hunt237.UnicodeStringIndexing}, "UnicodeRangeLoop": {funcName: "UnicodeRangeLoop", native: divergence_hunt237.UnicodeRangeLoop}, "UnicodeConcatenation": {funcName: "UnicodeConcatenation", native: divergence_hunt237.UnicodeConcatenation}, "MixedASCIIUnicode": {funcName: "MixedASCIIUnicode", native: divergence_hunt237.MixedASCIIUnicode}, "UnicodeComparison": {funcName: "UnicodeComparison", native: divergence_hunt237.UnicodeComparison}, "EmojiStringHandling": {funcName: "EmojiStringHandling", native: divergence_hunt237.EmojiStringHandling}, "UnicodeCaseConversion": {funcName: "UnicodeCaseConversion", native: divergence_hunt237.UnicodeCaseConversion}, "UnicodeByteAccess": {funcName: "UnicodeByteAccess", native: divergence_hunt237.UnicodeByteAccess}, "StringWithCombiningChars": {funcName: "StringWithCombiningChars", native: divergence_hunt237.StringWithCombiningChars},
	}})
}
func TestDivergenceHunt238(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt238Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SprintfBasicVerbs": {funcName: "SprintfBasicVerbs", native: divergence_hunt238.SprintfBasicVerbs}, "SprintfFloatVerbs": {funcName: "SprintfFloatVerbs", native: divergence_hunt238.SprintfFloatVerbs}, "SprintfWidthAndPrecision": {funcName: "SprintfWidthAndPrecision", native: divergence_hunt238.SprintfWidthAndPrecision}, "SprintfFloatPrecision": {funcName: "SprintfFloatPrecision", native: divergence_hunt238.SprintfFloatPrecision}, "SprintfVerbWidthPrecision": {funcName: "SprintfVerbWidthPrecision", native: divergence_hunt238.SprintfVerbWidthPrecision}, "SprintfHexOctalBinary": {funcName: "SprintfHexOctalBinary", native: divergence_hunt238.SprintfHexOctalBinary}, "SprintfPointerVerb": {funcName: "SprintfPointerVerb", native: divergence_hunt238.SprintfPointerVerb}, "SprintfGenericVerb": {funcName: "SprintfGenericVerb", native: divergence_hunt238.SprintfGenericVerb}, "SprintfTypeVerb": {funcName: "SprintfTypeVerb", native: divergence_hunt238.SprintfTypeVerb}, "SprintfPercentEscape": {funcName: "SprintfPercentEscape", native: divergence_hunt238.SprintfPercentEscape}, "SprintfMultipleValues": {funcName: "SprintfMultipleValues", native: divergence_hunt238.SprintfMultipleValues}, "SprintfStringPrecision": {funcName: "SprintfStringPrecision", native: divergence_hunt238.SprintfStringPrecision},
	}})
}
func TestDivergenceHunt239(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt239Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringerWithStruct": {funcName: "StringerWithStruct", native: divergence_hunt239.StringerWithStruct},
		"StringerWithPointerReceiver": {funcName: "StringerWithPointerReceiver", native: divergence_hunt239.StringerWithPointerReceiver},
		"StringerInSlice": {funcName: "StringerInSlice", native: divergence_hunt239.StringerInSlice},
		"StringerWithFmtSprintf": {funcName: "StringerWithFmtSprintf", native: divergence_hunt239.StringerWithFmtSprintf, knownIssue: true},
		"StringerWithFmtV": {funcName: "StringerWithFmtV", native: divergence_hunt239.StringerWithFmtV, knownIssue: true},
		"StringerInterfaceAssertion": {funcName: "StringerInterfaceAssertion", native: divergence_hunt239.StringerInterfaceAssertion},
		"StringerNilReceiver": {funcName: "StringerNilReceiver", native: divergence_hunt239.StringerNilReceiver},
		"StringerNestedCall": {funcName: "StringerNestedCall", native: divergence_hunt239.StringerNestedCall},
		"StringerWithSpecialChars": {funcName: "StringerWithSpecialChars", native: divergence_hunt239.StringerWithSpecialChars},
	}})
}
func TestDivergenceHunt240(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt240Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		// GoString() dispatch via %#v requires side-channel method invocation that
		// currently panics on methods touching struct receiver fields (known deep
		// bug in the method-resolution path). Cases that don't depend on it pass.
		"GoStringerBasic":            {funcName: "GoStringerBasic", native: divergence_hunt240.GoStringerBasic, knownIssue: true},
		"GoStringerStruct":           {funcName: "GoStringerStruct", native: divergence_hunt240.GoStringerStruct, knownIssue: true},
		"GoStringerWithColor":        {funcName: "GoStringerWithColor", native: divergence_hunt240.GoStringerWithColor, knownIssue: true},
		"GoStringerVsStringer":       {funcName: "GoStringerVsStringer", native: divergence_hunt240.GoStringerVsStringer},
		"GoStringerInSlice":          {funcName: "GoStringerInSlice", native: divergence_hunt240.GoStringerInSlice, knownIssue: true},
		"GoStringerNilCheck":         {funcName: "GoStringerNilCheck", native: divergence_hunt240.GoStringerNilCheck, knownIssue: true},
		"GoStringerEmptyStruct":      {funcName: "GoStringerEmptyStruct", native: divergence_hunt240.GoStringerEmptyStruct, knownIssue: true},
		"GoStringerNestedStruct":     {funcName: "GoStringerNestedStruct", native: divergence_hunt240.GoStringerNestedStruct, knownIssue: true},
		"GoStringerWithPointer":      {funcName: "GoStringerWithPointer", native: divergence_hunt240.GoStringerWithPointer},
		"GoStringerFormatComparison": {funcName: "GoStringerFormatComparison", native: divergence_hunt240.GoStringerFormatComparison, knownIssue: true},
	}})
}
func TestDivergenceHunt241(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt241Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NamedReturnDeferModify": {funcName: "NamedReturnDeferModify", native: divergence_hunt241.NamedReturnDeferModify}, "NamedReturnDeferIncrement": {funcName: "NamedReturnDeferIncrement", native: divergence_hunt241.NamedReturnDeferIncrement}, "NamedReturnDeferAdd": {funcName: "NamedReturnDeferAdd", native: divergence_hunt241.NamedReturnDeferAdd}, "NamedReturnDeferMultiply": {funcName: "NamedReturnDeferMultiply", native: divergence_hunt241.NamedReturnDeferMultiply}, "NamedReturnDeferChain": {funcName: "NamedReturnDeferChain", native: divergence_hunt241.NamedReturnDeferChain}, "NamedReturnDeferString": {funcName: "NamedReturnDeferString", native: divergence_hunt241.NamedReturnDeferString}, "NamedReturnDeferSlice": {funcName: "NamedReturnDeferSlice", native: divergence_hunt241.NamedReturnDeferSlice}, "NamedReturnDeferMap": {funcName: "NamedReturnDeferMap", native: divergence_hunt241.NamedReturnDeferMap}, "NamedReturnDeferStruct": {funcName: "NamedReturnDeferStruct", native: divergence_hunt241.NamedReturnDeferStruct}, "NamedReturnDeferPointer": {funcName: "NamedReturnDeferPointer", native: divergence_hunt241.NamedReturnDeferPointer}, "NamedReturnDeferWithPanicRecover": {funcName: "NamedReturnDeferWithPanicRecover", native: divergence_hunt241.NamedReturnDeferWithPanicRecover},
	}})
}
func TestDivergenceHunt242(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt242Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MultipleNamedReturnDefer": {funcName: "MultipleNamedReturnDefer", native: divergence_hunt242.MultipleNamedReturnDefer}, "MultipleNamedReturnDeferOne": {funcName: "MultipleNamedReturnDeferOne", native: divergence_hunt242.MultipleNamedReturnDeferOne}, "MultipleNamedReturnDeferSwap": {funcName: "MultipleNamedReturnDeferSwap", native: divergence_hunt242.MultipleNamedReturnDeferSwap}, "MultipleNamedReturnDeferIncrement": {funcName: "MultipleNamedReturnDeferIncrement", native: divergence_hunt242.MultipleNamedReturnDeferIncrement}, "MultipleNamedReturnDeferMultiply": {funcName: "MultipleNamedReturnDeferMultiply", native: divergence_hunt242.MultipleNamedReturnDeferMultiply}, "MultipleNamedReturnDeferChain": {funcName: "MultipleNamedReturnDeferChain", native: divergence_hunt242.MultipleNamedReturnDeferChain}, "MultipleNamedReturnDeferStringInt": {funcName: "MultipleNamedReturnDeferStringInt", native: divergence_hunt242.MultipleNamedReturnDeferStringInt}, "MultipleNamedReturnDeferBoolInt": {funcName: "MultipleNamedReturnDeferBoolInt", native: divergence_hunt242.MultipleNamedReturnDeferBoolInt}, "MultipleNamedReturnDeferSliceMap": {funcName: "MultipleNamedReturnDeferSliceMap", native: divergence_hunt242.MultipleNamedReturnDeferSliceMap}, "MultipleNamedReturnDeferPanicRecover": {funcName: "MultipleNamedReturnDeferPanicRecover", native: divergence_hunt242.MultipleNamedReturnDeferPanicRecover},
	}})
}
func TestDivergenceHunt243(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt243Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"DeferInForLoop": {funcName: "DeferInForLoop", native: divergence_hunt243.DeferInForLoop}, "DeferInReverseLoop": {funcName: "DeferInReverseLoop", native: divergence_hunt243.DeferInReverseLoop}, "DeferAccumulateInLoop": {funcName: "DeferAccumulateInLoop", native: divergence_hunt243.DeferAccumulateInLoop}, "DeferInNestedLoop": {funcName: "DeferInNestedLoop", native: divergence_hunt243.DeferInNestedLoop}, "DeferWithLoopVariableCapture": {funcName: "DeferWithLoopVariableCapture", native: divergence_hunt243.DeferWithLoopVariableCapture}, "DeferInRangeLoop": {funcName: "DeferInRangeLoop", native: divergence_hunt243.DeferInRangeLoop}, "DeferInRangeLoopIndex": {funcName: "DeferInRangeLoopIndex", native: divergence_hunt243.DeferInRangeLoopIndex}, "DeferInRangeMap": {funcName: "DeferInRangeMap", native: divergence_hunt243.DeferInRangeMap}, "DeferWithBreakLoop": {funcName: "DeferWithBreakLoop", native: divergence_hunt243.DeferWithBreakLoop}, "DeferWithContinueLoop": {funcName: "DeferWithContinueLoop", native: divergence_hunt243.DeferWithContinueLoop}, "DeferConditionalInLoop": {funcName: "DeferConditionalInLoop", native: divergence_hunt243.DeferConditionalInLoop},
	}})
}
func TestDivergenceHunt244(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt244Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicPanicRecover": {funcName: "BasicPanicRecover", native: divergence_hunt244.BasicPanicRecover}, "PanicWithIntRecover": {funcName: "PanicWithIntRecover", native: divergence_hunt244.PanicWithIntRecover}, "PanicWithStructRecover": {funcName: "PanicWithStructRecover", native: divergence_hunt244.PanicWithStructRecover}, "PanicWithPointerRecover": {funcName: "PanicWithPointerRecover", native: divergence_hunt244.PanicWithPointerRecover}, "RecoverReturnsNil": {funcName: "RecoverReturnsNil", native: divergence_hunt244.RecoverReturnsNil}, "MultipleRecoverCalls": {funcName: "MultipleRecoverCalls", native: divergence_hunt244.MultipleRecoverCalls}, "PanicInDeferRecovered": {funcName: "PanicInDeferRecovered", native: divergence_hunt244.PanicInDeferRecovered}, "RepanicAfterRecover": {funcName: "RepanicAfterRecover", native: divergence_hunt244.RepanicAfterRecover}, "PanicWithNilInterface": {funcName: "PanicWithNilInterface", native: divergence_hunt244.PanicWithNilInterface}, "PanicNilLiteral": {funcName: "PanicNilLiteral", native: divergence_hunt244.PanicNilLiteral},
	}})
}
func TestDivergenceHunt245(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt245Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NestedPanicRecoverSimple": {funcName: "NestedPanicRecoverSimple", native: divergence_hunt245.NestedPanicRecoverSimple}, "NestedPanicRecoverBothPanic": {funcName: "NestedPanicRecoverBothPanic", native: divergence_hunt245.NestedPanicRecoverBothPanic}, "TripleNestedPanicRecover": {funcName: "TripleNestedPanicRecover", native: divergence_hunt245.TripleNestedPanicRecover}, "NestedPanicRecoverPartial": {funcName: "NestedPanicRecoverPartial", native: divergence_hunt245.NestedPanicRecoverPartial}, "NestedDeferPanicRecover": {funcName: "NestedDeferPanicRecover", native: divergence_hunt245.NestedDeferPanicRecover}, "NestedPanicRecoverLoop": {funcName: "NestedPanicRecoverLoop", native: divergence_hunt245.NestedPanicRecoverLoop}, "DeepNestedRecoverChain": {funcName: "DeepNestedRecoverChain", native: divergence_hunt245.DeepNestedRecoverChain}, "NestedPanicDifferentTypes": {funcName: "NestedPanicDifferentTypes", native: divergence_hunt245.NestedPanicDifferentTypes}, "NestedPanicNamedReturn": {funcName: "NestedPanicNamedReturn", native: divergence_hunt245.NestedPanicNamedReturn},
	}})
}
func TestDivergenceHunt246(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt246Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BasicErrorWrap": {funcName: "BasicErrorWrap", native: divergence_hunt246.BasicErrorWrap}, "DoubleErrorWrap": {funcName: "DoubleErrorWrap", native: divergence_hunt246.DoubleErrorWrap}, "TripleErrorWrap": {funcName: "TripleErrorWrap", native: divergence_hunt246.TripleErrorWrap}, "ErrorUnwrapChain": {funcName: "ErrorUnwrapChain", native: divergence_hunt246.ErrorUnwrapChain}, "ErrorWrapNil": {funcName: "ErrorWrapNil", native: divergence_hunt246.ErrorWrapNil}, "ErrorUnwrapNonWrapped": {funcName: "ErrorUnwrapNonWrapped", native: divergence_hunt246.ErrorUnwrapNonWrapped}, "ErrorWrapWithFormat": {funcName: "ErrorWrapWithFormat", native: divergence_hunt246.ErrorWrapWithFormat}, "ErrorMultipleWrapSame": {funcName: "ErrorMultipleWrapSame", native: divergence_hunt246.ErrorMultipleWrapSame}, "ErrorWrapWithContext": {funcName: "ErrorWrapWithContext", native: divergence_hunt246.ErrorWrapWithContext}, "ErrorUnwrapNonWrapped2": {funcName: "ErrorUnwrapNonWrapped2", native: divergence_hunt246.ErrorUnwrapNonWrapped2}, "ErrorWrapAndError": {funcName: "ErrorWrapAndError", native: divergence_hunt246.ErrorWrapAndError},
	}})
}
func TestDivergenceHunt247(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt247Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ErrorInterfaceBasic": {funcName: "ErrorInterfaceBasic", native: divergence_hunt247.ErrorInterfaceBasic}, "ErrorInterfaceNil": {funcName: "ErrorInterfaceNil", native: divergence_hunt247.ErrorInterfaceNil}, "ErrorInterfacePointerNil": {funcName: "ErrorInterfacePointerNil", native: divergence_hunt247.ErrorInterfacePointerNil}, "ErrorInterfaceAssignability": {funcName: "ErrorInterfaceAssignability", native: divergence_hunt247.ErrorInterfaceAssignability, knownIssue: true}, "ErrorInterfaceMethodSet": {funcName: "ErrorInterfaceMethodSet", native: divergence_hunt247.ErrorInterfaceMethodSet}, "ErrorInterfaceAssertion": {funcName: "ErrorInterfaceAssertion", native: divergence_hunt247.ErrorInterfaceAssertion}, "ErrorInterfaceEmptyString": {funcName: "ErrorInterfaceEmptyString", native: divergence_hunt247.ErrorInterfaceEmptyString}, "ErrorInterfaceComparison": {funcName: "ErrorInterfaceComparison", native: divergence_hunt247.ErrorInterfaceComparison}, "ErrorInterfaceSlice": {funcName: "ErrorInterfaceSlice", native: divergence_hunt247.ErrorInterfaceSlice}, "ErrorInterfaceMap": {funcName: "ErrorInterfaceMap", native: divergence_hunt247.ErrorInterfaceMap}, "ErrorInterfaceFunction": {funcName: "ErrorInterfaceFunction", native: divergence_hunt247.ErrorInterfaceFunction},
	}})
}
func TestDivergenceHunt248(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt248Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"CustomErrorStruct": {funcName: "CustomErrorStruct", native: divergence_hunt248.CustomErrorStruct},
		"CustomErrorString": {funcName: "CustomErrorString", native: divergence_hunt248.CustomErrorString, knownIssue: true},
		"CustomErrorWithUnwrap": {funcName: "CustomErrorWithUnwrap", native: divergence_hunt248.CustomErrorWithUnwrap, knownIssue: true},
		"CustomErrorIs": {funcName: "CustomErrorIs", native: divergence_hunt248.CustomErrorIs, knownIssue: true},
		"CustomErrorAs": {funcName: "CustomErrorAs", native: divergence_hunt248.CustomErrorAs},
		"CustomErrorHierarchy": {funcName: "CustomErrorHierarchy", native: divergence_hunt248.CustomErrorHierarchy, knownIssue: true},
		"CustomErrorPointerValue": {funcName: "CustomErrorPointerValue", native: divergence_hunt248.CustomErrorPointerValue},
		"CustomErrorNilPointer": {funcName: "CustomErrorNilPointer", native: divergence_hunt248.CustomErrorNilPointer, knownIssue: true},
		"CustomErrorComposite": {funcName: "CustomErrorComposite", native: divergence_hunt248.CustomErrorComposite},
		"CustomErrorState": {funcName: "CustomErrorState", native: divergence_hunt248.CustomErrorState},
	}})
}
func TestDivergenceHunt249(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt249Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ErrorsIsBasic": {funcName: "ErrorsIsBasic", native: divergence_hunt249.ErrorsIsBasic},
		"ErrorsIsDeep": {funcName: "ErrorsIsDeep", native: divergence_hunt249.ErrorsIsDeep},
		"ErrorsIsNotFound": {funcName: "ErrorsIsNotFound", native: divergence_hunt249.ErrorsIsNotFound},
		"ErrorsIsNilTarget": {funcName: "ErrorsIsNilTarget", native: divergence_hunt249.ErrorsIsNilTarget},
		"ErrorsAsBasic": {funcName: "ErrorsAsBasic", native: divergence_hunt249.ErrorsAsBasic},
		"ErrorsAsWrapped": {funcName: "ErrorsAsWrapped", native: divergence_hunt249.ErrorsAsWrapped},
		"ErrorsAsInterface": {funcName: "ErrorsAsInterface", native: divergence_hunt249.ErrorsAsInterface, knownIssue: true},
		"ErrorsAsNotMatching": {funcName: "ErrorsAsNotMatching", native: divergence_hunt249.ErrorsAsNotMatching, knownIssue: true},
		"ErrorsIsSentinel": {funcName: "ErrorsIsSentinel", native: divergence_hunt249.ErrorsIsSentinel},
		"ErrorsIsAndAsTogether": {funcName: "ErrorsIsAndAsTogether", native: divergence_hunt249.ErrorsIsAndAsTogether, knownIssue: true},
		"ErrorsAsPointerValue": {funcName: "ErrorsAsPointerValue", native: divergence_hunt249.ErrorsAsPointerValue, knownIssue: true},
	}})
}
func TestDivergenceHunt250(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt250Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"JoinTwoErrors": {funcName: "JoinTwoErrors", native: divergence_hunt250.JoinTwoErrors, knownIssue: true},
		"JoinMultipleErrors": {funcName: "JoinMultipleErrors", native: divergence_hunt250.JoinMultipleErrors, knownIssue: true},
		"JoinWithNilErrors": {funcName: "JoinWithNilErrors", native: divergence_hunt250.JoinWithNilErrors, knownIssue: true},
		"JoinAllNilErrors": {funcName: "JoinAllNilErrors", native: divergence_hunt250.JoinAllNilErrors, knownIssue: true},
		"JoinErrorsIs": {funcName: "JoinErrorsIs", native: divergence_hunt250.JoinErrorsIs, knownIssue: true},
		"JoinErrorsAs": {funcName: "JoinErrorsAs", native: divergence_hunt250.JoinErrorsAs, knownIssue: true},
		"CollectErrorsAccumulates": {funcName: "CollectErrorsAccumulates", native: divergence_hunt250.CollectErrorsAccumulates, knownIssue: true},
		"FirstNonNilError": {funcName: "FirstNonNilError", native: divergence_hunt250.FirstNonNilError, knownIssue: true},
		"CombineErrorsInLoop": {funcName: "CombineErrorsInLoop", native: divergence_hunt250.CombineErrorsInLoop, knownIssue: true},
		"ErrorSliceToJoined": {funcName: "ErrorSliceToJoined", native: divergence_hunt250.ErrorSliceToJoined, knownIssue: true},
		"ProcessMultipleErrors": {funcName: "ProcessMultipleErrors", native: divergence_hunt250.ProcessMultipleErrors, knownIssue: true},
	}})
}
func TestDivergenceHunt251(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt251Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"IotaBasic": {funcName: "IotaBasic", native: divergence_hunt251.IotaBasic}, "IotaReset": {funcName: "IotaReset", native: divergence_hunt251.IotaReset}, "IotaBitshift": {funcName: "IotaBitshift", native: divergence_hunt251.IotaBitshift}, "IotaSkip": {funcName: "IotaSkip", native: divergence_hunt251.IotaSkip}, "IotaExpression": {funcName: "IotaExpression", native: divergence_hunt251.IotaExpression}, "IotaWithOffset": {funcName: "IotaWithOffset", native: divergence_hunt251.IotaWithOffset}, "IotaStringEnum": {funcName: "IotaStringEnum", native: divergence_hunt251.IotaStringEnum}, "IotaDaysOfWeek": {funcName: "IotaDaysOfWeek", native: divergence_hunt251.IotaDaysOfWeek}, "IotaWithParentheses": {funcName: "IotaWithParentheses", native: divergence_hunt251.IotaWithParentheses}, "IotaNegativeStep": {funcName: "IotaNegativeStep", native: divergence_hunt251.IotaNegativeStep}, "IotaMultiplication": {funcName: "IotaMultiplication", native: divergence_hunt251.IotaMultiplication},
	}})
}
func TestDivergenceHunt252(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt252Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"UntypedIntAssign": {funcName: "UntypedIntAssign", native: divergence_hunt252.UntypedIntAssign}, "UntypedFloatAssign": {funcName: "UntypedFloatAssign", native: divergence_hunt252.UntypedFloatAssign}, "TypedIntAssign": {funcName: "TypedIntAssign", native: divergence_hunt252.TypedIntAssign}, "UntypedStringAssign": {funcName: "UntypedStringAssign", native: divergence_hunt252.UntypedStringAssign}, "UntypedBoolAssign": {funcName: "UntypedBoolAssign", native: divergence_hunt252.UntypedBoolAssign}, "UntypedArithmetic": {funcName: "UntypedArithmetic", native: divergence_hunt252.UntypedArithmetic}, "TypedArithmetic": {funcName: "TypedArithmetic", native: divergence_hunt252.TypedArithmetic}, "MixedTypeArithmetic": {funcName: "MixedTypeArithmetic", native: divergence_hunt252.MixedTypeArithmetic}, "DefaultTypeInference": {funcName: "DefaultTypeInference", native: divergence_hunt252.DefaultTypeInference}, "LargeUntypedInteger": {funcName: "LargeUntypedInteger", native: divergence_hunt252.LargeUntypedInteger}, "UntypedRune": {funcName: "UntypedRune", native: divergence_hunt252.UntypedRune},
	}})
}
func TestDivergenceHunt253(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt253Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ConstAdd": {funcName: "ConstAdd", native: divergence_hunt253.ConstAdd}, "ConstSub": {funcName: "ConstSub", native: divergence_hunt253.ConstSub}, "ConstMul": {funcName: "ConstMul", native: divergence_hunt253.ConstMul}, "ConstDiv": {funcName: "ConstDiv", native: divergence_hunt253.ConstDiv}, "ConstMod": {funcName: "ConstMod", native: divergence_hunt253.ConstMod}, "ConstNeg": {funcName: "ConstNeg", native: divergence_hunt253.ConstNeg}, "ConstBitwiseAnd": {funcName: "ConstBitwiseAnd", native: divergence_hunt253.ConstBitwiseAnd}, "ConstBitwiseOr": {funcName: "ConstBitwiseOr", native: divergence_hunt253.ConstBitwiseOr}, "ConstBitwiseXor": {funcName: "ConstBitwiseXor", native: divergence_hunt253.ConstBitwiseXor}, "ConstBitShiftLeft": {funcName: "ConstBitShiftLeft", native: divergence_hunt253.ConstBitShiftLeft}, "ConstBitShiftRight": {funcName: "ConstBitShiftRight", native: divergence_hunt253.ConstBitShiftRight}, "ConstComplexExpr": {funcName: "ConstComplexExpr", native: divergence_hunt253.ConstComplexExpr},
	}})
}
func TestDivergenceHunt254(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt254Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"OverflowInt8Max": {funcName: "OverflowInt8Max", native: divergence_hunt254.OverflowInt8Max}, "OverflowInt8Min": {funcName: "OverflowInt8Min", native: divergence_hunt254.OverflowInt8Min}, "OverflowUint8Max": {funcName: "OverflowUint8Max", native: divergence_hunt254.OverflowUint8Max}, "OverflowInt16Max": {funcName: "OverflowInt16Max", native: divergence_hunt254.OverflowInt16Max}, "OverflowInt32Max": {funcName: "OverflowInt32Max", native: divergence_hunt254.OverflowInt32Max}, "OverflowUint32Max": {funcName: "OverflowUint32Max", native: divergence_hunt254.OverflowUint32Max}, "OverflowInt64Max": {funcName: "OverflowInt64Max", native: divergence_hunt254.OverflowInt64Max}, "OverflowUint64Max": {funcName: "OverflowUint64Max", native: divergence_hunt254.OverflowUint64Max}, "OverflowFloat32Max": {funcName: "OverflowFloat32Max", native: divergence_hunt254.OverflowFloat32Max}, "LargeShift": {funcName: "LargeShift", native: divergence_hunt254.LargeShift}, "UintOverflowBoundary": {funcName: "UintOverflowBoundary", native: divergence_hunt254.UintOverflowBoundary},
	}})
}
func TestDivergenceHunt255(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt255Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ShadowSimple": {funcName: "ShadowSimple", native: divergence_hunt255.ShadowSimple}, "ShadowInLoop": {funcName: "ShadowInLoop", native: divergence_hunt255.ShadowInLoop}, "ShadowIfStatement": {funcName: "ShadowIfStatement", native: divergence_hunt255.ShadowIfStatement}, "ShadowSwitchStatement": {funcName: "ShadowSwitchStatement", native: divergence_hunt255.ShadowSwitchStatement}, "ShadowForRange": {funcName: "ShadowForRange", native: divergence_hunt255.ShadowForRange}, "ShadowMultipleScopes": {funcName: "ShadowMultipleScopes", native: divergence_hunt255.ShadowMultipleScopes}, "ShadowFunctionParam": {funcName: "ShadowFunctionParam", native: divergence_hunt255.ShadowFunctionParam, knownIssue: true}, "ShadowWithShortDecl": {funcName: "ShadowWithShortDecl", native: divergence_hunt255.ShadowWithShortDecl}, "ShadowGlobal": {funcName: "ShadowGlobal", native: divergence_hunt255.ShadowGlobal}, "ShadowSameType": {funcName: "ShadowSameType", native: divergence_hunt255.ShadowSameType}, "ShadowDifferentType": {funcName: "ShadowDifferentType", native: divergence_hunt255.ShadowDifferentType},
	}})
}
func TestDivergenceHunt256(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt256Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ShortDeclBasic": {funcName: "ShortDeclBasic", native: divergence_hunt256.ShortDeclBasic}, "ShortDeclMultiple": {funcName: "ShortDeclMultiple", native: divergence_hunt256.ShortDeclMultiple}, "ShortDeclMixedTypes": {funcName: "ShortDeclMixedTypes", native: divergence_hunt256.ShortDeclMixedTypes}, "ShortDeclRedefine": {funcName: "ShortDeclRedefine", native: divergence_hunt256.ShortDeclRedefine}, "ShortDeclInIf": {funcName: "ShortDeclInIf", native: divergence_hunt256.ShortDeclInIf}, "ShortDeclInFor": {funcName: "ShortDeclInFor", native: divergence_hunt256.ShortDeclInFor}, "ShortDeclInSwitch": {funcName: "ShortDeclInSwitch", native: divergence_hunt256.ShortDeclInSwitch}, "ShortDeclSlice": {funcName: "ShortDeclSlice", native: divergence_hunt256.ShortDeclSlice}, "ShortDeclFunctionCall": {funcName: "ShortDeclFunctionCall", native: divergence_hunt256.ShortDeclFunctionCall}, "ShortDeclTypeInference": {funcName: "ShortDeclTypeInference", native: divergence_hunt256.ShortDeclTypeInference}, "ShortDeclBlank": {funcName: "ShortDeclBlank", native: divergence_hunt256.ShortDeclBlank}, "ShortDeclComposite": {funcName: "ShortDeclComposite", native: divergence_hunt256.ShortDeclComposite},
	}})
}
func TestDivergenceHunt257(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt257Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"MultiAssignBasic": {funcName: "MultiAssignBasic", native: divergence_hunt257.MultiAssignBasic}, "MultiAssignSwap": {funcName: "MultiAssignSwap", native: divergence_hunt257.MultiAssignSwap}, "MultiAssignTripleSwap": {funcName: "MultiAssignTripleSwap", native: divergence_hunt257.MultiAssignTripleSwap}, "MultiAssignFunctionReturn": {funcName: "MultiAssignFunctionReturn", native: divergence_hunt257.MultiAssignFunctionReturn}, "MultiAssignMapLookup": {funcName: "MultiAssignMapLookup", native: divergence_hunt257.MultiAssignMapLookup}, "MultiAssignTypeAssertion": {funcName: "MultiAssignTypeAssertion", native: divergence_hunt257.MultiAssignTypeAssertion}, "MultiAssignChannelRecv": {funcName: "MultiAssignChannelRecv", native: divergence_hunt257.MultiAssignChannelRecv}, "MultiAssignSameValue": {funcName: "MultiAssignSameValue", native: divergence_hunt257.MultiAssignSameValue}, "MultiAssignExpression": {funcName: "MultiAssignExpression", native: divergence_hunt257.MultiAssignExpression}, "MultiAssignArrayIndex": {funcName: "MultiAssignArrayIndex", native: divergence_hunt257.MultiAssignArrayIndex}, "MultiAssignMixedTypes": {funcName: "MultiAssignMixedTypes", native: divergence_hunt257.MultiAssignMixedTypes}, "MultiAssignReassign": {funcName: "MultiAssignReassign", native: divergence_hunt257.MultiAssignReassign},
	}})
}
func TestDivergenceHunt258(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt258Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"BlankInShortDecl": {funcName: "BlankInShortDecl", native: divergence_hunt258.BlankInShortDecl}, "BlankInMultiAssign": {funcName: "BlankInMultiAssign", native: divergence_hunt258.BlankInMultiAssign}, "BlankInRange": {funcName: "BlankInRange", native: divergence_hunt258.BlankInRange}, "BlankInRangeIndex": {funcName: "BlankInRangeIndex", native: divergence_hunt258.BlankInRangeIndex}, "BlankInMapRange": {funcName: "BlankInMapRange", native: divergence_hunt258.BlankInMapRange}, "BlankInTypeAssertion": {funcName: "BlankInTypeAssertion", native: divergence_hunt258.BlankInTypeAssertion}, "BlankInChannelRecv": {funcName: "BlankInChannelRecv", native: divergence_hunt258.BlankInChannelRecv}, "BlankInFuncParams": {funcName: "BlankInFuncParams", native: divergence_hunt258.BlankInFuncParams}, "BlankMultiple": {funcName: "BlankMultiple", native: divergence_hunt258.BlankMultiple}, "BlankInImport": {funcName: "BlankInImport", native: divergence_hunt258.BlankInImport}, "BlankWithSideEffects": {funcName: "BlankWithSideEffects", native: divergence_hunt258.BlankWithSideEffects},
	}})
}
func TestDivergenceHunt259(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt259Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InitOrder": {funcName: "InitOrder", native: divergence_hunt259.InitOrder, knownIssue: true},
		"InitWithVariable": {funcName: "InitWithVariable", native: divergence_hunt259.InitWithVariable, knownIssue: true},
		"InitWithMap": {funcName: "InitWithMap", native: divergence_hunt259.InitWithMap, knownIssue: true},
		"InitWithSlice": {funcName: "InitWithSlice", native: divergence_hunt259.InitWithSlice, knownIssue: true},
		"InitWithStruct": {funcName: "InitWithStruct", native: divergence_hunt259.InitWithStruct, knownIssue: true},
		"InitMultipleVars": {funcName: "InitMultipleVars", native: divergence_hunt259.InitMultipleVars, knownIssue: true},
		"InitComplexSetup": {funcName: "InitComplexSetup", native: divergence_hunt259.InitComplexSetup, knownIssue: true},
		"InitCounter": {funcName: "InitCounter", native: divergence_hunt259.InitCounter, knownIssue: true},
		"InitWithArray": {funcName: "InitWithArray", native: divergence_hunt259.InitWithArray, knownIssue: true},
		"InitConditional": {funcName: "InitConditional", native: divergence_hunt259.InitConditional, knownIssue: true},
		"InitLoop": {funcName: "InitLoop", native: divergence_hunt259.InitLoop, knownIssue: true},
	}})
}
func TestDivergenceHunt260(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt260Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"InitOrderDependency": {funcName: "InitOrderDependency", native: divergence_hunt260.InitOrderDependency, knownIssue: true},
		"InitOrderReverse": {funcName: "InitOrderReverse", native: divergence_hunt260.InitOrderReverse, knownIssue: true},
		"InitWithFunctionCall": {funcName: "InitWithFunctionCall", native: divergence_hunt260.InitWithFunctionCall, knownIssue: true},
		"InitWithExpression": {funcName: "InitWithExpression", native: divergence_hunt260.InitWithExpression, knownIssue: true},
		"InitWithStringConcat": {funcName: "InitWithStringConcat", native: divergence_hunt260.InitWithStringConcat, knownIssue: true},
		"InitWithSliceLiteral": {funcName: "InitWithSliceLiteral", native: divergence_hunt260.InitWithSliceLiteral, knownIssue: true},
		"InitWithMapLiteral": {funcName: "InitWithMapLiteral", native: divergence_hunt260.InitWithMapLiteral, knownIssue: true},
		"InitWithStructLiteral": {funcName: "InitWithStructLiteral", native: divergence_hunt260.InitWithStructLiteral, knownIssue: true},
		"InitWithArrayLiteral": {funcName: "InitWithArrayLiteral", native: divergence_hunt260.InitWithArrayLiteral, knownIssue: true},
		"InitMixedTypes": {funcName: "InitMixedTypes", native: divergence_hunt260.InitMixedTypes, knownIssue: true},
		"InitChainDependency": {funcName: "InitChainDependency", native: divergence_hunt260.InitChainDependency, knownIssue: true},
	}})
}
