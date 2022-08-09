// Package lru
// @author    : MuXiang123
// @time      : 2022/7/26 22:33
package lru

import (
	"fmt"
	"reflect"
	"testing"
)

type String string

func (d String) Length() int64 {
	return int64(len(d))
}

func onEvicted(string, Value) {
}

//测试get方法
func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

//当使用内存超过了设定值时，是否会触发“无用”节点的移除
func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	fmt.Println(lru.Get("key3"))
	fmt.Println(lru.Get("k3"))
	fmt.Println(lru.Get("key2"))
	if _, ok := lru.Get("key1"); ok || lru.Length() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

// 测试回调函数是否能够被调用
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 4)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1"}
	//比较非基础类型的两个值
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s +++ %s", expect, keys)
	}
}

func TestAdd(t *testing.T) {
	//callback := func(key string, value Value) {
	//	keys = append(keys, key)
	//}
	lru := New(int64(10), nil)
	lru.Add("key", String("1"))
	lru.Add("key", String("111"))

	if lru.nBytes != int64(len("key")+len("111")) {
		t.Fatal("expected 6 but got", lru.nBytes)
	}
}
