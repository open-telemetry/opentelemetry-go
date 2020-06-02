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

// TODO: This code takes a direct dependency on all the exporters and
// all the aggregators.  These should be indirect through the use of a
// registry and factory pattern in separate changes.

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/multi"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		Defaults     `mapstructure:"defaults"`
		Views        `mapstructure:"views"`
		Aggregations `mapstructure:"aggregations"`

		// Exporters `mapstructure:"exporters"`
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

	// Exporters map[string]Exporter

	Integrator struct {
		instDefault [][]*aggregation // always len=1
		views       map[string][]*aggregation
		state
	}

	aggregation struct {
		newFunc newFunc
		labels  []kv.Key
	}

	newFunc func(desc *metric.Descriptor) export.Aggregator

	stateKey struct {
		descriptor *metric.Descriptor
		distinct   label.Distinct
		resource   label.Distinct
	}

	stateValue struct {
		aggregator export.Aggregator
		labels     *label.Set
		resource   *resource.Resource
	}

	state struct {
		// RWMutex implements locking for the `CheckpointSet` interface.
		sync.RWMutex

		persistent map[stateKey]stateValue
		temporary  []export.Record
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

func (ci *Integrator) newInstDefault() [][]*aggregation {
	id := make([][]*aggregation, metric.NumKinds)

	addDef := func(mkind metric.Kind, _ aggregator.Kind, nf newFunc) {
		id[mkind] = []*aggregation{
			&aggregation{
				newFunc: nf,
			},
		}
	}

	addDef(metric.CounterKind, aggregator.SumKind, ci.sumAggregator)
	addDef(metric.UpDownCounterKind, aggregator.SumKind, ci.sumAggregator)
	addDef(metric.SumObserverKind, aggregator.SumKind, ci.sumAggregator)
	addDef(metric.UpDownSumObserverKind, aggregator.SumKind, ci.sumAggregator)
	addDef(metric.ValueRecorderKind, aggregator.MinMaxSumCountKind, ci.minmaxsumcountAggregator)
	addDef(metric.ValueObserverKind, aggregator.MinMaxSumCountKind, ci.minmaxsumcountAggregator)

	return id
}

func New(cfg Config) (*Integrator, error) {
	policies := map[string]*aggregation{}

	ci := &Integrator{
		views: map[string][]*aggregation{},
		state: state{
			persistent: map[stateKey]stateValue{},
		},
	}
	ci.instDefault = ci.newInstDefault()

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
			nf = ci.sumAggregator
		case strings.EqualFold("minmaxsumcount", agg.Aggregator):
			nf = ci.minmaxsumcountAggregator
		case strings.EqualFold("histogram", agg.Aggregator):
			nf = ci.histogramAggregator
		case strings.EqualFold("lastvalue", agg.Aggregator):
			nf = ci.lastvalueAggregator
		case strings.EqualFold("sketch", agg.Aggregator):
			nf = ci.sketchAggregator
		case strings.EqualFold("array", agg.Aggregator):
			nf = ci.arrayAggregator
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

		ci.instDefault[kind] = []*aggregation{agg}
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

			ci.views[instName] = append(ci.views[instName], agg)
		}
	}

	return ci, nil
}

func (ci *Integrator) sumAggregator(_ *metric.Descriptor) export.Aggregator {
	return sum.New()
}

func (ci *Integrator) minmaxsumcountAggregator(desc *metric.Descriptor) export.Aggregator {
	return minmaxsumcount.New(desc)
}

func (ci *Integrator) histogramAggregator(desc *metric.Descriptor) export.Aggregator {
	return histogram.New(desc, nil)
}

func (ci *Integrator) lastvalueAggregator(desc *metric.Descriptor) export.Aggregator {
	return lastvalue.New()
}

func (ci *Integrator) sketchAggregator(desc *metric.Descriptor) export.Aggregator {
	return ddsketch.New(desc, ddsketch.NewDefaultConfig())
}

func (ci *Integrator) arrayAggregator(desc *metric.Descriptor) export.Aggregator {
	return array.New()
}

func (ci *Integrator) aggregationFor(desc *metric.Descriptor) []*aggregation {
	views, ok := ci.views[desc.Name()]
	if !ok {
		return ci.instDefault[desc.MetricKind()]
	}
	return views
}

