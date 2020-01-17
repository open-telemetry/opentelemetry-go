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

package propagation

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/correlation"
)

// NoopTextFormat implements TextFormat that does nothing.
type NoopTextFormat struct{}

var _ TextFormat = NoopTextFormat{}

// Inject does nothing.
func (np NoopTextFormat) Inject(ctx context.Context, supplier Supplier) {
}

// Extract does nothing and returns an empty SpanContext
func (np NoopTextFormat) Extract(ctx context.Context, supplier Supplier) (core.SpanContext, correlation.Map) {
	return core.EmptySpanContext(), correlation.NewEmptyMap()
}

// GetAllKeys returns empty list of strings.
func (np NoopTextFormat) GetAllKeys() []string {
	return []string{}
}
