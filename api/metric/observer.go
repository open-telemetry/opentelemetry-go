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

package metric

// Observer is a base of typed-observers. Shouldn't be used directly.
type Observer struct {
	d *Descriptor
}

// Float64Observer is an observer that reports float64 values.
type Float64Observer struct {
	Observer
}

// Int64Observer is an observer that reports int64 values.
type Int64Observer struct {
	Observer
}

func newObserver(name string, valueKind ValueKind, mos ...GaugeOptionApplier) (o Observer) {
	o.d = registerDescriptor(name, ObserverKind, valueKind)
	for _, opt := range mos {
		opt.ApplyGaugeOption(o.d)
	}
	return
}

// NewFloat64Observer creates a new observer for float64.
func NewFloat64Observer(name string, mos ...GaugeOptionApplier) (o Float64Observer) {
	o.Observer = newObserver(name, Float64ValueKind, mos...)
	return
}

// NewInt64Observer creates a new observer for int64.
func NewInt64Observer(name string, mos ...GaugeOptionApplier) (o Int64Observer) {
	o.Observer = newObserver(name, Int64ValueKind, mos...)
	return
}

func (o Observer) Descriptor() *Descriptor {
	return o.d
}
