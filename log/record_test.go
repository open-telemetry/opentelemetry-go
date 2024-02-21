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

package log_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/log"
)

var y2k = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestRecordTimestamp(t *testing.T) {
	var r log.Record
	r.SetTimestamp(y2k)
	assert.Equal(t, y2k, r.Timestamp())
}

func TestRecordObservedTimestamp(t *testing.T) {
	var r log.Record
	r.SetObservedTimestamp(y2k)
	assert.Equal(t, y2k, r.ObservedTimestamp())
}

func TestRecordSeverity(t *testing.T) {
	var r log.Record
	r.SetSeverity(log.SeverityInfo)
	assert.Equal(t, log.SeverityInfo, r.Severity())
}

func TestRecordSeverityText(t *testing.T) {
	const text = "testing text"

	var r log.Record
	r.SetSeverityText(text)
	assert.Equal(t, text, r.SeverityText())
}

func TestRecordBody(t *testing.T) {
	body := log.StringValue("testing body value")

	var r log.Record
	r.SetBody(body)
	assert.Equal(t, body, r.Body())
}

func TestRecordAttributes(t *testing.T) {
	attrs := []log.KeyValue{
		log.String("k1", "str"),
		log.Float64("k2", 1.0),
		log.Int("k3", 2),
		log.Bool("k4", true),
		log.Bytes("k5", []byte{1}),
		log.Slice("k6", log.IntValue(3)),
		log.Map("k7", log.Bool("sub1", true)),
		log.String("k8", "str"),
		log.Float64("k9", 1.0),
		log.Int("k10", 2),
		log.Bool("k11", true),
		log.Bytes("k12", []byte{1}),
		log.Slice("k13", log.IntValue(3)),
		log.Map("k14", log.Bool("sub1", true)),
		{}, // Empty.
	}

	var r log.Record
	r.AddAttributes(attrs...)
	require.Equal(t, len(attrs), r.AttributesLen())

	t.Run("Correctness", func(t *testing.T) {
		var i int
		r.WalkAttributes(func(kv log.KeyValue) bool {
			assert.Equal(t, attrs[i], kv)
			i++
			return true
		})
	})

	t.Run("WalkAttributes/Filtering", func(t *testing.T) {
		for i := 1; i <= len(attrs); i++ {
			var j int
			r.WalkAttributes(func(log.KeyValue) bool {
				j++
				return j < i
			})
			assert.Equal(t, i, j, "number of attributes walked incorrect")
		}
	})
}

func BenchmarkRecord(b *testing.B) {
	var (
		tStamp time.Time
		sev    log.Severity
		text   string
		body   log.Value
		attr   log.KeyValue
		n      int
	)

	b.Run("Timestamp", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var r log.Record
			r.SetTimestamp(y2k)
			tStamp = r.Timestamp()
		}
	})

	b.Run("ObservedTimestamp", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var r log.Record
			r.SetObservedTimestamp(y2k)
			tStamp = r.ObservedTimestamp()
		}
	})

	b.Run("Severity", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var r log.Record
			r.SetSeverity(log.SeverityDebug)
			sev = r.Severity()
		}
	})

	b.Run("SeverityText", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var r log.Record
			r.SetSeverityText("text")
			text = r.SeverityText()
		}
	})

	bodyVal := log.BoolValue(true)
	b.Run("Body", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var r log.Record
			r.SetBody(bodyVal)
			body = r.Body()
		}
	})

	attrs10 := []log.KeyValue{
		log.Bool("b1", true),
		log.Int("i1", 324),
		log.Float64("f1", -230.213),
		log.String("s1", "value1"),
		log.Map("m1", log.Slice("slice1", log.BoolValue(true))),
		log.Bool("b2", false),
		log.Int("i2", 39847),
		log.Float64("f2", 0.382964329),
		log.String("s2", "value2"),
		log.Map("m2", log.Slice("slice2", log.BoolValue(false))),
	}
	attrs5 := attrs10[:5]
	walk := func(kv log.KeyValue) bool {
		attr = kv
		return true
	}
	b.Run("Attributes", func(b *testing.B) {
		b.Run("5", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var r log.Record
				r.AddAttributes(attrs5...)
				n = r.AttributesLen()
				r.WalkAttributes(walk)
			}
		})
		b.Run("10", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var r log.Record
				r.AddAttributes(attrs10...)
				n = r.AttributesLen()
				r.WalkAttributes(walk)
			}
		})
	})

	_, _, _, _, _, _ = tStamp, sev, text, body, attr, n
}
