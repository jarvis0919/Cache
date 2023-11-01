package lru

import "container/list"

// 缓存的结构体，最大缓存，当前缓存大小，双向链表指针，缓存，回调函数
// 双向链表存储数据并维护的是数据访问量的排序，最近访问的元素被移动到链表头部，
// 触发缓存大小限制时，则移除尾部元素。
// cache则存储每一个元素地址，方便检索数据
type Cache struct {
	maxBytes  int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

//装载缓存的结构体
type entry struct {
	key   string
	value Value
}

//缓存实现Value接口
type Value interface {
	Len() int
}

// 创建缓存的new方法
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//获取指定key的值，从map中查找到元素
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//添加节点
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

//移除最近最少访问的节点
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

//查看添加了多少元素
func (c *Cache) Len() int {
	return c.ll.Len()
}

// //查看缓存
// func (c *Cache) GetNbAndMb() (int64, int64) {
// 	return c.maxBytes, c.nbytes
// }
