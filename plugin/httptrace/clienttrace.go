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

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/trace"
)

type clientLevel struct {
	trace.Span
	ident string
}

var (
	HTTPStatus     = key.New("http.status")
	HTTPHeaderMIME = key.New("http.mime")
	HTTPRemoteAddr = key.New("http.remote")
	HTTPLocalAddr  = key.New("http.local")
	MessageKey     = key.New("message",
		key.WithDescription("message text: info, error, etc"),
	)
)

type clientTracer struct {
	context.Context
	httptrace.ClientTrace

	levels []clientLevel
}

func newClientTracer(ctx context.Context) *clientTracer {
	ct := &clientTracer{
		Context: ctx,
	}
	ct.open("http.request")
	return ct
}

func (ct *clientTracer) open(name string, attrs ...core.KeyValue) {
	_, sp := trace.Start(ct.Context, name, trace.WithAttributes(attrs...))
	ct.levels = append(ct.levels, clientLevel{
		Span:  sp,
		ident: name,
	})
}

func (ct *clientTracer) close(name string) {
	if len(ct.levels) == 0 {
		panic("remove me")
	}
	l := len(ct.levels)
	ct.levels[l-1].Finish()
	ct.levels = ct.levels[0 : l-1]
}

func (ct *clientTracer) current() trace.Span {
	return ct.levels[len(ct.levels)-1].Span
}

func (ct *clientTracer) currentName() string {
	if len(ct.levels) == 0 {
		return ""
	}
	return ct.levels[len(ct.levels)-1].ident
}

func (ct *clientTracer) getConn(host string) {
	ct.open("http.getconn", HostKey.String(host))
}

func (ct *clientTracer) gotConn(info httptrace.GotConnInfo) {
	ct.current().SetAttribute(HTTPRemoteAddr.String(info.Conn.RemoteAddr().String()))
	ct.current().SetAttribute(HTTPLocalAddr.String(info.Conn.LocalAddr().String()))

	ct.close("http.getconn")
}

func (ct *clientTracer) putIdleConn(err error) {
	if err != nil {
		ct.current().SetAttribute(MessageKey.String(err.Error()))
		ct.current().SetStatus(codes.Unknown)
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
	if ct.currentName() != "http.headers" {
		ct.open("http.headers")
	}
	ct.levels[0].SetAttribute(key.New("http." + strings.ToLower(k)).String(sa2s(v)))
}

func (ct *clientTracer) wroteHeaders() {
	ct.open("http.send")
}

func (ct *clientTracer) wroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		ct.levels[0].SetAttribute(MessageKey.String(info.Err.Error()))
		ct.levels[0].SetStatus(codes.Unknown)
	}
	ct.close("http.send")
}

func (ct *clientTracer) got100Continue() {
	ct.current().Event(ct.Context, "GOT 100 - Continue")
}

func (ct *clientTracer) wait100Continue() {
	ct.current().Event(ct.Context, "GOT 100 - Wait")
}

func (ct *clientTracer) got1xxResponse(code int, header textproto.MIMEHeader) error {
	ct.current().Event(ct.Context, "GOT 1xx",
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
