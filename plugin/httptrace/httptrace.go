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
	"encoding/binary"
	"net/http"

	"github.com/lightstep/tracecontext.go"
	"github.com/lightstep/tracecontext.go/tracestate"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/tag"
)

const (
	Vendor = "ot"
)

var (
	HostKey = key.New("http.host")
	URLKey  = key.New("http.url")

	encoding = binary.BigEndian
)

// Returns the Attributes, Context Tags, and SpanContext that were encoded by Inject.
func Extract(req *http.Request) ([]core.KeyValue, []core.KeyValue, core.SpanContext) {
	tc, err := tracecontext.FromHeaders(req.Header)

	if err != nil {
		return nil, nil, core.SpanContext{}
	}

	var sc core.SpanContext
	sc.SpanID = encoding.Uint64(tc.TraceParent.SpanID[0:8])
	sc.TraceID.High = encoding.Uint64(tc.TraceParent.TraceID[0:8])
	sc.TraceID.Low = encoding.Uint64(tc.TraceParent.TraceID[8:16])

	attrs := []core.KeyValue{
		URLKey.String(req.URL.String()),
		// Etc.
	}

	var tags []core.KeyValue

	for _, ts := range tc.TraceState {
		if ts.Vendor != Vendor {
			continue
		}
		// TODO: max-hops, type conversion questions answered,
		// case-conversion questions.
		tags = append(tags, key.New(ts.Tenant).String(ts.Value))
	}

	return attrs, tags, sc
}

type hinjector struct {
	*http.Request
}

func (h hinjector) Inject(sc core.SpanContext, tags tag.Map) {
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
			Tenant: kv.Key.Variable.Name,
			Value:  kv.Value.Emit(),
		})
		return true
	})

	tc.SetHeaders(h.Header)
}
