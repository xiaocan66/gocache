package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}
func TestGet(t *testing.T) {
	lru := New(int64(0), nil)

	lru.Add("key1", String("123141"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "123141" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}

}
func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "111", "222", "333"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))

	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	k1, k2, k3, k4 := "k1", "k2", "k3", "k4"
	v1, v2, v3, v4 := String("111"), String("222"), String("1234567890"), String("1234567890")
	keys := make([]string, 0)

	lru := New(int64(24), func(key string, v Value) {
		keys = append(keys, key)
	})
	lru.Add(k1, v1)
	lru.Add(k2, v2)
	lru.Add(k3, v3)
	lru.Add(k4, v4)

	expect := []string{"k1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", keys)
	}
}
