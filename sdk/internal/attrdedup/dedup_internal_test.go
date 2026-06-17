// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrdedup

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestInvalidArrayStorage(t *testing.T) {
	if got := arrayLen("invalid", valueType); got != 0 {
		t.Fatalf("arrayLen() = %d, want 0", got)
	}

	if got := arrayAt[attribute.Value]("invalid", valueType, 0); got.Type() != attribute.EMPTY {
		t.Fatalf("arrayAt() = %v, want empty value", got)
	}
	if got := arrayAt[attribute.Value](
		[1]attribute.Value{attribute.StringValue("value")},
		keyValueType,
		0,
	); got.Type() != attribute.EMPTY {
		t.Fatalf("arrayAt() = %v, want empty value", got)
	}
	if got := arrayAt[attribute.Value](
		[1]attribute.Value{attribute.StringValue("value")},
		valueType,
		-1,
	); got.Type() != attribute.EMPTY {
		t.Fatalf("arrayAt() = %v, want empty value", got)
	}
}
