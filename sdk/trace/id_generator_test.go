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

package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/trace"
)

func TestNewIDs(t *testing.T) {
	gen := defaultIDGenerator()
	n := 1000

	for i := 0; i < n; i++ {
		traceID, spanID := gen.NewIDs(context.Background())
		assert.Truef(t, traceID.IsValid(), "trace id: %s", traceID.String())
		assert.Truef(t, spanID.IsValid(), "span id: %s", spanID.String())
	}
}

func TestNewSpanID(t *testing.T) {
	gen := defaultIDGenerator()
	testTraceID := [16]byte{123, 123}
	n := 1000

	for i := 0; i < n; i++ {
		spanID := gen.NewSpanID(context.Background(), testTraceID)
		assert.Truef(t, spanID.IsValid(), "span id: %s", spanID.String())
	}
}

func TestNewSpanIDWithInvalidTraceID(t *testing.T) {
	gen := defaultIDGenerator()
	spanID := gen.NewSpanID(context.Background(), trace.TraceID{})
	assert.Truef(t, spanID.IsValid(), "span id: %s", spanID.String())
}
