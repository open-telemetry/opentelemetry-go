// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
)

func TestNewEntity(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		e := NewEntity(
			"service",
			[]attribute.KeyValue{
				attribute.String("service.name", "cart-service"),
			},
			WithEntityDescription(attribute.String("service.version", "1.0.0")),
			WithEntitySchemaURL("https://opentelemetry.io/schemas/1.20.0"),
		)

		assert.Equal(t, "service", e.Type())
		assert.Equal(t, []attribute.KeyValue{attribute.String("service.name", "cart-service")}, e.ID())
		assert.Equal(t, []attribute.KeyValue{attribute.String("service.version", "1.0.0")}, e.Description())
		assert.Equal(t, "https://opentelemetry.io/schemas/1.20.0", e.SchemaURL())
		assert.ElementsMatch(t, []attribute.KeyValue{
			attribute.String("service.name", "cart-service"),
			attribute.String("service.version", "1.0.0"),
		}, e.Attributes())
	})

	t.Run("DuplicateIDKeys", func(t *testing.T) {
		e := NewEntity("service", []attribute.KeyValue{
			attribute.String("service.name", "first"),
			attribute.String("service.name", "second"),
		})
		assert.Equal(t, []attribute.KeyValue{attribute.String("service.name", "second")}, e.ID())
	})

	t.Run("IDPrecedenceOverDescription", func(t *testing.T) {
		e := NewEntity(
			"service",
			[]attribute.KeyValue{attribute.String("service.name", "cart-service")},
			WithEntityDescription(
				attribute.String("service.name", "ignored"),
				attribute.String("service.version", "2.0.0"),
			),
		)
		assert.Equal(t, []attribute.KeyValue{attribute.String("service.name", "cart-service")}, e.ID())
		assert.Equal(t, []attribute.KeyValue{attribute.String("service.version", "2.0.0")}, e.Description())
	})

	t.Run("NilEntityMethods", func(t *testing.T) {
		var e *Entity
		assert.Empty(t, e.Type())
		assert.Nil(t, e.ID())
		assert.Nil(t, e.Description())
		assert.Empty(t, e.SchemaURL())
		assert.Nil(t, e.Attributes())
		assert.True(t, e.Equal(nil))
	})

	t.Run("Equal", func(t *testing.T) {
		e1 := NewEntity("service", []attribute.KeyValue{attribute.String("service.name", "cart")})
		e2 := NewEntity("service", []attribute.KeyValue{attribute.String("service.name", "cart")})
		e3 := NewEntity("service", []attribute.KeyValue{attribute.String("service.name", "other")})

		assert.True(t, e1.Equal(e2))
		assert.False(t, e1.Equal(e3))
		assert.False(t, e1.Equal(nil))
	})
}

func TestNewWithEntities(t *testing.T) {
	e1 := NewEntity("service", []attribute.KeyValue{
		attribute.String("service.name", "cart"),
	})
	e2 := NewEntity("host", []attribute.KeyValue{
		attribute.String("host.id", "host-123"),
	})

	r := NewWithEntities("https://opentelemetry.io/schemas/1.20.0", e1, e2)
	assert.Equal(t, "https://opentelemetry.io/schemas/1.20.0", r.SchemaURL())
	assert.Len(t, r.Entities(), 2)
	assert.ElementsMatch(t, []attribute.KeyValue{
		attribute.String("service.name", "cart"),
		attribute.String("host.id", "host-123"),
	}, r.Attributes())
	assert.Empty(t, r.UnassociatedAttributes())
}

func TestResourceWithEntitiesOption(t *testing.T) {
	e := NewEntity("service", []attribute.KeyValue{
		attribute.String("service.name", "payment"),
	})
	r, err := New(
		t.Context(),
		WithAttributes(attribute.String("custom.attr", "value")),
		WithEntities(e),
	)
	require.NoError(t, err)

	assert.Len(t, r.Entities(), 1)
	assert.Equal(t, []attribute.KeyValue{attribute.String("custom.attr", "value")}, r.UnassociatedAttributes())
	assert.ElementsMatch(t, []attribute.KeyValue{
		attribute.String("service.name", "payment"),
		attribute.String("custom.attr", "value"),
	}, r.Attributes())
}

