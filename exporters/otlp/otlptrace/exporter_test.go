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
	"testing"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

type noopClient struct {
}

var _ otlptrace.Client = (*noopClient)(nil)

func (m *noopClient) Start(_ context.Context) error {
	return nil
}

func (m *noopClient) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (m *noopClient) UploadTraces(_ context.Context, _ []*tracepb.ResourceSpans) error {
	return nil
}

func (m *noopClient) Reset() {
}

func TestInstallNewPipeline(t *testing.T) {
	ctx := context.Background()
	_, _, err := otlptrace.InstallNewPipeline(ctx, &noopClient{})
	assert.NoError(t, err)
	assert.IsType(t, &tracesdk.TracerProvider{}, otel.GetTracerProvider())
}

func TestNewExportPipeline(t *testing.T) {
	_, tp, err := otlptrace.NewExportPipeline(context.Background(), &noopClient{})
	assert.NoError(t, err)
	assert.NotEqual(t, tp, otel.GetTracerProvider())
}
