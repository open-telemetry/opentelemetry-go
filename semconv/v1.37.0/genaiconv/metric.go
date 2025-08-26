// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "gen_ai" namespace.
package genaiconv

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

// OperationNameAttr is an attribute conforming to the gen_ai.operation.name
// semantic conventions. It represents the name of the operation being performed.
type OperationNameAttr string

var (
	// OperationNameChat is the chat completion operation such as [OpenAI Chat API]
	// .
	//
	// [OpenAI Chat API]: https://platform.openai.com/docs/api-reference/chat
	OperationNameChat OperationNameAttr = "chat"
	// OperationNameGenerateContent is the multimodal content generation operation
	// such as [Gemini Generate Content].
	//
	// [Gemini Generate Content]: https://ai.google.dev/api/generate-content
	OperationNameGenerateContent OperationNameAttr = "generate_content"
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
	// OperationNameInvokeAgent is the invoke GenAI agent.
	OperationNameInvokeAgent OperationNameAttr = "invoke_agent"
	// OperationNameExecuteTool is the execute a tool.
	OperationNameExecuteTool OperationNameAttr = "execute_tool"
)

// ProviderNameAttr is an attribute conforming to the gen_ai.provider.name
// semantic conventions. It represents the Generative AI provider as identified
// by the client or server instrumentation.
type ProviderNameAttr string

var (
	// ProviderNameOpenAI is the [OpenAI].
	//
	// [OpenAI]: https://openai.com/
	ProviderNameOpenAI ProviderNameAttr = "openai"
	// ProviderNameGCPGenAI is the any Google generative AI endpoint.
	ProviderNameGCPGenAI ProviderNameAttr = "gcp.gen_ai"
	// ProviderNameGCPVertexAI is the [Vertex AI].
	//
	// [Vertex AI]: https://cloud.google.com/vertex-ai
	ProviderNameGCPVertexAI ProviderNameAttr = "gcp.vertex_ai"
	// ProviderNameGCPGemini is the [Gemini].
	//
	// [Gemini]: https://cloud.google.com/products/gemini
	ProviderNameGCPGemini ProviderNameAttr = "gcp.gemini"
	// ProviderNameAnthropic is the [Anthropic].
	//
	// [Anthropic]: https://www.anthropic.com/
	ProviderNameAnthropic ProviderNameAttr = "anthropic"
	// ProviderNameCohere is the [Cohere].
	//
	// [Cohere]: https://cohere.com/
	ProviderNameCohere ProviderNameAttr = "cohere"
	// ProviderNameAzureAIInference is the azure AI Inference.
	ProviderNameAzureAIInference ProviderNameAttr = "azure.ai.inference"
	// ProviderNameAzureAIOpenAI is the [Azure OpenAI].
	//
	// [Azure OpenAI]: https://azure.microsoft.com/products/ai-services/openai-service/
	ProviderNameAzureAIOpenAI ProviderNameAttr = "azure.ai.openai"
	// ProviderNameIBMWatsonxAI is the [IBM Watsonx AI].
	//
	// [IBM Watsonx AI]: https://www.ibm.com/products/watsonx-ai
	ProviderNameIBMWatsonxAI ProviderNameAttr = "ibm.watsonx.ai"
	// ProviderNameAWSBedrock is the [AWS Bedrock].
	//
	// [AWS Bedrock]: https://aws.amazon.com/bedrock
	ProviderNameAWSBedrock ProviderNameAttr = "aws.bedrock"
	// ProviderNamePerplexity is the [Perplexity].
	//
	// [Perplexity]: https://www.perplexity.ai/
	ProviderNamePerplexity ProviderNameAttr = "perplexity"
	// ProviderNameXAI is the [xAI].
	//
	// [xAI]: https://x.ai/
	ProviderNameXAI ProviderNameAttr = "x_ai"
	// ProviderNameDeepseek is the [DeepSeek].
	//
	// [DeepSeek]: https://www.deepseek.com/
	ProviderNameDeepseek ProviderNameAttr = "deepseek"
	// ProviderNameGroq is the [Groq].
	//
	// [Groq]: https://groq.com/
	ProviderNameGroq ProviderNameAttr = "groq"
	// ProviderNameMistralAI is the [Mistral AI].
	//
	// [Mistral AI]: https://mistral.ai/
	ProviderNameMistralAI ProviderNameAttr = "mistral_ai"
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

// ClientOperationDuration is an instrument used to record metric values
// conforming to the "gen_ai.client.operation.duration" semantic conventions. It
// represents the genAI operation duration.
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
		"gen_ai.client.operation.duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("GenAI operation duration."),
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
	return "gen_ai.client.operation.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientOperationDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientOperationDuration) Description() string {
	return "GenAI operation duration."
}

