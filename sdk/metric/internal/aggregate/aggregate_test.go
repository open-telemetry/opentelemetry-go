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

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

var (
	keyUser   = "user"
	userAlice = attribute.String(keyUser, "Alice")
	adminTrue = attribute.Bool("admin", true)

	alice = attribute.NewSet(userAlice, adminTrue)

	// Filtered.
	attrFltr = func(kv attribute.KeyValue) bool {
		return kv.Key == attribute.Key(keyUser)
	}
	fltrAlice = attribute.NewSet(userAlice)
)

type inputTester[N int64 | float64] struct {
	aggregator[N]

	value N
	attr  attribute.Set
}

func (it *inputTester[N]) Aggregate(v N, a attribute.Set) { it.value, it.attr = v, a }

func TestBuilderInput(t *testing.T) {
	t.Run("Int64", testBuilderInput[int64]())
	t.Run("Float64", testBuilderInput[float64]())
}

func testBuilderInput[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		value, attr := N(1), alice
		run := func(b Builder[N], wantA attribute.Set) func(*testing.T) {
			return func(t *testing.T) {
				t.Helper()

				it := &inputTester[N]{}
				meas := b.input(it)
				meas(context.Background(), value, attr)

				assert.Equal(t, value, it.value, "measured incorrect value")
				assert.Equal(t, wantA, it.attr, "measured incorrect attributes")
			}
		}

		t.Run("NoFilter", run(Builder[N]{}, attr))
		t.Run("Filter", run(Builder[N]{Filter: attrFltr}, fltrAlice))
	}
}
