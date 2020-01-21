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

// HTTPSupplier is an interface that specifies methods to retrieve and
// store a single value for a key to an associated carrier. It is
// implemented by http.Headers.
type HTTPSupplier interface {
	// Get method retrieves a single value for a given key.
	Get(key string) string
	// Set method stores a single value for a given key. Note that
	// this should not be appending a value to some array, but
	// rather overwrite the old value.
	Set(key string, value string)
}

// HTTPExtractor extracts information from a HTTPSupplier into a
// context.
type HTTPExtractor interface {
	// Extract method retrieves encoded information using supplier
	// from the associated carrier, decodes it and creates a new
	// context containing the decoded information.
	//
	// Information can be a correlation context or a remote span
	// context. In case of span context, the propagator should
	// store it in the context using
	// trace.ContextWithRemoteSpanContext. In case of correlation
	// context, the propagator should use correlation.WithMap to
	// store it in the context.
	Extract(context.Context, HTTPSupplier) context.Context
}

// HTTPInjector injects information into a HTTPSupplier.
type HTTPInjector interface {
	// Inject method retrieves information from the context,
	// encodes it into propagator specific format and then injects
	// the encoded information using supplier into an associated
	// carrier.
	Inject(context.Context, HTTPSupplier)
}

// HTTPPropagator is the interface to inject to and extract from
// HTTPSupplier.
type HTTPPropagator interface {
	HTTPInjector
	HTTPExtractor

	// GetAllKeys returns the HTTP header names used.
	GetAllKeys() []string
}
