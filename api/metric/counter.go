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

type Float64Counter struct {
	Instrument
}

type Int64Counter struct {
	Instrument
}

type Float64CounterHandle struct {
	Handle
}

type Int64CounterHandle struct {
	Handle
}

func NewFloat64Counter(name string, mos ...Option) (c Float64Counter) {
	registerInstrument(name, CounterKind, mos, &c.Instrument)
	return
}

func NewInt64Counter(name string, mos ...Option) (c Int64Counter) {
	registerInstrument(name, CounterKind, mos, &c.Instrument)
	return
}

func (c *Float64Counter) GetHandle(ctx context.Context, labels LabelSet) (h Float64CounterHandle) {
	h.Recorder = labels.Meter().RecorderFor(ctx, labels, c.Instrument)
	return
}

func (c *Int64Counter) GetHandle(ctx context.Context, labels LabelSet) (h Int64CounterHandle) {
	h.Recorder = labels.Meter().RecorderFor(ctx, labels, c.Instrument)
	return
}

func (g *Float64Counter) Add(ctx context.Context, value float64, labels LabelSet) {
	labels.Meter().RecordSingle(ctx, labels, Measurement{
		Instrument: g.Instrument,
		Value:      value,
	})
}

func (g *Int64Counter) Add(ctx context.Context, value int64, labels LabelSet) {
	labels.Meter().RecordSingle(ctx, labels, Measurement{
		Instrument: g.Instrument,
		Value:      float64(value),
	})
}

func (g *Float64CounterHandle) Add(ctx context.Context, value float64) {
	g.Recorder.Record(ctx, value)
}

func (g *Int64CounterHandle) Add(ctx context.Context, value int64) {
	g.Recorder.Record(ctx, float64(value))
}
