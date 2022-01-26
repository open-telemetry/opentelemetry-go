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

package syncint64 // import "go.opentelemetry.io/otel/metric/instrument/syncint64"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
)

func NewNoopInstruments() Instruments {
	return nonrecordingInstrument{}
}

type nonrecordingInstrument struct {
	instrument.Synchronous
}

var (
	_ Instruments   = nonrecordingInstrument{}
	_ Counter       = nonrecordingInstrument{}
	_ UpDownCounter = nonrecordingInstrument{}
	_ Histogram     = nonrecordingInstrument{}
)

func (n nonrecordingInstrument) Counter(name string, opts ...instrument.Option) (Counter, error) {
	return n, nil
}

func (n nonrecordingInstrument) UpDownCounter(name string, opts ...instrument.Option) (UpDownCounter, error) {
	return n, nil
}

func (n nonrecordingInstrument) Histogram(name string, opts ...instrument.Option) (Histogram, error) {
	return n, nil
}

func (nonrecordingInstrument) Add(context.Context, int64, ...attribute.KeyValue) {
}
func (nonrecordingInstrument) Record(context.Context, int64, ...attribute.KeyValue) {
}
