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

func TestAddLinks(t *testing.T) {
	ctx, tp := context.Background(), NewTracerProvider()
	defer func(ctx context.Context, tp *TracerProvider) {
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx, tp)

	ctx, span := tp.Tracer("test").Start(ctx, "test_add_links")
	defer span.End()

	links := []trace.Link{
		trace.LinkFromContext(ctx),
	}

	if s, ok := span.(trace.ExperimentalSpan); ok {
		s.AddLinks(links...)
	}

	assert.Equal(t, len(links), len(span.(ReadOnlySpan).Links()))

	for i, l := range span.(ReadOnlySpan).Links() {
		assert.Equal(t, links[i], trace.Link{
			Attributes:  l.Attributes,
			SpanContext: l.SpanContext,
		})
	}
}
