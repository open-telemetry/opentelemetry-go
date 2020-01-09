// Copyright 2019, OpenTelemetry Authors
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

	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
)

func TestMulitpleGlobalSetScope(t *testing.T) {

	s1 := scope.Empty().WithTracer(trace.NoopTracerSDK{}).WithMeter(metric.NoopMeterSDK{})
	s2 := scope.Empty().WithTracer(trace.NoopTracerSDK{}).WithMeter(metric.NoopMeterSDK{})

	if s1.Provider() == s2.Provider() {
		t.Fatal("impossible test condition")
	}

	global.SetScope(s1)

	defer func() { _ = recover() }()

	global.SetScope(s2)

	got := global.Scope().Provider()
	want := s1.Provider()
	if got != want {
		t.Fatalf("Provider: got %p, want %p\n", got, want)
	}
}
