// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

func verifyRing(t *testing.T, r *ring, num, sum int) {
	// Length.
	assert.Equal(t, num, r.Len(), "r.Len()")

	// Iteration.
	var n, s int
	r.Do(func(v Record) {
		n++
		body := v.Body()
		if body.Kind() != log.KindEmpty {
			s += int(body.AsInt64())
		}
	})
	assert.Equal(t, num, n, "number of forward iterations")
	if sum >= 0 {
		assert.Equal(t, sum, s, "forward ring sum")
	}

	if r == nil {
		return
	}

	// Connections.
	if r.next != nil {
		var p *ring // previous element.
		for q := r; p == nil || q != r; q = q.next {
			if p != nil {
				assert.Equalf(t, p, q.prev, "prev = %p, expected q.prev = %p", p, q.prev)
			}
			p = q
		}
		assert.Equalf(t, p, r.prev, "prev = %p, expected r.prev = %p", p, r.prev)
	}

	// Next, Prev.
	assert.Equal(t, r.next, r.Next(), "r.Next() != r.next")
	assert.Equal(t, r.prev, r.Prev(), "r.Prev() != r.prev")
}

func TestNewRing(t *testing.T) {
	for i := range 10 {
		// Empty value.
		r := newRing(i)
		verifyRing(t, r, i, -1)
	}

	for n := range 10 {
		r := newRing(n)
		for i := 1; i <= n; i++ {
			var rec Record
			rec.SetBody(log.IntValue(i))
			r.Value = rec
			r = r.Next()
		}

		sum := (n*n + n) / 2
		verifyRing(t, r, n, sum)
	}
}

func TestEmptyRing(t *testing.T) {
	var rNext, rPrev ring
	verifyRing(t, rNext.Next(), 1, 0)
	verifyRing(t, rPrev.Prev(), 1, 0)

	var rLen, rDo *ring
	assert.Equal(t, 0, rLen.Len(), "Len()")
	rDo.Do(func(Record) { assert.Fail(t, "Do func arg called") })
}
