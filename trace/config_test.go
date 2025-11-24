// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestNewSpanConfig(t *testing.T) {
	k1v1 := attribute.String("key1", "value1")
	k1v2 := attribute.String("key1", "value2")
	k2v2 := attribute.String("key2", "value2")

	timestamp0 := time.Unix(0, 0)
	timestamp1 := time.Unix(0, 0)

	link1 := Link{
		SpanContext: SpanContext{traceID: TraceID([16]byte{1, 1}), spanID: SpanID{3}},
		Attributes:  []attribute.KeyValue{k1v1},
	}
	link2 := Link{
		SpanContext: SpanContext{traceID: TraceID([16]byte{1, 1}), spanID: SpanID{3}},
		Attributes:  []attribute.KeyValue{k1v2, k2v2},
	}

	tests := []struct {
		options  []SpanStartOption
		expected SpanConfig
	}{
		{
			// No non-zero-values should be set.
			[]SpanStartOption{},
			SpanConfig{},
		},
		{
			[]SpanStartOption{
				WithAttributes(k1v1),
			},
			SpanConfig{
				attributes: []attribute.KeyValue{k1v1},
			},
		},
		{
			// Multiple calls should append not overwrite.
			[]SpanStartOption{
				WithAttributes(k1v1),
				WithAttributes(k1v2),
				WithAttributes(k2v2),
			},
			SpanConfig{
				// No uniqueness is guaranteed by the API.
				attributes: []attribute.KeyValue{k1v1, k1v2, k2v2},
			},
		},
		{
			[]SpanStartOption{
				WithAttributes(k1v1, k1v2, k2v2),
			},
			SpanConfig{
				// No uniqueness is guaranteed by the API.
				attributes: []attribute.KeyValue{k1v1, k1v2, k2v2},
			},
		},
		{
			[]SpanStartOption{
				WithTimestamp(timestamp0),
			},
			SpanConfig{
				timestamp: timestamp0,
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls overwrites with last-one-wins.
				WithTimestamp(timestamp0),
				WithTimestamp(timestamp1),
			},
			SpanConfig{
				timestamp: timestamp1,
			},
		},
		{
			[]SpanStartOption{
				WithLinks(link1),
			},
			SpanConfig{
				links: []Link{link1},
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls should append not overwrite.
				WithLinks(link1),
				WithLinks(link1, link2),
			},
			SpanConfig{
				// No uniqueness is guaranteed by the API.
				links: []Link{link1, link1, link2},
			},
		},
		{
			[]SpanStartOption{
				WithNewRoot(),
			},
			SpanConfig{
				newRoot: true,
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls should not change NewRoot state.
				WithNewRoot(),
				WithNewRoot(),
			},
			SpanConfig{
				newRoot: true,
			},
		},
		{
			[]SpanStartOption{
				WithSpanKind(SpanKindConsumer),
			},
			SpanConfig{
				spanKind: SpanKindConsumer,
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls overwrites with last-one-wins.
				WithSpanKind(SpanKindClient),
				WithSpanKind(SpanKindConsumer),
			},
			SpanConfig{
				spanKind: SpanKindConsumer,
			},
		},
		{
			[]SpanStartOption{
				WithProfileTask(true),
			},
			SpanConfig{
				profileTask: boolPtr(true),
			},
		},
		{
			[]SpanStartOption{
				WithProfileTask(false),
			},
			SpanConfig{
				profileTask: boolPtr(false),
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls overwrites with last-one-wins.
				WithProfileTask(true),
				WithProfileTask(false),
			},
			SpanConfig{
				profileTask: boolPtr(false),
			},
		},
		{
			[]SpanStartOption{
				WithProfileRegion(true),
			},
			SpanConfig{
				profileRegion: boolPtr(true),
			},
		},
		{
			[]SpanStartOption{
				WithProfileRegion(false),
			},
			SpanConfig{
				profileRegion: boolPtr(false),
			},
		},
		{
			[]SpanStartOption{
				// Multiple calls overwrites with last-one-wins.
				WithProfileRegion(true),
				WithProfileRegion(false),
			},
			SpanConfig{
				profileRegion: boolPtr(false),
			},
		},
		{
			[]SpanStartOption{
				WithProfileTask(true),
				WithProfileRegion(true),
			},
			SpanConfig{
				profileTask:   boolPtr(true),
				profileRegion: boolPtr(true),
			},
		},
		{
			// Everything should work together.
			[]SpanStartOption{
				WithAttributes(k1v1),
				WithTimestamp(timestamp0),
				WithLinks(link1, link2),
				WithNewRoot(),
				WithSpanKind(SpanKindConsumer),
				WithProfileTask(true),
				WithProfileRegion(false),
			},
			SpanConfig{
				attributes:    []attribute.KeyValue{k1v1},
				timestamp:     timestamp0,
				links:         []Link{link1, link2},
				newRoot:       true,
				spanKind:      SpanKindConsumer,
				profileTask:   boolPtr(true),
				profileRegion: boolPtr(false),
			},
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, NewSpanStartConfig(test.options...))
	}
}

