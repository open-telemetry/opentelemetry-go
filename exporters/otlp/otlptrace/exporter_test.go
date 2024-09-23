// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptrace_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

type client struct {
	uploadErr error
}

var _ otlptrace.Client = &client{}

func (c *client) Start(ctx context.Context) error {
	return nil
}

func (c *client) Stop(ctx context.Context) error {
	return nil
}

func (c *client) UploadTraces(ctx context.Context, protoSpans []*tracepb.ResourceSpans) error {
	return c.uploadErr
}

func TestExporterClientError(t *testing.T) {
	ctx := context.Background()
	exp, err := otlptrace.New(ctx, &client{
		uploadErr: context.Canceled,
	})
	assert.NoError(t, err)

	spans := tracetest.SpanStubs{{Name: "Span 0"}}.Snapshots()
	err = exp.ExportSpans(ctx, spans)

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.True(t, strings.HasPrefix(err.Error(), "traces export: "), err)

	assert.NoError(t, exp.Shutdown(ctx))
}
