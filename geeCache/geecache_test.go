package geeCache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetterFunc_Get(t *testing.T) {
	f := GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expected := []byte("key")
	if value, _ := f.Get("key"); !reflect.DeepEqual(value, expected) {
		t.Errorf("callback fails")
	}
}

var scoreDb = map[string]string{
	"Abe":   "100",
	"Bob":   "50",
	"Cathy": "90",
}

// TestGet - test callback if cache is empty, after callback test cache hit
func TestGet(t *testing.T) {
	callBackCounts := make(map[string]int, len(scoreDb))
	myGroup := NewGroup("score", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Printf("slow - Getting from local db - key: %s", key)
		if v, ok := scoreDb[key]; ok {
			if _, ok := callBackCounts[key]; !ok {
				// lazy init
				callBackCounts[key] = 0
			}
			callBackCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k, v := range scoreDb {
		// test callback
		if view, err := myGroup.Get(k); err != nil || view.AsString() != v {
			t.Fatalf("failed to get value from %s", k)
		}
		// test callback cannot be called more than once
		if _, err := myGroup.Get(k); err != nil || callBackCounts[k] > 1 {
			t.Fatalf("cache hit failed from %s", k)
		}
	}

	// test unknown key input
	if view, err := myGroup.Get("unknown"); err == nil {
		t.Fatalf("get trash from key: unknown: %s", view)
	}
}
