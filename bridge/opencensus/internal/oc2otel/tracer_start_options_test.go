// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oc2otel

import (
	"testing"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/trace"
)

func TestStartOptionsSpanKind(t *testing.T) {
	conv := map[int]trace.SpanKind{
		octrace.SpanKindClient:      trace.SpanKindClient,
		octrace.SpanKindServer:      trace.SpanKindServer,
		octrace.SpanKindUnspecified: trace.SpanKindUnspecified,
	}

	for oc, otel := range conv {
		ocOpts := []octrace.StartOption{octrace.WithSpanKind(oc)}
		otelOpts, err := StartOptions(ocOpts)
		if err != nil {
			t.Errorf("StartOptions errored: %v", err)
			continue
		}
		c := trace.NewSpanStartConfig(otelOpts...)
		if c.SpanKind() != otel {
			t.Errorf("conversion of SpanKind start option: got %v, want %v", c.SpanKind(), otel)
		}
	}
}

func TestStartOptionsSamplerErrors(t *testing.T) {
	ocOpts := []octrace.StartOption{octrace.WithSampler(octrace.AlwaysSample())}
	_, err := StartOptions(ocOpts)
	if err == nil {
		t.Error("StartOptions should error Sampler option")
	}
}
