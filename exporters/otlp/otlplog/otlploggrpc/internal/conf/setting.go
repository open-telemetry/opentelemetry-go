// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package conf // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/conf"

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
)

// Setting is a configuration setting value.
type Setting[T any] struct {
	Value T
	Set   bool
}

// NewSetting returns a new setting with the value set.
func NewSetting[T any](value T) Setting[T] {
	return Setting[T]{Value: value, Set: true}
}

// Resolver returns an updated setting after applying an resolution operation.
type Resolver[T any] func(Setting[T]) Setting[T]

// Resolve returns a resolved version of s.
//
// It will apply all the passed fn in the order provided, chaining together the
// return setting to the next input. The setting s is used as the initial
// argument to the first fn.
//
// Each fn needs to validate if it should apply given the Set state of the
// setting. This will not perform any checks on the set state when chaining
// function.
func (s Setting[T]) Resolve(fn ...Resolver[T]) Setting[T] {
	for _, f := range fn {
		s = f(s)
	}
	return s
}

// GetEnv returns a Resolver that will apply an environment variable value
// associated with the first set key to a setting value. The conv function is
// used to convert between the environment variable value and the setting type.
//
// If the input setting to the Resolver is set, the environment variable will
// not be applied.
//
// Any error returned from conv is sent to the OTel ErrorHandler and the
// setting will not be updated.
func GetEnv[T any](keys []string, conv func(string) (T, error)) Resolver[T] {
	return func(s Setting[T]) Setting[T] {
		if s.Set {
			// Passed, valid, options have precedence.
			return s
		}

		for _, key := range keys {
			if vStr := os.Getenv(key); vStr != "" {
				v, err := conv(vStr)
				if err == nil {
					s.Value = v
					s.Set = true
					break
				}
				otel.Handle(fmt.Errorf("invalid %s value %s: %w", key, vStr, err))
			}
		}
		return s
	}
}

// Fallback returns a resolve that will set a setting value to val if it is not
// already set.
//
// This is usually passed at the end of a resolver chain to ensure a default is
// applied if the setting has not already been set.
func Fallback[T any](val T) Resolver[T] {
	return func(s Setting[T]) Setting[T] {
		if !s.Set {
			s.Value = val
			s.Set = true
		}
		return s
	}
}
