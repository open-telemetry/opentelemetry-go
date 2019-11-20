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

package statsd

// See https://github.com/b/statsd_spec for the best-available statsd
// syntax specification.  See also
// https://github.com/statsd/statsd/edit/master/docs/metric_types.md

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/unit"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type (
	// Config supports common options that apply to statsd exporters.
	Config struct {
		// URL describes the destination for exporting statsd data.
		// e.g., udp://host:port
		//       tcp://host:port
		//       unix:///socket/path
		URL string

		// Writer is an alternate to providing a URL.  When Writer is
		// non-nil, URL will be ignored and the exporter will write to
		// the configured Writer interface.
		Writer io.Writer

		// MaxPacketSize this limits the packet size for packet-oriented transports.
		MaxPacketSize int

		// TODO support Dial and Write timeouts
	}

	// Exporter is common type meant to implement concrete statsd
	// exporters.
	Exporter struct {
		adapter Adapter
		config  Config
		conn    net.Conn
		writer  io.Writer
		buffer  bytes.Buffer
	}

	// Adapter supports statsd syntax variations, primarily plain
	// statsd vs. dogstatsd.
	Adapter interface {
		AppendName(export.Record, *bytes.Buffer)
		AppendTags(export.Record, *bytes.Buffer)
	}
)

const (
	formatCounter   = "c"
	formatHistogram = "h"
	formatGauge     = "g"
	formatTiming    = "ms"

	MaxPacketSize = 1 << 16
)

var (
	_ export.Exporter = &Exporter{}

	ErrInvalidScheme = fmt.Errorf("Invalid statsd transport")
)

// NewExport returns a common implementation for exporters that Export
// statsd syntax.
func NewExporter(config Config, adapter Adapter) (*Exporter, error) {
	if config.MaxPacketSize <= 0 {
		config.MaxPacketSize = MaxPacketSize
	}
	var writer io.Writer
	var conn net.Conn
	var err error
	if config.Writer != nil {
		writer = config.Writer
	} else {
		conn, err = dial(config.URL)
		if conn != nil {
			writer = conn
		}
	}
	// TODO: If err != nil, we return it _with_ a valid exporter; the
	// exporter should attempt to re-dial if it's retryable.  Add a
	// Start() and Stop() API.
	return &Exporter{
		adapter: adapter,
		config:  config,
		conn:    conn,
		writer:  writer,
	}, err
}

// dial connects to a statsd service using several common network
// types.  Presently "udp" and "unix" datagram socket connections are
// supported.
func dial(endpoint string) (net.Conn, error) {
	dest, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// TODO: Support tcp destination, need configurable timeouts first.

	scheme := dest.Scheme
	switch scheme {
	case "udp", "udp4", "udp6":
		udpAddr, err := net.ResolveUDPAddr(scheme, dest.Host)
		locAddr := &net.UDPAddr{}
		if err != nil {
			return nil, err
		}
		conn, err := net.DialUDP(scheme, locAddr, udpAddr)
		if err != nil {
			return nil, err
		}
		return conn, err
	case "unix", "unixgram":
		scheme = "unixgram"
		locAddr := &net.UnixAddr{}

		sockAddr, err := net.ResolveUnixAddr(scheme, dest.Path)
		if err != nil {
			return nil, err
		}
		conn, err := net.DialUnix(scheme, locAddr, sockAddr)
		if err != nil {
			return nil, err
		}
		return conn, err
	}
	return nil, ErrInvalidScheme
}

// Export is common code for any statsd-based metric.Exporter implementation.
func (e *Exporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
	buf := &e.buffer
	buf.Reset()

	var aggErr error
	var sendErr error

	checkpointSet.ForEach(func(rec export.Record) {
		before := buf.Len()

		if err := e.formatMetric(rec, buf); err != nil && aggErr == nil {
			aggErr = err
			return
		}

		if buf.Len() < e.config.MaxPacketSize {
			return
		}
		if before == 0 {
			// A single metric >= packet size
			if err := e.send(buf.Bytes()); err != nil && sendErr == nil {
				sendErr = err
			}
			buf.Reset()
			return
		}

		// Send and copy the leftover
		if err := e.send(buf.Bytes()[:before]); err != nil && sendErr == nil {
			sendErr = err
		}

		leftover := buf.Len() - before

		copy(buf.Bytes()[0:leftover], buf.Bytes()[before:])

		buf.Truncate(leftover)
	})
	if err := e.send(buf.Bytes()); err != nil && sendErr == nil {
		sendErr = err
	}
	if sendErr != nil {
		return sendErr
	}
	return aggErr
}

// send writes a complete buffer to the writer as a blocking call.
func (e *Exporter) send(buf []byte) error {
	for len(buf) != 0 {
		n, err := e.writer.Write(buf)
		if err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

// formatMetric formats an individual export record.  For some records
// this will emit a single statistic, for some it will emit more than
// one.
func (e *Exporter) formatMetric(rec export.Record, buf *bytes.Buffer) error {
	desc := rec.Descriptor()
	agg := rec.Aggregator()

	// TODO handle non-Points Distribution/MaxSumCount by
	// formatting individual quantiles, the sum, and the count as
	// single statistics.  For the dogstatsd variation, assuming
	// open-source systems like Veneur add support, figure out the
	// proper encoding for "d"-type distribution data.

	if pts, ok := agg.(aggregator.Points); ok {
		var format string
		if desc.Unit() == unit.Milliseconds {
			format = formatTiming
		} else {
			format = formatHistogram
		}
		points, err := pts.Points()
		if err != nil {
			return err
		}
		for _, pt := range points {
			e.formatSingleStat(rec, pt, format, buf)
		}

	} else if sum, ok := agg.(aggregator.Sum); ok {
		sum, err := sum.Sum()
		if err != nil {
			return err
		}
		e.formatSingleStat(rec, sum, formatCounter, buf)

	} else if lv, ok := agg.(aggregator.LastValue); ok {
		lv, _, err := lv.LastValue()
		if err != nil {
			return err
		}
		e.formatSingleStat(rec, lv, formatGauge, buf)
	}
	return nil
}

// formatSingleStat encodes a single item of statsd data followed by a
// newline.
func (e *Exporter) formatSingleStat(rec export.Record, val core.Number, fmtStr string, buf *bytes.Buffer) {
	e.adapter.AppendName(rec, buf)
	_, _ = buf.WriteRune(':')
	writeNumber(buf, val, rec.Descriptor().NumberKind())
	_, _ = buf.WriteRune('|')
	_, _ = buf.WriteString(fmtStr)
	e.adapter.AppendTags(rec, buf)
	_, _ = buf.WriteRune('\n')
}

func writeNumber(buf *bytes.Buffer, num core.Number, kind core.NumberKind) {
	var tmp [128]byte
	var conv []byte
	switch kind {
	case core.Int64NumberKind:
		conv = strconv.AppendInt(tmp[:0], num.AsInt64(), 10)
	case core.Float64NumberKind:
		conv = strconv.AppendFloat(tmp[:0], num.AsFloat64(), 'g', -1, 64)
	case core.Uint64NumberKind:
		conv = strconv.AppendUint(tmp[:0], num.AsUint64(), 10)

	}
	_, _ = buf.Write(conv)
}
