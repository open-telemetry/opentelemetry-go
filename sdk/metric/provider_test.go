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

package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/reader"
)

// TODO: incomplete
func TestOutputReuse(t *testing.T) {
	ctx := context.Background()
	rdr := metrictest.NewReader()
	provider := New(WithReader(rdr))
	producer := rdr.Producer()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	reuse := reader.Metrics{
		Scopes: make([]reader.Scope, 0, 25),
	}

	cntr.Add(ctx, 1, attribute.Int("K", 1))

	reuse = producer.Produce(ctx, &reuse)
	assert.Len(t, reuse.Scopes, 1)
	assert.Equal(t, 25, cap(reuse.Scopes))

	cntr.Add(ctx, 1, attribute.Int("K", 1))

	reuse = producer.Produce(ctx, &reuse)
	assert.Len(t, reuse.Scopes, 1)
	assert.Equal(t, 25, cap(reuse.Scopes))
}
