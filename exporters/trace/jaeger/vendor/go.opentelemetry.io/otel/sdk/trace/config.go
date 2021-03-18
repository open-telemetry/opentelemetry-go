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

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"go.opentelemetry.io/otel/sdk/resource"
)

// Config represents the global tracing configuration.
type Config struct {
	// DefaultSampler is the default sampler used when creating new spans.
	DefaultSampler Sampler

	// IDGenerator is for internal use only.
	IDGenerator IDGenerator

	// SpanLimits used to limit the number of attributes, events and links to a span.
	SpanLimits SpanLimits

	// Resource contains attributes representing an entity that produces telemetry.
	Resource *resource.Resource
}

// SpanLimits represents the limits of a span.
type SpanLimits struct {
	// AttributeCountLimit is the maximum allowed span attribute count.
	AttributeCountLimit int

	// EventCountLimit is the maximum allowed span event count.
	EventCountLimit int

	// LinkCountLimit is the maximum allowed span link count.
	LinkCountLimit int

	// AttributePerEventCountLimit is the maximum allowed attribute per span event count.
	AttributePerEventCountLimit int

	// AttributePerLinkCountLimit is the maximum allowed attribute per span link count.
	AttributePerLinkCountLimit int
}

const (
	// DefaultAttributeCountLimit is the default maximum allowed span attribute count.
	DefaultAttributeCountLimit = 128

	// DefaultEventCountLimit is the default maximum allowed span event count.
	DefaultEventCountLimit = 128

	// DefaultLinkCountLimit is the default maximum allowed span link count.
	DefaultLinkCountLimit = 128

	// DefaultAttributePerEventCountLimit is the default maximum allowed attribute per span event count.
	DefaultAttributePerEventCountLimit = 128

	// DefaultAttributePerLinkCountLimit is the default maximum allowed attribute per span link count.
	DefaultAttributePerLinkCountLimit = 128
)
