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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestManualReader(t *testing.T) {
	suite.Run(t, &readerTestSuite{Factory: func() Reader { return NewManualReader() }})
}

func BenchmarkManualReader(b *testing.B) {
	b.Run("Collect", benchReaderCollectFunc(NewManualReader()))
}

var deltaTemporalitySelector = func(InstrumentKind) metricdata.Temporality { return metricdata.DeltaTemporality }
var cumulativeTemporalitySelector = func(InstrumentKind) metricdata.Temporality { return metricdata.CumulativeTemporality }

func TestManualReaderTemporality(t *testing.T) {
	tests := []struct {
		name    string
		options []ManualReaderOption
		// Currently only testing constant temporality. This should be expanded
		// if we put more advanced selection in the SDK
		wantTemporality metricdata.Temporality
	}{
		{
			name:            "default",
			wantTemporality: metricdata.CumulativeTemporality,
		},
		{
			name: "delta",
			options: []ManualReaderOption{
				WithTemporalitySelector(deltaTemporalitySelector),
			},
			wantTemporality: metricdata.DeltaTemporality,
		},
		{
			name: "repeats overwrite",
			options: []ManualReaderOption{
				WithTemporalitySelector(deltaTemporalitySelector),
				WithTemporalitySelector(cumulativeTemporalitySelector),
			},
			wantTemporality: metricdata.CumulativeTemporality,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var undefinedInstrument InstrumentKind
			rdr := NewManualReader(tt.options...)
			assert.Equal(t, tt.wantTemporality, rdr.temporality(undefinedInstrument))
		})
	}
}

func TestManualReaderCollect(t *testing.T) {
	expiredCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1))
	defer cancel()

	tests := []struct {
		name string

		ctx             context.Context
		resourceMetrics *metricdata.ResourceMetrics

		expectedErr error
	}{
		{
			name: "with a valid context",

			ctx:             context.Background(),
			resourceMetrics: &metricdata.ResourceMetrics{},
		},
		{
			name: "with an expired context",

			ctx:             expiredCtx,
			resourceMetrics: &metricdata.ResourceMetrics{},

			expectedErr: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdr := NewManualReader()
			mp := NewMeterProvider(WithReader(rdr))
			meter := mp.Meter("test")

			// Ensure the pipeline has a callback setup
			testM, err := meter.Int64ObservableCounter("test")
			assert.NoError(t, err)
			_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
				return nil
			}, testM)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedErr, rdr.Collect(tt.ctx, tt.resourceMetrics))
		})
	}
}
