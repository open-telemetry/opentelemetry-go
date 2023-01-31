// Copyright The OpenTelemetry Authors
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

package otlptrace_test

import (
	"context"
	"errors"
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
	assert.True(t, errors.Is(err, context.Canceled))
	assert.True(t, strings.HasPrefix(err.Error(), "traces export: "), err)

	assert.NoError(t, exp.Shutdown(ctx))
}
