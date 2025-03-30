// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/otlpmetric/otest/collector.go.tmpl

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/otest"

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix" // nolint:depguard  // This is for testing.
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/oconf"
	collpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

// Collector is the collection target a Client sends metric uploads to.
type Collector interface {
	Collect() *Storage
}

type ExportResult struct {
	Response *collpb.ExportMetricsServiceResponse
	Err      error
}

// Storage stores uploaded OTLP metric data in their proto form.
type Storage struct {
	dataMu sync.Mutex
	data   []*mpb.ResourceMetrics
}

// NewStorage returns a configure storage ready to store received requests.
func NewStorage() *Storage {
	return &Storage{}
}

// Add adds the request to the Storage.
func (s *Storage) Add(request *collpb.ExportMetricsServiceRequest) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()
	s.data = append(s.data, request.ResourceMetrics...)
}

// Dump returns all added ResourceMetrics and clears the storage.
func (s *Storage) Dump() []*mpb.ResourceMetrics {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var data []*mpb.ResourceMetrics
	data, s.data = s.data, []*mpb.ResourceMetrics{}
	return data
}

// GRPCCollector is an OTLP gRPC server that collects all requests it receives.
type GRPCCollector struct {
	collpb.UnimplementedMetricsServiceServer

	headersMu sync.Mutex
	headers   metadata.MD
	storage   *Storage

	resultCh <-chan ExportResult
	listener net.Listener
	srv      *grpc.Server
}

// NewGRPCCollector returns a *GRPCCollector that is listening at the provided
// endpoint.
//
// If endpoint is an empty string, the returned collector will be listening on
// the localhost interface at an OS chosen port.
//
// If errCh is not nil, the collector will respond to Export calls with errors
// sent on that channel. This means that if errCh is not nil Export calls will
// block until an error is received.
func NewGRPCCollector(endpoint string, resultCh <-chan ExportResult) (*GRPCCollector, error) {
	if endpoint == "" {
		endpoint = "localhost:0"
	}

	c := &GRPCCollector{
		storage:  NewStorage(),
		resultCh: resultCh,
	}

	var err error
	c.listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	c.srv = grpc.NewServer()
	collpb.RegisterMetricsServiceServer(c.srv, c)
	go func() { _ = c.srv.Serve(c.listener) }()

	return c, nil
}

// Shutdown shuts down the gRPC server closing all open connections and
// listeners immediately.
func (c *GRPCCollector) Shutdown() { c.srv.Stop() }

// Addr returns the net.Addr c is listening at.
func (c *GRPCCollector) Addr() net.Addr {
	return c.listener.Addr()
}

// Collect returns the Storage holding all collected requests.
func (c *GRPCCollector) Collect() *Storage {
	return c.storage
}

// Headers returns the headers received for all requests.
func (c *GRPCCollector) Headers() map[string][]string {
	// Makes a copy.
	c.headersMu.Lock()
	defer c.headersMu.Unlock()
	return metadata.Join(c.headers)
}

// Export handles the export req.
func (c *GRPCCollector) Export(
	ctx context.Context,
	req *collpb.ExportMetricsServiceRequest,
) (*collpb.ExportMetricsServiceResponse, error) {
	c.storage.Add(req)

	if h, ok := metadata.FromIncomingContext(ctx); ok {
		c.headersMu.Lock()
		c.headers = metadata.Join(c.headers, h)
		c.headersMu.Unlock()
	}

	if c.resultCh != nil {
		r := <-c.resultCh
		if r.Response == nil {
			return &collpb.ExportMetricsServiceResponse{}, r.Err
		}
		return r.Response, r.Err
	}
	return &collpb.ExportMetricsServiceResponse{}, nil
}

var emptyExportMetricsServiceResponse = func() []byte {
	body := collpb.ExportMetricsServiceResponse{}
	r, err := proto.Marshal(&body)
	if err != nil {
		panic(err)
	}
	return r
}()

type HTTPResponseError struct {
	Err    error
	Status int
	Header http.Header
}