// Record records val to the current distribution for attrs.
//
// The operationName is the the name of the operation being performed.
//
// The providerName is the the Generative AI provider as identified by the client
// or server instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientOperationDuration) Record(
	ctx context.Context,
	val float64,
	operationName OperationNameAttr,
	providerName ProviderNameAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

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
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.provider.name", string(providerName)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ClientOperationDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientOperationDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrRequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ClientOperationDuration) AttrRequestModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.request.model", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ClientOperationDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrResponseModel returns an optional attribute for the
// "gen_ai.response.model" semantic convention. It represents the name of the
// model that generated the response.
func (ClientOperationDuration) AttrResponseModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.response.model", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the genAI server address.
func (ClientOperationDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientTokenUsage is an instrument used to record metric values conforming to
// the "gen_ai.client.token.usage" semantic conventions. It represents the number
// of input and output tokens used.
type ClientTokenUsage struct {
	metric.Int64Histogram
}

// NewClientTokenUsage returns a new ClientTokenUsage instrument.
func NewClientTokenUsage(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (ClientTokenUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientTokenUsage{noop.Int64Histogram{}}, nil
	}

	i, err := m.Int64Histogram(
		"gen_ai.client.token.usage",
		append([]metric.Int64HistogramOption{
			metric.WithDescription("Number of input and output tokens used."),
			metric.WithUnit("{token}"),
		}, opt...)...,
	)
	if err != nil {
	    return ClientTokenUsage{noop.Int64Histogram{}}, err
	}
	return ClientTokenUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientTokenUsage) Inst() metric.Int64Histogram {
	return m.Int64Histogram
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
	return "Number of input and output tokens used."
}

// Record records val to the current distribution for attrs.
//
// The operationName is the the name of the operation being performed.
//
// The providerName is the the Generative AI provider as identified by the client
// or server instrumentation.
//
// The tokenType is the the type of token being counted.
//
// All additional attrs passed are included in the recorded value.
func (m ClientTokenUsage) Record(
	ctx context.Context,
	val int64,
	operationName OperationNameAttr,
	providerName ProviderNameAttr,
	tokenType TokenTypeAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Histogram.Record(ctx, val)
		return
	}

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
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.provider.name", string(providerName)),
				attribute.String("gen_ai.token.type", string(tokenType)),
			)...,
		),
	)

	m.Int64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ClientTokenUsage) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Histogram.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Histogram.Record(ctx, val, *o...)
}

// AttrRequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ClientTokenUsage) AttrRequestModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.request.model", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ClientTokenUsage) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrResponseModel returns an optional attribute for the
// "gen_ai.response.model" semantic convention. It represents the name of the
// model that generated the response.
func (ClientTokenUsage) AttrResponseModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.response.model", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the genAI server address.
func (ClientTokenUsage) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ServerRequestDuration is an instrument used to record metric values conforming
// to the "gen_ai.server.request.duration" semantic conventions. It represents
// the generative AI server request duration such as time-to-last byte or last
// output token.
type ServerRequestDuration struct {
	metric.Float64Histogram
}

// NewServerRequestDuration returns a new ServerRequestDuration instrument.
func NewServerRequestDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ServerRequestDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerRequestDuration{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"gen_ai.server.request.duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Generative AI server request duration such as time-to-last byte or last output token."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return ServerRequestDuration{noop.Float64Histogram{}}, err
	}
	return ServerRequestDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerRequestDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
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
	return "Generative AI server request duration such as time-to-last byte or last output token."
}

