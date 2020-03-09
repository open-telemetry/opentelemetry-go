package jaeger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/apache/thrift/lib/go/thrift"

	gen "go.opentelemetry.io/otel/exporters/trace/jaeger/internal/gen-go/jaeger"
)

// batchUploader send a batch of spans to Jaeger
type batchUploader interface {
	upload(batch *gen.Batch) error
}

type EndpointOption func() (batchUploader, error)

// WithAgentEndpoint instructs exporter to send spans to jaeger-agent at this address.
// For example, localhost:6831.
func WithAgentEndpoint(agentEndpoint string) func() (batchUploader, error) {
	return func() (batchUploader, error) {
		if agentEndpoint == "" {
			return nil, errors.New("agentEndpoint must not be empty")
		}

		client, err := newAgentClientUDP(agentEndpoint, udpPacketMaxLength)
		if err != nil {
			return nil, err
		}

		return &agentUploader{client: client}, nil
	}
}

// WithCollectorEndpoint defines the full url to the Jaeger HTTP Thrift collector.
// For example, http://localhost:14268/api/traces
func WithCollectorEndpoint(collectorEndpoint string, options ...CollectorEndpointOption) func() (batchUploader, error) {
	return func() (batchUploader, error) {
		if collectorEndpoint == "" {
			return nil, errors.New("collectorEndpoint must not be empty")
		}

		o := &CollectorEndpointOptions{}
		for _, opt := range options {
			opt(o)
		}

		return &collectorUploader{
			endpoint: collectorEndpoint,
			username: o.username,
			password: o.password,
		}, nil
	}
}

type CollectorEndpointOption func(o *CollectorEndpointOptions)

type CollectorEndpointOptions struct {
	// username to be used if basic auth is required.
	username string

	// password to be used if basic auth is required.
	password string
}

// WithUsername sets the username to be used if basic auth is required.
func WithUsername(username string) func(o *CollectorEndpointOptions) {
	return func(o *CollectorEndpointOptions) {
		o.username = username
	}
}

// WithPassword sets the password to be used if basic auth is required.
func WithPassword(password string) func(o *CollectorEndpointOptions) {
	return func(o *CollectorEndpointOptions) {
		o.password = password
	}
}

// agentUploader implements batchUploader interface sending batches to
// Jaeger through the UDP agent.
type agentUploader struct {
	client *agentClientUDP
}

var _ batchUploader = (*agentUploader)(nil)

func (a *agentUploader) upload(batch *gen.Batch) error {
	return a.client.EmitBatch(batch)
}

// collectorUploader implements batchUploader interface sending batches to
// Jaeger through the collector http endpoint.
type collectorUploader struct {
	endpoint string
	username string
	password string
}

var _ batchUploader = (*collectorUploader)(nil)

func (c *collectorUploader) upload(batch *gen.Batch) error {
	body, err := serialize(batch)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.endpoint, body)
	if err != nil {
		return err
	}
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	req.Header.Set("Content-Type", "application/x-thrift")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	_, _ = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to upload traces; HTTP status code: %d", resp.StatusCode)
	}
	return nil
}

func serialize(obj thrift.TStruct) (*bytes.Buffer, error) {
	buf := thrift.NewTMemoryBuffer()
	if err := obj.Write(thrift.NewTBinaryProtocolTransport(buf)); err != nil {
		return nil, err
	}
	return buf.Buffer, nil
}
