// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel"
	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	cpb "go.opentelemetry.io/proto/otlp/common/v1"
	lpb "go.opentelemetry.io/proto/otlp/logs/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"

	"go.opentelemetry.io/otel/sdk/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	// Sat Jan 01 2000 00:00:00 GMT+0000.
	ts  = time.Date(2000, time.January, 0o1, 0, 0, 0, 0, time.FixedZone("GMT", 0))
	obs = ts.Add(30 * time.Second)

	kvAlice = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "alice"},
	}}
	kvBob = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "bob"},
	}}
	kvSrvName = &cpb.KeyValue{Key: "service.name", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "test server"},
	}}
	kvSrvVer = &cpb.KeyValue{Key: "service.version", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "v0.1.0"},
	}}

	pbSevA = lpb.SeverityNumber_SEVERITY_NUMBER_INFO
	pbSevB = lpb.SeverityNumber_SEVERITY_NUMBER_ERROR

	pbBodyA = &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{
			StringValue: "a",
		},
	}
	pbBodyB = &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{
			StringValue: "b",
		},
	}

	spanIDA  = []byte{0, 0, 0, 0, 0, 0, 0, 1}
	spanIDB  = []byte{0, 0, 0, 0, 0, 0, 0, 2}
	traceIDA = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	traceIDB = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	flagsA   = byte(1)
	flagsB   = byte(0)

	logRecords = []*lpb.LogRecord{
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevA,
			SeverityText:         "A",
			Body:                 pbBodyA,
			Attributes:           []*cpb.KeyValue{kvAlice},
			Flags:                uint32(flagsA),
			TraceId:              traceIDA,
			SpanId:               spanIDA,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevA,
			SeverityText:         "A",
			Body:                 pbBodyA,
			Attributes:           []*cpb.KeyValue{kvBob},
			Flags:                uint32(flagsA),
			TraceId:              traceIDA,
			SpanId:               spanIDA,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevB,
			SeverityText:         "B",
			Body:                 pbBodyB,
			Attributes:           []*cpb.KeyValue{kvAlice},
			Flags:                uint32(flagsB),
			TraceId:              traceIDB,
			SpanId:               spanIDB,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevB,
			SeverityText:         "B",
			Body:                 pbBodyB,
			Attributes:           []*cpb.KeyValue{kvBob},
			Flags:                uint32(flagsB),
			TraceId:              traceIDB,
			SpanId:               spanIDB,
		},
	}

	scope = &cpb.InstrumentationScope{
		Name:    "test/code/path",
		Version: "v0.1.0",
	}
	scopeLogs = []*lpb.ScopeLogs{
		{
			Scope:      scope,
			LogRecords: logRecords,
			SchemaUrl:  semconv.SchemaURL,
		},
	}

	res = &rpb.Resource{
		Attributes: []*cpb.KeyValue{kvSrvName, kvSrvVer},
	}
	resourceLogs = []*lpb.ResourceLogs{{
		Resource:  res,
		ScopeLogs: scopeLogs,
		SchemaUrl: semconv.SchemaURL,
	}}
)

type exportResult struct {
	Response *collogpb.ExportLogsServiceResponse
	Err      error
}

// storage stores uploaded OTLP log data in their proto form.
type storage struct {
	dataMu sync.Mutex
	data   []*lpb.ResourceLogs
}

// newStorage returns a configure storage ready to store received requests.
func newStorage() *storage {
	return &storage{}
}

// Add adds the request to the Storage.
func (s *storage) Add(request *collogpb.ExportLogsServiceRequest) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()
	s.data = append(s.data, request.ResourceLogs...)
}

// Dump returns all added ResourceLogs and clears the storage.
func (s *storage) Dump() []*lpb.ResourceLogs {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var data []*lpb.ResourceLogs
	data, s.data = s.data, []*lpb.ResourceLogs{}
	return data
}

var emptyExportLogsServiceResponse = func() []byte {
	body := collogpb.ExportLogsServiceResponse{}
	r, err := proto.Marshal(&body)
	if err != nil {
		panic(err)
	}
	return r
}()

type httpResponseError struct {
	Err    error
	Status int
	Header http.Header
}

func (e *httpResponseError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Err)
}

func (e *httpResponseError) Unwrap() error { return e.Err }

