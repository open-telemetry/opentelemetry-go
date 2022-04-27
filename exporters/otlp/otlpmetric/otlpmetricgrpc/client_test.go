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

package otlpmetricgrpc_test

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otlpmetrictest"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	oneRecord = otlpmetrictest.OneRecordReader()

	testResource = resource.Empty()
)

func TestNewExporterEndToEnd(t *testing.T) {
	tests := []struct {
		name           string
		additionalOpts []otlpmetricgrpc.Option
	}{
		{
			name: "StandardExporter",
		},
		{
			name: "WithCompressor",
			additionalOpts: []otlpmetricgrpc.Option{
				otlpmetricgrpc.WithCompressor(gzip.Name),
			},
		},
		{
			name: "WithServiceConfig",
			additionalOpts: []otlpmetricgrpc.Option{
				otlpmetricgrpc.WithServiceConfig("{}"),
			},
		},
		{
			name: "WithDialOptions",
			additionalOpts: []otlpmetricgrpc.Option{
				otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newExporterEndToEndTest(t, test.additionalOpts)
		})
	}
}

func newGRPCExporter(t *testing.T, ctx context.Context, endpoint string, additionalOpts ...otlpmetricgrpc.Option) *otlpmetric.Exporter {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithReconnectionPeriod(50 * time.Millisecond),
	}

	opts = append(opts, additionalOpts...)
	client := otlpmetricgrpc.NewClient(opts...)
	exp, err := otlpmetric.New(ctx, client)
	if err != nil {
		t.Fatalf("failed to create a new collector exporter: %v", err)
	}
	return exp
}

func newExporterEndToEndTest(t *testing.T, additionalOpts []otlpmetricgrpc.Option) {
	mc := runMockCollector(t)

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

	otlpmetrictest.RunEndToEndTest(ctx, t, exp, mc)
}

func TestExporterShutdown(t *testing.T) {
	mc := runMockCollector(t)
	defer func() {
		_ = mc.Stop()
	}()

	<-time.After(5 * time.Millisecond)

	otlpmetrictest.RunExporterShutdownTest(t, func() otlpmetric.Client {
		return otlpmetricgrpc.NewClient(
			otlpmetricgrpc.WithInsecure(),
			otlpmetricgrpc.WithEndpoint(mc.endpoint),
			otlpmetricgrpc.WithReconnectionPeriod(50*time.Millisecond),
		)
	})
}

func TestNewExporterInvokeStartThenStopManyTimes(t *testing.T) {
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

// This test takes a long time to run: to skip it, run tests using: -short.
func TestNewExporterCollectorOnBadConnection(t *testing.T) {
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

func TestNewExporterWithEndpoint(t *testing.T) {
	mc := runMockCollector(t)
	defer func() {
		_ = mc.stop()
	}()

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	_ = exp.Shutdown(ctx)
}

func TestNewExporterWithHeaders(t *testing.T) {
	mc := runMockCollector(t)
	defer func() {
		_ = mc.stop()
	}()

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlpmetricgrpc.WithHeaders(map[string]string{"header1": "value1"}))
	require.NoError(t, exp.Export(ctx, testResource, oneRecord))

	defer func() {
		_ = exp.Shutdown(ctx)
	}()

	headers := mc.getHeaders()
	require.Len(t, headers.Get("header1"), 1)
	assert.Equal(t, "value1", headers.Get("header1")[0])
}

func TestNewExporterWithTimeout(t *testing.T) {
	tts := []struct {
		name    string
		fn      func(exp *otlpmetric.Exporter) error
		timeout time.Duration
		metrics int
		spans   int
		code    codes.Code
		delay   bool
	}{
		{
			name: "Timeout Metrics",
			fn: func(exp *otlpmetric.Exporter) error {
				return exp.Export(context.Background(), testResource, oneRecord)
			},
			timeout: time.Millisecond * 100,
			code:    codes.DeadlineExceeded,
			delay:   true,
		},

		{
			name: "No Timeout Metrics",
			fn: func(exp *otlpmetric.Exporter) error {
				return exp.Export(context.Background(), testResource, oneRecord)
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
			exp := newGRPCExporter(t, ctx, mc.endpoint, otlpmetricgrpc.WithTimeout(tt.timeout), otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{Enabled: false}))
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

func TestStartErrorInvalidAddress(t *testing.T) {
	client := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		// Validate the connection in Start (which should return the error).
		otlpmetricgrpc.WithDialOption(
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
		),
		otlpmetricgrpc.WithEndpoint("invalid"),
		otlpmetricgrpc.WithReconnectionPeriod(time.Hour),
	)
	err := client.Start(context.Background())
	assert.EqualError(t, err, `connection error: desc = "transport: error while dialing: dial tcp: address invalid: missing port in address"`)
}

func TestEmptyData(t *testing.T) {
	mc := runMockCollector(t)

	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	defer func() {
		assert.NoError(t, exp.Shutdown(ctx))
	}()

	assert.NoError(t, exp.Export(ctx, testResource, otlpmetrictest.EmptyReader()))
}

func TestFailedMetricTransform(t *testing.T) {
	mc := runMockCollector(t)

	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	defer func() {
		assert.NoError(t, exp.Shutdown(ctx))
	}()

	assert.Error(t, exp.Export(ctx, testResource, otlpmetrictest.FailReader{}))
}
