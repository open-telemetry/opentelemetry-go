// Copyright 2019, OpenTelemetry Authors
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

package trace

import (
	"context"
	"sync"
	"time"
)

const (
	defaultMaxQueueSize = 2048
	defaultScheduledDelayMillis = time.Duration(5000 * time.Millisecond)
	defaultMaxExportBatchSize = 512
)
type BatchSpanProcessorOption struct {
	// MaxQueueSize is the maximum queue size. After the size is reached spans are dropped.
	// The default value is 2048.
	MaxQueueSize int

	// ScheduledDelayMillis is the delay interval in milliseconds between two consecutive exports.
	// The default value is 5000.
	ScheduledDelayMillis time.Duration

	// MaxExportBatchSize IS the maximum batch size of every export.
	// The default value is 512.
	MaxExportBatchSize int
}

// SimpleSpanProcessor implements SpanProcessor interfaces. It is used by
// exporters to receive SpanData synchronously when span is finished.
type BatchSpanProcessor struct {
	exporter BatchExporter
	o        BatchSpanProcessorOption

	qLock    sync.RWMutex
	queue    []*SpanData
	dropped  uint32

	stopCh             chan struct{}
	stopOnce           sync.Once

}

var _ SpanProcessor = (*BatchSpanProcessor)(nil)

// NewBatchSpanProcessor creates a new instance of BatchSpanProcessor
// for a given exporter.
func NewBatchSpanProcessor(exporter BatchExporter, opts BatchSpanProcessorOption) *BatchSpanProcessor {
	bsp := &BatchSpanProcessor{
		exporter: exporter,
		o: opts,
	}

	if bsp.o.ScheduledDelayMillis.Nanoseconds() == 0 {
		bsp.o.ScheduledDelayMillis = defaultScheduledDelayMillis
	}
	if bsp.o.MaxQueueSize <= 0 {
		bsp.o.MaxQueueSize = defaultMaxQueueSize
	}
	if bsp.o.MaxExportBatchSize <= 0 {
		bsp.o.MaxExportBatchSize = defaultMaxExportBatchSize
	}
	bsp.queue =    make([]*SpanData, 0, bsp.o.MaxQueueSize)

	//Start timer to export metrics
	ticker := time.NewTicker(bsp.o.ScheduledDelayMillis)
	go func(ctx context.Context) {
		defer ticker.Stop()
		for {
			select {
			case <-bsp.stopCh:
				return
			case <-ticker.C:
				sds := bsp.copyQueue()
				for i := 0; i < len(sds) ; i += bsp.o.MaxExportBatchSize {
					endIdx := i + bsp.o.MaxExportBatchSize
					if endIdx > len(sds) {
						endIdx = len(sds)
					}
					bsp.exporter.ExportSpans(sds[i:endIdx])
				}
			}
		}
	}(context.Background())

	return bsp
}

// OnStart method does nothing.
func (bsp *BatchSpanProcessor) OnStart(sd *SpanData) {
}

// OnEnd method exports SpanData using associated exporter.
func (bsp *BatchSpanProcessor) OnEnd(sd *SpanData) {
	bsp.qLock.Lock()
	defer bsp.qLock.Unlock()

	if len(bsp.queue) < bsp.o.MaxQueueSize {
		bsp.queue = append(bsp.queue, sd)
	} else {
		bsp.dropped++
	}
}

// Shutdown method does nothing. There is no data to cleanup.
func (bsp *BatchSpanProcessor) Shutdown() {
}

func (bsp *BatchSpanProcessor) copyQueue() []*SpanData {
	bsp.qLock.Lock()
	defer bsp.qLock.Unlock()

	sdToExport := bsp.queue[0:]
	bsp.queue = make([]*SpanData, 0, bsp.o.MaxQueueSize)
	return sdToExport
}