// Copyright 2024 Richard Northscope.  All rights reserved.
// Use of this source code is governed by the
// MIT license that can be found in the LICENSE file.
package tattle

// NOTE:
// You may need to adjust the Output stanza if line numbering changes.

import "fmt"

func ExampleTattler_String() {
	var tr struct {
		tat tattler // assuming 'type tattler tattle.Tattler'
		//other fields as needed
	}
	func() {
		func() {
			tr.tat.Latchf("Example Error") // Line 18
			fmt.Print(tr.tat.String())     // Line 19
		}() // Line 20
	}() // Line 21

	// Output:
	// Example Error
	//  Latched at:  exampleString_test.go:18 in github.com/rsnorthscope/tattle.ExampleTattler_String.ExampleTattler_String.func1.func2
	//  Called From: exampleString_test.go:20 in github.com/rsnorthscope/tattle.ExampleTattler_String.func1
	//  Called From: exampleString_test.go:21 in github.com/rsnorthscope/tattle.ExampleTattler_String

}
