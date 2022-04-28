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

package metric_test

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/nonrecording"
	"go.opentelemetry.io/otel/metric/unit"
)

//nolint:govet // Meter doesn't register for go vet
func ExampleMeter_synchronous() {
	// In a library or program this would be provided by otel.GetMeterProvider().
	meterProvider := nonrecording.NewNoopMeterProvider()

	workDuration, err := meterProvider.Meter("go.opentelemetry.io/otel/metric#SyncExample").SyncInt64().Histogram(
		"workDuration",
		instrument.WithUnit(unit.Milliseconds))
	if err != nil {
		fmt.Println("Failed to register instrument")
		panic(err)
	}

	startTime := time.Now()
	ctx := context.Background()
	// Do work
	// ...
	workDuration.Record(ctx, time.Since(startTime).Milliseconds())

}

//nolint:govet // Meter doesn't register for go vet
func ExampleMeter_asynchronous_single() {
	// In a library or program this would be provided by otel.GetMeterProvider().
	meterProvider := nonrecording.NewNoopMeterProvider()
	meter := meterProvider.Meter("go.opentelemetry.io/otel/metric#AsyncExample")

	memoryUsage, err := meter.AsyncInt64().Gauge(
		"MemoryUsage",
		instrument.WithUnit(unit.Bytes),
	)
	if err != nil {
		fmt.Println("Failed to register instrument")
		panic(err)
	}

	err = meter.RegisterCallback([]instrument.Asynchronous{memoryUsage},
		func(ctx context.Context) {
			// instrument.WithCallbackFunc(func(ctx context.Context) {
			//Do Work to get the real memoryUsage
			// mem := GatherMemory(ctx)
			mem := 75000

			memoryUsage.Observe(ctx, int64(mem))
		})
	if err != nil {
		fmt.Println("Failed to register callback")
		panic(err)
	}
}

//nolint:govet // Meter doesn't register for go vet
func ExampleMeter_asynchronous_multiple() {
	meterProvider := nonrecording.NewNoopMeterProvider()
	meter := meterProvider.Meter("go.opentelemetry.io/otel/metric#MultiAsyncExample")

	// This is just a sample of memory stats to record from the Memstats
	heapAlloc, _ := meter.AsyncInt64().UpDownCounter("heapAllocs")
	gcCount, _ := meter.AsyncInt64().Counter("gcCount")
	gcPause, _ := meter.SyncFloat64().Histogram("gcPause")

	err := meter.RegisterCallback([]instrument.Asynchronous{
		heapAlloc,
		gcCount,
	},
		func(ctx context.Context) {
			memStats := &runtime.MemStats{}
			// This call does work
			runtime.ReadMemStats(memStats)

			heapAlloc.Observe(ctx, int64(memStats.HeapAlloc))
			gcCount.Observe(ctx, int64(memStats.NumGC))

			// This function synchronously records the pauses
			computeGCPauses(ctx, gcPause, memStats.PauseNs[:])
		},
	)

	if err != nil {
		fmt.Println("Failed to register callback")
		panic(err)
	}
}

//This is just an example, see the the contrib runtime instrumentation for real implementation.
func computeGCPauses(ctx context.Context, recorder syncfloat64.Histogram, pauseBuff []uint64) {

}
