// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestWithUnsafeAttributes(t *testing.T) {
	kvs := []attribute.KeyValue{attribute.String("A", "B")}
	opt := WithUnsafeAttributes(kvs...)

	unsafeOpt, ok := opt.(*unsafeAttributesOption)
	if !ok {
		t.Fatalf("expected *unsafeAttributesOption")
	}

	attrs := unsafeOpt.RawAttributes()
	if len(attrs) != 1 {
		t.Errorf("expected 1 attribute, got %d", len(attrs))
	}
	if attrs[0].Key != "A" || attrs[0].Value.AsString() != "B" {
		t.Errorf("expected attribute A='B', got %v", attrs[0])
	}
}
