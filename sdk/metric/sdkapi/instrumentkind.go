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

//go:generate stringer -type=InstrumentKind

package sdkapi

// InstrumentKind describes the kind of instrument.
type InstrumentKind int8

const (
	// HistogramInstrumentKind indicates a Histogram instrument.
	HistogramInstrumentKind InstrumentKind = iota
	// GaugeObserverInstrumentKind indicates an GaugeObserver instrument.
	GaugeObserverInstrumentKind

	// CounterInstrumentKind indicates a Counter instrument.
	CounterInstrumentKind
	// UpDownCounterInstrumentKind indicates a UpDownCounter instrument.
	UpDownCounterInstrumentKind

	// CounterObserverInstrumentKind indicates a CounterObserver instrument.
	CounterObserverInstrumentKind
	// UpDownCounterObserverInstrumentKind indicates a UpDownCounterObserver
	// instrument.
	UpDownCounterObserverInstrumentKind
)

// Synchronous returns whether this is a synchronous kind of instrument.
func (k InstrumentKind) Synchronous() bool {
	switch k {
	case CounterInstrumentKind, UpDownCounterInstrumentKind, HistogramInstrumentKind:
		return true
	}
	return false
}

// HasTemporality
func (k InstrumentKind) HasTemporality() bool {
	return k != GaugeObserverInstrumentKind
}
