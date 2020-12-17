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

package transformjson_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	collectormetricspb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/metrics/v1"
	collectortracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/trace/v1"
	"go.opentelemetry.io/otel/exporters/otlp/internal/otlptest"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform/transformjson"
	metricexport "go.opentelemetry.io/otel/sdk/export/metric"
)

const snapshotJSON = `
{
  "resourceSpans": [
    {
      "instrumentationLibrarySpans": [
        {
          "instrumentationLibrary": {
            "name": "bar",
            "version": "0.0.0"
          },
          "spans": [
            {
              "endTimeUnixNano": "1606767840000000000",
              "kind": "SPAN_KIND_INTERNAL",
              "name": "foo",
              "parentSpanId": "0102030405060708",
              "spanId": "0304050607080900",
              "startTimeUnixNano":"1607458980000000000",
              "status": {
                "code": "STATUS_CODE_OK"
              },
              "traceId": "02030405060708090203040506070809"
            }
          ]
        }
      ],
      "resource": {
        "attributes": [
          {
            "key": "a",
            "value": {
              "stringValue": "b"
            }
          }
        ]
      }
    }
  ]
}
`

const metricsJSON = `
{
  "resourceMetrics": [
    {
      "instrumentationLibraryMetrics": [
        {
          "metrics": [
            {
              "intSum": {
                "aggregationTemporality": "AGGREGATION_TEMPORALITY_CUMULATIVE",
                "dataPoints": [
                  {
                    "labels": [
                      {
                        "key": "abc",
                        "value": "def"
                      },
                      {
                        "key": "one",
                        "value": "1"
                      }
                    ],
                    "startTimeUnixNano": "1607454900000000000",
                    "timeUnixNano": "1607454960000000000",
                    "value": "42"
                  }
                ],
                "isMonotonic": true
              },
              "name": "foo"
            }
          ]
        }
      ],
      "resource": {
        "attributes": [
          {
            "key": "a",
            "value": {
              "stringValue": "b"
            }
          }
        ]
      }
    }
  ]
}
`

func TestMarshalTraces(t *testing.T) {
	snapshotSlice := otlptest.SingleSpanSnapshot()
	protoSpans := transform.SpanData(snapshotSlice)
	request := &collectortracepb.ExportTraceServiceRequest{
		ResourceSpans: protoSpans,
	}
	rawRequest, err := transformjson.Marshal(request)
	assert.NoError(t, err)
	assert.JSONEq(t, snapshotJSON, (string)(rawRequest))
}

func TestMarshalMetrics(t *testing.T) {
	cps := otlptest.OneRecordCheckpointSet{}
	selector := metricexport.CumulativeExportKindSelector()
	rms, err := transform.CheckpointSet(context.Background(), selector, cps, 1)
	assert.NoError(t, err)
	request := &collectormetricspb.ExportMetricsServiceRequest{
		ResourceMetrics: rms,
	}
	rawRequest, err := transformjson.Marshal(request)
	assert.NoError(t, err)
	assert.JSONEq(t, metricsJSON, (string)(rawRequest))
}
