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

package otlpgrpc_test

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"

	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/internal/otlptest"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
)

func TestNew_endToEnd(t *testing.T) {
	tests := []struct {
		name           string
		additionalOpts []otlpgrpc.Option
	}{
		{
			name: "StandardExporter",
		},
		{
			name: "WithCompressor",
			additionalOpts: []otlpgrpc.Option{
				otlpgrpc.WithCompressor(gzip.Name),
			},
		},
		{
			name: "WithServiceConfig",
			additionalOpts: []otlpgrpc.Option{
				otlpgrpc.WithServiceConfig("{}"),
			},
		},
		{
			name: "WithDialOptions",
			additionalOpts: []otlpgrpc.Option{
				otlpgrpc.WithDialOption(grpc.WithBlock()),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newExporterEndToEndTest(t, test.additionalOpts)
		})
	}
}

func newGRPCExporter(t *testing.T, ctx context.Context, endpoint string, additionalOpts ...otlpgrpc.Option) *otlp.Exporter {
	opts := []otlpgrpc.Option{
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(endpoint),
		otlpgrpc.WithReconnectionPeriod(50 * time.Millisecond),
	}

	opts = append(opts, additionalOpts...)
	driver := otlpgrpc.NewDriver(opts...)
	exp, err := otlp.New(ctx, driver)
	if err != nil {
		t.Fatalf("failed to create a new collector exporter: %v", err)
	}
	return exp
}

func newExporterEndToEndTest(t *testing.T, additionalOpts []otlpgrpc.Option) {
	mc := runMockCollectorAtEndpoint(t, "localhost:56561")

	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint, additionalOpts...)
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := exp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	otlptest.RunEndToEndTest(ctx, t, exp, mc, mc)
}

func TestNew_invokeStartThenStopManyTimes(t *testing.T) {
	mc := runMockCollector(t)
	defer func() {
		_ = mc.stop()
	}()

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	defer func() {
		if err := exp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	// Invoke Start numerous times, should return errAlreadyStarted
	for i := 0; i < 10; i++ {
		if err := exp.Start(ctx); err == nil || !strings.Contains(err.Error(), "already started") {
			t.Fatalf("#%d unexpected Start error: %v", i, err)
		}
	}

	if err := exp.Shutdown(ctx); err != nil {
		t.Fatalf("failed to Shutdown the exporter: %v", err)
	}
	// Invoke Shutdown numerous times
	for i := 0; i < 10; i++ {
		if err := exp.Shutdown(ctx); err != nil {
			t.Fatalf(`#%d got error (%v) expected none`, i, err)
		}
	}
}

// This test takes a long time to run: to skip it, run tests using: -short
func TestNew_collectorOnBadConnection(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping this long running test")
	}

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to grab an available port: %v", err)
	}
	// Firstly close the "collector's" channel: optimistically this endpoint won't get reused ASAP
	// However, our goal of closing it is to simulate an unavailable connection
	_ = ln.Close()

	_, collectorPortStr, _ := net.SplitHostPort(ln.Addr().String())

	endpoint := fmt.Sprintf("localhost:%s", collectorPortStr)
	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, endpoint)
	_ = exp.Shutdown(ctx)
}

func TestNew_withEndpoint(t *testing.T) {
	mc := runMockCollector(t)
	defer func() {
		_ = mc.stop()
	}()

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	_ = exp.Shutdown(ctx)
}

func TestNew_WithTimeout(t *testing.T) {
	tts := []struct {
		name    string
		fn      func(exp *otlp.Exporter) error
		timeout time.Duration
		metrics int
		spans   int
		code    codes.Code
		delay   bool
	}{
		{
			name: "Timeout Metrics",
			fn: func(exp *otlp.Exporter) error {
				return exp.Export(context.Background(), otlptest.OneRecordCheckpointSet{})
			},
			timeout: time.Millisecond * 100,
			code:    codes.DeadlineExceeded,
			delay:   true,
		},
		{
			name: "No Timeout Metrics",
			fn: func(exp *otlp.Exporter) error {
				return exp.Export(context.Background(), otlptest.OneRecordCheckpointSet{})
			},
			timeout: time.Minute,
			metrics: 1,
			code:    codes.OK,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {

			mc := runMockCollector(t)
			if tt.delay {
				mc.metricSvc.delay = time.Second * 10
			}
			defer func() {
				_ = mc.stop()
			}()

			ctx := context.Background()
			exp := newGRPCExporter(t, ctx, mc.endpoint, otlpgrpc.WithTimeout(tt.timeout), otlpgrpc.WithRetry(otlp.RetrySettings{Enabled: false}))
			defer func() {
				_ = exp.Shutdown(ctx)
			}()

			err := tt.fn(exp)

			if tt.code == codes.OK {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			s := status.Convert(err)
			require.Equal(t, tt.code, s.Code())

			require.Len(t, mc.getMetrics(), tt.metrics)
		})
	}
}

func TestDisconnected(t *testing.T) {
	ctx := context.Background()
	// The endpoint is whatever, we want to be disconnected. But we
	// setting a blocking connection, so dialing to the invalid
	// endpoint actually fails.
	exp := newGRPCExporter(t, ctx, "invalid",
		otlpgrpc.WithReconnectionPeriod(time.Hour),
		otlpgrpc.WithDialOption(
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
		),
	)
	defer func() {
		assert.NoError(t, exp.Shutdown(ctx))
	}()

	assert.Error(t, exp.Export(ctx, otlptest.OneRecordCheckpointSet{}))
}

func TestEmptyData(t *testing.T) {
	mc := runMockCollectorAtEndpoint(t, "localhost:56561")

	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	defer func() {
		assert.NoError(t, exp.Shutdown(ctx))
	}()

	assert.NoError(t, exp.Export(ctx, otlptest.EmptyCheckpointSet{}))
}

func TestFailedMetricTransform(t *testing.T) {
	mc := runMockCollectorAtEndpoint(t, "localhost:56561")

	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	defer func() {
		assert.NoError(t, exp.Shutdown(ctx))
	}()

	assert.Error(t, exp.Export(ctx, otlptest.FailCheckpointSet{}))
}

func TestMultiConnectionDriver(t *testing.T) {
	mcTraces := runMockCollector(t)
	mcMetrics := runMockCollector(t)

	defer func() {
		_ = mcTraces.stop()
		_ = mcMetrics.stop()
	}()

	<-time.After(5 * time.Millisecond)

	commonOpts := []otlpgrpc.Option{
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithReconnectionPeriod(50 * time.Millisecond),
		otlpgrpc.WithDialOption(grpc.WithBlock()),
	}
	optsMetrics := append([]otlpgrpc.Option{
		otlpgrpc.WithEndpoint(mcMetrics.endpoint),
	}, commonOpts...)

	metricsDriver := otlpgrpc.NewDriver(optsMetrics...)
	driver := otlp.NewSplitDriver(otlp.WithMetricDriver(metricsDriver))
	ctx := context.Background()
	exp, err := otlp.New(ctx, driver)
	if err != nil {
		t.Fatalf("failed to create a new collector exporter: %v", err)
	}
	defer func() {
		assert.NoError(t, exp.Shutdown(ctx))
	}()
	otlptest.RunEndToEndTest(ctx, t, exp, mcTraces, mcMetrics)
}
