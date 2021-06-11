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

package otlpmetrictest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	exportmetric "go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

// RunEndToEndTest can be used by protocol driver tests to validate
// themselves.
func RunEndToEndTest(ctx context.Context, t *testing.T, exp *otlpmetric.Exporter, mcMetrics Collector) {
	selector := simple.NewWithInexpensiveDistribution()
	proc := processor.New(selector, exportmetric.StatelessExportKindSelector())
	cont := controller.New(proc, controller.WithExporter(exp))
	require.NoError(t, cont.Start(ctx))

	meter := cont.MeterProvider().Meter("test-meter")
	labels := []attribute.KeyValue{attribute.Bool("test", true)}

	type data struct {
		iKind metric.InstrumentKind
		nKind number.Kind
		val   int64
	}
	instruments := map[string]data{
		"test-int64-counter":         {metric.CounterInstrumentKind, number.Int64Kind, 1},
		"test-float64-counter":       {metric.CounterInstrumentKind, number.Float64Kind, 1},
		"test-int64-valuerecorder":   {metric.ValueRecorderInstrumentKind, number.Int64Kind, 2},
		"test-float64-valuerecorder": {metric.ValueRecorderInstrumentKind, number.Float64Kind, 2},
		"test-int64-valueobserver":   {metric.ValueObserverInstrumentKind, number.Int64Kind, 3},
		"test-float64-valueobserver": {metric.ValueObserverInstrumentKind, number.Float64Kind, 3},
	}
	for name, data := range instruments {
		data := data
		switch data.iKind {
		case metric.CounterInstrumentKind:
			switch data.nKind {
			case number.Int64Kind:
				metric.Must(meter).NewInt64Counter(name).Add(ctx, data.val, labels...)
			case number.Float64Kind:
				metric.Must(meter).NewFloat64Counter(name).Add(ctx, float64(data.val), labels...)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case metric.ValueRecorderInstrumentKind:
			switch data.nKind {
			case number.Int64Kind:
				metric.Must(meter).NewInt64ValueRecorder(name).Record(ctx, data.val, labels...)
			case number.Float64Kind:
				metric.Must(meter).NewFloat64ValueRecorder(name).Record(ctx, float64(data.val), labels...)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case metric.ValueObserverInstrumentKind:
			switch data.nKind {
			case number.Int64Kind:
				metric.Must(meter).NewInt64ValueObserver(name,
					func(_ context.Context, result metric.Int64ObserverResult) {
						result.Observe(data.val, labels...)
					},
				)
			case number.Float64Kind:
				callback := func(v float64) metric.Float64ObserverFunc {
					return metric.Float64ObserverFunc(func(_ context.Context, result metric.Float64ObserverResult) { result.Observe(v, labels...) })
				}(float64(data.val))
				metric.Must(meter).NewFloat64ValueObserver(name, callback)
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
		case metric.CounterInstrumentKind, metric.ValueObserverInstrumentKind:
			var dp []*metricpb.NumberDataPoint
			switch data.iKind {
			case metric.CounterInstrumentKind:
				require.NotNil(t, m.GetSum())
				dp = m.GetSum().GetDataPoints()
			case metric.ValueObserverInstrumentKind:
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
		case metric.ValueRecorderInstrumentKind:
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
