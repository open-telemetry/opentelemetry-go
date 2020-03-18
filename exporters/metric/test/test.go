package test

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

type CheckpointSet struct {
	encoder export.LabelEncoder
	records map[string]export.Record
	updates []export.Record
}

// NewCheckpointSet returns a test CheckpointSet that new records could be added.
// Records are grouped by their LabelSet.
func NewCheckpointSet(encoder export.LabelEncoder) *CheckpointSet {
	return &CheckpointSet{
		encoder: encoder,
		records: make(map[string]export.Record),
	}
}

func (p *CheckpointSet) Reset() {
	p.records = make(map[string]export.Record)
	p.updates = nil
}

// Add a new descriptor to a Checkpoint.
//
// If there is an existing record with the same descriptor and LabelSet
// the stored aggregator will be returned and should be merged.
func (p *CheckpointSet) Add(desc *metric.Descriptor, newAgg export.Aggregator, labels ...core.KeyValue) (agg export.Aggregator, added bool) {
	encoded := p.encoder.Encode(export.NewSliceLabelIterator(labels))
	elabels := export.NewLabels(export.LabelSlice(labels), encoded, p.encoder)

	key := desc.Name() + "_" + elabels.Encoded()
	if record, ok := p.records[key]; ok {
		return record.Aggregator(), false
	}

	rec := export.NewRecord(desc, elabels, newAgg)
	p.updates = append(p.updates, rec)
	p.records[key] = rec
	return newAgg, true
}

func createNumber(desc *metric.Descriptor, v float64) core.Number {
	if desc.NumberKind() == core.Float64NumberKind {
		return core.NewFloat64Number(v)
	}
	return core.NewInt64Number(int64(v))
}

func (p *CheckpointSet) AddLastValue(desc *metric.Descriptor, v float64, labels ...core.KeyValue) {
	p.updateAggregator(desc, lastvalue.New(), v, labels...)
}

func (p *CheckpointSet) AddCounter(desc *metric.Descriptor, v float64, labels ...core.KeyValue) {
	p.updateAggregator(desc, sum.New(), v, labels...)
}

func (p *CheckpointSet) AddMeasure(desc *metric.Descriptor, v float64, labels ...core.KeyValue) {
	p.updateAggregator(desc, array.New(), v, labels...)
}

func (p *CheckpointSet) updateAggregator(desc *metric.Descriptor, newAgg export.Aggregator, v float64, labels ...core.KeyValue) {
	ctx := context.Background()
	// Updates and checkpoint the new aggregator
	_ = newAgg.Update(ctx, createNumber(desc, v), desc)
	newAgg.Checkpoint(ctx, desc)

	// Try to add this aggregator to the CheckpointSet
	agg, added := p.Add(desc, newAgg, labels...)
	if !added {
		// An aggregator already exist for this descriptor and label set, we should merge them.
		_ = agg.Merge(newAgg, desc)
	}
}

func (p *CheckpointSet) ForEach(f func(export.Record) error) error {
	for _, r := range p.updates {
		if err := f(r); err != nil && !errors.Is(err, aggregator.ErrNoData) {
			return err
		}
	}
	return nil
}
