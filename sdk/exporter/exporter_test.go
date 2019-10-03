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

package exporter

import (
	"context"
	"testing"
)

type testSyncExporter struct {
	spans []string
}

func (t *testSyncExporter) ExportSpan(ctx context.Context, s string) {
	t.spans = append(t.spans, s)
}

func TestRegisterUnregister(t *testing.T) {
	var te testSyncExporter

	Register(ExporterTypeSync, &te)

	e := exporters.Load().(exportersMap)
	if len(e) != 1 {
		t.Errorf("Expected 1 exporter type. Got %d", len(e))
	}

	if len(e[0]) != 1 {
		t.Errorf("Expected 1 registered exporter. Got %d", len(e[0]))
	}

	Unregister(&te)

	e = exporters.Load().(exportersMap)
	if len(e) != 1 {
		t.Errorf("Expected 1 exporter type. Got %d", len(e))
	}

	if len(e[0]) != 0 {
		t.Errorf("Expected 0 registered exporters. Got %d", len(e[0]))
	}
}

func TestLoad(t *testing.T) {
	// Empty the exporters map so we can test loading an exporter for the first
	// time
	exporters.Store(make(exportersMap))

	var te testSyncExporter
	e := Load(ExporterTypeSync)

	if len(e) != 0 {
		t.Errorf("Expected no exporters. Got %d", len(e))
	}

	Register(ExporterTypeSync, &te)

	e = Load(ExporterTypeSync)
	if len(e) != 1 {
		t.Errorf("Expected 1 exporter. Got %d", len(e))
	}
}
