package test

import (
	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
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

func (p *CheckpointSet) ForEach(f func(export.Record)) {
	for _, r := range p.updates {
		f(r)
	}
}