func TestResourceMergeWithEntities(t *testing.T) {
	t.Run("DistinctEntityTypes", func(t *testing.T) {
		r1 := NewWithEntities("v1", NewEntity("service", []attribute.KeyValue{
			attribute.String("service.name", "s1"),
		}))
		r2 := NewWithEntities("v1", NewEntity("host", []attribute.KeyValue{
			attribute.String("host.id", "h1"),
		}))

		merged, err := Merge(r1, r2)
		require.NoError(t, err)
		assert.Len(t, merged.Entities(), 2)
		assert.ElementsMatch(t, []attribute.KeyValue{
			attribute.String("service.name", "s1"),
			attribute.String("host.id", "h1"),
		}, merged.Attributes())
	})

	t.Run("SameEntityTypeAndIDMergeDescriptions", func(t *testing.T) {
		e1 := NewEntity(
			"service",
			[]attribute.KeyValue{attribute.String("service.name", "cart")},
			WithEntityDescription(attribute.String("service.version", "1.0.0")),
		)
		e2 := NewEntity(
			"service",
			[]attribute.KeyValue{attribute.String("service.name", "cart")},
			WithEntityDescription(
				attribute.String("service.version", "2.0.0"),
				attribute.String("service.namespace", "prod"),
			),
		)

		r1 := NewWithEntities("", e1)
		r2 := NewWithEntities("", e2)

		merged, err := Merge(r1, r2)
		require.NoError(t, err)
		require.Len(t, merged.Entities(), 1)

		ent := merged.Entities()[0]
		assert.ElementsMatch(t, []attribute.KeyValue{
			attribute.String("service.version", "1.0.0"),
			attribute.String("service.namespace", "prod"),
		}, ent.Description())
	})

	t.Run("ConflictingEntityIdentityDropsNewEntity", func(t *testing.T) {
		e1 := NewEntity("service", []attribute.KeyValue{attribute.String("service.name", "cart")})
		e2 := NewEntity("service", []attribute.KeyValue{attribute.String("service.name", "checkout")})

		r1 := NewWithEntities("", e1)
		r2 := NewWithEntities("", e2)

		merged, err := Merge(r1, r2)
		require.NoError(t, err)
		require.Len(t, merged.Entities(), 1)
		assert.Equal(t, "cart", merged.Entities()[0].ID()[0].Value.AsString())
	})

	t.Run("RawAttributeOverridesEntityAttributeDemotesEntity", func(t *testing.T) {
		e1 := NewEntity(
			"service",
			[]attribute.KeyValue{attribute.String("service.name", "cart")},
			WithEntityDescription(attribute.String("service.version", "1.0.0")),
		)
		r1 := NewWithEntities("", e1)
		r2 := NewSchemaless(attribute.String("service.name", "overridden-cart"))

		merged, err := Merge(r1, r2)
		require.NoError(t, err)

		// The Service entity is dropped because its identifying attribute was overridden.
		assert.Empty(t, merged.Entities())
		assert.ElementsMatch(t, []attribute.KeyValue{
			attribute.String("service.name", "overridden-cart"),
			attribute.String("service.version", "1.0.0"),
		}, merged.UnassociatedAttributes())
	})

	t.Run("SchemaURLConflict", func(t *testing.T) {
		r1 := NewWithEntities("https://schema.1", NewEntity("service", []attribute.KeyValue{
			attribute.String("service.name", "s1"),
		}, WithEntitySchemaURL("https://schema.1")))
		r2 := NewWithEntities("https://schema.2", NewEntity("host", []attribute.KeyValue{
			attribute.String("host.id", "h1"),
		}, WithEntitySchemaURL("https://schema.2")))

		_, err := Merge(r1, r2)
		require.ErrorIs(t, err, ErrSchemaURLConflict)
	})
}

func BenchmarkNewEntity(b *testing.B) {
	id := []attribute.KeyValue{attribute.String("service.name", "cart")}
	desc := []attribute.KeyValue{attribute.String("service.version", "1.0.0")}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewEntity("service", id, WithEntityDescription(desc...))
	}
}

func BenchmarkMergeWithEntities(b *testing.B) {
	r1 := NewWithEntities("v1", NewEntity("service", []attribute.KeyValue{
		attribute.String("service.name", "cart"),
	}))
	r2 := NewWithEntities("v1", NewEntity("host", []attribute.KeyValue{
		attribute.String("host.id", "h-123"),
	}))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Merge(r1, r2)
	}
}
