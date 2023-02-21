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
var Dimensionless = Unit{code: "1", symbol: "1"}

// prefix is a numerical modifier of a unit quantity.
type prefix struct {
	// code is the case-sensitive code.
	code string
	// symbol is the print symbol of the prefix.
	symbol string
}

// Code returns the UCUM code defining p.
func (p prefix) Code() string {
	return p.code
}

// String returns the print symbol of p.
func (p prefix) String() string {
	return p.symbol
}

// Unit is a determinate standard quantity of measurement.
type Unit struct {
	prefix prefix

	code   string
	symbol string
}

// Option sets configuration for a Unit.
type Option interface {
	apply(u Unit) Unit
}

type fnOpt func(Unit) Unit

func (f fnOpt) apply(c Unit) Unit { return f(c) }

// WithPrintSymbol sets a Unit's print symbol. If this is not provided, the
// code will be used as a default.
func WithPrintSymbol(s string) Option {
	return fnOpt(func(c Unit) Unit {
		c.symbol = s
		return c
	})
}

// New returns a new Unit representing the unit code. The code should be the
// UCUM (https://ucum.org/ucum) case-sensitive code as this is what
// OpenTelemetry uses.
func New(code string, opts ...Option) Unit {
	u := Unit{code: code}
	for _, o := range opts {
		u = o.apply(u)
	}
	return u
}

// withPrefix returns a copy of u with prefix p.
func (u Unit) withPrefix(p prefix) Unit {
	u.prefix = p
	return u
}

// Code returns the UCUM code defining u.
func (u Unit) Code() string {
	return u.prefix.Code() + u.code
}

// String returns the print symbol of u if defined, otherwise it will return
// the UCUM code.
func (u Unit) String() string {
	if u.symbol != "" {
		return u.prefix.String() + u.symbol
	}
	return u.Code()
}

// MarshalText encodes the code of u into a textual form.
func (u Unit) MarshalText() ([]byte, error) {
	return []byte(u.code), nil
}

// UnmarshalText decodes the text into a Unit.
func (u *Unit) UnmarshalText(text []byte) error {
	u.code = string(text)
	return nil
}
