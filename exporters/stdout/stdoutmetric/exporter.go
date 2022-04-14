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

package stdoutmetric // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/reader"
)

type Exporter struct {
	config config
}

var _ reader.Exporter = &Exporter{}

// New creates an Exporter with the passed options.
func New(options ...Option) *Exporter {
	cfg := newConfig(options...)

	return &Exporter{
		config: cfg,
	}
}

type line struct {
	Name       string                   `json:"Name"`
	Attributes string                   `json:"Attributes,omitempty"`
	Resource   string                   `json:"Resource,omitempty"`
	Library    *instrumentation.Library `json:"InstrumentScope,omitempty"`
	Sum        interface{}              `json:"Sum,omitempty"`
	Count      interface{}              `json:"Count,omitempty"`
	Gauge      interface{}              `json:"Gauge,omitempty"`

	// Note: this is a pointer because omitempty doesn't work when time.IsZero()
	Timestamp *time.Time `json:"Timestamp,omitempty"`
}

func (e *Exporter) Export(_ context.Context, metrics reader.Metrics) error {
	var batch []line

	resource := metrics.Resource.String()
	for _, scope := range metrics.Scopes {

		for _, inst := range scope.Instruments {

			for _, point := range inst.Points {
				expose := line{
					Name:       inst.Descriptor.Name,
					Resource:   resource,
					Attributes: point.Attributes.Encoded(attribute.DefaultEncoder()),
				}
				switch agg := point.Aggregation.(type) {
				case aggregation.Histogram:
					expose.Sum = agg.Sum()
					expose.Count = agg.Count()
				case aggregation.Sum:
					expose.Sum = agg.Sum()
				case aggregation.Gauge:
					expose.Gauge = agg.Gauge()
				}
				batch = append(batch, expose)
			}
		}
	}
	buf, err := e.marshal(batch)
	fmt.Fprintln(e.config.Writer, string(buf))
	return err
}

func (e *Exporter) Shutdown(_ context.Context) error {
	return nil
}

func (e *Exporter) Flush(context.Context) error {
	return nil
}

// marshal v with appropriate indentation.
func (e *Exporter) marshal(v interface{}) ([]byte, error) {
	if e.config.PrettyPrint {
		return json.MarshalIndent(v, "", "\t")
	}
	return json.Marshal(v)
}
