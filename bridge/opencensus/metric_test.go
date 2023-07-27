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

package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	ocmetricdata "go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricproducer"
	ocresource "go.opencensus.io/resource"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
)

type fakeOCProducer struct {
	metrics []*ocmetricdata.Metric
}

func (f *fakeOCProducer) Read() []*ocmetricdata.Metric {
	return f.metrics
}

func TestPushMetricsExporter(t *testing.T) {
	now := time.Now()
	for _, tc := range []struct {
		desc         string
		input        []*ocmetricdata.Metric
		inputMetrics *metricdata.ResourceMetrics
		exportErr    error
		expected     *metricdata.ResourceMetrics
		expectErr    bool
	}{
		{
			desc: "empty batch isn't sent",
		},
		{
			desc:         "export error",
			exportErr:    fmt.Errorf("failed to export"),
			inputMetrics: &metricdata.ResourceMetrics{},
			input: []*ocmetricdata.Metric{
				{
					Resource: &ocresource.Resource{
						Labels: map[string]string{
							"R1": "V1",
							"R2": "V2",
						},
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							StartTime: now,
							Points: []ocmetricdata.Point{
								{Value: int64(123), Time: now},
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			desc: "success",
			input: []*ocmetricdata.Metric{
				{
					Resource: &ocresource.Resource{
						Labels: map[string]string{
							"R1": "V1",
							"R2": "V2",
						},
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							StartTime: now,
							Points: []ocmetricdata.Point{
								{Value: int64(123), Time: now},
							},
						},
					},
				},
			},
			inputMetrics: &metricdata.ResourceMetrics{
				Resource: resource.NewSchemaless(
					attribute.String("R1", "V1"),
					attribute.String("R2", "V2"),
				),
			},
			expected: &metricdata.ResourceMetrics{
				Resource: resource.NewSchemaless(
					attribute.String("R1", "V1"),
					attribute.String("R2", "V2"),
				),
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Scope: instrumentation.Scope{
							Name: scopeName,
						},
						Metrics: []metricdata.Metrics{
							{
								Name:        "",
								Description: "",
								Unit:        "",
								Data: metricdata.Gauge[int64]{
									DataPoints: []metricdata.DataPoint[int64]{
										{
											Attributes: attribute.NewSet(),
											StartTime:  now,
											Time:       now,
											Value:      123,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			fakeProducer := &fakeOCProducer{metrics: tc.input}
			metricproducer.GlobalManager().AddProducer(fakeProducer)
			defer metricproducer.GlobalManager().DeleteProducer(fakeProducer)
			fake := &fakeExporter{err: tc.exportErr}
			exporter := NewMetricExporter(fake)
			err := exporter.Export(context.Background(), tc.inputMetrics)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if tc.expected != nil {
				require.NotNil(t, fake.data)
				metricdatatest.AssertEqual(t, *tc.expected, *fake.data)
			} else {
				require.Nil(t, fake.data)
			}
		})
	}
}

type fakeExporter struct {
	metric.Exporter
	data *metricdata.ResourceMetrics
	err  error
}

func (f *fakeExporter) Export(ctx context.Context, data *metricdata.ResourceMetrics) error {
	if f.err == nil {
		f.data = data
	}
	return f.err
}
