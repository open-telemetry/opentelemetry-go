// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

// Severity represents a log record severity (also known as log level).
// Smaller numerical values correspond to less severe log records (such as debug events),
// larger numerical values correspond to more severe log records (such as errors and critical events).
type Severity int

// Severity values defined by OpenTelemetry.
const (
	// A fine-grained debugging log record. Typically disabled in default configurations.
	SeverityTrace1 Severity = iota + 1
	SeverityTrace2
	SeverityTrace3
	SeverityTrace4

	// A debugging log record.
	SeverityDebug1
	SeverityDebug2
	SeverityDebug3
	SeverityDebug4

	// An informational log record. Indicates that an event happened.
	SeverityInfo1
	SeverityInfo2
	SeverityInfo3
	SeverityInfo4

	// A warning log record. Not an error but is likely more important than an informational event.
	SeverityWarn1
	SeverityWarn2
	SeverityWarn3
	SeverityWarn4

	// An error log record. Something went wrong.
	SeverityError1
	SeverityError2
	SeverityError3
	SeverityError4

	// A fatal log record such as application or system crash.
	SeverityFatal1
	SeverityFatal2
	SeverityFatal3
	SeverityFatal4

	SeverityTrace = SeverityTrace1
	SeverityDebug = SeverityDebug1
	SeverityInfo  = SeverityInfo1
	SeverityWarn  = SeverityWarn1
	SeverityError = SeverityError1
	SeverityFatal = SeverityFatal1
)
