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
	"os"
	"os/signal"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	lemonsKey = attribute.Key("ex.com/lemons")
)

func main() {
	endpoint := "opentelemetry-collector:4318"

	client := otlpmetrichttp.NewClient(otlpmetrichttp.WithEndpoint(endpoint), otlpmetrichttp.WithInsecure())
	ctx := context.Background()
	exporter, err := otlpmetric.New(ctx, client)
	if err != nil {
		log.Fatalf("failed to create exporter: %s", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := exporter.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}()
	c := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(2*time.Second),
	)

	if err := c.Start(ctx); err != nil {
		log.Fatalf("could not start metric controller: %v", err)
	}
	global.SetMeterProvider(c)
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		// pushes any last exports to the receiver
		if err := c.Stop(ctx); err != nil {
			otel.Handle(err)
		}
	}()

	meter := global.Meter("ex.com/basic")

	observerLock := new(sync.RWMutex)
	observerValueToReport := new(float64)
	observerAttrsToReport := new([]attribute.KeyValue)

	gaugeObserver, err := meter.AsyncFloat64().Gauge("ex.com.one")
	if err != nil {
		log.Panicf("failed to initialize instrument: %v", err)
	}
	_ = meter.RegisterCallback([]instrument.Asynchronous{gaugeObserver}, func(ctx context.Context) {
		(*observerLock).RLock()
		value := *observerValueToReport
		attrs := *observerAttrsToReport
		(*observerLock).RUnlock()
		gaugeObserver.Observe(ctx, value, attrs...)
	})

	hist, err := meter.SyncFloat64().Histogram("ex.com.two")
	if err != nil {
		log.Panicf("failed to initialize instrument: %v", err)
	}
	counter, err := meter.SyncFloat64().Counter("ex.com.three")
	if err != nil {
		log.Panicf("failed to initialize instrument: %v", err)
	}

	commonAttrs := []attribute.KeyValue{lemonsKey.Int(10), attribute.String("A", "1"), attribute.String("B", "2"), attribute.String("C", "3")}
	notSoCommonAttrs := []attribute.KeyValue{lemonsKey.Int(13)}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	(*observerLock).Lock()
	*observerValueToReport = 1.0
	*observerAttrsToReport = commonAttrs
	(*observerLock).Unlock()

	hist.Record(ctx, 2.0, commonAttrs...)
	counter.Add(ctx, 12.0, commonAttrs...)

	time.Sleep(5 * time.Second)

	(*observerLock).Lock()
	*observerValueToReport = 1.0
	*observerAttrsToReport = notSoCommonAttrs
	(*observerLock).Unlock()
	hist.Record(ctx, 2.0, notSoCommonAttrs...)
	counter.Add(ctx, 22.0, notSoCommonAttrs...)

	time.Sleep(5 * time.Second)

	(*observerLock).Lock()
	*observerValueToReport = 13.0
	*observerAttrsToReport = commonAttrs
	(*observerLock).Unlock()
	hist.Record(ctx, 12.0, commonAttrs...)
	counter.Add(ctx, 13.0, commonAttrs...)

	fmt.Println("Example finished updating")
	<-ctx.Done()
}
