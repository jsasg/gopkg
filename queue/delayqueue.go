package queue

import (
	"container/heap"
	"sync"
	"time"

	"github.com/jsasg/gopkg/queue/priority"
)

type DelayQueue struct {
	C chan interface{}

	mu sync.Mutex
	pq priority.PriorityQueue

	wakeupC chan struct{}

	exitC chan struct{}
}

// NewDelayQueue 实例化一个延时队列
func NewDelayQueue(size int64) *DelayQueue {
	return &DelayQueue{
		C:       make(chan interface{}),
		pq:      priority.NewPriorityQueue(size),
		wakeupC: make(chan struct{}),
		exitC:   make(chan struct{}),
	}
}

// AfterFunc 给延时队列设置任务
func (dq *DelayQueue) AfterFunc(d time.Duration, fn func(v chan<- interface{})) {
	item := priority.NewPriorityQueueItem(time.Now().Add(d).UnixMilli(), fn)
	dq.mu.Lock()
	heap.Push(&dq.pq, item)
	index := item.Index()
	dq.mu.Unlock()
	if index == 0 {
		dq.wakeupC <- struct{}{}
	}
}

// Start 启动延时队列
func (dq *DelayQueue) Start() {
	go func() {
		for {
			dq.mu.Lock()
			callbackFunc, delay := dq.pq.Peek()
			delta := delay - time.Now().UnixMilli()
			dq.mu.Unlock()

			if callbackFunc == nil {
				select {
				case <-dq.wakeupC:
					continue
				case <-dq.exitC:
					return
				}
			}

			select {
			case <-time.After(time.Duration(delta) * time.Millisecond):
				// 执行
				callbackFunc(dq.C)
				heap.Remove(&dq.pq, 0)
			case <-dq.wakeupC:
				continue
			case <-dq.exitC:
				return
			}
		}
	}()
}

// Stop 停止延时队列
func (dq *DelayQueue) Stop() {
	dq.exitC <- struct{}{}
	close(dq.exitC)
	close(dq.wakeupC)
	close(dq.C)
}
