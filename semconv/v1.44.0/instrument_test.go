// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package semconv

import "testing"

func TestMetricInstrumentSelection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, got, want string
	}{
		{"double", "Float64Gauge", "Float64Gauge"},
		{"int", "Int64Gauge", "Int64Gauge"},
		{"unannotated_integer", "Int64Counter", "Int64Counter"},
		{"unannotated_seconds", "Float64Histogram", "Float64Histogram"},
		{"unannotated_milliseconds", "Float64Histogram", "Float64Histogram"},
		{"empty_annotations", "Float64UpDownCounter", "Float64UpDownCounter"},
		{"empty_code_generation", "Int64Gauge", "Int64Gauge"},
		{"float_override", "Float64ObservableCounter", "Float64ObservableCounter"},
		{"int_override", "Int64ObservableUpDownCounter", "Int64ObservableUpDownCounter"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if test.got != test.want {
				t.Errorf("instrument = %q, want %q", test.got, test.want)
			}
		})
	}
}
