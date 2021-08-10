package lru

import "fmt"

type LRUCache struct {
	Map     map[interface{}]*DLinkedNode
	Head    *DLinkedNode
	Tail    *DLinkedNode
	MaxSize uint32
	Size    uint32
}
type Key struct {
	DataP *string
}

// DLinkedNode double linked node for lru cache.
type DLinkedNode struct {
	NodeKey   interface{}
	NodeValue *interface{}
	left      *DLinkedNode
	right     *DLinkedNode
}

//NewLRUCache
//maxSize should greater than 1
func NewLRUCache(maxSize uint32) *LRUCache {

	return &LRUCache{
		Map:     make(map[interface{}]*DLinkedNode, 0),
		Size:    0,
		MaxSize: maxSize,
		Tail:    nil,
		Head:    nil,
	}
}

//Set set a kv
//1.map set
// 1.1 map contains ? update and skip step 2 : 1.2
// 1.2 (1) add node to head  (2)set map   (3)size ++ (4)to step 2
//2.check to size and rm node.
// 2.1 if size less or equal: skip
// 2.2 other: remove Tail node and Size --.
func (lru *LRUCache) Set(key, value interface{}) {
	v, contains := lru.Map[key]
	if contains {
		v.NodeValue = &value
	} else { //not contains
		//(1) add node to head
		node := &DLinkedNode{NodeKey: key, NodeValue: &value}
		node.right = lru.Head
		lru.Head = node
		if lru.Size >= 1 {
			lru.Head.left = node
		} else {
			//one element.
			lru.Tail = node
		}
		//(2)set map
		lru.Map[key] = node
		//(3)size ++
		lru.Size++
		//(4)if Size greater that MaxSize, remove Tail node and Size --.
		lru.checkSizeRemoveNode()
	}
}

//Get return *string value if kv exist.
//if contains kv.
//1. let value to be Head
func (lru *LRUCache) Get(key interface{}) (interface{}, error) {
	value, contains := lru.Map[key]
	if !contains {
		return nil, fmt.Errorf("error:%s", "not contains key")
	} else { //contains.
		//nodeLeft maybe is nil.
		nodeLeft := value.left
		// if is Head not node.
		if nodeLeft != nil {
			nodeRight := value.right
			value.right = lru.Head
			nodeLeft.right = nodeRight
			lru.Head = value
		} else {
			//is Head node,init Tail
			lru.Tail = value
		}
	}
	return *value.NodeValue, nil

}

//checkSizeRemoveNode if Size greater that MaxSize,
//rm map kv
//remove Tail node
//Size --.
func (lru *LRUCache) checkSizeRemoveNode() {
	if lru.Size > lru.MaxSize {
		////rm map kv
		delete(lru.Map, lru.Tail.NodeKey)
		//remove Tail node
		tailLeft := lru.Tail.left
		lru.Tail = tailLeft
		//2. Size --
		lru.Size--
	}
}
