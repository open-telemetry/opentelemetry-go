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
	"context"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/unit"
)

type MetricType int

const (
	Invalid    MetricType = iota
	Gauge                 // Supports Set()
	Cumulative            // Supports Inc()
)

type Meter interface {
	// TODO more Metric types
	GetFloat64Gauge(ctx context.Context, gauge *Float64GaugeHandle, labels ...core.KeyValue) Float64Gauge
}

type Float64Gauge interface {
	Set(ctx context.Context, value float64, labels ...core.KeyValue)
}

type Handle struct {
	Name        string
	Description string
	Unit        unit.Unit

	Type MetricType
	Keys []core.Key
}

type Option func(*Handle)

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return func(m *Handle) {
		m.Description = desc
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return func(m *Handle) {
		m.Unit = unit
	}
}

// WithKeys applies the provided dimension keys.
func WithKeys(keys ...core.Key) Option {
	return func(m *Handle) {
		m.Keys = keys
	}
}

func (mtype MetricType) String() string {
	switch mtype {
	case Gauge:
		return "gauge"
	case Cumulative:
		return "cumulative"
	default:
		return "unknown"
	}
}
