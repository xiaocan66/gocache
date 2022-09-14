package gocache

import (
	"fmt"
	"gocache/singleflight"
	"log"
	"sync"
)

//Getter get接口 当缓存不存在时调用这个接口
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	loader    *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {

	if getter == nil {
		panic("nil getter")

	}
	mu.Lock()

	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    new(singleflight.Group),
	}
	groups[name] = g
	return g
}
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[cache] hit ")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	val, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if val, err := g.getFromPeer(peer, key); err == nil {
					return val, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return val.(ByteView), nil
	}
	return

}

// getLocally 获取源数据
func (g *Group) getLocally(key string) (ByteView, error) {

	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil

}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("Register peerPicker  called more than once")

	}
	g.peers = peers
}

// populateCache 向缓存中添加值
func (g *Group) populateCache(key string, val ByteView) {
	g.mainCache.add(key, val)
}

// getFromPeer 从节点中获取值
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	by, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, nil
	}

	return ByteView{b: by}, nil

}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}
