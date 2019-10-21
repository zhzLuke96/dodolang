package skiplist

import (
	"math/rand"
	"sync"
	"time"
)

type SkipListNode struct {
	key  string
	data interface{}
	next []*SkipListNode
}

type SkipList struct {
	head   *SkipListNode
	tail   *SkipListNode
	length int
	level  int
	mut    *sync.RWMutex // safe async
	rand   *rand.Rand
}

func (list *SkipList) randomLevel() int {
	level := 1
	for ; level < list.level && list.rand.Uint32()&0x1 == 1; level++ {
	}
	return level
}

func NewSkipList(level int) *SkipList {
	list := &SkipList{}
	if level <= 0 {
		level = 32
	}
	list.level = level
	list.head = &SkipListNode{next: make([]*SkipListNode, level, level)}
	list.tail = &SkipListNode{}
	list.mut = &sync.RWMutex{}
	list.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	for index := range list.head.next {
		list.head.next[index] = list.tail
	}
	return list
}
func (list *SkipList) Set(key string, data interface{}) {
	list.mut.Lock()
	defer list.mut.Unlock()
	// append deepth
	level := list.randomLevel()
	// append idx
	update := make([]*SkipListNode, level, level)
	node := list.head
	for index := level - 1; index >= 0; index-- {
		for {
			node1 := node.next[index]
			if node1 == list.tail || node1.key > key { // insert idx
				update[index] = node
				break
			} else if node1.key == key { // modify idx
				node1.data = data
				return
			} else { // next find
				node = node1
			}
		}
	}
	// insert
	newNode := &SkipListNode{key, data, make([]*SkipListNode, level, level)}
	for index, node := range update {
		node.next[index], newNode.next[index] = newNode, node.next[index]
	}
	list.length++
}

func (list *SkipList) Del(key string) bool {
	list.mut.Lock()
	defer list.mut.Unlock()
	// query
	node := list.head
	remove := make([]*SkipListNode, list.level, list.level)
	var target *SkipListNode
	for index := len(node.next) - 1; index >= 0; index-- {
		for {
			node1 := node.next[index]
			if node1 == list.tail || node1.key > key {
				break
			} else if node1.key == key {
				remove[index] = node // got it
				target = node1
				break
			} else {
				node = node1
			}
		}
	}
	// delete
	if target != nil {
		for index, node1 := range remove {
			if node1 != nil {
				node1.next[index] = target.next[index]
			}
		}
		list.length--
		return true
	}
	return false
}

func (list *SkipList) Get(key string) (interface{}, bool) {
	list.mut.RLock()
	defer list.mut.RUnlock()
	node := list.head
	for index := len(node.next) - 1; index >= 0; index-- {
		for {
			node1 := node.next[index]
			if node1 == list.tail || node1.key > key {
				break
			} else if node1.key == key {
				return node1.data, true
			} else {
				node = node1
			}
		}
	}
	return nil, false
}

func (list *SkipList) Len() int {
	list.mut.RLock()
	defer list.mut.RUnlock()
	return list.length
}