func TestSpanStartConfigAttributeMutability(t *testing.T) {
	a := attribute.String("a", "val")
	b := attribute.String("b", "val")
	attrs := []attribute.KeyValue{a, b}
	conf := NewSpanStartConfig(WithAttributes(attrs...))

	// Mutating passed arg should not change configured attributes.
	attrs[0] = attribute.String("c", "val")

	want := SpanConfig{attributes: []attribute.KeyValue{a, b}}
	assert.Equal(t, want, conf)
}

func TestEndSpanConfig(t *testing.T) {
	timestamp := time.Unix(0, 0)

	tests := []struct {
		options  []SpanEndOption
		expected SpanConfig
	}{
		{
			[]SpanEndOption{},
			SpanConfig{},
		},
		{
			[]SpanEndOption{
				WithStackTrace(true),
			},
			SpanConfig{
				stackTrace: true,
			},
		},
		{
			[]SpanEndOption{
				WithTimestamp(timestamp),
			},
			SpanConfig{
				timestamp: timestamp,
			},
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, NewSpanEndConfig(test.options...))
	}
}

func TestTracerConfig(t *testing.T) {
	v1 := "semver:0.0.1"
	v2 := "semver:1.0.0"
	schemaURL := "https://opentelemetry.io/schemas/1.21.0"
	attrs := attribute.NewSet(
		attribute.String("user", "alice"),
		attribute.Bool("admin", true),
	)

	c := NewTracerConfig(
		// Multiple calls should overwrite.
		WithInstrumentationVersion(v1),
		WithInstrumentationVersion(v2),
		WithSchemaURL(schemaURL),
		WithInstrumentationAttributes(attrs.ToSlice()...),
	)

	assert.Equal(t, v2, c.InstrumentationVersion(), "instrumentation version")
	assert.Equal(t, schemaURL, c.SchemaURL(), "schema URL")
	assert.Equal(t, attrs, c.InstrumentationAttributes(), "instrumentation attributes")
}

func TestWithInstrumentationAttributesNotLazy(t *testing.T) {
	attrs := []attribute.KeyValue{
		attribute.String("service", "test"),
		attribute.Int("three", 3),
	}
	want := attribute.NewSet(attrs...)

	// WithInstrumentationAttributes is expected to immediately
	// create an immutable set from the attributes, so later changes
	// to attrs should not affect the config.
	opt := WithInstrumentationAttributes(attrs...)
	attrs[0] = attribute.String("service", "changed")

	c := NewTracerConfig(opt)
	assert.Equal(t, want, c.InstrumentationAttributes(), "instrumentation attributes")
}

func TestWithInstrumentationAttributeSet(t *testing.T) {
	attrs := attribute.NewSet(
		attribute.String("service", "test"),
		attribute.Int("three", 3),
	)

	c := NewTracerConfig(
		WithInstrumentationAttributeSet(attrs),
	)

	assert.Equal(t, attrs, c.InstrumentationAttributes(), "instrumentation attributes")
}

