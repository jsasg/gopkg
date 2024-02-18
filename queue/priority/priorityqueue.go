package priority

type item struct {
	fn       func(v chan<- interface{})
	priority int64
	index    int
}

func (it *item) Index() int {
	return it.index
}

func NewPriorityQueueItem(priority int64, fn func(v chan<- interface{})) *item {
	return &item{
		fn:       fn,
		priority: priority,
	}
}

type PriorityQueue []*item

func NewPriorityQueue(capacity int64) PriorityQueue {
	return make(PriorityQueue, 0, capacity)
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Peek() (func(v chan<- interface{}), int64) {
	if len(*pq) == 0 {
		return nil, 0
	}
	item := (*pq)[0]
	return item.fn, item.priority
}
