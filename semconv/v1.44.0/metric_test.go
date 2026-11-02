// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package semconv_test

import (
	"context"
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/semconv/v1.44.0/azureconv"
	"go.opentelemetry.io/otel/semconv/v1.44.0/containerconv"
	"go.opentelemetry.io/otel/semconv/v1.44.0/goconv"
	"go.opentelemetry.io/otel/semconv/v1.44.0/httpconv"
	"go.opentelemetry.io/otel/semconv/v1.44.0/hwconv"
	"go.opentelemetry.io/otel/semconv/v1.44.0/k8sconv"
	"go.opentelemetry.io/otel/semconv/v1.44.0/openshiftconv"
	"go.opentelemetry.io/otel/semconv/v1.44.0/processconv"
	"go.opentelemetry.io/otel/semconv/v1.44.0/rpcconv"
	"go.opentelemetry.io/otel/semconv/v1.44.0/systemconv"
)

func TestFractionalMetricInstruments(t *testing.T) {
	t.Parallel()
	m := noop.NewMeterProvider().Meter("test")
	charge, err := azureconv.NewCosmosDBClientOperationRequestCharge(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Float64Histogram)(nil), charge.Inst())
	cpu, err := containerconv.NewCPUUsage(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Float64Gauge)(nil), cpu.Inst())
	cpuObservable, err := containerconv.NewCPUUsageObservable(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Float64ObservableGauge)(nil), cpuObservable.Inst())
	limit, err := k8sconv.NewContainerCPULimitCurrent(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Float64UpDownCounter)(nil), limit.Inst())
	quota, err := openshiftconv.NewClusterquotaCPURequestUsedObservable(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Float64ObservableUpDownCounter)(nil), quota.Inst())
	utilization, err := processconv.NewCPUUtilization(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Float64Gauge)(nil), utilization.Inst())
	paging, err := systemconv.NewPagingUtilizationObservable(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Float64ObservableGauge)(nil), paging.Inst())
	temperature, err := hwconv.NewTemperature(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Float64Gauge)(nil), temperature.Inst())
	energy, err := hwconv.NewEnergyObservable(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Float64ObservableCounter)(nil), energy.Inst())
	gogc, err := goconv.NewConfigGogc(m)
	require.NoError(t, err)
	assert.Implements(t, (*metric.Int64ObservableUpDownCounter)(nil), gogc.Inst())
}

type chargeHistogram struct {
	noop.Float64Histogram
	values []float64
}

func (*chargeHistogram) Enabled(context.Context) bool { return true }

func (h *chargeHistogram) Record(_ context.Context, value float64, _ ...metric.RecordOption) {
	h.values = append(h.values, value)
}

type chargeMeter struct {
	noop.Meter
	histogram chargeHistogram
}

func (m *chargeMeter) Float64Histogram(string, ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return &m.histogram, nil
}

func TestFractionalChargeRecording(t *testing.T) {
	t.Parallel()
	m := new(chargeMeter)
	charge, err := azureconv.NewCosmosDBClientOperationRequestCharge(m)
	require.NoError(t, err)
	charge.Record(t.Context(), 46.18, "ReadItem")
	assert.Equal(t, []float64{46.18}, m.histogram.values)
}

type histogramMeter struct {
	noop.Meter
	config metric.Float64HistogramConfig
}

func (m *histogramMeter) Float64Histogram(
	_ string,
	opts ...metric.Float64HistogramOption,
) (metric.Float64Histogram, error) {
	m.config = metric.NewFloat64HistogramConfig(opts...)
	return noop.Float64Histogram{}, nil
}

func TestHistogramBucketOptions(t *testing.T) {
	t.Parallel()
	constructors := map[string]func(metric.Meter, ...metric.Float64HistogramOption) error{
		"rpc.client": func(m metric.Meter, opts ...metric.Float64HistogramOption) error {
			_, err := rpcconv.NewClientCallDuration(m, opts...)
			return err
		},
		"rpc.server": func(m metric.Meter, opts ...metric.Float64HistogramOption) error {
			_, err := rpcconv.NewServerCallDuration(m, opts...)
			return err
		},
		"http.client": func(m metric.Meter, opts ...metric.Float64HistogramOption) error {
			_, err := httpconv.NewClientRequestDuration(m, opts...)
			return err
		},
		"http.server": func(m metric.Meter, opts ...metric.Float64HistogramOption) error {
			_, err := httpconv.NewServerRequestDuration(m, opts...)
			return err
		},
	}
	for name, newInstrument := range constructors {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := new(histogramMeter)
			require.NoError(t, newInstrument(m))
			defaults := slices.Clone(m.config.ExplicitBucketBoundaries())
			require.NotEmpty(t, defaults)
			description := m.config.Description()
			require.NoError(t, newInstrument(m, metric.WithDescription("override"), metric.WithUnit("ms")))
			assert.Equal(t, defaults, m.config.ExplicitBucketBoundaries())
			assert.Equal(t, description, m.config.Description())
			assert.Equal(t, "s", m.config.Unit())
			for _, bounds := range [][]float64{{0.01, 0.1, 1}, {}} {
				require.NoError(t, newInstrument(m,
					metric.WithExplicitBucketBoundaries(99),
					metric.WithExplicitBucketBoundaries(bounds...),
				))
				assert.Equal(t, bounds, m.config.ExplicitBucketBoundaries())
			}
			backing := []metric.Float64HistogramOption{
				metric.WithExplicitBucketBoundaries(0.1, 1),
				metric.WithDescription("spare"),
				metric.WithUnit("spare"),
			}
			want := slices.Clone(backing)
			require.NoError(t, newInstrument(m, backing[:1]...))
			assert.Equal(t, want, backing)
			require.NoError(t, newInstrument(m))
			assert.Equal(t, defaults, m.config.ExplicitBucketBoundaries())
		})
	}
}

func TestConstructorPreservesCallerOptions(t *testing.T) {
	t.Parallel()
	backing := []metric.Int64CounterOption{
		metric.WithDescription("one"),
		metric.WithDescription("two"),
		metric.WithDescription("three"),
	}
	want := slices.Clone(backing)
	m := noop.NewMeterProvider().Meter("test")
	_, err := processconv.NewDiskOperations(m, backing[:1]...)
	require.NoError(t, err)
	assert.Equal(t, want, backing)
}

func TestConstructorConcurrentSafe(t *testing.T) {
	t.Parallel()
	opts := make([]metric.Int64CounterOption, 1, 3)
	opts[0] = metric.WithDescription("test")
	m := noop.NewMeterProvider().Meter("test")
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 100 {
				_, err := processconv.NewDiskOperations(m, opts...)
				assert.NoError(t, err)
			}
		})
	}
	wg.Wait()
}
