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
	*opts = append(*opts, metric.WithAttributes(attribute.String("key", "value")))
	backing := *opts

	PutAddOptions(opts)

	if got := len(*opts); got != 0 {
		t.Fatalf("PutAddOptions left slice length %d, want 0", got)
	}
	if backing[0] != nil {
		t.Fatal("PutAddOptions did not clear the slice backing array")
	}
}

func TestPutRecordOptionsClearsAndResetsSlice(t *testing.T) {
	opts := RecordOptions()
	*opts = append(*opts, metric.WithAttributes(attribute.String("key", "value")))
	backing := *opts

	PutRecordOptions(opts)

	if got := len(*opts); got != 0 {
		t.Fatalf("PutRecordOptions left slice length %d, want 0", got)
	}
	if backing[0] != nil {
		t.Fatal("PutRecordOptions did not clear the slice backing array")
	}
}
