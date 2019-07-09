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
	"sync/atomic"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/registry"
	"github.com/open-telemetry/opentelemetry-go/api/unit"
)

type MetricType int

const (
	Invalid    MetricType = iota
	Gauge                 // Supports Set()
	Cumulative            // Supports Inc()
)

type Meter interface {
	// TODO more Metric types
	GetFloat64Gauge(gauge *Float64GaugeRegistration, labels ...core.KeyValue) Float64Gauge
}

type Float64Gauge interface {
	Set(ctx context.Context, value float64, labels ...core.KeyValue)
}

type Registration struct {
	Variable registry.Variable

	Type MetricType
	Keys []core.Key
}

type noopMeter struct {
}

var (
	global atomic.Value
)

// GlobalMeter return meter registered with global registry.
// If no meter is registered then an instance of noop Meter is returned.
func GlobalMeter() Meter {
	if t := global.Load(); t != nil {
		return t.(Meter)
	}
	return noopMeter{}
}

// SetGlobalMeter sets provided meter as a global meter.
func SetGlobalMeter(t Meter) {
	global.Store(t)
}

type Option func(*Registration, *[]registry.Option)

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return func(_ *Registration, to *[]registry.Option) {
		*to = append(*to, registry.WithDescription(desc))
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return func(_ *Registration, to *[]registry.Option) {
		*to = append(*to, registry.WithUnit(unit))
	}
}

// WithKeys applies the provided dimension keys.
func WithKeys(keys ...core.Key) Option {
	return func(m *Registration, _ *[]registry.Option) {
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
