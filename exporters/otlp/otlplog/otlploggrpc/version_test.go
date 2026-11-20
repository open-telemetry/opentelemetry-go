// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal"
)

// versionRegex is based on https://github.com/Masterminds/semver/tree/v3.1.1.
var versionRegex = regexp.MustCompile(`^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

func TestVersionSemver(t *testing.T) {
	v := Version()
	assert.NotNil(t, versionRegex.FindStringSubmatch(v), "version is not semver: %s", v)
}

// TestVersionMatchesRegisteredReference guards against the public Version
// drifting from the module release again: versions.yaml publishes this module's
// version through ./internal/version.go, so that is the value the release
// tooling updates and the value the User-Agent must report.
func TestVersionMatchesRegisteredReference(t *testing.T) {
	assert.Equal(t, internal.Version, Version(),
		"Version() must report the module's registered version reference")
}
