package model_lru

import (
	"reflect"
	"testing"

	"cache/internal/model/lru"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	//设置最大缓存为0
	lru := lru.New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value lru.Value) {
		keys = append(keys, key)

	}
	lru := lru.New(int64(10), callback)
	lru.Add("key1", String("123456"))

	lru.Add("k2", String("k2"))

	lru.Add("k3", String("k3"))

	lru.Add("k4", String("k4"))
	expect := []string{"key1", "k2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}

}
