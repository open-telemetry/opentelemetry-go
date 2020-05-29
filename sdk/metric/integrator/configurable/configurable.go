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

type Config struct {
	Defaults []Default `json="default"`
	Views    []View    `json="view"`
}

type Default struct {
	InstrumentKind string   `json="instrument_kind"`
	Aggregator     string   `json="aggregator"`
	Labels         []string `json="labels"`
}

type View struct {
	InstrumentName string   `json="instrument_name"`
	Aggregator     string   `json="aggregator"`
	Labels         []string `json="labels"`
}

type Integrator struct {
	defaults map[string]Default
	views    map[string][]View
}

var _ export.Integrator = (*Integrator)(nil)

func New(cfg Config) *Integrator {
	defaults := map[string]Default{}
	views := map[string][]View{}

	for _, view := range cfg.Views {
		// TODO here parse the names, return an error
		views[view.InstrumentName] = append(views[view.InstrumentName], view)
	}
	for _, def := range cfg.Defaults {
		// TODO here same
		defaults[def.InstrumentKind] = def
	}

	return &Integrator{
		defaults: defaults,
		views:    views,
	}
}

func (ci *Integrator) AggregatorFor(desc *metric.Descriptor) export.Aggregator {
	views, ok := ci.views[desc.Name()]
	if !ok {
		if len(views) == 1 {
			return nil
		}
	}

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
