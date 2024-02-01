// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		{},
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

func TestRecordCloneConcurrentSafe(t *testing.T) {
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

	attrsBefore := attrsSlice(r1)

	// Changing this to r2 := r1 will make the test fail
	// and has a race condition.
	r2 := r1.Clone()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		r1.AddAttributes(Int("p", 1))
	}()
	go func() {
		defer wg.Done()
		r2.AddAttributes(Int("p", 2))
	}()
	wg.Wait()

	assert.Equal(t, append(attrsBefore, Int("p", 1)), attrsSlice(r1))
	assert.Equal(t, append(attrsBefore, Int("p", 2)), attrsSlice(r2))
}

func attrsSlice(r Record) []KeyValue {
	var attrs []KeyValue
	r.WalkAttributes(func(kv KeyValue) bool {
		attrs = append(attrs, kv)
		return true
	})
	return attrs
}
