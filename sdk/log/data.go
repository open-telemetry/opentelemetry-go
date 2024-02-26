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

import (
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

// ResourceRecords is a collection of ScopeMetrics and the associated Resource
// that created them.
type ResourceRecords struct {
	// Resource represents the entity that collected the metrics.
	Resource *resource.Resource
	// ScopeMetrics are the collection of metrics with unique Scopes.
	ScopeRecords []ScopeRecord
}

// ScopeRecords is a collection of log records produced by a Logger.
type ScopeRecord struct {
	// Scope is the Scope that the Logger was created with.
	Scope instrumentation.Scope
	// Records are a list of log records created by the Logger.
	Records []log.Record
}

// Record is log record emitted by the Logger.
type Record struct{}

// ProcessedRecord is a log record used for processing.
type ProcessedRecord struct{}
