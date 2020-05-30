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
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/multi"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
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
	Views map[string][]string

	Aggregations map[string]Aggregation

	Aggregation struct {
		Aggregator string   `mapstructure:"aggregator"`
		Labels     []string `mapstructure:"labels"`
	}

	Integrator struct {
		policies    map[string]*aggregation
		instDefault map[metric.Kind]*aggregation
		views       map[string][]*aggregation
	}

	newFunc func(desc *metric.Descriptor) export.Aggregator

	aggregation struct {
		newFunc newFunc
		labels  []kv.Key
	}
)

var _ export.Integrator = (*Integrator)(nil)

func ParseYamlData(data []byte) (cfg Config, err error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err = v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return
	}
	if err = v.UnmarshalExact(&cfg); err != nil {
		return
	}
	return
}

func s2k(ss ...string) (rr []kv.Key) {
	for _, s := range ss {
		rr = append(rr, kv.Key(s))
	}
	return
}

func New(cfg Config) (*Integrator, error) {
	policies := map[string]*aggregation{}
	instDefault := map[metric.Kind]*aggregation{}
	views := map[string][]*aggregation{}

	for policy, agg := range cfg.Aggregations {
		if agg.Aggregator == "" {
			return nil, fmt.Errorf("empty aggregation name")
		}
		for _, k := range agg.Labels {
			if k == "" {
				return nil, fmt.Errorf("empty aggregation key")
			}
		}
		var nf newFunc
		switch {
		case strings.EqualFold("sum", agg.Aggregator):
			nf = func(_ *metric.Descriptor) export.Aggregator {
				return sum.New()
			}
		case strings.EqualFold("minmaxsumcount", agg.Aggregator):
			nf = func(desc *metric.Descriptor) export.Aggregator {
				return minmaxsumcount.New(desc)
			}
		case strings.EqualFold("histogram", agg.Aggregator):
			nf = func(desc *metric.Descriptor) export.Aggregator {
				// TODO: boundaries
				return histogram.New(desc, nil)
			}
		case strings.EqualFold("lastvalue", agg.Aggregator):
			nf = func(desc *metric.Descriptor) export.Aggregator {
				return lastvalue.New()
			}
		case strings.EqualFold("sketch", agg.Aggregator):
			nf = func(desc *metric.Descriptor) export.Aggregator {
				// TODO: config
				return ddsketch.New(desc, ddsketch.NewDefaultConfig())
			}
		case strings.EqualFold("array", agg.Aggregator):
			nf = func(desc *metric.Descriptor) export.Aggregator {
				return array.New()
			}
		default:
			return nil, fmt.Errorf("unrecognized aggregator name: %s", agg.Aggregator)
		}

		agg := &aggregation{
			newFunc: nf,
			labels:  s2k(agg.Labels...),
		}

		policies[policy] = agg
	}

	for instKind, policy := range cfg.Defaults.Aggregation {
		agg, ok := policies[policy]

		if !ok {
			return nil, fmt.Errorf("undefined policy: %s", policy)
		}

		var kind metric.Kind
		switch {
		case strings.EqualFold("counter", instKind):
			kind = metric.CounterKind
		case strings.EqualFold("updowncounter", instKind):
			kind = metric.UpDownCounterKind
		case strings.EqualFold("valuerecorder", instKind):
			kind = metric.ValueRecorderKind
		case strings.EqualFold("sumobserver", instKind):
			kind = metric.SumObserverKind
		case strings.EqualFold("updownsumobserver", instKind):
			kind = metric.UpDownSumObserverKind
		case strings.EqualFold("valueobserver", instKind):
			kind = metric.ValueObserverKind
		default:
			return nil, fmt.Errorf("invalid instrument kind: %s", instKind)
		}

		instDefault[kind] = agg
	}

	for instName, list := range cfg.Views {
		// TODO: validate name
		if instName == "" {
			return nil, fmt.Errorf("empty instrument name")
		}

		for _, policy := range list {
			agg, ok := policies[policy]

			if !ok {
				return nil, fmt.Errorf("undefined policy: %s", policy)
			}

			views[instName] = append(views[instName], agg)
		}
	}

	return &Integrator{
		instDefault: instDefault,
		views:       views,
	}, nil
}

func (ci *Integrator) AggregatorFor(desc *metric.Descriptor) export.Aggregator {
	views, ok := ci.views[desc.Name()]
	if !ok {
		return ci.instDefault[desc.MetricKind()].newFunc(desc)
	}

	if len(views) == 1 {
		return views[0].newFunc(desc)
	}

	var aggs []export.Aggregator
	for _, v := range views {
		aggs = append(aggs, v.newFunc(desc))
	}

	return multi.New(aggs...)
}

func (ci *Integrator) Process(ctx context.Context, record export.Record) error {
	return nil
}
