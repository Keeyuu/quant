package priorityqueue

import (
	"container/heap"
	"errors"
)

// 队列饱和策略
type RejectedExecutionHandler interface {
	init(q *PriorityQueue)
	reject(item *priorityTask) error
}

type ReplaceRejectedStrategy struct {
	queue *PriorityQueue
}

func (handler *ReplaceRejectedStrategy) init(q *PriorityQueue) {
	handler.queue = q
}

// 替换最大的那个
func (handler *ReplaceRejectedStrategy) reject(item *priorityTask) error {
	max := handler.queue.maxItem()
	if max.Priority <= item.Priority {
		return errors.New("the new task is the max one, can not push")
	}
	heap.Pop(&(handler.queue.taskItems))
	heap.Push(&(handler.queue.taskItems), item)
	return nil
}
