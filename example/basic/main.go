// Copyright 2019, OpenTelemetry Authors
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

	"go.opentelemetry.io/api/distributedcontext"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/global"
)

var (
	tracer = global.TraceProvider().GetTracer("ex.com/basic")

	meter = metric.GlobalMeter() // TODO: should share resources ^^^?

	fooKey     = key.New("ex.com/foo")
	barKey     = key.New("ex.com/bar")
	lemonsKey  = key.New("ex.com/lemons")
	anotherKey = key.New("ex.com/another")
)

func main() {
	oneMetric := meter.NewFloat64Gauge("ex.com.one",
		metric.WithKeys(fooKey, barKey, lemonsKey),
		metric.WithDescription("A gauge set to 1.0"),
	)

	measureTwo := meter.NewFloat64Measure("ex.com.two")

	ctx := context.Background()

	ctx = distributedcontext.NewContext(ctx,
		distributedcontext.Insert(fooKey.String("foo1")),
		distributedcontext.Insert(barKey.String("bar1")),
	)

	commonLabels := meter.Labels(lemonsKey.Int(10))

	gauge := oneMetric.AcquireHandle(commonLabels)
	defer gauge.Release()

	measure := measureTwo.AcquireHandle(commonLabels)
	defer measure.Release()

	err := tracer.WithSpan(ctx, "operation", func(ctx context.Context) error {

		trace.CurrentSpan(ctx).AddEvent(ctx, "Nice operation!", key.New("bogons").Int(100))

		trace.CurrentSpan(ctx).SetAttributes(anotherKey.String("yes"))

		gauge.Set(ctx, 1)

		meter.RecordBatch(
			// Note: call-site variables added as context Entries:
			distributedcontext.NewContext(ctx,
				distributedcontext.Insert(anotherKey.String("xyz"))),
			commonLabels,

			oneMetric.Measurement(1.0),
			measureTwo.Measurement(2.0),
		)

		return tracer.WithSpan(
			ctx,
			"Sub operation...",
			func(ctx context.Context) error {
				trace.CurrentSpan(ctx).SetAttribute(lemonsKey.String("five"))

				trace.CurrentSpan(ctx).AddEvent(ctx, "Sub span event")

				measure.Record(ctx, 1.3)

				return nil
			},
		)
	})
	if err != nil {
		panic(err)
	}

	// TODO: How to flush?
	// loader.Flush()
}
