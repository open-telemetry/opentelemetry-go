package metric

import (
	"context"

	"go.opentelemetry.io/api/core"
)

type noopMeter struct{}
type noopHandle struct{}
type noopLabelSet struct{}

var _ Meter = noopMeter{}
var _ Handle = noopHandle{}
var _ LabelSet = noopLabelSet{}

func (noopHandle) RecordOne(context.Context, MeasurementValue) {
}

func (noopLabelSet) Meter() Meter {
	return noopMeter{}
}

func (noopMeter) DefineLabels(context.Context, ...core.KeyValue) LabelSet {
	return noopLabelSet{}
}

func (noopMeter) NewHandle(*Descriptor, LabelSet) Handle {
	return noopHandle{}
}

func (noopMeter) DeleteHandle(Handle) {
}

func (noopMeter) RecordBatch(context.Context, LabelSet, ...Measurement) {
}

func (noopMeter) RegisterObserver(Observer, ObserverCallback) {
}

func (noopMeter) UnregisterObserver(Observer) {
}
