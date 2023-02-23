// Copyright The OpenTelemetry Authors
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

package unit // import "go.opentelemetry.io/otel/metric/unit"

// Dimensionless is the UCUM dimensionless unit 1.
var Dimensionless = Unit{code: "1"}

// prefix is a numerical modifier of a unit quantity. Its value is the
// case-sensitive UCUM prefix code.
type prefix string

// Unit is a determinate standard quantity of measurement.
type Unit struct {
	prefix prefix

	code string
}

// New returns a new Unit representing the unit code. The code should be the
// UCUM (https://ucum.org/ucum) case-sensitive code as this is what
// OpenTelemetry uses.
func New(code string) Unit { return Unit{code: code} }

// withPrefix returns a copy of u with prefix p.
func (u Unit) withPrefix(p prefix) Unit {
	u.prefix = p
	return u
}

// String returns the string encoding of the UCUM unit defining u.
func (u Unit) String() string {
	return string(u.prefix) + u.code
}

// MarshalText encodes the code of u into a textual form.
func (u Unit) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText decodes the text into a Unit.
func (u *Unit) UnmarshalText(text []byte) error {
	u.code = string(text)
	return nil
}
