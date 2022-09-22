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

//go:build go1.18
// +build go1.18

package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	ocmetricdata "go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricproducer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func TestProducePartialError(t *testing.T) {
	badProducer := &fakeOCProducer{
		metrics: []*ocmetricdata.Metric{
			{
				Descriptor: ocmetricdata.Descriptor{
					Name:        "foo.com/bad-point",
					Description: "a bad type",
					Unit:        ocmetricdata.UnitDimensionless,
					Type:        ocmetricdata.TypeGaugeDistribution,
				},
			},
		},
	}
	metricproducer.GlobalManager().AddProducer(badProducer)
	defer metricproducer.GlobalManager().DeleteProducer(badProducer)

	end := time.Now()
	goodProducer := &fakeOCProducer{
		metrics: []*ocmetricdata.Metric{
			{
				Descriptor: ocmetricdata.Descriptor{
					Name:        "foo.com/gauge-a",
					Description: "an int testing gauge",
					Unit:        ocmetricdata.UnitBytes,
					Type:        ocmetricdata.TypeGaugeInt64,
				},
				TimeSeries: []*ocmetricdata.TimeSeries{
					{
						Points: []ocmetricdata.Point{
							ocmetricdata.NewInt64Point(end, 123),
						},
					},
				},
			},
		},
	}
	metricproducer.GlobalManager().AddProducer(goodProducer)
	defer metricproducer.GlobalManager().DeleteProducer(goodProducer)

	otelProducer := NewProducer()
	out, err := otelProducer.Produce(context.Background())
	assert.NotNil(t, err)
	expected := []metricdata.Metrics{
		{
			Name:        "foo.com/gauge-a",
			Description: "an int testing gauge",
			Unit:        unit.Bytes,
			Data: metricdata.Gauge[int64]{
				DataPoints: []metricdata.DataPoint[int64]{
					{
						Attributes: attribute.NewSet(),
						Time:       end,
						Value:      123,
					},
				},
			},
		},
	}
	assert.Equal(t, len(out), len(expected))
	for i := range out {
		metricdatatest.AssertEqual[metricdata.Metrics](t, out[i], expected[i])
	}
}

func TestInstrumentationScope(t *testing.T) {
	assert.Equal(t, NewProducer().InstrumentationScope(), instrumentation.Scope{Name: instrumentationName, Version: SemVersion()})
}

type fakeOCProducer struct {
	metrics []*ocmetricdata.Metric
}

func (f *fakeOCProducer) Read() []*ocmetricdata.Metric {
	return f.metrics
}
