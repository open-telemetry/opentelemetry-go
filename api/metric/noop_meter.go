package metric

import (
	"context"

	"go.opentelemetry.io/api/core"
)

type noopMeter struct{}
type noopRecorder struct{}
type noopLabelSet struct{}

var _ Meter = noopMeter{}
var _ Recorder = noopRecorder{}
var _ LabelSet = noopLabelSet{}

func (noopRecorder) Record(ctx context.Context, value float64) {
}

func (noopLabelSet) Meter() Meter {
	return noopMeter{}
}

func (noopMeter) DefineLabels(ctx context.Context, labels ...core.KeyValue) LabelSet {
	return noopLabelSet{}
}

func (noopMeter) RecordSingle(context.Context, LabelSet, Measurement) {
}

func (noopMeter) RecordBatch(context.Context, LabelSet, ...Measurement) {
}

func (noopMeter) RecorderFor(context.Context, LabelSet, Descriptor) Recorder {
	return noopRecorder{}
}
