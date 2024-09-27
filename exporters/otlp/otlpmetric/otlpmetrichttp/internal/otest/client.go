// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/otlpmetric/otest/client.go.tmpl

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/otest"

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	collpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	cpb "go.opentelemetry.io/proto/otlp/common/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"
)

var (
	// Sat Jan 01 2000 00:00:00 GMT+0000.
	start = time.Date(2000, time.January, 0o1, 0, 0, 0, 0, time.FixedZone("GMT", 0))
	end   = start.Add(30 * time.Second)

	kvAlice = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "alice"},
	}}
	kvBob = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "bob"},
	}}
	kvSrvName = &cpb.KeyValue{Key: "service.name", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "test server"},
	}}
	kvSrvVer = &cpb.KeyValue{Key: "service.version", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "v0.1.0"},
	}}

	min, max, sum = 2.0, 4.0, 90.0
	hdp           = []*mpb.HistogramDataPoint{
		{
			Attributes:        []*cpb.KeyValue{kvAlice},
			StartTimeUnixNano: uint64(start.UnixNano()),
			TimeUnixNano:      uint64(end.UnixNano()),
			Count:             30,
			Sum:               &sum,
			ExplicitBounds:    []float64{1, 5},
			BucketCounts:      []uint64{0, 30, 0},
			Min:               &min,
			Max:               &max,
		},
	}

	hist = &mpb.Histogram{
		AggregationTemporality: mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
		DataPoints:             hdp,
	}

	dPtsInt64 = []*mpb.NumberDataPoint{
		{
			Attributes:        []*cpb.KeyValue{kvAlice},
			StartTimeUnixNano: uint64(start.UnixNano()),
			TimeUnixNano:      uint64(end.UnixNano()),
			Value:             &mpb.NumberDataPoint_AsInt{AsInt: 1},
		},
		{
			Attributes:        []*cpb.KeyValue{kvBob},
			StartTimeUnixNano: uint64(start.UnixNano()),
			TimeUnixNano:      uint64(end.UnixNano()),
			Value:             &mpb.NumberDataPoint_AsInt{AsInt: 2},
		},
	}
	dPtsFloat64 = []*mpb.NumberDataPoint{
		{
			Attributes:        []*cpb.KeyValue{kvAlice},
			StartTimeUnixNano: uint64(start.UnixNano()),
			TimeUnixNano:      uint64(end.UnixNano()),
			Value:             &mpb.NumberDataPoint_AsDouble{AsDouble: 1.0},
		},
		{
			Attributes:        []*cpb.KeyValue{kvBob},
			StartTimeUnixNano: uint64(start.UnixNano()),
			TimeUnixNano:      uint64(end.UnixNano()),
			Value:             &mpb.NumberDataPoint_AsDouble{AsDouble: 2.0},
		},
	}

	sumInt64 = &mpb.Sum{
		AggregationTemporality: mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
		IsMonotonic:            true,
		DataPoints:             dPtsInt64,
	}
	sumFloat64 = &mpb.Sum{
		AggregationTemporality: mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
		IsMonotonic:            false,
		DataPoints:             dPtsFloat64,
	}

	gaugeInt64   = &mpb.Gauge{DataPoints: dPtsInt64}
	gaugeFloat64 = &mpb.Gauge{DataPoints: dPtsFloat64}

	metrics = []*mpb.Metric{
		{
			Name:        "int64-gauge",
			Description: "Gauge with int64 values",
			Unit:        "1",
			Data:        &mpb.Metric_Gauge{Gauge: gaugeInt64},
		},
		{
			Name:        "float64-gauge",
			Description: "Gauge with float64 values",
			Unit:        "1",
			Data:        &mpb.Metric_Gauge{Gauge: gaugeFloat64},
		},
		{
			Name:        "int64-sum",
			Description: "Sum with int64 values",
			Unit:        "1",
			Data:        &mpb.Metric_Sum{Sum: sumInt64},
		},
		{
			Name:        "float64-sum",
			Description: "Sum with float64 values",
			Unit:        "1",
			Data:        &mpb.Metric_Sum{Sum: sumFloat64},
		},
		{
			Name:        "histogram",
			Description: "Histogram",
			Unit:        "1",
			Data:        &mpb.Metric_Histogram{Histogram: hist},
		},
	}

	scope = &cpb.InstrumentationScope{
		Name:    "test/code/path",
		Version: "v0.1.0",
	}
	scopeMetrics = []*mpb.ScopeMetrics{
		{
			Scope:     scope,
			Metrics:   metrics,
			SchemaUrl: semconv.SchemaURL,
		},
	}

	res = &rpb.Resource{
		Attributes: []*cpb.KeyValue{kvSrvName, kvSrvVer},
	}
	resourceMetrics = &mpb.ResourceMetrics{
		Resource:     res,
		ScopeMetrics: scopeMetrics,
		SchemaUrl:    semconv.SchemaURL,
	}
)

