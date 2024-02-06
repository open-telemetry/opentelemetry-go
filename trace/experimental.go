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

package trace // import "go.opentelemetry.io/otel/trace"

// ExperimentalSpan exposes the same methods as trace.Span and additionally supports
// select features marked experimental by the specification and thus subject to breaking
// changes without warning.
//
// Warning: methods may be added to or removed from this interface in minor releases.
type ExperimentalSpan interface {
	Span

	// AddLinks adds links to the existing trace.Span after creation.
	AddLinks(links ...Link)
}
