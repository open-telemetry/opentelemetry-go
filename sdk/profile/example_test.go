// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile_test

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/profile"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

func ExampleNewProfilerProvider() {
	ctx := context.Background()
	reader := profile.NewManualReader()
	pp := profile.NewProfilerProvider(
		profile.WithProducer(profile.NewRuntimeProducer(
			profile.WithoutAllocsProfile(),
			profile.WithoutHeapProfile(),
			profile.WithGoroutineProfile(),
		)),
		profile.WithReader(reader),
	)
	defer func() { _ = pp.Shutdown(ctx) }()

	var rp profiledata.ResourceProfiles
	if err := reader.Collect(ctx, &rp); err != nil {
		panic(err)
	}
	fmt.Println(len(rp.ScopeProfiles) > 0)
	// Output: true
}
