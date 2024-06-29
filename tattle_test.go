// Copyright 2024 Richard Northscope.  All rights reserved.
// Use of this source code is governed by the
// MIT license that can be found in the LICENSE file.

package tattle

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type tattler = Tattler

var _ = fmt.Printf

func TestBasic(t *testing.T) {
	tat := Tattler{}
	string1 := "plugh"
	int2 := 5

	// Empty tattler

	if tat.Le() != nil {
		t.Errorf("Empty tattler has non-empty error")
	}
	if !tat.Ok() {
		t.Errorf("Empty tattler Ok() is false")
	}
	if tat.Led() {
		t.Errorf("Empty tattler Led() is true")
	}
	if tat.String() != "" {
		t.Errorf("Empty tattler String() is not \"\"")
	}

	// Tattler given 2 errors; first one latches.
	tat.Latchf("String1 bad value %s", string1)
	tat.Latchf("Int2 bad value %d", int2)
	if tat.Ok() {
		t.Fatalf("Tattler failed to latch on Latchf")
	}
	if !tat.Latch(nil) {
		t.Errorf("Latched tattler returns false when latching nil")
	}
	expected := "String1 bad value plugh"
	if tat.Le().Error() != expected {
		t.Errorf("Error %v not expected first error %s", tat.Le(), expected)
	}
	if tat.talep.missed != 1 {
		t.Errorf("Tattler reported %d missed errors not %d", tat.talep.missed, 1)
	}
	if !tat.Latch(tat.Le()) {
		t.Errorf("Latched tattler returns false when latching self")
	}
	if tat.talep.missed != 1 {
		t.Errorf("Latched tattler, count of missed latches changed when latching self")
	}

}

func TestFileAndLine(t *testing.T) {
	var tat1, tat2 tattler
	_, markFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Call to runtime.Caller(0) failed")
	}
	tat1.Latchf("Test Error")
	if tat1.Ok() {
		t.Fatalf("Latchf failure") // more in basic tests
	}
	tat2.Latch(fmt.Errorf("Alternate test error"))
	if tat2.Ok() {
		t.Fatalf("Latch failure")
	}
	got := tat1.talep.frames[0].File
	if got != markFile {
		t.Errorf("Trace back has '%s' expected '%s'", got, markFile)
	}

	tat1.talep.missed = 1234
	s := tat1.String()
	got = filepath.Base(s)
	if !strings.Contains(s, got) {
		t.Errorf("Log message does not contain file name '%s'", got)

	}
	if !strings.Contains(s, "1234") {
		t.Errorf("Log message does not contain repeat count 1234")
	}
	//
	// And a convenient place to test Import
	iTat := tattler{}
	iTat.Import(&tat1)
	if iTat.String() != tat1.String() {
		t.Errorf("Import generates inconsistent error '%v'/'%v'", &iTat, &tat1)
	}
	// Test reset
	tat1.Reset()
	if tat1.Led() {
		t.Errorf("tat.Reset() failed to clear tattler")
	}
}

// Test conventional deferred logf
func TestLogf(t *testing.T) {
	w := log.Writer()
	sw := new(strings.Builder)
	log.SetOutput(sw)
	(func() {
		tat := tattler{}
		defer tat.Logf("Header")
		tat.Latchf("Body")
	})()
	s := sw.String()
	if !strings.Contains(s, "Header") {
		t.Errorf("Error string '%v' lacks logf string 'Header'", s)
	}
	if !strings.Contains(s, "Body") {
		t.Errorf("Error string '%v' lacks Latchf string 'Body'", s)
	}

	sw.Reset()
	tat := tattler{}
	expected := "Error xyz"

	tat.Latchf(expected)
	tat.Log()
	s = sw.String()
	if !strings.Contains(s, expected) {
		t.Errorf("Log string '%s' lacks Latchf string '%s'", s, expected)
	}

	log.SetOutput(w)
}

func TestSetFrame(t *testing.T) {
	cf := traceFrames
	SetFrames(10000)
	if traceFrames != 100 {
		t.Errorf("SetFrames did not return expected result 100")
	}
	SetFrames(cf)
}
