package tattle

import (
	"fmt"
	"log"
	"testing"
)

type BenchRecord struct {
	ID   int
	Name string
	Age  int
	err  error   // for "hand coded" variant
	tat  tattler // assuming 'type tattler tattle.Tattler'
}

func Benchmark01CallOverhead(b *testing.B) {
	rp := &BenchRecord{}
	for n := 0; n < b.N; n++ {
		rp.Check01CallOverhead()
	}
}

// go:noinline
func (rp *BenchRecord) Check01CallOverhead() error {
	// Always defer a log call.
	//defer rp.tat.Logf("Record id %d", rp.ID)

	// If already latched, stop processing record.
	//if rp.tat.Led() {
	//	return rp.tat.Le()
	//}

	// Body of function starts here.  Here is a simple sample.
	//if len(rp.Name) > 1000 {
	//	return rp.tat.Latchf("name length %d exceeds maximum %d", len(rp.Name), 1000)
	//}
	//if rp.Age < 0 || rp.Age > 1000 {
	//	return rp.tat.Latchf("age %d not within range 0-%d", rp.Age, 1000)
	//}
	//return rp.tat.Le() // nil
	return nil

}

func Benchmark02TatErrorsOnly(b *testing.B) {
	rp := &BenchRecord{}
	for n := 0; n < b.N; n++ {
		rp.Check02TatErrorsOnly()
	}
}

// go:noinline
func (rp *BenchRecord) Check02TatErrorsOnly() error {
	// Always defer a log call.
	//defer rp.tat.Logf("Record id %d", rp.ID)

	// If already latched, stop processing record.
	//if rp.tat.Led() {
	//	return rp.tat.Le()
	//}

	// Body of function starts here.  Here is a simple sample.
	if len(rp.Name) > 1000 {
		return rp.tat.Latchf("name length %d exceeds maximum %d", len(rp.Name), 1000)
	}
	if rp.Age < 0 || rp.Age > 1000 {
		return rp.tat.Latchf("age %d not within range 0-%d", rp.Age, 1000)
	}
	return rp.tat.Le() // nil
}
func Benchmark03AddTaintCheck(b *testing.B) {
	rp := &BenchRecord{}
	for n := 0; n < b.N; n++ {
		rp.Check03AddTaintCheck()
	}
}

// go:noinline
func (rp *BenchRecord) Check03AddTaintCheck() error {
	// Always defer a log call.
	//defer rp.tat.Logf("Record id %d", rp.ID)

	// If already latched, stop processing record.
	if rp.tat.Led() {
		return rp.tat.Le()
	}

	// Body of function starts here.  Here is a simple sample.
	if len(rp.Name) > 1000 {
		return rp.tat.Latchf("name length %d exceeds maximum %d", len(rp.Name), 1000)
	}
	if rp.Age < 0 || rp.Age > 1000 {
		return rp.tat.Latchf("age %d not within range 0-%d", rp.Age, 1000)
	}
	return rp.tat.Le() // nil
}
func Benchmark04StandardTemplate(b *testing.B) {
	rp := &BenchRecord{}
	for n := 0; n < b.N; n++ {
		rp.Check04StandardTemplate()
	}
}

// go:noinline
func (rp *BenchRecord) Check04StandardTemplate() error {
	// Always defer a log call.
	defer rp.tat.Logf("Record id %d", rp.ID)

	// If already latched, stop processing record.
	if rp.tat.Led() {
		return rp.tat.Le()
	}

	// Body of function starts here.  Here is a simple sample.
	if len(rp.Name) > 1000 {
		return rp.tat.Latchf("name length %d exceeds maximum %d", len(rp.Name), 1000)
	}
	if rp.Age < 0 || rp.Age > 1000 {
		return rp.tat.Latchf("age %d not within range 0-%d", rp.Age, 1000)
	}
	return rp.tat.Le() // nil
}
func Benchmark05TemplateLogVsLogf(b *testing.B) {
	rp := &BenchRecord{}
	for n := 0; n < b.N; n++ {
		rp.Check05TemplateLogVsLogf()
	}
}

// go:noinline
func (rp *BenchRecord) Check05TemplateLogVsLogf() error {
	// Always defer a log call.
	defer rp.tat.Log()

	// If already latched, stop processing record.
	if rp.tat.Led() {
		return rp.tat.Le()
	}

	// Body of function starts here.  Here is a simple sample.
	if len(rp.Name) > 1000 {
		return rp.tat.Latchf("name length %d exceeds maximum %d", len(rp.Name), 1000)
	}
	if rp.Age < 0 || rp.Age > 1000 {
		return rp.tat.Latchf("age %d not within range 0-%d", rp.Age, 1000)
	}
	return rp.tat.Le() // nil
}

func Benchmark10HandCoded(b *testing.B) {
	rp := &BenchRecord{}
	for n := 0; n < b.N; n++ {
		rp.Check10HandCoded()
	}
}

// go:noinline
func (rp *BenchRecord) Check10HandCoded() error {
	// Always defer a log call.
	//defer rp.tat.Logf("Record id %d", rp.ID)
	defer func() {
		if rp.err != nil {
			log.Printf("%s", rp.err.Error())
		}

	}()
	//defer rp.tat.Logf("Record id %d", rp.ID)

	if rp.err != nil {
		return rp.err
	}

	// Body of function starts here.  Here is a simple sample.
	if len(rp.Name) > 1000 {
		// return rp.tat.Latchf("name length %d exceeds maximum %d", len(rp.Name), 1000)
		rp.err = fmt.Errorf("Name length %d exceeds maximum %d", len(rp.Name), 1000)
		return rp.err

	}
	if rp.Age < 0 || rp.Age > 1000 {
		// return rp.tat.Latchf("age %d not within range 0-%d", rp.Age, 1000)
		rp.err = fmt.Errorf("age %d not within range 0-%d", rp.Age, 1000)
		return rp.err

	}
	return rp.err

}
