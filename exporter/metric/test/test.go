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
	updates []export.Record
}

func NewCheckpointSet(encoder export.LabelEncoder) *CheckpointSet {
	return &CheckpointSet{
		encoder: encoder,
	}
}

func (p *CheckpointSet) Reset() {
	p.updates = nil
}

func (p *CheckpointSet) Add(desc *export.Descriptor, agg export.Aggregator, labels ...core.KeyValue) {
	encoded := p.encoder.Encode(labels)
	elabels := export.NewLabels(labels, encoded, p.encoder)

	p.updates = append(p.updates, export.NewRecord(desc, elabels, agg))
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
	_ = gagg.Update(ctx, createNumber(desc, v), desc)
	gagg.Checkpoint(ctx, desc)
	p.Add(desc, gagg, labels...)
}

func (p *CheckpointSet) AddCounter(desc *export.Descriptor, v float64, labels ...core.KeyValue) {
	ctx := context.Background()
	cagg := counter.New()
	_ = cagg.Update(ctx, createNumber(desc, v), desc)
	cagg.Checkpoint(ctx, desc)
	p.Add(desc, cagg, labels...)
}

func (p *CheckpointSet) AddMeasure(desc *export.Descriptor, v float64, labels ...core.KeyValue) {
	ctx := context.Background()
	magg := array.New()
	_ = magg.Update(ctx, createNumber(desc, v), desc)
	magg.Checkpoint(ctx, desc)
	p.Add(desc, magg, labels...)
}

func (p *CheckpointSet) ForEach(f func(export.Record)) {
	for _, r := range p.updates {
		f(r)
	}
}
