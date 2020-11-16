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

package global_test

import (
	"testing"

	"go.opentelemetry.io/otel/internal/global"
)

func TestResetsOfGlobalsPanic(t *testing.T) {
	global.ResetForTest()
	tests := map[string]func(){
		"SetTextMapPropagator": func() {
			global.SetTextMapPropagator(global.TextMapPropagator())
		},
		"SetTracerProvider": func() {
			global.SetTracerProvider(global.TracerProvider())
		},
		"SetMeterProvider": func() {
			global.SetMeterProvider(global.MeterProvider())
		},
	}

	for name, test := range tests {
		shouldPanic(t, name, test)
	}
}

func shouldPanic(t *testing.T, name string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("calling %s with default global did not panic", name)
		}
	}()

	f()
}
