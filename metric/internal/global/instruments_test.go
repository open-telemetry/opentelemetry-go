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

package global

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/nonrecording"
)

func Test_afCounter_setDelegate(t *testing.T) {
	delegate := afCounter{
		name: "testName",
		opts: []instrument.Option{},
	}

	go func() {
		for {
			delegate.Observe(context.Background(), 1)
		}
	}()

	delegate.setDelegate(nonrecording.NewNoopMeter())
}

type testCountingFloatInstrument struct {
	count int

	instrument.Asynchronous
	instrument.Synchronous
}

func (i *testCountingFloatInstrument) Observe(context.Context, float64, ...attribute.KeyValue) {
	i.count++
}
func (i *testCountingFloatInstrument) Add(context.Context, float64, ...attribute.KeyValue) {
	i.count++
}
func (i *testCountingFloatInstrument) Record(context.Context, float64, ...attribute.KeyValue) {
	i.count++
}

type testCountingIntInstrument struct {
	count int

	instrument.Asynchronous
	instrument.Synchronous
}

func (i *testCountingIntInstrument) Observe(context.Context, int64, ...attribute.KeyValue) {
	i.count++
}
func (i *testCountingIntInstrument) Add(context.Context, int64, ...attribute.KeyValue) {
	i.count++
}
func (i *testCountingIntInstrument) Record(context.Context, int64, ...attribute.KeyValue) {
	i.count++
}
