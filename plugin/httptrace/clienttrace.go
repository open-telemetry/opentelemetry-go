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
	"crypto/tls"
	"net/http/httptrace"
	"net/textproto"
	"strings"
	"sync"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/global"
)

var (
	HTTPStatus     = key.New("http.status")
	HTTPHeaderMIME = key.New("http.mime")
	HTTPRemoteAddr = key.New("http.remote")
	HTTPLocalAddr  = key.New("http.local")
	MessageKey     = key.New("message")
)

type clientTracer struct {
	context.Context
	httptrace.ClientTrace

	tr trace.Tracer

	activeHooks map[string]trace.Span
	root        trace.Span
	mtx         sync.Mutex
}

func newClientTracer(ctx context.Context) *clientTracer {
	ct := &clientTracer{
		Context:     ctx,
		activeHooks: make(map[string]trace.Span),
	}
	ct.tr = global.TraceProvider().GetTracer("go.opentelemetry.io/otel/plugin/httptrace")
	ct.start("http.request", "http.request")
	return ct
}

func (ct *clientTracer) start(hook, span string, attrs ...core.KeyValue) {
	_, sp := ct.tr.Start(ct.Context, span, trace.WithAttributes(attrs...), trace.WithSpanKind(trace.SpanKindClient))
	// TODO(paivagustavo): remove this for loop when `trace.WithAttributes(attrs...)` works.
	for _, attr := range attrs {
		sp.SetAttribute(attr)
	}
	ct.mtx.Lock()
	defer ct.mtx.Unlock()
	if ct.root == nil {
		ct.root = sp
	}
	if _, ok := ct.activeHooks[hook]; ok {
		// end was called before start is handled.
		sp.End()
		delete(ct.activeHooks, hook)
	} else {
		ct.activeHooks[hook] = sp
	}
}

func (ct *clientTracer) end(hook string) {
	ct.mtx.Lock()
	defer ct.mtx.Unlock()
	if s, ok := ct.activeHooks[hook]; ok {
		s.End()
		delete(ct.activeHooks, hook)
	} else {
		// start is not finished before end is called.
		ct.activeHooks[hook] = trace.NoopSpan{}
	}
}

func (ct *clientTracer) span(hook string) trace.Span {
	ct.mtx.Lock()
	defer ct.mtx.Unlock()
	return ct.activeHooks[hook]
}

func (ct *clientTracer) getConn(host string) {
	ct.start("http.getconn", "http.getconn", HostKey.String(host))
}

func (ct *clientTracer) gotConn(info httptrace.GotConnInfo) {
	ct.span("http.getconn").SetAttributes(
		HTTPRemoteAddr.String(info.Conn.RemoteAddr().String()),
		HTTPLocalAddr.String(info.Conn.LocalAddr().String()),
	)
	ct.end("http.getconn")
}

func (ct *clientTracer) putIdleConn(err error) {
	if err != nil {
		ct.span("http.receive").SetAttribute(MessageKey.String(err.Error()))
		ct.span("http.receive").SetStatus(codes.Unknown)
	}
	ct.end("http.receive")
}

func (ct *clientTracer) gotFirstResponseByte() {
	ct.start("http.receive", "http.receive")
}

func (ct *clientTracer) dnsStart(httptrace.DNSStartInfo) {
	ct.start("http.dns", "http.dns")
}

func (ct *clientTracer) dnsDone(httptrace.DNSDoneInfo) {
	ct.end("http.dns")
}

func (ct *clientTracer) connectStart(network, addr string) {
	ct.start("http.connect."+addr, "http.connect", HTTPRemoteAddr.String(addr))
}

func (ct *clientTracer) connectDone(network, addr string, err error) {
	ct.end("http.connect." + addr)
}

func (ct *clientTracer) tlsHandshakeStart() {
	ct.start("http.tls", "http.tls")
}

func (ct *clientTracer) tlsHandshakeDone(tls.ConnectionState, error) {
	ct.end("http.tls")
}

func (ct *clientTracer) wroteHeaderField(k string, v []string) {
	if ct.span("http.headers") == nil {
		ct.start("http.headers", "http.headers")
	}
	ct.root.SetAttribute(key.New("http." + strings.ToLower(k)).String(sa2s(v)))
}

func (ct *clientTracer) wroteHeaders() {
	ct.start("http.send", "http.send")
}

func (ct *clientTracer) wroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		ct.root.SetAttribute(MessageKey.String(info.Err.Error()))
		ct.root.SetStatus(codes.Unknown)
	}
	ct.end("http.send")
}

func (ct *clientTracer) got100Continue() {
	ct.span("http.receive").AddEvent(ct.Context, "GOT 100 - Continue")
}

func (ct *clientTracer) wait100Continue() {
	ct.span("http.receive").AddEvent(ct.Context, "GOT 100 - Wait")
}

func (ct *clientTracer) got1xxResponse(code int, header textproto.MIMEHeader) error {
	ct.span("http.receive").AddEvent(ct.Context, "GOT 1xx",
		HTTPStatus.Int(code),
		HTTPHeaderMIME.String(sm2s(header)),
	)
	return nil
}

func sa2s(value []string) string {
	if len(value) == 1 {
		return value[0]
	} else if len(value) == 0 {
		return "undefined"
	}
	return strings.Join(value, ",")
}

func sm2s(value map[string][]string) string {
	var buf strings.Builder
	for k, v := range value {
		if buf.Len() != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(sa2s(v))
	}
	return buf.String()
}
