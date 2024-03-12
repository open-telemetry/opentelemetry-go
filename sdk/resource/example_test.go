// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func ExampleNew() {
	res, err := resource.New(
		context.Background(),
		resource.WithFromEnv(),      // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
		resource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
		resource.WithProcess(),      // Discover and provide process information.
		resource.WithOS(),           // Discover and provide OS information.
		resource.WithContainer(),    // Discover and provide container information.
		resource.WithHost(),         // Discover and provide host information.
		resource.WithAttributes(attribute.String("foo", "bar")), // Add custom resource attributes.
		// resource.WithDetectors(thirdparty.Detector{}), // Bring your own external Detector implementation.
	)
	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		log.Println(err) // Log non-fatal issues.
	} else if err != nil {
		log.Fatalln(err) // The error may be fatal.
	}

	// Now, you can use the resource (e.g. pass it to a tracer or meter provider).
	fmt.Println(res.SchemaURL())
}
