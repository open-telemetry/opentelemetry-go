// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package global // import "go.opentelemetry.io/otel/internal/global"

/*
This file contains the forwarding implementation of the EntityEmitterProvider used as
the default global instance. Prior to initialization of an SDK, EntityEmitters
returned by the global EntityEmitterProvider will provide no-op functionality. This
means that all Span created prior to initialization are no-op Spans.

Once an SDK has been initialized, all provided no-op EntityEmitters are swapped for
EntityEmitters provided by the SDK defined EntityEmitterProvider. However, any Span started
prior to this initialization does not change its behavior. Meaning, the Span
remains a no-op Span.

The implementation to track and swap EntityEmitters locks all new EntityEmitter creation
until the swap is complete. This assumes that this operation is not
performance-critical. If that assumption is incorrect, be sure to configure an
SDK prior to any EntityEmitter creation.
*/

import (
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/entity"
	"go.opentelemetry.io/otel/entity/embedded"
)

// entityEmitterProvider is a placeholder for a configured SDK EntityEmitterProvider.
//
// All EntityEmitterProvider functionality is forwarded to a delegate once
// configured.
type entityEmitterProvider struct {
	embedded.EntityEmitterProvider

	mtx            sync.Mutex
	entityEmitters map[il]*entityEmitter
	delegate       entity.EntityEmitterProvider
}

// Compile-time guarantee that entityEmitterProvider implements the EntityEmitterProvider
// interface.
var _ entity.EntityEmitterProvider = &entityEmitterProvider{}

// setDelegate configures p to delegate all EntityEmitterProvider functionality to
// provider.
//
// All EntityEmitters provided prior to this function call are switched out to be
// EntityEmitters provided by provider.
//
// It is guaranteed by the caller that this happens only once.
func (p *entityEmitterProvider) setDelegate(provider entity.EntityEmitterProvider) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.delegate = provider

	if len(p.entityEmitters) == 0 {
		return
	}

	for _, t := range p.entityEmitters {
		t.setDelegate(provider)
	}

	p.entityEmitters = nil
}

// EntityEmitter implements EntityEmitterProvider.
func (p *entityEmitterProvider) EntityEmitter(name string, opts ...entity.EntityEmitterOption) entity.EntityEmitter {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.delegate != nil {
		return p.delegate.EntityEmitter(name, opts...)
	}

	// At this moment it is guaranteed that no sdk is installed, save the entityEmitter in the entityEmitters map.

	c := entity.NewEntityEmitterConfig(opts...)
	key := il{
		name:    name,
		version: c.InstrumentationVersion(),
	}

	if p.entityEmitters == nil {
		p.entityEmitters = make(map[il]*entityEmitter)
	}

	if val, ok := p.entityEmitters[key]; ok {
		return val
	}

	t := &entityEmitter{name: name, opts: opts, provider: p}
	p.entityEmitters[key] = t
	return t
}

// entityEmitter is a placeholder for a entity.EntityEmitter.
//
// All EntityEmitter functionality is forwarded to a delegate once configured.
// Otherwise, all functionality is forwarded to a NoopEntityEmitter.
type entityEmitter struct {
	embedded.EntityEmitter

	name     string
	opts     []entity.EntityEmitterOption
	provider *entityEmitterProvider

	delegate atomic.Value
}

// Compile-time guarantee that entityEmitter implements the entity.EntityEmitter interface.
var _ entity.EntityEmitter = &entityEmitter{}

// setDelegate configures t to delegate all EntityEmitter functionality to EntityEmitters
// created by provider.
//
// All subsequent calls to the EntityEmitter methods will be passed to the delegate.
//
// It is guaranteed by the caller that this happens only once.
func (t *entityEmitter) setDelegate(provider entity.EntityEmitterProvider) {
	t.delegate.Store(provider.EntityEmitter(t.name, t.opts...))
}

//// Start implements entity.EntityEmitter by forwarding the call to t.delegate if
//// set, otherwise it forwards the call to a NoopEntityEmitter.
//func (t *entityEmitter) Start(ctx context.Context, name string, opts ...entity.SpanStartOption) (
//	context.Context, entity.Span,
//) {
//	delegate := t.delegate.Load()
//	if delegate != nil {
//		return delegate.(entity.EntityEmitter).Start(ctx, name, opts...)
//	}
//
//	s := nonRecordingSpan{sc: entity.SpanContextFromContext(ctx), entityEmitter: t}
//	ctx = entity.ContextWithSpan(ctx, s)
//	return ctx, s
//}
