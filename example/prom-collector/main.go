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
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

func initMeter() {
	ctx := context.Background()

	res, err := resource.New(
		ctx,
		resource.WithAttributes(label.String("R", "V")),
	)
	if err != nil {
		log.Panic("could not initialize resource:", err)
	}

	otlpExporter, err := otlp.NewExporter(ctx,
		otlp.WithInsecure(),
		otlp.WithAddress("127.0.0.1:7001"),
		otlp.WithGRPCDialOption(grpc.WithBlock()), // useful for testing
	)
	if err != nil {
		log.Panic("could not initialize OTLP:", err)
	}

	options := []controller.Option{
		controller.WithResource(res),
		controller.WithExporter(otlpExporter),
	}

	cont := controller.New(
		processor.New(
			simple.NewWithHistogramDistribution([]float64{
				0.001, 0.01, 0.1, 1, 10, 100, 1000,
			}),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
		options...,
	)

	cont.Start()

	promExporter, err := prometheus.NewExporter(prometheus.Config{}, cont)
	if err != nil {
		log.Panic("could not initialize prometheus:", err)
	}
	http.HandleFunc("/", promExporter.ServeHTTP)
	go func() {
		_ = http.ListenAndServe(":2222", nil)
	}()

	otel.SetMeterProvider(cont.MeterProvider())

	fmt.Println("Prometheus server running on :2222")
	fmt.Println("Exporting OTLP to :7001")
}

func main() {
	initMeter()

	meter := otel.Meter("ex.com/basic")
	observerLock := new(sync.RWMutex)
	observerValueToReport := new(float64)
	observerLabelsToReport := new([]label.KeyValue)
	cb := func(_ context.Context, result metric.Float64ObserverResult) {
		(*observerLock).RLock()
		value := *observerValueToReport
		labels := *observerLabelsToReport
		(*observerLock).RUnlock()
		result.Observe(value, labels...)
	}
	_ = metric.Must(meter).NewFloat64ValueObserver("ex.com.one", cb,
		metric.WithDescription("A ValueObserver set to 1.0"),
	)

	valuerecorder := metric.Must(meter).NewFloat64ValueRecorder("ex.com.two")
	counter := metric.Must(meter).NewFloat64Counter("ex.com.three")

	commonLabels := []label.KeyValue{label.Int("I", 10), label.String("A", "1"), label.String("B", "2"), label.String("C", "3")}
	notSoCommonLabels := []label.KeyValue{label.Int("I", 13)}

	ctx := context.Background()

	(*observerLock).Lock()
	*observerValueToReport = 1.0
	*observerLabelsToReport = commonLabels
	(*observerLock).Unlock()
	meter.RecordBatch(
		ctx,
		commonLabels,
		valuerecorder.Measurement(2.0),
		counter.Measurement(12.0),
	)

	time.Sleep(5 * time.Second)

	(*observerLock).Lock()
	*observerValueToReport = 1.0
	*observerLabelsToReport = notSoCommonLabels
	(*observerLock).Unlock()
	meter.RecordBatch(
		ctx,
		notSoCommonLabels,
		valuerecorder.Measurement(2.0),
		counter.Measurement(22.0),
	)

	time.Sleep(5 * time.Second)

	(*observerLock).Lock()
	*observerValueToReport = 13.0
	*observerLabelsToReport = commonLabels
	(*observerLock).Unlock()
	meter.RecordBatch(
		ctx,
		commonLabels,
		valuerecorder.Measurement(12.0),
		counter.Measurement(13.0),
	)

	fmt.Println("Example finished updating, please visit :2222")

	select {}
}
