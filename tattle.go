// Copyright 2024 Richard Northscope.  All rights reserved.
// Use of this source code is governed by the
// MIT license that can be found in the LICENSE file.

/*
# Synopsis

Tattlers provide a mechanism to "latch" and log an error value associated
with some structure.  "Latching" is described further below.

One benefit of tattlers accrues from a "deferred logging" capability.

Tattlers are not complex.  There may be a good case for re-implementing
tattler-like ideas in one of your own packages. The "Implementation" section
below may be helpful.

# Introduction and Use Case

Package tattle was developed as part of a project that processes records
from many sources for storage in an immutable database. The processing
pipelines use Tattlers to store, or "latch" the first error encountered
during processing within the record that triggered the error.  If the error
cannot be resolved and cleared during processing, the record is rejected.

A Tattler does not implement the Error interface and is not intended
as a wrapper or alternative Error type.  It simply stores the first error
value.

Tattlers have logging methods, which do nothing if no error has been
latched.  The Tattler method Logf is typically used to provide context
information for the error, for example, a record ID.  This context is
usually known at the start of functions.  This affords a very helpful
capability: By deferring a call to a tattler Logf at the beginning of a
function, you automatically enable logging for any latched error that occurs
in any path through the function.

# Effective Use

Here is an effective way to use Tattlers:

 1. Include a Tattler variable, e.g., 'tat', in instances of types that have
    complex methods.
 2. Latch errors in the methods to the tat.
 3. Defer a tat Logf call in every method, so that errors get logged in
    close proximity to their occurrence.
 4. Once a tat is latched, treat the enclosing structure instance as
    tainted. If called with the tattler latched, it is fine to return the
    already-latched error.  (Return the error, not the tattler.)

A godoc example of the above follows this introduction.

# Implementation

Prior to an error, a tattler is a nil pointer to a concrete private type
known as a tale.

When a Latch method sees a non-nil error, it allocates the tale to store the
error and a few (default 3) frames from the stack.  This backtrace is only
accessible by the log functions; the error itself is unchanged. Only
the first non-nil error is latched. Subsequent different errors are counted
but their error values are not recorded by the tattler.

The log functions are low cost in the non-error case.  The printf-style
string in Logf isn't expanded unless there is an error.

The gnarly bits of code where stack trace information is gathered closely
follow the examples in the golang runtime documentation.

Tattlers have no inherent facility for concurrency.  If a structure contains
a Tattler (or an error variable) and there is a possibility of concurrent
access to the structure, the structure itself should be appropriately
protected.
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

// number of call frames logged
var traceFrames = uint32(3)

// SetFrames sets the number of trace back frames, default 3, included with
// errors reported by Log/Logf.  The requested number of frames is stopped
// to the range [0,100].
//
// SetFrames is meant to be called, if needed, during startup before
// multiple goroutines can see tattlers.  The setting applies at the process
// level.
func SetFrames(f uint32) {
	if f > 100 {
		f = 100
	}
	traceFrames = f
}

// fullLatch contains the latch logic.
// parameter b is the difference in frames
// between Latch and fulLatch.
func (tat *Tattler) fullLatch(b int, e error) bool {
	if e != nil {
		if tat.talep == nil {
			tp := new(tale)
			tat.talep = tp
			tp.latched = e

			// Capture backtrace frames.
			if traceFrames > 0 {
				pc := make([]uintptr, traceFrames)
				n := runtime.Callers(b+2, pc)
				pc = pc[:n] // truncate invalid entries

				var frame runtime.Frame
				frames := runtime.CallersFrames(pc)
				for more := true; more; {
					frame, more = frames.Next()
					tp.frames = append(tp.frames, frame)
				}
			}
		} else {
			if e != tat.talep.latched {
				tat.talep.missed++
			}
		}
	}
	return tat.talep != nil
}

// Import has semantics similar to Latch, except that its argument
// is a Tattler not an error.  Import is rarely used.
// An example use case is to taint a container
// with a tattle from an enclosed structure
// without generating an additional log message.
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
// Latch returns true if an error is, or was previously, latched.
func (tat *Tattler) Latch(e error) bool {
	if e != nil {
		return tat.fullLatch(1, e)
	}
	return tat.talep != nil
}

// Latchf creates a new error using printf-style arguments,
// passes it to Latch(), and returns the latched error.
// The error is created with fmt.Errorf so the %w verb
// is available in the format string.
// See fmt.Errorf for details.
//
// Latchf does not over-write a previously latched error.
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
// Typically the call to Log (though more often Logf below) is deferred at
// the beginning of a method, so that it will log any error that occurs
// during the method.  For example for a type Record with embedded mutex mux
// and tattler tat, you might write:
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
//	     defer p.tat.Logf("Record id %d", p.id)
//
//	     <body of Method() with various Latch() cases>
//	}
//
// Each Tattler instance is only logged once.  The encapsulated error
// may be logged again if it is extracted and latched into another Tattler instance.
func (tat *Tattler) Logf(s string, v ...interface{}) {
	t := tat.talep
	if t != nil && !t.logged {
		tat.fullLogf(s, v...)
	}
}

func (tat *Tattler) fullLogf(s string, v ...interface{}) {
	prefix := fmt.Sprintf(s, v...)
	log.Printf("%s%s", prefix, tat.String())
	tat.talep.logged = true
}

// Ok returns true if no error has been latched.
func (tat *Tattler) Ok() bool { return tat.talep == nil }

// Reset resets the tattler.
func (tat *Tattler) Reset() {
	tat.talep = nil
}

// String returns a string containing details of the latched error,
// including a limited trace back.
func (tat *Tattler) String() string {

	t := tat.talep
	if t != nil && t.latched != nil {
		sb := strings.Builder{}

		fmt.Fprintf(&sb, "%s\n", t.latched.Error())
		if len(t.frames) > 0 {
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