func (ci *Integrator) AggregatorFor(desc *metric.Descriptor) export.Aggregator {
	views := ci.aggregationFor(desc)

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
	// desc := record.Descriptor()
	// for _, view := range ci.aggregationFor(desc) {
	// 	keys := view.Labels

	// 	// Cache the mapping from Descriptor->Key->Index
	// 	ki, ok := b.descKeyIndex[desc]
	// 	if !ok {
	// 		ki = map[core.Key]int{}
	// 		b.descKeyIndex[desc] = ki

	// 		for i, k := range keys {
	// 			ki[k] = i
	// 		}
	// 	}

	// 	// Compute the value list.  Note: Unspecified values become
	// 	// empty strings.  TODO: pin this down, we have no appropriate
	// 	// Value constructor.
	// 	outputLabels := make([]core.KeyValue, len(keys))

	// 	for i, key := range keys {
	// 		outputLabels[i] = key.String("")
	// 	}

	// 	// Note also the possibility to speed this computation of
	// 	// "encoded" via "outputLabels" in the form of a (Descriptor,
	// 	// Labels)->(Labels, Encoded) cache.
	// 	iter := record.Labels().Iter()
	// 	for iter.Next() {
	// 		kv := iter.Label()
	// 		pos, ok := ki[kv.Key]
	// 		if !ok {
	// 			continue
	// 		}
	// 		outputLabels[pos].Value = kv.Value
	// 	}

	// 	// Compute an encoded lookup key.
	// 	elabels := export.NewSimpleLabels(b.labelEncoder, outputLabels...)
	// 	encoded := elabels.Encoded(b.labelEncoder)

	// 	// Merge this aggregator with all preceding aggregators that
	// 	// map to the same set of `outputLabels` labels.
	// 	agg := record.Aggregator()
	// 	key := batchKey{
	// 		descriptor: record.Descriptor(),
	// 		encoded:    encoded,
	// 	}
	// 	rag, ok := b.aggCheckpoint[key]
	// 	if ok {
	// 		// Combine the input aggregator with the current
	// 		// checkpoint state.
	// 		return rag.Aggregator().Merge(agg, desc)
	// 	}
	// 	// If this Batcher is stateful, create a copy of the
	// 	// Aggregator for long-term storage.  Otherwise the
	// 	// Meter implementation will checkpoint the aggregator
	// 	// again, overwriting the long-lived state.
	// 	if b.stateful {
	// 		tmp := agg
	// 		// Note: the call to AggregatorFor() followed by Merge
	// 		// is effectively a Clone() operation.
	// 		agg = b.AggregatorFor(desc)
	// 		if err := agg.Merge(tmp, desc); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	b.aggCheckpoint[key] = export.NewRecord(desc, elabels, agg)
	// }
	// return nil

	// Some of the former "defaultkeys batcher" goes here.

	// for _, v := range views {
	// 	ci.temporary = append(ci.temporary, record)

	// }

	// @@@ Decide which records are stateful and which are not.
	// Need to ask the exporters what their disposition:
	//   - pass-through
	//   - only delta
	//   - only cumulative
	// It depends on some kind of elective: e.g., GaugeHistogram
	// vs. Histogram Observer instruments are special because
	// inputs are cumulative, need to track start_time (but
	// probably not very special, consider the opposite case of
	// delta instruments reporting through cumulative
	// exporters...).
	//
	// For Prometheus this varies by the instrument.
	//
	// Store stateful records in a map; build a slice of stateless records.

	return nil
}

func (b *state) ForEach(f func(export.Record) error) error {
	// for key, value := range b.values {
	// 	if err := f(export.NewRecord(
	// 		key.descriptor,
	// 		value.labels,
	// 		value.resource,
	// 		value.aggregator,
	// 	)); err != nil && !errors.Is(err, aggregator.ErrNoData) {
	// 		return err
	// 	}
	// }
	return nil
}

// @@@

// func (b *Integrator) CheckpointSet() export.CheckpointSet {
// 	return &b.batch
// }

// func (b *Integrator) FinishedCollection() {
// 	if !b.stateful {
// 		b.batch.values = map[batchKey]batchValue{}
// 	}
// }
