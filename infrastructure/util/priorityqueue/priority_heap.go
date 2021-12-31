package priorityqueue

import "container/heap"

// 最基本的task结构
type priorityTask struct {
	Key   interface{}
	Value interface{}

	Priority int64
	Index    int
}

// 实现了heap的task的优先级堆
type priorityTaskHeap struct {
	items    []*priorityTask
	itemsMap map[interface{}]*priorityTask
}

func (tasks priorityTaskHeap) Len() int {
	return len(tasks.items)
}

func (tasks priorityTaskHeap) Less(i, j int) bool {
	return tasks.items[i].Priority > tasks.items[j].Priority
}

func (tasks priorityTaskHeap) Swap(i, j int) {
	tasks.items[i], tasks.items[j] = tasks.items[j], tasks.items[i]
	tasks.items[i].Index = j
	tasks.items[j].Index = i
	if tasks.itemsMap != nil {
		tasks.itemsMap[tasks.items[i].Key] = tasks.items[i]
		tasks.itemsMap[tasks.items[j].Key] = tasks.items[j]
	}
}

func (tasks *priorityTaskHeap) Push(x interface{}) {
	n := len(tasks.items)
	item := x.(*priorityTask)
	item.Index = n
	tasks.items = append(tasks.items, item)
	tasks.itemsMap[item.Key] = item
}

func (tasks *priorityTaskHeap) Pop() interface{} {
	oldTasks := tasks.items
	n := len(oldTasks)
	item := oldTasks[n-1]
	item.Index = -1
	delete(tasks.itemsMap, item.Key)
	tasks.items = oldTasks[0 : n-1]
	return item
}

func (tasks *priorityTaskHeap) Update(key interface{}, value interface{}, priority int64) {
	item := tasks.getByKey(key)
	if item != nil {
		tasks.updateItem(item, value, priority)
	}
}

func (tasks *priorityTaskHeap) getByKey(key interface{}) *priorityTask {
	if item, ok := tasks.itemsMap[key]; ok {
		return item
	}
	return nil
}

func (tasks *priorityTaskHeap) updateItem(item *priorityTask, value interface{}, priority int64) {
	item.Value = value
	item.Priority = priority
	heap.Fix(tasks, item.Index)
}
