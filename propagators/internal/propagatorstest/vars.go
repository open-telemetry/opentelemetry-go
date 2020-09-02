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

/*
Package propagatorstest contains common testing elements for the propagators package.
*/
package propagatorstest

import "go.opentelemetry.io/otel/api/trace"

const (
	TraceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736"
	SpanIDStr  = "00f067aa0ba902b7"
)

var (
	TraceID = MustTraceIDFromHex(TraceIDStr)
	SpanID  = MustSpanIDFromHex(SpanIDStr)
)

func MustTraceIDFromHex(s string) (t trace.ID) {
	var err error
	t, err = trace.IDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func MustSpanIDFromHex(s string) (t trace.SpanID) {
	var err error
	t, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}
