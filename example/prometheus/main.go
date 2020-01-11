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
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporter/metric/prometheus"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

var (
	fooKey    = key.New("ex.com/foo")
	barKey    = key.New("ex.com/bar")
	lemonsKey = key.New("ex.com/lemons")
)

func initMeter() *push.Controller {
	pusher, hf, err := prometheus.NewExportPipeline(prometheus.Config{})
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}
	http.HandleFunc("/", hf)
	go func() {
		_ = http.ListenAndServe(":2222", nil)
	}()
	global.SetScope(scope.WithMeterSDK(pusher.Meter()).WithNamespace("ex.com/basic"))
	return pusher
}

func main() {
	defer initMeter().Stop()

	oneMetric := metric.NewFloat64Gauge("one",
		metric.WithKeys(fooKey, barKey, lemonsKey),
		metric.WithDescription("A gauge set to 1.0"),
	)

	measureTwo := metric.NewFloat64Measure("two", metric.WithKeys(key.New("A")))
	measureThree := metric.NewFloat64Counter("three")

	ctx := global.Scope().AddResources(
		lemonsKey.Int(10),
		key.String("A", "1"),
		key.String("B", "2"),
		key.String("C", "3"),
	).InContext(context.Background())

	extraLabels := []core.KeyValue{
		barKey.Bool(false),
		lemonsKey.Int(13),
	}

	metric.RecordBatch(
		ctx,
		nil,
		oneMetric.Measurement(1.0),
		measureTwo.Measurement(2.0),
		measureThree.Measurement(12.0),
	)

	metric.RecordBatch(
		ctx,
		extraLabels,
		oneMetric.Measurement(1.0),
		measureTwo.Measurement(2.0),
		measureThree.Measurement(22.0),
	)

	time.Sleep(5 * time.Second)

	metric.RecordBatch(
		ctx,
		nil,
		oneMetric.Measurement(13.0),
		measureTwo.Measurement(12.0),
		measureThree.Measurement(13.0),
	)

	time.Sleep(100 * time.Second)
}
