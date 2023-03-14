package geeCache

import (
	"fmt"
	"log"
	"myMod/geeCache/singleflight"
	"sync"
)

// Group struct definition
type Group struct {
	name      string
	mainCache cache
	getter    Getter
	peers     PeerPicker
	// make sure each key is only fetched one
	loader *singleflight.Group
}

// RegisterPeers registers a PeerPicker for choosing remote peer
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// Getter - interface
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc - interface function declaration
// this func is used to pass in as param, not only func can be passed in
// struct that implement the interface can also be passed in
type GetterFunc func(key string) ([]byte, error)

// Get - implementation of the interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// var declaration, lock and groups
var (
	mutexLock sync.RWMutex
	groups    = make(map[string]*Group)
)

// NewGroup - Group Ctor
func NewGroup(name string, maxByte int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{maxByte: maxByte},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mutexLock.RLock() // multiple reader can acquire lock but no writer can acquire lock
	g := groups[name]
	mutexLock.RUnlock()
	return g
}

// Get value from a key
// 1. try from mainCache - if exist return the value
// 2. cache not exist - get from local (e.g. other peers)
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("empty key")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Printf("Cache hit - key: %s", key)
		return v, nil
	}
	log.Printf("[%s] Cache miss - key: %s", g.name, key)
	// not in cache, get from remote peers or local
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	res, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				} else {
					log.Println("[GeeCache] Failed to get from peer", err)
				}
			}
		}

		return g.getLocally(key)
	})

	if err == nil {
		return res.(ByteView), nil
	}
	return ByteView{}, err
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// use g's getter to get locally
	myByte, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	myView := ByteView{
		byteArr: make([]byte, len(myByte)),
	}
	copy(myView.byteArr, myByte)
	g.populateCache(key, myView)
	return myView, nil
}

// populateCache - push the missed cache into the mainCache
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{byteArr: bytes}, nil
}
