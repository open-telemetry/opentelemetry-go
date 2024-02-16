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
		log.List("k6", log.IntValue(3)),
		log.Map("k7", log.Bool("sub1", true)),
		log.String("k8", "str"),
		log.Float64("k9", 1.0),
		log.Int("k10", 2),
		log.Bool("k11", true),
		log.Bytes("k12", []byte{1}),
		log.List("k13", log.IntValue(3)),
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
