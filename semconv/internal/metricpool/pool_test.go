// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metricpool

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func TestPutAddOptionsClearsAndResetsSlice(t *testing.T) {
	opts := AddOptions()
	*opts = append(
		*opts,
		metric.WithAttributes(attribute.String("key1", "value1")),
		metric.WithAttributes(attribute.String("key2", "value2")),
	)
	backing := *opts

	PutAddOptions(opts)

	if got := len(*opts); got != 0 {
		t.Fatalf("PutAddOptions left slice length %d, want 0", got)
	}
	for i, v := range backing {
		if v != nil {
			t.Fatalf("PutAddOptions did not clear index %d", i)
		}
	}
}

func TestPutRecordOptionsClearsAndResetsSlice(t *testing.T) {
	opts := RecordOptions()
	*opts = append(
		*opts,
		metric.WithAttributes(attribute.String("key1", "value1")),
		metric.WithAttributes(attribute.String("key2", "value2")),
	)
	backing := *opts

	PutRecordOptions(opts)

	if got := len(*opts); got != 0 {
		t.Fatalf("PutRecordOptions left slice length %d, want 0", got)
	}
	for i, v := range backing {
		if v != nil {
			t.Fatalf("PutRecordOptions did not clear index %d", i)
		}
	}
}
