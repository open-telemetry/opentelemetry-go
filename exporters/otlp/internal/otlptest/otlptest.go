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

package otlptest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	exportmetric "go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

// RunEndToEndTest can be used by protocol driver tests to validate
// themselves.
func RunEndToEndTest(ctx context.Context, t *testing.T, exp *otlp.Exporter, mcTraces, mcMetrics Collector) {
	pOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(
			exp,
			// add following two options to ensure flush
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(10),
		),
	}
	tp1 := sdktrace.NewTracerProvider(append(pOpts,
		sdktrace.WithResource(resource.NewWithAttributes(
			attribute.String("rk1", "rv11)"),
			attribute.Int64("rk2", 5),
		)))...)

	tp2 := sdktrace.NewTracerProvider(append(pOpts,
		sdktrace.WithResource(resource.NewWithAttributes(
			attribute.String("rk1", "rv12)"),
			attribute.Float64("rk3", 6.5),
		)))...)

	tr1 := tp1.Tracer("test-tracer1")
	tr2 := tp2.Tracer("test-tracer2")
	// Now create few spans
	m := 4
	for i := 0; i < m; i++ {
		_, span := tr1.Start(ctx, "AlwaysSample")
		span.SetAttributes(attribute.Int64("i", int64(i)))
		span.End()

		_, span = tr2.Start(ctx, "AlwaysSample")
		span.SetAttributes(attribute.Int64("i", int64(i)))
		span.End()
	}

	selector := simple.NewWithInexpensiveDistribution()
	processor := processor.New(selector, exportmetric.StatelessExportKindSelector())
	cont := controller.New(processor, controller.WithExporter(exp))
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
	func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := tp1.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shut down a tracer provider 1: %v", err)
		}
		if err := tp2.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shut down a tracer provider 2: %v", err)
		}
	}()

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
	_ = mcTraces.Stop()
	_ = mcMetrics.Stop()

	// Now verify that we only got two resources
	rss := mcTraces.GetResourceSpans()
	if got, want := len(rss), 2; got != want {
		t.Fatalf("resource span count: got %d, want %d\n", got, want)
	}

	// Now verify spans and attributes for each resource span.
	for _, rs := range rss {
		if len(rs.InstrumentationLibrarySpans) == 0 {
			t.Fatalf("zero Instrumentation Library Spans")
		}
		if got, want := len(rs.InstrumentationLibrarySpans[0].Spans), m; got != want {
			t.Fatalf("span counts: got %d, want %d", got, want)
		}
		attrMap := map[int64]bool{}
		for _, s := range rs.InstrumentationLibrarySpans[0].Spans {
			if gotName, want := s.Name, "AlwaysSample"; gotName != want {
				t.Fatalf("span name: got %s, want %s", gotName, want)
			}
			attrMap[s.Attributes[0].Value.Value.(*commonpb.AnyValue_IntValue).IntValue] = true
		}
		if got, want := len(attrMap), m; got != want {
			t.Fatalf("span attribute unique values: got %d  want %d", got, want)
		}
		for i := 0; i < m; i++ {
			_, ok := attrMap[int64(i)]
			if !ok {
				t.Fatalf("span with attribute %d missing", i)
			}
		}
	}

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
		case metric.CounterInstrumentKind:
			switch data.nKind {
			case number.Int64Kind:
				if dp := m.GetIntSum().DataPoints; assert.Len(t, dp, 1) {
					assert.Equal(t, data.val, dp[0].Value, "invalid value for %q", m.Name)
				}
			case number.Float64Kind:
				if dp := m.GetDoubleSum().DataPoints; assert.Len(t, dp, 1) {
					assert.Equal(t, float64(data.val), dp[0].Value, "invalid value for %q", m.Name)
				}
			default:
				assert.Failf(t, "invalid number kind", data.nKind.String())
			}
		case metric.ValueObserverInstrumentKind:
			switch data.nKind {
			case number.Int64Kind:
				if dp := m.GetIntGauge().DataPoints; assert.Len(t, dp, 1) {
					assert.Equal(t, data.val, dp[0].Value, "invalid value for %q", m.Name)
				}
			case number.Float64Kind:
				if dp := m.GetDoubleGauge().DataPoints; assert.Len(t, dp, 1) {
					assert.Equal(t, float64(data.val), dp[0].Value, "invalid value for %q", m.Name)
				}
			default:
				assert.Failf(t, "invalid number kind", data.nKind.String())
			}
		case metric.ValueRecorderInstrumentKind:
			switch data.nKind {
			case number.Int64Kind:
				assert.NotNil(t, m.GetIntHistogram())
				if dp := m.GetIntHistogram().DataPoints; assert.Len(t, dp, 1) {
					count := dp[0].Count
					assert.Equal(t, uint64(1), count, "invalid count for %q", m.Name)
					assert.Equal(t, int64(data.val*int64(count)), dp[0].Sum, "invalid sum for %q (value %d)", m.Name, data.val)
				}
			case number.Float64Kind:
				assert.NotNil(t, m.GetDoubleHistogram())
				if dp := m.GetDoubleHistogram().DataPoints; assert.Len(t, dp, 1) {
					count := dp[0].Count
					assert.Equal(t, uint64(1), count, "invalid count for %q", m.Name)
					assert.Equal(t, float64(data.val*int64(count)), dp[0].Sum, "invalid sum for %q (value %d)", m.Name, data.val)
				}
			default:
				assert.Failf(t, "invalid number kind", data.nKind.String())
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
