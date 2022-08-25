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

package oc2otel

import (
	"testing"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/attribute"
)

func TestAttributes(t *testing.T) {
	in := []octrace.Attribute{
		octrace.BoolAttribute("bool", true),
		octrace.Int64Attribute("int64", 49),
		octrace.Float64Attribute("float64", 1.618),
		octrace.StringAttribute("key", "val"),
	}

	want := []attribute.KeyValue{
		attribute.Bool("bool", true),
		attribute.Int64("int64", 49),
		attribute.Float64("float64", 1.618),
		attribute.String("key", "val"),
	}
	got := Attributes(in)

	if len(got) != len(want) {
		t.Errorf("Attributes conversion failed: want %#v, got %#v", want, got)
	}
	for i := range got {
		if g, w := got[i], want[i]; g != w {
			t.Errorf("Attributes conversion: want %#v, got %#v", w, g)
		}
	}
}

func TestAttributeValueUnknown(t *testing.T) {
	got := AttributeValue([]byte{})
	if got != attribute.StringValue("unknown") {
		t.Errorf("AttributeValue of unknown wrong: %#v", got)
	}
}
