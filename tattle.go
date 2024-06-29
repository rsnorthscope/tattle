// Copyright 2019-2024 Richard Northscope.  All rights reserved.
// Use of this source code is governed by the
// MIT license that can be found in the LICENSE file.

/*
# Motivation

Go has great error handling. Paying attention to the errors can be a
challenge. This package helps you pay attention to the errors without
changing or extending the error concept.

A Tattler is a lightweight object that records, or Latches, the first error
it is given. Tattler defines a mechanism to log the error, including the
source code location where the error was latched.

Tattlers aren't pretty, but they can save you a lot of debugging time:

 1. Include a Tattler variable, which I usually call 'tat', in types
    (structures) that have complicated methods.
 2. Defer a Tattle Log() call in every method, so errors get logged in close
    proximity to their occurrence.
 3. Latch errors in the methods to the tat.
 4. If the tat latches with an unexpected error (n.b., most of them), treat
    the enclosing structure instance as tainted and stop doing anything with
    it. Simply return if called with the tat latched.

An example of the above follows this introduction.

In the debugging phase, one look at the log, plus knowledge of recent
changes, is enough to trigger that "oh shoot" moment when you realize what
you must do to address the error. Later, after you've forgotten what the
code does, a Tattler latch is still useful but it's far less satisfying.

A secondary use of a tat is to monitor a complicated sequence that is broken
into sub-functions, where an error doesn't call into question any particular
structure. In this case it may make sense to use a tat variable as a local
variable, and pass a pointer to it to the subfunctions, each of which should
defer a tat Log() call. and latch errors to the tat.  This is helpful but
less common scenario.

Full disclosure, tattlers aren't good for dealing with end-user errors.
End-user errors generally aren't inter-related so there is no benefit to
catching the first one.

Here is a tattler log message from an example included for the Tattler Logf
function. The messages will be considerably shorter for local packages :-) :

	2024/02/18 15:50:30 tr record 15:name size 8192 exceeds max 1000
	 Latched at:  exampleLogf_test.go:49 in github.com/rsnorthscope/tattle.ExampleTattler_Logf.func2
	 Called From: exampleLogf_test.go:35 in github.com/rsnorthscope/tattle.ExampleTattler_Logf.func1
	 Called From: exampleLogf_test.go:65 in github.com/rsnorthscope/tattle.ExampleTattler_Logf

# Implementation

Prior to an error, a tattler is a nil pointer to a concrete type known as
the tale.

When the Latch methods are given an error to latch, they allocate the tale
to store the error and 3 frames from the stack.  Only the first non-nil
error is latched; subsequent different non-nil errors are counted but the
error value is discarded.

The log functions are low cost in the non-error case.  The printf-style
string in a Logf isn't expanded unless there is an error.

A tattler doesn't implement the error interface. I tried it and decided it
was a mistake.

# Concurrency

Concurrency for tattlers present the same challenge as concurrency for error
variables.  In the above-mentioned project, the concurrency mechanisms for
the structures also protect the tattlers.

Lacking a use case, I haven't defined a "SerialTattler" type.

# Final Thoughts

  - If you have a Tattler, latch (almost) every error.
  - FWIW the most-used methods are Latch(), Latchf(), Le(), and Logf().
  - A deferred Log()/Logf() is your friend.
*/
package tattle

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

// A Tattler is used to record an error in a structure or within a call flow.
type Tattler struct {
	talep *tale
}

// A tale exists only for Tattlers that have latched an error
type tale struct {
	latched error
	frames  []runtime.Frame
	logged  bool
	missed  int
}

// fullLatch contains the latch logic.
// parameter b is the difference in frames
// between Latch and fulLatch, normally 1.
func (tat *Tattler) fullLatch(b int, e error) bool {
	if e != nil {
		if tat.talep == nil {
			tp := new(tale)
			tat.talep = tp
			tp.latched = e

			// Capture backtrace frames.  The following
			// make defines how many frames appear.
			pc := make([]uintptr, 3)
			n := runtime.Callers(b+2, pc)
			pc = pc[:n] // truncate invalid entries

			var frame runtime.Frame
			frames := runtime.CallersFrames(pc)
			for more := true; more; {
				frame, more = frames.Next()
				tp.frames = append(tp.frames, frame)
			}
		} else {
			if e != tat.talep.latched {
				tat.talep.missed++
			}
		}
	}
	return tat.talep != nil
}

