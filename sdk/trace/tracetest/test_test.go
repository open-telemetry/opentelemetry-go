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

package tracetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/trace"
)

// TestNoop tests only that the no-op does not crash in different scenarios.
func TestNoop(t *testing.T) {
	nsb := NewNoopExporter()

	require.NoError(t, nsb.ExportSpans(context.Background(), nil))
	require.NoError(t, nsb.ExportSpans(context.Background(), make([]*trace.SpanSnapshot, 10)))
	require.NoError(t, nsb.ExportSpans(context.Background(), make([]*trace.SpanSnapshot, 0, 10)))
}

func TestNewInMemoryExporter(t *testing.T) {
	imsb := NewInMemoryExporter()

	require.NoError(t, imsb.ExportSpans(context.Background(), nil))
	assert.Len(t, imsb.GetSpans(), 0)

	input := make([]*trace.SpanSnapshot, 10)
	for i := 0; i < 10; i++ {
		input[i] = new(trace.SpanSnapshot)
	}
	require.NoError(t, imsb.ExportSpans(context.Background(), input))
	sds := imsb.GetSpans()
	assert.Len(t, sds, 10)
	for i, sd := range sds {
		assert.Same(t, input[i], sd)
	}
	imsb.Reset()
	// Ensure that operations on the internal storage does not change the previously returned value.
	assert.Len(t, sds, 10)
	assert.Len(t, imsb.GetSpans(), 0)

	require.NoError(t, imsb.ExportSpans(context.Background(), input[0:1]))
	sds = imsb.GetSpans()
	assert.Len(t, sds, 1)
	assert.Same(t, input[0], sds[0])
}
