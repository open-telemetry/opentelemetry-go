// Copyright 2019, OpenTelemetry Authors
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

package stdout // import "go.opentelemetry.io/otel/exporter/metric/stdout"

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type Exporter struct {
	options Options
}

var _ export.Exporter = &Exporter{}

// Options are the options to be used when initializing a stdout export.
type Options struct {
	// File is the destination.  If not set, os.Stdout is used.
	File io.Writer

	// PrettyPrint will pretty the json representation of the span,
	// making it print "pretty". Default is false.
	PrettyPrint bool

	// DoNotPrintTime suppresses timestamp printing.  This is
	// useful to create deterministic test conditions.
	DoNotPrintTime bool

	// Quantiles are the desired aggregation quantiles for measure
	// metric data, used when the configured aggregator supports
	// quantiles.
	//
	// Note: this exporter is meant as a demonstration; a real
	// exporter may wish to configure quantiles on a per-metric
	// basis.
	Quantiles []float64
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

func New(options Options) (*Exporter, error) {
	if options.File == nil {
		options.File = os.Stdout
	}
	if options.Quantiles == nil {
		options.Quantiles = []float64{0.5, 0.9, 0.99}
	} else {
		for _, q := range options.Quantiles {
			if q < 0 || q > 1 {
				return nil, aggregator.ErrInvalidQuantile
			}
		}
	}
	return &Exporter{
		options: options,
	}, nil
}

func (e *Exporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
	// N.B. Only return one aggError, if any occur. They're likely
	// to be duplicates of the same error.
	var aggError error
	var batch expoBatch
	if !e.options.DoNotPrintTime {
		ts := time.Now()
		batch.Timestamp = &ts
	}
	checkpointSet.ForEach(func(record export.Record) {
		desc := record.Descriptor()
		agg := record.Aggregator()
		kind := desc.NumberKind()

		var expose expoLine

		if sum, ok := agg.(aggregator.Sum); ok {
			if value, err := sum.Sum(); err != nil {
				aggError = err
				expose.Sum = "NaN"
			} else {
				expose.Sum = value.AsInterface(kind)
			}
		}

		if mmsc, ok := agg.(aggregator.MinMaxSumCount); ok {
			if count, err := mmsc.Count(); err != nil {
				aggError = err
				expose.Count = "NaN"
			} else {
				expose.Count = count
			}

			if max, err := mmsc.Max(); err != nil {
				if err == aggregator.ErrEmptyDataSet {
					// This is a special case, indicates an aggregator that
					// was checkpointed before its first value was set.
					return
				}

				aggError = err
				expose.Max = "NaN"
			} else {
				expose.Max = max.AsInterface(kind)
			}

			if min, err := mmsc.Min(); err != nil {
				if err == aggregator.ErrEmptyDataSet {
					// This is a special case, indicates an aggregator that
					// was checkpointed before its first value was set.
					return
				}

				aggError = err
				expose.Min = "NaN"
			} else {
				expose.Min = min.AsInterface(kind)
			}

			if dist, ok := agg.(aggregator.Distribution); ok && len(e.options.Quantiles) != 0 {
				summary := make([]expoQuantile, len(e.options.Quantiles))
				expose.Quantiles = summary

				for i, q := range e.options.Quantiles {
					var vstr interface{}
					if value, err := dist.Quantile(q); err != nil {
						aggError = err
						vstr = "NaN"
					} else {
						vstr = value.AsInterface(kind)
					}
					summary[i] = expoQuantile{
						Q: q,
						V: vstr,
					}
				}
			}

		} else if lv, ok := agg.(aggregator.LastValue); ok {
			if value, timestamp, err := lv.LastValue(); err != nil {
				if err == aggregator.ErrNoLastValue {
					// This is a special case, indicates an aggregator that
					// was checkpointed before its first value was set.
					return
				}

				aggError = err
				expose.LastValue = "NaN"
			} else {
				expose.LastValue = value.AsInterface(kind)

				if !e.options.DoNotPrintTime {
					expose.Timestamp = &timestamp
				}
			}
		}

		var sb strings.Builder

		sb.WriteString(desc.Name())

		if labels := record.Labels(); labels.Len() > 0 {
			sb.WriteRune('{')
			sb.WriteString(labels.Encoded())
			sb.WriteRune('}')
		}

		expose.Name = sb.String()

		batch.Updates = append(batch.Updates, expose)
	})

	var data []byte
	var err error
	if e.options.PrettyPrint {
		data, err = json.MarshalIndent(batch, "", "\t")
	} else {
		data, err = json.Marshal(batch)
	}

	if err == nil {
		fmt.Fprintln(e.options.File, string(data))
	} else {
		return err
	}

	return aggError
}
