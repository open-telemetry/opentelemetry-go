// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/metric"

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

type attrConf interface {
	Attributes() attribute.Set
}

func TestConfigAttrs(t *testing.T) {
	t.Run("AddConfig", testConfAttr(func(mo ...MeasurementOption) attrConf {
		opts := make([]AddOption, len(mo))
		for i := range mo {
			opts[i] = mo[i].(AddOption)
		}
		return NewAddConfig(opts)
	}))

	t.Run("RecordConfig", testConfAttr(func(mo ...MeasurementOption) attrConf {
		opts := make([]RecordOption, len(mo))
		for i := range mo {
			opts[i] = mo[i].(RecordOption)
		}
		return NewRecordConfig(opts)
	}))

	t.Run("ObserveConfig", testConfAttr(func(mo ...MeasurementOption) attrConf {
		opts := make([]ObserveOption, len(mo))
		for i := range mo {
			opts[i] = mo[i].(ObserveOption)
		}
		return NewObserveConfig(opts)
	}))

	t.Run("FinishConfig", testFinishConfAttr(func(fo ...FinishOption) attrConf {
		return NewFinishConfig(fo)
	}))
}

func testConfAttr(newConf func(...MeasurementOption) attrConf) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("ZeroConfigEmpty", func(t *testing.T) {
			c := newConf()
			assert.Equal(t, *attribute.EmptySet(), c.Attributes())
		})

		t.Run("EmptySet", func(t *testing.T) {
			c := newConf(WithAttributeSet(*attribute.EmptySet()))
			assert.Equal(t, *attribute.EmptySet(), c.Attributes())
		})

		aliceAttr := attribute.String("user", "Alice")
		alice := attribute.NewSet(aliceAttr)
		t.Run("SingleWithAttributeSet", func(t *testing.T) {
			c := newConf(WithAttributeSet(alice))
			assert.Equal(t, alice, c.Attributes())
		})

		t.Run("SingleWithAttributes", func(t *testing.T) {
			c := newConf(WithAttributes(aliceAttr))
			assert.Equal(t, alice, c.Attributes())
		})

		bobAttr := attribute.String("user", "Bob")
		bob := attribute.NewSet(bobAttr)
		t.Run("MultiWithAttributeSet", func(t *testing.T) {
			c := newConf(WithAttributeSet(alice), WithAttributeSet(bob))
			assert.Equal(t, bob, c.Attributes())
		})

		t.Run("MergedWithAttributes", func(t *testing.T) {
			c := newConf(WithAttributes(aliceAttr, bobAttr))
			assert.Equal(t, bob, c.Attributes())
		})

		t.Run("MultiWithAttributeSet", func(t *testing.T) {
			c := newConf(WithAttributes(aliceAttr), WithAttributes(bobAttr))
			assert.Equal(t, bob, c.Attributes())
		})

		t.Run("MergedEmpty", func(t *testing.T) {
			c := newConf(WithAttributeSet(alice), WithAttributeSet(*attribute.EmptySet()))
			assert.Equal(t, alice, c.Attributes())
		})
	}
}

