package tattle

// NOTE:
// You may need to adjust the Output stanza if line numbering changes.

import (
	"fmt"
)

func ExampleSetFrames() {
	// This example is also the unit test for SetFrame.
	// The backtrace would happily flow into the test infrastructure.
	// The 2 function embedding is to ensure reproducibility.
	func() {
		func() {
			def := SetFrames(-100)
			fmt.Printf("SetFrame default: %d\n", def)
			lowStop := SetFrames(100000)
			fmt.Printf("SetFrame low stop: %d\n", lowStop)
			hiStop := SetFrames(0)
			fmt.Printf("SetFrame hi stop:  %d\n", hiStop)

			var tat = tattler{}
			SetFrames(0)
			tat.Latchf("ERROR")
			fmt.Printf("Example error, no trace back:%s", tat.String())

			tat = tattler{}
			SetFrames(1)
			tat.Latchf("ERROR")
			fmt.Printf("Example error, 1 frame trace back: %s", tat.String())

			tat = tattler{}
			SetFrames(def)
			tat.Latchf("ERROR")
			fmt.Printf("Example error, default trace back: %s", tat.String())

		}()
	}()
	// Output:
	// SetFrame default: 3
	// SetFrame low stop: 0
	// SetFrame hi stop:  100
	// Example error, no trace back:ERROR
	// Example error, 1 frame trace back: ERROR
	//  Latched at:  exampleSetFrames_test.go:30 in github.com/rsnorthscope/tattle.ExampleSetFrames.ExampleSetFrames.func1.func2
	// Example error, default trace back: ERROR
	//  Latched at:  exampleSetFrames_test.go:35 in github.com/rsnorthscope/tattle.ExampleSetFrames.ExampleSetFrames.func1.func2
	//  Called From: exampleSetFrames_test.go:38 in github.com/rsnorthscope/tattle.ExampleSetFrames.func1
	//  Called From: exampleSetFrames_test.go:39 in github.com/rsnorthscope/tattle.ExampleSetFrames
}
