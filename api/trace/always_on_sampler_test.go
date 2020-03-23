// Copyright 2020 OpenTelemetry Authors
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

package trace

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/core"
)

func TestAlwaysOnSamplerShouldSample(t *testing.T) {
	gotD := AlwaysOnSampler().ShouldSample(
		core.SpanContext{}, false, core.TraceID{}, core.SpanID{}, "span", SpanKindClient, []core.KeyValue{}, []Link{})
	wantD := Decision{Sampled: true}
	if diff := cmp.Diff(wantD, gotD); diff != "" {
		t.Errorf("Decision: +got, -want%v", diff)
	}
}

func TestAlwaysOnSamplerDescription(t *testing.T) {
	gotDesc := AlwaysOnSampler().Description()
	wantDesc := alwaysOnSamplerDescription
	if diff := cmp.Diff(wantDesc, gotDesc); diff != "" {
		t.Errorf("Description: +got, -want%v", diff)
	}
}
