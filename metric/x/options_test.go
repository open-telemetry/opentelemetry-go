// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestWithUnsafeAttributes(t *testing.T) {
	kvs := []attribute.KeyValue{attribute.String("A", "B")}
	opt := WithUnsafeAttributes(kvs...)

	unsafeOpt, ok := opt.(*unsafeAttributesOption)
	assert.True(t, ok, "expected *unsafeAttributesOption")
	assert.Equal(t, kvs, unsafeOpt.RawAttributes(), "expected stored attributes to match")
}
