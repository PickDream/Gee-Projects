package lru

import "container/list"

type Cache struct {
	maxBytes     int64
	currentBytes int64
	list         *list.List
	cacheMap     map[string]*list.Element
	OnEvicted    func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

//抽象出可计算的接口
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		list:      list.New(),
		cacheMap:  make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	//如果已经存在
	if element, ok := c.cacheMap[key]; ok {
		c.list.MoveToFront(element)
		sourceValue := element.Value.(*entry)
		c.currentBytes += int64(value.Len()) - int64(sourceValue.value.Len())
		sourceValue.value = value
	} else { //如果之前没有存在
		element = c.list.PushFront(&entry{key, value})
		c.cacheMap[key] = element
		c.maxBytes += int64(len(key)) + int64(value.Len())
	}
	//检查容量，维护LRU缓存
	for c.maxBytes != 0 && c.maxBytes < c.currentBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	element := c.list.Back()
	if element != nil {
		c.list.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cacheMap, kv.key)
		c.currentBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.list.Len()
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cacheMap[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}
