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
	File *os.File

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

	// Note: this is a pointer because omitempty doesn't work when time.IsZero()
	Timestamp *time.Time `json:"time,omitempty"`
}

var _ export.Exporter = &Exporter{}

func New(options Options) *Exporter {
	if options.File == nil {
		options.File = os.Stdout
	}
	return &Exporter{
		options: options,
	}
}

func (e *Exporter) Export(_ context.Context, producer export.Producer) {
	var batch expoBatch
	if !e.options.DoNotPrintTime {
		ts := time.Now()
		batch.Timestamp = &ts
	}
	producer.Foreach(func(record export.Record) {
		desc := record.Descriptor()
		labels := record.Labels()
		agg := record.Aggregator()

		var expose expoLine
		if sum, ok := agg.(aggregator.Sum); ok {
			expose.Sum = sum.Sum().Emit(desc.NumberKind())

		} else if lv, ok := agg.(aggregator.LastValue); ok {
			ts := lv.Timestamp()
			expose.LastValue = lv.LastValue().Emit(desc.NumberKind())
			expose.Timestamp = &ts

		} else if msc, ok := agg.(aggregator.MaxSumCount); ok {
			expose.Max = msc.Max().Emit(desc.NumberKind())
			expose.Sum = msc.Sum().Emit(desc.NumberKind())
			expose.Count = msc.Count().Emit(desc.NumberKind())

		} else if dist, ok := agg.(aggregator.Distribution); ok {
			expose.Max = dist.Max().Emit(desc.NumberKind())
			expose.Sum = dist.Sum().Emit(desc.NumberKind())
			expose.Count = dist.Count().Emit(desc.NumberKind())

			// TODO print one configured quantile per line
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

	if err != nil {
		fmt.Fprintf(e.options.File, "JSON encode error: %v\n", err)
		return
	}

	fmt.Fprintln(e.options.File, string(data))
}
