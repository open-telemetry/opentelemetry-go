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
	"math"
	"os"
	"strings"
	"time"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
)

type Exporter struct {
	options Options
}

// Options are the options to be used when initializing a stdout export.
type Options struct {
	// File is the destination.  If not set, os.Stdout is used.
	File io.Writer

	// PrettyPrint will pretty the json representation of the span,
	// making it print "pretty". Default is false.
	PrettyPrint bool

	// DoNotPrintTime suppresses timestamp printing.  This is
	// useful to create testable examples or if the are being
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
	Updates   []expoLine `json:"updates,omitempty"`
}

type expoLine struct {
	Name      string      `json:"name"`
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

var _ export.Exporter = &Exporter{}

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

func (e *Exporter) Export(_ context.Context, producer export.Producer) error {
	var batch expoBatch
	var errors []error
	if !e.options.DoNotPrintTime {
		ts := time.Now()
		batch.Timestamp = &ts
	}
	producer.Foreach(func(record export.Record) {
		desc := record.Descriptor()
		labels := record.Labels()
		agg := record.Aggregator()
		kind := desc.NumberKind()

		var expose expoLine
		if msc, ok := agg.(aggregator.MaxSumCount); ok {
			expose.Sum = msc.Sum().AsInterface(kind)
			expose.Count = msc.Count()

			if max, err := msc.Max(); err != nil {
				errors = append(errors, err)
				expose.Max = math.NaN()
			} else {
				expose.Max = max.AsInterface(kind)
			}

			if dist, ok := agg.(aggregator.Distribution); ok && len(e.options.Quantiles) != 0 {
				summary := make([]expoQuantile, len(e.options.Quantiles))
				expose.Quantiles = summary

				for i, q := range e.options.Quantiles {
					var vstr interface{}
					if value, err := dist.Quantile(q); err != nil {
						errors = append(errors, err)
						vstr = math.NaN()
					} else {
						vstr = value.AsInterface(kind)
					}
					summary[i] = expoQuantile{
						Q: q,
						V: vstr,
					}
				}
			}
		} else if sum, ok := agg.(aggregator.Sum); ok {
			expose.Sum = sum.Sum().AsInterface(kind)

		} else if lv, ok := agg.(aggregator.LastValue); ok {
			ts := lv.Timestamp()
			expose.LastValue = lv.LastValue().AsInterface(kind)

			if !e.options.DoNotPrintTime {
				expose.Timestamp = &ts
			}
		}

		var sb strings.Builder

		sb.WriteString(desc.Name())

		if labels.Len() > 0 {
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
		errors = append(errors, err)
	}

	switch len(errors) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("Stdout exporter: %w", errors[0])
	default:
		return fmt.Errorf("Stdout exporter: %v", errors)
	}
}
