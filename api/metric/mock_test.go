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

	"go.opentelemetry.io/api/core"
)

type (
	mockHandle struct {
		descriptor *Descriptor
		labelSet   *mockLabelSet
	}

	mockLabelSet struct {
		meter  *mockMeter
		labels map[core.Key]core.Value
	}

	batch struct {
		ctx          context.Context
		labelSet     *mockLabelSet
		measurements []Measurement
	}

	observerData struct {
		observer Observer
		callback ObserverCallback
	}

	observerMap map[DescriptorID]observerData

	mockMeter struct {
		measurementBatches []batch
		observers          observerMap
	}
)

var (
	_ Handle   = &mockHandle{}
	_ LabelSet = &mockLabelSet{}
	_ Meter    = &mockMeter{}
)

func (h *mockHandle) RecordOne(ctx context.Context, value MeasurementValue) {
	h.labelSet.meter.RecordBatch(ctx, h.labelSet, Measurement{
		Descriptor: h.descriptor,
		Value:      value,
	})
}

func (s *mockLabelSet) Meter() Meter {
	return s.meter
}

func newMockMeter() *mockMeter {
	return &mockMeter{}
}

func (m *mockMeter) DefineLabels(ctx context.Context, labels ...core.KeyValue) LabelSet {
	ul := make(map[core.Key]core.Value)
	for _, kv := range labels {
		ul[kv.Key] = kv.Value
	}
	return &mockLabelSet{
		meter:  m,
		labels: ul,
	}
}

func (m *mockMeter) RecordBatch(ctx context.Context, labels LabelSet, measurements ...Measurement) {
	ourLabelSet := labels.(*mockLabelSet)
	m.measurementBatches = append(m.measurementBatches, batch{
		ctx:          ctx,
		labelSet:     ourLabelSet,
		measurements: measurements,
	})
}

func (m *mockMeter) NewHandle(erm ExplicitReportingMetric, labels LabelSet) Handle {
	descriptor := erm.Descriptor()
	ourLabels := labels.(*mockLabelSet)

	return &mockHandle{
		descriptor: descriptor,
		labelSet:   ourLabels,
	}
}

func (m *mockMeter) DeleteHandle(Handle) {
}

func (m *mockMeter) RegisterObserver(o Observer, cb ObserverCallback) {
	id := o.Descriptor().ID()
	if _, ok := m.observers[id]; ok {
		return
	}
	data := observerData{
		observer: o,
		callback: cb,
	}
	if m.observers == nil {
		m.observers = observerMap{
			id: data,
		}
	} else {
		m.observers[id] = data
	}
}

func (m *mockMeter) UnregisterObserver(o Observer) {
	delete(m.observers, o.Descriptor().ID())
}

func (m *mockMeter) PerformObservations() {
	for _, data := range m.observers {
		o := data.observer
		descriptor := o.Descriptor()
		ocb := func(l LabelSet, v MeasurementValue) {
			m.RecordBatch(context.Background(), l, Measurement{
				Descriptor: descriptor,
				Value:      v,
			})
		}
		data.callback(m, o, ocb)
	}
}
