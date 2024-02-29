// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource_test

import (
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/resource"
)

func mockHostIDProvider() {
	resource.SetHostIDProvider(
		func() (string, error) { return "f2c668b579780554f70f72a063dc0864", nil },
	)
}

func mockHostIDProviderWithError() {
	resource.SetHostIDProvider(
		func() (string, error) { return "", assert.AnError },
	)
}

func restoreHostIDProvider() {
	resource.SetDefaultHostIDProvider()
}