func (e *HTTPResponseError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Err)
}

func (e *HTTPResponseError) Unwrap() error { return e.Err }

// HTTPCollector is an OTLP HTTP server that collects all requests it receives.
type HTTPCollector struct {
	plainTextResponse bool

	headersMu sync.Mutex
	headers   http.Header
	storage   *Storage

	resultCh <-chan ExportResult
	listener net.Listener
	srv      *http.Server
}

// NewHTTPCollector returns a *HTTPCollector that is listening at the provided
// endpoint.
//
// If endpoint is an empty string, the returned collector will be listening on
// the localhost interface at an OS chosen port, not use TLS, and listen at the
// default OTLP metric endpoint path ("/v1/metrics"). If the endpoint contains
// a prefix of "https" the server will generate weak self-signed TLS
// certificates and use them to server data. If the endpoint contains a path,
// that path will be used instead of the default OTLP metric endpoint path.
//
// If errCh is not nil, the collector will respond to HTTP requests with errors
// sent on that channel. This means that if errCh is not nil Export calls will
// block until an error is received.
func NewHTTPCollector(
	endpoint string,
	resultCh <-chan ExportResult,
	opts ...func(*HTTPCollector),
) (*HTTPCollector, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if u.Host == "" {
		u.Host = "localhost:0"
	}
	if u.Path == "" {
		u.Path = oconf.DefaultMetricsPath
	}

	c := &HTTPCollector{
		headers:  http.Header{},
		storage:  NewStorage(),
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
		cert, err := weakCertificate()
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

// WithHTTPCollectorRespondingPlainText makes the HTTPCollector return
// a plaintext, instead of protobuf, response.
func WithHTTPCollectorRespondingPlainText() func(*HTTPCollector) {
	return func(s *HTTPCollector) {
		s.plainTextResponse = true
	}
}

// Shutdown shuts down the HTTP server closing all open connections and
// listeners.
func (c *HTTPCollector) Shutdown(ctx context.Context) error {
	return c.srv.Shutdown(ctx)
}

// Addr returns the net.Addr c is listening at.
func (c *HTTPCollector) Addr() net.Addr {
	return c.listener.Addr()
}

// Collect returns the Storage holding all collected requests.
func (c *HTTPCollector) Collect() *Storage {
	return c.storage
}

// Headers returns the headers received for all requests.
func (c *HTTPCollector) Headers() map[string][]string {
	// Makes a copy.
	c.headersMu.Lock()
	defer c.headersMu.Unlock()
	return c.headers.Clone()
}

func (c *HTTPCollector) handler(w http.ResponseWriter, r *http.Request) {
	c.respond(w, c.record(r))
}

func (c *HTTPCollector) record(r *http.Request) ExportResult {
	// Currently only supports protobuf.
	if v := r.Header.Get("Content-Type"); v != "application/x-protobuf" {
		err := fmt.Errorf("content-type not supported: %s", v)
		return ExportResult{Err: err}
	}

	body, err := c.readBody(r)
	if err != nil {
		return ExportResult{Err: err}
	}
	pbRequest := &collpb.ExportMetricsServiceRequest{}
	err = proto.Unmarshal(body, pbRequest)
	if err != nil {
		return ExportResult{
			Err: &HTTPResponseError{
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
	return ExportResult{Err: err}
}

func (c *HTTPCollector) readBody(r *http.Request) (body []byte, err error) {
	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			_ = reader.Close()
			return nil, &HTTPResponseError{
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
			err = &HTTPResponseError{
				Err:    cErr,
				Status: http.StatusInternalServerError,
			}
		}
	}()
	body, err = io.ReadAll(reader)
	if err != nil {
		err = &HTTPResponseError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}
	return body, err
}

func (c *HTTPCollector) respond(w http.ResponseWriter, resp ExportResult) {
	if resp.Err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		var e *HTTPResponseError
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
		_, _ = w.Write(emptyExportMetricsServiceResponse)
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
func weakCertificate() (tls.Certificate, error) {
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
