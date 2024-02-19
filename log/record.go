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

package log // import "go.opentelemetry.io/otel/log"

import "time"

// Record represents a log record.
type Record struct{} // TODO (#4913): implement.

// Timestamp returns the time when the log record occurred.
func (r *Record) Timestamp() time.Time { return time.Time{} } // TODO (#4913): implement.

// SetTimestamp sets the time when the log record occurred.
func (r *Record) SetTimestamp(t time.Time) {} // TODO (#4913): implement.

// ObservedTimestamp returns the time when the log record was observed.
func (r *Record) ObservedTimestamp() time.Time { return time.Time{} } // TODO (#4913): implement.

// SetObservedTimestamp sets the time when the log record was observed.
func (r *Record) SetObservedTimestamp(t time.Time) {} // TODO (#4913): implement.

// Severity returns the [Severity] of the log record.
func (r *Record) Severity() Severity { return 0 } // TODO (#4913): implement.

// SetSeverity sets the [Severity] level of the log record.
func (r *Record) SetSeverity(level Severity) {} // TODO (#4913): implement.

// SeverityText returns severity (also known as log level) text. This is the
// original string representation of the severity as it is known at the source.
func (r *Record) SeverityText() string { return "" } // TODO (#4913): implement.

// SetSeverityText sets severity (also known as log level) text. This is the
// original string representation of the severity as it is known at the source.
func (r *Record) SetSeverityText(text string) {} // TODO (#4913): implement.

// Body returns the body of the log record.
func (r *Record) Body() Value { return Value{} } // TODO (#4913): implement.

// SetBody sets the body of the log record.
func (r *Record) SetBody(v Value) {} // TODO (#4913): implement.

// WalkAttributes walks all attributes the log record holds by calling f for
// each on each [KeyValue] in the [Record]. Iteration stops if f returns false.
func (r *Record) WalkAttributes(f func(KeyValue) bool) {} // TODO (#4913): implement.

// AddAttributes adds attributes to the log record.
func (r *Record) AddAttributes(attributes ...KeyValue) {} // TODO (#4913): implement.

// AttributesLen returns the number of attributes in the log record.
func (r *Record) AttributesLen() int { return 0 } // TODO (#4913): implement.
