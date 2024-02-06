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

package resource_test

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func Example() {
	res, err := resource.New(context.Background(),
		resource.WithFromEnv(),      // Pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
		resource.WithTelemetrySDK(), // Provide information about the OpenTelemetry SDK used.
		resource.WithProcess(),      // Discover and provide process information.
		resource.WithOS(),           // Discover and provide OS information.
		resource.WithContainer(),    // Discover and provide container information.
		resource.WithHost(),         // Discover and provide information.
		resource.WithAttributes(attribute.String("foo", "bar")), // Add custom resource attributes.
		// resource.WithDetectors(thirdparty.Detector{}), // Bring your own external Detector implementation.
	)
	if err != nil {
		log.Println(err) // Log issues during resource creation. Note that resource.New still returns a resource.
	}

	// Now, you can use the resource and e.g. pass it to a tracer or meter provider.
	fmt.Println(res.SchemaURL())

	// Output: https://opentelemetry.io/schemas/1.24.0
}
