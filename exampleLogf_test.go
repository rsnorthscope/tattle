package tattle

import (
	"fmt"
	"log"
	"strings"
)

// NOTE:
// You may need to adjust the Output stanza if line numbering changes.

func ExampleTattler_Logf() {
	type Record struct {
		ID   int
		Name string
		tat  tattler // assuming 'type tattler tattle.Tattler'
	}
	tr := &Record{ID: 15}
	// An example function call sequence follows, similar to what
	// you might see among a series of method calls for Record.  The
	// error will be detected in the level_2function.
	var level_2, level_1 func() error
	level_1 = func() error {
		// This function (level_1) will show a full set of tat interactions.
		// First up is to queue a log of any latches that occur in the
		// function.
		defer tr.tat.Log()
		// If the tat is already latched, then the record is tainted. Quite
		// frequently the best choice is to give up now.
		if tr.tat.Led() { // If already tattled,
			return tr.tat.Le() // return the tattle
		}
		// OK, record is currently clean.  Following runs level_2, latching
		// any error, and returning immediately if an error is latched.
		if tr.tat.Latch(level_2()) { // line 35, call to level_2
			return tr.tat.Le()
		}
		// Any additional checks here
		return tr.tat.Le()
	}
	// Level 2 is slightly abbreviated.
	level_2 = func() error {
		// Use of Logf to indicate which record is causing the error
		defer tr.tat.Logf("tr record %d:", tr.ID)
		// In this example we'll be checking the "Name" string in record ID 15
		// for an an imposed max length of 1000 bytes, say out of concern that
		// such a thing could be an exploit attempt.
		if len(tr.Name) > 1000 {
			tr.tat.Latchf("name size %d exceeds max 1000", len(tr.Name)) // line 49
		}
		return tr.tat.Le() // Return latched error.  Logf logs in any event.
	}
	tr.ID = 15
	// make a big name that looks like exploit exploit exploit.....
	tr.Name = "exploit "
	for i := 0; i < 10; i++ {
		tr.Name = tr.Name + tr.Name
	}
	// For the purposes of this test example, capture the log output as a
	// string.
	sb := &strings.Builder{}
	l := log.Writer()
	log.SetOutput(sb)
	// Now run the test
	level_1() // Line 65
	// reconnect the logger output for subsequent tests.
	log.SetOutput(l)
	// Replace actual date and time with fixed value
	fmt.Print("2024/02/18 15:50:30" + sb.String()[19:])
	// Output:
	// 2024/02/18 15:50:30 tr record 15:name size 8192 exceeds max 1000
	//  Latched at:  exampleLogf_test.go:49 in github.com/rsnorthscope/tattle.ExampleTattler_Logf.func2
	//  Called From: exampleLogf_test.go:35 in github.com/rsnorthscope/tattle.ExampleTattler_Logf.func1
	//  Called From: exampleLogf_test.go:65 in github.com/rsnorthscope/tattle.ExampleTattler_Logf
}
