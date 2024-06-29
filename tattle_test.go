// Copyright 2024 Richard Northscope.  All rights reserved.
// Use of this source code is governed by the
// MIT license that can be found in the LICENSE file.

package tattle

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"gotest.tools/assert"
)

type tattler = Tattler

var _ = fmt.Printf

func TestBasic(t *testing.T) {
	err := testf1(t)
	assert.Error(t, err, "String1 bad value plugh")
	err = testf2(t)
	assert.Error(t, err, "String1 bad value plugh")
	assert.Assert(t, testf3(t) == 99)

}

func testf1(t *testing.T) error {
	tat := Tattler{}

	string1 := "plugh"
	int2 := 5
	assert.NilError(t, tat.Le())
	assert.Equal(t, tat.Ok(), true)
	assert.Equal(t, tat.Led(), false)
	assert.Equal(t, tat.String(), "")
	tat.Latchf("String1 bad value %s", string1)
	tat.Latchf("Int2 bad value %d", int2)
	return tat.Le()
}

func testf2(t *testing.T) error {
	tat := tattler{}

	string1 := "plugh"
	int2 := 5
	tat.Latchf("String1 bad value %s", string1)
	tat.Latchf("Int2 bad value %d", int2)
	assert.Equal(t, tat.Ok(), false)
	assert.Equal(t, tat.Led(), true)
	return tat.Le()
}

func testf3(t *testing.T) int {
	tat := tattler{}

	for i := 0; i < 100; i++ {
		tat.Latchf("%d", i)
	}
	assert.Assert(t, tat.talep != nil)
	return int(tat.talep.missed)
}

func TestFileAndLine(t *testing.T) {
	var tat1, tat2 tattler
	_, markFile, markLine, ok := runtime.Caller(0) // Mark
	assert.Assert(t, ok)                           // Mark + 1
	tat1.Latchf("Test Error")                      // Mark + 2
	tat2.Latch(fmt.Errorf("Alternate test error")) // Mark + 3

	tat1.Latchf("Post-latch error")
	assert.Assert(t, tat1.talep != nil)
	assert.Assert(t, tat1.talep.latched != nil)
	assert.Assert(t, tat1.talep.frames[0].File == markFile)

	assert.Assert(t, tat1.talep.frames[0].Line == markLine+2) // Ref. Mark + 2
	assert.Assert(t, tat2.talep.frames[0].Line == markLine+3) // Ref. Mark + 3
	assert.Assert(t, tat1.talep.missed == 1)
	tat1.Latch(tat1.Le()) // Latching self shouldn't have any effect
	assert.Assert(t, tat1.talep.missed == 1)
	//
	// Convenient place to test String
	tat1.talep.missed = 1234
	s := tat1.String()
	assert.Assert(t, strings.Contains(s, filepath.Base(markFile)))
	assert.Assert(t, strings.Contains(s, "Test Error"))
	assert.Assert(t, strings.Contains(s, strconv.FormatInt(int64(markLine+2), 10)))
	assert.Assert(t, strings.Contains(s, "1234"))
	//
	// And a convenient place to test Import
	iTat := tattler{}
	iTat.Import(&tat1)
	assert.Assert(t, iTat.String() == tat1.String())
	// Nil tat latch
	iTat = tattler{}
	assert.Assert(t, !iTat.Latch(nil))

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
	assert.Assert(t, strings.Contains(s, "Header"))
	assert.Assert(t, strings.Contains(s, "Body"))

	sw.Reset()
	tat := tattler{}
	tat.Latchf("Error xyz")
	tat.Log()
	s = sw.String()
	assert.Assert(t, strings.Contains(s, "Error xyz"))

	log.SetOutput(w)
}

func TestSetFrame(t *testing.T) {
	cf := callFrames
	SetFrames(10000)
	assert.Assert(t, callFrames == 100)
	SetFrames(cf)
}
