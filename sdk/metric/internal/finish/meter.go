// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package finish

import "go.opentelemetry.io/otel/metric"

// Factory constructs Finish-aware synchronous instruments.
type Factory interface {
	Int64Counter(string, ...metric.Int64CounterOption) (metric.Int64Counter, error)
	Float64Counter(string, ...metric.Float64CounterOption) (metric.Float64Counter, error)
	Int64UpDownCounter(string, ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error)
	Float64UpDownCounter(string, ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error)
}

type meter struct {
	metric.Meter
	factory Factory
}

// NewMeter decorates base with Finish-aware instrument construction.
func NewMeter(base metric.Meter, factory Factory) metric.Meter {
	return &meter{Meter: base, factory: factory}
}

func (m *meter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return m.factory.Int64Counter(name, options...)
}

func (m *meter) Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	return m.factory.Float64Counter(name, options...)
}

func (m *meter) Int64UpDownCounter(
	name string,
	options ...metric.Int64UpDownCounterOption,
) (metric.Int64UpDownCounter, error) {
	return m.factory.Int64UpDownCounter(name, options...)
}

func (m *meter) Float64UpDownCounter(
	name string,
	options ...metric.Float64UpDownCounterOption,
) (metric.Float64UpDownCounter, error) {
	return m.factory.Float64UpDownCounter(name, options...)
}
