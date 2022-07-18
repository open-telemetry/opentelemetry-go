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

package exporttest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestSumsComparison(t *testing.T) {
	a := export.Sum{
		Temporality: export.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints: []export.DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("a", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      export.Int64(2),
			},
		},
	}

	b := export.Sum{
		Temporality: export.DeltaTemporality,
		IsMonotonic: false,
		DataPoints: []export.DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("b", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      export.Int64(1),
			},
		},
	}

	AssertSumsEqual(t, a, a)
	AssertSumsEqual(t, b, b)

	equal, explination := CompareSum(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 3, "Temporality, IsMonotonic, and DataPoints do not match")
}
