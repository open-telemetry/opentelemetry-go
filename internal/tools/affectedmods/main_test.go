// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"sort"
	"testing"
)

func TestOwningModule(t *testing.T) {
	// Mirrors the shape of `make print-sharded-benchmarks`: "." for the
	// root module, "./<dir>" for nested ones, and "./<mod>/." marking the
	// catch-all sub-shard of a split module. Sorted longest-first, as
	// owningModule expects.
	mods := []string{
		".",
		"./sdk",
		"./sdk/log",
		"./sdk/log/logtest",
		"./sdk/metric/.",
		"./sdk/metric/internal",
		"./bridge/opencensus",
		"./bridge/opencensus/test",
	}
	sort.Slice(mods, func(i, j int) bool { return len(mods[i]) > len(mods[j]) })

	tests := []struct {
		name    string
		path    string
		modules []string
		want    string
	}{
		{
			name: "direct hit on nested module",
			path: "sdk/log/exporter.go", modules: mods,
			want: "./sdk/log",
		},
		{
			name: "nested file picks longest prefix",
			path: "sdk/log/logtest/recorder.go", modules: mods,
			want: "./sdk/log/logtest",
		},
		{
			name: "shorter prefix when no longer match exists",
			path: "sdk/trace/span.go", modules: mods,
			want: "./sdk",
		},
		{
			name: "exact module path (no trailing slash)",
			path: "sdk", modules: mods,
			want: "./sdk",
		},
		{
			name: "top-level file falls through to root",
			path: "handler.go", modules: mods,
			want: ".",
		},
		{
			name: "nested file with no nested-module match falls through to root",
			path: "propagation/baggage.go", modules: mods,
			want: ".",
		},
		{
			name: "module-name prefix without slash boundary falls through to root",
			path: "sdk-extra/x.go", modules: mods,
			want: ".",
		},
		{
			name: "deep nesting under unrelated tree",
			path: "bridge/opencensus/test/exporter.go", modules: mods,
			want: "./bridge/opencensus/test",
		},
		{
			name: "split-module catch-all matches files at module root",
			path: "sdk/metric/meter.go", modules: mods,
			want: "./sdk/metric/.",
		},
		{
			name: "split-module catch-all matches files in non-internal subdir",
			path: "sdk/metric/exemplar/x.go", modules: mods,
			want: "./sdk/metric/.",
		},
		{
			name: "split-module sub-shard wins for its subtree",
			path: "sdk/metric/internal/x.go", modules: mods,
			want: "./sdk/metric/internal",
		},
		{
			name: "fallback empty when root not in list and no nested match",
			path: "handler.go", modules: []string{"./sdk"},
			want: "",
		},
		{
			name: "empty modules returns empty",
			path: "sdk/log/exporter.go", modules: nil,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := owningModule(tt.path, tt.modules)
			if got != tt.want {
				t.Errorf("owningModule(%q, %v) = %q, want %q",
					tt.path, tt.modules, got, tt.want)
			}
		})
	}
}
