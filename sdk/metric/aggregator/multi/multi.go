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

package multi // import "go.opentelemetry.io/otel/sdk/metric/aggregator/multi"

import (
	"context"

	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type Aggregator []export.Aggregator

// TODO: improve error handling (compared w/ below)?

var _ export.Aggregator = Aggregator{}
var _ aggregator.Multi = Aggregator{}

func New(aggs ...export.Aggregator) Aggregator {
	return Aggregator(aggs)
}

func (mu Aggregator) Aggregators() []export.Aggregator {
	return mu
}

func (mu Aggregator) Update(ctx context.Context, num metric.Number, desc *metric.Descriptor) error {
	var err error
	for _, agg := range mu {
		err1 := agg.Update(ctx, num, desc)
		if err1 != nil {
			err = err1
		}
	}
	return err
}

func (mu Aggregator) Checkpoint(ctx context.Context, desc *metric.Descriptor) {
	for _, agg := range mu {
		agg.Checkpoint(ctx, desc)
	}
}

func (mu Aggregator) Merge(agg export.Aggregator, desc *metric.Descriptor) error {
	var err error
	for _, agg := range mu {
		err1 := agg.Merge(agg, desc)
		if err1 != nil {
			err = err1
		}
	}
	return err
}
