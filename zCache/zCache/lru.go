package zCache

import "container/list"

type Cache struct {
	maxBytes  int64 //允许使用的最大内存
	usedBytes int64 //当前已使用的内存
	ll        *list.List
	cache     map[string]*list.Element //键是字符串，值是双向链表中对应节点的指针。
	//可选，当条目被清除时执行。
	onEvicted func(key string, value Value)
}

// 键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射。
type entry struct {
	key   string
	value Value
}

// Value 为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小。
type Value interface {
	Len() int
}

// New 用于实例化Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// Get 查找功能
// 如果键对应的链表节点存在，则将对应节点移动到队尾，并返回查找到的值。并返回查找到的值
// c.ll.MoveToFront(ele)，即将链表中的节点 ele 移动到队尾（双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾）
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 缓存淘汰，移除最近最少访问的节点（队首）
// c.ll.Back() 取到队首节点，从链表中删除。
// delete(c.cache, kv.key)，从字典中 c.cache 删除该节点的映射关系。
// 更新当前所用的内存 c.usedBytes。
// 如果回调函数 OnEvicted 不为 nil，则调用回调函数。
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key) //从cache中删除key
		c.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增修改
// 如果键存在，则更新对应节点的值，并将该节点移到队尾。
// 不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系。
// 更新 c.usedBytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len()) //如果键存在，更新value大小的差值即可
		kv.value = value

	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}
}

// Len 获取添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
