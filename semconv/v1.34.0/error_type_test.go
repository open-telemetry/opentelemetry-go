package semconv

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

func TestErrorType_BuiltinError(t *testing.T) {
	err := errors.New("something went wrong")

	kv := ErrorType(err)
	expected := attribute.String("error.type", "*errors.errorString")

	if kv.Value.AsString() != expected.Value.AsString() {
		t.Errorf("Expected %v, got %v", expected, kv)
	}
}

func TestErrorType_CustomError(t *testing.T) {
	err := CustomError{}

	kv := ErrorType(err)
	expectedType := reflect.TypeOf(err)
	expectedStr := fmt.Sprintf("%s.%s", expectedType.PkgPath(), expectedType.Name())

	if kv.Value.AsString() != expectedStr {
		t.Errorf("Expected %s, got %s", expectedStr, kv.Value.AsString())
	}
}

func TestErrorType_Nil(t *testing.T) {
	var err error = nil

	kv := ErrorType(err)

	if kv != ErrorTypeOther {
		t.Errorf("Expected ErrorTypeOther, got %v", kv)
	}
}
