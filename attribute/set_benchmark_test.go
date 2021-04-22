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

package attribute

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBoolKey(t *testing.T) {
	e1 := DefaultEncoder()
	e2 := newEncoderPrefix("1", e1)
	e3 := newEncoderPrefix("2", e1)

	s := NewSet(Bool("k1", true))

	done := make(chan struct{})

	runner := func(encoder Encoder, expected string) {
		for {
			require.Equal(t, expected, s.Encoded(encoder))
			select {
			case <-done:
				return
			default:
			}

		}
	}

	for i := 0; i < 100; i++ {
		go runner(e1, "k1=true")
		go runner(e2, "1k1=true")
		go runner(e3, "2k1=true")
	}

	<-time.After(time.Millisecond * 500)
	close(done)
}

func newEncoderAndSet() ([3]Encoder, Set) {
	e1 := DefaultEncoder()
	e2 := newEncoderPrefix("1", e1)
	e3 := newEncoderPrefix("2", e1)

	encoders := [3]Encoder{e1, e2, e3}

	s := NewSet(Bool("k1", true), String("k2", "v2"), Int("k3", 13))

	return encoders, s
}

func BenchmarkEncoded(b *testing.B) {
	b.Run("single encoder", func(b *testing.B) {
		b.ReportAllocs()

		encoders, s := newEncoderAndSet()
		for i := 0; i < b.N; i++ {
			_ = s.Encoded(encoders[0])
		}
	})

	b.Run("sequential encoders", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		encoders, s := newEncoderAndSet()
		for i := 0; i < b.N; i++ {
			_ = s.Encoded(encoders[i%3])
		}
	})

	b.Run("parallel single encoders", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			encoders, s := newEncoderAndSet()
			for pb.Next() {
				_ = s.Encoded(encoders[0])
			}
		})
	})

	b.Run("parallel encoders", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			encoders, s := newEncoderAndSet()
			for pb.Next() {
				_ = s.Encoded(encoders[i%3])
				i++
			}
		})
	})

}

type encoderPrefix struct {
	encoder Encoder
	prefix  string
	id      EncoderID
}

func (e encoderPrefix) Encode(iterator Iterator) string {
	return e.prefix + e.encoder.Encode(iterator)
}

func (e encoderPrefix) ID() EncoderID {
	return e.id
}

var _ Encoder = (*encoderPrefix)(nil)

func newEncoderPrefix(prefix string, encoder Encoder) Encoder {
	return &encoderPrefix{
		prefix:  prefix,
		encoder: encoder,
		id:      NewEncoderID(),
	}
}
