// Copyright 2024 Richard Northscope.  All rights reserved.
// Use of this source code is governed by the
// MIT license that can be found in the LICENSE file.package tattle
// type tattler tattle.Tattler already appears

package tattle

type Record struct {
	ID   int
	Name string
	Age  int
	tat  tattler // assuming 'type tattler tattle.Tattler'
}

// This example shows the typical method structure for
// a type that includes a tattler

func (rp *Record) Check() error {
	// Always defer a log call.
	defer rp.tat.Logf("Record id %d", rp.ID)

	// If already latched, stop processing record.
	if rp.tat.Led() {
		return rp.tat.Le()
	}

	// Body of function starts here
	if len(rp.Name) > 1000 {
		return rp.tat.Latchf("name length %d exceeds maximum %d", len(rp.Name), 1000)
	}
	if rp.Age < 0 || rp.Age > 1000 {
		return rp.tat.Latchf("age %d not within range 0-%d", rp.Age, 1000)
	}
	return rp.tat.Le() // nil

}

func Example() {

	rp := &Record{ID: 15, Age: 1001}
	//
	// The caller need not know about tattlers
	err := rp.Check()
	if err != nil {
		// error recovery
	}
}
