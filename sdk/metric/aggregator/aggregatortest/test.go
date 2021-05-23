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

package aggregatortest // import "go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"sort"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
)

const Magnitude = 1000

type Profile struct {
	NumberKind number.Kind
	Random     func(sign int) number.Number
}

type NoopAggregator struct{}
type NoopAggregation struct{}

var _ export.Aggregator = NoopAggregator{}
var _ aggregation.Aggregation = NoopAggregation{}

func newProfiles() []Profile {
	rnd := rand.New(rand.NewSource(rand.Int63()))
	return []Profile{
		{
			NumberKind: number.Int64Kind,
			Random: func(sign int) number.Number {
				return number.NewInt64Number(int64(sign) * int64(rnd.Intn(Magnitude+1)))
			},
		},
		{
			NumberKind: number.Float64Kind,
			Random: func(sign int) number.Number {
				return number.NewFloat64Number(float64(sign) * rnd.Float64() * Magnitude)
			},
		},
	}
}

func NewAggregatorTest(mkind metric.InstrumentKind, nkind number.Kind) *metric.Descriptor {
	desc := metric.NewDescriptor("test.name", mkind, nkind)
	return &desc
}

func RunProfiles(t *testing.T, f func(*testing.T, Profile)) {
	for _, profile := range newProfiles() {
		t.Run(profile.NumberKind.String(), func(t *testing.T) {
			f(t, profile)
		})
	}
}

// TestMain ensures local struct alignment prior to running tests.
func TestMain(m *testing.M) {
	fields := []ottest.FieldOffset{
		{
			Name:   "Numbers.numbers",
			Offset: unsafe.Offsetof(Numbers{}.numbers),
		},
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

type Numbers struct {
	// numbers has to be aligned for 64-bit atomic operations.
	numbers []number.Number
	kind    number.Kind
}

func NewNumbers(kind number.Kind) Numbers {
	return Numbers{
		kind: kind,
	}
}

func (n *Numbers) Append(v number.Number) {
	n.numbers = append(n.numbers, v)
}

func (n *Numbers) Sort() {
	sort.Sort(n)
}

func (n *Numbers) Less(i, j int) bool {
	return n.numbers[i].CompareNumber(n.kind, n.numbers[j]) < 0
}

func (n *Numbers) Len() int {
	return len(n.numbers)
}

func (n *Numbers) Swap(i, j int) {
	n.numbers[i], n.numbers[j] = n.numbers[j], n.numbers[i]
}

func (n *Numbers) Sum() number.Number {
	var sum number.Number
	for _, num := range n.numbers {
		sum.AddNumber(n.kind, num)
	}
	return sum
}

func (n *Numbers) Count() uint64 {
	return uint64(len(n.numbers))
}

func (n *Numbers) Min() number.Number {
	return n.numbers[0]
}

func (n *Numbers) Max() number.Number {
	return n.numbers[len(n.numbers)-1]
}

func (n *Numbers) Points() []number.Number {
	return n.numbers
}

// CheckedUpdate performs the same range test the SDK does on behalf of the aggregator.
func CheckedUpdate(t *testing.T, agg export.Aggregator, number number.Number, descriptor *metric.Descriptor) {
	ctx := context.Background()

	// Note: Aggregator tests are written assuming that the SDK
	// has performed the RangeTest. Therefore we skip errors that
	// would have been detected by the RangeTest.
	err := aggregator.RangeTest(number, descriptor)
	if err != nil {
		return
	}

	if err := agg.Update(ctx, number, descriptor); err != nil {
		t.Error("Unexpected Update failure", err)
	}
}

func CheckedMerge(t *testing.T, aggInto, aggFrom export.Aggregator, descriptor *metric.Descriptor) {
	if err := aggInto.Merge(aggFrom, descriptor); err != nil {
		t.Error("Unexpected Merge failure", err)
	}
}

func (NoopAggregation) Kind() aggregation.Kind {
	return aggregation.Kind("Noop")
}

func (NoopAggregator) Aggregation() aggregation.Aggregation {
	return NoopAggregation{}
}

func (NoopAggregator) Update(context.Context, number.Number, *metric.Descriptor) error {
	return nil
}

func (NoopAggregator) SynchronizedMove(export.Aggregator, *metric.Descriptor) error {
	return nil
}

func (NoopAggregator) Merge(export.Aggregator, *metric.Descriptor) error {
	return nil
}

func SynchronizedMoveResetTest(t *testing.T, mkind metric.InstrumentKind, nf func(*metric.Descriptor) export.Aggregator) {
	t.Run("reset on nil", func(t *testing.T) {
		// Ensures that SynchronizedMove(nil, descriptor) discards and
		// resets the aggregator.
		RunProfiles(t, func(t *testing.T, profile Profile) {
			descriptor := NewAggregatorTest(
				mkind,
				profile.NumberKind,
			)
			agg := nf(descriptor)

			for i := 0; i < 10; i++ {
				x1 := profile.Random(+1)
				CheckedUpdate(t, agg, x1, descriptor)
			}

			require.NoError(t, agg.SynchronizedMove(nil, descriptor))

			if count, ok := agg.(aggregation.Count); ok {
				c, err := count.Count()
				require.Equal(t, uint64(0), c)
				require.NoError(t, err)
			}

			if sum, ok := agg.(aggregation.Sum); ok {
				s, err := sum.Sum()
				require.Equal(t, number.Number(0), s)
				require.NoError(t, err)
			}

			if lv, ok := agg.(aggregation.LastValue); ok {
				v, _, err := lv.LastValue()
				require.Equal(t, number.Number(0), v)
				require.Error(t, err)
				require.True(t, errors.Is(err, aggregation.ErrNoData))
			}
		})
	})

	t.Run("no reset on incorrect type", func(t *testing.T) {
		// Ensures that SynchronizedMove(wrong_type, descriptor) does not
		// reset the aggregator.
		RunProfiles(t, func(t *testing.T, profile Profile) {
			descriptor := NewAggregatorTest(
				mkind,
				profile.NumberKind,
			)
			agg := nf(descriptor)

			var input number.Number
			const inval = 100
			if profile.NumberKind == number.Int64Kind {
				input = number.NewInt64Number(inval)
			} else {
				input = number.NewFloat64Number(inval)
			}

			CheckedUpdate(t, agg, input, descriptor)

			err := agg.SynchronizedMove(NoopAggregator{}, descriptor)
			require.Error(t, err)
			require.True(t, errors.Is(err, aggregation.ErrInconsistentType))

			// Test that the aggregator was not reset

			if count, ok := agg.(aggregation.Count); ok {
				c, err := count.Count()
				require.Equal(t, uint64(1), c)
				require.NoError(t, err)
			}

			if sum, ok := agg.(aggregation.Sum); ok {
				s, err := sum.Sum()
				require.Equal(t, input, s)
				require.NoError(t, err)
			}

			if lv, ok := agg.(aggregation.LastValue); ok {
				v, _, err := lv.LastValue()
				require.Equal(t, input, v)
				require.NoError(t, err)
			}

		})
	})

}
