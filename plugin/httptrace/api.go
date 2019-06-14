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
