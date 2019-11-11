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

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/unit"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
)

type (
	Config struct {
		// URL describes the
		//
		// e.g., udp://host:port
		//       tcp://host:port
		//       unix:///socket/path
		URL string

		// MaxPacketSize this limits the packet size for packet-oriented transports.
		MaxPacketSize int

		// TODO support Dial and Write timeouts
	}

	Exporter struct {
		adapter Adapter
		config  Config
		conn    net.Conn
		buffer  bytes.Buffer
	}

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

func NewExporter(config Config, adapter Adapter) (*Exporter, error) {
	if config.MaxPacketSize <= 0 {
		config.MaxPacketSize = MaxPacketSize
	}

	conn, err := dial(config.URL)
	return &Exporter{
		adapter: adapter,
		config:  config,
		conn:    conn,
	}, err
}

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

func (e *Exporter) Export(_ context.Context, producer export.Producer) error {
	buf := &e.buffer
	buf.Reset()

	var retErr error

	producer.Foreach(func(rec export.Record) {
		before := buf.Len()

		e.format(rec, buf)

		if buf.Len() < e.config.MaxPacketSize {
			return
		}
		if before == 0 {
			// A single metric >= packet size
			if err := e.send(buf.Bytes()); err != nil && retErr == nil {
				retErr = err
			}
			buf.Reset()
			return
		}

		// Send and copy the leftover
		if err := e.send(buf.Bytes()[:before]); err != nil && retErr == nil {
			retErr = err
		}

		leftover := buf.Len() - before

		copy(buf.Bytes()[0:leftover], buf.Bytes()[before:])

		buf.Truncate(leftover)
	})
	if err := e.send(buf.Bytes()); err != nil && retErr == nil {
		retErr = err
	}
	return retErr
}

func (e *Exporter) send(buf []byte) error {
	for len(buf) != 0 {
		n, err := e.conn.Write(buf)
		if err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

// For basic statsd, see
// https://github.com/statsd/statsd/edit/master/docs/metric_types.md
func (e *Exporter) format(rec export.Record, buf *bytes.Buffer) {
	desc := rec.Descriptor()
	kind := desc.NumberKind()
	agg := rec.Aggregator()

	if pts, ok := agg.(aggregator.Points); ok {
		var format string
		if desc.Unit() == unit.Milliseconds {
			format = formatHistogram
		} else {
			format = formatTiming
		}
		for _, pt := range pts.Points() {
			e.adapter.AppendName(rec, buf)
			_, _ = buf.WriteRune(':')
			writeNumber(buf, pt, kind)
			_, _ = buf.WriteRune('|')
			_, _ = buf.WriteString(format)
		}

	} else if sum, ok := agg.(aggregator.Sum); ok {
		e.adapter.AppendName(rec, buf)
		_, _ = buf.WriteRune(':')
		writeNumber(buf, sum.Sum(), kind)
		_, _ = buf.WriteRune('|')
		_, _ = buf.WriteString(formatCounter)

	} else if lv, ok := agg.(aggregator.LastValue); ok {
		e.adapter.AppendName(rec, buf)
		_, _ = buf.WriteRune(':')
		writeNumber(buf, lv.LastValue(), kind)
		_, _ = buf.WriteRune('|')
		_, _ = buf.WriteString(formatGauge)
	}

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
