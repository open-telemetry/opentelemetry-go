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

func TestLimiterAttributes(t *testing.T) {
	m := map[attribute.Set]struct{}{alice: {}}
	t.Run("NoLimit", func(t *testing.T) {
		l := newLimiter[struct{}](0)
		assert.Equal(t, alice, l.Attributes(alice, m))
		assert.Equal(t, bob, l.Attributes(bob, m))
	})

	t.Run("NotAtLimit/Exists", func(t *testing.T) {
		l := newLimiter[struct{}](3)
		assert.Equal(t, alice, l.Attributes(alice, m))
	})

	t.Run("NotAtLimit/DoesNotExist", func(t *testing.T) {
		l := newLimiter[struct{}](3)
		assert.Equal(t, bob, l.Attributes(bob, m))
	})

	t.Run("AtLimit/Exists", func(t *testing.T) {
		l := newLimiter[struct{}](2)
		assert.Equal(t, alice, l.Attributes(alice, m))
	})

	t.Run("AtLimit/DoesNotExist", func(t *testing.T) {
		l := newLimiter[struct{}](2)
		assert.Equal(t, overflowSet, l.Attributes(bob, m))
	})
}

var limitedAttr attribute.Set

func BenchmarkLimiterAttributes(b *testing.B) {
	m := map[attribute.Set]struct{}{alice: {}}
	l := newLimiter[struct{}](2)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		limitedAttr = l.Attributes(alice, m)
		limitedAttr = l.Attributes(bob, m)
	}
}
