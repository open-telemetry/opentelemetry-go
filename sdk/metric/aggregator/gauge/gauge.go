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

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/sdk/export"
)

// Note: This aggregator enforces the behavior of monotonic gauges to
// the best of its ability, but it will not retain any memory of
// infrequently used gauges.  Exporters may wish to enforce this, or
// they may simply treat monotonic as a semantic hint.

type (

	// Aggregator aggregates gauge events.
	Aggregator struct {
		// data is an atomic pointer to *gaugeData.  It is set
		// to `nil` if the gauge has not been set since the
		// last collection.
		live unsafe.Pointer

		// N.B. Export is not called when save is nil
		save unsafe.Pointer
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

// New returns a new gauge aggregator.  This aggregator retains the
// last value and timestamp that were recorded.
func New() *Aggregator {
	return &Aggregator{}
}

// AsInt64 returns the recorded gauge value as an int64.
func (g *Aggregator) AsNumber() core.Number {
	return (*gaugeData)(g.save).value.AsNumber()
}

// Timestamp returns the timestamp of the alst recorded gauge value.
func (g *Aggregator) Timestamp() time.Time {
	return (*gaugeData)(g.save).timestamp
}

// Collect saves the current value (atomically) and exports it.
func (g *Aggregator) Collect(ctx context.Context, rec export.MetricRecord, exp export.MetricBatcher) {
	g.save = atomic.SwapPointer(&g.live, nil)

	if g.save == nil {
		// There is no current value. This indicates a harmless race
		// involving Collect() and DeleteHandle().
		return
	}

	exp.Export(ctx, rec, g)
}

// Collect updates the current value (atomically) for later export.
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
	atomic.StorePointer(&g.live, unsafe.Pointer(ngd))
}

func (g *Aggregator) updateMonotonic(number core.Number, desc export.Descriptor) {
	ngd := &gaugeData{
		timestamp: time.Now(),
	}
	kind := desc.NumberKind()

	for {
		gd := (*gaugeData)(atomic.LoadPointer(&g.live))

		if gd != nil {
			if gd.value.CompareNumber(kind, number) > 0 {
				// TODO warn
				return
			}
		}
		ngd.value = number

		if atomic.CompareAndSwapPointer(&g.live, unsafe.Pointer(gd), unsafe.Pointer(ngd)) {
			return
		}
	}
}
