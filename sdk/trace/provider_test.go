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

package trace

import (
	"context"
	"testing"

	export "go.opentelemetry.io/otel/sdk/export/trace"
)

type basicSpanProcesor struct {
	running bool
}

func (t *basicSpanProcesor) Shutdown(context.Context) error {
	t.running = false
	return nil
}

func (t *basicSpanProcesor) OnStart(s *export.SpanData) {}
func (t *basicSpanProcesor) OnEnd(s *export.SpanData)   {}
func (t *basicSpanProcesor) ForceFlush()                {}

func TestShutdownTraceProvider(t *testing.T) {
	stp := NewTracerProvider()
	sp := &basicSpanProcesor{}
	stp.RegisterSpanProcessor(sp)

	sp.running = true

	_ = stp.Shutdown(context.Background())

	if sp.running != false {
		t.Errorf("Error shutdown basicSpanProcesor\n")
	}
}
