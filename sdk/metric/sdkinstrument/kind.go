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

//go:generate stringer -type=Kind

package sdkinstrument

// Kind describes the kind of instrument.
type Kind int8

const (
	// CounterKind indicates a Counter instrument.
	CounterKind Kind = iota
	// UpDownCounterKind indicates a UpDownCounter instrument.
	UpDownCounterKind
	// HistogramKind indicates a Histogram instrument.
	HistogramKind

	// CounterObserverKind indicates a CounterObserver instrument.
	CounterObserverKind
	// UpDownCounterObserverKind indicates a UpDownCounterObserver
	// instrument.
	UpDownCounterObserverKind
	// GaugeObserverKind indicates an GaugeObserver instrument.
	GaugeObserverKind

	// NumKinds is the size of an array, useful for indexing by instrument kind.
	NumKinds
)

// Synchronous returns whether this is a synchronous kind of instrument.
func (k Kind) Synchronous() bool {
	switch k {
	case CounterKind, UpDownCounterKind, HistogramKind:
		return true
	}
	return false
}

// HasTemporality
func (k Kind) HasTemporality() bool {
	return k != GaugeObserverKind
}