// httpCollector is an OTLP HTTP server that collects all requests it receives.
type httpCollector struct {
	plainTextResponse bool

	headersMu sync.Mutex
	headers   http.Header
	storage   *storage

	resultCh <-chan exportResult
	listener net.Listener
	srv      *http.Server
}

// newHTTPCollector returns a *HTTPCollector that is listening at the provided
// endpoint.
//
// If endpoint is an empty string, the returned collector will be listening on
// the localhost interface at an OS chosen port, not use TLS, and listen at the
// default OTLP log endpoint path ("/v1/logs"). If the endpoint contains a
// prefix of "https" the server will generate weak self-signed TLS certificates
// and use them to server data. If the endpoint contains a path, that path will
// be used instead of the default OTLP metric endpoint path.
//
// If errCh is not nil, the collector will respond to HTTP requests with errors
// sent on that channel. This means that if errCh is not nil Export calls will
// block until an error is received.
func newHTTPCollector(endpoint string, resultCh <-chan exportResult, opts ...func(*httpCollector)) (*httpCollector, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if u.Host == "" {
		u.Host = "localhost:0"
	}
	if u.Path == "" {
		u.Path = defaultPath
	}

	c := &httpCollector{
		headers:  http.Header{},
		storage:  newStorage(),
		resultCh: resultCh,
	}
	for _, opt := range opts {
		opt(c)
	}

	c.listener, err = net.Listen("tcp", u.Host)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle(u.Path, http.HandlerFunc(c.handler))
	c.srv = &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if u.Scheme == "https" {
		cert, err := newWeakCertificate()
		if err != nil {
			return nil, err
		}
		c.srv.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		go func() { _ = c.srv.ServeTLS(c.listener, "", "") }()
	} else {
		go func() { _ = c.srv.Serve(c.listener) }()
	}
	return c, nil
}

// withHTTPCollectorRespondingPlainText makes the HTTPCollector return
// a plaintext, instead of protobuf, response.
func withHTTPCollectorRespondingPlainText() func(*httpCollector) {
	return func(s *httpCollector) {
		s.plainTextResponse = true
	}
}

// Shutdown shuts down the HTTP server closing all open connections and
// listeners.
func (c *httpCollector) Shutdown(ctx context.Context) error {
	return c.srv.Shutdown(ctx)
}

// Addr returns the net.Addr c is listening at.
func (c *httpCollector) Addr() net.Addr {
	return c.listener.Addr()
}

// Collect returns the Storage holding all collected requests.
func (c *httpCollector) Collect() *storage {
	return c.storage
}

// Headers returns the headers received for all requests.
func (c *httpCollector) Headers() map[string][]string {
	// Makes a copy.
	c.headersMu.Lock()
	defer c.headersMu.Unlock()
	return c.headers.Clone()
}

func (c *httpCollector) handler(w http.ResponseWriter, r *http.Request) {
	c.respond(w, c.record(r))
}

func (c *httpCollector) record(r *http.Request) exportResult {
	// Currently only supports protobuf.
	if v := r.Header.Get("Content-Type"); v != "application/x-protobuf" {
		err := fmt.Errorf("content-type not supported: %s", v)
		return exportResult{Err: err}
	}

	body, err := c.readBody(r)
	if err != nil {
		return exportResult{Err: err}
	}
	pbRequest := &collogpb.ExportLogsServiceRequest{}
	err = proto.Unmarshal(body, pbRequest)
	if err != nil {
		return exportResult{
			Err: &httpResponseError{
				Err:    err,
				Status: http.StatusInternalServerError,
			},
		}
	}
	c.storage.Add(pbRequest)

	c.headersMu.Lock()
	for k, vals := range r.Header {
		for _, v := range vals {
			c.headers.Add(k, v)
		}
	}
	c.headersMu.Unlock()

	if c.resultCh != nil {
		return <-c.resultCh
	}
	return exportResult{Err: err}
}

func (c *httpCollector) readBody(r *http.Request) (body []byte, err error) {
	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			_ = reader.Close()
			return nil, &httpResponseError{
				Err:    err,
				Status: http.StatusInternalServerError,
			}
		}
	default:
		reader = r.Body
	}

	defer func() {
		cErr := reader.Close()
		if err == nil && cErr != nil {
			err = &httpResponseError{
				Err:    cErr,
				Status: http.StatusInternalServerError,
			}
		}
	}()
	body, err = io.ReadAll(reader)
	if err != nil {
		err = &httpResponseError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}
	return body, err
}

