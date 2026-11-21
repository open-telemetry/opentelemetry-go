// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetrichttp

import (
	"os/exec"
	"strings"
	"testing"
)

// TestNoGRPCDependency ensures the OTLP/HTTP exporter packages do not depend
// on gRPC. The HTTP transport must not pull in google.golang.org/grpc.
//
// See https://github.com/open-telemetry/opentelemetry-go/issues/2579.
func TestNoGRPCDependency(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("go command not available: %v", err)
	}
	out, err := exec.CommandContext(t.Context(), goBin, "list", "-deps", "./...").Output()
	if err != nil {
		t.Fatalf("go list -deps ./... failed: %v", err)
	}
	for line := range strings.SplitSeq(string(out), "\n") {
		if dep := strings.TrimSpace(line); strings.Contains(dep, "google.golang.org/grpc") {
			t.Errorf("OTLP/HTTP exporter must not depend on gRPC, found %q in go list -deps ./...", dep)
		}
	}
}
