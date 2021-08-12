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

package oteltest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
)

func TestTraceStateFromKeyValues(t *testing.T) {
	ts, err := TraceStateFromKeyValues()
	require.NoError(t, err)
	assert.Equal(t, 0, ts.Len(), "empty attributes creats zero value TraceState")

	ts, err = TraceStateFromKeyValues(
		attribute.String("key0", "string"),
		attribute.Bool("key1", true),
		attribute.Int64("key2", 1),
		attribute.Float64("key3", 1.1),
	)
	require.NoError(t, err)
	expected := "key0=string,key1=true,key2=1,key3=1.1"
	assert.Equal(t, expected, ts.String())
}
