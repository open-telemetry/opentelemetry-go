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
)

type Float64Measure struct {
	Instrument
}

type Int64Measure struct {
	Instrument
}

type Float64MeasureHandle struct {
	Handle
}

type Int64MeasureHandle struct {
	Handle
}

func NewFloat64Measure(name string, mos ...Option) (m Float64Measure) {
	registerInstrument(name, MeasureKind, mos, &m.Instrument)
	return
}

func NewInt64Measure(name string, mos ...Option) (m Int64Measure) {
	registerInstrument(name, MeasureKind, mos, &m.Instrument)
	return
}

func (m *Float64Measure) GetHandle(ctx context.Context, labels LabelSet) (h Float64MeasureHandle) {
	h.Recorder = labels.Meter().RecorderFor(ctx, labels, m.Instrument)
	return
}

func (m *Int64Measure) GetHandle(ctx context.Context, labels LabelSet) (h Int64MeasureHandle) {
	h.Recorder = labels.Meter().RecorderFor(ctx, labels, m.Instrument)
	return
}

func (g *Float64Measure) Record(ctx context.Context, value float64, labels LabelSet) {
	labels.Meter().RecordSingle(ctx, labels, Measurement{
		Instrument: g.Instrument,
		Value:      value,
	})
}

func (g *Int64Measure) Record(ctx context.Context, value int64, labels LabelSet) {
	labels.Meter().RecordSingle(ctx, labels, Measurement{
		Instrument: g.Instrument,
		Value:      float64(value),
	})
}

func (g *Float64MeasureHandle) Record(ctx context.Context, value float64) {
	g.Recorder.Record(ctx, value)
}

func (g *Int64MeasureHandle) Record(ctx context.Context, value int64) {
	g.Recorder.Record(ctx, float64(value))
}
