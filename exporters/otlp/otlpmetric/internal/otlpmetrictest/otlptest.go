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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/metric/instrument"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

// RunEndToEndTest can be used by protocol driver tests to validate
// themselves.
func RunEndToEndTest(ctx context.Context, t *testing.T, exp *otlpmetric.Exporter, mcMetrics Collector) {
	selector := simple.NewWithHistogramDistribution()
	proc := processor.NewFactory(selector, aggregation.StatelessTemporalitySelector())
	cont := controller.New(proc, controller.WithExporter(exp))
	require.NoError(t, cont.Start(ctx))

	meter := cont.Meter("test-meter")
	attrs := []attribute.KeyValue{attribute.Bool("test", true)}

	type data struct {
		iKind sdkapi.InstrumentKind
		nKind number.Kind
		val   int64
	}
	instruments := map[string]data{
		"test-int64-counter":         {sdkapi.CounterInstrumentKind, number.Int64Kind, 1},
		"test-float64-counter":       {sdkapi.CounterInstrumentKind, number.Float64Kind, 1},
		"test-int64-histogram":       {sdkapi.HistogramInstrumentKind, number.Int64Kind, 2},
		"test-float64-histogram":     {sdkapi.HistogramInstrumentKind, number.Float64Kind, 2},
		"test-int64-gaugeobserver":   {sdkapi.GaugeObserverInstrumentKind, number.Int64Kind, 3},
		"test-float64-gaugeobserver": {sdkapi.GaugeObserverInstrumentKind, number.Float64Kind, 3},
	}
	for name, data := range instruments {
		data := data
		switch data.iKind {
		case sdkapi.CounterInstrumentKind:
			switch data.nKind {
			case number.Int64Kind:
				c, _ := meter.SyncInt64().Counter(name)
				c.Add(ctx, data.val, attrs...)
			case number.Float64Kind:
				c, _ := meter.SyncFloat64().Counter(name)
				c.Add(ctx, float64(data.val), attrs...)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case sdkapi.HistogramInstrumentKind:
			switch data.nKind {
			case number.Int64Kind:
				c, _ := meter.SyncInt64().Histogram(name)
				c.Record(ctx, data.val, attrs...)
			case number.Float64Kind:
				c, _ := meter.SyncFloat64().Histogram(name)
				c.Record(ctx, float64(data.val), attrs...)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case sdkapi.GaugeObserverInstrumentKind:
			switch data.nKind {
			case number.Int64Kind:
				g, _ := meter.AsyncInt64().Gauge(name)
				_ = meter.RegisterCallback([]instrument.Asynchronous{g}, func(ctx context.Context) {
					g.Observe(ctx, data.val, attrs...)
				})
			case number.Float64Kind:
				g, _ := meter.AsyncFloat64().Gauge(name)
				_ = meter.RegisterCallback([]instrument.Asynchronous{g}, func(ctx context.Context) {
					g.Observe(ctx, float64(data.val), attrs...)
				})
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		default:
			assert.Failf(t, "unsupported metrics testing kind", data.iKind.String())
		}
	}

	// Flush and close.
	require.NoError(t, cont.Stop(ctx))

	// Wait >2 cycles.
	<-time.After(40 * time.Millisecond)

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
		case sdkapi.CounterInstrumentKind, sdkapi.GaugeObserverInstrumentKind:
			var dp []*metricpb.NumberDataPoint
			switch data.iKind {
			case sdkapi.CounterInstrumentKind:
				require.NotNil(t, m.GetSum())
				dp = m.GetSum().GetDataPoints()
			case sdkapi.GaugeObserverInstrumentKind:
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
		case sdkapi.HistogramInstrumentKind:
			require.NotNil(t, m.GetHistogram())
			if dp := m.GetHistogram().DataPoints; assert.Len(t, dp, 1) {
				count := dp[0].Count
				assert.Equal(t, uint64(1), count, "invalid count for %q", m.Name)
				require.NotNil(t, dp[0].Sum)
				assert.Equal(t, float64(data.val*int64(count)), *dp[0].Sum, "invalid sum for %q (value %d)", m.Name, data.val)
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
