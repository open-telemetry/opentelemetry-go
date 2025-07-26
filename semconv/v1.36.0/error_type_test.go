// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package semconv // import "go.opentelemetry.io/otel/semconv/v1.36.0"

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

type CustomError struct{}

func (CustomError) Error() string {
	return "custom error"
}

func TestErrorType(t *testing.T) {
	customErr := CustomError{}
	builtinErr := errors.New("something went wrong")
	var nilErr error

	wantCustomType := reflect.TypeOf(customErr)
	wantCustomStr := fmt.Sprintf("%s.%s", wantCustomType.PkgPath(), wantCustomType.Name())

	tests := []struct {
		name string
		err  error
		want attribute.KeyValue
	}{
		{
			name: "BuiltinError",
			err:  builtinErr,
			want: attribute.String("error.type", "*errors.errorString"),
		},
		{
			name: "CustomError",
			err:  customErr,
			want: attribute.String("error.type", wantCustomStr),
		},
		{
			name: "NilError",
			err:  nilErr,
			want: ErrorTypeOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorType(tt.err)
			if got != tt.want {
				t.Errorf("ErrorType(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
