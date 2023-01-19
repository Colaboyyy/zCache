package zCache

import (
	"fmt"
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1 = 1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

// 测试，当使用内存超过了设定值时，是否会触发“无用”节点的移除
func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	fmt.Println("cap:", cap)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	// cap只能刚好容纳k1，k2，v1，v2，k3添加后，key1被淘汰
	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("RemoveOldest ke1 failed!")
	}
}

// 测试回调函数能否被调用
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0) //用于记录被删除的key
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456")) //cap刚好放下，放k2时由于大于cap触发callback，删除key1，将key1添加至keys
	fmt.Println("lru.usedBytes", lru.usedBytes)
	fmt.Println("keys:", keys)
	lru.Add("k2", String("v2"))
	fmt.Println("lru.usedBytes", lru.usedBytes)
	fmt.Println("keys:", keys)
	lru.Add("k3", String("v3"))
	fmt.Println("lru.usedBytes", lru.usedBytes)
	fmt.Println("keys:", keys)
	lru.Add("K4", String("v4")) //放不下了，将k2删除，添加至keys

	expect := []string{"key1", "k2"}
	fmt.Println("lru.usedBytes", lru.usedBytes)
	fmt.Println("keys:", keys)
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
