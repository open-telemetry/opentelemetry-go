// Copyright 2020, OpenTelemetry Authors
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

package observermmsc // import "go.opentelemetry.io/otel/sdk/metric/aggregator/observermmsc"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
)

type (
	Aggregator struct {
		g    *gauge.Aggregator
		mmsc *minmaxsumcount.Aggregator
	}
)

var (
	_ export.Aggregator         = &Aggregator{}
	_ aggregator.LastValue      = &Aggregator{}
	_ aggregator.MinMaxSumCount = &Aggregator{}
)

func New(desc *export.Descriptor) *Aggregator {
	return &Aggregator{
		g:    gauge.New(),
		mmsc: minmaxsumcount.New(desc),
	}
}

func (a *Aggregator) Update(ctx context.Context, number core.Number, descriptor *export.Descriptor) error {
	if err := a.g.Update(ctx, number, descriptor); err != nil {
		return err
	}
	a.mmsc.UpdateMMSC(number, descriptor)
	return nil
}

func (a *Aggregator) Checkpoint(ctx context.Context, descriptor *export.Descriptor) {
	a.g.Checkpoint(ctx, descriptor)
	a.mmsc.Checkpoint(ctx, descriptor)
}

func (a *Aggregator) Merge(oa export.Aggregator, descriptor *export.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(a, oa)
	}

	a.g.MergeGaugeAggregator(o.g, descriptor)
	a.mmsc.MergeMMSCAggregator(o.mmsc, descriptor)
	return nil
}

func (a *Aggregator) LastValue() (core.Number, time.Time, error) {
	return a.g.LastValue()
}

func (a *Aggregator) Min() (core.Number, error) {
	return a.mmsc.Min()
}

func (a *Aggregator) Max() (core.Number, error) {
	return a.mmsc.Max()
}

func (a *Aggregator) Sum() (core.Number, error) {
	return a.mmsc.Sum()
}

func (a *Aggregator) Count() (int64, error) {
	return a.mmsc.Count()
}
