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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestLimitAttr(t *testing.T) {
	m := map[attribute.Set]struct{}{alice: {}}

	t.Run("NoLimit", func(t *testing.T) {
		assert.Equal(t, alice, limitAttr(alice, m, 0))
		assert.Equal(t, bob, limitAttr(bob, m, 0))
	})

	t.Run("NotAtLimit/Exists", func(t *testing.T) {
		assert.Equal(t, alice, limitAttr(alice, m, 3))
	})

	t.Run("NotAtLimit/DoesNotExist", func(t *testing.T) {
		assert.Equal(t, bob, limitAttr(bob, m, 3))
	})

	t.Run("AtLimit/Exists", func(t *testing.T) {
		assert.Equal(t, alice, limitAttr(alice, m, 2))
	})

	t.Run("AtLimit/DoesNotExist", func(t *testing.T) {
		assert.Equal(t, overflowSet, limitAttr(bob, m, 2))
	})
}
