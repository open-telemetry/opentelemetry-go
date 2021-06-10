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

	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

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

type stubTransformingProtocolDriver struct {
	rm []*metricpb.ResourceMetrics
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

func (m *stubTransformingProtocolDriver) Reset() {
	m.rm = nil
}

func newExporter(t *testing.T, opts ...otlp.ExporterOption) (*otlp.Exporter, *stubTransformingProtocolDriver) {
	driver := &stubTransformingProtocolDriver{}
	exp, err := otlp.New(context.Background(), driver, opts...)
	require.NoError(t, err)
	return exp, driver
}

func TestExporterShutdownHonorsTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e := otlp.NewUnstarted(&stubProtocolDriver{})
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

	e := otlp.NewUnstarted(&stubProtocolDriver{})
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

	e := otlp.NewUnstarted(&stubProtocolDriver{})
	if err := e.Start(ctx); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	if err := e.Shutdown(ctx); err != nil {
		t.Errorf("shutdown errored: expected nil, got %v", err)
	}
}

func TestExporterShutdownManyTimes(t *testing.T) {
	ctx := context.Background()
	e, err := otlp.New(ctx, &stubProtocolDriver{})
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

func TestSplitDriver(t *testing.T) {
	recordCount := 5
	assertExport := func(t testing.TB, ctx context.Context, driver otlp.ProtocolDriver) {
		t.Helper()
		assert.NoError(t, driver.ExportMetrics(ctx, stubCheckpointSet{recordCount}, metricsdk.StatelessExportKindSelector()))
	}

	t.Run("with just metric driver", func(t *testing.T) {
		driverMetrics := &stubProtocolDriver{}

		driver := otlp.NewSplitDriver(otlp.WithMetricDriver(driverMetrics))
		ctx := context.Background()
		assert.NoError(t, driver.Start(ctx))

		assert.Equal(t, 1, driverMetrics.started)
		assert.Equal(t, 0, driverMetrics.stopped)
		assert.Equal(t, 0, driverMetrics.tracesExported)
		assert.Equal(t, 0, driverMetrics.metricsExported)

		assertExport(t, ctx, driver)
		assert.Len(t, driverMetrics.rm, recordCount)
		assert.Equal(t, 0, driverMetrics.tracesExported)
		assert.Equal(t, 1, driverMetrics.metricsExported)

		assert.NoError(t, driver.Stop(ctx))
		assert.Equal(t, 1, driverMetrics.started)
		assert.Equal(t, 1, driverMetrics.stopped)
		assert.Equal(t, 0, driverMetrics.tracesExported)
		assert.Equal(t, 1, driverMetrics.metricsExported)
	})

	t.Run("with no drivers configured", func(t *testing.T) {
		driver := otlp.NewSplitDriver()
		ctx := context.Background()
		assert.NoError(t, driver.Start(ctx))

		assert.NoError(t, driver.ExportMetrics(ctx, stubCheckpointSet{recordCount}, metricsdk.StatelessExportKindSelector()))
		assert.NoError(t, driver.Stop(ctx))
	})

}

func TestSplitDriverFail(t *testing.T) {
	ctx := context.Background()
	for i := 0; i < 16; i++ {
		var (
			errStartMetric error
			errStopMetric  error
		)
		if (i & 4) != 0 {
			errStartMetric = errors.New("metric start failed")
		}
		if (i & 8) != 0 {
			errStopMetric = errors.New("metric stop failed")
		}
		shouldStartFail := errStartMetric != nil
		shouldStopFail := errStopMetric != nil

		driverMetrics := &stubProtocolDriver{
			injectedStartError: errStartMetric,
			injectedStopError:  errStopMetric,
		}
		driver := otlp.NewSplitDriver(otlp.WithMetricDriver(driverMetrics))
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
