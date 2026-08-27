// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/x"

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var (
	// ErrReaderNotRegistered is returned when Collect is called before a
	// Reader is registered with a MeterProvider.
	ErrReaderNotRegistered = errors.New("reader is not registered")
	// ErrReaderShutdown is returned after a Reader has been shut down.
	ErrReaderShutdown = errors.New("reader is shutdown")
)

// Reader is the private-boundary reader contract used by this experimental
// module. Stable SDK Readers cannot be mixed into this provider.
type Reader interface {
	register(*pipeline)
	temporality(InstrumentKind) metricdata.Temporality
	aggregation(InstrumentKind) Aggregation
	cardinalityLimit(InstrumentKind) (int, bool)
	Collect(context.Context, *metricdata.ResourceMetrics) error
	Shutdown(context.Context) error
}

// TemporalitySelector selects temporality for an instrument kind.
type TemporalitySelector func(InstrumentKind) metricdata.Temporality

// AggregationSelector selects aggregation for an instrument kind.
type AggregationSelector func(InstrumentKind) Aggregation

// CardinalityLimitSelector selects a limit and whether the provider limit is
// used as a fallback.
type CardinalityLimitSelector func(InstrumentKind) (limit int, fallback bool)

// CumulativeTemporalitySelector selects cumulative temporality.
func CumulativeTemporalitySelector(InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

// DeltaTemporalitySelector selects delta temporality for counters.
func DeltaTemporalitySelector(kind InstrumentKind) metricdata.Temporality {
	if kind == InstrumentKindCounter {
		return metricdata.DeltaTemporality
	}
	return metricdata.CumulativeTemporality
}

// DefaultTemporalitySelector selects cumulative temporality.
func DefaultTemporalitySelector(kind InstrumentKind) metricdata.Temporality {
	return CumulativeTemporalitySelector(kind)
}

// ManualReader collects metric data on demand.
type ManualReader struct {
	mu       sync.Mutex
	pipeline *pipeline
	shutdown bool
	cfg      manualReaderConfig
}

type manualReaderConfig struct {
	temporality TemporalitySelector
	aggregation AggregationSelector
	limit       CardinalityLimitSelector
}

// ManualReaderOption configures a ManualReader.
type ManualReaderOption interface {
	applyManual(manualReaderConfig) manualReaderConfig
}

type manualReaderOptionFunc func(manualReaderConfig) manualReaderConfig

func (f manualReaderOptionFunc) applyManual(cfg manualReaderConfig) manualReaderConfig { return f(cfg) }

// NewManualReader returns a Reader that is collected directly.
func NewManualReader(options ...ManualReaderOption) *ManualReader {
	cfg := manualReaderConfig{
		temporality: DefaultTemporalitySelector,
		aggregation: DefaultAggregationSelector,
		limit:       func(InstrumentKind) (int, bool) { return 0, true },
	}
	for _, option := range options {
		cfg = option.applyManual(cfg)
	}
	return &ManualReader{cfg: cfg}
}

// WithTemporalitySelector configures a ManualReader temporality selector.
func WithTemporalitySelector(selector TemporalitySelector) ManualReaderOption {
	return manualReaderOptionFunc(func(cfg manualReaderConfig) manualReaderConfig {
		cfg.temporality = selector
		return cfg
	})
}

// WithAggregationSelector configures a ManualReader aggregation selector.
func WithAggregationSelector(selector AggregationSelector) ManualReaderOption {
	return manualReaderOptionFunc(func(cfg manualReaderConfig) manualReaderConfig {
		cfg.aggregation = selector
		return cfg
	})
}

// WithCardinalityLimitSelector configures per-kind cardinality limits.
func WithCardinalityLimitSelector(selector CardinalityLimitSelector) ManualReaderOption {
	return manualReaderOptionFunc(func(cfg manualReaderConfig) manualReaderConfig {
		cfg.limit = selector
		return cfg
	})
}

func (r *ManualReader) register(p *pipeline) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.pipeline == nil && !r.shutdown {
		r.pipeline = p
	}
}

func (r *ManualReader) temporality(kind InstrumentKind) metricdata.Temporality {
	return r.cfg.temporality(kind)
}

func (r *ManualReader) aggregation(kind InstrumentKind) Aggregation {
	return r.cfg.aggregation(kind)
}

func (r *ManualReader) cardinalityLimit(kind InstrumentKind) (int, bool) {
	return r.cfg.limit(kind)
}

// Collect stores the current collection in destination.
func (r *ManualReader) Collect(ctx context.Context, destination *metricdata.ResourceMetrics) error {
	if destination == nil {
		return errors.New("manual reader: *metricdata.ResourceMetrics is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.shutdown {
		return ErrReaderShutdown
	}
	if r.pipeline == nil {
		return ErrReaderNotRegistered
	}
	r.pipeline.collect(destination)
	return nil
}

// Shutdown releases the ManualReader.
func (r *ManualReader) Shutdown(context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.shutdown {
		return ErrReaderShutdown
	}
	r.shutdown = true
	r.pipeline = nil
	return nil
}
