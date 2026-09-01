// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package goconv

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

type countingMeter struct {
	noop.Meter
	histogram *countingHistogram
}

func (m countingMeter) Float64Histogram(
	string,
	...metric.Float64HistogramOption,
) (metric.Float64Histogram, error) {
	return m.histogram, nil
}

type countingHistogram struct {
	noop.Float64Histogram
	records int
}

func (h *countingHistogram) Record(context.Context, float64, ...metric.RecordOption) {
	h.records++
}

func TestScheduleDurationRecordWithoutAttributes(t *testing.T) {
	histogram := new(countingHistogram)
	instrument, err := NewScheduleDuration(countingMeter{histogram: histogram})
	if err != nil {
		t.Fatalf("NewScheduleDuration() error = %v", err)
	}

	instrument.Record(t.Context(), 0.1)
	if histogram.records != 1 {
		t.Fatalf("Record calls = %d, want 1", histogram.records)
	}
}
