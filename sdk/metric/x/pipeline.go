// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/x"

import (
	"slices"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/x/internal/boundsum"
	"go.opentelemetry.io/otel/sdk/resource"
)

type pipeline struct {
	resource       *resource.Resource
	reader         Reader
	views          []View
	exemplarFilter exemplar.Filter
	limit          int

	mu      sync.Mutex
	streams []*counterStream
}

func newPipeline(
	res *resource.Resource,
	reader Reader,
	views []View,
	filter exemplar.Filter,
	limit int,
) *pipeline {
	return &pipeline{
		resource:       res,
		reader:         reader,
		views:          slices.Clone(views),
		exemplarFilter: filter,
		limit:          limit,
	}
}

type counterStream struct {
	scope       instrumentation.Scope
	name        string
	description string
	unit        string
	filter      attribute.Filter
	store       *boundsum.Store
}

func (p *pipeline) counterStreams(inst Instrument, allowed []attribute.Key) []*counterStream {
	p.mu.Lock()
	defer p.mu.Unlock()

	var streams []*counterStream
	matched := false
	for _, view := range p.views {
		stream, ok := view(inst)
		if !ok {
			continue
		}
		matched = true
		if !sumAggregation(stream.Aggregation, p.reader.aggregation(inst.Kind)) {
			continue
		}
		streams = append(streams, p.newCounterStream(inst, stream))
	}
	if !matched {
		stream := Stream{Name: inst.Name, Description: inst.Description, Unit: inst.Unit}
		if allowed != nil {
			stream.AttributeFilter = attribute.NewAllowKeysFilter(allowed...)
		}
		if sumAggregation(nil, p.reader.aggregation(inst.Kind)) {
			streams = append(streams, p.newCounterStream(inst, stream))
		}
	}
	p.streams = append(p.streams, streams...)
	return streams
}

func sumAggregation(configured, readerDefault Aggregation) bool {
	if configured == nil {
		configured = readerDefault
	}
	switch configured.(type) {
	case nil, AggregationDefault, AggregationSum:
		return true
	default:
		return false
	}
}

func (p *pipeline) newCounterStream(inst Instrument, stream Stream) *counterStream {
	if stream.Name == "" {
		stream.Name = inst.Name
	}
	if stream.Description == "" {
		stream.Description = inst.Description
	}
	if stream.Unit == "" {
		stream.Unit = inst.Unit
	}
	selector := stream.ExemplarReservoirProviderSelector
	if selector == nil {
		selector = DefaultExemplarReservoirProviderSelector
	}
	return &counterStream{
		scope:       inst.Scope,
		name:        stream.Name,
		description: stream.Description,
		unit:        stream.Unit,
		filter:      stream.AttributeFilter,
		store: boundsum.New(boundsum.Config{
			Temporality: p.reader.temporality(inst.Kind),
			Limit:       p.limit,
			Filter:      p.exemplarFilter,
			Reservoir:   selector(AggregationSum{}),
		}),
	}
}

func (s *counterStream) attributes(attrs attribute.Set) (attribute.Set, []attribute.KeyValue) {
	if s.filter == nil {
		return attrs, nil
	}
	return attrs.Filter(s.filter)
}

func (p *pipeline) collect(destination *metricdata.ResourceMetrics) {
	p.mu.Lock()
	defer p.mu.Unlock()
	destination.Resource = p.resource
	destination.ScopeMetrics = destination.ScopeMetrics[:0]
	now := time.Now()
	for _, stream := range p.streams {
		points := stream.store.Collect(now)
		if len(points) == 0 {
			continue
		}
		idx := -1
		for i := range destination.ScopeMetrics {
			if destination.ScopeMetrics[i].Scope == stream.scope {
				idx = i
				break
			}
		}
		if idx < 0 {
			destination.ScopeMetrics = append(destination.ScopeMetrics, metricdata.ScopeMetrics{Scope: stream.scope})
			idx = len(destination.ScopeMetrics) - 1
		}
		destination.ScopeMetrics[idx].Metrics = append(destination.ScopeMetrics[idx].Metrics, metricdata.Metrics{
			Name:        stream.name,
			Description: stream.description,
			Unit:        stream.unit,
			Data: metricdata.Sum[int64]{
				DataPoints:  points,
				Temporality: stream.store.Temporality(),
				IsMonotonic: true,
			},
		})
	}
}
