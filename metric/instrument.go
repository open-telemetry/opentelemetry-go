// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/metric"

import (
	"go.opentelemetry.io/otel/attribute"
)

// Observable is used as a grouping mechanism for all instruments that are
// updated within a Callback.
type Observable interface {
	observable()
}

// InstrumentOption applies options to all instruments.
type InstrumentOption interface {
	Int64CounterOption
	Int64UpDownCounterOption
	Int64HistogramOption
	Int64GaugeOption
	Int64ObservableCounterOption
	Int64ObservableUpDownCounterOption
	Int64ObservableGaugeOption

	Float64CounterOption
	Float64UpDownCounterOption
	Float64HistogramOption
	Float64GaugeOption
	Float64ObservableCounterOption
	Float64ObservableUpDownCounterOption
	Float64ObservableGaugeOption
}

// HistogramOption applies options to histogram instruments.
type HistogramOption interface {
	Int64HistogramOption
	Float64HistogramOption
}

type descOpt string

func (o descOpt) applyFloat64Counter(c Float64CounterConfig) Float64CounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64UpDownCounter(c Float64UpDownCounterConfig) Float64UpDownCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64Histogram(c Float64HistogramConfig) Float64HistogramConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64Gauge(c Float64GaugeConfig) Float64GaugeConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64ObservableCounter(c Float64ObservableCounterConfig) Float64ObservableCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64ObservableUpDownCounter(
	c Float64ObservableUpDownCounterConfig,
) Float64ObservableUpDownCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64ObservableGauge(c Float64ObservableGaugeConfig) Float64ObservableGaugeConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64Counter(c Int64CounterConfig) Int64CounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64UpDownCounter(c Int64UpDownCounterConfig) Int64UpDownCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64Histogram(c Int64HistogramConfig) Int64HistogramConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64Gauge(c Int64GaugeConfig) Int64GaugeConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64ObservableCounter(c Int64ObservableCounterConfig) Int64ObservableCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64ObservableUpDownCounter(
	c Int64ObservableUpDownCounterConfig,
) Int64ObservableUpDownCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64ObservableGauge(c Int64ObservableGaugeConfig) Int64ObservableGaugeConfig {
	c.description = string(o)
	return c
}

// WithDescription sets the instrument description.
func WithDescription(desc string) InstrumentOption { return descOpt(desc) }

type unitOpt string

func (o unitOpt) applyFloat64Counter(c Float64CounterConfig) Float64CounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64UpDownCounter(c Float64UpDownCounterConfig) Float64UpDownCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64Histogram(c Float64HistogramConfig) Float64HistogramConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64Gauge(c Float64GaugeConfig) Float64GaugeConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64ObservableCounter(c Float64ObservableCounterConfig) Float64ObservableCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64ObservableUpDownCounter(
	c Float64ObservableUpDownCounterConfig,
) Float64ObservableUpDownCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64ObservableGauge(c Float64ObservableGaugeConfig) Float64ObservableGaugeConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64Counter(c Int64CounterConfig) Int64CounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64UpDownCounter(c Int64UpDownCounterConfig) Int64UpDownCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64Histogram(c Int64HistogramConfig) Int64HistogramConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64Gauge(c Int64GaugeConfig) Int64GaugeConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64ObservableCounter(c Int64ObservableCounterConfig) Int64ObservableCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64ObservableUpDownCounter(
	c Int64ObservableUpDownCounterConfig,
) Int64ObservableUpDownCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64ObservableGauge(c Int64ObservableGaugeConfig) Int64ObservableGaugeConfig {
	c.unit = string(o)
	return c
}

// WithUnit sets the instrument unit.
//
// The unit u should be defined using the appropriate [UCUM](https://ucum.org) case-sensitive code.
func WithUnit(u string) InstrumentOption { return unitOpt(u) }

// WithExplicitBucketBoundaries sets the instrument explicit bucket boundaries.
//
// This option is considered "advisory", and may be ignored by API implementations.
func WithExplicitBucketBoundaries(bounds ...float64) HistogramOption { return bucketOpt(bounds) }

type bucketOpt []float64

func (o bucketOpt) applyFloat64Histogram(c Float64HistogramConfig) Float64HistogramConfig {
	c.explicitBucketBoundaries = o
	return c
}

func (o bucketOpt) applyInt64Histogram(c Int64HistogramConfig) Int64HistogramConfig {
	c.explicitBucketBoundaries = o
	return c
}

// AddOption applies options to an addition measurement. See
// [MeasurementOption] for other options that can be used as an AddOption.
type AddOption interface {
	applyAdd(AddConfig) AddConfig
}

// AddConfig contains options for an addition measurement.
type AddConfig struct {
	attrs attribute.Set
}

