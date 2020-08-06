package metric

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
)

func ResolveNumberByKind(t *testing.T, kind metric.NumberKind, value float64) metric.Number {
	t.Helper()
	switch kind {
	case metric.Int64NumberKind:
		return metric.NewInt64Number(int64(value))
	case metric.Float64NumberKind:
		return metric.NewFloat64Number(value)
	}
	panic("invalid number kind")
}

// TODO: ADD doc
func CheckSyncBatches(t *testing.T, ctx context.Context, labels []kv.KeyValue, mock *MeterImpl, nkind metric.NumberKind, mkind metric.Kind, instrument metric.InstrumentImpl, expected ...float64) {
	t.Helper()
	if len(mock.MeasurementBatches) != 3 {
		t.Errorf("Expected 3 recorded measurement batches, got %d", len(mock.MeasurementBatches))
	}
	ourInstrument := instrument.Implementation().(*Sync)
	for i, got := range mock.MeasurementBatches {
		if got.Ctx != ctx {
			d := func(c context.Context) string {
				return fmt.Sprintf("(ptr: %p, ctx %#v)", c, c)
			}
			t.Errorf("Wrong recorded context in batch %d, expected %s, got %s", i, d(ctx), d(got.Ctx))
		}
		if !assert.Equal(t, got.Labels, labels) {
			t.Errorf("Wrong recorded label set in batch %d, expected %v, got %v", i, labels, got.Labels)
		}
		if len(got.Measurements) != 1 {
			t.Errorf("Expected 1 measurement in batch %d, got %d", i, len(got.Measurements))
		}
		minMLen := 1
		if minMLen > len(got.Measurements) {
			minMLen = len(got.Measurements)
		}
		for j := 0; j < minMLen; j++ {
			measurement := got.Measurements[j]
			require.Equal(t, mkind, measurement.Instrument.Descriptor().MetricKind())

			if measurement.Instrument.Implementation() != ourInstrument {
				d := func(iface interface{}) string {
					i := iface.(*Instrument)
					return fmt.Sprintf("(ptr: %p, instrument %#v)", i, i)
				}
				t.Errorf("Wrong recorded instrument in measurement %d in batch %d, expected %s, got %s", j, i, d(ourInstrument), d(measurement.Instrument.Implementation()))
			}
			expect := ResolveNumberByKind(t, nkind, expected[i])
			if measurement.Number.CompareNumber(nkind, expect) != 0 {
				t.Errorf("Wrong recorded value in measurement %d in batch %d, expected %s, got %s", j, i, expect.Emit(nkind), measurement.Number.Emit(nkind))
			}
		}
	}
}
