package event

import "go.opentelemetry.io/otel/propagation"

type Event struct {
	EventType  EventType
	Attributes map[string]string
	Body       []byte
}

type EventType string

const (
	OrderEvent EventType = "Order"
)

var _ propagation.TextMapCarrier = (*EventCarrier)(nil)

type EventCarrier struct {
	event *Event
}

func NewEventCarrier(event *Event) *EventCarrier {
	return &EventCarrier{event: event}
}

func (e EventCarrier) Get(key string) string {
	if v, ok := e.event.Attributes[key]; ok {
		return v
	}
	return ""
}

func (e EventCarrier) Set(key string, value string) {
	e.event.Attributes[key] = value
}

func (e EventCarrier) Keys() []string {
	keys := make([]string, 0, len(e.event.Attributes))
	for _, key := range e.event.Attributes {
		keys = append(keys, key)
	}
	return keys
}
