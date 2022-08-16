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

//go:build go1.18
// +build go1.18

package otlpmetrichttp // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
)

// New constructs a new Exporter and starts it.
func New(ctx context.Context, opts ...Option) (*otlpmetric.Exporter, error) {
	return otlpmetric.New(ctx, NewClient(opts...))
}
