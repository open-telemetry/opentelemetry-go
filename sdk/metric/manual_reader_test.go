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

package metric // import "go.opentelemetry.io/otel/sdk/metric/reader"

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestManualReader(t *testing.T) {
	suite.Run(t, &readerTestSuite{Factory: func() Reader { return NewManualReader() }})
}

func BenchmarkManualReader(b *testing.B) {
	b.Run("Collect", benchReaderCollectFunc(NewManualReader()))
}

var deltaTemporalitySelector = func(InstrumentKind) Temporality { return DeltaTemporality }
var cumulativeTemporalitySelector = func(InstrumentKind) Temporality { return CumulativeTemporality }

func TestManualReaderTemporality(t *testing.T) {
	tests := []struct {
		name    string
		options []ManualReaderOption
		// Currently only testing constant temporality. This should be expanded
		// if we put more advanced selection in the SDK
		wantTemporality Temporality
	}{
		{
			name:            "default",
			wantTemporality: CumulativeTemporality,
		},
		{
			name: "delta",
			options: []ManualReaderOption{
				WithTemporality(deltaTemporalitySelector),
			},
			wantTemporality: DeltaTemporality,
		},
		{
			name: "repeats overwrite",
			options: []ManualReaderOption{
				WithTemporality(deltaTemporalitySelector),
				WithTemporality(cumulativeTemporalitySelector),
			},
			wantTemporality: CumulativeTemporality,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdr := NewManualReader(tt.options...)
			assert.Equal(t, tt.wantTemporality, rdr.temporality(undefinedInstrument))
		})
	}
}
