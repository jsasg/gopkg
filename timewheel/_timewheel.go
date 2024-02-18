package timewheel

import (
	"container/list"
	"sync"
	"time"
)

// Task 任务
type Task struct {
	expired int64  // 执行时间
	runable func() // 任务回调函数
}

// Bucket 时间轮桶
type Bucket struct {
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
	pos               int64      // 当前时间刻度
	bucketTickMs      int64      // 时间刻度
	bucketSize        int64      // 时间轮中桶的数量
	bucketTotalTickMs int64      // 时间跨度
	buckets           []*Bucket  // 时间轮的桶位
	upperTimeWheel    *TimeWheel // 上一级时间轮引用
	lowerTimeWheel    *TimeWheel // 下一级时间轮引用

	ticker *time.Ticker
}

// New 实例时间轮
func New(tick time.Duration, size int64) *TimeWheel {
	timeWheel := newTimeWheel(tick, size)
	timeWheel.ticker = time.NewTicker(tick)
	return timeWheel
}
func newTimeWheel(tick time.Duration, size int64) *TimeWheel {
	buckets := make([]*Bucket, size)
	for i := range buckets {
		buckets[i] = &Bucket{
			tasks: list.New(),
		}
	}
	tickMs := tick.Milliseconds()
	return &TimeWheel{
		pos:               0,
		bucketTickMs:      tickMs,
		bucketSize:        size,
		bucketTotalTickMs: size * tickMs,
		buckets:           buckets,
	}
}

// Start 开始时间轮
func (tw *TimeWheel) Start() {
	go func() {
		for range tw.ticker.C {
			tw.handle()
		}
	}()
}

// Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	tw.ticker.Stop()
}

// AfterFunc 设置任务
func (tw *TimeWheel) AfterFunc(d time.Duration, fn func()) {
	task := Task{
		expired: d.Milliseconds() + tw.pos*tw.bucketTickMs,
		runable: fn,
	}
	tw.tryInsert(task)
}

// tryInsert 新增任务
func (tw *TimeWheel) tryInsert(task Task) {
	// 执行时间已到，开始执行
	if task.expired <= 0 {
		go task.runable()
	} else {
		if task.expired >= tw.bucketTickMs && task.expired < tw.bucketTotalTickMs {
			// 延迟时间未在本时间轮逸出，加入本时间轮中
			vid := task.expired / tw.bucketTickMs
			slot := vid % tw.bucketSize
			task.expired = task.expired - slot*tw.bucketTickMs
			bucket := tw.buckets[slot]
			bucket.PushBack(task)
		} else {
			// 延迟时间逸出，构建上级时间轮
			if tw.upperTimeWheel == nil {
				tick := time.Duration(tw.bucketTotalTickMs) * time.Millisecond
				tw.upperTimeWheel = newTimeWheel(tick, tw.bucketSize)
				tw.upperTimeWheel.lowerTimeWheel = tw
			}
			tw.upperTimeWheel.tryInsert(task)
		}
	}
}

// tryExecute 执行任务
func (tw *TimeWheel) tryExecute(task Task) {
	// 执行时间已到，开始执行
	if task.expired <= 0 {
		go task.runable()
	} else {
		if task.expired >= tw.bucketTickMs && task.expired < tw.bucketTotalTickMs {
			// 延迟时间未在本时间轮逸出，加入本时间轮中
			vid := task.expired / tw.bucketTickMs
			slot := vid % tw.bucketSize
			task.expired = task.expired - slot*tw.bucketTickMs
			bucket := tw.buckets[slot]
			bucket.PushBack(task)
		} else {
			// 延迟时间逸出，构建上级时间轮
			if tw.lowerTimeWheel != nil {
				tw.lowerTimeWheel.tryExecute(task)
			}
		}
	}
}

// handle 驱动时间轮转动
func (tw *TimeWheel) handle() {
	tw.pos++
	if tw.pos > tw.bucketSize-1 {
		tw.pos = 0
		if tw.upperTimeWheel != nil {
			tw.upperTimeWheel.handle()
		}
	}
	bucket := tw.buckets[tw.pos]
	for item := bucket.tasks.Front(); item != nil; {
		if task, ok := item.Value.(Task); ok {
			tw.tryExecute(task)
			next := item.Next()
			bucket.tasks.Remove(item)
			item = next
		}
	}
}
