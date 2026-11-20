// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc/internal"
)

var versionRegex = regexp.MustCompile(`^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

func TestVersionSemver(t *testing.T) {
	v := otlpprofilegrpc.Version()
	assert.NotNil(t, versionRegex.FindStringSubmatch(v), "version is not semver: %s", v)
}

func TestVersionMatchesRegisteredReference(t *testing.T) {
	assert.Equal(t, internal.Version, otlpprofilegrpc.Version(),
		"Version() must report the module's registered version reference")
}
