package ratelimit

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestSplitProb(t *testing.T) {
	require.Equal(t, -1, expFromFloat64(0.6))
	require.Equal(t, -2, expFromFloat64(0.4))
	require.Equal(t, 0.5, expToFloat64(-1))
	require.Equal(t, 0.25, expToFloat64(-2))

	for _, tc := range []struct {
		in   float64
		low  int
		frac float64
	}{
		// Probability 0.75 corresponds with choosing S=1 (the
		// "low" probability) 50% of the time and S=0 (the
		// "high" probability) 50% of the time.
		{0.75, 1, 0.5},
		{0.6, 1, 0.8},
		{0.9, 1, 0.2},

		// Powers of 2 exactly
		{1, 0, 1},
		{0.5, 1, 1},
		{0.25, 2, 1},

		// Smaller numbers
		{0.1, 4, 0.4}, // 0.1 == 0.4 * 1/16 + 0.6 * 1/8
	} {
		low, high, frac := splitProb(tc.in)
		require.Equal(t, tc.low, low, "got %v want %v", low, tc.low)
		require.Equal(t, tc.low-1, high, "got %v want %v", high, tc.low-1)
		require.InEpsilon(t, tc.frac, frac, 1e-6, "got %v want %v", frac, tc.frac)
	}
}

type testExporter struct {
	lock  sync.Mutex
	spans []trace.ReadOnlySpan
}

func (t *testExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.spans = append(t.spans, spans...)
	return nil
}

func (t *testExporter) Shutdown(ctx context.Context) error {
	return nil
}

func TestRateLimitBasic(t *testing.T) {
	const (
		testInterval = time.Second
		testRate     = 10
		initialRate  = 100
	)
	var (
		simuTime = time.Now()
		round    = 0
	)
	ctx := context.Background()
	nowfunc := func() time.Time {
		return simuTime.Add(time.Duration(round) * testInterval)
	}

	rs := NewSampler(testRate, WithInterval(testInterval), WithNowFunc(nowfunc))
	te := &testExporter{}
	provider := trace.NewTracerProvider(trace.WithSyncer(te), trace.WithSampler(rs))
	tracer := provider.Tracer("test")

	var created int64

	// Simulate a linear rising request rate.
	// Start at 10x the planned sampling rate
	for ; round < 10; round++ {
		for i := 0; i < ((round + 1) * initialRate); i++ {
			created++
			_, span := tracer.Start(ctx, "span")
			span.End()
		}
		t.Log("round", round, "wrote", ((round + 1) * initialRate), "exporter has", len(te.spans))
	}

	for ; round < 20; round++ {
		for i := 0; i < ((20 - round) * initialRate); i++ {
			created++
			_, span := tracer.Start(ctx, "span")
			span.End()
		}
		t.Log("round", round, "wrote", ((20 - round) * initialRate), "exporter has", len(te.spans))
	}

	// Sum the adjusted counts.
	var estimatedCount int64
	for _, sp := range te.spans {
		thisCnt := int64(1)
		for _, attr := range sp.Attributes() {
			if attr.Key == "sampler.adjusted_count" {
				thisCnt = attr.Value.AsInt64()
				break
			}
		}

		estimatedCount += thisCnt
	}

	// The estimated-count error is less than 6%.
	require.InEpsilon(t, created, estimatedCount, 0.06)

	// We had 100 spans in the first round, unconditionally.
	spanCount := len(te.spans) - 100
	avgRate := spanCount / 19

	// The average-rate error is less than 15%.
	require.InEpsilon(t, testRate, avgRate, 0.15)
}
