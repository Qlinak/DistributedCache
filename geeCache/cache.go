package geeCache

import (
	"myMod/geeCache/lru"
	"sync"
)

type cache struct {
	mutexLock sync.Mutex
	lru       *lru.Cache // a pointer of a lru Cache
	maxByte   int64
}

func (myCache *cache) add(key string, value ByteView) {
	myCache.mutexLock.Lock()
	defer myCache.mutexLock.Unlock() // executes right after add returns
	// lazy init - init when in need
	if myCache.lru == nil {
		myCache.lru = lru.New(myCache.maxByte, nil)
	}
	myCache.lru.Add(key, value)
}

func (myCache *cache) get(key string) (value ByteView, ok bool) {
	myCache.mutexLock.Lock()
	defer myCache.mutexLock.Unlock()
	if myCache.lru == nil {
		return
	}

	if value, ok := myCache.lru.Get(key); ok {
		return value.(ByteView), true
	}

	return
}
