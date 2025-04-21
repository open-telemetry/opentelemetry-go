// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/gen_ai"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

// OperationNameAttr is an attribute conforming to the gen_ai.operation.name
// semantic conventions. It represents the name of the operation being performed.
type OperationNameAttr string

var (
	// OperationNameChat is the chat completion operation such as [OpenAI Chat API]
	// .
	//
	// [OpenAI Chat API]: https://platform.openai.com/docs/api-reference/chat
	OperationNameChat OperationNameAttr = "chat"
	// OperationNameTextCompletion is the text completions operation such as
	// [OpenAI Completions API (Legacy)].
	//
	// [OpenAI Completions API (Legacy)]: https://platform.openai.com/docs/api-reference/completions
	OperationNameTextCompletion OperationNameAttr = "text_completion"
	// OperationNameEmbeddings is the embeddings operation such as
	// [OpenAI Create embeddings API].
	//
	// [OpenAI Create embeddings API]: https://platform.openai.com/docs/api-reference/embeddings/create
	OperationNameEmbeddings OperationNameAttr = "embeddings"
	// OperationNameCreateAgent is the create GenAI agent.
	OperationNameCreateAgent OperationNameAttr = "create_agent"
	// OperationNameExecuteTool is the execute a tool.
	OperationNameExecuteTool OperationNameAttr = "execute_tool"
)

// SystemAttr is an attribute conforming to the gen_ai.system semantic
// conventions. It represents the Generative AI product as identified by the
// client or server instrumentation.
type SystemAttr string

var (
	// SystemOpenai is the openAI.
	SystemOpenai SystemAttr = "openai"
	// SystemVertexAI is the vertex AI.
	SystemVertexAI SystemAttr = "vertex_ai"
	// SystemGemini is the gemini.
	SystemGemini SystemAttr = "gemini"
	// SystemAnthropic is the anthropic.
	SystemAnthropic SystemAttr = "anthropic"
	// SystemCohere is the cohere.
	SystemCohere SystemAttr = "cohere"
	// SystemAzAIInference is the azure AI Inference.
	SystemAzAIInference SystemAttr = "az.ai.inference"
	// SystemAzAIOpenai is the azure OpenAI.
	SystemAzAIOpenai SystemAttr = "az.ai.openai"
	// SystemIbmWatsonxAI is the IBM Watsonx AI.
	SystemIbmWatsonxAI SystemAttr = "ibm.watsonx.ai"
	// SystemAWSBedrock is the AWS Bedrock.
	SystemAWSBedrock SystemAttr = "aws.bedrock"
	// SystemPerplexity is the perplexity.
	SystemPerplexity SystemAttr = "perplexity"
	// SystemXai is the xAI.
	SystemXai SystemAttr = "xai"
	// SystemDeepseek is the deepSeek.
	SystemDeepseek SystemAttr = "deepseek"
	// SystemGroq is the groq.
	SystemGroq SystemAttr = "groq"
	// SystemMistralAI is the mistral AI.
	SystemMistralAI SystemAttr = "mistral_ai"
)

// TokenTypeAttr is an attribute conforming to the gen_ai.token.type semantic
// conventions. It represents the type of token being counted.
type TokenTypeAttr string

var (
	// TokenTypeInput is the input tokens (prompt, input, etc.).
	TokenTypeInput TokenTypeAttr = "input"
	// TokenTypeOutput is the output tokens (completion, response, etc.).
	TokenTypeOutput TokenTypeAttr = "output"
)

// GenAIClientOperationDuration is an instrument used to record metric values
// conforming to the "gen_ai.client.operation.duration" semantic conventions. It
// represents the genAI operation duration.
type ClientOperationDuration struct {
	inst metric.Float64Histogram
}

// NewClientOperationDuration returns a new ClientOperationDuration instrument.
func NewClientOperationDuration(m metric.Meter) (ClientOperationDuration, error) {
	i, err := m.Float64Histogram(
	    "gen_ai.client.operation.duration",
	    metric.WithDescription("GenAI operation duration"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ClientOperationDuration{}, err
	}
	return ClientOperationDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientOperationDuration) Name() string {
	return "gen_ai.client.operation.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientOperationDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientOperationDuration) Description() string {
	return "GenAI operation duration"
}

// Record records incr to the existing count.
//
// The genAiOperationName is the the name of the operation being performed.
//
// The genAiSystem is the the Generative AI product as identified by the client
// or server instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientOperationDuration) Record(
    ctx context.Context,
    val float64,
	operationName OperationNameAttr,
	system SystemAttr,
	attrs ...ClientOperationDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.system", string(system)),
			)...,
		),
	)
}

