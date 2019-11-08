package test

import (
	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type Producer struct {
	encoder export.LabelEncoder
	updates []export.Record
}

func NewProducer(encoder export.LabelEncoder) *Producer {
	return &Producer{
		encoder: encoder,
	}
}

func (p *Producer) Add(desc *export.Descriptor, agg export.Aggregator, labels ...core.KeyValue) {
	encoded := p.encoder.EncodeLabels(labels)
	elabels := export.NewLabels(labels, encoded, p.encoder)

	p.updates = append(p.updates, export.NewRecord(desc, elabels, agg))
}

func (p *Producer) Foreach(f func(export.Record)) {
	for _, r := range p.updates {
		f(r)
	}
}
