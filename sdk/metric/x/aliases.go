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
)