// Save benchmark results to a file level var to avoid the compiler optimizing
// away the actual work.
var (
	spanConfig  SpanConfig
	eventConfig EventConfig
)

func BenchmarkNewTracerConfig(b *testing.B) {
	for _, bb := range []struct {
		name    string
		options []TracerOption
	}{
		{
			name: "with no options",
		},
		{
			name: "with an instrumentation version",
			options: []TracerOption{
				WithInstrumentationVersion("testing version"),
			},
		},
		{
			name: "with a schema url",
			options: []TracerOption{
				WithSchemaURL("testing URL"),
			},
		},
		{
			name: "with instrumentation attribute",
			options: []TracerOption{
				WithInstrumentationAttributes(attribute.String("key", "value")),
			},
		},
		{
			name: "with instrumentation attribute set",
			options: []TracerOption{
				WithInstrumentationAttributeSet(attribute.NewSet(attribute.String("key", "value"))),
			},
		},
	} {
		b.Run(bb.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				NewTracerConfig(bb.options...)
			}
		})
	}
}

func BenchmarkNewSpanStartConfig(b *testing.B) {
	for _, bb := range []struct {
		name    string
		options []SpanStartOption
	}{
		{
			name: "with no options",
		},
		{
			name: "with attributes",
			options: []SpanStartOption{
				WithAttributes(attribute.Bool("key", true)),
			},
		},
		{
			name: "with attributes set multiple times",
			options: []SpanStartOption{
				WithAttributes(attribute.Bool("key", true)),
				WithAttributes(attribute.Bool("secondKey", false)),
			},
		},
		{
			name: "with a timestamp",
			options: []SpanStartOption{
				WithTimestamp(time.Now()),
			},
		},
		{
			name: "with links",
			options: []SpanStartOption{
				WithLinks(Link{}),
			},
		},
		{
			name: "with links set multiple times",
			options: []SpanStartOption{
				WithLinks(Link{}),
				WithLinks(Link{}),
			},
		},
		{
			name: "with new root",
			options: []SpanStartOption{
				WithNewRoot(),
			},
		},
		{
			name: "with span kind",
			options: []SpanStartOption{
				WithSpanKind(SpanKindClient),
			},
		},
		{
			name: "with profile task",
			options: []SpanStartOption{
				WithProfileTask(true),
			},
		},
		{
			name: "with profile region",
			options: []SpanStartOption{
				WithProfileRegion(true),
			},
		},
		{
			name: "with profile task and region",
			options: []SpanStartOption{
				WithProfileTask(true),
				WithProfileRegion(false),
			},
		},
	} {
		b.Run(bb.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				spanConfig = NewSpanStartConfig(bb.options...)
			}
		})
	}
}

func BenchmarkNewSpanEndConfig(b *testing.B) {
	for _, bb := range []struct {
		name    string
		options []SpanEndOption
	}{
		{
			name: "with no options",
		},
		{
			name: "with a timestamp",
			options: []SpanEndOption{
				WithTimestamp(time.Now()),
			},
		},
		{
			name: "with stack trace",
			options: []SpanEndOption{
				WithStackTrace(true),
			},
		},
	} {
		b.Run(bb.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				spanConfig = NewSpanEndConfig(bb.options...)
			}
		})
	}
}

func BenchmarkNewEventConfig(b *testing.B) {
	for _, bb := range []struct {
		name    string
		options []EventOption
	}{
		{
			name: "with no options",
		},
		{
			name: "with attributes",
			options: []EventOption{
				WithAttributes(attribute.Bool("key", true)),
			},
		},
		{
			name: "with attributes set multiple times",
			options: []EventOption{
				WithAttributes(attribute.Bool("key", true)),
				WithAttributes(attribute.Bool("secondKey", false)),
			},
		},
		{
			name: "with a timestamp",
			options: []EventOption{
				WithTimestamp(time.Now()),
			},
		},
		{
			name: "with a stacktrace",
			options: []EventOption{
				WithStackTrace(true),
			},
		},
	} {
		b.Run(bb.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				eventConfig = NewEventConfig(bb.options...)
			}
		})
	}
}

