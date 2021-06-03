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

package baggage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/internal/baggage"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, Baggage{}, FromContext(ctx))

	b := Baggage{list: baggage.List{"key": baggage.Item{Value: "val"}}}
	ctx = ContextWithBaggage(ctx, b)
	assert.Equal(t, b, FromContext(ctx))

	ctx = ContextWithoutBaggage(ctx)
	assert.Equal(t, Baggage{}, FromContext(ctx))
}
