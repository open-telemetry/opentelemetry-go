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

package stdoutmetric_test // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
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
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

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

func TestExporterHonorsContextErrors(t *testing.T) {
	t.Run("Shutdown", testCtxErrHonored(func(t *testing.T) func(context.Context) error {
		exp, err := stdoutmetric.New(testEncoderOption())
		require.NoError(t, err)
		return exp.Shutdown
	}))

	t.Run("ForceFlush", testCtxErrHonored(func(t *testing.T) func(context.Context) error {
		exp, err := stdoutmetric.New(testEncoderOption())
		require.NoError(t, err)
		return exp.ForceFlush
	}))

	t.Run("Export", testCtxErrHonored(func(t *testing.T) func(context.Context) error {
		exp, err := stdoutmetric.New(testEncoderOption())
		require.NoError(t, err)
		return func(ctx context.Context) error {
			data := new(metricdata.ResourceMetrics)
			return exp.Export(ctx, data)
		}
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

func TestTemporalitySelector(t *testing.T) {
	exp, err := stdoutmetric.New(
		testEncoderOption(),
		stdoutmetric.WithTemporalitySelector(deltaSelector),
	)
	require.NoError(t, err)

	var unknownKind metric.InstrumentKind
	assert.Equal(t, metricdata.DeltaTemporality, exp.Temporality(unknownKind))
}

func dropSelector(metric.InstrumentKind) aggregation.Aggregation {
	return aggregation.Drop{}
}

func TestAggregationSelector(t *testing.T) {
	exp, err := stdoutmetric.New(
		testEncoderOption(),
		stdoutmetric.WithAggregationSelector(dropSelector),
	)
	require.NoError(t, err)

	var unknownKind metric.InstrumentKind
	assert.Equal(t, aggregation.Drop{}, exp.Aggregation(unknownKind))
}
