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

	tr trace.Tracer

	activeHooks map[string]trace.Span
	root        trace.Span
	mtx         sync.Mutex
}

func NewClientTrace(ctx context.Context) *httptrace.ClientTrace {
	ct := &clientTracer{
		Context:     ctx,
		activeHooks: make(map[string]trace.Span),
	}

	ct.tr = global.TraceProvider().GetTracer("go.opentelemetry.io/otel/plugin/httptrace")
	ct.start("http.request", "http.request")

	return &httptrace.ClientTrace{
		GetConn:              ct.getConn,
		GotConn:              ct.gotConn,
		PutIdleConn:          ct.putIdleConn,
		GotFirstResponseByte: ct.gotFirstResponseByte,
		Got100Continue:       ct.got100Continue,
		Got1xxResponse:       ct.got1xxResponse,
		DNSStart:             ct.dnsStart,
		DNSDone:              ct.dnsDone,
		ConnectStart:         ct.connectStart,
		ConnectDone:          ct.connectDone,
		TLSHandshakeStart:    ct.tlsHandshakeStart,
		TLSHandshakeDone:     ct.tlsHandshakeDone,
		WroteHeaderField:     ct.wroteHeaderField,
		WroteHeaders:         ct.wroteHeaders,
		Wait100Continue:      ct.wait100Continue,
		WroteRequest:         ct.wroteRequest,
	}
}

func (ct *clientTracer) start(hook, spanName string, attrs ...core.KeyValue) {
	_, sp := ct.tr.Start(ct.Context, spanName, trace.WithAttributes(attrs...), trace.WithSpanKind(trace.SpanKindClient))
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

func (ct *clientTracer) end(hook string, err error, attrs ...core.KeyValue) {
	ct.mtx.Lock()
	defer ct.mtx.Unlock()
	if span, ok := ct.activeHooks[hook]; ok {
		if err != nil {
			span.SetStatus(codes.Unknown)
			span.SetAttribute(MessageKey.String(err.Error()))
		}
		span.SetAttributes(attrs...)
		span.End()
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
	ct.end("http.getconn",
		nil,
		HTTPRemoteAddr.String(info.Conn.RemoteAddr().String()),
		HTTPLocalAddr.String(info.Conn.LocalAddr().String()),
	)
}

func (ct *clientTracer) putIdleConn(err error) {
	ct.end("http.receive", err)
}

func (ct *clientTracer) gotFirstResponseByte() {
	ct.start("http.receive", "http.receive")
}

func (ct *clientTracer) dnsStart(info httptrace.DNSStartInfo) {
	ct.start("http.dns", "http.dns", HostKey.String(info.Host))
}

func (ct *clientTracer) dnsDone(info httptrace.DNSDoneInfo) {
	ct.end("http.dns", info.Err)
}

func (ct *clientTracer) connectStart(network, addr string) {
	ct.start("http.connect."+addr, "http.connect", HTTPRemoteAddr.String(addr))
}

func (ct *clientTracer) connectDone(network, addr string, err error) {
	ct.end("http.connect."+addr, err)
}

func (ct *clientTracer) tlsHandshakeStart() {
	ct.start("http.tls", "http.tls")
}

func (ct *clientTracer) tlsHandshakeDone(_ tls.ConnectionState, err error) {
	ct.end("http.tls", err)
}

func (ct *clientTracer) wroteHeaderField(k string, v []string) {
	if ct.span("http.headers") == nil {
		ct.start("http.headers", "http.headers")
	}
	ct.root.SetAttribute(key.String("http."+strings.ToLower(k), sliceToString(v)))
}

func (ct *clientTracer) wroteHeaders() {
	ct.start("http.send", "http.send")
}

func (ct *clientTracer) wroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		ct.root.SetAttribute(MessageKey.String(info.Err.Error()))
		ct.root.SetStatus(codes.Unknown)
	}
	ct.end("http.send", info.Err)
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

func sliceToString(value []string) string {
	if len(value) == 0 {
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
		buf.WriteString(sliceToString(v))
	}
	return buf.String()
}