// Import has semantics similar to Latch, except that it's argument
// is a Tattler not an error.  Import is rarely used.
// An example use case is to taint a container
// with a tattle from an enclosed structure.
func (tat *Tattler) Import(itp *Tattler) bool {
	if itp.talep != nil && tat.talep == nil {
		tat.talep = new(tale)
		*tat.talep = *itp.talep
	}
	return tat.talep != nil
}

// Latch latches error e if e is not nil and if no prior error has been
// latched.
//
// If e is not nil and a prior error has been latched that differs from e,
// then a count of post-latch errors is incremented.
//
// Latch returns true if an error is or was previously latched.
func (tat *Tattler) Latch(e error) bool {
	if e != nil {
		return tat.fullLatch(1, e)
	}
	return false
}

// Latchf creates a new error using printf-style arguments,
// passes it to Latch(), and returns the latched error.
// The error is created with fmt.Errorf so the %w verb
// is available in the format string.
// See fmt.Errorf for details.
func (tat *Tattler) Latchf(s string, v ...interface{}) error {
	tat.fullLatch(1, fmt.Errorf(s, v...))
	return tat.Le()

}

// Le returns the latched error, nil if none.  Mnemonics for Le
// are Latched error, or the punny tat.Le() "tattle".
func (tat *Tattler) Le() error {
	if tat.talep == nil {
		return nil
	}
	return tat.talep.latched
}

// Led (mnemonics Latched, or tattled) returns true if an error has been latched.
func (tat *Tattler) Led() bool { return tat.talep != nil }

// Log logs a latched error.  Log is a no-op if the tattler is not latched,
// or if the error was previously logged.
//
// Typically the call to Logf is deferred at the beginning of a method, so
// that it will log any error that occurs during the method.  For example
// for a type Record with embedded mutex mux and tattler tat, you might
// write:
//
//	func (p *Record) Method() {
//	     p.mux.Lock()         // common. Locks record,
//	     defer p.mux.Unlock() // including tattler.
//	     defer p.tat.Log()
//
//	     <body of Method() with various Latch() cases>
//	}
//
// Each Tattler instance is only logged once.  The encapsulated error may be
// logged again if it is extracted and latched into a different Tattler
// instance.
func (tat *Tattler) Log() {
	tat.Logf("")
}

// Logf logs a latched error, using the provided arguments as a printf-style
// prefix.  Logf is a no-op if the tattler is not latched, or if it was
// previously logged.
//
// Typically the call to Logf is deferred at the beginning of a method, so that
// it will log any error that occurs during the method.  For example for a type
// Record with embedded mutex mux and tattler tat, you might write:
//
//	func (p *Record) Method() {
//	     p.mux.Lock()         // common. Locks record,
//	     defer p.mux.Unlock() // including tattler.
//	     defer p.tat.Logf("Record.Method")
//
//	     <body of Method() with various Latch() cases>
//	}
//
// In this case Logf will log the
// error prefixed with "Record.Method:"
//
// Each Tattler instance is only logged once.  The encapsulated error
// may be logged again if it is extracted and latched into another Tattler instance.
func (tat *Tattler) Logf(s string, v ...interface{}) {
	t := tat.talep
	if t != nil && !t.logged {
		prefix := fmt.Sprintf(s, v...)
		log.Printf("%s%s", prefix, tat.String())
		t.logged = true
	}
}

// Ok returns true if no error has been latched.
func (tat *Tattler) Ok() bool { return tat.talep == nil }

// String returns a string containing details of the latched error,
// including, including a limited trace back.
func (tat *Tattler) String() string {

	t := tat.talep
	if t != nil && t.latched != nil {
		sb := strings.Builder{}

		fmt.Fprintf(&sb, "%s\n", t.latched.Error())
		if len(t.frames) > 0 { // Yes, always is >0, but I'm paranoid.
			fmt.Fprintf(&sb, " Latched at:  %s:%d in %s\n",
				filepath.Base(t.frames[0].File), t.frames[0].Line, t.frames[0].Function)
			for _, frame := range t.frames[1:] {
				fmt.Fprintf(&sb, " Called From: %s:%d in %s\n",
					filepath.Base(frame.File), frame.Line, frame.Function)
			}
		}
		if t.missed != 0 {
			fmt.Fprintf(&sb, " %d post-latch errors\n", t.missed)
		}
		return sb.String()
	}
	return ""
}
