package test

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
)

type CheckpointSet struct {
	encoder export.LabelEncoder
	records map[string]export.Record
	updates []export.Record
}

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

func (p *CheckpointSet) Add(desc *export.Descriptor, agg export.Aggregator, labels ...core.KeyValue) export.Aggregator {
	encoded := p.encoder.Encode(labels)
	elabels := export.NewLabels(labels, encoded, p.encoder)

	key := desc.Name() + "_" + elabels.Encoded()
	if record, ok := p.records[key]; ok {
		return record.Aggregator()
	}

	rec := export.NewRecord(desc, elabels, agg)
	p.updates = append(p.updates, rec)
	p.records[key] = rec
	return agg
}

func createNumber(desc *export.Descriptor, v float64) core.Number {
	if desc.NumberKind() == core.Float64NumberKind {
		return core.NewFloat64Number(v)
	}
	return core.NewInt64Number(int64(v))
}

func (p *CheckpointSet) AddGauge(desc *export.Descriptor, v float64, labels ...core.KeyValue) {
	ctx := context.Background()
	gagg := gauge.New()
	agg := p.Add(desc, gagg, labels...)
	_ = gagg.Update(ctx, createNumber(desc, v), desc)
	gagg.Checkpoint(ctx, desc)
	if agg != gagg {
		_ = agg.Merge(gagg, desc)
	}
}

func (p *CheckpointSet) AddCounter(desc *export.Descriptor, v float64, labels ...core.KeyValue) {
	ctx := context.Background()
	cagg := counter.New()
	agg := p.Add(desc, cagg, labels...)
	_ = cagg.Update(ctx, createNumber(desc, v), desc)
	cagg.Checkpoint(ctx, desc)
	if agg != cagg {
		_ = agg.Merge(cagg, desc)
	}
}

func (p *CheckpointSet) AddMeasure(desc *export.Descriptor, v float64, labels ...core.KeyValue) {
	ctx := context.Background()
	magg := array.New()
	agg := p.Add(desc, magg, labels...)
	_ = magg.Update(ctx, createNumber(desc, v), desc)
	magg.Checkpoint(ctx, desc)
	if agg != magg {
		_ = agg.Merge(magg, desc)
	}
}

func (p *CheckpointSet) ForEach(f func(export.Record)) {
	for _, r := range p.updates {
		f(r)
	}
}
