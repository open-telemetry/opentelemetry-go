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

package otlp_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func stubSpanSnapshot(count int) []*tracesdk.SpanSnapshot {
	spans := make([]*tracesdk.SpanSnapshot, 0, count)
	for i := 0; i < count; i++ {
		spans = append(spans, new(tracesdk.SpanSnapshot))
	}
	return spans
}

type stubCheckpointSet struct {
	limit int
}

var _ metricsdk.CheckpointSet = stubCheckpointSet{}

func (s stubCheckpointSet) ForEach(kindSelector metricsdk.ExportKindSelector, recordFunc func(metricsdk.Record) error) error {
	for i := 0; i < s.limit; i++ {
		if err := recordFunc(metricsdk.Record{}); err != nil {
			return err
		}
	}
	return nil
}

func (stubCheckpointSet) Lock()    {}
func (stubCheckpointSet) Unlock()  {}
func (stubCheckpointSet) RLock()   {}
func (stubCheckpointSet) RUnlock() {}

type stubProtocolDriver struct {
	started         int
	stopped         int
	tracesExported  int
	metricsExported int

	injectedStartError error
	injectedStopError  error

	rm []metricsdk.Record
	rs []tracesdk.SpanSnapshot
}

var _ otlp.ProtocolDriver = (*stubProtocolDriver)(nil)

func (m *stubProtocolDriver) Start(ctx context.Context) error {
	m.started++
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return m.injectedStartError
	}
}

func (m *stubProtocolDriver) Stop(ctx context.Context) error {
	m.stopped++
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return m.injectedStopError
	}
}

