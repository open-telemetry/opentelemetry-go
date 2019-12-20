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

// Float64Measure is a metric that records float64 values.
type Float64Measure struct {
	commonMetric
}

// Int64Measure is a metric that records int64 values.
type Int64Measure struct {
	commonMetric
}

// Float64MeasureBoundInstrument is a bound instrument for Float64Measure.
//
// It inherits the Release function from commonBoundInstrument.
type Float64MeasureBoundInstrument struct {
	commonBoundInstrument
}

// Int64MeasureBoundInstrument is a bound instrument for Int64Measure.
//
// It inherits the Release function from commonBoundInstrument.
type Int64MeasureBoundInstrument struct {
	commonBoundInstrument
}

// AcquireBoundInstrument creates a boundInstrument for this measure. The labels should
// contain the keys and values for each key specified in the measure
// with the WithKeys option.
//
// If the labels do not contain a value for the key specified in the
// measure with the WithKeys option, then the missing value will be
// treated as unspecified.
func (c *Float64Measure) AcquireBoundInstrument(labels LabelSet) (h Float64MeasureBoundInstrument) {
	h.commonBoundInstrument = c.acquireCommonBoundInstrument(labels)
	return
}

// AcquireBoundInstrument creates a boundInstrument for this measure. The labels should
// contain the keys and values for each key specified in the measure
// with the WithKeys option.
//
// If the labels do not contain a value for the key specified in the
// measure with the WithKeys option, then the missing value will be
// treated as unspecified.
func (c *Int64Measure) AcquireBoundInstrument(labels LabelSet) (h Int64MeasureBoundInstrument) {
	h.commonBoundInstrument = c.acquireCommonBoundInstrument(labels)
	return
}

// Measurement creates a Measurement object to use with batch
// recording.
func (c *Float64Measure) Measurement(value float64) Measurement {
	return c.float64Measurement(value)
}

// Measurement creates a Measurement object to use with batch
// recording.
func (c *Int64Measure) Measurement(value int64) Measurement {
	return c.int64Measurement(value)
}

// Record adds a new value to the list of measure's records. The
// labels should contain the keys and values for each key specified in
// the measure with the WithKeys option.
//
// If the labels do not contain a value for the key specified in the
// measure with the WithKeys option, then the missing value will be
// treated as unspecified.
func (c *Float64Measure) Record(ctx context.Context, value float64, labels LabelSet) {
	c.recordOne(ctx, core.NewFloat64Number(value), labels)
}

// Record adds a new value to the list of measure's records. The
// labels should contain the keys and values for each key specified in
// the measure with the WithKeys option.
//
// If the labels do not contain a value for the key specified in the
// measure with the WithKeys option, then the missing value will be
// treated as unspecified.
func (c *Int64Measure) Record(ctx context.Context, value int64, labels LabelSet) {
	c.recordOne(ctx, core.NewInt64Number(value), labels)
}

// Record adds a new value to the list of measure's records.
func (h *Float64MeasureBoundInstrument) Record(ctx context.Context, value float64) {
	h.recordOne(ctx, core.NewFloat64Number(value))
}

// Record adds a new value to the list of measure's records.
func (h *Int64MeasureBoundInstrument) Record(ctx context.Context, value int64) {
	h.recordOne(ctx, core.NewInt64Number(value))
}
