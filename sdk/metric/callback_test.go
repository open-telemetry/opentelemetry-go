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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
)

func TestCallbackCollectError(t *testing.T) {
	cBack := callback[int64]{
		observe: func(context.Context, int64, ...attribute.KeyValue) {
			assert.Fail(t, "observe should not be called")
		},
		newIter: newInt64Iter(func(context.Context) ([]asyncint64.Observation, error) {
			return nil, assert.AnError
		}),
	}

	assert.NotPanics(t, func() {
		assert.ErrorIs(t, cBack.collect(context.Background()), assert.AnError)
	})
}

func TestCallbackCollectInt64(t *testing.T) {
	const val int64 = 42
	attrs := []attribute.KeyValue{attribute.String("name", "alice")}

	cBack := callback[int64]{
		observe: func(_ context.Context, v int64, a ...attribute.KeyValue) {
			assert.Equal(t, val, v, "recorded value")
			assert.Equal(t, attrs, a, "recorded attribute")
		},
		newIter: newInt64Iter(func(context.Context) ([]asyncint64.Observation, error) {
			return []asyncint64.Observation{{Attributes: attrs, Value: val}}, nil
		}),
	}

	assert.NotPanics(t, func() {
		assert.NoError(t, cBack.collect(context.Background()))
	})
}

func TestCallbackCollectFloat64(t *testing.T) {
	const val float64 = 42.
	attrs := []attribute.KeyValue{attribute.String("name", "alice")}

	cBack := callback[float64]{
		observe: func(_ context.Context, v float64, a ...attribute.KeyValue) {
			assert.Equal(t, val, v, "recorded value")
			assert.Equal(t, attrs, a, "recorded attribute")
		},
		newIter: newFloat64Iter(func(context.Context) ([]asyncfloat64.Observation, error) {
			return []asyncfloat64.Observation{{Attributes: attrs, Value: val}}, nil
		}),
	}

	assert.NotPanics(t, func() {
		assert.NoError(t, cBack.collect(context.Background()))
	})
}
