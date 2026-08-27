// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/x"

import (
	"context"
	"slices"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	metricx "go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/x/internal/boundsum"
)

type counterID struct {
	name, description, unit string
}

type meter struct {
	noop.Meter
	scope instrumentation.Scope
	pipes []*pipeline

	mu       sync.Mutex
	counters map[counterID]*int64Counter
}

func (m *meter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	cfg := metric.NewInt64CounterConfig(options...)
	id := counterID{strings.ToLower(name), cfg.Description(), cfg.Unit()}
	m.mu.Lock()
	defer m.mu.Unlock()
	if counter, ok := m.counters[id]; ok {
		return counter, nil
	}
	inst := Instrument{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindCounter,
		Scope:       m.scope,
	}
	allowed := defaultAttributes(options)
	counter := &int64Counter{}
	for _, pipe := range m.pipes {
		counter.streams = append(counter.streams, pipe.counterStreams(inst, allowed)...)
	}
	m.counters[id] = counter
	return counter, nil
}

func defaultAttributes[T any](options []T) []attribute.Key {
	var keys []attribute.Key
	found := false
	for _, option := range options {
		if experimental, ok := any(option).(interface{ AllowedKeys() []attribute.Key }); ok {
			found = true
			keys = append(keys, experimental.AllowedKeys()...)
		}
	}
	if found && keys == nil {
		return []attribute.Key{}
	}
	return keys
}

type int64Counter struct {
	noop.Int64Counter
	streams []*counterStream
}

var (
	_ metric.Int64Counter        = (*int64Counter)(nil)
	_ metricx.Int64CounterBinder = (*int64Counter)(nil)
	_ metricx.Finisher           = (*int64Counter)(nil)
)

func (c *int64Counter) Enabled(context.Context) bool { return len(c.streams) > 0 }

func (c *int64Counter) Add(ctx context.Context, value int64, options ...metric.AddOption) {
	attrs := resolveAddAttributes(options)
	for _, stream := range c.streams {
		filtered, dropped := stream.attributes(attrs)
		stream.store.Measure(ctx, value, filtered, dropped)
	}
}

func (c *int64Counter) Bind(attrs ...attribute.KeyValue) metric.Int64Counter {
	set := attribute.NewSet(slices.Clone(attrs)...)
	bound := &boundInt64Counter{instrument: c, attrs: set}
	for _, stream := range c.streams {
		filtered, dropped := stream.attributes(set)
		bound.handles = append(bound.handles, stream.store.Bind(filtered, dropped))
	}
	return bound
}

func (c *int64Counter) Finish(_ context.Context, attrs ...attribute.KeyValue) {
	set := attribute.NewSet(slices.Clone(attrs)...)
	for _, stream := range c.streams {
		filtered, _ := stream.attributes(set)
		stream.store.Finish(filtered)
	}
}

type boundInt64Counter struct {
	noop.Int64Counter
	instrument *int64Counter
	attrs      attribute.Set
	handles    []*boundsum.Handle
}

var _ metric.Int64Counter = (*boundInt64Counter)(nil)

func (c *boundInt64Counter) Enabled(ctx context.Context) bool {
	return c.instrument.Enabled(ctx)
}

func (c *boundInt64Counter) Add(ctx context.Context, value int64, options ...metric.AddOption) {
	if len(options) == 0 {
		for _, handle := range c.handles {
			handle.Add(ctx, value)
		}
		return
	}
	extra := resolveAddAttributes(options)
	if extra.Len() == 0 {
		for _, handle := range c.handles {
			handle.Add(ctx, value)
		}
		return
	}
	merged := append(c.attrs.ToSlice(), extra.ToSlice()...)
	set := attribute.NewSet(merged...)
	for _, stream := range c.instrument.streams {
		filtered, dropped := stream.attributes(set)
		stream.store.Measure(ctx, value, filtered, dropped)
	}
}

func resolveAddAttributes(options []metric.AddOption) attribute.Set {
	cfg := metric.NewAddConfig(options)
	set := cfg.Attributes()
	attrs := set.ToSlice()
	for _, option := range options {
		if raw, ok := any(option).(interface{ RawAttributes() []attribute.KeyValue }); ok {
			attrs = append(attrs, raw.RawAttributes()...)
		}
	}
	return attribute.NewSet(attrs...)
}
