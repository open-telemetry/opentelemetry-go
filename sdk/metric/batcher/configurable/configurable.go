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

package configurable // import "go.opentelemetry.io/otel/sdk/metric/batcher/configurable"

import (
	"context"

	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type Config []Item

type Item struct {
	Name       string
	Aggregator string
}

type Integrator struct {
}

var _ export.Batcher = (*Integrator)(nil)

func New(cfg Config) *Integrator {
	return &Integrator{}
}

func (ci *Integrator) AggregatorFor(*metric.Descriptor) export.Aggregator {
	return nil
}

func (ci *Integrator) Process(ctx context.Context, record export.Record) error {
	return nil
}

func (ci *Integrator) CheckpointSet() export.CheckpointSet {
	return nil
}

func (ci *Integrator) FinishedCollection() {
}