func (m ClientOperationDuration) conv(in []ClientOperationDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientOperationDurationAttr()
	}
	return out
}

// ClientOperationDurationAttr is an optional attribute for the
// ClientOperationDuration instrument.
type ClientOperationDurationAttr interface {
    clientOperationDurationAttr() attribute.KeyValue
}

type clientOperationDurationAttr struct {
	kv attribute.KeyValue
}

func (a clientOperationDurationAttr) clientOperationDurationAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientOperationDuration) ErrorType(val ErrorTypeAttr) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("error.type", string(val))}
}

// RequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ClientOperationDuration) RequestModel(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("gen_ai.request.model", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ClientOperationDuration) ServerPort(val int) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.Int("server.port", val)}
}

// ResponseModel returns an optional attribute for the "gen_ai.response.model"
// semantic convention. It represents the name of the model that generated the
// response.
func (ClientOperationDuration) ResponseModel(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("gen_ai.response.model", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the genAI server address.
func (ClientOperationDuration) ServerAddress(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("server.address", val)}
}

// GenAIClientTokenUsage is an instrument used to record metric values conforming
// to the "gen_ai.client.token.usage" semantic conventions. It represents the
// measures number of input and output tokens used.
type ClientTokenUsage struct {
	inst metric.Int64Histogram
}

// NewClientTokenUsage returns a new ClientTokenUsage instrument.
func NewClientTokenUsage(m metric.Meter) (ClientTokenUsage, error) {
	i, err := m.Int64Histogram(
	    "gen_ai.client.token.usage",
	    metric.WithDescription("Measures number of input and output tokens used"),
	    metric.WithUnit("{token}"),
	)
	if err != nil {
	    return ClientTokenUsage{}, err
	}
	return ClientTokenUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientTokenUsage) Name() string {
	return "gen_ai.client.token.usage"
}

// Unit returns the semantic convention unit of the instrument
func (ClientTokenUsage) Unit() string {
	return "{token}"
}

// Description returns the semantic convention description of the instrument
func (ClientTokenUsage) Description() string {
	return "Measures number of input and output tokens used"
}

// Record records incr to the existing count.
//
// The genAiOperationName is the the name of the operation being performed.
//
// The genAiSystem is the the Generative AI product as identified by the client
// or server instrumentation.
//
// The genAiTokenType is the the type of token being counted.
//
// All additional attrs passed are included in the recorded value.
func (m ClientTokenUsage) Record(
    ctx context.Context,
    val int64,
	operationName OperationNameAttr,
	system SystemAttr,
	tokenType TokenTypeAttr,
	attrs ...ClientTokenUsageAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.system", string(system)),
				attribute.String("gen_ai.token.type", string(tokenType)),
			)...,
		),
	)
}

func (m ClientTokenUsage) conv(in []ClientTokenUsageAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientTokenUsageAttr()
	}
	return out
}

// ClientTokenUsageAttr is an optional attribute for the ClientTokenUsage
// instrument.
type ClientTokenUsageAttr interface {
    clientTokenUsageAttr() attribute.KeyValue
}

type clientTokenUsageAttr struct {
	kv attribute.KeyValue
}

func (a clientTokenUsageAttr) clientTokenUsageAttr() attribute.KeyValue {
    return a.kv
}

// RequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ClientTokenUsage) RequestModel(val string) ClientTokenUsageAttr {
	return clientTokenUsageAttr{kv: attribute.String("gen_ai.request.model", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ClientTokenUsage) ServerPort(val int) ClientTokenUsageAttr {
	return clientTokenUsageAttr{kv: attribute.Int("server.port", val)}
}

// ResponseModel returns an optional attribute for the "gen_ai.response.model"
// semantic convention. It represents the name of the model that generated the
// response.
func (ClientTokenUsage) ResponseModel(val string) ClientTokenUsageAttr {
	return clientTokenUsageAttr{kv: attribute.String("gen_ai.response.model", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the genAI server address.
func (ClientTokenUsage) ServerAddress(val string) ClientTokenUsageAttr {
	return clientTokenUsageAttr{kv: attribute.String("server.address", val)}
}

// GenAIServerRequestDuration is an instrument used to record metric values
// conforming to the "gen_ai.server.request.duration" semantic conventions. It
// represents the generative AI server request duration such as time-to-last byte
// or last output token.
type ServerRequestDuration struct {
	inst metric.Float64Histogram
}

// NewServerRequestDuration returns a new ServerRequestDuration instrument.
func NewServerRequestDuration(m metric.Meter) (ServerRequestDuration, error) {
	i, err := m.Float64Histogram(
	    "gen_ai.server.request.duration",
	    metric.WithDescription("Generative AI server request duration such as time-to-last byte or last output token"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ServerRequestDuration{}, err
	}
	return ServerRequestDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerRequestDuration) Name() string {
	return "gen_ai.server.request.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ServerRequestDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ServerRequestDuration) Description() string {
	return "Generative AI server request duration such as time-to-last byte or last output token"
}

// Record records incr to the existing count.
//
// The genAiOperationName is the the name of the operation being performed.
//
// The genAiSystem is the the Generative AI product as identified by the client
// or server instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ServerRequestDuration) Record(
    ctx context.Context,
    val float64,
	operationName OperationNameAttr,
	system SystemAttr,
	attrs ...ServerRequestDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.system", string(system)),
			)...,
		),
	)
}

func (m ServerRequestDuration) conv(in []ServerRequestDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.serverRequestDurationAttr()
	}
	return out
}

// ServerRequestDurationAttr is an optional attribute for the
// ServerRequestDuration instrument.
type ServerRequestDurationAttr interface {
    serverRequestDurationAttr() attribute.KeyValue
}

type serverRequestDurationAttr struct {
	kv attribute.KeyValue
}

func (a serverRequestDurationAttr) serverRequestDurationAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ServerRequestDuration) ErrorType(val ErrorTypeAttr) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("error.type", string(val))}
}

// RequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ServerRequestDuration) RequestModel(val string) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("gen_ai.request.model", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ServerRequestDuration) ServerPort(val int) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.Int("server.port", val)}
}

// ResponseModel returns an optional attribute for the "gen_ai.response.model"
// semantic convention. It represents the name of the model that generated the
// response.
func (ServerRequestDuration) ResponseModel(val string) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("gen_ai.response.model", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the genAI server address.
func (ServerRequestDuration) ServerAddress(val string) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("server.address", val)}
}

// GenAIServerTimePerOutputToken is an instrument used to record metric values
// conforming to the "gen_ai.server.time_per_output_token" semantic conventions.
// It represents the time per output token generated after the first token for
// successful responses.
type ServerTimePerOutputToken struct {
	inst metric.Float64Histogram
}

// NewServerTimePerOutputToken returns a new ServerTimePerOutputToken instrument.
func NewServerTimePerOutputToken(m metric.Meter) (ServerTimePerOutputToken, error) {
	i, err := m.Float64Histogram(
	    "gen_ai.server.time_per_output_token",
	    metric.WithDescription("Time per output token generated after the first token for successful responses"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ServerTimePerOutputToken{}, err
	}
	return ServerTimePerOutputToken{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerTimePerOutputToken) Name() string {
	return "gen_ai.server.time_per_output_token"
}

// Unit returns the semantic convention unit of the instrument
func (ServerTimePerOutputToken) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ServerTimePerOutputToken) Description() string {
	return "Time per output token generated after the first token for successful responses"
}

// Record records incr to the existing count.
//
// The genAiOperationName is the the name of the operation being performed.
//
// The genAiSystem is the the Generative AI product as identified by the client
// or server instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ServerTimePerOutputToken) Record(
    ctx context.Context,
    val float64,
	operationName OperationNameAttr,
	system SystemAttr,
	attrs ...ServerTimePerOutputTokenAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.system", string(system)),
			)...,
		),
	)
}

func (m ServerTimePerOutputToken) conv(in []ServerTimePerOutputTokenAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.serverTimePerOutputTokenAttr()
	}
	return out
}

// ServerTimePerOutputTokenAttr is an optional attribute for the
// ServerTimePerOutputToken instrument.
type ServerTimePerOutputTokenAttr interface {
    serverTimePerOutputTokenAttr() attribute.KeyValue
}