type Client interface {
	UploadMetrics(context.Context, *mpb.ResourceMetrics) error
	ForceFlush(context.Context) error
	Shutdown(context.Context) error
}

// ClientFactory is a function that when called returns a
// Client implementation that is connected to also returned
// Collector implementation. The Client is ready to upload metric data to the
// Collector which is ready to store that data.
//
// If resultCh is not nil, the returned Collector needs to use the responses
// from that channel to send back to the client for every export request.
type ClientFactory func(resultCh <-chan ExportResult) (Client, Collector)

// RunClientTests runs a suite of Client integration tests. For example:
//
//	t.Run("Integration", RunClientTests(factory))
func RunClientTests(f ClientFactory) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("ClientHonorsContextErrors", func(t *testing.T) {
			t.Run("Shutdown", testCtxErrs(func() func(context.Context) error {
				c, _ := f(nil)
				return c.Shutdown
			}))

			t.Run("ForceFlush", testCtxErrs(func() func(context.Context) error {
				c, _ := f(nil)
				return c.ForceFlush
			}))

			t.Run("UploadMetrics", testCtxErrs(func() func(context.Context) error {
				c, _ := f(nil)
				return func(ctx context.Context) error {
					return c.UploadMetrics(ctx, nil)
				}
			}))
		})

		t.Run("ForceFlushFlushes", func(t *testing.T) {
			ctx := context.Background()
			client, collector := f(nil)
			require.NoError(t, client.UploadMetrics(ctx, resourceMetrics))

			require.NoError(t, client.ForceFlush(ctx))
			rm := collector.Collect().Dump()
			// Data correctness is not important, just it was received.
			require.NotEmpty(t, rm, "no data uploaded")

			require.NoError(t, client.Shutdown(ctx))
			rm = collector.Collect().Dump()
			assert.Empty(t, rm, "client did not flush all data")
		})

		t.Run("UploadMetrics", func(t *testing.T) {
			ctx := context.Background()
			client, coll := f(nil)

			require.NoError(t, client.UploadMetrics(ctx, resourceMetrics))
			require.NoError(t, client.Shutdown(ctx))
			got := coll.Collect().Dump()
			require.Len(t, got, 1, "upload of one ResourceMetrics")
			diff := cmp.Diff(got[0], resourceMetrics, cmp.Comparer(proto.Equal))
			if diff != "" {
				t.Fatalf("unexpected ResourceMetrics:\n%s", diff)
			}
		})

		t.Run("PartialSuccess", func(t *testing.T) {
			const n, msg = 2, "bad data"
			rCh := make(chan ExportResult, 3)
			rCh <- ExportResult{
				Response: &collpb.ExportMetricsServiceResponse{
					PartialSuccess: &collpb.ExportMetricsPartialSuccess{
						RejectedDataPoints: n,
						ErrorMessage:       msg,
					},
				},
			}
			rCh <- ExportResult{
				Response: &collpb.ExportMetricsServiceResponse{
					PartialSuccess: &collpb.ExportMetricsPartialSuccess{
						// Should not be logged.
						RejectedDataPoints: 0,
						ErrorMessage:       "",
					},
				},
			}
			rCh <- ExportResult{
				Response: &collpb.ExportMetricsServiceResponse{},
			}

			ctx := context.Background()
			client, _ := f(rCh)

			defer func(orig otel.ErrorHandler) {
				otel.SetErrorHandler(orig)
			}(otel.GetErrorHandler())

			errs := []error{}
			eh := otel.ErrorHandlerFunc(func(e error) { errs = append(errs, e) })
			otel.SetErrorHandler(eh)

			require.NoError(t, client.UploadMetrics(ctx, resourceMetrics))
			require.NoError(t, client.UploadMetrics(ctx, resourceMetrics))
			require.NoError(t, client.UploadMetrics(ctx, resourceMetrics))
			require.NoError(t, client.Shutdown(ctx))

			require.Len(t, errs, 1)
			want := fmt.Sprintf("%s (%d metric data points rejected)", msg, n)
			assert.ErrorContains(t, errs[0], want)
		})
	}
}

func testCtxErrs(factory func() func(context.Context) error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		t.Run("DeadlineExceeded", func(t *testing.T) {
			innerCtx, innerCancel := context.WithTimeout(ctx, time.Nanosecond)
			t.Cleanup(innerCancel)
			<-innerCtx.Done()

			f := factory()
			assert.ErrorIs(t, f(innerCtx), context.DeadlineExceeded)
		})

		t.Run("Canceled", func(t *testing.T) {
			innerCtx, innerCancel := context.WithCancel(ctx)
			innerCancel()

			f := factory()
			assert.ErrorIs(t, f(innerCtx), context.Canceled)
		})
	}
}
