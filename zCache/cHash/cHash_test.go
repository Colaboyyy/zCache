package cHash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	}) // 需要知道每一个传入的key的哈希值，这里哈希函数是直接返回key的数字

	hash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	hash.Add("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}

func TestMap_Add(t *testing.T) {
	type fields struct {
		hash     Hash
		replicas int
		keys     []int
		hashMap  map[int]string
	}
	type args struct {
		keys []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "1",
			fields: fields{
				hash: func(key []byte) uint32 {
					i, _ := strconv.Atoi(string(key))
					return uint32(i)
				},
				replicas: 3,
				keys:     nil,
				hashMap:  make(map[int]string),
			},
			args: args{keys: []string{"6", "4", "2"}}},
	}
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				hash:     tt.fields.hash,
				replicas: tt.fields.replicas,
				keys:     tt.fields.keys,
				hashMap:  tt.fields.hashMap,
			}
			m.Add(tt.args.keys...)
			for k, v := range testCases {
				if m.Get(k) != v {
					t.Errorf("Asking for %s, should have yielded %s", k, v)
				}
			}
			m.Add("8")
			testCases["27"] = "8"
			for k, v := range testCases {
				if m.Get(k) != v {
					t.Errorf("Asking for %s, should have yielded %s", k, v)
				}
			}
		})
	}

}
