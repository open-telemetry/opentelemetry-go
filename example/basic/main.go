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

package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/stdout"
)

var (
	fooKey     = kv.Key("ex.com/foo")
	barKey     = kv.Key("ex.com/bar")
	lemonsKey  = kv.Key("ex.com/lemons")
	anotherKey = kv.Key("ex.com/another")
)

func main() {
	pusher, err := stdout.InstallNewPipeline([]stdout.Option{
		stdout.WithQuantiles([]float64{0.5, 0.9, 0.99}),
		stdout.WithPrettyPrint(),
	}, nil)
	if err != nil {
		log.Fatalf("failed to initialize stdout export pipeline: %v", err)
	}
	defer pusher.Stop()

	tracer := global.Tracer("ex.com/basic")
	meter := global.Meter("ex.com/basic")

	commonLabels := []kv.KeyValue{lemonsKey.Int(10), kv.String("A", "1"), kv.String("B", "2"), kv.String("C", "3")}

	oneMetricCB := func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, commonLabels...)
	}
	_ = metric.Must(meter).NewFloat64ValueObserver("ex.com.one", oneMetricCB,
		metric.WithDescription("A ValueObserver set to 1.0"),
	)

	valuerecorderTwo := metric.Must(meter).NewFloat64ValueRecorder("ex.com.two")

	ctx := context.Background()

	ctx = correlation.NewContext(ctx,
		fooKey.String("foo1"),
		barKey.String("bar1"),
	)

	valuerecorder := valuerecorderTwo.Bind(commonLabels...)
	defer valuerecorder.Unbind()

	err = tracer.WithSpan(ctx, "operation", func(ctx context.Context) error {

		trace.SpanFromContext(ctx).AddEvent(ctx, "Nice operation!", kv.Key("bogons").Int(100))

		trace.SpanFromContext(ctx).SetAttributes(anotherKey.String("yes"))

		meter.RecordBatch(
			// Note: call-site variables added as context Entries:
			correlation.NewContext(ctx, anotherKey.String("xyz")),
			commonLabels,

			valuerecorderTwo.Measurement(2.0),
		)

		return tracer.WithSpan(
			ctx,
			"Sub operation...",
			func(ctx context.Context) error {
				trace.SpanFromContext(ctx).SetAttributes(lemonsKey.String("five"))

				trace.SpanFromContext(ctx).AddEvent(ctx, "Sub span event")

				valuerecorder.Record(ctx, 1.3)

				return nil
			},
		)
	})
	if err != nil {
		panic(err)
	}
}
