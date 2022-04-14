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

package stdoutmetric_test

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/reader"
)

const (
	instrumentationName    = "github.com/instrumentron"
	instrumentationVersion = "v0.1.0"
)

var (
	// TODO Bring back Global package
	// meter = global.GetMeterProvider().Meter(
	// 	instrumentationName,
	// 	metric.WithInstrumentationVersion(instrumentationVersion),
	// )
	meter metric.Meter

	loopCounter syncint64.Counter
	paramValue  syncint64.Histogram

	nameKey = attribute.Key("function.name")
)

func add(ctx context.Context, x, y int64) int64 {
	nameKV := nameKey.String("add")

	loopCounter.Add(ctx, 1, nameKV)
	paramValue.Record(ctx, x, nameKV)
	paramValue.Record(ctx, y, nameKV)

	return x + y
}

func multiply(ctx context.Context, x, y int64) int64 {
	nameKV := nameKey.String("multiply")

	loopCounter.Add(ctx, 1, nameKV)
	paramValue.Record(ctx, x, nameKV)
	paramValue.Record(ctx, y, nameKV)

	return x * y
}

func InstallExportPipeline(ctx context.Context) *reader.ManualReader {
	exporter := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	reader := reader.NewManualReader(exporter)

	mp := sdkmetric.New(
		sdkmetric.WithReader(reader),
	)
	// TODO Bring back Global package
	// global.SetMeterProvider(pusher)
	meter = mp.Meter(instrumentationName, metric.WithInstrumentationVersion(instrumentationVersion))

	var err error

	loopCounter, err = meter.SyncInt64().Counter("function.loops")
	if err != nil {
		log.Fatalf("creating instrument: %v", err)
	}
	paramValue, err = meter.SyncInt64().Histogram("function.param")
	if err != nil {
		log.Fatalf("creating instrument: %v", err)
	}

	return reader
}

func Example() {
	ctx := context.Background()

	// TODO: Registers a meter Provider globally.
	reader := InstallExportPipeline(ctx)

	log.Println("the answer is", add(ctx, multiply(ctx, multiply(ctx, 2, 2), 10), 2))

	reader.Collect(ctx, nil)
}
