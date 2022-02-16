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
import "go.opentelemetry.io/otel/sdk/internal/env"

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

func (sl *SpanLimits) ensureDefault() {
	if sl.EventCountLimit <= 0 {
		sl.EventCountLimit = DefaultEventCountLimit
	}
	if sl.AttributeCountLimit <= 0 {
		sl.AttributeCountLimit = DefaultAttributeCountLimit
	}
	if sl.LinkCountLimit <= 0 {
		sl.LinkCountLimit = DefaultLinkCountLimit
	}
	if sl.AttributePerEventCountLimit <= 0 {
		sl.AttributePerEventCountLimit = DefaultAttributePerEventCountLimit
	}
	if sl.AttributePerLinkCountLimit <= 0 {
		sl.AttributePerLinkCountLimit = DefaultAttributePerLinkCountLimit
	}
}

func (sl *SpanLimits) parsePotentialEnvConfigs() {
	sl.AttributeCountLimit = env.IntEnvOr(env.SpanAttributesCountKey, sl.AttributeCountLimit)
	sl.LinkCountLimit = env.IntEnvOr(env.SpanLinkCountKey, sl.LinkCountLimit)
	sl.EventCountLimit = env.IntEnvOr(env.SpanEventCountKey, sl.EventCountLimit)
}

const (
	// DefaultAttributeCountLimit is the default maximum allowed span attribute count.
	// If not specified via WithSpanLimits, will try to retrieve the value from
	// environment variable `OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT`.
	// If Invalid value (negative or zero) is provided, the default value 128 will be used.
	DefaultAttributeCountLimit = 128

	// DefaultEventCountLimit is the default maximum allowed span event count.
	// If not specified via WithSpanLimits, will try to retrieve the value from
	// environment variable `OTEL_SPAN_EVENT_COUNT_LIMIT`.
	// If Invalid value (negative or zero) is provided, the default value 128 will be used.
	DefaultEventCountLimit = 128

	// DefaultLinkCountLimit is the default maximum allowed span link count.
	// If the value is not specified via WithSpanLimits, will try to retrieve the value from
	// environment variable `OTEL_SPAN_LINK_COUNT_LIMIT`.
	// If Invalid value (negative or zero) is provided, the default value 128 will be used.
	DefaultLinkCountLimit = 128

	// DefaultAttributePerEventCountLimit is the default maximum allowed attribute per span event count.
	DefaultAttributePerEventCountLimit = 128

	// DefaultAttributePerLinkCountLimit is the default maximum allowed attribute per span link count.
	DefaultAttributePerLinkCountLimit = 128
)
