// Copyright 2019, OpenTelemetry Authors
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
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/registry"
	"go.opentelemetry.io/api/trace"
)

var (
	lemonsKey  = key.New("ex.com/lemons", registry.WithDescription("A Lemons var"))
)

// SubOperation is simply an example to demonstrate the use of named tracer.
// It creates a named tracer with its package path.
func SubOperation(ctx context.Context) error {
	return trace.GlobalManager().Tracer("example/namedtracer/foo").WithSpan(
		ctx,
		"Sub operation...",
		func(ctx context.Context) error {
			trace.CurrentSpan(ctx).SetAttribute(lemonsKey.String("five"))

			trace.CurrentSpan(ctx).Event(ctx, "Sub span event")

			return nil
		},
	)
}
