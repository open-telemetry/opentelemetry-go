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

	"go.opentelemetry.io/otel/api/core"
)

// Float64Gauge is a metric that stores the last float64 value.
type Float64Gauge struct {
	commonMetric
}

// Int64Gauge is a metric that stores the last int64 value.
type Int64Gauge struct {
	commonMetric
}

// BoundFloat64Gauge is a bound instrument for Float64Gauge.
//
// It inherits the Unbind function from commonBoundInstrument.
type BoundFloat64Gauge struct {
	commonBoundInstrument
}

// BoundInt64Gauge is a bound instrument for Int64Gauge.
//
// It inherits the Unbind function from commonBoundInstrument.
type BoundInt64Gauge struct {
	commonBoundInstrument
}

// Bind creates a bound instrument for this gauge. The labels should
// contain the keys and values for each key specified in the gauge
// with the WithKeys option.
//
// If the labels do not contain a value for the key specified in the
// gauge with the WithKeys option, then the missing value will be
// treated as unspecified.
func (g *Float64Gauge) Bind(labels LabelSet) (h BoundFloat64Gauge) {
	h.commonBoundInstrument = g.bind(labels)
	return
}

// Bind creates a bound instrument for this gauge. The labels should
// contain the keys and values for each key specified in the gauge
// with the WithKeys option.
//
// If the labels do not contain a value for the key specified in the
// gauge with the WithKeys option, then the missing value will be
// treated as unspecified.
func (g *Int64Gauge) Bind(labels LabelSet) (h BoundInt64Gauge) {
	h.commonBoundInstrument = g.bind(labels)
	return
}

// Measurement creates a Measurement object to use with batch
// recording.
func (g *Float64Gauge) Measurement(value float64) Measurement {
	return g.float64Measurement(value)
}

// Measurement creates a Measurement object to use with batch
// recording.
func (g *Int64Gauge) Measurement(value int64) Measurement {
	return g.int64Measurement(value)
}

// Set assigns the passed value to the value of the gauge. The labels
// should contain the keys and values for each key specified in the
// gauge with the WithKeys option.
//
// If the labels do not contain a value for the key specified in the
// gauge with the WithKeys option, then the missing value will be
// treated as unspecified.
func (g *Float64Gauge) Set(ctx context.Context, value float64, labels LabelSet) {
	g.directRecord(ctx, core.NewFloat64Number(value), labels)
}

// Set assigns the passed value to the value of the gauge. The labels
// should contain the keys and values for each key specified in the
// gauge with the WithKeys option.
//
// If the labels do not contain a value for the key specified in the
// gauge with the WithKeys option, then the missing value will be
// treated as unspecified.
func (g *Int64Gauge) Set(ctx context.Context, value int64, labels LabelSet) {
	g.directRecord(ctx, core.NewInt64Number(value), labels)
}

// Set assigns the passed value to the value of the gauge.
func (b *BoundFloat64Gauge) Set(ctx context.Context, value float64) {
	b.directRecord(ctx, core.NewFloat64Number(value))
}

// Set assigns the passed value to the value of the gauge.
func (b *BoundInt64Gauge) Set(ctx context.Context, value int64) {
	b.directRecord(ctx, core.NewInt64Number(value))
}
