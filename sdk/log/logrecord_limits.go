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

package log // import "go.opentelemetry.io/otel/sdk/log"

import "go.opentelemetry.io/otel/sdk/internal/env"

const (
	// DefaultAttributeValueLengthLimit is the default maximum allowed
	// attribute value length, unlimited.
	DefaultAttributeValueLengthLimit = -1

	// DefaultAttributeCountLimit is the default maximum number of attributes
	// a span can have.
	DefaultAttributeCountLimit = 128
)

// LogRecordLimits represents the limits of a span.
type LogRecordLimits struct {
	// AttributeValueLengthLimit is the maximum allowed attribute value length.
	//
	// This limit only applies to string and string slice attribute values.
	// Any string longer than this value will be truncated to this length.
	//
	// Setting this to a negative value means no limit is applied.
	AttributeValueLengthLimit int

	// AttributeCountLimit is the maximum allowed span attribute count. Any
	// attribute added to a span once this limit is reached will be dropped.
	//
	// Setting this to zero means no attributes will be recorded.
	//
	// Setting this to a negative value means no limit is applied.
	AttributeCountLimit int
}

// NewLogRecordLimits returns a LogRecordLimits with all limits set to the value their
// corresponding environment variable holds, or the default if unset.
//
// • AttributeValueLengthLimit: OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT
// (default: unlimited)
//
// • AttributeCountLimit: OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT (default: 128)
func NewLogRecordLimits() LogRecordLimits {
	return LogRecordLimits{
		AttributeValueLengthLimit: env.SpanAttributeValueLength(DefaultAttributeValueLengthLimit),
		AttributeCountLimit:       env.SpanAttributeCount(DefaultAttributeCountLimit),
	}
}
