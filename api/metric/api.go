// Copyright 2019, OpenTelemetry Authors
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

package metric

import (
	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
	"github.com/open-telemetry/opentelemetry-go/api/unit"
)

type (
	Metric interface {
		Measure() core.Measure

		DefinitionID() core.EventID

		Type() MetricType
		Fields() []core.Key
		Err() error

		base() *baseMetric
	}

	MetricType int
)

const (
	Invalid MetricType = iota
	GaugeInt64
	GaugeFloat64
	DerivedGaugeInt64
	DerivedGaugeFloat64
	CumulativeInt64
	CumulativeFloat64
	DerivedCumulativeInt64
	DerivedCumulativeFloat64
)

type (
	Option func(*baseMetric, *[]tag.Option)
)

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return func(_ *baseMetric, to *[]tag.Option) {
		*to = append(*to, tag.WithDescription(desc))
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return func(_ *baseMetric, to *[]tag.Option) {
		*to = append(*to, tag.WithUnit(unit))
	}
}

// WithKeys applies the provided dimension keys.
func WithKeys(keys ...core.Key) Option {
	return func(bm *baseMetric, _ *[]tag.Option) {
		bm.keys = keys
	}
}
