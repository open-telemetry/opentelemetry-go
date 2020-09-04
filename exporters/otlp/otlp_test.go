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

package otlp

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExporterShutdownHonorsTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Set the reconnect period to be longer than this test should last so the
	// exporter will be in a continuous state of trying to reconnect and will
	// not shutdown.
	e := NewUnstartedExporter(WithReconnectionPeriod(1 * time.Minute))
	if err := e.Start(); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	var innerCancel context.CancelFunc
	ctx, innerCancel = context.WithCancel(ctx)
	innerCancel()
	if err := e.Shutdown(ctx); err == nil {
		t.Error("expected context canceled error, got nil")
	} else if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context canceled error, got %v", err)
	}
}
