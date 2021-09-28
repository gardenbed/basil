package factory_test

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gardenbed/basil/factory"
)

func ExampleString() {
	s := factory.String()
	fmt.Printf("%s\n", s)
}

func ExampleStringPtr() {
	s := factory.StringPtr()
	fmt.Printf("%s\n", *s)
}

func ExampleStringSlice() {
	s := factory.StringSlice()
	fmt.Printf("%s\n", s)
}

func ExampleName() {
	name := factory.Name()
	fmt.Printf("%s\n", name)
}

func ExampleEmail() {
	email := factory.Email()
	fmt.Printf("%s\n", email)
}

func ExampleTime() {
	t := factory.Time()
	fmt.Printf("%s\n", t)
}

func ExampleTimeBefore() {
	t := factory.TimeBefore(time.Now())
	fmt.Printf("%s\n", t)
}

func ExampleTimeAfter() {
	t := factory.TimeAfter(time.Now())
	fmt.Printf("%s\n", t)
}

func ExampleURL() {
	u := factory.URL()
	fmt.Printf("%s\n", &u)
}

func ExamplePopulate() {
	object := struct {
		String     string
		Bool       bool
		Int        int
		Uint       uint
		Float64    float64
		Complex128 complex128
		Nested     struct {
			Duration time.Duration
			Time     *time.Time
			URL      *url.URL
		}
	}{}

	if err := factory.Populate(&object, false); err != nil {
		log.Fatalf("populate error: %s", err)
	}

	fmt.Printf("%+v\n", object)
}
