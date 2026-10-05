// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package transform

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
)

func resetAttrCache() {
	attrCacheMu.Lock()
	defer attrCacheMu.Unlock()
	clear(attrCache)
}

func TestCachedAttrsMatchesAttrIter(t *testing.T) {
	resetAttrCache()
	set := attribute.NewSet(
		attribute.String("env", "prod"),
		attribute.Int("n", 7),
		attribute.Bool("ok", true),
	)

	want := AttrIter(set.Iter())
	assert.Equal(t, want, cachedAttrs(set), "cache miss must equal AttrIter")
	assert.Equal(t, want, cachedAttrs(set), "cache hit must equal AttrIter")
}

func TestCachedAttrsHitReusesSlice(t *testing.T) {
	resetAttrCache()
	set := attribute.NewSet(attribute.String("env", "prod"))

	first := cachedAttrs(set)
	second := cachedAttrs(set)
	require.NotEmpty(t, first)
	assert.Same(t, &first[0], &second[0], "cache hit must reuse the slice backing array")
}

func TestCachedAttrsEmptySet(t *testing.T) {
	resetAttrCache()
	assert.Nil(t, cachedAttrs(*attribute.EmptySet()))
}

func TestCachedAttrsClearsWhenFull(t *testing.T) {
	resetAttrCache()
	for i := range maxCachedAttributeSets {
		cachedAttrs(attribute.NewSet(attribute.Int("series", i)))
	}
	require.Len(t, attrCache, maxCachedAttributeSets)

	// Adding another entry clears the cache before storing the new entry.
	cachedAttrs(attribute.NewSet(attribute.Int("series", maxCachedAttributeSets)))
	assert.Len(t, attrCache, 1, "cache should reset to just the new entry")
}

func TestCachedAttrsConcurrent(t *testing.T) {
	resetAttrCache()

	const n = 50
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			set := attribute.NewSet(attribute.Int("series", i))
			assert.Equal(t, AttrIter(set.Iter()), cachedAttrs(set))
		})
	}
	wg.Wait()
}

func TestMaxCachedAttributeSetsFromEnv(t *testing.T) {
	const envKey = "OTEL_GO_X_OTLP_ATTRIBUTE_CACHE_SIZE"

	tests := []struct {
		name string
		env  string
		want int
	}{
		{name: "unset", env: "", want: defaultMaxCachedAttributeSets},
		{name: "valid", env: "500", want: 500},
		{name: "zero", env: "0", want: defaultMaxCachedAttributeSets},
		{name: "negative", env: "-1", want: defaultMaxCachedAttributeSets},
		{name: "not a number", env: "abc", want: defaultMaxCachedAttributeSets},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(envKey, tt.env)
			assert.Equal(t, tt.want, maxCachedAttributeSetsFromEnv())
		})
	}
}
