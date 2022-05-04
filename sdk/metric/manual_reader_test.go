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

package metric // import "go.opentelemetry.io/otel/sdk/metric/reader"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestManualReaderNotRegistered(t *testing.T) {
	rdr := &ManualReader{}

	_, err := rdr.Collect(context.Background())
	require.ErrorIs(t, err, ErrReaderNotRegistered)
}

type testProducer struct{}

var testMetrics = export.Metrics{
	Scopes: []export.Scope{
		{
			Library: instrumentation.Library{
				Name:    "TestLibrary",
				Version: "0.0.1-beta1",
			},
		},
	},
}

func (p testProducer) produce(context.Context) export.Metrics {
	return testMetrics
}

func TestManualReaderProducer(t *testing.T) {
	rdr := &ManualReader{}
	rdr.register(testProducer{})

	m, err := rdr.Collect(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, testMetrics, m)
}
