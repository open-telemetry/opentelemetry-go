// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestLimiterAttributes(t *testing.T) {
	m := map[attribute.Distinct]struct{}{alice.Equivalent(): {}}
	t.Run("NoLimit", func(t *testing.T) {
		l := newLimiter[struct{}](0)
		assert.Equal(t, alice, l.Attributes(alice, len(m), m))
		assert.Equal(t, bob, l.Attributes(bob, len(m), m))
	})

	t.Run("NotAtLimit/Exists", func(t *testing.T) {
		l := newLimiter[struct{}](3)
		assert.Equal(t, alice, l.Attributes(alice, len(m), m))
	})

	t.Run("NotAtLimit/DoesNotExist", func(t *testing.T) {
		l := newLimiter[struct{}](3)
		assert.Equal(t, bob, l.Attributes(bob, len(m), m))
	})

	t.Run("AtLimit/Exists", func(t *testing.T) {
		l := newLimiter[struct{}](2)
		assert.Equal(t, alice, l.Attributes(alice, len(m), m))
	})

	t.Run("AtLimit/DoesNotExist", func(t *testing.T) {
		l := newLimiter[struct{}](2)
		assert.Equal(t, overflowSet, l.Attributes(bob, len(m), m))
	})
}

var limitedAttr attribute.Set

func BenchmarkLimiterAttributes(b *testing.B) {
	m := map[attribute.Distinct]struct{}{alice.Equivalent(): {}}
	l := newLimiter[struct{}](2)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		limitedAttr = l.Attributes(alice, len(m), m)
		limitedAttr = l.Attributes(bob, len(m), m)
	}
}
