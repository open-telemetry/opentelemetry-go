// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log // import "go.opentelemetry.io/otel/log"

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// Record represents a log record.
type Record struct {
	Timestamp         time.Time
	ObservedTimestamp time.Time
	Severity          Severity
	SeverityText      string
	Body              string
	Attributes        []attribute.KeyValue
}

// Severity represents a log record severity.
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
)

const (
	SeverityTrace = SeverityTrace1
	SeverityDebug = SeverityDebug1
	SeverityInfo  = SeverityInfo1
	SeverityWarn  = SeverityWarn1
	SeverityError = SeverityError1
	SeverityFatal = SeverityFatal1
)
