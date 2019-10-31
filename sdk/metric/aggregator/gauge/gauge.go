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

package gauge

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/sdk/export"
)

// Note: This aggregator enforces the behavior of monotonic gauges to
// the best of its ability, but it will not retain any memory of
// infrequently used gauges.  Exporters may wish to enforce this, or
// they may simply treat monotonic as a semantic hint.

type (

	// Aggregator aggregates gauge events.
	Aggregator struct {
		// current is an atomic pointer to *gaugeData.  It is never nil.
		current unsafe.Pointer

		// checkpoint is a copy of the current value taken in Collect()
		checkpoint unsafe.Pointer
	}

	// gaugeData stores the current value of a gauge along with
	// a sequence number to determine the winner of a race.
	gaugeData struct {
		// value is the int64- or float64-encoded Set() data
		value core.Number

		// timestamp indicates when this record was submitted.
		// this can be used to pick a winner when multiple
		// records contain gauge data for the same labels due
		// to races.
		timestamp time.Time
	}
)

var _ export.MetricAggregator = &Aggregator{}

// An unset gauge has zero timestamp and zero value.
var unsetGauge = &gaugeData{}

// New returns a new gauge aggregator.  This aggregator retains the
// last value and timestamp that were recorded.
func New() *Aggregator {
	return &Aggregator{
		current:    unsafe.Pointer(unsetGauge),
		checkpoint: unsafe.Pointer(unsetGauge),
	}
}

// AsNumber returns the recorded gauge value as an int64.
func (g *Aggregator) AsNumber() core.Number {
	return (*gaugeData)(g.checkpoint).value.AsNumber()
}

// Timestamp returns the timestamp of the alst recorded gauge value.
func (g *Aggregator) Timestamp() time.Time {
	return (*gaugeData)(g.checkpoint).timestamp
}

// Collect checkpoints the current value (atomically) and exports it.
func (g *Aggregator) Collect(ctx context.Context, rec export.MetricRecord, exp export.MetricBatcher) {
	g.checkpoint = atomic.LoadPointer(&g.current)

	exp.Export(ctx, rec, g)
}

// Update modifies the current value (atomically) for later export.
func (g *Aggregator) Update(_ context.Context, number core.Number, rec export.MetricRecord) {
	desc := rec.Descriptor()
	if !desc.Alternate() {
		g.updateNonMonotonic(number)
	} else {
		g.updateMonotonic(number, desc)
	}
}

func (g *Aggregator) updateNonMonotonic(number core.Number) {
	ngd := &gaugeData{
		value:     number,
		timestamp: time.Now(),
	}
	atomic.StorePointer(&g.current, unsafe.Pointer(ngd))
}

func (g *Aggregator) updateMonotonic(number core.Number, desc *export.Descriptor) {
	ngd := &gaugeData{
		timestamp: time.Now(),
		value:     number,
	}
	kind := desc.NumberKind()

	for {
		gd := (*gaugeData)(atomic.LoadPointer(&g.current))

		if gd.value.CompareNumber(kind, number) > 0 {
			// TODO warn
			return
		}

		if atomic.CompareAndSwapPointer(&g.current, unsafe.Pointer(gd), unsafe.Pointer(ngd)) {
			return
		}
	}
}

func (g *Aggregator) Merge(oa export.MetricAggregator, desc *export.Descriptor) {
	o, _ := oa.(*Aggregator)
	if o == nil {
		// TODO warn
		return
	}

	ggd := (*gaugeData)(atomic.LoadPointer(&g.checkpoint))
	ogd := (*gaugeData)(atomic.LoadPointer(&o.checkpoint))

	if desc.Alternate() {
		// Monotonic: use the greater value
		cmp := ggd.value.CompareNumber(desc.NumberKind(), ogd.value)

		if cmp > 0 {
			return
		}

		if cmp < 0 {
			g.checkpoint = unsafe.Pointer(ogd)
			return
		}
	}
	// Non-monotonic gauge or equal values
	if ggd.timestamp.After(ogd.timestamp) {
		return
	}

	g.checkpoint = unsafe.Pointer(ogd)
}
