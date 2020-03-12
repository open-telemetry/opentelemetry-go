// Copyright 2020, OpenTelemetry Authors
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

	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"

	"go.opentelemetry.io/otel/api/core"
	metricapi "go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/otlp"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/metric/batcher/ungrouped"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
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

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(exp, // add following two options to ensure flush
			sdktrace.WithScheduleDelayMillis(15),
			sdktrace.WithMaxExportBatchSize(10),
		))
	assert.NoError(t, err)

	//global.SetTraceProvider(tp)

	tr := tp.Tracer("test-tracer")
	// Now create few spans
	m := 4
	for i := 0; i < m; i++ {
		_, span := tr.Start(context.Background(), "AlwaysSample")
		span.SetAttributes(core.Key("i").Int64(int64(i)))
		span.End()
	}

	selector := simple.NewWithExactMeasure()
	batcher := ungrouped.New(selector, true)
	pusher := push.New(batcher, exp, 60*time.Second)
	pusher.Start()

	ctx := context.Background()
	meter := pusher.Meter("test-meter")
	labels := meter.Labels(core.Key("test").Bool(true))

	// TODO: support observers
	type data struct {
		iKind metricsdk.Kind
		nKind core.NumberKind
		val   int64
	}
	instruments := map[string]data{
		"test-int64-counter":    {metricsdk.CounterKind, core.Int64NumberKind, 1},
		"test-float64-counter":  {metricsdk.CounterKind, core.Float64NumberKind, 1},
		"test-int64-measure":    {metricsdk.MeasureKind, core.Int64NumberKind, 2},
		"test-float64-measure":  {metricsdk.MeasureKind, core.Float64NumberKind, 2},
		"test-int64-observer":   {metricsdk.ObserverKind, core.Int64NumberKind, 3},
		"test-float64-observer": {metricsdk.ObserverKind, core.Float64NumberKind, 3},
	}
	for name, data := range instruments {
		switch data.iKind {
		case metricsdk.CounterKind:
			switch data.nKind {
			case core.Int64NumberKind:
				metricapi.Must(meter).NewInt64Counter(name).Add(ctx, data.val, labels)
			case core.Float64NumberKind:
				metricapi.Must(meter).NewFloat64Counter(name).Add(ctx, float64(data.val), labels)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case metricsdk.MeasureKind:
			switch data.nKind {
			case core.Int64NumberKind:
				metricapi.Must(meter).NewInt64Measure(name).Record(ctx, data.val, labels)
			case core.Float64NumberKind:
				metricapi.Must(meter).NewFloat64Measure(name).Record(ctx, float64(data.val), labels)
			default:
				assert.Failf(t, "unsupported number testing kind", data.nKind.String())
			}
		case metricsdk.ObserverKind:
			switch data.nKind {
			case core.Int64NumberKind:
				callback := func(v int64) metricapi.Int64ObserverCallback {
					return metricapi.Int64ObserverCallback(func(result metricapi.Int64ObserverResult) { result.Observe(v, labels) })
				}(data.val)
				metricapi.Must(meter).RegisterInt64Observer(name, callback)
			case core.Float64NumberKind:
				callback := func(v float64) metricapi.Float64ObserverCallback {
					return metricapi.Float64ObserverCallback(func(result metricapi.Float64ObserverResult) { result.Observe(v, labels) })
				}(float64(data.val))
				metricapi.Must(meter).RegisterFloat64Observer(name, callback)
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

	spans := mc.getSpans()

	// Now verify that we received all spans.
	if got, want := len(spans), m; got != want {
		t.Fatalf("span counts: got %d, want %d", got, want)
	}
	for i := 0; i < 4; i++ {
		if gotName, want := spans[i].Name, "AlwaysSample"; gotName != want {
			t.Fatalf("span name: got %s, want %s", gotName, want)
		}
		if got, want := spans[i].Attributes[0].IntValue, int64(i); got != want {
			t.Fatalf("span attribute value: got %d, want %d", got, want)
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
		case metricsdk.CounterKind:
			switch data.nKind {
			case core.Int64NumberKind:
				assert.Equal(t, metricpb.MetricDescriptor_COUNTER_INT64.String(), desc.GetType().String())
				if dp := m.GetInt64Datapoints(); assert.Len(t, dp, 1) {
					assert.Equal(t, data.val, dp[0].Value, "invalid value for %q", desc.Name)
				}
			case core.Float64NumberKind:
				assert.Equal(t, metricpb.MetricDescriptor_COUNTER_DOUBLE.String(), desc.GetType().String())
				if dp := m.GetDoubleDatapoints(); assert.Len(t, dp, 1) {
					assert.Equal(t, float64(data.val), dp[0].Value, "invalid value for %q", desc.Name)
				}
			default:
				assert.Failf(t, "invalid number kind", data.nKind.String())
			}
		case metricsdk.MeasureKind, metricsdk.ObserverKind:
			assert.Equal(t, metricpb.MetricDescriptor_SUMMARY.String(), desc.GetType().String())
			m.GetSummaryDatapoints()
			if dp := m.GetSummaryDatapoints(); assert.Len(t, dp, 1) {
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

		exp.ExportSpans(context.Background(), []*export.SpanData{{Name: "in the midst"}})

		// Now resurrect the collector by making a new one but reusing the
		// old address, and the collector should reconnect automatically.
		nmc := runMockColAtAddr(t, mc.address)

		// Give the exporter sometime to reconnect
		<-time.After(reconnectionPeriod * 4)

		n := 10
		for i := 0; i < n; i++ {
			exp.ExportSpans(context.Background(), []*export.SpanData{{Name: "Resurrected"}})
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
