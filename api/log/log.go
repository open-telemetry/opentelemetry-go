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

package log

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/scope"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
)

type (
	Interface interface {
		Log(ctx context.Context, msg string, fields ...core.KeyValue)
		Logf(ctx context.Context, fmt string, args ...interface{})
	}

	Logger struct {
		scope.Scope
	}
)

func With(scope scope.Scope) Logger {
	return Logger{scope}
}

func Log(ctx context.Context, msg string, fields ...core.KeyValue) {
	With(scope.Active(ctx)).Log(ctx, msg, fields...)
}

func Logf(ctx context.Context, fmt string, args ...interface{}) {
	With(scope.Active(ctx)).Logf(ctx, fmt, args...)
}

func (l Logger) Log(ctx context.Context, msg string, fields ...core.KeyValue) {
	observer.Record(observer.Event{
		Type:       observer.LOG_EVENT,
		Scope:      l.ScopeID(),
		String:     msg,
		Attributes: fields,
		Context:    ctx,
	})
}

func (l Logger) Logf(ctx context.Context, fmt string, args ...interface{}) {
	observer.Record(observer.Event{
		Type:      observer.LOGF_EVENT,
		Scope:     l.ScopeID(),
		String:    fmt,
		Arguments: args,
		Context:   ctx,
	})
}
