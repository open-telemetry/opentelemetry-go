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

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
)

//nolint:govet // Meter doesn't register for go vet
func ExampleMeter_synchronous() {
	// In a library or program this would be provided by otel.GetMeterProvider().
	meterProvider := metric.NewNoopMeterProvider()

	workDuration, err := meterProvider.Meter("go.opentelemetry.io/otel/metric#SyncExample").Int64Histogram(
		"workDuration",
		metric.WithUnit(unit.Milliseconds),
	)
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
	meterProvider := metric.NewNoopMeterProvider()
	meter := meterProvider.Meter("go.opentelemetry.io/otel/metric#AsyncExample")

	var memoryUsage metric.Int64ObservableGauge
	memoryUsage, err := meter.Int64ObservableGauge(
		"MemoryUsage",
		metric.WithUnit(unit.Bytes),
		metric.WithCallback(func(ctx context.Context) error {
			// Do Work to get the real memoryUsage.
			// mem := GatherMemory(ctx)
			mem := 75000

			memoryUsage.Observe(ctx, int64(mem))
			return nil
		}),
	)
	if err != nil {
		fmt.Println("Failed to register instrument")
		panic(err)
	}
}

//nolint:govet // Meter doesn't register for go vet
func ExampleMeter_asynchronous_multiple() {
	meterProvider := metric.NewNoopMeterProvider()
	meter := meterProvider.Meter("go.opentelemetry.io/otel/metric#MultiAsyncExample")

	heapAlloc, _ := meter.Int64ObservableUpDownCounter("heapAllocs")
	gcCount, _ := meter.Int64ObservableCounter("gcCount")
	gcPause, _ := meter.Float64Histogram("gcPause")

	_, err := meter.RegisterCallback(
		func(ctx context.Context) error {
			memStats := &runtime.MemStats{}
			runtime.ReadMemStats(memStats)

			heapAlloc.Observe(ctx, int64(memStats.HeapAlloc))
			gcCount.Observe(ctx, int64(memStats.NumGC))

			// Synchronously records the GC pauses.
			computeGCPauses(ctx, gcPause, memStats.PauseNs[:])

			return nil
		},
		heapAlloc,
		gcCount,
	)
	if err != nil {
		fmt.Println("Failed to register callback")
		panic(err)
	}
}

// This is just an example, see the the contrib runtime instrumentation for real implementation.
func computeGCPauses(ctx context.Context, recorder metric.Float64Histogram, pauseBuff []uint64) {}
