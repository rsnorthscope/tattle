package tattle

import (
	"fmt"
	"log"
	"strings"
)

// NOTE:
// You may need to adjust the Output stanza if line numbering changes.

func ExampleTattler_Logf() {
	// This example focuses on Logf.  See full example after
	// the intro to see the typical template of a method.
	type Record struct {
		ID   int
		Name string
		tat  tattler // assuming 'type tattler tattle.Tattler'
	}
	rp := &Record{ID: 15}
	rp.Name = "exploit "
	for i := 0; i < 10; i++ {
		rp.Name = rp.Name + rp.Name
	}
	// For the purposes of this test example, capture the log output as a
	// string.
	sb := &strings.Builder{}
	l := log.Writer()
	log.SetOutput(sb)
	// Now run the test
	func() {
		func() error {
			// Use of Logf to indicate which record is causing the error
			defer rp.tat.Logf("tr record %d:", rp.ID)
			if len(rp.Name) > 1000 {
				rp.tat.Latchf("name size %d exceeds max 1000", len(rp.Name)) // line 36
			}
			return rp.tat.Le() // Return latched error.  Logf logs in any event.
		}() // Line 39
	}() // Line 40
	// reconnect the logger output for subsequent tests.
	log.SetOutput(l)
	// Replace actual date and time with fixed value
	fmt.Print("2024/02/18 15:50:30" + sb.String()[19:])
	// Output:
	// 2024/02/18 15:50:30 tr record 15:name size 8192 exceeds max 1000
	//  Latched at:  exampleLogf_test.go:36 in github.com/rsnorthscope/tattle.ExampleTattler_Logf.ExampleTattler_Logf.func1.func2
	//  Called From: exampleLogf_test.go:39 in github.com/rsnorthscope/tattle.ExampleTattler_Logf.func1
	//  Called From: exampleLogf_test.go:40 in github.com/rsnorthscope/tattle.ExampleTattler_Logf

}
