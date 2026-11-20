// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilehttp_test

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp"
	"go.opentelemetry.io/otel/sdk/profile"
)

func Example() {
	ctx := context.Background()
	exp, err := otlpprofilehttp.New(ctx, otlpprofilehttp.WithInsecure())
	if err != nil {
		panic(err)
	}

	reader := profile.NewPeriodicReader(exp)
	provider := profile.NewProfilerProvider(
		profile.WithProducer(profile.NewRuntimeProducer()),
		profile.WithReader(reader),
	)
	// Application runs while PeriodicReader collects and exports profiles.
	// Shutdown performs a final flush before stopping the exporter.
	if err := provider.Shutdown(ctx); err != nil {
		panic(err)
	}
}
