package timewheel

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jsasg/gopkg/queue"
)

// Task 任务
type Task struct {
	expired int64  // 执行时间
	runable func() // 任务回调函数
}

// Bucket 时间轮桶
type Bucket struct {
	expired int64

	tasks *list.List
	mu    sync.Mutex
}

// PushBack 往时间轮桶中加入任务
func (b *Bucket) PushBack(t Task) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tasks.PushBack(t)
}

type TimeWheel struct {
	currentTime       int64
	bucketTickMs      int64      // 时间刻度
	bucketSize        int64      // 时间轮中桶的数量
	bucketTotalTickMs int64      // 时间跨度
	buckets           []*Bucket  // 时间轮的桶位
	upperTimeWheel    *TimeWheel // 上一级时间轮引用

	queue *queue.DelayQueue
}

// New 实例时间轮
func New(tick time.Duration, size int64) *TimeWheel {
	startMs := time.Now().UnixMilli()
	return newTimeWheel(tick, size, startMs, queue.NewDelayQueue(size))
}

func newTimeWheel(tick time.Duration, size int64, startMs int64, queue *queue.DelayQueue) *TimeWheel {
	buckets := make([]*Bucket, size)
	for i := range buckets {
		buckets[i] = &Bucket{
			tasks: list.New(),
		}
	}
	tickMs := tick.Milliseconds()
	currentTime := startMs - startMs%tickMs
	return &TimeWheel{
		currentTime:       currentTime,
		bucketTickMs:      tickMs,
		bucketSize:        size,
		bucketTotalTickMs: size * tickMs,
		buckets:           buckets,

		queue: queue,
	}
}

// Start 开始时间轮
func (tw *TimeWheel) Start() {
	tw.queue.Start()
	go func() {
		for item := range tw.queue.C {
			if b, ok := item.(*Bucket); ok {
				tw.clock(b.expired)

				for t := b.tasks.Front(); t != nil; {
					if task, ok := t.Value.(Task); ok {
						tw.tryExecute(task)
						n := t.Next()
						b.tasks.Remove(t)
						t = n
					}
				}
			}
		}
	}()
}

// Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	tw.queue.Stop()
}

// AfterFunc 设置任务
func (tw *TimeWheel) AfterFunc(d time.Duration, fn func()) {
	task := Task{
		expired: time.Now().Add(d).UnixMilli(),
		runable: fn,
	}
	tw.tryExecute(task)
}

// tryExecute 新增任务，未执行时间任务则加入时间轮对应槽位中等待执行
func (tw *TimeWheel) tryExecute(task Task) {
	// 执行时间已到，开始执行
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if task.expired < currentTime+tw.bucketTickMs {
		go task.runable()
	} else {
		if task.expired < currentTime+tw.bucketTotalTickMs {
			// 延迟时间未在本时间轮逸出，加入本时间轮中
			vid := task.expired / tw.bucketTickMs
			slot := vid % tw.bucketSize
			bucket := tw.buckets[slot]
			bucket.PushBack(task)

			// 延时队列驱动任务执行
			bucket.expired = vid * tw.bucketTickMs
			delta := time.Until(time.UnixMilli(bucket.expired))
			tw.queue.AfterFunc(delta, func(v chan<- interface{}) { v <- bucket })
		} else {
			// 延迟时间逸出，构建上级时间轮
			if tw.upperTimeWheel == nil {
				tick := time.Duration(tw.bucketTotalTickMs) * time.Millisecond
				tw.upperTimeWheel = newTimeWheel(tick, tw.bucketSize, tw.currentTime, tw.queue)
			}
			tw.upperTimeWheel.tryExecute(task)
		}
	}
}

// clock 驱动时间变化
func (tw *TimeWheel) clock(expired int64) {
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if expired >= currentTime+tw.bucketTickMs {
		currentTime = expired - expired%tw.bucketTickMs
		atomic.StoreInt64(&tw.currentTime, currentTime)

		if tw.upperTimeWheel != nil {
			tw.upperTimeWheel.clock(currentTime)
		}
	}
}
