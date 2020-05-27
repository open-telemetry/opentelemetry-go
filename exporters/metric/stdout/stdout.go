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

package stdout // import "go.opentelemetry.io/otel/exporters/metric/stdout"

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/label"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type Exporter struct {
	config Config
}

var _ export.Exporter = &Exporter{}

// Config is the configuration to be used when initializing a stdout export.
type Config struct {
	// Writer is the destination.  If not set, os.Stdout is used.
	Writer io.Writer

	// PrettyPrint will pretty the json representation of the span,
	// making it print "pretty". Default is false.
	PrettyPrint bool

	// DoNotPrintTime suppresses timestamp printing.  This is
	// useful to create deterministic test conditions.
	DoNotPrintTime bool

	// Quantiles are the desired aggregation quantiles for distribution
	// summaries, used when the configured aggregator supports
	// quantiles.
	//
	// Note: this exporter is meant as a demonstration; a real
	// exporter may wish to configure quantiles on a per-metric
	// basis.
	Quantiles []float64

	// LabelEncoder encodes the labels
	LabelEncoder label.Encoder
}

type expoBatch struct {
	Timestamp *time.Time `json:"time,omitempty"`
	Updates   []expoLine `json:"updates"`
}

type expoLine struct {
	Name      string      `json:"name"`
	Min       interface{} `json:"min,omitempty"`
	Max       interface{} `json:"max,omitempty"`
	Sum       interface{} `json:"sum,omitempty"`
	Count     interface{} `json:"count,omitempty"`
	LastValue interface{} `json:"last,omitempty"`

	Quantiles interface{} `json:"quantiles,omitempty"`

	// Note: this is a pointer because omitempty doesn't work when time.IsZero()
	Timestamp *time.Time `json:"time,omitempty"`
}

type expoQuantile struct {
	Q interface{} `json:"q"`
	V interface{} `json:"v"`
}

// NewRawExporter creates a stdout Exporter for use in a pipeline.
func NewRawExporter(config Config) (*Exporter, error) {
	if config.Writer == nil {
		config.Writer = os.Stdout
	}
	if config.Quantiles == nil {
		config.Quantiles = []float64{0.5, 0.9, 0.99}
	} else {
		for _, q := range config.Quantiles {
			if q < 0 || q > 1 {
				return nil, aggregator.ErrInvalidQuantile
			}
		}
	}
	if config.LabelEncoder == nil {
		config.LabelEncoder = label.DefaultEncoder()
	}
	return &Exporter{
		config: config,
	}, nil
}

// InstallNewPipeline instantiates a NewExportPipeline and registers it globally.
// Typically called as:
//
// 	pipeline, err := stdout.InstallNewPipeline(stdout.Config{...})
// 	if err != nil {
// 		...
// 	}
// 	defer pipeline.Stop()
// 	... Done
func InstallNewPipeline(config Config, options ...push.Option) (*push.Controller, error) {
	controller, err := NewExportPipeline(config, options...)
	if err != nil {
		return controller, err
	}
	global.SetMeterProvider(controller.Provider())
	return controller, err
}

// NewExportPipeline sets up a complete export pipeline with the
// recommended setup, chaining a NewRawExporter into the recommended
// selectors and integrators.
//
// The pipeline is configured with a stateful integrator unless the
// push.WithStateful(false) option is used.
func NewExportPipeline(config Config, options ...push.Option) (*push.Controller, error) {
	exporter, err := NewRawExporter(config)
	if err != nil {
		return nil, err
	}
	pusher := push.New(
		simple.NewWithExactDistribution(),
		exporter,
		append([]push.Option{push.WithStateful(true)}, options...)...,
	)
	pusher.Start()

	return pusher, nil
}

func (e *Exporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
	var aggError error
	var batch expoBatch
	if !e.config.DoNotPrintTime {
		ts := time.Now()
		batch.Timestamp = &ts
	}
	aggError = checkpointSet.ForEach(func(record export.Record) error {
		desc := record.Descriptor()
		agg := record.Aggregator()
		kind := desc.NumberKind()
		encodedResource := record.Resource().Encoded(e.config.LabelEncoder)

		var expose expoLine

		if sum, ok := agg.(aggregator.Sum); ok {
			value, err := sum.Sum()
			if err != nil {
				return err
			}
			expose.Sum = value.AsInterface(kind)
		}

		if mmsc, ok := agg.(aggregator.MinMaxSumCount); ok {
			count, err := mmsc.Count()
			if err != nil {
				return err
			}
			expose.Count = count

			max, err := mmsc.Max()
			if err != nil {
				return err
			}
			expose.Max = max.AsInterface(kind)

			min, err := mmsc.Min()
			if err != nil {
				return err
			}
			expose.Min = min.AsInterface(kind)

			if dist, ok := agg.(aggregator.Distribution); ok && len(e.config.Quantiles) != 0 {
				summary := make([]expoQuantile, len(e.config.Quantiles))
				expose.Quantiles = summary

				for i, q := range e.config.Quantiles {
					var vstr interface{}
					value, err := dist.Quantile(q)
					if err != nil {
						return err
					}
					vstr = value.AsInterface(kind)
					summary[i] = expoQuantile{
						Q: q,
						V: vstr,
					}
				}
			}
		} else if lv, ok := agg.(aggregator.LastValue); ok {
			value, timestamp, err := lv.LastValue()
			if err != nil {
				return err
			}
			expose.LastValue = value.AsInterface(kind)

			if !e.config.DoNotPrintTime {
				expose.Timestamp = &timestamp
			}
		}

		var encodedLabels string
		iter := record.Labels().Iter()
		if iter.Len() > 0 {
			encodedLabels = record.Labels().Encoded(e.config.LabelEncoder)
		}

		var sb strings.Builder

		sb.WriteString(desc.Name())

		if len(encodedLabels) > 0 || len(encodedResource) > 0 {
			sb.WriteRune('{')
			sb.WriteString(encodedResource)
			if len(encodedLabels) > 0 && len(encodedResource) > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(encodedLabels)
			sb.WriteRune('}')
		}

		expose.Name = sb.String()

		batch.Updates = append(batch.Updates, expose)
		return nil
	})

	var data []byte
	var err error
	if e.config.PrettyPrint {
		data, err = json.MarshalIndent(batch, "", "\t")
	} else {
		data, err = json.Marshal(batch)
	}

	if err == nil {
		fmt.Fprintln(e.config.Writer, string(data))
	} else {
		return err
	}

	return aggError
}
