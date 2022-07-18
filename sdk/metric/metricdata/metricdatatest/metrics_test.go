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

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestMetricsComparison(t *testing.T) {
	a := metricdata.Metrics{
		Name:        "a",
		Description: "a desc",
		Unit:        unit.Dimensionless,
	}

	b := metricdata.Metrics{
		Name:        "b",
		Description: "b desc",
		Unit:        unit.Bytes,
		Data: metricdata.Gauge{
			DataPoints: []metricdata.DataPoint{
				{
					Attributes: attribute.NewSet(attribute.Bool("b", true)),
					StartTime:  time.Now(),
					Time:       time.Now(),
					Value:      metricdata.Int64(1),
				},
			},
		},
	}

	AssertMetricsEqual(t, a, a)
	AssertMetricsEqual(t, b, b)

	equal, explanation := CompareMetrics(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explanation, 5, "Name, Description, Unit, and Data (Aggregation) do not match")
}
