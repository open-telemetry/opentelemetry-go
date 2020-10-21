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

package metrictest // import "go.opentelemetry.io/otel/sdk/export/metric/metrictest"

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"
)

type mapkey struct {
	desc     *otel.Descriptor
	distinct label.Distinct
}

// CheckpointSet is useful for testing Exporters.
// TODO(#872): Uses of this can be replaced by processortest.Output.
type CheckpointSet struct {
	sync.RWMutex
	records  map[mapkey]export.Record
	updates  []export.Record
	resource *resource.Resource
}

// NoopAggregator is useful for testing Exporters.
type NoopAggregator struct{}

var _ export.Aggregator = (*NoopAggregator)(nil)

// Update implements export.Aggregator.
func (NoopAggregator) Update(context.Context, otel.Number, *otel.Descriptor) error {
	return nil
}

// SynchronizedMove implements export.Aggregator.
func (NoopAggregator) SynchronizedMove(export.Aggregator, *otel.Descriptor) error {
	return nil
}

// Merge implements export.Aggregator.
func (NoopAggregator) Merge(export.Aggregator, *otel.Descriptor) error {
	return nil
}

// Aggregation returns an interface for reading the state of this aggregator.
func (NoopAggregator) Aggregation() aggregation.Aggregation {
	return NoopAggregator{}
}

// Kind implements aggregation.Aggregation.
func (NoopAggregator) Kind() aggregation.Kind {
	return aggregation.Kind("Noop")
}

// NewCheckpointSet returns a test CheckpointSet that new records could be added.
// Records are grouped by their encoded labels.
func NewCheckpointSet(resource *resource.Resource) *CheckpointSet {
	return &CheckpointSet{
		records:  make(map[mapkey]export.Record),
		resource: resource,
	}
}

// Reset clears the Aggregator state.
func (p *CheckpointSet) Reset() {
	p.records = make(map[mapkey]export.Record)
	p.updates = nil
}

// Add a new record to a CheckpointSet.
//
// If there is an existing record with the same descriptor and labels,
// the stored aggregator will be returned and should be merged.
func (p *CheckpointSet) Add(desc *otel.Descriptor, newAgg export.Aggregator, labels ...label.KeyValue) (agg export.Aggregator, added bool) {
	elabels := label.NewSet(labels...)

	key := mapkey{
		desc:     desc,
		distinct: elabels.Equivalent(),
	}
	if record, ok := p.records[key]; ok {
		return record.Aggregation().(export.Aggregator), false
	}

	rec := export.NewRecord(desc, &elabels, p.resource, newAgg.Aggregation(), time.Time{}, time.Time{})
	p.updates = append(p.updates, rec)
	p.records[key] = rec
	return newAgg, true
}

// ForEach does not use ExportKindSelected: use a real Processor to
// test ExportKind functionality.
func (p *CheckpointSet) ForEach(_ export.ExportKindSelector, f func(export.Record) error) error {
	for _, r := range p.updates {
		if err := f(r); err != nil && !errors.Is(err, aggregation.ErrNoData) {
			return err
		}
	}
	return nil
}

// Takes a slice of []some.Aggregator and returns a slice of []export.Aggregator
func Unslice2(sl interface{}) (one, two export.Aggregator) {
	slv := reflect.ValueOf(sl)
	if slv.Type().Kind() != reflect.Slice {
		panic("Invalid Unslice2")
	}
	if slv.Len() != 2 {
		panic("Invalid Unslice2: length > 2")
	}
	one = slv.Index(0).Addr().Interface().(export.Aggregator)
	two = slv.Index(1).Addr().Interface().(export.Aggregator)
	return
}
