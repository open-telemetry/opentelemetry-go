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

package httptrace

import (
	"context"
	"encoding/binary"
	"net/http"

	"go.opentelemetry.io/api/trace"

	"github.com/lightstep/tracecontext.go"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/propagation"
)

const (
	Vendor = "ot"
)

var (
	HostKey = key.New("http.host")
	URLKey  = key.New("http.url")

	encoding   = binary.BigEndian
	propagator = propagation.HttpTraceContextPropagator()
)

// Returns the Attributes, Context Tags, and SpanContext that were encoded by Inject.
func Extract(ctx context.Context, req *http.Request) ([]core.KeyValue, []core.KeyValue, core.SpanContext) {
	sc := propagator.Extract(ctx, req.Header)

	attrs := []core.KeyValue{
		URLKey.String(req.URL.String()),
		// Etc.
	}

	return attrs, nil, sc
}

func Inject(ctx context.Context, req *http.Request) {
	propagator.Inject(ctx, req.Header)
	sc := trace.CurrentSpan(ctx).SpanContext()
	var tc tracecontext.TraceContext
	var sid [8]byte
	var tid [16]byte

	encoding.PutUint64(sid[0:8], sc.SpanID)
	encoding.PutUint64(tid[0:8], sc.TraceID.High)
	encoding.PutUint64(tid[8:16], sc.TraceID.Low)

	tc.TraceParent.Version = tracecontext.Version
	tc.TraceParent.TraceID = tid
	tc.TraceParent.SpanID = sid
	tc.TraceParent.Flags.Recorded = true // Note: not implemented.

	tags.Foreach(func(kv core.KeyValue) bool {
		// TODO: implement MaxHops
		tc.TraceState = append(tc.TraceState, tracestate.Member{
			Vendor: Vendor,
			Tenant: kv.Key.Name,
			Value:  kv.Value.Emit(),
		})
		return true
	})

	tc.SetHeaders(req.Header)
}
