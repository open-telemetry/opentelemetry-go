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

package statsd // import "go.opentelemetry.io/otel/exporter/metric/statsd"

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/sdk/export"
)

type (
	Config struct {
		// e.g., udp://host:port
		//       tcp://host:port
		//       unix:///socket/path
		URL string

		// For packet-oriented transports, this limits the packet size.
		MaxPacketSize int
	}

	Exporter struct {
		config Config

		buffer bytes.Buffer
	}

	formatCode string
)

var (
	_ export.MetricExporter     = &Exporter{}
	_ export.MetricLabelEncoder = &Exporter{}
)

func New(config Config) *Exporter {

}

func (e *Exporter) EncodeLabels([]core.KeyValue) string {
	return ""
}

func (e *Exporter) Export(_ context.Context, producer export.MetricProducer) {
	producer.Foreach(e.exportOne)
}

func (e *Exporter) exportOne(MetricAggregator, ProducedRecord) {
}

// @@@ DEAD CODE BELOW

const (
	formatCounter   formatCode = "c"
	formatHistogram formatCode = "h"
	formatGauge     formatCode = "g"
	formatTiming    formatCode = "ms"
)

func (sink *statsdSink) sendToFlusher() {
	buf := sink.statBufferPool.Get()
	// metric:value|type|#tag1:value1,tag2:value2
	// we use buf.WriteString instead of Fprintf because it's faster
	// as per documentation, WriteString never returns an error, so we ignore it here
	_, _ = buf.WriteString(sink.namespace)
	_, _ = buf.WriteString(name)
	_, _ = buf.WriteString(":")
	writeFloat(buf, value)
	_, _ = buf.WriteString("|")
	_, _ = buf.WriteString(string(statsdFormat))

	if len(tags) > 0 && len(tagString) > 0 {
		log.Warning(ctx, "Tags provided in multiple formats, one set will be dropped")
	}

	if len(tags) > 0 {
		sink.serializeTagsToBuffer(tags, buf)
	} else if len(tagString) > 0 {
		_, _ = buf.WriteString(tagString)
	}

}

func writeFloat(buf *bytes.Buffer, f float64) {
	var bs [128]byte
	fmtd := strconv.AppendFloat(bs[:0], f, 'g', -1, 64)
	_, _ = buf.Write(fmtd)
}

func (sink *statsdSink) flusher(ctx context.Context) error {
	defer close(sink.flushes)

	var currentConn net.Conn
	defer func() {
		var err error
		if currentConn != nil {
			if err = currentConn.Close(); err != nil {
				log.Errorf(ctx, "error while closing connection to statsd: %v", err)
			}
		}
	}()
	nextFlush := time.After(sink.flushInterval)

	buffer := &bytes.Buffer{}
	flushBuffer := func() {
		ctx := context.Background()
		var err error
		if sink.config.AggregateMetrics.VerboseTracing {
			var finish trace.FinishFunc
			ctx, finish = trace.Start(ctx, "statsdSink.flusher.flushBuffer")
			defer finish(&err)
		}

		// Make sure to reset the buffer before returning early to avoid leaks
		data := buffer.Next(buffer.Len())
		buffer.Reset()

		if currentConn == nil {
			var conn net.Conn
			conn, err = sink.dialer()
			if err != nil {
				log.Errorf(ctx, "error while connecting to statsd: %v", err)
				return
			}
			currentConn = conn
		}

		_ = currentConn.SetWriteDeadline(time.Now().Add(writeTimeout))
		var written int
		for written = 0; written < len(data); {
			var n int
			n, err = sink.writeWithTracing(ctx, currentConn, data[written:])

			if err != nil {
				log.Errorf(ctx, "error while writing to statsd: %v", err)
				currentConn = nil
				return
			}
			written += n
		}
		// This was done to avoid high cardinality tags as the desired information was to distinguish
		// between batches that require multiple sends vs. batches that only require a single send.
		if sink.config.AggregateMetrics.VerboseTracing {
			approxBytes := func(bytes int) string {
				operator := "<="
				if bytes > batchSizeBytes {
					operator = ">"
				}
				return fmt.Sprintf("%v %v", operator, batchSizeBytes)
			}
			trace.SetTag(ctx, tag.String("metrics_bytes_written", approxBytes(written)))
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case stat, ok := <-sink.toFlusher:
			if !ok {
				if buffer.Len() > 0 {
					flushBuffer()
				}
				return ctx.Err()
			}

			_, _ = stat.WriteTo(buffer)
			sink.statBufferPool.Put(stat)
			_, _ = buffer.WriteString("\n")

			if buffer.Len() >= batchSizeBytes {
				flushBuffer()
			}
			// TODO: change this to use a jittered ticker that drops ticks?
		case <-nextFlush:
			flushBuffer()
			nextFlush = time.After(sink.flushInterval)
		}
	}
}

func (sink *statsdSink) serializeTagsToString(tags []instrument.Tags) string {
	buf := sink.tagBufferPool.Get()
	defer sink.tagBufferPool.Put(buf)
	sink.serializeTagsToBuffer(tags, buf)
	return buf.String()
}

func (sink *statsdSink) serializeTagsToBuffer(tags []instrument.Tags, buf *bytes.Buffer) {
	// This function should always return a non-nil value, so get a new buffer from the shared pool if needed
	if len(tags) == 0 {
		return
	}

	excludedTags := sink.LoadExcluded()
	delimiter := "|#"
	for _, tags := range tags {
		for _, t := range tags {
			if !excludedTags.IsExcluded(t) {
				_, _ = buf.WriteString(delimiter)
				t.WriteTo(buf)
				delimiter = ","
			}
		}
	}
}

// TODO: dial async so if statsd is down, we don't error
func NewStatsdSinkUDS(namespace string, cfg DatadogConfig) Sink {
	addr := cfg.Address()
	if addr == "" {
		return &nullSink{}, nil
	}
	var dialer func() (net.Conn, error)
	if strings.HasPrefix(addr, "unix://") {
		addr = strings.TrimPrefix(addr, "unix://")
		dialer = func() (net.Conn, error) {
			return net.Dial("unixgram", addr)
		}
	} else {
		dialer = func() (net.Conn, error) {
			return net.Dial("udp", addr)
		}
	}

	return newStatsdSinkWithDialer(context.Background(), namespace, cfg, clockwork.NewRealClock(), dialer)
}
