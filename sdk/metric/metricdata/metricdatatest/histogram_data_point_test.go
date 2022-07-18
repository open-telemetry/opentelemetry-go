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
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestHistogramDataPointsComparison(t *testing.T) {
	a := metricdata.HistogramDataPoint{
		Attributes:   attribute.NewSet(attribute.Bool("a", true)),
		StartTime:    time.Now(),
		Time:         time.Now(),
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Sum:          2,
	}

	max, min := 99.0, 3.
	b := metricdata.HistogramDataPoint{
		Attributes:   attribute.NewSet(attribute.Bool("b", true)),
		StartTime:    time.Now(),
		Time:         time.Now(),
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          &max,
		Min:          &min,
		Sum:          3,
	}

	AssertHistogramDataPointsEqual(t, a, a)
	AssertHistogramDataPointsEqual(t, b, b)

	equal, explination := CompareHistogramDataPoint(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 9, "Attributes, StartTime, Time, Count, Bounds, BucketCounts, Max, Min, and Sum do not match")
}
