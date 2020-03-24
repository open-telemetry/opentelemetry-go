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
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

var (
	fooKey    = key.New("ex.com/foo")
	barKey    = key.New("ex.com/bar")
	lemonsKey = key.New("ex.com/lemons")
)

func initMeter() *push.Controller {
	pusher, hf, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}
	http.HandleFunc("/", hf)
	go func() {
		_ = http.ListenAndServe(":2222", nil)
	}()

	return pusher
}

func main() {
	defer initMeter().Stop()

	meter := global.Meter("ex.com/basic")
	observerLock := new(sync.RWMutex)
	observerValueToReport := new(float64)
	observerLabelSetToReport := new(metric.LabelSet)
	cb := func(result metric.Float64ObserverResult) {
		(*observerLock).RLock()
		value := *observerValueToReport
		labelset := *observerLabelSetToReport
		(*observerLock).RUnlock()
		result.Observe(value, labelset)
	}
	_ = metric.Must(meter).RegisterFloat64Observer("ex.com.one", cb,
		metric.WithKeys(fooKey, barKey, lemonsKey),
		metric.WithDescription("A measure set to 1.0"),
	)

	measureTwo := metric.Must(meter).NewFloat64Measure("ex.com.two", metric.WithKeys(key.New("A")))
	measureThree := metric.Must(meter).NewFloat64Counter("ex.com.three")

	commonLabels := meter.Labels(lemonsKey.Int(10), key.String("A", "1"), key.String("B", "2"), key.String("C", "3"))
	notSoCommonLabels := meter.Labels(lemonsKey.Int(13))

	ctx := context.Background()

	(*observerLock).Lock()
	*observerValueToReport = 1.0
	*observerLabelSetToReport = &commonLabels
	(*observerLock).Unlock()
	meter.RecordBatch(
		ctx,
		commonLabels,
		measureTwo.Measurement(2.0),
		measureThree.Measurement(12.0),
	)

	time.Sleep(5 * time.Second)

	(*observerLock).Lock()
	*observerValueToReport = 1.0
	*observerLabelSetToReport = &notSoCommonLabels
	(*observerLock).Unlock()
	meter.RecordBatch(
		ctx,
		notSoCommonLabels,
		measureTwo.Measurement(2.0),
		measureThree.Measurement(22.0),
	)

	time.Sleep(5 * time.Second)

	(*observerLock).Lock()
	*observerValueToReport = 13.0
	*observerLabelSetToReport = &commonLabels
	(*observerLock).Unlock()
	meter.RecordBatch(
		ctx,
		commonLabels,
		measureTwo.Measurement(12.0),
		measureThree.Measurement(13.0),
	)

	time.Sleep(100 * time.Second)
}
