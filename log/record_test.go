// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/cmplxattr"
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
	body := cmplxattr.StringValue("testing body value")

	var r log.Record
	r.SetBody(body)
	assert.Equal(t, body, r.Body())
}

func TestRecordAttributes(t *testing.T) {
	attrs := []cmplxattr.KeyValue{
		cmplxattr.String("k1", "str"),
		cmplxattr.Float64("k2", 1.0),
		cmplxattr.Int("k3", 2),
		cmplxattr.Bool("k4", true),
		cmplxattr.Bytes("k5", []byte{1}),
		cmplxattr.Slice("k6", cmplxattr.IntValue(3)),
		cmplxattr.Map("k7", cmplxattr.Bool("sub1", true)),
		cmplxattr.String("k8", "str"),
		cmplxattr.Float64("k9", 1.0),
		cmplxattr.Int("k10", 2),
		cmplxattr.Bool("k11", true),
		cmplxattr.Bytes("k12", []byte{1}),
		cmplxattr.Slice("k13", cmplxattr.IntValue(3)),
		cmplxattr.Map("k14", cmplxattr.Bool("sub1", true)),
		{}, // Empty.
	}

	var r log.Record
	r.AddAttributes(attrs...)
	require.Equal(t, len(attrs), r.AttributesLen())

	t.Run("Correctness", func(t *testing.T) {
		var i int
		r.WalkAttributes(func(kv cmplxattr.KeyValue) bool {
			assert.Equal(t, attrs[i], kv)
			i++
			return true
		})
	})

	t.Run("WalkAttributes/Filtering", func(t *testing.T) {
		for i := 1; i <= len(attrs); i++ {
			var j int
			r.WalkAttributes(func(cmplxattr.KeyValue) bool {
				j++
				return j < i
			})
			assert.Equal(t, i, j, "number of attributes walked incorrect")
		}
	})
}

func TestRecordAllocationLimits(t *testing.T) {
	const runs = 5

	// Assign testing results to external scope so the compiler doesn't
	// optimize away the testing statements.
	var (
		tStamp time.Time
		sev    log.Severity
		text   string
		body   cmplxattr.Value
		n      int
		attr   cmplxattr.KeyValue
	)

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var r log.Record
		r.SetTimestamp(y2k)
		tStamp = r.Timestamp()
	}), "Timestamp")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var r log.Record
		r.SetObservedTimestamp(y2k)
		tStamp = r.ObservedTimestamp()
	}), "ObservedTimestamp")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var r log.Record
		r.SetSeverity(log.SeverityDebug)
		sev = r.Severity()
	}), "Severity")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var r log.Record
		r.SetSeverityText("severity text")
		text = r.SeverityText()
	}), "SeverityText")

	bodyVal := cmplxattr.BoolValue(true)
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var r log.Record
		r.SetBody(bodyVal)
		body = r.Body()
	}), "Body")

	attrVal := []cmplxattr.KeyValue{cmplxattr.Bool("k", true), cmplxattr.Int("i", 1)}
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		var r log.Record
		r.AddAttributes(attrVal...)
		n = r.AttributesLen()
		r.WalkAttributes(func(kv cmplxattr.KeyValue) bool {
			attr = kv
			return true
		})
	}), "Attributes")

	// Convince the linter these values are used.
	_, _, _, _, _, _ = tStamp, sev, text, body, n, attr
}
