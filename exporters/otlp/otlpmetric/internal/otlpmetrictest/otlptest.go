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

package otlpmetrictest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otlpmetrictest"

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/metric/instrument"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

type client struct {
	lock            sync.Mutex
	resourceMetrics *metricpb.ResourceMetrics
	uploadCount     uint64
	startCount      uint64
	stopCount       uint64
}

var _ otlpmetric.Client = &client{}

func (c *client) Start(ctx context.Context) error {
	atomic.AddUint64(&c.startCount, 1)
	return nil
}
func (c *client) Stop(ctx context.Context) error {
	atomic.AddUint64(&c.stopCount, 1)
	return nil
}
func (c *client) UploadMetrics(ctx context.Context, protoMetrics *metricpb.ResourceMetrics) error {
	atomic.AddUint64(&c.uploadCount, 1)
	c.lock.Lock()
	defer c.lock.Unlock()

	c.resourceMetrics = protoMetrics

	return nil
}

// RunEndToEndTest can be used by protocol driver tests to validate
// themselves.
func RunEndToEndTest(ctx context.Context, t *testing.T, exp *otlpmetric.Exporter, mcMetrics Collector) {
	rdr := reader.NewManualReader(exp)
	mp := sdkmetric.New(
		sdkmetric.WithReader(rdr),
	)

	meter := mp.Meter("test-meter")
	labels := []attribute.KeyValue{attribute.Bool("test", true)}

	type data struct {
		iKind sdkinstrument.Kind
		nKind number.Kind
		val   int64
	}
	instruments := map[string]data{
		"test-int64-counter":         {sdkinstrument.CounterKind, number.Int64Kind, 1},
		"test-float64-counter":       {sdkinstrument.CounterKind, number.Float64Kind, 1},
		"test-int64-gaugeobserver":   {sdkinstrument.GaugeObserverKind, number.Int64Kind, 3},
		"test-float64-gaugeobserver": {sdkinstrument.GaugeObserverKind, number.Float64Kind, 3},
	}
	for name, data := range instruments {
		data := data
		switch data.iKind {
		case sdkinstrument.CounterKind:
			switch data.nKind {
			case number.Int64Kind:
				c, _ := meter.SyncInt64().Counter(name)
				c.Add(ctx, data.val, labels...)
			case number.Float64Kind:
				c, _ := meter.SyncFloat64().Counter(name)
				c.Add(ctx, float64(data.val), labels...)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case sdkinstrument.HistogramKind:
			switch data.nKind {
			case number.Int64Kind:
				c, _ := meter.SyncInt64().Histogram(name)
				c.Record(ctx, data.val, labels...)
			case number.Float64Kind:
				c, _ := meter.SyncFloat64().Histogram(name)
				c.Record(ctx, float64(data.val), labels...)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case sdkinstrument.GaugeObserverKind:
			switch data.nKind {
			case number.Int64Kind:
				g, _ := meter.AsyncInt64().Gauge(name)
				_ = meter.RegisterCallback([]instrument.Asynchronous{g}, func(ctx context.Context) {
					g.Observe(ctx, data.val, labels...)
				})
			case number.Float64Kind:
				g, _ := meter.AsyncFloat64().Gauge(name)
				_ = meter.RegisterCallback([]instrument.Asynchronous{g}, func(ctx context.Context) {
					g.Observe(ctx, float64(data.val), labels...)
				})
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		default:
			assert.Failf(t, "unsupported metrics testing kind", data.iKind.String())
		}
	}

	// Collect
	err := rdr.Collect(ctx, nil)
	assert.NoError(t, err)

	// Now shutdown the exporter
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := exp.Shutdown(ctx); err != nil {
		t.Fatalf("failed to stop the exporter: %v", err)
	}

	// Shutdown the collector too so that we can begin
	// verification checks of expected data back.
	_ = mcMetrics.Stop()

	metrics := mcMetrics.GetMetrics()
	assert.Len(t, metrics, len(instruments), "not enough metrics exported")
	seen := make(map[string]struct{}, len(instruments))
	for _, m := range metrics {
		data, ok := instruments[m.Name]
		if !ok {
			assert.Failf(t, "unknown metrics", m.Name)
			continue
		}
		seen[m.Name] = struct{}{}

		switch data.iKind {
		case sdkinstrument.CounterKind, sdkinstrument.GaugeObserverKind:
			var dp []*metricpb.NumberDataPoint
			switch data.iKind {
			case sdkinstrument.CounterKind:
				require.NotNil(t, m.GetSum())
				dp = m.GetSum().GetDataPoints()
			case sdkinstrument.GaugeObserverKind:
				require.NotNil(t, m.GetGauge())
				dp = m.GetGauge().GetDataPoints()
			}
			if assert.Len(t, dp, 1) {
				switch data.nKind {
				case number.Int64Kind:
					v := &metricpb.NumberDataPoint_AsInt{AsInt: data.val}
					assert.Equal(t, v, dp[0].Value, "invalid value for %q", m.Name)
				case number.Float64Kind:
					v := &metricpb.NumberDataPoint_AsDouble{AsDouble: float64(data.val)}
					assert.Equal(t, v, dp[0].Value, "invalid value for %q", m.Name)
				}
			}
		case sdkinstrument.HistogramKind:
			require.NotNil(t, m.GetSummary())
			if dp := m.GetSummary().DataPoints; assert.Len(t, dp, 1) {
				count := dp[0].Count
				assert.Equal(t, uint64(1), count, "invalid count for %q", m.Name)
				assert.Equal(t, float64(data.val*int64(count)), dp[0].Sum, "invalid sum for %q (value %d)", m.Name, data.val)
			}
		default:
			assert.Failf(t, "invalid metrics kind", data.iKind.String())
		}
	}

	for i := range instruments {
		if _, ok := seen[i]; !ok {
			assert.Fail(t, fmt.Sprintf("no metric(s) exported for %q", i))
		}
	}
}