func (c *httpCollector) respond(w http.ResponseWriter, resp exportResult) {
	if resp.Err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		var e *httpResponseError
		if errors.As(resp.Err, &e) {
			for k, vals := range e.Header {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			w.WriteHeader(e.Status)
			fmt.Fprintln(w, e.Error())
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, resp.Err.Error())
		}
		return
	}

	if c.plainTextResponse {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
		return
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusOK)
	if resp.Response == nil {
		_, _ = w.Write(emptyExportLogsServiceResponse)
	} else {
		r, err := proto.Marshal(resp.Response)
		if err != nil {
			panic(err)
		}
		_, _ = w.Write(r)
	}
}

// Based on https://golang.org/src/crypto/tls/generate_cert.go,
// simplified and weakened.
func newWeakCertificate() (tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour)
	m := new(big.Int).Lsh(big.NewInt(1), 128)
	sn, err := rand.Int(rand.Reader, m)
	if err != nil {
		return tls.Certificate{}, err
	}
	tmpl := x509.Certificate{
		SerialNumber:          sn,
		Subject:               pkix.Name{Organization: []string{"otel-go"}},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv6loopback, net.IPv4(127, 0, 0, 1)},
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}
	var certBuf bytes.Buffer
	err = pem.Encode(&certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return tls.Certificate{}, err
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return tls.Certificate{}, err
	}
	var privBuf bytes.Buffer
	err = pem.Encode(&privBuf, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.X509KeyPair(certBuf.Bytes(), privBuf.Bytes())
}

func TestClient(t *testing.T) {
	factory := func(rCh <-chan exportResult) (*client, *httpCollector) {
		coll, err := newHTTPCollector("", rCh)
		require.NoError(t, err)

		addr := coll.Addr().String()
		opts := []Option{WithEndpoint(addr), WithInsecure()}
		cfg := newConfig(opts)
		client, err := newHTTPClient(cfg)
		require.NoError(t, err)
		return client, coll
	}

	t.Run("ClientHonorsContextErrors", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		t.Run("DeadlineExceeded", func(t *testing.T) {
			innerCtx, innerCancel := context.WithTimeout(ctx, time.Nanosecond)
			t.Cleanup(innerCancel)
			<-innerCtx.Done()

			c, _ := factory(nil)
			assert.ErrorIs(t, c.uploadLogs(innerCtx, nil), context.DeadlineExceeded)
		})

		t.Run("Canceled", func(t *testing.T) {
			innerCtx, innerCancel := context.WithCancel(ctx)
			innerCancel()

			c, _ := factory(nil)
			assert.ErrorIs(t, c.uploadLogs(innerCtx, nil), context.Canceled)
		})
	})

	t.Run("uploadLogs", func(t *testing.T) {
		ctx := context.Background()
		client, coll := factory(nil)

		require.NoError(t, client.uploadLogs(ctx, resourceLogs))
		got := coll.Collect().Dump()
		require.Len(t, got, 1, "upload of one ResourceLogs")
		diff := cmp.Diff(got[0], resourceLogs[0], cmp.Comparer(proto.Equal))
		if diff != "" {
			t.Fatalf("unexpected ResourceLogs:\n%s", diff)
		}
	})

	t.Run("PartialSuccess", func(t *testing.T) {
		const n, msg = 2, "bad data"
		rCh := make(chan exportResult, 3)
		rCh <- exportResult{
			Response: &collogpb.ExportLogsServiceResponse{
				PartialSuccess: &collogpb.ExportLogsPartialSuccess{
					RejectedLogRecords: n,
					ErrorMessage:       msg,
				},
			},
		}
		rCh <- exportResult{
			Response: &collogpb.ExportLogsServiceResponse{
				PartialSuccess: &collogpb.ExportLogsPartialSuccess{
					// Should not be logged.
					RejectedLogRecords: 0,
					ErrorMessage:       "",
				},
			},
		}
		rCh <- exportResult{
			Response: &collogpb.ExportLogsServiceResponse{},
		}

		ctx := context.Background()
		client, _ := factory(rCh)

		defer func(orig otel.ErrorHandler) {
			otel.SetErrorHandler(orig)
		}(otel.GetErrorHandler())

		errs := []error{}
		eh := otel.ErrorHandlerFunc(func(e error) { errs = append(errs, e) })
		otel.SetErrorHandler(eh)

		require.NoError(t, client.UploadLogs(ctx, resourceLogs))
		require.NoError(t, client.UploadLogs(ctx, resourceLogs))
		require.NoError(t, client.UploadLogs(ctx, resourceLogs))

		require.Len(t, errs, 1)
		want := fmt.Sprintf("%s (%d log records rejected)", msg, n)
		assert.ErrorContains(t, errs[0], want)
	})
}