func TestWithInstrumentationAttributesMerge(t *testing.T) {
	aliceAttr := attribute.String("user", "Alice")
	bobAttr := attribute.String("user", "Bob")
	adminAttr := attribute.Bool("admin", true)

	alice := attribute.NewSet(aliceAttr)
	bob := attribute.NewSet(bobAttr)
	aliceAdmin := attribute.NewSet(aliceAttr, adminAttr)
	bobAdmin := attribute.NewSet(bobAttr, adminAttr)

	t.Run("SameKey", func(t *testing.T) {
		c := NewTracerConfig(
			WithInstrumentationAttributes(aliceAttr),
			WithInstrumentationAttributes(bobAttr),
		)
		assert.Equal(t, bob, c.InstrumentationAttributes(),
			"Later values for the same key should overwrite earlier ones.")
	})

	t.Run("DifferentKeys", func(t *testing.T) {
		c := NewTracerConfig(
			WithInstrumentationAttributes(aliceAttr),
			WithInstrumentationAttributes(adminAttr),
		)
		assert.Equal(t, aliceAdmin, c.InstrumentationAttributes(),
			"Different keys should be merged")
	})

	t.Run("Mixed", func(t *testing.T) {
		c := NewTracerConfig(
			WithInstrumentationAttributes(aliceAttr, adminAttr),
			WithInstrumentationAttributes(bobAttr),
		)
		assert.Equal(t, bobAdmin, c.InstrumentationAttributes(),
			"Combination of same and different keys should be merged.")
	})

	t.Run("MergedEmpty", func(t *testing.T) {
		c := NewTracerConfig(
			WithInstrumentationAttributes(aliceAttr),
			WithInstrumentationAttributes(),
		)
		assert.Equal(t, alice, c.InstrumentationAttributes(),
			"Empty attributes should not affect existing ones.")
	})

	t.Run("SameKeyWithSet", func(t *testing.T) {
		c := NewTracerConfig(
			WithInstrumentationAttributeSet(alice),
			WithInstrumentationAttributeSet(bob),
		)
		assert.Equal(t, bob, c.InstrumentationAttributes(),
			"Later values for the same key should overwrite earlier ones.")
	})

	t.Run("DifferentKeysWithSet", func(t *testing.T) {
		c := NewTracerConfig(
			WithInstrumentationAttributeSet(alice),
			WithInstrumentationAttributeSet(attribute.NewSet(adminAttr)),
		)
		assert.Equal(t, aliceAdmin, c.InstrumentationAttributes(),
			"Different keys should be merged.")
	})

	t.Run("MixedWithSet", func(t *testing.T) {
		c := NewTracerConfig(
			WithInstrumentationAttributeSet(aliceAdmin),
			WithInstrumentationAttributeSet(bob),
		)
		assert.Equal(t, bobAdmin, c.InstrumentationAttributes(),
			"Combination of same and different keys should be merged.")
	})

	t.Run("MergedEmptyWithSet", func(t *testing.T) {
		c := NewTracerConfig(
			WithInstrumentationAttributeSet(alice),
			WithInstrumentationAttributeSet(attribute.NewSet()),
		)
		assert.Equal(t, alice, c.InstrumentationAttributes(),
			"Empty attribute set should not affect existing ones.")
	})

	t.Run("MixedAttributesAndSet", func(t *testing.T) {
		c := NewTracerConfig(
			WithInstrumentationAttributes(aliceAttr),
			WithInstrumentationAttributeSet(attribute.NewSet(bobAttr, adminAttr)),
		)
		assert.Equal(t, bobAdmin, c.InstrumentationAttributes(),
			"Attributes and attribute sets should be merged together.")
	})
}
