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
)

func TestContextWithBaggage(t *testing.T) {
	ctx := context.Background()
	b := Baggage{list: map[string]value{"foo": {v: "1"}}}

	nCtx := ContextWithBaggage(ctx, b)
	assert.Equal(t, b, nCtx.Value(baggageKey))
	assert.Nil(t, ctx.Value(baggageKey))
}

func TestContextWithoutBaggage(t *testing.T) {
	b := Baggage{list: map[string]value{"foo": {v: "1"}}}

	ctx := context.Background()
	ctx = context.WithValue(ctx, baggageKey, b)

	nCtx := ContextWithoutBaggage(ctx)
	assert.Nil(t, nCtx.Value(baggageKey))
	assert.Equal(t, b, ctx.Value(baggageKey))
}

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, Baggage{}, FromContext(ctx))

	b := Baggage{list: map[string]value{"foo": {v: "1"}}}
	ctx = context.WithValue(ctx, baggageKey, b)
	assert.Equal(t, b, FromContext(ctx))
}
