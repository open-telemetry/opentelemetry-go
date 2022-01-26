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

package asyncfloat64 // import "go.opentelemetry.io/otel/metric/instrument/asyncfloat64"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
)

func NewNoopInstruments() Instruments {
	return nonrecordingInstrument{}
}

type nonrecordingInstrument struct {
	instrument.Asynchronous
}

var (
	_ Instruments   = nonrecordingInstrument{}
	_ Counter       = nonrecordingInstrument{}
	_ UpDownCounter = nonrecordingInstrument{}
	_ Gauge         = nonrecordingInstrument{}
)

func (n nonrecordingInstrument) Counter(name string, opts ...instrument.Option) (Counter, error) {
	return n, nil
}

func (n nonrecordingInstrument) UpDownCounter(name string, opts ...instrument.Option) (UpDownCounter, error) {
	return n, nil
}

func (n nonrecordingInstrument) Gauge(name string, opts ...instrument.Option) (Gauge, error) {
	return n, nil
}

func (nonrecordingInstrument) Observe(context.Context, float64, ...attribute.KeyValue) {

}
