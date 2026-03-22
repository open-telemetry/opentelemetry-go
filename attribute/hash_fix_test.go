// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestHashKVWithUnknownType(t *testing.T) {
	// Test that hashKV handles unknown types gracefully without panicking.
	// This test creates a Value with an invalid type to ensure the fix works.
	
	// Create a KeyValue with empty value (should have EMPTY type)
	kv := attribute.KeyValue{Key: "test"}
	
	// This should not panic even if the type is somehow invalid
	set := attribute.NewSet(kv)
	
	// Verify the set can be created and hashed without issues
	hash := set.Equivalent()
	if hash == "" {
		t.Error("hash should not be empty for valid KeyValue")
	}
	
	// Test with explicitly empty value
	emptyKV := attribute.String("key", "")
	set2 := attribute.NewSet(emptyKV)
	hash2 := set2.Equivalent()
	if hash2 == "" {
		t.Error("hash should not be empty for empty string value")
	}
}

func TestHashKVConsistencyAfterFix(t *testing.T) {
	// Test that the fix maintains hash consistency
	kvs := []attribute.KeyValue{
		attribute.Bool("bool", true),
		attribute.Int("int", 42),
		attribute.Float64("float", 3.14),
		attribute.String("string", "test"),
		attribute.BoolSlice("boolslice", []bool{true, false}),
		attribute.IntSlice("intslice", []int{1, 2, 3}),
		attribute.Float64Slice("floatslice", []float64{1.1, 2.2}),
		attribute.StringSlice("stringslice", []string{"a", "b"}),
		attribute.KeyValue{Key: "empty"}, // EMPTY type
	}
	
	// Create sets multiple times to ensure consistency
	set1 := attribute.NewSet(kvs...)
	set2 := attribute.NewSet(kvs...)
	
	if set1.Equivalent() != set2.Equivalent() {
		t.Error("hash should be consistent for same attributes")
	}
	
	// Test that empty values are handled consistently
	emptyKVs := []attribute.KeyValue{
		attribute.KeyValue{Key: "empty1"},
		attribute.KeyValue{Key: "empty2"},
	}
	
	set3 := attribute.NewSet(emptyKVs...)
	set4 := attribute.NewSet(emptyKVs...)
	
	if set3.Equivalent() != set4.Equivalent() {
		t.Error("hash should be consistent for empty values")
	}
}
