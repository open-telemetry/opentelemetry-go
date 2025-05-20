// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "messaging" namespace.
package messagingconv

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}
	recOptPool = &sync.Pool{New: func() any { return &[]metric.RecordOption{} }}
)

// ErrorTypeAttr is an attribute conforming to the error.type semantic
// conventions. It represents the describes a class of error the operation ended
// with.
type ErrorTypeAttr string

var (
	// ErrorTypeOther is a fallback error value to be used when the instrumentation
	// doesn't define a custom value.
	ErrorTypeOther ErrorTypeAttr = "_OTHER"
)

// OperationTypeAttr is an attribute conforming to the messaging.operation.type
// semantic conventions. It represents a string identifying the type of the
// messaging operation.
type OperationTypeAttr string

var (
	// OperationTypeCreate is a message is created. "Create" spans always refer to a
	// single message and are used to provide a unique creation context for messages
	// in batch sending scenarios.
	OperationTypeCreate OperationTypeAttr = "create"
	// OperationTypeSend is the one or more messages are provided for sending to an
	// intermediary. If a single message is sent, the context of the "Send" span can
	// be used as the creation context and no "Create" span needs to be created.
	OperationTypeSend OperationTypeAttr = "send"
	// OperationTypeReceive is the one or more messages are requested by a consumer.
	// This operation refers to pull-based scenarios, where consumers explicitly
	// call methods of messaging SDKs to receive messages.
	OperationTypeReceive OperationTypeAttr = "receive"
	// OperationTypeProcess is the one or more messages are processed by a consumer.
	OperationTypeProcess OperationTypeAttr = "process"
	// OperationTypeSettle is the one or more messages are settled.
	OperationTypeSettle OperationTypeAttr = "settle"
)

// SystemAttr is an attribute conforming to the messaging.system semantic
// conventions. It represents the messaging system as identified by the client
// instrumentation.
type SystemAttr string

var (
	// SystemActiveMQ is the apache ActiveMQ.
	SystemActiveMQ SystemAttr = "activemq"
	// SystemAWSSQS is the amazon Simple Queue Service (SQS).
	SystemAWSSQS SystemAttr = "aws_sqs"
	// SystemEventGrid is the azure Event Grid.
	SystemEventGrid SystemAttr = "eventgrid"
	// SystemEventHubs is the azure Event Hubs.
	SystemEventHubs SystemAttr = "eventhubs"
	// SystemServiceBus is the azure Service Bus.
	SystemServiceBus SystemAttr = "servicebus"
	// SystemGCPPubSub is the google Cloud Pub/Sub.
	SystemGCPPubSub SystemAttr = "gcp_pubsub"
	// SystemJMS is the java Message Service.
	SystemJMS SystemAttr = "jms"
	// SystemKafka is the apache Kafka.
	SystemKafka SystemAttr = "kafka"
	// SystemRabbitMQ is the rabbitMQ.
	SystemRabbitMQ SystemAttr = "rabbitmq"
	// SystemRocketMQ is the apache RocketMQ.
	SystemRocketMQ SystemAttr = "rocketmq"
	// SystemPulsar is the apache Pulsar.
	SystemPulsar SystemAttr = "pulsar"
)

// ClientConsumedMessages is an instrument used to record metric values
// conforming to the "messaging.client.consumed.messages" semantic conventions.
// It represents the number of messages that were delivered to the application.
type ClientConsumedMessages struct {
	metric.Int64Counter
}

// NewClientConsumedMessages returns a new ClientConsumedMessages instrument.
func NewClientConsumedMessages(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientConsumedMessages, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConsumedMessages{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"messaging.client.consumed.messages",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Number of messages that were delivered to the application."),
			metric.WithUnit("{message}"),
		}, opt...)...,
	)
	if err != nil {
	    return ClientConsumedMessages{noop.Int64Counter{}}, err
	}
	return ClientConsumedMessages{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConsumedMessages) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientConsumedMessages) Name() string {
	return "messaging.client.consumed.messages"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConsumedMessages) Unit() string {
	return "{message}"
}

// Description returns the semantic convention description of the instrument
func (ClientConsumedMessages) Description() string {
	return "Number of messages that were delivered to the application."
}

// Add adds incr to the existing count.
//
// The operationName is the the system-specific name of the messaging operation.
//
// The system is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
//
// Records the number of messages pulled from the broker or number of messages
// dispatched to the application in push-based scenarios.
// The metric SHOULD be reported once per message delivery. For example, if
// receiving and processing operations are both instrumented for a single message
// delivery, this counter is incremented when the message is received and not
// reported when it is processed.
func (m ClientConsumedMessages) Add(
	ctx context.Context,
	incr int64,
	operationName string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientConsumedMessages) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrConsumerGroupName returns an optional attribute for the
// "messaging.consumer.group.name" semantic convention. It represents the name of
// the consumer group with which a consumer is associated.
func (ClientConsumedMessages) AttrConsumerGroupName(val string) attribute.KeyValue {
	return attribute.String("messaging.consumer.group.name", val)
}

// AttrDestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ClientConsumedMessages) AttrDestinationName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.name", val)
}

// AttrDestinationSubscriptionName returns an optional attribute for the
// "messaging.destination.subscription.name" semantic convention. It represents
// the name of the destination subscription from which a message is consumed.
func (ClientConsumedMessages) AttrDestinationSubscriptionName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.subscription.name", val)
}

// AttrDestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ClientConsumedMessages) AttrDestinationTemplate(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.template", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientConsumedMessages) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrDestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ClientConsumedMessages) AttrDestinationPartitionID(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.partition.id", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConsumedMessages) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ClientOperationDuration is an instrument used to record metric values
// conforming to the "messaging.client.operation.duration" semantic conventions.
// It represents the duration of messaging operation initiated by a producer or
// consumer client.
type ClientOperationDuration struct {
	metric.Float64Histogram
}

