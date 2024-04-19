// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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

func TestAttributesFromMap(t *testing.T) {
	in := map[string]interface{}{
		"bool":    true,
		"int64":   int64(49),
		"float64": float64(1.618),
		"key":     "val",
	}

	want := []attribute.KeyValue{
		attribute.Bool("bool", true),
		attribute.Int64("int64", 49),
		attribute.Float64("float64", 1.618),
		attribute.String("key", "val"),
	}
	got := AttributesFromMap(in)

	gotAttributeSet := attribute.NewSet(got...)
	wantAttributeSet := attribute.NewSet(want...)
	if !gotAttributeSet.Equals(&wantAttributeSet) {
		t.Errorf("Attributes conversion want %v, got %v", wantAttributeSet.Encoded(attribute.DefaultEncoder()), gotAttributeSet.Encoded(attribute.DefaultEncoder()))
	}
}

func TestAttributeValueUnknown(t *testing.T) {
	got := AttributeValue([]byte{})
	if got != attribute.StringValue("unknown") {
		t.Errorf("AttributeValue of unknown wrong: %#v", got)
	}
}
