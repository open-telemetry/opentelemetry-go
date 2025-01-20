// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import "go.opentelemetry.io/otel/attribute"

// ConvertAttributeValue converts [attribute.Value) to [Value].
func ConvertAttributeValue(v attribute.Value) Value {
	return Value{}
}

// ConvertAttributeKeyValue converts [attribute.KeyValue) to [KeyValue].
func ConvertAttributeKeyValue(kv attribute.KeyValue) KeyValue {
	return KeyValue{}
}
