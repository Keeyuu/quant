package priorityqueue

import (
	"container/heap"
	"errors"
	"sync"
)

// 不重复task的优先级队列，目前写死了是从大到小排列，pop时取最后一个
type PriorityQueue struct {
	taskItems      priorityTaskHeap
	rejectStrategy RejectedExecutionHandler // 饱和策略
	maxSize        int                      // 队列最大长度，由饱和策略控制, <=0 时不限长度
	lock           sync.RWMutex
}

func (q *PriorityQueue) Init(maxSize int, handler RejectedExecutionHandler) error {
	q.taskItems.items = make([]*priorityTask, 0, 16)
	q.taskItems.itemsMap = make(map[interface{}]*priorityTask, 16)
	if handler != nil {
		if maxSize <= 0 {
			return errors.New("invalid max size")
		}
		q.maxSize = maxSize
		q.rejectStrategy = handler
		handler.init(q)
	}
	return nil
}

func (q *PriorityQueue) Len() int {
	q.lock.RLock()
	size := q.taskItems.Len()
	q.lock.RUnlock()
	return size
}

func (q *PriorityQueue) maxItem() *priorityTask {
	taskLen := q.taskItems.Len()
	if taskLen == 0 {
		return nil
	}
	return q.taskItems.items[0]
}

func (q *PriorityQueue) Push(key, value interface{}, priority int64) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	size := q.taskItems.Len()
	item := q.taskItems.getByKey(key)
	if size > 0 && item != nil {
		q.taskItems.updateItem(item, value, priority)
		return true
	}
	item = &priorityTask{
		Key:      key,
		Value:    value,
		Priority: priority,
		Index:    -1,
	}
	// 超过队列长度，执行饱和策略
	if q.rejectStrategy != nil && size >= q.maxSize {
		if err := q.rejectStrategy.reject(item); err != nil {
			return false
		}
		return true
	}
	heap.Push(&(q.taskItems), item)
	return true
}

func (q *PriorityQueue) Pop() *priorityTask {
	q.lock.Lock()
	defer q.lock.Unlock()
	size := q.taskItems.Len()
	if size > 0 {
		return heap.Pop(&(q.taskItems)).(*priorityTask)
	}
	return nil
}
