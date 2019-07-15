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

package stats

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/registry"
)

type MeasureHandle struct {
	Variable registry.Variable
}

type Measure interface {
	V() registry.Variable
	M(value float64) Measurement
}

type Measurement struct {
	Measure Measure
	Value   float64
}

type Recorder interface {
	// TODO: Note as in rfc 0001, allow raw Measures to have pre-defined labels:
	GetMeasure(ctx context.Context, measure *MeasureHandle, labels ...core.KeyValue) Measure

	Record(ctx context.Context, m ...Measurement)
	RecordSingle(ctx context.Context, m Measurement)
}

type noopRecorder struct{}
type noopMeasure struct{}

var global atomic.Value

// GlobalRecorder return meter registered with global registry.
// If no meter is registered then an instance of noop Recorder is returned.
func GlobalRecorder() Recorder {
	if t := global.Load(); t != nil {
		return t.(Recorder)
	}
	return noopRecorder{}
}

// SetGlobalRecorder sets provided meter as a global meter.
func SetGlobalRecorder(t Recorder) {
	global.Store(t)
}

func Record(ctx context.Context, m ...Measurement) {
	GlobalRecorder().Record(ctx, m...)
}

func RecordSingle(ctx context.Context, m Measurement) {
	GlobalRecorder().RecordSingle(ctx, m)
}

type AnyStatistic struct{}

func (AnyStatistic) String() string {
	return "AnyStatistic"
}

var (
	WithDescription = registry.WithDescription
	WithUnit        = registry.WithUnit
)

func NewMeasure(name string, opts ...registry.Option) *MeasureHandle {
	return &MeasureHandle{
		Variable: registry.Register(name, AnyStatistic{}, opts...),
	}
}

func (m *MeasureHandle) M(value float64) Measurement {
	return Measurement{
		Measure: m,
		Value:   value,
	}
}

func (m *MeasureHandle) V() registry.Variable {
	return m.Variable
}

func (noopRecorder) Record(ctx context.Context, m ...Measurement) {
}

func (noopRecorder) RecordSingle(ctx context.Context, m Measurement) {
}

func (noopRecorder) GetMeasure(ctx context.Context, handle *MeasureHandle, labels ...core.KeyValue) Measure {
	return noopMeasure{}
}

func (noopMeasure) M(float64) Measurement {
	return Measurement{}
}

func (noopMeasure) V() registry.Variable {
	return registry.Variable{}
}
