// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc_test

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc"
	"go.opentelemetry.io/otel/sdk/profile"
)

func Example() {
	ctx := context.Background()
	exp, err := otlpprofilegrpc.New(ctx, otlpprofilegrpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	reader := profile.NewPeriodicReader(exp)
	provider := profile.NewProfilerProvider(
		profile.WithProducer(profile.NewRuntimeProducer()),
		profile.WithReader(reader),
	)
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	if err := provider.ForceFlush(ctx); err != nil {
		panic(err)
	}
}
