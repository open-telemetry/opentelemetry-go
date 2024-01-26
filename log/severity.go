// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate stringer -type=Severity -linecomment

package log // import "go.opentelemetry.io/otel/log"

// Severity represents a log record severity (also known as log level).
// Smaller numerical values correspond to less severe log records (such as debug events),
// larger numerical values correspond to more severe log records (such as errors and critical events).
type Severity int

// Severity values defined by OpenTelemetry.
const (
	// A fine-grained debugging log record. Typically disabled in default configurations.
	SeverityTrace1 Severity = iota + 1 // TRACE
	SeverityTrace2                     // TRACE2
	SeverityTrace3                     // TRACE3
	SeverityTrace4                     // TRACE4

	// A debugging log record.
	SeverityDebug1 // DEBUG
	SeverityDebug2 // DEBUG1
	SeverityDebug3 // DEBUG2
	SeverityDebug4 // DEBUG3

	// An informational log record. Indicates that an event happened.
	SeverityInfo1 // INFO
	SeverityInfo2 // INFO1
	SeverityInfo3 // INFO2
	SeverityInfo4 // INFO3

	// A warning log record. Not an error but is likely more important than an informational event.
	SeverityWarn1 // WARN
	SeverityWarn2 // WARN2
	SeverityWarn3 // WARN3
	SeverityWarn4 // WARN4

	// An error log record. Something went wrong.
	SeverityError1 // ERROR
	SeverityError2 // ERROR2
	SeverityError3 // ERROR3
	SeverityError4 // ERROR4

	// A fatal log record such as application or system crash.
	SeverityFatal1 // FATAL
	SeverityFatal2 // FATAL2
	SeverityFatal3 // FATAL3
	SeverityFatal4 // FATAL4

	SeverityTrace = SeverityTrace1
	SeverityDebug = SeverityDebug1
	SeverityInfo  = SeverityInfo1
	SeverityWarn  = SeverityWarn1
	SeverityError = SeverityError1
	SeverityFatal = SeverityFatal1
)
