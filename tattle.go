// Copyright 2019-2024 Richard Northscope.  All rights reserved.

/*
Package tattle provides a simple facility to attach an error to the
structure that caused it.

# High level summary

What tattle does could  be done using error values:

	type Record struct {
	    <record contents>
	    err error
	}

Tattlers add the source code location of the error, and a logging
capability.

	type tattler tattle.Tattler // if you like
	...
	type Record struct {
	   <record contents>
	   tat tattler
	}

Prior to an error, a tattler is a nil pointer to a concrete structure - so
it is smaller than even a nil error, out of sensitivity to applications that
have extraordinary object counts.

The tattler's Latch method, tat.Latch(<error>), allocates a bit of memory to
store the error and its source code latch location.  Only the first non-nil
error is latched; subsequent non-nil errors are counted but discarded.

A tattler doesn't implement the error interface. Tattlers are a mousetrap,
not a mouse.  The latched error is available through method Le(). So,
tat.Le().  Yeah, it's a coding pun.

# Experience

Tattlers are in use in a developing database / http server project.
Currently at 14K lines of go code, the project defines 59
tattlers, with 183 Latch cases. The vast majority of Latch cases detect
programmatic, not end-user, errors. The tattlers are helpful.

# Concurrency

Concurrency for tattlers present the same challenge as
concurrency for error variables.  In the above-mentioned project, the
concurrency mechanisms for the structures also protect the tattlers.

Lacking a use case,
I haven't defined a "SerialTattler" type.

# Advice

  - Return errors, not tattlers.
  - The most-used methods are Latch(), Latchf(), Le(), and Logf().
  - A deferred Logf() is your friend.
*/
package tattle

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

// A Tattler is used to record an error value within a structure.
type Tattler struct {
	talep *tale
}

// A tale exists only for Tattlers that have latched an error
type tale struct {
	latched error
	trace   string // backtrace string
	file    string // used in test
	line    int    // used in test
	logged  bool
	missed  int
}

func (tat *Tattler) fullLatch(b int, e error) bool {
	if e != nil {
		if tat.talep == nil {
			tp := new(tale)
			tat.talep = tp
			tp.latched = e
			pc := make([]uintptr, 3)
			n := runtime.Callers(b+2, pc)
			pc = pc[:n] // truncate invalid entries
			sb := strings.Builder{}
			if n <= 0 {
				fmt.Fprintf(&sb, "\n\tNo back trace available\n")
			} else {
				frames := runtime.CallersFrames(pc)
				frame, more := frames.Next()
				tp.file = filepath.Base(frame.File)
				tp.line = frame.Line
				fmt.Fprintf(&sb, "\tLatched at:  %s:%d in %s\n",
					tp.file, tp.line, frame.Function)
				for more {
					frame, more = frames.Next()
					fmt.Fprintf(&sb, "\tCalled From: %s:%d in %s\n",
						filepath.Base(frame.File), frame.Line, frame.Function)
				}
			}
			tp.trace = sb.String()
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
// passes it to Latch(), and returns the latched error.  The
// returned error may have been latched prior to Latchf().
func (tat *Tattler) Latchf(s string, v ...interface{}) error {
	tat.fullLatch(1, fmt.Errorf(s, v...))
	return tat.Le()

}

// Le (mnemonics Latched Error, or tattle) returns the latched error, nil if none
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
// Typically the call to Logf is deferred at the beginning of a method, so that
// it will log any error that occurs during the method.  For example for a type
// Record with embedded mutex mux and tattler tat, you might write:
//
//	func (p *Record) Method() {
//	     p.mux.Lock()         // common. Locks record,
//	     defer p.mux.Unlock() // including tattler.
//	     defer p.tat.Log()
//
//	     <body of Method() with various Latch() cases>
//	}
//
// Each Tattler instance is only logged once.  The encapsulated error
// may be logged again if it is extracted and latched into another Tattler instance.
func (tat *Tattler) Log() {
	tat.Logf("Latched")
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
		log.Printf("%s: %s", prefix, tat.String())
		t.logged = true
	}
}

// Ok returns true if no error has been latched.
func (tat *Tattler) Ok() bool { return tat.talep == nil }

// String returns a string containing details of the latched error,
// including the source code file and line number of the latch.
func (tat *Tattler) String() string {

	t := tat.talep
	if t != nil && t.latched != nil {
		sb := strings.Builder{}

		fmt.Fprintf(&sb, "%s\n%s", t.latched.Error(), t.trace)
		if t.missed != 0 {
			fmt.Fprintf(&sb, "%d post-latch errors\n", t.missed)
		}
		return sb.String()
	}
	return ""
}
