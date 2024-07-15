// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestValidateSpanKind(t *testing.T) {
	tests := []struct {
		in   SpanKind
		want SpanKind
	}{
		{
			SpanKindUnspecified,
			SpanKindInternal,
		},
		{
			SpanKindInternal,
			SpanKindInternal,
		},
		{
			SpanKindServer,
			SpanKindServer,
		},
		{
			SpanKindClient,
			SpanKindClient,
		},
		{
			SpanKindProducer,
			SpanKindProducer,
		},
		{
			SpanKindConsumer,
			SpanKindConsumer,
		},
	}
	for _, test := range tests {
		if got := ValidateSpanKind(test.in); got != test.want {
			t.Errorf("ValidateSpanKind(%#v) = %#v, want %#v", test.in, got, test.want)
		}
	}
}

func TestSpanKindString(t *testing.T) {
	tests := []struct {
		in   SpanKind
		want string
	}{
		{
			SpanKindUnspecified,
			"unspecified",
		},
		{
			SpanKindInternal,
			"internal",
		},
		{
			SpanKindServer,
			"server",
		},
		{
			SpanKindClient,
			"client",
		},
		{
			SpanKindProducer,
			"producer",
		},
		{
			SpanKindConsumer,
			"consumer",
		},
	}
	for _, test := range tests {
		if got := test.in.String(); got != test.want {
			t.Errorf("%#v.String() = %#v, want %#v", test.in, got, test.want)
		}
	}
}

func TestLinkFromContext(t *testing.T) {
	k1v1 := attribute.String("key1", "value1")
	spanCtx := SpanContext{traceID: TraceID([16]byte{1}), remote: true}

	receiverCtx := ContextWithRemoteSpanContext(context.Background(), spanCtx)
	link := LinkFromContext(receiverCtx, k1v1)

	if !assertSpanContextEqual(link.SpanContext, spanCtx) {
		t.Fatalf("LinkFromContext: Unexpected context created: %s", cmp.Diff(link.SpanContext, spanCtx))
	}
	assert.Equal(t, link.Attributes[0], k1v1)
}
