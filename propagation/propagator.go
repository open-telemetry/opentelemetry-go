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

// Package tracecontext contains HTTP propagator for TraceContext standard.
// See https://github.com/w3c/distributed-tracing for more information.
package propagation // import "go.opentelemetry.io/propagation"

import (
	"net/http"

	"go.opentelemetry.io/api/trace"
)

// HTTPPropagator is an interface that specifies methods to create Extractor
// and Injector objects for an http request. Typically, an http plugin uses
// this interface to allow user to configure appropriate propagators.
type HTTPPropagator interface {
	Extractor(req *http.Request) trace.Extractor
	Injector(req *http.Request) trace.Injector
}
