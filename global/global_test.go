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

	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/global"
)

type TestProvider1 struct {
}

var _ trace.Provider = &TestProvider1{}

func (tp *TestProvider1) GetTracer(name string) trace.Tracer {
	return &trace.NoopTracer{}
}

func TestMulitpleGlobalProvider(t *testing.T) {

	p1 := TestProvider1{}
	p2 := trace.NoopTraceProvider{}
	global.SetTraceProvider(&p1)
	global.SetTraceProvider(&p2)

	got := global.TraceProvider()
	want := &p2
	if got != want {
		t.Fatalf("Provider: got %p, want %p\n", got, want)
	}
}