// Record records val to the current distribution for attrs.
//
// The operationName is the the name of the operation being performed.
//
// The providerName is the the Generative AI provider as identified by the client
// or server instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ServerRequestDuration) Record(
	ctx context.Context,
	val float64,
	operationName OperationNameAttr,
	providerName ProviderNameAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

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
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.provider.name", string(providerName)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ServerRequestDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ServerRequestDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrRequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ServerRequestDuration) AttrRequestModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.request.model", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ServerRequestDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrResponseModel returns an optional attribute for the
// "gen_ai.response.model" semantic convention. It represents the name of the
// model that generated the response.
func (ServerRequestDuration) AttrResponseModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.response.model", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the genAI server address.
func (ServerRequestDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ServerTimePerOutputToken is an instrument used to record metric values
// conforming to the "gen_ai.server.time_per_output_token" semantic conventions.
// It represents the time per output token generated after the first token for
// successful responses.
type ServerTimePerOutputToken struct {
	metric.Float64Histogram
}

// NewServerTimePerOutputToken returns a new ServerTimePerOutputToken instrument.
func NewServerTimePerOutputToken(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ServerTimePerOutputToken, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerTimePerOutputToken{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"gen_ai.server.time_per_output_token",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Time per output token generated after the first token for successful responses."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return ServerTimePerOutputToken{noop.Float64Histogram{}}, err
	}
	return ServerTimePerOutputToken{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerTimePerOutputToken) Inst() metric.Float64Histogram {
	return m.Float64Histogram
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
	return "Time per output token generated after the first token for successful responses."
}

// Record records val to the current distribution for attrs.
//
// The operationName is the the name of the operation being performed.
//
// The providerName is the the Generative AI provider as identified by the client
// or server instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ServerTimePerOutputToken) Record(
	ctx context.Context,
	val float64,
	operationName OperationNameAttr,
	providerName ProviderNameAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

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
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.provider.name", string(providerName)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ServerTimePerOutputToken) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrRequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ServerTimePerOutputToken) AttrRequestModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.request.model", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ServerTimePerOutputToken) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrResponseModel returns an optional attribute for the
// "gen_ai.response.model" semantic convention. It represents the name of the
// model that generated the response.
func (ServerTimePerOutputToken) AttrResponseModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.response.model", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the genAI server address.
func (ServerTimePerOutputToken) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ServerTimeToFirstToken is an instrument used to record metric values
// conforming to the "gen_ai.server.time_to_first_token" semantic conventions. It
// represents the time to generate first token for successful responses.
type ServerTimeToFirstToken struct {
	metric.Float64Histogram
}

// NewServerTimeToFirstToken returns a new ServerTimeToFirstToken instrument.
func NewServerTimeToFirstToken(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ServerTimeToFirstToken, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerTimeToFirstToken{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"gen_ai.server.time_to_first_token",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Time to generate first token for successful responses."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return ServerTimeToFirstToken{noop.Float64Histogram{}}, err
	}
	return ServerTimeToFirstToken{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerTimeToFirstToken) Inst() metric.Float64Histogram {
	return m.Float64Histogram
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
	return "Time to generate first token for successful responses."
}

// Record records val to the current distribution for attrs.
//
// The operationName is the the name of the operation being performed.
//
// The providerName is the the Generative AI provider as identified by the client
// or server instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ServerTimeToFirstToken) Record(
	ctx context.Context,
	val float64,
	operationName OperationNameAttr,
	providerName ProviderNameAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

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
				attribute.String("gen_ai.operation.name", string(operationName)),
				attribute.String("gen_ai.provider.name", string(providerName)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ServerTimeToFirstToken) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrRequestModel returns an optional attribute for the "gen_ai.request.model"
// semantic convention. It represents the name of the GenAI model a request is
// being made to.
func (ServerTimeToFirstToken) AttrRequestModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.request.model", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the genAI server port.
func (ServerTimeToFirstToken) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrResponseModel returns an optional attribute for the
// "gen_ai.response.model" semantic convention. It represents the name of the
// model that generated the response.
func (ServerTimeToFirstToken) AttrResponseModel(val string) attribute.KeyValue {
	return attribute.String("gen_ai.response.model", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the genAI server address.
func (ServerTimeToFirstToken) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}