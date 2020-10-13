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

package foo

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

var (
	lemonsKey = label.Key("ex.com/lemons")
)

// SubOperation is an example to demonstrate the use of named tracer.
// It creates a named tracer with its package path.
func SubOperation(ctx context.Context) error {
	// Using global provider. Alternative is to have application provide a getter
	// for its component to get the instance of the provider.
	tr := global.Tracer("example/namedtracer/foo")

	var span otel.Span
	_, span = tr.Start(ctx, "Sub operation...")
	defer span.End()
	span.SetAttributes(lemonsKey.String("five"))
	span.AddEvent("Sub span event")

	return nil
}
