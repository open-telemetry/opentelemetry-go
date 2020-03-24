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

package dogstatsd_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/dogstatsd"
	"go.opentelemetry.io/otel/exporters/metric/internal/statsd"
	"go.opentelemetry.io/otel/exporters/metric/test"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

// TestDogstatsLabels that labels are formatted in the correct style,
// whether or not the provided labels were encoded by a statsd label
// encoder.
func TestDogstatsLabels(t *testing.T) {
	encoder := statsd.NewLabelEncoder()
	ctx := context.Background()
	checkpointSet := test.NewCheckpointSet(encoder)

	desc := metric.NewDescriptor("test.name", metric.CounterKind, core.Int64NumberKind)
	cagg := sum.New()
	_ = cagg.Update(ctx, core.NewInt64Number(123), &desc)
	cagg.Checkpoint(ctx, &desc)

	checkpointSet.Add(&desc, cagg, key.New("A").String("B"))

	var buf bytes.Buffer
	exp, err := dogstatsd.NewRawExporter(dogstatsd.Config{
		Writer: &buf,
	})
	require.Nil(t, err)

	err = exp.Export(ctx, checkpointSet)
	require.Nil(t, err)

	require.Equal(t, "test.name:123|c|#A:B\n", buf.String())
}
