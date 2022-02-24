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

const (
	// DefaultAttributeValueLengthLimit is the default maximum allowed
	// attribute value length, unlimited.
	DefaultAttributeValueLengthLimit = -1

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

// SpanLimits represents the limits of a span.
type SpanLimits struct {
	// AttributeValueLengthLimit is the maximum allowed attribute value length.
	//
	// This limit only applies to string and string slice attribute values.
	//
	// Setting this to zero means no string or string slice attributes we be
	// recorded.
	//
	// Setting this value a negative value means no limit is applied.
	AttributeValueLengthLimit int

	// AttributeCountLimit is the maximum allowed span attribute count.
	//
	// Setting this to zero means no attributes we be recorded.
	//
	// Setting this value a negative value means no limit is applied.
	AttributeCountLimit int

	// EventCountLimit is the maximum allowed span event count.
	//
	// Setting this to zero means no events we be recorded.
	//
	// Setting this value a negative value means no limit is applied.
	EventCountLimit int

	// LinkCountLimit is the maximum allowed span link count.
	//
	// Setting this to zero means no links we be recorded.
	//
	// Setting this value a negative value means no limit is applied.
	LinkCountLimit int

	// AttributePerEventCountLimit is the maximum allowed attribute per span event count.
	//
	// Setting this to zero means no attributes will be recorded for events.
	//
	// Setting this value a negative value means no limit is applied.
	AttributePerEventCountLimit int

	// AttributePerLinkCountLimit is the maximum allowed attribute per span link count.
	//
	// Setting this to zero means no attributes will be recorded for links.
	//
	// Setting this value a negative value means no limit is applied.
	AttributePerLinkCountLimit int
}

// NewSpanLimits returns a SpanLimits with all limits set to defaults.
func NewSpanLimits() SpanLimits {
	return SpanLimits{
		AttributeValueLengthLimit:   DefaultAttributeValueLengthLimit,
		AttributeCountLimit:         DefaultAttributeCountLimit,
		EventCountLimit:             DefaultEventCountLimit,
		LinkCountLimit:              DefaultLinkCountLimit,
		AttributePerEventCountLimit: DefaultAttributePerEventCountLimit,
		AttributePerLinkCountLimit:  DefaultAttributePerLinkCountLimit,
	}
}

// newEnvSpanLimits returns a SpanLimits with all limits set to the values
// defined by related environment variables, or defaults otherwise.
func newEnvSpanLimits() SpanLimits {
	return SpanLimits{
		AttributeValueLengthLimit:   env.SpanAttributeValueLength(DefaultAttributeValueLengthLimit),
		AttributeCountLimit:         env.SpanAttributeCount(DefaultAttributeCountLimit),
		EventCountLimit:             env.SpanEventCount(DefaultEventCountLimit),
		LinkCountLimit:              env.SpanLinkCount(DefaultLinkCountLimit),
		AttributePerEventCountLimit: env.SpanEventAttributeCount(DefaultAttributePerEventCountLimit),
		AttributePerLinkCountLimit:  env.SpanLinkAttributeCount(DefaultAttributePerLinkCountLimit),
	}
}
