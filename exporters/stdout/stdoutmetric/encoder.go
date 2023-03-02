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
	"encoding/json"
	"errors"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Encoder encodes and outputs OpenTelemetry metric data-types as human
// readable text.
type Encoder interface {
	// Encode handles the encoding and writing of OpenTelemetry metric data.
	Encode(v any) error
}

// encoderHolder is the concrete type used to wrap an Encoder so it can be
// used as a atomic.Value type.
type encoderHolder struct {
	encoder Encoder
}

func (e encoderHolder) Encode(v any) error { return e.encoder.Encode(v) }

// shutdownEncoder is used when the exporter is shutdown. It always returns
// errShutdown when Encode is called.
type shutdownEncoder struct{}

var errShutdown = errors.New("exporter shutdown")

func (shutdownEncoder) Encode(any) error { return errShutdown }

type encoderIgnoreTimestamp struct {
	encoder Encoder
}

// NewEncoderIgnoreTimestamp return encoderIgnoreTimestamp which wrap a Encoder,
// It redact timestamp to a zero value before Encode .
func NewEncoderIgnoreTimestamp(encoder *json.Encoder) Encoder {
	return &encoderIgnoreTimestamp{
		encoder: encoder,
	}
}

// Encode redact timestamp to a zero value before Encode .
func (e *encoderIgnoreTimestamp) Encode(v any) error {
	rm := v.(metricdata.ResourceMetrics)
	for i, sm := range rm.ScopeMetrics {
		for j, m := range sm.Metrics {
			switch v := m.Data.(type) {
			case metricdata.Sum[float64]:
				for k := range v.DataPoints {
					v.DataPoints[k].StartTime = time.Time{}
					v.DataPoints[k].Time = time.Time{}
				}
				rm.ScopeMetrics[i].Metrics[j].Data = v
			case metricdata.Sum[int64]:
				for k := range v.DataPoints {
					v.DataPoints[k].StartTime = time.Time{}
					v.DataPoints[k].Time = time.Time{}
				}
				rm.ScopeMetrics[i].Metrics[j].Data = v
			case metricdata.Gauge[float64]:
				for k := range v.DataPoints {
					v.DataPoints[k].StartTime = time.Time{}
					v.DataPoints[k].Time = time.Time{}
				}
				rm.ScopeMetrics[i].Metrics[j].Data = v
			case metricdata.Gauge[int64]:
				for k := range v.DataPoints {
					v.DataPoints[k].StartTime = time.Time{}
					v.DataPoints[k].Time = time.Time{}
				}
				rm.ScopeMetrics[i].Metrics[j].Data = v
			case metricdata.Histogram:
				for k := range v.DataPoints {
					v.DataPoints[k].StartTime = time.Time{}
					v.DataPoints[k].Time = time.Time{}
				}
				rm.ScopeMetrics[i].Metrics[j].Data = v
			default:
				panic("invalid Aggregation")
			}
		}
	}

	return e.encoder.Encode(v)
}
