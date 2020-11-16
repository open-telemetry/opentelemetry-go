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

package otel

import (
	"testing"

	"go.opentelemetry.io/otel/internal/trace/noop"
	"go.opentelemetry.io/otel/trace"
)

type testTracerProvider struct{}

var _ trace.TracerProvider = &testTracerProvider{}

func (*testTracerProvider) Tracer(_ string, _ ...trace.TracerOption) trace.Tracer {
	return noop.Tracer
}

func TestMultipleGlobalTracerProvider(t *testing.T) {
	p1 := testTracerProvider{}
	p2 := trace.NewNoopTracerProvider()
	SetTracerProvider(&p1)
	SetTracerProvider(p2)

	got := GetTracerProvider()
	want := p2
	if got != want {
		t.Fatalf("TracerProvider: got %p, want %p\n", got, want)
	}
}
