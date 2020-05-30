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

package configurable // import "go.opentelemetry.io/otel/sdk/metric/integrator/configurable"

import (
	"bytes"
	"context"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type (
	AggregationType int
)

const (
	Sum AggregationType = iota
	MinMaxSumCount
	Histogram
	LastValue
	Sketch
	Exact
)

type (
	Config struct {
		Defaults     `mapstructure:"defaults"`
		Views        `mapstructure:"views"`
		Aggregations `mapstructure:"aggregations"`
	}

	Defaults struct {
		// Instrument kind name to aggregation policy
		Aggregation map[string]string `mapstructure:"aggregation"`
	}

	// Instrument name to aggregation policy
	Views map[string]string

	Aggregations map[string]Aggregation

	Aggregation struct {
		Aggregator string   `mapstructure:"aggregator"`
		Labels     []string `mapstructure:"labels"`
	}

	Integrator struct {
		defaults map[string]*aggregation
		views    map[string][]*aggregation
	}

	aggregation struct {
	}
)

var _ export.Integrator = (*Integrator)(nil)

func ParseYamlData(data []byte) (cfg Config, err error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err = v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return
	}

	// Check for valid toplevel fields
	if err = v.UnmarshalExact(&cfg); err != nil {
		return
	}

	return
}

func New(cfg Config) *Integrator {
	defaults := map[string]*aggregation{}
	views := map[string][]*aggregation{}

	for _, view := range cfg.Views {
		// TODO here parse the names, return an error
		// views[view.InstrumentName] = append(views[view.InstrumentName], view)
		_ = view
	}
	for _, def := range cfg.Defaults.Aggregation {
		// TODO here same
		// defaults[def.InstrumentKind] = def
		_ = def
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
