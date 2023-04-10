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

package instrument // import "go.opentelemetry.io/otel/metric/instrument"

import "go.opentelemetry.io/otel/attribute"

// Observable is used as a grouping mechanism for all instruments that are
// updated within a Callback.
type Observable interface {
	observable()
}

// Option applies options to all instruments.
type Option interface {
	Int64CounterOption
	Int64UpDownCounterOption
	Int64HistogramOption
	Int64ObservableCounterOption
	Int64ObservableUpDownCounterOption
	Int64ObservableGaugeOption

	Float64CounterOption
	Float64UpDownCounterOption
	Float64HistogramOption
	Float64ObservableCounterOption
	Float64ObservableUpDownCounterOption
	Float64ObservableGaugeOption
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

func (o descOpt) applyFloat64ObservableCounter(c Float64ObservableCounterConfig) Float64ObservableCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64ObservableUpDownCounter(c Float64ObservableUpDownCounterConfig) Float64ObservableUpDownCounterConfig {
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

func (o descOpt) applyInt64ObservableCounter(c Int64ObservableCounterConfig) Int64ObservableCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64ObservableUpDownCounter(c Int64ObservableUpDownCounterConfig) Int64ObservableUpDownCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64ObservableGauge(c Int64ObservableGaugeConfig) Int64ObservableGaugeConfig {
	c.description = string(o)
	return c
}

// WithDescription sets the instrument description.
func WithDescription(desc string) Option { return descOpt(desc) }

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

func (o unitOpt) applyFloat64ObservableCounter(c Float64ObservableCounterConfig) Float64ObservableCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64ObservableUpDownCounter(c Float64ObservableUpDownCounterConfig) Float64ObservableUpDownCounterConfig {
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

func (o unitOpt) applyInt64ObservableCounter(c Int64ObservableCounterConfig) Int64ObservableCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64ObservableUpDownCounter(c Int64ObservableUpDownCounterConfig) Int64ObservableUpDownCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64ObservableGauge(c Int64ObservableGaugeConfig) Int64ObservableGaugeConfig {
	c.unit = string(o)
	return c
}

// WithUnit sets the instrument unit.
func WithUnit(u string) Option { return unitOpt(u) }

// Int64AddOption applies options to an addition measurement. See
// [MeasurementOption] for other options that can be used as a Int64AddOption.
type Int64AddOption interface {
	applyInt64Add(Int64AddConfig) Int64AddConfig
}

// Int64AddConfig contains options for an int64 addition measurement.
type Int64AddConfig struct {
	attrs attribute.Set
}

// NewInt64AddConfig returns a new [Int64AddConfig] with all opts applied.
func NewInt64AddConfig(opts []Int64AddOption) Int64AddConfig {
	config := Int64AddConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyInt64Add(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c Int64AddConfig) Attributes() attribute.Set {
	return c.attrs
}

// Float64AddOption applies options to an addition measurement. See
// [MeasurementOption] for other options that can be used as a
// Float64AddOption.
type Float64AddOption interface {
	applyFloat64Add(Float64AddConfig) Float64AddConfig
}

// Float64AddConfig contains options for an float64 addition measurement.
type Float64AddConfig struct {
	attrs attribute.Set
}

// NewFloat64AddConfig returns a new [Float64AddConfig] with all opts applied.
func NewFloat64AddConfig(opts []Float64AddOption) Float64AddConfig {
	config := Float64AddConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyFloat64Add(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c Float64AddConfig) Attributes() attribute.Set {
	return c.attrs
}

// Int64RecordOption applies options to an addition measurement. See
// [MeasurementOption] for other options that can be used as a
// Int64RecordOption.
type Int64RecordOption interface {
	applyInt64Record(Int64RecordConfig) Int64RecordConfig
}

// Int64RecordConfig contains options for an int64 recorded measurement.
type Int64RecordConfig struct {
	attrs attribute.Set
}

// NewInt64RecordConfig returns a new [Int64RecordConfig] with all opts
// applied.
func NewInt64RecordConfig(opts []Int64RecordOption) Int64RecordConfig {
	config := Int64RecordConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyInt64Record(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c Int64RecordConfig) Attributes() attribute.Set {
	return c.attrs
}

// Float64RecordOption applies options to an addition measurement. See
// [MeasurementOption] for other options that can be used as a
// Float64RecordOption.
type Float64RecordOption interface {
	applyFloat64Record(Float64RecordConfig) Float64RecordConfig
}

// Float64RecordConfig contains options for an float64 recorded measurement.
type Float64RecordConfig struct {
	attrs attribute.Set
}

// NewFloat64RecordConfig returns a new [Float64RecordConfig] with all opts
// applied.
func NewFloat64RecordConfig(opts []Float64RecordOption) Float64RecordConfig {
	config := Float64RecordConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyFloat64Record(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c Float64RecordConfig) Attributes() attribute.Set {
	return c.attrs
}

// Int64ObserveOption applies options to an addition measurement. See
// [MeasurementOption] for other options that can be used as a
// Int64ObserveOption.
type Int64ObserveOption interface {
	applyInt64Observe(Int64ObserveConfig) Int64ObserveConfig
}

// Int64ObserveConfig contains options for an int64 observed measurement.
type Int64ObserveConfig struct {
	attrs attribute.Set
}

// NewInt64ObserveConfig returns a new [Int64ObserveConfig] with all opts
// applied.
func NewInt64ObserveConfig(opts []Int64ObserveOption) Int64ObserveConfig {
	config := Int64ObserveConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyInt64Observe(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c Int64ObserveConfig) Attributes() attribute.Set {
	return c.attrs
}

// Float64ObserveOption applies options to an addition measurement. See
// [MeasurementOption] for other options that can be used as a
// Float64ObserveOption.
type Float64ObserveOption interface {
	applyFloat64Observe(Float64ObserveConfig) Float64ObserveConfig
}

// Float64ObserveConfig contains options for an float64 observed measurement.
type Float64ObserveConfig struct {
	attrs attribute.Set
}

// NewFloat64ObserveConfig returns a new [Float64ObserveConfig] with all opts
// applied.
func NewFloat64ObserveConfig(opts []Float64ObserveOption) Float64ObserveConfig {
	config := Float64ObserveConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyFloat64Observe(config)
	}
	return config
}

// Attributes returns the configured attribute set.
func (c Float64ObserveConfig) Attributes() attribute.Set {
	return c.attrs
}

// MeasurementOption applies options to all instrument measurement.
type MeasurementOption interface {
	Int64AddOption
	Float64AddOption
	Int64RecordOption
	Float64RecordOption
	Int64ObserveOption
	Float64ObserveOption
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

func (o attrOpt) applyInt64Add(c Int64AddConfig) Int64AddConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyFloat64Add(c Float64AddConfig) Float64AddConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyInt64Record(c Int64RecordConfig) Int64RecordConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyFloat64Record(c Float64RecordConfig) Float64RecordConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyInt64Observe(c Int64ObserveConfig) Int64ObserveConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyFloat64Observe(c Float64ObserveConfig) Float64ObserveConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

// WithAttributeSet sets the attribute Set associated with a measurement is
// made with.
//
// If multiple WithAttributeSet or WithAttributes options are passed the
// attributes will be merged together in the order they are passed. Attributes
// with duplicate keys will use the last value passed.
func WithAttributeSet(attributes attribute.Set) MeasurementOption {
	return attrOpt{set: attributes}
}

// WithAttributeSet converts attributes into an attribute Set and sets the Set
// to be associated with a measurement. This is shorthand for:
//
//	WithAttributes(attribute.NewSet(attributes...))
//
// See [WithAttributeSet] for how multiple WithAttributes are merged.
func WithAttributes(attributes ...attribute.KeyValue) MeasurementOption {
	cp := make([]attribute.KeyValue, len(attributes))
	copy(cp, attributes)
	return attrOpt{set: attribute.NewSet(cp...)}
}