func TestClientWithHTTPCollectorRespondingPlainText(t *testing.T) {
	ctx := context.Background()
	coll, err := newHTTPCollector("", nil, withHTTPCollectorRespondingPlainText())
	require.NoError(t, err)

	addr := coll.Addr().String()
	opts := []Option{WithEndpoint(addr), WithInsecure()}
	cfg := newConfig(opts)
	client, err := newHTTPClient(cfg)
	require.NoError(t, err)

	require.NoError(t, client.uploadLogs(ctx, make([]*lpb.ResourceLogs, 1)))
	got := coll.Collect().Dump()
	require.Len(t, got, 1, "upload of one ResourceLogs")
}

func TestNewWithInvalidEndpoint(t *testing.T) {
	ctx := context.Background()
	exp, err := New(ctx, WithEndpoint("host:invalid-port"))
	assert.Error(t, err)
	assert.Nil(t, exp)
}

func TestConfig(t *testing.T) {
	factoryFunc := func(ePt string, rCh <-chan exportResult, o ...Option) (log.Exporter, *httpCollector) {
		coll, err := newHTTPCollector(ePt, rCh)
		require.NoError(t, err)

		opts := []Option{WithEndpoint(coll.Addr().String())}
		if !strings.HasPrefix(strings.ToLower(ePt), "https") {
			opts = append(opts, WithInsecure())
		}
		opts = append(opts, o...)

		ctx := context.Background()
		exp, err := New(ctx, opts...)
		require.NoError(t, err)
		return exp, coll
	}

	t.Run("WithEndpointURL", func(t *testing.T) {
		coll, err := newHTTPCollector("", nil)
		require.NoError(t, err)
		ctx := context.Background()

		target := "http://" + coll.Addr().String() + defaultPath
		exp, err := New(ctx, WithEndpointURL(target))
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

		assert.NoError(t, exp.Export(ctx, make([]log.Record, 1)))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithHeaders", func(t *testing.T) {
		key := http.CanonicalHeaderKey("my-custom-header")
		headers := map[string]string{key: "custom-value"}
		exp, coll := factoryFunc("", nil, WithHeaders(headers))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		require.NoError(t, exp.Export(ctx, make([]log.Record, 1)))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := coll.Headers()
		require.Regexp(t, "OTel Go OTLP over HTTP/protobuf logs exporter/[01]\\..*", got)
		require.Contains(t, got, key)
		assert.Equal(t, []string{headers[key]}, got[key])
	})

	t.Run("WithTimeout", func(t *testing.T) {
		// Do not send on rCh so the Collector never responds to the client.
		rCh := make(chan exportResult)
		exp, coll := factoryFunc(
			"",
			rCh,
			WithTimeout(time.Millisecond),
			WithRetry(RetryConfig{Enabled: false}),
		)
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		// Push this after Shutdown so the HTTP server doesn't hang.
		t.Cleanup(func() { close(rCh) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		err := exp.Export(ctx, make([]log.Record, 1))
		assert.ErrorAs(t, err, new(retryableError))
	})

	t.Run("WithCompressionGZip", func(t *testing.T) {
		exp, coll := factoryFunc("", nil, WithCompression(GzipCompression))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, make([]log.Record, 1)))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithRetry", func(t *testing.T) {
		emptyErr := errors.New("")
		rCh := make(chan exportResult, 5)
		header := http.Header{http.CanonicalHeaderKey("Retry-After"): {"10"}}
		// All retryable errors.
		rCh <- exportResult{Err: &httpResponseError{
			Status: http.StatusServiceUnavailable,
			Err:    emptyErr,
			Header: header,
		}}
		rCh <- exportResult{Err: &httpResponseError{
			Status: http.StatusTooManyRequests,
			Err:    emptyErr,
		}}
		rCh <- exportResult{Err: &httpResponseError{
			Status: http.StatusGatewayTimeout,
			Err:    emptyErr,
		}}
		rCh <- exportResult{Err: &httpResponseError{
			Status: http.StatusBadGateway,
			Err:    emptyErr,
		}}
		rCh <- exportResult{}
		exp, coll := factoryFunc("", rCh, WithRetry(RetryConfig{
			Enabled:         true,
			InitialInterval: time.Nanosecond,
			MaxInterval:     time.Millisecond,
			MaxElapsedTime:  time.Minute,
		}))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		// Push this after Shutdown so the HTTP server doesn't hang.
		t.Cleanup(func() { close(rCh) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, make([]log.Record, 1)), "failed retry")
		assert.Empty(t, rCh, "failed HTTP responses did not occur")
	})

	t.Run("WithRetryAndExporterErr", func(t *testing.T) {
		exporterErr := errors.New("rpc error: code = Unavailable desc = service.name not found in resource attributes")
		rCh := make(chan exportResult, 1)
		rCh <- exportResult{Err: &httpResponseError{
			Status: http.StatusTooManyRequests,
			Err:    exporterErr,
		}}
		exp, coll := factoryFunc("", rCh, WithRetry(RetryConfig{
			Enabled: false,
		}))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		// Push this after Shutdown so the HTTP server doesn't hang.
		t.Cleanup(func() { close(rCh) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		err := exp.Export(ctx, make([]log.Record, 1))
		assert.ErrorContains(t, err, exporterErr.Error())

		// To test the `Unwrap` and `As` function of retryable error
		var retryErr *retryableError
		assert.ErrorAs(t, err, &retryErr)
		assert.ErrorIs(t, err, *retryErr)
	})

	t.Run("WithURLPath", func(t *testing.T) {
		path := "/prefix/v2/logs"
		ePt := fmt.Sprintf("http://localhost:0%s", path)
		exp, coll := factoryFunc(ePt, nil, WithURLPath(path))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, make([]log.Record, 1)))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithTLSClientConfig", func(t *testing.T) {
		ePt := "https://localhost:0"
		tlsCfg := &tls.Config{InsecureSkipVerify: true}
		exp, coll := factoryFunc(ePt, nil, WithTLSClientConfig(tlsCfg))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, make([]log.Record, 1)))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithCustomUserAgent", func(t *testing.T) {
		key := http.CanonicalHeaderKey("user-agent")
		headers := map[string]string{key: "custom-user-agent"}
		exp, coll := factoryFunc("", nil, WithHeaders(headers))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		require.NoError(t, exp.Export(ctx, make([]log.Record, 1)))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := coll.Headers()
		require.Contains(t, got, key)
		assert.Equal(t, []string{headers[key]}, got[key])
	})

	t.Run("WithProxy", func(t *testing.T) {
		headerKeySetInProxy := http.CanonicalHeaderKey("X-Using-Proxy")
		headerValueSetInProxy := "true"
		exp, coll := factoryFunc("", nil, WithProxy(func(r *http.Request) (*url.URL, error) {
			r.Header.Set(headerKeySetInProxy, headerValueSetInProxy)
			return r.URL, nil
		}))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		require.NoError(t, exp.Export(ctx, make([]log.Record, 1)))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := coll.Headers()
		require.Contains(t, got, headerKeySetInProxy)
		assert.Equal(t, []string{headerValueSetInProxy}, got[headerKeySetInProxy])
	})

	t.Run("non-retryable errors are propagated", func(t *testing.T) {
		exporterErr := errors.New("missing required attribute aaaa")
		rCh := make(chan exportResult, 1)
		rCh <- exportResult{Err: &httpResponseError{
			Status: http.StatusBadRequest,
			Err:    exporterErr,
		}}

		exp, coll := factoryFunc("", rCh, WithRetry(RetryConfig{
			Enabled: false,
		}))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		// Push this after Shutdown so the HTTP server doesn't hang.
		t.Cleanup(func() { close(rCh) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		err := exp.Export(ctx, make([]log.Record, 1))
		assert.ErrorContains(t, err, exporterErr.Error())
	})
}
