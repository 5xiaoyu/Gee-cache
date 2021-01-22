package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes 	int64
	nBytes	 	int64
	ll 			*list.List
	cache 		map[string]*list.Element
	//
	OnEvicted	func(key string, value Value)
}

type entry struct {
	key string
	value Value
}

type  Value interface {
	Len() int
}

// constructor
func New(maxBytes int64, onEvicted func (string ,Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//  Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	// 如果键存在，更新对应的节点值，并将节点移动到最前
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry) // todo 1类型转换吗
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value}) // list 的节点 element指针
		c.cache[key] = ele // map中key对应的value 是 list的element 指针
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// get look ups a key's value
// 如果链表中存在 key对应的节点，将节点移动，返回找到的值

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele,ok := c.cache[key];ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// remove the oldest item
// 缓存淘汰 最近做少访问的 节点
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 回调函数，回调函数不为空，调用
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// cache entries
func (c *Cache) Len() int {
	return  c.ll.Len()
}