// NewClientOperationDuration returns a new ClientOperationDuration instrument.
func NewClientOperationDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientOperationDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientOperationDuration{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"messaging.client.operation.duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Duration of messaging operation initiated by a producer or consumer client."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return ClientOperationDuration{noop.Float64Histogram{}}, err
	}
	return ClientOperationDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientOperationDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientOperationDuration) Name() string {
	return "messaging.client.operation.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientOperationDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientOperationDuration) Description() string {
	return "Duration of messaging operation initiated by a producer or consumer client."
}

// Record records val to the current distribution.
//
// The operationName is the the system-specific name of the messaging operation.
//
// The system is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
//
// This metric SHOULD NOT be used to report processing duration - processing
// duration is reported in `messaging.process.duration` metric.
func (m ClientOperationDuration) Record(
	ctx context.Context,
	val float64,
	operationName string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientOperationDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrConsumerGroupName returns an optional attribute for the
// "messaging.consumer.group.name" semantic convention. It represents the name of
// the consumer group with which a consumer is associated.
func (ClientOperationDuration) AttrConsumerGroupName(val string) attribute.KeyValue {
	return attribute.String("messaging.consumer.group.name", val)
}

// AttrDestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ClientOperationDuration) AttrDestinationName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.name", val)
}

// AttrDestinationSubscriptionName returns an optional attribute for the
// "messaging.destination.subscription.name" semantic convention. It represents
// the name of the destination subscription from which a message is consumed.
func (ClientOperationDuration) AttrDestinationSubscriptionName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.subscription.name", val)
}

// AttrDestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ClientOperationDuration) AttrDestinationTemplate(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.template", val)
}

// AttrOperationType returns an optional attribute for the
// "messaging.operation.type" semantic convention. It represents a string
// identifying the type of the messaging operation.
func (ClientOperationDuration) AttrOperationType(val OperationTypeAttr) attribute.KeyValue {
	return attribute.String("messaging.operation.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientOperationDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrDestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ClientOperationDuration) AttrDestinationPartitionID(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.partition.id", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientOperationDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ClientSentMessages is an instrument used to record metric values conforming to
// the "messaging.client.sent.messages" semantic conventions. It represents the
// number of messages producer attempted to send to the broker.
type ClientSentMessages struct {
	metric.Int64Counter
}

// NewClientSentMessages returns a new ClientSentMessages instrument.
func NewClientSentMessages(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientSentMessages, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientSentMessages{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"messaging.client.sent.messages",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Number of messages producer attempted to send to the broker."),
			metric.WithUnit("{message}"),
		}, opt...)...,
	)
	if err != nil {
	    return ClientSentMessages{noop.Int64Counter{}}, err
	}
	return ClientSentMessages{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientSentMessages) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientSentMessages) Name() string {
	return "messaging.client.sent.messages"
}

// Unit returns the semantic convention unit of the instrument
func (ClientSentMessages) Unit() string {
	return "{message}"
}

// Description returns the semantic convention description of the instrument
func (ClientSentMessages) Description() string {
	return "Number of messages producer attempted to send to the broker."
}

// Add adds incr to the existing count.
//
// The operationName is the the system-specific name of the messaging operation.
//
// The system is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
//
// This metric MUST NOT count messages that were created but haven't yet been
// sent.
func (m ClientSentMessages) Add(
	ctx context.Context,
	incr int64,
	operationName string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientSentMessages) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrDestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ClientSentMessages) AttrDestinationName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.name", val)
}

// AttrDestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ClientSentMessages) AttrDestinationTemplate(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.template", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientSentMessages) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrDestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ClientSentMessages) AttrDestinationPartitionID(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.partition.id", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientSentMessages) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ProcessDuration is an instrument used to record metric values conforming to
// the "messaging.process.duration" semantic conventions. It represents the
// duration of processing operation.
type ProcessDuration struct {
	metric.Float64Histogram
}

// NewProcessDuration returns a new ProcessDuration instrument.
func NewProcessDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ProcessDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ProcessDuration{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"messaging.process.duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Duration of processing operation."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return ProcessDuration{noop.Float64Histogram{}}, err
	}
	return ProcessDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ProcessDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ProcessDuration) Name() string {
	return "messaging.process.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ProcessDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ProcessDuration) Description() string {
	return "Duration of processing operation."
}

// Record records val to the current distribution.
//
// The operationName is the the system-specific name of the messaging operation.
//
// The system is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
//
// This metric MUST be reported for operations with `messaging.operation.type`
// that matches `process`.
func (m ProcessDuration) Record(
	ctx context.Context,
	val float64,
	operationName string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ProcessDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrConsumerGroupName returns an optional attribute for the
// "messaging.consumer.group.name" semantic convention. It represents the name of
// the consumer group with which a consumer is associated.
func (ProcessDuration) AttrConsumerGroupName(val string) attribute.KeyValue {
	return attribute.String("messaging.consumer.group.name", val)
}

// AttrDestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ProcessDuration) AttrDestinationName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.name", val)
}

// AttrDestinationSubscriptionName returns an optional attribute for the
// "messaging.destination.subscription.name" semantic convention. It represents
// the name of the destination subscription from which a message is consumed.
func (ProcessDuration) AttrDestinationSubscriptionName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.subscription.name", val)
}

// AttrDestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ProcessDuration) AttrDestinationTemplate(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.template", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ProcessDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrDestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ProcessDuration) AttrDestinationPartitionID(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.partition.id", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ProcessDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}