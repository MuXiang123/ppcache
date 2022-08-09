// Package ppcache
// @author    : MuXiang123
// @time      : 2022/7/27 22:45
package ppcache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

//模拟耗时的数据库
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetGroup(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	pp := NewGroup("scores", 2<<10, GetterFunc(
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
		if view, err := pp.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		} // load from callback function
		if _, err := pp.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := pp.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

func TestGetterFunc_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		f       GetterFunc
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_Get(t *testing.T) {
	type fields struct {
		name      string
		getter    Getter
		mainCache cache
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ByteView
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				name:      tt.fields.name,
				getter:    tt.fields.getter,
				mainCache: tt.fields.mainCache,
			}
			got, err := g.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_getLocally(t *testing.T) {
	type fields struct {
		name      string
		getter    Getter
		mainCache cache
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ByteView
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				name:      tt.fields.name,
				getter:    tt.fields.getter,
				mainCache: tt.fields.mainCache,
			}
			got, err := g.getLocally(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLocally() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLocally() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_load(t *testing.T) {
	type fields struct {
		name      string
		getter    Getter
		mainCache cache
	}
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue ByteView
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				name:      tt.fields.name,
				getter:    tt.fields.getter,
				mainCache: tt.fields.mainCache,
			}
			gotValue, err := g.load(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestGroup_populateCache(t *testing.T) {
	type fields struct {
		name      string
		getter    Getter
		mainCache cache
	}
	type args struct {
		key   string
		value ByteView
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				name:      tt.fields.name,
				getter:    tt.fields.getter,
				mainCache: tt.fields.mainCache,
			}
			g.populateCache(tt.args.key, tt.args.value)
		})
	}
}

func TestNewGroup(t *testing.T) {
	type args struct {
		name       string
		cacheBytes int64
		getter     Getter
	}
	tests := []struct {
		name string
		args args
		want *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGroup(tt.args.name, tt.args.cacheBytes, tt.args.getter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
