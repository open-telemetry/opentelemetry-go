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

package metric

import (
	"context"
	"sync/atomic"
)

var (
	descriptorID uint64
)

// TODO: Maybe unexport that and document very _very_ clearly, that
// you can still get a descriptor with NewInt64Counter(…).Descriptor

// CommonMetric holds a descriptor. It is used mostly to implement the
// common parts for every metric kind.
type CommonMetric struct {
	*Descriptor
}

func (m CommonMetric) getHandle(labels LabelSet) Handle {
	return labels.Meter().NewHandle(m.Descriptor, labels)
}

func (m CommonMetric) float64Measurement(value float64) Measurement {
	return Measurement{
		Descriptor: m.Descriptor,
		Value:      NewFloat64MeasurementValue(value),
	}
}

func (m CommonMetric) int64Measurement(value int64) Measurement {
	return Measurement{
		Descriptor: m.Descriptor,
		Value:      NewInt64MeasurementValue(value),
	}
}

func (m CommonMetric) recordOne(ctx context.Context, value MeasurementValue, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, Measurement{
		Descriptor: m.Descriptor,
		Value:      value,
	})
}

func registerCommonMetric(name string, kind Kind, valueKind ValueKind) CommonMetric {
	return CommonMetric{
		Descriptor: registerDescriptor(name, kind, valueKind),
	}
}

func registerDescriptor(name string, kind Kind, valueKind ValueKind) *Descriptor {
	return &Descriptor{
		name:      name,
		kind:      kind,
		valueKind: valueKind,
		id:        DescriptorID(atomic.AddUint64(&descriptorID, 1)),
	}
}
