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
)

// NoopHTTPPropagator implements HTTPPropagator that does nothing.
type NoopHTTPPropagator struct{}

var _ HTTPPropagator = NoopHTTPPropagator{}

// Inject does nothing.
func (NoopHTTPPropagator) Inject(ctx context.Context, supplier HTTPSupplier) {
}

// Extract does nothing.
func (NoopHTTPPropagator) Extract(ctx context.Context, supplier HTTPSupplier) context.Context {
	return ctx
}

// GetAllKeys returns empty list of strings.
func (NoopHTTPPropagator) GetAllKeys() []string {
	return []string{}
}
