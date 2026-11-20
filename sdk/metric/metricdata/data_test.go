// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metricdata_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestExtremaValue(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		testCases := []struct {
			name        string
			extrema     metricdata.Extrema[int64]
			wantVal     int64
			wantDefined bool
		}{
			{
				name:        "zero-value",
				extrema:     metricdata.Extrema[int64]{},
				wantVal:     0,
				wantDefined: false,
			},
			{
				name:        "positive",
				extrema:     metricdata.NewExtrema[int64](42),
				wantVal:     42,
				wantDefined: true,
			},
			{
				name:        "negative",
				extrema:     metricdata.NewExtrema[int64](-100),
				wantVal:     -100,
				wantDefined: true,
			},
			{
				name:        "explicit zero",
				extrema:     metricdata.NewExtrema[int64](0),
				wantVal:     0,
				wantDefined: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				val, ok := tc.extrema.Value()
				assert.Equal(t, tc.wantVal, val)
				assert.Equal(t, tc.wantDefined, ok)
			})
		}
	})

	t.Run("float64", func(t *testing.T) {
		testCases := []struct {
			name        string
			extrema     metricdata.Extrema[float64]
			wantVal     float64
			wantDefined bool
		}{
			{
				name:        "zero-value",
				extrema:     metricdata.Extrema[float64]{},
				wantVal:     0.0,
				wantDefined: false,
			},
			{
				name:        "positive",
				extrema:     metricdata.NewExtrema[float64](3.14159),
				wantVal:     3.14159,
				wantDefined: true,
			},
			{
				name:        "negative",
				extrema:     metricdata.NewExtrema[float64](-273.15),
				wantVal:     -273.15,
				wantDefined: true,
			},
			{
				name:        "explicit zero",
				extrema:     metricdata.NewExtrema[float64](0.0),
				wantVal:     0.0,
				wantDefined: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				val, ok := tc.extrema.Value()
				assert.Equal(t, tc.wantVal, val)
				assert.Equal(t, tc.wantDefined, ok)
			})
		}
	})
}

func TestExtremaMarshalText(t *testing.T) {
	testCases := []struct {
		name    string
		marshal func() ([]byte, error)
		want    string
	}{
		{
			name: "int64 zero-value returns null",
			marshal: func() ([]byte, error) {
				var e metricdata.Extrema[int64]
				return e.MarshalText()
			},
			want: "null",
		},
		{
			name: "int64 defined value returns number",
			marshal: func() ([]byte, error) {
				e := metricdata.NewExtrema[int64](100)
				return e.MarshalText()
			},
			want: "100",
		},
		{
			name: "float64 zero-value returns null",
			marshal: func() ([]byte, error) {
				var e metricdata.Extrema[float64]
				return e.MarshalText()
			},
			want: "null",
		},
		{
			name: "float64 defined value returns number",
			marshal: func() ([]byte, error) {
				e := metricdata.NewExtrema[float64](12.5)
				return e.MarshalText()
			},
			want: "12.5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.marshal()
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(b))
		})
	}
}

func TestExtremaMarshalJSON(t *testing.T) {
	type wrapperInt struct {
		Min *metricdata.Extrema[int64] `json:"min"`
		Max *metricdata.Extrema[int64] `json:"max"`
	}

	type wrapperFloat struct {
		Min *metricdata.Extrema[float64] `json:"min"`
		Max *metricdata.Extrema[float64] `json:"max"`
	}

	t.Run("int64 json roundtrip", func(t *testing.T) {
		zero := metricdata.Extrema[int64]{}
		val := metricdata.NewExtrema[int64](50)
		w := wrapperInt{
			Min: &zero,
			Max: &val,
		}

		b, err := json.Marshal(w)
		require.NoError(t, err)
		assert.JSONEq(t, `{"min":null,"max":50}`, string(b))
	})

	t.Run("float64 json roundtrip", func(t *testing.T) {
		val := metricdata.NewExtrema[float64](-1.5)
		zero := metricdata.Extrema[float64]{}
		w := wrapperFloat{
			Min: &val,
			Max: &zero,
		}

		b, err := json.Marshal(w)
		require.NoError(t, err)
		assert.JSONEq(t, `{"min":-1.5,"max":null}`, string(b))
	})

	t.Run("pointer receiver direct MarshalJSON", func(t *testing.T) {
		testCases := []struct {
			name    string
			marshal func() ([]byte, error)
			want    string
		}{
			{
				name: "int64 defined",
				marshal: func() ([]byte, error) {
					e := metricdata.NewExtrema[int64](999)
					return (&e).MarshalJSON()
				},
				want: "999",
			},
			{
				name: "int64 zero-value",
				marshal: func() ([]byte, error) {
					var zero metricdata.Extrema[int64]
					return (&zero).MarshalJSON()
				},
				want: "null",
			},
			{
				name: "float64 defined",
				marshal: func() ([]byte, error) {
					e := metricdata.NewExtrema[float64](-45.75)
					return (&e).MarshalJSON()
				},
				want: "-45.75",
			},
			{
				name: "float64 zero-value",
				marshal: func() ([]byte, error) {
					var zero metricdata.Extrema[float64]
					return (&zero).MarshalJSON()
				},
				want: "null",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				b, err := tc.marshal()
				require.NoError(t, err)
				assert.Equal(t, tc.want, string(b))
			})
		}
	})
}
