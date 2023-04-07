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

// Option applies options to all instruments.
type Option[N int64 | float64] interface {
	CounterOption[N]
	UpDownCounterOption[N]
	HistogramOption[N]
	ObservableCounterOption[N]
	ObservableUpDownCounterOption[N]
	ObservableGaugeOption[N]
}

type descOpt[N int64 | float64] struct {
	description string
}

func (o descOpt[N]) applyCounter(c CounterConfig[N]) CounterConfig[N] { // nolint: unused
	c.description = o.description
	return c
}

func (o descOpt[N]) applyUpDownCounter(c UpDownCounterConfig[N]) UpDownCounterConfig[N] { // nolint: unused
	c.description = o.description
	return c
}

func (o descOpt[N]) applyHistogram(c HistogramConfig[N]) HistogramConfig[N] { // nolint: unused
	c.description = o.description
	return c
}

func (o descOpt[N]) applyObservableCounter(c ObservableCounterConfig[N]) ObservableCounterConfig[N] { // nolint: unused
	c.description = o.description
	return c
}

func (o descOpt[N]) applyObservableUpDownCounter(c ObservableUpDownCounterConfig[N]) ObservableUpDownCounterConfig[N] { // nolint: unused
	c.description = o.description
	return c
}

func (o descOpt[N]) applyObservableGauge(c ObservableGaugeConfig[N]) ObservableGaugeConfig[N] { // nolint: unused
	c.description = o.description
	return c
}

// WithDescription sets the instrument description.
func WithDescription[N int64 | float64](desc string) Option[N] { return descOpt[N]{desc} }

type unitOpt[N int64 | float64] struct {
	unit string
}

func (o unitOpt[N]) applyCounter(c CounterConfig[N]) CounterConfig[N] { // nolint: unused
	c.unit = o.unit
	return c
}

func (o unitOpt[N]) applyUpDownCounter(c UpDownCounterConfig[N]) UpDownCounterConfig[N] { // nolint: unused
	c.unit = o.unit
	return c
}

func (o unitOpt[N]) applyHistogram(c HistogramConfig[N]) HistogramConfig[N] { // nolint: unused
	c.unit = o.unit
	return c
}

func (o unitOpt[N]) applyObservableCounter(c ObservableCounterConfig[N]) ObservableCounterConfig[N] { // nolint: unused
	c.unit = o.unit
	return c
}

func (o unitOpt[N]) applyObservableUpDownCounter(c ObservableUpDownCounterConfig[N]) ObservableUpDownCounterConfig[N] { // nolint: unused
	c.unit = o.unit
	return c
}

func (o unitOpt[N]) applyObservableGauge(c ObservableGaugeConfig[N]) ObservableGaugeConfig[N] { // nolint: unused
	c.unit = o.unit
	return c
}

// WithUnit sets the instrument unit.
func WithUnit[N int64 | float64](u string) Option[N] { return unitOpt[N]{u} }
