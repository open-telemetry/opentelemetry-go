// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"
)

func BenchmarkEvent(b *testing.B) {
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
			var e log.Event
			e.SetTimestamp(y2k)
			tStamp = e.Timestamp()
		}
	})

	b.Run("ObservedTimestamp", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var e log.Event
			e.SetObservedTimestamp(y2k)
			tStamp = e.ObservedTimestamp()
		}
	})

	b.Run("Severity", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var e log.Event
			e.SetSeverity(log.SeverityDebug)
			sev = e.Severity()
		}
	})

	b.Run("SeverityText", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var e log.Event
			e.SetSeverityText("text")
			text = e.SeverityText()
		}
	})

	bodyVal := log.BoolValue(true)
	b.Run("Body", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			var e log.Event
			e.SetBody(bodyVal)
			body = e.Body()
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
				var e log.Event
				e.AddAttributes(attrs5...)
				n = e.AttributesLen()
				e.WalkAttributes(walk)
			}
		})
		b.Run("10", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var e log.Event
				e.AddAttributes(attrs10...)
				n = e.AttributesLen()
				e.WalkAttributes(walk)
			}
		})
	})

	// Convince the linter these values are used.
	_, _, _, _, _, _ = tStamp, sev, text, body, attr, n
}