type serverTimePerOutputTokenAttr struct {
	kv attribute.KeyValue
}

func (a serverTimePerOutputTokenAttr) serverTimePerOutputTokenAttr() attribute.KeyValue {
    return a.kv
}

// RequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ServerTimePerOutputToken) RequestModel(val string) ServerTimePerOutputTokenAttr {
	return serverTimePerOutputTokenAttr{kv: attribute.String("gen_ai.request.model", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ServerTimePerOutputToken) ServerPort(val int) ServerTimePerOutputTokenAttr {
	return serverTimePerOutputTokenAttr{kv: attribute.Int("server.port", val)}
}

// ResponseModel returns an optional attribute for the "gen_ai.response.model"
// semantic convention. It represents the name of the model that generated the
// response.
func (ServerTimePerOutputToken) ResponseModel(val string) ServerTimePerOutputTokenAttr {
	return serverTimePerOutputTokenAttr{kv: attribute.String("gen_ai.response.model", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the genAI server address.
func (ServerTimePerOutputToken) ServerAddress(val string) ServerTimePerOutputTokenAttr {
	return serverTimePerOutputTokenAttr{kv: attribute.String("server.address", val)}
}

// GenAIServerTimeToFirstToken is an instrument used to record metric values
// conforming to the "gen_ai.server.time_to_first_token" semantic conventions. It
// represents the time to generate first token for successful responses.
type ServerTimeToFirstToken struct {
	inst metric.Float64Histogram
}

// NewServerTimeToFirstToken returns a new ServerTimeToFirstToken instrument.
func NewServerTimeToFirstToken(m metric.Meter) (ServerTimeToFirstToken, error) {
	i, err := m.Float64Histogram(
	    "gen_ai.server.time_to_first_token",
	    metric.WithDescription("Time to generate first token for successful responses"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ServerTimeToFirstToken{}, err
	}
	return ServerTimeToFirstToken{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerTimeToFirstToken) Name() string {
	return "gen_ai.server.time_to_first_token"
}

// Unit returns the semantic convention unit of the instrument
func (ServerTimeToFirstToken) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ServerTimeToFirstToken) Description() string {
	return "Time to generate first token for successful responses"
}

// Record records incr to the existing count.
//
// The genAiOperationName is the the name of the operation being performed.
//
// The genAiSystem is the the Generative AI product as identified by the client
// or server instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ServerTimeToFirstToken) Record(
    ctx context.Context,
    val float64,
	operationName OperationNameAttr,
	system SystemAttr,
	attrs ...ServerTimeToFirstTokenAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.system", string(system)),
			)...,
		),
	)
}

func (m ServerTimeToFirstToken) conv(in []ServerTimeToFirstTokenAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.serverTimeToFirstTokenAttr()
	}
	return out
}

// ServerTimeToFirstTokenAttr is an optional attribute for the
// ServerTimeToFirstToken instrument.
type ServerTimeToFirstTokenAttr interface {
    serverTimeToFirstTokenAttr() attribute.KeyValue
}

type serverTimeToFirstTokenAttr struct {
	kv attribute.KeyValue
}

func (a serverTimeToFirstTokenAttr) serverTimeToFirstTokenAttr() attribute.KeyValue {
    return a.kv
}

// RequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ServerTimeToFirstToken) RequestModel(val string) ServerTimeToFirstTokenAttr {
	return serverTimeToFirstTokenAttr{kv: attribute.String("gen_ai.request.model", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ServerTimeToFirstToken) ServerPort(val int) ServerTimeToFirstTokenAttr {
	return serverTimeToFirstTokenAttr{kv: attribute.Int("server.port", val)}
}

// ResponseModel returns an optional attribute for the "gen_ai.response.model"
// semantic convention. It represents the name of the model that generated the
// response.
func (ServerTimeToFirstToken) ResponseModel(val string) ServerTimeToFirstTokenAttr {
	return serverTimeToFirstTokenAttr{kv: attribute.String("gen_ai.response.model", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the genAI server address.
func (ServerTimeToFirstToken) ServerAddress(val string) ServerTimeToFirstTokenAttr {
	return serverTimeToFirstTokenAttr{kv: attribute.String("server.address", val)}
}