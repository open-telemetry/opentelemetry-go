// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func BenchmarkInstrument(b *testing.B) {
	attr := func(id int) attribute.Set {
		return attribute.NewSet(
			attribute.String("user", "Alice"),
			attribute.Bool("admin", true),
			attribute.Int("id", id),
		)
	}

	b.Run("instrumentImpl/aggregate", func(b *testing.B) {
		build := aggregate.Builder[int64]{}
		var meas []aggregate.Measure[int64]

		build.Temporality = metricdata.CumulativeTemporality
		in, _ := build.LastValue()
		meas = append(meas, in)

		build.Temporality = metricdata.DeltaTemporality
		in, _ = build.LastValue()
		meas = append(meas, in)

		build.Temporality = metricdata.CumulativeTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		build.Temporality = metricdata.DeltaTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		inst := int64Inst{measures: meas}
		ctx := b.Context()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			inst.aggregate(ctx, int64(i), attr(i))
		}
	})

	b.Run("observable/observe", func(b *testing.B) {
		build := aggregate.Builder[int64]{}
		var meas []aggregate.Measure[int64]

		in, _ := build.PrecomputedLastValue()
		meas = append(meas, in)

		build.Temporality = metricdata.CumulativeTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		build.Temporality = metricdata.DeltaTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		o := observable[int64]{measures: meas}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			o.observe(int64(i), attr(i))
		}
	})
}

