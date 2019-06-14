package httptrace

import (
	"context"
	"net/http"
	"net/http/httptrace"

	"github.com/open-telemetry/opentelemetry-go/api/trace"
)

// Client
func W3C(ctx context.Context, req *http.Request) (context.Context, *http.Request, trace.Injector) {
	t := newClientTracer(ctx)

	t.GetConn = t.getConn
	t.GotConn = t.gotConn
	t.PutIdleConn = t.putIdleConn
	t.GotFirstResponseByte = t.gotFirstResponseByte
	t.Got100Continue = t.got100Continue
	t.Got1xxResponse = t.got1xxResponse
	t.DNSStart = t.dnsStart
	t.DNSDone = t.dnsDone
	t.ConnectStart = t.connectStart
	t.ConnectDone = t.connectDone
	t.TLSHandshakeStart = t.tlsHandshakeStart
	t.TLSHandshakeDone = t.tlsHandshakeDone
	t.WroteHeaderField = t.wroteHeaderField
	t.WroteHeaders = t.wroteHeaders
	t.Wait100Continue = t.wait100Continue
	t.WroteRequest = t.wroteRequest

	ctx = httptrace.WithClientTrace(ctx, &t.ClientTrace)
	req = req.WithContext(ctx)
	return ctx, req, hinjector{req}
}
