// Copyright The OpenTelemetry Authors
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

//go:build go1.18
// +build go1.18

package otest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otest"

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/oconf"
	collpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

var emptyExportMetricsServiceResponse = func() []byte {
	body := collpb.ExportMetricsServiceResponse{}
	r, err := proto.Marshal(&body)
	if err != nil {
		panic(err)
	}
	return r
}()

// Collector is the collection target a Client sends metric uploads to.
type Collector interface {
	Collect() *Storage
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
	headersMu sync.Mutex
	headers   http.Header
	storage   *Storage

	errCh    <-chan error
	listener net.Listener
	srv      *http.Server
}

// NewHTTPCollector returns a *HTTPCollector that is listening at the provided
// endpoint.
//
// If endpoint is an empty string, the returned collector will be listeing on
// the localhost interface at an OS chosen port.
//
// If errCh is not nil, the collector will respond to HTTP requests with errors
// sent on that channel. This means that if errCh is not nil Export calls will
// block until an error is received.
func NewHTTPCollector(endpoint string, errCh <-chan error) (*HTTPCollector, error) {
	if endpoint == "" {
		endpoint = "localhost:0"
	}

	c := &HTTPCollector{
		headers: http.Header{},
		storage: NewStorage(),
		errCh:   errCh,
	}

	var err error
	c.listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle(oconf.DefaultMetricsPath, http.HandlerFunc(c.handler))
	c.srv = &http.Server{Handler: mux}
	go func() { _ = c.srv.Serve(c.listener) }()
	return c, nil
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

func (c *HTTPCollector) record(r *http.Request) error {
	// Currently only supports protobuf.
	if v := r.Header.Get("Content-Type"); v != "application/x-protobuf" {
		return fmt.Errorf("content-type not supported: %s", v)
	}

	body, err := c.readBody(r)
	if err != nil {
		return err
	}
	pbRequest := &collpb.ExportMetricsServiceRequest{}
	err = proto.Unmarshal(body, pbRequest)
	if err != nil {
		return &HTTPResponseError{
			Err:    err,
			Status: http.StatusInternalServerError,
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

	if c.errCh != nil {
		err = <-c.errCh
	}
	return err
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

func (c *HTTPCollector) respond(w http.ResponseWriter, err error) {
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		var e *HTTPResponseError
		if errors.As(err, &e) {
			for k, vals := range e.Header {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			w.WriteHeader(e.Status)
			fmt.Fprintln(w, e.Error())
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(emptyExportMetricsServiceResponse)
}
