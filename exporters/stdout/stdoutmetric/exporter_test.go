// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutmetric_test // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func testEncoderOption() stdoutmetric.Option {
	// Discard export output for testing.
	enc := json.NewEncoder(io.Discard)
	return stdoutmetric.WithEncoder(enc)
}

func testCtxErrHonored(factory func(*testing.T) func(context.Context) error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		ctx := context.Background()

		t.Run("DeadlineExceeded", func(t *testing.T) {
			innerCtx, innerCancel := context.WithTimeout(ctx, time.Nanosecond)
			t.Cleanup(innerCancel)
			<-innerCtx.Done()

			f := factory(t)
			assert.ErrorIs(t, f(innerCtx), context.DeadlineExceeded)
		})

		t.Run("Canceled", func(t *testing.T) {
			innerCtx, innerCancel := context.WithCancel(ctx)
			innerCancel()

			f := factory(t)
			assert.ErrorIs(t, f(innerCtx), context.Canceled)
		})

		t.Run("NoError", func(t *testing.T) {
			f := factory(t)
			assert.NoError(t, f(ctx))
		})
	}
}

func testCtxErrIgnored(factory func(*testing.T) func(context.Context) error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		ctx := context.Background()

		t.Run("Canceled Ignored", func(t *testing.T) {
			innerCtx, innerCancel := context.WithCancel(ctx)
			innerCancel()

			f := factory(t)
			assert.NoError(t, f(innerCtx))
		})

		t.Run("NoError", func(t *testing.T) {
			f := factory(t)
			assert.NoError(t, f(ctx))
		})
	}
}

func TestExporterExportHonorsContextErrors(t *testing.T) {
	t.Run("Export", testCtxErrHonored(func(t *testing.T) func(context.Context) error {
		exp, err := stdoutmetric.New(testEncoderOption())
		require.NoError(t, err)
		return func(ctx context.Context) error {
			data := new(metricdata.ResourceMetrics)
			return exp.Export(ctx, data)
		}
	}))
}

func TestExporterForceFlushIgnoresContextErrors(t *testing.T) {
	t.Run("ForceFlush", testCtxErrIgnored(func(t *testing.T) func(context.Context) error {
		exp, err := stdoutmetric.New(testEncoderOption())
		require.NoError(t, err)
		return exp.ForceFlush
	}))
}

func TestExporterShutdownIgnoresContextErrors(t *testing.T) {
	t.Run("Shutdown", testCtxErrIgnored(func(t *testing.T) func(context.Context) error {
		exp, err := stdoutmetric.New(testEncoderOption())
		require.NoError(t, err)
		return exp.Shutdown
	}))
}

func TestShutdownExporterReturnsShutdownErrorOnExport(t *testing.T) {
	var (
		data     = new(metricdata.ResourceMetrics)
		ctx      = context.Background()
		exp, err = stdoutmetric.New(testEncoderOption())
	)
	require.NoError(t, err)
	require.NoError(t, exp.Shutdown(ctx))
	assert.EqualError(t, exp.Export(ctx, data), "exporter shutdown")
}

func deltaSelector(metric.InstrumentKind) metricdata.Temporality {
	return metricdata.DeltaTemporality
}

func TestExportWithOptions(t *testing.T) {
	var (
		data = new(metricdata.ResourceMetrics)
		ctx  = context.Background()
	)

	for _, tt := range []struct {
		name string
		opts []stdoutmetric.Option

		expectedData string
	}{
		{
			name:         "with no options",
			expectedData: "{\"Resource\":null,\"ScopeMetrics\":null}\n",
		},
		{
			name: "with pretty print",
			opts: []stdoutmetric.Option{
				stdoutmetric.WithPrettyPrint(),
			},
			expectedData: "{\n\t\"Resource\": null,\n\t\"ScopeMetrics\": null\n}\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			opts := append(tt.opts, stdoutmetric.WithWriter(&b))

			exp, err := stdoutmetric.New(opts...)
			require.NoError(t, err)
			require.NoError(t, exp.Export(ctx, data))

			assert.Equal(t, tt.expectedData, b.String())
		})
	}
}

func TestTemporalitySelector(t *testing.T) {
	exp, err := stdoutmetric.New(
		testEncoderOption(),
		stdoutmetric.WithTemporalitySelector(deltaSelector),
	)
	require.NoError(t, err)

	var unknownKind metric.InstrumentKind
	assert.Equal(t, metricdata.DeltaTemporality, exp.Temporality(unknownKind))
}

func dropSelector(metric.InstrumentKind) metric.Aggregation {
	return metric.AggregationDrop{}
}

func TestAggregationSelector(t *testing.T) {
	exp, err := stdoutmetric.New(
		testEncoderOption(),
		stdoutmetric.WithAggregationSelector(dropSelector),
	)
	require.NoError(t, err)

	var unknownKind metric.InstrumentKind
	assert.Equal(t, metric.AggregationDrop{}, exp.Aggregation(unknownKind))
}
