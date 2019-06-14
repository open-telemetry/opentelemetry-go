package httptrace

import (
	"encoding/binary"
	"net/http"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
	"github.com/lightstep/tracecontext.go"
	"github.com/lightstep/tracecontext.go/tracestate"
)

const (
	Vendor = "ot"
)

type (
	hinjector struct {
		*http.Request
	}
)

var (
	HostKey = tag.New("http.host")
	URLKey  = tag.New("http.url")

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
	sc.TraceIDHigh = encoding.Uint64(tc.TraceParent.TraceID[0:8])
	sc.TraceIDLow = encoding.Uint64(tc.TraceParent.TraceID[8:16])

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
		tags = append(tags, tag.New(ts.Tenant).String(ts.Value))
	}

	return attrs, tags, sc
}

func (h hinjector) Inject(sc core.SpanContext, tags tag.Map) {
	var tc tracecontext.TraceContext
	var sid [8]byte
	var tid [16]byte

	encoding.PutUint64(sid[0:8], sc.SpanID)
	encoding.PutUint64(tid[0:8], sc.TraceIDHigh)
	encoding.PutUint64(tid[8:16], sc.TraceIDLow)

	tc.TraceParent.Version = tracecontext.Version
	tc.TraceParent.TraceID = tid
	tc.TraceParent.SpanID = sid
	tc.TraceParent.Flags.Recorded = true // Note: not implemented.

	tags.Foreach(func(kv core.KeyValue) bool {
		// TODO: implement MaxHops
		tc.TraceState = append(tc.TraceState, tracestate.Member{
			Vendor: Vendor,
			Tenant: kv.Key.Name(),
			Value:  kv.Value.Emit(),
		})
		return true
	})

	tc.SetHeaders(h.Header)
}
