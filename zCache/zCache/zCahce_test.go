package zCache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed!")
	}
}

// 模拟一个数据库
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGroup_Get(t *testing.T) {
	loadCounts := make(map[string]int, len(db)) // 统计某个键调用回调函数的次数，如果次数大于1，则表示调用了多次回调函数，没有缓存。
	// 使用强制类型转换，将匿名func转换为GetterFunc，函数类型是不能直接调用的，函数实例才可以
	group := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	for k, v := range db {
		if view, err := group.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", k)
		}
		// 从回调函数加载,表示没有缓存
		if _, err := group.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := group.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown should be empty, but got %s", view)
	}

	if view, err := group.Get("Tom"); err != nil {
		t.Fatalf("the value of Tom should be empty, but got %s", view)
	}
}