func (m *stubProtocolDriver) ExportMetrics(parent context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error {
	m.metricsExported++
	return cps.ForEach(selector, func(record metricsdk.Record) error {
		m.rm = append(m.rm, record)
		return nil
	})
}

func (m *stubProtocolDriver) ExportTraces(ctx context.Context, ss []*tracesdk.SpanSnapshot) error {
	m.tracesExported++
	for _, rs := range ss {
		if rs == nil {
			continue
		}
		m.rs = append(m.rs, *rs)
	}
	return nil
}

type stubTransformingProtocolDriver struct {
	rm []*metricpb.ResourceMetrics
	rs []*tracepb.ResourceSpans
}

var _ otlp.ProtocolDriver = (*stubTransformingProtocolDriver)(nil)

func (m *stubTransformingProtocolDriver) Start(ctx context.Context) error {
	return nil
}

func (m *stubTransformingProtocolDriver) Stop(ctx context.Context) error {
	return nil
}

func (m *stubTransformingProtocolDriver) ExportMetrics(parent context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error {
	rms, err := transform.CheckpointSet(parent, selector, cps, 1)
	if err != nil {
		return err
	}
	for _, rm := range rms {
		if rm == nil {
			continue
		}
		m.rm = append(m.rm, rm)
	}
	return nil
}

func (m *stubTransformingProtocolDriver) ExportTraces(ctx context.Context, ss []*tracesdk.SpanSnapshot) error {
	for _, rs := range transform.SpanData(ss) {
		if rs == nil {
			continue
		}
		m.rs = append(m.rs, rs)
	}
	return nil
}

func (m *stubTransformingProtocolDriver) Reset() {
	m.rm = nil
	m.rs = nil
}

func newExporter(t *testing.T, opts ...otlp.ExporterOption) (*otlp.Exporter, *stubTransformingProtocolDriver) {
	driver := &stubTransformingProtocolDriver{}
	exp, err := otlp.NewExporter(context.Background(), driver, opts...)
	require.NoError(t, err)
	return exp, driver
}

func TestExporterShutdownHonorsTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e := otlp.NewUnstartedExporter(&stubProtocolDriver{})
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

	e := otlp.NewUnstartedExporter(&stubProtocolDriver{})
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

	e := otlp.NewUnstartedExporter(&stubProtocolDriver{})
	if err := e.Start(ctx); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	if err := e.Shutdown(ctx); err != nil {
		t.Errorf("shutdown errored: expected nil, got %v", err)
	}
}

func TestExporterShutdownManyTimes(t *testing.T) {
	ctx := context.Background()
	e, err := otlp.NewExporter(ctx, &stubProtocolDriver{})
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
	_, _, _, err := otlp.InstallNewPipeline(ctx, &stubProtocolDriver{})
	assert.NoError(t, err)
	assert.IsType(t, &tracesdk.TracerProvider{}, otel.GetTracerProvider())
}

func TestNewExportPipeline(t *testing.T) {
	testCases := []struct {
		name             string
		expOpts          []otlp.ExporterOption
		testSpanSampling bool
	}{
		{
			name: "simple pipeline",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, tp, _, err := otlp.NewExportPipeline(
				context.Background(),
				&stubProtocolDriver{},
				tc.expOpts...,
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

func TestSplitDriver(t *testing.T) {
	driverTraces := &stubProtocolDriver{}
	driverMetrics := &stubProtocolDriver{}
	config := otlp.SplitConfig{
		ForMetrics: driverMetrics,
		ForTraces:  driverTraces,
	}
	driver := otlp.NewSplitDriver(config)
	ctx := context.Background()
	assert.NoError(t, driver.Start(ctx))
	assert.Equal(t, 1, driverTraces.started)
	assert.Equal(t, 1, driverMetrics.started)
	assert.Equal(t, 0, driverTraces.stopped)
	assert.Equal(t, 0, driverMetrics.stopped)
	assert.Equal(t, 0, driverTraces.tracesExported)
	assert.Equal(t, 0, driverTraces.metricsExported)
	assert.Equal(t, 0, driverMetrics.tracesExported)
	assert.Equal(t, 0, driverMetrics.metricsExported)

	recordCount := 5
	spanCount := 7
	assert.NoError(t, driver.ExportMetrics(ctx, stubCheckpointSet{recordCount}, metricsdk.StatelessExportKindSelector()))
	assert.NoError(t, driver.ExportTraces(ctx, stubSpanSnapshot(spanCount)))
	assert.Len(t, driverTraces.rm, 0)
	assert.Len(t, driverTraces.rs, spanCount)
	assert.Len(t, driverMetrics.rm, recordCount)
	assert.Len(t, driverMetrics.rs, 0)
	assert.Equal(t, 1, driverTraces.tracesExported)
	assert.Equal(t, 0, driverTraces.metricsExported)
	assert.Equal(t, 0, driverMetrics.tracesExported)
	assert.Equal(t, 1, driverMetrics.metricsExported)

	assert.NoError(t, driver.Stop(ctx))
	assert.Equal(t, 1, driverTraces.started)
	assert.Equal(t, 1, driverMetrics.started)
	assert.Equal(t, 1, driverTraces.stopped)
	assert.Equal(t, 1, driverMetrics.stopped)
	assert.Equal(t, 1, driverTraces.tracesExported)
	assert.Equal(t, 0, driverTraces.metricsExported)
	assert.Equal(t, 0, driverMetrics.tracesExported)
	assert.Equal(t, 1, driverMetrics.metricsExported)
}

func TestSplitDriverFail(t *testing.T) {
	ctx := context.Background()
	for i := 0; i < 16; i++ {
		var (
			errStartMetric error
			errStartTrace  error
			errStopMetric  error
			errStopTrace   error
		)
		if (i & 1) != 0 {
			errStartTrace = errors.New("trace start failed")
		}
		if (i & 2) != 0 {
			errStopTrace = errors.New("trace stop failed")
		}
		if (i & 4) != 0 {
			errStartMetric = errors.New("metric start failed")
		}
		if (i & 8) != 0 {
			errStopMetric = errors.New("metric stop failed")
		}
		shouldStartFail := errStartTrace != nil || errStartMetric != nil
		shouldStopFail := errStopTrace != nil || errStopMetric != nil

		driverTraces := &stubProtocolDriver{
			injectedStartError: errStartTrace,
			injectedStopError:  errStopTrace,
		}
		driverMetrics := &stubProtocolDriver{
			injectedStartError: errStartMetric,
			injectedStopError:  errStopMetric,
		}
		config := otlp.SplitConfig{
			ForMetrics: driverMetrics,
			ForTraces:  driverTraces,
		}
		driver := otlp.NewSplitDriver(config)
		errStart := driver.Start(ctx)
		if shouldStartFail {
			assert.Error(t, errStart)
		} else {
			assert.NoError(t, errStart)
		}
		errStop := driver.Stop(ctx)
		if shouldStopFail {
			assert.Error(t, errStop)
		} else {
			assert.NoError(t, errStop)
		}
	}
}
