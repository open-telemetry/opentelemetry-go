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
	"fmt"
	"net/http/httptrace"
	"net/textproto"
	"strings"
	"sync"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/trace"
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

	levels map[string]trace.Span
	root   trace.Span
	mtx    sync.Mutex
}

func newClientTracer(ctx context.Context) *clientTracer {
	ct := &clientTracer{
		Context: ctx,
		levels:  make(map[string]trace.Span),
	}
	ct.tr = trace.GlobalProvider().GetTracer("go.opentelemetry.io/plugin/httptrace")
	ct.open("http.request")
	return ct
}

func (ct *clientTracer) open(name string, attrs ...core.KeyValue) {
	_, sp := ct.tr.Start(ct.Context, name, trace.WithAttributes(attrs...), trace.WithSpanKind(trace.SpanKindClient))
	ct.mtx.Lock()
	defer ct.mtx.Unlock()
	if ct.root == nil {
		ct.root = sp
	}
	ct.levels[name] = sp
}

func (ct *clientTracer) close(name string) {
	ct.mtx.Lock()
	defer ct.mtx.Unlock()
	if s, ok := ct.levels[name]; ok {
		s.End()
		delete(ct.levels, name)
	} else {
		panic(fmt.Sprintf("failed to find span %s in levels.", name))
	}
}

func (ct *clientTracer) span(name string) trace.Span {
	ct.mtx.Lock()
	defer ct.mtx.Unlock()
	return ct.levels[name]
}

func (ct *clientTracer) getConn(host string) {
	ct.open("http.getconn", HostKey.String(host))
}

func (ct *clientTracer) gotConn(info httptrace.GotConnInfo) {
	ct.span("http.getconn").SetAttribute(HTTPRemoteAddr.String(info.Conn.RemoteAddr().String()))
	ct.span("http.getconn").SetAttribute(HTTPLocalAddr.String(info.Conn.LocalAddr().String()))

	ct.close("http.getconn")
}

func (ct *clientTracer) putIdleConn(err error) {
	if err != nil {
		ct.span("http.receive").SetAttribute(MessageKey.String(err.Error()))
		ct.span("http.receive").SetStatus(codes.Unknown)
	}
	ct.close("http.receive")
}

func (ct *clientTracer) gotFirstResponseByte() {
	ct.open("http.receive")
}

func (ct *clientTracer) dnsStart(httptrace.DNSStartInfo) {
	ct.open("http.dns")
}

func (ct *clientTracer) dnsDone(httptrace.DNSDoneInfo) {
	ct.close("http.dns")
}

func (ct *clientTracer) connectStart(network, addr string) {
	ct.open("http.connect")
}

func (ct *clientTracer) connectDone(network, addr string, err error) {
	ct.close("http.connect")
}

func (ct *clientTracer) tlsHandshakeStart() {
	ct.open("http.tls")
}

func (ct *clientTracer) tlsHandshakeDone(tls.ConnectionState, error) {
	ct.close("http.tls")
}

func (ct *clientTracer) wroteHeaderField(k string, v []string) {
	if ct.span("http.headers") == nil {
		ct.open("http.headers")
	}
	ct.root.SetAttribute(key.New("http." + strings.ToLower(k)).String(sa2s(v)))
}

func (ct *clientTracer) wroteHeaders() {
	ct.open("http.send")
}

func (ct *clientTracer) wroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		ct.root.SetAttribute(MessageKey.String(info.Err.Error()))
		ct.root.SetStatus(codes.Unknown)
	}
	ct.close("http.send")
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