// NewAddConfig returns a new [AddConfig] with all opts applied.
func NewAddConfig(opts []AddOption) AddConfig {
	config := AddConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		if _, ok := o.(experimentalOption); ok {
			continue
		}
		config = o.applyAdd(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c AddConfig) Attributes() attribute.Set {
	return c.attrs
}

// RecordOption applies options to an addition measurement. See
// [MeasurementOption] for other options that can be used as a RecordOption.
type RecordOption interface {
	applyRecord(RecordConfig) RecordConfig
}

// RecordConfig contains options for a recorded measurement.
type RecordConfig struct {
	attrs attribute.Set
}

// NewRecordConfig returns a new [RecordConfig] with all opts applied.
func NewRecordConfig(opts []RecordOption) RecordConfig {
	config := RecordConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		if _, ok := o.(experimentalOption); ok {
			continue
		}
		config = o.applyRecord(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c RecordConfig) Attributes() attribute.Set {
	return c.attrs
}

// ObserveOption applies options to an addition measurement. See
// [MeasurementOption] for other options that can be used as a ObserveOption.
type ObserveOption interface {
	applyObserve(ObserveConfig) ObserveConfig
}

// ObserveConfig contains options for an observed measurement.
type ObserveConfig struct {
	attrs attribute.Set
}

// NewObserveConfig returns a new [ObserveConfig] with all opts applied.
func NewObserveConfig(opts []ObserveOption) ObserveConfig {
	config := ObserveConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		if _, ok := o.(experimentalOption); ok {
			continue
		}
		config = o.applyObserve(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c ObserveConfig) Attributes() attribute.Set {
	return c.attrs
}

// MeasurementOption applies options to all instrument measurement.
type MeasurementOption interface {
	AddOption
	RecordOption
	ObserveOption
}

// FinishOption applies options to a synchronous instrument finish operation.
type FinishOption interface {
	applyFinish(FinishConfig) FinishConfig
}

// MatchAttributesFunc returns whether attrs matches a finish selector.
type MatchAttributesFunc func(attribute.Set) bool

// AttributeSetOption applies attribute options to measurement and finish calls.
type AttributeSetOption interface {
	MeasurementOption
	FinishOption
}

// FinishConfig contains options for a finish operation.
type FinishConfig struct {
	attrs      attribute.Set
	hasAttrs   bool
	matchAttrs MatchAttributesFunc
}

// NewFinishConfig returns a new [FinishConfig] with all opts applied.
func NewFinishConfig(opts []FinishOption) FinishConfig {
	config := FinishConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyFinish(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c FinishConfig) Attributes() attribute.Set {
	return c.attrs
}

// HasAttributes reports whether an exact attribute set selector was configured.
func (c FinishConfig) HasAttributes() bool {
	return c.hasAttrs
}

// MatchAttributes returns the configured finish matcher, if any.
func (c FinishConfig) MatchAttributes() MatchAttributesFunc {
	return c.matchAttrs
}

// Matcher returns the configured finish selector as a predicate.
func (c FinishConfig) Matcher() MatchAttributesFunc {
	switch {
	case c.hasAttrs && c.matchAttrs != nil:
		return func(attrs attribute.Set) bool {
			return attrs.Equals(&c.attrs) && c.matchAttrs(attrs)
		}
	case c.hasAttrs:
		return func(attrs attribute.Set) bool {
			return attrs.Equals(&c.attrs)
		}
	case c.matchAttrs != nil:
		return c.matchAttrs
	default:
		empty := *attribute.EmptySet()
		return func(attrs attribute.Set) bool {
			return attrs.Equals(&empty)
		}
	}
}

type attrOpt struct {
	set attribute.Set
}

// mergeSets returns the union of keys between a and b. Any duplicate keys will
// use the value associated with b.
func mergeSets(a, b attribute.Set) attribute.Set {
	// NewMergeIterator uses the first value for any duplicates.
	iter := attribute.NewMergeIterator(&b, &a)
	merged := make([]attribute.KeyValue, 0, a.Len()+b.Len())
	for iter.Next() {
		merged = append(merged, iter.Attribute())
	}
	return attribute.NewSet(merged...)
}

func (o attrOpt) applyAdd(c AddConfig) AddConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyRecord(c RecordConfig) RecordConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyObserve(c ObserveConfig) ObserveConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyFinish(c FinishConfig) FinishConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
		c.hasAttrs = true
	default:
		c.attrs = mergeSets(c.attrs, o.set)
		c.hasAttrs = true
	}
	return c
}

type matchAttrOpt struct {
	match MatchAttributesFunc
}

func (o matchAttrOpt) applyFinish(c FinishConfig) FinishConfig {
	if o.match == nil {
		return c
	}
	if c.matchAttrs == nil {
		c.matchAttrs = o.match
		return c
	}
	prev := c.matchAttrs
	c.matchAttrs = func(attrs attribute.Set) bool {
		return prev(attrs) && o.match(attrs)
	}
	return c
}

// WithAttributeSet sets the attribute Set associated with a measurement is
// made with.
//
// If multiple WithAttributeSet or WithAttributes options are passed the
// attributes will be merged together in the order they are passed. Attributes
// with duplicate keys will use the last value passed.
func WithAttributeSet(attributes attribute.Set) AttributeSetOption {
	return attrOpt{set: attributes}
}

// WithAttributes converts attributes into an attribute Set and sets the Set to
// be associated with a measurement. This is shorthand for:
//
//	cp := make([]attribute.KeyValue, len(attributes))
//	copy(cp, attributes)
//	WithAttributeSet(attribute.NewSet(cp...))
//
// [attribute.NewSet] may modify the passed attributes so this will make a copy
// of attributes before creating a set in order to ensure this function is
// concurrent safe. This makes this option function less optimized in
// comparison to [WithAttributeSet]. Therefore, [WithAttributeSet] should be
// preferred for performance sensitive code.
//
// See [WithAttributeSet] for information about how multiple WithAttributes are
// merged.
func WithAttributes(attributes ...attribute.KeyValue) AttributeSetOption {
	cp := make([]attribute.KeyValue, len(attributes))
	copy(cp, attributes)
	return attrOpt{set: attribute.NewSet(cp...)}
}

// WithMatchAttributes sets a finish matcher that decides whether a series
// should be finished based on its attributes.
//
// This option is only used by Finish calls. If multiple matchers are passed,
// all of them need to match for a series to be finished.
//
// If combined with WithAttributeSet or WithAttributes, both the exact
// attribute selector and the matcher need to match.
func WithMatchAttributes(match MatchAttributesFunc) FinishOption {
	return matchAttrOpt{match: match}
}
