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
	// ValueRecorderKind indicates a ValueRecorder instrument.
	ValueRecorderKind Kind = iota
	// ValueObserverKind indicates an ValueObserver instrument.
	ValueObserverKind

	// CounterKind indicates a Counter instrument.
	CounterKind
	// UpDownCounterKind indicates a UpDownCounter instrument.
	UpDownCounterKind

	// SumObserverKind indicates a SumObserver instrument.
	SumObserverKind
	// UpDownSumObserverKind indicates a UpDownSumObserver instrument.
	UpDownSumObserverKind
)
