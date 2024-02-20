// Copyright 2024 Richard Northscope.  All rights reserved.
// Use of this source code is governed by the
// MIT license that can be found in the LICENSE file.

package tattle

// NOTE:
// You may need to adjust the Output stanza if line numbering changes.

import (
	"fmt"
)

func ExampleSetFrames() {
	// The backtrace would happily flow into the test infrastructure.
	// The 2 function embedding is to ensure reproducibility.
	func() {
		func() {
			tat := tattler{}
			tat.Latchf("ERROR")
			fmt.Printf("Example error, default trace back: %s", tat.String())

			tat = tattler{}
			SetFrames(0)
			tat.Latchf("ERROR")
			fmt.Printf("Example error, no trace back:%s", tat.String())

			tat = tattler{}
			SetFrames(1)
			tat.Latchf("ERROR")
			fmt.Printf("Example error, 1 frame trace back: %s", tat.String())

			SetFrames(3)
		}()
	}()
	// Output:
	// Example error, default trace back: ERROR
	//  Latched at:  exampleSetFrames_test.go:20 in github.com/rsnorthscope/tattle.ExampleSetFrames.ExampleSetFrames.func1.func2
	//  Called From: exampleSetFrames_test.go:34 in github.com/rsnorthscope/tattle.ExampleSetFrames.func1
	//  Called From: exampleSetFrames_test.go:35 in github.com/rsnorthscope/tattle.ExampleSetFrames
	// Example error, no trace back:ERROR
	// Example error, 1 frame trace back: ERROR
	//  Latched at:  exampleSetFrames_test.go:30 in github.com/rsnorthscope/tattle.ExampleSetFrames.ExampleSetFrames.func1.func2
}
