// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptrace_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type client struct {
	uploadErr error
}

var _ otlptrace.Client = &client{}

func (*client) Start(context.Context) error {
	return nil
}

func (*client) Stop(context.Context) error {
	return nil
}

func (c *client) UploadTraces(context.Context, []*tracepb.ResourceSpans) error {
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
	assert.True(t, strings.HasPrefix(err.Error(), "traces export: "), "%+v", err)

	assert.NoError(t, exp.Shutdown(ctx))
}
