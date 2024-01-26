// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
)

var (
	testTime     = time.Date(1988, time.November, 17, 0, 0, 0, 0, time.UTC)
	testSeverity = SeverityInfo
	testString   = "message"
	testFloat    = 1.2345
	testInt      = 32768
	testBool     = true
)

func TestRecordTimestamp(t *testing.T) {
	r := Record{}

	r.SetTimestamp(testTime)

	assert.Equal(t, testTime, r.Timestamp())
}

func TestRecordObservedTimestamp(t *testing.T) {
	r := Record{}

	r.SetObservedTimestamp(testTime)

	assert.Equal(t, testTime, r.ObservedTimestamp())
}

func TestRecordSeverity(t *testing.T) {
	r := Record{}

	r.SetSeverity(testSeverity)

	assert.Equal(t, testSeverity, r.Severity())
}

func TestRecordSeverityText(t *testing.T) {
	r := Record{}

	r.SetSeverityText(testString)

	assert.Equal(t, testString, r.SeverityText())
}

func TestRecordBody(t *testing.T) {
	r := Record{}
	body := StringValue(testString)

	r.SetBody(body)

	assert.Equal(t, body, r.Body())
}

func TestRecordAttributes(t *testing.T) {
	r := Record{}
	attrs := []KeyValue{
		String("k1", testString),
		Float64("k2", testFloat),
		Int("k3", testInt),
		Bool("k4", testBool),
		String("k5", testString),
		Float64("k6", testFloat),
		Int("k7", testInt),
		Bool("k8", testBool),
	}
	r.AddAttributes(attrs...)

	assert.Equal(t, len(attrs), r.AttributesLen())

	var got []KeyValue
	r.WalkAttributes(func(kv KeyValue) bool {
		got = append(got, kv)
		return true
	})
	assert.Equal(t, attrs, got)

	testCases := []struct {
		name  string
		index int
	}{
		{
			name:  "front",
			index: 2,
		},
		{
			name:  "back",
			index: 6,
		},
	}
	for _, tc := range testCases {
		i := 0
		r.WalkAttributes(func(kv KeyValue) bool {
			i++
			return i < tc.index
		})
		assert.Equal(t, tc.index, i, "WalkAttributes early return for %s", tc.name)
	}
}

func TestRecordAttributesInvalid(t *testing.T) {
	r := Record{}
	attrs := []KeyValue{
		String("k1", testString),
		{},
		Int("k3", testInt),
		Bool("k4", testBool),
		String("k5", testString),
		Float64("k6", testFloat),
		Int("k7", testInt),
		{},
	}
	r.AddAttributes(attrs...)

	assert.Equal(t, len(attrs)-2, r.AttributesLen())
}

func TestRecordAliasingAndClone(t *testing.T) {
	defer func(orig otel.ErrorHandler) {
		otel.SetErrorHandler(orig)
	}(otel.GetErrorHandler())
	var errs []error
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		errs = append(errs, err)
	}))

	// Create a record whose Attrs overflow the inline array,
	// creating a slice in r.back.
	r1 := Record{}
	for i := 0; i < attributesInlineCount+1; i++ {
		r1.AddAttributes(Int("k", i))
	}

	// Ensure that r1.back's capacity exceeds its length.
	b := make([]KeyValue, len(r1.back), len(r1.back)+1)
	copy(b, r1.back)
	r1.back = b

	// Make a copy that shares state.
	// Adding to both should emit an special error for the second call.
	r2 := r1
	r1AttrsBefore := attrsSlice(r1)
	r1.AddAttributes(Int("p", 0))
	assert.Zero(t, errs)
	r2.AddAttributes(Int("p", 1))
	assert.Equal(t, []error{errUnsafeAddAttrs}, errs, "sends an error via ErrorHandler when a dirty AddAttribute is detected")
	errs = nil
	assert.Equal(t, append(r1AttrsBefore, Int("p", 0)), attrsSlice(r1))
	assert.Equal(t, append(r1AttrsBefore, Int("p", 1)), attrsSlice(r2))

	// Adding to a clone is fine.
	r1Attrs := attrsSlice(r1)
	r3 := r1.Clone()
	assert.Equal(t, r1Attrs, attrsSlice(r3))
	r3.AddAttributes(Int("p", 2))
	assert.Zero(t, errs)
	assert.Equal(t, r1Attrs, attrsSlice(r1), "r1 is unchanged")
	assert.Equal(t, append(r1Attrs, Int("p", 2)), attrsSlice(r3))
	r3.SetSeverity(SeverityDebug)
	assert.NotEqual(t, r3.Severity(), r1.Severity(), "r1 is unchanged")
}

func attrsSlice(r Record) []KeyValue {
	var attrs []KeyValue
	r.WalkAttributes(func(kv KeyValue) bool {
		attrs = append(attrs, kv)
		return true
	})
	return attrs
}
