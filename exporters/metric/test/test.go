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

package test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"
)

type mapkey struct {
	desc     *metric.Descriptor
	distinct label.Distinct
}

type ckptRecord struct {
	export.Record
	export.Aggregator
}

type CheckpointSet struct {
	sync.RWMutex
	records  map[mapkey]ckptRecord
	updates  []export.Record
	resource *resource.Resource
}

// NoopAggregator is useful for testing Exporters.
type NoopAggregator struct{}

var _ export.Aggregator = (*NoopAggregator)(nil)

// Update implements export.Aggregator.
func (*NoopAggregator) Update(context.Context, metric.Number, *metric.Descriptor) error {
	return nil
}

// SynchronizedCopy implements export.Aggregator.
func (*NoopAggregator) SynchronizedCopy(export.Aggregator, *metric.Descriptor) error {
	return nil
}

// Merge implements export.Aggregator.
func (*NoopAggregator) Merge(export.Aggregator, *metric.Descriptor) error {
	return nil
}

// Kind implements aggregation.Aggregation.
func (*NoopAggregator) Kind() aggregation.Kind {
	return aggregation.NoopKind
}

// NewCheckpointSet returns a test CheckpointSet that new records could be added.
// Records are grouped by their encoded labels.
func NewCheckpointSet(resource *resource.Resource) *CheckpointSet {
	return &CheckpointSet{
		records:  map[mapkey]ckptRecord{},
		resource: resource,
	}
}

func (p *CheckpointSet) Reset() {
	p.records = map[mapkey]ckptRecord{}
	p.updates = nil
}

// Add a new descriptor to a Checkpoint.
//
// If there is an existing record with the same descriptor and labels,
// the stored aggregator will be returned and should be merged.
func (p *CheckpointSet) Add(desc *metric.Descriptor, newAgg export.Aggregator, labels ...kv.KeyValue) (agg export.Aggregator, added bool) {
	elabels := label.NewSet(labels...)

	key := mapkey{
		desc:     desc,
		distinct: elabels.Equivalent(),
	}
	if record, ok := p.records[key]; ok {
		return record.Aggregator, false
	}

	rec := export.NewRecord(desc, &elabels, p.resource, newAgg, time.Time{}, time.Time{})
	p.updates = append(p.updates, rec)
	p.records[key] = ckptRecord{
		Record:     rec,
		Aggregator: newAgg,
	}
	return newAgg, true
}

// ForEach exposes the records in this checkpoint. Note that this test
// does not make use of the ExporterKind argument: use a real Integrator
// for such testing.
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
