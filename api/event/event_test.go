package event

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
)

func TestEventFields(t *testing.T) {
	k1v1 := key.New("key1").String("value1")
	k2v2 := key.New("key2").String("value2")

	newEvent := NewEvent("msg")
	newEvent.AddAttribute(k1v1)
	newEvent.AddAttribute(k2v2)

	if message := newEvent.Message(); message != "msg" {
		t.Errorf("error with event message: expected %s, given %s", "msg", message)
	}

	want := []core.KeyValue{k1v1, k2v2}
	attributes := newEvent.Attributes()

	if diff := cmp.Diff(want, attributes); diff != "" {
		t.Errorf("error with event attributes: %s", diff)
	}

	// Check if attributes are immutable
	attributes[0].Value = core.Value{String: "value0"}
	if diff := cmp.Diff(want, newEvent.Attributes()); diff != "" {
		t.Errorf("event attributes are not immutable: %s", diff)
	}
}
