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

package reader

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	// Config includes Reader-specific options that the user
	// configures, which are generally per-instrument-kind defaults.
	Config struct {
		DefaultAggregationKindFunc
		DefaultAggregationTemporalityFunc
		DefaultAggregationConfigFunc
	}

	// Option is a functional option pattern for generating Configs.
	Option func(*Config)

	// DefaultAggregationKindFunc is a per-instrument-kind aggregation.Kind choice.
	DefaultAggregationKindFunc func(sdkinstrument.Kind) aggregation.Kind

	// DefaultAggregationKindFunc is a per-instrument-kind aggregation.Temporality choice.
	DefaultAggregationTemporalityFunc func(sdkinstrument.Kind) aggregation.Temporality

	// DefaultAggregationKindFunc is a per-instrument-kind, per-number-kind aggregator.Config choice.
	DefaultAggregationConfigFunc func(sdkinstrument.Kind) (int64Config, float64Config aggregator.Config)

	// ReaderConfig represents the connection between an Exporter and
	// the MeterProvider.  Readers give the internal View compiler
	// all the necessary information to construct a metrics pipeline.
	ReaderConfig struct {
		// reader is the output terminal of a metrics pipeline.
		reader Reader

		// These four arrays are the evaluted
		// per-instrument-kind choices described in the
		// Config.

		defAggr   [sdkinstrument.NumKinds]aggregation.Kind
		defTempo  [sdkinstrument.NumKinds]aggregation.Temporality
		defI64Cfg [sdkinstrument.NumKinds]aggregator.Config
		defF64Cfg [sdkinstrument.NumKinds]aggregator.Config
	}

	// Sequence provides the three relevant timestamps to the Reader
	Sequence struct {
		// Start is the time when the MeterProvider was initialized.
		Start time.Time
		// Last is the time when the previous collection
		// happened.  If there was no previous collection,
		// this will match Start.
		Last time.Time
		// Now is the moment the current collection began.  This value
		// will be used as the subsequent value for Last.
		Now time.Time
	}

	// Metrics is the result of a single collection.
	//
	// This struct supports re-use of the nested memory structure
	// in its Scopes slice such that repeated calls Produce will
	// not reallocate the same quantity of memory again and again.
	//
	// To re-use the memory from a previous Metrics value, pass a
	// pointer pointer to Produce(). This is safe for push-based
	// exporters that perform sequential collection.
	Metrics struct {
		// Resource is the MeterProvider's configured Resource.
		Resource *resource.Resource

		// Scopes is a slice of metric data, one per Meter.
		Scopes []Scope
	}

	// Scope is the result of a single collection for a single Meter.
	//
	// See the comments on Metrics about re-use of slices in this struct.
	Scope struct {
		// Library describes the instrumentation scope.
		Library instrumentation.Library

		// Instruments is a slice of metric data, one per Instrument
		// in the scope.
		Instruments []Instrument
	}

	// Instrument is the result of a single collection for a single Instrument.
	//
	// See the comments on Metrics about re-use of slices in this struct.
	Instrument struct {
		// Descriptor describes an instrument created through a View,
		// including name, unit, description, instrument and number kinds.
		Descriptor sdkinstrument.Descriptor

		// Temporality is the temporality preference used to
		// generate this point.  For Gauge instruments, which
		// do not support aggregation temporality, this
		// setting indicates whether to maintain state of
		// prior gauge values (Cumulative) or only to report
		// gauge values that change (Delta).
		Temporality aggregation.Temporality

		// Points is a slice of metric data, one per attribute.Set value.
		Points []Point
	}

	// Point is a timeseries data point resulting from a single collection.
	Point struct {
		// Attributes are the coordinates of this series.
		Attributes attribute.Set

		// Aggregation determines the kind of data point
		// recorded in this series.
		Aggregation aggregation.Aggregation

		// Start indicates the start of the collection
		// interval reflected in this series, which is set
		// according to the configured temporality.
		Start time.Time

		// End indicates the moment at which the collection
		// was performed.
		End time.Time
	}

	// Reader is the consumer of metrics associated with a single Reader.
	Reader interface {
		// Registeree supports connecting the MeterProvider
		// with the Exporter.
		Registeree

		// Flush is called when ... ? (TODO)
		Flush(context.Context) error

		// Shutdown is called when the SDK is shutting down. (TODO)
		Shutdown(context.Context) error
	}

	// Registeree is the mechanism for binding a MeterProvider
	// with an Exporter during initialization.
	Registeree interface {
		// Register is called exactly during MeterProvider
		// initialization.  The supplied Producer gives the
		// Exporter and/or its Reader the interface used to
		// collect metrics.
		Register(Producer)
	}

	// Producer is the interface used to perform collection.
	Producer interface {
		// Produce returns metrics from a single collection.
		//
		// Produce may be called concurrently,
		//
		// The `in` parameter supports re-use of memory from
		// one collection to the next.  Callers that pass `in`
		// will write metrics into the same slices and structs.
		//
		// When `in` is nil, a new Metrics object is returned.
		Produce(ctx context.Context, in *Metrics) Metrics
	}
)

// WithDefaultAggregationKindFunc configures the default
// aggregation.Kind to use with each kind of instrument.
func WithDefaultAggregationKindFunc(d DefaultAggregationKindFunc) Option {
	return func(cfg *Config) {
		cfg.DefaultAggregationKindFunc = d
	}
}

