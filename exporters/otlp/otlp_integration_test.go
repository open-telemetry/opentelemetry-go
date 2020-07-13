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

package otlp_test

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/common/v1"
	metricpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/metrics/v1"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	metricapi "go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/otlp"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	exporttrace "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestNewExporter_endToEnd(t *testing.T) {
	tests := []struct {
		name           string
		additionalOpts []otlp.ExporterOption
	}{
		{
			name: "StandardExporter",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newExporterEndToEndTest(t, test.additionalOpts)
		})
	}
}

func newExporterEndToEndTest(t *testing.T, additionalOpts []otlp.ExporterOption) {
	mc := runMockColAtAddr(t, "localhost:56561")

	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	opts := []otlp.ExporterOption{
		otlp.WithInsecure(),
		otlp.WithAddress(mc.address),
		otlp.WithReconnectionPeriod(50 * time.Millisecond),
	}

	opts = append(opts, additionalOpts...)
	exp, err := otlp.NewExporter(opts...)
	if err != nil {
		t.Fatalf("failed to create a new collector exporter: %v", err)
	}
	defer func() {
		_ = exp.Stop()
	}()

	pOpts := []sdktrace.ProviderOption{
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(exp, // add following two options to ensure flush
			sdktrace.WithBatchTimeout(15),
			sdktrace.WithMaxExportBatchSize(10),
		),
	}
	tp1, err := sdktrace.NewProvider(append(pOpts,
		sdktrace.WithResource(resource.New(
			kv.String("rk1", "rv11)"),
			kv.Int64("rk2", 5),
		)))...)
	assert.NoError(t, err)

	tp2, err := sdktrace.NewProvider(append(pOpts,
		sdktrace.WithResource(resource.New(
			kv.String("rk1", "rv12)"),
			kv.Float32("rk3", 6.5),
		)))...)
	assert.NoError(t, err)

	tr1 := tp1.Tracer("test-tracer1")
	tr2 := tp2.Tracer("test-tracer2")
	// Now create few spans
	m := 4
	for i := 0; i < m; i++ {
		_, span := tr1.Start(context.Background(), "AlwaysSample")
		span.SetAttributes(kv.Int64("i", int64(i)))
		span.End()

		_, span = tr2.Start(context.Background(), "AlwaysSample")
		span.SetAttributes(kv.Int64("i", int64(i)))
		span.End()
	}

	selector := simple.NewWithExactDistribution()
	processor := processor.New(selector, metricsdk.PassThroughExporter)
	pusher := push.New(processor, exp)
	pusher.Start()

	ctx := context.Background()
	meter := pusher.Provider().Meter("test-meter")
	labels := []kv.KeyValue{kv.Bool("test", true)}

	type data struct {
		iKind metric.Kind
		nKind metricapi.NumberKind
		val   int64
	}
	instruments := map[string]data{
		"test-int64-counter":         {metric.CounterKind, metricapi.Int64NumberKind, 1},
		"test-float64-counter":       {metric.CounterKind, metricapi.Float64NumberKind, 1},
		"test-int64-valuerecorder":   {metric.ValueRecorderKind, metricapi.Int64NumberKind, 2},
		"test-float64-valuerecorder": {metric.ValueRecorderKind, metricapi.Float64NumberKind, 2},
		"test-int64-valueobserver":   {metric.ValueObserverKind, metricapi.Int64NumberKind, 3},
		"test-float64-valueobserver": {metric.ValueObserverKind, metricapi.Float64NumberKind, 3},
	}
	for name, data := range instruments {
		switch data.iKind {
		case metric.CounterKind:
			switch data.nKind {
			case metricapi.Int64NumberKind:
				metricapi.Must(meter).NewInt64Counter(name).Add(ctx, data.val, labels...)
			case metricapi.Float64NumberKind:
				metricapi.Must(meter).NewFloat64Counter(name).Add(ctx, float64(data.val), labels...)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case metric.ValueRecorderKind:
			switch data.nKind {
			case metricapi.Int64NumberKind:
				metricapi.Must(meter).NewInt64ValueRecorder(name).Record(ctx, data.val, labels...)
			case metricapi.Float64NumberKind:
				metricapi.Must(meter).NewFloat64ValueRecorder(name).Record(ctx, float64(data.val), labels...)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case metric.ValueObserverKind:
			switch data.nKind {
			case metricapi.Int64NumberKind:
				callback := func(v int64) metricapi.Int64ObserverCallback {
					return metricapi.Int64ObserverCallback(func(_ context.Context, result metricapi.Int64ObserverResult) { result.Observe(v, labels...) })
				}(data.val)
				metricapi.Must(meter).NewInt64ValueObserver(name, callback)
			case metricapi.Float64NumberKind:
				callback := func(v float64) metricapi.Float64ObserverCallback {
					return metricapi.Float64ObserverCallback(func(_ context.Context, result metricapi.Float64ObserverResult) { result.Observe(v, labels...) })
				}(float64(data.val))
				metricapi.Must(meter).NewFloat64ValueObserver(name, callback)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		default:
			assert.Failf(t, "unsupported metrics testing kind", data.iKind.String())
		}
	}

	// Flush and close.
	pusher.Stop()

	// Wait >2 cycles.
	<-time.After(40 * time.Millisecond)

	// Now shutdown the exporter
	if err := exp.Stop(); err != nil {
		t.Fatalf("failed to stop the exporter: %v", err)
	}

	// Shutdown the collector too so that we can begin
	// verification checks of expected data back.
	_ = mc.stop()

	// Now verify that we only got two resources
	rss := mc.getResourceSpans()
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

	metrics := mc.getMetrics()
	assert.Len(t, metrics, len(instruments), "not enough metrics exported")
	seen := make(map[string]struct{}, len(instruments))
	for _, m := range metrics {
		desc := m.GetMetricDescriptor()
		data, ok := instruments[desc.Name]
		if !ok {
			assert.Failf(t, "unknown metrics", desc.Name)
			continue
		}
		seen[desc.Name] = struct{}{}

		switch data.iKind {
		case metric.CounterKind:
			switch data.nKind {
			case metricapi.Int64NumberKind:
				assert.Equal(t, metricpb.MetricDescriptor_INT64.String(), desc.GetType().String())
				if dp := m.GetInt64DataPoints(); assert.Len(t, dp, 1) {
					assert.Equal(t, data.val, dp[0].Value, "invalid value for %q", desc.Name)
				}
			case metricapi.Float64NumberKind:
				assert.Equal(t, metricpb.MetricDescriptor_DOUBLE.String(), desc.GetType().String())
				if dp := m.GetDoubleDataPoints(); assert.Len(t, dp, 1) {
					assert.Equal(t, float64(data.val), dp[0].Value, "invalid value for %q", desc.Name)
				}
			default:
				assert.Failf(t, "invalid number kind", data.nKind.String())
			}
		case metric.ValueRecorderKind, metric.ValueObserverKind:
			assert.Equal(t, metricpb.MetricDescriptor_SUMMARY.String(), desc.GetType().String())
			m.GetSummaryDataPoints()
			if dp := m.GetSummaryDataPoints(); assert.Len(t, dp, 1) {
				count := dp[0].Count
				assert.Equal(t, uint64(1), count, "invalid count for %q", desc.Name)
				assert.Equal(t, float64(data.val*int64(count)), dp[0].Sum, "invalid sum for %q (value %d)", desc.Name, data.val)
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

func TestNewExporter_invokeStartThenStopManyTimes(t *testing.T) {
	mc := runMockCol(t)
	defer func() {
		_ = mc.stop()
	}()

	exp, err := otlp.NewExporter(otlp.WithInsecure(),
		otlp.WithReconnectionPeriod(50*time.Millisecond),
		otlp.WithAddress(mc.address))
	if err != nil {
		t.Fatalf("error creating exporter: %v", err)
	}
	defer func() {
		_ = exp.Stop()
	}()

	// Invoke Start numerous times, should return errAlreadyStarted
	for i := 0; i < 10; i++ {
		if err := exp.Start(); err == nil || !strings.Contains(err.Error(), "already started") {
			t.Fatalf("#%d unexpected Start error: %v", i, err)
		}
	}

	_ = exp.Stop()
	// Invoke Stop numerous times
	for i := 0; i < 10; i++ {
		if err := exp.Stop(); err != nil {
			t.Fatalf(`#%d got error (%v) expected none`, i, err)
		}
	}
}

func TestNewExporter_collectorConnectionDiesThenReconnects(t *testing.T) {
	mc := runMockCol(t)

	reconnectionPeriod := 20 * time.Millisecond
	exp, err := otlp.NewExporter(otlp.WithInsecure(),
		otlp.WithAddress(mc.address),
		otlp.WithReconnectionPeriod(reconnectionPeriod))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer func() {
		_ = exp.Stop()
	}()

	// We'll now stop the collector right away to simulate a connection
	// dying in the midst of communication or even not existing before.
	_ = mc.stop()

	// In the test below, we'll stop the collector many times,
	// while exporting traces and test to ensure that we can
	// reconnect.
	for j := 0; j < 3; j++ {

		exp.ExportSpans(context.Background(), []*exporttrace.SpanData{{Name: "in the midst"}})

		// Now resurrect the collector by making a new one but reusing the
		// old address, and the collector should reconnect automatically.
		nmc := runMockColAtAddr(t, mc.address)

		// Give the exporter sometime to reconnect
		<-time.After(reconnectionPeriod * 4)

		n := 10
		for i := 0; i < n; i++ {
			exp.ExportSpans(context.Background(), []*exporttrace.SpanData{{Name: "Resurrected"}})
		}

		nmaSpans := nmc.getSpans()
		// Expecting 10 spanData that were sampled, given that
		if g, w := len(nmaSpans), n; g != w {
			t.Fatalf("Round #%d: Connected collector: spans: got %d want %d", j, g, w)
		}

		dSpans := mc.getSpans()
		// Expecting 0 spans to have been received by the original but now dead collector
		if g, w := len(dSpans), 0; g != w {
			t.Fatalf("Round #%d: Disconnected collector: spans: got %d want %d", j, g, w)
		}
		_ = nmc.stop()
	}
}

// This test takes a long time to run: to skip it, run tests using: -short
func TestNewExporter_collectorOnBadConnection(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping this long running test")
	}

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to grab an available port: %v", err)
	}
	// Firstly close the "collector's" channel: optimistically this address won't get reused ASAP
	// However, our goal of closing it is to simulate an unavailable connection
	_ = ln.Close()

	_, collectorPortStr, _ := net.SplitHostPort(ln.Addr().String())

	address := fmt.Sprintf("localhost:%s", collectorPortStr)
	exp, err := otlp.NewExporter(otlp.WithInsecure(),
		otlp.WithReconnectionPeriod(50*time.Millisecond),
		otlp.WithAddress(address))
	if err != nil {
		t.Fatalf("Despite an indefinite background reconnection, got error: %v", err)
	}
	_ = exp.Stop()
}

func TestNewExporter_withAddress(t *testing.T) {
	mc := runMockCol(t)
	defer func() {
		_ = mc.stop()
	}()

	exp := otlp.NewUnstartedExporter(
		otlp.WithInsecure(),
		otlp.WithReconnectionPeriod(50*time.Millisecond),
		otlp.WithAddress(mc.address))

	defer func() {
		_ = exp.Stop()
	}()

	if err := exp.Start(); err != nil {
		t.Fatalf("Unexpected Start error: %v", err)
	}
}

func TestNewExporter_withHeaders(t *testing.T) {
	mc := runMockCol(t)
	defer func() {
		_ = mc.stop()
	}()

	exp, _ := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithReconnectionPeriod(50*time.Millisecond),
		otlp.WithAddress(mc.address),
		otlp.WithHeaders(map[string]string{"header1": "value1"}),
	)
	exp.ExportSpans(context.Background(), []*exporttrace.SpanData{{Name: "in the midst"}})

	defer func() {
		_ = exp.Stop()
	}()

	headers := mc.getHeaders()
	require.Len(t, headers.Get("header1"), 1)
	assert.Equal(t, "value1", headers.Get("header1")[0])
}
