// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetric_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
)

// regex taken from https://github.com/Masterminds/semver/tree/v3.1.1
var versionRegex = regexp.MustCompile(`^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

func TestVersionSemver(t *testing.T) {
	v := otlpmetric.Version()
	assert.NotNil(t, versionRegex.FindStringSubmatch(v), "version is not semver: %s", v)
}
