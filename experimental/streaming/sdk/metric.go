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

package sdk

import (
	"context"
	"time"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/experimental/streaming/exporter"
)

type metricHandle struct {
	descriptor *metric.Descriptor
	labels     metricLabels
}

var _ metric.Handle = &metricHandle{}

type metricLabels struct {
	sdk   *sdk
	scope exporter.ScopeID
}

var _ metric.LabelSet = &metricLabels{}

func (h *metricHandle) RecordOne(ctx context.Context, value metric.MeasurementValue) {
	h.labels.sdk.exporter.Record(exporter.Event{
		Type:    exporter.SINGLE_METRIC,
		Context: ctx,
		Scope:   h.labels.scope,
		Measurement: metric.Measurement{
			Descriptor: h.descriptor,
			Value:      value,
		},
	})
}

func (m metricLabels) Meter() metric.Meter {
	return m.sdk
}

func (s *sdk) DefineLabels(ctx context.Context, labels ...core.KeyValue) metric.LabelSet {
	return metricLabels{
		sdk:   s,
		scope: s.exporter.NewScope(exporter.ScopeID{}, labels...),
	}
}

func (s *sdk) NewHandle(descriptor *metric.Descriptor, labels metric.LabelSet) metric.Handle {
	mlabels, _ := labels.(metricLabels)

	return &metricHandle{
		descriptor: descriptor,
		labels:     mlabels,
	}
}

func (s *sdk) DeleteHandle(handle metric.Handle) {
}

func (s *sdk) RecordBatch(ctx context.Context, labels metric.LabelSet, ms ...metric.Measurement) {
	eventType := exporter.BATCH_METRIC
	if len(ms) == 1 {
		eventType = exporter.SINGLE_METRIC
	}
	oms := make([]metric.Measurement, len(ms))
	mlabels, _ := labels.(metricLabels)

	copy(oms, ms)

	s.exporter.Record(exporter.Event{
		Type:         eventType,
		Context:      ctx,
		Scope:        mlabels.scope,
		Measurements: oms,
	})
}

func (s *sdk) RegisterObserver(observer metric.Observer, callback metric.ObserverCallback) {
	if s.insertNewObserver(observer, callback) {
		go s.observersRoutine()
	}
}

func (s *sdk) insertNewObserver(observer metric.Observer, callback metric.ObserverCallback) bool {
	s.observersLock.Lock()
	defer s.observersLock.Unlock()
	old := s.loadObserversMap()
	id := observer.Descriptor.ID()
	if _, ok := old[id]; ok {
		return false
	}
	observers := make(observersMap)
	for oid, data := range old {
		observers[oid] = data
	}
	observers[id] = observerData{
		observer: observer,
		callback: callback,
	}
	s.storeObserversMap(observers)
	return old == nil
}

func (s *sdk) UnregisterObserver(observer metric.Observer) {
	s.observersLock.Lock()
	defer s.observersLock.Unlock()
	old := s.loadObserversMap()
	id := observer.Descriptor.ID()
	if _, ok := old[id]; !ok {
		return
	}
	if len(old) == 1 {
		s.storeObserversMap(nil)
		return
	}
	observers := make(observersMap)
	for oid, data := range old {
		if oid != id {
			observers[oid] = data
		}
	}
	s.storeObserversMap(observers)
}

func (s *sdk) observersRoutine() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		m := s.loadObserversMap()
		if m == nil {
			return
		}
		for _, data := range m {
			ocb := s.getObservationCallback(data.observer.Descriptor)
			data.callback(s, data.observer, ocb)
		}
	}
}

func (s *sdk) getObservationCallback(descriptor *metric.Descriptor) metric.ObservationCallback {
	return func(l metric.LabelSet, v metric.MeasurementValue) {
		s.RecordBatch(context.Background(), l, metric.Measurement{
			Descriptor: descriptor,
			Value:      v,
		})
	}
}

func (s *sdk) loadObserversMap() observersMap {
	i := s.observers.Load()
	if i == nil {
		return nil
	}
	m := i.(observersMap)
	return m
}

func (s *sdk) storeObserversMap(m observersMap) {
	s.observers.Store(m)
}
