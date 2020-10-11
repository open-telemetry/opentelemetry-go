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

package metric

// InstrumentKind describes the kind of instrument.
type InstrumentKind int8

const (
	// ValueRecorderInstrumentKind indicates a ValueRecorder instrument.
	ValueRecorderInstrumentKind InstrumentKind = iota
	// ValueObserverInstrumentKind indicates an ValueObserver instrument.
	ValueObserverInstrumentKind

	// CounterInstrumentKind indicates a Counter instrument.
	CounterInstrumentKind
	// UpDownCounterInstrumentKind indicates a UpDownCounter instrument.
	UpDownCounterInstrumentKind

	// SumObserverInstrumentKind indicates a SumObserver instrument.
	SumObserverInstrumentKind
	// UpDownSumObserverInstrumentKind indicates a UpDownSumObserver
	// instrument.
	UpDownSumObserverInstrumentKind
)

// Synchronous returns whether this is a synchronous kind of instrument.
func (k InstrumentKind) Synchronous() bool {
	switch k {
	case CounterInstrumentKind, UpDownCounterInstrumentKind, ValueRecorderInstrumentKind:
		return true
	}
	return false
}

// Asynchronous returns whether this is an asynchronous kind of instrument.
func (k InstrumentKind) Asynchronous() bool {
	return !k.Synchronous()
}

// Adding returns whether this kind of instrument adds its inputs (as opposed to Grouping).
func (k InstrumentKind) Adding() bool {
	switch k {
	case CounterInstrumentKind, UpDownCounterInstrumentKind, SumObserverInstrumentKind, UpDownSumObserverInstrumentKind:
		return true
	}
	return false
}

// Grouping returns whether this kind of instrument groups its inputs (as opposed to Adding).
func (k InstrumentKind) Grouping() bool {
	return !k.Adding()
}

// Monotonic returns whether this kind of instrument exposes a non-decreasing sum.
func (k InstrumentKind) Monotonic() bool {
	switch k {
	case CounterInstrumentKind, SumObserverInstrumentKind:
		return true
	}
	return false
}

// PrecomputedSum returns whether this kind of instrument receives precomputed sums.
func (k InstrumentKind) PrecomputedSum() bool {
	return k.Adding() && k.Asynchronous()
}
