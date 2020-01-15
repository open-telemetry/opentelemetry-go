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

package aggregator_test // import "go.opentelemetry.io/otel/sdk/export/metric/aggregator"

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
)

func TestInconsistentMergeErr(t *testing.T) {
	err := aggregator.NewInconsistentMergeError(counter.New(), gauge.New())
	require.Equal(
		t,
		"cannot merge *counter.Aggregator with *gauge.Aggregator: inconsistent aggregator types",
		err.Error(),
	)
	require.True(t, errors.Is(err, aggregator.ErrInconsistentType))
}

func testRangeNaN(t *testing.T, desc *export.Descriptor) {
	// If the descriptor uses int64 numbers, this won't register as NaN
	nan := core.NewFloat64Number(math.NaN())
	err := aggregator.RangeTest(nan, desc)

	if desc.NumberKind() == core.Float64NumberKind {
		require.Equal(t, aggregator.ErrNaNInput, err)
	} else {
		require.Nil(t, err)
	}
}

func testRangeNegative(t *testing.T, alt bool, desc *export.Descriptor) {
	var neg, pos core.Number

	if desc.NumberKind() == core.Float64NumberKind {
		pos = core.NewFloat64Number(+1)
		neg = core.NewFloat64Number(-1)
	} else {
		pos = core.NewInt64Number(+1)
		neg = core.NewInt64Number(-1)
	}

	posErr := aggregator.RangeTest(pos, desc)
	negErr := aggregator.RangeTest(neg, desc)

	require.Nil(t, posErr)

	if desc.MetricKind() == export.GaugeKind {
		require.Nil(t, negErr)
	} else {
		require.Equal(t, negErr == nil, alt)
	}
}

func TestRangeTest(t *testing.T) {
	for _, nkind := range []core.NumberKind{core.Float64NumberKind, core.Int64NumberKind} {
		t.Run(nkind.String(), func(t *testing.T) {
			for _, mkind := range []export.MetricKind{
				export.CounterKind,
				export.GaugeKind,
				export.MeasureKind,
			} {
				t.Run(mkind.String(), func(t *testing.T) {
					for _, alt := range []bool{true, false} {
						t.Run(fmt.Sprint(alt), func(t *testing.T) {
							desc := export.NewDescriptor(
								"name",
								mkind,
								nil,
								"",
								"",
								nkind,
								alt,
							)
							testRangeNaN(t, desc)
							testRangeNegative(t, alt, desc)
						})
					}
				})
			}
		})
	}
}