func testFinishConfAttr(newConf func(...FinishOption) attrConf) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("ZeroConfigEmpty", func(t *testing.T) {
			c := newConf()
			assert.Equal(t, *attribute.EmptySet(), c.Attributes())
		})

		t.Run("EmptySet", func(t *testing.T) {
			c := newConf(WithAttributeSet(*attribute.EmptySet()))
			assert.Equal(t, *attribute.EmptySet(), c.Attributes())
		})

		aliceAttr := attribute.String("user", "Alice")
		alice := attribute.NewSet(aliceAttr)
		t.Run("SingleWithAttributeSet", func(t *testing.T) {
			c := newConf(WithAttributeSet(alice))
			assert.Equal(t, alice, c.Attributes())
		})

		t.Run("SingleWithAttributes", func(t *testing.T) {
			c := newConf(WithAttributes(aliceAttr))
			assert.Equal(t, alice, c.Attributes())
		})

		bobAttr := attribute.String("user", "Bob")
		bob := attribute.NewSet(bobAttr)
		t.Run("MultiWithAttributeSet", func(t *testing.T) {
			c := newConf(WithAttributeSet(alice), WithAttributeSet(bob))
			assert.Equal(t, bob, c.Attributes())
		})

		t.Run("MergedWithAttributes", func(t *testing.T) {
			c := newConf(WithAttributes(aliceAttr, bobAttr))
			assert.Equal(t, bob, c.Attributes())
		})

		t.Run("MultiWithAttributes", func(t *testing.T) {
			c := newConf(WithAttributes(aliceAttr), WithAttributes(bobAttr))
			assert.Equal(t, bob, c.Attributes())
		})

		t.Run("MergedEmpty", func(t *testing.T) {
			c := newConf(WithAttributeSet(alice), WithAttributeSet(*attribute.EmptySet()))
			assert.Equal(t, alice, c.Attributes())
		})
	}
}

func TestWithAttributesConcurrentSafe(*testing.T) {
	attrs := []attribute.KeyValue{
		attribute.String("user", "Alice"),
		attribute.Bool("admin", true),
		attribute.String("user", "Bob"),
	}

	var wg sync.WaitGroup
	wg.Go(func() {
		opt := []AddOption{WithAttributes(attrs...)}
		_ = NewAddConfig(opt)
	})
	wg.Go(func() {
		opt := []AddOption{WithAttributes(attrs...)}
		_ = NewAddConfig(opt)
	})

	wg.Go(func() {
		opt := []RecordOption{WithAttributes(attrs...)}
		_ = NewRecordConfig(opt)
	})
	wg.Go(func() {
		opt := []RecordOption{WithAttributes(attrs...)}
		_ = NewRecordConfig(opt)
	})

	wg.Go(func() {
		opt := []ObserveOption{WithAttributes(attrs...)}
		_ = NewObserveConfig(opt)
	})
	wg.Go(func() {
		opt := []ObserveOption{WithAttributes(attrs...)}
		_ = NewObserveConfig(opt)
	})

	wg.Wait()
}

func TestFinishConfigMatcher(t *testing.T) {
	containerA := attribute.NewSet(
		attribute.String("container.id", "a"),
		attribute.String("pod", "api-0"),
	)
	containerB := attribute.NewSet(
		attribute.String("container.id", "b"),
		attribute.String("pod", "api-1"),
	)
	empty := *attribute.EmptySet()

	t.Run("DefaultMatchesEmptySet", func(t *testing.T) {
		c := NewFinishConfig(nil)
		assert.True(t, c.Matcher()(empty))
		assert.False(t, c.Matcher()(containerA))
	})

	t.Run("ExactAttributesOnly", func(t *testing.T) {
		c := NewFinishConfig([]FinishOption{WithAttributeSet(containerA)})
		assert.True(t, c.Matcher()(containerA))
		assert.False(t, c.Matcher()(containerB))
	})

	t.Run("MatcherOnly", func(t *testing.T) {
		c := NewFinishConfig([]FinishOption{
			WithMatchAttributes(func(attrs attribute.Set) bool {
				v, ok := (&attrs).Value("container.id")
				return ok && v.AsString() == "a"
			}),
		})
		assert.True(t, c.Matcher()(containerA))
		assert.False(t, c.Matcher()(containerB))
		assert.False(t, c.Matcher()(empty))
	})

	t.Run("ExactAndMatcherBothNeedToMatch", func(t *testing.T) {
		c := NewFinishConfig([]FinishOption{
			WithAttributeSet(containerA),
			WithMatchAttributes(func(attrs attribute.Set) bool {
				v, ok := (&attrs).Value("pod")
				return ok && v.AsString() == "api-0"
			}),
		})
		assert.True(t, c.Matcher()(containerA))
		assert.False(t, c.Matcher()(containerB))
	})
}
