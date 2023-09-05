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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func ExampleMeter_synchronous() {
	// Create a histogram using the global MeterProvider.
	workDuration, err := otel.Meter("go.opentelemetry.io/otel/metric#SyncExample").Int64Histogram(
		"workDuration",
		metric.WithUnit("ms"))
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

func ExampleMeter_asynchronous_single() {
	meter := otel.Meter("go.opentelemetry.io/otel/metric#AsyncExample")

	_, err := meter.Int64ObservableGauge(
		"DiskUsage",
		metric.WithUnit("By"),
		metric.WithInt64Callback(func(_ context.Context, obsrv metric.Int64Observer) error {
			// Do the real work here to get the real disk usage. For example,
			//
			//   usage, err := GetDiskUsage(diskID)
			//   if err != nil {
			//   	if retryable(err) {
			//   		// Retry the usage measurement.
			//   	} else {
			//   		return err
			//   	}
			//   }
			//
			// For demonstration purpose, a static value is used here.
			usage := 75000
			obsrv.Observe(int64(usage), metric.WithAttributes(attribute.Int("disk.id", 3)))
			return nil
		}),
	)
	if err != nil {
		fmt.Println("failed to register instrument")
		panic(err)
	}
}

func ExampleMeter_asynchronous_multiple() {
	meter := otel.Meter("go.opentelemetry.io/otel/metric#MultiAsyncExample")

	// This is just a sample of memory stats to record from the Memstats
	heapAlloc, err := meter.Int64ObservableUpDownCounter("heapAllocs")
	if err != nil {
		fmt.Println("failed to register updown counter for heapAllocs")
		panic(err)
	}
	gcCount, err := meter.Int64ObservableCounter("gcCount")
	if err != nil {
		fmt.Println("failed to register counter for gcCount")
		panic(err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			memStats := &runtime.MemStats{}
			// This call does work
			runtime.ReadMemStats(memStats)

			o.ObserveInt64(heapAlloc, int64(memStats.HeapAlloc))
			o.ObserveInt64(gcCount, int64(memStats.NumGC))

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
