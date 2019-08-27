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
	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/tag"
)

// NoopTextFormatPropagator implements TextFormatPropagator that does nothing.
type NoopTextFormatPropagator struct{}

var _ TextFormatPropagator = NoopTextFormatPropagator{}

// CarrierExtractor returns NoopExtractor
func (ntp NoopTextFormatPropagator) CarrierExtractor(carrier interface{}) Extractor {
	return NoopExtractor{}
}

// CarrierInjector returns NoopInjector
func (ntp NoopTextFormatPropagator) CarrierInjector(carrier interface{}) Injector {
	return NoopInjector{}
}

// NoopInjector implements Injector interface that does nothing.
type NoopInjector struct{}

var _ Injector = NoopInjector{}

func (ni NoopInjector) Inject(sc core.SpanContext, tm tag.Map) {
}

// NoopExtractor implements Extractor interface that does nothing.
type NoopExtractor struct{}

var _ Extractor = NoopExtractor{}

// Extract method always returns Invalid SpanContext and empty tag.Map
func (ne NoopExtractor) Extract() (core.SpanContext, tag.Map) {
	return core.EmptySpanContext(), tag.NewEmptyMap()
}
