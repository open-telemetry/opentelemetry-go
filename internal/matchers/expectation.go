// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package matchers

import (
	"testing"
)

type Expectation struct {
	t       *testing.T
	negated bool
	actual  interface{}
}

func (e *Expectation) ToEqual(expected interface{}) {
	e.verifyExpectedNotNil(expected)

	if e.actual != expected {
		e.t.Fatalf("Expected\n\t%v\nto equal\n\t%v\n", e.actual, expected)
	}
}

func (e *Expectation) NotToEqual(expected interface{}) {
	e.verifyExpectedNotNil(expected)

	if e.actual == expected {
		e.t.Fatalf("Expected\n\t%v\nnot to equal\n\t%v\n", e.actual, expected)
	}
}

func (e *Expectation) ToBeNil() {
	if e.actual != nil {
		e.t.Fatalf("Expected\n\t%v\nto be nil\n", e.actual)
	}
}

func (e *Expectation) NotToBeNil() {
	if e.actual == nil {
		e.t.Fatalf("Expected\n\t%v\nnot to be nil\n", e.actual)
	}
}

func (e *Expectation) ToBeTrue() {
	switch a := e.actual.(type) {
	case bool:
		if e.actual == false {
			e.t.Fatalf("Expected\n\t%v\nto be true\n", e.actual)
		}
	default:
		e.t.Fatalf("Cannot check if non-bool value\n\t%v\nis truthy\n", a)
	}
}

func (e *Expectation) ToBeFalse() {
	switch a := e.actual.(type) {
	case bool:
		if e.actual == true {
			e.t.Fatalf("Expected\n\t%v\nto be false\n", e.actual)
		}
	default:
		e.t.Fatalf("Cannot check if non-bool value\n\t%v\nis truthy\n", a)
	}
}

func (e *Expectation) verifyExpectedNotNil(expected interface{}) {
	if expected == nil {
		e.t.Fatal("Refusing to compare with <nil>. Use `ToBeNil` or `NotToBeNil` instead.")
	}
}
