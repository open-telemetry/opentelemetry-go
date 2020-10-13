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

package aggregator_test // import "go.opentelemetry.io/otel/sdk/metric/aggregator"

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

func TestInconsistentAggregatorErr(t *testing.T) {
	err := aggregator.NewInconsistentAggregatorError(&sum.New(1)[0], &lastvalue.New(1)[0])
	require.Equal(
		t,
		"inconsistent aggregator types: *sum.Aggregator and *lastvalue.Aggregator",
		err.Error(),
	)
	require.True(t, errors.Is(err, aggregation.ErrInconsistentType))
}

func testRangeNaN(t *testing.T, desc *otel.Descriptor) {
	// If the descriptor uses int64 numbers, this won't register as NaN
	nan := otel.NewFloat64Number(math.NaN())
	err := aggregator.RangeTest(nan, desc)

	if desc.NumberKind() == otel.Float64NumberKind {
		require.Equal(t, aggregation.ErrNaNInput, err)
	} else {
		require.Nil(t, err)
	}
}

func testRangeNegative(t *testing.T, desc *otel.Descriptor) {
	var neg, pos otel.Number

	if desc.NumberKind() == otel.Float64NumberKind {
		pos = otel.NewFloat64Number(+1)
		neg = otel.NewFloat64Number(-1)
	} else {
		pos = otel.NewInt64Number(+1)
		neg = otel.NewInt64Number(-1)
	}

	posErr := aggregator.RangeTest(pos, desc)
	negErr := aggregator.RangeTest(neg, desc)

	require.Nil(t, posErr)
	require.Equal(t, negErr, aggregation.ErrNegativeInput)
}

func TestRangeTest(t *testing.T) {
	// Only Counters implement a range test.
	for _, nkind := range []otel.NumberKind{otel.Float64NumberKind, otel.Int64NumberKind} {
		t.Run(nkind.String(), func(t *testing.T) {
			desc := otel.NewDescriptor(
				"name",
				otel.CounterInstrumentKind,
				nkind,
			)
			testRangeNegative(t, &desc)
		})
	}
}

func TestNaNTest(t *testing.T) {
	for _, nkind := range []otel.NumberKind{otel.Float64NumberKind, otel.Int64NumberKind} {
		t.Run(nkind.String(), func(t *testing.T) {
			for _, mkind := range []otel.InstrumentKind{
				otel.CounterInstrumentKind,
				otel.ValueRecorderInstrumentKind,
				otel.ValueObserverInstrumentKind,
			} {
				desc := otel.NewDescriptor(
					"name",
					mkind,
					nkind,
				)
				testRangeNaN(t, &desc)
			}
		})
	}
}
