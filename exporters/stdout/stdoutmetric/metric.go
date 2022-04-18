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
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/resource"
)

type metricExporter struct {
	config config
}

var _ export.Exporter = &metricExporter{}

type line struct {
	Name      string      `json:"Name"`
	Sum       interface{} `json:"Sum,omitempty"`
	Count     interface{} `json:"Count,omitempty"`
	LastValue interface{} `json:"Last,omitempty"`

	// Note: this is a pointer because omitempty doesn't work when time.IsZero()
	Timestamp *time.Time `json:"Timestamp,omitempty"`
}

func (e *metricExporter) TemporalityFor(desc *sdkapi.Descriptor, kind aggregation.Kind) aggregation.Temporality {
	return aggregation.StatelessTemporalitySelector().TemporalityFor(desc, kind)
}

func (e *metricExporter) Export(_ context.Context, res *resource.Resource, reader export.InstrumentationLibraryReader) error {
	var aggError error
	var batch []line
	aggError = reader.ForEach(func(lib instrumentation.Library, mr export.Reader) error {

		var instAttrs []attribute.KeyValue
		if name := lib.Name; name != "" {
			instAttrs = append(instAttrs, attribute.String("instrumentation.name", name))
			if version := lib.Version; version != "" {
				instAttrs = append(instAttrs, attribute.String("instrumentation.version", version))
			}
			if schema := lib.SchemaURL; schema != "" {
				instAttrs = append(instAttrs, attribute.String("instrumentation.schema_url", schema))
			}
		}
		instSet := attribute.NewSet(instAttrs...)
		encodedInstAttrs := instSet.Encoded(e.config.Encoder)

		return mr.ForEach(e, func(record export.Record) error {
			desc := record.Descriptor()
			agg := record.Aggregation()
			kind := desc.NumberKind()
			encodedResource := res.Encoded(e.config.Encoder)

			var expose line

			if sum, ok := agg.(aggregation.Sum); ok {
				value, err := sum.Sum()
				if err != nil {
					return err
				}
				expose.Sum = value.AsInterface(kind)
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

			var encodedAttrs string
			iter := record.Attributes().Iter()
			if iter.Len() > 0 {
				encodedAttrs = record.Attributes().Encoded(e.config.Encoder)
			}

			var sb strings.Builder

			sb.WriteString(desc.Name())

			if len(encodedAttrs) > 0 || len(encodedResource) > 0 || len(encodedInstAttrs) > 0 {
				sb.WriteRune('{')
				sb.WriteString(encodedResource)
				if len(encodedInstAttrs) > 0 && len(encodedResource) > 0 {
					sb.WriteRune(',')
				}
				sb.WriteString(encodedInstAttrs)
				if len(encodedAttrs) > 0 && (len(encodedInstAttrs) > 0 || len(encodedResource) > 0) {
					sb.WriteRune(',')
				}
				sb.WriteString(encodedAttrs)
				sb.WriteRune('}')
			}

			expose.Name = sb.String()

			batch = append(batch, expose)
			return nil
		})
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

// marshal v with appropriate indentation.
func (e *metricExporter) marshal(v interface{}) ([]byte, error) {
	if e.config.PrettyPrint {
		return json.MarshalIndent(v, "", "\t")
	}
	return json.Marshal(v)
}
