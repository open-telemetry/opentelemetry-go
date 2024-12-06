// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/log"
)

func TestEventTimestamp(t *testing.T) {
	var e log.Event
	e.SetTimestamp(y2k)
	assert.Equal(t, y2k, e.Timestamp())
}

func TestEventObservedTimestamp(t *testing.T) {
	var e log.Event
	e.SetObservedTimestamp(y2k)
	assert.Equal(t, y2k, e.ObservedTimestamp())
}

func TestEventSeverity(t *testing.T) {
	var e log.Event
	e.SetSeverity(log.SeverityInfo)
	assert.Equal(t, log.SeverityInfo, e.Severity())
}

func TestEventSeverityText(t *testing.T) {
	const text = "testing text"

	var e log.Event
	e.SetSeverityText(text)
	assert.Equal(t, text, e.SeverityText())
}

func TestEventBody(t *testing.T) {
	body := log.StringValue("testing body value")

	var e log.Event
	e.SetBody(body)
	assert.Equal(t, body, e.Body())
}

func TestEventAttributes(t *testing.T) {
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

	var e log.Event
	e.AddAttributes(attrs...)
	require.Equal(t, len(attrs), e.AttributesLen())

	t.Run("Correctness", func(t *testing.T) {
		var i int
		e.WalkAttributes(func(kv log.KeyValue) bool {
			assert.Equal(t, attrs[i], kv)
			i++
			return true
		})
	})

	t.Run("WalkAttributes/Filtering", func(t *testing.T) {
		for i := 1; i <= len(attrs); i++ {
			var j int
			e.WalkAttributes(func(log.KeyValue) bool {
				j++
				return j < i
			})
			assert.Equal(t, i, j, "number of attributes walked incorrect")
		}
	})
}

func TestEventAllocationLimits(t *testing.T) {
	const runs = 5

	// Assign testing results to external scope so the compiler doesn't
	// optimize away the testing statements.
	var (
		tStamp time.Time
		sev    log.Severity
		text   string
		body   log.Value
		n      int
		attr   log.KeyValue
	)

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var e log.Event
		e.SetTimestamp(y2k)
		tStamp = e.Timestamp()
	}), "Timestamp")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var e log.Event
		e.SetObservedTimestamp(y2k)
		tStamp = e.ObservedTimestamp()
	}), "ObservedTimestamp")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var e log.Event
		e.SetSeverity(log.SeverityDebug)
		sev = e.Severity()
	}), "Severity")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var e log.Event
		e.SetSeverityText("severity text")
		text = e.SeverityText()
	}), "SeverityText")

	bodyVal := log.BoolValue(true)
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var e log.Event
		e.SetBody(bodyVal)
		body = e.Body()
	}), "Body")

	attrVal := []log.KeyValue{log.Bool("k", true), log.Int("i", 1)}
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var e log.Event
		e.AddAttributes(attrVal...)
		n = e.AttributesLen()
		e.WalkAttributes(func(kv log.KeyValue) bool {
			attr = kv
			return true
		})
	}), "Attributes")

	// Convince the linter these values are used.
	_, _, _, _, _, _ = tStamp, sev, text, body, n, attr
}
