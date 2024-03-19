// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"fmt"
	"os"
	"strconv"

	"go.opentelemetry.io/otel"
)

type setting[T any] struct {
	Value T
	Set   bool
}

func newSetting[T any](value T) setting[T] {
	return setting[T]{Value: value, Set: true}
}

func (s setting[T]) Resolve(fn ...func(setting[T]) setting[T]) setting[T] {
	for _, f := range fn {
		s = f(s)
	}
	return s
}

func clearLessThanOne[T ~int | ~int64]() func(setting[T]) setting[T] {
	return func(s setting[T]) setting[T] {
		if s.Value < 1 {
			s.Value = 0
			s.Set = false
		}
		return s
	}
}

func getenv[T ~int | ~int64](key string) func(setting[T]) setting[T] {
	return func(s setting[T]) setting[T] {
		if s.Set {
			// Passed, valid, options have precedence.
			return s
		}

		if v := os.Getenv(key); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				otel.Handle(fmt.Errorf("invalid %s value %s: %w", key, v, err))
			} else {
				s.Value = T(n)
				s.Set = true
			}
		}
		return s
	}
}

func fallback[T any](val T) func(setting[T]) setting[T] {
	return func(s setting[T]) setting[T] {
		if !s.Set {
			s.Value = val
			s.Set = true
		}
		return s
	}
}
