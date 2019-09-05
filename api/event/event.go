package event // import "go.opentelemetry.io/api/event"

import (
	"go.opentelemetry.io/api/core"
	"time"
)

// Event describes an event with a message string, timestamp and set of Attributes.
type Event struct {
	Message    string
	Time       time.Time
	Attributes []core.KeyValue
}
