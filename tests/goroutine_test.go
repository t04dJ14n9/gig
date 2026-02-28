package tests

import (
	_ "embed"
	"testing"

	"gig"
	_ "gig/stdlib/packages"
	"gig/tests/testdata/goroutine"
)

//go:embed testdata/goroutine/main.go
var goroutineSrc string

// goroutineTests maps test names to native functions for comparison.
var goroutineTests = map[string]testCase{
	"goroutine/BasicSpawn":                  {goroutineSrc, "BasicSpawn", func() any { return goroutine.BasicSpawn() }},
	"goroutine/ChannelCommunication":        {goroutineSrc, "ChannelCommunication", func() any { return goroutine.ChannelCommunication() }},
	"goroutine/WithArguments":               {goroutineSrc, "WithArguments", func() any { return goroutine.WithArguments() }},
	"goroutine/WithStruct":                  {goroutineSrc, "WithStruct", func() any { return goroutine.WithStruct() }},
	"goroutine/DifferentTypes":              {goroutineSrc, "DifferentTypes", func() any { return goroutine.DifferentTypes() }},
	"goroutine/GlobalsSharing":              {goroutineSrc, "GlobalsSharing", func() any { return goroutine.GlobalsSharing() }},
	"goroutine/MultipleSends":               {goroutineSrc, "MultipleSends", func() any { return goroutine.MultipleSends() }},
	"goroutine/ParallelExecution":           {goroutineSrc, "ParallelExecution", func() any { return goroutine.ParallelExecution() }},
	"goroutine/ClosureCapture":              {goroutineSrc, "ClosureCapture", func() any { return goroutine.ClosureCapture() }},
	"goroutine/ClosureCaptureMultiple":      {goroutineSrc, "ClosureCaptureMultiple", func() any { return goroutine.ClosureCaptureMultiple() }},
	"goroutine/SelectStatement":             {goroutineSrc, "SelectStatement", func() any { return goroutine.SelectStatement() }},
	"goroutine/SelectDefault":               {goroutineSrc, "SelectDefault", func() any { return goroutine.SelectDefault() }},
	"goroutine/SelectSend":                  {goroutineSrc, "SelectSend", func() any { return goroutine.SelectSend() }},
	"goroutine/RangeOverChannel":            {goroutineSrc, "RangeOverChannel", func() any { return goroutine.RangeOverChannel() }},
	"goroutine/RangeOverChannelWithBuiltin": {goroutineSrc, "RangeOverChannelWithBuiltin", func() any { return goroutine.RangeOverChannelWithBuiltin() }},
}

// TestGoroutine runs all goroutine tests using the same pattern as TestAllStdlib.
func TestGoroutine(t *testing.T) {
	for name, tc := range goroutineTests {
		t.Run(name, func(t *testing.T) {
			src := toMainPackage(tc.src)
			prog, err := gig.Build(src)
			if err != nil {
				t.Fatalf("Build error: %v", err)
			}

			result, err := prog.Run(tc.funcName)
			if err != nil {
				t.Fatalf("Run error: %v", err)
			}

			expected := tc.native()
			compareResults(t, result, expected)
		})
	}
}
