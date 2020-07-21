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
	"strings"
	"time"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

type metricExporter struct {
	config Config
}

var _ export.Exporter = &metricExporter{}

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

func (e *metricExporter) ExportKindFor(*metric.Descriptor, aggregation.Kind) export.ExportKind {
	return export.PassThroughExporter
}

func (e *metricExporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
	var aggError error
	var batch expoBatch
	if e.config.Timestamps {
		ts := time.Now()
		batch.Timestamp = &ts
	}
	aggError = checkpointSet.ForEach(e, func(record export.Record) error {
		desc := record.Descriptor()
		agg := record.Aggregation()
		kind := desc.NumberKind()
		encodedResource := record.Resource().Encoded(e.config.LabelEncoder)

		var instLabels []kv.KeyValue
		if name := desc.InstrumentationName(); name != "" {
			instLabels = append(instLabels, kv.String("instrumentation.name", name))
			if version := desc.InstrumentationVersion(); version != "" {
				instLabels = append(instLabels, kv.String("instrumentation.version", version))
			}
		}
		instSet := label.NewSet(instLabels...)
		encodedInstLabels := instSet.Encoded(e.config.LabelEncoder)

		var expose expoLine

		if sum, ok := agg.(aggregation.Sum); ok {
			value, err := sum.Sum()
			if err != nil {
				return err
			}
			expose.Sum = value.AsInterface(kind)
		}

		if mmsc, ok := agg.(aggregation.MinMaxSumCount); ok {
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

			if dist, ok := agg.(aggregation.Distribution); ok && len(e.config.Quantiles) != 0 {
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
		} else if lv, ok := agg.(aggregation.LastValue); ok {
			value, timestamp, err := lv.LastValue()
			if err != nil {
				return err
			}
			expose.LastValue = value.AsInterface(kind)

			if e.config.Timestamps {
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

		if len(encodedLabels) > 0 || len(encodedResource) > 0 || len(encodedInstLabels) > 0 {
			sb.WriteRune('{')
			sb.WriteString(encodedResource)
			if len(encodedInstLabels) > 0 && len(encodedResource) > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(encodedInstLabels)
			if len(encodedLabels) > 0 && (len(encodedInstLabels) > 0 || len(encodedResource) > 0) {
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
