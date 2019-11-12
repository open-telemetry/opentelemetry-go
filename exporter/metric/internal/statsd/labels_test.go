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
)

func TestLabelSytnax(t *testing.T) {
	encoder := statsd.NewLabelEncoder()

	require.Equal(t, `|#A:B,C:D,E:1.5`, encoder.EncodeLabels([]core.KeyValue{
		key.New("A").String("B"),
		key.New("C").String("D"),
		key.New("E").Float64(1.5),
	}))
}
