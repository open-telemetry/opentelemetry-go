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
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

type stubClient struct {
	rs []*tracepb.ResourceSpans
}

var _ otlptrace.Client = (*stubClient)(nil)

func (m *stubClient) Start(_ context.Context) error {
	return nil
}

func (m *stubClient) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (m *stubClient) UploadTraces(ctx context.Context, protoSpans []*tracepb.ResourceSpans) error {
	m.rs = append(m.rs, protoSpans...)
	return nil
}

func (m *stubClient) Reset() {
	m.rs = nil
}

func TestExporterShutdownHonorsTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e := otlptrace.NewUnstartedExporter(&stubClient{})
	if err := e.Start(ctx); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	innerCtx, innerCancel := context.WithTimeout(ctx, time.Microsecond)
	<-time.After(time.Second)
	if err := e.Shutdown(innerCtx); err == nil {
		t.Error("expected context DeadlineExceeded error, got nil")
	} else if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context DeadlineExceeded error, got %v", err)
	}
	innerCancel()
}

func TestExporterShutdownHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e := otlptrace.NewUnstartedExporter(&stubClient{})
	if err := e.Start(ctx); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	var innerCancel context.CancelFunc
	ctx, innerCancel = context.WithCancel(ctx)
	innerCancel()
	if err := e.Shutdown(ctx); err == nil {
		t.Error("expected context canceled error, got nil")
	} else if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context canceled error, got %v", err)
	}
}

func TestExporterShutdownNoError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e := otlptrace.NewUnstartedExporter(&stubClient{})
	if err := e.Start(ctx); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	if err := e.Shutdown(ctx); err != nil {
		t.Errorf("shutdown errored: expected nil, got %v", err)
	}
}

func TestExporterShutdownManyTimes(t *testing.T) {
	ctx := context.Background()
	e, err := otlptrace.NewExporter(ctx, &stubClient{})
	if err != nil {
		t.Fatalf("failed to start an exporter: %v", err)
	}
	ch := make(chan struct{})
	wg := sync.WaitGroup{}
	const num int = 20
	wg.Add(num)
	errs := make([]error, num)
	for i := 0; i < num; i++ {
		go func(idx int) {
			defer wg.Done()
			<-ch
			errs[idx] = e.Shutdown(ctx)
		}(i)
	}
	close(ch)
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			t.Fatalf("failed to shutdown exporter: %v", err)
		}
	}
}

func TestInstallNewPipeline(t *testing.T) {
	ctx := context.Background()
	_, _, err := otlptrace.InstallNewPipeline(ctx, &stubClient{})
	assert.NoError(t, err)
	assert.IsType(t, &tracesdk.TracerProvider{}, otel.GetTracerProvider())
}

func TestNewExportPipeline(t *testing.T) {
	testCases := []struct {
		name             string
		testSpanSampling bool
	}{
		{
			name: "simple pipeline",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, tp, err := otlptrace.NewExportPipeline(
				context.Background(),
				&stubClient{},
			)

			assert.NoError(t, err)
			assert.NotEqual(t, tp, otel.GetTracerProvider())

			_, span := tp.Tracer("otlp test").Start(context.Background(), tc.name)
			spanCtx := span.SpanContext()
			assert.Equal(t, true, spanCtx.IsSampled())
			span.End()
		})
	}
}
