// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/trace"
)

func TestNewIDs(t *testing.T) {
	gen := defaultIDGenerator()
	n := 1000

	for i := 0; i < n; i++ {
		traceID, spanID := gen.NewIDs(context.Background())
		assert.Truef(t, traceID.IsValid(), "trace id: %s", traceID.String())
		assert.Truef(t, spanID.IsValid(), "span id: %s", spanID.String())
	}
}

func TestNewSpanID(t *testing.T) {
	gen := defaultIDGenerator()
	testTraceID := [16]byte{123, 123}
	n := 1000

	for i := 0; i < n; i++ {
		spanID := gen.NewSpanID(context.Background(), testTraceID)
		assert.Truef(t, spanID.IsValid(), "span id: %s", spanID.String())
	}
}

func TestNewSpanIDWithInvalidTraceID(t *testing.T) {
	gen := defaultIDGenerator()
	spanID := gen.NewSpanID(context.Background(), trace.TraceID{})
	assert.Truef(t, spanID.IsValid(), "span id: %s", spanID.String())
}
