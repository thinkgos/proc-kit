package log_flow_collector

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type LogFlowStore[T any] interface {
	Flush(sync bool, data []*T)
}

type LogFlowCollector[T any] struct {
	batchSize     int
	flushInterval time.Duration
	store         LogFlowStore[T]
	entry         chan *T
	closed        atomic.Bool
	wg            sync.WaitGroup
	cancel        context.CancelFunc
	miss          atomic.Int64
}

type Option[T any] func(*LogFlowCollector[T])

func WithBatchSize[T any](size int) Option[T] {
	return func(lc *LogFlowCollector[T]) {
		lc.batchSize = size
	}
}

func WithFlushInterval[T any](interval time.Duration) Option[T] {
	return func(lc *LogFlowCollector[T]) {
		lc.flushInterval = interval
	}
}

// NewLogFlowCollector 创建一个新的 DataflowCollector
// 数据收集采用尽力而为, 可观测.
func NewLogFlowCollector[T any](store LogFlowStore[T], opts ...Option[T]) *LogFlowCollector[T] {
	ctx, cancel := context.WithCancel(context.Background())
	lc := &LogFlowCollector[T]{
		batchSize:     500,
		flushInterval: time.Second * 30,
		store:         store,
		closed:        atomic.Bool{},
		wg:            sync.WaitGroup{},
		cancel:        cancel,
		entry:         make(chan *T),
		miss:          atomic.Int64{},
	}
	for _, f := range opts {
		f(lc)
	}

	lc.entry = make(chan *T, lc.batchSize+lc.batchSize/2)
	lc.wg.Add(1)
	go lc.worker(ctx)
	return lc
}

// Close 在调用Close之前, 要保证没有请求进来调用Send, 以保证数据正常落盘
func (lc *LogFlowCollector[T]) Close() error {
	if lc.closed.CompareAndSwap(false, true) {
		lc.cancel()
		lc.wg.Wait()
	}
	return nil
}

// MissTotal 返回丢弃的条目数
func (lc *LogFlowCollector[T]) MissTotal() int64 {
	return lc.miss.Load()
}

func (lc *LogFlowCollector[T]) Send(entry *T) {
	if lc.closed.Load() {
		return
	}
	select {
	case lc.entry <- entry:
	default:
		// 通道已满，丢弃该条目
		lc.miss.Add(1)
	}
}

func (lc *LogFlowCollector[T]) worker(ctx context.Context) {
	defer lc.wg.Done()

	ticker := time.NewTicker(lc.flushInterval)
	defer ticker.Stop()
	buffer := make([]*T, 0, lc.batchSize)
	for {
		select {
		case entry := <-lc.entry:
			buffer = append(buffer, entry)
			// 达到指定数量直接批量刷库
			if len(buffer) >= lc.batchSize {
				lc.flush(false, buffer)
				ticker.Reset(lc.flushInterval)
				buffer = make([]*T, 0, lc.batchSize)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				lc.flush(false, buffer)
				buffer = make([]*T, 0, lc.batchSize)
			}
		case <-ctx.Done():
			// 清空 chan 所有数据进行刷盘
			for {
				select {
				case entry := <-lc.entry:
					buffer = append(buffer, entry)
					if len(buffer) >= lc.batchSize {
						lc.flush(true, buffer)
						buffer = make([]*T, 0, lc.batchSize)
					}
				default:
					if len(buffer) > 0 {
						lc.flush(true, buffer)
					}
					return
				}
			}
		}
	}
}

func (lc *LogFlowCollector[T]) flush(sync bool, buffer []*T) {
	defer func() {
		if e := recover(); e != nil {
			slog.Error("[log_flow_collector] flush cause panic", "error", e)
		}
	}()
	lc.store.Flush(sync, buffer)
}