// WithDefaultAggregationTemporalityFunc configures the default
// aggregation.Temporality to use with each kind of instrument.
func WithDefaultAggregationTemporalityFunc(d DefaultAggregationTemporalityFunc) Option {
	return func(cfg *Config) {
		cfg.DefaultAggregationTemporalityFunc = d
	}
}

// WithDefaultAggregationConfigFunc configures the default
// aggregator.Config to use with each kind of instrument.
func WithDefaultAggregationConfigFunc(d DefaultAggregationConfigFunc) Option {
	return func(cfg *Config) {
		cfg.DefaultAggregationConfigFunc = d
	}
}

// StandardAggregation is the specified default aggregation Kind for
// each instrument Kind.
func StandardAggregationKind(ik sdkinstrument.Kind) aggregation.Kind {
	switch ik {
	case sdkinstrument.HistogramKind:
		return aggregation.HistogramKind
	case sdkinstrument.GaugeObserverKind:
		return aggregation.GaugeKind
	default:
		return aggregation.SumKind
	}
}

// StandardTemporality returns the specified default Cumulative
// temporality for all instrument kinds.
func StandardTemporality(ik sdkinstrument.Kind) aggregation.Temporality {
	return aggregation.CumulativeTemporality
}

// DeltaPreferredTemporality returns the specified Delta temporality
// for all instrument kinds except UpDownCounter, which remain Cumulative.
func DeltaPreferredTemporality(ik sdkinstrument.Kind) aggregation.Temporality {
	switch ik {
	case sdkinstrument.UpDownCounterKind, sdkinstrument.UpDownCounterObserverKind:
		return aggregation.CumulativeTemporality
	default:
		return aggregation.DeltaTemporality
	}
}

// StandardConfig returns two default aggregator.Configs.
func StandardConfig(ik sdkinstrument.Kind) (ints, floats aggregator.Config) {
	return aggregator.Config{}, aggregator.Config{}
}

// NewConfig returns a new Reader configured for `exporter` with provided
// optional configuration.
func NewConfig(exporter Reader, opts ...Option) *ReaderConfig {
	cfg := Config{
		DefaultAggregationKindFunc:        StandardAggregationKind,
		DefaultAggregationTemporalityFunc: StandardTemporality,
		DefaultAggregationConfigFunc:      StandardConfig,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	r := &ReaderConfig{
		reader: exporter,
	}
	for i := sdkinstrument.Kind(0); i < sdkinstrument.NumKinds; i++ {
		r.defAggr[i] = cfg.DefaultAggregationKindFunc(i)
		r.defTempo[i] = cfg.DefaultAggregationTemporalityFunc(i)
		r.defI64Cfg[i], r.defF64Cfg[i] = cfg.DefaultAggregationConfigFunc(i)
	}
	return r
}

// DefaultAggregation returns the default aggregation.Kind for each instrument kind.
func (r *ReaderConfig) DefaultAggregation(k sdkinstrument.Kind) aggregation.Kind {
	return r.defAggr[k]
}

// DefaultTemporality returns the default aggregation.Temporality for each instrument kind.
func (r *ReaderConfig) DefaultTemporality(k sdkinstrument.Kind) aggregation.Temporality {
	return r.defTempo[k]
}

// DefaultAggregationConfig returns the default aggregation.Temporality for each instrument kind.
func (r *ReaderConfig) DefaultAggregationConfig(k sdkinstrument.Kind, nk number.Kind) aggregator.Config {
	if nk == number.Int64Kind {
		return r.defI64Cfg[k]
	}
	return r.defF64Cfg[k]
}

// Reader returns the Reader's associated Reader.
func (r *ReaderConfig) Reader() Reader {
	return r.reader
}

func (m *Metrics) Reset() {
	resetScopes(&m.Scopes)
}

func (s *Scope) Reset() {
	resetInstruments(&s.Instruments)
}

func (in *Instrument) Reset() {
	resetPoints(&in.Points)
}

func resetScopes(ss *[]Scope) {
	for i := range *ss {
		(*ss)[i].Reset()
	}
	(*ss) = (*ss)[0:0:cap((*ss))]
}

func resetInstruments(is *[]Instrument) {
	for i := range *is {
		(*is)[i].Reset()
	}
	(*is) = (*is)[0:0:cap((*is))]
}

func resetPoints(ps *[]Point) {
	(*ps) = (*ps)[0:0:cap((*ps))]
}

func Reallocate[T any](p *[]T) *T {
	if len(*p) < cap(*p) {
		(*p) = (*p)[0 : len(*p)+1 : cap(*p)]
	} else {
		var empty T
		(*p) = append(*p, empty)
	}
	return &(*p)[len(*p)-1]
}

type Exporter interface {
	Export(context.Context, Metrics) error
	Flush(context.Context) error
	Shutdown(context.Context) error
}

type ManualReader struct {
	lock     sync.Mutex
	producer Producer

	exporter Exporter
}

func NewManualReader(exp Exporter) *ManualReader {
	return &ManualReader{
		exporter: exp,
	}
}

func (m *ManualReader) Collect(ctx context.Context, in *Metrics) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.producer == nil || m.exporter == nil {
		return nil
	}

	return m.exporter.Export(ctx, m.producer.Produce(ctx, in))
}

func (m *ManualReader) Flush(ctx context.Context) error {
	if m.exporter == nil {
		return nil
	}
	return m.exporter.Flush(ctx)
}
func (m *ManualReader) Shutdown(ctx context.Context) error {
	if m.exporter == nil {
		return nil
	}
	return m.exporter.Shutdown(ctx)
}
func (m *ManualReader) Register(prod Producer) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.producer = prod
}
