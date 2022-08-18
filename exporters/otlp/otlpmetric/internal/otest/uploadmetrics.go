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

//go:build go1.18
// +build go1.18

package otest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otest"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/unit"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	cpb "go.opentelemetry.io/proto/otlp/common/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"
)

var (
	// Sat Jan 01 2000 00:00:00 GMT+0000.
	start = time.Date(2000, time.January, 01, 0, 0, 0, 0, time.FixedZone("GMT", 0))
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
	hdp           = []*mpb.HistogramDataPoint{{
		Attributes:        []*cpb.KeyValue{kvAlice},
		StartTimeUnixNano: uint64(start.UnixNano()),
		TimeUnixNano:      uint64(end.UnixNano()),
		Count:             30,
		Sum:               &sum,
		ExplicitBounds:    []float64{1, 5},
		BucketCounts:      []uint64{0, 30, 0},
		Min:               &min,
		Max:               &max,
	}}

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
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Gauge{Gauge: gaugeInt64},
		},
		{
			Name:        "float64-gauge",
			Description: "Gauge with float64 values",
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Gauge{Gauge: gaugeFloat64},
		},
		{
			Name:        "int64-sum",
			Description: "Sum with int64 values",
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Sum{Sum: sumInt64},
		},
		{
			Name:        "float64-sum",
			Description: "Sum with float64 values",
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Sum{Sum: sumFloat64},
		},
		{
			Name:        "histogram",
			Description: "Histogram",
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Histogram{Histogram: hist},
		},
	}

	scope = &cpb.InstrumentationScope{
		Name:    "test/code/path",
		Version: "v0.1.0",
	}
	scopeMetrics = []*mpb.ScopeMetrics{{
		Scope:     scope,
		Metrics:   metrics,
		SchemaUrl: semconv.SchemaURL,
	}}

	res = &rpb.Resource{
		Attributes: []*cpb.KeyValue{kvSrvName, kvSrvVer},
	}
	resourceMetrics = &mpb.ResourceMetrics{
		Resource:     res,
		ScopeMetrics: scopeMetrics,
		SchemaUrl:    semconv.SchemaURL,
	}
)

func testUploadMetrics(f ClientFactory) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			ctx := context.Background()
			client, coll := f()

			emptyRM := &mpb.ResourceMetrics{
				Resource:  res,
				SchemaUrl: semconv.SchemaURL,
			}
			require.NoError(t, client.UploadMetrics(ctx, emptyRM))

			emptySM := &mpb.ResourceMetrics{
				Resource: res,
				ScopeMetrics: []*mpb.ScopeMetrics{{
					Scope:     scope,
					SchemaUrl: semconv.SchemaURL,
				}},
				SchemaUrl: semconv.SchemaURL,
			}
			require.NoError(t, client.UploadMetrics(ctx, emptySM))

			require.NoError(t, client.Shutdown(ctx))
			got := coll.Collect().dump()
			assert.Contains(t, got, emptyRM)
			assert.Contains(t, got, emptySM)
		})

		t.Run("All", func(t *testing.T) {
			ctx := context.Background()
			client, coll := f()

			require.NoError(t, client.UploadMetrics(ctx, resourceMetrics))
			require.NoError(t, client.Shutdown(ctx))
			got := coll.Collect().dump()
			assert.Contains(t, got, resourceMetrics)
		})
	}
}
