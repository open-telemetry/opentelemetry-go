package event

import "testing"

func TestEventCarrier_Get(t *testing.T) {
	e := Event{
		EventType:  OrderEvent,
		Attributes: map[string]string{},
		Body:       nil,
	}
	carrier := NewEventCarrier(&e)

	carrier.Set("foo", "bar")
	if bar := carrier.Get("foo"); bar != "bar" {
		t.Error("value not found")
	}
}
