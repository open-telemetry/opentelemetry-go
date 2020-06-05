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

package metric

// Kind describes the kind of instrument.
type Kind int8

const (
	// CounterKind indicates a Counter instrument.
	CounterKind Kind = iota

	// UpDownCounterKind indicates a UpDownCounter instrument.
	UpDownCounterKind

	// ValueRecorderKind indicates a ValueRecorder instrument.
	ValueRecorderKind

	// SumObserverKind indicates a SumObserver instrument.
	SumObserverKind

	// UpDownSumObserverKind indicates a UpDownSumObserver instrument.
	UpDownSumObserverKind

	// ValueObserverKind indicates an ValueObserver instrument.
	ValueObserverKind

	// NumKinds is the number of metric kinds.
	NumKinds = 6
)

// Synchronous returns whether this is a synchronous kind of instrument.
func (k Kind) Synchronous() bool {
	switch k {
	case CounterKind, UpDownCounterKind, ValueRecorderKind:
		return true
	}
	return false
}

// Asynchronous returns whether this is an asynchronous kind of instrument.
func (k Kind) Asynchronous() bool {
	return !k.Synchronous()
}

// Cumulative returns whether this instruments inputs are measured
// from the start of the process.
func (k Kind) Cumulative() bool {
	switch k {
	case SumObserverKind, UpDownSumObserverKind:
		return true
	}
	return false
}
