// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/x"

import sdkmetric "go.opentelemetry.io/otel/sdk/metric"

// Aggregation and the related aliases reuse stable SDK configuration types.
type (
	Aggregation = sdkmetric.Aggregation
	// AggregationDefault uses the stable default aggregation configuration.
	AggregationDefault = sdkmetric.AggregationDefault
	// AggregationDrop uses the stable drop aggregation configuration.
	AggregationDrop = sdkmetric.AggregationDrop
	// AggregationSum uses the stable Sum aggregation configuration.
	AggregationSum = sdkmetric.AggregationSum
	// AggregationLastValue uses the stable LastValue aggregation configuration.
	AggregationLastValue = sdkmetric.AggregationLastValue
	// AggregationExplicitBucketHistogram uses the stable explicit histogram configuration.
	AggregationExplicitBucketHistogram = sdkmetric.AggregationExplicitBucketHistogram
	// AggregationBase2ExponentialHistogram uses the stable exponential histogram configuration.
	AggregationBase2ExponentialHistogram = sdkmetric.AggregationBase2ExponentialHistogram
	// Instrument uses the stable SDK instrument descriptor.
	Instrument = sdkmetric.Instrument
	// InstrumentKind uses the stable SDK instrument kind.
	InstrumentKind = sdkmetric.InstrumentKind
	// Stream uses the stable SDK stream descriptor.
	Stream = sdkmetric.Stream
	// View uses the stable SDK view contract.
	View = sdkmetric.View
	// ExemplarReservoirProviderSelector uses the stable selector contract.
	ExemplarReservoirProviderSelector = sdkmetric.ExemplarReservoirProviderSelector
	// Option configures an experimental MeterProvider.
	Option = sdkmetric.Option
	// Reader is the stable SDK Reader contract used by the experimental facade.
	Reader = sdkmetric.Reader
	// ManualReader collects metric data on demand.
	ManualReader = sdkmetric.ManualReader
	// ManualReaderOption configures a ManualReader.
	ManualReaderOption = sdkmetric.ManualReaderOption
	// TemporalitySelector selects temporality for an instrument kind.
	TemporalitySelector = sdkmetric.TemporalitySelector
	// AggregationSelector selects aggregation for an instrument kind.
	AggregationSelector = sdkmetric.AggregationSelector
	// CardinalityLimitSelector selects cardinality for an instrument kind.
	CardinalityLimitSelector = sdkmetric.CardinalityLimitSelector
)

// InstrumentKindCounter and the related constants reuse stable instrument kinds.
const (
	InstrumentKindCounter = sdkmetric.InstrumentKindCounter
	// InstrumentKindUpDownCounter identifies a synchronous UpDownCounter.
	InstrumentKindUpDownCounter = sdkmetric.InstrumentKindUpDownCounter
	// InstrumentKindHistogram identifies a synchronous Histogram.
	InstrumentKindHistogram = sdkmetric.InstrumentKindHistogram
	// InstrumentKindObservableCounter identifies an asynchronous Counter.
	InstrumentKindObservableCounter = sdkmetric.InstrumentKindObservableCounter
	// InstrumentKindObservableUpDownCounter identifies an asynchronous UpDownCounter.
	InstrumentKindObservableUpDownCounter = sdkmetric.InstrumentKindObservableUpDownCounter
	// InstrumentKindObservableGauge identifies an asynchronous Gauge.
	InstrumentKindObservableGauge = sdkmetric.InstrumentKindObservableGauge
	// InstrumentKindGauge identifies a synchronous Gauge.
	InstrumentKindGauge = sdkmetric.InstrumentKindGauge
)

// NewView and the related helpers reuse stable SDK configuration functions.
var (
	NewView = sdkmetric.NewView
	// DefaultAggregationSelector uses the stable SDK aggregation defaults.
	DefaultAggregationSelector = sdkmetric.DefaultAggregationSelector
	// DefaultExemplarReservoirProviderSelector uses the stable SDK exemplar defaults.
	DefaultExemplarReservoirProviderSelector = sdkmetric.DefaultExemplarReservoirProviderSelector
	// NewManualReader returns a stable ManualReader for the experimental provider.
	NewManualReader = sdkmetric.NewManualReader
	// WithTemporalitySelector configures reader temporality.
	WithTemporalitySelector = sdkmetric.WithTemporalitySelector
	// WithAggregationSelector configures reader aggregation.
	WithAggregationSelector = sdkmetric.WithAggregationSelector
	// WithCardinalityLimitSelector configures reader cardinality.
	WithCardinalityLimitSelector = sdkmetric.WithCardinalityLimitSelector
	// CumulativeTemporalitySelector selects cumulative temporality.
	CumulativeTemporalitySelector = sdkmetric.CumulativeTemporalitySelector
	// DeltaTemporalitySelector selects delta temporality.
	DeltaTemporalitySelector = sdkmetric.DeltaTemporalitySelector
	// DefaultTemporalitySelector selects the default temporality.
	DefaultTemporalitySelector = sdkmetric.DefaultTemporalitySelector
	// WithResource associates a resource with a MeterProvider.
	WithResource = sdkmetric.WithResource
	// WithReader associates a Reader with a MeterProvider.
	WithReader = sdkmetric.WithReader
	// WithView associates views with a MeterProvider.
	WithView = sdkmetric.WithView
	// WithExemplarFilter configures exemplar filtering.
	WithExemplarFilter = sdkmetric.WithExemplarFilter
	// WithCardinalityLimit configures the provider cardinality limit.
	WithCardinalityLimit = sdkmetric.WithCardinalityLimit
)

var (
	// ErrReaderNotRegistered is returned before a Reader is registered.
	ErrReaderNotRegistered = sdkmetric.ErrReaderNotRegistered
	// ErrReaderShutdown is returned after a Reader is shut down.
	ErrReaderShutdown = sdkmetric.ErrReaderShutdown
)
