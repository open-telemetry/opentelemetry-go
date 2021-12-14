// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oc2otel // import "go.opentelemetry.io/otel/bridge/opencensus/internal/oc2otel"

import (
	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/attribute"
)

func Attributes(attr []octrace.Attribute) []attribute.KeyValue {
	otelAttr := make([]attribute.KeyValue, len(attr))
	for i, a := range attr {
		otelAttr[i] = attribute.KeyValue{
			Key:   attribute.Key(a.Key()),
			Value: AttributeValue(a.Value()),
		}
	}
	return otelAttr
}

func AttributeValue(ocval interface{}) attribute.Value {
	switch v := ocval.(type) {
	case bool:
		return attribute.BoolValue(v)
	case int64:
		return attribute.Int64Value(v)
	case float64:
		return attribute.Float64Value(v)
	case string:
		return attribute.StringValue(v)
	default:
		return attribute.StringValue("unknown")
	}
}
