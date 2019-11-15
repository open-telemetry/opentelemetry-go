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

package gauge // import "go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
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

		// checkpoint is a copy of the current value taken in Checkpoint()
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

var _ export.Aggregator = &Aggregator{}
var _ aggregator.LastValue = &Aggregator{}

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

// LastValue returns the last-recorded gauge value and the
// corresponding timestamp.  The error value aggregator.ErrNoLastValue
// will be returned if (due to a race condition) the checkpoint was
// computed before the first value was set.
func (g *Aggregator) LastValue() (core.Number, time.Time, error) {
	gd := (*gaugeData)(g.checkpoint)
	if gd == unsetGauge {
		return core.Number(0), time.Time{}, aggregator.ErrNoLastValue
	}
	return gd.value.AsNumber(), gd.timestamp, nil
}

// Checkpoint atomically saves the current value.
func (g *Aggregator) Checkpoint(ctx context.Context, _ *export.Descriptor) {
	g.checkpoint = atomic.LoadPointer(&g.current)
}

// Update atomically sets the current "last" value.
func (g *Aggregator) Update(_ context.Context, number core.Number, desc *export.Descriptor) error {
	if !desc.Alternate() {
		g.updateNonMonotonic(number)
		return nil
	}
	return g.updateMonotonic(number, desc)
}

func (g *Aggregator) updateNonMonotonic(number core.Number) {
	ngd := &gaugeData{
		value:     number,
		timestamp: time.Now(),
	}
	atomic.StorePointer(&g.current, unsafe.Pointer(ngd))
}

func (g *Aggregator) updateMonotonic(number core.Number, desc *export.Descriptor) error {
	ngd := &gaugeData{
		timestamp: time.Now(),
		value:     number,
	}
	kind := desc.NumberKind()

	for {
		gd := (*gaugeData)(atomic.LoadPointer(&g.current))

		if gd.value.CompareNumber(kind, number) > 0 {
			return aggregator.ErrNonMonotoneInput
		}

		if atomic.CompareAndSwapPointer(&g.current, unsafe.Pointer(gd), unsafe.Pointer(ngd)) {
			return nil
		}
	}
}

// Merge combines state from two aggregators.  If the gauge is
// declared as monotonic, the greater value is chosen.  If the gauge
// is declared as non-monotonic, the most-recently set value is
// chosen.
func (g *Aggregator) Merge(oa export.Aggregator, desc *export.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(g, oa)
	}

	ggd := (*gaugeData)(atomic.LoadPointer(&g.checkpoint))
	ogd := (*gaugeData)(atomic.LoadPointer(&o.checkpoint))

	if desc.Alternate() {
		// Monotonic: use the greater value
		cmp := ggd.value.CompareNumber(desc.NumberKind(), ogd.value)

		if cmp > 0 {
			return nil
		}

		if cmp < 0 {
			g.checkpoint = unsafe.Pointer(ogd)
			return nil
		}
	}
	// Non-monotonic gauge or equal values
	if ggd.timestamp.After(ogd.timestamp) {
		return nil
	}

	g.checkpoint = unsafe.Pointer(ogd)
	return nil
}
