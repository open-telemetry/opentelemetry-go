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

package stdout

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	apimetric "go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

type metricExporter struct {
	config Config
}

var _ metric.Exporter = &metricExporter{}

type line struct {
	Name      string      `json:"Name"`
	Min       interface{} `json:"Min,omitempty"`
	Max       interface{} `json:"Max,omitempty"`
	Sum       interface{} `json:"Sum,omitempty"`
	Count     interface{} `json:"Count,omitempty"`
	LastValue interface{} `json:"Last,omitempty"`

	Quantiles []quantile `json:"Quantiles,omitempty"`

	Histogram histogram `json:"histogram,omitempty"`

	// Note: this is a pointer because omitempty doesn't work when time.IsZero()
	Timestamp *time.Time `json:"Timestamp,omitempty"`
}

type (
	boundary struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	}

	bucket struct {
		boundary `json:"bounds"`
		Count    float64 `json:"count"`
	}

	histogram []bucket

	quantile struct {
		Quantile interface{} `json:"Quantile"`
		Value    interface{} `json:"Value"`
	}
)

func (hist histogram) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(hist)
	if err != nil {
		return nil, err
	}
	return append(bytes, []byte("\n"+printHistogram(hist))...), nil
}

func (e *metricExporter) ExportKindFor(*apimetric.Descriptor, aggregation.Kind) metric.ExportKind {
	return metric.PassThroughExporter
}

func (e *metricExporter) Export(_ context.Context, checkpointSet metric.CheckpointSet) error {
	if e.config.DisableMetricExport {
		return nil
	}
	var aggError error
	var batch []line
	aggError = checkpointSet.ForEach(e, func(record metric.Record) error {
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

		var expose line

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
				summary := make([]quantile, len(e.config.Quantiles))
				expose.Quantiles = summary

				for i, q := range e.config.Quantiles {
					value, err := dist.Quantile(q)
					if err != nil {
						return err
					}
					summary[i] = quantile{
						Quantile: q,
						Value:    value.AsInterface(kind),
					}
				}
			}

			if hist, ok := agg.(aggregation.Histogram); ok {
				val, err := hist.Histogram()
				if err != nil {
					return err
				}
				buckets := make([]bucket, 0, len(val.Boundaries))
				lowerBound := float64(-math.MaxFloat64)
				upperBound := val.Boundaries[0]
				for i := 0; i < len(val.Boundaries); i++ {
					buckets = append(buckets, bucket{
						boundary: boundary{Min: lowerBound, Max: upperBound},
						Count:    val.Counts[i],
					})
					lowerBound = val.Boundaries[i]
					if i == len(val.Boundaries)-1 {
						upperBound = float64(math.MaxFloat64)
					} else {
						upperBound = val.Boundaries[i+1]
					}
				}
				expose.Histogram = buckets
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

		batch = append(batch, expose)
		return nil
	})
	if len(batch) == 0 {
		return aggError
	}

	data, err := e.marshal(batch)
	if err != nil {
		return err
	}
	fmt.Fprintln(e.config.Writer, string(data))

	return aggError
}

// marshal v with approriate indentation.
func (e *metricExporter) marshal(v interface{}) ([]byte, error) {
	if e.config.PrettyPrint {
		return json.MarshalIndent(v, "", "\t")
	}
	return json.Marshal(v)
}

// printHistogram prints histogram in horizontal tabular form
func printHistogram(hist histogram) string {
	if len(hist) == 0 {
		return ""
	}
	var totalCount float64
	for _, buck := range hist {
		totalCount += buck.Count
	}
	str := fmt.Sprintf("Total Count: %f\n", totalCount)

	for _, buck := range hist {
		str += fmt.Sprintf("[%e\t-\t%e)\t\t", buck.Min, buck.Max)
		for i := 0; i < int((buck.Count/totalCount)*100); i++ {
			str += "*"
		}
		str += "\n"
	}
	return str
}
