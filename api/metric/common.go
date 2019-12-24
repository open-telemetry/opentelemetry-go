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

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type commonMetric struct {
	instrument InstrumentImpl
}

type commonBoundInstrument struct {
	boundInstrument BoundInstrumentImpl
}

func (m commonMetric) bind(labels LabelSet) commonBoundInstrument {
	return newCommonBoundInstrument(m.instrument.Bind(labels))
}

func (m commonMetric) float64Measurement(value float64) Measurement {
	return newMeasurement(m.instrument, core.NewFloat64Number(value))
}

func (m commonMetric) int64Measurement(value int64) Measurement {
	return newMeasurement(m.instrument, core.NewInt64Number(value))
}

func (m commonMetric) directRecord(ctx context.Context, number core.Number, labels LabelSet) {
	m.instrument.RecordOne(ctx, number, labels)
}

func (m commonMetric) Impl() InstrumentImpl {
	return m.instrument
}

func (h commonBoundInstrument) directRecord(ctx context.Context, number core.Number) {
	h.boundInstrument.RecordOne(ctx, number)
}

func (h commonBoundInstrument) Unbind() {
	h.boundInstrument.Unbind()
}

func newCommonMetric(instrument InstrumentImpl) commonMetric {
	return commonMetric{
		instrument: instrument,
	}
}

func newCommonBoundInstrument(boundInstrument BoundInstrumentImpl) commonBoundInstrument {
	return commonBoundInstrument{
		boundInstrument: boundInstrument,
	}
}

func newMeasurement(instrument InstrumentImpl, number core.Number) Measurement {
	return Measurement{
		instrument: instrument,
		number:     number,
	}
}
