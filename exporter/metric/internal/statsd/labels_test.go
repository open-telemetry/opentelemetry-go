// Copyright 2019, OpenTelemetry Authors
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

package statsd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/exporter/metric/internal/statsd"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
)

var testLabels = []core.KeyValue{
	key.New("A").String("B"),
	key.New("C").String("D"),
	key.New("E").Float64(1.5),
}

func TestLabelSyntax(t *testing.T) {
	encoder := statsd.NewLabelEncoder()

	require.Equal(t, `|#A:B,C:D,E:1.5`, encoder.Encode(testLabels))

	require.Equal(t, `|#A:B`, encoder.Encode([]core.KeyValue{
		key.New("A").String("B"),
	}))

	require.Equal(t, "", encoder.Encode(nil))
}

func TestLabelForceEncode(t *testing.T) {
	defaultLabelEncoder := sdk.NewDefaultLabelEncoder()
	statsdLabelEncoder := statsd.NewLabelEncoder()

	exportLabelsDefault := export.NewLabels(testLabels, defaultLabelEncoder.Encode(testLabels), defaultLabelEncoder)
	exportLabelsStatsd := export.NewLabels(testLabels, statsdLabelEncoder.Encode(testLabels), statsdLabelEncoder)

	statsdEncoding := exportLabelsStatsd.Encoded()
	require.NotEqual(t, statsdEncoding, exportLabelsDefault.Encoded())

	forced, repeat := statsdLabelEncoder.ForceEncode(exportLabelsDefault)
	require.Equal(t, statsdEncoding, forced)
	require.True(t, repeat)

	forced, repeat = statsdLabelEncoder.ForceEncode(exportLabelsStatsd)
	require.Equal(t, statsdEncoding, forced)
	require.False(t, repeat)

	// Check that this works for an embedded implementation.
	exportLabelsEmbed := export.NewLabels(testLabels, statsdEncoding, struct {
		*statsd.LabelEncoder
	}{LabelEncoder: statsdLabelEncoder})

	forced, repeat = statsdLabelEncoder.ForceEncode(exportLabelsEmbed)
	require.Equal(t, statsdEncoding, forced)
	require.False(t, repeat)
}