func TestExtractRawKVs(t *testing.T) {
	k1 := attribute.String("k1", "v1")
	k2 := attribute.String("k2", "v2")
	k3 := attribute.String("k3", "v3")

	tests := []struct {
		name string
		opts []metric.AddOption
		want []attribute.KeyValue
	}{
		{
			name: "Empty",
			opts: nil,
			want: nil,
		},
		{
			name: "NoRawAttributes",
			opts: []metric.AddOption{metric.WithAttributes(k1)},
			want: nil,
		},
		{
			name: "OneRawAttributes",
			opts: []metric.AddOption{x.WithUnsafeAttributes(k1, k2)},
			want: []attribute.KeyValue{k1, k2},
		},
		{
			name: "MultipleRawAttributes",
			opts: []metric.AddOption{
				x.WithUnsafeAttributes(k1),
				x.WithUnsafeAttributes(k2, k3),
			},
			want: []attribute.KeyValue{k1, k2, k3},
		},
		{
			name: "Mixed",
			opts: []metric.AddOption{
				x.WithUnsafeAttributes(k1),
				metric.WithAttributes(k2),
				x.WithUnsafeAttributes(k3),
			},
			want: []attribute.KeyValue{k1, k3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRawKVs(tt.opts)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestExtractRawKVs_NoModifyOriginalSlice(t *testing.T) {
	k1 := attribute.String("k1", "v1")
	k2 := attribute.String("k2", "v2")

	s1 := []attribute.KeyValue{k1}
	// Give it some extra capacity to ensure append would overwrite if it reused the backing array.
	s1 = append(s1, attribute.KeyValue{})[:1]

	opt1 := x.WithUnsafeAttributes(s1...)
	opt2 := x.WithUnsafeAttributes(k2)

	opts := []metric.AddOption{opt1, opt2}

	got := extractRawKVs(opts)
	require.Equal(t, []attribute.KeyValue{k1, k2}, got)

	// Check that s1 was not modified.
	require.Equal(t, []attribute.KeyValue{k1}, s1[:1])
	fullS1 := s1[:cap(s1)]
	if len(fullS1) > 1 {
		require.NotEqual(t, k2, fullS1[1], "Original slice was modified by extractRawKVs")
	}
}

func TestResolveAttributes(t *testing.T) {
	k1 := attribute.String("k1", "v1")
	k2 := attribute.String("k2", "v2")
	k3 := attribute.String("k3", "v3")
	k1Alt := attribute.String("k1", "v1_alt")

	tests := []struct {
		name        string
		configAttrs attribute.Set
		rawKVs      []attribute.KeyValue
		want        attribute.Set
	}{
		{
			name:        "Empty",
			configAttrs: *attribute.EmptySet(),
			rawKVs:      nil,
			want:        *attribute.EmptySet(),
		},
		{
			name:        "OnlyConfig",
			configAttrs: attribute.NewSet(k1, k2),
			rawKVs:      nil,
			want:        attribute.NewSet(k1, k2),
		},
		{
			name:        "OnlyRaw",
			configAttrs: *attribute.EmptySet(),
			rawKVs:      []attribute.KeyValue{k1, k2},
			want:        attribute.NewSet(k1, k2),
		},
		{
			name:        "MergeNoOverlap",
			configAttrs: attribute.NewSet(k1),
			rawKVs:      []attribute.KeyValue{k2, k3},
			want:        attribute.NewSet(k1, k2, k3),
		},
		{
			name:        "MergeWithOverlap_RawOverrides",
			configAttrs: attribute.NewSet(k1, k2),
			rawKVs:      []attribute.KeyValue{k1Alt, k3},
			want:        attribute.NewSet(k1Alt, k2, k3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveAttributes(tt.configAttrs, tt.rawKVs)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestMeterWithUnsafeAttributes(t *testing.T) {
	k1 := attribute.Key("k1")
	k2 := attribute.Key("k2")

	combined := attribute.NewSet(k1.String("alice"), k2.String("bob"))

	testCases := []struct {
		name     string
		instName string
		record   func(t *testing.T, m metric.Meter)
		wantData func(attrs attribute.Set) metricdata.Aggregation
	}{
		{
			name:     "Int64Counter",
			instName: "sint",
			record: func(t *testing.T, m metric.Meter) {
				ctr, err := m.Int64Counter("sint")
				require.NoError(t, err)
				ctr.Add(t.Context(), 3, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints:  []metricdata.DataPoint[int64]{{Attributes: attrs, Value: 3}},
				}
			},
		},
		{
			name:     "Float64Counter",
			instName: "sfloat",
			record: func(t *testing.T, m metric.Meter) {
				ctr, err := m.Float64Counter("sfloat")
				require.NoError(t, err)
				ctr.Add(t.Context(), 3.5, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints:  []metricdata.DataPoint[float64]{{Attributes: attrs, Value: 3.5}},
				}
			},
		},
		{
			name:     "Int64ObservableCounter",
			instName: "aint",
			record: func(t *testing.T, m metric.Meter) {
				_, err := m.Int64ObservableCounter("aint",
					metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
						o.Observe(4, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
						return nil
					}),
				)
				require.NoError(t, err)
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints:  []metricdata.DataPoint[int64]{{Attributes: attrs, Value: 4}},
				}
			},
		},
		{
			name:     "Float64ObservableCounter",
			instName: "afloat",
			record: func(t *testing.T, m metric.Meter) {
				_, err := m.Float64ObservableCounter("afloat",
					metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
						o.Observe(4.5, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
						return nil
					}),
				)
				require.NoError(t, err)
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints:  []metricdata.DataPoint[float64]{{Attributes: attrs, Value: 4.5}},
				}
			},
		},
		{
			name:     "Int64UpDownCounter",
			instName: "sudint",
			record: func(t *testing.T, m metric.Meter) {
				ctr, err := m.Int64UpDownCounter("sudint")
				require.NoError(t, err)
				ctr.Add(t.Context(), -3, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints:  []metricdata.DataPoint[int64]{{Attributes: attrs, Value: -3}},
				}
			},
		},
		{
			name:     "Float64UpDownCounter",
			instName: "sudfloat",
			record: func(t *testing.T, m metric.Meter) {
				ctr, err := m.Float64UpDownCounter("sudfloat")
				require.NoError(t, err)
				ctr.Add(t.Context(), -3.5, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints:  []metricdata.DataPoint[float64]{{Attributes: attrs, Value: -3.5}},
				}
			},
		},
		{
			name:     "Int64Histogram",
			instName: "shist",
			record: func(t *testing.T, m metric.Meter) {
				hist, err := m.Int64Histogram("shist")
				require.NoError(t, err)
				hist.Record(t.Context(), 5, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Histogram[int64]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[int64]{
						{
							Attributes: attrs,
							Count:      1,
							Sum:        5,
							Bounds: []float64{
								0,
								5,
								10,
								25,
								50,
								75,
								100,
								250,
								500,
								750,
								1000,
								2500,
								5000,
								7500,
								10000,
							},
							BucketCounts: []uint64{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Min:          metricdata.NewExtrema[int64](5),
							Max:          metricdata.NewExtrema[int64](5),
						},
					},
				}
			},
		},
		{
			name:     "Float64Histogram",
			instName: "sfhist",
			record: func(t *testing.T, m metric.Meter) {
				hist, err := m.Float64Histogram("sfhist")
				require.NoError(t, err)
				hist.Record(t.Context(), 5.5, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Histogram[float64]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[float64]{
						{
							Attributes: attrs,
							Count:      1,
							Sum:        5.5,
							Bounds: []float64{
								0,
								5,
								10,
								25,
								50,
								75,
								100,
								250,
								500,
								750,
								1000,
								2500,
								5000,
								7500,
								10000,
							},
							BucketCounts: []uint64{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Min:          metricdata.NewExtrema[float64](5.5),
							Max:          metricdata.NewExtrema[float64](5.5),
						},
					},
				}
			},
		},
		{
			name:     "Int64ObservableUpDownCounter",
			instName: "audint",
			record: func(t *testing.T, m metric.Meter) {
				_, err := m.Int64ObservableUpDownCounter("audint",
					metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
						o.Observe(-4, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
						return nil
					}),
				)
				require.NoError(t, err)
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints:  []metricdata.DataPoint[int64]{{Attributes: attrs, Value: -4}},
				}
			},
		},
		{
			name:     "Float64ObservableUpDownCounter",
			instName: "audfloat",
			record: func(t *testing.T, m metric.Meter) {
				_, err := m.Float64ObservableUpDownCounter("audfloat",
					metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
						o.Observe(-4.5, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
						return nil
					}),
				)
				require.NoError(t, err)
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints:  []metricdata.DataPoint[float64]{{Attributes: attrs, Value: -4.5}},
				}
			},
		},
		{
			name:     "Int64ObservableGauge",
			instName: "agint",
			record: func(t *testing.T, m metric.Meter) {
				_, err := m.Int64ObservableGauge("agint",
					metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
						o.Observe(10, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
						return nil
					}),
				)
				require.NoError(t, err)
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Gauge[int64]{
					DataPoints: []metricdata.DataPoint[int64]{{Attributes: attrs, Value: 10}},
				}
			},
		},
		{
			name:     "Float64ObservableGauge",
			instName: "agfloat",
			record: func(t *testing.T, m metric.Meter) {
				_, err := m.Float64ObservableGauge("agfloat",
					metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
						o.Observe(10.5, x.WithUnsafeAttributes(k1.String("alice"), k2.String("bob")))
						return nil
					}),
				)
				require.NoError(t, err)
			},
			wantData: func(attrs attribute.Set) metricdata.Aggregation {
				return metricdata.Gauge[float64]{
					DataPoints: []metricdata.DataPoint[float64]{{Attributes: attrs, Value: 10.5}},
				}
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rdr := NewManualReader()
			m := NewMeterProvider(WithReader(rdr)).Meter("test")
			tt.record(t, m)

			rm := metricdata.ResourceMetrics{}
			err := rdr.Collect(t.Context(), &rm)
			require.NoError(t, err)

			require.Len(t, rm.ScopeMetrics, 1)
			sm := rm.ScopeMetrics[0]
			require.Len(t, sm.Metrics, 1)
			got := sm.Metrics[0]

			want := metricdata.Metrics{
				Name: tt.instName,
				Data: tt.wantData(combined),
			}
			metricdatatest.AssertEqual(t, want, got, metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreExemplars())
		})
	}
}
