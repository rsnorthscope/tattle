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
			tr.tat.Latchf("Example Error") // Line 15
			fmt.Print(tr.tat.String())     // Line 16
		}() // Line 17
	}() // Line 18

	// Output:
	// Example Error
	//  Latched at:  ExampleString_test.go:15 in github.com/rsnorthscope/tattle.ExampleTattler_String.ExampleTattler_String.func1.func2
	//  Called From: ExampleString_test.go:17 in github.com/rsnorthscope/tattle.ExampleTattler_String.func1
	//  Called From: ExampleString_test.go:18 in github.com/rsnorthscope/tattle.ExampleTattler_String

}
