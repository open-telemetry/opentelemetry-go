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

//go:build go1.17
// +build go1.17

package metric // import "go.opentelemetry.io/otel/sdk/metric"

// Temporality defines the window that an aggregation was calculated over.
type Temporality uint8

const (
	// undefinedTemporality represents an unset Temporality.
	//nolint:deadcode,unused,varcheck
	undefinedTemporality Temporality = iota

	// CumulativeTemporality defines a measurement interval that continues to
	// expand forward in time from a starting point. New measurements are
	// added to all previous measurements since a start time.
	CumulativeTemporality

	// DeltaTemporality defines a measurement interval that resets each cycle.
	// Measurements from one cycle are recorded independently, measurements
	// from other cycles do not affect them.
	DeltaTemporality
)

// WithTemporality uses the selector to determine the Temporality measurements
// from instrument should be recorded with.
func WithTemporality(selector func(instrument InstrumentKind) Temporality) ReaderOption {
	return temporalitySelectorOption{selector: selector}
}

type temporalitySelectorOption struct {
	selector func(instrument InstrumentKind) Temporality
}

// applyManual returns a manualReaderConfig with option applied.
func (t temporalitySelectorOption) applyManual(mrc manualReaderConfig) manualReaderConfig {
	mrc.temporalitySelector = t.selector
	return mrc
}

// applyPeriodic returns a periodicReaderConfig with option applied.
func (t temporalitySelectorOption) applyPeriodic(prc periodicReaderConfig) periodicReaderConfig {
	prc.temporalitySelector = t.selector
	return prc
}

func defaultTemporalitySelector(_ InstrumentKind) Temporality {
	return CumulativeTemporality
}
