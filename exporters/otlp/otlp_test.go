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

package otlp

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	metricpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/metrics/v1"
	tracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/trace/v1"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"

	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/export/trace"
)

type stubProtocolDriver struct {
	rm []metricpb.ResourceMetrics
	rs []tracepb.ResourceSpans
}

var _ ProtocolDriver = (*stubProtocolDriver)(nil)

func (m *stubProtocolDriver) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (m *stubProtocolDriver) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (m *stubProtocolDriver) ExportMetrics(parent context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error {
	rms, err := transform.CheckpointSet(parent, selector, cps, 1)
	if err != nil {
		return err
	}
	for _, rm := range rms {
		if rm == nil {
			continue
		}
		m.rm = append(m.rm, *rm)
	}
	return nil
}

func (m *stubProtocolDriver) ExportTraces(ctx context.Context, ss []*tracesdk.SpanSnapshot) error {
	for _, rs := range transform.SpanData(ss) {
		if rs == nil {
			continue
		}
		m.rs = append(m.rs, *rs)
	}
	return nil
}

func (m *stubProtocolDriver) Reset() {
	m.rm = nil
	m.rs = nil
}

func newExporter(t *testing.T, opts ...ExporterOption) (*Exporter, *stubProtocolDriver) {
	driver := &stubProtocolDriver{}
	exp, err := NewExporter(context.Background(), driver, opts...)
	require.NoError(t, err)
	return exp, driver
}

func TestExporterShutdownHonorsTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e := NewUnstartedExporter(&stubProtocolDriver{})
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

	e := NewUnstartedExporter(&stubProtocolDriver{})
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

	e := NewUnstartedExporter(&stubProtocolDriver{})
	if err := e.Start(ctx); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	if err := e.Shutdown(ctx); err != nil {
		t.Errorf("shutdown errored: expected nil, got %v", err)
	}
}

func TestExporterShutdownManyTimes(t *testing.T) {
	ctx := context.Background()
	e, err := NewExporter(ctx, &stubProtocolDriver{})
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
